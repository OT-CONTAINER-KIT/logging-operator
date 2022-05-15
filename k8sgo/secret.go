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

// SecretsParameters is an interface for secret input
type SecretsParameters struct {
	Name        string
	OwnerDef    metav1.OwnerReference
	Namespace   string
	SecretsMeta metav1.ObjectMeta
	SecretKey   string
	SecretValue []byte
}

// GenerateSecret is a method that will generate a secret interface
func GenerateSecret(params SecretsParameters) *corev1.Secret {
	secret := &corev1.Secret{
		TypeMeta:   GenerateMetaInformation("Secret", "v1"),
		ObjectMeta: params.SecretsMeta,
		Data: map[string][]byte{
			params.SecretKey: params.SecretValue,
		},
	}
	AddOwnerRefToObject(secret, params.OwnerDef)
	return secret
}

// createSecret is a method to create Kubernetes secrets
func CreateSecret(namespace string, secret *corev1.Secret) error {
	logger := LogGenerator(secret.Name, namespace, "Secret")
	_, err := GenerateK8sClient().CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Secret creation is failed")
		return err
	}
	logger.Info("Secret creation is successful")
	return nil
}

//nolint:gosimple
// GetSecret is a method to check secret exists
func GetSecret(name, namespace string) (*corev1.Secret, error) {
	secretInfo, err := GenerateK8sClient().CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secretInfo, nil
}
