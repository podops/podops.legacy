.PHONY: all
all: api cdn cli

.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cdn
cdn:
	cd cmd/cdn && gcloud app deploy . --quiet

.PHONY: cli
cli:
	cd cmd/cli && go install po.go