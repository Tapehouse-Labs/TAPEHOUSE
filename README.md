# Tapehouse

[![CI](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml/badge.svg)](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml)

Portfolio margin for Stock Tokens on Robinhood Chain.

## Repository layout

- `apps/landing` — the website at tapehouse.xyz (Next.js)
- `contracts` — Solidity contracts (Foundry)
- `stylus` — Stylus programs (Rust), starting with `band`, the price band
- `deployments` — contract addresses per chain, one JSON file per chain ID

## Requirements

| Tool | Version | Install |
|---|---|---|
| Node.js | 24.x | [nodejs.org](https://nodejs.org) or `nvm install` (reads `.nvmrc`) |
| pnpm | 12.5.1, from `packageManager` | `corepack enable pnpm` |
| Foundry | 1.8.3 | `foundryup --install v1.8.3` |
| Slither | 0.11.6 | `pipx install --force slither-analyzer==0.11.6` |
| Rust | 1.91.0, from `stylus/rust-toolchain.toml` | [rustup](https://rustup.rs) installs it on first use |
| cargo-stylus | 0.10.9 | `cargo install --locked --force cargo-stylus@0.10.9` |
| Docker | any recent | [docker.com](https://www.docker.com) — only for `make devnode` |
| Python 3 and jq | any recent | preinstalled on macOS — only for the dev-node suite |
| GNU Make | 3.81 or newer | preinstalled on macOS and most Linux distributions |

## Build and test

```bash
git clone https://github.com/Tapehouse-Labs/TAPEHOUSE.git tapehouse
cd tapehouse
make build
make test
```

`make build` initialises git submodules and installs JavaScript dependencies on first run. It stops with a message if Node, Foundry or cargo-stylus does not match the versions above.

| Target | What it does |
|---|---|
| `make build` | Builds every component |
| `make test` | Runs every test suite |
| `make lint` | Formatting, lint and static analysis |
| `make coverage` | Solidity coverage, failing below 95% of lines or branches |
| `make gas` | Contract sizes, gas snapshots and WASM sizes, failing on any change |
| `make snapshot` | Regenerates the gas snapshots and `stylus/.wasm-size`; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-stylus` · `test-stylus` · `lint-stylus` · `gas-stylus` · `snapshot-stylus` | Stylus programs only |
| `make check-activation` | `cargo stylus check` of every program against Robinhood Chain, its testnet and Arbitrum One |
| `make build-apps` · `test-apps` · `lint-apps` | Apps only |
| `make devnode deploy-stylus-devnode test-stylus-devnode` | Deploys `band` to a local dev node and writes live RedStone prices through it |
| `make initcode-stylus` | Builds every program's initcode reproducibly, in the image `cargo stylus verify` uses |
| `make verify-stylus CHAIN=<id> TX=<hash>` | Verifies a deployment against the checked-out source |

## Networks

| Network | Chain ID | Foundry alias | RPC variable | Default |
|---|---|---|---|---|
| Robinhood Chain | 4663 | `robinhood` | `ROBINHOOD_RPC_URL` | `https://robinhood.drpc.org` |
| Robinhood Chain testnet | 46630 | `robinhood-testnet` | `ROBINHOOD_TESTNET_RPC_URL` | `https://robinhood-testnet.drpc.org` |
| Arbitrum One | 42161 | `arbitrum` | `ARBITRUM_RPC_URL` | `https://arbitrum.gateway.tenderly.co` |
| Local dev node | 412346 | `devnode` | — | `http://127.0.0.1:8547` |

The defaults are public archive endpoints and need no key. Set a variable to use your own endpoint. The Makefile exports the defaults, so run fork tests through `make`. To call `forge` or `cast` directly, export the variables first. `make devnode` starts a local Nitro node on ArbOS 61, the same version as Robinhood Chain and Arbitrum One, and `make devnode-stop` removes it.

## Addresses

`deployments/<chainId>.json` lists every contract this repository uses on that chain. Fork tests check each entry against the chain at a pinned block.

## Deployer keys

No key is stored in this repository. Deployments sign with an encrypted Foundry keystore:

```bash
cast wallet import tapehouse-deployer --interactive
```

`forge` uses it with `--account tapehouse-deployer`. `cargo stylus` uses it with `--keystore-path ~/.foundry/keystores/tapehouse-deployer --keystore-password-path <file>`.

## The band program

`stylus/contracts/band` verifies signed RedStone data packages on-chain with the same rules as `PrimaryProdDataServiceConsumerBase` in `@redstone-finance/evm-connector` 1.0.0: five authorised signers, three unique signers per feed, the median of their values, and a package at most 180 s old and 60 s ahead. Verification errors keep the reference contract's names and selectors.

- `writePrices(bytes32[] feedIds, bytes payload)` verifies the payload and stores each value. Anyone may call it; a package no newer than the stored one reverts with `PackageNotNewer`.
- `price(bytes32 feedId)` returns the value, its package timestamp in milliseconds and the block timestamp it was written at, and never reverts.
- `asset(bytes32 symbol)` returns the asset's Chainlink feed and RedStone feed ID, both zero for an unknown symbol.
- `legs(bytes32 symbol)` returns the Chainlink answer and its `updatedAt` in seconds, then the stored RedStone 24/7 value and its package timestamp in milliseconds. Zero for a leg that is unset or unreadable; it never reverts and never judges staleness.

The asset configuration is set once, in the constructor, and cannot change. Each asset has a Chainlink feed, a RedStone feed ID or both. Every configured feed must report 8 decimals. The program is deployed through [StylusDeployer](https://github.com/OffchainLabs/nitro-contracts/blob/main/src/stylus/StylusDeployer.sol) at `0xcEcba2F1DC234f70Dd89F2041029807F8D03A990`, which deploys, activates and runs the constructor in one transaction, so nobody else can call the constructor first:

```bash
make initcode-stylus
args=$(stylus/scripts/band-args.sh deployments/<chainId>.json) && read -r symbols feeds feed_ids <<<"$args"
stylus/scripts/deploy-initcode.sh <rpc> stylus/target/reproducible/band.initcode.hex \
  "0x5585258d$(cast abi-encode 'f(bytes32[],address[],bytes32[])' "$symbols" "$feeds" "$feed_ids" | cut -c3-)" <signer flags>
stylus/scripts/check-band.sh <rpc> <band address> deployments/<chainId>.json
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

`0x5585258d` is the selector of every Stylus constructor. For development, `cargo stylus deploy --no-verify … --constructor-args …` deploys the same way in one command, but a program built that way can never be verified; `--constructor-args` must then be the last flag. `band-args.sh` builds the arguments for the launch assets from the chain's registry file, and `check-band.sh` checks a deployed program against it: every configured asset, each feed's description, and that the launch assets it leaves out are unconfigured. `make devnode` deploys StylusDeployer on the dev node at its canonical address, from `stylus/scripts/stylus-deployer.hex`: the salt and initcode of its deployment on Arbitrum One.

`stylus/scripts/redstone-payload.py` builds a payload from the latest packages, for example `python3 stylus/scripts/redstone-payload.py NVDA---24_7`. It reads the public gateways, or the main gateway when `REDSTONE_API_KEY` is set.

## Verification

Every deployed Stylus program can be rebuilt from this repository and compared byte for byte with the code on chain. `make initcode-stylus` and `make verify-stylus` run `cargo-stylus` in the Docker image that `cargo stylus deploy` and `cargo stylus verify` use for reproducible builds (`offchainlabs/cargo-stylus-base` plus the toolchain in `stylus/rust-toolchain.toml`, linux/amd64), on the staged content of `stylus/`, so untracked files and unstaged edits never enter the build.

To verify a deployment yourself you need only Docker, make and git (`make initcode-stylus` also needs Foundry's `cast`):

```bash
git checkout <commit>
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

It passes only when `cargo-stylus` prints `Verification successful`. Running `cargo stylus verify` directly also works, from `stylus/contracts/band`, but it exits 0 even when verification fails: read its last line. The *Stylus verification* workflow runs the same check on GitHub for any commit and deployment.

Verification covers the program's code. The configuration it was constructed with is checked by `stylus/scripts/check-band.sh`, which only reads the chain.

Each deployment below verifies from the commit that added its row: `git log -1 --format=%H -S <address> -- README.md`.

| Chain | Program | Address | Deployment |
|---|---|---|---|
| 46630 | band | `0x8571fc20dD9323AF25E0D5c3F4795D8954f95498` | `0x559061ce397294f3e829ba15551fdada499b4f666d503f7f245b527a3102de08` |
