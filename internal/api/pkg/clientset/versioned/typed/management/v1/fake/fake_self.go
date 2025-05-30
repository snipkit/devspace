// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "dev.khulnasoft.com/api/pkg/apis/management/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSelves implements SelfInterface
type FakeSelves struct {
	Fake *FakeManagementV1
}

var selvesResource = v1.SchemeGroupVersion.WithResource("selves")

var selvesKind = v1.SchemeGroupVersion.WithKind("Self")

// Get takes name of the self, and returns the corresponding self object, and an error if there is any.
func (c *FakeSelves) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Self, err error) {
	emptyResult := &v1.Self{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(selvesResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.Self), err
}

// List takes label and field selectors, and returns the list of Selves that match those selectors.
func (c *FakeSelves) List(ctx context.Context, opts metav1.ListOptions) (result *v1.SelfList, err error) {
	emptyResult := &v1.SelfList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(selvesResource, selvesKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.SelfList{ListMeta: obj.(*v1.SelfList).ListMeta}
	for _, item := range obj.(*v1.SelfList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested selves.
func (c *FakeSelves) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(selvesResource, opts))
}

// Create takes the representation of a self and creates it.  Returns the server's representation of the self, and an error, if there is any.
func (c *FakeSelves) Create(ctx context.Context, self *v1.Self, opts metav1.CreateOptions) (result *v1.Self, err error) {
	emptyResult := &v1.Self{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(selvesResource, self, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.Self), err
}

// Update takes the representation of a self and updates it. Returns the server's representation of the self, and an error, if there is any.
func (c *FakeSelves) Update(ctx context.Context, self *v1.Self, opts metav1.UpdateOptions) (result *v1.Self, err error) {
	emptyResult := &v1.Self{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(selvesResource, self, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.Self), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSelves) UpdateStatus(ctx context.Context, self *v1.Self, opts metav1.UpdateOptions) (result *v1.Self, err error) {
	emptyResult := &v1.Self{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(selvesResource, "status", self, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.Self), err
}

// Delete takes name of the self and deletes it. Returns an error if one occurs.
func (c *FakeSelves) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(selvesResource, name, opts), &v1.Self{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSelves) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(selvesResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1.SelfList{})
	return err
}

// Patch applies the patch and returns the patched self.
func (c *FakeSelves) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Self, err error) {
	emptyResult := &v1.Self{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(selvesResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.Self), err
}
