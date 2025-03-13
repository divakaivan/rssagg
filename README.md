### RSS Aggregator API

- This project develops an API that allows users to authenticate, scrape RSS feeds, follow feeds of their choice, and view posts from the feeds they follow
- The API is tested, dockerised and available on [Docker Hub](https://hub.docker.com/repository/docker/timee98642/rss-agg-api/general)
- Then the API is used in a local Kubernetes application, ending with dashboards for monitoring Kubernetes and the API

Tech: Go, SQL, GitHub Actions, Docker, Kubernetes, Prometheus, Grafana

Below I go in more detail about each section.

### API 

![api-docs](project-info/api.svg)




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