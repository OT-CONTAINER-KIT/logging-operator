package labels

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateElasticMetaInformation generates the meta information
func GenerateElasticMetaInformation(resourceKind string, apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       resourceKind,
		APIVersion: apiVersion,
	}
}

// GenerateElasticObjectMetaInformation generates the object meta information
func GenerateElasticObjectMetaInformation(name string, namespace string, labels map[string]string, annotations map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      labels,
		Annotations: annotations,
	}
}

// ElasticAsOwner generates and returns object refernece
func ElasticAsOwner(cr *loggingv1alpha1.Elasticsearch) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: cr.APIVersion,
		Kind:       cr.Kind,
		Name:       cr.Name,
		UID:        cr.UID,
		Controller: &trueVar,
	}
}

// GenerateElasticAnnotations generates and returns statefulsets annotations
func GenerateElasticAnnotations() map[string]string {
	return map[string]string{
		"logging.opstreelabs.in": "true",
		"logging.opstreelabs.in/kind": "Elasticsearch"
		"prometheus.io/scrape": "true",
		"prometheus.io/port":   "9121",
	}
}
