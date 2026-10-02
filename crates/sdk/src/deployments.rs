// SPDX-License-Identifier: MIT OR Apache-2.0
use std::collections::BTreeMap;
use std::str::FromStr;

use alloy::primitives::{Address, B256};
use serde_json::Value;

use crate::{Error, Result};

/// The chains whose Chainlink feeds price the share rather than the Stock Token.
pub const SHARE_PRICE_CHAINS: &[u64] = &[42161];

/// The margin accounts' cross position.
pub const CROSS: B256 = B256::ZERO;

/// A chain's address registry, `deployments/<chainId>.json`.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Deployments {
    /// `.chainId`.
    pub chain_id: u64,
    /// `.tokens`: USDG, WETH and the Stock Tokens.
    pub tokens: BTreeMap<String, Address>,
    /// `.chainlink`: the Chainlink feeds, which price the token or the share as the chain does.
    pub chainlink: BTreeMap<String, Address>,
    /// `.bandFeeds`: each asset's `BandFeed`.
    pub band_feeds: BTreeMap<String, Address>,
    /// `.tapehouse`: Tapehouse's own contracts, but for the stock lending vaults.
    pub tapehouse: BTreeMap<String, Address>,
    /// `.tapehouse.StockLending`: each Stock Token's lending vault.
    pub stock_lending: BTreeMap<String, Address>,
    /// `.uniswapV3`: the factory, the router, the quoter and the pools.
    pub uniswap_v3: BTreeMap<String, Address>,
}

/// A Chainlink feed that prices the Stock Token, as Robinhood Chain's do.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct TokenPriceFeed(pub Address);

/// A Chainlink feed that prices the share, as Arbitrum One's do.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct SharePriceFeed(pub Address);

/// An asset's symbol, or a RedStone feed ID, as the contracts take it: its bytes, right-padded to 32. A name longer
/// than 32 bytes is refused.
pub fn to_bytes32(name: &str) -> Result<B256> {
    if name.len() > B256::len_bytes() {
        return Err(Error::Argument(format!("{name:?} is longer than 32 bytes")));
    }
    Ok(B256::right_padding_from(name.as_bytes()))
}

impl Deployments {
    /// Reads a chain's registry from `path`.
    pub fn load(path: impl AsRef<std::path::Path>) -> Result<Self> {
        let text = std::fs::read_to_string(path.as_ref())
            .map_err(|e| Error::Registry(format!("{}: {e}", path.as_ref().display())))?;
        Self::parse(&text)
    }

    /// Parses a chain's registry. An address in mixed case must carry its checksum.
    pub fn parse(json: &str) -> Result<Self> {
        let registry: Value =
            serde_json::from_str(json).map_err(|e| Error::Registry(e.to_string()))?;
        let chain_id = registry["chainId"]
            .as_u64()
            .filter(|id| *id != 0)
            .ok_or_else(|| Error::Registry("the registry has no chainId".into()))?;
        let mut tapehouse = registry["tapehouse"].clone();
        let stock_lending = tapehouse
            .as_object_mut()
            .and_then(|group| group.remove("StockLending"))
            .unwrap_or(Value::Null);
        Ok(Self {
            chain_id,
            tokens: addresses(&registry["tokens"], ".tokens")?,
            chainlink: addresses(&registry["chainlink"], ".chainlink")?,
            band_feeds: addresses(&registry["bandFeeds"], ".bandFeeds")?,
            tapehouse: addresses(&tapehouse, ".tapehouse")?,
            stock_lending: addresses(&stock_lending, ".tapehouse.StockLending")?,
            uniswap_v3: addresses(&registry["uniswapV3"], ".uniswapV3")?,
        })
    }

    /// Whether the chain's Chainlink feeds price the share rather than the Stock Token.
    pub fn share_prices(&self) -> bool {
        SHARE_PRICE_CHAINS.contains(&self.chain_id)
    }

    /// The Chainlink feed `name` of a chain whose feeds price the Stock Token.
    pub fn token_price_feed(&self, name: &str) -> Result<TokenPriceFeed> {
        let feed = entry(&self.chainlink, name, ".chainlink")?;
        if self.share_prices() {
            return Err(Error::Registry(format!(
                ".chainlink.{name} prices the share on chain {}",
                self.chain_id
            )));
        }
        Ok(TokenPriceFeed(feed))
    }

    /// The Chainlink feed `name` of a chain whose feeds price the share.
    pub fn share_price_feed(&self, name: &str) -> Result<SharePriceFeed> {
        let feed = entry(&self.chainlink, name, ".chainlink")?;
        if !self.share_prices() {
            return Err(Error::Registry(format!(
                ".chainlink.{name} prices the Stock Token on chain {}",
                self.chain_id
            )));
        }
        Ok(SharePriceFeed(feed))
    }
}

pub(crate) fn entry(group: &BTreeMap<String, Address>, name: &str, path: &str) -> Result<Address> {
    group
        .get(name)
        .copied()
        .ok_or_else(|| Error::Registry(format!("the registry has no {path}.{name}")))
}

fn addresses(group: &Value, path: &str) -> Result<BTreeMap<String, Address>> {
    let Some(group) = group.as_object() else {
        return match group {
            Value::Null => Ok(BTreeMap::new()),
            _ => Err(Error::Registry(format!("{path} is not an object"))),
        };
    };
    group
        .iter()
        .map(|(name, value)| {
            let invalid = || Error::Registry(format!("{path}.{name} is not an address"));
            let text = value.as_str().ok_or_else(invalid)?;
            let address = if text == text.to_lowercase() {
                Address::from_str(text).map_err(|_| invalid())?
            } else {
                Address::parse_checksummed(text, None).map_err(|_| invalid())?
            };
            Ok((name.clone(), address))
        })
        .collect()
}
