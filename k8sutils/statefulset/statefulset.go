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

package statefulset

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/util/intstr"
	"strings"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	graceTime = 10
)

var log = logf.Log.WithName("elastic_statefulset")

// SyncStatefulSet will sync the statefulset in Kubernetes
func SyncStatefulSet(cr *loggingv1alpha1.Elasticsearch, statefulset *appsv1.StatefulSet, nodeType string) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Node.Type", nodeType,
	)

	elasticStatefuleName := cr.ObjectMeta.Name + "-" + nodeType
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for statefulset")
	}

	statefulObject, err := k8sClient.AppsV1().StatefulSets(cr.Namespace).Get(context.TODO(), elasticStatefuleName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating elasticsearch setup", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.AppsV1().StatefulSets(cr.Namespace).Create(context.TODO(), statefulset, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if statefulObject != statefulset {
		reqLogger.Info("Updating elasticsearch setup", "Name", cr.ObjectMeta.Name)
		k8sClient.AppsV1().StatefulSets(cr.Namespace).Update(context.TODO(), statefulset, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Elasticsearch nodes are in sync")
	}
}

// GeneratePVCTemplate is a method to generate the PVC template
func GeneratePVCTemplate(cr *loggingv1alpha1.Elasticsearch, nodeType string, storageSpec *loggingv1alpha1.Storage) corev1.PersistentVolumeClaim {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Node.Type", nodeType,
	)

	pvcTemplate := storageSpec.VolumeClaimTemplate
	elasticStatefuleName := cr.ObjectMeta.Name + "-" + nodeType

	if storageSpec == nil {
		reqLogger.Info("No storage is defined for elasticsearch", "Elasticsearch.Name", cr.ObjectMeta.Name)
	} else {
		pvcTemplate.CreationTimestamp = metav1.Time{}
		pvcTemplate.Name = elasticStatefuleName
		if storageSpec.VolumeClaimTemplate.Spec.AccessModes == nil {
			pvcTemplate.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		} else {
			pvcTemplate.Spec.AccessModes = storageSpec.VolumeClaimTemplate.Spec.AccessModes
		}
		pvcTemplate.Spec.Resources = storageSpec.VolumeClaimTemplate.Spec.Resources
		pvcTemplate.Spec.Selector = storageSpec.VolumeClaimTemplate.Spec.Selector
		pvcTemplate.Spec.Selector = storageSpec.VolumeClaimTemplate.Spec.Selector
	}
	reqLogger.Info("Successfully generated the PVC template")
	return pvcTemplate
}

// SysctlInitContainer will generate the initContainer for system params
func SysctlInitContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {
	var privileged = true
	var runasUser int64 = 0
	return corev1.Container{
		Name:            "sysctl-init",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Command:         []string{"sysctl", "-w", "vm.max_map_count=262144"},
		SecurityContext: &corev1.SecurityContext{
			Privileged: &privileged,
			RunAsUser:  &runasUser,
		},
	}
}

// PluginsInitContainer will generate the initContainer for plugins installation
func PluginsInitContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {
	var privileged = true
	var runasUser int64 = 0
	plugins := []string{"sh", "-c"}

	command := []string{"bin/elasticsearch-plugin install --batch"}
	for _, plugin := range cr.Spec.Plugins {
		command = append(command, *plugin)
	}

	plugins = append(plugins, strings.Join(command, " "))

	return corev1.Container{
		Name:            "plugins-install",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Command:         plugins,
		SecurityContext: &corev1.SecurityContext{
			Privileged: &privileged,
			RunAsUser:  &runasUser,
		},
		VolumeMounts: []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      "plugin-volume",
				MountPath: "/usr/share/elasticsearch/plugins",
			},
		},
	}
}

// ElasticContainer will generate the elastic container interface
func ElasticContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {
	var containerDefinition corev1.Container

	containerDefinition = corev1.Container{
		Name:            "elastic",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Env:             []corev1.EnvVar{},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{}, Requests: corev1.ResourceList{},
		},
		VolumeMounts: []corev1.VolumeMount{},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: graceTime,
			PeriodSeconds:       10,
			FailureThreshold:    3,
			SuccessThreshold:    3,
			TimeoutSeconds:      5,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{
						"sh",
						"-c",
						readinessScript,
					},
				},
			},
		},
	}
	return containerDefinition
}

// StatefulSetObject is for generating the statefulset definition
func StatefulSetObject(cr *loggingv1alpha1.Elasticsearch, nodeType string, labels map[string]string, replicas *int32) *appsv1.StatefulSet {
	var runasUser int64 = 1000
	var fsGroup int64 = 1000
	var serviceLink = true
	statefulset := &appsv1.StatefulSet{
		TypeMeta:   identifier.GenerateMetaInformation("StatefulSet", "apps/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name+"-"+nodeType, cr.Namespace, labels, identifier.GenerateElasticAnnotations()),
		Spec: appsv1.StatefulSetSpec{
			Selector:            identifier.LabelSelectors(labels),
			ServiceName:         cr.ObjectMeta.Name + "-" + nodeType + "-headless",
			Replicas:            replicas,
			PodManagementPolicy: appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type:          appsv1.RollingUpdateStatefulSetStrategyType,
				RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{},
					Containers:     []corev1.Container{},
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:   &fsGroup,
						RunAsUser: &runasUser,
					},
					EnableServiceLinks: &serviceLink,
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "plugin-volume",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
	identifier.AddOwnerRefToObject(statefulset, identifier.ElasticAsOwner(cr))
	return statefulset
}
