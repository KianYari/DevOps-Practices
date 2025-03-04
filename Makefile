create-cluster:
	kind create cluster --name $(n)
apply:
	kubectl apply -f k8s/${f}

build:
	docker build -t ghcr.io/kianyari/$(n) .

push:
	docker push ghcr.io/kianyari/$(n)

base64:
		echo -n $(t) | base64

des:
	kubectl describe $(r)


# this creates the secret nedded to pull the image from ghcr
ghcr-secret:
	kubectl create secret docker-registry github-registry-secret \
  --docker-server=ghcr.io \
  --docker-username=$(username) \
  --docker-password=$(token) \
  --dry-run=client -o yaml > k8s/ghcr-secret.yml