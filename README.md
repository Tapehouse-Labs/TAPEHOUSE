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
| REUSE | 6.2.0 | `pipx install --force 'reuse[charset-normalizer]==6.2.0'` — for `make lint` |
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
| `make lint` | Formatting, lint, static analysis, ShellCheck of the scripts and `reuse lint` |
| `make coverage` | Solidity coverage, failing below 95% of lines or branches |
| `make gas` | Contract sizes, gas snapshots and WASM sizes, failing on any change; the Stylus fragment count needs network access to Robinhood Chain |
| `make snapshot` | Regenerates the gas snapshots and `stylus/.wasm-size`; commit the result |
| `make build-contracts` · `test-contracts` · `lint-contracts` · `coverage-contracts` · `gas-contracts` · `snapshot-contracts` | Contracts only |
| `make build-stylus` · `test-stylus` · `lint-stylus` · `gas-stylus` · `snapshot-stylus` | Stylus programs only |
| `make check-activation` | `cargo stylus check` of every program against Robinhood Chain, its testnet and Arbitrum One |
| `make build-apps` · `test-apps` · `lint-apps` | Apps only |
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
| `make deploy-contracts-devnode test-contracts-devnode` | Deploys the supply vault, the margin accounts, the liquidator, the gap backstop and the reopening auction to the dev node and checks loans, a liquidation, a cover and the auction's wiring through them against the band and margin programs |

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

`margin` is the Stylus program that margins a portfolio of Stock Tokens as one position. It holds the risk parameters it generates its scenarios from: each asset's daily volatility and weekend gap, each pair's correlation, and a hard floor under every one of them. It also holds each asset's liquidation depth, whose first value is its ceiling, and the Uniswap v3 pool it is liquidated in.

- **Units.** Volatility is the standard deviation of daily log returns, in centi-basis-points (1e-6): NVDA's 31,352 is 3.1352% a day. A weekend gap is a move from the last close before the market shuts for two days or more to the next open, also in centi-basis-points. Correlation is in basis points (1e-4), and positive: every pair among the launch assets has a floor above zero, and an asset that moves against the others needs a new program. A depth is in whole USD.
- **Floors.** They are set in the constructor and never change. They are modelled on the anti-procyclicality tools EMIR gives clearing houses for margin requirements (RTS 153/2013, Article 28), applied here to the parameters beneath the margin.
  - A volatility floor is the asset's 10-year volatility, after the Article's 10-year volatility floor.
  - A correlation floor gives 25% weight to stressed observations, after the Article's second option: 75% the pair's 10-year correlation and 25% its correlation over the 60 sessions with SPY's worst return in those 10 years, 26 December 2019 to 23 March 2020 (SPY −30.2%). Correlations rise in a crash, so a floor keeps a calm year from granting a diversification benefit that disappears when it is needed.
  - A weekend-gap floor is the mean of the largest 1% of the asset's 521 such moves in the ten years, up or down, from adjusted closes and opens: 11.86% for NVDA, 5.48% for SPY. Weekend gaps are not a fixed multiple of daily volatility: their 99th percentile runs from 2.1 daily standard deviations for GOOGL to 3.2 for AAPL.
  - The first values are the floors, or, for volatility and correlation, the last 252 sessions' figures where those are higher. All come from Yahoo Finance's daily adjusted closes over the ten years to 25 September 2026, and are in `stylus/contracts/margin/parameters.json`.
- **How the parameters move.** The owner replaces them all at once with `setParameters`. Each value stays at or above its floor, or for a depth above zero and at or below its ceiling. Each moves by at most ×1.5 or ÷1.5, but a depth may fall by any amount, since a smaller depth only raises requirements; modelled on the factor of about 1.5 OCC uses when it reviews day-over-day changes in its margin coverage; a bound on each parameter is not by itself a bound on a portfolio's margin. An update comes at least a day after the one before; the first may come at any time.
- **Events.** Every value the constructor sets and every value an update changes emits `VolatilitySet`, `CorrelationSet`, `GapSet` or `DepthSet`, so the values' history rebuilds from events alone. The floors and ceilings never change and are read with the views.
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
| `requirement(3)` | 284,014 | 3,755,130 | 13.2× |
| `requirement(6)` | 425,824 | 7,171,396 | 16.8× |
| `currentRequirement(3)` | 348,196 | 3,816,682 | 11.0× |
| `currentRequirement(1 of 6)` | 453,033 | 7,189,041 | 15.9× |
| `assets()` | 94,636 | 42,216 | 0.4× |
| `volatility(NVDA)` | 99,212 | 46,665 | 0.5× |
| `correlation(NVDA,SPY)` | 106,019 | 54,471 | 0.5× |

- **Stylus wins wherever the scenario set is walked,** from 4.4× for the 32-point lattice to 20.5× for the six launch assets' full set, and 13.2× and 16.8× for the requirement of three and six assets.
  - Every call pays about 95,000 to enter the uncached, three-fragment program. That is the cost of `assets()`, intrinsic gas included.
  - Net of it, the three-asset requirement is about 20× cheaper.
  - Generating the 778 rows, `scenarioDigest(256)`, takes about three quarters of the three-asset requirement's gas in Stylus and five sixths in Solidity. The P&L and aggregation over them are the rest, with the pool reads.
- **Stylus loses the small views,** at about twice the gas: `assets()`, `volatility` and `correlation` read a few words, and entering the program costs more than reading them. On chains with Arbitrum's CacheManager, caching the program lowers that entry cost. One scenario, the smallest walk, is already 2.0× cheaper in Stylus.
- **A real pool read costs the same work from either side:** both run the pool's own Solidity. In the fork, at block 75,093,578, reading and pricing the six launch pools and ETH/USD adds 458,971 gas to the reference's requirement, 76,495 per pool. The dev node's stub pools cost far less.
- **`wasm-opt` is worth 2% of size and 3% to 5% of gas.** Deployed without the `[wasm-opt]` tables on the dev node, `margin` compresses to 53,446 bytes against 52,282, and costs 294,039 for the three-asset requirement against 284,014, 219,433 for `scenarioDigest(256)` against 209,547, and 100,055 for `assets()` against 94,636.

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
- **The guardian,** set by the owner, can pause new borrowing and withdrawals from positions in debt; deposits, repayments, seizures and write-offs stay open.
- **The weekend premium.** While the 24/5 session is closed, over weekends and holidays from the post-market's end to the reopening, each position's debt carries a premium of `premiumRate` basis points a year on top of the vault's interest. The vault never lends or sees it: it is owed to the accounts beside the debt, counts against the position's equity but not against the caps, and a position owing only premium is in debt like any other. A repayment pays the debt first and the premium after, so lenders are paid before the backstop. Of each premium paid, `reserveShare` goes to the fee reserve and the rest is set aside for the backstop, which the owner sets once and which claims it with `claimPremium`. The owner sets the rate, at most 100% a year, from the change on, closed time before it accruing at the old rate, and sends the reserve on with `withdrawReserve`. A write-off forgives the premium.
- **Premium timing.** `accruePremium`, which anyone may call and every loan, repayment, write-off, withdrawal in debt and rate change calls first, records each closure from the band's `session()`: its close while ahead, from the regular open on the last trading day, and its reopening while closed. A closure recorded at both ends accrues exactly from its close to its reopening, whenever the accounts are called. Otherwise the premium counts only time seen closed: a closure first seen closed accrues from then, one whose reopening is never seen stops at its last accrual while closed, and a spell the band cannot tell the session counts only once the same closure is seen closed after it. A closed session is taken as the recorded closure only if it reopens within 72 hours of that closure's close, the longest a weekend and a holiday run; a longer closure accrues from when it is seen.
- **Fixed at deployment.** The band, the engine, the vault, the caps and the reserve's share cannot be changed. A new band or engine, such as one with a new asset, means new accounts and, since the vault's borrower is set once, a new vault.

Reads: `health(account, position)` (equity and requirement in USD with 18 decimals, the engine's missing bits and regime; the equity net of the premium), `stocks`, `debt`, `debtShares`, `totalDebtShares`, `premium`, `premiumIndex`, `collateral`, `leverage`, `liquidationPrice`, `isAuthorized`, `holding`, `debtCap`, `weekendDebtCap`, `liquidator`, `guardian`, `borrowingPaused`, `premiumRate`, `reserveShare`, `closure`, `closureStart`, `reserve`, `backstopPremium` and `backstop`. The views count a burn once it is synced.

## The liquidator

`Liquidator` is the margin accounts' liquidator: it sells the collateral of a position that falls short, in a Dutch auction that follows the market. Anyone may start, stop and buy.

- **The market.** It is open while the band's 24/5 session is open and NYSE is in regular hours or between two trading days. From the regular close before a weekend or holiday to the regular open after it, the post-market and the 24/5 session's reopening included, and while the session is not known, it is closed, so a position the band still explains is not sold in the thin books before the regular open.
- **When a position falls short.** While the market is open, when its equity with each Stock Token at its band's low edge is below the engine's current requirement, as `health` reports it; WETH counts for 84% of its Chainlink price, Aave v3's liquidation threshold, so a position borrowed to the accounts' 80% has room. While it is closed, only when the position falls short at both edges of its bands: the engine's liquidity add-on grows with the price above the pool's, so the high edge alone is not the most favourable. A shortfall the band can still explain waits. An unknown session is judged against the open requirement with its 25% buffer, `requirement(…, 172800, false)` × 5 / 4, never the closed one a keeper's outage steps up to.
- **When nothing is judged.** While a held asset's band is halted or its Stock Token's multiplier change is not yet confirmed, while WETH is held and its Chainlink price is missing, non-positive, from the future or more than a day and a minute old, and while the L2 sequencer is not settled: such a position waits, since what it holds has no price, and a halt has no time limit. A frozen Stock Token, paused by its issuer or with the accounts blocklisted, is still valued at its band but cannot be sold; the position's other collateral can, and its USDG settles against its debt with `settleCash`.
- **The auction.** `start(account, position)` syncs the position's assets and opens an auction if it falls short. A Stock Token's price starts at its band's high edge, and WETH's at its Chainlink price, and falls a basis point a second while the market is open and a basis point every four seconds while it is closed, down to 10% below the low edge or WETH's price. An auction runs an hour while open and four hours while closed, long enough to reach its floor from a band at its 15% cap; after that, or once the market changes state, it starts again from the top. `stop` ends it once the position no longer falls short, and anyone may call it, the borrower included. `price(account, position, token)` gives the ask, in USD with 8 decimals per token.
- **A purchase.** `buy(account, position, token, amount, maxCost, receiver)` takes USDG from the buyer at the ask, at most `maxCost`, and sends the tokens to `receiver`. Half a percent of it goes to the accounts' reserve, unless the position is worth less than it owes; the rest repays the position, its debt first and its premium after. A purchase repays at most half of what the position owes, or all of it once that is 2,000 USDG or less, as Aave v3 does, and the next purchase judges the position again. While the market is closed, at most a tenth of a position's holding of a Stock Token is sold an hour, an allowance that refills continuously (`hourlyAllowance`). The contract does not rank positions: keepers start the deepest first, by `shortfall(account, position)`, which gives each position's equity and requirement.
- **Losses.** `writeOff(account, position)` writes off what a position still owes once it holds nothing, or only holdings a burn left worthless: lenders bear it, and its premium is forgiven. If the accounts refuse it and the position is worth less than it owes, holdings worth less than a USDG (`DUST`), which anyone could leave in it to block its write-off, are swept to the caller first: Stock Tokens at their band's low edge, WETH only at a fresh Chainlink price. A holding in a halted asset is not swept. Once the accounts have a backstop, only it may write off, after covering what it can.
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
- **Lots.** `enroll(account, position, symbol)`, which anyone may call, enters a position that falls short at its bands' low edges, as the accounts' `health` has it, into the round of one of its Stock Tokens; a cross position may be enrolled in the round of each Stock Token it holds. Its lot is what repays all it owes at the round's floor, fee included, at most its holding. A round keeps 32 lots; past that, a larger lot takes the smallest one's place. While the market is open, the liquidator holds an enrolled position out of its Dutch auction until the round's clearing window closes.
- **The floor.** 10% below the low edge of the band the asset's `BandFeed` sealed at the reopen (`seals`), the liquidator's floor; without a seal, or with a halted one, below the band's low edge when the round opens. No bid below it is taken.
- **Bids.** `commit(symbol, commitment, deposit)` takes a deposit of at least 100 USDG, which may exceed the bid to hide its size; the commitment is `keccak256(abi.encode(bidder, symbol, openMs, quantity, price, salt))`, so no one else can reveal it. `reveal(symbol, openMs, quantity, price, salt)` keeps from the deposit what the bid is worth at its price, its escrow, and returns the rest; a bid below the floor, worth less than 100 USDG at the floor, or worth more than its deposit is refused. A deposit not revealed loses a tenth of itself, at least 100 USDG, to the accounts' reserve (`forfeit`). A round keeps 64 bids; past that, a higher bid takes the lowest one's place, whose escrow goes back, and a bid that finds no place gets its whole deposit back (`Outbid`).
- **Clearing.** Anyone submits the clearing price with `clear(symbol, openMs, price)`, and the contract checks it in one pass over the bids: the highest bid price at which the bids at or above it take the whole supply; when all of them together take less, the lowest bid price; with no bid, the floor. Every bid above it fills, bids at it share what is left pro rata, and each lot's share of what is taken is sold through the liquidator at that one price; the gap backstop buys the rest of the lot at the same price, at most the band's low edge, so a high bid cannot set what the backstop pays; a backstop that cannot buy, its USDG wiped or its address blocked by an issuer, leaves the rest in the position and never stops the clearing. `claim(symbol, openMs, index)` then sends each bidder its tokens and the rest of its escrow. A round nobody clears in its hour refunds every bid, and its positions return to the Dutch auction.
- **What is not sold.** What neither the bids nor the backstop take stays in the positions, which return to the Dutch auction, and what an emptied position still owes goes to the gap backstop's cover. An asset halted over the week reopens into the Dutch auction: no round is keyed on a halt's end.

Reads: `round`, `lots`, `bids`, `enrolled`, `commitments`, `feeds`, `lastOpenMs`, `phase` and the constants.

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
