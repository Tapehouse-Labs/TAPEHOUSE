# Concepts

**Band.** A Stylus program prices each asset from Chainlink, RedStone's 24/7 feed and Uniswap v3 TWAPs and says how wide its doubt is: its centre is the live 24/7 leg, else Chainlink, never an average with a sleeping leg. `band_quote` returns:

| `state` | Meaning |
| --- | --- |
| `halted` | No live leg, or trading halted: every figure is zero and the asset counts for nothing |
| `degraded` | Fewer live legs than expected, or Chainlink's session not known: a wider band |
| `closed` | Chainlink's 24/5 session is closed |
| `open` | Chainlink's 24/5 session is open and the legs are live |

`halfWidthBps` is the band's half-width; `low` and `high` are its edges. Lending uses the low edge.

**Margin engine.** A second Stylus program walks a lattice of price scenarios with pairwise correlations, each asset's weekend gap and the depth of its Uniswap pool, and returns one requirement for the whole portfolio. `accounts_health` returns equity, requirement and the difference.

**Weekend.** The requirement rises towards a close and the leverage cap is 5x while the market is closed. A shortfall falls on the borrower's collateral first, then on the gap backstop up to its exposure limit per position, then on lenders.

**Reopening auction.** Positions that fall short over a closure are sold in a sealed-bid, uniform-price auction by lots at the regular open: commit, reveal, clear, claim.

**Recovery gate.** After an asset was halted or had a multiplier change pending, positions holding it cannot borrow or withdraw against it for 30 minutes. `accounts_recovery` says until when; the revert is `AssetRecovering`.

**Reverts you will meet.** `Unauthorized(caller, account)`: the caller was never authorized or given a mandate. `AssetRecovering`: see above. `MandateExpired`, `MandateDebtLimit`, `MandateLeverageLimit`, `MandateReceiver`, `InvalidMandate`: see the `tapehouse-mandate` skill. `SalesClosed`: gap cover is on sale only on the last trading day before a weekend or holiday.

**USDG.** Paxos' Global Dollar is the only asset Tapehouse lends and borrows, and the unit of the auction, the backstop and gap cover.
