NODE_MAJOR := 24
FOUNDRY_VERSION := 1.8.3
FOUNDRY_VERSION_RE := $(subst .,\.,$(FOUNDRY_VERSION))

.PHONY: all build test check-toolchains check-node check-foundry \
	build-apps test-apps build-contracts test-contracts

all: build

build: build-contracts build-apps

test: test-contracts test-apps

check-toolchains: check-node check-foundry

check-node:
	@node --version | grep -Eq '^v$(NODE_MAJOR)\.' || \
		{ echo "Node $(NODE_MAJOR).x is required, found $$(node --version)."; exit 1; }
	@command -v pnpm >/dev/null || \
		{ echo "pnpm is not installed. Run: corepack enable pnpm"; exit 1; }

check-foundry:
	@forge --version | head -n 1 | grep -Eq '^forge Version: $(FOUNDRY_VERSION_RE)(-|$$)' || \
		{ echo "Foundry $(FOUNDRY_VERSION) is required, found: $$(forge --version | head -n 1)."; \
		  echo "Run: foundryup --install v$(FOUNDRY_VERSION)"; exit 1; }

node_modules/.modules.yaml: package.json pnpm-workspace.yaml pnpm-lock.yaml $(wildcard apps/*/package.json)
	pnpm install --frozen-lockfile
	@touch $@

build-apps: check-node node_modules/.modules.yaml
	pnpm -r run build

test-apps: check-node node_modules/.modules.yaml
	pnpm -r run --if-present test

build-contracts: check-foundry
	git submodule update --init --recursive
	cd contracts && forge build

test-contracts: check-foundry
	git submodule update --init --recursive
	cd contracts && forge test
