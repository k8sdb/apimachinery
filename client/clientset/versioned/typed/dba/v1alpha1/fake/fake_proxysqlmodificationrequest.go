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

package fake

import (
	v1alpha1 "kubedb.dev/apimachinery/apis/dba/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeProxySQLModificationRequests implements ProxySQLModificationRequestInterface
type FakeProxySQLModificationRequests struct {
	Fake *FakeDbaV1alpha1
	ns   string
}

var proxysqlmodificationrequestsResource = schema.GroupVersionResource{Group: "dba.kubedb.com", Version: "v1alpha1", Resource: "proxysqlmodificationrequests"}

var proxysqlmodificationrequestsKind = schema.GroupVersionKind{Group: "dba.kubedb.com", Version: "v1alpha1", Kind: "ProxySQLModificationRequest"}

// Get takes name of the proxySQLModificationRequest, and returns the corresponding proxySQLModificationRequest object, and an error if there is any.
func (c *FakeProxySQLModificationRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.ProxySQLModificationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(proxysqlmodificationrequestsResource, c.ns, name), &v1alpha1.ProxySQLModificationRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ProxySQLModificationRequest), err
}

// List takes label and field selectors, and returns the list of ProxySQLModificationRequests that match those selectors.
func (c *FakeProxySQLModificationRequests) List(opts v1.ListOptions) (result *v1alpha1.ProxySQLModificationRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(proxysqlmodificationrequestsResource, proxysqlmodificationrequestsKind, c.ns, opts), &v1alpha1.ProxySQLModificationRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ProxySQLModificationRequestList{ListMeta: obj.(*v1alpha1.ProxySQLModificationRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.ProxySQLModificationRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested proxySQLModificationRequests.
func (c *FakeProxySQLModificationRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(proxysqlmodificationrequestsResource, c.ns, opts))

}

// Create takes the representation of a proxySQLModificationRequest and creates it.  Returns the server's representation of the proxySQLModificationRequest, and an error, if there is any.
func (c *FakeProxySQLModificationRequests) Create(proxySQLModificationRequest *v1alpha1.ProxySQLModificationRequest) (result *v1alpha1.ProxySQLModificationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(proxysqlmodificationrequestsResource, c.ns, proxySQLModificationRequest), &v1alpha1.ProxySQLModificationRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ProxySQLModificationRequest), err
}

// Update takes the representation of a proxySQLModificationRequest and updates it. Returns the server's representation of the proxySQLModificationRequest, and an error, if there is any.
func (c *FakeProxySQLModificationRequests) Update(proxySQLModificationRequest *v1alpha1.ProxySQLModificationRequest) (result *v1alpha1.ProxySQLModificationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(proxysqlmodificationrequestsResource, c.ns, proxySQLModificationRequest), &v1alpha1.ProxySQLModificationRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ProxySQLModificationRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeProxySQLModificationRequests) UpdateStatus(proxySQLModificationRequest *v1alpha1.ProxySQLModificationRequest) (*v1alpha1.ProxySQLModificationRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(proxysqlmodificationrequestsResource, "status", c.ns, proxySQLModificationRequest), &v1alpha1.ProxySQLModificationRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ProxySQLModificationRequest), err
}

// Delete takes name of the proxySQLModificationRequest and deletes it. Returns an error if one occurs.
func (c *FakeProxySQLModificationRequests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(proxysqlmodificationrequestsResource, c.ns, name), &v1alpha1.ProxySQLModificationRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProxySQLModificationRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(proxysqlmodificationrequestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ProxySQLModificationRequestList{})
	return err
}

// Patch applies the patch and returns the patched proxySQLModificationRequest.
func (c *FakeProxySQLModificationRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ProxySQLModificationRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(proxysqlmodificationrequestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ProxySQLModificationRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ProxySQLModificationRequest), err
}
