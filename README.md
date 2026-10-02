# Tapehouse

[![CI](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml/badge.svg)](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml)

Portfolio margin for Stock Tokens on Robinhood Chain.

## Repository layout

- `apps/landing` — the website at tapehouse.xyz (Next.js)
- `contracts` — Solidity contracts (Foundry)
- `stylus` — Stylus programs (Rust): `band`, the price band, and `margin`, the portfolio margin engine, with the code they share in `stylus/crates`
- `deployments` — contract addresses per chain, one JSON file per chain ID
- `packages/sdk` — the TypeScript SDK (viem)
- `packages/mcp` — the MCP server for AI agents, on the TypeScript SDK
- `services` — the Go module: the Go SDK in `services/sdk` (go-ethereum), and the indexer and public API in `services/cmd/indexer`
- `crates` — Rust crates: the Rust SDK in `crates/sdk` (alloy)

## Requirements

| Tool | Version | Install |
|---|---|---|
| Node.js | 24.x | [nodejs.org](https://nodejs.org) or `nvm install` (reads `.nvmrc`) |
| pnpm | 12.5.1, from `packageManager` | `corepack enable pnpm` |
| Foundry | 1.8.3 | `foundryup --install v1.8.3` |
| Slither | 0.11.6 | `pipx install --force slither-analyzer==0.11.6` |
| REUSE | 6.2.0 | `pipx install --force 'reuse[charset-normalizer]==6.2.0'` — for `make lint` |
| Rust | 1.91.0 for `stylus`, 1.99.0 for `crates`, from each `rust-toolchain.toml` | [rustup](https://rustup.rs) installs them on first use |
| Go | 1.27.1, from `services/go.mod` | [go.dev](https://go.dev/dl): Go 1.21 or newer downloads it on first use |
| golangci-lint | 2.14.0 | `curl -sSfL https://golangci-lint.run/install.sh \| sh -s -- -b $(go env GOPATH)/bin v2.14.0` — for `make lint` |
| cargo-stylus | 0.10.9 | `cargo install --locked --force cargo-stylus@0.10.9` |
| Binaryen | 133, from `stylus/Stylus.toml` | `make` downloads it on first use into `$XDG_CACHE_HOME/binaryen` (default `~/.cache/binaryen`), checked against a pinned SHA-256 |
| Docker | any recent | [docker.com](https://www.docker.com) — for `make devnode`, `deploy-stylus` and `verify-stylus` |
| Python 3 and jq | any recent | preinstalled on macOS — only for the dev-node suite |
| uv | 0.12.21 | `curl -LsSf https://astral.sh/uv/0.12.21/install.sh \| sh` — only for `make test-mcp-devnode`, which runs the Python MCP client |
| ShellCheck | any recent | `brew install shellcheck`; preinstalled on GitHub's Ubuntu runners — for `make lint` |
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
| `make lint` | Formatting, lint, static analysis, ShellCheck of the scripts and `reuse lint` |
| `make coverage` | Solidity coverage, failing below 95% of lines or branches |
| `make gas` | Contract sizes, gas snapshots and WASM sizes, failing on any change; the Stylus fragment count needs network access to Robinhood Chain |
| `make snapshot` | Regenerates the gas snapshots and `stylus/.wasm-size`; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-stylus` · `test-stylus` · `lint-stylus` · `gas-stylus` · `snapshot-stylus` | Stylus programs only |
| `make check-activation` | `cargo stylus check` of every program against Robinhood Chain, its testnet and Arbitrum One |
| `make build-apps` · `test-apps` · `lint-apps` | Apps, the TypeScript SDK and the MCP server only |
| `make build-services` · `test-services` · `lint-services` | The Go module only; `test-services` reports coverage |
| `make build-crates` · `test-crates` · `lint-crates` | The Rust crates only |
| `make bindings` | Builds the contracts and regenerates the SDKs' bindings from their ABIs; commit the result |
| `make check-bindings` | Regenerates the bindings and fails if they differ from what is committed |
| `make devnode deploy-stylus-devnode test-stylus-devnode` | Deploys `band` and `margin` to a local dev node, writes live RedStone prices through the band, reads it through a `BandFeed`, margins a portfolio against stub pools and the band's session, checks the margin engine's Solidity reference against it, and updates the margin engine's parameters |
| `make gas-stylus-devnode` · `snapshot-stylus-devnode` | Compares the dev-node suite's L2 gas with `stylus/.gas-devnode`, failing on a move over 0.5%, or regenerates it |
| `make gas-table` | Prints the L2 gas of the margin program and of its Solidity reference on the same calls, from the last dev-node run |
| `make deploy-stylus CHAIN=<id> SIGNER='<flags>' [CONTRACT=margin]` | Deploys `band`, or the program `CONTRACT` names, reproducibly, configured from `deployments/<id>.json`, and prints the transaction and address |
| `make verify-stylus CHAIN=<id> TX=<hash> [CONTRACT=margin]` | Verifies a deployment of `band`, or of the program `CONTRACT` names, against the checked-out source |
| `make deploy-band-feeds CHAIN=<id> SIGNER='<flags>'` | Deploys a `BandFeed` for every launch asset the chain's band configures and prints each address |
| `make verify-band-feeds CHAIN=<id>` | Verifies every feed in the registry's `.bandFeeds` on Sourcify |
| `make simulate-supply-vault` | Simulates the supply vault's deployment on a pinned fork of Robinhood Chain |
| `make deploy-supply-vault CHAIN=<id> SIGNER='<flags>'` | Deploys the supply vault with `forge script`, configured from `deployments/<id>.json`, and records it there |
| `make verify-supply-vault CHAIN=<id>` | Verifies the registry's supply vault on Sourcify |
| `DEBT_CAP=<units> WEEKEND_DEBT_CAP=<units> PREMIUM_RATE=<bps a year> RESERVE_SHARE=<bps> make deploy-margin-accounts CHAIN=<id> SIGNER='<flags>'` | Deploys the margin accounts with `forge create`, each asset capped at its selling depth and with the weekend premium's rate and the reserve's share, records them in `deployments/<id>.json`, makes them the supply vault's borrower and the owner their guardian; the signer owns the vault |
| `make verify-margin-accounts CHAIN=<id>` | Verifies the registry's margin accounts on Sourcify |
| `make deploy-liquidator CHAIN=<id> SIGNER='<flags>'` | Deploys the liquidator with `forge create`, makes it the margin accounts' liquidator and records it in `deployments/<id>.json`; the signer owns the accounts |
| `make verify-liquidator CHAIN=<id>` | Verifies the registry's liquidator on Sourcify |
| `EXPOSURE_LIMITS=<cross>,<asset>,... SEED=<units> make deploy-gap-backstop CHAIN=<id> SIGNER='<flags>'` | Deploys the gap backstop with `forge create` and its exposure limits, makes it the margin accounts' backstop, records it in `deployments/<id>.json` and deposits the team's seed; the signer owns the accounts |
| `make verify-gap-backstop CHAIN=<id>` | Verifies the registry's gap backstop on Sourcify |
| `make deploy-reopening-auction CHAIN=<id> SIGNER='<flags>'` | Deploys the reopening auction with `forge create` and the chain's band feeds, makes it the liquidator's auction and records it in `deployments/<id>.json`; the signer owns the accounts |
| `make verify-reopening-auction CHAIN=<id>` | Verifies the registry's reopening auction on Sourcify |
| `RATE_MODEL=<optimal>,<base>,<slope1>,<slope2> FEE_SHARE=<bps> make deploy-stock-lending CHAIN=<id> SIGNER='<flags>'` | Deploys a stock lending vault for each of the margin accounts' Stock Tokens with `forge create`, makes the accounts its depositor and it their lending vault for the asset, and records it in `deployments/<id>.json`; the signer owns the accounts |
| `make verify-stock-lending CHAIN=<id>` | Verifies the registry's stock lending vaults on Sourcify |
| `make deploy-short-positions CHAIN=<id> SIGNER='<flags>'` | Deploys the short positions with `forge create`, each Stock Token shortable through its one `<ASSET>_USDG_<fee>` pool in `.uniswapV3`, makes them the borrower of every lending vault of the accounts that has none, and records them in `deployments/<id>.json`; the signer owns the accounts |
| `make verify-short-positions CHAIN=<id>` | Verifies the registry's short positions on Sourcify |
| `make deploy-morpho-oracles CHAIN=<id> SIGNER='<flags>'` | Deploys a `MorphoBandOracle` with `forge create` for every launch asset the chain's band names a Stock Token for, priced in the registry's USDG and owned by `.tapehouse.Owner`, and records each in `deployments/<id>.json` |
| `make verify-morpho-oracles CHAIN=<id>` | Verifies the registry's Morpho oracles on Sourcify |
| `BASKET=<key> NAME=<name> SYMBOL=<symbol> ASSETS=<asset>,... UNITS=<raw units>,... make deploy-basket CHAIN=<id> SIGNER='<flags>'` | Deploys a basket with `forge create`, `UNITS` raw units of each of `ASSETS` per share as its first target, owned by `.tapehouse.Owner`, records it under `.tapehouse.Baskets.<key>` and makes the margin accounts take its shares; the signer owns the accounts |
| `make verify-baskets CHAIN=<id>` | Verifies the registry's baskets on Sourcify, before any target change |
| `make deploy-gap-cover CHAIN=<id> SIGNER='<flags>'` | Deploys the gap cover with `forge create` over the registry's margin program and USDG, and records it in `deployments/<id>.json` |
| `make verify-gap-cover CHAIN=<id>` | Verifies the registry's gap cover on Sourcify |
| `make deploy-contracts-devnode test-contracts-devnode` | Deploys the supply vault, the margin accounts, the liquidator, the gap backstop, the reopening auction, a stock lending vault for each Stock Token, the short positions, a Morpho oracle for each Stock Token, the basket `PAIR` of NVDA and SPY, the gap cover and the band's feeds to the dev node and checks loans, a liquidation, a cover, the auction's wiring, a loan and recall of SPY, a short of SPY, each oracle's price the basket margined as its Stock Tokens and SPY gap cover priced from the margin program's weekend gap through them against the band and margin programs |
| `make test-sdks-devnode` | Runs each SDK's example against that deployment: SPY's band and feed, stale RedStone packages refused, a short of SPY sold and bought back by an address the account authorized, a share of the basket `PAIR` minted, deposited, read as its Stock Tokens and unwrapped, and SPY gap cover quoted and bought, or refused outside the sales |
| `make test-mcp-devnode` | Drives the MCP server against that deployment with the official TypeScript client over Streamable HTTP and the official Python client over stdio: bands, Morpho oracles and the basket `PAIR` read, a short of SPY sold and bought back, a share of `PAIR` minted, deposited and unwrapped, and gap cover written and SPY cover bought, or refused with `SalesClosed` outside the sales, through calls the server prepares and the example signs |
| `make test-indexer-devnode` | Indexes that deployment from its first block, checks every log `eth_getLogs` returns against the index, decoded as the bindings decode it, reads views, halts, Morpho oracles, debts and shorts against the SDK at the same blocks, writes the RedStone relay's packages through the band and receives them on the stream, rebuilds the index from the chain alone, and loads the API |
| `make run-indexer CHAIN=<id> [REGISTRY=<file>]` | Runs the indexer and its API on the chain's registry, with the chain's RPC URL from the table below |

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

Trading halts are signed by a separate key, whose address each chain's registry records as `.tapehouse.HaltSigner`:

```bash
cast wallet new ~/.foundry/keystores tapehouse-halt-signer
```

## The band program

`stylus/contracts/band` verifies signed RedStone data packages on-chain with the same rules as `PrimaryProdDataServiceConsumerBase` in `@redstone-finance/evm-connector` 1.0.0: five authorised signers, three unique signers per feed, the median of their values, and a package at most 180 s old and 60 s ahead. Verification errors keep the reference contract's names and selectors.

- `writePrices(bytes32[] feedIds, bytes payload)` verifies the payload and stores each value. Anyone may call it; a package no newer than the stored one reverts with `PackageNotNewer`, and the three market-status feeds must be written together (`IncompleteStatus`).
- `price(bytes32 feedId)` returns the value, its package timestamp in milliseconds and the block timestamp it was written at, and never reverts.
- `asset(bytes32 symbol)` returns the asset's Chainlink feed, RedStone feed ID, index feed ID and Stock Token, all zero for an unknown symbol.
- `legs(bytes32 symbol)` returns the Chainlink answer and its `updatedAt` in seconds, then the 24/7 price and its package timestamp in milliseconds, both in the token's terms. Zero for a leg that is unset, unreadable or not known in the current terms; it never reverts and never judges staleness.
- `quote(bytes32 symbol)` returns the band: its state (0 halted, 1 degraded, 2 closed, 3 open), how many legs are live, its centre, half-width in basis points, and lower and upper bounds.
- `corporateAction(bytes32 symbol)` returns the Stock Token's multiplier change as it affects the band: its status (0 none, 1 scheduled, 2 not yet confirmed), when it takes effect, and the multipliers before and after it.
- `syncMultiplier(bytes32 symbol)` confirms a halted change (`MultiplierConfirmed`), then records the token's latest one (`MultiplierRecorded`). Anyone may call it.
- `writeHalt(bytes32 symbol, bool halted, uint64 issuedAt, uint64 expiresAt, bytes signature)` writes a trading halt signed by the halt signer (`HaltWritten`). Anyone may call it.
- `halt(bytes32 symbol)` returns whether a signed halt holds, until when, when its last message was issued, and whether the issuer has paused the Stock Token's oracle.
- `haltSigner()` returns the address whose halts the band accepts.
- The events are declared for the SDKs in `contracts/src/interfaces/IBandPrices.sol`, which no contract imports, and a unit test checks them against the program's.
- `owner()`, `pendingOwner()`, `transferOwnership(address)`, `acceptOwnership()` and `renounceOwnership()` follow OpenZeppelin's `Ownable2Step`, with its events and errors. The owner can only rotate the halt signer, so in effect it can halt an asset, never price one.
- `setHaltSigner(address)` replaces the halt signer (`HaltSignerUpdated`). Owner only, never zero.
- `chainConfig()` returns the L2 sequencer-uptime feed the band follows and whether its Chainlink feeds follow NYSE regular hours.
- `sequencerSettled()` returns whether that sequencer has been up for more than an hour; always true without a sequencer-uptime feed. Every band is at most degraded while it is false, and a consumer that must not act during the grace reads it here, because degraded has other causes too.
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

### Trading halts

A trading halt on the underlying halts the asset's band: `quote` returns all zeros while either source holds it.

- **The signed halt.** Robinhood reports whether each Stock Token's underlying is halted (`isTradingHalt` in `https://api.robinhood.com/rhj/prices`), but not on-chain. The band accepts that state per asset when it is signed by Tapehouse's halt signer and written with `writeHalt`, which anyone may call. It is the one input the band takes on Tapehouse's own signature; every other price is signed by RedStone or Chainlink.
- **The issuer's oracle pause.** The Stock Token's `oraclePaused()` is set by the issuer's oracle pauser, as it was around CRWD's 4:1 split on 2 July 2026. The band reads it at every quote, and `syncMultiplier` confirms no multiplier step while it is set. A configured token must implement it, and a token that cannot be read counts as paused.

A signed halt is the EIP-712 message `HaltState(bytes32 symbol,bool halted,uint64 issuedAt,uint64 expiresAt)` under the domain `EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)` with name "Tapehouse Band", version "1", the chain and the band's address, so it is valid for one band only. The signature is 65 bytes, `r ‖ s ‖ v`, with the same checks as RedStone's: `v` of 27 or 28, a low `s`. A bad or foreign signature reverts with RedStone's `InvalidSignature` or `SignerNotAuthorised`. The band accepts a message when:

- it recovers to `haltSigner()`;
- it was issued no later than the block, expires after it, and spans at most an hour;
- it was issued after the last message written for the asset.

A halted message holds the asset until `expiresAt`; a message with `halted` false lifts the halt. A halt that is not refreshed lapses within an hour, and until a message lifts it the band is at most degraded: `halt` then reports a non-zero `until` in the past. The halt signer is set in the constructor, cannot be zero, and is rotated by the owner with `setHaltSigner`. A message is bound to one band, so a new band holds no halt until the signer signs for it.

### Arbitrum One

The same program is configured for Arbitrum One from `deployments/42161.json`. There it prices the **share**, not the Stock Token: Arbitrum One's Chainlink feeds report share prices, and the registry names no Stock Tokens, so no multiplier applies.

- **Chainlink follows NYSE regular hours.** Arbitrum One's feeds print from about 12:00 UTC, before the open, through regular hours and at times after the close, and never overnight; their after-hours prints can swing (on 24 August 2026 NVDA alternated eight times between two prices about 148 bps apart after the close), and their weekend heartbeat can carry a price that differs from Friday's last round. The band there counts the Chainlink leg only while the signed status says NYSE is in regular hours, and anchors SPY's index leg only then. Outside regular hours an asset has its 24/7 leg alone, so its band is degraded where Robinhood Chain's has two legs.
- **The L2 sequencer.** The band follows Arbitrum One's sequencer-uptime feed, as Chainlink advises for L2s: while the sequencer is down, and for an hour after it comes back, every band is at most degraded.

Robinhood Chain has no sequencer-uptime feed, and its Chainlink feeds follow the 24/5 session. `band-args.sh` sets both from the chain's registry file: `.chainlinkSequencer.Uptime`, and regular hours on chain 42161.

### Ownership and redeploys

The owner is the constructor's `initialOwner`, the registry's `.tapehouse.Owner`, not `msg.sender`, which is StylusDeployer. Ownership moves in two steps, through `stylus/crates/ownable`, which every program shares, and `check-band.sh` checks the owner, because verification covers the code, not the constructor's arguments. The per-asset configuration can never change; a new configuration is a new band.

A redeployed band is a new address, and every consumer holding the old one is re-pointed. Before that, a new band needs:
- an hour of keeper writes, because its variance starts at zero;
- every 24/7 feed written, then `syncMultiplier` on every asset with a Stock Token, because the constructor confirms no multiplier step;
- the halt signer's messages for any halt in progress, because a message names one band.

### SPY

SPY has no RedStone 24/7 feed. Its 24/7 leg is its last Chainlink print moved by the S&P 500 index (`USA500.Y---24_7`) since that print: `print × index / index at the print`. The index price at the print is anchored when an index package signed within 120 s of a new Chainlink round is written, so the leg never assumes a fixed ratio between the fund and the index. The index's variance stands in for SPY's, and the leg carries 25 bps of extra half-width for the gap between the fund and its index. SPY has no 24/7 leg until its first anchor, and none where it has no Chainlink feed. A Chainlink round that repeats the anchored price, such as a heartbeat over a weekend, keeps the anchor.

Two things follow from the anchor. While the session is open SPY's two legs are not independent: the index leg restarts from each new print, so it cross-checks the index's move since that print, not the print itself. And whoever writes the index chooses which package within 120 s of the print becomes the anchor, so the leg can carry up to 240 s of the index's move as a bias.

The asset configuration is set once, in the constructor, and cannot change. Each asset has a Chainlink feed, a 24/7 leg, or both. The 24/7 leg is a RedStone feed ID or an index feed ID, never both; an index needs the Chainlink feed that anchors it and backs one asset. An asset may name its Stock Token. Every configured feed must report 8 decimals. The halt signer and the owner are the registry's `.tapehouse.HaltSigner` and `.tapehouse.Owner`.

`band` is larger than one contract's code limit, so it is deployed in code fragments that a small root contract points to, as Arbitrum allows (up to four fragments). The root is deployed through [StylusDeployer](https://github.com/OffchainLabs/nitro-contracts/blob/main/src/stylus/StylusDeployer.sol) at `0xcEcba2F1DC234f70Dd89F2041029807F8D03A990`, which deploys, activates and runs the constructor in one transaction, so nobody else can call the constructor first:

```bash
make deploy-stylus CHAIN=<chainId> SIGNER='--account <name> --password-file <file>'
stylus/scripts/check-band.sh <rpc> <band address> deployments/<chainId>.json
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

`make deploy-stylus` runs `cargo stylus deploy` in the reproducible build image, on the staged content of `stylus/`, with the constructor arguments `band-args.sh` builds for the launch assets from the chain's registry file, and prints the deployment transaction and the program's address. `SIGNER` is `--account <name> --password-file <file>` for a Foundry keystore, or `--private-key-path <file>`; the files are mounted read-only and never appear on a command line. cargo-stylus sends each transaction with a fee cap equal to the gas price it just read, so a base fee that rises before the transaction lands rejects it; `MAX_FEE_PER_GAS_GWEI=<gwei>` sets the cap instead, and on Arbitrum chains the price paid is still the base fee. `check-band.sh` checks a deployed program against the registry: every configured asset, each feed's description, that the launch assets it leaves out are unconfigured, and the halt signer, owner and chain configuration. For development, `cargo stylus deploy --no-verify … --constructor-args …` deploys the same way outside Docker, but a program built that way can never be verified; `--constructor-args` must then be the last flag. `make devnode` deploys StylusDeployer on the dev node at its canonical address, from `stylus/scripts/stylus-deployer.hex`: the salt and initcode of its deployment on Arbitrum One.

Every program is optimised after the build by Binaryen's `wasm-opt`, with the version and flags pinned in the `[wasm-opt]` table of `stylus/Stylus.toml`. cargo-stylus folds that recipe into the program's project hash and applies it in reproducible builds, so verification rebuilds the same bytes. `make` puts that `wasm-opt` first on `PATH`; to run `cargo stylus` directly, add `${XDG_CACHE_HOME:-~/.cache}/binaryen/version_133/bin` to `PATH` yourself.

`stylus/scripts/redstone-payload.py` builds a payload from the latest packages, for example `python3 stylus/scripts/redstone-payload.py NVDA---24_7`. Several packages that share a timestamp go in one payload: `python3 stylus/scripts/redstone-payload.py NVDA---24_7 USA500.Y---24_7 NY_MARKET_STATUS`. It reads RedStone's authenticated gateways first, with the API keys in `REDSTONE_API_KEY` and `REDSTONE_BACKUP_API_KEY` where set, and falls back to the keyless public gateways.

## Band feeds

`contracts/src/BandFeed.sol` serves one asset's band to Solidity consumers. It reads `band` on every call, holds no owner and has nothing to configure after deployment.

- **`AggregatorV3Interface`.** `latestRoundData()` answers one side of the band (low, mid or high, chosen at deployment; lending collateral reads the low side) with 8 decimals. The band is recomputed on every read and its state already judges staleness, so the round ID, `startedAt` and `updatedAt` are the current block's timestamp, and `getRoundData` answers only that round (`NoRound`).
- **No answer is not zero.** A halted band reverts with `NoAnswer`: an unconfirmed multiplier step, a signed halt, or a paused Stock Token oracle. Where the band follows a sequencer-uptime feed, the feed also reverts with `SequencerNotSettled` while the sequencer is down or came back an hour ago or less. A consumer that treats a revert as "no price" stops new borrowing against the asset and never liquidates on a price of zero.
- **`latestBand()`** returns the whole band: state, live legs, centre, half-width and bounds; the variance behind the half-width (the asset's 24/7 feed, or its index for SPY); the session; the signed halt and the Stock Token's oracle pause; whether the sequencer is settled; and the token's own market.
- **The token's market, beside the band.** `twap` is the Stock Token's 30-minute Uniswap v3 TWAP in its pool's quote token, USDG on Robinhood Chain, with 8 decimals, and `premiumBps` is its distance from the band's centre, with USDG taken at par. It is never part of the band. Without a pool, or while the pool's history is shorter than 30 minutes, `twapValid` is false and both are zero.
- **`seal()`** records the band before the session reopens. Anyone may call it while the session is closed and its reopen is at most ten minutes away; the latest call before the reopen stands in `seals(reopenMs)` and every call emits `Sealed` with the whole band, so the record rebuilds from events alone. On Arbitrum One the reopen is still the 24/5 session's, and the band sealed then has its 24/7 leg alone, because the Chainlink leg counts only from NYSE's regular open.

`description()` says what is priced: the Robinhood Stock Token where the band names one, otherwise the share (Arbitrum One, the testnet).

```bash
make deploy-band-feeds CHAIN=<chainId> SIGNER='--account <name> --password-file <file>'
make verify-band-feeds CHAIN=<chainId>
```

`make deploy-band-feeds` deploys a low-side feed for every launch asset the registry's band configures, with the asset's `<ASSET>_USDG_*` pool from `.uniswapV3` when it has a Stock Token (it stops if the registry names two), and names each launch asset it skips. `BAND_ASSETS='<ASSET> …'` limits it to the assets named, for example to finish a deployment that stopped partway. Record each address in the registry's `.bandFeeds` group; fork tests check every entry against its band. `make verify-band-feeds` rebuilds each feed's constructor arguments from the chain and verifies the feed on [Sourcify](https://sourcify.dev), which supports all three chains and the compiler this repository pins.

## Morpho oracles

`contracts/src/MorphoBandOracle.sol` lets a Morpho Blue market take a Robinhood Stock Token as collateral priced by its band. It implements Morpho's `IOracle`: `price()` is one unit of the token in units of the market's loan token, scaled by 1e36.

- **The low edge, open or closed.** `price()` is the band's low edge in every state that has a price: open, closed and degraded, a lapsed signed halt included. While the market is closed the low edge is the conservative bound the band draws around the 24/7 price; while it is open it still keeps the band's half-width, which grows when the two legs disagree, below the centre. The loan token is taken at par with the US dollar.
- **The scale.** The band prices a whole token with 8 decimals, and Morpho counts raw units, so the oracle multiplies by `scaleFactor`, 10^(36 + loan decimals − collateral decimals − 8): 1e16 for an 18-decimal Stock Token against USDG's 6. The band's price on Robinhood Chain is already the token's, its ERC-8056 multiplier applied, so the oracle applies no multiplier of its own.
- **No answer is not zero.** A band with no price reverts with `NoAnswer`: the band is halted, or has no live leg. No live leg is the most common cause. While the market is closed the band's one leg is the 24/7 price, live for 120 seconds after its RedStone package, so a band nobody has written to for two minutes has no price until someone writes a fresh package. The band is also halted by an unconfirmed multiplier step, the issuer's oracle pause, or a halt signed by Tapehouse's halt signer, the one input the band takes on Tapehouse's own signature. Where the band follows a sequencer-uptime feed, the oracle also reverts with `SequencerNotSettled` while the sequencer is down or came back an hour ago or less; no band that names a Stock Token follows one today. Morpho Blue then refuses to borrow, to withdraw collateral against debt and to liquidate; supplying, withdrawing what is not lent, repaying and adding collateral go on. `halt()` tells the halt's sources apart.
- **One governed band.** A market's oracle is fixed when the market is created, so the oracle holds the band as state. Its owner, the registry's `.tapehouse.Owner`, re-points it with `setBand(band)` (`BandSet`). `setBand` checks one thing, that the new band names the same Stock Token for the asset (`AssetMismatch`), which keeps the collateral and its decimals; the price's 8 decimals are the band interface's, and nothing checks the new band's prices. The new band's sequencer configuration is read again. A band redeploy re-points every oracle, and no market has to move. Ownership moves in two steps and cannot be renounced (`OwnershipCannotBeRenounced`).

```bash
make deploy-morpho-oracles CHAIN=<chainId> SIGNER='--account <name> --password-file <file>'
make verify-morpho-oracles CHAIN=<chainId>
```

`make deploy-morpho-oracles` deploys an oracle for every launch asset whose Stock Token the registry's band names, priced in `.tokens.USDG` and owned by `.tapehouse.Owner`, and records it in `.morphoOracles`; it skips an asset the band prices as a share and stops if the band names a token other than the registry's. `MORPHO_ASSETS='<ASSET> …'` limits it to the assets named. Verify each oracle before it is re-pointed or its ownership moves: `make verify-morpho-oracles` reads its constructor arguments back from it.

### For curators

- **Creating a market.** Call Morpho Blue's `createMarket((loanToken, collateralToken, oracle, irm, lltv))` with `.tokens.USDG`, the oracle's `collateralToken()`, `.morphoOracles.<ASSET>`, `.morpho.AdaptiveCurveIrm` and an enabled LLTV, on `.morpho.Blue`. Check first that the oracle's `band()` is the registry's `.tapehouse.Band`, its `loanToken()` the market's, its `owner()` a timelock, that `price()` answers, and that the oracle is verified on Sourcify.
- **Who can change the price.** The oracle's `owner()` can set any price. `setBand` takes any contract that names the same Stock Token, whatever quote it returns, so the owner could raise the price and borrow a market's liquidity, lower it and liquidate every borrower, or point the oracle at a band that reverts and stop every liquidation, in one transaction. The first owner is the deployer's key: do not create a market on an oracle until its `owner()` is a timelock. Every re-point then waits out the timelock's `getMinDelay()`, and `BandSet` announces it.
- **A halt.** While the band has no price, nothing in the market can be liquidated, and a fall in the price meanwhile lands on the market's suppliers. At a 62.5% LLTV a loan at the limit becomes bad debt once the price has fallen about 29.6%, where its collateral no longer covers the debt times Morpho's incentive of 1.1268, and Morpho Blue writes bad debt off against the suppliers' deposits. A signed halt lasts as long as Tapehouse's halt signer keeps renewing it, an hour at a time, and the band's owner can replace the signer; the issuer's oracle pause halts the band too. A market on this oracle trusts the halt signer and the band's owner with that.
- **The weekend.** The band widens when the market closes, and Morpho Blue has no close factor: once a loan is past the LLTV a liquidator may repay all of it and seize collateral worth 1.1268 times the debt at a 62.5% LLTV, a loss to the borrower of about 12.7% of the debt. Over the weekend of 19 and 20 September 2026 NVDA's half-width went from 30 to 55 basis points at the close, and a loan taken at the limit in Friday's post-market stood at 62.60% at Saturday's low edge, past the LLTV while the centre rose; no loan came near bad debt. The widening comes every Friday, so liquidators can wait for it. Borrowers need headroom below the LLTV of at least the widening from the open half-width to the closed one, 25 basis points on that weekend and up to 1,500, the band's widest half-width, at worst, and room for the move on top of it; size the LLTV, and tell borrowers, accordingly.
- **Not a `BandFeed`.** A market could also price a Stock Token through Morpho's `MorphoChainlinkOracleV2` over a low-side `BandFeed`, which has 8 decimals. A `BandFeed` has no owner, so a band redeploy is a new set of feeds, and such a market could never follow it. A market that must follow a band redeploy prices through this oracle.
- **Tapehouse's own lending.** The stock lending vaults never supply Morpho markets: a lender's tokens supplied to a market could not be recalled on notice.

### For liquidators

- **Write the price first.** While the market is closed the band has a price only while someone writes RedStone packages to it, and `writePrices` is open to anyone. Liquidate from a contract that calls the band's `writePrices(feedIds, payload)` with a fresh signed package and then Morpho Blue's `liquidate`, in the same transaction. A package no newer than the one stored reverts with `PackageNotNewer(feedId, storedTimestampMs, packageTimestampMs)`: someone wrote first, so the contract catches it and liquidates against the stored price. The asset's 24/7 feed ID is the second value of the band's `asset(symbol)`, or the third, its index, where the second is zero; `stylus/scripts/redstone-payload.py` builds a payload from RedStone's gateways.
- **Which package.** The band takes any signed package up to 180 seconds old that is newer than the one stored, so whoever writes before acting, a liquidator or a borrower, chooses among the last few packages. The low edge's half-width, 30 basis points at least, covers what the price moves in that time.

## The margin engine

`margin` is the Stylus program that margins a portfolio of Stock Tokens as one position. It holds the risk parameters it generates its scenarios from: each asset's daily volatility and weekend gap, each pair's correlation, and a hard floor under every one of them. It also holds each asset's liquidation depth, whose first value is its ceiling, and the Uniswap v3 pool it is liquidated in.

- **Units.** Volatility is the standard deviation of daily log returns, in centi-basis-points (1e-6): NVDA's 31,352 is 3.1352% a day. A weekend gap is a move from the last close before the market shuts for two days or more to the next open, also in centi-basis-points. Correlation is in basis points (1e-4), and positive: every pair among the launch assets has a floor above zero, and an asset that moves against the others needs a new program. A depth is in whole USD.
- **Floors.** They are set in the constructor and never change. They are modelled on the anti-procyclicality tools EMIR gives clearing houses for margin requirements (RTS 153/2013, Article 28), applied here to the parameters beneath the margin.
  - A volatility floor is the asset's 10-year volatility, after the Article's 10-year volatility floor.
  - A correlation floor gives 25% weight to stressed observations, after the Article's second option: 75% the pair's 10-year correlation and 25% its correlation over the 60 sessions with SPY's worst return in those 10 years, 26 December 2019 to 23 March 2020 (SPY −30.2%). Correlations rise in a crash, so a floor keeps a calm year from granting a diversification benefit that disappears when it is needed.
  - A weekend-gap floor is the mean of the largest 1% of the asset's 521 such moves in the ten years, up or down, from adjusted closes and opens: 11.86% for NVDA, 5.48% for SPY. Weekend gaps are not a fixed multiple of daily volatility: their 99th percentile runs from 2.1 daily standard deviations for GOOGL to 3.2 for AAPL.
  - The first values are the floors, or, for volatility and correlation, the last 252 sessions' figures where those are higher. All come from Yahoo Finance's daily adjusted closes over the ten years to 25 September 2026, and are in `stylus/contracts/margin/parameters.json`.
- **How the parameters move.** The owner replaces them all at once with `setParameters`. Each value stays at or above its floor, or for a depth above zero and at or below its ceiling. Each moves by at most ×1.5 or ÷1.5, but a depth may fall by any amount, since a smaller depth only raises requirements; modelled on the factor of about 1.5 OCC uses when it reviews day-over-day changes in its margin coverage; a bound on each parameter is not by itself a bound on a portfolio's margin. An update comes at least a day after the one before; the first may come at any time.
- **Events.** Every value the constructor sets and every value an update changes emits `VolatilitySet`, `CorrelationSet`, `GapSet` or `DepthSet`, so the values' history rebuilds from events alone. The floors and ceilings never change and are read with the views. `contracts/src/interfaces/IMarginParameters.sol` declares the events and the views of the parameters for the SDKs, and a unit test checks its events against the program's.
- **Valid matrices only.** The correlation matrix must stay positive definite, as the correlation matrix of assets that are not linear combinations of one another is. The check is exact: Bareiss's fraction-free elimination yields every leading principal minor as an integer, and Sylvester's criterion asks that each be positive. Hadamard's inequality keeps every intermediate within 256 bits for up to 9 assets; the engine takes at most 8.
- **Ownership** follows OpenZeppelin's `Ownable2Step`, from the same `stylus/crates/ownable` as the band's, with an explicit `initialOwner`: the registry's owner, as for the band.

Reads: `assets()`, `volatility(symbol)`, `correlation(symbol, other)` and `weekendGap(symbol)`, each with its floor, `depth(symbol)` with its ceilings, `pool(symbol)`, `ethUsdFeed()`, `band()`, `weekendLeverage()`, `market()` and `lastUpdate()`.

### Scenarios

The engine draws a deterministic scenario set from its parameters on every call. Nothing is signed off chain, and anyone can rebuild the set.

| Scenarios | How each asset moves |
|---|---|
| 256 joint draws | Its daily volatility × √(horizon in days) × its row of the Cholesky factor of the correlation matrix, applied to 256 standard normal points |
| 256 with every correlation 1 | Its volatility-scaled share of one normal draw |
| 256 with every correlation 0 | Its volatility-scaled share of its own normal draw |
| market down, market up | β × the market moving by its expected shortfall at 99% (2.665 σ) over the horizon |
| weekend gap down, weekend gap up | Its weekend gap, together with every other asset |
| each asset alone, gap down and up | Its own weekend gap while every other asset stays put, one asset at a time in symbol order |

- **The normal points.** They come from a rank-1 Korobov lattice (generator 13 for 256 points), shifted in each dimension by ⌊256 · frac(d · φ)⌋, φ = (√5 − 1)/2, so the first point is not the corner of the cube. Each point stands for its 1/256 slice of the normal: it is the slice's mean, so the worst k points average to the normal's expected shortfall at level k/256, to within the rounding of each mean to a millionth. Each dimension takes every slice once.
- **The market.** It is SPY where the band prices it, and the equal-weighted portfolio of the assets where it does not, as on the testnet. β is each asset's covariance with the market over the market's variance.
- **Order.** Scenarios are drawn in the ascending order of the symbols, so an asset's scenarios do not depend on where it sits in the list.
- **Arithmetic.** It is integer: the factor in 1e12, returns in millionths, every division truncated toward zero. No return falls below −100%.
- **Horizon.** It is in seconds and scales the draws and the market move by √(horizon / 1 day); the weekend gap does not scale.

`scenario(index, horizon)` returns one scenario, and `scenarioDigest(size, horizon)` the keccak-256 of a whole set of 32, 64, 128 or 256 lattice points, in symbol order. The digest does not depend on the order of the assets. `stylus/contracts/margin/testdata/scenario-vectors.json` holds 144 digests computed independently of this program, and the dev-node suite checks the program's digest against one of them.

How closely the joint draws track the model, over 200 long-only portfolios of the launch assets at their first values: the portfolio's expected shortfall at 99% across the joint draws, over the same figure for a normal distribution with the stored covariance.

| Lattice points | 5th percentile | Median | 95th percentile |
|---|---|---|---|
| 32 | 0.698 | 0.814 | 0.948 |
| 64 | 0.720 | 0.808 | 0.951 |
| 128 | 0.865 | 0.936 | 1.011 |
| 256 | 0.934 | 0.978 | 1.033 |

### Requirement

`requirement(quantities, prices, horizon, spansClosure)` returns the margin a portfolio needs, in USD with 18 decimals, and a bit for every asset whose liquidity input is missing. It takes any portfolio: signed token amounts with 18 decimals, negative for a short, and prices in USD with 8 decimals, both in the order of `assets()`. The caller gives the horizon in seconds and says whether it spans a market closure.

| Part | What it is |
|---|---|
| Expected shortfall | The mean loss in the worst 1% of the 256 joint draws, or of the 256 draws with every correlation 0 where that is larger. 1% of 256 is 2.56 scenarios: the two worst count in full and the third at 0.56. |
| Diversification cap | Diversification takes off at most 80% of the gap between that and the sum of each asset's own expected shortfall over the draws with every correlation 1, after EMIR's cap on margin offsets (RTS 153/2013, Article 27). |
| Stress | The loss when the market moves down or up and, across a closure, when every asset gaps at once or one asset gaps alone. |
| Single-position floor | The loss when one asset alone moves by the range of FINRA Rule 4210(g)'s portfolio margin, ∓15% for a stock and −8% or +6% for the market asset, whichever asset loses most. FINRA adds up every underlying's greatest loss; this floor takes the largest one alone, so it never overrides the diversification the expected shortfall allows. |
| Liquidity | What liquidating each position in its pool costs beyond its value. |

The requirement is the largest of the first four, plus the liquidity add-on.

- **Liquidity.** Each asset is liquidated in its deepest Uniswap v3 pool on Robinhood Chain: SPY's is SPY/WETH, priced in USD with Chainlink's ETH/USD, and USDG counts at par. The second swap, from WETH to USDG, is left out.
  - The add-on is the pool's fee, the pool's discount to the valuation price when it works against the position, and the price impact.
  - The discount counts up to the asset's weekend gap: whoever moves a pool's 30-minute mean price cannot push requirements further than the asset's largest weekend move.
  - The impact grows linearly to a 10% move at the pool's depth, the USD it absorbs within that move on the liquidation's side, and whatever lies beyond the depth counts in full. Up to a 50% move the pools absorb only 1.01 to 1.61 times what they absorb within 10%.
  - The pool's price and liquidity are its own oracle's over the last 30 minutes, the time-weighted mean tick and harmonic-mean liquidity. The constructor takes a pool only if it trades the asset's Stock Token against USDG or WETH and keeps more than 1,800 observations, so 30 minutes of history cannot be overwritten.
- **Depth.** The engine takes the smaller of two depths.
  - The governed one is set with the other parameters, and its first value is its ceiling. It is measured off chain by walking the pool's initialized ticks: a walk costs one read per tick, and the 0.05% pools hold one about every 10 ticks. The first values are the lower of a Monday and a Sunday in September 2026, below.
  - The pool's own is the depth at its time-weighted liquidity, as if that liquidity held across the whole move. On these pools that gives 0.7 to 4.3 times the walked depth within a 2% move and 1.05 to 22 times at 10 to 20%, so it binds only when liquidity leaves a pool's current range.
  - The pool's own lowers the governed one by half at most. A harmonic mean counts each second with no liquidity in range as a liquidity of one, so a pool pushed out of its liquidity for a single second would read as nearly empty for 30 minutes. A real withdrawal of more than half is followed by lowering the governed depth, which may fall at once.
  - Linear impact to 10% held on every walked curve: within 1, 2 and 5% each pool absorbs at least its share of its 10% depth.
  - Every account is charged as if it had the pool to itself.
- **Missing input.** An asset held without a pool, as on Arbitrum One and the testnet, adds no liquidity add-on and sets its bit. One whose pool lacks 30 minutes of history, or needs an ETH/USD answer more than 86,460 s old, is charged its pool's fee and governed depth, with no discount, and sets its bit.

| Asset | Pool | Sells within 10% | Buys within 10% |
|---|---|---|---|
| NVDA | NVDA/USDG 0.05% | 3,561,773 | 1,377,156 |
| TSLA | TSLA/USDG 0.3% | 164,067 | 166,576 |
| AAPL | AAPL/USDG 0.05% | 169,957 | 136,784 |
| MSFT | MSFT/USDG 0.3% | 373,481 | 143,815 |
| GOOGL | GOOGL/USDG 0.05% | 386,624 | 354,428 |
| SPY | SPY/WETH 0.05% | 320,861 | 877,772 |

Over the 200 portfolios above, at 100,000 USD per unit of weight and two days, the expected shortfall with the cap runs at a median 1.020 of the model's (5th percentile 0.976, 95th 1.064): the cap covers the joint draws' shortfall. The single-position floor sets 15 of the 200 requirements. Across a closure, over three days, the gap scenarios set 20.

`stylus/contracts/margin/testdata/requirement-vectors.json` holds 120 requirements computed separately from this program, in Python, from the pools' state at Robinhood Chain block 75,093,578: with the pools read, unread and absent. The dev-node suite checks the program's requirement against one more.

### The current requirement

`currentRequirement(quantities, prices)` is the requirement an account needs now. It reads the band's session and returns the requirement, the missing bits and the regime: 0 unknown, 1 closed, 2 open, 3 closing.

| Regime | Requirement |
|---|---|
| Open | The requirement over two days, plus a 25% buffer |
| Closing: from three hours before a weekend or holiday close to the end of its four-hour post-market | Rises in a straight line from the open value to the closed one |
| Closed | The largest of the buffered open requirement, the requirement across the closure, which counts the weekend-gap scenarios, and a fifth of the gross exposure plus the liquidity add-on |
| Unknown: the band cannot tell, as when its signed status is more than an hour old, or no band is set | As closed |

- **The horizon is two days in every regime.** Over ten years, the move from a Friday close to the Monday open has 0.45 to 0.73 times the variance of one trading day, but its 99% expected shortfall matches that of 1.8 to 4.1 days of normal diffusion: the closure adds a jump, not a longer diffusion. The draws keep EMIR's two-day margin period of risk, and the weekend-gap scenarios carry the jump.
- **The buffer covers the weekend first.** It is collected all week and released into the closure: across a closure an account needs the largest of the buffered requirement, the closed one and the leverage floor below, not their sum. Over the 200 long-only portfolios above, before the liquidity add-on, the buffered requirement already covers the closed one for 194; the largest shortfall is 5.1%.
- **The rest is added while the market is liquid.** The ramp starts three hours before the regular close, so an account that must add margin or reduce can do so against the regular session's book, and ends when the session closes. The band marks the last trading day before a weekend or holiday close from its open, so the ramp needs no clock or calendar of its own.
- **At most 5× across a closure.** Leverage, gross exposure over margin, is capped at 5× across a weekend or holiday closure. Across a closure, and whenever the session is unknown, the current requirement is at least a fifth of the gross exposure, rounded up, plus the liquidity add-on; the ramp before the close brings that floor in, and an open session never sees it. `weekendLeverage()` returns the cap in basis points.
  - **Why.** At the model's own weekend limit before the liquidity add-on, a median 10.7×, 193 of the 200 long-only portfolios above would have lost more than their margin on their worst weekend or holiday close of the last ten years, every one on 16 March 2020. None does at 5×, which leaves at least a quarter more than the largest close-to-open move of any launch asset in that time, TSLA's 15.5%, and the add-on on top pays the liquidation.
  - **What it costs.** Across the weekend, what those portfolios can borrow against their stock falls from a median 90.6% of its value to 80.0%, before the liquidity add-on. An account at its weekday limit must cut its gross exposure by about half, or add margin, before the close; the ramp gives it seven hours, three of them in the regular session. Shorts count in the gross exposure, so both legs of a hedge pay the cap.
- **Isolated positions.** A position margined alone is the portfolio that holds only it: `currentRequirement` with every other quantity zero, at the cost of a full call. It gets no diversification and keeps its own floor, liquidity add-on, gap scenarios and leverage cap, and nothing else held changes it. For positions in distinct assets, isolating never needs less margin: over random portfolios in every regime, the requirement of a portfolio never exceeds the sum of its positions' own, to within a few wei of rounding per position. Split across positions, one asset's slices would together pay a smaller liquidity add-on than the whole, since the add-on grows faster than the position, so an account must hold each asset in one position only: the engine prices what it is given and does not check it.
- **No scheduled change of regime raises a requirement at once.** The ramp starts where the open requirement stands and ends where the closed one starts, and a reopen only lowers it. The tests replay a month of real sessions, every 15 minutes, and find the requirement rising only through the ramp. Three unscheduled changes can step it up: a session that turns unknown, a band that missed the regular close, and a holiday eve signed as a short close, whose ramp can only start at the close.

A caller cannot choose the unknown regime by starving the band's `session()` of gas: a failed call that leaves no more than the 64th of the gas a call holds back reverts with `InsufficientGas`.

```bash
make deploy-stylus CHAIN=<chainId> SIGNER='--account <name> --password-file <file>' CONTRACT=margin
stylus/scripts/check-margin.sh <rpc> <address> deployments/<chainId>.json
```

`stylus/scripts/margin-args.sh` builds the constructor arguments: the assets the chain's band prices, in `band-args.sh`'s order, each with its parameters from `parameters.json`, its pool from the registry's `.uniswapV3` entry `parameters.json` names and its Stock Token, zero where absent; the correlations; the file's market where the band prices it; the registry's USDG, WETH and ETH/USD feed and its band, zero where absent; and its `.tapehouse.Owner`. It refuses an asset, a pair or a pool name the file lacks. `check-margin.sh` checks the assets, the floors and ceilings, the market, the pools, the feed, the band and the owner, that the registry's band prices every asset, and the values themselves until the first update.

### Stylus versus Solidity

`contracts/test/reference/MarginReference.sol` is the margin engine in Solidity: the same scenario set, requirement, pool reads, session and weekend cap, over the same configuration, with Uniswap's own `TickMath` and `FullMath`. It is a test double for Solidity tests, not a deployment: it holds its configuration without the program's checks, parameter updates or ownership transfers.

It answers every view of the program to the bit. Its tests reproduce the 144 scenario digests and 14 launch scenarios in `stylus/contracts/margin/testdata`, the 120 requirements against Robinhood Chain's pools at block 75,093,578 in a fork test, the band's 38 sessions and a month of real ones. On the dev node it is deployed with the program's constructor arguments, and every view is compared call by call. Pools and feeds that answer as Uniswap and Chainlink do give the same result in both; malformed return data, which no real pool or feed gives, reverts the reference where the program reads no answer.

It is compiled with this repository's settings, solc 0.8.37 and the optimizer at 200 runs, and its hot loops are `unchecked`, as the program's release build is. Checked arithmetic would nearly double its gas, and `via_ir` at 10,000 runs costs 11% more on the requirement than these settings.

The dev-node suite measures both on the same calls, in L2 gas. `make devnode deploy-stylus-devnode test-stylus-devnode gas-table` prints the run's table; this one is from `stylus/.gas-devnode`, and `make gas` fails if they differ:

| Call | Stylus | Solidity | Solidity / Stylus |
|---|---|---|---|
| `scenario(0)` | 118,013 | 232,187 | 2.0× |
| `scenarioDigest(32)` | 128,910 | 561,633 | 4.4× |
| `scenarioDigest(64)` | 139,628 | 933,948 | 6.7× |
| `scenarioDigest(128)` | 163,408 | 1,678,505 | 10.3× |
| `scenarioDigest(256)` | 209,547 | 3,167,787 | 15.1× |
| `scenarioDigest(256,launch)` | 317,868 | 6,502,649 | 20.5× |
| `requirement(3)` | 287,025 | 3,758,425 | 13.1× |
| `requirement(6)` | 428,836 | 7,174,691 | 16.7× |
| `currentRequirement(3)` | 351,236 | 3,820,422 | 10.9× |
| `currentRequirement(1 of 6)` | 453,061 | 7,189,769 | 15.9× |
| `assets()` | 94,636 | 42,216 | 0.4× |
| `volatility(NVDA)` | 99,212 | 46,665 | 0.5× |
| `correlation(NVDA,SPY)` | 106,019 | 54,471 | 0.5× |

- **Stylus wins wherever the scenario set is walked,** from 4.4× for the 32-point lattice to 20.5× for the six launch assets' full set, and 13.1× and 16.7× for the requirement of three and six assets.
  - Every call pays about 95,000 to enter the uncached, three-fragment program. That is the cost of `assets()`, intrinsic gas included.
  - Net of it, the three-asset requirement is about 20× cheaper.
  - Generating the 778 rows, `scenarioDigest(256)`, takes about three quarters of the three-asset requirement's gas in Stylus and five sixths in Solidity. The P&L and aggregation over them are the rest, with the pool reads.
- **Stylus loses the small views,** at about twice the gas: `assets()`, `volatility` and `correlation` read a few words, and entering the program costs more than reading them. On chains with Arbitrum's CacheManager, caching the program lowers that entry cost. One scenario, the smallest walk, is already 2.0× cheaper in Stylus.
- **A real pool read costs the same work from either side:** both run the pool's own Solidity. In the fork, at block 75,093,578, reading and pricing the six launch pools and ETH/USD adds 458,971 gas to the reference's requirement, 76,495 per pool. The dev node's stub pools cost far less.
- **`wasm-opt` is worth 2% of size and 3% to 5% of gas.** Deployed without the `[wasm-opt]` tables on the dev node, `margin` compresses to 53,446 bytes against 52,282, and costs about 3.5% more for the three-asset requirement, 219,433 for `scenarioDigest(256)` against 209,547, and 100,055 for `assets()` against 94,636.

## The supply vault

`SupplyVault` holds the USDG that margin accounts borrow. Lenders deposit USDG and receive `thUSDG` shares, an [ERC-4626](https://eips.ethereum.org/EIPS/eip-4626) vault built on OpenZeppelin Contracts 5.7.0, and earn all the interest the borrower pays.

- **One borrower, set once.** The owner sets the borrower once, Tapehouse's margin accounts, and can never change it. The vault takes no deposit until then, so no lender's USDG can be lent to a borrower chosen after it was deposited. The vault lends only USDG that lenders deposited here.
- **Its own count.** The vault's assets are the USDG it holds by its own count plus what the borrower owes, interest included. USDG sent to it directly changes no share's value, and shares carry six more decimals than USDG, OpenZeppelin's virtual offset, so a first depositor cannot inflate the share price against the next.
- **Interest.** The debt grows by a borrow index with 27 decimals, compounded at every deposit, withdrawal, loan and repayment, so no amount of accruals rounds interest away. It is always rounded up, in lenders' favour.
- **Liquidity.** A lender withdraws up to the USDG the vault holds; the rest returns as the borrower repays. Only the borrower repays, so what each of its accounts owes adds up to the vault's debt, but for the few units Morpho's virtual shares hold.
- **Losses.** Lenders bear last whatever the borrower cannot repay and the Gap Backstop does not cover: the margin accounts write that debt off, and every share loses its part at once. Until then shares keep their price, so a lender who withdraws before a write-off loses nothing and those who stay bear all of it.
- **The rate.** The borrow rate rises with utilisation, the share of the vault's assets lent, in two straight lines: from 0% at no utilisation to 6% a year at 90%, then to 46% at full utilisation. Lenders earn the borrow rate on the share lent: 5.4% a year at 90%. The owner may change the model at once, until ownership moves to a timelock, within the bounds Aave v3.1 and later set for their own two-slope rates: the kink between 1% and 99% utilisation, the second slope at least the first, and at most 1,000% a year in all. Interest accrues at the old rate up to the change.
- **Why these numbers.** On 29 September 2026, Morpho's four largest USDG markets against Stock Tokens on Robinhood Chain, about $700,000 in all and mostly one borrower's, ran at 99.98% to 100% utilisation and paid lenders 4.8% to 7.9%. USDG's largest vault, Steakhouse's, paid 3.7% on $507 million. 5.4% at 90% sits between the two and keeps a tenth of the vault liquid.
- **USDG's issuer.** Paxos can pause USDG or freeze an address, and then wipe a frozen address's balance; one address holds all three powers, and USDG's timelock delays only upgrades (24 hours on Robinhood Chain, one on its testnet). While USDG is paused, the vault frozen, or its balance below its own count, every `max*` view returns zero, so deposits and withdrawals revert with ERC-4626's `ERC4626ExceededMax…`, and loans and repayments with USDG's own `ContractPaused()` or `AddressFrozen()`. If the issuer wipes the vault, it stays closed until `sync()`, which anyone may call, brings its count down to its balance: lenders' shares lose what was wiped, and no later deposit pays for it.
- **Permits.** `depositWithPermit` deposits with a USDG permit. USDG's EIP-712 domain is "Global Dollar", version 1, and it has no `eip712Domain()`: wallets build the domain from those two values. A permit already spent, as by a front-runner, does not stop the deposit.

Reads: `debt()`, `idle()`, `utilization()`, `borrowRate()`, `supplyRate()`, `borrowIndex()` and `rateModel()`, with the ERC-4626 views; rates are a year, with 18 decimals.

`make deploy-supply-vault` deploys with `forge script`, then a separate `Register` script records the address in `deployments/<id>.json` from Foundry's broadcast files, only for a deployment whose transaction succeeded.

## The margin accounts

`MarginAccounts` is where a borrower holds collateral and borrows USDG from the supply vault, the vault's one borrower.

- **Positions.** Each address has a cross position and, for any asset it chooses, an isolated position named by the asset's symbol. A position holds Stock Tokens, USDG and WETH and has its own debt. An asset sits in one position of an account at a time: in slices, it would pay less liquidity add-on than whole. An isolated position is margined alone, so it gets no diversification and its loss never reaches the cross position.
- **Prices.** Stock Tokens are held in raw units and valued at the low edge of their band, which already carries the token's multiplier. A degraded band, such as the hours after a reopen before Chainlink prints, backs debt at its low edge like any other. A halted band counts for nothing. WETH counts for 80% of its Chainlink price, Aave v3's loan-to-value for WETH on Arbitrum, and for nothing when that price is over a day old. USDG counts at par, the unit debt is owed in.
- **The requirement.** Borrowing, and withdrawing from a position in debt, need the position's equity, its collateral less its debt, to meet the engine's `currentRequirement` over its Stock Tokens at those prices. A position without Stock Tokens needs only non-negative equity. A position holding USDG does not borrow: it repays with that USDG first, with `repayWithCollateral`, so no position borrows against the currency it owes.
- **What closes.** Deposits and repayments stay open whatever the band, the session and the guardian; a deposit still stops at its cap or the issuer's pause or blocklist, and a repayment while USDG is paused. Borrowing stops while USDG's issuer has frozen the accounts, since no repayment could come back. Borrowing and withdrawing from a position in debt close for a position holding an asset whose band is halted, whose Stock Token has a multiplier change pending, or whose liquidity the engine cannot read; for any position while the L2 sequencer is not settled; and for a position holding Stock Tokens while the band cannot tell the session.
- **The weekend cap.** From the regular open on the last trading day before a weekend or holiday close until the reopening, a loan or withdrawal must leave the position's gross exposure to Stock Tokens at most five times its equity, `weekendLeverage`; a position holding only USDG and WETH has none. `leverage(account, position)` gives it, in basis points.
- **Debt.** The vault holds all the accounts' debt as one. Each position owns borrow shares of it, with Morpho Blue's virtual offsets, so interest accrues once, in the vault. Only the accounts repay the vault, so its debt falls only as shares are retired; anyone may repay a position through the accounts. Collateral USDG is never lent.
- **Liquidation price.** `liquidationPrice(account, position, symbol, borrowing)` searches, through the engine, for the price of one asset at which the position stops meeting its requirement once it owes `borrowing` more, every other price and pool fixed, to within a 2^14th of the current price, so a borrower sees it before signing. While the band cannot tell the session, the requirement it meets is the closed one: the price at which borrowing stops.
- **Losses.** The liquidator, set once by the owner, seizes a position's collateral and writes off what the emptied position still owes: its borrow shares are retired and the vault writes the same amount off, so lenders bear it and no other position's debt moves. Its fees go to the reserve (`collectFee`).
- **Authorization.** An account may let another address, such as a router, borrow, withdraw and repay with collateral for it, and deposit Stock Tokens and USDG into it (`setAuthorization`). Anyone may deposit WETH into an account and repay its debt.
- **The issuer's controls.** Robinhood can pause every Stock Token at once, or one, blocklist an address, the accounts' included, and burn from any address. A paused Stock Token, or one whose issuer has blocklisted the accounts, makes its own deposits, withdrawals and seizures revert with the issuer's `IsPaused()` or `Blocked(address)`; it is still valued at its band, but since it could not be seized it backs no new risk: a position holding it borrows nothing and withdraws nothing while in debt (`AssetFrozen(symbol)`), as a frozen reserve takes no new loans in Aave. Other tokens, and positions without it, carry on. An account the issuer has blocklisted cannot move its Stock Tokens through the accounts either, nor borrow or withdraw in debt from a position holding them (`Blocked(account)`), though it can repay, and a liquidator can still seize them.
- **Burns.** Each position's Stock Tokens are units of the accounts' holding of that asset. If the issuer burns from the accounts, `sync(symbol)`, which anyone may call and every call touching the asset calls first, brings every holding of it down in proportion: the accounts cannot tell whose tokens the issuer meant, so that asset's holders share the burn, rather than the last to withdraw bearing all of it. A position the burn leaves short goes through the loss cascade like any other: its collateral, then the backstop, then lenders; no other asset's holders are touched. A holding never rises: tokens sent to the accounts are credited to no one and only absorb a later burn. A holding a burn leaves worth nothing gates nothing, and anyone may clear it (`clear(account, position, symbol)`), as a write-off does; an asset burned entirely takes deposits again once every holding of it is cleared.
- **Caps.** Set in the constructor: each asset's holding, at most its selling depth in the engine over the token's price when deployed, since the engine charges every position as if it had the pool to itself; the positions' total debt; and, from the regular open on the last trading day before a closure until the reopening, or while the band cannot tell the session, the total debt a new loan may reach; debt taken before then still crosses the closure, so `debtCap` is what bounds it. A cap in tokens moves with the token's price. `holding(symbol)`, `debtCap` and `weekendDebtCap` give them.
- **Lending.** A position may lend its Stock Tokens through the asset's stock lending vault, set once by the owner (`setLending`) and owned by the accounts' owner, and earn its fee: `lend(position, token, amount, account)` moves them from its holding into the vault for the vault's shares, and `unlend` brings them back, as far as the vault can return them to the position: what has come back for its recalls, and the vault's tokens no recall waits for. A position recalls with `recall(position, token, amount, account)`: what the vault holds free comes back at once, and the rest is queued as a ticket in the vault's queue (`recalls`, `claim`), at least a hundredth of a token unless it is the position's first and whole recall (`MIN_RECALL`); its owner may keep three open, and the liquidator a fourth (`MAX_RECALLS`). What comes back for it is the position's, taken in its turn, and anyone may move it into the position's holding with `settle(account, position, token)`, which gives up the position's tickets instead once nothing it lent is worth anything. What the shares are worth (`lent`) still counts for the asset in the engine, and at `RECALL_HAIRCUT`, 5%, less in the position's equity: half the 10% move an asset's cap allows, the average premium a buy-in of a capped holding pays to bring it back. A deposit's cap counts what the accounts lent too. A paused or blocklisted Stock Token gates the position whether it holds it or lent it, and so does a vault the issuer blocklisted or burnt from before it synced (`AssetFrozen`); a write-off waits until nothing is left lent, though lent dust worth less than a USDG is swept with the rest. `sellable` gives what the liquidator may take now: the holding, and what the vault can return of what was lent; a seizure takes from the holding first and the rest straight from the vault to the buyer.
- **Permits.** `depositWithPermit` deposits with an EIP-2612 permit for the token. A Stock Token's permit domain is its current `name()`, such as "NVIDIA • Robinhood Token", and version "1", as its `eip712Domain()` reads; the issuer can rename a token, so a signer reads the domain each time rather than keep it.
- **The guardian,** set by the owner, can pause new borrowing and withdrawals from positions in debt; deposits, repayments, seizures and write-offs stay open.
- **The weekend premium.** While the 24/5 session is closed, over weekends and holidays from the post-market's end to the reopening, each position's debt carries a premium of `premiumRate` basis points a year on top of the vault's interest. The vault never lends or sees it: it is owed to the accounts beside the debt, counts against the position's equity but not against the caps, and a position owing only premium is in debt like any other. A repayment pays the debt first and the premium after, so lenders are paid before the backstop. Of each premium paid, `reserveShare` goes to the fee reserve and the rest is set aside for the backstop, which the owner sets once and which claims it with `claimPremium`. The owner sets the rate, at most 100% a year, from the change on, closed time before it accruing at the old rate, and sends the reserve on with `withdrawReserve`. A write-off forgives the premium.
- **Premium timing.** `accruePremium`, which anyone may call and every loan, repayment, write-off, withdrawal in debt and rate change calls first, records each closure from the band's `session()`: its close while ahead, from the regular open on the last trading day, and its reopening while closed. A closure recorded at both ends accrues exactly from its close to its reopening, whenever the accounts are called. Otherwise the premium counts only time seen closed: a closure first seen closed accrues from then, one whose reopening is never seen stops at its last accrual while closed, and a spell the band cannot tell the session counts only once the same closure is seen closed after it. A closed session is taken as the recorded closure only if it reopens within 72 hours of that closure's close, the longest a weekend and a holiday run; a longer closure accrues from when it is seen.
- **Baskets.** The cross position also takes the shares of the baskets the owner adds (`addBasket`, at most `MAX_BASKETS`, eight), each priced by the accounts' band, holding only their Stock Tokens in the engine's order and owned by their owner. The engine margins a position's shares as the Stock Tokens they redeem for (`inBaskets(account, position)`, rounded down), together with what it holds directly, so their correlations count and no basket is an asset of its own. Those Stock Tokens count against their caps, as the shares credited to positions redeem for, and a position holding a basket takes no new loan or withdrawal in debt while a Stock Token of its baskets is past its cap, however a gift to the basket or a rebalance put it there (`AssetCapExceeded`). They sit in the cross position: a basket goes in only while the account holds none of them in isolated positions, and none of them goes into an isolated position while a basket in the cross position holds it. A basket one of whose Stock Tokens the issuer has paused, or whose issuer has blocklisted the accounts or the basket, backs no new loan or withdrawal in debt (`BasketFrozen`). `unwrap(account, basket, shares)` redeems shares into the Stock Tokens they hold, which the position then holds in their place at the same equity and requirement: the account, an address it authorized and the liquidator may, and anyone once the position falls short at its bands' low edges, as `health` shows. A position holding shares is not written off.
- **Fixed at deployment.** The band, the engine, the vault, the caps and the reserve's share cannot be changed. A new band or engine, such as one with a new asset, means new accounts and, since the vault's borrower is set once, a new vault.

Reads: `health(account, position)` (equity and requirement in USD with 18 decimals, the engine's missing bits and regime; the equity net of the premium and of the recall haircut on what is lent), `stocks`, `debt`, `debtShares`, `totalDebtShares`, `premium`, `premiumIndex`, `collateral`, `lent`, `sellable`, `claim`, `recalls`, `lending`, `leverage`, `liquidationPrice`, `isAuthorized`, `holding`, `inBaskets`, `baskets`, `debtCap`, `weekendDebtCap`, `liquidator`, `guardian`, `borrowingPaused`, `premiumRate`, `reserveShare`, `closure`, `closureStart`, `reserve`, `backstopPremium` and `backstop`. The views count a burn once it is synced.

## The liquidator

`Liquidator` is the margin accounts' liquidator: it sells the collateral of a position that falls short, in a Dutch auction that follows the market. Anyone may start, stop and buy.

- **The market.** It is open while the band's 24/5 session is open and NYSE is in regular hours or between two trading days. From the regular close before a weekend or holiday to the regular open after it, the post-market and the 24/5 session's reopening included, and while the session is not known, it is closed, so a position the band still explains is not sold in the thin books before the regular open.
- **When a position falls short.** While the market is open, when its equity with each Stock Token at its band's low edge is below the engine's current requirement, as `health` reports it; WETH counts for 84% of its Chainlink price, Aave v3's liquidation threshold, so a position borrowed to the accounts' 80% has room. While it is closed, only when the position falls short at both edges of its bands: the engine's liquidity add-on grows with the price above the pool's, so the high edge alone is not the most favourable. A shortfall the band can still explain waits. An unknown session is judged against the open requirement with its 25% buffer, `requirement(…, 172800, false)` × 5 / 4, never the closed one a keeper's outage steps up to.
- **When nothing is judged.** While a held asset's band is halted or its Stock Token's multiplier change is not yet confirmed, while WETH is held and its Chainlink price is missing, non-positive, from the future or more than a day and a minute old, and while the L2 sequencer is not settled: such a position waits, since what it holds has no price, and a halt has no time limit. A frozen Stock Token, paused by its issuer or with the accounts blocklisted, is still valued at its band but cannot be sold; the position's other collateral can, and its USDG settles against its debt with `settleCash`.
- **The auction.** `start(account, position)` syncs the position's assets and opens an auction if it falls short. A Stock Token's price starts at its band's high edge, and WETH's at its Chainlink price, and falls a basis point a second while the market is open and a basis point every four seconds while it is closed, down to 10% below the low edge or WETH's price. An auction runs an hour while open and four hours while closed, long enough to reach its floor from a band at its 15% cap; after that, or once the market changes state, it starts again from the top. `stop` ends it once the position no longer falls short, and anyone may call it, the borrower included. `price(account, position, token)` gives the ask, in USD with 8 decimals per token.
- **A purchase.** `buy(account, position, token, amount, maxCost, receiver)` takes USDG from the buyer at the ask, at most `maxCost`, and sends the tokens to `receiver`. Half a percent of it goes to the accounts' reserve, unless the position is worth less than it owes; the rest repays the position, its debt first and its premium after. A purchase repays at most half of what the position owes, or all of it once that is 2,000 USDG or less, as Aave v3 does, and the next purchase judges the position again. While the market is closed, at most a tenth of a position's holding of a Stock Token is sold an hour, an allowance that refills continuously (`hourlyAllowance`). The contract does not rank positions: keepers start the deepest first, by `shortfall(account, position)`, which gives each position's equity and requirement.
- **Lent collateral.** What a position lent through a stock lending vault counts at the recall haircut less in its equity, as the accounts count it to within rounding, and a purchase, a reopening auction's lot and the dust swept before a write-off take it back from the vault as far as the vault can return it to the position (`sellable`). For the rest, anyone may call `recall(account, position, token)` once the position falls short: it recalls all the position lent and has not recalled yet, which comes back within the vault's notice or by a buy-in, and the position's write-off waits for it.
- **Losses.** `writeOff(account, position)` writes off what a position still owes once it holds nothing, or only holdings a burn left worthless: lenders bear it, and its premium is forgiven. If the accounts refuse it and the position is worth less than it owes, holdings worth less than a USDG (`DUST`), which anyone could leave in it to block its write-off, are swept to the caller first: Stock Tokens at their band's low edge, WETH only at a fresh Chainlink price. A holding in a halted asset is not swept. Once the accounts have a backstop, only it may write off, after covering what it can.
- **Baskets.** A basket counts as the Stock Tokens its shares redeem for. Once a position falls short, every judgement (`start`, `buy`, `settle`, `settleCash`, `recall`, `writeOff`) first redeems every basket share it holds into those Stock Tokens in the position (the accounts' `unwrap`), whatever they come to, and they are then sold as any other; a basket whose Stock Tokens the issuer has frozen stays, and the rest is sold.
- **The reopening auction.** The accounts' owner sets it once (`setAuction`). While the market is open, the positions it will sell are held out of `start` and `buy` until its clearing window closes (`heldUntil`); while it is closed, the Dutch auction still sells them. The auction, and the backstop for what the bids leave, buy through `settle(account, position, token, amount, price)` at the auction's uniform price, with the same fee and seniority as a purchase and no close factor; a position that no longer falls short sells nothing, and a rounding excess goes to the reserve. The deposits its bidders forfeit reach the reserve through `collect`.

## The gap backstop

`GapBackstop` is the junior tranche behind the supply vault's lenders, an ERC-4626 vault over USDG. Its depositors earn the weekend premium the margin accounts set aside for it and cover, first, what a position emptied by liquidation still owes. Losses fall on the borrower's collateral first, then on the backstop, then on lenders.

- **The premium.** `claim()`, which anyone may call, takes the premium the accounts set aside (`claimPremium`) and adds it to the shares at once. It arrives in lumps, as borrowers repay; a write-off forgives a position's premium, the backstop's loss.
- **Covering a shortfall.** `cover(account, position)`, which anyone may call, repays what an emptied position owes with the backstop's USDG, then has the liquidator write off the rest to lenders: once the accounts name a backstop, only it may. It reverts unless the position holds nothing, or only holdings a burn left worthless. It repays the debt alone, since the position's premium would come back to the backstop less the reserve's share. While USDG's issuer freezes the backstop it repays nothing, and lenders bear the shortfall rather than wait on the freeze.
- **Exposure limits.** The backstop covers at most `exposureLimit(position)` USDG in a closure: one limit for the cross positions together, and one for each Stock Token's isolated positions. A cross position is margined as one portfolio, so its loss has no split by asset; an isolated position holds one Stock Token, with any WETH and USDG. A closure is counted from the close it began with (`closureStart()`, which a closure longer than 72 hours keeps while the accounts record it again) until a closure begins at least 72 hours later, so a long closure is counted once; `closureMs` is the one the backstop counts now. `exposureLeft(position)` gives what it may still cover, at most the USDG it holds; `covered(position, closesMs)` what it covered. A shortfall past a limit is written off to lenders. The limits are sized against the accounts' `debtCap`, not their `weekendDebtCap`: debt taken before a closure crosses it.
- **Changing a limit.** The owner sets a new limit with `setExposureLimit`; it applies only to closures that close 13 days or more later, the longest a depositor who wanted out before it could need, so no limit changes under a closure already known. `exposureLimits(position)` gives the limit in force, the next one and when the next one starts.
- **Leaving.** A depositor starts a cooldown with `startCooldown()` for all its shares; they may be redeemed from a week later, for six days. Redemptions are shut while the market is closed, from the regular close before a weekend or holiday to the regular open after it, or while the session is unknown, and from the moment the accounts record a coming close, on the last trading day before it, until a day after its reopening, so its shortfalls are covered first, even on a trading day between two closures. A closure the accounts never recorded, because nobody touched them before or during it, sets no lock: keepers accrue the accounts on each last trading day and during each closure. Shares sent away leave the cooldown. Six days always hold an open day, even around a long weekend.
- **What the bids leave.** At the reopening auction's clearing, the backstop buys at the clearing price, at most the band's low edge, what the bids did not take of each lot, within the position's exposure limit for the closure and the USDG it holds, and only while it has shares and is not frozen. The Stock Tokens it buys are its depositors' at that moment, in proportion to their shares, and deposits are shut from the moment the accounts record a coming close until a day after its reopening, so no one buys into gains the closure already made: `gains(owner, token)` gives what each may take and `claimGains(token)` sends it. Shares sent away keep the gains already earned.
- **Seed.** The team makes the first deposit when it deploys the backstop, from the owner's address, a public `Deposit`. It leaves like any other share.
- **The issuer's controls.** The backstop counts the USDG it holds itself, so USDG sent to it directly changes no share's value. While USDG is paused, the backstop frozen or its USDG wiped, deposits and redemptions stop; `sync()` counts a wipe as a loss.

Reads: `held`, `closureMs`, `exposureLimit`, `exposureLimits`, `exposureLeft`, `covered`, `cooldowns`, `gains`, `gainsPerShare`, `gainsEpoch`, `accounts` and the ERC-4626 views.

## The reopening auction

`ReopeningAuction` sells, at NYSE's regular open, the collateral of positions that fell short over a weekend or holiday closure, in one sealed-bid, uniform-price batch per Stock Token, instead of a race of liquidation bots.

- **Phases.** They follow the band's `session()`. From the 24/5 session's reopen before a regular open that follows a weekend or holiday, until half an hour before that open, positions are enrolled and bids committed; in that last half hour, bids are revealed; in the first hour of regular trading, the round clears. A round is keyed by its Stock Token and the regular open, in milliseconds (`round(symbol, openMs)`); `phase()` gives the open and whether bids are being revealed. After a holiday the session names the open until the midnight that ends it; from then on, the open of the last round opened, if a round opened before that midnight (`lastOpenMs`).
- **Lots.** `enroll(account, position, symbol)`, which anyone may call, enters a position that falls short at its bands' low edges, as the accounts' `health` has it, into the round of one of its Stock Tokens; a cross position may be enrolled in the round of each Stock Token it holds. Its baskets are first redeemed into their Stock Tokens (the accounts' `unwrap`), so its lot counts them. Its lot is what repays all it owes at the round's floor, fee included, at most what it holds and what its lending vault can return of what it lent (`sellable`). A round keeps 32 lots; past that, a larger lot takes the smallest one's place. While the market is open, the liquidator holds an enrolled position out of its Dutch auction until the round's clearing window closes.
- **The floor.** 10% below the low edge of the band the asset's `BandFeed` sealed at the reopen (`seals`), the liquidator's floor; without a seal, or with a halted one, below the band's low edge when the round opens. No bid below it is taken.
- **Bids.** `commit(symbol, commitment, deposit)` takes a deposit of at least 100 USDG, which may exceed the bid to hide its size; the commitment is `keccak256(abi.encode(bidder, symbol, openMs, quantity, price, salt))`, so no one else can reveal it. `reveal(symbol, openMs, quantity, price, salt)` keeps from the deposit what the bid is worth at its price, its escrow, and returns the rest; a bid below the floor, worth less than 100 USDG at the floor, or worth more than its deposit is refused. A deposit not revealed loses a tenth of itself, at least 100 USDG, to the accounts' reserve (`forfeit`). A round keeps 64 bids; past that, a higher bid takes the lowest one's place, whose escrow goes back, and a bid that finds no place gets its whole deposit back (`Outbid`).
- **Clearing.** Anyone submits the clearing price with `clear(symbol, openMs, price)`, and the contract checks it in one pass over the bids: the highest bid price at which the bids at or above it take the whole supply; when all of them together take less, the lowest bid price; with no bid, the floor. Every bid above it fills, bids at it share what is left pro rata, and each lot's share of what is taken is sold through the liquidator at that one price; the gap backstop buys the rest of the lot at the same price, at most the band's low edge, so a high bid cannot set what the backstop pays; a backstop that cannot buy, its USDG wiped or its address blocked by an issuer, leaves the rest in the position and never stops the clearing. `claim(symbol, openMs, index)` then sends each bidder its tokens and the rest of its escrow. A round nobody clears in its hour refunds every bid, and its positions return to the Dutch auction.
- **What is not sold.** What neither the bids nor the backstop take stays in the positions, which return to the Dutch auction, and what an emptied position still owes goes to the gap backstop's cover. An asset halted over the week reopens into the Dutch auction: no round is keyed on a halt's end.

Reads: `round`, `lots`, `bids`, `enrolled`, `commitments`, `feeds`, `lastOpenMs`, `phase` and the constants.

## Stock lending

`StockLendingVault` lends one Stock Token: a holder's tokens earn a fee while they still back the holder's loan. There is one vault per Stock Token, an ERC-4626 vault over the token in its raw units, so dividends and splits, which move only the token's ERC-8056 multiplier, pass through every loan.

- **Who lends.** The margin accounts alone, set once as the vault's depositor (`setDepositor`), for the positions that choose to lend (The margin accounts, *Lending*). The vault's shares stay with the accounts, credited to each position.
- **Who borrows.** One borrower, set once by the owner (`setBorrower`): the short positions. It borrows open-term (`borrow`, `repay`) and may take at most `MAX_UTILIZATION`, 90%, of the vault when it borrows; returns, and a debt that grows with its fee, can bring the idle part below a tenth. It must implement `buyIn(assets)` (`IStockLendingBorrower`), or no buy-in can run. What it cannot return it writes off (`writeOff`), a loss to the lenders, as the vault's assets fall by it. The vault's owner must be the accounts' owner, so no one else chooses who borrows the tokens that back their loans.
- **The fee.** It accrues on the debt in the token itself, by a borrow index with 27 decimals, at a rate set by utilisation on a two-slope curve the owner may change within its bounds (`setRateModel`, as the supply vault's): `base` at none lent, `base + slope1` at `optimal`, and `base + slope1 + slope2` at all lent, which the vault reaches only once lenders take back what the borrower left idle. Nothing accrues while nothing is lent. The owner takes `feeShare` of it, at most `MAX_FEE_SHARE`, 20%, as shares minted at each accrual; every conversion counts those pending, so no deposit or redemption takes part of them. `borrowRate`, `supplyRate`, `utilization`, `borrowable` and `debt` give the curve and the vault now.
- **The issuer's controls.** While the token is paused, the vault blocklisted, or its balance below its count, every `max*` view returns zero: nothing enters or leaves. A burn from the vault falls on this vault's lenders, in proportion to their shares: `sync()`, which anyone may call and the accounts' `sync(symbol)` calls, brings its count down to its balance, and no other asset's holders, nor the accounts' own holding, lose anything. Until it is synced, a position that lent through it borrows nothing, as from a blocklisted vault. Tokens sent to the vault directly change no share's value.
- **Recall.** What lenders want back comes, in order, from the tokens the vault holds, from new lending and repayments, from the borrower within a notice, and last from a forced buy-in. Each recall is a ticket at the end of one queue (`recall`, `ticket`). Every token that comes into the vault is held for the tickets in their order (`locked`, up to `assigned`) before any is free to borrow or withdraw (`idle`), and tickets are taken one at a time, in turn (`head`, `claimable`, `reachable`). Each ticket has its own notice, one day (`NOTICE`); once it runs out, anyone may call `buyIn(ticket)`, which has the borrower buy in and return all the vault does not yet hold of the queue up to that ticket's end, as far as it owes it, or reverts with `BuyInFailed`; it waits while the token is paused, the vault blocklisted or its tokens burnt (`Unavailable`). A burn takes the free tokens first, then those held for the latest tickets, which wait again. The depositor takes what came back with `reclaim` and gives up the rest of the ticket next in turn with `forfeit`. Raw units are recalled, so a dividend or a split during the notice changes nothing owed.
- **Leaving.** A lender leaves by `unlend` and then the accounts' `withdraw`, to its own address, so each lender holds its own tokens should the issuer pay only eligible holders, which a vault cannot be; what the borrower holds comes back through a recall.

Reads: `depositor`, `borrower`, `idle`, `scaledDebt`, `borrowIndex`, `lastAccrual`, `rateModel`, `feeShare`, `locked`, `tickets`, `head`, `requested`, `assigned`, `served`, `ticket`, `claimable`, `reachable`, `NOTICE` and the ERC-4626 views.

## Short positions

`ShortPositions` lets an address short a Stock Token on its own margin: it borrows the token from the accounts' lending vault for that asset, sells it through the token's Uniswap v3 pool with USDG, and keeps the USDG. It is each lending vault's borrower, so the fee a lender earns is what the shorts pay.

- **One short per token.** Each address holds, for each Stock Token, one short of that token alone, with USDG as its only collateral: the margin it adds (`deposit`, `withdraw`) and the proceeds of what it sold. Nothing here deposits into or borrows from the supply vault, and no position of the margin accounts is touched. Authorization is the accounts' own (`setAuthorization`): an address an account authorized there may also sell, buy back and withdraw its shorts here. The accounts' guardian pauses new shorts and withdrawals from open ones as it pauses loans.
- **Selling and buying back.** `sell(symbol, amount, minProceeds, account)` borrows `amount` and sells it through Uniswap's `SwapRouter02`, at no less than the band's low edge. `cover(symbol, amount, maxCost, account)` buys back up to what the short owes, at no more than `maxCost` and its USDG, and repays it; a buy-back of part of what it owes also pays no more than the band's high edge, so it never leaves the short weaker, and buying back all it owes, or so much that it would retire all its borrow shares, closes the short. The fee accrues in the vault, in the token, so a short owes its share of the vault's debt as borrow shares (`position`); a dividend or a split moves what that is worth, not how many tokens are owed.
- **Margin.** The engine margins a short as a negative quantity at the band's high edge (`health`): it loses on up-moves, meets the asset's single-position floor, and is charged on its pool's buying depth. Its equity is its USDG less what it owes at the high edge; a sale or a withdrawal must leave it at or above its requirement and within the weekend leverage cap, while the band is not halted, no multiplier change is pending, the session is known and the pool's depth read.
- **A fall of 10%.** As Regulation SHO Rule 201 does, once the band's centre falls 10% from the asset's last close, a short sells at no less than the band's centre for the rest of that day and the next, in UTC (`restriction`). The close is the last centre the shorts recorded on an earlier day: every sale records one, and anyone may (`mark`).
- **Recalls.** When a recall's notice runs out and anyone calls the vault's `buyIn(ticket)`, the shorts buy in what is due through the pool, at no more than `MAX_PREMIUM_BPS`, 1%, above the band's centre, nor its high edge, and with the token's shorts' USDG alone, and repay it; every short of the token pays its part of the cost in proportion to its borrow shares, when it is next touched (`book`). A buy-in that would leave the vault owed a remnant too small for the shorts' borrow shares buys in all of it. While the issuer has paused the token or blocklisted the shorts, nothing is bought in and notices keep running.
- **Books.** Each token's shorts share a book: their borrow shares, their USDG and the cost index of the buy-ins. Once the vault is owed nothing while borrow shares remain, as after a buy-in of all it lent, the next sale opens a new book (`epoch`, `NewEpoch`): the earlier shorts owe nothing and take out what they hold.
- **Liquidation.** Anyone may `liquidate` a short whose equity is below its requirement, at the band's high edge while the market is open and at both edges while it is not: it buys back all the short owes, at no more than `MAX_PREMIUM_BPS` above the band's centre, nor its high edge, and pays the caller `LIQUIDATION_BONUS_BPS`, half a percent, of the cost from the short's USDG; it waits while the pool asks more (`AboveLimit`).
- **Who bears a loss.** A short whose USDG cannot buy back all it owes spends all of it, and the rest is written off on the vault, a loss to its lenders. A short can also owe more of the buy-ins than it holds, its part having been paid with the book's USDG: at its liquidation the book's other shorts are charged that and owe that much less of the token, written off too, so the lenders bear it as far as those shorts still owe the vault, never writing off so much that the vault is owed less than the book's borrow shares stand for; past that, those shorts bear the rest. Until it is liquidated, such a short closes only by liquidation (`Insolvent`), and the book cannot pay out what it no longer holds (`DeficitOpen`).

Reads: `accounts`, `band`, `engine`, `usdg`, `router`, `weekendLeverage`, `fee`, `position`, `health`, `epoch`, `book`, `restriction`.

## Index baskets

`Basket` is a share of a fixed set of up to eight Stock Tokens, minted and redeemed in kind: one ERC-20 that holds a portfolio, margined by the engine as that portfolio.

- **In kind.** `mint(shares, receiver, maxAssets)` takes each token's part of what the basket holds, rounded up, at most `maxAssets`; `redeem(shares, receiver, owner)` gives each token's part, rounded down. Shares follow ERC-4626's rounding, so neither ever takes from the other holders, and a share is asked for, never priced, so no donation can make a mint cost more than a raw unit of each token above its part. The first mint, and the first after every share is gone, pays the target. `previewMint` and `previewRedeem` give the amounts.
- **Raw units.** The basket holds and counts each token in its raw units, so dividends and splits, which move only the token's ERC-8056 multiplier, pass through, and mints and redemptions stay open over a multiplier change. A share is a claim on its part of what the basket holds, read from its balances: a burn by the issuer from the basket falls at once on every holder in proportion, and a token sent to it is every holder's.
- **The issuer's controls.** A share moves only as its tokens could: no transfer, mint or redemption while any of them is paused (`ComponentPaused`), and none to, from or by an address the issuer has blocklisted in any of their registries (`Blocked`). A paused token, or a blocklisted basket, stops every redemption: shares are never redeemed in part.
- **The target and its notice.** The owner proposes a new target, in raw units of each token per share, with `proposeTarget`; it takes effect `NOTICE`, seven days, later (`target`, `pendingTarget`), so a holder who disagrees redeems in kind before it binds. The set of tokens never changes: a target can only weigh them anew or leave one out.
- **Rebalancing.** Anyone may move the basket toward the target in effect with `rebalance(assetsIn, assetsOut, receiver)`: giving it tokens it holds too few of and taking those it holds too many of, each without passing its target, none both ways, and only if what comes in, at each band's low edge, is worth at least what goes out, at each band's high edge. So at every price the band allows, the basket is worth no less after a rebalance, and no rebalance moves a token away from the target in effect. Each token that moves needs a band with a price and no multiplier change pending, and the L2 sequencer settled. Whoever rebalances pays the bands' width: holders never do.
- **Margin.** The margin accounts take the shares of the baskets their owner adds and margin them as the Stock Tokens they redeem for (The margin accounts, *Baskets*); a short position's baskets are sold as those tokens (The liquidator, *Baskets*).

Reads: `components`, `totalAssets`, `target`, `pendingTarget`, `previewMint`, `previewRedeem`, `paused`, `isBlocked`, `band`, `owner` and the ERC-20 views.

## Weekend gap cover

`GapCover` sells cover against a Stock Token's fall over a market closure, paid in USDG, and is an ERC-4626 vault over USDG for the writers who back it.

- **What a cover pays.** `buy(symbol, notional, deductibleBps, limitBps, maxPremium, holder)` buys, for the coming closure, cover that pays `notional` times the fall of the asset's price from the closure's reference to its reopening price, beyond the deductible and up to the limit: a layer, as reinsurance writes one. A rise pays nothing. The price is the Stock Token's, as the band and its Chainlink feed quote it, so a dividend or a split pays nothing by itself; the cover is settled in USDG and needs no token held. Each purchase is its own cover, held by `holder` (`covers(id)`), and is not transferable.
- **When it is sold.** From the moment the band's session shows the coming close, on the last trading day before a weekend or holiday, until that close: the 24/5 session's close, or the regular close where the band's Chainlink feeds follow NYSE's regular hours (`regularHours`), as on Arbitrum One. Never while the band is halted for the asset, a multiplier change is pending or the sequencer is not settled.
- **The price.** A weekend gap scales a generalized Pareto tail of shape ½ fitted to ten years of the launch assets' weekend falls and scaled to the period the cover pays, which on Robinhood Chain ends at the 24/5 session's reopen, Sunday 20:00 ET, rather than at Monday's 09:30 open: a fall is in the tail past `TAIL_THRESHOLD_PPM` millionths of the gap, which happens on `TAIL_PROBABILITY_PPM` millionths of weekends, with a scale of `TAIL_SCALE_PPM` millionths of the gap. The deductible is at least the threshold. The premium (`quote`) is the layer's expected payout under that tail, in closed form, plus `LOADING_BPS`, 50%, rounded up, so a larger gap raises every price.
- **The gap follows the week.** During the sales the gap is `MOVE_TO_GAP_BPS` of the asset's realised move over the trading week before the close, twice it, at least `MIN_GAP_BPS` of the engine's weekend gap of the asset, `weekendGap(symbol)`, 55%, and at most `MAX_GAP_MULTIPLE` times it, eight (`pricingGap`). The move is the root of the sum of the squared daily changes, in millionths, of the asset's Chainlink answers at the eight days that end a day before the close, each the answer of the last round started before it in the phase of the feed's latest round or, where that phase has none, in the latest earlier phase that has one, so an aggregator upgrade during the week leaves it readable. The close fixes the days and the feed the rounds, so no caller chooses them. Once the week's last day has passed, the series keeps the move for the rest of its covers: its first cover keeps it, or anyone does before with `measure(symbol)`, which saves the first buyer the search. A day without such a round or a positive answer leaves the week unread, and the cover prices at the most. A calm week prices below the engine's gap and a turbulent one above it; outside the sales a quote is at the engine's gap. Fitted on 2016–21 so that premiums come to 1.5 times the payouts of calm and of turbulent weeks alike, the price keeps writers ahead on every launch asset over 2021–26, whether buyers spread evenly or crowd into the weekends after turbulent weeks.
- **A reference that has not caught up.** The feed's last round becomes the reference only if no other prints before the sales end, and it can stand up to Chainlink's deviation threshold away from the market. Where it stands above the band's centre at a purchase, the layer is priced as measured from the centre, deductible and limit both lowered by the difference, and the deductible so measured must reach the tail's threshold (`StaleReference`); `minDeductible` is the threshold plus that difference.
- **Writers.** Writers deposit USDG and receive shares. Every cover reserves its whole possible payout, `notional × (limit − deductible)`, from their USDG until it is released (`reserved`, `capacity`); a purchase beyond what is free reverts with `NoCapacity`. A cover's premium is held apart (`premiums`) and joins the writers' USDG, less what the cover pays, when it is released. Deposits are open while every outstanding cover is over the coming closure and its sales have not ended; redemptions only while no cover is outstanding.
- **Settlement.** A closure's reference is its feed's last Chainlink round started before the sales ended. Its reopening price is the first round that starts at or after the reopen, the 24/5 session's reopen or, where the feeds follow regular hours, the regular open after it, when that round lands within the settlement window inside the band recorded in its minute; otherwise it is the median of the band centres recorded in the window, and the series is flagged. Anyone records the band each minute (`observe`) and settles a series once its window ends (`settle(symbol, closesMs, referenceRound, lastRound)`), naming the feed's last rounds started before the sales ended and before the reopen, which the cover checks against the feed. Every input is a Chainlink round or an `Observed` event, so anyone can rebuild a settlement.
- **The window.** It is the 15 one-minute slots from the reopen (`WINDOW_SLOTS`), and no one chooses it: a slot is recorded only within it and only from the band itself. A minute in which the band cannot vouch for the asset (halted, closed, or the sequencer not settled), recorded in the window's first slot, moves the window by one slot; a band that cannot vouch for it for more than 15 slots after the reopen (`MAX_DELAY_SLOTS`) leaves no window, and the series is void. The median settles a series only when at least 8 of the window's slots hold the band (`QUORUM_SLOTS`), so a record that starts late, or one that keeps only the minutes it likes, cannot settle it alone: without the first round and that quorum, the series is void.
- **Void.** A void series refunds its covers' premiums. Settlement voids one whose window never opened, one without a usable round or a quorum, and one whose Stock Token's multiplier changed materially, as the band counts it, after its reference round started and before its window ended, since the reference and the reopening price would be in different terms. A series no one settles within a week of its close is void too (`void`).
- **Paying out.** Anyone may `release(id)` a cover once its series has settled or is void: it credits the holder its payout or refund (`payouts`), frees its reserve and adds its premium less its payout to the writers' USDG. The holder takes what it is owed with `claim(receiver)`.
- **Recording the closure.** The cover records the session's next close while the session is open and the reopen while it is closed (`record`, which buying calls first). A closure keeps the close first recorded for it as the key of its series (`lastCloseMs`), even when the session revises the close before it is reached; the sales end, and the reference is taken, at the close as revised. A series settles only once its reopen has been recorded while the session was closed, so keepers record it during each closure. `sales()` names the closure on sale and when its sales end, and returns zeros while no cover is sold.
- **On the dev node.** CI runs every Friday at 16:00 UTC as well as on every change, inside the sales in summer and in winter, so the dev node buys cover through the Stylus engine each week; at other times it sees the purchase refused with `SalesClosed`.
- **The issuer's controls.** The cover counts its USDG itself, so USDG sent to it directly changes no share's value. While USDG is paused, the cover frozen or its USDG wiped, deposits, redemptions and purchases stop; settlement and release go on, and claims wait. `sync()` counts a wipe as a loss to the writers, after the premiums and credits the cover holds.

Reads: `engine`, `band`, `regularHours`, `held`, `reserved`, `premiums`, `owed`, `outstanding`, `outstandingIn`, `lastCloseMs`, `reopenOf`, `sales`, `series`, `covers`, `coverCount`, `payouts`, `feed`, `capacity`, `quote`, `minDeductible`, `pricingGap` and the ERC-4626 views.

## SDKs

Three SDKs read the band, its feeds, the margin accounts, the short positions, the baskets and the gap cover, and build their transactions:

| SDK | Directory | Package | Bindings |
|---|---|---|---|
| TypeScript | `packages/sdk` | `@tapehouse/sdk` on npm, over viem 2 | `as const` ABIs in `src/generated.ts` |
| Go | `services/sdk` | `github.com/tapehouse/tapehouse/services/sdk`, in the services module | abigen v2, one package per contract in `sdk/bindings` |
| Rust | `crates/sdk` | `tapehouse-sdk` on crates.io, over alloy 2 | alloy's `sol!` over `abi/*.json` |

- **Generated from the ABIs.** `make bindings` builds the contracts and writes each SDK's bindings from the ABIs of `BandFeed`, `SupplyVault`, `MarginAccounts`, `Liquidator`, `GapBackstop`, `ReopeningAuction`, `StockLendingVault`, `ShortPositions`, `MorphoBandOracle`, `Basket`, `GapCover`, the band as a RedStone relayer calls it (`IBandPrices`: `IBand`'s reads, the band's events, `writePrices`, `price` and `haltSigner`), the margin engine's risk parameters (`IMarginParameters`), Uniswap's QuoterV2, USDG, the Stock Tokens (`IStockToken`), Morpho Blue (`IMorpho`), both as the conformance tests pin them, and Chainlink's aggregator. The bindings are committed, and `make check-bindings` fails in CI when they differ from the build. Each SDK adds by hand only what the ABIs cannot say.
- **Addresses** come from `deployments/<chainId>.json`, read at runtime: no SDK compiles in an address. A mixed-case address must carry its checksum. A `.chainlink` feed is typed by what it prices: the Stock Token on Robinhood Chain, the share on Arbitrum One.
- **Morpho oracles.** Each SDK reads an asset's `MorphoBandOracle` from `.morphoOracles` (beside `.morpho.Blue` and `.morpho.AdaptiveCurveIrm`, and Morpho Blue's markets, `.morpho.Markets`, read by their 32-byte ids as `morphoMarkets`): its `price()`, `halt()`, `band()`, `symbol()`, `collateralToken()`, `loanToken()`, `scaleFactor()` and `owner()`. A `price()` that reverts with `NoAnswer` or `SequencerNotSettled` is no price, never zero, and the SDK says why: `halted` where the band holds a signed halt, the issuer's pause or an unconfirmed multiplier step for the asset, `stale` otherwise, a band with no live leg, which a fresh `writePrices` ends, and `sequencerNotSettled`. `morpho.price` in TypeScript, `MorphoOracles().Price` in Go and `oracle_price` in Rust read all of it at one block.
- **Authorization.** `setAuthorization(authorized, allowed)` lets another address, such as a router, borrow, withdraw and deposit for an account in the margin accounts, and act for its shorts; until then both revert with `Unauthorized(caller, account)`.
- **Errors.** Every revert of the supply vault, the margin accounts, the liquidator, the gap backstop, the reopening auction, the stock lending vaults, the shorts, the gap cover, the feeds, the band, the Stock Tokens and USDG decodes by name and arguments, and so do Solidity's `Error(string)`, which carries the router's `Too little received` and `Too much requested`, and `Panic(uint256)`.
- **Sales and buy-backs.** A sale's `minProceeds` and a buy-back's `maxCost` come from a quote of the asset's pool through Uniswap's QuoterV2, `.uniswapV3.QuoterV2`, less or plus the slippage the caller names, from 0 to 10,000 basis points; never from the contract.
- **Repayments and positions.** A repayment pays a position's debt first and its premium after; `repayment` reads both at one block, and a repayment of their sum clears the position at that block. Each SDK also reads a position's `collateral`, `leverage` and `liquidationPrice`.
- **RedStone.** The band's `writePrices` takes its signed packages from a source the integrator supplies (`PackageSource`): its own gateway client, cache or relay. No SDK holds an API key or calls a gateway.
- **Baskets.** A registry's baskets are read from `.tapehouse.Baskets.<key>`, and each SDK mints and redeems them by key, with the `maxAssets` of a mint from `previewMint`, reads their `components`, `previewRedeem`, `target` and `pendingTarget`, and unwraps a position's shares. `inBaskets` reads what a position holds of each Stock Token through its baskets, by asset, at the block it reads the accounts' Stock Tokens, beside `health`: the margin a basket is held at. A basket's address goes through `deposit` and `withdraw` as any collateral, in the cross position.
- **Gap cover.** Each SDK reads the closure on sale and when its sales end, quotes cover with the writers' USDG it reserves, builds a purchase with its `maxPremium`, reads a series as open, settled with its reference and reopening prices, or void, computes a settled cover's payout, and builds the measurements, observations, settlement, releases and claims keepers and holders send.
- **Examples.** Each SDK's `devnode` example reads SPY's band and its feed, has the band refuse stale signed packages, and has a fresh address the account authorizes sell 1 SPY short for it at a quote and buy it back; the account then mints that address a share of the basket `PAIR`, which it deposits in its cross position, reads as the NVDA and SPY the share holds, and unwraps; then it writes gap cover and buys SPY cover at its quote while the cover's sales are open, or has the cover refuse it with `SalesClosed`. Each reads the account owner's key from `PRIVATE_KEY` in its environment, never from its command line. `make test-sdks-devnode` runs all three against the dev node, whose registry `make deploy-contracts-devnode` gives band feeds and the dev node a stub QuoterV2.
- **Publishing.** Each package is at version 0.1.0, declares `MIT OR Apache-2.0` and carries both license texts. None is published yet: the npm scope, the crate name and the Go module path, which needs this repository at `github.com/tapehouse/tapehouse`, are published when Tapehouse launches.

## MCP server

`packages/mcp`, `@tapehouse/mcp`, is a [Model Context Protocol](https://modelcontextprotocol.io) server through which an AI agent reads the band, its feeds, Chainlink's feeds, the Morpho oracles, the baskets, the margin accounts, the short positions and the gap cover, and prepares their transactions. It is built on the official TypeScript SDK for the 2026-07-28 specification, serves clients of the 2025 revisions too, and reads and builds everything through `@tapehouse/sdk`; its package README lists the tools.

```bash
make build-apps
TAPEHOUSE_RPC_URL=https://robinhood.drpc.org TAPEHOUSE_DEPLOYMENTS=deployments/4663.json node packages/mcp/dist/main.js
```

- **No key.** The server never holds or reads a key. A transaction tool returns unsigned calls, an approval first where an allowance is short, for the account, or an address it authorized with `accounts_set_authorization`, to check and sign in its own wallet.
- **Configuration from the environment.** `TAPEHOUSE_RPC_URL`, which the server never shows a client, `TAPEHOUSE_DEPLOYMENTS`, the registry, which must be for the chain the endpoint serves, `TAPEHOUSE_MAX_SLIPPAGE_BPS`, the most slippage a sale or buy-back may state, 100 basis points unless set, and `TAPEHOUSE_MAX_CALLS_PER_MINUTE`, the most tool calls it answers in a minute, 120 unless set.
- **Transports.** stdio by default, for a host that launches it; `--http <port>` serves Streamable HTTP on `127.0.0.1` only, refusing any other `Host` or `Origin`. Any process on the machine can reach it; a hosted endpoint needs TLS, authentication and a front of its own.
- **Validated input, data out.** Every argument is checked against its schema before any call; every result is structured content with an output schema, its block named, amounts as decimal strings of base units, reverts by name and arguments. What a tool returns is data from the chain or the registry, not instructions, and the server's instructions tell the agent so.
- **Tested with two clients.** `make test-mcp-devnode` drives it with the official TypeScript client over Streamable HTTP and the official Python client over stdio, against the dev node.

## Indexer and public API

`services/cmd/indexer` indexes a chain's Tapehouse contracts and serves them over REST and a WebSocket stream, with RedStone's signed packages relayed beside them. `services/internal/api/openapi.yaml` documents the API in OpenAPI 3.1, and the server serves it at `/v1/openapi.yaml`.

```bash
make run-indexer CHAIN=42161
curl http://127.0.0.1:8080/v1/status
```

- **What it indexes.** Every contract under the registry's `.tapehouse`, its stock lending vaults and baskets, each `.bandFeeds` and `.morphoOracles` entry, every event of each; and each Stock Token's `OraclePaused`, `OracleUnpaused`, `Paused`, `Unpaused` and `UIMultiplierUpdated`, which halt the band or move its multiplier. Each log is decoded by the ABI of the SDKs' binding of its contract, so a contract the bindings do not know is refused, never skipped. Morpho Blue, `.morpho.Blue`, is read for the markets `.morpho.Markets` records, not indexed.
- **From the chain alone.** Indexing starts at the earliest block at which one of Tapehouse's contracts has code, found by bisection over `eth_getCode`. The index is a SQLite file, `TAPEHOUSE_DB`, which nothing else writes: deleting it loses nothing, and a registry that names other contracts rebuilds it. Nothing is served that the chain or a signed package cannot reproduce: events with their block, transaction and log index, views read at a named block, RedStone packages with their signatures.
- **Reorganisations.** Every event carries `finalized`, true once its block is at or below the block the RPC reports finalized (`eth_getBlockByNumber("finalized")`, which on an Arbitrum chain is final on its parent chain). Above it, the index checks its head's hash at every step, and every log against the hash of its block; when the chain replaced a block, the index rewinds to the newest block whose hash the chain still has, and the stream sends `rewind`. A dev node reports no finalized block, so nothing there is final.
- **Reads.** `/v1/events` pages events in order, by contract, event, block range and arguments; `/v1/contracts/{contract}/views/{function}` calls any view of a contract the registry names, at the index's head or a named block, and answers a revert by name. Composite reads serve the keepers and the app at the head: `/v1/accounts/debts`, each position with debt and its health; `/v1/shorts`, each open short; `/v1/halts`, each asset's halt, labelled as the one value Tapehouse signs; `/v1/morpho/oracles`, each oracle's band, owner and price or why it has none, as the Go SDK reads them; `/v1/morpho/markets`; and `/v1/gap-cover/series/{symbol}/{closesMs}`, a series with the week's move measured for it, every band its window recorded and its settlement.
- **The stream.** `/v1/stream` is a WebSocket of events as they are indexed, of one contract or one event where asked, with the head after every step. A client that falls 256 messages behind is disconnected and resumes from `/v1/events`.
- **Tiers.** The public tier allows 10 requests a second, bursts of 20 and 2 streams per client address; behind a proxy listed in `TAPEHOUSE_TRUSTED_PROXIES`, the address is the last one in `X-Forwarded-For` that is not a proxy. The keyed tier, for integrators, allows 1,000 requests a second, bursts of 2,000 and 50 streams per key, sent as `X-API-Key`. The operator issues a key and configures only its SHA-256 digest, in `TAPEHOUSE_API_KEYS`, separated by commas, so the server holds no key: `printf %s "$KEY" | shasum -a 256`.
- **The RedStone relay.** `/v1/redstone/payload?feeds=NVDA---24_7,NY_MARKET_STATUS` returns the latest signed packages of those data package IDs and their payload, read from RedStone's keyed gateways with `REDSTONE_API_KEY` and `REDSTONE_BACKUP_API_KEY` where set, then its public ones, at most once every two seconds a gateway. Every package is served only when its signature recovers to one of the five signers the band accepts, with the three signatures the band requires, all at one timestamp. `payload` is what an SDK's `PackageSource` returns for the same IDs, so no browser or integrator needs RedStone's key.
- **Configuration from the environment.** `TAPEHOUSE_DEPLOYMENTS`, the registry; `TAPEHOUSE_RPC_URL`, which must serve the registry's chain and which the indexer never logs or serves; `TAPEHOUSE_DB`, `tapehouse-<chainId>.db` unless set; `TAPEHOUSE_LISTEN`, `127.0.0.1:8080` unless set; `TAPEHOUSE_RPC_RATE`, the most requests a second it sends the RPC, 20 unless set; and the keys and proxies above. `make run-indexer` takes the RPC URL from the chain's variable in *Networks*.
- **Latency.** With 64 clients on 8 keys, each key at 95% of the keyed tier's rate, about 7,600 reads a second of every kind, the 99th percentile is within 25 ms; `make test-indexer-devnode` measures it.

## Verification

Every deployed Stylus program can be rebuilt from this repository and compared byte for byte with the code on chain, fragment by fragment. `make deploy-stylus` and `make verify-stylus` run `cargo-stylus` in the Docker image that `cargo stylus deploy` and `cargo stylus verify` use for reproducible builds (`offchainlabs/cargo-stylus-base` plus the toolchain in `stylus/rust-toolchain.toml`, linux/amd64), on the staged content of `stylus/`, so untracked files and unstaged edits never enter the build.

To verify a deployment yourself you need only Docker, make and git:

```bash
git checkout <commit>
make verify-stylus CHAIN=<chainId> TX=<deployment tx>
```

It passes only when `cargo-stylus` prints `Verification successful`. Both targets run `cargo-stylus` in the program's own directory, `stylus/contracts/<program>`: in a workspace with more than one program, `cargo stylus verify` finds the constructor there and nowhere else. Running `cargo stylus verify` directly also works, from that directory, but it exits 0 even when verification fails: read its last line. The *Stylus verification* workflow runs the same check on GitHub for any commit and deployment.

Verification covers the program's code. The configuration it was constructed with is checked by `stylus/scripts/check-band.sh` and `stylus/scripts/check-margin.sh`, which only read the chain.

Each deployment below verifies from the commit that added its row: `git log --reverse --format=%H -S <address> -- README.md | head -n 1`.

| Chain | Program | Address | Deployment |
|---|---|---|---|
| 46630 | band | `0x8571fc20dD9323AF25E0D5c3F4795D8954f95498` | `0x559061ce397294f3e829ba15551fdada499b4f666d503f7f245b527a3102de08` |
| 46630 | band | `0xB8Ed13E7A695e53376C6f31E93AF44EA44386e8E` | `0x17973e0b59a1ed3069d59294842ad4e08b8551c5b627c978d7589c968b397ad8` |
| 46630 | band | `0x7c53aAa661185943F4dE195210406949fA5618C6` | `0xecaec7d4ae1997b08e034bc8c20866615e083be7856f71f96793e333e8d7b8da` |
| 46630 | band | `0x3Bf0F882d0edB6C06cBF0D78684F7464B72Fb985` | `0x03f46afb9767c866a36f3b6b53f972b7113ce54a6c8233856430576550677b96` |
| 46630 | band | `0xa6FebD4225232E71A6A46209ADB46fD3dE1f5BDA` | `0xad35314edf2339d7ea91a75fb9618ca07906e1e7ac67064cf6827f50def1d934` |
| 46630 | band | `0xA5896f75679F3D7c3aAe31FAd94C5C1B7FfB9A3f` | `0x4a5d7f7c48174bce43c1d3788720546a86901282cf3accc7f0e3afffbd7cc3c1` |
| 46630 | band | `0xa70118d3324D90532E7D2854627b13CacE305641` | `0x7f003392328803b2fd7331f6dcc2d96f387e8e9b811cdc74b01e9c77ba5e1310` |
| 42161 | band | `0xa0c0Cb25F5504395fB977DA7186a7B68cBcfa8Eb` | `0x48f8a4eb8629970093c70e250c14037f360b2f380282d11f6e113a17bd408a12` |
| 46630 | margin | `0x0FB6856c36c25e01190d6a8f2eBbE28aCA05a341` | `0xb6192ef3b615fdc2254eb718178ea74c0d83e7670b4301bd0fabf1c980eb0a38` |
| 46630 | margin | `0x9B2DB8135222d7B05aEA29B54aE0317E8640D6B0` | `0x311698b2648cf30aa5a1ab9fcca23aeb1f6fb8db0a780c83e844e4cbe66cd5db` |

## License

Licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or <https://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](LICENSE-MIT) or <https://opensource.org/licenses/MIT>)

at your option.

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this repository by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.

The repository follows the [REUSE](https://reuse.software) specification: every file states its license, in an SPDX header or in `REUSE.toml`, which also records the origin and license of the few files that come from others. The landing's generated images are CC0-1.0: machine-generated images may carry no copyright, and whatever rights the contributors hold are waived. `make lint` runs `reuse lint`.

The libraries in `contracts/lib` keep their own licenses. `BandFeed` compiles in Uniswap v3's `OracleLibrary`, `TickMath` and `IUniswapV3Pool`, which are GPL-2.0-or-later. Taking `BandFeed.sol` under MIT, its deployed bytecode is a combined work under GPL-2.0-or-later, whose source is this repository; each deployed feed's exact source is the commit that deployed it.
