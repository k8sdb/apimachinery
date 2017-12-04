/*
Copyright 2017 The KubeDB Authors.

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

package fake

import (
	v1alpha1 "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMongoDBs implements MongoDBInterface
type FakeMongoDBs struct {
	Fake *FakeKubedbV1alpha1
	ns   string
}

var mongodbsResource = schema.GroupVersionResource{Group: "kubedb.com", Version: "v1alpha1", Resource: "mongodbs"}

var mongodbsKind = schema.GroupVersionKind{Group: "kubedb.com", Version: "v1alpha1", Kind: "MongoDB"}

// Get takes name of the mongoDB, and returns the corresponding mongoDB object, and an error if there is any.
func (c *FakeMongoDBs) Get(name string, options v1.GetOptions) (result *v1alpha1.MongoDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mongodbsResource, c.ns, name), &v1alpha1.MongoDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoDB), err
}

// List takes label and field selectors, and returns the list of MongoDBs that match those selectors.
func (c *FakeMongoDBs) List(opts v1.ListOptions) (result *v1alpha1.MongoDBList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mongodbsResource, mongodbsKind, c.ns, opts), &v1alpha1.MongoDBList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MongoDBList{}
	for _, item := range obj.(*v1alpha1.MongoDBList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mongoDBs.
func (c *FakeMongoDBs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mongodbsResource, c.ns, opts))

}

// Create takes the representation of a mongoDB and creates it.  Returns the server's representation of the mongoDB, and an error, if there is any.
func (c *FakeMongoDBs) Create(mongoDB *v1alpha1.MongoDB) (result *v1alpha1.MongoDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mongodbsResource, c.ns, mongoDB), &v1alpha1.MongoDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoDB), err
}

// Update takes the representation of a mongoDB and updates it. Returns the server's representation of the mongoDB, and an error, if there is any.
func (c *FakeMongoDBs) Update(mongoDB *v1alpha1.MongoDB) (result *v1alpha1.MongoDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mongodbsResource, c.ns, mongoDB), &v1alpha1.MongoDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoDB), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMongoDBs) UpdateStatus(mongoDB *v1alpha1.MongoDB) (*v1alpha1.MongoDB, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mongodbsResource, "status", c.ns, mongoDB), &v1alpha1.MongoDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoDB), err
}

// Delete takes name of the mongoDB and deletes it. Returns an error if one occurs.
func (c *FakeMongoDBs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mongodbsResource, c.ns, name), &v1alpha1.MongoDB{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMongoDBs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mongodbsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.MongoDBList{})
	return err
}

// Patch applies the patch and returns the patched mongoDB.
func (c *FakeMongoDBs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mongodbsResource, c.ns, name, data, subresources...), &v1alpha1.MongoDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MongoDB), err
}
