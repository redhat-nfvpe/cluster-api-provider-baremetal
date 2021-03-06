package utils

import (
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/rest"

	apiv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	machineactuator "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/actuators/machine"
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/record"
)

func CreateActuator(machine *machinev1.Machine, userData *apiv1.Secret) (*machineactuator.Actuator, error) {
	objList := []runtime.Object{}
	if userData != nil {
		objList = append(objList, userData)
	}

	fakeClient := fake.NewFakeClient(objList...)
	fakeKubeClient := kubernetesfake.NewSimpleClientset(objList...)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		glog.Fatal(err)
	}

	// Create a fake config
	config := rest.Config{
		Host: "api.test.something.com",
		TLSClientConfig: rest.TLSClientConfig{
			ServerName: "https://api.test.something.com:6443",
		},
	}

	params := machineactuator.ActuatorParams{
		Client:        fakeClient,
		Config:        config,
		Codec:         codec,
		KubeClient:    fakeKubeClient,
		EventRecorder: &record.FakeRecorder{},
	}

	actuator, err := machineactuator.NewActuator(params)

	if err != nil {
		return nil, err
	}

	return actuator, nil
}

func ReadClusterResources(clusterLoc, machineLoc, userDataLoc string) (*machinev1.Cluster, *machinev1.Machine, *apiv1.Secret, error) {
	machine := &machinev1.Machine{}
	{
		bytes, err := ioutil.ReadFile(machineLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read machine manifest %q: %v", machineLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &machine); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal machine manifest %q: %v", machineLoc, err)
		}
	}

	cluster := &machinev1.Cluster{}
	{
		bytes, err := ioutil.ReadFile(clusterLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read cluster manifest %q: %v", clusterLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &cluster); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal cluster manifest %q: %v", clusterLoc, err)
		}
	}

	var userDataSecret *apiv1.Secret
	if userDataLoc != "" {
		userDataSecret = &apiv1.Secret{}
		bytes, err := ioutil.ReadFile(userDataLoc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read user data manifest %q: %v", userDataLoc, err)
		}

		if err = yaml.Unmarshal(bytes, &userDataSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal user data manifest %q: %v", userDataLoc, err)
		}
	}

	return cluster, machine, userDataSecret, nil
}
