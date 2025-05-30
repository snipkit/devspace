// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"

	v1 "dev.khulnasoft.com/api/pkg/apis/storage/v1"
	scheme "dev.khulnasoft.com/api/pkg/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// DevSpaceEnvironmentTemplatesGetter has a method to return a DevSpaceEnvironmentTemplateInterface.
// A group's client should implement this interface.
type DevSpaceEnvironmentTemplatesGetter interface {
	DevSpaceEnvironmentTemplates() DevSpaceEnvironmentTemplateInterface
}

// DevSpaceEnvironmentTemplateInterface has methods to work with DevSpaceEnvironmentTemplate resources.
type DevSpaceEnvironmentTemplateInterface interface {
	Create(ctx context.Context, devSpaceEnvironmentTemplate *v1.DevSpaceEnvironmentTemplate, opts metav1.CreateOptions) (*v1.DevSpaceEnvironmentTemplate, error)
	Update(ctx context.Context, devSpaceEnvironmentTemplate *v1.DevSpaceEnvironmentTemplate, opts metav1.UpdateOptions) (*v1.DevSpaceEnvironmentTemplate, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.DevSpaceEnvironmentTemplate, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.DevSpaceEnvironmentTemplateList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.DevSpaceEnvironmentTemplate, err error)
	DevSpaceEnvironmentTemplateExpansion
}

// devSpaceEnvironmentTemplates implements DevSpaceEnvironmentTemplateInterface
type devSpaceEnvironmentTemplates struct {
	*gentype.ClientWithList[*v1.DevSpaceEnvironmentTemplate, *v1.DevSpaceEnvironmentTemplateList]
}

// newDevSpaceEnvironmentTemplates returns a DevSpaceEnvironmentTemplates
func newDevSpaceEnvironmentTemplates(c *StorageV1Client) *devSpaceEnvironmentTemplates {
	return &devSpaceEnvironmentTemplates{
		gentype.NewClientWithList[*v1.DevSpaceEnvironmentTemplate, *v1.DevSpaceEnvironmentTemplateList](
			"devspaceenvironmenttemplates",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1.DevSpaceEnvironmentTemplate { return &v1.DevSpaceEnvironmentTemplate{} },
			func() *v1.DevSpaceEnvironmentTemplateList { return &v1.DevSpaceEnvironmentTemplateList{} }),
	}
}
