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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IndexTemplateSpec defines the desired state of IndexTemplate
type IndexTemplateSpec struct {
	Rollover Rollover `json:"rollover,omitempty"`
	Delete   Delete   `json:"delete,omitempty"`
}

// IndexTemplateStatus defines the observed state of IndexTemplate
type IndexTemplateStatus struct {
	Rollover Rollover `json:"rollover,omitempty"`
	Delete   Delete   `json:"delete,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IndexTemplate is the Schema for the indextemplates API
type IndexTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IndexTemplateSpec   `json:"spec,omitempty"`
	Status IndexTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IndexTemplateList contains a list of IndexTemplate
type IndexTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IndexTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IndexTemplate{}, &IndexTemplateList{})
}
