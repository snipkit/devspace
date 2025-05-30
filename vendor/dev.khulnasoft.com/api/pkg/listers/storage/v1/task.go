// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// TaskLister helps list Tasks.
// All objects returned here must be treated as read-only.
type TaskLister interface {
	// List lists all Tasks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Task, err error)
	// Get retrieves the Task from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Task, error)
	TaskListerExpansion
}

// taskLister implements the TaskLister interface.
type taskLister struct {
	listers.ResourceIndexer[*v1.Task]
}

// NewTaskLister returns a new TaskLister.
func NewTaskLister(indexer cache.Indexer) TaskLister {
	return &taskLister{listers.New[*v1.Task](indexer, v1.Resource("task"))}
}
