VERSION ?= latest
REGISTRY ?= europe-docker.pkg.dev/srlinux/eu.gcr.io
PROJECT ?= mtls

# Image URL to use all building/pushing image targets
IMG ?= $(REGISTRY)/${PROJECT}:$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all :$(generate)

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate
generate: protoc-gen-go protoc-gen-gofast protoc-gen-go-grpc
	echo $(protoc-gen-go-grpc)
	echo $(LOCALBIN)
	protoc -I . $(shell find ./api/greeter/greeterpb -name '*.proto') --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/mtls main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool binaries
PROTOC_GO ?= $(LOCALBIN)/protoc-gen-go
PROTOC_GO_FAST ?= $(LOCALBIN)/protoc-gen-gofast
PROTOC_GO_GRPC ?= $(LOCALBIN)/protoc-gen-go-grpc

## Tool Versions
PROTOC_GO_VERSION ?= latest
PROTOC_GO_FAST_VERSION ?= latest
PROTOC_GO_GRPC_VERSION ?= latest

.PHONY: protoc-gen-go
protoc-gen-gofast: $(PROTOC_GO) ## Download protoc-gen-gofast locally if necessary.
$(PROTOC_GO): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-go || GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)


.PHONY: protoc-gen-gofast
protoc-gen-gofast: $(PROTOC_GO_FAST) ## Download protoc-gen-gofast locally if necessary.
$(PROTOC_GO_FAST): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-gofast || GOBIN=$(LOCALBIN) go install -v github.com/gogo/protobuf/protoc-gen-gofast@$(PROTOC_GO_FAST_VERSION)

.PHONY: protoc-gen-go-grpc
protoc-gen-gogrpc: $(PROTOC_GO_GRPC) ## Download protoc-gen-golang-grpc locally if necessary.
$(PROTOC_GO_GRPC): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-go-grpc || GOBIN=$(LOCALBIN) go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GO_GRPC_VERSION)
