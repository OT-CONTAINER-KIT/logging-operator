<p align="center">
  <img src="./static/logging-operator-logo.svg" height="220" width="220">
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/OT-CONTAINER-KIT/logging-operator">
    <img src="https://goreportcard.com/badge/github.com/OT-CONTAINER-KIT/logging-operator" alt="GoReportCard">
  </a>
  <a href="http://golang.org">
    <img src="https://img.shields.io/github/go-mod/go-version/OT-CONTAINER-KIT/logging-operator" alt="GitHub go.mod Go version (subdirectory of monorepo)">
  </a>
  <a href="https://quay.io/repository/opstree/logging-operator">
    <img src="https://img.shields.io/badge/container-ready-green" alt="Docker">
  </a>
  <a href="https://github.com/OT-CONTAINER-KIT/logging-operator/master/LICENSE">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  </a>
</p>

Logging Operator is an operator created in Golang to setup and manage EFK(Elasticsearch, Fluentd, and Kibana) cluster inside Kubernetes and Openshift environment. This operator is capable of setting up each individual component of EFK cluster separately.

For documentation, please refer to [https://ot-logging-operator.netlify.app/](https://ot-logging-operator.netlify.app/)

## Architecture


```shell
$ kubectl apply -f config/crd/bases/
$ kubectl apply -f config/manager/manager.yaml
$ kubectl apply -f config/rbac/service_account.yaml
$ kubectl apply -f config/rbac/role.yaml
$ kubectl apply -f config/rbac/role_binding.yaml
```
