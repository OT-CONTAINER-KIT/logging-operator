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

package k8sgo

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRoleParameters is an interface for clusterroles input
type ClusterRoleParameters struct {
	Name            string
	OwnerDef        metav1.OwnerReference
	Namespace       string
	ClusterRoleMeta metav1.ObjectMeta
	Rules           []rbacv1.PolicyRule
}

// GenerateClusterRoles is a method that will generate a clusterroles interface
func GenerateClusterRoles(params ClusterRoleParameters) *rbacv1.ClusterRole {
	clusterRole := &rbacv1.ClusterRole{
		TypeMeta:   GenerateMetaInformation("ClusterRole", "rbac.authorization.k8s.io/v1"),
		ObjectMeta: params.ClusterRoleMeta,
		Rules:      params.Rules,
	}
	AddOwnerRefToObject(clusterRole, params.OwnerDef)
	return clusterRole
}

// CreateClusterRole is a method to create Kubernetes clusterroles
func CreateClusterRole(clusterRole *rbacv1.ClusterRole) error {
	logger := LogGenerator(clusterRole.Name, "global", "ClusterRole")
	_, err := GenerateK8sClient().RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "ClusterRole creation is failed")
		return err
	}
	logger.Info("ClusterRole creation is successful")
	return nil
}

// GetClusterRole is a method to check clusterrole exists
//
//nolint:gosimple
func GetClusterRole(name string) (*rbacv1.ClusterRole, error) {
	clusterRoleInfo, err := GenerateK8sClient().RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clusterRoleInfo, nil
}
