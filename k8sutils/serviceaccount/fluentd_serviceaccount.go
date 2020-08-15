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

package serviceaccount

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

var log = logf.Log.WithName("fluentd_serviceaccount")

func generateServiceAccount(cr *loggingv1alpha1.Fluentd, labels map[string]string) *corev1.ServiceAccount {

	serviceAccount := &corev1.ServiceAccount{
		TypeMeta:   identifier.GenerateMetaInformation("ServiceAccount", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
	}
	identifier.AddOwnerRefToObject(serviceAccount, identifier.FluentdAsOwner(cr))
	return serviceAccount
}

// SyncFluentdServiceAccount will create and update the serviceaccount
func SyncFluentdServiceAccount(cr *loggingv1alpha1.Fluentd, account *corev1.ServiceAccount) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "ServiceAccount",
	)

	serviceAccountName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for serviceAccount")
	}

	serviceAccountObject, err := k8sClient.CoreV1().ServiceAccounts(cr.Namespace).Get(context.TODO(), serviceAccountName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating serviceaccount for fluentd", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.CoreV1().ServiceAccounts(cr.Namespace).Create(context.TODO(), account, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if serviceAccountObject != account {
		reqLogger.Info("Updating serviceaccount for fluentd", "Name", cr.ObjectMeta.Name)
		k8sClient.CoreV1().ServiceAccounts(cr.Namespace).Update(context.TODO(), account, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Fluentd serviceaccount are in sync")
	}
}

// CreateFluentdServiceAccount creates the serviceaccount for fluentd configuration
func CreateFluentdServiceAccount(cr *loggingv1alpha1.Fluentd) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}

	config := generateServiceAccount(cr, labels)
	SyncFluentdServiceAccount(cr, config)
}
