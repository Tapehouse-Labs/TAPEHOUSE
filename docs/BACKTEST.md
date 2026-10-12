<!-- SPDX-License-Identifier: CC0-1.0 -->

# The margin against sixteen years of weekends

Tapehouse's weekend parameters are not opinions. The margin engine, the weekend leverage cap, the gap backstop's limits and the gap cover's price are each set from, and then tested against, every closure of the market since 2010. This page is the result of one run, with the figures it quotes in [`docs/backtest-2026-10-09.json`](backtest-2026-10-09.json), and the commands to repeat it.

```bash
ROBINHOOD_RPC_URL=… ARBITRUM_RPC_URL=… make backtest   # writes services/backtest.json
```

## What was run

| | |
|---|---|
| Assets | NVDA, TSLA, AAPL, MSFT, GOOGL, SPY |
| History | 4,096 sessions, 29 June 2010 to 9 October 2026, adjusted for splits and dividends (Yahoo Finance's chart API) |
| Closures | 890 weekends and holidays, from the last close to the next open; 14 of them since Robinhood Chain launched on 1 July 2026 |
| Portfolios | 200 long-only and 200 long-short, 100,000 USD a unit, margined by the engine's own arithmetic at its first parameters |
| Fidelity | The Go code that evaluates the band and the engine reads the Stylus programs' own test vectors, so the run is the programs' arithmetic, not a model of it |

## What the closures did

The worst move of each asset from a close to the next open, in 890 closures:

| Asset | Worst closure | Date |
|---|---|---|
| NVDA | −14.7% | 28 Jan 2019 |
| TSLA | +16.8% | 17 Jan 2012 |
| AAPL | −13.0% | 16 Mar 2020 |
| MSFT | −11.9% | 16 Mar 2020 |
| GOOGL | −10.3% | 16 Mar 2020 |
| SPY | −10.4% | 16 Mar 2020 |

The market is closed 31% of the time. The longest closure was 96 hours (31 October 2012, Hurricane Sandy).

## Did the margin cover them?

- **Across every portfolio and every closure the requirement in force covered the loss.** Over the 200 long-only portfolios and the 890 closures, the requirement at the start of the ramp that rises towards a weekend covers 99.8% of the cases, and at its midpoint and its end, 100%. For the long-short portfolios, 99.99%, 100% and 100%.
- **The ramp is needed, and the page does not hide it.** The weekday requirement alone, with its 25% buffer, covers the closed requirement for 193 of 200 long-only portfolios and 198 of 200 long-short ones. Against each portfolio's own worst closure it is enough for only 7 of the long-only and 178 of the long-short; the worst shortfall is 4.4%, which is what the rise towards the closed requirement before a weekend is for.
- **Bad debt in sixteen years: none.** At the margin accounts' requirement across a closure, the engine's loss rate over the 890 closures is 0, and so is Regulation T's flat 50%'s. The difference is capacity: the engine lets a portfolio borrow up to 80% of its collateral's value against 50%.

## The weekend leverage cap

- **The rule.** The cap is the largest whole leverage L with 1/L at least 1.25 times the worst move of any launch asset. With a worst move of 16.8%, the rule gives **4×**.
- **What is deployed.** The accounts hold **5×**. At 5×, 0 of the 400 portfolios lose more than their margin in any of the 890 closures, and the capacity is 80%, against 75% at 4×. The rule is the conservative figure, and the 5× deployed cap is a choice the history supports and does not prove.
- **Without any cap,** 193 long-only and 22 long-short portfolios would have lost more than their margin, in four closures: 24 August 2015, 9 March 2020, 16 March 2020 and 5 August 2024.
- **16 March 2020.** Counted as a scenario in the engine, it leaves 88.1% capacity and costs the long-short portfolios no more margin; handled by the cap alone, capacity is 80.0% and the long-short portfolios pay 13.2% more margin.

## Is the two-day horizon enough?

The engine diffuses prices over two days. A closure's expected shortfall at 99%, squared, measured in days of a normal diffusion at each asset's volatility:

| NVDA | TSLA | AAPL | MSFT | GOOGL | SPY |
|---|---|---|---|---|---|
| 1.85 | 1.33 | 2.70 | 1.71 | 1.38 | 2.68 |

The tails of AAPL and SPY are longer than two days of diffusion. That is why the weekend-gap scenarios and the ramp exist beside the horizon, and the coverage above, over every closure, is the test of all of them together.

## The gap backstop and the gap cover

- **The backstop.** For a debt of 100,000 USD, the exposure limits the history asks for are 954 USD for the cross positions and 4,079, 4,260, 2,079 and 838 USD for NVDA, TSLA, AAPL and MSFT; GOOGL and SPY need none. In 5.07 years of closed market the backstop would have absorbed 1,635 USD, so a weekend premium of 36 basis points a year pays for it.
- **The gap cover's price,** checked on weekends it was not fitted on. In 61% of the 890 weekends the gap the cover prices came from the prior week's move and not from its floor. The 24/5 reopening's realised move exceeded the priced gap in 0.17% of asset-weekends under uniform demand, 0.06% when demand follows the prior week's move, and 0.02% when it follows its square; the engine's gap alone would have been exceeded 0.09% of the time. On the 14 weekends since launch, never.

## The band against the reopening

For each closure since launch, the band as it would have quoted when the closure's last public 24/7 value was written, and whether the first Chainlink round after the reopening fell inside it. RedStone keeps public 24/7 values for 30 days, so only 24 of the 84 asset-closures have one; the rest are counted, not quoted.

| Chain | Reopenings evaluated | Inside the band | Median error | Median lead |
|---|---|---|---|---|
| Robinhood Chain | 24 of 84 | 19 | 15.4 bps | 3,580 s |
| Arbitrum One | 24 of 84 | 11 | 40.2 bps | 1,796 s |

That is 79% on Robinhood Chain and 46% on Arbitrum One, on small samples. The band is the price with a declared doubt, written minutes before the close; what lies beyond it is what the margin, the backstop and the reopening auction are for. From deployment on, the band's own `PriceWritten` events keep every closure, so the sample grows with each weekend.

## What it costs to liquidate

- **The decay.** The liquidator's price falls 60 basis points a minute with the market open, reaching its floor in 16.7 minutes, and 15 a minute with it closed, in 66.7.
- **The pools.** At block 87,026,892 on Robinhood Chain, the engine's linear liquidity charge against what each pool's quoter would have cost, at a tenth, half and all of the depth:

| Asset | Engine | Pool |
|---|---|---|
| NVDA | 0.55% / 2.55% / 5.05% | 0.96% / 32.51% / 66.26% |
| TSLA | 0.80% / 2.80% / 5.30% | 0.50% / 1.57% / 2.94% |
| AAPL | 0.55% / 2.55% / 5.05% | 0.16% / 0.86% / 1.90% |
| MSFT | 0.80% / 2.80% / 5.30% | 1.05% / 3.59% / 6.53% |
| GOOGL | 0.55% / 2.55% / 5.05% | 0.11% / 0.99% / 2.06% |
| SPY | 0.55% / 2.55% / 5.05% | 0.28% / 0.97% / 1.74% |

The charge covers the pools of TSLA, AAPL, GOOGL and SPY. It does **not** cover NVDA's, which cost far more, nor MSFT's at any of the three depths, by 0.25 to 1.2 points, at that block. That is the figure the calibrator re-measures, and why a depth falls to its measured value at once while it rises at most ×1.5 a step; Robinhood Chain's mainnet pools are not ours, and the testnet's are small.

## What this does not show

- **A backtest is history.** A closure worse than every one since 2010 would still be a bad debt; the backstop and the auction exist for it.
- **The portfolios are synthetic.** Two hundred of each kind, at 100,000 USD a unit, not real users.
- **Six assets, the launch set.** The testnet runs two of them.
- **The band's accuracy is measured on 24 asset-closures per chain,** not on a decade.

The same run feeds `make calibrate`: a proposal for new parameters is refused if, over this history, it leaves more portfolios losing more than their margin, or a higher bad-debt rate, than the parameters in force.
