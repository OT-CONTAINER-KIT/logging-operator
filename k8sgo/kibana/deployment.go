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

package k8skibana

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sort"

	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// CreateKibanaSetup is a method to create deployment for Kibana
func CreateKibanaSetup(cr *loggingv1beta1.Kibana) error {
	appName := cr.ObjectMeta.Name
	labels := map[string]string{
		"app":     cr.ObjectMeta.Name,
		"service": "kibana",
	}
	deploymentParams := k8sgo.DeploymentParameters{
		Replicas:       cr.Spec.Replicas,
		OwnerDef:       k8sgo.KibanaAsOwner(cr),
		Namespace:      cr.Namespace,
		DeploymentMeta: k8sgo.GenerateObjectMetaInformation(appName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		ContainerParams: k8sgo.ContainerParams{
			Name:           "kibana",
			Image:          fmt.Sprintf("docker.elastic.co/kibana/kibana:%s", cr.Spec.ElasticConfig.ESVersion),
			VolumeMount:    getVolumeMounts(cr),
			EnvVar:         generateEnvVariables(cr),
			ReadinessProbe: createProbeInfo(),
		},
		Labels:      labels,
		Annotations: k8sgo.GenerateAnnotations(),
		Volumes:     getVolumes(cr),
	}
	if cr.Spec.KubernetesConfig != nil {
		deploymentParams.Affinity = cr.Spec.KubernetesConfig.Affinity
		deploymentParams.NodeSelector = cr.Spec.KubernetesConfig.NodeSelector
		deploymentParams.PriorityClassName = cr.Spec.KubernetesConfig.PriorityClassName
		deploymentParams.Tolerations = cr.Spec.KubernetesConfig.Tolerations
		deploymentParams.ContainerParams.Resources = cr.Spec.KubernetesConfig.Resources
	} else {
		deploymentParams.Affinity = &corev1.Affinity{}
		deploymentParams.NodeSelector = map[string]string{}
		deploymentParams.PriorityClassName = nil
		deploymentParams.Tolerations = &[]corev1.Toleration{}
		deploymentParams.ContainerParams.Resources = &corev1.ResourceRequirements{}
	}
	err := k8sgo.CreateOrUpdateDeployment(deploymentParams)
	if err != nil {
		return err
	}
	return nil
}

// getVolumes is a method to define addtional volumes
func getVolumes(cr *loggingv1beta1.Kibana) *[]corev1.Volume {
	var volumes []corev1.Volume
	if cr.Spec.Security != nil {
		if *cr.Spec.Security.TLSEnabled && cr.Spec.Security.TLSEnabled != nil {
			volumes = append(volumes, corev1.Volume{
				Name: "tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: fmt.Sprintf("%s-tls-cert", cr.Spec.ElasticConfig.ClusterName),
					},
				},
			})
		}
	}
	return &volumes
}

// getVolumes is a method to define volumes mount
func getVolumeMounts(cr *loggingv1beta1.Kibana) *[]corev1.VolumeMount {
	var volumeMounts []corev1.VolumeMount
	if cr.Spec.Security != nil {
		if *cr.Spec.Security.TLSEnabled && cr.Spec.Security.TLSEnabled != nil {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      "tls",
				MountPath: "/usr/share/kibana/config/certs",
			})
		}
	}
	return &volumeMounts
}

// generateEnvVariables is a method to create environment variable for Kibana
func generateEnvVariables(cr *loggingv1beta1.Kibana) []corev1.EnvVar {
	kibanaEnvVars := []corev1.EnvVar{
		{Name: "ELASTICSEARCH_HOSTS", Value: *cr.Spec.ElasticConfig.Host},
		{Name: "SERVER_HOST", Value: "0.0.0.0"},
		{Name: "SERVER_NAME", Value: "kibana"},
	}
	if cr.Spec.Security != nil {
		if *cr.Spec.Security.TLSEnabled && cr.Spec.Security.TLSEnabled != nil {
			kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{Name: "ELASTIC_USERNAME", Value: "elastic"})
			kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{
				Name: "ELASTIC_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s-password", cr.Spec.ElasticConfig.ClusterName),
						},
						Key: "password",
					},
				},
			})
			kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{
				Name: "ELASTICSEARCH_SERVICEACCOUNTTOKEN",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: *cr.Spec.Security.ExistingSecret,
						},
						Key: "token",
					},
				},
			})
			kibanaEnvVars = append(kibanaEnvVars, corev1.EnvVar{Name: "ELASTICSEARCH_SSL_VERIFICATIONMODE", Value: "none"})
		}
	}
	sort.SliceStable(kibanaEnvVars, func(i, j int) bool {
		return kibanaEnvVars[i].Name < kibanaEnvVars[j].Name
	})
	return kibanaEnvVars
}

// createProbeInfo is a method to create probe for k8s
func createProbeInfo() *corev1.Probe {
	return &corev1.Probe{
		InitialDelaySeconds: 15,
		PeriodSeconds:       15,
		FailureThreshold:    5,
		TimeoutSeconds:      5,
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: []string{"bash", "-c", healthCheckScript},
			},
		},
	}
}
