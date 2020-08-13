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

package clientservice

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/service"
)

// ClientElasticSearchService is to create client service of elasticsearch
func ClientElasticSearchService(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-client",
		"role":                        "client",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	clientService := service.GenerateService(cr, labels, cr.ObjectMeta.Name, "client")

	clientHeadlessService := service.GenerateHeadlessService(cr, labels, cr.ObjectMeta.Name, "client")

	service.SyncService(cr, clientService, "client")
	service.SyncService(cr, clientHeadlessService, "client-headless")
}
