package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaremetalMachineProviderConfig
type BaremetalMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	IgnKey string `json:"ignKey"`

	Spec   BaremetalMachineProviderConfigSpec   `json:"spec"`
	Status BaremetalMachineProviderConfigStatus `json:"status"`
}

// BaremetalMachineProviderConfigSpec
type BaremetalMachineProviderConfigSpec struct {
	Something string `json:"something"`
}

// BaremetalMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
type BaremetalMachineProviderConfigStatus struct {
	metav1.TypeMeta `json:",inline"`

	Status string `json:"status"`
}

// BaremetalClusterProviderConfig is the type that will be embedded in a Cluster.Spec.ProviderSpec field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BaremetalClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaremetalProviderConfigList contains a list of BaremetalProviderConfig
type BaremetalMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BaremetalMachineProviderConfig `json:"items"`
}
