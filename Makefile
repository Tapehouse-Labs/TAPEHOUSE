# SPDX-License-Identifier: MIT OR Apache-2.0
NODE_MAJOR := 24
FOUNDRY_VERSION := 1.8.3
FOUNDRY_VERSION_RE := $(subst .,\.,$(FOUNDRY_VERSION))
SLITHER_VERSION := 0.11.6
SLITHER_VERSION_RE := $(subst .,\.,$(SLITHER_VERSION))
REUSE_VERSION := 6.2.0
REUSE_VERSION_RE := $(subst .,\.,$(REUSE_VERSION))
CARGO_STYLUS_VERSION := 0.10.9
CARGO_STYLUS_VERSION_RE := $(subst .,\.,$(CARGO_STYLUS_VERSION))
BINARYEN_VERSION := $(shell awk '/^\[/ { table = $$0 } table == "[wasm-opt]" && $$1 == "version" { gsub(/"/, "", $$3); print $$3 }' stylus/Stylus.toml)
BINARYEN_HOME := $(or $(XDG_CACHE_HOME),$(HOME)/.cache)/binaryen/version_$(BINARYEN_VERSION)
export PATH := $(BINARYEN_HOME)/bin:$(PATH)
NITRO_IMAGE := offchainlabs/nitro-node:v3.11.4-7d5ac27-slim-stripped
DEVNODE_RPC_URL := http://127.0.0.1:8547
DEVNODE_KEY := 0xb6b15c8cb491557369f3c7d2c287b053eb229daa9c22138887752191c9520659
DEVNODE_ACCOUNT := 0x3f1Eae7D46d88F08fc2F8ed27FCb2AB183EB2d0E
COVERAGE_MIN := 95
STYLUS_MAX_FRAGMENTS := 4
DEVNODE_GAS_TOLERANCE_BPS := 50
STYLUS_SIZE_RUSTFLAGS := --remap-path-prefix=$(or $(CARGO_HOME),$(HOME)/.cargo)/registry/src=/cargo
CONTRACT ?= band
RPC_URL_4663 = $(ROBINHOOD_RPC_URL)
RPC_URL_46630 = $(ROBINHOOD_TESTNET_RPC_URL)
RPC_URL_42161 = $(ARBITRUM_RPC_URL)
RPC_URL_412346 = $(DEVNODE_RPC_URL)
ROBINHOOD_FORK_BLOCK := 69922505
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
	check-toolchains check-node check-foundry check-slither check-reuse check-stylus check-docker submodules \
	build-apps test-apps lint-apps \
	build-contracts test-contracts lint-contracts coverage-contracts gas-contracts snapshot-contracts \
	build-stylus test-stylus lint-stylus gas-stylus snapshot-stylus check-activation deploy-stylus verify-stylus \
	devnode devnode-stop deploy-stylus-devnode test-stylus-devnode gas-stylus-devnode snapshot-stylus-devnode gas-table \
	deploy-band-feeds verify-band-feeds simulate-supply-vault deploy-supply-vault verify-supply-vault \
	deploy-margin-accounts verify-margin-accounts deploy-liquidator verify-liquidator deploy-contracts-devnode test-contracts-devnode lint-scripts lint-licenses

all: build

build: build-contracts build-stylus build-apps

test: test-contracts test-stylus test-apps

lint: lint-contracts lint-stylus lint-apps lint-scripts lint-licenses

coverage: coverage-contracts

gas: gas-contracts gas-stylus

snapshot: snapshot-contracts snapshot-stylus

check-toolchains: check-node check-foundry check-slither check-reuse check-stylus

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

check-reuse:
	@reuse --version 2>/dev/null | head -n 1 | grep -Eq '^reuse, version $(REUSE_VERSION_RE)$$' || \
		{ echo "REUSE $(REUSE_VERSION) is required, found: $$(reuse --version 2>/dev/null | head -n 1 || echo none)."; \
		  echo "Run: pipx install --force 'reuse[charset-normalizer]==$(REUSE_VERSION)'"; exit 1; }

check-stylus:
	@command -v rustup >/dev/null || \
		{ echo "rustup is required. Install it from https://rustup.rs"; exit 1; }
	@cargo stylus --version 2>/dev/null | grep -Eq '^stylus $(CARGO_STYLUS_VERSION_RE)$$' || \
		{ echo "cargo-stylus $(CARGO_STYLUS_VERSION) is required, found: $$(cargo stylus --version 2>/dev/null || echo none)."; \
		  echo "Run: cargo install --locked --force cargo-stylus@$(CARGO_STYLUS_VERSION)"; exit 1; }
	@stylus/scripts/binaryen.sh $(BINARYEN_VERSION) $(BINARYEN_HOME)

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
	@git submodule update --init --recursive

lint-scripts:
	@shellcheck contracts/script/*.sh stylus/scripts/*.sh

lint-licenses: check-reuse
	reuse lint

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
		set -- $$(cd $$dir && cargo stylus check -e $(ROBINHOOD_RPC_URL) 2>&1 | \
			perl -ne 's/\e\[[0-9;]*m//g; /contract size: .*\((\d+) bytes\)(?: \((\d+) fragments\))?/ and print "$$1 ", $$2 // 1, "\n"'); \
		[ $$# -eq 2 ] || { echo "cargo stylus check on ROBINHOOD_RPC_URL did not report the size of $$dir"; exit 1; }; \
		size=$$1 fragments=$$2; \
		echo "$$(basename $$dir): $$size bytes compressed in $$fragments fragments, limit $(STYLUS_MAX_FRAGMENTS)"; \
		[ $$fragments -le $(STYLUS_MAX_FRAGMENTS) ] || exit 1; done
	@test -z "$$(git ls-files --others --exclude-standard -- stylus/.wasm-size)" || \
		{ echo "Untracked stylus/.wasm-size. Run: make snapshot-stylus, then git add it."; exit 1; }
	@git diff --quiet -- stylus/.wasm-size || \
		{ echo "WASM size changed and is not staged. Run: make snapshot-stylus, then git add it."; exit 1; }
	@table=$$(mktemp) && stylus/scripts/gas-table.sh > $$table && \
		sed -n '/^| Call | Stylus | Solidity/,/^$$/p' README.md | sed '$$d' | cmp -s - $$table; code=$$?; rm -f $$table; \
		[ $$code -eq 0 ] || { echo "README.md's gas table differs from stylus/.gas-devnode. Regenerate it with stylus/scripts/gas-table.sh."; exit 1; }

check-activation: check-stylus
	@for rpc in $(ROBINHOOD_RPC_URL) $(ROBINHOOD_TESTNET_RPC_URL) $(ARBITRUM_RPC_URL); do \
		for dir in stylus/contracts/*/; do \
			(cd $$dir && cargo stylus check -e $$rpc) || exit 1; done; done

gas-stylus-devnode:
	@test -s stylus/target/devnode-gas.txt || { echo "No dev-node gas report. Run: make test-stylus-devnode"; exit 1; }
	@awk -v tol=$(DEVNODE_GAS_TOLERANCE_BPS) '{ gas = $$NF; name = $$0; sub(/ [0-9]+$$/, "", name) } \
		NR == FNR { now[name] = gas; next } \
		{ seen[name] = 1; if (!(name in now)) { bad = 1; printf "%s: not measured, snapshot %s\n", name, gas; next } \
			d = now[name] - gas; if ((d < 0 ? -d : d) * 10000 > gas * tol) { bad = 1; \
			printf "%s: %s L2 gas, snapshot %s\n", name, now[name], gas } } \
		END { for (k in now) if (!(k in seen)) { bad = 1; printf "%s: %s L2 gas, not in the snapshot\n", k, now[k] } \
			if (bad) { print "Dev-node gas moved more than " tol / 100 "%. Run: make snapshot-stylus-devnode"; exit 1 } \
			print "Dev-node gas within " tol / 100 "% of stylus/.gas-devnode." }' \
		stylus/target/devnode-gas.txt stylus/.gas-devnode

snapshot-stylus-devnode:
	@test -s stylus/target/devnode-gas.txt || { echo "No dev-node gas report. Run: make test-stylus-devnode"; exit 1; }
	@LC_ALL=C sort stylus/target/devnode-gas.txt > stylus/.gas-devnode && cat stylus/.gas-devnode

gas-table:
	@test -s stylus/target/devnode-gas.txt || { echo "No dev-node gas report. Run: make test-stylus-devnode"; exit 1; }
	@stylus/scripts/gas-table.sh stylus/target/devnode-gas.txt

deploy-stylus: check-docker check-foundry
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(SIGNER)" && test -x stylus/scripts/$(CONTRACT)-args.sh || \
		{ echo "Usage: make deploy-stylus CHAIN=<4663|46630|42161|412346> SIGNER='<signer flags>' [CONTRACT=<band|margin>] [REGISTRY=<file>]"; exit 1; }
	@set -f; args=$$(stylus/scripts/$(CONTRACT)-args.sh $(or $(REGISTRY),deployments/$(CHAIN).json)) && \
		CARGO_STYLUS_VERSION=$(CARGO_STYLUS_VERSION) stylus/scripts/reproducible.sh deploy $(RPC_URL_$(CHAIN)) $(CONTRACT) \
		$(SIGNER) -- $$args

deploy-band-feeds: check-foundry submodules
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(SIGNER)" || \
		{ echo "Usage: make deploy-band-feeds CHAIN=<4663|46630|42161|412346> SIGNER='<forge wallet flags>' [REGISTRY=<file>]"; exit 1; }
	@contracts/script/deploy-band-feeds.sh $(RPC_URL_$(CHAIN)) $(or $(REGISTRY),deployments/$(CHAIN).json) $(SIGNER)

simulate-supply-vault: check-foundry submodules
	@mkdir -p stylus/target && jq --arg owner "$$(jq -r .tapehouse.Owner deployments/46630.json)" \
		'.tapehouse.Owner //= $$owner' deployments/4663.json > stylus/target/4663-registry.json
	cd contracts && REGISTRY=../stylus/target/4663-registry.json forge script script/DeploySupplyVault.s.sol \
		--fork-url $(ROBINHOOD_RPC_URL) --fork-block-number $(ROBINHOOD_FORK_BLOCK)

deploy-supply-vault: check-foundry submodules
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(SIGNER)" || \
		{ echo "Usage: make deploy-supply-vault CHAIN=<4663|46630|412346> SIGNER='<forge wallet flags>' [REGISTRY=<file>]"; exit 1; }
	cd contracts && REGISTRY=$(abspath $(or $(REGISTRY),deployments/$(CHAIN).json)) forge script script/DeploySupplyVault.s.sol \
		--rpc-url $(RPC_URL_$(CHAIN)) $(SIGNER) --broadcast $(if $(filter 412346,$(CHAIN)),--skip-simulation)
	cd contracts && REGISTRY=$(abspath $(or $(REGISTRY),deployments/$(CHAIN).json)) forge script script/Register.s.sol \
		--rpc-url $(RPC_URL_$(CHAIN))

deploy-margin-accounts: check-foundry submodules
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(SIGNER)" || \
		{ echo "Usage: DEBT_CAP=<USDG units> WEEKEND_DEBT_CAP=<USDG units> PREMIUM_RATE=<bps a year> RESERVE_SHARE=<bps> make deploy-margin-accounts CHAIN=<4663|46630|412346> SIGNER='<forge wallet flags of the vault owner>' [REGISTRY=<file>]"; exit 1; }
	@contracts/script/deploy-margin-accounts.sh $(RPC_URL_$(CHAIN)) $(or $(REGISTRY),deployments/$(CHAIN).json) $(SIGNER)

deploy-liquidator: check-foundry submodules
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(SIGNER)" || \
		{ echo "Usage: make deploy-liquidator CHAIN=<4663|46630|412346> SIGNER='<forge wallet flags of the accounts' owner>' [REGISTRY=<file>]"; exit 1; }
	@contracts/script/deploy-liquidator.sh $(RPC_URL_$(CHAIN)) $(or $(REGISTRY),deployments/$(CHAIN).json) $(SIGNER)

deploy-contracts-devnode:
	rm -rf contracts/broadcast/*/412346
	$(MAKE) deploy-supply-vault CHAIN=412346 SIGNER="--private-key $(DEVNODE_KEY)" REGISTRY=stylus/target/devnode-registry.json
	DEBT_CAP=1000000000000 WEEKEND_DEBT_CAP=500000000000 PREMIUM_RATE=500 RESERVE_SHARE=1000 \
		$(MAKE) deploy-margin-accounts CHAIN=412346 SIGNER="--private-key $(DEVNODE_KEY)" \
		REGISTRY=stylus/target/devnode-registry.json
	$(MAKE) deploy-liquidator CHAIN=412346 SIGNER="--private-key $(DEVNODE_KEY)" REGISTRY=stylus/target/devnode-registry.json

test-contracts-devnode:
	contracts/script/devnode-supply-vault-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) stylus/target/devnode-registry.json
	contracts/script/devnode-margin-accounts-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) stylus/target/devnode-registry.json
	contracts/script/devnode-liquidator-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) stylus/target/devnode-registry.json

verify-supply-vault: check-foundry submodules
	@case "$(CHAIN)" in 4663|46630) test -n "$(RPC_URL_$(CHAIN))" ;; *) false ;; esac || \
		{ echo "Usage: make verify-supply-vault CHAIN=<4663|46630>, with the chain's RPC URL set"; exit 1; }
	cd contracts && forge verify-contract --chain-id $(CHAIN) --verifier sourcify --rpc-url $(RPC_URL_$(CHAIN)) \
		--creation-transaction-hash $$(jq -r --arg vault "$$(jq -r .tapehouse.SupplyVault ../deployments/$(CHAIN).json)" \
			'.transactions[] | select((.contractAddress // "" | ascii_downcase) == ($$vault | ascii_downcase)) | .hash' \
			broadcast/DeploySupplyVault.s.sol/$(CHAIN)/run-latest.json) \
		--watch $$(jq -r .tapehouse.SupplyVault ../deployments/$(CHAIN).json) \
		src/SupplyVault.sol:SupplyVault

verify-margin-accounts: check-foundry submodules
	@case "$(CHAIN)" in 4663|46630) test -n "$(RPC_URL_$(CHAIN))" ;; *) false ;; esac || \
		{ echo "Usage: make verify-margin-accounts CHAIN=<4663|46630>, with the chain's RPC URL set"; exit 1; }
	@contracts/script/verify-margin-accounts.sh $(RPC_URL_$(CHAIN)) deployments/$(CHAIN).json

verify-liquidator: check-foundry submodules
	@case "$(CHAIN)" in 4663|46630) test -n "$(RPC_URL_$(CHAIN))" ;; *) false ;; esac || \
		{ echo "Usage: make verify-liquidator CHAIN=<4663|46630>, with the chain's RPC URL set"; exit 1; }
	@contracts/script/verify-liquidator.sh $(RPC_URL_$(CHAIN)) deployments/$(CHAIN).json

verify-band-feeds: check-foundry submodules
	@case "$(CHAIN)" in 4663|46630|42161) test -n "$(RPC_URL_$(CHAIN))" ;; *) false ;; esac || \
		{ echo "Usage: make verify-band-feeds CHAIN=<4663|46630|42161>, with the chain's RPC URL set"; exit 1; }
	@contracts/script/verify-band-feeds.sh $(RPC_URL_$(CHAIN)) deployments/$(CHAIN).json

verify-stylus: check-docker
	@test -n "$(RPC_URL_$(CHAIN))" && test -n "$(TX)" || \
		{ echo "Usage: make verify-stylus CHAIN=<4663|46630|42161|412346> TX=<deployment tx> [CONTRACT=<band|margin>]"; exit 1; }
	@CARGO_STYLUS_VERSION=$(CARGO_STYLUS_VERSION) stylus/scripts/reproducible.sh verify $(RPC_URL_$(CHAIN)) $(TX) $(CONTRACT)

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
	stylus/scripts/devnode-reference-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) \
		$$(cat stylus/target/devnode-margin) stylus/target/devnode-registry.json
	stylus/scripts/devnode-margin-e2e.sh $(DEVNODE_RPC_URL) $(DEVNODE_KEY) \
		$$(cat stylus/target/devnode-margin) stylus/target/devnode-registry.json
