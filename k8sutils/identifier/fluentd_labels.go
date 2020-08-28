package identifier

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FluentdAsOwner generates and returns object reference
func FluentdAsOwner(cr *loggingv1alpha1.Fluentd) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: cr.APIVersion,
		Kind:       cr.Kind,
		Name:       cr.Name,
		UID:        cr.UID,
		Controller: &trueVar,
	}
}

// GenerateFluentdAnnotations generates and returns fluentd annotations
func GenerateFluentdAnnotations() map[string]string {
	return map[string]string{
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}
}
