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
	"sort"

	corev1 "k8s.io/api/core/v1"
	loggingv1beta1 "logging-operator/api/v1beta1"
)

// SetupElasticSearchIngestion is a method to setup elastic data
func SetupElasticSearchIngestion(cr *loggingv1beta1.Elasticsearch) error {
	nodeParams := loggingv1beta1.NodeSpecificConfig{
		KubernetesConfig: cr.Spec.ESIngestion.KubernetesConfig,
		Replicas:         cr.Spec.ESIngestion.Replicas,
		CustomConfig:     cr.Spec.ESIngestion.CustomConfig,
		Storage:          cr.Spec.ESIngestion.Storage,
		JvmMaxMemory:     cr.Spec.ESIngestion.JvmMaxMemory,
		JvmMinMemory:     cr.Spec.ESIngestion.JvmMinMemory,
	}
	envVars := generateEnvVariables(cr, nodeParams)
	envVars = append(envVars, corev1.EnvVar{Name: "discovery.seed_hosts", Value: fmt.Sprintf("%s-master-headless", cr.ObjectMeta.Name)})
	envVars = append(envVars, corev1.EnvVar{Name: "network.host", Value: "0.0.0.0"})
	envVars = append(envVars, corev1.EnvVar{Name: "cluster.name", Value: cr.Spec.ClusterName})
	envVars = append(envVars, corev1.EnvVar{Name: "node.roles", Value: "ingest"})

	if cr.Spec.Security != nil {
		if cr.Spec.Security.TLSEnabled != nil && *cr.Spec.Security.TLSEnabled {
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.enabled", Value: "true"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.enabled", Value: "true"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.verification_mode", Value: "certificate"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.http.ssl.enabled", Value: "true"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.http.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
			envVars = append(envVars, corev1.EnvVar{Name: "xpack.security.http.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		}
	}
	sort.SliceStable(envVars, func(i, j int) bool {
		return envVars[i].Name < envVars[j].Name
	})
	err := CreateElasticsearchStatefulSet(cr, &nodeParams, "ingestion", envVars)
	if err != nil {
		return err
	}
	return nil
}
