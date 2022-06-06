---
title: "Elasticsearch Setup"
weight: 3
linkTitle: "Elasticsearch Setup"
description: >
    Elasticsearch setup and management using logging operator
---

The operator is capable for setting up elasticsearch cluster with all the best practices in terms of security, performance and reliability.

There are different elasticsearch nodes supported by this operator:-

- **Master Node:** A node that has the master role (default), which makes it eligible to be elected as the master node, which controls the cluster.
- **Data Node:** A node that has the data role (default). Data nodes hold data and perform data related operations such as CRUD, search, and aggregations.
- **Ingestion Node:** A node that has ingest role (default). Ingest nodes are able to apply an ingest pipeline to a document in order to transform and enrich the document before indexing. With a heavy ingest load, it makes sense to use dedicated ingest nodes and to not include ingest role from nodes that have the master or data roles.
- **Client or Coordinator Node:** Requests like search requests or bulk-indexing requests may involve data held on different data nodes. A search request, for example, is executed in two phases which are coordinated by the node which receives the client request the coordinating node.

There are few additional functionalities supported in the elasticsearch CRD.

- TLS support and xpack support
- Multi node cluster setup - master, data, ingestion, client
- Custom configuration for each type of elasticsearch node

<div align="center">
    <img src="https://github.com/OT-CONTAINER-KIT/logging-operator/blob/master/static/es-architecture.png?raw=true">
</div>

## Setup using Helm (Deployment Tool)

Add the helm repository, so that Elasticsearch chart can be available for the installation. The repository can be added by:-

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

Once all these things have completed, we can install Elasticsearch cluster by using:-

```shell
# Install the helm chart of Elasticsearch
$ helm install elasticsearch ot-helm/elasticsearch --namespace ot-operators \
  --set esMaster.storage.storageClass=do-block-storage \
  --set esData.storage.storageClass=do-block-storage
...
NAME: elasticsearch
LAST DEPLOYED: Mon Jun  6 15:06:45 2022
NAMESPACE:     ot-operators
STATUS:        deployed
REVISION:      1
TEST SUITE:    None
NOTES:
  CHART NAME:    elasticsearch
  CHART VERSION: 0.3.1
  APP VERSION:   0.3.0

The helm chart for Elasticsearch setup has been deployed.

Get the list of pods by executing:
    kubectl get pods --namespace ot-operators -l 'role in (master,data,ingestion,client)'

For getting the credential for admin user:
    kubectl get secrets -n ot-operators elasticsearch-password -o jsonpath="{.data.password}" | base64 -d
```

Verify the pod status and secret value by using:-

```shell
# Verify the status of the pods
$ kubectl get pods --namespace ot-operators -l 'role in (master,data,ingestion,client)'
...
NAME                     READY   STATUS    RESTARTS   AGE
elasticsearch-data-0     1/1     Running   0          77s
elasticsearch-data-1     1/1     Running   0          77s
elasticsearch-data-2     1/1     Running   0          77s
elasticsearch-master-0   1/1     Running   0          77s
elasticsearch-master-1   1/1     Running   0          77s
elasticsearch-master-2   1/1     Running   0          77s
```

```shell
# Verify the secret value
$ kubectl get secrets -n ot-operators elasticsearch-password -o jsonpath="{.data.password}" | base64 -d
...
EuDyr4A105EjqaNW
```

Elasticsearch cluster can be listed and verify using kubectl cli as well.

```shell
$ kubectl get elasticsearch -n ot-operators
...
NAME            VERSION   STATE   SHARDS   INDICES
elasticsearch   7.17.0    green   2        2
```

## Setup by Kubectl (Kubernetes CLI)

It is not a recommended way for setting for Elasticsearch cluster, it can be used for the POC and learning of Logging operator deployment.

All the kubectl related manifest are located inside the [example](https://github.com/OT-CONTAINER-KIT/logging-operator/tree/master/examples/elasticsearch) folder which can be applied using `kubectl apply -f`.

For an example:-

```shell
$ kubectl apply -f examples/elasticsearch/basic-cluster/basic-elastic.yaml -n ot-operators
...
elasticsearch/elasticsearch is created
```

## Validation of Elasticsearch

To validate the state of Elasticsearch cluster, we can take the shell access of the Elasticsearch pod and verify elasticsearch version and details using `curl` command.

```shell
# Verify endpoint of elasticsearch
$ export ELASTIC_PASSWORD=$(kubectl get secrets -n ot-operators \
  elasticsearch-password -o jsonpath="{.data.password}" | base64 -d)

$ kubectl exec -it elasticsearch-master-0 -c elastic -n ot-operators \
  -- curl -u elastic:$ELASTIC_PASSWORD -k https://localhost:9200
...
{
  "name" : "elasticsearch-master-0",
  "cluster_name" : "elastic-prod",
  "cluster_uuid" : "vPtAZQt9SEWsl8NSfNVYzw",
  "version" : {
    "number" : "7.17.0",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "bee86328705acaa9a6daede7140defd4d9ec56bd",
    "build_date" : "2022-01-28T08:36:04.875279988Z",
    "build_snapshot" : false,
    "lucene_version" : "8.11.1",
    "minimum_wire_compatibility_version" : "6.8.0",
    "minimum_index_compatibility_version" : "6.0.0-beta1"
  },
  "tagline" : "You Know, for Search"
}
```

### Node status and health

Once the version details are verified we can list down the nodes connected to elasticsearch cluster and their health status. Also, we can verify the status health of complete elasticsearch cluster.

```shell
# Cluster health of elasticsearch cluster
$ kubectl exec -it elasticsearch-master-0 -c elastic -n ot-operators \
  -- curl -u elastic:$ELASTIC_PASSWORD -k https://localhost:9200/_cluster/health
...
{
  "cluster_name": "elastic-prod",
  "status": "green",
  "timed_out": false,
  "number_of_nodes": 6,
  "number_of_data_nodes": 3,
  "active_primary_shards": 1,
  "active_shards": 2,
  "relocating_shards": 0,
  "initializing_shards": 0,
  "unassigned_shards": 0,
  "delayed_unassigned_shards": 0,
  "number_of_pending_tasks": 0,
  "number_of_in_flight_fetch": 0,
  "task_max_waiting_in_queue_millis": 0,
  "active_shards_percent_as_number": 100
}
```

```shell
# Node status of elasticsearch
$ kubectl exec -it elasticsearch-master-0 -c elastic -n ot-operators \
  -- curl -u elastic:$ELASTIC_PASSWORD -k https://localhost:9200/_cat/nodes
...
ip           heap.percent ram.percent cpu load_1m load_5m load_15m node.role master name
10.244.1.69            54          19   0    0.00    0.00     0.01 m         -      elasticsearch-master-2
10.244.0.82            43          20   0    0.00    0.00     0.00 d         -      elasticsearch-data-2
10.244.0.150           28          19   0    0.00    0.12     0.12 d         -      elasticsearch-data-1
10.244.0.13            57          19   1    0.00    0.00     0.00 m         -      elasticsearch-master-0
10.244.1.72            13          20   0    0.00    0.00     0.01 d         -      elasticsearch-data-0
10.244.0.161           61          20   2    0.00    0.12     0.12 m         *      elasticsearch-master-1
```