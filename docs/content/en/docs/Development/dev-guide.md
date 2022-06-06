---
title: "Development Guide"
weight: 2
linkTitle: "Development Guide"
description: >
   Development guide for Logging Operator
---

## Pre-requisites

**Access to Kubernetes cluster**

First, you will need access to a Kubernetes cluster. The easiest way to start is minikube.

- [Virtualbox](https://www.virtualbox.org/wiki/Downloads)
- [Minikube](https://kubernetes.io/docs/setup/minikube/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

**Tools to build an Operator**

Apart from kubernetes cluster, there are some tools which are needed to build and test the Logging Operator.

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/)
- [Docker](https://docs.docker.com/install/)
- [Operator SDK](https://github.com/operator-framework/operator-sdk/blob/v0.8.1/doc/user/install-operator-sdk.md)
- [Make](https://www.gnu.org/software/make/manual/make.html)

## Building Operator

To build the operator on local system, we can use `make` command.

```shell
$ make manager
...
go build -o bin/manager main.go
```

MongoDB operator gets packaged as a container image for running on the Kubernetes cluster.

```shell
$ make docker-build
...
[+] Building 124.8s (19/19) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                         0.1s
 => => transferring dockerfile: 866B                                                                                                                                         0.0s
 => [internal] load .dockerignore                                                                                                                                            0.1s
 => => transferring context: 171B                                                                                                                                            0.0s
 => [internal] load metadata for gcr.io/distroless/static:nonroot                                                                                                            1.6s
 => [internal] load metadata for docker.io/library/golang:1.17                                                                                                               0.0s
 => CACHED [stage-1 1/3] FROM gcr.io/distroless/static:nonroot@sha256:2556293984c5738fc75208cce52cf0a4762c709cf38e4bf8def65a61992da0ad                                       0.0s
 => [internal] load build context                                                                                                                                            0.1s
 => => transferring context: 265.64kB                                                                                                                                        0.1s
 => [builder  1/11] FROM docker.io/library/golang:1.17                                                                                                                       0.0s
 => CACHED [builder  2/11] WORKDIR /workspace                                                                                                                                0.0s
 => [builder  3/11] COPY go.mod go.mod                                                                                                                                       0.1s
 => [builder  4/11] COPY go.sum go.sum                                                                                                                                       0.1s
 => [builder  5/11] RUN go mod download                                                                                                                                     32.5s
 => [builder  6/11] COPY main.go main.go                                                                                                                                     0.0s
 => [builder  7/11] COPY api/ api/                                                                                                                                           0.0s
 => [builder  8/11] COPY controllers/ controllers/                                                                                                                           0.0s
 => [builder  9/11] COPY k8sgo/ k8sgo/                                                                                                                                       0.0s
 => [builder 10/11] COPY elasticgo/ elasticgo/                                                                                                                               0.0s
 => [builder 11/11] RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go                                                                                89.2s
 => [stage-1 2/3] COPY --from=builder /workspace/manager .                                                                                                                   0.2s
 => exporting to image                                                                                                                                                       0.3s
 => => exporting layers                                                                                                                                                      0.3s
 => => writing image sha256:0875d1dd92839e2722f50d9f6b0be6fbe60ac56f3e3aa13ecad3b1c6a5862330                                                                                 0.0s
 => => naming to quay.io/opstree/logging-operator:v0.3.1
```

If you want to play it on Kubernetes. You can use a minikube.

```shell
$ minikube start --vm-driver virtualbox
...
😄  minikube v1.0.1 on linux (amd64)
🤹  Downloading Kubernetes v1.14.1 images in the background ...
🔥  Creating kvm2 VM (CPUs=2, Memory=2048MB, Disk=20000MB) ...
📶  "minikube" IP address is 192.168.39.240
🐳  Configuring Docker as the container runtime ...
🐳  Version of container runtime is 18.06.3-ce
⌛  Waiting for image downloads to complete ...
✨  Preparing Kubernetes environment ...
🚜  Pulling images required by Kubernetes v1.14.1 ...
🚀  Launching Kubernetes v1.14.1 using kubeadm ... 
⌛  Waiting for pods: apiserver proxy etcd scheduler controller dns
🔑  Configuring cluster permissions ...
🤔  Verifying component health .....
💗  kubectl is now configured to use "minikube"
🏄  Done! Thank you for using minikube!
```

```shell
$ make test
```
