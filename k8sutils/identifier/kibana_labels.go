package identifier

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KibanaAsOwner generates and returns object reference
func KibanaAsOwner(cr *loggingv1alpha1.Kibana) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: cr.APIVersion,
		Kind:       cr.Kind,
		Name:       cr.Name,
		UID:        cr.UID,
		Controller: &trueVar,
	}
}

// GenerateKibanaAnnotations generates and returns kibana annotations
func GenerateKibanaAnnotations() map[string]string {
	return map[string]string{
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Kibana",
	}
}
