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

package k8selastic

import (
	"encoding/base64"
	"fmt"
	"github.com/thanhpk/randstr"
	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

const (
	elasticsearchCertificateData = "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"
)

// CreateElasticAutoSecret is a method to generate automatic credentials
func CreateElasticAutoSecret(cr *loggingv1beta1.Elasticsearch) error {
	secretName := fmt.Sprintf("%s-password", cr.ObjectMeta.Name)
	labels := map[string]string{
		"app": cr.ObjectMeta.Name,
	}

	secretParams := k8sgo.SecretsParameters{
		Name:        secretName,
		OwnerDef:    k8sgo.ElasticAsOwner(cr),
		Namespace:   cr.Namespace,
		SecretsMeta: k8sgo.GenerateObjectMetaInformation(secretName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		SecretKey:   "password",
		SecretValue: []byte(randstr.String(16)),
	}
	err := k8sgo.CreateSecret(cr.Namespace, k8sgo.GenerateSecret(secretParams))

	if err != nil {
		return err
	}

	return nil
}

// CreateElasticTLSSecret is a method to generate automatic credentials
func CreateElasticTLSSecret(cr *loggingv1beta1.Elasticsearch) error {
	secretName := fmt.Sprintf("%s-tls-cert", cr.ObjectMeta.Name)
	labels := map[string]string{
		"app": cr.ObjectMeta.Name,
	}
	decodedSecret, err := base64.StdEncoding.DecodeString(elasticsearchCertificateData)
	if err != nil {
		return err
	}
	secretParams := k8sgo.SecretsParameters{
		Name:        secretName,
		OwnerDef:    k8sgo.ElasticAsOwner(cr),
		Namespace:   cr.Namespace,
		SecretsMeta: k8sgo.GenerateObjectMetaInformation(secretName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		SecretKey:   "elastic-certificates.p12",
		SecretValue: decodedSecret,
	}
	err = k8sgo.CreateSecret(cr.Namespace, k8sgo.GenerateSecret(secretParams))

	if err != nil {
		return err
	}

	return nil
}
