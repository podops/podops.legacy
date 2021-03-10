.PHONY: all
all: build_test local cli web cdn api

VERSION_TAG = 0.9.6
PLATFORM_LINUX = GOARCH=amd64 GOOS=linux
PLATFORM_MAC = GOARCH=amd64 GOOS=darwin
PLATFORM_WINDOWS = GOARCH=amd64 GOOS=windows

.PHONY: build_test
build_test:
	cd cmd/cli && go build po.go && rm po
	cd cmd/api && go build main.go && rm main
	cd cmd/cdn && go build main.go && rm main
	cd examples/simple && go build main.go && rm main

.PHONY: web
web:
	cd ../podops.dev && gridsome build
	rm -rf cmd/cdn/public
	cp -R ../podops.dev/dist cmd/cdn/public

.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cdn
cdn:
	cd cmd/cdn && gcloud app deploy . --quiet

.PHONY: local
local:
	cd cmd/cli && go mod verify && go mod tidy && go install po.go

.PHONY: cli
cli:
	rm -f build/po-*
	cd cmd/cli && ${PLATFORM_LINUX} go build -o ../../build/po-linux-${VERSION_TAG} po.go
	cd cmd/cli && ${PLATFORM_MAC} go build -o ../../build/po-mac-${VERSION_TAG} po.go
	cd cmd/cli && ${PLATFORM_WINDOWS} go build -o ../../build/po-windows-${VERSION_TAG} po.go
	cp build/po-linux-${VERSION_TAG} build/po-linux
	cp build/po-mac-${VERSION_TAG} build/po-mac
	cp build/po-windows-${VERSION_TAG} build/po-windows
	gsutil rsync build gs://cdn.podops.dev/downloads