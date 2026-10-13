# The tapehouse server's tools

Every read names the block it was made at. A transaction tool returns unsigned calls and names who signs them; none sends anything.

**Reads**

| Tool | Answers |
| --- | --- |
| `registry` | Every address the other tools use |
| `band_quote`, `band_session`, `band_halt` | An asset's band; the 24/5 session and NYSE's state; a trading halt |
| `band_latest_band`, `band_sealed` | The whole band from the asset's `BandFeed`; the band it sealed before a reopen |
| `chainlink_round` | A Chainlink feed's latest round |
| `morpho_oracle` | The price a Morpho Blue oracle answers, or why it has none |
| `baskets_components`, `baskets_target`, `baskets_pending_target`, `baskets_preview_mint`, `baskets_preview_redeem` | A basket's tokens, target and its previews |
| `accounts_health`, `accounts_leverage`, `accounts_collateral`, `accounts_in_baskets` | A position's equity, requirement, leverage and what it holds |
| `accounts_liquidation_price`, `accounts_repayment` | The price at which a position meets its requirement once it borrows more; what repays it |
| `accounts_recovery` | Until when a halted asset's positions cannot borrow or withdraw |
| `accounts_is_authorized`, `accounts_mandate` | Whether an address may act for an account, and where an agent stands under its mandate |
| `shorts_position`, `shorts_health`, `shorts_restriction` | A short, its health, and the short-sale restriction |
| `gap_cover_sales`, `gap_cover_quote`, `gap_cover_series`, `gap_cover_position`, `gap_cover_vault` | Weekend gap cover |

**Transactions**

| Tool | Prepares |
| --- | --- |
| `accounts_deposit`, `accounts_withdraw`, `accounts_borrow`, `accounts_repay`, `accounts_unwrap` | The margin accounts' calls, with an approval first where the allowance is short |
| `accounts_set_authorization` | Lets an address do everything for the account, until it signs the reverse. Prefer `accounts_set_mandate` |
| `accounts_set_mandate` | Gives an address a bounded authority until a time |
| `baskets_mint`, `baskets_redeem` | A basket mint or redemption in kind |
| `shorts_deposit`, `shorts_withdraw`, `shorts_sell`, `shorts_cover` | A short's USDG, a sale and a buy-back, with slippage limits the server caps |
| `gap_cover_buy`, `gap_cover_release`, `gap_cover_claim`, `gap_cover_deposit`, `gap_cover_redeem`, `gap_cover_measure` | Buying, releasing and claiming cover, and a writer's deposit and redemption |

The server answers at most 120 tool calls a minute, over every client together: plan the reads, and do not poll.
