# Tapehouse

[![CI](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml/badge.svg)](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml)

Portfolio margin for Stock Tokens on Robinhood Chain.

## Repository layout

- `apps/landing` — the website at tapehouse.xyz (Next.js)
- `contracts` — Solidity contracts (Foundry)

## Requirements

| Tool | Version | Install |
|---|---|---|
| Node.js | 24.x | [nodejs.org](https://nodejs.org) or `nvm install` (reads `.nvmrc`) |
| pnpm | 12.5.1, from `packageManager` | `corepack enable pnpm` |
| Foundry | 1.8.3 | `foundryup --install v1.8.3` |
| Slither | 0.11.6 | `pipx install --force slither-analyzer==0.11.6` |
| GNU Make | 3.81 or newer | preinstalled on macOS and most Linux distributions |

## Build and test

```bash
git clone https://github.com/Tapehouse-Labs/TAPEHOUSE.git tapehouse
cd tapehouse
make build
make test
```

`make build` initialises git submodules and installs JavaScript dependencies on first run. It stops with a message if Node or Foundry does not match the versions above.

| Target | What it does |
|---|---|
| `make build` | Builds every component |
| `make test` | Runs every test suite |
| `make lint` | Formatting, lint and static analysis |
| `make coverage` | Solidity coverage, failing below 95% of lines or branches |
| `make gas` | Contract sizes and gas snapshots, failing on any change |
| `make snapshot` | Regenerates the gas snapshots; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-apps` · `test-apps` · `lint-apps` | Apps only |
