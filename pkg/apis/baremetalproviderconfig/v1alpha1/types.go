package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaremetalProviderConfig
type BaremetalProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	IgnKey string `json:"ignKey"`

	Spec   BaremetalProviderConfigSpec   `json:"spec"`
	Status BaremetalProviderConfigStatus `json:"status"`
}

// BaremetalProviderConfigSpec
type BaremetalProviderConfigSpec struct {
	Something string `json:"something"`
}

// BaremetalProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
type BaremetalProviderConfigStatus struct {
	metav1.TypeMeta `json:",inline"`

	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaremetalProviderConfigList contains a list of BaremetalProviderConfig
type BaremetalProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BaremetalProviderConfig `json:"items"`
}
