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

package clusterrole

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
	"logging-operator/k8sutils/client"
	"logging-operator/k8sutils/identifier"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("fluentd_cluster_role")

func generateClusterRole(cr *loggingv1alpha1.Fluentd, labels map[string]string) *rbacv1.ClusterRole {

	clusterRole := &rbacv1.ClusterRole{
		TypeMeta:   identifier.GenerateMetaInformation("ClusterRole", "rbac.authorization.k8s.io/v1beta1"),
		ObjectMeta: identifier.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, identifier.GenerateFluentdAnnotations()),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"pods", "namespaces"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
	identifier.AddOwnerRefToObject(clusterRole, identifier.FluentdAsOwner(cr))

	return clusterRole
}

// SyncFluentdClusterRole will create and update the clusterrole
func SyncFluentdClusterRole(cr *loggingv1alpha1.Fluentd, role *rbacv1.ClusterRole) {
	reqLogger := log.WithValues(
		"Request.Namespace", cr.Namespace,
		"Request.Name", cr.ObjectMeta.Name,
		"Service.Type", "ClusteRole",
	)

	clusterRoleName := cr.ObjectMeta.Name
	k8sClient, err := client.GenerateK8sClient()
	if err != nil {
		reqLogger.Error(err, "Unable to generate K8s client for cluster role")
	}

	clusterRoleObject, err := k8sClient.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metav1.GetOptions{})

	if err != nil {
		reqLogger.Info("Creating cluster role for fluentd", "Name", cr.ObjectMeta.Name)
		_, err := k8sClient.RbacV1().ClusterRoles().Create(context.TODO(), role, metav1.CreateOptions{})
		if err != nil {
			reqLogger.Error(err, "Got an error please check")
		}
	} else if clusterRoleObject != role {
		reqLogger.Info("Updating cluster role for fluentd", "Name", cr.ObjectMeta.Name)
		k8sClient.RbacV1().ClusterRoles().Update(context.TODO(), role, metav1.UpdateOptions{})
	} else {
		reqLogger.Info("Fluentd cluster role are in sync")
	}
}

// CreateFluentdClusterRole creates the clusterrole for fluentd configuration
func CreateFluentdClusterRole(cr *loggingv1alpha1.Fluentd) {
	labels := map[string]string{
		"app":                         cr.ObjectMeta.Name,
		"logging.opstreelabs.in":      "true",
		"logging.opstreelabs.in/kind": "Fluentd",
	}

	config := generateClusterRole(cr, labels)
	SyncFluentdClusterRole(cr, config)
}
