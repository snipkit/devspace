package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RunnerConfig holds the config the runner retrieves from Khulnasoft
// +subresource-request
type RunnerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// TokenCaCert is the certificate authority the Khulnasoft tokens will
	// be signed with
	// +optional
	TokenCaCert []byte `json:"tokenCaCert,omitempty"`
}
