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

// FakeLicenseTokens implements LicenseTokenInterface
type FakeLicenseTokens struct {
	Fake *FakeManagementV1
}

var licensetokensResource = v1.SchemeGroupVersion.WithResource("licensetokens")

var licensetokensKind = v1.SchemeGroupVersion.WithKind("LicenseToken")

// Get takes name of the licenseToken, and returns the corresponding licenseToken object, and an error if there is any.
func (c *FakeLicenseTokens) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.LicenseToken, err error) {
	emptyResult := &v1.LicenseToken{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(licensetokensResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.LicenseToken), err
}

// List takes label and field selectors, and returns the list of LicenseTokens that match those selectors.
func (c *FakeLicenseTokens) List(ctx context.Context, opts metav1.ListOptions) (result *v1.LicenseTokenList, err error) {
	emptyResult := &v1.LicenseTokenList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(licensetokensResource, licensetokensKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.LicenseTokenList{ListMeta: obj.(*v1.LicenseTokenList).ListMeta}
	for _, item := range obj.(*v1.LicenseTokenList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested licenseTokens.
func (c *FakeLicenseTokens) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(licensetokensResource, opts))
}

// Create takes the representation of a licenseToken and creates it.  Returns the server's representation of the licenseToken, and an error, if there is any.
func (c *FakeLicenseTokens) Create(ctx context.Context, licenseToken *v1.LicenseToken, opts metav1.CreateOptions) (result *v1.LicenseToken, err error) {
	emptyResult := &v1.LicenseToken{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(licensetokensResource, licenseToken, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.LicenseToken), err
}

// Update takes the representation of a licenseToken and updates it. Returns the server's representation of the licenseToken, and an error, if there is any.
func (c *FakeLicenseTokens) Update(ctx context.Context, licenseToken *v1.LicenseToken, opts metav1.UpdateOptions) (result *v1.LicenseToken, err error) {
	emptyResult := &v1.LicenseToken{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(licensetokensResource, licenseToken, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.LicenseToken), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLicenseTokens) UpdateStatus(ctx context.Context, licenseToken *v1.LicenseToken, opts metav1.UpdateOptions) (result *v1.LicenseToken, err error) {
	emptyResult := &v1.LicenseToken{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(licensetokensResource, "status", licenseToken, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.LicenseToken), err
}

// Delete takes name of the licenseToken and deletes it. Returns an error if one occurs.
func (c *FakeLicenseTokens) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(licensetokensResource, name, opts), &v1.LicenseToken{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLicenseTokens) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(licensetokensResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1.LicenseTokenList{})
	return err
}

// Patch applies the patch and returns the patched licenseToken.
func (c *FakeLicenseTokens) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.LicenseToken, err error) {
	emptyResult := &v1.LicenseToken{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(licensetokensResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.LicenseToken), err
}
