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

// KibanaSpec defines the desired state of Kibana
type KibanaSpec struct {
	// +kubebuilder:default:=1
	Replicas         *int32            `json:"replicas,omitempty"`
	ElasticConfig    ElasticConfig     `json:"esCluster"`
	Security         *Security         `json:"esSecurity,omitempty"`
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`
}

// KibanaStatus defines the observed state of Kibana
type KibanaStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type=string,priority=0,JSONPath=`.spec.esCluster.esVersion`
// +kubebuilder:printcolumn:name="Es Cluster",type=string,priority=0,JSONPath=`.spec.esCluster.clusterName`
// Kibana is the Schema for the kibanas API
type Kibana struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KibanaSpec   `json:"spec,omitempty"`
	Status KibanaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KibanaList contains a list of Kibana
type KibanaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kibana `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kibana{}, &KibanaList{})
}
