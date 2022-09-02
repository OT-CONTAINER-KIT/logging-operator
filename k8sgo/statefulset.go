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
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/iamabhishek-dubey/k8s-objectmatcher/patch"
)

// StatefulSetParameters is the input struct for statefulset
type StatefulSetParameters struct {
	StatefulSetMeta   metav1.ObjectMeta
	OwnerDef          metav1.OwnerReference
	Namespace         string
	ContainerParams   ContainerParams
	Labels            map[string]string
	Annotations       map[string]string
	Replicas          *int32
	PVCParameters     PVCParameters
	Affinity          *corev1.Affinity
	NodeSelector      map[string]string
	Tolerations       *[]corev1.Toleration
	PriorityClassName *string
	SecurityContext   *corev1.PodSecurityContext
	ExtraVolumes      *[]corev1.Volume
	ESPlugins         *[]string
}

// PVCParameters is a struct to pass arguments for PVC
type PVCParameters struct {
	Name             string
	Namespace        string
	Labels           map[string]string
	Annotations      map[string]string
	AccessModes      []corev1.PersistentVolumeAccessMode
	StorageClassName *string
	StorageSize      string
}

// CreateOrUpdateStateFul method will create or update StatefulSet
func CreateOrUpdateStateFul(params StatefulSetParameters) error {
	logger := LogGenerator(params.StatefulSetMeta.Name, params.Namespace, "StatefulSet")
	storedStateful, err := GetStateFulSet(params.Namespace, params.StatefulSetMeta.Name)
	statefulSetDef := generateStatefulSetDef(params)
	if err != nil {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(statefulSetDef); err != nil {
			logger.Error(err, "Unable to patch statefulset with comparison object")
			return err
		}
		if errors.IsNotFound(err) {
			return createStateFulSet(params.Namespace, statefulSetDef)
		}
		return err
	}
	return patchStateFulSet(storedStateful, statefulSetDef, params.Namespace)
}

// patchStateFulSet will patch Statefulset
func patchStateFulSet(storedStateful *appsv1.StatefulSet, newStateful *appsv1.StatefulSet, namespace string) error {
	logger := LogGenerator(storedStateful.Name, namespace, "StatefulSet")
	// adding meta information
	newStateful.ResourceVersion = storedStateful.ResourceVersion
	newStateful.CreationTimestamp = storedStateful.CreationTimestamp
	newStateful.ManagedFields = storedStateful.ManagedFields
	patchResult, err := patch.DefaultPatchMaker.Calculate(storedStateful, newStateful,
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
		patch.IgnorePersistenVolumeFields(),
		patch.IgnoreField("kind"),
		patch.IgnoreField("apiVersion"),
		patch.IgnoreField("metadata"),
	)
	if err != nil {
		logger.Error(err, "Unable to patch statefulset with comparison object")
		return err
	}
	if !patchResult.IsEmpty() {
		logger.Info("Changes in statefulset Detected, Updating...", "patch", string(patchResult.Patch))
		for key, value := range storedStateful.Annotations {
			if _, present := newStateful.Annotations[key]; !present {
				newStateful.Annotations[key] = value
			}
		}
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newStateful); err != nil {
			logger.Error(err, "Unable to patch statefulset with comparison object")
			return err
		}
		return updateStateFulSet(namespace, newStateful)
	}
	logger.Info("Reconciliation Complete, no Changes required.")
	return nil
}

// createStateFulSet is a method to create statefulset in Kubernetes
func createStateFulSet(namespace string, stateful *appsv1.StatefulSet) error {
	logger := LogGenerator(stateful.Name, namespace, "StatefulSet")
	_, err := GenerateK8sClient().AppsV1().StatefulSets(namespace).Create(context.TODO(), stateful, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Statefulset creation failed")
		return err
	}
	logger.Info("Statefulset successfully created")
	return nil
}

// updateStateFulSet is a method to update statefulset in Kubernetes
func updateStateFulSet(namespace string, stateful *appsv1.StatefulSet) error {
	logger := LogGenerator(stateful.Name, namespace, "StatefulSet")
	_, err := GenerateK8sClient().AppsV1().StatefulSets(namespace).Update(context.TODO(), stateful, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err, "Statefulset update failed")
		return err
	}
	logger.Info("Statefulset successfully updated")
	return nil
}

// GetStateFulSet is a method to get statefulset in Kubernetes
func GetStateFulSet(namespace string, stateful string) (*appsv1.StatefulSet, error) {
	logger := LogGenerator(stateful, namespace, "StatefulSet")
	statefulInfo, err := GenerateK8sClient().AppsV1().StatefulSets(namespace).Get(context.TODO(), stateful, metav1.GetOptions{})
	if err != nil {
		logger.Info("Statefulset get action failed")
		return nil, err
	}
	logger.Info("Statefulset get action was successful")
	return statefulInfo, err
}

// generateStatefulSetDef is a method to generate statefulset definition
func generateStatefulSetDef(params StatefulSetParameters) *appsv1.StatefulSet {
	var serviceLink = true
	var runasUser int64 = 1000
	var fsGroup int64 = 1000
	statefulset := &appsv1.StatefulSet{
		TypeMeta:   GenerateMetaInformation("StatefulSet", "apps/v1"),
		ObjectMeta: params.StatefulSetMeta,
		Spec: appsv1.StatefulSetSpec{
			Selector:            LabelSelectors(params.Labels),
			ServiceName:         params.StatefulSetMeta.Name,
			Replicas:            params.Replicas,
			PodManagementPolicy: appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type:          appsv1.RollingUpdateStatefulSetStrategyType,
				RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: params.Labels},
				Spec: corev1.PodSpec{
					Containers:   generateContainerDef(params.ContainerParams),
					NodeSelector: params.NodeSelector,
					Affinity:     params.Affinity,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:   &fsGroup,
						RunAsUser: &runasUser,
					},
					InitContainers:     []corev1.Container{getInitContainer(params.ContainerParams)},
					EnableServiceLinks: &serviceLink,
				},
			},
		},
	}

	if params.ESPlugins != nil {
		statefulset.Spec.Template.Spec.InitContainers = append(statefulset.Spec.Template.Spec.InitContainers, getPluginInitContainers(params))
	}
	if params.ExtraVolumes != nil {
		statefulset.Spec.Template.Spec.Volumes = *params.ExtraVolumes
	}
	if params.PriorityClassName != nil {
		statefulset.Spec.Template.Spec.PriorityClassName = *params.PriorityClassName
	}
	if params.Tolerations != nil {
		statefulset.Spec.Template.Spec.Tolerations = *params.Tolerations
	}
	statefulset.Spec.VolumeClaimTemplates = append(statefulset.Spec.VolumeClaimTemplates, generatePersistentVolumeTemplate(params.PVCParameters))
	AddOwnerRefToObject(statefulset, params.OwnerDef)
	return statefulset
}

// generatePersistentVolumeTemplate is a method to create the persistent volume claim template
func generatePersistentVolumeTemplate(params PVCParameters) corev1.PersistentVolumeClaim {
	return corev1.PersistentVolumeClaim{
		TypeMeta:   GenerateMetaInformation("PersistentVolumeClaim", "v1"),
		ObjectMeta: metav1.ObjectMeta{Name: params.Name},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: params.AccessModes,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(params.StorageSize),
				},
			},
			StorageClassName: params.StorageClassName,
		},
	}
}

// getInitContainer is a method to create Init Container
func getInitContainer(params ContainerParams) corev1.Container {
	var privileged = true
	var runasUser int64 = 0
	return corev1.Container{
		Name:    "sysctl-init",
		Image:   params.Image,
		Command: []string{"sysctl", "-w", "vm.max_map_count=262144"},
		SecurityContext: &corev1.SecurityContext{
			Privileged: &privileged,
			RunAsUser:  &runasUser,
		},
	}
}

// getPluginInitContainers is a method to create plugins init container
func getPluginInitContainers(params StatefulSetParameters) corev1.Container {
	shellCommand := []string{"sh", "-c"}
	command := []string{"bin/elasticsearch-plugin install --batch"}
	command = append(command, *params.ESPlugins...)

	shellCommand = append(shellCommand, strings.Join(command, " "))
	return corev1.Container{
		Name:    "plugins",
		Image:   params.ContainerParams.Image,
		Command: shellCommand,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "plugin-volume",
				MountPath: "/usr/share/elasticsearch/plugins",
			},
		},
	}
}
