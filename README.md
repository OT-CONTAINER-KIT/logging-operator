<p align="left">
  <img src="./static/logging-operator-logo.svg" height="180" width="180">
</p>

## Logging Operator

A golang based CRD operator to setup and manage logging (Elasticsearch, Fluentd and Kibana) in Kubernetes cluster. It helps to setup each component of EFK stack separately.

> The K8s API name is "logging.opstreelabs.in/v1alpha1"

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

