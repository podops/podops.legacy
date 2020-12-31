.PHONY: all
all: cli

.PHONY: cli
cli:
	cd cmd/cli && go install po.go