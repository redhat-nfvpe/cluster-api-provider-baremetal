package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BaremetalMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	IgnKey   string    `json:"ignKey"`
	Ignition *Ignition `json:"ignition"`
}

// Ignition contains location of ignition to be run during bootstrapping
type Ignition struct {
	// Ignition config to be run during bootstrapping
	UserDataSecret string `json:"userDataSecret"`
}

type BaremetalMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// BaremetalMachineID is a unique identifier for the baremetal machine
	BaremetalMachineID *string `json:"baremetalMachineID"`

	// BaremetalMachineState is the current state of the baremetal machine
	// TODO: define a state machine for baremetal
	BaremetalMachineState *string `json:"baremetalMachineState"`
}

// BaremetalMachineProviderConfigList contains a list of BaremetalMachineProviderConfig
type BaremetalMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BaremetalMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BaremetalMachineProviderConfig{}, &BaremetalMachineProviderConfigList{}, &BaremetalMachineProviderStatus{})
}
