package dataservice

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/service"
)

// DataElasticSearchService is to create data service of elasticsearch
func DataElasticSearchService(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-data",
		"role":                        "data",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	dataService := service.GenerateService(cr, labels, cr.ObjectMeta.Name, "data")

	dataHeadlessService := service.GenerateHeadlessService(cr, labels, cr.ObjectMeta.Name, "data")

	service.SyncService(cr, dataService, "data")
	service.SyncService(cr, dataHeadlessService, "data-headless")
}
