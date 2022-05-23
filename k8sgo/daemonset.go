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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/iamabhishek-dubey/k8s-objectmatcher/patch"
)

// DaemonSetParameters is a struct for daemonset inputs
type DaemonSetParameters struct {
	DaemonSetMeta     metav1.ObjectMeta
	OwnerDef          metav1.OwnerReference
	Namespace         string
	ContainerParams   ContainerParams
	Labels            map[string]string
	Annotations       map[string]string
	NodeSelector      map[string]string
	Affinity          *corev1.Affinity
	Tolerations       *[]corev1.Toleration
	PriorityClassName *string
	SecurityContext   *corev1.PodSecurityContext
	Volumes           *[]corev1.Volume
}

// CreateOrUpdateDaemonSet method will create or update DaemonSet
func CreateOrUpdateDaemonSet(params DaemonSetParameters) error {
	logger := LogGenerator(params.DaemonSetMeta.Name, params.Namespace, "DaemonSet")
	storedDaemonSet, err := getDaemonSet(params.Namespace, params.DaemonSetMeta.Name)
	daemonSet := generateDaemonSet(params)
	if err != nil {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(daemonSet); err != nil {
			logger.Error(err, "Unable to patch daemonset with comparison object")
			return err
		}
		if errors.IsNotFound(err) {
			return createDaemonSet(params.Namespace, daemonSet)
		}
		return err
	}
	return patchDaemonSet(storedDaemonSet, daemonSet, params.Namespace)
}

// patchDaemonSet will patch DaemonSet
func patchDaemonSet(storedDaemonSet *appsv1.DaemonSet, newDaemonSet *appsv1.DaemonSet, namespace string) error {
	logger := LogGenerator(storedDaemonSet.Name, namespace, "DaemonSet")
	// adding meta information
	newDaemonSet.ResourceVersion = storedDaemonSet.ResourceVersion
	newDaemonSet.CreationTimestamp = storedDaemonSet.CreationTimestamp
	newDaemonSet.ManagedFields = storedDaemonSet.ManagedFields
	patchResult, err := patch.DefaultPatchMaker.Calculate(storedDaemonSet, newDaemonSet,
		patch.IgnoreStatusFields(),
		patch.IgnoreField("kind"),
		patch.IgnoreField("apiVersion"),
		patch.IgnoreField("metadata"),
	)
	if err != nil {
		logger.Error(err, "Unable to patch DaemonSet with comparison object")
		return err
	}
	if !patchResult.IsEmpty() {
		logger.Info("Changes in daemonset Detected, Updating...", "patch", string(patchResult.Patch))
		for key, value := range storedDaemonSet.Annotations {
			if _, present := newDaemonSet.Annotations[key]; !present {
				newDaemonSet.Annotations[key] = value
			}
		}
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newDaemonSet); err != nil {
			logger.Error(err, "Unable to patch daemonset with comparison object")
			return err
		}
		return updateDaemonSet(namespace, newDaemonSet)
	}
	logger.Info("Reconciliation Complete, no Changes required for DaemonSet")
	return nil
}

// generateDaemonSet is a method to generate description of daemonset
func generateDaemonSet(params DaemonSetParameters) *appsv1.DaemonSet {
	daemonSet := &appsv1.DaemonSet{
		TypeMeta:   GenerateMetaInformation("DaemonSet", "apps/v1"),
		ObjectMeta: params.DaemonSetMeta,
		Spec: appsv1.DaemonSetSpec{
			Selector: LabelSelectors(params.Labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: params.Labels},
				Spec: corev1.PodSpec{
					ServiceAccountName: params.DaemonSetMeta.Name,
					Containers:         generateContainerDef(params.ContainerParams),
					NodeSelector:       params.NodeSelector,
					Affinity:           params.Affinity,
				},
			},
		},
	}
	if params.Volumes != nil {
		daemonSet.Spec.Template.Spec.Volumes = *params.Volumes
	}
	if params.PriorityClassName != nil {
		daemonSet.Spec.Template.Spec.PriorityClassName = *params.PriorityClassName
	}
	if params.Tolerations != nil {
		daemonSet.Spec.Template.Spec.Tolerations = *params.Tolerations
	}
	AddOwnerRefToObject(daemonSet, params.OwnerDef)
	return daemonSet
}

// getDaemonSet is a method to get daemonset in Kubernetes
func getDaemonSet(namespace string, daemonSet string) (*appsv1.DaemonSet, error) {
	logger := LogGenerator(daemonSet, namespace, "DaemonSet")
	daemonSetInfo, err := GenerateK8sClient().AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSet, metav1.GetOptions{})
	if err != nil {
		logger.Info("DaemonSet get action failed")
		return nil, err
	}
	logger.Info("DaemonSet get action was successful")
	return daemonSetInfo, err
}

// createDaemonSet is a method to create daemonset in Kubernetes
func createDaemonSet(namespace string, daemonSet *appsv1.DaemonSet) error {
	logger := LogGenerator(daemonSet.Name, namespace, "DaemonSet")
	_, err := GenerateK8sClient().AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "DaemonSet creation failed")
		return err
	}
	logger.Info("DaemonSet successfully created")
	return nil
}

// updateDaemonSet is a method to update daemonset in Kubernetes
func updateDaemonSet(namespace string, daemonSet *appsv1.DaemonSet) error {
	logger := LogGenerator(daemonSet.Name, namespace, "DaemonSet")
	_, err := GenerateK8sClient().AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err, "DaemonSet update failed")
		return err
	}
	logger.Info("DaemonSet successfully updated")
	return nil
}
