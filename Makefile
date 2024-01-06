## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

WIRE_TAGS = "oss"

-include local/Makefile

GO = go
GO_FILES ?= ./pkg/...
SH_FILES ?= $(shell find ./scripts -name *.sh)
GO_BUILD_FLAGS += $(if $(GO_BUILD_DEV),-dev)
GO_BUILD_FLAGS += $(if $(GO_BUILD_TAGS),-build-tags=$(GO_BUILD_TAGS))

targets := $(shell echo '$(sources)' | tr "," " ")

$(NGALERT_SPEC_TARGET):
	+$(MAKE) -C pkg/services/ngalert/api/tooling api.json

gen-go: $(WIRE)
	@echo "generate go files"
	/Users/rlaisqls/go/bin/wire-v0.5.0 gen -tags $(WIRE_TAGS) ./pkg/server

build: ## Build all Go binaries.
	@echo "build go files"
	$(GO) run build.go $(GO_BUILD_FLAGS) build

run: $(BRA) ## Build and run web server on filesystem changes.
	$(BRA) run

##@ Linting
golangci-lint: $(GOLANGCI_LINT)
	@echo "lint via golangci-lint"
	$(GOLANGCI_LINT) run \
		--config .golangci.toml \
		$(GO_FILES)

lint-go: golangci-lint ## Run all code checks for backend. You can use GO_FILES to specify exact files to check