.PHONY: all
all: cli web cdn

.PHONY: build_test
build_test:
	go mod verify && go mod tidy
	cd cmd/cli && go build po.go && rm po
	cd cmd/podops-cdn && go build main.go && rm main
	cd examples/simple && go build main.go && rm main

.PHONY: web
web:
	cd ../podops.dev && gridsome build
	rm -rf cmd/cdn/public
	cp -R ../podops.dev/dist cmd/cdn/public

.PHONY: cli
cli:
	cd cmd/cli && go mod verify && go mod tidy && go install po.go

.PHONY: cdn
cdn:
	cd cmd/cdn && gcloud app deploy . --quiet