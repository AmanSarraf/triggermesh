/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/triggermesh/triggermesh/pkg/apis/sources/v1alpha1"
	scheme "github.com/triggermesh/triggermesh/pkg/client/generated/clientset/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GoogleCloudBillingSourcesGetter has a method to return a GoogleCloudBillingSourceInterface.
// A group's client should implement this interface.
type GoogleCloudBillingSourcesGetter interface {
	GoogleCloudBillingSources(namespace string) GoogleCloudBillingSourceInterface
}

// GoogleCloudBillingSourceInterface has methods to work with GoogleCloudBillingSource resources.
type GoogleCloudBillingSourceInterface interface {
	Create(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.CreateOptions) (*v1alpha1.GoogleCloudBillingSource, error)
	Update(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.UpdateOptions) (*v1alpha1.GoogleCloudBillingSource, error)
	UpdateStatus(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.UpdateOptions) (*v1alpha1.GoogleCloudBillingSource, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.GoogleCloudBillingSource, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.GoogleCloudBillingSourceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.GoogleCloudBillingSource, err error)
	GoogleCloudBillingSourceExpansion
}

// googleCloudBillingSources implements GoogleCloudBillingSourceInterface
type googleCloudBillingSources struct {
	client rest.Interface
	ns     string
}

// newGoogleCloudBillingSources returns a GoogleCloudBillingSources
func newGoogleCloudBillingSources(c *SourcesV1alpha1Client, namespace string) *googleCloudBillingSources {
	return &googleCloudBillingSources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the googleCloudBillingSource, and returns the corresponding googleCloudBillingSource object, and an error if there is any.
func (c *googleCloudBillingSources) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.GoogleCloudBillingSource, err error) {
	result = &v1alpha1.GoogleCloudBillingSource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GoogleCloudBillingSources that match those selectors.
func (c *googleCloudBillingSources) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.GoogleCloudBillingSourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.GoogleCloudBillingSourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested googleCloudBillingSources.
func (c *googleCloudBillingSources) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a googleCloudBillingSource and creates it.  Returns the server's representation of the googleCloudBillingSource, and an error, if there is any.
func (c *googleCloudBillingSources) Create(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.CreateOptions) (result *v1alpha1.GoogleCloudBillingSource, err error) {
	result = &v1alpha1.GoogleCloudBillingSource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(googleCloudBillingSource).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a googleCloudBillingSource and updates it. Returns the server's representation of the googleCloudBillingSource, and an error, if there is any.
func (c *googleCloudBillingSources) Update(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.UpdateOptions) (result *v1alpha1.GoogleCloudBillingSource, err error) {
	result = &v1alpha1.GoogleCloudBillingSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		Name(googleCloudBillingSource.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(googleCloudBillingSource).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *googleCloudBillingSources) UpdateStatus(ctx context.Context, googleCloudBillingSource *v1alpha1.GoogleCloudBillingSource, opts v1.UpdateOptions) (result *v1alpha1.GoogleCloudBillingSource, err error) {
	result = &v1alpha1.GoogleCloudBillingSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		Name(googleCloudBillingSource.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(googleCloudBillingSource).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the googleCloudBillingSource and deletes it. Returns an error if one occurs.
func (c *googleCloudBillingSources) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *googleCloudBillingSources) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched googleCloudBillingSource.
func (c *googleCloudBillingSources) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.GoogleCloudBillingSource, err error) {
	result = &v1alpha1.GoogleCloudBillingSource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("googlecloudbillingsources").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}