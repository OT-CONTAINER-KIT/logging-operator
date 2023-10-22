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

// FluentdSpec defines the desired state of Fluentd
type FluentdSpec struct {
	ElasticConfig    ElasticConfig     `json:"esCluster"`
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`
	Security         *Security         `json:"esSecurity,omitempty"`
	// +kubebuilder:default:=namespace_name
	// +kubebuilder:validation:Pattern=`namespace_name$|pod_name$`
	IndexNameStrategy *string `json:"indexNameStrategy,omitempty"`
	CustomConfig      *string `json:"customConfig,omitempty"`
	AdditionalConfig  *string `json:"additionalConfig,omitempty"`
}

// ElasticConfig is a method for elasticsearch configuration
type ElasticConfig struct {
	Host        *string `json:"host"`
	ClusterName string  `json:"clusterName,omitempty"`
	ESVersion   string  `json:"esVersion,omitempty"`
}

// FluentdStatus defines the observed state of Fluentd
type FluentdStatus struct {
	TotalAgents *int32 `json:"totalAgents,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Elasticsearch Host",type=string,priority=0,JSONPath=`.spec.esCluster.host`
// +kubebuilder:printcolumn:name="Total Agents",type=string,priority=0,JSONPath=`.status.totalAgents`
// Fluentd is the Schema for the fluentds API
type Fluentd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluentdSpec   `json:"spec,omitempty"`
	Status FluentdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FluentdList contains a list of Fluentd
type FluentdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fluentd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fluentd{}, &FluentdList{})
}
