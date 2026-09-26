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
| Binaryen | 133, from `stylus/Stylus.toml` | `make` downloads it on first use into `$XDG_CACHE_HOME/binaryen` (default `~/.cache/binaryen`), checked against a pinned SHA-256 |
| Docker | any recent | [docker.com](https://www.docker.com) — for `make devnode`, `deploy-stylus` and `verify-stylus` |
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
| `make gas` | Contract sizes, gas snapshots and WASM sizes, failing on any change; the Stylus fragment count needs network access to Robinhood Chain |
| `make snapshot` | Regenerates the gas snapshots and `stylus/.wasm-size`; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-stylus` · `test-stylus` · `lint-stylus` · `gas-stylus` · `snapshot-stylus` | Stylus programs only |
| `make check-activation` | `cargo stylus check` of every program against Robinhood Chain, its testnet and Arbitrum One |
| `make build-apps` · `test-apps` · `lint-apps` | Apps only |
| `make devnode deploy-stylus-devnode test-stylus-devnode` | Deploys `band` to a local dev node and writes live RedStone prices through it |
| `make deploy-stylus CHAIN=<id> SIGNER='<flags>'` | Deploys `band` reproducibly, configured from `deployments/<id>.json`, and prints the transaction and address |
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

- `writePrices(bytes32[] feedIds, bytes payload)` verifies the payload and stores each value. Anyone may call it; a package no newer than the stored one reverts with `PackageNotNewer`, and the three market-status feeds must be written together (`IncompleteStatus`).
- `price(bytes32 feedId)` returns the value, its package timestamp in milliseconds and the block timestamp it was written at, and never reverts.
- `asset(bytes32 symbol)` returns the asset's Chainlink feed, RedStone feed ID, index feed ID and Stock Token, all zero for an unknown symbol.
- `legs(bytes32 symbol)` returns the Chainlink answer and its `updatedAt` in seconds, then the 24/7 price and its package timestamp in milliseconds, both in the token's terms. Zero for a leg that is unset, unreadable or not known in the current terms; it never reverts and never judges staleness.
- `quote(bytes32 symbol)` returns the band: its state (0 halted, 1 degraded, 2 closed, 3 open), how many legs are live, its centre, half-width in basis points, and lower and upper bounds.
- `corporateAction(bytes32 symbol)` returns the Stock Token's multiplier change as it affects the band: its status (0 none, 1 scheduled, 2 not yet confirmed), when it takes effect, and the multipliers before and after it.
- `syncMultiplier(bytes32 symbol)` confirms a halted change (`MultiplierConfirmed`), then records the token's latest one (`MultiplierRecorded`). Anyone may call it.
- `session()` returns Chainlink's 24/5 session (0 not known, 1 closed, 2 open), NYSE's state and its next state (0 not known, 1 regular hours, 2 short close, 3 weekend or holiday), when NYSE changes state, and the session's next boundary: when an open session closes or a closed one reopens. Times are in milliseconds, zero when not known.
- `anchor(bytes32 symbol)` returns the Chainlink print, index price and print time that an asset priced from an index is anchored to. Each new anchor emits `Anchored`.
- `variance(bytes32 feedId)` returns the EWMA variance of a configured asset's 24/7 feed, in centi-basis-points squared per minute.
  - λ = 0.94 per sample.
  - A feed's first written price is its first sample. After that, a sample is a written price at least 50 s of package time after the previous sample, and its return is normalised to one minute.
  - A zero price, or one above 2^64 − 1, is never sampled.
  - A gap never resets the variance.

The band itself, in `stylus/contracts/band/src/quote.rs`, is integer arithmetic over both legs:

| Term | Rule |
|---|---|
| Live legs | The 24/7 leg is live up to 120 s old. Chainlink is live while its session is open and it is at most 86,460 s old (the heartbeat plus 60 s). |
| Centre | The live 24/7 leg, else Chainlink. Never an average with a sleeping leg. |
| Half-width | At least 30 bps, and at least z = 3 volatility over a 3-minute latency. It adds the disagreement above Chainlink's 50 bps deviation when both legs are live, 25 bps with one leg, 50 bps without the 24/7 leg, and 25 bps while the 24/7 leg is made from an index. It is capped at 1,500 bps. |
| State | Open or closed with Chainlink's session, degraded with fewer live legs than expected or while the session is not known, halted with none. |

### The session

Chainlink's equity feeds on Robinhood Chain follow a 24/5 session, and the band uses the same session on every chain: from 20:00 ET on the evening before each trading day to 20:00 ET on it, closed over weekends and NYSE holidays. The band derives it from RedStone's signed New York market status (`NY_MARKET_CURRENT_STATUS`, `NY_MARKET_NEXT_STATUS` and `NY_MARKET_NEXT_CHANGE_TIME`, one package), never from a calendar in code:

| NYSE, as signed | 24/5 session |
|---|---|
| Regular hours, or a short close between trading days | Open |
| A short close followed by a weekend or holiday close (the evening before a holiday) | Open only for the post-market after the regular close |
| A weekend or holiday close | Open for 4 hours after the regular close (post-market), then closed until 20:00 ET before the next trading day: 13.5 hours before a regular open, or 4 hours before the midnight that ends a holiday |

Each boundary is a fixed distance from a signed time, and none of those distances spans a daylight-saving change, so the rule needs no time zone. The three values must come from one package and be at most an hour old; otherwise the session is not known, Chainlink is treated as asleep, and every band that is not halted is degraded. The regular close is recorded when a status written during regular hours announces it, so a band deployed after a Friday close treats that Friday's post-market as closed.

### Multipliers

A Robinhood Stock Token's price is its share's price times the token's ERC-8056 multiplier, rounded down as the token rounds: `floor(share × uiMultiplier / 1e18)`. Chainlink's feeds on Robinhood Chain already price the token; RedStone's 24/7 feeds price the share. The band therefore prices the 24/7 leg as `floor(share × uiMultiplier / 1e18)`, reading the multiplier from the token at every call, so a change takes effect exactly at its `effectiveAt`.

A multiplier change is a step. The token does not expose the multiplier it replaced, so the band records each change it sees: in the constructor, and whenever `syncMultiplier` is called. Only a change recorded before its step has a known size; the constructor records the token's current multiplier too.

| Change | Before `effectiveAt` | After `effectiveAt` |
|---|---|---|
| Under 50 bps, known size (a reinvested dividend) | Nothing changes | A Chainlink round from before the step is scaled to the new terms, because the share did not move |
| 50 bps or more, or unknown size (a split, a large dividend) | `corporateAction` reports it as scheduled | The band is halted until `syncMultiplier` sees a Chainlink round that started at or after the step fall inside the band the 24/7 leg draws alone and, when the size is known, closer to the 24/7 price in the new terms than in the old, because the share may or may not have moved with the step |

A halted asset stays halted until it is confirmed: a later change does not replace an unconfirmed one. A change that took effect before deployment has no known size, so a new band starts halted for that asset until its first `syncMultiplier` after a keeper write. Each call reads the token's scheduled multiplier and `effectiveAt`, so a change nobody recorded, or one replaced without notice, is treated as unknown in size rather than priced wrongly; a current multiplier other than the one the recorded change implies is a missed step, halted from the moment it is seen. A token that cannot be read halts its asset. A 24/7 variance sample never spans a recorded step. For an asset priced from an index, the anchor's print is a Chainlink round like any other: scaled across a step under 50 bps, and no leg after a material step until a round after it re-anchors the index. Its confirmation compares Chainlink with a leg anchored to Chainlink, so it checks only the index's move since that print.

### SPY

SPY has no RedStone 24/7 feed. Its 24/7 leg is its last Chainlink print moved by the S&P 500 index (`USA500.Y---24_7`) since that print: `print × index / index at the print`. The index price at the print is anchored when an index package signed within 120 s of a new Chainlink round is written, so the leg never assumes a fixed ratio between the fund and the index. The index's variance stands in for SPY's, and the leg carries 25 bps of extra half-width for the gap between the fund and its index. SPY has no 24/7 leg until its first anchor, and none where it has no Chainlink feed. On Arbitrum One it keeps its Chainlink leg alone: those feeds follow regular hours and repeat the close in heartbeats, so an anchor there would pair a stale close with a live index. A Chainlink round that repeats the anchored price, such as a heartbeat over a weekend, keeps the anchor.

Two things follow from the anchor. While the session is open SPY's two legs are not independent: the index leg restarts from each new print, so it cross-checks the index's move since that print, not the print itself. And whoever writes the index chooses which package within 120 s of the print becomes the anchor, so the leg can carry up to 240 s of the index's move as a bias.

The asset configuration is set once, in the constructor, and cannot change. Each asset has a Chainlink feed, a 24/7 leg, or both. The 24/7 leg is a RedStone feed ID or an index feed ID, never both; an index needs the Chainlink feed that anchors it and backs one asset. An asset may name its Stock Token. Every configured feed must report 8 decimals.

`band` is larger than one contract's code limit, so it is deployed in code fragments that a small root contract points to, as Arbitrum allows (up to four fragments). The root is deployed through [StylusDeployer](https://github.com/OffchainLabs/nitro-contracts/blob/main/src/stylus/StylusDeployer.sol) at `0xcEcba2F1DC234f70Dd89F2041029807F8D03A990`, which deploys, activates and runs the constructor in one transaction, so nobody else can call the constructor first:

```bash
make deploy-stylus CHAIN=<chainId> SIGNER='--account <name> --password-file <file>'
stylus/scripts/check-band.sh <rpc> <band address> deployments/<chainId>.json
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

`make deploy-stylus` runs `cargo stylus deploy` in the reproducible build image, on the staged content of `stylus/`, with the constructor arguments `band-args.sh` builds for the launch assets from the chain's registry file, and prints the deployment transaction and the program's address. `SIGNER` is `--account <name> --password-file <file>` for a Foundry keystore, or `--private-key-path <file>`; the files are mounted read-only and never appear on a command line. `check-band.sh` checks a deployed program against the registry: every configured asset, each feed's description, and that the launch assets it leaves out are unconfigured. For development, `cargo stylus deploy --no-verify … --constructor-args …` deploys the same way outside Docker, but a program built that way can never be verified; `--constructor-args` must then be the last flag. `make devnode` deploys StylusDeployer on the dev node at its canonical address, from `stylus/scripts/stylus-deployer.hex`: the salt and initcode of its deployment on Arbitrum One.

Every program is optimised after the build by Binaryen's `wasm-opt`, with the version and flags pinned in the `[wasm-opt]` table of `stylus/Stylus.toml`. cargo-stylus folds that recipe into the program's project hash and applies it in reproducible builds, so verification rebuilds the same bytes. `make` puts that `wasm-opt` first on `PATH`; to run `cargo stylus` directly, add `${XDG_CACHE_HOME:-~/.cache}/binaryen/version_133/bin` to `PATH` yourself.

`stylus/scripts/redstone-payload.py` builds a payload from the latest packages, for example `python3 stylus/scripts/redstone-payload.py NVDA---24_7`. Several packages that share a timestamp go in one payload: `python3 stylus/scripts/redstone-payload.py NVDA---24_7 USA500.Y---24_7 NY_MARKET_STATUS`. It reads the public gateways, or the main gateway when `REDSTONE_API_KEY` is set.

## Verification

Every deployed Stylus program can be rebuilt from this repository and compared byte for byte with the code on chain, fragment by fragment. `make deploy-stylus` and `make verify-stylus` run `cargo-stylus` in the Docker image that `cargo stylus deploy` and `cargo stylus verify` use for reproducible builds (`offchainlabs/cargo-stylus-base` plus the toolchain in `stylus/rust-toolchain.toml`, linux/amd64), on the staged content of `stylus/`, so untracked files and unstaged edits never enter the build.

To verify a deployment yourself you need only Docker, make and git:

```bash
git checkout <commit>
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

It passes only when `cargo-stylus` prints `Verification successful`. Running `cargo stylus verify` directly also works, from `stylus/contracts/band`, but it exits 0 even when verification fails: read its last line. The *Stylus verification* workflow runs the same check on GitHub for any commit and deployment.

Verification covers the program's code. The configuration it was constructed with is checked by `stylus/scripts/check-band.sh`, which only reads the chain.

Each deployment below verifies from the commit that added its row: `git log --reverse --format=%H -S <address> -- README.md | head -n 1`.

| Chain | Program | Address | Deployment |
|---|---|---|---|
| 46630 | band | `0x8571fc20dD9323AF25E0D5c3F4795D8954f95498` | `0x559061ce397294f3e829ba15551fdada499b4f666d503f7f245b527a3102de08` |
| 46630 | band | `0xB8Ed13E7A695e53376C6f31E93AF44EA44386e8E` | `0x17973e0b59a1ed3069d59294842ad4e08b8551c5b627c978d7589c968b397ad8` |
| 46630 | band | `0x7c53aAa661185943F4dE195210406949fA5618C6` | `0xecaec7d4ae1997b08e034bc8c20866615e083be7856f71f96793e333e8d7b8da` |
| 46630 | band | `0x3Bf0F882d0edB6C06cBF0D78684F7464B72Fb985` | `0x03f46afb9767c866a36f3b6b53f972b7113ce54a6c8233856430576550677b96` |
| 46630 | band | `0xa6FebD4225232E71A6A46209ADB46fD3dE1f5BDA` | `0xad35314edf2339d7ea91a75fb9618ca07906e1e7ac67064cf6827f50def1d934` |
