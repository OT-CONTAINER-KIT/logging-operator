package ingestionservice

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/service"
)

// IngestionElasticSearchService is to create ingestion service of elasticsearch
func IngestionElasticSearchService(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-ingestion",
		"role":                        "ingestion",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	ingestionService := service.GenerateService(cr, labels, cr.ObjectMeta.Name, "ingestion")

	ingestionHeadlessService := service.GenerateHeadlessService(cr, labels, cr.ObjectMeta.Name, "ingestion")

	service.SyncService(cr, ingestionService, "ingestion")
	service.SyncService(cr, ingestionHeadlessService, "ingestion-headless")
}
