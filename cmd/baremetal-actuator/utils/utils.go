package utils

import (
	"fmt"
	"io/ioutil"

	apiv1 "k8s.io/api/core/v1"

	machinev1 "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	machineactuator "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/actuators/machine"
	test "github.com/redhat-nfvpe/cluster-api-provider-baremetal/test"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/record"
)

func CreateActuator(machine *machinev1.Machine, userData *apiv1.Secret) *machineactuator.Actuator {
	objList := []runtime.Object{}
	if userData != nil {
		objList = append(objList, userData)
	}
	fakeKubeClient := kubernetesfake.NewSimpleClientset(objList...)

	params := machineactuator.ActuatorParams{
		ClusterClient:       test.NewSimpleClientset(machine),
		KubeClient:          fakeKubeClient,
		EventRecorder:       &record.FakeRecorder{},
		ServerListenAddress: "localhost:8081",
	}

	actuator, _ := machineactuator.NewActuator(params)
	return actuator
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
