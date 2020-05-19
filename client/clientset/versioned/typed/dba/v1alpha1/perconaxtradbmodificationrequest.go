/*
Copyright The KubeDB Authors.

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
	"time"

	v1alpha1 "kubedb.dev/apimachinery/apis/dba/v1alpha1"
	scheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PerconaXtraDBModificationRequestsGetter has a method to return a PerconaXtraDBModificationRequestInterface.
// A group's client should implement this interface.
type PerconaXtraDBModificationRequestsGetter interface {
	PerconaXtraDBModificationRequests(namespace string) PerconaXtraDBModificationRequestInterface
}

// PerconaXtraDBModificationRequestInterface has methods to work with PerconaXtraDBModificationRequest resources.
type PerconaXtraDBModificationRequestInterface interface {
	Create(*v1alpha1.PerconaXtraDBModificationRequest) (*v1alpha1.PerconaXtraDBModificationRequest, error)
	Update(*v1alpha1.PerconaXtraDBModificationRequest) (*v1alpha1.PerconaXtraDBModificationRequest, error)
	UpdateStatus(*v1alpha1.PerconaXtraDBModificationRequest) (*v1alpha1.PerconaXtraDBModificationRequest, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PerconaXtraDBModificationRequest, error)
	List(opts v1.ListOptions) (*v1alpha1.PerconaXtraDBModificationRequestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PerconaXtraDBModificationRequest, err error)
	PerconaXtraDBModificationRequestExpansion
}

// perconaXtraDBModificationRequests implements PerconaXtraDBModificationRequestInterface
type perconaXtraDBModificationRequests struct {
	client rest.Interface
	ns     string
}

// newPerconaXtraDBModificationRequests returns a PerconaXtraDBModificationRequests
func newPerconaXtraDBModificationRequests(c *DbaV1alpha1Client, namespace string) *perconaXtraDBModificationRequests {
	return &perconaXtraDBModificationRequests{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the perconaXtraDBModificationRequest, and returns the corresponding perconaXtraDBModificationRequest object, and an error if there is any.
func (c *perconaXtraDBModificationRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.PerconaXtraDBModificationRequest, err error) {
	result = &v1alpha1.PerconaXtraDBModificationRequest{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PerconaXtraDBModificationRequests that match those selectors.
func (c *perconaXtraDBModificationRequests) List(opts v1.ListOptions) (result *v1alpha1.PerconaXtraDBModificationRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PerconaXtraDBModificationRequestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested perconaXtraDBModificationRequests.
func (c *perconaXtraDBModificationRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a perconaXtraDBModificationRequest and creates it.  Returns the server's representation of the perconaXtraDBModificationRequest, and an error, if there is any.
func (c *perconaXtraDBModificationRequests) Create(perconaXtraDBModificationRequest *v1alpha1.PerconaXtraDBModificationRequest) (result *v1alpha1.PerconaXtraDBModificationRequest, err error) {
	result = &v1alpha1.PerconaXtraDBModificationRequest{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		Body(perconaXtraDBModificationRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a perconaXtraDBModificationRequest and updates it. Returns the server's representation of the perconaXtraDBModificationRequest, and an error, if there is any.
func (c *perconaXtraDBModificationRequests) Update(perconaXtraDBModificationRequest *v1alpha1.PerconaXtraDBModificationRequest) (result *v1alpha1.PerconaXtraDBModificationRequest, err error) {
	result = &v1alpha1.PerconaXtraDBModificationRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		Name(perconaXtraDBModificationRequest.Name).
		Body(perconaXtraDBModificationRequest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *perconaXtraDBModificationRequests) UpdateStatus(perconaXtraDBModificationRequest *v1alpha1.PerconaXtraDBModificationRequest) (result *v1alpha1.PerconaXtraDBModificationRequest, err error) {
	result = &v1alpha1.PerconaXtraDBModificationRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		Name(perconaXtraDBModificationRequest.Name).
		SubResource("status").
		Body(perconaXtraDBModificationRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the perconaXtraDBModificationRequest and deletes it. Returns an error if one occurs.
func (c *perconaXtraDBModificationRequests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *perconaXtraDBModificationRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched perconaXtraDBModificationRequest.
func (c *perconaXtraDBModificationRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PerconaXtraDBModificationRequest, err error) {
	result = &v1alpha1.PerconaXtraDBModificationRequest{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("perconaxtradbmodificationrequests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
