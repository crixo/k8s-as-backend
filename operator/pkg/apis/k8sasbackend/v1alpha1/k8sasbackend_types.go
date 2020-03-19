package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// K8sAsBackendSpec defines the desired state of K8sAsBackend
type K8sAsBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32 `json:"size"`
}

// K8sAsBackendStatus defines the observed state of K8sAsBackend
type K8sAsBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +optional
	// +listType=set
	AdmissionWebhookPems []string `json:"pems,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// K8sAsBackend is the Schema for the k8sasbackends API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=k8sasbackends,scope=Namespaced
type K8sAsBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   K8sAsBackendSpec   `json:"spec,omitempty"`
	Status K8sAsBackendStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// K8sAsBackendList contains a list of K8sAsBackend
type K8sAsBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []K8sAsBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&K8sAsBackend{}, &K8sAsBackendList{})
}
