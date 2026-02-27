.PHONY: help build install dev check test fmt vet lint clean run

# Default target
.DEFAULT_GOAL := help

##@ Development

dev: build install ## Build and install locally (use after code changes)

run: ## Run without building
	go run .

##@ Build

build: ## Build binary to ./p
	go build -o p .

install: ## Install to $GOPATH/bin
	go install .

##@ Quality

check: fmt vet test ## Run all checks (fmt, vet, test)

test: ## Run tests
	go test ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run

##@ Cleanup

clean: ## Remove build artifacts
	rm -f p

##@ Help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
