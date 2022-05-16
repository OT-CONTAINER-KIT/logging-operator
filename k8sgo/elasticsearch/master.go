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
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	loggingv1beta1 "logging-operator/api/v1beta1"
)

// SetupElasticSearchMaster is a method to setup elastic master
func SetupElasticSearchMaster(cr *loggingv1beta1.Elasticsearch) error {
	var nodes []string
	nodeParams := loggingv1beta1.NodeSpecificConfig{
		KubernetesConfig: cr.Spec.ESMaster.KubernetesConfig,
		Replicas:         cr.Spec.ESMaster.Replicas,
		CustomConfig:     cr.Spec.ESMaster.CustomConfig,
		Storage:          cr.Spec.ESMaster.Storage,
		JvmMaxMemory:     cr.Spec.ESMaster.JvmMaxMemory,
		JvmMinMemory:     cr.Spec.ESMaster.JvmMinMemory,
	}
	envVars := generateEnvVariables(cr, nodeParams)
	for count := 1; count <= int(*cr.Spec.ESMaster.Replicas); count++ {
		nodes = append(nodes, fmt.Sprintf("%s-master-%s,", cr.ObjectMeta.Name, strconv.Itoa(count)))
	}
	envVars = append(envVars, corev1.EnvVar{Name: "cluster.initial_master_nodes", Value: strings.Join(nodes, "")})
	envVars = append(envVars, corev1.EnvVar{Name: "discovery.seed_hosts", Value: fmt.Sprintf("%s-master-headless", cr.ObjectMeta.Name)})
	envVars = append(envVars, corev1.EnvVar{Name: "network.host", Value: "0.0.0.0"})
	envVars = append(envVars, corev1.EnvVar{Name: "cluster.name", Value: cr.Spec.ClusterName})
	envVars = append(envVars, corev1.EnvVar{Name: "node.data", Value: "false"})
	envVars = append(envVars, corev1.EnvVar{Name: "node.ingest", Value: "false"})
	envVars = append(envVars, corev1.EnvVar{Name: "node.master", Value: "true"})

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
	err := CreateElasticsearchStatefulSet(cr, nodeParams, "master", envVars)
	if err != nil {
		return err
	}
	return nil
}
