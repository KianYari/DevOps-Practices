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