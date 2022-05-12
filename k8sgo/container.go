package k8sgo

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type containerParams struct {
	Name        string
	Image       string
	Resources   *corev1.ResourceRequirements
	VolumeMount *corev1.VolumeMount
	EnvVar      []corev1.EnvVar
}

// generateContainerDef is a method to create container definition
func generateContainerDef(params containerParams) []corev1.Container {
	containerDef := []corev1.Container{
		Name:         params.Name,
		Image:        params.Image,
		VolumeMounts: params.VolumeMount,
		Env:          params.EnvVar,
		Resources:    params.Resources,
	}
	return containerDef
}
