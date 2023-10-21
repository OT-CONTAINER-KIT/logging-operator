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

// ClusterRoleBindingParameters is an interface for clusterrolebindings input
type ClusterRoleBindingParameters struct {
	Name                   string
	OwnerDef               metav1.OwnerReference
	Namespace              string
	ClusterRoleBindingMeta metav1.ObjectMeta
}

// GenerateClusterRoleBinding is a method that will generate a clusterrolebindings interface
func GenerateClusterRoleBinding(params ClusterRoleBindingParameters) *rbacv1.ClusterRoleBinding {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		TypeMeta:   GenerateMetaInformation("ClusterRoleBinding", "rbac.authorization.k8s.io/v1"),
		ObjectMeta: params.ClusterRoleBindingMeta,
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     params.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      params.Name,
				Namespace: params.Namespace,
			},
		},
	}
	AddOwnerRefToObject(clusterRoleBinding, params.OwnerDef)
	return clusterRoleBinding
}

// CreateClusterRoleBinding is a method to create Kubernetes clusterrolebindings
func CreateClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	logger := LogGenerator(clusterRoleBinding.Name, "global", "ClusterRoleBinding")
	_, err := GenerateK8sClient().RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "ClusterRoleBinding creation is failed")
		return err
	}
	logger.Info("ClusterRoleBinding creation is successful")
	return nil
}

// GetClusterRoleBinding is a method to check clusterrolebindings exists
//
//nolint:gosimple
func GetClusterRoleBinding(name string) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBindingsInfo, err := GenerateK8sClient().RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return clusterRoleBindingsInfo, nil
}
