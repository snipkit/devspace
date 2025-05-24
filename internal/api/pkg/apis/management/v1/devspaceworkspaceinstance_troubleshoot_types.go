package v1

import (
	storagev1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +subresource-request
type DevSpaceWorkspaceInstanceTroubleshoot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// State holds the workspaces state as given by 'devspace export'
	// +optional
	State string `json:"state,omitempty"`

	// Workspace holds the workspace's instance object data
	// +optional
	Workspace *DevSpaceWorkspaceInstance `json:"workspace,omitempty"`

	// Template holds the workspace instance's template used to create it.
	// This is the raw template, not the rendered one.
	// +optional
	Template *storagev1.DevSpaceWorkspaceTemplate `json:"template,omitempty"`

	// Pods is a list of pod objects that are linked to the workspace.
	// +optional
	Pods []corev1.Pod `json:"pods,omitempty"`

	// PVCs is a list of PVC objects that are linked to the workspace.
	// +optional
	PVCs []corev1.PersistentVolumeClaim `json:"pvcs,omitempty"`

	// Errors is a list of errors that occurred while trying to collect
	// informations for troubleshooting.
	// +optional
	Errors []string `json:"errors,omitempty"`
}
