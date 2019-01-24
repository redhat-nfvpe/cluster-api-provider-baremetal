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
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	clusterclient "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
)

var MachineActuator *Actuator

// Actuator is responsible for performing machine reconciliation
type Actuator struct {
	clusterClient   clusterclient.Interface
	kubeClient      kubernetes.Interface
	codec           codec
	eventRecorder   record.EventRecorder
	baremetalServer BaremetalServer
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClusterClient       clusterclient.Interface
	KubeClient          kubernetes.Interface
	Codec               codec
	EventRecorder       record.EventRecorder
	ServerListenAddress string
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {

	actuator := &Actuator{
		clusterClient: params.ClusterClient,
		kubeClient:    params.KubeClient,
		codec:         params.Codec,
		eventRecorder: params.EventRecorder,
	}

	baremetalServerParams := BaremetalServerParams{
		ListenAddress:   params.ServerListenAddress,
		GetIgnitionFunc: actuator.getIgnition,
	}

	actuator.baremetalServer = server.NewBaremetalServer(baremetalServerParams)

	return actuator, nil
}

const (
	createEventAction = "Create"
	deleteEventAction = "Delete"
	noEventAction     = ""
)

func (a *Actuator) getIgnition(signature string) (string, error) {
	// TODO: Use k8s library to get ignition file path from etcd,
	// then read actual ignition file from...somewhere?
	return "{'fake': 'ignition'}", nil
}
