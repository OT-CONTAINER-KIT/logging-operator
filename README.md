<p align="center">
  <img src="./static/logging-operator-logo.svg" height="220" width="220">
</p>

```shell
$ kubectl apply -f config/crd/bases/
$ kubectl apply -f config/manager/manager.yaml
$ kubectl apply -f config/rbac/service_account.yaml
$ kubectl apply -f config/rbac/role.yaml
$ kubectl apply -f config/rbac/role_binding.yaml
```
