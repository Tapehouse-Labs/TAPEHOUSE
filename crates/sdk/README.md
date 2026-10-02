# tapehouse-sdk

Typed reads and transactions for [Tapehouse](https://github.com/Tapehouse-Labs/TAPEHOUSE)'s price band, its feeds, the margin accounts, the short positions and the Morpho oracles, on [alloy](https://alloy.rs).

```toml
[dependencies]
tapehouse-sdk = "0.1.0"
```

`Tapehouse` takes a provider and a chain's registry, `deployments/<chainId>.json` from the Tapehouse repository, loaded with `Deployments::load`: the crate compiles in no address. Each accessor returns alloy's contract instance, generated with `sol!` from the contracts' ABIs, whose methods build typed calls and transactions.

```rust
use alloy::primitives::U256;
use tapehouse_sdk::{Deployments, Tapehouse, to_bytes32};

let tapehouse = Tapehouse::new(provider, Deployments::load("deployments/4663.json")?);
let quote = tapehouse.band()?.quote(to_bytes32("NVDA")?).call().await?;
let one = U256::from(10u64.pow(18));
let sale = tapehouse.quote_sale("NVDA", one, 50).await?;
let receipt = tapehouse
    .shorts()?
    .sell(to_bytes32("NVDA")?, one, sale.min_proceeds, account)
    .send()
    .await?
    .get_receipt()
    .await?;
```

- `band`, `band_feed`, `accounts`, `shorts`, `usdg`, `stock_token`, `token_price` and `share_price`: the contracts, at the registry's addresses. A `.chainlink` feed is typed by what it prices: `TokenPriceFeed` on Robinhood Chain, `SharePriceFeed` on Arbitrum One.
- `morpho_oracle` and `oracle_price`: an asset's Morpho oracle from `.morphoOracles`, and its `price()` as `OraclePrice`: a price, or no price, never zero, where it reverts with `NoAnswer` or `SequencerNotSettled`, with the reason (`NoPrice::Stale`, `Halted` or `SequencerNotSettled`) read at one block.
- `write_prices`: the band's `writePrices`, with the signed RedStone packages a `PackageSource` of your own supplies. The crate holds no API key.
- `quote_sale` and `quote_cover`: the limits of a sale and a buy-back, from Uniswap's QuoterV2, with a slippage of 0 to 10,000 basis points.
- `repayment`: the debt and premium a full repayment takes, both read at one block; `collateral`, `leverage` and `liquidation_price` of a position.
- `Error`: `Registry`, `Argument` (a slippage above 10,000 basis points, or a name longer than 32 bytes), `Revert`, `Contract` and `Source`, with the underlying error as its `source()`.
- `decode_revert` and `Error::Revert`: the name and arguments of any revert of the contracts, the band, the Stock Tokens, USDG, or Solidity's `Error(string)` and `Panic`.

`examples/devnode.rs` runs all of it against a local Nitro dev node; it reads the account owner's key from `PRIVATE_KEY` in its environment, never from its command line.

## License

Licensed under either of [Apache License, Version 2.0](LICENSE-APACHE) or [MIT license](LICENSE-MIT) at your option.
