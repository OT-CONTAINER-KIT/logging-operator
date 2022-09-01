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

package k8sgo

import (
	corev1 "k8s.io/api/core/v1"
)

type ContainerParams struct {
	Name           string
	Image          string
	Resources      *corev1.ResourceRequirements
	InitResources  *corev1.ResourceRequirements
	VolumeMount    *[]corev1.VolumeMount
	EnvVar         []corev1.EnvVar
	EnvVarFrom     []corev1.EnvFromSource
	ReadinessProbe *corev1.Probe
	LivenessProbe  *corev1.Probe
	Lifecycle      *corev1.Lifecycle
}

// generateContainerDef is a method to create container definition
func generateContainerDef(params ContainerParams) []corev1.Container {
	containerDef := []corev1.Container{
		{
			Name:           params.Name,
			Image:          params.Image,
			VolumeMounts:   *params.VolumeMount,
			Env:            params.EnvVar,
			LivenessProbe:  params.LivenessProbe,
			ReadinessProbe: params.ReadinessProbe,
		},
	}
	if params.Resources != nil {
		containerDef[0].Resources = *params.Resources
	}
	if params.EnvVarFrom != nil {
		containerDef[0].EnvFrom = params.EnvVarFrom
	}
	if params.Lifecycle != nil {
		containerDef[0].Lifecycle = params.Lifecycle
	}
	return containerDef
}
