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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IndexLifeCycleSpec defines the desired state of IndexLifeCycle
type IndexLifeCycleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of IndexLifeCycle. Edit indexlifecycle_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// IndexLifeCycleStatus defines the observed state of IndexLifeCycle
type IndexLifeCycleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IndexLifeCycle is the Schema for the indexlifecycles API
type IndexLifeCycle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IndexLifeCycleSpec   `json:"spec,omitempty"`
	Status IndexLifeCycleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IndexLifeCycleList contains a list of IndexLifeCycle
type IndexLifeCycleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IndexLifeCycle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IndexLifeCycle{}, &IndexLifeCycleList{})
}
