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

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/k8sgo"
	"logging-operator/k8sgo/fluentd"
)

// FluentdReconciler reconciles a Fluentd object
type FluentdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=fluentds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=fluentds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=fluentds/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=serviceaccounts;pods;namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FluentdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &loggingv1beta1.Fluentd{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Second * 10}, nil
		}
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	err = setupFluentdRBAC(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	err = k8sfluentd.CreateFluentdConfigMap(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	err = k8sfluentd.CreateFluentdDaemonSet(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

// setupFluentdRBAC is a method to setup RBAC access for Fluentd
func setupFluentdRBAC(instance *loggingv1beta1.Fluentd) error {
	_, err := k8sgo.GetServiceAccount(instance.ObjectMeta.Name, instance.Namespace)
	if err != nil {
		err = k8sfluentd.CreateFluentdServiceAccount(instance)
		if err != nil {
			return err
		}
	}
	_, err = k8sgo.GetClusterRole(instance.ObjectMeta.Name)
	if err != nil {
		err = k8sfluentd.CreateFluentdClusterRole(instance)
		if err != nil {
			return err
		}
	}
	_, err = k8sgo.GetClusterRoleBinding(instance.ObjectMeta.Name)
	if err != nil {
		err = k8sfluentd.CreateFluentdClusterRoleBinding(instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluentdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Fluentd{}).
		Complete(r)
}
