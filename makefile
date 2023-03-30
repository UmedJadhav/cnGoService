SHELL := /bin/bash

tidy:
	go mod tidy
	go mod vendor

run:
	go run app/services/sales-api/main.go

build:
	go build -ldflags "-X main.build=local"

VERSION := 1.0

all: sales-api

sales-api:
	docker build -f config/docker/dockerfile.sales-api -t sales-api-amd64:$(VERSION) \
			--build-arg BUILD_REF=$(VERSION) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%S"Z` \
			.

KIND_CLUSTER := test-api-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1 \
		--name $(KIND_CLUSTER) \
		--config config/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd config/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build config/k8s/kind/sales-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment sales-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-status-service:
	kubectl get pods -o wide --watch

kind-describe:
	kubectl describe node
	kubectl describe svc
	kubectl describe pod -l app=sales

kind-describe-deployment:
	kubectl describe deployment sales-pod

kind-describe-replicaset:
	kubectl get rs
	kubectl describe rs -l app=sales

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
