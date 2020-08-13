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

package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elastic_service")

// SyncService will sync the services of elasticsearch
func SyncService(cr *loggingv1alpha1.Elasticsearch, service *corev1.Service, nodeType string) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Elasticsearch.Name", cr.ObjectMeta.Name, "Node.Type", nodeType)

	elasticServiceName := cr.ObjectMeta.Name + "-" + nodeType
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for service")
	}

	getService, err := k8sClient.CoreV1().Services(cr.Namespace).Get(context.TODO(), elasticServiceName, metav1.GetOptions{})
	if err != nil {
		reqLogger.Info("Creating elasticsearch service", "Elasticsearch.Service.Name", elasticServiceName)
		k8sClient.CoreV1().Services(cr.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	} else if getService != service {
		reqLogger.Info("Updating elasticsearch service", "Elasticsearch.Service.Name", elasticServiceName)
		k8sClient.CoreV1().Services(cr.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Elasticsearch service is already in sync", "Elasticsearch.Service.Name", elasticServiceName)
	}
}

// GenerateHeadlessService generate service definition
func GenerateHeadlessService(cr *loggingv1alpha1.Elasticsearch, labels map[string]string, serviceName string, nodeType string) *corev1.Service {

	serviceAnnotations := identifier.GenerateElasticAnnotations()
	serviceAnnotations["service.alpha.kubernetes.io/tolerate-unready-endpoints"] = "true"
	headlessServiceName := serviceName + "-" + nodeType + "-headless"
	service := &corev1.Service{
		TypeMeta:   identifier.GenerateMetaInformation("Service", "core/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(headlessServiceName, cr.Namespace, labels, serviceAnnotations),
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       9200,
					TargetPort: intstr.FromInt(int(9200)),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "transport",
					Port:       9300,
					TargetPort: intstr.FromInt(int(9200)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
	identifier.AddOwnerRefToObject(service, identifier.ElasticAsOwner(cr))
	return service
}

// GenerateService generate service definition
func GenerateService(cr *loggingv1alpha1.Elasticsearch, labels map[string]string, serviceName string, nodeType string) *corev1.Service {
	ServiceName := serviceName + "-" + nodeType
	service := &corev1.Service{
		TypeMeta:   identifier.GenerateMetaInformation("Service", "core/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(ServiceName, cr.Namespace, labels, identifier.GenerateElasticAnnotations()),
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       9200,
					TargetPort: intstr.FromInt(int(9200)),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "transport",
					Port:       9300,
					TargetPort: intstr.FromInt(int(9200)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
	identifier.AddOwnerRefToObject(service, identifier.ElasticAsOwner(cr))
	return service
}
