```bash

docker login ghcr.io

docker compose build
docker compose push


kind create cluster --name $(name)

kubectl create secret docker-registry github-registry \
	--docker-server=ghcr.io \
	--docker-username=${username} \
	--docker-password=${token} \
	--docker-email=${email}

kubectl apply -f k8s/secret.yml
kubectl apply -f k8s/configmap.yml

kubectl apply -f k8s/pg-pv.yml
kubectl apply -f k8s/pg-pvc.yml
kubectl apply -f k8s/pg.yml

kubectl apply -f k8s/app.yml


helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

helm template ingress-nginx ingress-nginx \
--repo https://kubernetes.github.io/ingress-nginx \
--version 4.12.0 \
--namespace ingress-nginx > ./k8s/ingress/controller/nginx/manifests/nginx-ingress.1.12.0.yml


kubectl create namespace ingress-nginx

kubectl apply -f ./k8s/ingress/controller/nginx/manifests/nginx-ingress.1.12.0.yml

sudo kubectl -n ingress-nginx port-forward svc/ingress-nginx-controller 443

kubectl apply -f ./k8s/ingress/controller/nginx/features/app.yml

```