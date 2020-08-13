package master

import (
	// "context"
	"strconv"
	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	// "logging-operator/k8sutils/client"
	"logging-operator/k8sutils/statefulset"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elastic_master")

func generateMasterContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {

	reqLogger := log.WithValues("Namespace", cr.Namespace, "Elasticsearch.Name", cr.ObjectMeta.Name, "Node.Type", "master")
	var nodes []string
	containerDefinition := statefulset.ElasticContainer(cr)

	if cr.Spec.Master.Resources != nil {
		containerDefinition.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Master.Resources.ResourceLimits.CPU)
		containerDefinition.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(cr.Spec.Master.Resources.ResourceRequests.CPU)
		containerDefinition.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Master.Resources.ResourceLimits.Memory)
		containerDefinition.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(cr.Spec.Master.Resources.ResourceRequests.Memory)
	}

	if cr.Spec.Master.Storage != nil {
		VolumeMounts := corev1.VolumeMount{
			Name:      cr.ObjectMeta.Name + "-master",
			MountPath: "/usr/share/elasticsearch/data",
		}
		containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, VolumeMounts)
	}

	pluginVolume := corev1.VolumeMount{
		Name:      "plugin-volume",
		MountPath: "/usr/share/elasticsearch/plugins",
	}
	containerDefinition.VolumeMounts = append(containerDefinition.VolumeMounts, pluginVolume)

	for count := 1; count <= int(*cr.Spec.Master.Count); count++ {
		nodes = append(nodes, cr.ObjectMeta.Name+"-master-"+strconv.Itoa(count)+",")
	}

	masterEnvVars := []corev1.EnvVar{
		corev1.EnvVar{Name: "cluster.initial_master_nodes", Value: strings.Join(nodes, "")},
		corev1.EnvVar{Name: "discovery.seed_hosts", Value: cr.ObjectMeta.Name + "-master-headless"},
		corev1.EnvVar{Name: "network.host", Value: "0.0.0.0"},
		corev1.EnvVar{Name: "cluster.name", Value: cr.Spec.ClusterName},
		corev1.EnvVar{Name: "ES_JAVA_OPTS", Value: "-Xmx" + cr.Spec.Master.JVMOptions.Max + " " + "-Xms" + cr.Spec.Master.JVMOptions.Min},
		corev1.EnvVar{Name: "node.data", Value: "false"},
		corev1.EnvVar{Name: "node.ingest", Value: "false"},
		corev1.EnvVar{Name: "node.master", Value: "true"},
		corev1.EnvVar{Name: "node.name", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
	}

	if *cr.Spec.Security.TLSEnabled != false {
		masterEnvVars = append(masterEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "https"})
	} else {
		masterEnvVars = append(masterEnvVars, corev1.EnvVar{Name: "SCHEME", Value: "http"})
	}

	containerDefinition.Env = masterEnvVars

	reqLogger.Info("Successfully generated the contiainer definition for elasticsearch master")
	return containerDefinition
}

// ElasticSearchMaster creates the elasticsearch master statefulset
func ElasticSearchMaster(cr *loggingv1alpha1.Elasticsearch) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name + "-master",
		"role":                        "master",
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Elasticsearch",
	}

	statefulsetObject := statefulset.StatefulSetObject(cr, "master", labels, cr.Spec.Master.Count)

	statefulsetObject.Spec.Template.Spec.Containers = append(statefulsetObject.Spec.Template.Spec.Containers, generateMasterContainer(cr))

	statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.SysctlInitContainer(cr))

	if cr.Spec.Plugins != nil && len(cr.Spec.Plugins) != 0 {
		statefulsetObject.Spec.Template.Spec.InitContainers = append(statefulsetObject.Spec.Template.Spec.InitContainers, statefulset.PluginsInitContainer(cr))
	}

	if cr.Spec.Master.Storage != nil {
		statefulsetObject.Spec.VolumeClaimTemplates = append(statefulsetObject.Spec.VolumeClaimTemplates, statefulset.GeneratePVCTemplate(cr, "master", cr.Spec.Master.Storage))
	}

	if cr.Spec.Master.Affinity != nil {
		statefulsetObject.Spec.Template.Spec.Affinity = cr.Spec.Master.Affinity
	}
	statefulset.SyncStatefulSet(cr, statefulsetObject, "master")
}
