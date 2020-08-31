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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IndexLifecycleSpec defines the desired state of IndexLifecycle
type IndexLifecycleSpec struct {
	Rollover Rollover `json:"rollover,omitempty"`
	Delete   Delete   `json:"delete,omitempty"`
}

// Rollover is the struct for index roll overing
type Rollover struct {
	MaxSize string `json:"maxSize,omitempty"`
	MaxAge  string `json:"maxAge,omitempty"`
}

// Delete is the struct for index deletion
type Delete struct {
	MinAge string `json:"minAge,omitempty"`
}

// IndexLifecycleStatus defines the observed state of IndexLifecycle
type IndexLifecycleStatus struct {
	Rollover Rollover `json:"rollover,omitempty"`
	Delete   Delete   `json:"delete,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IndexLifecycle is the Schema for the indexlifecycles API
type IndexLifecycle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IndexLifecycleSpec   `json:"spec,omitempty"`
	Status IndexLifecycleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IndexLifecycleList contains a list of IndexLifecycle
type IndexLifecycleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IndexLifecycle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IndexLifecycle{}, &IndexLifecycleList{})
}
