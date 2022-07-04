---
title: "Fluentd Config"
weight: 3
linkTitle: "Fluentd Config"
description: >
    Elasticsearch configuration paramaters with respect to logging operator
---

Fluentd configuration is easily customizable using `helm` as well `kubectl`. Since all the configurations are in the form YAML file, it can be easily changed and customized.

The values.yaml file for Fluentd setup can be found [here](https://github.com/OT-CONTAINER-KIT/helm-charts/tree/main/charts/fluentd). But if the setup is not done using Helm, in that case Kubernetes manifests needs to be customized.

## Helm chart parameters

| **Name**                         | **Values**             | **Description**                                                 |
|----------------------------------|------------------------|-----------------------------------------------------------------|
| elasticSearchHost                | elasticsearch-master   | Hostname or URL of the elasticsearch server                     |
| indexNameStrategy                | namespace_name         | Strategy for creating indexes like:- namespace_name or pod_name |
| resources                        | {}                     | Resources for fluentd daemonset pods                            |
| nodeSelectors                    | {}                     | Nodeselectors map key-values for fluentd daemonset pods         |
| affinity                         | {}                     | Affinity and anit-affinity for fluentd daemonset pods           |
| tolerations                      | {}                     | Tolerations and taints for fluentd daemonset pods               |
| customConfiguration              | {}                     | Custom configuration parameters for fluentd                     |
| additionalConfiguration          | {}                     | Additional configuration parameters for fluentd                 |
| esSecurity.enabled               | true                   | To enabled the xpack security of fluentd                        |
| esSecurity.elasticSearchPassword | elasticsearch-password | Credentials for elasticsearch authentication                    |

## CRD object definition parameters

These are the parameters that are currently supported by the Logging Operator for the Fluentd setup:-

- esCluster
- indexNameStrategy
- esSecurity
- customConfig
- additionalConfig
- kubernetesConfig

### esCluster

`esCluster` is a general parameter of Fluentd CRD for providing the information about Elasticsearch nodes.

```yaml
  esCluster:
    host: elasticsearch-master
```

### indexNameStrategy

`indexNameStrategy` naming standard for the indexes created inside the Elasticsearch cluster, It could be based on namespace like `kubernetes-marketing-2022-07-04` or based on application/pod name `kubernetes-gateway-application-2022-07-04`.

```yaml
  indexNameStrategy: namespace_name
```

### esSecurity

`esSecurity` s the security specification for Fluentd CRD. If we want to enable authentication and TLS, in that case, we can enable this configuration. To enable the authentication we need to provide secret reference in Kubernetes.

```yaml
  esSecurity:
    tlsEnabled: true
    existingSecret: elasticsearch-password
```

### customConfig

`customConfig` is a field of Fluentd definition through which existing configuration of Fluentd can be overwritten, but be cautious while making this change because it can break the Fluentd.

```yaml
  customConfig: fluentd-custom-config
```

### additionalConfig

`additionalConfig` is a field of Fluentd definition through which additional configuration can be mounted inside the Fluentd log-shipper. Additional configmap will be part of fluentd configuration.

```yaml
  additionalConfig: fluentd-additional-config
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
