DOCKER_REGISTRY ?= localhost
DOCKER_NAMESPACE ?= vlad
VERSION ?= dev

PROJECT ?= simple-api
IMG_NAME ?= $(DOCKER_REGISTRY)/$(DOCKER_NAMESPACE)/$(PROJECT)
IMG ?= $(IMG_NAME):$(VERSION)

.PHONY: docker-build
docker-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/simple-api main.go
	docker build -t ${IMG} .