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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"logging-operator/k8sutils/clusterrole"
	"logging-operator/k8sutils/clusterrolebindings"
	"logging-operator/k8sutils/configmap"
	"logging-operator/k8sutils/daemonset"
	"logging-operator/k8sutils/serviceaccount"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
)

// FluentdReconciler reconciles a Fluentd object
type FluentdReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.opstreelabs.in,resources=fluentds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.opstreelabs.in,resources=fluentds/status,verbs=get;update;patch

func (r *FluentdReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("fluentd", req.NamespacedName)

	instance := &loggingv1alpha1.Fluentd{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Second * 10}, nil
		}
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	if err := controllerutil.SetControllerReference(instance, instance, r.Scheme); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	configmap.CreateFluentdConfigMap(instance)
	serviceaccount.CreateFluentdServiceAccount(instance)
	clusterrole.CreateFluentdClusterRole(instance)
	clusterrolebindings.CreateFluentdClusterRoleBinding(instance)
	daemonset.CreateFluentdDaemonset(instance)

	reqLogger.Info("Will reconcile after 10 seconds", "Fluentd.Namespace", instance.Namespace, "Fluentd.Name", instance.Name)
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

func (r *FluentdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1alpha1.Fluentd{}).
		Complete(r)
}
