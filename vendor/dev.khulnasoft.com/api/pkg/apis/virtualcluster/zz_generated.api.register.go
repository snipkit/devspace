// Code generated by generator. DO NOT EDIT.

package virtualcluster

import (
	"context"
	"fmt"

	clusterv1 "dev.khulnasoft.com/agentapi/pkg/apis/khulnasoft/cluster/v1"
	"dev.khulnasoft.com/api/pkg/managerfactory"
	"dev.khulnasoft.com/apiserver/pkg/builders"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

type NewRESTFunc func(factory managerfactory.SharedManagerFactory) rest.Storage

var (
	VirtualclusterHelmReleaseStorage = builders.NewApiResourceWithStorage( // Resource status endpoint
		InternalHelmRelease,
		func() runtime.Object { return &HelmRelease{} },     // Register versioned resource
		func() runtime.Object { return &HelmReleaseList{} }, // Register versioned resource list
		NewHelmReleaseREST,
	)
	NewHelmReleaseREST = func(getter generic.RESTOptionsGetter) rest.Storage {
		return NewHelmReleaseRESTFunc(Factory)
	}
	NewHelmReleaseRESTFunc NewRESTFunc
	InternalHelmRelease    = builders.NewInternalResource(
		"helmreleases",
		"HelmRelease",
		func() runtime.Object { return &HelmRelease{} },
		func() runtime.Object { return &HelmReleaseList{} },
	)
	InternalHelmReleaseStatus = builders.NewInternalResourceStatus(
		"helmreleases",
		"HelmReleaseStatus",
		func() runtime.Object { return &HelmRelease{} },
		func() runtime.Object { return &HelmReleaseList{} },
	)
	// Registered resources and subresources
	ApiVersion = builders.NewApiGroup("virtualcluster.khulnasoft.com").WithKinds(
		InternalHelmRelease,
		InternalHelmReleaseStatus,
	)

	// Required by code generated by go2idl
	AddToScheme = (&runtime.SchemeBuilder{
		ApiVersion.SchemeBuilder.AddToScheme,
		RegisterDefaults,
	}).AddToScheme
	SchemeBuilder      = ApiVersion.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// +genclient
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HelmReleaseSpec   `json:"spec,omitempty"`
	Status            HelmReleaseStatus `json:"status,omitempty"`
}

type HelmReleaseSpec struct {
	clusterv1.HelmReleaseSpec `json:",inline"`
}

type HelmReleaseStatus struct {
	clusterv1.HelmReleaseStatus `json:",inline"`
}

// HelmRelease Functions and Structs
//
// +k8s:deepcopy-gen=false
type HelmReleaseStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type HelmReleaseStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

func (HelmRelease) NewStatus() interface{} {
	return HelmReleaseStatus{}
}

func (pc *HelmRelease) GetStatus() interface{} {
	return pc.Status
}

func (pc *HelmRelease) SetStatus(s interface{}) {
	pc.Status = s.(HelmReleaseStatus)
}

func (pc *HelmRelease) GetSpec() interface{} {
	return pc.Spec
}

func (pc *HelmRelease) SetSpec(s interface{}) {
	pc.Spec = s.(HelmReleaseSpec)
}

func (pc *HelmRelease) GetObjectMeta() *metav1.ObjectMeta {
	return &pc.ObjectMeta
}

func (pc *HelmRelease) SetGeneration(generation int64) {
	pc.ObjectMeta.Generation = generation
}

func (pc HelmRelease) GetGeneration() int64 {
	return pc.ObjectMeta.Generation
}

// Registry is an interface for things that know how to store HelmRelease.
// +k8s:deepcopy-gen=false
type HelmReleaseRegistry interface {
	ListHelmReleases(ctx context.Context, options *internalversion.ListOptions) (*HelmReleaseList, error)
	GetHelmRelease(ctx context.Context, id string, options *metav1.GetOptions) (*HelmRelease, error)
	CreateHelmRelease(ctx context.Context, id *HelmRelease) (*HelmRelease, error)
	UpdateHelmRelease(ctx context.Context, id *HelmRelease) (*HelmRelease, error)
	DeleteHelmRelease(ctx context.Context, id string) (bool, error)
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched types will panic.
func NewHelmReleaseRegistry(sp builders.StandardStorageProvider) HelmReleaseRegistry {
	return &storageHelmRelease{sp}
}

// Implement Registry
// storage puts strong typing around storage calls
// +k8s:deepcopy-gen=false
type storageHelmRelease struct {
	builders.StandardStorageProvider
}

func (s *storageHelmRelease) ListHelmReleases(ctx context.Context, options *internalversion.ListOptions) (*HelmReleaseList, error) {
	if options != nil && options.FieldSelector != nil && !options.FieldSelector.Empty() {
		return nil, fmt.Errorf("field selector not supported yet")
	}
	st := s.GetStandardStorage()
	obj, err := st.List(ctx, options)
	if err != nil {
		return nil, err
	}
	return obj.(*HelmReleaseList), err
}

func (s *storageHelmRelease) GetHelmRelease(ctx context.Context, id string, options *metav1.GetOptions) (*HelmRelease, error) {
	st := s.GetStandardStorage()
	obj, err := st.Get(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return obj.(*HelmRelease), nil
}

func (s *storageHelmRelease) CreateHelmRelease(ctx context.Context, object *HelmRelease) (*HelmRelease, error) {
	st := s.GetStandardStorage()
	obj, err := st.Create(ctx, object, nil, &metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*HelmRelease), nil
}

func (s *storageHelmRelease) UpdateHelmRelease(ctx context.Context, object *HelmRelease) (*HelmRelease, error) {
	st := s.GetStandardStorage()
	obj, _, err := st.Update(ctx, object.Name, rest.DefaultUpdatedObjectInfo(object), nil, nil, false, &metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}
	return obj.(*HelmRelease), nil
}

func (s *storageHelmRelease) DeleteHelmRelease(ctx context.Context, id string) (bool, error) {
	st := s.GetStandardStorage()
	_, sync, err := st.Delete(ctx, id, nil, &metav1.DeleteOptions{})
	return sync, err
}
