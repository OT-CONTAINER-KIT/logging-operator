package elasticutils

import (
	// "bytes"
	// "crypto/tls"
	// "encoding/json"
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	// "net/http"
)

type IndexTemplateType struct {
	IndexPatterns []string `json:"index_patterns"`
	Settings      struct {
		NumberOfShards              int32  `json:"number_of_shards"`
		NumberOfReplicas            int32  `json:"number_of_replicas"`
		IndexLifecycleName          string `json:"index.lifecycle.name"`
		IndexLifecycleRolloverAlias string `json:"index.lifecycle.rollover_alias"`
	} `json:"settings"`
}

func generateIndexTemplate(cr *loggingv1alpha1.IndexTemplate) IndexTemplateType {

	indexTemplate := IndexTemplateType{}

	indexTemplate.IndexPatterns = cr.Spec.IndexPatterns
	indexTemplate.Settings.NumberOfShards = cr.Spec.IndexTemplateSettings.Shards
	indexTemplate.Settings.NumberOfReplicas = cr.Spec.IndexTemplateSettings.Replicas
	indexTemplate.Settings.IndexLifecycleName = cr.Spec.IndexTemplateSettings.IndexLifecycleName
	indexTemplate.Settings.IndexLifecycleRolloverAlias = cr.Spec.IndexTemplateSettings.RollOverAlias

	return indexTemplate
}
