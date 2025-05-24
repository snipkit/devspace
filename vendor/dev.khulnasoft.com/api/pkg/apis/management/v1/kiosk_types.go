package v1

import (
	clusterv1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster/v1"
	agentstoragev1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/storage/v1"
	uiv1 "dev.khulnasoft.com/api/pkg/apis/ui/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This file is just used as a collector for kiosk objects we want to generate stuff for

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kiosk holds the kiosk types
// +k8s:openapi-gen=true
// +resource:path=kiosk,rest=KioskREST
type Kiosk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KioskSpec   `json:"spec,omitempty"`
	Status KioskStatus `json:"status,omitempty"`
}

type KioskSpec struct {
	// cluster.khulnasoft.com
	HelmRelease     clusterv1.HelmRelease     `json:"helmRelease,omitempty"`
	SleepModeConfig clusterv1.SleepModeConfig `json:"sleepModeConfig,omitempty"`
	ChartInfo       clusterv1.ChartInfo       `json:"chartInfo,omitempty"`

	// storage.khulnasoft.com
	StorageClusterQuota agentstoragev1.ClusterQuota `json:"storageClusterQuota,omitempty"`

	// ui.khulnasoft.com
	UISettings uiv1.UISettings `json:"UISettings,omitempty"`

	License License `json:"license,omitempty"`
}

type KioskStatus struct {
}
