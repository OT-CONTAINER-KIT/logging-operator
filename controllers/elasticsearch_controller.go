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
	"logging-operator/k8sutils/secret"
	clientservice "logging-operator/k8sutils/service/client"
	dataservice "logging-operator/k8sutils/service/data"
	ingestionservice "logging-operator/k8sutils/service/ingestion"
	masterservice "logging-operator/k8sutils/service/master"
	clientnode "logging-operator/k8sutils/statefulset/client"
	"logging-operator/k8sutils/statefulset/data"
	"logging-operator/k8sutils/statefulset/ingestion"
	"logging-operator/k8sutils/statefulset/master"
	elasticutils "logging-operator/utils/elasticsearch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	loggingv1alpha1 "logging-operator/api/v1alpha1"
)

// ElasticsearchReconciler reconciles a Elasticsearch object
type ElasticsearchReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.opstreelabs.in,resources=elasticsearches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.opstreelabs.in,resources=elasticsearches/status,verbs=get;update;patch

func (r *ElasticsearchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// _ = context.Background()
	reqLogger := r.Log.WithValues("elasticsearch", req.NamespacedName)

	instance := &loggingv1alpha1.Elasticsearch{}
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

	if *instance.Spec.Security.TLSEnabled != false && instance.Spec.Security.TLSEnabled != nil {
		tlsCert := secret.GenerateTLSecret(instance)
		secret.CreateAndUpdateSecret(instance, tlsCert)

		password := secret.GenerateElasticPassword(instance)
		secret.CreateAndUpdateSecret(instance, password)
	}

	if instance.Spec.Master.Enabled != false {
		master.ElasticSearchMaster(instance)
		masterservice.MasterElasticSearchService(instance)
		instance.Status.Master = instance.Spec.Master.Count
	}

	if instance.Spec.Data.Enabled != false {
		data.ElasticSearchData(instance)
		dataservice.DataElasticSearchService(instance)
		instance.Status.Data = instance.Spec.Data.Count
	}

	if instance.Spec.Ingestion.Enabled != false {
		ingestion.ElasticSearchIngestion(instance)
		ingestionservice.IngestionElasticSearchService(instance)
		instance.Status.Ingestion = instance.Spec.Ingestion.Count
	}

	if instance.Spec.Client.Enabled != false {
		clientnode.ElasticSearchClient(instance)
		clientservice.ClientElasticSearchService(instance)
		instance.Status.Client = instance.Spec.Client.Count
	}
	instance.Status.ClusterName = instance.Spec.ClusterName

	clusterStatus, err := elasticutils.GetElasticHealth(instance)

	if err != nil {
		instance.Status.ClusterState = "Not Ready"
	} else {
		instance.Status.ClusterState = *clusterStatus
	}

	if err := r.Status().Update(context.TODO(), instance); err != nil {
		if errors.IsConflict(err) {
			reqLogger.Error(err, "Conflict updating Elasticsearch status, requeueing")
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	reqLogger.Info("Will reconcile after 10 seconds", "Elasticsearch.Namespace", instance.Namespace, "Elasticsearch.Name", instance.Name)
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

func (r *ElasticsearchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1alpha1.Elasticsearch{}).
		Complete(r)
}
