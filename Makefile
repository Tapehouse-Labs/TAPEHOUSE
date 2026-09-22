NODE_MAJOR := 24
FOUNDRY_VERSION := 1.8.3
FOUNDRY_VERSION_RE := $(subst .,\.,$(FOUNDRY_VERSION))
SLITHER_VERSION := 0.11.6
SLITHER_VERSION_RE := $(subst .,\.,$(SLITHER_VERSION))
COVERAGE_MIN := 95
SNAPSHOT_FILTER := --no-match-test '^(testFuzz|invariant|statefulFuzz)'

.PHONY: all build test lint coverage gas snapshot \
	check-toolchains check-node check-foundry check-slither submodules \
	build-apps test-apps lint-apps \
	build-contracts test-contracts lint-contracts coverage-contracts gas-contracts snapshot-contracts

all: build

build: build-contracts build-apps

test: test-contracts test-apps

lint: lint-contracts lint-apps

coverage: coverage-contracts

gas: gas-contracts

snapshot: snapshot-contracts

check-toolchains: check-node check-foundry check-slither

check-node:
	@node --version | grep -Eq '^v$(NODE_MAJOR)\.' || \
		{ echo "Node $(NODE_MAJOR).x is required, found $$(node --version)."; exit 1; }
	@command -v pnpm >/dev/null || \
		{ echo "pnpm is not installed. Run: corepack enable pnpm"; exit 1; }

check-foundry:
	@forge --version | head -n 1 | grep -Eq '^forge Version: $(FOUNDRY_VERSION_RE)(-|$$)' || \
		{ echo "Foundry $(FOUNDRY_VERSION) is required, found: $$(forge --version | head -n 1)."; \
		  echo "Run: foundryup --install v$(FOUNDRY_VERSION)"; exit 1; }

check-slither:
	@slither --version 2>/dev/null | grep -Eq '^$(SLITHER_VERSION_RE)$$' || \
		{ echo "Slither $(SLITHER_VERSION) is required, found: $$(slither --version 2>/dev/null || echo none)."; \
		  echo "Run: pipx install --force slither-analyzer==$(SLITHER_VERSION)"; exit 1; }

node_modules/.modules.yaml: package.json pnpm-workspace.yaml pnpm-lock.yaml $(wildcard apps/*/package.json)
	pnpm install --frozen-lockfile
	@touch $@

build-apps: check-node node_modules/.modules.yaml
	pnpm -r run build

test-apps: check-node node_modules/.modules.yaml
	pnpm -r run --if-present test

lint-apps: check-node node_modules/.modules.yaml
	pnpm -r run --if-present lint

submodules:
	git submodule update --init --recursive

build-contracts: check-foundry submodules
	cd contracts && forge build

test-contracts: check-foundry submodules
	cd contracts && forge test --gas-snapshot-emit false

lint-contracts: check-foundry check-slither submodules
	cd contracts && forge fmt --check
	cd contracts && forge lint --deny warnings
	@if [ -z "$$(find contracts/src -name '*.sol' 2>/dev/null)" ]; then \
		echo "No Solidity sources: Slither skipped."; \
	else cd contracts && slither .; fi

coverage-contracts: check-foundry submodules
	cd contracts && rm -f lcov.info && forge coverage --report summary --report lcov \
		--gas-snapshot-emit false --no-match-coverage '^(test|script)/'
	@if [ ! -f contracts/lcov.info ]; then echo "No Solidity sources: coverage skipped."; else \
		awk -F: -v min=$(COVERAGE_MIN) ' \
			/^LF:/ { lf += $$2 } /^LH:/ { lh += $$2 } /^BRF:/ { bf += $$2 } /^BRH:/ { bh += $$2 } \
			END { if (lf == 0) { print "No Solidity sources: coverage skipped."; exit 0 } \
				l = 100 * lh / lf; b = bf ? 100 * bh / bf : 100; \
				printf "Coverage: lines %.2f%% (%d/%d), branches %.2f%% (%d/%d), minimum %d%%\n", l, lh, lf, b, bh, bf, min; \
				if (l < min || b < min) exit 1 }' contracts/lcov.info; fi

gas-contracts: check-foundry submodules
	cd contracts && forge build --sizes
	cd contracts && FORGE_SNAPSHOT_CHECK=true forge snapshot --check $(SNAPSHOT_FILTER)
	@test -z "$$(git ls-files --others --exclude-standard -- contracts/snapshots)" || \
		{ echo "Untracked gas snapshots in contracts/snapshots. Run: make snapshot, then git add them."; exit 1; }
	@git diff --quiet -- contracts/.gas-snapshot contracts/snapshots || \
		{ echo "Gas snapshots changed and are not staged. Run: make snapshot, then git add them."; exit 1; }

snapshot-contracts: check-foundry submodules
	cd contracts && forge snapshot $(SNAPSHOT_FILTER)
