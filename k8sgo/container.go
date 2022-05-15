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

type containerParams struct {
	Name           string
	Image          string
	Resources      *corev1.ResourceRequirements
	VolumeMount    *[]corev1.VolumeMount
	EnvVar         []corev1.EnvVar
	ReadinessProbe *corev1.Probe
	LivenessProbe  *corev1.Probe
}

// generateContainerDef is a method to create container definition
func generateContainerDef(params containerParams) []corev1.Container {
	containerDef := []corev1.Container{
		{
			Name:           params.Name,
			Image:          params.Image,
			VolumeMounts:   *params.VolumeMount,
			Env:            params.EnvVar,
			Resources:      *params.Resources,
			LivenessProbe:  params.LivenessProbe,
			ReadinessProbe: params.ReadinessProbe,
		},
	}
	return containerDef
}
