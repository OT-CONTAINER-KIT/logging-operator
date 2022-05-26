/*
Copyright 2022 Opstree Solutions.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// ElasticsearchSpec defines the desired state of Elasticsearch
type ElasticsearchSpec struct {
	ClusterName string    `json:"esClusterName"`
	ESVersion   string    `json:"esVersion"`
	Security    *Security `json:"esSecurity,omitempty"`
	// +kubebuilder:validation:default:={esMaster:{replicas: 3}}
	// +kubebuilder:default:={storage:{accessModes: {ReadWriteOnce}, storageSize: "1Gi"},jvmMaxMemory: "1g", jvmMinMemory: "1g", replicas: 3}
	ESMaster    *NodeSpecificConfig `json:"esMaster,omitempty"`
	ESData      *NodeSpecificConfig `json:"esData,omitempty"`
	ESIngestion *NodeSpecificConfig `json:"esIngestion,omitempty"`
	ESClient    *NodeSpecificConfig `json:"esClient,omitempty"`
}

// NodeSpecificConfig defines the properties for elasticsearch nodes
type NodeSpecificConfig struct {
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`
	Replicas         *int32            `json:"replicas,omitempty"`
	CustomConfig     *string           `json:"customConfig,omitempty"`
	Storage          *Storage          `json:"storage,omitempty"`
	JvmMaxMemory     *string           `json:"jvmMaxMemory,omitempty"`
	JvmMinMemory     *string           `json:"jvmMinMemory,omitempty"`
}

// Security defines the security config of Elasticsearch
type Security struct {
	ExistingSecret       *string `json:"existingSecret,omitempty"`
	TLSEnabled           *bool   `json:"tlsEnabled,omitempty"`
	AutoGeneratePassword *bool   `json:"autoGeneratePassword,omitempty"`
}

//+kubebuilder:subresource:status
// ElasticsearchStatus defines the observed state of Elasticsearch
type ElasticsearchStatus struct {
	ESVersion    string `json:"esVersion,omitempty"`
	ClusterState string `json:"esClusterState,omitempty"`
	ActiveShards *int32 `json:"activeShards,omitempty"`
	Indices      *int32 `json:"indices,omitempty"`
	ESMaster     *int32 `json:"esMaster,omitempty"`
	ESData       *int32 `json:"esData,omitempty"`
	ESClient     *int32 `json:"esClient,omitempty"`
	ESIngestion  *int32 `json:"esIngestion,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type=string,priority=0,JSONPath=`.status.esVersion`
// +kubebuilder:printcolumn:name="State",type=string,priority=0,JSONPath=`.status.esClusterState`
// +kubebuilder:printcolumn:name="Shards",type=integer,priority=0,JSONPath=`.status.activeShards`
// +kubebuilder:printcolumn:name="Indices",type=integer,priority=0,JSONPath=`.status.indices`
// +kubebuilder:printcolumn:name="Master",type=integer,priority=1,JSONPath=`.status.esMaster`
// +kubebuilder:printcolumn:name="Data",type=integer,priority=1,JSONPath=`.status.esClient`
// +kubebuilder:printcolumn:name="Client",type=integer,priority=1,JSONPath=`.status.esMaster`
// +kubebuilder:printcolumn:name="Ingestion",type=integer,priority=1,JSONPath=`.status.esIngestion`
// Elasticsearch is the Schema for the elasticsearches API
type Elasticsearch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchSpec   `json:"spec,omitempty"`
	Status ElasticsearchStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ElasticsearchList contains a list of Elasticsearch
type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Elasticsearch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Elasticsearch{}, &ElasticsearchList{})
}
