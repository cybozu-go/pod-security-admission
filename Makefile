CONTROLLER_TOOLS_VERSION = 0.7.0
KUSTOMIZE_VERSION = 4.4.1
ENVTEST_K8S_VERSION = 1.22.1

# Set the shell used to bash for better error handling.
SHELL = /bin/bash
.SHELLFLAGS = -e -o pipefail -c
BIN_DIR = ./bin
INSTALL_YAML = build/install.yaml

KUSTOMIZE = $(BIN_DIR)/kustomize
CONTROLLER_GEN = $(BIN_DIR)/controller-gen
STATICCHECK = $(BIN_DIR)/staticcheck
NILERR = $(BIN_DIR)/nilerr
SETUP_ENVTEST = $(BIN_DIR)/setup-envtest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: $(CONTROLLER_GEN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=pod-security-admission webhook paths="./..."

.PHONY: generate
generate: $(CONTROLLER_GEN) ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: check-generate
check-generate:
	$(MAKE) manifests
	$(MAKE) generate
	git diff --exit-code --name-only

.PHONY: lint
lint: $(STATICCHECK) $(NILERR)
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	$(STATICCHECK) ./...
	go vet ./...

.PHONY: test
test: manifests generate ## Run tests.
	{ \
	source <($(SETUP_ENVTEST) use -p env $(ENVTEST_K8S_VERSION)) && \
	go test -v -count=1 ./... -coverprofile cover.out ; \
	}

##@ Build

.PHONY: build
build: ## Build binary.
	CGO_ENABLED=0 go build -o bin/pod-security-admission -ldflags="-w -s" main.go

$(INSTALL_YAML): $(KUSTOMIZE)
	mkdir -p build
	$(KUSTOMIZE) build ./config/default > $@

$(CONTROLLER_GEN): ## Download controller-gen locally if necessary.
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_TOOLS_VERSION))

$(KUSTOMIZE): ## Download kustomize locally if necessary.
	mkdir -p $(BIN_DIR)
	curl -sSLf https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv$(KUSTOMIZE_VERSION)/kustomize_v$(KUSTOMIZE_VERSION)_linux_amd64.tar.gz | tar -xz -C $(BIN_DIR)

$(STATICCHECK):
	$(call go-install-tool,$(STATICCHECK),honnef.co/go/tools/cmd/staticcheck@latest)

$(NILERR):
	$(call go-install-tool,$(NILERR),github.com/gostaticanalysis/nilerr/cmd/nilerr@latest)

$(SETUP_ENVTEST):
	$(call go-install-tool,$(SETUP_ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

.PHONY: setup
setup: $(STATICCHECK) $(NILERR) $(KUSTOMIZE) $(CONTROLLER_GEN) $(SETUP_ENVTEST)

.PHONY: clean
clean:
	rm -rf bin
	rm -f config/crd/bases/*

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
