# Build configuration
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
SKIP_INSTALL := false
BINARY_NAME := devspace
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_TREE_STATE := $(if $(shell git status --porcelain),dirty,clean)
LDFLAGS := -X main.version=$(VERSION) \
           -X main.commit=$(GIT_COMMIT) \
           -X main.date=$(BUILD_DATE) \
           -X main.treeState=$(GIT_TREE_STATE)

# Platform configuration
PLATFORM_HOST ?= localhost:8080
NAMESPACE ?= khulnasoft

# Help message
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: install
deps: ## Install dependencies
	go mod download
	cd desktop && yarn install --frozen-lockfile

.PHONY: build
build: deps ## Build the CLI and Desktop
	SKIP_INSTALL=$(SKIP_INSTALL) BUILD_PLATFORMS=$(GOOS) BUILD_ARCHS=$(GOARCH) ./hack/rebuild.sh

.PHONY: run-desktop
run-desktop: build ## Run the desktop app
	cd desktop && yarn desktop:dev:debug

.PHONY: run-daemon
run-daemon: build ## Run the daemon against khulnasoft host
	$(BINARY_NAME) pro daemon start --host $(PLATFORM_HOST)

##@ Testing

.PHONY: test
TEST_PACKAGES ?= ./...
test: ## Run tests
	go test -v -coverprofile=coverage.out $(TEST_PACKAGES)

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html

##@ Linting

.PHONY: lint
lint: ## Run linters
	golangci-lint run ./...

##@ Build

.PHONY: clean
clean: ## Clean build artifacts
	go clean
	rm -rf ./bin
	rm -rf ./coverage.*

##@ Deployment

.PHONY: cp-to-platform
cp-to-platform: ## Copy the devspace binary to the platform pod
	@echo "Building for linux/$(GOARCH)..."
	@SKIP_INSTALL=true BUILD_PLATFORMS=linux BUILD_ARCHS=$(GOARCH) ./hack/rebuild.sh
	@POD=$$(kubectl get pod -n $(NAMESPACE) -l app=khulnasoft,release=khulnasoft -o jsonpath='{.items[0].metadata.name}'); \
	echo "Copying ./test/devspace-linux-$(GOARCH) to pod $$POD"; \
	kubectl cp -n $(NAMESPACE) ./test/devspace-linux-$(GOARCH) $$POD:/usr/local/bin/$(BINARY_NAME)

##@ Release

.PHONY: release
release: ## Create a new release
	@echo "Creating release $(VERSION)"
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

##@ Helpers

.PHONY: version
version: ## Display version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Default target
.DEFAULT_GOAL := help
