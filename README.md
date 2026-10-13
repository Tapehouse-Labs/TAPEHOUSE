<p align="center">
  <img src="apps/landing/app/icon.svg" alt="Tapehouse" width="96" />
</p>

<h1 align="center">Tapehouse</h1>

<p align="center">
  <strong>Portfolio margin for Stock Tokens on Robinhood Chain</strong>
</p>

<p align="center">
  Your whole portfolio borrows more, and your collateral earns a fee while it backs the loan.
</p>

<p align="center">
  <a href="https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml"><img src="https://github.com/Tapehouse-Labs/TAPEHOUSE/actions/workflows/ci.yml/badge.svg" alt="CI" /></a>
</p>

<p align="center">
  <a href="https://app.tapehouse.xyz">Live app</a> · <a href="#how-it-works">How It Works</a> · <a href="JUDGES.md">For the judges</a> · <a href="#stylus-versus-solidity">Stylus</a> · <a href="#deployed-contracts">Contracts</a> · <a href="#getting-started">Run Locally</a>
</p>

---

## The Problem

Stock Tokens trade around the clock, but the prices that lend against them do not.

- Chainlink's equity feeds on Robinhood Chain go dark at the weekend, and their heartbeat is suspended with them.
- The weekend is still a full trading period: about 30% of SPY's Stock Token transfers on the chain happen then.
- Lending protocols answer by pricing every asset on its own, at its last close, with conservative loan-to-value ratios. A portfolio of five assets is treated as five separate risks.

| | Today | Tapehouse |
|---|---|---|
| Weekend price | Last close, or a pause | A deterministic band from several sources, valid all week |
| Margin | One asset at a time | One requirement for the whole portfolio, with correlations |
| Gap loss | Socialised after the fact | Priced in, then absorbed by a declared backstop with a stated limit |
| Liquidation | Bots racing at the reopening | A sealed-bid auction by lots when the market reopens |
| Collateral | Idle | Earns a lending fee while it backs the loan |

---

## How It Works

```
   Chainlink · RedStone · Uniswap v3 TWAP
                    │
                    ▼
        ┌───────────────────────┐        ┌────────────────────────────┐
        │  Band  (Stylus)       │        │  Keepers · Indexer · API   │
        │  price ± width,       │◄───────│  prices, halts, history    │
        │  session, regime      │        └────────────────────────────┘
        └──────────┬────────────┘
                   ▼
        ┌───────────────────────┐
        │  Margin engine        │  scenario lattice, correlations,
        │  (Stylus)             │  weekend gap, pool depth
        └──────────┬────────────┘
                   ▼
   ┌─────────────────────────────────────────────────────────────┐
   │  Margin accounts  ·  Supply vault (USDG)  ·  Stock lending  │
   └──────────┬──────────────────────────────────┬───────────────┘
              ▼                                  ▼
   ┌─────────────────────┐            ┌────────────────────────┐
   │  Liquidator         │───────────►│  Gap backstop          │
   │  Reopening auction  │            │  junior capital with a │
   │  (sealed bids)      │            │  declared exposure cap │
   └─────────────────────┘            └────────────────────────┘
```

1. **The band**, a Stylus program, turns the sources into a price, a width and a regime for each asset, including the 24/5 session and trading halts.
2. **The margin engine**, a Stylus program, walks a lattice of price scenarios for the whole portfolio, with pairwise correlations, the asset's weekend gap and the depth of its pool, and returns one requirement.
3. **Margin accounts** hold the positions, the **supply vault** lends USDG against them, and **stock lending** pays the collateral's owner a fee.
4. When an account falls short, the **liquidator** sells it; at a reopening it goes to a **sealed-bid auction by lots**.
5. A shortfall beyond the collateral falls to the **gap backstop**, whose exposure is capped per position.
6. **Keepers** push prices and halts, the **indexer and API** serve history, and **smart accounts** pay the gas of a user's first operations.

---

## Live Demo

**<https://app.tapehouse.xyz>**, on Robinhood Chain testnet (46630).

| Page | What it shows |
|---|---|
| `/` Account | Live bands for every asset, deposit, borrow and repay, and the liquidation price before anything is signed |
| `/earn` | The supply vault, the gap backstop and the other places to put USDG to work |
| `/risk` | The backstop's exposure against its limit, the reopening auction, and each band's TWAP and premium |
| `/simulator` | What a position needs to survive a move |

Test USDG comes from [Paxos' faucet](https://faucet.paxos.com/?network=robinhood), linked from the app's header. A smart account pays no gas for its creation and its first three operations.

---

## Judge Walkthrough

[`JUDGES.md`](JUDGES.md) is the five-minute path, with commands that read the live testnet, and states what is live, what is tested and what is not. The table below says where each claim is verified.

| Claim | Where to verify it |
|---|---|
| Stylus is justified | [The table below](#stylus-versus-solidity): the margin program against a bit-exact Solidity reference, on the same dev-node calls |
| The band and the engine are deterministic | Fixed vectors, and a Solidity reference that agrees to the bit: `make test-stylus` |
| The deployed programs are the code in this repository | [Verification](#verification): rebuild a deployment and compare it byte for byte |
| It runs on a public network | The [live app](https://app.tapehouse.xyz) and the [addresses below](#deployed-contracts) |
| Every claim is tested | The CI badge above, and `make test` |
| The weekend parameters come from history | [`docs/BACKTEST.md`](docs/BACKTEST.md): 890 closures since 2010, with where the model falls short |
| It can be integrated | [SDKs](#services-and-sdks) in TypeScript, Rust and Go, and an MCP server for agents |

---

## Stylus versus Solidity

`contracts/test/reference/MarginReference.sol` is the margin engine in Solidity: the same scenario set, requirement, session and weekend cap. It answers every view of the Stylus program to the bit. The dev-node suite measures both on the same calls, in L2 gas. `make devnode deploy-stylus-devnode test-stylus-devnode gas-table` prints the run's table; this one is from `stylus/.gas-devnode`, and `make gas` fails if they differ:

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

- **Stylus wins wherever the scenario set is walked**: 13.1× cheaper for the requirement of three assets and 16.7× for six, and up to 20.5× for the six launch assets' full set.
- **Stylus loses the small views**, at about twice the gas: entering an uncached program costs about 95,000, more than reading a few words. Net of that entry cost the three-asset requirement is about 20× cheaper.
- The method, the compiler settings and the cost of a real pool read are in the [reference](docs/REFERENCE.md#stylus-versus-solidity).
- **Check it live.** The Solidity reference is deployed on the testnet, and `make compare-margin CHAIN=46630` calls it and the margin program with the same arguments, fails if a single bit differs, and prints the gas of each.

---

## Deployed Contracts

Robinhood Chain testnet (46630). Each address links to its page on the [explorer](https://explorer.testnet.chain.robinhood.com).

| Contract | Address |
|---|---|
| Band (Stylus) | [`0xCAB25d6C9460BEE96997A1Df664b14daC2c49359`](https://explorer.testnet.chain.robinhood.com/address/0xCAB25d6C9460BEE96997A1Df664b14daC2c49359) |
| Margin engine (Stylus) | [`0x85cbe267eF0c5dEC5Fc2aFc0b439cc62caB94982`](https://explorer.testnet.chain.robinhood.com/address/0x85cbe267eF0c5dEC5Fc2aFc0b439cc62caB94982) |
| Margin accounts | [`0xcE8C3CB4E1A8b944BC607353055cE5a3D07eC101`](https://explorer.testnet.chain.robinhood.com/address/0xcE8C3CB4E1A8b944BC607353055cE5a3D07eC101) |
| Supply vault | [`0xDa62E570B6d1A038C1b3b011555BdA1fFd82038F`](https://explorer.testnet.chain.robinhood.com/address/0xDa62E570B6d1A038C1b3b011555BdA1fFd82038F) |
| Liquidator | [`0x53579Ef2b4F36E8E84F9776239db29156e4465aE`](https://explorer.testnet.chain.robinhood.com/address/0x53579Ef2b4F36E8E84F9776239db29156e4465aE) |
| Gap backstop | [`0x32acef0BBCB78b88d82b672DFE3F325105dAEd8E`](https://explorer.testnet.chain.robinhood.com/address/0x32acef0BBCB78b88d82b672DFE3F325105dAEd8E) |
| Reopening auction | [`0x8fcE1e7d2D5e470A5ae7569F803d2957b3A771CB`](https://explorer.testnet.chain.robinhood.com/address/0x8fcE1e7d2D5e470A5ae7569F803d2957b3A771CB) |
| Stock lending, TSLA | [`0xD7921e0C3D475647a384735a64aaB56C30fdEd03`](https://explorer.testnet.chain.robinhood.com/address/0xD7921e0C3D475647a384735a64aaB56C30fdEd03) |
| Stock lending, AMZN | [`0xF12432A8B6263D206D5d482ac1FA4bE29744D1A1`](https://explorer.testnet.chain.robinhood.com/address/0xF12432A8B6263D206D5d482ac1FA4bE29744D1A1) |
| Margin engine, Solidity reference | [`0x68B73e2c70cd61B9d94cCF7c0F3502F9DE93749c`](https://explorer.testnet.chain.robinhood.com/address/0x68B73e2c70cd61B9d94cCF7c0F3502F9DE93749c) |
| Short positions | [`0xadc174ca32EBD4B27579E3759903F43b8F53D97F`](https://explorer.testnet.chain.robinhood.com/address/0xadc174ca32EBD4B27579E3759903F43b8F53D97F) |
| Basket `thPAIR` (TSLA and AMZN) | [`0x31B777cCC84236B6f53882890ABC7CD19b4c9575`](https://explorer.testnet.chain.robinhood.com/address/0x31B777cCC84236B6f53882890ABC7CD19b4c9575) |
| Gap cover | [`0x482434d135eeA48797c1d960DF0f3A01E434ad17`](https://explorer.testnet.chain.robinhood.com/address/0x482434d135eeA48797c1d960DF0f3A01E434ad17) |
| Sponsor paymaster | [`0x3721b8136593BCC78787943664cF9404d12Ecbd2`](https://explorer.testnet.chain.robinhood.com/address/0x3721b8136593BCC78787943664cF9404d12Ecbd2) |
| Wrapped TSLA (`wTSLA`) | [`0x588ecCac8c70cdE770b57BA9A192CfF8D636534e`](https://explorer.testnet.chain.robinhood.com/address/0x588ecCac8c70cdE770b57BA9A192CfF8D636534e) |
| Wrapped AMZN (`wAMZN`) | [`0x04D3d2Bb86BbeE4b4AeD29474fB30D7b639Aa381`](https://explorer.testnet.chain.robinhood.com/address/0x04D3d2Bb86BbeE4b4AeD29474fB30D7b639Aa381) |
| TSLA price feed (RedStone round feed) | [`0x64C39E726207689A6cdf5E688520F850E70cE617`](https://explorer.testnet.chain.robinhood.com/address/0x64C39E726207689A6cdf5E688520F850E70cE617) |
| AMZN price feed (RedStone round feed) | [`0x9eF3F38921CeB70969dE21A812727A59b01dF644`](https://explorer.testnet.chain.robinhood.com/address/0x9eF3F38921CeB70969dE21A812727A59b01dF644) |
| Uniswap v3 factory | [`0x843788502e38cD2BeCf3B1A09f8CF30353f3a983`](https://explorer.testnet.chain.robinhood.com/address/0x843788502e38cD2BeCf3B1A09f8CF30353f3a983) |
| Uniswap v3 swap router | [`0x0A8c14340c3f59aF1B1e2C834135e1477bae2C64`](https://explorer.testnet.chain.robinhood.com/address/0x0A8c14340c3f59aF1B1e2C834135e1477bae2C64) |
| Uniswap v3 quoter | [`0x905a4ca58Fed98Bc21639DC651db2271F7E3B9e0`](https://explorer.testnet.chain.robinhood.com/address/0x905a4ca58Fed98Bc21639DC651db2271F7E3B9e0) |
| wTSLA / USDG pool, 0.05% | [`0x3B5a8161d8A352b6C2cE2f358eE3413eDbdAdAfF`](https://explorer.testnet.chain.robinhood.com/address/0x3B5a8161d8A352b6C2cE2f358eE3413eDbdAdAfF) |
| wAMZN / USDG pool, 0.05% | [`0x951f2226CfBcc7cEC0347a7be10a78082D3078D2`](https://explorer.testnet.chain.robinhood.com/address/0x951f2226CfBcc7cEC0347a7be10a78082D3078D2) |

`deployments/<chainId>.json` lists every contract the repository uses on each chain, and fork tests check each entry against the chain.

### What is not live yet

The testnet's Stock Tokens predate the `oraclePaused` view the band reads from a Stock Token, so the testnet's collateral is the wrapped token: `wTSLA` and `wAMZN` wrap Robinhood's testnet TSLA and AMZN one to one, forward their multiplier, pause and issuer registry, and add `oraclePaused`. The app wraps a deposit's Stock Token first, from the wallet. On Robinhood Chain the Stock Tokens are deposited as they are.

The testnet has no Uniswap v3, so the repository deploys Uniswap's factory, a single-pool swap router and quoter, and a pool for each asset against USDG, seeded with the USDG the team holds. The margin engine reads each pool's depth, which is small, so the testnet's loans are small.

The band and the margin engine on the testnet are second instances of the programs already active on the chain, deployed with `make clone-stylus`. Arbitrum paused new Stylus activations on 2 October 2026, so the programs the repository now holds, which read `tokenPaused` when a token has no `oraclePaused`, wait for activations to reopen; the wrapper makes the testnet run without them.

Stock lending, short positions and the basket are deployed on the testnet and work on chain, at the small scale the pools' depth allows. The app offers deposits, loans and the supply vault; it does not offer them yet. The gap cover is deployed, and the testnet has no Chainlink feed to settle it on, so the band's Chainlink feeds are RedStone round feeds: each follows RedStone's regular-hours price of the asset, signed by the data service that signs the band's 24/7 feeds, checked on chain, and scaled by the Stock Token's multiplier as Robinhood Chain's Chainlink feeds are. The band and the cover read them as they read Arbitrum One's regular-hours Chainlink feeds, which they are written for. They are not independent of RedStone's signers as Chainlink's are, and on Robinhood Chain the band and the cover take Chainlink's. Cover is on sale from the last trading day before a weekend or holiday, so it opens on the testnet on Friday. The testnet's pools charge 0.05%: a short sells only at the band's low edge or above, which a band as narrow as 0.3% leaves no room for in a 0.3% pool. Nothing but the pools keeper arbitrages them, so their price follows the band's only as closely as it trades them back. The app is not served to embargoed countries and regions, by Vercel's geolocation (451), and sends nothing for an address on OFAC's sanctions list, taken from the Treasury's file; the issuer's blocks on a connected address are shown on connecting. The mainnet deployment follows the testnet's.

---

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

---

## Services and SDKs

| Component | Path | What it is |
|---|---|---|
| TypeScript SDK | `packages/sdk` | Typed reads and transactions on viem, published as [`@tapehouse/sdk`](https://www.npmjs.com/package/@tapehouse/sdk) |
| Rust SDK | `crates/sdk` | The same surface on alloy |
| Go SDK | `services/sdk` | The same surface on go-ethereum |
| MCP server | `packages/mcp` | Lets an AI agent read the band and the engine |
| Claude Code plugin | `plugins/tapehouse` | The MCP server and three skills for an agent: read, act under a mandate, pay the API per request |
| Indexer and API | `services/cmd/indexer` | History, public API, WebSocket stream and RedStone relay |
| Keepers | `services/cmd/keeper` | Prices, halts and the weekend's routine work |
| Reference bidder | `services/cmd/bidder` | A bidder for the reopening auction, as an example for others |
| Sponsor | `services/cmd/sponsor` | Signs the user operations the paymaster pays for |
| x402 facilitator | `services/cmd/facilitator` | Settles in USDG the payments of clients that read the API beyond its free allowance, with no key |
| Alerts | `services/cmd/alerts` | Warns holders before a weekend and before a liquidation |
| Backtest and calibration | `services/cmd/backtest`, `services/cmd/calibrate` | Replays history and proposes parameters |

---

## Getting Started

```bash
git clone https://github.com/Tapehouse-Labs/TAPEHOUSE.git tapehouse
cd tapehouse
make build
make test
```

`make build` initialises the submodules and installs JavaScript dependencies on first run, and stops with a message if Node, Foundry or cargo-stylus does not match the pinned versions. The requirements, every `make` target and each component's commands are in the [reference](docs/REFERENCE.md).

To use Tapehouse from Claude Code, install its plugin, which starts the MCP server and teaches the agent to read a position, act under a mandate and pay the API per request:

```bash
/plugin install tapehouse --marketplace Tapehouse-Labs/TAPEHOUSE
```

```bash
make devnode                   # a local Nitro node, ArbOS 61
make deploy-contracts-devnode  # the contracts and programs, on it
make test-app-devnode          # the app, driven in Chromium against that deployment
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Price band and margin engine | Rust on Arbitrum Stylus, `cargo-stylus` 0.10.9 |
| Contracts | Solidity 0.8.37, Foundry, OpenZeppelin |
| Oracles and pools | Chainlink Data Feeds, RedStone, Uniswap v3 |
| Smart accounts | ERC-4337 (EntryPoint v0.7), Alto bundler, a signed sponsor paymaster |
| App | Next.js, wagmi, viem, permissionless, next-intl (English, Spanish, Portuguese) |
| Services | Go, go-ethereum, SQLite |
| Chains | Robinhood Chain (4663) and its testnet (46630), Arbitrum One (42161) |

---

## Repository Layout

- `apps/landing` — the website at tapehouse.xyz (Next.js)
- `apps/docs` — the documentation at docs.tapehouse.xyz (Next.js): guides, SDK and contract reference, and the indexer API reference generated from its OpenAPI specification
- `apps/app` — the app (Next.js, wagmi)
- `packages/brand` — the design tokens and base styles the landing and the app share
- `contracts` — Solidity contracts (Foundry)
- `stylus` — Stylus programs (Rust): `band`, the price band, and `margin`, the portfolio margin engine, with the code they share in `stylus/crates`
- `deployments` — contract addresses per chain, one JSON file per chain ID
- `packages/sdk` — the TypeScript SDK (viem)
- `packages/mcp` — the MCP server for AI agents, on the TypeScript SDK
- `plugins/tapehouse` — the Claude Code plugin: the MCP server and three skills
- `services` — the Go module: the Go SDK in `services/sdk` (go-ethereum), the indexer and public API in `services/cmd/indexer`, the backtester in `services/cmd/backtest`, a reference bidder in `services/cmd/bidder`, the liquidation-risk alerts in `services/cmd/alerts`, the parameter calibrator in `services/cmd/calibrate`, the sponsor service in `services/cmd/sponsor` and the x402 facilitator in `services/cmd/facilitator`
- `crates` — Rust crates: the Rust SDK in `crates/sdk` (alloy)

The [reference](docs/REFERENCE.md) documents each component in turn.

---

## License

Licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or <https://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](LICENSE-MIT) or <https://opensource.org/licenses/MIT>)

at your option.

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this repository by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.

The repository follows the [REUSE](https://reuse.software) specification: every file states its license, in an SPDX header or in `REUSE.toml`, which also records the origin and license of the few files that come from others. The landing's generated images are CC0-1.0: machine-generated images may carry no copyright, and whatever rights the contributors hold are waived. `make lint` runs `reuse lint`.

The libraries in `contracts/lib` keep their own licenses. eth-infinitism's `account-abstraction`, GPL-3.0, is compiled only by the tests and by the dev node's deployment of the EntryPoint and its factory; `SponsorPaymaster` builds on OpenZeppelin's `Paymaster` and does not include it. `BandFeed` compiles in Uniswap v3's `OracleLibrary`, `TickMath` and `IUniswapV3Pool`, which are GPL-2.0-or-later. Taking `BandFeed.sol` under MIT, its deployed bytecode is a combined work under GPL-2.0-or-later, whose source is this repository; each deployed feed's exact source is the commit that deployed it.
