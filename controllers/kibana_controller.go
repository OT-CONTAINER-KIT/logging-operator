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
	"logging-operator/k8sgo/kibana"
)

// KibanaReconciler reconciles a Kibana object
type KibanaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=kibanas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=kibanas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=kibanas/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *KibanaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &loggingv1beta1.Kibana{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Second * 10}, nil
		}
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	err = k8skibana.CreateKibanaSetup(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	err = k8skibana.CreateKibanaService(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KibanaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Kibana{}).
		Complete(r)
}
