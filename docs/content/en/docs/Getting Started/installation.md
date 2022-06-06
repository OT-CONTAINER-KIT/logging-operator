---
title: "Installation"
weight: 2
linkTitle: "Installation"
description: >
    Logging Operator installation, upgrade guide
---

Logging operator is based on the CRD framework of Kubernetes, for more information about the CRD framework please refer to the [official documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). In a nutshell, CRD is a feature through which we can develop our own custom API's inside Kubernetes.

The API versions for Logging Operator available are:-

- ElasticSearch
- Fluentd
- Kibana

Logging Operator requires a Kubernetes cluster of version >=1.16.0. If you have just started with the CRD and Operators, its highly recommended using the latest version of Kubernetes.

Setup of Logging operator can be easily done by using simple [helm](https://helm.sh) and [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) commands.

{{< alert title="Note" >}}The recommded of way of installation is helm.{{< /alert >}}

## Setup using Helm tool

The setup can be done by using helm. The logging-operator can easily get installed using helm commands.

```shell
# Add the helm chart
$ helm repo add ot-helm https://ot-container-kit.github.io/helm-charts/
...
"ot-helm" has been added to your repositories
```

```shell
# Deploy the Logging Operator
$ helm upgrade logging-operator ot-helm/logging-operator \
  --install --namespace ot-operators
...
Release "logging-operator" does not exist. Installing it now.
NAME:          logging-operator
LAST DEPLOYED: Sun May 29 01:06:58 2022
NAMESPACE:     ot-operators
STATUS:        deployed
REVISION:      1
```

After the deployment, verify the installation of operator.

```shell
# Testing Operator
$ helm test logging-operator --namespace ot-operators
...
NAME:           logging-operator
LAST DEPLOYED:  Sun May 29 01:06:58 2022
NAMESPACE:      ot-operators
STATUS:         deployed
REVISION:       1
TEST SUITE:     logging-operator-test-connection
Last Started:   Sun May 29 01:07:56 2022
Last Completed: Sun May 29 01:08:02 2022
Phase:          Succeeded
```

Verify the deployment of Logging Operator using `kubectl` command.

```shell
# List the pod and status of logging-operator
$ kubectl get pods -n ot-operators -l name=logging-operator
...
NAME                               READY   STATUS    RESTARTS   AGE
logging-operator-fc88b45b5-8rmtj   1/1     Running   0          21d
```

## Setup using Kubectl

In any case using helm chart is not a possiblity, the Logging operator can be installed by `kubectl` commands as well.

As a first step, we need to set up a namespace and then deploy the CRD definitions inside Kubernetes.

```shell
# Setup of CRDS
$ kubectl apply -f https://raw.githubusercontent.com/OT-CONTAINER-KIT/logging-operator/master/config/crd/bases/logging.logging.opstreelabs.in_elasticsearches.yaml
$ kubectl apply -f https://raw.githubusercontent.com/OT-CONTAINER-KIT/logging-operator/master/config/crd/bases/logging.logging.opstreelabs.in_fluentds.yaml
$ kubectl apply -f https://raw.githubusercontent.com/OT-CONTAINER-KIT/logging-operator/master/config/crd/bases/logging.logging.opstreelabs.in_kibanas.yaml
$ kubectl apply -f https://github.com/OT-CONTAINER-KIT/logging-operator/raw/master/config/crd/bases/logging.logging.opstreelabs.in_indextemplates.yaml
$ kubectl apply -f https://github.com/OT-CONTAINER-KIT/logging-operator/raw/master/config/crd/bases/logging.logging.opstreelabs.in_indexlifecycles.yaml
```

Once we have namespace in the place, we need to set up the RBAC related stuff like:- ClusterRoleBindings, ClusterRole, Serviceaccount.

```shell
# Setup of RBAC account
$ kubectl apply -f https://raw.githubusercontent.com/OT-CONTAINER-KIT/logging-operator/main/config/rbac/service_account.yaml
$ kubectl apply -f https://raw.githubusercontent.com/OT-CONTAINER-KIT/logging-operator/main/config/rbac/role.yaml
$ kubectl apply -f https://github.com/OT-CONTAINER-KIT/logging-operator/blob/main/config/rbac/role_binding.yaml
```

As last part of the setup, now we can deploy the Logging Operator as deployment of Kubernetes.

```shell
# Deployment for MongoDB Operator
$ kubectl apply -f https://github.com/OT-CONTAINER-KIT/logging-operator/raw/main/config/manager/manager.yaml
```

Verify the deployment of Logging Operator using `kubectl` command.

```shell
# List the pod and status of logging-operator
$ kubectl get pods -n ot-operators -l name=logging-operator
...
NAME                               READY   STATUS    RESTARTS   AGE
logging-operator-fc88b45b5-8rmtj   1/1     Running   0          21d
```
