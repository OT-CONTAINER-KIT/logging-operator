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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	// 	"sigs.k8s.io/controller-runtime/pkg/log"

	loggingv1beta1 "logging-operator/api/v1beta1"
	"logging-operator/elasticgo"
	"logging-operator/k8sgo"
	"logging-operator/k8sgo/elasticsearch"
)

// ElasticsearchReconciler reconciles a Elasticsearch object
type ElasticsearchReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=elasticsearches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=elasticsearches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=logging.logging.opstreelabs.in,resources=elasticsearches/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps;events;services;secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ElasticsearchReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &loggingv1beta1.Elasticsearch{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Second * 10}, nil
		}
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	err = secretManager(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	err = k8selastic.CreateElasticSearchService(instance, "master")
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	err = k8selastic.SetupElasticSearchMaster(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	if instance.Spec.ESData != nil {
		err = k8selastic.SetupElasticSearchData(instance)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
		err = k8selastic.CreateElasticSearchService(instance, "data")
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}

	if instance.Spec.ESIngestion != nil {
		err = k8selastic.SetupElasticSearchIngestion(instance)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
		err = k8selastic.CreateElasticSearchService(instance, "ingestion")
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}

	if instance.Spec.ESClient != nil {
		err = k8selastic.SetupElasticSearchClient(instance)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
		err = k8selastic.CreateElasticSearchService(instance, "client")
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}

	if err := controllerutil.SetControllerReference(instance, instance, r.Scheme); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	instance.Status.ESVersion = instance.Spec.ESVersion
	instance.Status.ESMaster = instance.Spec.ESMaster.Replicas
	if instance.Spec.ESData != nil {
		instance.Status.ESData = instance.Spec.ESData.Replicas
	}
	if instance.Spec.ESIngestion != nil {
		instance.Status.ESIngestion = instance.Spec.ESIngestion.Replicas
	}
	if instance.Spec.ESClient != nil {
		instance.Status.ESClient = instance.Spec.ESClient.Replicas
	}

	clusterInfo, err := elasticgo.GetElasticClusterDetails(instance)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 10}, err
	}

	instance.Status.ClusterState = clusterInfo.ClusterState
	instance.Status.ActiveShards = &clusterInfo.Shards
	instance.Status.Indices = &clusterInfo.Shards

	if clusterInfo.ClusterState == "green" {
		err = serviceAccountSecretManager(instance)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 10}, err
		}
	}
	if err := r.Status().Update(context.TODO(), instance); err != nil {
		if errors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

// secretManager is a method to create and manage secrets
func secretManager(instance *loggingv1beta1.Elasticsearch) error {
	if instance.Spec.Security != nil {
		if instance.Spec.Security.AutoGeneratePassword != nil && *instance.Spec.Security.AutoGeneratePassword {
			secretName := fmt.Sprintf("%s-%s", instance.ObjectMeta.Name, "password")
			_, err := k8sgo.GetSecret(secretName, instance.Namespace)

			if err != nil {
				err = k8selastic.CreateElasticAutoSecret(instance)
				if err != nil {
					return err
				}
			}
		}
	}

	if instance.Spec.Security != nil {
		if instance.Spec.Security.TLSEnabled != nil && *instance.Spec.Security.TLSEnabled {
			tlsSecretName := fmt.Sprintf("%s-%s", instance.ObjectMeta.Name, "tls-cert")
			_, err := k8sgo.GetSecret(tlsSecretName, instance.Namespace)

			if err != nil {
				err = k8selastic.CreateElasticTLSSecret(instance)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// serviceAccountSecretManager is a method for service account
func serviceAccountSecretManager(instance *loggingv1beta1.Elasticsearch) error {
	if instance.Spec.Security != nil {
		tokenSecretName := fmt.Sprintf("%s-sa-token", instance.ObjectMeta.Name)
		_, err := k8sgo.GetSecret(tokenSecretName, instance.Namespace)
		if err != nil {
			err = k8selastic.CreateServiceAccountToken(instance)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ElasticsearchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1beta1.Elasticsearch{}).
		Complete(r)
}
