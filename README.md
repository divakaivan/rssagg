

### Setup Kubernetes

```bash
kubectl apply -f kubernetes/.

# check everything is good 
kubectl get all

# expose API to local
kubectl port-forward api-pod-name 8080
```

### Setup Monitoring

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

helm repo update

helm install my-prometheus prometheus-community/prometheus

helm repo add grafana https://grafana.github.io/helm-charts

helm repo update

helm install my-grafana grafana-grafana

# check pods/svc with kubectl are running...

# get grafana pwd and save it for login
kubectl get secret -n default my-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# port-forward grafana to localhost
kubectl port-forward service/my-grafana 3000:80 --address='0.0.0.0' 

# visit 0.0.0.0:3000 and use admin/pwd
```

### Future

- send API logs to a dedicated database
- use an ingress controller
- implement the repository pattern
- add caching 
- 