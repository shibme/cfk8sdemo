# Cloudflare K8s Demo

A sample app that reads and writes data into K8s Config Map. This app is built as demo for K8s and Cloudflare integration.

## Steps to test CF with K8s

### Deploy app and cloudflared

- Create a namespace in your K8s cluster
```bash
kubectl create namespace cfk8sdemo
```
- Create secrets in your namespace
```bash
kubectl create secret -n cfk8sdemo generic api --from-literal=token=<api_token>
kubectl create secret -n cfk8sdemo generic cf --from-literal=token=<cloudflare_token>
```
- Create a test config map in your namespace
```bash
kubectl create configmap -n cfk8sdemo test --from-literal=hello=world
```
- Deploy the app and cloudflared in your namespace
```bash
kubectl apply -n cfk8sdemo -f k8s/deploy.yaml
```

### Configure tunnel in cloudflared

- Open cloudflare tunnels page and create a new tunnel
- Configure a public hostname with your domain and route traffic to `http://api-service:8888`
- Open the page `https://<your_domain>/api/map/test`
- You should see the content of the test config map