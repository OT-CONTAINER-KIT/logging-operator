---
title: "Elasticsearch Config"
weight: 2
linkTitle: "Elasticsearch Config"
description: >
    Elasticsearch configuration paramaters with respect to logging operator
---

Elasticsearch configuration is easily customizable using `helm` as well `kubectl`. Since all the configurations are in the form YAML file, it can be easily changed and customized.

The values.yaml file for Elasticsearch setup can be found [here](https://github.com/OT-CONTAINER-KIT/helm-charts/tree/main/charts/elasticsearch). But if the setup is not done using Helm, in that case Kubernetes manifests needs to be customized.

## Helm Chart Parameters

| **Name**                         | **Value**       | **Description**                                                    |
|----------------------------------|-----------------|--------------------------------------------------------------------|
| clusterName                      | elastic-prod    | Name of the elasticsearch cluster                                  |
| esVersion                        | 7.17.0          | Major and minor version of elaticsearch                            |
| esPlugins                        | []              | Plugins list to install inside elasticsearch                       |
| esKeystoreSecret                 | -               | Keystore secret to include in elasticsearch cluster                |
| customConfiguration              | {}              | Additional configuration parameters for elasticsearch              |
| esSecurity.enabled               | true            | To enabled the xpack security of elasticsearch                     |
| esMaster.replicas                | 3               | Number of replicas for elasticsearch master node                   |
| esMaster.storage.storageSize     | 20Gi            | Size of the elasticsearch persistent volume for master             |
| esMaster.storage.accessModes     | [ReadWriteOnce] | Access modes of the elasticsearch persistent volume for master     |
| esMaster.storage.storageClass    | default         | Storage class of the elasticsearch persistent volume for master    |
| esMaster.jvmMaxMemory            | 1Gi             | Java max memory for elasticsearch master node                      |
| esMaster.jvmMinMemory            | 1Gi             | Java min memory for elasticsearch master node                      |
| esMaster.resources               | {}              | Resources for elasticsearch master pods                            |
| esMaster.nodeSelectors           | {}              | Nodeselectors map key-values for elasticsearch master pods         |
| esMaster.affinity                | {}              | Affinity and anit-affinity for elasticsearch master pods           |
| esMaster.tolerations             | {}              | Tolerations and taints for elasticsearch master pods               |
| esData.replicas                  | 3               | Number of replicas for elasticsearch data node                     |
| esData.storage.storageSize       | 50Gi            | Size of the elasticsearch persistent volume for data               |
| esData.storage.accessModes       | [ReadWriteOnce] | Access modes of the elasticsearch persistent volume for data       |
| esData.storage.storageClass      | default         | Storage class of the elasticsearch persistent volume for data      |
| esData.jvmMaxMemory              | 1Gi             | Java max memory for elasticsearch data node                        |
| esData.jvmMinMemory              | 1Gi             | Java min memory for elasticsearch data node                        |
| esData.resources                 | {}              | Resources for elasticsearch data pods                              |
| esData.nodeSelectors             | {}              | Nodeselectors map key-values for elasticsearch data pods           |
| esData.affinity                  | {}              | Affinity and anit-affinity for elasticsearch data pods             |
| esData.tolerations               | {}              | Tolerations and taints for elasticsearch data pods                 |
| esIngestion.replicas             | -               | Number of replicas for elasticsearch ingestion node                |
| esIngestion.storage.storageSize  | -               | Size of the elasticsearch persistent volume for ingestion          |
| esIngestion.storage.accessModes  | -               | Access modes of the elasticsearch persistent volume for ingestion  |
| esIngestion.storage.storageClass | -               | Storage class of the elasticsearch persistent volume for ingestion |
| esIngestion.jvmMaxMemory         | -               | Java max memory for elasticsearch ingestion node                   |
| esIngestion.jvmMinMemory         | -               | Java min memory for elasticsearch ingestion node                   |
| esIngestion.resources            | -               | Resources for elasticsearch ingestion pods                         |
| esIngestion.nodeSelectors        | -               | Nodeselectors map key-values for elasticsearch ingestion pods      |
| esIngestion.affinity             | -               | Affinity and anit-affinity for elasticsearch ingestion pods        |
| esIngestion.tolerations          | -               | Tolerations and taints for elasticsearch ingestion pods            |
| esClient.replicas                | -               | Number of replicas for elasticsearch ingestion node                |
| esClient.storage.storageSize     | -               | Size of the elasticsearch persistent volume for client             |
| esClient.storage.accessModes     | -               | Access modes of the elasticsearch persistent volume for client     |
| esClient.storage.storageClass    | -               | Storage class of the elasticsearch persistent volume for client    |
| esClient.jvmMaxMemory            | -               | Java max memory for elasticsearch client node                      |
| esClient.jvmMinMemory            | -               | Java min memory for elasticsearch client node                      |
| esClient.resources               | -               | Resources for elasticsearch client pods                            |
| esClient.nodeSelectors           | -               | Nodeselectors map key-values for elasticsearch client pods         |
| esClient.affinity                | -               | Affinity and anit-affinity for elasticsearch client pods           |
| esClient.tolerations             | -               | Tolerations and taints for elasticsearch client pods               |

## CRD Object Definition Parameters

These are the parameters that are currently supported by the Logging Operator for the Elastisearch setup:-

- esClusterName
- esVersion
- esMaster
- esData
- esIngestion
- esClient
- esSecurity
- customConfig

### esClusterName

`esClusterName` is a parameter to define the name of elasticsearch cluster.

```yaml
esClusterName: "prod"
```

### esVersion

`esVersion` is a CRD option through which we can define the version of elasticsearch.

```yaml
esVersion: "7.16.0"
```

### esPlugins

`esPlugins` is a CRD parameter through which we can define the list of plugins that needs to install inside elasticsearch cluster.

```yaml
esPlugins: ["respository-s3", "repository-gcs"]
```

### esKeystoreSecret

`esKeystoreSecret` is a CRD parameter through which we can define the keystore related secret to include in elasticsearch cluster.

```yaml
esKeystoreSecret: keystore-secret
```

### esMaster

`esMaster` is a general configuration parameter for Elasticsearch CRD for defining the configuration of Elasticsearch Master node. This includes Kubernetes related configurations and Elasticsearch properties related configurations.

```yaml
  esMaster:
    replicas: 2
    storage:
      storageSize: 2Gi
      accessModes: [ReadWriteOnce]
      storageClass: do-block-storage
    jvmMaxMemory: "512m"
    jvmMinMemory: "512m"
    kubernetesConfig:
      elasticAffinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: beta.kubernetes.io/os
                    operator: In
                    values:
                      - linux
      nodeSelectors:
        kubernetes.io/os: linux
      priorityClassName: system-node-critical
      resources:
        requests:
          cpu: 101m
          memory: 512Mi
        limits:
          cpu: 2000m
          memory: 2Gi
      tolerations:
        - key: "example-key"
          operator: "Exists"
          effect: "NoSchedule"
```

**Note:- All properties defined under kubernetesConfig can be used for other elasticsearch node types as well.**

### esData

`esData` is a general configuration parameter for Elasticsearch CRD for defining the configuration of Elasticsearch Data node. This includes Kubernetes related configurations and Elasticsearch properties related configurations.

```yaml
  esData:
    replicas: 2
    storage:
      storageSize: 2Gi
      accessModes: [ReadWriteOnce]
      storageClass: do-block-storage
    jvmMaxMemory: "512m"
    jvmMinMemory: "512m"
```

### esIngestion

`esIngestion` is a general configuration parameter for Elasticsearch CRD for defining the configuration of Elasticsearch Ingestion node. This includes Kubernetes related configurations and Elasticsearch properties related configurations.

```yaml
  esIngestion:
    replicas: 2
    storage:
      storageSize: 2Gi
      accessModes: [ReadWriteOnce]
      storageClass: do-block-storage
    jvmMaxMemory: "512m"
    jvmMinMemory: "512m"
```

### esClient

`esClient` is a general configuration parameter for Elasticsearch CRD for defining the configuration of Elasticsearch Client node. This includes Kubernetes related configurations and Elasticsearch properties related configurations.

```yaml
  esClient:
    replicas: 2
    storage:
      storageSize: 2Gi
      accessModes: [ReadWriteOnce]
      storageClass: do-block-storage
    jvmMaxMemory: "512m"
    jvmMinMemory: "512m"
```

### esSecurity

`esSecurity` s the security specification for Elasticsearch CRD. If we want to enable authentication and TLS, in that case, we can enable this configuration. To enable the authentication we need to provide secret reference in Kubernetes.

```yaml
  esSecurity:
    autoGeneratePassword: true
    tlsEnabled: true  
#   existingSecret: elastic-custom-password
```

### customConfig

`customConfig` is a Elasticsearch config file parameter through which we can provide custom configuration to elasticsearch nodes. This property is applicable for all types of nodes in elasticsearch.

```yaml
  esMaster:
    replicas: 3
    storage:
      storageSize: 2Gi
      accessModes: [ReadWriteOnce]
      storageClass: do-block-storage
    customConfig: elastic-additional-config
```
