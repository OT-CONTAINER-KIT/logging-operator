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

package configmap

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	configmapData string
)

var log = logf.Log.WithName("configmap")

func generateConfigMap(cr *loggingv1alpha1.Fluentd, labels map[string]string) *corev1.ConfigMap {

	if *cr.Spec.LogPrefix == "namespace" {
		configmapData = ConfigMapContentNamespace
	} else {
		configmapData = ConfigMapContentPod
	}

	config := &corev1.ConfigMap{
		TypeMeta:   identifier.GenerateMetaInformation("ConfigMap", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
		Data: map[string]string{
			"fluent.conf": configmapData,
		},
	}
	identifier.AddOwnerRefToObject(config, identifier.FluentdAsOwner(cr))
	return config
}

func generateExtraConfigMap(cr *loggingv1alpha1.Fluentd, labels map[string]string) *corev1.ConfigMap {
	config := &corev1.ConfigMap{
		TypeMeta:   identifier.GenerateMetaInformation("ConfigMap", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name+"-extra-config", cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
		Data:       *cr.Spec.CustomConfiguration,
	}
	identifier.AddOwnerRefToObject(config, identifier.FluentdAsOwner(cr))
	return config
}

// SyncConfigMap will sync the configmap in Kubernetes
func SyncConfigMap(cr *loggingv1alpha1.Fluentd, config *corev1.ConfigMap, configMapName string) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "configmap",
	)

	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for configmap")
	}

	configMapObject, err := k8sClient.CoreV1().ConfigMaps(cr.Namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating configmap for fluentd", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.CoreV1().ConfigMaps(cr.Namespace).Create(context.TODO(), config, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if configMapObject != config {
		reqLogger.Info("Updating configmap for fluentd", "Name", cr.ObjectMeta.Name)
		k8sClient.CoreV1().ConfigMaps(cr.Namespace).Update(context.TODO(), config, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Fluentd configmap are in sync")
	}
}

// CreateFluentdConfigMap creates the configmap for fluentd configuration
func CreateFluentdConfigMap(cr *loggingv1alpha1.Fluentd) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}

	config := generateConfigMap(cr, labels)
	SyncConfigMap(cr, config, cr.ObjectMeta.Name)
}

// CreateFluentdExtraConfigMap will create the extra configurations for fluentd
func CreateFluentdExtraConfigMap(cr *loggingv1alpha1.Fluentd) {

	configName := cr.ObjectMeta.Name + "-extra-config"
	labels := map[string]string{
		"app":                         configName,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}
	config := generateExtraConfigMap(cr, labels)
	SyncConfigMap(cr, config, configName)
}
