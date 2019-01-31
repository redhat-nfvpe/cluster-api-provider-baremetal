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

	"github.com/golang/glog"
	server "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clusterclient "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

var MachineActuator *Actuator

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	clusterClient   clusterclient.Interface
	kubeClient      kubernetes.Interface
	eventRecorder   record.EventRecorder
	baremetalServer server.BaremetalServer
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClusterClient       clusterclient.Interface
	KubeClient          kubernetes.Interface
	EventRecorder       record.EventRecorder
	ServerListenAddress string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {

	actuator := &Actuator{
		clusterClient: params.ClusterClient,
		kubeClient:    params.KubeClient,
		eventRecorder: params.EventRecorder,
	}

	baremetalServerParams := server.BaremetalServerParams{
		ListenAddress:   params.ServerListenAddress,
		GetIgnitionFunc: actuator.getIgnition,
	}

	var err error

	actuator.baremetalServer, err = server.NewBaremetalServer(baremetalServerParams)

	return actuator, err
}

const (
	createEventAction = "Create"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

// Create creates a machine and is invoked by the Machine Controller
func (a *Actuator) Create(context context.Context, cluster *clusterv1.Cluster, machine *clusterv1.Machine) error {
	glog.Infof("Creating machine %q for cluster %q.", machine.Name, cluster.Name)

	return nil
}

func (a *Actuator) getIgnition(signature string) (string, error) {
	// TODO: Use k8s library to get ignition file path from etcd,
	// then read actual ignition file from...somewhere?
	return "{'fake': 'ignition'}", nil
}
