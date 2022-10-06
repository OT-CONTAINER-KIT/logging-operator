---
title: "Kibana Setup"
weight: 4
linkTitle: "Kibana Setup"
description: >
    Kibana setup and management using logging operator
---

The operator is capable for setting up Kibana as a visualization and dashboard tool for elasticsearch cluster. There are few additional functionalities added to this CRD.

<div align="center">
    <img src="https://github.com/OT-CONTAINER-KIT/logging-operator/blob/master/static/kibana-architecture.png?raw=true">
</div>

## Setup using Helm (Deployment Tool)

Add the helm repository, so that Kibana chart can be available for the installation. The repository can be added by:-

```shell
# Adding helm repository
$ helm repo add ot-helm https://ot-container-kit.github.io/helm-charts/
...
"ot-helm" has been added to your repositories
```

If the repository is added make sure you have updated it with the latest information.

```shell
# Updating ot-helm repository
$ helm repo update
```

Once all these things have completed, we can install Kibana cluster by using:-

```shell
# Install the helm chart of Kibana
$ helm upgrade kibana ot-helm/kibana --install --namespace ot-operators
...
NAME: kibana
LAST DEPLOYED: Sat Aug  6 23:51:28 2022
NAMESPACE: ot-operators
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
  CHART NAME: kibana
  CHART VERSION: 0.3.2
  APP VERSION: 0.3.0

The helm chart for Kibana setup has been deployed.

Get the list of pods by executing:
    kubectl get pods --namespace ot-operators -l 'app=kibana'

For getting the credential for admin user:
    kubectl get kibana kibana -n ot-operators
```

Verify the pod status value by using:-

```shell
# Verify the status of the pods
$ kubectl get pods --namespace ot-operators -l 'app=kibana'
...
NAME                      READY   STATUS    RESTARTS   AGE
kibana-7b649df777-nkr2p   1/1     Running   0          3m27s
```

Kibana deployment can be listed and verify using `kubectl cli` as well.

```shell
$ kubectl get kibana -n ot-operators
NAME     VERSION   ES CLUSTER
kibana   7.17.0    elasticsearch
```

## Setup by Kubectl (Kubernetes CLI)

It is not a recommended way for setting for Kibana, it can be used for the POC and learning of Logging operator deployment.

All the kubectl related manifest are located inside the [example](https://github.com/OT-CONTAINER-KIT/logging-operator/tree/master/examples/kibana) folder which can be applied using `kubectl apply -f`.

For an example:-

```shell
$ kubectl apply -f examples/kibana/basic/kibana.yaml -n ot-operators
...
kibana.logging.logging.opstreelabs.in/kibana is created
```

## Validation of Kibana

To validate the state of Kibana, we can verify the log status of kibana pods managed by deployment.

```shell
# Validation of kibana logs
$ kubectl logs kibana-7bc5cd8747-pgtzc -n ot-operators
...
{"type":"log","@timestamp":"2022-08-06T18:22:04+00:00","tags":["info","plugins-service"],"pid":8,"message":"Plugin \"metricsEntities\" is disabled."}
{"type":"log","@timestamp":"2022-08-06T18:22:04+00:00","tags":["info","http","server","Preboot"],"pid":8,"message":"http server running at http://0.0.0.0:5601"}
{"type":"log","@timestamp":"2022-08-06T18:22:04+00:00","tags":["warning","config","deprecation"],"pid":8,"message":"Starting in 8.0, the Kibana logging format will be changing. This may affect you if you are doing any special handling of your Kibana logs, such as ingesting logs into Elasticsearch for further analysis. If you are using the new logging configuration, you are already receiving logs in both old and new formats, and the old format will simply be going away. If you are not yet using the new logging configuration, the log format will change upon upgrade to 8.0. Beginning in 8.0, the format of JSON logs will be ECS-compatible JSON, and the default pattern log format will be configurable with our new logging system. Please refer to the documentation for more information about the new logging format."}
```

Also, the UI can be accessed on `5601` port for validation.

![](https://github.com/OT-CONTAINER-KIT/logging-operator/blob/master/static/kibana-ui.png?raw=true)
