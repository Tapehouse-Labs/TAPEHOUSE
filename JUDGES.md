<!-- SPDX-License-Identifier: CC0-1.0 -->

# Tapehouse, for the judges

Tapehouse lends USDG against a whole portfolio of Robinhood Chain Stock Tokens, and keeps its price live while the stock market is closed. It is on Robinhood Chain testnet (46630) today, with Paxos' USDG as the only asset it lends and borrows.

**The problem.** Stock Tokens trade around the clock, but Chainlink's equity feeds go dark at the weekend, so every lender prices each asset alone, at its last close, with a conservative loan-to-value. A weekend gap is then either ignored or socialised after the fact.

**What Tapehouse does about it.** A *band* (a Stylus program) prices every asset all week from several sources and says how wide its doubt is. A *margin engine* (a Stylus program) turns the band into one requirement for the whole portfolio. A gap that exceeds it falls first on the borrower's collateral, then on a backstop with a declared limit, and a liquidation at a reopening goes to a sealed-bid auction by lots instead of a race of bots.

## Five minutes, nothing to install but Foundry

Every command reads the live testnet. Set `R=https://rpc.testnet.chain.robinhood.com`.

**1. See the band.** `quote` returns the band's state (0 halted, 1 degraded, 2 closed, 3 open), how many price legs are live, its centre, its half-width in basis points, and its lower and upper bounds, prices in USD with 8 decimals.

```bash
cast call 0xCAB25d6C9460BEE96997A1Df664b14daC2c49359 \
  "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint64)" $(cast format-bytes32-string TSLA) --rpc-url $R
```

During regular hours it answers state 3 with two legs live: Chainlink's feed and RedStone's 24/7 price. When Chainlink sleeps, at the weekend or outside regular hours, the band does not stop: it answers state 1 with the one leg that is live, and a width that says so.

**2. See the session the band reads.** The 24/5 session (0 not known, 1 closed, 2 open), NYSE's state and its next state (1 regular hours, 2 short close, 3 weekend or holiday), the time NYSE changes state and the session's next boundary, in milliseconds.

```bash
cast call 0xCAB25d6C9460BEE96997A1Df664b14daC2c49359 \
  "session()(uint8,uint8,uint8,uint64,uint64)" --rpc-url $R
```

**3. Use it.** Open <https://app.tapehouse.xyz>, connect a wallet, take test USDG from [Paxos' faucet](https://faucet.paxos.com/?network=robinhood) and a Stock Token from [Robinhood's](https://faucet.testnet.chain.robinhood.com/). The app wraps the testnet's Stock Token, deposits it, and shows the **liquidation price before anything is signed**. A smart account pays no gas for its creation and its first three operations.

**4. See that USDG is the unit.** The supply vault lends USDG, and the margin accounts borrow it.

```bash
cast call 0xDa62E570B6d1A038C1b3b011555BdA1fFd82038F "asset()(address)" --rpc-url $R   # 0x7E955252…802F, Paxos USDG
cast call 0xDa62E570B6d1A038C1b3b011555BdA1fFd82038F "totalAssets()(uint256)" --rpc-url $R
```

**5. See the reopening auction's phase.** It reads zero outside the window between the 24/5 session's reopen and the regular open (see [the table](#what-is-live-what-is-tested-what-is-not)).

```bash
cast call 0x8fcE1e7d2D5e470A5ae7569F803d2957b3A771CB "phase()(uint64,bool)" --rpc-url $R
```

**6. Check that the program on chain is the code here.** Needs Docker.

```bash
git checkout <the commit of the row in the README's verification table>
make verify-stylus CHAIN=46630 TX=0x7f003392328803b2fd7331f6dcc2d96f387e8e9b811cdc74b01e9c77ba5e1310
```

The testnet's band and margin engine are second instances of programs verified this way (`0xa70118d3…5641` and `0x9B2DB813…D6B0`), cloned while Arbitrum pauses new Stylus activations.

**7. See why Stylus, live.** The margin engine is also written in Solidity, bit for bit, and deployed beside the program (`0x68B73e2c70cd61B9d94cCF7c0F3502F9DE93749c`). One command, which reads the chain and needs no key, calls both with the same arguments, stops at the first answer that differs by a single bit, and prints the L2 gas of each:

```bash
make compare-margin CHAIN=46630
```

On the testnet's two assets it compares 32 calls (the views, scenarios, the digest of every lattice size and the requirement at three horizons) and the Stylus program is 12.3× cheaper for the requirement and 13.5× for the 256-point scenario set. On the dev node the same suite gives 13.1× for three assets, 16.7× for six and up to 20.5× for the full launch set ([README](README.md#stylus-versus-solidity)); `make devnode deploy-stylus-devnode test-stylus-devnode gas-table` reproduces it. To call one yourself: `cast call $MARGIN "requirement(int256[],uint256[],uint64,bool)(uint256,uint8)" "[1000000000000000000,-500000000000000000]" "[37800000000,26000000000]" 172800 false`, once on the program (`0x85cbe267eF0c5dEC5Fc2aFc0b439cc62caB94982`) and once on the reference.

**8. Read the evidence behind the parameters.** [`docs/BACKTEST.md`](docs/BACKTEST.md) runs the band and the margin engine over all 890 weekends and holidays since 2010 and says what happened: the requirement in force covered every closure in 400 synthetic portfolios, and no portfolio lost more than its margin at the deployed 5× weekend cap. It also states where the model falls short, such as NVDA's pool depth and the band's hit rate on the reopening. `make backtest` repeats it.

**9. Run the tests.** `make test`. Fork tests need an archive RPC for each chain (`ROBINHOOD_RPC_URL`, `ROBINHOOD_TESTNET_RPC_URL`, `ARBITRUM_RPC_URL`; see the [reference](docs/REFERENCE.md#networks)).

## What is live, what is tested, what is not

| Piece | On testnet | How you can check it |
|---|---|---|
| Band, margin engine | Live, fed every minute by a keeper | Steps 1 and 2; `quote` changes with the market |
| Margin accounts, supply vault (USDG) | Live | The app; steps 3 and 4 |
| Liquidator, gap backstop | Live | The app's Risk page; `cast call` on the backstop |
| Stock lending, short positions, baskets | Live, at the small scale the pools' depth allows | `/earn` and `/trade` in the app; the SDK's `shorts`, `stockLending`, `baskets` |
| Gap cover | Live; sales open only on the last trading day before a weekend or holiday | The cover's `sales` view (see the reference) |
| Reopening auction | Deployed; its phases need a real reopening | Step 5. Tested with `vm.warp` in the contract suite, and end to end in the app against a node front that states each phase (`apps/app/e2e/auction.spec.ts`) |
| Sponsor paymaster, bundler, smart accounts | Live | The app: a wallet with no ether creates its account and deposits for free |
| Indexer and public API | Live | <https://indexer-production-f5dc.up.railway.app/v1/status> |
| Pay per request over x402 | Live: beyond the free allowance a read is answered 402 and paid in USDG, with no key. Settled on the testnet in [`0x8e6651bc…df9345`](https://explorer.testnet.chain.robinhood.com/tx/0x8e6651bc96582ef77329e3a209982d15a2c54c8b6dcefedbb562709d26df9345) by Tapehouse's own facilitator, since the chain has none | `curl -si` the API over 10 requests a second for the 402; [the reference](docs/REFERENCE.md#paying-per-request) |
| Claude Code plugin | In the repository: the MCP server and three skills; `@tapehouse/mcp` is not on npm yet, so the plugin cannot start until it is | `claude plugin validate plugins/tapehouse` |
| SDKs, MCP server | The TypeScript SDK is published as `@tapehouse/sdk`; the Rust and Go SDKs and the MCP server are in the repository | `npm view @tapehouse/sdk` |
| **Robinhood Chain mainnet** | **Not live** | Arbitrum paused new Stylus activations on 2 October 2026, which blocks deploying the band and margin programs on any network until it lifts |
| Third-party audit | **None yet** | The reference's slither run, the fork tests and the reproducible builds are what stands in for it |

## What is simplified on the testnet, stated plainly

- **Collateral is wrapped.** The testnet's Stock Tokens predate the `oraclePaused` view the band reads, so `wTSLA` and `wAMZN` wrap Robinhood's testnet TSLA and AMZN one to one and add it. The app wraps a deposit for you.
- **Uniswap v3 is ours.** The testnet has none, so the repository deploys Uniswap's own contracts and seeds a pool for each asset against USDG. The engine reads those pools' depth, which is small, so loans are small.
- **Two assets.** TSLA and AMZN on the testnet, because those are the Stock Tokens it has. The code is written for the launch set.
- **No mocks in shipped code.** Test doubles live only in tests. The integrations are proved with fork tests against the real chain.

## Tests

| Suite | Tests |
|---|---|
| Solidity contracts, with invariants and fork tests (`forge test`) | 566 |
| Uniswap deployment for the testnet (`forge test`) | 7 |
| Stylus programs (`cargo test`) | 208 |
| Rust SDK | 33 |
| Go: services and SDK | 316 |
| TypeScript: SDK, MCP server, app and website | 285 |

The app's end-to-end tests (`make test-app-devnode`, desktop and phone Chromium, with axe's WCAG 2.2 AA checks) and the services' dev-node tests run against a local Nitro node. One fork test needs a logs RPC that serves wide ranges, and fails without it.

## Where to read next

- [`README.md`](README.md): the problem, the architecture and every deployed address.
- [`docs/REFERENCE.md`](docs/REFERENCE.md): each component in turn, with its commands.
- `contracts/src`, `stylus/contracts/band`, `stylus/contracts/margin`: the code the steps above exercise.
