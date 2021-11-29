# Envoy test

## Prerequisites

[Common tools](doc/MACOS-COMMON-TOOLS.md)
[Create K8S cluster with KinD](doc/CREATE-K8S-CLUSTER-KIND.md)

## :computer: Build

Install docker images to the registry:
```shell
#export DOCKER_REGISTRY="localhost:5000/";
export DOCKER_REGISTRY="";

docker build docker/service-a -t ${DOCKER_REGISTRY}service-a:latest
docker build docker/service-b -t ${DOCKER_REGISTRY}service-b:latest
```

## :rocket: Deploy

Deploy images locally

```shell
# in case we are working with local images
export CLUSTER_NAME="envoy-playground";
kind load docker-image service-a:latest --name ${CLUSTER_NAME} 
kind load docker-image service-b:latest --name ${CLUSTER_NAME} 
```

Install services:
```shell
helm install envoy-playground helm --values helm/values.yaml
helm upgrade --install --force envoy-playground helm --values helm/values.yaml
```

Uninstall
```shell
helm uninstall envoy-playground
```

## :chart: Monitor

Deploy prometheus
```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
#helm install prometheus prometheus-community/prometheus
helm install prometheus prometheus-community/kube-prometheus-stack -f conf/helm-prometheus.yml
```

Expose grafana
```shell
kubectl apply -f conf/k8s-prometheus-ingress.yml
# https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus#configuration
```

## :white_check_mark: Test
```shell
# TEST

# health checks
curl -v http://localhost/service-a/health/liveness
curl -v http://localhost/service-b/health/liveness

curl -v http://localhost:8080
curl -v http://localhost:8081

# Service a Logs
kubectl logs --follow `kubectl get pods -l app=envoy-playground-service-a -o=jsonpath='{.items[0].metadata.name}'` --container envoy
kubectl logs --follow `kubectl get pods -l app=envoy-playground-service-a -o=jsonpath='{.items[0].metadata.name}'` --container service
kubectl get event --field-selector involvedObject.name=`kubectl get pods -l app=envoy-playground-service-a -o=jsonpath='{.items[0].metadata.name}'`

kubectl logs --follow `kubectl get pods -l app=envoy-playground-service-b -o=jsonpath='{.items[0].metadata.name}'` --container envoy
kubectl logs --follow `kubectl get pods -l app=envoy-playground-service-b -o=jsonpath='{.items[0].metadata.name}'` --container service
```