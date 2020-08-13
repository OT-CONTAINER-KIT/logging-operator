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
