---
title: "Continous Integration Pipeline"
weight: 3
linkTitle: "Continous Integration Pipeline"
description: >
    Continous Integration Pipeline for Logging Operator
---

We are using Azure DevOps pipeline for the Continous Integration in the Logging Operator. It checks all the important checks for the corresponding Pull Request. Also, this pipeline is capable of making releases on Quay, Dockerhub, and GitHub.

The pipeline definition can be edited inside the [.azure-pipelines](https://github.com/OT-CONTAINER-KIT/logging-operator/tree/main/.azure-pipelines).

![](https://github.com/OT-CONTAINER-KIT/mongodb-operator/blob/main/static/mongodb-ci-pipeline.png?raw=true)

Tools used for CI process:-

- **Golang ---> https://go.dev/**
- **Golang CI Lint ---. https://github.com/golangci/golangci-lint**
- **Hadolint ---> https://github.com/hadolint/hadolint**
- **GoSec ---> https://github.com/securego/gosec**
- **Trivy ---> https://github.com/aquasecurity/trivy**

![](https://deep-image.ai/images/2022-06-06/c822bf05-860e-49e0-9a80-514155f968c1.png)
