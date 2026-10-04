---
version: 1
slug: "account"
primary_target: "account"
related_targets: ["app"]
---

# Account and borrow

Scope: the account surface at `/` once a wallet connects, and the public margin simulator at `/simulator`, both themes. Visitor mode: Operate.

Audience and job: the holder who will not sell, with Stock Tokens in their wallet and no ether, and the collateral optimiser. The job is to see portfolio, capacity and liquidation price, then deposit, borrow, repay and withdraw knowing every number before signing. The margin account is the smart account their wallet owns; Tapehouse's paymaster pays its first operations. Proof is each figure read from the chain at a block; none is illustrative.

Constraints: the thread's world unchanged. Every read and transaction through @tapehouse/sdk; every transaction simulated before the wallet is asked; a revert explained in words. Regimes in words: unknown (counted as closed), closed, open (25% buffer), closing. No ET clock: times are hours to the band's boundary.

## Direction contract

THESIS: The stem is the clock. The account's week runs down it, now, the ramp, the close, the 24/5 reopen, the regular open, with the requirement against equity at each, so the weekend is a place the visitor sees before it arrives. Refuses the tabbed borrow widget with a health gauge.
OWN-WORLD: The thread: bone ground, warm ink by alpha, registry navy as 2px stem and the requirement fills' edge, liquidation red only on the figure, amber only for the ramp. Pill controls, tabular figures, the bracketed seal beside every money figure.
STORY: The visitor reads their liquidation price, types a loan and watches the figure and every point of the week move, signs once, and later repays, lends or recalls from the same column.
FIRST VIEWPORT: Seal top right; label and liquidation price at display size under the bar; one composer line, action switch, amount and the primary button; then the week markers hanging off the stem, each with a requirement bar against equity; positions below.
FORM: The week, structure 4 of 7 on the grounded list, dealt by seed 7b0fae1c.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
