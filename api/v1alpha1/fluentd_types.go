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

// FluentdSpec defines the desired state of Fluentd
type FluentdSpec struct {
	FluentdElasticsearch FluentdElasticsearch `json:"elasticsearch,omitempty"`
	Image                string               `json:"image"`
	ImagePullPolicy      corev1.PullPolicy    `json:"imagePullPolicy,omitempty"`
	NodeSelector         *map[string]string   `json:"nodeSelector,omitempty"`
	LogPrefix            *string              `json:"logPrefix,omitempty"`
	CustomConfiguration  *string              `json:"customConfiguration,omitempty"`
	Resources            *Resources           `json:"resources,omitempty"`
}

// FluentdElasticsearch is the struct for elasticsearch configuration for fluentd
type FluentdElasticsearch struct {
	Host       string `json:"host,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	TLSEnabled bool   `json:"tlsEnabled,omitempty"`
}

// FluentdStatus defines the observed state of Fluentd
type FluentdStatus struct {
	Elasticsearch string `json:"elasticsearch,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Fluentd is the Schema for the fluentds API
type Fluentd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluentdSpec   `json:"spec,omitempty"`
	Status FluentdStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FluentdList contains a list of Fluentd
type FluentdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fluentd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fluentd{}, &FluentdList{})
}
