.PHONY: all
all: build_test cli api

TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/podops

.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cli
cli:
	cd cmd/cli && go build -o po cli.go && mv po /Users/turing/devel/go/bin/po

.PHONY: build_test
build_test:
	cd cmd/cli && go build cli.go && rm cli
	cd cmd/api && go build main.go && rm main
	cd cmd/cdnapi && go build main.go && rm main
	cd cmd/cdn && go build main.go && rm main

.PHONY: test
test:
	go test
	cd apiv1 && go test
	cd internal/metadata && go test
	cd graphql && go test
	cd internal/platform && go test
	cd auth && go test
	
build_cdn: cmd/cdn/main.go
	cd cmd/cdn && ${TARGET_LINUX} go build -o svc main.go && docker build -t ${CONTAINER_REGISTRY}/cdn . && docker push ${CONTAINER_REGISTRY}/cdn

build_cdnapi: cmd/cdn/main.go
	cd cmd/cdnapi && ${TARGET_LINUX} go build -o svc main.go && docker build -t ${CONTAINER_REGISTRY}/cdnapi . && docker push ${CONTAINER_REGISTRY}/cdnapi
