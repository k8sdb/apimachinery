/*
Copyright AppsCode Inc. and Contributors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
)

// EtcdLister helps list Etcds.
type EtcdLister interface {
	// List lists all Etcds in the indexer.
	List(selector labels.Selector) (ret []*v1alpha2.Etcd, err error)
	// Etcds returns an object that can list and get Etcds.
	Etcds(namespace string) EtcdNamespaceLister
	EtcdListerExpansion
}

// etcdLister implements the EtcdLister interface.
type etcdLister struct {
	indexer cache.Indexer
}

// NewEtcdLister returns a new EtcdLister.
func NewEtcdLister(indexer cache.Indexer) EtcdLister {
	return &etcdLister{indexer: indexer}
}

// List lists all Etcds in the indexer.
func (s *etcdLister) List(selector labels.Selector) (ret []*v1alpha2.Etcd, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.Etcd))
	})
	return ret, err
}

// Etcds returns an object that can list and get Etcds.
func (s *etcdLister) Etcds(namespace string) EtcdNamespaceLister {
	return etcdNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EtcdNamespaceLister helps list and get Etcds.
type EtcdNamespaceLister interface {
	// List lists all Etcds in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha2.Etcd, err error)
	// Get retrieves the Etcd from the indexer for a given namespace and name.
	Get(name string) (*v1alpha2.Etcd, error)
	EtcdNamespaceListerExpansion
}

// etcdNamespaceLister implements the EtcdNamespaceLister
// interface.
type etcdNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Etcds in the indexer for a given namespace.
func (s etcdNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.Etcd, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.Etcd))
	})
	return ret, err
}

// Get retrieves the Etcd from the indexer for a given namespace and name.
func (s etcdNamespaceLister) Get(name string) (*v1alpha2.Etcd, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("etcd"), name)
	}
	return obj.(*v1alpha2.Etcd), nil
}
