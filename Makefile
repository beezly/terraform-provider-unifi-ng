# terraform-provider-unifi — OpenAPI-driven build pipeline
# Spec source: https://beez.ly/unifi-apis/

BINARY_NAME     = terraform-provider-unifi
SPEC_VERSION   ?= $(shell cat specs/latest-version.txt 2>/dev/null || echo "unknown")
RAW_SPEC        = specs/network-$(SPEC_VERSION).json
NORMALIZED_SPEC = specs/network-$(SPEC_VERSION)-normalized.json
PROVIDER_SPEC   = provider-code-spec.json
GENERATED_DIR   = internal/provider
GENERATOR_CFG   = generator_config.yml
GOBIN          ?= $(shell go env GOPATH)/bin

.PHONY: all fetch normalize generate build clean update check-tools

## all: Full pipeline — fetch, normalize, generate, build
all: fetch normalize generate build

## fetch: Download the latest UniFi OpenAPI spec
fetch:
	@bash scripts/fetch-spec.sh

## normalize: Fix schema names for codegen compatibility
normalize:
	@echo "Normalizing spec $(SPEC_VERSION)..."
	@python3 scripts/normalize-spec.py $(RAW_SPEC) $(NORMALIZED_SPEC)

## generate: Run codegen tools to produce Go source
generate: check-tools
	@echo "Generating provider code spec from OpenAPI..."
	$(GOBIN)/tfplugingen-openapi generate \
		--config $(GENERATOR_CFG) \
		--output $(PROVIDER_SPEC) \
		$(NORMALIZED_SPEC)
	@echo "Generating Go source from provider code spec..."
	$(GOBIN)/tfplugingen-framework generate all \
		--input $(PROVIDER_SPEC) \
		--output $(GENERATED_DIR)
	@echo "Generation complete."

## build: Compile the provider binary
build:
	@echo "Building $(BINARY_NAME)..."
	go mod tidy
	go build -o $(BINARY_NAME) .
	@echo "Built $(BINARY_NAME)"

## update: Full pipeline — fetch new spec (if any), regenerate, rebuild
update: fetch normalize generate build
	@echo "Update complete. Spec version: $(SPEC_VERSION)"

## check-tools: Verify codegen tools are installed
check-tools:
	@which $(GOBIN)/tfplugingen-openapi > /dev/null 2>&1 || \
		(echo "Installing tfplugingen-openapi..." && \
		 go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest)
	@which $(GOBIN)/tfplugingen-framework > /dev/null 2>&1 || \
		(echo "Installing tfplugingen-framework..." && \
		 go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest)

## clean: Remove generated files and binary
clean:
	rm -f $(BINARY_NAME)
	rm -f $(PROVIDER_SPEC)
	rm -f specs/*-normalized.json
	find $(GENERATED_DIR) -name "*_gen.go" -delete
	@echo "Cleaned generated files."

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //'
