// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"

	v1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	scheme "dev.khulnasoft.com/api/pkg/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// EventsGetter has a method to return a EventInterface.
// A group's client should implement this interface.
type EventsGetter interface {
	Events() EventInterface
}

// EventInterface has methods to work with Event resources.
type EventInterface interface {
	Create(ctx context.Context, event *v1.Event, opts metav1.CreateOptions) (*v1.Event, error)
	Update(ctx context.Context, event *v1.Event, opts metav1.UpdateOptions) (*v1.Event, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, event *v1.Event, opts metav1.UpdateOptions) (*v1.Event, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Event, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.EventList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Event, err error)
	EventExpansion
}

// events implements EventInterface
type events struct {
	*gentype.ClientWithList[*v1.Event, *v1.EventList]
}

// newEvents returns a Events
func newEvents(c *ManagementV1Client) *events {
	return &events{
		gentype.NewClientWithList[*v1.Event, *v1.EventList](
			"events",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1.Event { return &v1.Event{} },
			func() *v1.EventList { return &v1.EventList{} }),
	}
}
