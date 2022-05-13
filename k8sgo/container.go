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
