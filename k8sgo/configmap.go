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
	"context"
	"github.com/banzaicloud/k8s-objectmatcher/patch"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapParameters is an interface for configmap input
type ConfigMapParameters struct {
	Name           string
	OwnerDef       metav1.OwnerReference
	Namespace      string
	ConfigMapMeta  metav1.ObjectMeta
	ConfigMapKey   string
	ConfigMapValue string
}

// CreateOrUpdateConfigMap method will create or update configMap
func CreateOrUpdateConfigMap(params ConfigMapParameters) error {
	logger := LogGenerator(params.ConfigMapMeta.Name, params.Namespace, "ConfigMap")
	configMapDef := generateConfigMap(params)
	storedConfigMap, err := getConfigMap(params.ConfigMapMeta.Name, params.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(configMapDef); err != nil {
				logger.Error(err, "Unable to patch configmap with compare annotations")
			}
			return createConfigMap(params.Namespace, configMapDef)
		}
		return err
	}
	return patchConfigMap(storedConfigMap, configMapDef, params.Namespace)
}

// generateConfigMap is a method that will generate a configmap interface
func generateConfigMap(params ConfigMapParameters) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		TypeMeta:   GenerateMetaInformation("ConfigMap", "v1"),
		ObjectMeta: params.ConfigMapMeta,
		Data: map[string]string{
			params.ConfigMapKey: params.ConfigMapValue,
		},
	}
	AddOwnerRefToObject(configMap, params.OwnerDef)
	return configMap
}

// updateConfigMap is a method to update Kubernetes configMap
func updateConfigMap(namespace string, configMap *corev1.ConfigMap) error {
	logger := LogGenerator(configMap.Name, namespace, "ConfigMap")
	_, err := GenerateK8sClient().CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err, "ConfigMap updation is failed")
		return err
	}
	logger.Info("ConfigMap updation is successful")
	return nil
}

// createConfigMap is a method to create Kubernetes configMap
func createConfigMap(namespace string, configMap *corev1.ConfigMap) error {
	logger := LogGenerator(configMap.Name, namespace, "ConfigMap")
	_, err := GenerateK8sClient().CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "ConfigMap creation is failed")
		return err
	}
	logger.Info("ConfigMap creation is successful")
	return nil
}

// getConfigMap is a method to check configmap exists
//
//nolint:gosimple
func getConfigMap(name, namespace string) (*corev1.ConfigMap, error) {
	configMapInfo, err := GenerateK8sClient().CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return configMapInfo, nil
}

// patchConfigMap will patch Kubernetes configMap
func patchConfigMap(storedConfigMap *corev1.ConfigMap, newConfigMap *corev1.ConfigMap, namespace string) error {
	logger := LogGenerator(storedConfigMap.Name, namespace, "ConfigMap")
	// adding meta fields
	newConfigMap.ResourceVersion = storedConfigMap.ResourceVersion
	newConfigMap.CreationTimestamp = storedConfigMap.CreationTimestamp
	newConfigMap.ManagedFields = storedConfigMap.ManagedFields

	patchResult, err := patch.DefaultPatchMaker.Calculate(storedConfigMap, newConfigMap,
		patch.IgnoreStatusFields(),
		patch.IgnoreField("kind"),
		patch.IgnoreField("apiVersion"),
		patch.IgnoreField("metadata"),
	)
	if err != nil {
		logger.Error(err, "Unable to patch ConfigMap with comparison object")
		return err
	}
	if !patchResult.IsEmpty() {
		for key, value := range storedConfigMap.Annotations {
			if _, present := newConfigMap.Annotations[key]; !present {
				newConfigMap.Annotations[key] = value
			}
		}
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newConfigMap); err != nil {
			logger.Error(err, "Unable to patch ConfigMap with comparison object")
			return err
		}
		logger.Info("Syncing service with defined properties")
		return updateConfigMap(namespace, newConfigMap)
	}
	logger.Info("ConfigMap is already in-sync")
	return nil
}
