TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/podops

.PHONY: all
all: build_test cli api

.PHONY: deploy_all
deploy_all: build_test cli build_cdn build_cdnapi api deploy_services

.PHONY: deploy_services
deploy_services:
	cd ../podops-infra && ansible-playbook -i inventory/podops.dev.yml playbooks/update_services.yml

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

.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v cmd` -coverprofile=coverage.txt -covermode=atomic

.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cli
cli:
	cd cmd/cli && go build -o po cli.go && mv po /Users/turing/devel/go/bin/po
	
.PHONY: build_cdn
build_cdn:
	cd cmd/cdn && ${TARGET_LINUX} go build -o svc main.go && docker build -t ${CONTAINER_REGISTRY}/cdn . && docker push ${CONTAINER_REGISTRY}/cdn

.PHONY: build_cdnapi
build_cdnapi:
	cd cmd/cdnapi && ${TARGET_LINUX} go build -o svc main.go && docker build -t ${CONTAINER_REGISTRY}/cdnapi . && docker push ${CONTAINER_REGISTRY}/cdnapi
