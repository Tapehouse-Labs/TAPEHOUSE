NODE_MAJOR := 24
FOUNDRY_VERSION := 1.8.3
FOUNDRY_VERSION_RE := $(subst .,\.,$(FOUNDRY_VERSION))
SLITHER_VERSION := 0.11.6
SLITHER_VERSION_RE := $(subst .,\.,$(SLITHER_VERSION))
CARGO_STYLUS_VERSION := 0.10.9
CARGO_STYLUS_VERSION_RE := $(subst .,\.,$(CARGO_STYLUS_VERSION))
NITRO_IMAGE := offchainlabs/nitro-node:v3.11.4-7d5ac27-slim-stripped
DEVNODE_RPC_URL := http://127.0.0.1:8547
DEVNODE_KEY := 0xb6b15c8cb491557369f3c7d2c287b053eb229daa9c22138887752191c9520659
DEVNODE_ACCOUNT := 0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E
COVERAGE_MIN := 95
STYLUS_MAX_COMPRESSED_BYTES := 24576
STYLUS_SIZE_RUSTFLAGS := --remap-path-prefix=$(or $(CARGO_HOME),$(HOME)/.cargo)/registry/src=/cargo
CONTRACT ?= band
RPC_URL_4663 = $(ROBINHOOD_RPC_URL)
RPC_URL_46630 = $(ROBINHOOD_TESTNET_RPC_URL)
RPC_URL_42161 = $(ARBITRUM_RPC_URL)
RPC_URL_412346 = $(DEVNODE_RPC_URL)
SNAPSHOT_FILTER := --no-match-test '^(testFuzz|invariant|statefulFuzz)'
ifeq ($(strip $(ROBINHOOD_RPC_URL)),)
ROBINHOOD_RPC_URL := https://robinhood.drpc.org
endif
ifeq ($(strip $(ROBINHOOD_TESTNET_RPC_URL)),)
ROBINHOOD_TESTNET_RPC_URL := https://robinhood-testnet.drpc.org
endif
ifeq ($(strip $(ARBITRUM_RPC_URL)),)
ARBITRUM_RPC_URL := https://arbitrum.gateway.tenderly.co
endif
ifeq ($(strip $(ROBINHOOD_LOGS_RPC_URL)),)
ROBINHOOD_LOGS_RPC_URL := https://rpc.mainnet.chain.robinhood.com
endif
export ROBINHOOD_RPC_URL ROBINHOOD_TESTNET_RPC_URL ARBITRUM_RPC_URL ROBINHOOD_LOGS_RPC_URL

.PHONY: all build test lint coverage gas snapshot \
	check-toolchains check-node check-foundry check-slither check-stylus check-docker submodules \
	build-apps test-apps lint-apps \
	build-contracts test-contracts lint-contracts coverage-contracts gas-contracts snapshot-contracts \
	build-stylus test-stylus lint-stylus gas-stylus snapshot-stylus check-activation initcode-stylus verify-stylus \
	devnode devnode-stop deploy-stylus-devnode test-stylus-devnode

all: build

build: build-contracts build-stylus build-apps

test: test-contracts test-stylus test-apps

lint: lint-contracts lint-stylus lint-apps

coverage: coverage-contracts

gas: gas-contracts gas-stylus

snapshot: snapshot-contracts snapshot-stylus

check-toolchains: check-node check-foundry check-slither check-stylus

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

check-stylus:
	@command -v rustup >/dev/null || \
		{ echo "rustup is required. Install it from https://rustup.rs"; exit 1; }
	@cargo stylus --version 2>/dev/null | grep -Eq '^stylus $(CARGO_STYLUS_VERSION_RE)$$' || \
		{ echo "cargo-stylus $(CARGO_STYLUS_VERSION) is required, found: $$(cargo stylus --version 2>/dev/null || echo none)."; \
		  echo "Run: cargo install --locked --force cargo-stylus@$(CARGO_STYLUS_VERSION)"; exit 1; }

check-docker:
	@docker info >/dev/null 2>&1 || \
		{ echo "Docker is required for the dev node and reproducible builds. Start Docker and retry."; exit 1; }

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

build-stylus: check-stylus
	cd stylus && cargo build --locked --release --target wasm32-unknown-unknown --lib

test-stylus: check-stylus
	cd stylus && cargo test --locked

lint-stylus: check-stylus
	cd stylus && cargo fmt --check
	cd stylus && cargo clippy --locked --all-targets -- -D warnings
	cd stylus && cargo clippy --locked --release --target wasm32-unknown-unknown --lib -- -D warnings

snapshot-stylus: check-stylus
	cd stylus && RUSTFLAGS='$(STYLUS_SIZE_RUSTFLAGS)' cargo build --locked --release \
		--target wasm32-unknown-unknown --lib --target-dir target/size
	@for wasm in stylus/target/size/wasm32-unknown-unknown/release/*.wasm; do \
		echo "$$(basename $$wasm .wasm) $$(wc -c < $$wasm | tr -d ' ')"; done > stylus/.wasm-size

gas-stylus: snapshot-stylus
	@for dir in stylus/contracts/*/; do \
		(cd $$dir && cargo stylus get-initcode --output ../../target/initcode.hex >/dev/null) || exit 1; \
		size=$$(($$(tr -d '\n' < stylus/target/initcode.hex | wc -c) / 2 - 43)); \
		echo "$$(basename $$dir): $$size bytes compressed, limit $(STYLUS_MAX_COMPRESSED_BYTES)"; \
		[ $$size -le $(STYLUS_MAX_COMPRESSED_BYTES) ] || exit 1; done
	@test -z "$$(git ls-files --others --exclude-standard -- stylus/.wasm-size)" || \
		{ echo "Untracked stylus/.wasm-size. Run: make snapshot-stylus, then git add it."; exit 1; }
	@git diff --quiet -- stylus/.wasm-size || \
		{ echo "WASM size changed and is not staged. Run: make snapshot-stylus, then git add it."; exit 1; }

check-activation: check-stylus
	@for rpc in $(ROBINHOOD_RPC_URL) $(ROBINHOOD_TESTNET_RPC_URL) $(ARBITRUM_RPC_URL); do \
		for dir in stylus/contracts/*/; do \
			(cd $$dir && cargo stylus check -e $$rpc) || exit 1; done; done

initcode-stylus: check-docker check-foundry
	CARGO_STYLUS_VERSION=$(CARGO_STYLUS_VERSION) stylus/scripts/reproducible.sh initcode

verify-stylus: check-docker
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(TX)" || \
		{ echo "Usage: make verify-stylus CHAIN=<4663|46630|42161|412346> TX=<deployment tx> [CONTRACT=band]"; exit 1; }
	CARGO_STYLUS_VERSION=$(CARGO_STYLUS_VERSION) stylus/scripts/reproducible.sh verify $(RPC_URL_$(CHAIN)) $(TX) $(CONTRACT)

devnode: check-docker check-foundry
	@docker rm -f tapehouse-devnode >/dev/null 2>&1 || true
	docker run -d --name tapehouse-devnode -p 127.0.0.1:8547:8547 $(NITRO_IMAGE) \
		--dev --http.addr 0.0.0.0 --http.api=net,web3,eth,debug
	@i=0; until cast chain-id --rpc-url $(DEVNODE_RPC_URL) >/dev/null 2>&1; do \
		i=$$((i + 1)); [ $$i -lt 150 ] || { docker logs --tail 20 tapehouse-devnode; \
		echo "The dev node did not answer on $(DEVNODE_RPC_URL)."; exit 1; }; \
		sleep 0.2; done
	@cast send --rpc-url $(DEVNODE_RPC_URL) --private-key $(DEVNODE_KEY) \
		0x0000000000000000000000000000000000000070 "scheduleArbOSUpgrade(uint64,uint64)" 61 0 >/dev/null
	@cast send --rpc-url $(DEVNODE_RPC_URL) --private-key $(DEVNODE_KEY) --value 0 $(DEVNODE_ACCOUNT) >/dev/null
	@test "$$(cast call --rpc-url $(DEVNODE_RPC_URL) 0x0000000000000000000000000000000000000064 'arbOSVersion()(uint256)')" = 116 || \
		{ echo "The dev node did not reach ArbOS 61."; exit 1; }
	@stylus/scripts/deploy-stylus-deployer.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY)
	@echo "Dev node ready on $(DEVNODE_RPC_URL): chain $$(cast chain-id --rpc-url $(DEVNODE_RPC_URL)), ArbOS 61."

devnode-stop:
	docker rm -f tapehouse-devnode

deploy-stylus-devnode: check-stylus check-foundry submodules
	stylus/scripts/devnode-deploy.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY)

test-stylus-devnode: check-stylus check-foundry submodules
	stylus/scripts/devnode-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) \
		$$(cat stylus/target/devnode-band) stylus/target/devnode-registry.json
