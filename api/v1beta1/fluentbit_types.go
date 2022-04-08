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

// FluentBitSpec defines the desired state of FluentBit
type FluentBitSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of FluentBit. Edit fluentbit_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// FluentBitStatus defines the observed state of FluentBit
type FluentBitStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FluentBit is the Schema for the fluentbits API
type FluentBit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluentBitSpec   `json:"spec,omitempty"`
	Status FluentBitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FluentBitList contains a list of FluentBit
type FluentBitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FluentBit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluentBit{}, &FluentBitList{})
}
