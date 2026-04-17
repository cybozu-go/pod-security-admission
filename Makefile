ENVTEST_K8S_VERSION = 1.34.x

# Set the shell used to bash for better error handling.
SHELL = /bin/bash
.SHELLFLAGS = -e -o pipefail -c

STATICCHECK = go tool staticcheck
SETUP_ENVTEST = go tool setup-envtest

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
manifests: setup ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	controller-gen rbac:roleName=pod-security-admission webhook paths="./..."

.PHONY: generate
generate: setup ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: check-generate
check-generate:
	$(MAKE) manifests
	$(MAKE) generate
	git diff --exit-code --name-only

.PHONY: lint
lint:
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	$(STATICCHECK) ./...
	go vet ./...

.PHONY: test
test: manifests generate ## Run tests.
	source <($(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) -p env); \
		go test -v -count 1 -race ./... -ginkgo.v -ginkgo.fail-fast

##@ Build

.PHONY: build
build: ## Build binary.
	CGO_ENABLED=0 go build -o bin/pod-security-admission -ldflags="-w -s" ./cmd

.PHONY: setup
setup:
	aqua install --only-link

.PHONY: clean
clean:
	rm -rf bin
	rm -f config/crd/bases/*
