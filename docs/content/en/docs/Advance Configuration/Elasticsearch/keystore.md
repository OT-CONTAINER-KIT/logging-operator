---
title: "Keystore Integration"
linkTitle: "Keystore Integration"
weight: 6
description: >
    Keystore integration configuration for elasticsearch
---

## Keystore integation

Keystore is a recommended way of integrating different credentials like:- AWS, GCP, Azure and Slack, etc. to elasticsearch cluster. We simply need to create a Kubernetes secret and the operator can take care the integration of Kubernetes secret to elasticsearch keystore.

```shell
$ kubectl create secret generic slack-hook \
--from-literal=xpack.notification.slack.account.monitoring.secure_url='https://hooks.slack.com/services/asdasdasd/asda'
```

or yaml file is also one of the way for creating the secret.

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: encryption-key
data:
  xpack.notification.slack.account.monitoring.secure_url: aHR0cHM6Ly9ob29rcy5zbGFjay5jb20vc2VydmljZXMvYXNkYXNkYXNkL2FzZGFzZGFzL2FzZGFzZA==
# other secrets key value pairs can be defined here
type: Opaque
```

Then simply we can define the keystore secret name in CRD definition.

```yaml
---
apiVersion: logging.logging.opstreelabs.in/v1beta1
kind: Elasticsearch
metadata:
  name: elasticsearch
spec:
  esClusterName: "prod"
  esVersion: "7.16.0"
  esKeystoreSecret: encryption-key
```

Validation of keystore can be done using `elasticsearch-keystore` command.

```shell
$ kubectl exec -it elasticsearch-master-0 -n ot-operators -- ./bin/elasticsearch-keystore list
...
keystore.seed
xpack.notification.slack.account.monitoring.secure_url
```

## Helm Configuration

Keystore integration can also be done using helm chart of elasticsearch. We just need to define the keystore secret name in the values file of elasticsearch helm chart.

https://github.com/OT-CONTAINER-KIT/helm-charts/blob/main/charts/elasticsearch/values.yaml#L9

```shell
$ helm upgrade elasticsearch ot-helm/elasticsearch --namespace ot-operators \
  --set esMaster.storage.storageClass=do-block-storage \
  --set esData.storage.storageClass=do-block-storage --install
```
