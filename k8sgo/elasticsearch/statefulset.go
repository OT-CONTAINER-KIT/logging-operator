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

package k8selastic

import (
	"fmt"
	"logging-operator/k8sgo"

	corev1 "k8s.io/api/core/v1"
	loggingv1beta1 "logging-operator/api/v1beta1"
)

// CreateElasticsearchStatefulSet is a method to create elasticsearch statefulset
func CreateElasticsearchStatefulSet(cr *loggingv1beta1.Elasticsearch, nodeConfig *loggingv1beta1.NodeSpecificConfig, role string, envVars []corev1.EnvVar) error {
	appName := fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, role)
	labels := map[string]string{
		"app":  appName,
		"role": role,
	}
	statefulsetParams := k8sgo.StatefulSetParameters{
		OwnerDef:        k8sgo.ElasticAsOwner(cr),
		StatefulSetMeta: k8sgo.GenerateObjectMetaInformation(appName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		Namespace:       cr.Namespace,
		ContainerParams: k8sgo.ContainerParams{
			Name:           "elastic",
			Image:          fmt.Sprintf("docker.elastic.co/elasticsearch/elasticsearch:%s", cr.Spec.ESVersion),
			VolumeMount:    getVolumeMounts(cr, role),
			EnvVar:         envVars,
			ReadinessProbe: createProbeInfo(),
		},
		Labels:      labels,
		Annotations: k8sgo.GenerateAnnotations(),
		Replicas:    nodeConfig.Replicas,
		PVCParameters: k8sgo.PVCParameters{
			Name:             appName,
			Namespace:        cr.Namespace,
			Labels:           labels,
			Annotations:      k8sgo.GenerateAnnotations(),
			AccessModes:      nodeConfig.Storage.AccessModes,
			StorageSize:      nodeConfig.Storage.StorageSize,
			StorageClassName: nodeConfig.Storage.StorageClassName,
		},
	}
	statefulsetParams.ExtraVolumes = getVolumes(cr)

	if nodeConfig != nil {
		if nodeConfig.CustomConfig != nil {
			statefulsetParams.ContainerParams.EnvVarFrom = []corev1.EnvFromSource{
				{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: *nodeConfig.CustomConfig,
						},
					},
				},
			}
		}
		if nodeConfig.KubernetesConfig != nil {
			statefulsetParams.Affinity = nodeConfig.KubernetesConfig.Affinity
			statefulsetParams.NodeSelector = nodeConfig.KubernetesConfig.NodeSelector
			statefulsetParams.PriorityClassName = nodeConfig.KubernetesConfig.PriorityClassName
			statefulsetParams.Tolerations = nodeConfig.KubernetesConfig.Tolerations
			statefulsetParams.ContainerParams.Resources = nodeConfig.KubernetesConfig.Resources
		} else {
			statefulsetParams.Affinity = &corev1.Affinity{}
			statefulsetParams.NodeSelector = map[string]string{}
			statefulsetParams.PriorityClassName = nil
			statefulsetParams.Tolerations = &[]corev1.Toleration{}
			statefulsetParams.ContainerParams.Resources = &corev1.ResourceRequirements{}
			statefulsetParams.ContainerParams.InitResources = &corev1.ResourceRequirements{}
		}
	}
	err := k8sgo.CreateOrUpdateStateFul(statefulsetParams)
	if err != nil {
		return err
	}
	return nil
}

// getVolumeMounts is a method to get volume mounts for statefulset
func getVolumeMounts(cr *loggingv1beta1.Elasticsearch, role string) *[]corev1.VolumeMount {
	appName := fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, role)
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      appName,
			MountPath: "/usr/share/elasticsearch/data",
		},
	}
	if cr.Spec.Security != nil {
		if cr.Spec.Security.TLSEnabled != nil && *cr.Spec.Security.TLSEnabled {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      "tls-cert",
				MountPath: "/usr/share/elasticsearch/config/certs",
			})
		}
	}
	return &volumeMounts
}

// getVolumes is a method to define addtional volumes
func getVolumes(cr *loggingv1beta1.Elasticsearch) *[]corev1.Volume {
	var volume []corev1.Volume
	if cr.Spec.Security != nil {
		if cr.Spec.Security.TLSEnabled != nil && *cr.Spec.Security.TLSEnabled {
			volume = append(volume, corev1.Volume{
				Name: "tls-cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: fmt.Sprintf("%s-tls-cert", cr.ObjectMeta.Name),
					},
				},
			})
		}
	}
	return &volume
}

// generateEnvVariables is a method to create environment variables
func generateEnvVariables(cr *loggingv1beta1.Elasticsearch, nodeConfig loggingv1beta1.NodeSpecificConfig) []corev1.EnvVar {
	var javaOpts string
	envVars := []corev1.EnvVar{{Name: "ELASTIC_USERNAME", Value: "elastic"}}
	if cr.Spec.Security != nil {
		if cr.Spec.Security.AutoGeneratePassword != nil && *cr.Spec.Security.AutoGeneratePassword {
			envVars = append(envVars, corev1.EnvVar{
				Name: "ELASTIC_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, "password"),
						},
						Key: "password",
					},
				},
			})
		}
		if cr.Spec.Security.ExistingSecret != nil {
			envVars = append(envVars, corev1.EnvVar{
				Name: "ELASTIC_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: *cr.Spec.Security.ExistingSecret,
						},
						Key: "password",
					},
				},
			})
		}
		if cr.Spec.Security.TLSEnabled != nil && *cr.Spec.Security.TLSEnabled {
			envVars = append(envVars, corev1.EnvVar{Name: "SCHEME", Value: "https"})
		} else {
			envVars = append(envVars, corev1.EnvVar{Name: "SCHEME", Value: "http"})
		}
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "SCHEME", Value: "http"})
	}
	if nodeConfig.JvmMaxMemory != nil && nodeConfig.JvmMinMemory != nil {
		javaOpts = fmt.Sprintf("-Xmx%s -Xms%s", *nodeConfig.JvmMaxMemory, *nodeConfig.JvmMinMemory)
	} else {
		javaOpts = fmt.Sprintf("-Xmx1g -Xms1g")
	}
	envVars = append(envVars, corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: javaOpts})
	return envVars
}

// createProbeInfo is a method to create probe for elasticsearch
func createProbeInfo() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 15,
		PeriodSeconds:       15,
		FailureThreshold:    5,
		TimeoutSeconds:      5,
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"bash", "-c", healthCheckScript},
			},
		},
	}
}
