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

// ElasticsearchSpec defines the desired state of Elasticsearch
type ElasticsearchSpec struct {
	ClusterName string    `json:"esClusterName"`
	ESVersion   string    `json:"esVersion"`
	Security    *Security `json:"esSecurity,omitempty"`
	// +kubebuilder:default:={esMaster:{replicas: 3}}
	ESMaster    *NodeSpecificConfig `json:"esMaster,omitempty"`
	ESData      *NodeSpecificConfig `json:"esData,omitempty"`
	ESIngestion *NodeSpecificConfig `json:"esIngestion,omitempty"`
	ESClient    *NodeSpecificConfig `json:"esClient,omitempty"`
}

// NodeSpecificConfig defines the properties for elasticsearch nodes
type NodeSpecificConfig struct {
	KubernetesConfig   *KubernetesConfig  `json:"kubernetesConfig,omitempty"`
	Replicas           *int32             `json:"replicas,omitempty"`
	CustomEnvVariables *map[string]string `json:"customEnvVariables,omitempty"`
	// +kubebuilder:default:={storage:{accessModes: [ReadWriteOnce], storageSize: "1Gi"}}
	Storage *Storage `json:"storage,omitempty"`
	// +kubebuilder:default:="1g"
	JvmMaxMemory *string `json:"jvmMaxMemory,omitempty"`
	// +kubebuilder:default:="1g"
	JvmMinMemory *string `json:"jvmMinMemory,omitempty"`
}

// Security defines the security config of Elasticsearch
type Security struct {
	ExistingSecret       *string `json:"existingSecret,omitempty"`
	TLSEnabled           *bool   `json:"tlsEnabled,omitempty"`
	AutoGeneratePassword *bool   `json:"autoGeneratePassword,omitempty"`
}

// ElasticsearchStatus defines the observed state of Elasticsearch
type ElasticsearchStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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
