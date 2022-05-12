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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesConfig will define the Kubernetes specific properties
type KubernetesConfig struct {
	Resources         *corev1.ResourceRequirements `json:"resources,omitempty"`
	Affinity          *corev1.Affinity             `json:"mongoAffinity,omitempty"`
	Tolerations       *[]corev1.Toleration         `json:"tolerations,omitempty"`
	PriorityClassName *string                      `json:"priorityClassName,omitempty"`
	SecurityContext   *corev1.PodSecurityContext   `json:"securityContext,omitempty"`
}

// Storage is the inteface to add pvc and pv support in MongoDB
type Storage struct {
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,1,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
	StorageClassName *string                             `json:"storageClass,omitempty" protobuf:"bytes,5,opt,name=storageClassName"`
	StorageSize      string                              `json:"storageSize,omitempty" protobuf:"bytes,5,opt,name=storageClassName"`
}
