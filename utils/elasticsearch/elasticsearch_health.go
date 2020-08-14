/*
Copyright 2020 Opstree Solutions.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
