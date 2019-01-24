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
