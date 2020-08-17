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
)

var (
	kibanaConfigMapData string
)

func generateKibanaConfigMap(cr *loggingv1alpha1.Kibana, labels map[string]string) *corev1.ConfigMap {

	if cr.Spec.KibanaElasticsearch.TLSEnabled != false {
		kibanaConfigMapData = kiabanConfigTLSData
	} else {
		kibanaConfigMapData = kiabanConfigData
	}
	config := &corev1.ConfigMap{
		TypeMeta:   identifier.GenerateMetaInformation("ConfigMap", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateKibanaAnnotations()),
		Data: map[string]string{
			"kibana.yml": kibanaConfigMapData,
		},
	}
	identifier.AddOwnerRefToObject(config, identifier.KibanaAsOwner(cr))
	return config
}

// SyncKibanaConfigMap will sync the configmap in Kubernetes
func SyncKibanaConfigMap(cr *loggingv1alpha1.Kibana, config *corev1.ConfigMap) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "configmap",
	)

	configMapName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for configmap")
	}

	configMapObject, err := k8sClient.CoreV1().ConfigMaps(cr.Namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating configmap for kibana", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.CoreV1().ConfigMaps(cr.Namespace).Create(context.TODO(), config, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if configMapObject != config {
		reqLogger.Info("Updating configmap for kibana", "Name", cr.ObjectMeta.Name)
		k8sClient.CoreV1().ConfigMaps(cr.Namespace).Update(context.TODO(), config, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Kibana configmap are in sync")
	}
}

// CreateKibanaConfigMap creates the configmap for kibana configuration
func CreateKibanaConfigMap(cr *loggingv1alpha1.Kibana) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Kibana",
	}

	config := generateKibanaConfigMap(cr, labels)
	SyncKibanaConfigMap(cr, config)
}
