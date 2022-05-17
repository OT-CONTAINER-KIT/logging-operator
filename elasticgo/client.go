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
	"crypto/tls"
	"net"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// generateElasticClient is a method to generate client for elasticsearch
func generateElasticClient(cr *loggingv1beta1.Elasticsearch) error {
	var urlScheme, elasticPassword string
	if cr.Spec.Security != nil {
		urlScheme = "https"
	} else {
		urlScheme = "http"
	}
	elasticURL := fmt.Sprintf("%s://%s-master:9200", urlScheme, cr.ObjectMeta.Name)
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticURL,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
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
		return nil, err
	}
	return err, nil
}
