// Package v1alpha1 contains API Schema definitions for the config v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=storage.khulnasoft.com
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "storage.khulnasoft.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme

	// SchemeGroupVersion is a shim that expect this to be present in the api package
	SchemeGroupVersion = GroupVersion
)

type AccessAccessor interface {
	GetAccess() []Access
	SetAccess(access []Access)

	GetOwner() *UserOrTeam
	SetOwner(userOrTeam *UserOrTeam)
}

type VersionsAccessor interface {
	GetVersions() []VersionAccessor
}

type VersionAccessor interface {
	GetVersion() string
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
