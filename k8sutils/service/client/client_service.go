package clientservice

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/service"
)

// ClientElasticSearchService is to create client service of elasticsearch
func ClientElasticSearchService(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-client",
		"role":                        "client",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	clientService := service.GenerateService(cr, labels, cr.ObjectMeta.Name, "client")

	clientHeadlessService := service.GenerateHeadlessService(cr, labels, cr.ObjectMeta.Name, "client")

	service.SyncService(cr, clientService, "client")
	service.SyncService(cr, clientHeadlessService, "client-headless")
}
