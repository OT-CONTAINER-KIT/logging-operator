package elasticutils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"net/http"
)

type IndexLifeCycleType struct {
	Policy struct {
		Phases struct {
			Hot struct {
				Actions struct {
					Rollover struct {
						MaxSize string `json:"max_size"`
						MaxAge  string `json:"max_age"`
					} `json:"rollover"`
				} `json:"actions"`
			} `json:"hot"`
			Delete struct {
				MinAge  string `json:"min_age"`
				Actions struct {
					Delete struct {
					} `json:"delete"`
				} `json:"actions"`
			} `json:"delete"`
		} `json:"phases"`
	} `json:"policy"`
}

// CreateUpdateIndexLifeCyclePolicy will create and update the index lifecycle policy
func CreateUpdateIndexLifeCyclePolicy(cr *loggingv1alpha1.IndexLifecycle) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Policy.Name", cr.ObjectMeta.Name, "Log.Type", "IndexLifeCycle Policy")
	policyData, err := json.Marshal(generateLifeCyclePolicy(cr))

	if err != nil {
		reqLogger.Error(err, "Error while generating lifecycle policy data")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	requestBody := bytes.NewReader(policyData)

	req, err := http.NewRequest("PUT", *cr.Spec.Elasticsearch.Host+"/_ilm/policy/"+cr.ObjectMeta.Name+"?pretty", requestBody)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while creating index lifecycle policy")
	}
	defer resp.Body.Close()
	reqLogger.Info("Successfully created the index lifecycle policy")
}

// DeleteIndexLifeCyclePolicy will delete the index lifecycle policy
func DeleteIndexLifeCyclePolicy(cr *loggingv1alpha1.IndexLifecycle) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Policy.Name", cr.ObjectMeta.Name, "Log.Type", "IndexLifeCycle Policy")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", *cr.Spec.Elasticsearch.Host+"/_ilm/policy/"+cr.ObjectMeta.Name, nil)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while deleting index lifecycle policy")
	}
	defer resp.Body.Close()
	reqLogger.Info("Successfully deleted the index lifecycle policy")
}

// CompareandUpdatePolicy will compare and create the index lifecycle policy
func CompareandUpdatePolicy(cr *loggingv1alpha1.IndexLifecycle) {
	reqLogger := log.WithValues("Namespace", cr.Namespace, "Policy.Name", cr.ObjectMeta.Name, "Log.Type", "IndexLifeCycle Policy")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", *cr.Spec.Elasticsearch.Host+"/_ilm/policy/"+cr.ObjectMeta.Name, nil)
	if err != nil {
		reqLogger.Error(err, "Error while generating request information")
	}

	if cr.Spec.Elasticsearch.Username != nil && cr.Spec.Elasticsearch.Password != nil {
		req.SetBasicAuth(*cr.Spec.Elasticsearch.Username, *cr.Spec.Elasticsearch.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		reqLogger.Error(err, "Request failed while getting index lifecycle policy")
		CreateUpdateIndexLifeCyclePolicy(cr)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var existingData IndexLifeCycleType
	err = decoder.Decode(&existingData)

	if existingData != generateLifeCyclePolicy(cr) {
		CreateUpdateIndexLifeCyclePolicy(cr)
	}
	return
}

func generateLifeCyclePolicy(cr *loggingv1alpha1.IndexLifecycle) IndexLifeCycleType {

	lifeCyclePolicy := IndexLifeCycleType{}

	lifeCyclePolicy.Policy.Phases.Hot.Actions.Rollover.MaxSize = cr.Spec.Rollover.MaxSize
	lifeCyclePolicy.Policy.Phases.Hot.Actions.Rollover.MaxAge = cr.Spec.Rollover.MaxAge

	lifeCyclePolicy.Policy.Phases.Delete.MinAge = cr.Spec.Delete.MinAge

	return lifeCyclePolicy
}
