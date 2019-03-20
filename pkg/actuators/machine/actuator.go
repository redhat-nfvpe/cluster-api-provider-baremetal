/*
Copyright 2018 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package machine

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	server "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/server"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	apierrors "github.com/openshift/cluster-api/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	providerconfigv1 "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"

	goipmi "github.com/vmware/goipmi"
)

var MachineActuator *Actuator

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	client             client.Client
	kubeClient         kubernetes.Interface
	eventRecorder      record.EventRecorder
	codec              codec
	baremetalAPIServer *server.APIServer
}

type codec interface {
	DecodeFromProviderSpec(machinev1.ProviderSpec, runtime.Object) error
	DecodeProviderStatus(*runtime.RawExtension, runtime.Object) error
	EncodeProviderStatus(runtime.Object) (*runtime.RawExtension, error)
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	Client              client.Client
	Codec               codec
	KubeClient          kubernetes.Interface
	EventRecorder       record.EventRecorder
	ServerListenAddress string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {

	actuator := &Actuator{
		client:        params.Client,
		codec:         params.Codec,
		kubeClient:    params.KubeClient,
		eventRecorder: params.EventRecorder,
	}

	var err error

	bms := server.BaremetalServer{}
	handler := server.NewServerAPIHandler(bms)

	actuator.baremetalAPIServer = server.NewAPIServer(handler, 8081, true, "TODO", "TODO")

	return actuator, err
}

const (
	createEventAction = "Create"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

const (
	powerOnState    = "Powering on"
	powerOffState   = "Powering off"
	provioningState = "Provisioning"
	runningState    = "Running"
)

// Create creates a machine and is invoked by the Machine Controller
func (a *Actuator) Create(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("Creating machine %q for cluster %q.", machine.Name, cluster.Name)

	err := a.CreateMachine(cluster, machine)

	if err != nil {
		return errors.Errorf("failed to create machine: %+v", err)
	}

	return nil
}

// CreateMachine should extract data from spec and start the target machine via goimpi
func (a *Actuator) CreateMachine(cluster *machinev1.Cluster, machine *machinev1.Machine) error {

	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), createEventAction)
	}

	hostAddress := machineProviderConfig.Ipmi.HostAddress
	username := machineProviderConfig.Ipmi.Username
	password := machineProviderConfig.Ipmi.Password
	bootDevice := machineProviderConfig.Ipmi.BootDevice

	c := &goipmi.Connection{
		Hostname:  hostAddress,
		Username:  username,
		Password:  password,
		Interface: "lanplus",
	}

	ipmiClient, err := goipmi.NewClient(c)

	if err != nil {
		glog.Errorf("Error connecting to machine via IPMI: %v", err)
		return err
	}

	err = ipmiClient.Open()

	if err != nil {
		glog.Errorf("Error opening connection to machine via IPMI: %v", err)
		return err
	}

	if bootDevice == "pxe" {
		err = ipmiClient.SetBootDevice(goipmi.BootDevicePxe)

		if err != nil {
			glog.Errorf("Error setting machine to PXE boot via IPMI: %v", err)
			return err
		}
	} else {
		// Otherwise just do remote CD-ROM for now
		err = ipmiClient.SetBootDevice(goipmi.BootDeviceRemoteCdrom)

		if err != nil {
			glog.Errorf("Error setting machine to remote CD-ROM boot via IPMI: %v", err)
			return err
		}
	}

	// Power cycle machine
	err = ipmiClient.Control(goipmi.ControlPowerCycle)

	if err != nil {
		// Try just powering-up instead
		err = ipmiClient.Control(goipmi.ControlPowerUp)

		a.updateStatus(machine, powerOnState)

		if err != nil {
			glog.Errorf("Error powering-up machine via IPMI: %v", err)
			return err
		}
	}

	time.Sleep(10 * time.Second)
	a.updateStatus(machine, provioningState)

	defer ipmiClient.Close()

	return nil
}

func (a *Actuator) getIgnition(signature string) (string, error) {
	// TODO: Use k8s library to get ignition file path from etcd,
	// then read actual ignition file from...somewhere?
	return "{'fake': 'ignition'}", nil
}

// ProviderConfigMachine gets the machine provider config MachineSetSpec from the
// specified cluster-api MachineSpec.
func ProviderConfigMachine(codec codec, ms *machinev1.MachineSpec) (*providerconfigv1.BaremetalMachineProviderSpec, error) {
	providerSpec := ms.ProviderSpec
	if providerSpec.Value == nil {
		return nil, fmt.Errorf("no Value in ProviderConfig")
	}

	var config providerconfigv1.BaremetalMachineProviderSpec
	if err := codec.DecodeFromProviderSpec(providerSpec, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Set corresponding event based on error. It also returns the original error
// for convenience, so callers can do "return handleMachineError(...)".
func (a *Actuator) handleMachineError(machine *machinev1.Machine, err *apierrors.MachineError, eventAction string) error {
	if eventAction != noEventAction {
		a.eventRecorder.Eventf(machine, corev1.EventTypeWarning, "Failed"+eventAction, "%v", err.Reason)
	}

	glog.Errorf("Machine error: %v", err.Message)
	return err
}

// Delete : empty method
func (a *Actuator) Delete(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	glog.Infof("Deleting machine %q for cluster %q.", machine.Name, cluster.Name)

	err := a.DeleteMachine(cluster, machine)

	if err != nil {
		return errors.Errorf("failed to delete machine: %+v", err)
	}

	return nil
}

// DeleteMachine should extract data from etc and power off the target machine via goimpi
func (a *Actuator) DeleteMachine(cluster *machinev1.Cluster, machine *machinev1.Machine) error {

	machineProviderConfig, err := ProviderConfigMachine(a.codec, &machine.Spec)
	if err != nil {
		return a.handleMachineError(machine, apierrors.InvalidMachineConfiguration("error getting machineProviderConfig from spec: %v", err), createEventAction)
	}

	hostAddress := machineProviderConfig.Ipmi.HostAddress
	username := machineProviderConfig.Ipmi.Username
	password := machineProviderConfig.Ipmi.Password

	c := &goipmi.Connection{
		Hostname:  hostAddress,
		Username:  username,
		Password:  password,
		Interface: "lanplus",
	}

	ipmiClient, err := goipmi.NewClient(c)

	if err != nil {
		glog.Errorf("Error connecting to machine via IPMI: %v", err)
		return err
	}

	err = ipmiClient.Open()

	if err != nil {
		glog.Errorf("Error opening connection to machine via IPMI: %v", err)
		return err
	}

	// Power off machine
	err = ipmiClient.Control(goipmi.ControlPowerDown)

	a.updateStatus(machine, powerOffState)
	time.Sleep(10 * time.Second)

	if err != nil {
		glog.Errorf("Error powering off machine via IPMI: %v", err)
		return err
	}

	defer ipmiClient.Close()

	return nil
}

// Update : empty method
func (a *Actuator) Update(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) error {
	return nil
}

// Exists : empty method
func (a *Actuator) Exists(context context.Context, cluster *machinev1.Cluster, machine *machinev1.Machine) (bool, error) {
	return false, nil
}

// updateStatus updates a machine object's status.
func (a *Actuator) updateStatus(machine *machinev1.Machine, state string) error {
	glog.Infof("Updating status for %s", machine.Name)

	status, err := ProviderStatusFromMachine(a.codec, machine)
	if err != nil {
		glog.Errorf("Unable to get provider status from machine: %v", err)
		return err
	}

	// Update the libvirt provider status in-place.
	if err := UpdateProviderStatus(status, state); err != nil {
		glog.Errorf("Unable to update provider status: %v", err)
		return err
	}

	return nil
}

// UpdateProviderStatus updates the provider status in-place
func UpdateProviderStatus(status *providerconfigv1.BaremetalMachineProviderStatus, state string) error {

	status.Status = state

	return nil
}

// EncodeProviderStatus encodes a libvirt provider
// status as a runtime.RawExtension for inclusion in a MachineStatus
// object.
func EncodeProviderStatus(codec codec, status *providerconfigv1.BaremetalMachineProviderStatus) (*runtime.RawExtension, error) {
	return codec.EncodeProviderStatus(status)
}

// ProviderStatusFromMachine deserializes a libvirt provider status
// from a machine object.
func ProviderStatusFromMachine(codec codec, machine *machinev1.Machine) (*providerconfigv1.BaremetalMachineProviderStatus, error) {
	status := &providerconfigv1.BaremetalMachineProviderStatus{}
	var err error
	if machine.Status.ProviderStatus != nil {
		err = codec.DecodeProviderStatus(machine.Status.ProviderStatus, status)
	}

	return status, err
}
