# Envoy test

Port forwarding (if needed)
```shell

# Forward ports 
kubectl port-forward service/envoy-playground-service-a 8080:http
kubectl port-forward service/envoy-playground-service-b 8081:http

# or with kubepfm

kubepfm <<EOF
service/envoy-playground-service-a:8080:http
service/envoy-playground-service-b:8081:http
EOF
```

Debug cluster
```shell
kubectl run -it debug-pod --image=curlimages/curl --restart=Never -- sh

kubectl delete debug-pod


service_a_ip="`kubectl get pods -l=app=envoy-playground-service-a -o=jsonpath='{.items[0].status.podIP}'`"
service_b_ip="`kubectl get pods -l=app=envoy-playground-service-b -o=jsonpath='{.items[0].status.podIP}'`"

echo "curl -v http://${service_a_ip}:9901/stats/prometheus"
echo "curl -v http://${service_a_ip}:8081/metrics"
```