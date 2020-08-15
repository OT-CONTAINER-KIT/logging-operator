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

package clusterrolebindings

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("fluentd_cluster_rolebindings")

func generateClusterRoleBindings(cr *loggingv1alpha1.Fluentd, labels map[string]string) *rbacv1.ClusterRoleBinding {

	clusterRoleBindings := &rbacv1.ClusterRoleBinding{
		TypeMeta:   identifier.GenerateMetaInformation("ClusterRoleBinding", "rbac.authorization.k8s.io/v1beta1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     cr.ObjectMeta.Name,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      cr.ObjectMeta.Name,
				Namespace: cr.Namespace,
			},
		},
	}
	identifier.AddOwnerRefToObject(clusterRoleBindings, identifier.FluentdAsOwner(cr))
	return clusterRoleBindings
}

// SyncFluentdClusterRoleBinding will create and update the rolebinding
func SyncFluentdClusterRoleBinding(cr *loggingv1alpha1.Fluentd, rolebinding *rbacv1.ClusterRoleBinding) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "ClusteRoleBinding",
	)

	clusterRoleBindingName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for cluster rolebinding")
	}

	clusterRoleBindingObject, err := k8sClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating cluster role binding for fluentd", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.RbacV1().ClusterRoleBindings().Create(context.TODO(), rolebinding, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if clusterRoleBindingObject != rolebinding {
		reqLogger.Info("Updating cluster role binding for fluentd", "Name", cr.ObjectMeta.Name)
		k8sClient.RbacV1().ClusterRoleBindings().Update(context.TODO(), rolebinding, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Fluentd cluster role binding are in sync")
	}
}

// CreateFluentdClusterRoleBinding creates the clusterrolebindings for fluentd configuration
func CreateFluentdClusterRoleBinding(cr *loggingv1alpha1.Fluentd) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}

	config := generateClusterRoleBindings(cr, labels)
	SyncFluentdClusterRoleBinding(cr, config)
}
