# Walkthrough

You'll need Helm to run a bunch of things:
Package Management: https://cloudnative.oracle.com/template.html#application-development/package-management/helm/readme.md

## Up and running:

### Logging
Read more about EFK logging, and it's install:
https://cloudnative.oracle.com/template.html#observability-and-analysis/logging/efk-stack

Download and install:
```
curl -o /tmp/efk.tar.gz  https://raw.githubusercontent.com/oracle/cloudnative/master/observability-and-analysis/logging/efk-stack/efk.tar.gz

tar -xzvf /tmp/efk.tar.gz
```

Get Kibana the URL and port:
```
kubectl get services | grep kibana-logging
```

### Monitoring
Read more about how to setup prometheus: https://cloudnative.oracle.com/template.html#observability-and-analysis/telemetry/prometheus

Install:
```
helm install stable/prometheus --name prom-demo -f ./src/prometheus-values.yaml
```

Get the URL:
Node port service (http://(node-ip):30001). To get the node-ip:
```
kubectl get nodes
```

### Tracing
Install Jaeger via instructions: [https://github.com/jaegertracing/jaeger-kubernetes](https://github.com/jaegertracing/jaeger-kubernetes)

*Note: This is just the dev setup; don't use for prod*

```
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-kubernetes/master/all-in-one/jaeger-all-in-one-template.yml

```

## Demo it with an app

Demo ecommerce application.

Deploy:
```
kubectl apply -f microservices.yaml
```

Look at logging with products, monitoring with users and telemetry with cart

## Read More

Jaeger Helm chart: [https://github.com/helm/charts/tree/master/incubator/jaeger](https://github.com/helm/charts/tree/master/incubator/jaeger)

Opentracing Tutorial: [https://github.com/yurishkuro/opentracing-tutorial](https://github.com/yurishkuro/opentracing-tutorial)
