.PHONY: all
all: cli

.PHONY: cli
cli:
	cd cmd/cli && go mod verify && go mod tidy && go install po.go