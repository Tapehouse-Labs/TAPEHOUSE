---
version: 1
slug: "trade"
primary_target: "trade"
related_targets: ["app"]
---

# Trade

Scope: the Trade surface at `/trade`, both themes: shorting a Stock Token on USDG, buying gap cover over the coming closure where it is deployed, and minting, redeeming and holding baskets in the cross position. Visitor mode: Operate.

Audience and job: the collateral optimiser and the holder hedging a weekend. The job is to read every line of an order before signing, the quote, the limit the slippage sets, the fees, the band's bounds and the margin before and after, then send it, and later buy back, claim or redeem. The smart account trades; withdrawals, payouts and redeemed tokens go to the wallet. Proof is each figure read from the chain at a block, quotes from the pool's QuoterV2.

Constraints: the thread's world unchanged. Every read and transaction through @tapehouse/sdk; no raw ABI. Slippage starts at 50 bps, editable from 0 to 10,000; the SDK chooses none and a failed quote is an error, never a zero limit. Gap cover only where the registry deploys GapCover and during its sales. A basket enters the cross position only.

## Direction contract

THESIS: The ticket. Every order is a broker's ticket spelled out line by line in a fixed sequence, what you sell or buy, what the pool pays or asks, the limit at your slippage, the band's bound, the fees, the margin after, with the button at its foot, so nothing about the order is learned after signing. Refuses the swap widget with a single number and a gear icon.
OWN-WORLD: The thread: bone ground, warm ink by alpha, registry navy as the 2px stem and the ticket's 2px head rule, 1px rules between its lines, amber for restrictions and closed sales. Label-type line names, strong tabular figures right-aligned, pill controls, the bracketed seal.
STORY: The visitor picks short, gap cover or baskets, types one amount, reads the ticket fill in from the chain, signs once for free, and finds the trade in its book below to buy back, claim or redeem.
FIRST VIEWPORT: Seal top right; under the bar, the label and the ticket's kind switch; the ticket's first line with its amount input, then the quote and the limit as the first strong lines; the button at the ticket's foot within the viewport on desktop.
FORM: The trade ticket, structure 1 of 7 on the grounded list, dealt by seed 9f96e981.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
