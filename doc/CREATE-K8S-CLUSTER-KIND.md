# Create a K8S Cluster With KindD

## Pre Prerequisites

For examples docker + KinD

[KinD Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager)

The commands in MacOS would be:
```shell
brew cask install docker
brew install kind@0.11.1
```

## Create the cluster
```shell
export CLUSTER_NAME="envoy-test";
cat <<EOF | kind create cluster --name ${CLUSTER_NAME}  --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
# Install ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
# and wait
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# check connectivity
kubectl cluster-info --context kind-envoy-test

docker port envoy-test-control-plane
# 6443/tcp -> 127.0.0.1:55313

kind get clusters
kubectl config use-context kind-envoy-test
kubectl cluster-info
```

Cluster deletion

```shell
kind delete cluster --name=envoy-test
# or
docker stop envoy-test-control-plane
docker rm envoy-test-control-plane

# check there's no "envoy-test-control-plane" container
docker container ls
```