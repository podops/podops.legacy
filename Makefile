.PHONY: all
all: build_test cli web cdn api

VERSION_TAG = 0.9.7
PLATFORM_LINUX = GOARCH=amd64 GOOS=linux
PLATFORM_MAC = GOARCH=amd64 GOOS=darwin
PLATFORM_WINDOWS = GOARCH=amd64 GOOS=windows

.PHONY: build_test
build_test:
	cd cmd/cli && go build cli.go && rm cli
	cd cmd/api && go build main.go && rm main
	cd cmd/cdn && go build main.go && rm main
	#cd examples/simple && go build main.go && rm main

.PHONY: test
test:
	cd apiv1 && go test
	cd pkg/auth && go test
	
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

.PHONY: cli
cli:
	cd cmd/cli && go build -o po cli.go && mv po /Users/turing/devel/go/bin/po
	
#.PHONY: cli
#cli:
#	rm -f build/po-*
#	cd cmd/cli && ${PLATFORM_LINUX} go build -o ../../build/po-linux-${VERSION_TAG} po.go
#	cd cmd/cli && ${PLATFORM_MAC} go build -o ../../build/po-mac-${VERSION_TAG} po.go
#	cd cmd/cli && ${PLATFORM_WINDOWS} go build -o ../../build/po-windows-${VERSION_TAG} po.go
#	cd build && chmod 755 po* && gzip po-linux-${VERSION_TAG} && gzip po-mac-${VERSION_TAG} && gzip po-windows-${VERSION_TAG}
#	gsutil rsync build gs://cdn.podops.dev/downloads