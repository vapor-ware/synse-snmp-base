#
# Synse SNMP Plugin Base
#

VERSION    := 0.1.1
BIN_NAME   := synse-snmp-base

GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2> /dev/null || true)
GIT_TAG    ?= $(shell git describe --tags 2> /dev/null || true)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%T 2> /dev/null)
GO_VERSION := $(shell go version | awk '{ print $$3 }')

PKG_CTX := github.com/vapor-ware/synse-sdk/sdk
LDFLAGS := -w \
	-X ${PKG_CTX}.BuildDate=${BUILD_DATE} \
	-X ${PKG_CTX}.GitCommit=${GIT_COMMIT} \
	-X ${PKG_CTX}.GitTag=${GIT_TAG} \
	-X ${PKG_CTX}.GoVersion=${GO_VERSION} \
	-X ${PKG_CTX}.PluginVersion=${PLUGIN_VERSION}


.PHONY: build
build:  ## Build the binary
	go build -ldflags "${LDFLAGS}" -o ${BIN_NAME}

.PHONY: build-linux
build-linux:  ## Build the binarry for linux amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} .

.PHONY: clean
clean:  ## Remove temporary files
	go clean -v
	rm -rf dist

.PHONY: cover
cover: unit-test  ## Run tests and open an HTML coverage report
	go tool cover -html=coverage.out

.PHONY: dep
dep:  ## Verify and tidy gomod dependencies
	go mod verify
	go mod tidy

.PHONY: fmt
fmt:  ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current version
	git tag -a ${VERSION} -m "${BIN_NAME} version ${VERSION}"
	git push -u origin ${VERSION}

.PHONY: lint
lint:  ## Lint project source files
	golint -set_exit_status ./pkg/...

.PHONY: test-all
test-all: unit-test integration-test  ## Run all project tests

.PHONY: unit-test
unit-test:  ## Run project unit tests
	@ # Note: this requires go1.10+ in order to do multi-package coverage reports
	@echo "--> Running Unit Tests"
	go test -v -short -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: integration-test
integration-test:  ## Run project integrationt tests
	@echo "--> Running Integration Tests"
	@./scripts/integration.sh

.PHONY: version
version:  ## Print the version of the plugin
	@echo "${VERSION}"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help


# Targets for Jenkins CI

.PHONY: test
test: unit-test

.PHONY: ci-integration-test
ci-integration-test:
	go test -run Integration -coverprofile=coverage.out -covermode=atomic ./...
