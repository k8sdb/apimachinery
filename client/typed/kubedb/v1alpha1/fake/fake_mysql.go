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
	v1alpha1 "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMysqls implements MysqlInterface
type FakeMysqls struct {
	Fake *FakeKubedbV1alpha1
	ns   string
}

var mysqlsResource = schema.GroupVersionResource{Group: "kubedb.com", Version: "v1alpha1", Resource: "mysqls"}

var mysqlsKind = schema.GroupVersionKind{Group: "kubedb.com", Version: "v1alpha1", Kind: "Mysql"}

func (c *FakeMysqls) Create(mysql *v1alpha1.Mysql) (result *v1alpha1.Mysql, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mysqlsResource, c.ns, mysql), &v1alpha1.Mysql{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Mysql), err
}

func (c *FakeMysqls) Update(mysql *v1alpha1.Mysql) (result *v1alpha1.Mysql, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mysqlsResource, c.ns, mysql), &v1alpha1.Mysql{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Mysql), err
}

func (c *FakeMysqls) UpdateStatus(mysql *v1alpha1.Mysql) (*v1alpha1.Mysql, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(mysqlsResource, "status", c.ns, mysql), &v1alpha1.Mysql{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Mysql), err
}

func (c *FakeMysqls) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mysqlsResource, c.ns, name), &v1alpha1.Mysql{})

	return err
}

func (c *FakeMysqls) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mysqlsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.MysqlList{})
	return err
}

func (c *FakeMysqls) Get(name string, options v1.GetOptions) (result *v1alpha1.Mysql, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mysqlsResource, c.ns, name), &v1alpha1.Mysql{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Mysql), err
}

func (c *FakeMysqls) List(opts v1.ListOptions) (result *v1alpha1.MysqlList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mysqlsResource, mysqlsKind, c.ns, opts), &v1alpha1.MysqlList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MysqlList{}
	for _, item := range obj.(*v1alpha1.MysqlList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mysqls.
func (c *FakeMysqls) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mysqlsResource, c.ns, opts))

}

// Patch applies the patch and returns the patched mysql.
func (c *FakeMysqls) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Mysql, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mysqlsResource, c.ns, name, data, subresources...), &v1alpha1.Mysql{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Mysql), err
}
