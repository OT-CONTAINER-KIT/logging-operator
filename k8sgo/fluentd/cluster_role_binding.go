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

package k8sfluentd

import (
	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
)

// CreateFluentdClusterRoleBinding is a method to create clusterrolebinding for fluentd
func CreateFluentdClusterRoleBinding(cr *loggingv1beta1.Fluentd) error {
	labels := map[string]string{
		"app": cr.ObjectMeta.Name,
	}
	clusterRoleBindingParams := k8sgo.ClusterRoleBindingParameters{
		Name:                   cr.ObjectMeta.Name,
		OwnerDef:               k8sgo.FluentdAsOwner(cr),
		Namespace:              cr.Namespace,
		ClusterRoleBindingMeta: k8sgo.GenerateObjectMetaInformation(cr.ObjectMeta.Name, cr.Namespace, labels, k8sgo.GenerateAnnotations()),
	}
	err := k8sgo.CreateClusterRoleBinding(k8sgo.GenerateClusterRoleBinding(clusterRoleBindingParams))
	if err != nil {
		return err
	}
	return nil
}
