/*
Copyright 2020 Opstree Solutions.

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

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KibanaSpec defines the desired state of Kibana
type KibanaSpec struct {
	Replicas            *int32              `json:"replicas"`
	Image               string              `json:"image"`
	ImagePullPolicy     corev1.PullPolicy   `json:"imagePullPolicy,omitempty"`
	KibanaElasticsearch KibanaElasticsearch `json:"elasticsearch,omitempty"`
	ElasticSecretName   *string             `json:"elasticSecretName,omitempty"`
	Resources           *Resources          `json:"resources,omitempty"`
	Affinity            *corev1.Affinity    `json:"affinity,omitempty"`
}

// KibanaElasticsearch is the struct for elasticsearch configuration for fluentd
type KibanaElasticsearch struct {
	Host       string `json:"host,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	TLSEnabled bool   `json:"tlsEnabled,omitempty"`
}

// KibanaStatus defines the observed state of Kibana
type KibanaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Elasticsearch string `json:"elasticsearch,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Kibana is the Schema for the kibanas API
type Kibana struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KibanaSpec   `json:"spec,omitempty"`
	Status KibanaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KibanaList contains a list of Kibana
type KibanaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kibana `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kibana{}, &KibanaList{})
}
