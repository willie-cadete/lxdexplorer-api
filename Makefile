SEMVER := `git describe --tags --abbrev=0 2> /dev/null || git rev-parse --short HEAD`
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

.PHONY: all lint race msan test build docker help

all: help

fmt: ## Format code
	@go fmt $(go list ./... | grep -v /vendor/)

lint: ## Lint the files
	@golangci-lint run

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

build: ## Build the binary file without specify OS
	@go build -o dist/ ./...

vet: ## Run go static analysis tool
	@go vet ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

clean: ## Remove previous build
	@rm -rvf dist/*

docker:
	@docker build -t lxdexplorer-api:$(SEMVER) .

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
