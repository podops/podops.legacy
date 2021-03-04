.PHONY: all
all: web cli

.PHONY: build_test
build_test:
	go mod verify && go mod tidy
	cd cmd/cli && go build po.go && rm po
	cd cmd/podops-cdn && go build main.go && rm main
	cd examples/simple && go build main.go && rm main

.PHONY: web
web:
	cd ../podops.dev && gridsome build
	rm -rf cmd/podops-cdn/public
	cp -R ../podops.dev/dist cmd/podops-cdn/public

.PHONY: cli
cli:
	cd cmd/cli && go mod verify && go mod tidy && go install po.go