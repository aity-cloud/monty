package v1beta1

import (
	montycorev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type Role struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              *montycorev1.Role `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Role `json:"items"`
}

// +kubebuilder:object:root=true
type RoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              *montycorev1.RoleBinding `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
type RoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&Role{}, &RoleList{},
		&RoleBinding{}, &RoleBindingList{},
	)
}
