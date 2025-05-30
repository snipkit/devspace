// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// FeatureLister helps list Features.
// All objects returned here must be treated as read-only.
type FeatureLister interface {
	// List lists all Features in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Feature, err error)
	// Get retrieves the Feature from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Feature, error)
	FeatureListerExpansion
}

// featureLister implements the FeatureLister interface.
type featureLister struct {
	listers.ResourceIndexer[*v1.Feature]
}

// NewFeatureLister returns a new FeatureLister.
func NewFeatureLister(indexer cache.Indexer) FeatureLister {
	return &featureLister{listers.New[*v1.Feature](indexer, v1.Resource("feature"))}
}
