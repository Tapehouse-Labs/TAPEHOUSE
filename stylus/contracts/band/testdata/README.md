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
