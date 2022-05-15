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

package k8selastic

import (
	"fmt"

	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// CreateElasticSearchService is a method to create service for ES
func CreateElasticSearchService(cr *loggingv1beta1.Elasticsearch, role string) error {
	appName := fmt.Sprintf("%s-%s", cr.ObjectMeta.Name, role)
	labels := map[string]string{
		"app":  appName,
		"role": role,
	}
	serviceParams := k8sgo.ServiceParameters{
		ServiceMeta:              k8sgo.GenerateObjectMetaInformation(appName, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		OwnerDef:                 k8sgo.ElasticAsOwner(cr),
		Labels:                   labels,
		Annotations:              k8sgo.GenerateAnnotations(),
		Namespace:                cr.Namespace,
		HeadlessService:          false,
		PublishNotReadyAddresses: false,
		Port: []k8sgo.PortInfo{
			{
				PortName: "http",
				Port:     9200,
			},
			{
				PortName: "transport",
				Port:     9300,
			},
		},
	}

	err := k8sgo.CreateOrUpdateService(serviceParams)
	if err != nil {
		return err
	}
	serviceParams = k8sgo.ServiceParameters{
		ServiceMeta:              k8sgo.GenerateObjectMetaInformation(fmt.Sprintf("%s-%s", appName, "headless"), cr.Namespace, labels, k8sgo.GenerateAnnotations()),
		OwnerDef:                 k8sgo.ElasticAsOwner(cr),
		Labels:                   labels,
		Annotations:              k8sgo.GenerateAnnotations(),
		Namespace:                cr.Namespace,
		HeadlessService:          true,
		PublishNotReadyAddresses: true,
		Port: []k8sgo.PortInfo{
			{
				PortName: "http",
				Port:     9200,
			},
			{
				PortName: "transport",
				Port:     9300,
			},
		},
	}
	err = k8sgo.CreateOrUpdateService(serviceParams)
	if err != nil {
		return err
	}
	return nil
}
