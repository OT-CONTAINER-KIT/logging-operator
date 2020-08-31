#!/bin/bash

lint_chart() {
    echo "--------------Linting Helm Chart--------------"
    helm lint ./helm-charts/logging-operator/
}

deploy_chart() {
    echo "--------------Deploying Helm Chart--------------"
    helm upgrade logging-operator \
        ./helm-charts/logging-operator/ \
        -f ./helm-charts/logging-operator/values.yaml \
        --namespace logging-operator --install
}

validate_chart() {
    echo "--------------Testing Helm Chart--------------"
    helm test logging-operator --namespace logging-operator
}

validate_container_state() {
    echo "--------------Validating Deployment Status--------------"
    output=$(kubectl get pods -n logging-operator -l app.kubernetes.io/name=logging-operator \
    -o jsonpath="{.items[*]['status.phase']}")
    if [ "${output}" != "Running" ] && [ "${output}" != "" ]
    then
        echo "Container is not healthy"
        exit 1
    else
        echo "Container is running fine"
    fi
}

main_function() {
    lint_chart
    deploy_chart
    validate_chart
    sleep 30s
    validate_container_state
}

main_function
