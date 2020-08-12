package secret

import (
	"context"
	"encoding/base64"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elastic_secret")

// GenerateTLSecret is a method to generate the secret for Elasticsearch TLS
func GenerateTLSecret(cr *loggingv1alpha1.Elasticsearch) *corev1.Secret {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Secret.Name", cr.ObjectMeta.Name)

	decoded, err := base64.StdEncoding.DecodeString(elasticsearchCertificateData)
	if err != nil {
		reqLogger.Error(err, "Unable to decode certificate")
	}
	labels := map[string]string{
		"name":                        cr.ObjectMeta.Name + "-tls",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	secret := &corev1.Secret{
		TypeMeta:   identifier.GenerateMetaInformation("Secret", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name+"-tls", cr.Namespace, labels, identifier.GenerateElasticAnnotations()),
		Data: map[string][]byte{
			"elastic-certificates.p12": decoded,
		},
	}
	identifier.AddOwnerRefToObject(secret, identifier.ElasticAsOwner(cr))
	return secret
}

// CreateAndUpdateSecret is a method for creating secret
func CreateAndUpdateSecret(cr *loggingv1alpha1.Elasticsearch, secretBody *corev1.Secret) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Secret.Name", cr.ObjectMeta.Name)

	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for secret")
	}

	secretName, err := k8sClient.CoreV1().Secrets(cr.Namespace).Get(context.TODO(), cr.ObjectMeta.Name, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating secret for elasticsearch", "Secret.Name", cr.ObjectMeta.Name)
		k8sClient.CoreV1().Secrets(cr.Namespace).Create(context.TODO(), secretBody, metav1.CreateOptions{})
	} else if secretBody != secretName {
		reqLogger.Info("Reconciling secret for elasticsearch", "Secret.Name", cr.ObjectMeta.Name)
		k8sClient.CoreV1().Secrets(cr.Namespace).Update(context.TODO(), secretBody, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Elasticsearch secret is already synced", "Secret.Name", cr.ObjectMeta.Name)
	}
}

// GenerateElasticPassword will generate the passowrd field for elasticsearch
func GenerateElasticPassword(cr *loggingv1alpha1.Elasticsearch) *corev1.Secret {
	password := []byte(cr.Spec.Security.Password)
	labels := map[string]string{
		"name":                        cr.ObjectMeta.Name + "-password",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}
	secret := &corev1.Secret{
		TypeMeta:   identifier.GenerateMetaInformation("Secret", "v1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name+"-password", cr.Namespace, labels, identifier.GenerateElasticAnnotations()),
		Data: map[string][]byte{
			"password": password,
		},
	}
	identifier.AddOwnerRefToObject(secret, identifier.ElasticAsOwner(cr))
	return secret
}
