build:
	docker build -t ghcr.io/kianyari/$(n) .

push:
	docker push ghcr.io/kianyari/$(n)

up:
	docker compose down
	docker compose up --build -d
	docker image prune -f

ps:
	docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"




base64:
		echo -n "$(t)" | base64





create-cluster:
	kind create cluster --name $(n)

apply:
	kubectl apply -f k8s/${f}
des:
	kubectl describe $(r)
all:
	kubectl get all
	kubectl get pv
	kubectl get pvc
	kubectl get storageclass
	kubectl get secret

github-registry:
	kubectl create secret docker-registry github-registry \
	--docker-server=ghcr.io \
	--docker-username=${username} \
	--docker-password=${token} \
	--docker-email=${email}

del:
	kubectl delete ${r}

restart:
	kubectl rollout restart deployment $(d)

port-forward:
	kubectl port-forward service/app 8080:8080

install-nginx-ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0/deploy/static/provider/cloud/deploy.yaml

wait-for-nginx-ingress:
	kubectl wait --namespace ingress-nginx \
  	--for=condition=ready pod \
 	--selector=app.kubernetes.io/component=controller \
  	--timeout=120s

clean-cluster:
	kubectl delete all --all
	kubectl delete pvc --all
	kubectl delete pv --all
	kubectl delete storageclass --all
	kubectl get secret -o name | grep -v "github-registry" | xargs kubectl delete
add-nginx-ingress-to-helm:
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

CHART_VERSION="4.12.0"
APP_VERSION="1.12.0"

get-nginx-template:
	helm template ingress-nginx ingress-nginx \
	--repo https://kubernetes.github.io/ingress-nginx \
	--version ${CHART_VERSION} \
	--namespace ingress-nginx \
	> ./k8s/ingress/controller/nginx/manifests/nginx-ingress.${APP_VERSION}.yml

start-nginx-ingress:
	kubectl create namespace ingress-nginx
	kubectl apply -f ./k8s/ingress/controller/nginx/manifests/nginx-ingress.${APP_VERSION}.yml

port-forward-nginx:
	kubectl -n ingress-nginx port-forward svc/ingress-nginx-controller 443

