#![cfg_attr(not(any(test, feature = "export-abi")), no_main)]
#![cfg_attr(not(any(test, feature = "export-abi")), no_std)]

#[macro_use]
extern crate alloc;

pub mod error;
pub mod matrix;
pub mod scenario;

use alloc::vec::Vec;

use ownable::{IOwnable2Step, Ownable2Step};
use stylus_sdk::alloy_primitives::{Address, B256, U8, U16, U32, U64};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::prelude::*;
use stylus_sdk::storage::{
    StorageB256, StorageMap, StorageU8, StorageU16, StorageU32, StorageU64, StorageVec,
};

use crate::error::{
    AssetCount, CorrelationStepTooLarge, DuplicateAsset, GapStepTooLarge, InvalidCorrelation,
    InvalidGap, InvalidVolatility, LengthMismatch, MarginError, NotPositiveDefinite,
    ScenarioOutOfRange, UnknownAsset, UnsupportedScenarioSize, UpdateTooSoon,
    VolatilityStepTooLarge, ZeroSymbol,
};
use crate::matrix::{ONE, pair, positive_definite};
use crate::scenario::{Parameters, SIZES, Set, count};

/// Most assets the engine takes: the exact positive-definiteness check stays within 256 bits up to 9.
pub const MAX_ASSETS: usize = 8;
/// First bit of the weekend gaps in `set_parameters`'s mask of changed values, after the volatilities' and the
/// correlations'.
const GAP_BITS: usize = MAX_ASSETS + MAX_ASSETS * (MAX_ASSETS - 1) / 2;
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
    market: StorageU8,
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
    /// `market` is the asset that stands for the market, or zero for the equal-weighted portfolio of the
    /// assets. `initial_owner` may update the parameters.
    #[constructor]
    #[allow(clippy::too_many_arguments)]
    pub fn constructor(
        &mut self,
        symbols: Vec<B256>,
        volatility_floors: Vec<u32>,
        volatilities: Vec<u32>,
        correlation_floors: Vec<u16>,
        correlations: Vec<u16>,
        gap_floors: Vec<u32>,
        gaps: Vec<u32>,
        market: B256,
        initial_owner: Address,
    ) -> Result<(), MarginError> {
        self.ownable.initialize::<MarginError>(initial_owner)?;
        let n = symbols.len();
        if n == 0 || n > MAX_ASSETS {
            return Err(MarginError::AssetCount(AssetCount {
                count: stylus_sdk::alloy_primitives::U256::from(n),
            }));
        }
        let pairs = n * (n - 1) / 2;
        if volatility_floors.len() != n
            || volatilities.len() != n
            || correlation_floors.len() != pairs
            || correlations.len() != pairs
            || gap_floors.len() != n
            || gaps.len() != n
        {
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
        for i in 0..n {
            let floor = volatility_floors[i];
            if floor == 0 {
                return Err(invalid_volatility(symbols[i], volatilities[i], floor));
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
        for i in 0..n {
            if gap_floors[i] == 0 {
                return Err(invalid_gap(symbols[i], gaps[i], gap_floors[i]));
            }
            self.gap_floors.push(U32::from(gap_floors[i]));
        }
        if market != B256::ZERO {
            let i = self.position(market)?;
            self.market.set(U8::from(i + 1));
        }
        self.check(&symbols, &volatilities, &correlations, &gaps)?;
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
        Ok(())
    }

    /// Replaces every parameter, in the constructor's order. Each value stays at or above its floor and
    /// moves by at most ×1.5 or ÷1.5, the correlation matrix stays positive definite, and an update comes at
    /// least a day after the one before. Emits an event for every value that changes. Owner only.
    pub fn set_parameters(
        &mut self,
        volatilities: Vec<u32>,
        correlations: Vec<u16>,
        gaps: Vec<u32>,
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
        if volatilities.len() != n || correlations.len() != n * (n - 1) / 2 || gaps.len() != n {
            return Err(MarginError::LengthMismatch(LengthMismatch {}));
        }
        self.check(&symbols, &volatilities, &correlations, &gaps)?;
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
        AssetCount, CorrelationStepTooLarge, DuplicateAsset, GapStepTooLarge, InvalidCorrelation,
        InvalidGap, InvalidVolatility, LengthMismatch, NotPositiveDefinite, ScenarioOutOfRange,
        UnknownAsset, UnsupportedScenarioSize, UpdateTooSoon, VolatilityStepTooLarge, ZeroSymbol,
    };
    use ownable::{OwnableInvalidOwner, OwnableUnauthorizedAccount, OwnershipTransferred};
    use proptest::prelude::*;
    use stylus_sdk::alloy_primitives::keccak256;
    use stylus_sdk::alloy_primitives::{U256, address, b256};
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

    struct Config {
        symbols: Vec<B256>,
        volatility_floors: Vec<u32>,
        volatilities: Vec<u32>,
        correlation_floors: Vec<u16>,
        correlations: Vec<u16>,
        gap_floors: Vec<u32>,
        gaps: Vec<u32>,
        market: B256,
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
            market: SPY,
            owner: OWNER,
        }
    }

    fn deploy(vm: &TestVM, c: Config) -> Result<Margin, MarginError> {
        let mut margin = Margin::from(vm);
        margin.constructor(
            c.symbols,
            c.volatility_floors,
            c.volatilities,
            c.correlation_floors,
            c.correlations,
            c.gap_floors,
            c.gaps,
            c.market,
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
            ]
        );
        assert_eq!(margin.weekend_gap(TSLA), Ok((135_345, 135_345)));
        assert_eq!(margin.market(), SPY);
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
            },
            count(9),
        );
        reject(
            &|c| c.volatilities.pop().map(drop).unwrap(),
            MarginError::LengthMismatch(LengthMismatch {}),
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
                market: B256::ZERO,
                owner: OWNER,
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
                market: SPY,
                owner: OWNER,
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
            )
            .unwrap();
        assert_eq!(margin.last_update(), T0);
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL - 1);
        assert_eq!(
            margin.set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec()),
            Err(MarginError::UpdateTooSoon(UpdateTooSoon {
                nextUpdateAt: T0 + UPDATE_INTERVAL
            }))
        );
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        margin
            .set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec())
            .unwrap();
        assert_eq!(margin.last_update(), T0 + UPDATE_INTERVAL);
        assert_eq!(margin.correlation(NVDA, TSLA), Ok((4_637, 4_637)));
    }

    #[test]
    fn a_value_moves_by_at_most_one_and_a_half_times_either_way() {
        let rise = |vols: [u32; 3], corrs: [u16; 3]| {
            let vm = TestVM::default();
            let mut margin = deployed(&vm);
            margin.set_parameters(vols.to_vec(), corrs.to_vec(), GAPS.to_vec())
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
            )
            .unwrap();
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        assert_eq!(
            margin.set_parameters(
                vec![31_352, 37_436, 11_334 + 1],
                vec![4_637, 7_203, 6_231],
                GAPS.to_vec()
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
                GAPS.to_vec()
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
                GAPS.to_vec()
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
                GAPS.to_vec()
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
                GAPS.to_vec()
            ),
            Err(MarginError::NotPositiveDefinite(NotPositiveDefinite {}))
        );
        assert_eq!(
            margin.set_parameters(VOLATILITIES.to_vec(), vec![4_637; 2], GAPS.to_vec()),
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
            margin.set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec()),
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
            margin.set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec()),
            unauthorized(OWNER)
        );
        vm.set_sender(next);
        margin
            .set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec())
            .unwrap();
        margin.renounce_ownership().unwrap();
        vm.set_block_timestamp(T0 + UPDATE_INTERVAL);
        assert_eq!(
            margin.set_parameters(VOLATILITIES.to_vec(), CORRELATIONS.to_vec(), GAPS.to_vec()),
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
        reject(
            &|c| c.gaps.push(1),
            MarginError::LengthMismatch(LengthMismatch {}),
        );

        let vm = TestVM::default();
        let mut margin = deployed(&vm);
        let before = vm.get_emitted_logs().len();
        assert_eq!(
            margin.set_parameters(
                VOLATILITIES.to_vec(),
                CORRELATIONS.to_vec(),
                vec![177_902, 135_345, 54_840]
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
                market: SPY,
                owner: OWNER,
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
                match margin.set_parameters(volatilities.clone(), correlations.clone(), GAPS.to_vec()) {
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
}
