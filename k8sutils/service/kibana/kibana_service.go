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

package kibanaservice

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

var log = logf.Log.WithName("kibana_service")

// generateKibanaService generate service definition for kibana
func generateKibanaService(cr *loggingv1alpha1.Kibana, labels map[string]string) *corev1.Service {
	ServiceName := cr.ObjectMeta.Name
	service := &corev1.Service{
		TypeMeta:   identifier.GenerateMetaInformation("Service", "core/v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(ServiceName, cr.Namespace, labels, identifier.GenerateKibanaAnnotations()),
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       5601,
					TargetPort: intstr.FromInt(int(5601)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
	identifier.AddOwnerRefToObject(service, identifier.KibanaAsOwner(cr))
	return service
}

// SyncKibanaService will sync the services of kibana
func SyncKibanaService(cr *loggingv1alpha1.Kibana, service *corev1.Service) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Kibana.Name", cr.ObjectMeta.Name)

	kibanaServiceName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for service")
	}

	getService, err := k8sClient.CoreV1().Services(cr.Namespace).Get(context.TODO(), kibanaServiceName, metav1.GetOptions{})
	if err != nil {
		reqLogger.Info("Creating kibana service", "Kibana.Service.Name", kibanaServiceName)
		k8sClient.CoreV1().Services(cr.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	} else if getService != service {
		reqLogger.Info("Updating kibana service", "Kibana.Service.Name", kibanaServiceName)
		k8sClient.CoreV1().Services(cr.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Kibana service is already in sync", "Kibana.Service.Name", kibanaServiceName)
	}
}

// CreateKibanaService creates the configmap for kibana deployment
func CreateKibanaService(cr *loggingv1alpha1.Kibana) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Kibana",
	}

	config := generateKibanaService(cr, labels)
	SyncKibanaService(cr, config)
}
