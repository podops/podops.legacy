.PHONY: all
all: cli

.PHONY: cli
cli:
	cd cli && go mod verify && go mod tidy && go install po.go