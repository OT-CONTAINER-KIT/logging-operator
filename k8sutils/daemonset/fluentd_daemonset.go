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

package daemonset

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("fluentd_daemonset")

func generateFluentdContainer(cr *loggingv1alpha1.Fluentd) *corev1.Container {

	containerDefinition := &corev1.Container{
		Name:            "fluentd",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Env:             []corev1.EnvVar{},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{}, Requests: corev1.ResourceList{},
		},
		VolumeMounts: []corev1.VolumeMount{},
	}

	if cr.Spec.Resources != nil {
		containerDefinition.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Resources.ResourceLimits.CPU)
		containerDefinition.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Resources.ResourceRequests.CPU)
		containerDefinition.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Resources.ResourceLimits.Memory)
		containerDefinition.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Resources.ResourceRequests.Memory)
	}

	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{Name: "varlogs", MountPath: "/var/log"},
		corev1.VolumeMount{Name: "varlibdockercontainers", MountPath: "/var/lib/docker/containers", ReadOnly: true},
		corev1.VolumeMount{Name: "fluentd", MountPath: "/fluentd/etc/fluent.conf", SubPath: "fluent.conf"},
	}

	fluentdEnvVars := []corev1.EnvVar{
		corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_HOST", Value: cr.Spec.FluentdElasticsearch.Host},
		corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_PORT", Value: "9200"},
		corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SED_DISABLE", Value: "true"},
	}

	if cr.Spec.FluentdElasticsearch.TLSEnabled != false {
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_USER", Value: cr.Spec.FluentdElasticsearch.Username})
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_PASSWORD", Value: cr.Spec.FluentdElasticsearch.Password})
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SSL_VERIFY", Value: "false"})
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SSL_VERSION", Value: "TLSv1_2"})
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SCHEME", Value: "https"})
	} else {
		fluentdEnvVars = append(fluentdEnvVars, corev1.EnvVar{Name: "FLUENT_ELASTICSEARCH_SCHEME", Value: "http"})
	}
	containerDefinition.VolumeMounts = volumeMounts
	containerDefinition.Env = fluentdEnvVars

	return containerDefinition
}

func generateDaemonSet(cr *loggingv1alpha1.Fluentd, labels map[string]string) *appsv1.DaemonSet {

	fluentdContainer := generateFluentdContainer(cr)
	daemonsetObject := &appsv1.DaemonSet{
		TypeMeta:   identifier.GenerateMetaInformation("DaemonSet", "apps/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
		Spec: appsv1.DaemonSetSpec{
			Selector: identifier.LabelSelectors(labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.ObjectMeta.Name,
					Containers: []corev1.Container{
						*fluentdContainer,
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "varlogs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
								},
							},
						},
						corev1.Volume{
							Name: "varlibdockercontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/docker/containers",
								},
							},
						},
						corev1.Volume{
							Name: "fluentd",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cr.ObjectMeta.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if cr.Spec.NodeSelector != nil {
		daemonsetObject.Spec.Template.Spec.NodeSelector = *cr.Spec.NodeSelector
	}
	identifier.AddOwnerRefToObject(daemonsetObject, identifier.FluentdAsOwner(cr))

	return daemonsetObject
}

// SyncFluentdDaemonset will sync the daemonset in Kubernetes
func SyncFluentdDaemonset(cr *loggingv1alpha1.Fluentd, daemon *appsv1.DaemonSet) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "Daemonset",
	)

	daemonsetName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for daemonset")
	}

	daemonsetObject, err := k8sClient.AppsV1().DaemonSets(cr.Namespace).Get(context.TODO(), daemonsetName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating fluentd setup", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.AppsV1().DaemonSets(cr.Namespace).Create(context.TODO(), daemon, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if daemonsetObject != daemon {
		reqLogger.Info("Updating fluentd setup", "Name", cr.ObjectMeta.Name)
		k8sClient.AppsV1().DaemonSets(cr.Namespace).Update(context.TODO(), daemon, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Fluentd daemonsets are in sync")
	}
}

// CreateFluentdDaemonset creates the daemonset for fluentd
func CreateFluentdDaemonset(cr *loggingv1alpha1.Fluentd) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}

	config := generateDaemonSet(cr, labels)
	SyncFluentdDaemonset(cr, config)
}
