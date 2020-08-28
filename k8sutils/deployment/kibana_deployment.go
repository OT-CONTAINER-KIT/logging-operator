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

package deployment

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

var log = logf.Log.WithName("kibana_deployment")

func generateKibanaContainer(cr *loggingv1alpha1.Kibana) *corev1.Container {

	var runasUser int64 = 1000
	var runasNonRoot = true

	containerDefinition := &corev1.Container{
		Name:            "kibana",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Env:             []corev1.EnvVar{},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{}, Requests: corev1.ResourceList{},
		},
		VolumeMounts: []corev1.VolumeMount{},
		SecurityContext: &corev1.SecurityContext{
			RunAsNonRoot: &runasNonRoot,
			RunAsUser:    &runasUser,
		},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: 10,
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

	if cr.Spec.Resources != nil {
		containerDefinition.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Resources.ResourceLimits.CPU)
		containerDefinition.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Resources.ResourceRequests.CPU)
		containerDefinition.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Resources.ResourceLimits.Memory)
		containerDefinition.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Resources.ResourceRequests.Memory)
	}

	volumeMounts := []corev1.VolumeMount{
		{Name: "kibanaconfig", MountPath: "/usr/share/kibana/config/kibana.yml", SubPath: "kibana.yml"},
	}

	kibanaEnvVars := []corev1.EnvVar{
		{Name: "ELASTICSEARCH_HOSTS", Value: cr.Spec.KibanaElasticsearch.Host},
		{Name: "SERVER_HOST", Value: "0.0.0.0"},
	}

	if cr.Spec.KibanaElasticsearch.TLSEnabled != false {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "tls",
			MountPath: "/usr/share/kibana/config/certs",
		})
		kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{Name: "ELASTICSEARCH_USERNAME", Value: cr.Spec.KibanaElasticsearch.Username})
		kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{Name: "ELASTICSEARCH_PASSWORD", Value: cr.Spec.KibanaElasticsearch.Password})
	}

	containerDefinition.VolumeMounts = volumeMounts
	containerDefinition.Env = kibanaEnvVars

	return containerDefinition
}

func generateKibanaDeployment(cr *loggingv1alpha1.Kibana, labels map[string]string) *appsv1.Deployment {

	kibanaContainer := generateKibanaContainer(cr)

	deploymentObject := &appsv1.Deployment{
		TypeMeta:   identifier.GenerateMetaInformation("Deployment", "apps/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateKibanaAnnotations()),
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.Replicas,
			Selector: identifier.LabelSelectors(labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*kibanaContainer,
					},
					Volumes: []corev1.Volume{
						{
							Name: "kibanaconfig",
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

	if cr.Spec.KibanaElasticsearch.TLSEnabled != false {
		deploymentObject.Spec.Template.Spec.Volumes = append(deploymentObject.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: *cr.Spec.ElasticSecretName,
				},
			},
		})
	}

	if cr.Spec.Affinity != nil {
		deploymentObject.Spec.Template.Spec.Affinity = cr.Spec.Affinity
	}
	identifier.AddOwnerRefToObject(deploymentObject, identifier.KibanaAsOwner(cr))
	return deploymentObject
}

// SyncKibanaDeployment will sync the deployment in Kubernetes
func SyncKibanaDeployment(cr *loggingv1alpha1.Kibana, deploy *appsv1.Deployment) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "Deployment",
	)

	deploymentName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for deployment")
	}

	deploymentObject, err := k8sClient.AppsV1().Deployments(cr.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating deployment setup", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.AppsV1().Deployments(cr.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if deploymentObject != deploy {
		reqLogger.Info("Updating deployment setup Kibana", "Name", cr.ObjectMeta.Name)
		k8sClient.AppsV1().Deployments(cr.Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Kibana deployment are in sync")
	}
}

// CreateKibanaDeployment creates the configmap for kibana deployment
func CreateKibanaDeployment(cr *loggingv1alpha1.Kibana) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Kibana",
	}

	config := generateKibanaDeployment(cr, labels)
	SyncKibanaDeployment(cr, config)
}
