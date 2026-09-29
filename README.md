# Tapehouse

[![CI](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml/badge.svg)](https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml)

Portfolio margin for Stock Tokens on Robinhood Chain.

## Repository layout

- `apps/landing` — the website at tapehouse.xyz (Next.js)
- `contracts` — Solidity contracts (Foundry)
- `stylus` — Stylus programs (Rust): `band`, the price band, and `margin`, the portfolio margin engine, with the code they share in `stylus/crates`
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
| `make lint` | Formatting, lint, static analysis and ShellCheck of the scripts |
| `make coverage` | Solidity coverage, failing below 95% of lines or branches |
| `make gas` | Contract sizes, gas snapshots and WASM sizes, failing on any change; the Stylus fragment count needs network access to Robinhood Chain |
| `make snapshot` | Regenerates the gas snapshots and `stylus/.wasm-size`; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-stylus` · `test-stylus` · `lint-stylus` · `gas-stylus` · `snapshot-stylus` | Stylus programs only |
| `make check-activation` | `cargo stylus check` of every program against Robinhood Chain, its testnet and Arbitrum One |
| `make build-apps` · `test-apps` · `lint-apps` | Apps only |
| `make devnode deploy-stylus-devnode test-stylus-devnode` | Deploys `band` and `margin` to a local dev node, writes live RedStone prices through the band, reads it through a `BandFeed`, and updates the margin engine's parameters |
| `make gas-stylus-devnode` · `snapshot-stylus-devnode` | Compares the dev-node suite's L2 gas with `stylus/.gas-devnode`, failing on a move over 0.5%, or regenerates it |
| `make deploy-stylus CHAIN=<id> SIGNER='<flags>' [CONTRACT=margin]` | Deploys `band`, or the program `CONTRACT` names, reproducibly, configured from `deployments/<id>.json`, and prints the transaction and address |
| `make verify-stylus CHAIN=<id> TX=<hash> [CONTRACT=margin]` | Verifies a deployment of `band`, or of the program `CONTRACT` names, against the checked-out source |
| `make deploy-band-feeds CHAIN=<id> SIGNER='<flags>'` | Deploys a `BandFeed` for every launch asset the chain's band configures and prints each address |
| `make verify-band-feeds CHAIN=<id>` | Verifies every feed in the registry's `.bandFeeds` on Sourcify |

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

## The margin engine

`margin` is the Stylus program that margins a portfolio of Stock Tokens as one position. It starts with the risk parameters it generates its scenarios from: each asset's daily volatility, each pair's correlation, and a hard floor under every one of them.

- **Units.** Volatility is the standard deviation of daily log returns, in centi-basis-points (1e-6): NVDA's 31,352 is 3.1352% a day. Correlation is in basis points (1e-4), and positive: every pair among the launch assets has a floor above zero, and an asset that moves against the others needs a new program.
- **Floors.** They are set in the constructor and never change. They are modelled on the anti-procyclicality tools EMIR gives clearing houses for margin requirements (RTS 153/2013, Article 28), applied here to the parameters beneath the margin.
  - A volatility floor is the asset's 10-year volatility, after the Article's 10-year volatility floor.
  - A correlation floor gives 25% weight to stressed observations, after the Article's second option: 75% the pair's 10-year correlation and 25% its correlation over the 60 sessions with SPY's worst return in those 10 years, 26 December 2019 to 23 March 2020 (SPY −30.2%). Correlations rise in a crash, so a floor keeps a calm year from granting a diversification benefit that disappears when it is needed.
  - The first values are the floors, or the last 252 sessions' figures where those are higher. All come from Yahoo Finance's daily adjusted closes over the ten years to 25 September 2026, and are in `stylus/contracts/margin/parameters.json`.
- **How the parameters move.** The owner replaces them all at once with `setParameters`. Each value stays at or above its floor and moves by at most ×1.5 or ÷1.5, modelled on the factor of about 1.5 OCC uses when it reviews day-over-day changes in its margin coverage; a bound on each parameter is not by itself a bound on a portfolio's margin. An update comes at least a day after the one before; the first may come at any time. Later values come from an off-chain calibrator, through this same path.
- **Events.** Every value the constructor sets and every value an update changes emits `VolatilitySet` or `CorrelationSet`, so the values' history rebuilds from events alone. The floors never change and are read with the views.
- **Valid matrices only.** The correlation matrix must stay positive definite, as the correlation matrix of assets that are not linear combinations of one another is. The check is exact: Bareiss's fraction-free elimination yields every leading principal minor as an integer, and Sylvester's criterion asks that each be positive. Hadamard's inequality keeps every intermediate within 256 bits for up to 9 assets; the engine takes at most 8.
- **Ownership** follows OpenZeppelin's `Ownable2Step`, from the same `stylus/crates/ownable` as the band's, with an explicit `initialOwner`. It moves to the timelock with the band's.

Reads: `assets()`, `volatility(symbol)` and `correlation(symbol, other)`, each with its floor, and `lastUpdate()`.

```bash
make deploy-stylus CHAIN=<chainId> SIGNER='--account <name> --password-file <file>' CONTRACT=margin
stylus/scripts/check-margin.sh <rpc> <address> deployments/<chainId>.json
```

`stylus/scripts/margin-args.sh` builds the constructor arguments: the assets the chain's band prices, in `band-args.sh`'s order, with their parameters from `parameters.json`, and the registry's `.tapehouse.Owner`. It refuses an asset or a pair the file lacks. `check-margin.sh` checks the assets, the floors and the owner, that the registry's band prices every asset, and the values themselves until the first update.

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
