<p align="left">
  <img src="./static/logging-operator-logo.svg" height="180" width="180">
</p>

## Logging Operator

A golang based CRD operator to setup and manage logging (Elasticsearch, Fluentd and Kibana) in Kubernetes cluster. It helps to setup each component of EFK stack separately.

> The K8s API name is "logging.opstreelabs.in/v1alpha1"

Our roadmap is present in [ROADMAP](ROADMAP.md)

### Supported Features

The "Logging Operator" includes these features:-

- Elasticsearch different node types, like:-
  - Master Node
  - Data Node
  - Ingestion Node
  - Client/Cordinator Node
- Elasticsearch setup with/without TLS
- Customizable elasticsearch configuration and Heap size
- Fluentd as a log-shipper which already has JSON logs support
- Kibana integration with elasticsearch for logs visulization
- Seamless upgrade for Elasticsearch, Fluentd, and Kibana
- Inculcated best practices for Kubernetes setup like `SecurityContext` and `Privilege Control`
- Loosely coupled setup, i.e. Elasticsearch, Fluentd, and Kibana can be setup individually as well.

### Architecture

<div align="center">
    <img src="./static/logging-operator-arch.png">
</div>

### Purpose

The purpose behind creating this CRD operator was to provide an easy and yet production grade logging setup on Kubernetes. But it doesn't mean this can only be used for logging setup only.

> This operator blocks Elasticsearch, Fluentd, and Kibana are loosely-couples so they can be setup individually as well. For example:- If we need elasticsearch for application database we can setup only elasticsearch as well by using this operator.

### Prerequisites

The "Logging Operator" needs a Kubernetes/Openshift cluster of version `>=1.8.0`. If you have just started using Operatorss, its highly recommend to use latest version of Kubernetes.

The cluster size selection should be done on the basis of requirement and resources.

### Logging Operator Installation

For "Logging Operator" installation, we have categorized the steps in 3 parts:-

- Namespace Setup for operator
- CRD setup in kubernetes cluster
- RBAC setup for operator to create resources in Kubernetes
- Operator deployment and validation

#### Namespace setup

Since we are going to use pre-baked manifests of Kubernetes in that case we need to setup the namespace with a specific name called "logging-operator".

```shell
kubectl create ns logging-operator
```

#### CRD Setup

So we have already pre-configured CRD in [config/crd](./config/crd) directory. We just have to run a magical `kubectl` commands.

```shell
kubectl apply -f config/crd/
```

#### RBAC setup

Similar like CRD, we have pre-baked RBAC config files as well inside [config/crd](./config/rbac) which can be installed and configured by `kubectl`

```shell
kubectl apply -f config/rbac/
```

#### Operator Deployment and Validation

Once all the initial steps are done, we can create the deployment for "Logging Operator". The deployment manifests for operator is present inside [config/manager/manager.yaml](./config/manager/manager.yaml) file.

```shell
kubectl apply -f config/manager/manager.yaml
```

### Examples

All the examples are present inside the [config/samples/](./config/samples/) directory. These manifests can be applied by `kubectl` command line. These configuration have some dummy values which can be changed and customized by the individuals as per needs and requirements.

### Contact Information

This project is managed by [OpsTree Solutions](https://opstree.com). If you have any queries or suggestions, mail us at opensource@opstree.com
