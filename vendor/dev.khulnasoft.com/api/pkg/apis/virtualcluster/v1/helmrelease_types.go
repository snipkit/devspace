package v1

import (
	clusterv1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=helmreleases,rest=HelmReleaseREST
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

type HelmReleaseSpec struct {
	clusterv1.HelmReleaseSpec `json:",inline"`
}

type HelmReleaseStatus struct {
	clusterv1.HelmReleaseStatus `json:",inline"`
}
