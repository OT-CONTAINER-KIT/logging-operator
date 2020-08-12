package statefulset

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	graceTime = 20
)

var log = logf.Log.WithName("elastic_statefulset")

// SyncStatefulSet will sync the statefulset in Kubernetes
func SyncStatefulSet(cr *loggingv1alpha1.Elasticsearch, statefulset *appsv1.StatefulSet, nodeType string) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Node.Type", nodeType,
	)

	elasticStatefuleName := cr.ObjectMeta.Name + "-" + nodeType
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for statefulset")
	}

	statefulObject, err := k8sClient.AppsV1().StatefulSets(cr.Namespace).Get(context.TODO(), elasticStatefuleName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating elasticsearch setup", "Name", cr.ObjectMeta.Name)
		k8sClient.AppsV1().StatefulSets(cr.Namespace).Create(context.TODO(), statefulset, metav1.CreateOptions{})
	} else if statefulObject != statefulset {
		reqLogger.Info("Updating elasticsearch setup", "Name", cr.ObjectMeta.Name)
		k8sClient.AppsV1().StatefulSets(cr.Namespace).Update(context.TODO(), statefulset, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Elasticsearch nodes are in sync")
	}
}

// GeneratePVCTemplate is a method to generate the PVC template
func GeneratePVCTemplate(cr *loggingv1alpha1.Elasticsearch, nodeType string, storageSpec *loggingv1alpha1.Storage) corev1.PersistentVolumeClaim {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Node.Type", nodeType,
	)

	pvcTemplate := storageSpec.VolumeClaimTemplate
	elasticStatefuleName := cr.ObjectMeta.Name + "-" + nodeType

	if storageSpec == nil {
		reqLogger.Info("No storage is defined for elasticsearch", "Elasticsearch.Name", cr.ObjectMeta.Name)
	} else {
		pvcTemplate.CreationTimestamp = metav1.Time{}
		pvcTemplate.Name = elasticStatefuleName
		if storageSpec.VolumeClaimTemplate.Spec.AccessModes == nil {
			pvcTemplate.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		} else {
			pvcTemplate.Spec.AccessModes = storageSpec.VolumeClaimTemplate.Spec.AccessModes
		}
		pvcTemplate.Spec.Resources = storageSpec.VolumeClaimTemplate.Spec.Resources
		pvcTemplate.Spec.Selector = storageSpec.VolumeClaimTemplate.Spec.Selector
		pvcTemplate.Spec.Selector = storageSpec.VolumeClaimTemplate.Spec.Selector
	}
	reqLogger.Info("Successfully generated the PVC template")
	return pvcTemplate
}

func GenerateElasticContainer(cr *loggingv1alpha1.Elasticsearch) corev1.Container {
	var containerDefinition corev1.Container

	containerDefinition = corev1.Container{
		Name:            "elastic",
		Image:           cr.Spec.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Env:             []corev1.EnvVar{},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{}, Requests: corev1.ResourceList{},
		},
		VolumeMounts: []corev1.VolumeMount{},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: graceTime,
			PeriodSeconds:       15,
			FailureThreshold:    5,
			TimeoutSeconds:      5,
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt(9200),
				},
			},
		},
		LivenessProbe: &corev1.Probe{
			InitialDelaySeconds: graceTime,
			TimeoutSeconds:      5,
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt(9200),
				},
			},
		},
	}
	return containerDefinition
}
