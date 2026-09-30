// SPDX-License-Identifier: MIT OR Apache-2.0
//! The scenario set: deterministic joint, comonotone and independent draws from the stored parameters,
//! a common market shock, and weekend gap shocks for every asset at once and for each asset alone.

use alloc::vec::Vec;

use crate::matrix::{ONE, pair};

/// Scale of the correlation factor: 1e12.
const S: i128 = 1_000_000_000_000;
/// Returns, volatilities and gaps are in millionths.
const PPM: i128 = 1_000_000;
const DAY: u128 = 86_400;
/// A pivot of the factor below this counts as zero, which keeps every entry within `i128`.
const PIVOT_FLOOR: i128 = S / PPM;
/// The standard normal's expected shortfall at 99%, in millionths.
pub const Z_MARKET: i128 = 2_665_214;
/// Lattice sizes the engine can draw.
pub const SIZES: [usize; 4] = [32, 64, 128, 256];

const MEANS_32: [i32; 16] = [
    -2252211, -1683275, -1420533, -1231294, -1078398, -947376, -830935, -724828, -626336, -533591,
    -445235, -360235, -277767, -197152, -117800, -39186,
];
const MEANS_64: [i32; 32] = [
    -2510185, -1994237, -1764207, -1602344, -1474296, -1366770, -1273110, -1189477, -1113440,
    -1043356, -978060, -916692, -858599, -803270, -750300, -699356, -650167, -602505, -556174,
    -511008, -466862, -423608, -381134, -339336, -298123, -257411, -217121, -177182, -137522,
    -98078, -58787, -19586,
];
const MEANS_128: [i32; 64] = [
    -2747787, -2272583, -2065885, -1922589, -1810660, -1717753, -1637708, -1566980, -1503333,
    -1445259, -1391692, -1341848, -1295133, -1251087, -1209345, -1169610, -1131642, -1095239,
    -1060233, -1026480, -993859, -962262, -931598, -901787, -872756, -844443, -816791, -789750,
    -763275, -737324, -711861, -686851, -662264, -638071, -614246, -590764, -567604, -544744,
    -522166, -499851, -477782, -455943, -434320, -412897, -391663, -370604, -349708, -328964,
    -308360, -287886, -267533, -247289, -227147, -207096, -187128, -167235, -147407, -127637,
    -107918, -88239, -68595, -48978, -29379, -9792,
];
const MEANS_256: [i32; 128] = [
    -2969102, -2526473, -2337436, -2207729, -2107268, -2024501, -1953679, -1891500, -1835885,
    -1785435, -1739163, -1696343, -1656427, -1618989, -1583693, -1550268, -1518490, -1488176,
    -1459172, -1431346, -1404587, -1378797, -1353893, -1329802, -1306459, -1283807, -1261796,
    -1240379, -1219517, -1199172, -1179313, -1159907, -1140929, -1122354, -1104158, -1086320,
    -1068821, -1051644, -1034772, -1018189, -1001881, -985836, -970041, -954484, -939154, -924042,
    -909139, -894434, -879921, -865591, -851436, -837449, -823625, -809957, -796438, -783063,
    -769826, -756723, -743749, -730899, -718169, -705553, -693050, -680653, -668360, -656168,
    -644072, -632070, -620158, -608333, -596593, -584935, -573355, -561852, -550423, -539065,
    -527777, -516555, -505398, -494303, -483269, -472294, -461375, -450511, -439700, -428940,
    -418229, -407566, -396949, -386377, -375848, -365360, -354913, -344504, -334132, -323796,
    -313494, -303226, -292989, -282783, -272607, -262459, -252337, -242242, -232171, -222123,
    -212098, -202094, -192111, -182146, -172199, -162270, -152357, -142458, -132573, -122702,
    -112842, -102993, -93154, -83325, -73503, -63688, -53880, -44076, -34277, -24481, -14688,
    -4896,
];

/// The parameters scenarios are drawn from, in the stored order of the assets. The program's own bounds keep
/// every intermediate within `i128`: at most 8 assets, volatilities and gaps of at most 1,000,000, and a positive
/// definite correlation matrix.
pub struct Parameters<'a> {
    /// The assets' symbols, distinct.
    pub symbols: &'a [[u8; 32]],
    /// Each asset's daily volatility, in millionths.
    pub volatilities: &'a [u32],
    /// Each pair's correlation in basis points, for each pair `(i, j)`, `i < j`, row by row.
    pub correlations: &'a [u16],
    /// Each asset's weekend gap, in millionths.
    pub gaps: &'a [u32],
    /// The position of the market asset, or `None` for the equal-weighted portfolio.
    pub market: Option<usize>,
}

/// Everything a row needs, computed once per set.
pub struct Set<'a> {
    parameters: &'a Parameters<'a>,
    size: usize,
    order: Vec<usize>,
    rank: Vec<usize>,
    factor: Vec<i128>,
    root: i128,
    moves: Vec<i128>,
}

/// Number of scenarios in a set of lattice size `size` over `n` assets: three draws of `size`, four shocks
/// for the market and every asset at once, and two for each asset alone.
pub fn count(size: usize, n: usize) -> usize {
    3 * size + 4 + 2 * n
}

fn generator(size: usize) -> usize {
    match size {
        32 => 5,
        64 => 11,
        _ => 13,
    }
}

fn mean(size: usize, j: usize) -> i128 {
    let half = size / 2;
    let table: &[i32] = match size {
        32 => &MEANS_32,
        64 => &MEANS_64,
        128 => &MEANS_128,
        _ => &MEANS_256,
    };
    if j < half {
        table[j].into()
    } else {
        -i128::from(table[size - 1 - j])
    }
}

/// Cranley–Patterson shift of dimension `d`: ⌊size · frac(d · φ)⌋, φ = (√5 − 1)/2.
fn shift(size: usize, d: usize) -> usize {
    ((size as u128 * ((d as u128 * 618_033_988_749_894_848) % 1_000_000_000_000_000_000))
        / 1_000_000_000_000_000_000) as usize
}

/// Standard normal draw of point `j` in dimension `d` of the shifted Korobov lattice, in millionths.
fn draw(size: usize, j: usize, d: usize) -> i128 {
    let mut g = 1;
    for _ in 0..d {
        g = g * generator(size) % size;
    }
    mean(size, (j * g + shift(size, d)) % size)
}

/// Integer square root, rounded down.
pub fn isqrt(n: u128) -> u128 {
    if n == 0 {
        return 0;
    }
    let mut x = n;
    let mut y = x.div_ceil(2);
    while y < x {
        x = y;
        y = (x + n / x) / 2;
    }
    x
}

impl<'a> Set<'a> {
    /// Prepares the set for lattice size `size` (one of [`SIZES`]) over a horizon of `horizon` seconds.
    pub fn new(parameters: &'a Parameters<'a>, size: usize, horizon: u64) -> Self {
        debug_assert!(SIZES.contains(&size), "unsupported lattice size {size}");
        let n = parameters.symbols.len();
        let mut order: Vec<usize> = (0..n).collect();
        for i in 1..n {
            let mut k = i;
            while k > 0 && parameters.symbols[order[k - 1]] > parameters.symbols[order[k]] {
                order.swap(k - 1, k);
                k -= 1;
            }
        }
        let mut rank = vec![0; n];
        for (r, &i) in order.iter().enumerate() {
            rank[i] = r;
        }
        let rho = |a: usize, b: usize| -> i128 {
            if a == b {
                ONE.into()
            } else {
                parameters.correlations[pair(n, a.min(b), a.max(b))].into()
            }
        };
        let mut factor = vec![0i128; n * n];
        for j in 0..n {
            let acc = S * S
                - (0..j)
                    .map(|k| factor[j * n + k] * factor[j * n + k])
                    .sum::<i128>();
            let pivot = isqrt(acc.max(0) as u128) as i128;
            factor[j * n + j] = if pivot >= PIVOT_FLOOR { pivot } else { 0 };
            for i in j + 1..n {
                let acc = rho(order[i], order[j]) * (S / i128::from(ONE)) * S
                    - (0..j)
                        .map(|k| factor[i * n + k] * factor[j * n + k])
                        .sum::<i128>();
                factor[i * n + j] = if factor[j * n + j] > 0 {
                    (acc / factor[j * n + j]).clamp(-S, S)
                } else {
                    0
                };
            }
        }
        let root = isqrt(u128::from(horizon) * 1_000_000_000_000 / DAY) as i128;
        let weight = |k: usize| i128::from(parameters.market.is_none_or(|m| m == k));
        let vol = |k: usize| i128::from(parameters.volatilities[k]);
        let variance: i128 = (0..n)
            .flat_map(|a| (0..n).map(move |b| (a, b)))
            .map(|(a, b)| weight(a) * weight(b) * vol(a) * vol(b) * rho(a, b))
            .sum();
        let sigma = isqrt((variance * i128::from(ONE)) as u128) as i128;
        let moves = (0..n)
            .map(|i| {
                let covariance: i128 = (0..n)
                    .map(|k| weight(k) * vol(i) * vol(k) * rho(i, k))
                    .sum();
                Z_MARKET * root * covariance / (1_000_000_000_000 * sigma)
            })
            .collect();
        Set {
            parameters,
            size,
            order,
            rank,
            factor,
            root,
            moves,
        }
    }

    /// Scenario `index` of [`count`], as returns in millionths in the stored order of the assets.
    pub fn row(&self, index: usize) -> Vec<i32> {
        let p = self.parameters;
        let n = p.symbols.len();
        let size = self.size;
        let scaled =
            |x: i128, i: usize| clamp(x * i128::from(p.volatilities[i]) * self.root / (PPM * PPM));
        let mut row = vec![0; n];
        if index < size {
            let draws: Vec<i128> = (0..n).map(|d| draw(size, index, d)).collect();
            for (c, &i) in self.order.iter().enumerate() {
                let x: i128 = (0..=c)
                    .map(|k| self.factor[c * n + k] * draws[k])
                    .sum::<i128>()
                    / S;
                row[i] = scaled(x, i);
            }
        } else if index < 2 * size {
            let x = draw(size, index - size, 0);
            for (i, r) in row.iter_mut().enumerate() {
                *r = scaled(x, i);
            }
        } else if index < 3 * size {
            for (i, r) in row.iter_mut().enumerate() {
                *r = scaled(draw(size, index - 2 * size, self.rank[i]), i);
            }
        } else if index < 3 * size + 4 {
            let shock = index - 3 * size;
            for (i, r) in row.iter_mut().enumerate() {
                let magnitude = if shock < 2 {
                    self.moves[i]
                } else {
                    p.gaps[i].into()
                };
                *r = clamp(if shock.is_multiple_of(2) {
                    -magnitude
                } else {
                    magnitude
                });
            }
        } else {
            let alone = index - 3 * size - 4;
            let i = self.order[alone / 2];
            let gap = i128::from(p.gaps[i]);
            row[i] = clamp(if alone.is_multiple_of(2) { -gap } else { gap });
        }
        row
    }

    /// The lattice size of the set.
    pub fn size(&self) -> usize {
        self.size
    }

    /// The parameters the set is drawn from.
    pub fn parameters(&self) -> &Parameters<'a> {
        self.parameters
    }

    /// Every scenario in the canonical order of the assets, ascending by symbol, as big-endian int32:
    /// the same bytes for any order the assets were stored in.
    pub fn encoded(&self) -> Vec<u8> {
        let n = self.order.len();
        let mut out = Vec::with_capacity(count(self.size, n) * n * 4);
        for index in 0..count(self.size, n) {
            let row = self.row(index);
            for &i in &self.order {
                out.extend_from_slice(&row[i].to_be_bytes());
            }
        }
        out
    }
}

/// No return below −100%, and none beyond the int32 range.
fn clamp(r: i128) -> i32 {
    r.clamp(-PPM, i32::MAX.into()) as i32
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;
    use stylus_sdk::alloy_primitives::keccak256;

    #[derive(Debug)]
    struct Config {
        symbols: Vec<[u8; 32]>,
        volatilities: Vec<u32>,
        correlations: Vec<u16>,
        gaps: Vec<u32>,
        market: Option<usize>,
    }

    impl Config {
        fn parameters(&self) -> Parameters<'_> {
            Parameters {
                symbols: &self.symbols,
                volatilities: &self.volatilities,
                correlations: &self.correlations,
                gaps: &self.gaps,
                market: self.market,
            }
        }
    }

    fn symbol(name: &str) -> [u8; 32] {
        let mut word = [0u8; 32];
        word[..name.len()].copy_from_slice(name.as_bytes());
        word
    }

    fn numbers<T: TryFrom<u64>>(value: &serde_json::Value) -> Vec<T> {
        value
            .as_array()
            .unwrap()
            .iter()
            .map(|v| T::try_from(v.as_u64().unwrap()).ok().unwrap())
            .collect()
    }

    fn vectors() -> serde_json::Value {
        serde_json::from_str(include_str!("../testdata/scenario-vectors.json")).unwrap()
    }

    fn config(c: &serde_json::Value) -> Config {
        let names: Vec<&str> = c["names"]
            .as_array()
            .unwrap()
            .iter()
            .map(|n| n.as_str().unwrap())
            .collect();
        Config {
            symbols: names.iter().map(|n| symbol(n)).collect(),
            volatilities: numbers(&c["volatilities"]),
            correlations: numbers(&c["correlations"]),
            gaps: numbers(&c["gaps"]),
            market: c["market"]
                .as_str()
                .map(|m| names.iter().position(|n| *n == m).unwrap()),
        }
    }

    #[test]
    fn every_set_is_the_reference_set() {
        let v = vectors();
        let digests = v["digests"].as_array().unwrap();
        assert_eq!(digests.len(), 144);
        for d in digests {
            let c = config(&v["configs"][d["config"].as_str().unwrap()]);
            let p = c.parameters();
            let set = Set::new(
                &p,
                d["size"].as_u64().unwrap() as usize,
                d["horizon"].as_u64().unwrap(),
            );
            assert_eq!(
                keccak256(set.encoded()).to_string(),
                d["digest"].as_str().unwrap(),
                "{d}"
            );
        }
        let c = config(&v["configs"]["launch"]);
        let p = c.parameters();
        let set = Set::new(&p, 256, 172_800);
        for row in v["launchRows"]["rows"].as_array().unwrap() {
            let expected: Vec<i32> = row["returns"]
                .as_array()
                .unwrap()
                .iter()
                .map(|r| r.as_i64().unwrap() as i32)
                .collect();
            assert_eq!(set.row(row["index"].as_u64().unwrap() as usize), expected);
        }
    }

    #[test]
    fn each_dimension_of_the_lattice_takes_every_stratum_once() {
        for size in SIZES {
            for d in 0..8 {
                let mut draws: Vec<i128> = (0..size).map(|j| draw(size, j, d)).collect();
                draws.sort();
                let mut means: Vec<i128> = (0..size).map(|j| mean(size, j)).collect();
                means.sort();
                assert_eq!(draws, means, "size {size}, dimension {d}");
                assert_eq!(means.iter().sum::<i128>(), 0);
            }
        }
    }

    #[test]
    #[should_panic]
    fn a_lattice_size_outside_the_supported_ones_is_refused() {
        let c = Config {
            symbols: vec![symbol("SPY")],
            volatilities: vec![11_335],
            correlations: vec![],
            gaps: vec![54_840],
            market: Some(0),
        };
        Set::new(&c.parameters(), 100, 172_800);
    }

    #[test]
    fn a_single_asset_draws_the_same_scenarios_in_every_set() {
        let c = Config {
            symbols: vec![symbol("SPY")],
            volatilities: vec![11_335],
            correlations: vec![],
            gaps: vec![54_840],
            market: Some(0),
        };
        let p = c.parameters();
        let set = Set::new(&p, 64, 172_800);
        for j in 0..64 {
            assert_eq!(set.row(j), set.row(64 + j));
            assert_eq!(set.row(j), set.row(128 + j));
        }
    }

    fn strategy() -> impl Strategy<Value = Config> {
        (1usize..=8)
            .prop_flat_map(|n| {
                (
                    Just(n),
                    prop::collection::vec(1u32..=1_000_000, n),
                    prop::collection::vec(1u16..=9_000, n * (n - 1) / 2),
                    prop::collection::vec(1u32..=1_000_000, n),
                    prop::option::of(0..n),
                    prop::collection::btree_set(1u32..=u32::MAX, n),
                )
            })
            .prop_map(
                |(_, volatilities, correlations, gaps, market, ids)| Config {
                    symbols: ids
                        .into_iter()
                        .map(|id| {
                            let mut word = [0u8; 32];
                            word[..4].copy_from_slice(&id.to_be_bytes());
                            word
                        })
                        .collect(),
                    volatilities,
                    correlations,
                    gaps,
                    market,
                },
            )
    }

    fn reordered(c: &Config, order: &[usize]) -> Config {
        let n = c.symbols.len();
        let rho = |a: usize, b: usize| c.correlations[pair(n, a.min(b), a.max(b))];
        Config {
            symbols: order.iter().map(|&i| c.symbols[i]).collect(),
            volatilities: order.iter().map(|&i| c.volatilities[i]).collect(),
            correlations: (0..n)
                .flat_map(|i| (i + 1..n).map(move |j| (i, j)))
                .map(|(i, j)| rho(order[i], order[j]))
                .collect(),
            gaps: order.iter().map(|&i| c.gaps[i]).collect(),
            market: c
                .market
                .map(|m| order.iter().position(|&i| i == m).unwrap()),
        }
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, cases: 64, ..ProptestConfig::default() })]

        #[test]
        fn reordering_the_assets_only_permutes_each_scenario(
            c in strategy(),
            shuffle in Just((0..8).collect::<Vec<usize>>()).prop_shuffle(),
            horizon in 0u64..=864_000,
        ) {
            let n = c.symbols.len();
            let order: Vec<usize> = shuffle.into_iter().filter(|&i| i < n).collect();
            let r = reordered(&c, &order);
            let (p, q) = (c.parameters(), r.parameters());
            let (a, b) = (Set::new(&p, 32, horizon), Set::new(&q, 32, horizon));
            for index in 0..count(32, n) {
                let (x, y) = (a.row(index), b.row(index));
                for (k, &i) in order.iter().enumerate() {
                    prop_assert_eq!(x[i], y[k]);
                }
            }
            prop_assert_eq!(a.encoded(), b.encoded());
        }

        #[test]
        fn the_market_moves_itself_by_its_expected_shortfall(c in strategy(), horizon in 0u64..=864_000) {
            if let Some(m) = c.market {
                let p = c.parameters();
                let set = Set::new(&p, 32, horizon);
                let expected = Z_MARKET * set.root * i128::from(c.volatilities[m]) / (PPM * PPM);
                prop_assert_eq!(i128::from(set.row(3 * 32 + 1)[m]), expected.min(i32::MAX.into()));
                prop_assert_eq!(i128::from(set.row(3 * 32)[m]), (-expected).max(-PPM));
            }
        }

        #[test]
        fn four_days_move_twice_as_far_as_one(c in strategy()) {
            let p = c.parameters();
            let (day, four) = (Set::new(&p, 32, 86_400), Set::new(&p, 32, 345_600));
            for index in 0..3 * 32 {
                for (x, y) in day.row(index).into_iter().zip(four.row(index)) {
                    if y > -1_000_000 {
                        prop_assert!((i64::from(y) - 2 * i64::from(x)).abs() <= 2);
                    }
                }
            }
        }

        #[test]
        fn no_return_falls_below_minus_one(c in strategy(), horizon in any::<u64>()) {
            let p = c.parameters();
            let set = Set::new(&p, 32, horizon);
            for index in 0..count(32, c.symbols.len()) {
                prop_assert!(set.row(index).into_iter().all(|r| r >= -1_000_000));
            }
        }

        #[test]
        fn each_asset_gaps_alone_in_symbol_order(c in strategy()) {
            let n = c.symbols.len();
            let p = c.parameters();
            let set = Set::new(&p, 32, 86_400);
            let mut order: Vec<usize> = (0..n).collect();
            order.sort_by_key(|&i| c.symbols[i]);
            for (r, &i) in order.iter().enumerate() {
                let (down, up) = (set.row(3 * 32 + 4 + 2 * r), set.row(3 * 32 + 5 + 2 * r));
                for k in 0..n {
                    let gap = if k == i { c.gaps[i] as i32 } else { 0 };
                    prop_assert_eq!((down[k], up[k]), (-gap.min(1_000_000), gap));
                }
            }
        }

        #[test]
        fn the_square_root_is_exact(n in any::<u128>()) {
            let r = isqrt(n);
            prop_assert!(r * r <= n);
            prop_assert!((r + 1).checked_mul(r + 1).is_none_or(|s| s > n));
        }
    }
}
