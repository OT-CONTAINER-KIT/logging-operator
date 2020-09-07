package elasticutils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"net/http"
	"reflect"
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

// CreateUpdateIndexTemplate will create and update the index template
func CreateUpdateIndexTemplate(cr *loggingv1alpha1.IndexTemplate) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Template.Name", cr.ObjectMeta.Name, "Resource.Type", "Index Template")
	indexTemplateData, err := json.Marshal(generateIndexTemplate(cr))

	if err != nil {
		reqLogger.Error(err, "Error while generating index template data")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	requestBody := bytes.NewReader(indexTemplateData)

	req, err := http.NewRequest("PUT", *cr.Spec.Elasticsearch.Host+"/_template/"+cr.ObjectMeta.Name+"?pretty", requestBody)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while creating index template")
	}
	defer resp.Body.Close()
	reqLogger.Info("Successfully created the index template")
}

// DeleteIndexTemplate will delete the index template
func DeleteIndexTemplate(cr *loggingv1alpha1.IndexTemplate) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Template.Name", cr.ObjectMeta.Name, "Resource.Type", "Index Template")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", *cr.Spec.Elasticsearch.Host+"/_template/"+cr.ObjectMeta.Name, nil)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while deleting index template")
	}
	defer resp.Body.Close()
	reqLogger.Info("Successfully deleted the index template")
}

// CompareandUpdateIndexTemplate will compare and create the index template
func CompareandUpdateIndexTemplate(cr *loggingv1alpha1.IndexTemplate) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Template.Name", cr.ObjectMeta.Name, "Resource.Type", "Index Template")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", *cr.Spec.Elasticsearch.Host+"/_template/"+cr.ObjectMeta.Name, nil)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while getting index template")
		CreateUpdateIndexTemplate(cr)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var existingData IndexTemplateType
	err = decoder.Decode(&existingData)

	compareResult := reflect.DeepEqual(existingData, generateIndexTemplate(cr))
	if compareResult != true {
		CreateUpdateIndexTemplate(cr)
	}
	return
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
