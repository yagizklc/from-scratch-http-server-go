
# Binary name
BINARY_NAME=main

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOBINPATH=$(GOBIN)/$(BINARY_NAME)
GOFILES=$(GOBASE)/app/*.go

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean run air deps vet fmt test integration builder doctor integration-doctor

all: build

build: deps
	@echo "Building..."
	@go build -o $(GOBINPATH) $(GOFILES)

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf $(GOBIN)

run: build
	@echo "Running..."
	@$(GOBINPATH)

air: build
	@air --build.cmd "go build -o $(GOBINPATH) $(GOFILES)" --build.bin "$(GOBINPATH)"

deps:
	@echo "Ensuring dependencies are up to date..."
	@go mod tidy

vet:
	@echo "Running go vet..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

test:
	@echo "Running tests..."
	@go test ./...

integration:
	@echo "Running integration tests..."
	@go test -tags integration ./...

doctor: fmt vet test build
	@echo "All checks passed!"

integration-doctor: fmt vet test integration build
	@echo "All checks and integration tests passed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  clean            - Remove binary and clear cache"
	@echo "  run              - Build and run the application"
	@echo "  air              - Deamon running application"
	@echo "  deps             - Ensure dependencies are up to date"
	@echo "  vet              - Run go vet"
	@echo "  fmt              - Format the code"
	@echo "  test             - Run unit tests"
	@echo "  integration      - Run integration tests"
	@echo "  doctor           - Run vet, test, and build"
	@echo "  integration-doctor - Run vet, test, integration tests, and build"
	@echo "  help             - Print this help message"