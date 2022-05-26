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

package k8sfluentd

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sort"

	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// CreateFluentdDaemonSet is a method to create daemonset for Fluentd
func CreateFluentdDaemonSet(cr *loggingv1beta1.Fluentd) error {
	appName := cr.ObjectMeta.Name
	labels := map[string]string{
		"app": cr.ObjectMeta.Name,
	}
	daemonSetParams := k8sgo.DaemonSetParameters{
		Namespace:     cr.Namespace,
		OwnerDef:      k8sgo.FluentdAsOwner(cr),
		DaemonSetMeta: k8sgo.GenerateObjectMetaInformation(appName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		ContainerParams: k8sgo.ContainerParams{
			Name:        "fluentd",
			Image:       "fluent/fluentd-kubernetes-daemonset:v1-debian-elasticsearch",
			VolumeMount: generateVolumeMounts(cr),
			EnvVar:      generateEnvVariables(cr),
		},
		Labels:      labels,
		Annotations: k8sgo.GenerateAnnotations(),
		Volumes:     getVolumes(cr),
	}
	if cr.Spec.KubernetesConfig != nil {
		daemonSetParams.Affinity = cr.Spec.KubernetesConfig.Affinity
		daemonSetParams.NodeSelector = cr.Spec.KubernetesConfig.NodeSelector
		daemonSetParams.PriorityClassName = cr.Spec.KubernetesConfig.PriorityClassName
		daemonSetParams.Tolerations = cr.Spec.KubernetesConfig.Tolerations
		daemonSetParams.ContainerParams.Resources = cr.Spec.KubernetesConfig.Resources
	} else {
		daemonSetParams.Affinity = &corev1.Affinity{}
		daemonSetParams.NodeSelector = map[string]string{}
		daemonSetParams.PriorityClassName = nil
		daemonSetParams.Tolerations = &[]corev1.Toleration{}
		daemonSetParams.ContainerParams.Resources = &corev1.ResourceRequirements{}
	}
	err := k8sgo.CreateOrUpdateDaemonSet(daemonSetParams)
	if err != nil {
		return err
	}
	return nil
}

// generateEnvVariables is a method to create environment variable for Fluentd
func generateEnvVariables(cr *loggingv1beta1.Fluentd) []corev1.EnvVar {
	fluentdEnvVars := []corev1.EnvVar{
		{Name: "FLUENT_ELASTICSEARCH_HOST", Value: *cr.Spec.ElasticConfig.Host},
		{Name: "FLUENT_ELASTICSEARCH_PORT", Value: "9200"},
		{Name: "FLUENT_CONTAINER_TAIL_PARSER_TYPE", Value: "/^(?<time>.+) (?<stream>stdout|stderr)( (?<logtag>.))? (?<log>.*)$/"},
		{Name: "FLUENTD_SYSTEMD_CONF", Value: "disable"},
	}
	if cr.Spec.Security != nil {
		if *cr.Spec.Security.TLSEnabled {
			fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_USER", Value: "elastic"})
			fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SSL_VERIFY", Value: "false"})
			fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SSL_VERSION", Value: "TLSv1_2"})
			fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SCHEME", Value: "https"})
			fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{
				Name: "FLUENT_ELASTICSEARCH_PASSWORD",
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
	}
	if cr.Spec.IndexNameStrategy != nil {
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_LOGSTASH_PREFIX", Value: fmt.Sprintf("kubernetes-${record['kubernetes']['%s']}", *cr.Spec.IndexNameStrategy)})
	}
	sort.SliceStable(fluentdEnvVars, func(i, j int) bool {
		return fluentdEnvVars[i].Name < fluentdEnvVars[j].Name
	})
	return fluentdEnvVars
}

// getVolumes is a method to define addtional volumes
func getVolumes(cr *loggingv1beta1.Fluentd) *[]corev1.Volume {
	volume := []corev1.Volume{
		{
			Name: "varlogs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log",
				},
			},
		},
		{
			Name: "varlibdockercontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/docker/containers",
				},
			},
		},
	}
	if cr.Spec.CustomConfig == nil {
		volume = append(volume, corev1.Volume{
			Name: "fluentd",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cr.ObjectMeta.Name,
					},
				},
			},
		})
	} else {
		volume = append(volume, corev1.Volume{
			Name: "fluentd",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: *cr.Spec.CustomConfig,
					},
				},
			},
		})
	}
	if cr.Spec.AdditionalConfig != nil {
		volume = append(volume, corev1.Volume{
			Name: "fluentd-additional",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: *cr.Spec.AdditionalConfig,
					},
				},
			},
		})
	}
	return &volume
}

// generateVolumeMounts is a method to create Volume Mounts
func generateVolumeMounts(cr *loggingv1beta1.Fluentd) *[]corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{Name: "varlogs", MountPath: "/var/log"},
		{Name: "varlibdockercontainers", MountPath: "/var/lib/docker/containers", ReadOnly: true},
		{Name: "fluentd", MountPath: "/fluentd/etc/fluent.conf", SubPath: "fluent.conf"},
	}
	if cr.Spec.AdditionalConfig != nil {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "fluentd-additional",
			MountPath: "/fluentd/etc/conf.d/additional-config/",
		})
	}
	return &volumeMounts
}
