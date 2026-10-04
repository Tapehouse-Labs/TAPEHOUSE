// SPDX-License-Identifier: MIT OR Apache-2.0
//! Bindings of the contracts' ABIs, `abi/*.json`, which `make bindings` exports from the Solidity build.
#![allow(missing_docs, clippy::too_many_arguments)]

use alloy::sol;

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Aggregator,
    "abi/Aggregator.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Band,
    "abi/Band.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Basket,
    "abi/Basket.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    BandFeed,
    "abi/BandFeed.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    GapBackstop,
    "abi/GapBackstop.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    GapCover,
    "abi/GapCover.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Liquidator,
    "abi/Liquidator.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Margin,
    "abi/Margin.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    MarginAccounts,
    "abi/MarginAccounts.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    MorphoBandOracle,
    "abi/MorphoBandOracle.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    MorphoBlue,
    "abi/MorphoBlue.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    QuoterV2,
    "abi/QuoterV2.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    ReopeningAuction,
    "abi/ReopeningAuction.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    ShortPositions,
    "abi/ShortPositions.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    SponsorPaymaster,
    "abi/SponsorPaymaster.json"
);

/// `StockLendingVault`'s ABI names `SupplyVault`'s `RateModel`, for which `sol!` declares a `SupplyVault` module of its
/// own, so the binding is declared apart from `SupplyVault`'s.
mod stock_lending_vault {
    use alloy::sol;

    sol!(
        #[sol(rpc)]
        #[derive(Debug)]
        StockLendingVault,
        "abi/StockLendingVault.json"
    );
}

pub use stock_lending_vault::StockLendingVault;

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    StockToken,
    "abi/StockToken.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    SupplyVault,
    "abi/SupplyVault.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Usdg,
    "abi/Usdg.json"
);
