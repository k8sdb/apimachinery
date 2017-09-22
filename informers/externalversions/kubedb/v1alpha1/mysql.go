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

// This file was automatically generated by informer-gen

package v1alpha1

import (
	kubedb_v1alpha1 "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	client "github.com/k8sdb/apimachinery/client"
	internalinterfaces "github.com/k8sdb/apimachinery/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/k8sdb/apimachinery/listers/kubedb/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// MySQLInformer provides access to a shared informer and lister for
// MySQLs.
type MySQLInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.MySQLLister
}

type mySQLInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

func newMySQLInformer(client client.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	sharedIndexInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.KubedbV1alpha1().MySQLs(v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.KubedbV1alpha1().MySQLs(v1.NamespaceAll).Watch(options)
			},
		},
		&kubedb_v1alpha1.MySQL{},
		resyncPeriod,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	return sharedIndexInformer
}

func (f *mySQLInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kubedb_v1alpha1.MySQL{}, newMySQLInformer)
}

func (f *mySQLInformer) Lister() v1alpha1.MySQLLister {
	return v1alpha1.NewMySQLLister(f.Informer().GetIndexer())
}