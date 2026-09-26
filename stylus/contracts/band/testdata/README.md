# Test data

Real signed data packages from the `redstone-primary-prod` data service, taken from `https://oracle-gateway-1.a.redstone.finance/data-packages/latest/redstone-primary-prod` and serialised as on-chain payloads by `stylus/scripts/redstone-payload.py`. Nothing is re-signed: tamper tests alter a copy.

The expected values are what `PrimaryProdDataServiceConsumerBase` from `@redstone-finance/evm-connector` 1.0.0 returns for the same payload at the same block timestamp.

| File | Data packages | Timestamp (ms) | Values |
|---|---|---|---|
| `nvda-24_7.hex` | `NVDA---24_7`, 5 signers | 1790119750000 | 22846110842 |
| `status.hex` | `NY_MARKET_CURRENT_STATUS`, `NY_MARKET_NEXT_STATUS`, `NY_MARKET_NEXT_CHANGE_TIME`, 5 signers each | 1790119750000 | 101000000, 100000000, 179017020000000000000 |
| `ny-market-status.hex` | `NY_MARKET_STATUS`, 5 signers, three data points each | 1790122370000 | 101000000, 100000000, 179017020000000000000 |

`quote-vectors.json` holds the band's expected outputs. Integers are decimal strings, because several exceed 64 bits.

- `quotes`: 12 cases from one-minute captures of the RedStone gateway and of Chainlink on Robinhood Chain (NVDA and TSLA, 18–21 Sep 2026), 4 synthetic cases, 10 boundary cases, and 5 SPY cases priced from the S&P 500 index, built on a capture of 25 Sep 2026.
- `ewma`: 11 variance updates: rises, falls, gaps, zero prices and the return cap.
- `index_legs` and `anchors`: SPY's 24/7 leg from its anchor, and when a new Chainlink print re-anchors it, around SPY's print of 25 Sep 2026 16:03:00 UTC.

`session-vectors.json` holds Chainlink's 24/5 session derived from the signed New York market status.

- `sessions`: 38 cases. Real status values signed around the weekend of 18–21 Sep 2026, Labor Day 2026 and a weekday night, and cases marked synthetic where no capture exists yet: both daylight-saving changes, Thanksgiving, an early close and a holiday eve signed as a short close. Also stale, mixed and malformed status.
- `closes`: which written status records the regular close.

Over 30 days of the signed status (27 Aug–25 Sep 2026, four weekends including Labor Day), the same rule places all 7,826 Chainlink rounds of the 35 equity feeds on Robinhood Chain inside the session or within 2 minutes after its close, and the first round after each reopen lands within 25 s.

Every expected value was computed in unbounded integer arithmetic, independently of this crate.

`multiplier-vectors.json` holds the ERC-8056 multiplier rules. Real cases come from Robinhood Chain: NVDA's multiplier applied to the Arbitrum One share price reproduces the Robinhood Chain answer exactly, and SPY's step of 18 Sep 2026 (`effectiveAt` 1789690233, multiplier 1.001717991187472003) turns the share price 760.2447 into the first round after it, 761.55079369. Synthetic cases, labelled as such by their names, cover a 4:1 split, a split followed by a dividend before Chainlink confirms it, a step nobody recorded, a 2.15% dividend, a 2% special dividend, and the edges of the arithmetic.
