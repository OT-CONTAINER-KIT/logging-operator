---
title: "Fluentd Setup"
weight: 3
linkTitle: "Fluentd Setup"
description: >
    Fluentd setup and management using logging operator
---

The operator is capable for setting up fluentd as a log shipper to trace, collect and ship logs to elasticsearch cluster. There are few additional functionalities added to this CRD.

- Namespace and application name based indexes
- Custom and additional configuration support
- TLS and auth support for authentication

<div align="center">
    <img src="https://github.com/OT-CONTAINER-KIT/logging-operator/blob/master/static/fluentd-architecture-sharpened_upscaled_x2.png?raw=true">
</div>

## Setup using Helm (Deployment Tool)

Add the helm repository, so that Fluentd chart can be available for the installation. The repository can be added by:-

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

Once all these things have completed, we can install Fluentd cluster by using:-

```shell
# Install the helm chart of Fluentd
$ helm install fluentd ot-helm/fluentd --namespace ot-operators
...
NAME:          fluentd
LAST DEPLOYED: Mon Jun  6 19:37:11 2022
NAMESPACE:     ot-operators
STATUS:        deployed
REVISION:      1
TEST SUITE:    None
NOTES:
  CHART NAME:    fluentd
  CHART VERSION: 0.3.0
  APP VERSION:   0.3.0

The helm chart for Fluentd setup has been deployed.

Get the list of pods by executing:
    kubectl get pods --namespace ot-operators -l 'app=fluentd'

For getting the credential for admin user:
    kubectl get fluentd fluentd -n ot-operators
```

Verify the pod status and secret value by using:-

```shell
# Verify the status of the pods
$ kubectl get pods --namespace ot-operators -l 'app=fluentd'.
...
NAME            READY   STATUS    RESTARTS   AGE
fluentd-7w48q   1/1     Running   0          3m9s
fluentd-dgcwx   1/1     Running   0          3m9s
fluentd-kq52c   1/1     Running   0          3m9s
```

Fluentd daemonset can be listed and verify using `kubectl cli` as well.

```shell
$ kubectl get fluentd -n ot-operators
...
NAME      ELASTICSEARCH HOST     TOTAL AGENTS
fluentd   elasticsearch-master   3
```

## Setup by Kubectl (Kubernetes CLI)

It is not a recommended way for setting for Fluentd, it can be used for the POC and learning of Logging operator deployment.

All the kubectl related manifest are located inside the [example](https://github.com/OT-CONTAINER-KIT/logging-operator/tree/master/examples/fluentd) folder which can be applied using `kubectl apply -f`.

For an example:-

```shell
$ kubectl apply -f examples/fluentd/basic/fluentd.yaml -n ot-operators
...
fluentd/fluentd is created
```

## Validation of Fluentd

To validate the state of Fluentd, we can verify the log status of fluentd pods managed by daemonset.

```shell
# Validation of fluentd logs
$ kubectl logs fluentd-7w48q -n ot-operators
...
2022-06-06 14:07:28 +0000 [info]: #0 [in_tail_container_logs] following tail of /var/log/containers/fluentd-7w48q_ot-operators_fluentd-f49b48f7f447d05139819861b8b17c30e2bf2de094e25e23d1e9c5a274fd3d7e.log
2022-06-06 14:07:28 +0000 [info]: #0 fluentd worker is now running worker=0
2022-06-06 14:07:54 +0000 [info]: #0 [filter_kube_metadata] stats - namespace_cache_size: 5, pod_cache_size: 32, namespace_cache_api_updates: 16, pod_cache_api_updates: 16, id_cache_miss: 16, pod_cache_host_updates: 32, namespace_cache_host_updates: 5
```

Also, we can list down the indices using the `curl` command from the elasticsearch pod/container. If indices are available inside the elasticsearch that means fluentd is shipping the logs to elasticsearch without any issues.

```shell
$ export ELASTIC_PASSWORD=$(kubectl get secrets -n ot-operators \
  elasticsearch-password -o jsonpath="{.data.password}" | base64 -d)

$ kubectl exec -it elasticsearch-master-0 -c elastic -n ot-operators \
  -- curl -u elastic:$ELASTIC_PASSWORD -k "https://localhost:9200/_cat/indices?v"
...
health status index                              uuid                   pri rep docs.count docs.deleted store.size pri.store.size
green  open   .geoip_databases                   _GEkcekFSr2KY1Z4jFmWRQ   1   1         40            0     76.4mb         38.2mb
green  open   kubernetes-ot-operators-2022.06.06 QlS_dyjzQ8qIXQi2PgpABA   1   1      20665            0      7.9mb          3.9mb
green  open   kubernetes-kube-system-2022.06.06  vWQ5IzoHQWW9zl8bQk0jlw   1   1      12006            0        7mb          4.3mb
```
