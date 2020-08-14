package elasticutils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	elasticURL string
)

var log = logf.Log.WithName("elastic_healthcheck")

// GetElasticHealth will return the health of elasticsearch service
func GetElasticHealth(cr *loggingv1alpha1.Elasticsearch) (*string, error) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Elasticsearch.Name", cr.ObjectMeta.Name, "Log.Type", "healthcheck")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if *cr.Spec.Security.TLSEnabled != false {
		elasticURL = "https://" + cr.ObjectMeta.Name + "-master:9200"
	} else {
		elasticURL = "http://" + cr.ObjectMeta.Name + "-master:9200"
	}

	req, err := http.NewRequest("GET", elasticURL+"/_cluster/health", nil)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
		return nil, err
	}

	if cr.Spec.Security.Password != "" {
		req.SetBasicAuth("elastic", cr.Spec.Security.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Error while capturing the response")
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		reqLogger.Error(err, "Error reading response from elasticsearch")
		return nil, err
	}
	var responseFilter map[string]interface{}
	json.Unmarshal(responseData, &responseFilter)

	status := fmt.Sprintf("%v", responseFilter["status"])
	return &status, nil
}
