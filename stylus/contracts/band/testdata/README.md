# Test data

Real signed data packages from the `redstone-primary-prod` data service, taken from `https://oracle-gateway-1.a.redstone.finance/data-packages/latest/redstone-primary-prod` and serialised as on-chain payloads by `stylus/scripts/redstone-payload.py`. Nothing is re-signed: tamper tests alter a copy.

The expected values are what `PrimaryProdDataServiceConsumerBase` from `@redstone-finance/evm-connector` 1.0.0 returns for the same payload at the same block timestamp.

| File | Data packages | Timestamp (ms) | Values |
|---|---|---|---|
| `nvda-24_7.hex` | `NVDA---24_7`, 5 signers | 1790119750000 | 22846110842 |
| `status.hex` | `NY_MARKET_CURRENT_STATUS`, `NY_MARKET_NEXT_STATUS`, `NY_MARKET_NEXT_CHANGE_TIME`, 5 signers each | 1790119750000 | 101000000, 100000000, 179017020000000000000 |
| `ny-market-status.hex` | `NY_MARKET_STATUS`, 5 signers, three data points each | 1790122370000 | 101000000, 100000000, 179017020000000000000 |

`quote-vectors.json` holds the band's expected outputs. Integers are decimal strings, because several exceed 64 bits.

- `quotes`: 12 cases from one-minute captures of the RedStone gateway and of Chainlink on Robinhood Chain (NVDA and TSLA, 18–21 Sep 2026), 4 synthetic cases, and 10 boundary cases.
- `ewma`: 11 variance updates: rises, falls, gaps, zero prices and the return cap.

Every expected value was computed in unbounded integer arithmetic, independently of this crate.
