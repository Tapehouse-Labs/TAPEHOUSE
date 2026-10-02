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
    GapCover,
    "abi/GapCover.json"
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
    QuoterV2,
    "abi/QuoterV2.json"
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
    StockToken,
    "abi/StockToken.json"
);

sol!(
    #[sol(rpc)]
    #[derive(Debug)]
    Usdg,
    "abi/Usdg.json"
);
