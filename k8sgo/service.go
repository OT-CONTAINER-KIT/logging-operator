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

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServiceParameters is a structure for service inputs
type ServiceParameters struct {
	ServiceMeta              metav1.ObjectMeta
	OwnerDef                 metav1.OwnerReference
	Labels                   map[string]string
	Annotations              map[string]string
	Namespace                string
	HeadlessService          bool
	Port                     []PortInfo
	PublishNotReadyAddresses bool
}

// PortInfo is a structure for port information
type PortInfo struct {
	PortName string
	Port     int32
}

// CreateOrUpdateService method will create or update service
func CreateOrUpdateService(params ServiceParameters) error {
	logger := LogGenerator(params.ServiceMeta.Name, params.Namespace, "Service")
	serviceDef := generateServiceDef(params)
	storedService, err := getService(params.Namespace, params.ServiceMeta.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(serviceDef); err != nil {
				logger.Error(err, "Unable to patch service with compare annotations")
			}
			return createService(params.Namespace, serviceDef)
		}
		return err
	}
	return patchService(storedService, serviceDef, params.Namespace)
}

// patchService will patch Kubernetes service
func patchService(storedService *corev1.Service, newService *corev1.Service, namespace string) error {
	logger := LogGenerator(storedService.Name, namespace, "Service")
	// adding meta fields
	newService.ResourceVersion = storedService.ResourceVersion
	newService.CreationTimestamp = storedService.CreationTimestamp
	newService.ManagedFields = storedService.ManagedFields
	newService.Spec.ClusterIP = storedService.Spec.ClusterIP

	patchResult, err := patch.DefaultPatchMaker.Calculate(storedService, newService,
		patch.IgnoreStatusFields(),
		patch.IgnoreField("kind"),
		patch.IgnoreField("apiVersion"),
	)
	if err != nil {
		logger.Error(err, "Unable to patch service with comparison object")
		return err
	}
	if !patchResult.IsEmpty() {
		for key, value := range storedService.Annotations {
			if _, present := newService.Annotations[key]; !present {
				newService.Annotations[key] = value
			}
		}
		if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(newService); err != nil {
			logger.Error(err, "Unable to patch service with comparison object")
			return err
		}
		logger.Info("Syncing service with defined properties")
		return updateService(namespace, newService)
	}
	logger.Info("Service is already in-sync")
	return nil
}

// createService is a method to create service
func createService(namespace string, service *corev1.Service) error {
	logger := LogGenerator(service.Name, namespace, "Service")
	_, err := GenerateK8sClient().CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(err, "Service creation is failed")
		return err
	}
	logger.Info("Service creation is successful")
	return nil
}

// updateService is a method to update service
func updateService(namespace string, service *corev1.Service) error {
	logger := LogGenerator(service.Name, namespace, "Service")
	_, err := GenerateK8sClient().CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(err, "Service updation is failed")
		return err
	}
	logger.Info("Service updation is successful")
	return nil
}

// getService is a method to get service
func getService(namespace string, service string) (*corev1.Service, error) {
	logger := LogGenerator(service, namespace, "Service")
	serviceInfo, err := GenerateK8sClient().CoreV1().Services(namespace).Get(context.TODO(), service, metav1.GetOptions{})
	if err != nil {
		logger.Info("Service get action is failed")
		return nil, err
	}
	logger.Info("Service get action is successful")
	return serviceInfo, nil
}

// generateServiceDef is a method to generate service definition
func generateServiceDef(params ServiceParameters) *corev1.Service {
	service := &corev1.Service{
		TypeMeta:   GenerateMetaInformation("Service", "core/v1"),
		ObjectMeta: params.ServiceMeta,
		Spec: corev1.ServiceSpec{
			Selector:                 params.Labels,
			Ports:                    []corev1.ServicePort{},
			PublishNotReadyAddresses: params.PublishNotReadyAddresses,
		},
	}

	for _, portInfo := range params.Port {
		service.Spec.Ports = append(service.Spec.Ports, corev1.ServicePort{
			Name:       portInfo.PortName,
			Port:       portInfo.Port,
			TargetPort: intstr.FromInt(int(portInfo.Port)),
			Protocol:   corev1.ProtocolTCP,
		})
	}

	if params.HeadlessService {
		service.Spec.ClusterIP = "None"
	}
	AddOwnerRefToObject(service, params.OwnerDef)
	return service
}
