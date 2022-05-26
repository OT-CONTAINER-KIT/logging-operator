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

// DeploymentParameters is a struct for deployment inputs
type DeploymentParameters struct {
	Replicas          *int32
	DeploymentMeta    metav1.ObjectMeta
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

// CreateOrUpdateDeployment method will create or update deployment
func CreateOrUpdateDeployment(params DeploymentParameters) error {
	logger := LogGenerator(params.DeploymentMeta.Name, params.Namespace, "Deployment")
	storedDeployment, err := getDeployment(params.Namespace, params.DeploymentMeta.Name)
	deployment := generateDeploymentParams(params)
	if err != nil {
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(deployment); err != nil {
			logger.Error(err, "Unable to patch deployment with comparison object")
			return err
		}
		if errors.IsNotFound(err) {
			return createDeployment(params.Namespace, deployment)
		}
		return err
	}
	return patchDeployment(storedDeployment, deployment, params.Namespace)
}

// patchDeployment will patch Deployment
func patchDeployment(storedDeployment *appsv1.Deployment, newDeployment *appsv1.Deployment, namespace string) error {
	logger := LogGenerator(storedDeployment.Name, namespace, "Deployment")
	// adding meta information
	newDeployment.ResourceVersion = storedDeployment.ResourceVersion
	newDeployment.CreationTimestamp = storedDeployment.CreationTimestamp
	newDeployment.ManagedFields = storedDeployment.ManagedFields
	patchResult, err := patch.DefaultPatchMaker.Calculate(storedDeployment, newDeployment,
		patch.IgnoreStatusFields(),
		patch.IgnoreField("kind"),
		patch.IgnoreField("apiVersion"),
		patch.IgnoreField("metadata"),
	)
	if err != nil {
		logger.Error(err, "Unable to patch Deployment with comparison object")
		return err
	}
	if !patchResult.IsEmpty() {
		logger.Info("Changes in deployment Detected, Updating...", "patch", string(patchResult.Patch))
		for key, value := range storedDeployment.Annotations {
			if _, present := newDeployment.Annotations[key]; !present {
				newDeployment.Annotations[key] = value
			}
		}
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newDeployment); err != nil {
			logger.Error(err, "Unable to patch deployment with comparison object")
			return err
		}
		return updateDeployment(namespace, newDeployment)
	}
	logger.Info("Reconciliation Complete, no Changes required for Deployment")
	return nil
}

// generateDeploymentParams is a method to generate description of deployment
func generateDeploymentParams(params DeploymentParameters) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta:   GenerateMetaInformation("Deployment", "apps/v1"),
		ObjectMeta: params.DeploymentMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: params.Replicas,
			Selector: LabelSelectors(params.Labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: params.Labels},
				Spec: corev1.PodSpec{
					Containers:   generateContainerDef(params.ContainerParams),
					NodeSelector: params.NodeSelector,
					Affinity:     params.Affinity,
				},
			},
		},
	}
	if params.Volumes != nil {
		deployment.Spec.Template.Spec.Volumes = *params.Volumes
	}
	if params.PriorityClassName != nil {
		deployment.Spec.Template.Spec.PriorityClassName = *params.PriorityClassName
	}
	if params.Tolerations != nil {
		deployment.Spec.Template.Spec.Tolerations = *params.Tolerations
	}
	AddOwnerRefToObject(deployment, params.OwnerDef)
	return deployment
}

// getDeployment is a method to get deployment in Kubernetes
func getDeployment(namespace string, deployment string) (*appsv1.Deployment, error) {
	logger := LogGenerator(deployment, namespace, "Deployment")
	deploymentInfo, err := GenerateK8sClient().AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		logger.Info("Deployment get action failed")
		return nil, err
	}
	logger.Info("Deployment get action was successful")
	return deploymentInfo, err
}

// createDeployment is a method to create deployment in Kubernetes
func createDeployment(namespace string, deployment *appsv1.Deployment) error {
	logger := LogGenerator(deployment.Name, namespace, "Deployment")
	_, err := GenerateK8sClient().AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Deployment creation failed")
		return err
	}
	logger.Info("Deployment successfully created")
	return nil
}

// updateDeployment is a method to update deployment in Kubernetes
func updateDeployment(namespace string, deployment *appsv1.Deployment) error {
	logger := LogGenerator(deployment.Name, namespace, "Deployment")
	_, err := GenerateK8sClient().AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err, "Deployment update failed")
		return err
	}
	logger.Info("Deployment successfully updated")
	return nil
}
