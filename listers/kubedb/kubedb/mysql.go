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

// This file was automatically generated by lister-gen

package kubedb

import (
	kubedb "github.com/k8sdb/apimachinery/apis/kubedb"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MySQLLister helps list MySQLs.
type MySQLLister interface {
	// List lists all MySQLs in the indexer.
	List(selector labels.Selector) (ret []*kubedb.MySQL, err error)
	// MySQLs returns an object that can list and get MySQLs.
	MySQLs(namespace string) MySQLNamespaceLister
	MySQLListerExpansion
}

// mySQLLister implements the MySQLLister interface.
type mySQLLister struct {
	indexer cache.Indexer
}

// NewMySQLLister returns a new MySQLLister.
func NewMySQLLister(indexer cache.Indexer) MySQLLister {
	return &mySQLLister{indexer: indexer}
}

// List lists all MySQLs in the indexer.
func (s *mySQLLister) List(selector labels.Selector) (ret []*kubedb.MySQL, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*kubedb.MySQL))
	})
	return ret, err
}

// MySQLs returns an object that can list and get MySQLs.
func (s *mySQLLister) MySQLs(namespace string) MySQLNamespaceLister {
	return mySQLNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MySQLNamespaceLister helps list and get MySQLs.
type MySQLNamespaceLister interface {
	// List lists all MySQLs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*kubedb.MySQL, err error)
	// Get retrieves the MySQL from the indexer for a given namespace and name.
	Get(name string) (*kubedb.MySQL, error)
	MySQLNamespaceListerExpansion
}

// mySQLNamespaceLister implements the MySQLNamespaceLister
// interface.
type mySQLNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MySQLs in the indexer for a given namespace.
func (s mySQLNamespaceLister) List(selector labels.Selector) (ret []*kubedb.MySQL, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*kubedb.MySQL))
	})
	return ret, err
}

// Get retrieves the MySQL from the indexer for a given namespace and name.
func (s mySQLNamespaceLister) Get(name string) (*kubedb.MySQL, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(kubedb.Resource("mysql"), name)
	}
	return obj.(*kubedb.MySQL), nil
}
