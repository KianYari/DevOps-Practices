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


pg:
	kubectl apply -f k8s/pg-pv.yml
	kubectl apply -f k8s/pg-pvc.yml
	kubectl apply -f k8s/pg.yml

pg-del:
	kubectl delete -f k8s/pg.yml
	kubectl delete -f k8s/pg-pvc.yml
	kubectl delete -f k8s/pg-pv.yml


port-forward:
	kubectl port-forward service/app 8080:8080