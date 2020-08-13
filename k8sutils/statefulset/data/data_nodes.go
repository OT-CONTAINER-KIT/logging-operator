package data

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/statefulset"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elastic_data")

func generateDataContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {

	reqLogger := log.WithValues("Namespace", cr.Namespace, "Elasticsearch.Name", cr.ObjectMeta.Name, "Node.Type", "data")
	containerDefinition := statefulset.ElasticContainer(cr)

	if cr.Spec.Data.Resources != nil {
		containerDefinition.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Data.Resources.ResourceLimits.CPU)
		containerDefinition.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Data.Resources.ResourceRequests.CPU)
		containerDefinition.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Data.Resources.ResourceLimits.Memory)
		containerDefinition.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Data.Resources.ResourceRequests.Memory)
	}

	if cr.Spec.Data.Storage != nil {
		VolumeMounts := corev1.VolumeMount{Name: cr.ObjectMeta.Name + "-data", MountPath: "/usr/share/elasticsearch/data"}
		containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, VolumeMounts)
	}

	pluginVolume := corev1.VolumeMount{Name: "plugin-volume", MountPath: "/usr/share/elasticsearch/plugins"}
	containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, pluginVolume)

	if *cr.Spec.Security.TLSEnabled != false {
		containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, corev1.VolumeMount{Name: "tls-certificates", MountPath: "/usr/share/elasticsearch/config/certs"})
	}

	dataEnvVars := []corev1.EnvVar{
		corev1.EnvVar{Name: "discovery.seed_hosts", Value: cr.ObjectMeta.Name + "-master-headless"},
		corev1.EnvVar{Name: "network.host", Value: "0.0.0.0"},
		corev1.EnvVar{Name: "cluster.name", Value: cr.Spec.ClusterName},
		corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: "-Xmx" + cr.Spec.Data.JVMOptions.Max + " " + "-Xms" + cr.Spec.Data.JVMOptions.Min},
		corev1.EnvVar{Name: "node.data", Value: "true"},
		corev1.EnvVar{Name: "node.ingest", Value: "false"},
		corev1.EnvVar{Name: "node.master", Value: "false"},
		corev1.EnvVar{Name: "node.name", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
	}

	if *cr.Spec.Security.TLSEnabled != false {
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "https"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "ELASTIC_PASSWORD", Value: cr.Spec.Security.Password})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "ELASTIC_USERNAME", Value: "elastic"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.enabled", Value: "true"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.enabled", Value: "true"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.verification_mode", Value: "certificate"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.transport.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.enabled", Value: "true"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.truststore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "xpack.security.http.ssl.keystore.path", Value: "/usr/share/elasticsearch/config/certs/elastic-certificates.p12"})
	} else {
		dataEnvVars = append(dataEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "http"})
	}

	containerDefinition.Env = dataEnvVars

	reqLogger.Info("Successfully generated the contiainer definition for elasticsearch data")
	return containerDefinition
}

// ElasticSearchData creates the elasticsearch data statefulset
func ElasticSearchData(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-data",
		"role":                        "data",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	statefulsetObject := statefulset.StatefulSetObject(cr, "data", labels, cr.Spec.Data.Count)

	statefulsetObject.Spec.Template.Spec.Containers = append(statefulsetObject.Spec.Template.Spec.Containers, generateDataContainer(cr))

	statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.SysctlInitContainer(cr))

	if cr.Spec.Plugins != nil && len(cr.Spec.Plugins) != 0 {
		statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.PluginsInitContainer(cr))
	}

	if cr.Spec.Data.Storage != nil {
		statefulsetObject.Spec.VolumeClaimTemplates = append(statefulsetObject.Spec.VolumeClaimTemplates, statefulset.GeneratePVCTemplate(cr, "data", cr.Spec.Data.Storage))
	}

	if cr.Spec.Data.Affinity != nil {
		statefulsetObject.Spec.Template.Spec.Affinity = cr.Spec.Data.Affinity
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
	statefulset.SyncStatefulSet(cr, statefulsetObject, "data")
}
