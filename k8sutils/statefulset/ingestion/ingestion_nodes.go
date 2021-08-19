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

package ingestion

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sort"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/statefulset"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elastic_ingestion")

func generateIngestionContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {

	reqLogger := log.WithValues("Namespace", cr.Namespace, "Elasticsearch.Name", cr.ObjectMeta.Name, "Node.Type", "ingestion")
	containerDefinition := statefulset.ElasticContainer(cr)

	if cr.Spec.Ingestion.Resources != nil {
		containerDefinition.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Ingestion.Resources.ResourceLimits.CPU)
		containerDefinition.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Ingestion.Resources.ResourceRequests.CPU)
		containerDefinition.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Ingestion.Resources.ResourceLimits.Memory)
		containerDefinition.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Ingestion.Resources.ResourceRequests.Memory)
	}

	if cr.Spec.Ingestion.Storage != nil {
		VolumeMounts := corev1.VolumeMount{Name: cr.ObjectMeta.Name + "-ingestion", MountPath: "/usr/share/elasticsearch/data"}
		containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, VolumeMounts)
	}

	pluginVolume := corev1.VolumeMount{Name: "plugin-volume", MountPath: "/usr/share/elasticsearch/plugins"}
	containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, pluginVolume)

	if *cr.Spec.Security.TLSEnabled != false {
		containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, corev1.VolumeMount{Name: "tls-certificates", MountPath: "/usr/share/elasticsearch/config/certs"})
	}

	ingestionEnvVars := []corev1.EnvVar{
		{Name: "discovery.seed_hosts", Value: cr.ObjectMeta.Name + "-master-headless"},
		{Name: "network.host", Value: "0.0.0.0"},
		{Name: "cluster.name", Value: cr.Spec.ClusterName},
		{Name: "ES_JAVA_OPTS", Value: "-Xmx" + cr.Spec.Ingestion.JVMOptions.Max + " " + "-Xms" + cr.Spec.Ingestion.JVMOptions.Min},
		{Name: "node.data", Value: "false"},
		{Name: "node.ingest", Value: "true"},
		{Name: "node.master", Value: "false"},
		{Name: "node.name", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
	}

	if *cr.Spec.Security.TLSEnabled != false {
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "https"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "ELASTIC_PASSWORD", Value: cr.Spec.Security.Password})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "ELASTIC_USERNAME", Value: "elastic"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.enabled", Value: "true"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.enabled", Value: "true"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.verification_mode", Value: "certificate"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.enabled", Value: "true"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
	} else {
		ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "http"})
	}

	if cr.Spec.Ingestion.ExtraEnvVariables != nil {
		for envName, envValue := range *cr.Spec.Ingestion.ExtraEnvVariables {
			ingestionEnvVars = append(ingestionEnvVars, corev1.EnvVar{Name: envName, Value: envValue})
		}
	}

	sort.SliceStable(ingestionEnvVars, func(i, j int) bool {
		return ingestionEnvVars[i].Name < ingestionEnvVars[j].Name
	})
	containerDefinition.Env = ingestionEnvVars

	reqLogger.Info("Successfully generated the contiainer definition for elasticsearch ingestion")
	return containerDefinition
}

// ElasticSearchIngestion creates the elasticsearch ingestion statefulset
func ElasticSearchIngestion(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-ingestion",
		"role":                        "ingestion",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	statefulsetObject := statefulset.StatefulSetObject(cr, "ingestion", labels, cr.Spec.Ingestion.Count)

	statefulsetObject.Spec.Template.Spec.Containers = append(statefulsetObject.Spec.Template.Spec.Containers, generateIngestionContainer(cr))

	statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.SysctlInitContainer(cr))

	if cr.Spec.Plugins != nil && len(cr.Spec.Plugins) != 0 {
		statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.PluginsInitContainer(cr))
	}

	if cr.Spec.Ingestion.Storage != nil {
		statefulsetObject.Spec.VolumeClaimTemplates = append(statefulsetObject.Spec.VolumeClaimTemplates, statefulset.GeneratePVCTemplate(cr, "ingestion", cr.Spec.Ingestion.Storage))
	}

	if cr.Spec.Ingestion.Affinity != nil {
		statefulsetObject.Spec.Template.Spec.Affinity = cr.Spec.Ingestion.Affinity
	}

	tlsSecretVolume := corev1.Volume{
		Name: "tls-certificates",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: cr.ObjectMeta.Name + "-tls",
			},
		},
	}
	if *cr.Spec.Security.TLSEnabled != false {
		statefulsetObject.Spec.Template.Spec.Volumes = append(statefulsetObject.Spec.Template.Spec.Volumes, tlsSecretVolume)
	}
	statefulset.SyncStatefulSet(cr, statefulsetObject, "ingestion")
}
