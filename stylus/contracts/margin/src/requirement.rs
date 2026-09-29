//! The portfolio requirement: expected shortfall at 99% over the scenario set, capped diversification,
//! stress rows, reference floors and the liquidity add-on. Amounts are USD with 18 decimals.

use alloc::vec::Vec;

use stylus_sdk::alloy_primitives::{I256, U256};

use crate::scenario::{Set, count};

const PPM: i128 = 1_000_000;
/// Largest exposure to one asset: 10 billion USD, with 18 decimals.
pub const MAX_EXPOSURE: i128 = 10_000_000_000 * WAD;
const WAD: i128 = 1_000_000_000_000_000_000;
/// Single-stock move of the reference floor, in millionths: FINRA 4210(g)'s ±15%.
pub const STOCK_MOVE: i32 = 150_000;
/// Broad-index moves of the reference floor, in millionths: FINRA 4210(g)'s −8% and +6%.
pub const INDEX_DOWN: i32 = 80_000;
/// See [`INDEX_DOWN`].
pub const INDEX_UP: i32 = 60_000;
/// The worst losses the expected shortfall at 99% needs, for any lattice size under 300.
const TAIL: usize = 3;

/// What a pool says about liquidating an asset: its time-weighted price in USD with 8 decimals, the USD with
/// 18 decimals it pays out for a 10% fall and takes in for a 10% rise, and its fee in millionths.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct PoolTerms {
    pub price: U256,
    pub selling: U256,
    pub buying: U256,
    pub fee: u32,
}

/// What an asset's pool gives the liquidity add-on.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Pool {
    /// The asset has no pool.
    Absent,
    /// Its pool's oracle, or the ETH/USD answer it needs, cannot be read: the pool's fee and the governed
    /// depth still count.
    Unread { fee: u32 },
    /// What its pool says.
    Read(PoolTerms),
}

/// The worst [`TAIL`] losses seen so far, largest first.
#[derive(Clone, Copy)]
struct Worst([i128; TAIL]);

impl Worst {
    const fn new() -> Self {
        Self([i128::MIN; TAIL])
    }

    fn push(&mut self, loss: i128) {
        let mut loss = loss;
        for slot in &mut self.0 {
            if loss > *slot {
                core::mem::swap(slot, &mut loss);
            }
        }
    }

    /// The mean of the worst 1% of `size` scenarios: whole scenarios, then a fraction of the next one.
    fn shortfall(&self, size: usize) -> i128 {
        let whole = size / 100;
        let part = (size % 100) as i128;
        let sum: i128 = self.0[..whole].iter().sum();
        (100 * sum + part * self.0[whole]) / size as i128
    }
}

/// Each position's value in USD with 18 decimals: its signed `quantities`, with 18 decimals, at its `prices`,
/// with 8, rounded toward zero. `Err` holds the first asset whose value exceeds [`MAX_EXPOSURE`].
pub fn exposures(quantities: &[I256], prices: &[U256]) -> Result<Vec<i128>, usize> {
    quantities
        .iter()
        .zip(prices)
        .enumerate()
        .map(|(i, (&q, &p))| {
            I256::try_from(p)
                .ok()
                .and_then(|p| q.checked_mul(p))
                .map(|v| v / I256::try_from(100_000_000).unwrap())
                .and_then(|v| i128::try_from(v).ok())
                .filter(|v| v.unsigned_abs() <= MAX_EXPOSURE.unsigned_abs())
                .ok_or(i)
        })
        .collect()
}

/// The loss of `exposures` in a scenario of `returns`, both in the order of the assets.
pub fn loss(exposures: &[i128], returns: &[i32]) -> i128 {
    exposures
        .iter()
        .zip(returns)
        .map(|(&e, &r)| -(e * i128::from(r)) / PPM)
        .sum()
}

/// The requirement from the scenarios of `set` for `exposures`, in the order of the assets:
/// - the expected shortfall at 99% of the joint draws, or of the draws with every correlation 0 if larger;
/// - raised so that diversification takes off at most 80% of the gap between that and the sum of each
///   asset's own expected shortfall, over the draws with every correlation 1;
/// - at least the loss when the market moves down or up, and when `spans_closure`, when every asset gaps
///   down or up at once or one asset gaps alone;
/// - at least the loss when one asset alone moves by FINRA 4210(g)'s portfolio-margin range, ∓15% for a
///   stock and −8% or +6% for the market asset, whichever asset loses most.
///
/// Never negative.
pub fn scenario_requirement(set: &Set, exposures: &[i128], spans_closure: bool) -> i128 {
    let (open, closed) = scenario_requirements(set, exposures);
    if spans_closure { closed } else { open }
}

/// [`scenario_requirement`] over the open market and across a closure, from one pass over the set.
#[inline(never)]
pub fn scenario_requirements(set: &Set, exposures: &[i128]) -> (i128, i128) {
    let market = set.parameters().market;
    let size = set.size();
    let n = exposures.len();
    let mut joint = Worst::new();
    let mut apart = Worst::new();
    let mut alone = vec![Worst::new(); n];
    let mut stress = 0;
    let mut gaps = 0;
    for index in 0..count(size, n) {
        let returns = set.row(index);
        if index < size {
            joint.push(loss(exposures, &returns));
        } else if index < 2 * size {
            for (i, worst) in alone.iter_mut().enumerate() {
                worst.push(-(exposures[i] * i128::from(returns[i])) / PPM);
            }
        } else if index < 3 * size {
            apart.push(loss(exposures, &returns));
        } else if index < 3 * size + 2 {
            stress = stress.max(loss(exposures, &returns));
        } else {
            gaps = gaps.max(loss(exposures, &returns));
        }
    }
    let dependent = joint.shortfall(size).max(apart.shortfall(size));
    let standalone: i128 = alone.iter().map(|w| w.shortfall(size)).sum();
    let capped = dependent + (standalone - dependent).max(0) / 5;
    let floor = exposures
        .iter()
        .enumerate()
        .map(|(i, &e)| {
            let adverse = match (Some(i) == market, e > 0) {
                (true, true) => INDEX_DOWN,
                (true, false) => INDEX_UP,
                (false, _) => STOCK_MOVE,
            };
            e.abs() * i128::from(adverse) / PPM
        })
        .max()
        .unwrap_or(0);
    let open = capped.max(stress).max(floor).max(0);
    (open, open.max(gaps))
}

/// The requirement of `exposures` valued at `prices`: [`scenario_requirement`] plus [`liquidity`].
pub fn requirement(
    set: &Set,
    exposures: &[i128],
    prices: &[U256],
    spans_closure: bool,
    depths: &[u32],
    pools: &[Pool],
) -> (U256, u8) {
    let (open, closed, missing) = requirements(set, exposures, prices, depths, pools);
    (if spans_closure { closed } else { open }, missing)
}

/// [`requirement`] over the open market and across a closure, and the missing bits.
#[inline(never)]
pub fn requirements(
    set: &Set,
    exposures: &[i128],
    prices: &[U256],
    depths: &[u32],
    pools: &[Pool],
) -> (U256, U256, u8) {
    let (open, closed) = scenario_requirements(set, exposures);
    let (addon, missing) = liquidity(exposures, prices, set.parameters().gaps, depths, pools);
    (
        U256::from(open) + addon,
        U256::from(closed) + addon,
        missing,
    )
}

/// The [`liquidity_addon`] of every position, on the side a liquidation takes, and a bit for every asset
/// held without a read pool. `depths` holds each asset's selling then buying depth in whole USD. A read
/// pool counts at the smaller of the governed depth and its own, and its own lowers the governed depth by
/// half at most: one second with no liquidity in range empties a harmonic mean. An unread pool counts at the
/// governed depth, with no discount.
pub fn liquidity(
    exposures: &[i128],
    prices: &[U256],
    gaps: &[u32],
    depths: &[u32],
    pools: &[Pool],
) -> (U256, u8) {
    let mut total = U256::ZERO;
    let mut missing = 0;
    for (i, &exposure) in exposures.iter().enumerate() {
        if exposure == 0 {
            continue;
        }
        let long = exposure > 0;
        let governed = U256::from(depths[if long { 2 * i } else { 2 * i + 1 }]) * U256::from(WAD);
        let value = U256::from(exposure.unsigned_abs());
        let (pool_price, depth, fee) = match pools[i] {
            Pool::Absent => {
                missing |= 1 << i;
                continue;
            }
            Pool::Unread { fee } => {
                missing |= 1 << i;
                (prices[i], governed, fee)
            }
            Pool::Read(pool) => (
                pool.price,
                governed.min(
                    (if long { pool.selling } else { pool.buying }).max(governed / U256::from(2)),
                ),
                pool.fee,
            ),
        };
        total += liquidity_addon(value, long, prices[i], pool_price, depth, fee, gaps[i]);
    }
    (total, missing)
}

/// What liquidating one position costs beyond its value, from its pool: the fee, the pool's discount to
/// the valuation price when it works against the position, at most the asset's weekend `gap` in millionths,
/// and the price impact. The impact grows
/// linearly to a 10% move at `depth`, the value the pool absorbs within that move on the position's side,
/// and whatever lies beyond `depth` counts in full. `value` and `depth` are USD with 18 decimals; prices
/// are USD with 8 decimals; `fee` is in millionths.
pub fn liquidity_addon(
    value: U256,
    long: bool,
    price: U256,
    pool_price: U256,
    depth: U256,
    fee: u32,
    gap: u32,
) -> U256 {
    let adverse = if long {
        price.saturating_sub(pool_price)
    } else {
        pool_price.saturating_sub(price)
    }
    .min(price * U256::from(gap) / U256::from(PPM));
    let discount = if price.is_zero() {
        U256::ZERO
    } else {
        value * adverse / price
    };
    let fee = value * U256::from(fee) / U256::from(PPM);
    let impact = if value <= depth {
        if depth.is_zero() {
            U256::ZERO
        } else {
            value * value / (U256::from(20) * depth)
        }
    } else {
        depth / U256::from(20) + (value - depth)
    };
    discount + fee + impact
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::scenario::Parameters;
    use crate::uniswap::{classify, mean, terms};
    use proptest::prelude::*;
    use stylus_sdk::alloy_primitives::{Address, aliases::U160, hex};
    use stylus_sdk::alloy_sol_types::SolValue;

    const PARAMETERS: &str = include_str!("../parameters.json");
    const VECTORS: &str = include_str!("../testdata/requirement-vectors.json");
    const DAYS_2: u64 = 172_800;
    const SYMBOLS: [[u8; 32]; 3] = [
        *b"NVDA\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0",
        *b"TSLA\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0",
        *b"SPY\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0",
    ];
    const VOLATILITIES: [u32; 3] = [31_352, 37_436, 11_335];
    const CORRELATIONS: [u16; 3] = [4_637, 7_203, 6_232];
    const GAPS: [u32; 3] = [118_601, 135_345, 54_840];
    const DOLLAR: i128 = WAD;

    fn three() -> Parameters<'static> {
        Parameters {
            symbols: &SYMBOLS,
            volatilities: &VOLATILITIES,
            correlations: &CORRELATIONS,
            gaps: &GAPS,
            market: Some(2),
        }
    }

    #[test]
    fn the_shortfall_takes_two_and_fifty_six_hundredths_of_the_worst_of_256() {
        let mut worst = Worst::new();
        for loss in [5, 900, 100, 700, 800, -3] {
            worst.push(loss);
        }
        assert_eq!(worst.0, [900, 800, 700]);
        assert_eq!(worst.shortfall(256), (100 * (900 + 800) + 56 * 700) / 256);
        assert_eq!(worst.shortfall(32), 900);
        assert_eq!(worst.shortfall(128), (100 * 900 + 28 * 800) / 128);
    }

    #[test]
    fn the_requirement_matches_the_reference() {
        let p: serde_json::Value = serde_json::from_str(PARAMETERS).unwrap();
        let v: serde_json::Value = serde_json::from_str(VECTORS).unwrap();
        let names: Vec<&str> = v["symbols"]
            .as_array()
            .unwrap()
            .iter()
            .map(|s| s.as_str().unwrap())
            .collect();
        let n = names.len();
        let symbols: Vec<[u8; 32]> = names
            .iter()
            .map(|name| {
                let mut s = [0u8; 32];
                s[..name.len()].copy_from_slice(name.as_bytes());
                s
            })
            .collect();
        let int = |x: &serde_json::Value| x.as_u64().unwrap();
        let volatilities: Vec<u32> = names
            .iter()
            .map(|k| int(&p["volatility"][k]["initial"]) as u32)
            .collect();
        let gaps: Vec<u32> = names
            .iter()
            .map(|k| int(&p["weekendGap"][k]["initial"]) as u32)
            .collect();
        let correlations: Vec<u16> = (0..n)
            .flat_map(|i| (i + 1..n).map(move |j| (i, j)))
            .map(|(i, j)| {
                int(&p["correlation"][format!("{}/{}", names[i], names[j])]["initial"]) as u16
            })
            .collect();
        let depths: Vec<u32> = names
            .iter()
            .flat_map(|k| {
                p["depth"][k]["initial"]
                    .as_array()
                    .unwrap()
                    .iter()
                    .map(|x| x.as_u64().unwrap() as u32)
                    .collect::<Vec<_>>()
            })
            .collect();
        let market = names
            .iter()
            .position(|&k| k == p["market"].as_str().unwrap());
        let parameters = Parameters {
            symbols: &symbols,
            volatilities: &volatilities,
            correlations: &correlations,
            gaps: &gaps,
            market,
        };
        let address = |x: &serde_json::Value| x.as_str().unwrap().parse::<Address>().unwrap();
        let bytes = |x: &serde_json::Value| {
            hex::decode(x.as_str().unwrap().trim_start_matches("0x")).unwrap()
        };
        let round =
            <(U256, I256, U256, U256, U256)>::abi_decode_params(&bytes(&v["ethUsdRound"])).unwrap();
        let pools: Vec<Pool> = names
            .iter()
            .map(|k| {
                let pool = &v["pools"][k];
                let (stock_is_token0, quote_is_weth) = classify(
                    address(&pool["token0"]),
                    address(&pool["token1"]),
                    address(&v["tokens"][k]),
                    address(&v["usdg"]),
                    address(&v["weth"]),
                )
                .unwrap();
                let (ticks, seconds) =
                    <(Vec<I256>, Vec<U256>)>::abi_decode_params(&bytes(&pool["observe"])).unwrap();
                let (tick, liquidity) = mean(
                    ticks[0].as_i64(),
                    ticks[1].as_i64(),
                    U160::from(seconds[0]),
                    U160::from(seconds[1]),
                )
                .unwrap();
                let decimals = [int(&pool["decimals0"]) as u8, int(&pool["decimals1"]) as u8];
                let (ds, dq) = if stock_is_token0 {
                    (decimals[0], decimals[1])
                } else {
                    (decimals[1], decimals[0])
                };
                let usd = if quote_is_weth {
                    round.1.into_raw()
                } else {
                    U256::from(100_000_000)
                };
                let (price, selling, buying) =
                    terms(stock_is_token0, ds, dq, tick, liquidity, usd).unwrap();
                Pool::Read(PoolTerms {
                    price,
                    selling,
                    buying,
                    fee: int(&pool["fee"]) as u32,
                })
            })
            .collect();
        let prices: Vec<U256> = v["prices"]
            .as_array()
            .unwrap()
            .iter()
            .map(|x| U256::from(x.as_u64().unwrap()))
            .collect();
        for vector in v["vectors"].as_array().unwrap() {
            let quantities: Vec<I256> = vector["quantities"]
                .as_array()
                .unwrap()
                .iter()
                .map(|q| q.as_str().unwrap().parse().unwrap())
                .collect();
            let exposures = exposures(&quantities, &prices).unwrap();
            let set = Set::new(&parameters, 256, int(&vector["horizon"]));
            let with: Vec<Pool> = match vector["pools"].as_str().unwrap() {
                "read" => pools.clone(),
                "unread" => pools
                    .iter()
                    .map(|p| match p {
                        Pool::Read(terms) => Pool::Unread { fee: terms.fee },
                        other => *other,
                    })
                    .collect(),
                _ => vec![Pool::Absent; n],
            };
            let closure = vector["spansClosure"].as_bool().unwrap();
            let expected = (
                vector["requirement"]
                    .as_str()
                    .unwrap()
                    .parse::<U256>()
                    .unwrap(),
                int(&vector["missing"]) as u8,
            );
            assert_eq!(
                requirement(&set, &exposures, &prices, closure, &depths, &with),
                expected,
                "{}",
                vector["label"]
            );
        }
    }

    #[test]
    fn the_dev_node_reference_holds() {
        let v: serde_json::Value = serde_json::from_str(VECTORS).unwrap();
        let p: serde_json::Value = serde_json::from_str(PARAMETERS).unwrap();
        let dev = &v["devnode"];
        let quantities: Vec<I256> = dev["quantities"]
            .as_array()
            .unwrap()
            .iter()
            .map(|q| q.as_str().unwrap().parse().unwrap())
            .collect();
        let prices: Vec<U256> = dev["prices"]
            .as_array()
            .unwrap()
            .iter()
            .map(|x| U256::from(x.as_u64().unwrap()))
            .collect();
        let depths: Vec<u32> = ["NVDA", "TSLA", "SPY"]
            .iter()
            .flat_map(|k| {
                p["depth"][k]["initial"]
                    .as_array()
                    .unwrap()
                    .iter()
                    .map(|x| x.as_u64().unwrap() as u32)
                    .collect::<Vec<_>>()
            })
            .collect();
        let stub =
            |stock_is_token0: bool, quote_decimals: u8, tick: i64, liquidity: u128, usd: u64| {
                let per_second = (U256::from(1) << 128) / U256::from(liquidity);
                let (tick, liquidity) = mean(
                    0,
                    tick * 1_800,
                    U160::ZERO,
                    U160::from(per_second * U256::from(1_800)),
                )
                .unwrap();
                let (price, selling, buying) = terms(
                    stock_is_token0,
                    18,
                    quote_decimals,
                    tick,
                    liquidity,
                    U256::from(usd),
                )
                .unwrap();
                Pool::Read(PoolTerms {
                    price,
                    selling,
                    buying,
                    fee: 500,
                })
            };
        let nvda = stub(false, 6, 221_989, 11_245_526_858_841_909_681, 100_000_000);
        let spy = stub(
            true,
            18,
            -12_513,
            16_029_297_629_534_329_325_587,
            268_330_550_000,
        );
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let exposures = exposures(&quantities, &prices).unwrap();
        assert_eq!(
            requirement(
                &set,
                &exposures,
                &prices,
                false,
                &depths,
                &[nvda, Pool::Absent, spy]
            ),
            (
                dev["requirement"]
                    .as_str()
                    .unwrap()
                    .parse::<U256>()
                    .unwrap(),
                dev["missing"].as_u64().unwrap() as u8
            )
        );
    }

    #[test]
    fn diversified_beats_concentrated_in_one_stock() {
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let split = scenario_requirement(&set, &[100_000 * DOLLAR; 3], false);
        let stocks = scenario_requirement(&set, &[150_000 * DOLLAR, 150_000 * DOLLAR, 0], false);
        for stock in 0..2 {
            let mut all = [0; 3];
            all[stock] = 300_000 * DOLLAR;
            let concentrated = scenario_requirement(&set, &all, false);
            assert!(split < concentrated && stocks < concentrated, "{stock}");
        }
    }

    #[test]
    fn the_gap_rows_count_only_across_a_closure() {
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let split = [100_000 * DOLLAR; 3];
        let weekday = scenario_requirement(&set, &split, false);
        let weekend = scenario_requirement(&set, &split, true);
        assert!(weekend > weekday);
        let gaps_together: i128 = GAPS
            .iter()
            .map(|&g| 100_000 * DOLLAR * i128::from(g) / PPM)
            .sum();
        assert_eq!(weekend, gaps_together);
    }

    #[test]
    fn the_market_move_counts_over_the_open_market_and_across_a_closure() {
        let parameters = three();
        let set = Set::new(&parameters, 256, 5 * DAYS_2);
        let short = [0, 0, -100_000 * DOLLAR];
        let up = loss(&short, &set.row(3 * 256 + 1));
        assert!(up > 6_000 * DOLLAR);
        assert_eq!(scenario_requirements(&set, &short), (up, up));
    }

    #[test]
    fn a_hedge_still_pays_for_one_leg_gapping_alone_across_a_closure() {
        let gaps = [200_000, 135_345, 54_840];
        let parameters = Parameters {
            gaps: &gaps,
            ..three()
        };
        let set = Set::new(&parameters, 256, DAYS_2);
        let hedge = [100_000 * DOLLAR, -100_000 * DOLLAR, 0];
        assert!(scenario_requirement(&set, &hedge, false) < 20_000 * DOLLAR);
        assert_eq!(scenario_requirement(&set, &hedge, true), 20_000 * DOLLAR);
    }

    #[test]
    fn a_single_stock_needs_its_fifteen_percent_and_the_index_its_eight_or_six() {
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let alone = |i: usize, e: i128| {
            let mut exposures = [0; 3];
            exposures[i] = e;
            scenario_requirement(&set, &exposures, false)
        };
        assert_eq!(alone(0, 100_000 * DOLLAR), 15_000 * DOLLAR);
        assert_eq!(alone(0, -100_000 * DOLLAR), 15_000 * DOLLAR);
        assert_eq!(alone(2, 100_000 * DOLLAR), 8_000 * DOLLAR);
        assert_eq!(alone(2, -100_000 * DOLLAR), 6_000 * DOLLAR);
    }

    #[test]
    fn the_shallower_of_the_two_depths_counts_on_the_liquidation_side() {
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let dollars = |d: u64| U256::from(d) * U256::from(WAD);
        let pool = |selling: u64, buying: u64| {
            Pool::Read(PoolTerms {
                price: U256::from(10_000_000_000u64),
                selling: dollars(selling),
                buying: dollars(buying),
                fee: 0,
            })
        };
        let prices = [U256::from(10_000_000_000u64); 3];
        let depths = [1_000_000, 2_000_000, 1, 1, 1, 1];
        let base = U256::from(scenario_requirement(&set, &[100_000 * DOLLAR, 0, 0], false));
        let addon = |exposure: i128, terms: Pool| {
            requirement(
                &set,
                &[exposure, 0, 0],
                &prices,
                false,
                &depths,
                &[terms, Pool::Absent, Pool::Absent],
            )
            .0
        };
        let impact =
            |depth: u64| dollars(100_000) * dollars(100_000) / (U256::from(20) * dollars(depth));
        assert_eq!(
            addon(100_000 * DOLLAR, pool(5_000_000, 5_000_000)) - base,
            impact(1_000_000)
        );
        assert_eq!(
            addon(100_000 * DOLLAR, pool(500_000, 5_000_000)) - base,
            impact(500_000)
        );
        let short = U256::from(scenario_requirement(
            &set,
            &[-100_000 * DOLLAR, 0, 0],
            false,
        ));
        assert_eq!(
            addon(-100_000 * DOLLAR, pool(5_000_000, 5_000_000)) - short,
            impact(2_000_000)
        );
        assert_eq!(
            addon(-100_000 * DOLLAR, pool(5_000_000, 1_200_000)) - short,
            impact(1_200_000)
        );
        assert_eq!(
            addon(-100_000 * DOLLAR, Pool::Unread { fee: 0 }) - short,
            impact(2_000_000)
        );
    }

    #[test]
    fn a_pool_emptied_for_a_moment_lowers_the_depth_by_half_at_most() {
        let parameters = three();
        let set = Set::new(&parameters, 256, DAYS_2);
        let dollars = |d: u64| U256::from(d) * U256::from(WAD);
        let emptied = Pool::Read(PoolTerms {
            price: U256::from(10_000_000_000u64),
            selling: U256::ZERO,
            buying: dollars(1_799),
            fee: 0,
        });
        let prices = [U256::from(10_000_000_000u64); 3];
        let depths = [1_000_000, 2_000_000, 1, 1, 1, 1];
        let impact =
            |depth: u64| dollars(100_000) * dollars(100_000) / (U256::from(20) * dollars(depth));
        for (exposure, depth) in [(100_000 * DOLLAR, 500_000), (-100_000 * DOLLAR, 1_000_000)] {
            let base = U256::from(scenario_requirement(&set, &[exposure, 0, 0], false));
            let (total, missing) = requirement(
                &set,
                &[exposure, 0, 0],
                &prices,
                false,
                &depths,
                &[emptied, Pool::Absent, Pool::Absent],
            );
            assert_eq!(total - base, impact(depth));
            assert_eq!(missing, 0);
        }
    }

    #[test]
    fn the_largest_positions_over_the_longest_horizon_fit_in_i128() {
        let gaps = [1_000_000; 3];
        let volatilities = [1_000_000; 3];
        let parameters = Parameters {
            volatilities: &volatilities,
            gaps: &gaps,
            ..three()
        };
        let set = Set::new(&parameters, 256, u64::MAX);
        for exposures in [
            [MAX_EXPOSURE; 3],
            [-MAX_EXPOSURE; 3],
            [MAX_EXPOSURE, -MAX_EXPOSURE, MAX_EXPOSURE],
        ] {
            assert!(scenario_requirement(&set, &exposures, true) > 0);
        }
    }

    #[test]
    fn no_position_is_worth_more_than_ten_billion_dollars() {
        let price = [U256::from(100_000_000u64)];
        let at = |dollars: i128| exposures(&[I256::try_from(dollars * DOLLAR).unwrap()], &price);
        assert_eq!(at(10_000_000_000), Ok(vec![MAX_EXPOSURE]));
        assert_eq!(at(-10_000_000_000), Ok(vec![-MAX_EXPOSURE]));
        assert_eq!(at(10_000_000_001), Err(0));
        assert_eq!(exposures(&[I256::ONE], &[U256::MAX]), Err(0));
    }

    #[test]
    fn a_position_past_its_depth_counts_in_full_beyond_it() {
        let depth = U256::from(1_000_000) * U256::from(WAD);
        let at = |value: U256| {
            liquidity_addon(value, true, U256::from(100), U256::from(100), depth, 0, 0)
        };
        assert_eq!(at(depth), depth / U256::from(20));
        assert_eq!(at(depth / U256::from(2)), depth / U256::from(80));
        assert_eq!(at(depth * U256::from(2)), depth / U256::from(20) + depth);
        assert_eq!(at(U256::ZERO), U256::ZERO);
        assert_eq!(
            liquidity_addon(
                depth,
                true,
                U256::from(100),
                U256::from(100),
                U256::ZERO,
                0,
                0
            ),
            depth
        );
    }

    #[test]
    fn only_a_pool_price_against_the_position_and_the_fee_cost_more() {
        let value = U256::from(100_000) * U256::from(WAD);
        let depth = U256::MAX / U256::from(4);
        let cost = |long: bool, pool: u64, fee: u32| {
            liquidity_addon(
                value,
                long,
                U256::from(10_000_000_000u64),
                U256::from(pool),
                depth,
                fee,
                118_601,
            )
        };
        let impact = cost(true, 10_000_000_000, 0);
        assert_eq!(
            cost(true, 9_900_000_000, 0) - impact,
            value / U256::from(100)
        );
        assert_eq!(cost(true, 10_100_000_000, 0), impact);
        assert_eq!(
            cost(false, 10_100_000_000, 0) - impact,
            value / U256::from(100)
        );
        assert_eq!(cost(false, 9_900_000_000, 0), impact);
        assert_eq!(
            cost(true, 10_000_000_000, 3_000) - impact,
            value * U256::from(3) / U256::from(1_000)
        );
        assert_eq!(
            cost(true, 5_000_000_000, 0) - impact,
            value * U256::from(118_601) / U256::from(1_000_000)
        );
    }

    proptest! {
        #[test]
        fn no_requirement_falls_as_a_portfolio_grows(
            exposures in proptest::array::uniform3(-1_000_000 * DOLLAR..1_000_000 * DOLLAR),
            times in 1i128..100,
            closure in any::<bool>(),
        ) {
            let parameters = three();
            let set = Set::new(&parameters, 256, DAYS_2);
            let larger = exposures.map(|e| e * times);
            prop_assert!(scenario_requirement(&set, &larger, closure) >= scenario_requirement(&set, &exposures, closure));
        }

        #[test]
        fn a_short_loses_on_every_up_move(asset in 0usize..3, index in 0usize..count(256, 3), dollars in 1i128..1_000_000) {
            let parameters = three();
            let set = Set::new(&parameters, 256, DAYS_2);
            let returns = set.row(index);
            let mut exposures = [0; 3];
            exposures[asset] = -dollars * DOLLAR;
            prop_assert_eq!(loss(&exposures, &returns) > 0, returns[asset] > 0);
            prop_assert_eq!(loss(&exposures, &returns) < 0, returns[asset] < 0);
        }

        #[test]
        fn the_addon_never_falls_as_a_position_grows(
            smaller in 0u128..10_000_000_000_000_000_000_000_000,
            extra in 0u128..10_000_000_000_000_000_000_000_000,
            depth in 1u128..10_000_000_000_000_000_000_000_000,
            long in any::<bool>(),
        ) {
            let at = |v: u128| liquidity_addon(U256::from(v), long, U256::from(100), U256::from(99), U256::from(depth), 500, 118_601);
            prop_assert!(at(smaller + extra) >= at(smaller));
        }
    }
}
