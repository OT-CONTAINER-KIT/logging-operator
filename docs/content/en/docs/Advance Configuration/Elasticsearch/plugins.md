---
title: "Plugins Management"
linkTitle: "Plugins Management"
weight: 5
description: >
    Plugins related configuration for elasticsearch
---

## Plugins Installation 

Plugins installation has been simplified using the Logging Operator. To install the plugins inside the elasticsearch, we just need to define the list of plugins inside the `esPlugins` section of elasticsearch CRD.

For example:-

```yaml
---
apiVersion: logging.logging.opstreelabs.in/v1beta1
kind: Elasticsearch
metadata:
  name: elasticsearch
spec:
  esClusterName: "prod"
  esVersion: "7.16.0"
  esPlugins: ["repository-s3"]
```

Validation of plugin installation can be done using `elasticsearch-plugin` or `curl` command.

```shell
$ kubectl exec -it elasticsearch-master-0 -n ot-operators -- ./bin/elasticsearch-plugin list
...
repository-s3
```

```shell
$ kubectl exec -it elasticsearch-master-0 -n ot-operators -- curl http://localhost:9200/_cat/plugins
...
elasticsearch-master-1 repository-s3 7.16.0
elasticsearch-master-2 repository-s3 7.16.0
elasticsearch-master-0 repository-s3 7.16.0
```

## Helm Configuration

Plugin installation can also be done using helm chart of elasticsearch. We just need to define the plugins list in the values file of elasticsearch helm chart.

https://github.com/OT-CONTAINER-KIT/helm-charts/blob/main/charts/elasticsearch/values.yaml#L8

```shell
$ helm upgrade elasticsearch ot-helm/elasticsearch --namespace ot-operators \
  --set esMaster.storage.storageClass=do-block-storage \
  --set esData.storage.storageClass=do-block-storage --install
```
