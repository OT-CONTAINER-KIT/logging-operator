package masterservice

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/service"
)

// MasterElasticSearchService is to create master service of elasticsearch
func MasterElasticSearchService(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-master",
		"role":                        "master",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	masterService := service.GenerateService(cr, labels, cr.ObjectMeta.Name, "master")

	masterHeadlessService := service.GenerateHeadlessService(cr, labels, cr.ObjectMeta.Name, "master")

	service.SyncService(cr, masterService, "master")
	service.SyncService(cr, masterHeadlessService, "master-headless")
}
