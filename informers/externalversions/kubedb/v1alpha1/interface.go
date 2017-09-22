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
	internalinterfaces "github.com/k8sdb/apimachinery/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// DormantDatabases returns a DormantDatabaseInformer.
	DormantDatabases() DormantDatabaseInformer
	// Elasticsearchs returns a ElasticsearchInformer.
	Elasticsearchs() ElasticsearchInformer
	// MySQLs returns a MySQLInformer.
	MySQLs() MySQLInformer
	// Postgreses returns a PostgresInformer.
	Postgreses() PostgresInformer
	// Snapshots returns a SnapshotInformer.
	Snapshots() SnapshotInformer
}

type version struct {
	internalinterfaces.SharedInformerFactory
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory) Interface {
	return &version{f}
}

// DormantDatabases returns a DormantDatabaseInformer.
func (v *version) DormantDatabases() DormantDatabaseInformer {
	return &dormantDatabaseInformer{factory: v.SharedInformerFactory}
}

// Elasticsearchs returns a ElasticsearchInformer.
func (v *version) Elasticsearchs() ElasticsearchInformer {
	return &elasticsearchInformer{factory: v.SharedInformerFactory}
}

// MySQLs returns a MySQLInformer.
func (v *version) MySQLs() MySQLInformer {
	return &mySQLInformer{factory: v.SharedInformerFactory}
}

// Postgreses returns a PostgresInformer.
func (v *version) Postgreses() PostgresInformer {
	return &postgresInformer{factory: v.SharedInformerFactory}
}

// Snapshots returns a SnapshotInformer.
func (v *version) Snapshots() SnapshotInformer {
	return &snapshotInformer{factory: v.SharedInformerFactory}
}
