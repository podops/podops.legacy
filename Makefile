.PHONY: all
all: build_test cli cdn api

PLATFORM_LINUX = GOARCH=amd64 GOOS=linux
PLATFORM_MAC = GOARCH=amd64 GOOS=darwin
PLATFORM_WINDOWS = GOARCH=amd64 GOOS=windows

.PHONY: build_test
build_test:
	cd cmd/cli && go build cli.go && rm cli
	cd cmd/api && go build main.go && rm main
	cd cmd/cdn && go build main.go && rm main

.PHONY: test
test:
	go test
	cd internal/api && go test
	cd pkg/auth && go test
	
.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cdn
cdn:
	cd cmd/cdn && gcloud app deploy . --quiet

.PHONY: cli
cli:
	cd cmd/cli && go build -o po cli.go && mv po /Users/turing/devel/go/bin/po
