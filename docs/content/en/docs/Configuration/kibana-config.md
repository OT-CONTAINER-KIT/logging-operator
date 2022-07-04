---
title: "Kibana Config"
weight: 4
linkTitle: "Kibana Config"
description: >
    Kibana configuration paramaters with respect to logging operator
---

Kibana configuration is easily customizable using `helm` as well `kubectl`. Since all the configurations are in the form YAML file, it can be easily changed and customized.

The values.yaml file for Kibana setup can be found [here](https://github.com/OT-CONTAINER-KIT/helm-charts/tree/main/charts/kibana). But if the setup is not done using Helm, in that case Kubernetes manifests needs to be customized.

## Helm chart parameters

| **Name**                          | **Value**                         | **Description**                                            |
|-----------------------------------|-----------------------------------|------------------------------------------------------------|
| replicas                          | 1                                 | Number of deployment replicas for kibana                   |
| esCluster.esURL                   | https://elasticsearch-master:9200 | Hostname or URL of the elasticsearch server                |
| esCluster.esVersion               | 7.17.0                            | Version of the kibana in pair with elasticsearch           |
| esCluster.clusterName             | elasticsearch                     | Name of the elasticsearch created by elasticsearch crd     |
| resources                         | {}                                | Resources for kibana visualization pods                    |
| nodeSelectors                     | {}                                | Nodeselectors map key-values for kibana visualization pods |
| affinity                          | {}                                | Affinity and anit-affinity for kibana visualization pods   |
| tolerations                       | {}                                | Tolerations and taints for kibana visualization pods       |
| esSecurity.enabled                | true                              | To enabled the xpack security of kibana                    |
| esSecurity.elasticSearchPassword  | elasticsearch-password            | Credentials for elasticsearch authentication               |
| externalService.enabled           | false                             | To create a LoadBalancer service of kibana                 |
| ingress.enabled                   | false                             | To enable the ingress resource for kibana                  |
| ingress.host                      | kibana.opstree.com                | Hostname or URL on which kibana will be exposed            |
| ingress.tls.enabled               | false                             | To enable SSL on kibana ingress resource                   |
| ingress.tls.secret                | tls-secret                        | SSL certificate for kibana ingress resource                |

## CRD object definition parameters

These are the parameters that are currently supported by the Logging Operator for the Kibana setup:-

- replicas
- esCluster
- esSecurity
- kubernetesConfig

### replicas

`replicas` is field definition of Kibana CRD in which we can define how many replicas/instances of Kibana we would like to run. Similar field like replicas in deployment and replicasets.

```yaml
  replicas: 1
```

### esCluster

`esCluster` is a general parameter of Fluentd CRD for providing the information about Elasticsearch nodes.

```yaml
  esCluster:
    host: https://elasticsearch-master:9200
    esVersion: 7.16.0
    clusterName: elasticsearch
```

### esSecurity

`esSecurity` s the security specification for Fluentd CRD. If we want to enable authentication and TLS, in that case, we can enable this configuration. To enable the authentication we need to provide secret reference in Kubernetes.

```yaml
  esSecurity:
    tlsEnabled: true
    existingSecret: elasticsearch-password
```

### kubernetesConfig

`kubernetesConfig` is the general configuration paramater for Fluentd CRD in which we are defining the Kubernetes related configuration details like- image, tag, imagePullPolicy, and resources.

```yaml
  kubernetesConfig:
    resources:
      requests:
        cpu: 100m
        memory: 100Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```
