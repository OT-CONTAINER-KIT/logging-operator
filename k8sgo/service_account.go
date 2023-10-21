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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceAccountParameters is an interface for serviceaccount input
type ServiceAccountParameters struct {
	Name               string
	OwnerDef           metav1.OwnerReference
	Namespace          string
	ServiceAccountMeta metav1.ObjectMeta
}

// GenerateServiceAccount is a method that will generate a serviceaccount interface
func GenerateServiceAccount(params ServiceAccountParameters) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta:   GenerateMetaInformation("ServiceAccount", "v1"),
		ObjectMeta: params.ServiceAccountMeta,
	}
	AddOwnerRefToObject(serviceAccount, params.OwnerDef)
	return serviceAccount
}

// CreateServiceAccount is a method to create Kubernetes serviceaccount
func CreateServiceAccount(namespace string, serviceAccount *corev1.ServiceAccount) error {
	logger := LogGenerator(serviceAccount.Name, namespace, "ServiceAccount")
	_, err := GenerateK8sClient().CoreV1().ServiceAccounts(namespace).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "ServiceAccount creation is failed")
		return err
	}
	logger.Info("ServiceAccount creation is successful")
	return nil
}

// GetServiceAccount is a method to check serviceaccount exists
//
//nolint:gosimple
func GetServiceAccount(name, namespace string) (*corev1.ServiceAccount, error) {
	serviceAccountInfo, err := GenerateK8sClient().CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return serviceAccountInfo, nil
}
