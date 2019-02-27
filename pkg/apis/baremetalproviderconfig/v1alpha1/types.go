package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaremetalMachineProviderSpec
type BaremetalMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	IgnKey string `json:"ignKey"`

	Ipmi *Ipmi `json:"ipmi"`
}

// Image contains the info for the actuator to know what OS to place (and
// how to place it) on the actual baremetal machine
// type Image struct {
// 	Type string `json:"type"`

// }

// Ipmi contains the info for the actuator to control the actual baremetal machine
type Ipmi struct {
	HostAddress string `json:"hostAddress"`
	// FIXME: store these as a secret?
	Username   string `json:"username"`
	Password   string `json:"password"`
	BootDevice string `json:"bootDevice"`
}

// BaremetalMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BaremetalMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// BaremetalProviderConfigList contains a list of BaremetalProviderConfig
type BaremetalMachineProviderSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BaremetalMachineProviderSpec `json:"items"`
}

// BaremetalClusterProviderSpec is the type that will be embedded in a Cluster.Spec.ProviderSpec field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BaremetalClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BaremetalClusterProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}
