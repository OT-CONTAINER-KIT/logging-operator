/*
Copyright 2022 Opstree Solutions.

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

package elasticgo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// ESClusterDetails is a method for return ESClusterDetails
type ESClusterDetails struct {
	ClusterState string `json:"status"`
	Shards       int32  `json:"active_shards"`
}

// ElasticsearchToken is a interface for elasticsearch token
type ElasticsearchToken struct {
	Created bool `json:"created"`
	Token   struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"token"`
}

// generateElasticClient is a method to generate client for elasticsearch
func generateElasticClient(cr *loggingv1beta1.Elasticsearch) (esapi.Transport, error) {
	logger := k8sgo.LogGenerator(cr.ObjectMeta.Name, cr.Namespace, "Elasticsearch")
	var urlScheme, elasticPassword string
	if cr.Spec.Security != nil {
		urlScheme = "https"
	} else {
		urlScheme = "http"
	}
	elasticURL := fmt.Sprintf("%s://%s-master.%s:9200", urlScheme, cr.ObjectMeta.Name, cr.ObjectMeta.Namespace)
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticURL,
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	if cr.Spec.Security != nil {
		if cr.Spec.Security.ExistingSecret != nil {
			elasticPassword = k8sgo.GetElasticDBPassword(*cr.Spec.Security.ExistingSecret, cr.Namespace)
		} else {
			elasticPassword = k8sgo.GetElasticDBPassword(fmt.Sprintf("%s-password", cr.ObjectMeta.Name), cr.Namespace)
		}

		cfg.Username = "elastic"
		cfg.Password = elasticPassword
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error(err, "Failed in generating elasticsearch client")
		return nil, err
	}
	return es, nil
}

// GetElasticClusterDetails is a method to get health of elastic
func GetElasticClusterDetails(cr *loggingv1beta1.Elasticsearch) (ESClusterDetails, error) {
	var clusterInfo ESClusterDetails
	logger := k8sgo.LogGenerator(cr.ObjectMeta.Name, cr.Namespace, "Elasticsearch")
	esClient, err := generateElasticClient(cr)
	if err != nil {
		logger.Error(err, "Failed in generating elasticsearch client")
		return clusterInfo, err
	}
	req := esapi.ClusterHealthRequest{}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		logger.Error(err, "Error while making request to elasticsearch")
		return clusterInfo, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&clusterInfo)
	if err != nil {
		return clusterInfo, err
	}
	return clusterInfo, nil
}

// CreateServiceAccountTokenKibana is a method to get serviceaccount token for Kibana
func CreateServiceAccountTokenKibana(cr *loggingv1beta1.Elasticsearch) (ElasticsearchToken, error) {
	var tokenInfo ElasticsearchToken
	logger := k8sgo.LogGenerator(cr.ObjectMeta.Name, cr.Namespace, "Elasticsearch")
	esClient, err := generateElasticClient(cr)
	if err != nil {
		logger.Error(err, "Failed in generating elasticsearch client")
		return tokenInfo, err
	}
	req := esapi.SecurityCreateServiceTokenRequest{Namespace: "elastic", Service: "kibana", Name: "token-sa"}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		logger.Error(err, "Error while making request to elasticsearch")
		return tokenInfo, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&tokenInfo)
	if err != nil {
		return tokenInfo, err
	}
	return tokenInfo, nil
}
