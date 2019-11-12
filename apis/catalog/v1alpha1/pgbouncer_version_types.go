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
package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodePgBouncerVersion     = "pbversion"
	ResourceKindPgBouncerVersion     = "PgBouncerVersion"
	ResourceSingularPgBouncerVersion = "pgbouncerversion"
	ResourcePluralPgBouncerVersion   = "pgbouncerversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerVersion defines a PgBouncer database version.
type PgBouncerVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PgBouncerVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// PgBouncerVersionSpec is the spec for pgbouncer version
type PgBouncerVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Database Image
	Server PgBouncerServerVersion `json:"server" protobuf:"bytes,2,opt,name=server"`
	// Exporter Image
	Exporter PgBouncerVersionExporter `json:"exporter" protobuf:"bytes,3,opt,name=exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,4,opt,name=deprecated"`
}

// PgBouncerServerVersion is the PgBouncer Database image
type PgBouncerServerVersion struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// PostgresVersionExporter is the image for the Postgres exporter
type PgBouncerVersionExporter struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerVersionList is a list of PgBouncerVersions
type PgBouncerVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of PgBouncerVersion CRD objects
	Items []PgBouncerVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
