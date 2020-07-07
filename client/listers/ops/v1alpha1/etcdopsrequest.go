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

package v1alpha1

import (
	v1alpha1 "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EtcdOpsRequestLister helps list EtcdOpsRequests.
type EtcdOpsRequestLister interface {
	// List lists all EtcdOpsRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.EtcdOpsRequest, err error)
	// EtcdOpsRequests returns an object that can list and get EtcdOpsRequests.
	EtcdOpsRequests(namespace string) EtcdOpsRequestNamespaceLister
	EtcdOpsRequestListerExpansion
}

// etcdOpsRequestLister implements the EtcdOpsRequestLister interface.
type etcdOpsRequestLister struct {
	indexer cache.Indexer
}

// NewEtcdOpsRequestLister returns a new EtcdOpsRequestLister.
func NewEtcdOpsRequestLister(indexer cache.Indexer) EtcdOpsRequestLister {
	return &etcdOpsRequestLister{indexer: indexer}
}

// List lists all EtcdOpsRequests in the indexer.
func (s *etcdOpsRequestLister) List(selector labels.Selector) (ret []*v1alpha1.EtcdOpsRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EtcdOpsRequest))
	})
	return ret, err
}

// EtcdOpsRequests returns an object that can list and get EtcdOpsRequests.
func (s *etcdOpsRequestLister) EtcdOpsRequests(namespace string) EtcdOpsRequestNamespaceLister {
	return etcdOpsRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EtcdOpsRequestNamespaceLister helps list and get EtcdOpsRequests.
type EtcdOpsRequestNamespaceLister interface {
	// List lists all EtcdOpsRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.EtcdOpsRequest, err error)
	// Get retrieves the EtcdOpsRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.EtcdOpsRequest, error)
	EtcdOpsRequestNamespaceListerExpansion
}

// etcdOpsRequestNamespaceLister implements the EtcdOpsRequestNamespaceLister
// interface.
type etcdOpsRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EtcdOpsRequests in the indexer for a given namespace.
func (s etcdOpsRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.EtcdOpsRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EtcdOpsRequest))
	})
	return ret, err
}

// Get retrieves the EtcdOpsRequest from the indexer for a given namespace and name.
func (s etcdOpsRequestNamespaceLister) Get(name string) (*v1alpha1.EtcdOpsRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("etcdopsrequest"), name)
	}
	return obj.(*v1alpha1.EtcdOpsRequest), nil
}
