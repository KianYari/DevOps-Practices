#### Create Cluster
```bash
kind create cluster --name $(name)
```
#### Create ghcr secret inorder to pull form private registry inside k8s cluster
```bash
kubectl create secret docker-registry github-registry \
	--docker-server=ghcr.io \
	--docker-username=${username} \
	--docker-password=${token} \
	--docker-email=${email}
```
#### Create secrets and configmap
```bash
kubectl apply -f k8s/secret.yml
kubectl apply -f k8s/configmap.yml
```
#### Setup Database
```bash
kubectl apply -f k8s/pg-pv.yml
kubectl apply -f k8s/pg-pvc.yml
kubectl apply -f k8s/pg.yml
```
#### Build and Setup deployment
```bash
docker build -t ghcr.io/${username}/app:latest
```
```bash
kubectl apply -f k8s/app.yml
```
#### Setup Ingress and Nginx Ingress Controller
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
```
```bash
helm template ingress-nginx ingress-nginx \
--repo https://kubernetes.github.io/ingress-nginx \
--version 4.12.0 \
--namespace ingress-nginx > ./k8s/ingress/controller/nginx/manifests/nginx-ingress.1.12.0.yml
```
```bash
kubectl create namespace ingress-nginx
kubectl apply -f ./k8s/ingress/controller/nginx/manifests/nginx-ingress.1.12.0.yml
```
###### Setup port forwarding for Nginx service
```bash
sudo kubectl -n ingress-nginx port-forward svc/ingress-nginx-controller 443
```
###### Setup Ingress for app
```bash
kubectl apply -f ./k8s/ingress/controller/nginx/features/app.yml
```


## using helm
```bash
helm create app
```
```bash
helm install ${releaseName} ${path} --values=${path}
```
```bash
helm upgrade ${releaseName} ${path} --values=${path}
```
```bash
helm package app/
```
```bash
helm push app-version oci://ghcr.io/$(username)
```
