---
version: 1
slug: "earn"
primary_target: "earn"
related_targets: ["app"]
---

# Earn

Scope: the Earn surface at `/earn`, both themes: supplying USDG, the gap backstop, lending Stock Tokens from the margin account, and writing gap cover where it is deployed. Visitor mode: Operate.

Audience and job: the holder with idle USDG or Stock Tokens, and the collateral optimiser comparing yields and risks. The job is to compare where money can work, read its loss order, exposure limits, exit terms and recall terms, then deposit and later leave, knowing every term before signing. Shares are held by the smart account; withdrawals and claimed Stock Tokens go to the wallet. Proof is each figure read from the chain at a block.

Constraints: the thread's world unchanged. Every read and transaction through @tapehouse/sdk; no raw ABI. The loss cascade stated plainly: borrower's collateral, then the backstop, then lenders. No lending yield named while no one shorts. Gap cover writing only where the registry deploys GapCover; elsewhere it says so.

## Direction contract

THESIS: One ledger of terms. Every place money can work is a row answering the same four questions, what it earns, how it leaves, what it bears and what you hold, so the trade between yield, exit and loss is read across, not hunted through four pages. Refuses the grid of yield cards with an APY headline each.
OWN-WORLD: The thread: bone ground, warm ink by alpha, registry navy as the 2px stem and the 2px rule over the open row, amber only for closed doors and warnings. A real table in tabular figures, label-type headings, pill controls, the bracketed seal.
STORY: The visitor scans the ledger, opens a row, reads its notes before the composer, deposits once for free, and later starts a cooldown, redeems, recalls or claims gains from the same row.
FIRST VIEWPORT: Seal top right; label and the visitor's total across places at display size under the bar with one sentence, the bar clear of the header as on the account; then the ledger off the stem, its header and open row starting in the first viewport on desktop, the row's facts, notes and composer following below it.
FORM: The ledger of terms, structure 5 of 7 on the grounded list, dealt by seed 0fde52d9.
FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
