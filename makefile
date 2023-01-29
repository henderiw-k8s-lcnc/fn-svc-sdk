.PHONY: all
all: generate

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate
generate: protoc-gen-gofast protoc-gen-go-grpc 
	protoc -I . $(shell find ./pkg/ -name '*.proto') --gofast_out=. --gofast_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: generate fmt vet  ## Run tests.
	 go test ./... -coverprofile cover.out

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
PROTOC_GO_FAST ?= $(LOCALBIN)/protoc-gen-gofast
PROTOC_GO_GRPC ?= $(LOCALBIN)/protoc-gen-go-grpc

## Tool Versions
PROTOC_GO_FAST_VERSION ?= latest
PROTOC_GO_GRPC_VERSION ?= latest

.PHONY: protoc-gen-gofast
protoc-gen-gofast: $(PROTOC_GO_FAST) ## Download protoc-gen-gofast locally if necessary.
$(PROTOC_GO_FAST): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-gofast || GOBIN=$(LOCALBIN) go install -v github.com/gogo/protobuf/protoc-gen-gofast@$(PROTOC_GO_FAST_VERSION)

.PHONY: protoc-gen-go-grpc
protoc-gen-gogrpc: $(PROTOC_GO_GRPC) ## Download protoc-gen-golang-grpc locally if necessary.
$(PROTOC_GO_GRPC): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-go-grpc || GOBIN=$(LOCALBIN) go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GO_GRPC_VERSION)
