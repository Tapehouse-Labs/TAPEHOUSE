#![cfg_attr(not(any(test, feature = "export-abi")), no_main)]
#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

pub mod error;
pub mod matrix;
pub mod requirement;
pub mod scenario;
pub mod session;
pub mod uniswap;

use alloc::vec::Vec;

use ownable::{IOwnable2Step, Ownable2Step};
use stylus_sdk::alloy_primitives::{Address, B256, I256, U8, U16, U32, U64, U256};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::{
    StorageAddress, StorageB256, StorageBool, StorageMap, StorageU8, StorageU16, StorageU32,
    StorageU64, StorageVec,
};

use crate::error::{
    AssetCount, CorrelationStepTooLarge, DepthStepTooLarge, DuplicateAsset, ExposureTooLarge,
    GapStepTooLarge, InsufficientGas, InvalidCorrelation, InvalidDepth, InvalidFeed, InvalidGap,
    InvalidPool, InvalidVolatility, LengthMismatch, MarginError, NotPositiveDefinite,
    ScenarioOutOfRange, UnknownAsset, UnsupportedScenarioSize, UpdateTooSoon,
    VolatilityStepTooLarge, ZeroSymbol,
};
use crate::matrix::{ONE, pair, positive_definite};
use crate::requirement::{Pool, PoolTerms, exposures, requirements};
use crate::scenario::{Parameters, SIZES, Set, count};
use crate::session::{HORIZON, WEEKEND_LEVERAGE, current, regime, starved};

/// Most assets the engine takes: the exact positive-definiteness check stays within 256 bits up to 9.
pub const MAX_ASSETS: usize = 8;
/// First bit of the weekend gaps in `set_parameters`'s mask of changed values, after the volatilities' and the
/// correlations'.
const GAP_BITS: usize = MAX_ASSETS + MAX_ASSETS * (MAX_ASSETS - 1) / 2;
/// First bit of the depths in the same mask, one per asset for either side, after the weekend gaps'.
const DEPTH_BITS: usize = GAP_BITS + MAX_ASSETS;
/// Highest daily volatility, in centi-basis-points: 100% a day.
pub const MAX_VOLATILITY: u32 = 1_000_000;
/// Largest weekend gap, in centi-basis-points: 100%.
pub const MAX_GAP: u32 = 1_000_000;
/// Shortest time between two updates, in seconds.
pub const UPDATE_INTERVAL: u64 = 86_400;
/// Lattice size of the scenario set the engine margins with.
pub const SCENARIO_SIZE: usize = 256;

sol! {
    event VolatilitySet(bytes32 indexed symbol, uint32 value);
    event CorrelationSet(bytes32 indexed symbol, bytes32 indexed other, uint16 value);
    event GapSet(bytes32 indexed symbol, uint32 value);
    event DepthSet(bytes32 indexed symbol, uint32 selling, uint32 buying);

}

/// An asset as the constructor takes it: its symbol, volatility floor, volatility, weekend-gap floor, weekend
/// gap, selling and buying depths, which are also their ceilings, its Uniswap v3 pool, and its Stock Token.
pub type Asset = (B256, u32, u32, u32, u32, u32, u32, Address, Address);

/// Where an asset is liquidated: its Uniswap v3 pool, which of the pool's tokens is the asset, whether it
/// is quoted in WETH rather than USDG, both tokens' decimals and the pool's fee in millionths.
#[storage]
pub struct Source {
    pool: StorageAddress,
    stock_is_token0: StorageBool,
    quote_is_weth: StorageBool,
    stock_decimals: StorageU8,
    quote_decimals: StorageU8,
    fee: StorageU32,
}

#[storage]
#[entrypoint]
pub struct Margin {
    ownable: Ownable2Step,
    assets: StorageVec<StorageB256>,
    positions: StorageMap<B256, StorageU8>,
    volatilities: StorageVec<StorageU32>,
    volatility_floors: StorageVec<StorageU32>,
    correlations: StorageVec<StorageU16>,
    correlation_floors: StorageVec<StorageU16>,
    gaps: StorageVec<StorageU32>,
    gap_floors: StorageVec<StorageU32>,
    depths: StorageVec<StorageU32>,
    depth_ceilings: StorageVec<StorageU32>,
    market: StorageU8,
    sources: StorageMap<B256, Source>,
    eth_usd: StorageAddress,
    band: StorageAddress,
    last_update: StorageU64,
}

#[public]
#[implements(IOwnable2Step<Error = MarginError>)]
impl Margin {
    /// Sets the assets, their risk parameters and the hard floors under them, which never change.
    /// Volatilities are daily, and each weekend gap is the move from the last close before a closure to the next
    /// open, both in centi-basis-points; correlations are in basis points, listed for each pair `(i, j)`, `i < j`,
    /// row by row. Every floor is positive, every value sits at or above its floor, and the correlation matrix is
    /// positive definite.
    /// Depths are the USD a liquidation can sell, then buy, within a 10% move of each asset's pool, and their
    /// first values are ceilings that never change. `market` is the asset that stands for the market, or zero
    /// for the equal-weighted portfolio of the assets. Each asset's pool trades its Stock Token against `usdg`
    /// or `weth`, or is zero for none; `eth_usd` prices WETH. `band`'s session sets the current
    /// requirement's regime. `initial_owner` may update the parameters.
    #[constructor]
    #[allow(clippy::too_many_arguments)]
    pub fn constructor(
        &mut self,
        assets: Vec<Asset>,
        correlation_floors: Vec<u16>,
        correlations: Vec<u16>,
        market: B256,
        usdg: Address,
        weth: Address,
        eth_usd: Address,
        band: Address,
        initial_owner: Address,
    ) -> Result<(), MarginError> {
        self.ownable.initialize::<MarginError>(initial_owner)?;
        let symbols: Vec<B256> = assets.iter().map(|&(symbol, ..)| symbol).collect();
        let volatilities: Vec<u32> = assets
            .iter()
            .map(|&(_, _, volatility, ..)| volatility)
            .collect();
        let gaps: Vec<u32> = assets.iter().map(|&(_, _, _, _, gap, ..)| gap).collect();
        let depths: Vec<u32> = assets
            .iter()
            .flat_map(|&(.., selling, buying, _, _)| [selling, buying])
            .collect();
        let n = symbols.len();
        if n == 0 || n > MAX_ASSETS {
            return Err(MarginError::AssetCount(AssetCount {
                count: stylus_sdk::alloy_primitives::U256::from(n),
            }));
        }
        let pairs = n * (n - 1) / 2;
        if correlation_floors.len() != pairs || correlations.len() != pairs {
            return Err(MarginError::LengthMismatch(LengthMismatch {}));
        }
        for (i, &symbol) in symbols.iter().enumerate() {
            if symbol == B256::ZERO {
                return Err(MarginError::ZeroSymbol(ZeroSymbol {}));
            }
            if self.positions.get(symbol) != U8::ZERO {
                return Err(MarginError::DuplicateAsset(DuplicateAsset { symbol }));
            }
            self.positions.setter(symbol).set(U8::from(i + 1));
            self.assets.push(symbol);
        }
        for &(symbol, floor, volatility, ..) in &assets {
            if floor == 0 {
                return Err(invalid_volatility(symbol, volatility, floor));
            }
            self.volatility_floors.push(U32::from(floor));
        }
        for i in 0..n {
            for j in i + 1..n {
                let floor = correlation_floors[pair(n, i, j)];
                if floor == 0 {
                    return Err(invalid_correlation(
                        symbols[i],
                        symbols[j],
                        correlations[pair(n, i, j)],
                        floor,
                    ));
                }
                self.correlation_floors.push(U16::from(floor));
            }
        }
        for &(symbol, _, _, floor, gap, ..) in &assets {
            if floor == 0 {
                return Err(invalid_gap(symbol, gap, floor));
            }
            self.gap_floors.push(U32::from(floor));
        }
        for &depth in &depths {
            self.depth_ceilings.push(U32::from(depth));
        }
        if market != B256::ZERO {
            let i = self.position(market)?;
            self.market.set(U8::from(i + 1));
        }
        self.eth_usd.set(eth_usd);
        self.band.set(band);
        for &(symbol, .., pool, token) in &assets {
            self.set_source(symbol, pool, token, usdg, weth, eth_usd)?;
        }
        self.check(&symbols, &volatilities, &correlations, &gaps, &depths)?;
        for (i, &value) in volatilities.iter().enumerate() {
            self.volatilities.push(U32::from(value));
            self.vm().log(VolatilitySet {
                symbol: symbols[i],
                value,
            });
        }
        for i in 0..n {
            for j in i + 1..n {
                let value = correlations[pair(n, i, j)];
                self.correlations.push(U16::from(value));
                self.vm().log(CorrelationSet {
                    symbol: symbols[i],
                    other: symbols[j],
                    value,
                });
            }
        }
        for (i, &value) in gaps.iter().enumerate() {
            self.gaps.push(U32::from(value));
            self.vm().log(GapSet {
                symbol: symbols[i],
                value,
            });
        }
        for (i, pair) in depths.chunks(2).enumerate() {
            self.depths.push(U32::from(pair[0]));
            self.depths.push(U32::from(pair[1]));
            self.vm().log(DepthSet {
                symbol: symbols[i],
                selling: pair[0],
                buying: pair[1],
            });
        }
        Ok(())
    }

    /// Replaces every parameter, in the constructor's order. Each value stays at or above its floor, or for a
    /// depth above zero and at or below its ceiling. Each moves by at most ×1.5 or ÷1.5, except that a depth
    /// may fall by any amount. The correlation matrix stays positive definite, and an update comes at least a
    /// day after the one before. Emits an event for every value that changes. Owner only.
    pub fn set_parameters(
        &mut self,
        volatilities: Vec<u32>,
        correlations: Vec<u16>,
        gaps: Vec<u32>,
        depths: Vec<u32>,
    ) -> Result<(), MarginError> {
        self.ownable.only_owner::<MarginError>()?;
        let now = self.vm().block_timestamp();
        let last = self.last_update.get().to::<u64>();
        if now < last + UPDATE_INTERVAL {
            return Err(MarginError::UpdateTooSoon(UpdateTooSoon {
                nextUpdateAt: last + UPDATE_INTERVAL,
            }));
        }
        let symbols = self.assets();
        let n = symbols.len();
        if volatilities.len() != n
            || correlations.len() != n * (n - 1) / 2
            || gaps.len() != n
            || depths.len() != 2 * n
        {
            return Err(MarginError::LengthMismatch(LengthMismatch {}));
        }
        self.check(&symbols, &volatilities, &correlations, &gaps, &depths)?;
        let mut changed = 0u64;
        for (i, &value) in volatilities.iter().enumerate() {
            let previous = self.volatilities.get(i).unwrap().to::<u32>();
            changed |= u64::from(previous != value) << i;
            if !within_step(previous.into(), value.into()) {
                return Err(MarginError::VolatilityStepTooLarge(
                    VolatilityStepTooLarge {
                        symbol: symbols[i],
                        previous,
                        value,
                    },
                ));
            }
        }
        for i in 0..n {
            for j in i + 1..n {
                let k = pair(n, i, j);
                let previous = self.correlations.get(k).unwrap().to::<u16>();
                changed |= u64::from(previous != correlations[k]) << (MAX_ASSETS + k);
                if !within_step(previous.into(), correlations[k].into()) {
                    return Err(MarginError::CorrelationStepTooLarge(
                        CorrelationStepTooLarge {
                            symbol: symbols[i],
                            other: symbols[j],
                            previous,
                            value: correlations[k],
                        },
                    ));
                }
            }
        }
        for (i, &value) in gaps.iter().enumerate() {
            let previous = self.gaps.get(i).unwrap().to::<u32>();
            changed |= u64::from(previous != value) << (GAP_BITS + i);
            if !within_step(previous.into(), value.into()) {
                return Err(MarginError::GapStepTooLarge(GapStepTooLarge {
                    symbol: symbols[i],
                    previous,
                    value,
                }));
            }
        }
        for (k, &value) in depths.iter().enumerate() {
            let previous = self.depths.get(k).unwrap().to::<u32>();
            changed |= u64::from(previous != value) << (DEPTH_BITS + k / 2);
            if 2 * u64::from(value) > 3 * u64::from(previous) {
                return Err(MarginError::DepthStepTooLarge(DepthStepTooLarge {
                    symbol: symbols[k / 2],
                    previous,
                    value,
                }));
            }
        }
        for (i, &value) in volatilities.iter().enumerate() {
            if changed >> i & 1 == 1 {
                self.volatilities.setter(i).unwrap().set(U32::from(value));
                self.vm().log(VolatilitySet {
                    symbol: symbols[i],
                    value,
                });
            }
        }
        for i in 0..n {
            for j in i + 1..n {
                let k = pair(n, i, j);
                let value = correlations[k];
                if changed >> (MAX_ASSETS + k) & 1 == 1 {
                    self.correlations.setter(k).unwrap().set(U16::from(value));
                    self.vm().log(CorrelationSet {
                        symbol: symbols[i],
                        other: symbols[j],
                        value,
                    });
                }
            }
        }
        for (i, &value) in gaps.iter().enumerate() {
            if changed >> (GAP_BITS + i) & 1 == 1 {
                self.gaps.setter(i).unwrap().set(U32::from(value));
                self.vm().log(GapSet {
                    symbol: symbols[i],
                    value,
                });
            }
        }
        for (i, pair) in depths.chunks(2).enumerate() {
            let (selling, buying) = (pair[0], pair[1]);
            if changed >> (DEPTH_BITS + i) & 1 == 1 {
                self.depths.setter(2 * i).unwrap().set(U32::from(selling));
                self.depths
                    .setter(2 * i + 1)
                    .unwrap()
                    .set(U32::from(buying));
                self.vm().log(DepthSet {
                    symbol: symbols[i],
                    selling,
                    buying,
                });
            }
        }
        self.last_update.set(U64::from(now));
        Ok(())
    }

    /// The assets, in the order of every parameter list.
    pub fn assets(&self) -> Vec<B256> {
        (0..self.assets.len())
            .map(|i| self.assets.get(i).unwrap())
            .collect()
    }

    /// The daily volatility of `symbol` and its floor, in centi-basis-points.
    pub fn volatility(&self, symbol: B256) -> Result<(u32, u32), MarginError> {
        let i = self.position(symbol)?;
        Ok((
            self.volatilities.get(i).unwrap().to::<u32>(),
            self.volatility_floors.get(i).unwrap().to::<u32>(),
        ))
    }

    /// The correlation of `symbol` and `other` and its floor, in basis points; 10 000 for an asset with
    /// itself.
    pub fn correlation(&self, symbol: B256, other: B256) -> Result<(u16, u16), MarginError> {
        let (i, j) = (self.position(symbol)?, self.position(other)?);
        if i == j {
            return Ok((ONE, ONE));
        }
        let k = pair(self.assets.len(), i.min(j), i.max(j));
        Ok((
            self.correlations.get(k).unwrap().to::<u16>(),
            self.correlation_floors.get(k).unwrap().to::<u16>(),
        ))
    }

    /// The weekend gap of `symbol` and its floor, in centi-basis-points.
    pub fn weekend_gap(&self, symbol: B256) -> Result<(u32, u32), MarginError> {
        let i = self.position(symbol)?;
        Ok((
            self.gaps.get(i).unwrap().to::<u32>(),
            self.gap_floors.get(i).unwrap().to::<u32>(),
        ))
    }

    /// The USD a liquidation of `symbol` can sell, then buy, within a 10% move of its pool, and the ceilings
    /// over both.
    pub fn depth(&self, symbol: B256) -> Result<(u32, u32, u32, u32), MarginError> {
        let i = self.position(symbol)?;
        let at = |v: &StorageVec<StorageU32>, k: usize| v.get(k).unwrap().to::<u32>();
        Ok((
            at(&self.depths, 2 * i),
            at(&self.depths, 2 * i + 1),
            at(&self.depth_ceilings, 2 * i),
            at(&self.depth_ceilings, 2 * i + 1),
        ))
    }

    /// The Uniswap v3 pool `symbol` is liquidated in; zero for none.
    pub fn pool(&self, symbol: B256) -> Result<Address, MarginError> {
        self.position(symbol)?;
        Ok(self.sources.get(symbol).pool.get())
    }

    /// The Chainlink feed that prices WETH; zero for none.
    pub fn eth_usd_feed(&self) -> Address {
        self.eth_usd.get()
    }

    /// The band whose session sets the current requirement's regime.
    pub fn band(&self) -> Address {
        self.band.get()
    }

    /// The asset that stands for the market; zero for the equal-weighted portfolio of the assets.
    pub fn market(&self) -> B256 {
        match self.market.get().to::<usize>() {
            0 => B256::ZERO,
            p => self.assets.get(p - 1).unwrap(),
        }
    }

    /// When the parameters were last updated; zero before the first update, which may come at any time.
    pub fn last_update(&self) -> u64 {
        self.last_update.get().to::<u64>()
    }

    /// Scenario `index` of the engine's set over a horizon of `horizon` seconds: returns in millionths, in
    /// the order of `assets()`. The set holds 256 joint draws with the stored correlations, 256 with every
    /// correlation 1, 256 with every correlation 0, the market down and up, every asset's weekend gap down
    /// and up, then each asset's gap alone, down and up, in the ascending order of the symbols.
    pub fn scenario(&self, index: u16, horizon: u64) -> Result<Vec<i32>, MarginError> {
        if usize::from(index) >= count(SCENARIO_SIZE, self.assets.len()) {
            return Err(MarginError::ScenarioOutOfRange(ScenarioOutOfRange {
                index,
            }));
        }
        let (symbols, volatilities, correlations, gaps, market) = self.parameters();
        let parameters = Parameters {
            symbols: &symbols,
            volatilities: &volatilities,
            correlations: &correlations,
            gaps: &gaps,
            market,
        };
        Ok(Set::new(&parameters, SCENARIO_SIZE, horizon).row(index.into()))
    }

    /// Keccak-256 of every scenario of a set with `size` lattice points (32, 64, 128 or 256) over a
    /// horizon of `horizon` seconds, each as int32 returns in the ascending order of the symbols. It does
    /// not depend on the order the assets were given in.
    pub fn scenario_digest(&self, size: u16, horizon: u64) -> Result<B256, MarginError> {
        if !SIZES.contains(&usize::from(size)) {
            return Err(MarginError::UnsupportedScenarioSize(
                UnsupportedScenarioSize { size },
            ));
        }
        let (symbols, volatilities, correlations, gaps, market) = self.parameters();
        let parameters = Parameters {
            symbols: &symbols,
            volatilities: &volatilities,
            correlations: &correlations,
            gaps: &gaps,
            market,
        };
        let encoded = Set::new(&parameters, size.into(), horizon).encoded();
        Ok(self.vm().native_keccak256(&encoded))
    }

    /// The margin a portfolio needs over a horizon of `horizon` seconds, in USD with 18 decimals, and a bit
    /// for every asset, in the order of `assets()`, whose liquidity input is missing. `quantities` are signed
    /// token amounts with 18 decimals, negative for a short, and `prices` USD with 8 decimals, both in the
    /// order of `assets()`. `spans_closure` counts the weekend-gap scenarios. The requirement is the larger of
    /// the expected shortfall at 99% over the scenario set, with diversification capped at 80%, the stress
    /// scenarios and the reference floor, plus what liquidating each position in its pool costs. An asset
    /// without a pool adds nothing and sets its bit. One whose pool lacks 30 minutes of history, or needs an
    /// ETH/USD answer that is stale, is charged its pool's fee and governed depth, and sets its bit.
    pub fn requirement(
        &self,
        quantities: Vec<I256>,
        prices: Vec<U256>,
        horizon: u64,
        spans_closure: bool,
    ) -> Result<(U256, u8), MarginError> {
        let (open, closed, _, missing) = self.requirements(&quantities, &prices, horizon)?;
        Ok((if spans_closure { closed } else { open }, missing))
    }

    /// The margin a portfolio needs now, in USD with 18 decimals, the bits of `requirement`, and the regime
    /// the band's session puts it in: 0 unknown, 1 closed, 2 open, 3 closing. The horizon is two days in
    /// every regime. The open market adds a 25% buffer. A closure, or a session the band cannot tell,
    /// needs the largest of that, the requirement across a closure, which counts the weekend-gap scenarios,
    /// and the gross value over `weekendLeverage()` plus the liquidity add-on; the seven hours before a
    /// weekend or holiday session close rise to it in a straight line. A position margined alone, as an
    /// isolated one is, is the portfolio that holds only it.
    /// A call that leaves the band's `session()` too little gas reverts with `InsufficientGas` rather than
    /// read as an unknown session.
    pub fn current_requirement(
        &self,
        quantities: Vec<I256>,
        prices: Vec<U256>,
    ) -> Result<(U256, u8, u8), MarginError> {
        let (open, closed, floor, missing) = self.requirements(&quantities, &prices, HORIZON)?;
        let before = self.vm().evm_gas_left();
        let session = session::read(self.vm(), self.band.get());
        if session.is_none() && starved(before, self.vm().evm_gas_left()) {
            return Err(MarginError::InsufficientGas(InsufficientGas {}));
        }
        let regime = regime(session, self.vm().block_timestamp().saturating_mul(1_000));
        Ok((current(open, closed, floor, regime), missing, regime.code()))
    }

    /// The most gross exposure the current requirement allows per unit of margin across a closure, in basis
    /// points.
    pub fn weekend_leverage(&self) -> u32 {
        WEEKEND_LEVERAGE
    }
}

#[public]
impl IOwnable2Step for Margin {
    type Error = MarginError;

    fn owner(&self) -> Address {
        self.ownable.owner()
    }

    fn pending_owner(&self) -> Address {
        self.ownable.pending_owner()
    }

    fn transfer_ownership(&mut self, new_owner: Address) -> Result<(), MarginError> {
        self.ownable.transfer_ownership(new_owner)
    }

    fn accept_ownership(&mut self) -> Result<(), MarginError> {
        self.ownable.accept_ownership()
    }

    fn renounce_ownership(&mut self) -> Result<(), MarginError> {
        self.ownable.renounce_ownership()
    }
}

impl Margin {
    /// The requirement of a portfolio over the open market and across a closure, its leverage floor with the
    /// liquidity add-on, and its missing bits.
    #[inline(never)]
    fn requirements(
        &self,
        quantities: &[I256],
        prices: &[U256],
        horizon: u64,
    ) -> Result<(U256, U256, U256, u8), MarginError> {
        let (symbols, volatilities, correlations, gaps, market) = self.parameters();
        let n = symbols.len();
        if quantities.len() != n || prices.len() != n {
            return Err(MarginError::LengthMismatch(LengthMismatch {}));
        }
        let exposures = exposures(quantities, prices).map_err(|i| {
            MarginError::ExposureTooLarge(ExposureTooLarge {
                symbol: B256::from(symbols[i]),
            })
        })?;
        if exposures.iter().all(|&e| e == 0) {
            return Ok((U256::ZERO, U256::ZERO, U256::ZERO, 0));
        }
        let parameters = Parameters {
            symbols: &symbols,
            volatilities: &volatilities,
            correlations: &correlations,
            gaps: &gaps,
            market,
        };
        let set = Set::new(&parameters, SCENARIO_SIZE, horizon);
        let now = self.vm().block_timestamp();
        let mut eth_usd = None;
        let pools: Vec<Pool> = exposures
            .iter()
            .zip(&symbols)
            .map(|(&e, &symbol)| {
                let source = self.sources.get(B256::from(symbol));
                if e == 0 || source.pool.get() == Address::ZERO {
                    return Pool::Absent;
                }
                match self.pool_terms(&source, &mut eth_usd, now) {
                    Some(terms) => Pool::Read(terms),
                    None => Pool::Unread {
                        fee: source.fee.get().to(),
                    },
                }
            })
            .collect();
        let depths: Vec<u32> = (0..2 * n)
            .map(|k| self.depths.get(k).unwrap().to())
            .collect();
        Ok(requirements(&set, &exposures, prices, &depths, &pools))
    }

    #[inline(never)]
    fn set_source(
        &mut self,
        symbol: B256,
        pool: Address,
        stock: Address,
        usdg: Address,
        weth: Address,
        eth_usd: Address,
    ) -> Result<(), MarginError> {
        if pool == Address::ZERO {
            return Ok(());
        }
        let invalid = || MarginError::InvalidPool(InvalidPool { symbol, pool });
        let (token0, token1, fee) = uniswap::tokens(self.vm(), pool).ok_or_else(invalid)?;
        if uniswap::cardinality(self.vm(), pool).is_none_or(|c| u32::from(c) <= uniswap::WINDOW) {
            return Err(invalid());
        }
        let (stock_is_token0, quote_is_weth) =
            uniswap::classify(token0, token1, stock, usdg, weth).ok_or_else(invalid)?;
        let (stock, quoted_in) = if stock_is_token0 {
            (token0, token1)
        } else {
            (token1, token0)
        };
        let stock_decimals = uniswap::decimals(self.vm(), stock)
            .filter(|&d| d <= 36)
            .ok_or_else(invalid)?;
        let quote_decimals = uniswap::decimals(self.vm(), quoted_in)
            .filter(|&d| d <= 18)
            .ok_or_else(invalid)?;
        if quote_is_weth
            && (eth_usd == Address::ZERO || !uniswap::has_feed_decimals(self.vm(), eth_usd))
        {
            return Err(MarginError::InvalidFeed(InvalidFeed { feed: eth_usd }));
        }
        let mut source = self.sources.setter(symbol);
        source.pool.set(pool);
        source.stock_is_token0.set(stock_is_token0);
        source.quote_is_weth.set(quote_is_weth);
        source.stock_decimals.set(U8::from(stock_decimals));
        source.quote_decimals.set(U8::from(quote_decimals));
        source.fee.set(U32::from(fee));
        Ok(())
    }

    /// What the asset's pool says about liquidating it, from its time-weighted tick and liquidity. `None`
    /// without 30 minutes of pool history, or without a live ETH/USD answer for a WETH pool.
    #[inline(never)]
    fn pool_terms(
        &self,
        source: &Source,
        eth_usd: &mut Option<Option<U256>>,
        now: u64,
    ) -> Option<PoolTerms> {
        let (tick, liquidity) = uniswap::consult(self.vm(), source.pool.get())?;
        let usd = if source.quote_is_weth.get() {
            (*eth_usd.get_or_insert_with(|| uniswap::answer(self.vm(), self.eth_usd.get(), now)))?
        } else {
            U256::from(100_000_000)
        };
        let (price, selling, buying) = uniswap::terms(
            source.stock_is_token0.get(),
            source.stock_decimals.get().to(),
            source.quote_decimals.get().to(),
            tick,
            liquidity,
            usd,
        )?;
        Some(PoolTerms {
            price,
            selling,
            buying,
            fee: source.fee.get().to(),
        })
    }

    fn position(&self, symbol: B256) -> Result<usize, MarginError> {
        match self.positions.get(symbol).to::<usize>() {
            0 => Err(MarginError::UnknownAsset(UnknownAsset { symbol })),
            p => Ok(p - 1),
        }
    }

    #[allow(clippy::type_complexity)]
    fn parameters(&self) -> (Vec<[u8; 32]>, Vec<u32>, Vec<u16>, Vec<u32>, Option<usize>) {
        let n = self.assets.len();
        let pairs = n * (n - 1) / 2;
        (
            (0..n).map(|i| self.assets.get(i).unwrap().0).collect(),
            (0..n)
                .map(|i| self.volatilities.get(i).unwrap().to())
                .collect(),
            (0..pairs)
                .map(|k| self.correlations.get(k).unwrap().to())
                .collect(),
            (0..n).map(|i| self.gaps.get(i).unwrap().to()).collect(),
            self.market.get().to::<usize>().checked_sub(1),
        )
    }

    fn check(
        &self,
        symbols: &[B256],
        volatilities: &[u32],
        correlations: &[u16],
        gaps: &[u32],
        depths: &[u32],
    ) -> Result<(), MarginError> {
        let n = symbols.len();
        for (i, &value) in volatilities.iter().enumerate() {
            let floor = self.volatility_floors.get(i).unwrap().to::<u32>();
            if value < floor || value > MAX_VOLATILITY {
                return Err(invalid_volatility(symbols[i], value, floor));
            }
        }
        for i in 0..n {
            for j in i + 1..n {
                let k = pair(n, i, j);
                let floor = self.correlation_floors.get(k).unwrap().to::<u16>();
                if correlations[k] < floor || correlations[k] > ONE {
                    return Err(invalid_correlation(
                        symbols[i],
                        symbols[j],
                        correlations[k],
                        floor,
                    ));
                }
            }
        }
        for (i, &value) in gaps.iter().enumerate() {
            let floor = self.gap_floors.get(i).unwrap().to::<u32>();
            if value < floor || value > MAX_GAP {
                return Err(invalid_gap(symbols[i], value, floor));
            }
        }
        for (k, &value) in depths.iter().enumerate() {
            let ceiling = self.depth_ceilings.get(k).unwrap().to::<u32>();
            if value == 0 || value > ceiling {
                return Err(invalid_depth(symbols[k / 2], value, ceiling));
            }
        }
        if !positive_definite(n, correlations) {
            return Err(MarginError::NotPositiveDefinite(NotPositiveDefinite {}));
        }
        Ok(())
    }
}

fn invalid_volatility(symbol: B256, value: u32, floor: u32) -> MarginError {
    MarginError::InvalidVolatility(InvalidVolatility {
        symbol,
        value,
        floor,
    })
}

fn invalid_depth(symbol: B256, value: u32, ceiling: u32) -> MarginError {
    MarginError::InvalidDepth(InvalidDepth {
        symbol,
        value,
        ceiling,
    })
}

fn invalid_gap(symbol: B256, value: u32, floor: u32) -> MarginError {
    MarginError::InvalidGap(InvalidGap {
        symbol,
        value,
        floor,
    })
}

fn invalid_correlation(symbol: B256, other: B256, value: u16, floor: u16) -> MarginError {
    MarginError::InvalidCorrelation(InvalidCorrelation {
        symbol,
        other,
        value,
        floor,
    })
}

/// Whether `value` is within ×1.5 of `previous` either way.
fn within_step(previous: u64, value: u64) -> bool {
    2 * value <= 3 * previous && 3 * value >= 2 * previous
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::error::{
        AssetCount, CorrelationStepTooLarge, DepthStepTooLarge, DuplicateAsset, ExposureTooLarge,
        GapStepTooLarge, InsufficientGas, InvalidCorrelation, InvalidDepth, InvalidGap,
        InvalidPool, InvalidVolatility, LengthMismatch, NotPositiveDefinite, ScenarioOutOfRange,
        UnknownAsset, UnsupportedScenarioSize, UpdateTooSoon, VolatilityStepTooLarge, ZeroSymbol,
    };
    use ownable::{OwnableInvalidOwner, OwnableUnauthorizedAccount, OwnershipTransferred};
    use proptest::prelude::*;
    use stylus_sdk::alloy_primitives::keccak256;
    use stylus_sdk::alloy_primitives::{U256, address, b256, hex};
    use stylus_sdk::alloy_sol_types::{SolEvent, SolValue};
    use stylus_sdk::testing::*;

    const OWNER: Address = address!("0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac");
    const NVDA: B256 = b256!("0x4e56444100000000000000000000000000000000000000000000000000000000");
    const TSLA: B256 = b256!("0x54534c4100000000000000000000000000000000000000000000000000000000");
    const SPY: B256 = b256!("0x5350590000000000000000000000000000000000000000000000000000000000");
    const T0: u64 = 1_790_000_000;
    const VOLATILITY_FLOORS: [u32; 3] = [31_352, 37_436, 11_335];
    const VOLATILITIES: [u32; 3] = [31_352, 37_436, 11_335];
    const CORRELATION_FLOORS: [u16; 3] = [4_637, 7_203, 5_178];
    const CORRELATIONS: [u16; 3] = [4_637, 7_203, 6_232];
    const GAP_FLOORS: [u32; 3] = [118_601, 135_345, 54_840];
    const GAPS: [u32; 3] = [118_601, 135_345, 54_840];
    const DEPTHS: [u32; 6] = [3_561_773, 1_377_156, 164_067, 166_576, 320_861, 877_772];

    struct Config {
        symbols: Vec<B256>,
        volatility_floors: Vec<u32>,
        volatilities: Vec<u32>,
        correlation_floors: Vec<u16>,
        correlations: Vec<u16>,
        gap_floors: Vec<u32>,
        gaps: Vec<u32>,
        depths: Vec<u32>,
        market: B256,
        pools: Vec<Address>,
        tokens: Vec<Address>,
        usdg: Address,
        weth: Address,
        eth_usd: Address,
        band: Address,
        owner: Address,
    }

    fn config() -> Config {
        Config {
            symbols: vec![NVDA, TSLA, SPY],
            volatility_floors: VOLATILITY_FLOORS.to_vec(),
            volatilities: VOLATILITIES.to_vec(),
            correlation_floors: CORRELATION_FLOORS.to_vec(),
            correlations: CORRELATIONS.to_vec(),
            gap_floors: GAP_FLOORS.to_vec(),
            gaps: GAPS.to_vec(),
            depths: DEPTHS.to_vec(),
            market: SPY,
            pools: vec![Address::ZERO; 3],
            tokens: vec![Address::ZERO; 3],
            usdg: Address::ZERO,
            weth: Address::ZERO,
            eth_usd: Address::ZERO,
            band: Address::ZERO,
            owner: OWNER,
        }
    }

    fn deploy(vm: &TestVM, c: Config) -> Result<Margin, MarginError> {
        let mut margin = Margin::from(vm);
        let assets = (0..c.symbols.len())
            .map(|i| {
                (
                    c.symbols[i],
                    c.volatility_floors[i],
                    c.volatilities[i],
                    c.gap_floors[i],
                    c.gaps[i],
                    c.depths[2 * i],
                    c.depths[2 * i + 1],
                    c.pools[i],
                    c.tokens[i],
                )
            })
            .collect();
        margin.constructor(
            assets,
            c.correlation_floors,
            c.correlations,
            c.market,
            c.usdg,
            c.weth,
            c.eth_usd,
            c.band,
            c.owner,
        )?;
        Ok(margin)
    }

    fn deployed(vm: &TestVM) -> Margin {
        vm.set_block_timestamp(T0);
        vm.set_sender(OWNER);
        deploy(vm, config()).unwrap()
    }

    fn symbol(name: &str) -> B256 {
        let mut word = [0u8; 32];
        word[..name.len()].copy_from_slice(name.as_bytes());
        B256::from(word)
    }

    fn volatility_set(symbol: B256, value: u32) -> (Vec<B256>, Vec<u8>) {
        (
            vec![VolatilitySet::SIGNATURE_HASH, symbol],
            value.abi_encode(),
        )
    }

    fn depth_set(symbol: B256, selling: u32, buying: u32) -> (Vec<B256>, Vec<u8>) {
        let event = DepthSet {
            symbol,
            selling,
            buying,
        };
        (vec![DepthSet::SIGNATURE_HASH, symbol], event.encode_data())
    }

    fn gap_set(symbol: B256, value: u32) -> (Vec<B256>, Vec<u8>) {
        (vec![GapSet::SIGNATURE_HASH, symbol], value.abi_encode())
    }

    fn correlation_set(symbol: B256, other: B256, value: u16) -> (Vec<B256>, Vec<u8>) {
        (
            vec![CorrelationSet::SIGNATURE_HASH, symbol, other],
            value.abi_encode(),
        )
    }

    #[test]
    fn the_constructor_stores_the_parameters_and_their_floors() {
        let vm = TestVM::default();
        let margin = deployed(&vm);
        assert_eq!(margin.assets(), vec![NVDA, TSLA, SPY]);
        assert_eq!(margin.volatility(NVDA), Ok((31_352, 31_352)));
        assert_eq!(margin.volatility(SPY), Ok((11_335, 11_335)));
        assert_eq!(margin.correlation(NVDA, TSLA), Ok((4_637, 4_637)));
        assert_eq!(margin.correlation(SPY, TSLA), Ok((6_232, 5_178)));
        assert_eq!(margin.correlation(TSLA, SPY), margin.correlation(SPY, TSLA));
        assert_eq!(margin.correlation(SPY, SPY), Ok((ONE, ONE)));
        assert_eq!(margin.last_update(), 0);
        assert_eq!(margin.owner(), OWNER);
        let unknown = symbol("AMZN");
        let err = MarginError::UnknownAsset(UnknownAsset { symbol: unknown });
        assert_eq!(margin.volatility(unknown), Err(err.clone()));
        assert_eq!(margin.correlation(NVDA, unknown), Err(err));
        let logs = vm.get_emitted_logs();
        assert_eq!(logs[0].0[0], OwnershipTransferred::SIGNATURE_HASH);
        assert_eq!(
            logs[1..].to_vec(),
            vec![
                volatility_set(NVDA, 31_352),
                volatility_set(TSLA, 37_436),
                volatility_set(SPY, 11_335),
                correlation_set(NVDA, TSLA, 4_637),
                correlation_set(NVDA, SPY, 7_203),
                correlation_set(TSLA, SPY, 6_232),
                gap_set(NVDA, 118_601),
                gap_set(TSLA, 135_345),
                gap_set(SPY, 54_840),
                depth_set(NVDA, 3_561_773, 1_377_156),
                depth_set(TSLA, 164_067, 166_576),
                depth_set(SPY, 320_861, 877_772),
            ]
        );
        assert_eq!(margin.weekend_gap(TSLA), Ok((135_345, 135_345)));
        assert_eq!(margin.market(), SPY);
        assert_eq!(margin.depth(TSLA), Ok((164_067, 166_576, 164_067, 166_576)));
        assert_eq!(margin.pool(TSLA), Ok(Address::ZERO));
        assert_eq!(margin.eth_usd_feed(), Address::ZERO);
    }

    #[test]
    fn the_constructor_rejects_an_invalid_configuration() {
        let reject = |change: &dyn Fn(&mut Config), expected: MarginError| {
            let vm = TestVM::default();
            let mut c = config();
            change(&mut c);
            assert_eq!(deploy(&vm, c).err(), Some(expected));
        };
        let count = |n: usize| {
            MarginError::AssetCount(AssetCount {
                count: U256::from(n),
            })
        };
        let volatility = |symbol, value, floor| {
            MarginError::InvalidVolatility(InvalidVolatility {
                symbol,
                value,
                floor,
            })
        };
        let correlation = |symbol, other, value, floor| {
            MarginError::InvalidCorrelation(InvalidCorrelation {
                symbol,
                other,
                value,
                floor,
            })
        };
        reject(
            &|c| c.owner = Address::ZERO,
            MarginError::OwnableInvalidOwner(OwnableInvalidOwner {
                owner: Address::ZERO,
            }),
        );
        reject(
            &|c| {
                c.symbols.clear();
                c.volatility_floors.clear();
                c.volatilities.clear();
                c.correlation_floors.clear();
                c.correlations.clear();
            },
            count(0),
        );
        reject(
            &|c| {
                c.symbols = (0..9).map(|i| symbol(&format!("A{i}"))).collect();
                c.volatility_floors = vec![1; 9];
                c.volatilities = vec![1; 9];
                c.correlation_floors = vec![1; 36];
                c.correlations = vec![1; 36];
                c.gap_floors = vec![1; 9];
                c.gaps = vec![1; 9];
                c.depths = vec![1; 18];
                c.pools = vec![Address::ZERO; 9];
                c.tokens = vec![Address::ZERO; 9];
            },
            count(9),
        );
        reject(
            &|c| c.correlation_floors.push(1),
            MarginError::LengthMismatch(LengthMismatch {}),
        );
        reject(
            &|c| c.symbols[1] = B256::ZERO,
            MarginError::ZeroSymbol(ZeroSymbol {}),
        );
        reject(
            &|c| c.symbols[2] = NVDA,
            MarginError::DuplicateAsset(DuplicateAsset { symbol: NVDA }),
        );
        reject(&|c| c.volatility_floors[1] = 0, volatility(TSLA, 37_436, 0));
        reject(
            &|c| c.volatility_floors[1] = MAX_VOLATILITY + 1,
            volatility(TSLA, 37_436, MAX_VOLATILITY + 1),
        );
        reject(
            &|c| c.volatilities[2] = 11_334,
            volatility(SPY, 11_334, 11_335),
        );
        reject(
            &|c| {
                c.volatility_floors[2] = MAX_VOLATILITY;
                c.volatilities[2] = MAX_VOLATILITY + 1;
            },
            volatility(SPY, MAX_VOLATILITY + 1, MAX_VOLATILITY),
        );
        reject(
            &|c| c.correlation_floors[0] = 0,
            correlation(NVDA, TSLA, 4_637, 0),
        );
        reject(
            &|c| c.correlation_floors[0] = ONE + 1,
            correlation(NVDA, TSLA, 4_637, ONE + 1),
        );
        reject(
            &|c| c.correlations[2] = 5_177,
            correlation(TSLA, SPY, 5_177, 5_178),
        );
        reject(
            &|c| c.correlations[2] = ONE + 1,
            correlation(TSLA, SPY, ONE + 1, 5_178),
        );
        reject(
            &|c| {
                c.correlation_floors = vec![1_000, 1_000, 1_000];
                c.correlations = vec![9_000, 9_000, 1_000];
            },
            MarginError::NotPositiveDefinite(NotPositiveDefinite {}),
        );
    }

    #[test]
    fn the_engine_takes_up_to_eight_assets() {
        let vm = TestVM::default();
        let symbols: Vec<B256> = (0..8).map(|i| symbol(&format!("A{i}"))).collect();
        let margin = deploy(
            &vm,
            Config {
                symbols: symbols.clone(),
                volatility_floors: vec![1; 8],
                volatilities: vec![MAX_VOLATILITY; 8],
                correlation_floors: vec![1; 28],
                correlations: vec![9_999; 28],
                gap_floors: vec![1; 8],
                gaps: vec![MAX_GAP; 8],
                depths: vec![u32::MAX; 16],
                market: B256::ZERO,
                pools: vec![Address::ZERO; 8],
                tokens: vec![Address::ZERO; 8],
                ..config()
            },
        )
        .unwrap();
        assert_eq!(margin.assets(), symbols);
        assert_eq!(margin.correlation(symbols[7], symbols[6]), Ok((9_999, 1)));
    }

    #[test]
    fn a_single_asset_needs_no_correlation() {
        let vm = TestVM::default();
        let margin = deploy(
            &vm,
            Config {
                symbols: vec![SPY],
                volatility_floors: vec![11_335],
                volatilities: vec![11_335],
                correlation_floors: vec![],
                correlations: vec![],
                gap_floors: vec![54_840],
                gaps: vec![54_840],
                depths: DEPTHS[4..].to_vec(),
                pools: vec![Address::ZERO],
                tokens: vec![Address::ZERO],
                ..config()
            },
        )
        .unwrap();
        assert_eq!(margin.correlation(SPY, SPY), Ok((ONE, ONE)));
    }

    #[test]
    fn the_first_update_may_come_at_once_and_each_next_a_day_later() {
        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                vec![4_700, 7_203, 6_232],
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
            .unwrap();
        assert_eq!(margin.last_update(), T0);
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL - 1);
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::UpdateTooSoon(UpdateTooSoon {
                nextUpdateAt: T0 + UPDATE_INTERVAL
            }))
        );
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
            .unwrap();
        assert_eq!(margin.last_update(), T0 + UPDATE_INTERVAL);
        assert_eq!(margin.correlation(NVDA, TSLA), Ok((4_637, 4_637)));
    }

    #[test]
    fn a_value_moves_by_at_most_one_and_a_half_times_either_way() {
        let rise = |vols: [u32; 3], corrs: [u16; 3]| {
            let vm = TestVM::default();
            let mut margin = deployed(&vm);
            margin.set_parameters(
                vols.to_vec(),
                corrs.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
        };
        assert_eq!(rise([47_028, 37_436, 11_335], CORRELATIONS), Ok(()));
        assert_eq!(
            rise([47_029, 37_436, 11_335], CORRELATIONS),
            Err(MarginError::VolatilityStepTooLarge(
                VolatilityStepTooLarge {
                    symbol: NVDA,
                    previous: 31_352,
                    value: 47_029
                }
            ))
        );
        assert_eq!(rise(VOLATILITIES, [6_955, 7_203, 6_232]), Ok(()));
        assert_eq!(
            rise(VOLATILITIES, [6_956, 7_203, 6_232]),
            Err(MarginError::CorrelationStepTooLarge(
                CorrelationStepTooLarge {
                    symbol: NVDA,
                    other: TSLA,
                    previous: 4_637,
                    value: 6_956
                }
            ))
        );

        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        margin
            .set_parameters(
                vec![31_352, 37_436, 17_002],
                vec![4_637, 7_203, 9_348],
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
            .unwrap();
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        assert_eq!(
            margin.set_parameters(
                vec![31_352, 37_436, 11_334 + 1],
                vec![4_637, 7_203, 6_231],
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::CorrelationStepTooLarge(
                CorrelationStepTooLarge {
                    symbol: TSLA,
                    other: SPY,
                    previous: 9_348,
                    value: 6_231
                }
            ))
        );
        assert_eq!(
            margin.set_parameters(
                vec![31_352, 37_436, 11_334 + 1],
                vec![4_637, 7_203, 6_232],
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Ok(())
        );
        assert_eq!(margin.volatility(SPY), Ok((11_335, 11_335)));
    }

    #[test]
    fn an_update_keeps_the_floors_and_a_positive_definite_matrix() {
        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        assert_eq!(
            margin.set_parameters(
                vec![31_351, 37_436, 11_335],
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::InvalidVolatility(InvalidVolatility {
                symbol: NVDA,
                value: 31_351,
                floor: 31_352
            }))
        );
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                vec![4_637, 7_202, 6_232],
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::InvalidCorrelation(InvalidCorrelation {
                symbol: NVDA,
                other: SPY,
                value: 7_202,
                floor: 7_203
            }))
        );
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                vec![4_637, 9_990, 5_178],
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::NotPositiveDefinite(NotPositiveDefinite {}))
        );
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                vec![4_637; 2],
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            Err(MarginError::LengthMismatch(LengthMismatch {}))
        );
        assert_eq!(margin.last_update(), 0);
        assert_eq!(margin.correlation(NVDA, SPY), Ok((7_203, 7_203)));
    }

    #[test]
    fn an_update_logs_only_the_values_that_change() {
        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        let before = vm.get_emitted_logs().len();
        margin
            .set_parameters(
                vec![31_352, 40_000, 11_335],
                vec![4_637, 7_203, 7_000],
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
            .unwrap();
        assert_eq!(
            vm.get_emitted_logs()[before..].to_vec(),
            vec![
                volatility_set(TSLA, 40_000),
                correlation_set(TSLA, SPY, 7_000)
            ]
        );
    }

    #[test]
    fn only_the_owner_updates_and_ownership_moves_in_two_steps() {
        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        let next = address!("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC");
        let unauthorized = |account| {
            Err(MarginError::OwnableUnauthorizedAccount(
                OwnableUnauthorizedAccount { account },
            ))
        };
        vm.set_sender(next);
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            unauthorized(next)
        );
        assert_eq!(margin.transfer_ownership(next), unauthorized(next));
        vm.set_sender(OWNER);
        margin.transfer_ownership(next).unwrap();
        assert_eq!(margin.pending_owner(), next);
        vm.set_sender(next);
        margin.accept_ownership().unwrap();
        assert_eq!(
            (margin.owner(), margin.pending_owner()),
            (next, Address::ZERO)
        );
        vm.set_sender(OWNER);
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            unauthorized(OWNER)
        );
        vm.set_sender(next);
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec(),
            )
            .unwrap();
        margin.renounce_ownership().unwrap();
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                DEPTHS.to_vec()
            ),
            unauthorized(next)
        );
    }

    #[test]
    fn the_weekend_gap_keeps_the_same_rules() {
        let reject = |change: &dyn Fn(&mut Config), expected: MarginError| {
            let vm = TestVM::default();
            let mut c = config();
            change(&mut c);
            assert_eq!(deploy(&vm, c).err(), Some(expected));
        };
        let gap = |symbol, value, floor| {
            MarginError::InvalidGap(InvalidGap {
                symbol,
                value,
                floor,
            })
        };
        reject(&|c| c.gap_floors[0] = 0, gap(NVDA, 118_601, 0));
        reject(&|c| c.gaps[1] = 135_344, gap(TSLA, 135_344, 135_345));
        reject(
            &|c| {
                c.gap_floors[2] = MAX_GAP;
                c.gaps[2] = MAX_GAP + 1;
            },
            gap(SPY, MAX_GAP + 1, MAX_GAP),
        );

        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        let before = vm.get_emitted_logs().len();
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                vec![177_902, 135_345, 54_840],
                DEPTHS.to_vec()
            ),
            Err(MarginError::GapStepTooLarge(GapStepTooLarge {
                symbol: NVDA,
                previous: 118_601,
                value: 177_902
            }))
        );
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                vec![177_901, 135_345, 54_840],
                DEPTHS.to_vec(),
            )
            .unwrap();
        assert_eq!(
            vm.get_emitted_logs()[before..].to_vec(),
            vec![gap_set(NVDA, 177_901)]
        );
        assert_eq!(margin.weekend_gap(NVDA), Ok((177_901, 118_601)));
    }

    #[test]
    fn the_market_is_one_of_the_assets_or_none() {
        let vm = TestVM::default();
        let mut c = config();
        c.market = symbol("AMZN");
        assert_eq!(
            deploy(&vm, c).err(),
            Some(MarginError::UnknownAsset(UnknownAsset {
                symbol: symbol("AMZN")
            }))
        );
        let vm = TestVM::default();
        let mut c = config();
        c.market = B256::ZERO;
        assert_eq!(deploy(&vm, c).unwrap().market(), B256::ZERO);
    }

    #[test]
    fn scenarios_come_by_index_and_as_a_digest_whatever_the_order_of_the_assets() {
        let vm = TestVM::default();
        let margin = deployed(&vm);
        let symbols: Vec<[u8; 32]> = [NVDA, TSLA, SPY].iter().map(|s| s.0).collect();
        let parameters = Parameters {
            symbols: &symbols,
            volatilities: &VOLATILITIES,
            correlations: &CORRELATIONS,
            gaps: &GAPS,
            market: Some(2),
        };
        let set = Set::new(&parameters, SCENARIO_SIZE, 172_800);
        for index in [0u16, 255, 256, 767, 768, 771, 772, 777] {
            assert_eq!(margin.scenario(index, 172_800), Ok(set.row(index.into())));
        }
        assert_eq!(
            margin.scenario(778, 172_800),
            Err(MarginError::ScenarioOutOfRange(ScenarioOutOfRange {
                index: 778
            }))
        );
        assert_eq!(
            margin.scenario_digest(100, 172_800),
            Err(MarginError::UnsupportedScenarioSize(
                UnsupportedScenarioSize { size: 100 }
            ))
        );
        let digest = margin.scenario_digest(256, 172_800).unwrap();
        assert_eq!(digest, keccak256(set.encoded()));

        let vm = TestVM::default();
        vm.set_sender(OWNER);
        let reordered = deploy(
            &vm,
            Config {
                symbols: vec![SPY, NVDA, TSLA],
                volatility_floors: vec![11_335, 31_352, 37_436],
                volatilities: vec![11_335, 31_352, 37_436],
                correlation_floors: vec![7_203, 5_178, 4_637],
                correlations: vec![7_203, 6_232, 4_637],
                gap_floors: vec![54_840, 118_601, 135_345],
                gaps: vec![54_840, 118_601, 135_345],
                depths: [&DEPTHS[4..], &DEPTHS[..4]].concat(),
                ..config()
            },
        )
        .unwrap();
        assert_eq!(reordered.scenario_digest(256, 172_800), Ok(digest));
        assert_eq!(
            reordered.scenario(0, 172_800).unwrap()[1],
            margin.scenario(0, 172_800).unwrap()[0]
        );
    }

    /// Per-mille factors on the current volatilities and correlations.
    fn factors(
        volatility: core::ops::RangeInclusive<u64>,
        correlation: core::ops::RangeInclusive<u64>,
    ) -> impl Strategy<Value = (Vec<u64>, Vec<u64>)> {
        (
            prop::collection::vec(volatility, 3),
            prop::collection::vec(correlation, 3),
        )
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, ..ProptestConfig::default() })]

        #[test]
        fn an_accepted_update_keeps_every_bound_and_a_rejected_one_changes_nothing(
            first in factors(1_000..=1_200, 1_000..=1_050),
            rest in prop::collection::vec(factors(600..=1_600, 800..=1_300), 0..5),
        ) {
            let vm = TestVM::default();
            let mut margin = deployed(&vm);
            let mut now = T0;
            let mut accepted = 0;
            for (up, across) in core::iter::once(first).chain(rest) {
                let volatilities: Vec<u32> = [NVDA, TSLA, SPY]
                    .iter()
                    .zip(&up)
                    .map(|(&s, &f)| (u64::from(margin.volatility(s).unwrap().0) * f / 1_000).max(1) as u32)
                    .collect();
                let correlations: Vec<u16> = [(NVDA, TSLA), (NVDA, SPY), (TSLA, SPY)]
                    .iter()
                    .zip(&across)
                    .map(|(&(a, b), &f)| (u64::from(margin.correlation(a, b).unwrap().0) * f / 1_000).clamp(1, u64::from(ONE)) as u16)
                    .collect();
                let before = (margin.volatility(NVDA), margin.volatility(TSLA), margin.volatility(SPY),
                    margin.correlation(NVDA, TSLA), margin.correlation(NVDA, SPY), margin.correlation(TSLA, SPY));
                let previous: Vec<u32> = [NVDA, TSLA, SPY].iter().map(|&s| margin.volatility(s).unwrap().0).collect();
                let last = margin.last_update();
                match margin.set_parameters(volatilities.clone(), correlations.clone(), GAPS.to_vec(), DEPTHS.to_vec()) {
                    Ok(()) => {
                        for (i, &s) in [NVDA, TSLA, SPY].iter().enumerate() {
                            let (value, floor) = margin.volatility(s).unwrap();
                            prop_assert!(value >= floor && 2 * value <= 3 * previous[i] && 3 * value >= 2 * previous[i]);
                        }
                        let (a, b, c) = (margin.correlation(NVDA, TSLA).unwrap(), margin.correlation(NVDA, SPY).unwrap(),
                            margin.correlation(TSLA, SPY).unwrap());
                        prop_assert!(a.0 >= a.1 && b.0 >= b.1 && c.0 >= c.1);
                        prop_assert!(positive_definite(3, &[a.0, b.0, c.0]));
                        prop_assert_eq!(margin.last_update(), now);
                        accepted += 1;
                    }
                    Err(_) => {
                        prop_assert_eq!(before, (margin.volatility(NVDA), margin.volatility(TSLA), margin.volatility(SPY),
                            margin.correlation(NVDA, TSLA), margin.correlation(NVDA, SPY), margin.correlation(TSLA, SPY)));
                        prop_assert_eq!(margin.last_update(), last);
                    }
                }
                now += UPDATE_INTERVAL;
                vm.set_block_timestamp(now);
            }
            prop_assert!(accepted > 0);
        }
    }

    #[test]
    fn the_constructor_refuses_a_pool_it_cannot_read() {
        let vm = TestVM::default();
        let pool = address!("0x0000000000000000000000000000000000000d05");
        vm.mock_static_call(pool, hex::decode("0dfe1681").unwrap(), Err(vec![]));
        assert_eq!(
            deploy(
                &vm,
                Config {
                    pools: vec![pool, Address::ZERO, Address::ZERO],
                    ..config()
                }
            )
            .err(),
            Some(MarginError::InvalidPool(InvalidPool { symbol: NVDA, pool }))
        );
    }

    #[test]
    fn the_first_depths_are_their_ceilings_and_none_is_zero() {
        let mut zero = DEPTHS;
        zero[3] = 0;
        assert_eq!(
            deploy(
                &TestVM::default(),
                Config {
                    depths: zero.to_vec(),
                    ..config()
                }
            )
            .err(),
            Some(MarginError::InvalidDepth(InvalidDepth {
                symbol: TSLA,
                value: 0,
                ceiling: 0
            }))
        );
        let mut halved = DEPTHS;
        halved[0] /= 2;
        let vm = TestVM::default();
        let margin = deploy(
            &vm,
            Config {
                depths: halved.to_vec(),
                ..config()
            },
        )
        .unwrap();
        assert_eq!(
            margin.depth(NVDA),
            Ok((1_780_886, 1_377_156, 1_780_886, 1_377_156))
        );
    }

    #[test]
    fn a_depth_falls_at_once_but_rises_by_at_most_half_again_a_day() {
        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        vm.set_sender(OWNER);
        vm.set_block_timestamp(T0);
        let mut fallen = DEPTHS;
        fallen[1] = 137_715;
        let before = vm.get_emitted_logs().len();
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                fallen.to_vec(),
            )
            .unwrap();
        assert_eq!(
            margin.depth(NVDA),
            Ok((3_561_773, 137_715, 3_561_773, 1_377_156))
        );
        assert_eq!(
            vm.get_emitted_logs()[before..].to_vec(),
            vec![depth_set(NVDA, 3_561_773, 137_715)]
        );
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        let mut risen = fallen;
        risen[1] = 137_715 * 3 / 2 + 1;
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                risen.to_vec()
            ),
            Err(MarginError::DepthStepTooLarge(DepthStepTooLarge {
                symbol: NVDA,
                previous: 137_715,
                value: 137_715 * 3 / 2 + 1
            }))
        );
        let mut above = fallen;
        above[2] = 164_068;
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                above.to_vec()
            ),
            Err(MarginError::InvalidDepth(InvalidDepth {
                symbol: TSLA,
                value: 164_068,
                ceiling: 164_067
            }))
        );
        risen[1] -= 1;
        margin
            .set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                GAPS.to_vec(),
                risen.to_vec(),
            )
            .unwrap();
        assert_eq!(margin.depth(NVDA).unwrap().1, 206_572);
    }

    #[test]
    fn the_requirement_takes_one_quantity_and_price_per_asset() {
        let vm = TestVM::default();
        let margin = deployed(&vm);
        let prices = vec![U256::from(22_886_000_000u64); 3];
        assert_eq!(
            margin.requirement(vec![I256::ZERO; 3], prices.clone(), 172_800, true),
            Ok((U256::ZERO, 0))
        );
        assert_eq!(
            margin.requirement(vec![I256::ZERO; 2], prices.clone(), 172_800, false),
            Err(MarginError::LengthMismatch(LengthMismatch {}))
        );
        assert_eq!(
            margin.requirement(vec![I256::ZERO; 3], prices[..2].to_vec(), 172_800, false),
            Err(MarginError::LengthMismatch(LengthMismatch {}))
        );
        let huge = I256::try_from(10u128.pow(38)).unwrap();
        assert_eq!(
            margin.requirement(
                vec![I256::ZERO, huge, I256::ZERO],
                prices.clone(),
                172_800,
                false
            ),
            Err(MarginError::ExposureTooLarge(ExposureTooLarge {
                symbol: TSLA
            }))
        );
        assert_eq!(
            margin.requirement(
                vec![I256::MIN, I256::ZERO, I256::ZERO],
                prices,
                172_800,
                false
            ),
            Err(MarginError::ExposureTooLarge(ExposureTooLarge {
                symbol: NVDA
            }))
        );
    }

    #[test]
    fn an_asset_held_without_a_pool_sets_its_bit() {
        let vm = TestVM::default();
        let margin = deployed(&vm);
        let one = I256::try_from(1_000_000_000_000_000_000u128).unwrap();
        let prices = vec![
            U256::from(22_886_000_000u64),
            U256::from(35_802_000_000u64),
            U256::from(76_779_000_000u64),
        ];
        let (value, missing) = margin
            .requirement(vec![one, I256::ZERO, -one], prices, 172_800, false)
            .unwrap();
        assert_eq!(missing, 0b101);
        assert_eq!(
            value,
            U256::from(767_790_000_000_000_000_000u128) * U256::from(6) / U256::from(100)
        );
    }

    #[test]
    fn the_current_requirement_follows_the_band_session() {
        let band = address!("0x0000000000000000000000000000000000000b01");
        let vm = TestVM::default();
        let margin = deploy(&vm, Config { band, ..config() }).unwrap();
        assert_eq!(margin.band(), band);
        let one = I256::try_from(1_000_000_000_000_000_000u128).unwrap();
        let prices = vec![
            U256::from(22_886_000_000u64),
            U256::from(35_802_000_000u64),
            U256::from(76_779_000_000u64),
        ];
        let holding = vec![one; 3];
        let (open, _) = margin
            .requirement(holding.clone(), prices.clone(), session::HORIZON, false)
            .unwrap();
        let (closed, _) = margin
            .requirement(holding.clone(), prices.clone(), session::HORIZON, true)
            .unwrap();
        let buffered = (open * U256::from(5)).div_ceil(U256::from(4));
        assert!(closed > buffered);
        assert_eq!(margin.weekend_leverage(), 50_000);
        let gross = prices.iter().fold(U256::ZERO, |sum, &p| {
            sum + p * U256::from(10_000_000_000u64)
        });
        let across = closed.max((gross * U256::from(10_000)).div_ceil(U256::from(50_000)));
        assert!(across > closed);
        let now = 1_790_000_000u64;
        vm.set_block_timestamp(now);
        let call = hex::decode("5e3568b8").unwrap();
        let at = |s: (u8, u8, u8, u64, u64)| {
            let words = [s.0.into(), s.1.into(), s.2.into(), s.3, s.4].map(U256::from);
            vm.mock_static_call(band, call.clone(), Ok(words.abi_encode_params()));
            margin
                .current_requirement(holding.clone(), prices.clone())
                .unwrap()
        };
        assert_eq!(at((2, 1, 2, 0, 0)), (buffered, 0b111, 2));
        assert_eq!(
            at((1, 3, 1, 0, now * 1_000 + 3_600_000)),
            (across, 0b111, 1)
        );
        let (value, _, code) = at((2, 3, 1, 0, now * 1_000 + 12_600_000));
        assert_eq!(code, 3);
        assert_eq!(value, buffered + (across - buffered) / U256::from(2));
        vm.mock_static_call(band, call.clone(), Err(vec![]));
        assert_eq!(
            margin
                .current_requirement(holding.clone(), prices.clone())
                .unwrap(),
            (across, 0b111, 0)
        );
        vm.set_gas_left(0);
        assert_eq!(
            margin.current_requirement(holding.clone(), prices.clone()),
            Err(MarginError::InsufficientGas(InsufficientGas {}))
        );
        assert_eq!(at((2, 1, 2, 0, 0)), (buffered, 0b111, 2));
        vm.set_gas_left(u64::MAX);
        let unset = TestVM::default();
        let without = deploy(&unset, config()).unwrap();
        assert_eq!(without.band(), Address::ZERO);
        assert_eq!(
            without.current_requirement(holding, prices).unwrap(),
            (across, 0b111, 0)
        );
    }
}
