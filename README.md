<p align="center">
  <img src="./static/logging-operator-logo2.svg" height="220" width="220">
</p>

<p align="center">
  <a href="https://circleci.com/gh/OT-CONTAINER-KIT/logging-operator/tree/master">
    <img src="https://circleci.com/gh/OT-CONTAINER-KIT/logging-operator/tree/master.svg?style=shield" alt="CircleCI">
  </a>

  <a href="https://goreportcard.com/report/github.com/OT-CONTAINER-KIT/logging-operator">
    <img src="https://goreportcard.com/badge/github.com/OT-CONTAINER-KIT/logging-operator" alt="Go Report Card">
  </a>

  <a href="https://quay.io/repository/opstree/logging-operator">
    <img src="https://img.shields.io/badge/container-ready-green" alt="Docker Repository on Quay">
  </a>

  <a href="https://github.com/OT-CONTAINER-KIT/logging-operator/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="Apache License">
  </a>

  <a href="https://codeclimate.com/github/OT-CONTAINER-KIT/logging-operator/maintainability">
    <img src="https://api.codeclimate.com/v1/badges/f9e99ffcba997de51eaa/maintainability" alt="Maintainability">
  </a>

  <a href="https://github.com/OT-CONTAINER-KIT/logging-operator/releases">
    <img src="https://img.shields.io/github/v/release/OT-CONTAINER-KIT/logging-operator" alt="GitHub release (latest by date)">
  </a>
</p>

## Logging Operator

A golang based CRD operator to setup and manage logging stack (Elasticsearch, Fluentd, and Kibana) in the Kubernetes cluster. It helps to setup each component of the EFK stack separately.

> The K8s API name is "logging.opstreelabs.in/v1alpha1"

### Documentation

[Documentation](https://docs.opstreelabs.in/logging-operator)

### Supported Features

The "Logging Operator" includes these features:-

- Elasticsearch different node types, like:-
  - **Master Node** => A node that has the master role (default), which makes it eligible to be elected as the master node, which controls the cluster.
  - **Data Node** => A node that has the data role (default). Data nodes hold data and perform data related operations such as CRUD, search, and aggregations.
  - **Ingestion Node** => A node that has the ingest role (default). Ingest nodes are able to apply an ingest pipeline to a document in order to transform and enrich the document before indexing. With a heavy ingest load, it makes sense to use dedicated ingest nodes and to not include the ingest role from nodes that have the master or data roles.
  - **Client or Coordinator Node** => Requests like search requests or bulk-indexing requests may involve data held on different data nodes. A search request, for example, is executed in two phases which are coordinated by the node which receives the client request — the coordinating node.
- Elasticsearch setup with or without TLS on Transport and HTTP Layer
- Customizable elasticsearch configuration and configurable heap size
- Fluentd as a lightweight log-shipper and JSON field seperation support
- Kibana integration with elasticsearch for logs visualization
- Seamless upgrade for Elasticsearch, Fluentd, and Kibana stack
- Inculcated best practices for Kubernetes setup like `SecurityContext` and `Privilege Control`
- Loosely coupled setup, i.e. Elasticsearch, Fluentd, and Kibana setup can be done individually as well
- Index Lifecycle support to manage rollover and cleanup of indexes
- Index template support for configuring index settings like:- policy, replicas, shards etc.

### Architecture

<div align="center">
    <img src="./static/logging-operator-arch.png">
</div>

### Purpose

The purpose behind creating this CRD operator was to provide an easy and yet production grade logging setup on Kubernetes. But it doesn't mean this can only be used for logging setup only.

> This operator blocks Elasticsearch, Fluentd, and Kibana are loosely-couples so they can be setup individually as well. For example:- If we need elasticsearch for application database we can setup only elasticsearch as well by using this operator.

### Prerequisites

The "Logging Operator" needs a Kubernetes/Openshift cluster of version `>=1.8.0`. If you have just started using Operatorss, it's highly recommended to use the latest version of Kubernetes.

The cluster size selection should be done on the basis of requirements and resources.

### Logging Operator Installation

For the "Logging Operator" installation, we have categorized the steps in 3 parts:-

- Namespace Setup for operator
- CRD setup in Kubernetes cluster
- RBAC setup for an operator to create resources in Kubernetes
- Operator deployment and validation

### Examples

All the examples are present inside the [config/samples/](./config/samples/) directory. These manifests can be applied by the `kubectl` command line. These configurations have some dummy values which can be changed and customized by the individuals as per needs and requirements.

### Contact Information

This project is managed by [OpsTree Solutions](https://opstree.com). If you have any queries or suggestions, mail us at opensource@opstree.com
