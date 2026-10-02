// SPDX-License-Identifier: MIT OR Apache-2.0
use std::collections::BTreeMap;
use std::future::Future;

use alloy::contract::SolCallBuilder;
use alloy::eips::BlockId;
use alloy::primitives::{Address, B256, Bytes, U256, aliases::U24};
use alloy::providers::Provider;

use crate::bindings::{
    Aggregator, Band, BandFeed, Basket, IQuoterV2, MarginAccounts, MorphoBandOracle, QuoterV2,
    ShortPositions, StockToken, Usdg,
};
use crate::deployments::{Deployments, SharePriceFeed, TokenPriceFeed, entry, to_bytes32};
use crate::{Error, Result, Revert};

/// Where signed RedStone data packages come from: the integrator's own gateway client, cache or relay. It returns the
/// payload the band's `writePrices` verifies, the packages for `feed_ids` serialised as RedStone's EVM connector
/// appends them. The SDK holds no API key and calls no gateway.
pub trait PackageSource {
    /// The payload of the latest signed packages of `feed_ids`.
    fn payload(
        &self,
        feed_ids: &[&str],
    ) -> impl Future<Output = std::result::Result<Bytes, Box<dyn std::error::Error + Send + Sync>>> + Send;
}

/// What a sale pays through the asset's pool now, by QuoterV2's quote, and the least to accept.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Sale {
    /// What the pool pays now, in USDG.
    pub proceeds: U256,
    /// The quote less the slippage: the sale's `minProceeds`.
    pub min_proceeds: U256,
}

/// What a buy-back costs through the asset's pool now, by QuoterV2's quote, and the most to pay.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct BuyBack {
    /// What the pool asks now, in USDG.
    pub cost: U256,
    /// The quote plus the slippage: the cover's `maxCost`.
    pub max_cost: U256,
}

/// What repays all a position owes: its debt, its premium, and their sum, in USDG.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Repayment {
    /// The position's debt.
    pub debt: U256,
    /// The premium it owes.
    pub premium: U256,
    /// Both: the `assets` of a repayment that clears it now.
    pub assets: U256,
}

/// Why a Morpho oracle has no price.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum NoPrice {
    /// The band has no live leg, which a fresh `writePrices` ends.
    Stale,
    /// The band holds a signed halt, the issuer's pause or an unconfirmed multiplier step for the asset.
    Halted,
    /// The L2 sequencer is down or came back an hour ago or less.
    SequencerNotSettled,
}

/// A Morpho oracle's answer. Never a price of zero.
#[derive(Clone, Debug, PartialEq)]
pub enum OraclePrice {
    /// The price of one whole collateral token in the loan token, with Morpho's 36 + loan decimals − collateral
    /// decimals.
    Price(U256),
    /// No price, why, and the revert that said so.
    NoPrice {
        /// Why the oracle has no price.
        reason: NoPrice,
        /// The oracle's revert.
        revert: Revert,
    },
}

/// The assets a basket holds and their Stock Tokens, in the order of every list of amounts.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Components {
    /// The assets' names.
    pub assets: Vec<String>,
    /// Their Stock Tokens.
    pub tokens: Vec<Address>,
}

/// A basket's target proposed and not yet in effect: raw units of each Stock Token per share, and when it takes effect,
/// in seconds. Empty and zero if none.
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct PendingTarget {
    /// The units per share.
    pub units: Vec<U256>,
    /// When they take effect.
    pub effective_at: u64,
}

const BPS: u64 = 10_000;

fn asset_name(symbol: &B256) -> String {
    String::from_utf8_lossy(symbol.as_slice())
        .trim_end_matches('\0')
        .to_string()
}

fn check_slippage(slippage_bps: u64) -> Result<()> {
    if slippage_bps > BPS {
        return Err(Error::Argument(format!(
            "slippage_bps {slippage_bps} is outside 0 to 10000"
        )));
    }
    Ok(())
}

/// Tapehouse's contracts on one chain, at the addresses of its registry. Each accessor returns alloy's contract
/// instance, whose methods build typed calls and transactions.
#[derive(Clone, Debug)]
pub struct Tapehouse<P> {
    provider: P,
    deployments: Deployments,
}

impl<P: Provider + Clone> Tapehouse<P> {
    /// Tapehouse's contracts as `deployments` names them, over `provider`.
    pub fn new(provider: P, deployments: Deployments) -> Self {
        Self {
            provider,
            deployments,
        }
    }

    /// The registry the addresses come from.
    pub fn deployments(&self) -> &Deployments {
        &self.deployments
    }

    /// The band, `.tapehouse.Band`: its reads and `writePrices`.
    pub fn band(&self) -> Result<Band::BandInstance<P>> {
        let address = entry(&self.deployments.tapehouse, "Band", ".tapehouse")?;
        Ok(Band::new(address, self.provider.clone()))
    }

    /// `asset`'s `BandFeed`, from `.bandFeeds`.
    pub fn band_feed(&self, asset: &str) -> Result<BandFeed::BandFeedInstance<P>> {
        let address = entry(&self.deployments.band_feeds, asset, ".bandFeeds")?;
        Ok(BandFeed::new(address, self.provider.clone()))
    }

    /// The margin accounts, `.tapehouse.MarginAccounts`.
    pub fn accounts(&self) -> Result<MarginAccounts::MarginAccountsInstance<P>> {
        let address = entry(&self.deployments.tapehouse, "MarginAccounts", ".tapehouse")?;
        Ok(MarginAccounts::new(address, self.provider.clone()))
    }

    /// The short positions, `.tapehouse.ShortPositions`.
    pub fn shorts(&self) -> Result<ShortPositions::ShortPositionsInstance<P>> {
        let address = entry(&self.deployments.tapehouse, "ShortPositions", ".tapehouse")?;
        Ok(ShortPositions::new(address, self.provider.clone()))
    }

    /// The basket `key`, from `.tapehouse.Baskets`: `mint` takes, for each Stock Token in its order, at most
    /// `maxAssets`, its part rounded up as `previewMint` reads it; `redeem` gives each part rounded down.
    pub fn basket(&self, key: &str) -> Result<Basket::BasketInstance<P>> {
        let address = entry(&self.deployments.baskets, key, ".tapehouse.Baskets")?;
        Ok(Basket::new(address, self.provider.clone()))
    }

    /// The margin accounts' `unwrap` of `shares` of the basket `key` in `account`'s cross position, for its Stock
    /// Tokens, which the position then holds. The account or an address it authorized may, and anyone once the position
    /// falls short. A basket's shares themselves go in and out of the cross position through `deposit` and `withdraw`.
    pub fn unwrap(
        &self,
        account: Address,
        key: &str,
        shares: U256,
    ) -> Result<SolCallBuilder<P, MarginAccounts::unwrapCall>> {
        let basket = entry(&self.deployments.baskets, key, ".tapehouse.Baskets")?;
        Ok(self
            .accounts()?
            .unwrap(account, basket, shares)
            .with_cloned_provider())
    }

    /// The assets the basket `key` holds and their Stock Tokens.
    pub async fn basket_components(&self, key: &str) -> Result<Components> {
        let out = self.basket(key)?.components().call().await?;
        Ok(Components {
            assets: out.symbols.iter().map(asset_name).collect(),
            tokens: out.tokens,
        })
    }

    /// What minting `shares` of the basket `key` takes of each Stock Token, in raw units, rounded up.
    pub async fn preview_mint(&self, key: &str, shares: U256) -> Result<Vec<U256>> {
        Ok(self.basket(key)?.previewMint(shares).call().await?)
    }

    /// What redeeming `shares` of the basket `key` gives of each Stock Token, in raw units, rounded down.
    pub async fn preview_redeem(&self, key: &str, shares: U256) -> Result<Vec<U256>> {
        Ok(self.basket(key)?.previewRedeem(shares).call().await?)
    }

    /// The target in effect of the basket `key`: raw units of each Stock Token per share.
    pub async fn basket_target(&self, key: &str) -> Result<Vec<U256>> {
        Ok(self.basket(key)?.target().call().await?)
    }

    /// The target of the basket `key` proposed and not yet in effect.
    pub async fn pending_target(&self, key: &str) -> Result<PendingTarget> {
        let out = self.basket(key)?.pendingTarget().call().await?;
        Ok(PendingTarget {
            units: out.units,
            effective_at: out.effectiveAt,
        })
    }

    /// What `account`'s `position` holds of each Stock Token through its baskets, by asset, as its shares would
    /// redeem now: the engine margins them as those Stock Tokens. Both reads are at the latest block.
    pub async fn in_baskets(
        &self,
        account: Address,
        position: B256,
    ) -> Result<BTreeMap<String, U256>> {
        let accounts = self.accounts()?;
        let block = BlockId::number(
            self.provider
                .get_block_number()
                .await
                .map_err(alloy::contract::Error::from)?,
        );
        let stocks = accounts.stocks().block(block).call().await?;
        let amounts = accounts
            .inBaskets(account, position)
            .block(block)
            .call()
            .await?;
        if amounts.len() != stocks.symbols.len() {
            return Err(Error::Registry(format!(
                "the accounts name {} assets and {} amounts",
                stocks.symbols.len(),
                amounts.len()
            )));
        }
        Ok(stocks.symbols.iter().map(asset_name).zip(amounts).collect())
    }

    /// USDG, `.tokens.USDG`.
    pub fn usdg(&self) -> Result<Usdg::UsdgInstance<P>> {
        let address = entry(&self.deployments.tokens, "USDG", ".tokens")?;
        Ok(Usdg::new(address, self.provider.clone()))
    }

    /// `asset`'s Stock Token, from `.tokens`.
    pub fn stock_token(&self, asset: &str) -> Result<StockToken::StockTokenInstance<P>> {
        let address = entry(&self.deployments.tokens, asset, ".tokens")?;
        Ok(StockToken::new(address, self.provider.clone()))
    }

    /// `asset`'s Morpho oracle, from `.morphoOracles`: its `price()`, `halt()`, `band()`, `symbol()`,
    /// `collateralToken()`, `loanToken()`, `scaleFactor()` and `owner()`.
    pub fn morpho_oracle(
        &self,
        asset: &str,
    ) -> Result<MorphoBandOracle::MorphoBandOracleInstance<P>> {
        let address = entry(&self.deployments.morpho_oracles, asset, ".morphoOracles")?;
        Ok(MorphoBandOracle::new(address, self.provider.clone()))
    }

    /// The price `asset`'s Morpho oracle answers Morpho Blue, read at the latest block. A revert with `NoAnswer` or
    /// `SequencerNotSettled` is no price, never zero; the oracle's halt and the band's corporate action, read at the
    /// same block, tell a stale band from a halt.
    pub async fn oracle_price(&self, asset: &str) -> Result<OraclePrice> {
        let oracle = self.morpho_oracle(asset)?;
        let block = BlockId::number(
            self.provider
                .get_block_number()
                .await
                .map_err(alloy::contract::Error::from)?,
        );
        let revert = match oracle
            .price()
            .block(block)
            .call()
            .await
            .map_err(Error::from)
        {
            Ok(price) => return Ok(OraclePrice::Price(price)),
            Err(Error::Revert(revert)) if revert.name == "SequencerNotSettled" => {
                return Ok(OraclePrice::NoPrice {
                    reason: NoPrice::SequencerNotSettled,
                    revert,
                });
            }
            Err(Error::Revert(revert)) if revert.name == "NoAnswer" => revert,
            Err(error) => return Err(error),
        };
        let halt = oracle.halt().block(block).call().await?;
        let band = oracle.band().block(block).call().await?;
        let symbol = oracle.symbol().block(block).call().await?;
        let step = Band::new(band, self.provider.clone())
            .corporateAction(symbol)
            .block(block)
            .call()
            .await?;
        let halted = halt.signedHalt || halt.oraclePaused || step.status == 2;
        Ok(OraclePrice::NoPrice {
            reason: if halted {
                NoPrice::Halted
            } else {
                NoPrice::Stale
            },
            revert,
        })
    }

    /// A Chainlink feed that prices the Stock Token.
    pub fn token_price(&self, feed: TokenPriceFeed) -> Aggregator::AggregatorInstance<P> {
        Aggregator::new(feed.0, self.provider.clone())
    }

    /// A Chainlink feed that prices the share.
    pub fn share_price(&self, feed: SharePriceFeed) -> Aggregator::AggregatorInstance<P> {
        Aggregator::new(feed.0, self.provider.clone())
    }

    /// The band's `writePrices` of `feed_ids`, with the payload `source` signs for them.
    pub async fn write_prices(
        &self,
        source: &impl PackageSource,
        feed_ids: &[&str],
    ) -> Result<SolCallBuilder<P, Band::writePricesCall>> {
        let ids = feed_ids
            .iter()
            .map(|id| to_bytes32(id))
            .collect::<Result<_>>()?;
        let payload = source.payload(feed_ids).await.map_err(Error::Source)?;
        Ok(self
            .band()?
            .writePrices(ids, payload)
            .with_cloned_provider())
    }

    /// What repays all `account`'s `position` owes now: its debt and its premium, both read at the latest block.
    pub async fn repayment(&self, account: Address, position: B256) -> Result<Repayment> {
        let accounts = self.accounts()?;
        let block = BlockId::number(
            self.provider
                .get_block_number()
                .await
                .map_err(alloy::contract::Error::from)?,
        );
        let debt = accounts.debt(account, position).block(block).call().await?;
        let premium = accounts
            .premium(account, position)
            .block(block)
            .call()
            .await?;
        Ok(Repayment {
            debt,
            premium,
            assets: debt + premium,
        })
    }

    /// What `account`'s `position` holds of `token`: USDG, WETH or a Stock Token.
    pub async fn collateral(
        &self,
        account: Address,
        position: B256,
        token: Address,
    ) -> Result<U256> {
        Ok(self
            .accounts()?
            .collateral(account, position, token)
            .call()
            .await?)
    }

    /// The gross exposure of `account`'s `position` over its equity, in basis points; `U256::MAX` for a position
    /// without equity.
    pub async fn leverage(&self, account: Address, position: B256) -> Result<U256> {
        Ok(self.accounts()?.leverage(account, position).call().await?)
    }

    /// The price of `asset`, in USD with 8 decimals, at or above which `account`'s `position` meets its requirement
    /// once it owes `borrowing` more USDG: its band's low edge where it falls short already, zero where no price
    /// leaves it short.
    pub async fn liquidation_price(
        &self,
        account: Address,
        position: B256,
        asset: &str,
        borrowing: U256,
    ) -> Result<U256> {
        let symbol = to_bytes32(asset)?;
        Ok(self
            .accounts()?
            .liquidationPrice(account, position, symbol, borrowing)
            .call()
            .await?)
    }

    /// What selling `amount` of `asset` pays through its pool now, from Uniswap's QuoterV2, and that less
    /// `slippage_bps`, at most 10000.
    pub async fn quote_sale(&self, asset: &str, amount: U256, slippage_bps: u64) -> Result<Sale> {
        check_slippage(slippage_bps)?;
        let (quoter, token, usdg, fee) = self.pool(asset).await?;
        let params = IQuoterV2::QuoteExactInputSingleParams {
            tokenIn: token,
            tokenOut: usdg,
            amountIn: amount,
            fee,
            sqrtPriceLimitX96: Default::default(),
        };
        let proceeds = quoter.quoteExactInputSingle(params).call().await?.amountOut;
        Ok(Sale {
            proceeds,
            min_proceeds: proceeds * U256::from(BPS - slippage_bps) / U256::from(BPS),
        })
    }

    /// What buying back `amount` of `asset` costs through its pool now, from Uniswap's QuoterV2, and that plus
    /// `slippage_bps`, at most 10000, rounded up.
    pub async fn quote_cover(
        &self,
        asset: &str,
        amount: U256,
        slippage_bps: u64,
    ) -> Result<BuyBack> {
        check_slippage(slippage_bps)?;
        let (quoter, token, usdg, fee) = self.pool(asset).await?;
        let params = IQuoterV2::QuoteExactOutputSingleParams {
            tokenIn: usdg,
            tokenOut: token,
            amount,
            fee,
            sqrtPriceLimitX96: Default::default(),
        };
        let cost = quoter.quoteExactOutputSingle(params).call().await?.amountIn;
        Ok(BuyBack {
            cost,
            max_cost: (cost * U256::from(BPS + slippage_bps)).div_ceil(U256::from(BPS)),
        })
    }

    async fn pool(
        &self,
        asset: &str,
    ) -> Result<(QuoterV2::QuoterV2Instance<P>, Address, Address, U24)> {
        let quoter = entry(&self.deployments.uniswap_v3, "QuoterV2", ".uniswapV3")?;
        let token = entry(&self.deployments.tokens, asset, ".tokens")?;
        let usdg = entry(&self.deployments.tokens, "USDG", ".tokens")?;
        let fee = self.shorts()?.fee(to_bytes32(asset)?).call().await?;
        Ok((
            QuoterV2::new(quoter, self.provider.clone()),
            token,
            usdg,
            fee,
        ))
    }
}
