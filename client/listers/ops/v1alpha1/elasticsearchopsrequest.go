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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "kubedb.dev/apimachinery/apis/ops/v1alpha1"
)

// ElasticsearchOpsRequestLister helps list ElasticsearchOpsRequests.
type ElasticsearchOpsRequestLister interface {
	// List lists all ElasticsearchOpsRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchOpsRequest, err error)
	// ElasticsearchOpsRequests returns an object that can list and get ElasticsearchOpsRequests.
	ElasticsearchOpsRequests(namespace string) ElasticsearchOpsRequestNamespaceLister
	ElasticsearchOpsRequestListerExpansion
}

// elasticsearchOpsRequestLister implements the ElasticsearchOpsRequestLister interface.
type elasticsearchOpsRequestLister struct {
	indexer cache.Indexer
}

// NewElasticsearchOpsRequestLister returns a new ElasticsearchOpsRequestLister.
func NewElasticsearchOpsRequestLister(indexer cache.Indexer) ElasticsearchOpsRequestLister {
	return &elasticsearchOpsRequestLister{indexer: indexer}
}

// List lists all ElasticsearchOpsRequests in the indexer.
func (s *elasticsearchOpsRequestLister) List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchOpsRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ElasticsearchOpsRequest))
	})
	return ret, err
}

// ElasticsearchOpsRequests returns an object that can list and get ElasticsearchOpsRequests.
func (s *elasticsearchOpsRequestLister) ElasticsearchOpsRequests(namespace string) ElasticsearchOpsRequestNamespaceLister {
	return elasticsearchOpsRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ElasticsearchOpsRequestNamespaceLister helps list and get ElasticsearchOpsRequests.
type ElasticsearchOpsRequestNamespaceLister interface {
	// List lists all ElasticsearchOpsRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchOpsRequest, err error)
	// Get retrieves the ElasticsearchOpsRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ElasticsearchOpsRequest, error)
	ElasticsearchOpsRequestNamespaceListerExpansion
}

// elasticsearchOpsRequestNamespaceLister implements the ElasticsearchOpsRequestNamespaceLister
// interface.
type elasticsearchOpsRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ElasticsearchOpsRequests in the indexer for a given namespace.
func (s elasticsearchOpsRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ElasticsearchOpsRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ElasticsearchOpsRequest))
	})
	return ret, err
}

// Get retrieves the ElasticsearchOpsRequest from the indexer for a given namespace and name.
func (s elasticsearchOpsRequestNamespaceLister) Get(name string) (*v1alpha1.ElasticsearchOpsRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("elasticsearchopsrequest"), name)
	}
	return obj.(*v1alpha1.ElasticsearchOpsRequest), nil
}
