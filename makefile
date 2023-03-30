SHELL := /bin/bash

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

VERSION := 0.1

all: service

service:
	docker build -f config/docker/dockerfile -t test-service-amd64:$(VERSION) \
			--build-arg BUILD_REF=$(VERSION) \
			--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%S"Z` \
			.
