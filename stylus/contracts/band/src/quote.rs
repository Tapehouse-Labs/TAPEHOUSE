//! The price band: centre, half-width, bounds and state from both legs, in integer arithmetic.
//!
//! Prices carry 8 decimals, widths are in basis points, and the variance is in centi-basis-points
//! squared per minute.

use crate::chainlink::{CL_DEV_BPS, CL_MAX_AGE_S};

/// Half-width floor, in basis points.
pub const FLOOR_BPS: u128 = 30;
/// Volatility multiplier z, times 10.
pub const Z_X10: u128 = 30;
/// Minutes a price may age before it is used on-chain.
pub const LATENCY_MIN: u128 = 3;
/// Surcharge when only one leg is live, in basis points.
pub const SINGLE_SOURCE_BPS: u128 = 25;
/// Maximum age of a live 24/7 leg, in seconds.
pub const LIVE_MAX_AGE_S: u64 = 120;
/// Half-width cap, in basis points.
pub const MAX_HALF_BPS: u128 = 1500;
/// EWMA lambda numerator, per one-minute sample.
pub const EWMA_NUM: u128 = 94;
/// EWMA lambda denominator.
pub const EWMA_DEN: u128 = 100;
/// Cap on one sample's return, in centi-basis-points, so the variance fits in `u128`.
pub const R_CAP_CPB: u128 = 1 << 57;
/// Minimum package-time gap between two variance samples, in milliseconds.
pub const SAMPLE_MIN_GAP_MS: u64 = 50_000;

const SQRT_LATENCY_X100: u128 = (LATENCY_MIN * 10_000).isqrt();

/// What the band reports about its legs, from the most to the least restrictive.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
#[repr(u8)]
pub enum State {
    /// No leg is live.
    Halted = 0,
    /// Fewer legs are live than the session expects, or the session is not known.
    Degraded = 1,
    /// Chainlink's session is closed and the 24/7 leg is live.
    Closed = 2,
    /// Chainlink's session is open and both legs are live.
    Open = 3,
}

/// Both legs of one asset, the Chainlink session and the variance of the 24/7 leg.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Inputs {
    /// The 24/7 price, zero when the leg is absent.
    pub live247_px: u64,
    /// Seconds since the 24/7 price's package.
    pub live247_age_s: u64,
    /// The Chainlink answer, zero when the leg is absent.
    pub cl_px: u64,
    /// Seconds since the Chainlink answer's `updatedAt`.
    pub cl_age_s: u64,
    /// Whether Chainlink's 24/5 session is open.
    pub cl_session_open: bool,
    /// The variance of the 24/7 leg, below 2^120.
    pub var_cpb2: u128,
    /// Extra half-width while the 24/7 leg is live, in basis points: non-zero when it is made from an
    /// index.
    pub basis_bps: u64,
}

/// The band: its state, how many legs are live, its centre, half-width and bounds.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Quote {
    /// What the band reports about its legs.
    pub state: State,
    /// How many legs are live.
    pub live: u8,
    /// The centre.
    pub mid: u64,
    /// The half-width, in basis points.
    pub half_bps: u64,
    /// The lower bound, rounded down.
    pub low: u64,
    /// The upper bound, rounded down. It can exceed `u64` when the centre is near `u64::MAX`.
    pub high: u128,
}

/// Computes the band. The centre is the live 24/7 leg, else Chainlink, and never an average with a
/// sleeping leg.
pub fn compute(inputs: &Inputs) -> Quote {
    let live247 = inputs.live247_px > 0 && inputs.live247_age_s <= LIVE_MAX_AGE_S;
    let cl_live = inputs.cl_px > 0 && inputs.cl_session_open && inputs.cl_age_s <= CL_MAX_AGE_S;
    if !live247 && !cl_live {
        return Quote {
            state: State::Halted,
            live: 0,
            mid: 0,
            half_bps: 0,
            low: 0,
            high: 0,
        };
    }
    let mid = if live247 {
        inputs.live247_px
    } else {
        inputs.cl_px
    };

    let vol_bps = Z_X10 * inputs.var_cpb2.isqrt() * SQRT_LATENCY_X100 / (10 * 100 * 100);
    let mut half = FLOOR_BPS.max(vol_bps);
    let live = u8::from(live247) + u8::from(cl_live);
    if live == 2 {
        let disagreement = u128::from(inputs.cl_px.abs_diff(inputs.live247_px)) * 10_000
            / u128::from(inputs.live247_px);
        half += disagreement.saturating_sub(u128::from(CL_DEV_BPS));
    } else {
        half += SINGLE_SOURCE_BPS;
    }
    if live247 {
        half += u128::from(inputs.basis_bps);
    } else {
        half += u128::from(CL_DEV_BPS);
    }
    let half = half.min(MAX_HALF_BPS);

    let expected = if inputs.cl_session_open { 2 } else { 1 };
    let state = match (live < expected, inputs.cl_session_open) {
        (true, _) => State::Degraded,
        (false, true) => State::Open,
        (false, false) => State::Closed,
    };
    let centre = u128::from(mid);
    Quote {
        state,
        live,
        mid,
        half_bps: half as u64,
        low: (centre * (10_000 - half) / 10_000) as u64,
        high: centre * (10_000 + half) / 10_000,
    }
}

/// Folds one sample into the variance: the return from `prev_px` to `px`, normalised to one minute
/// over `dt_s` seconds. A fall rounds toward minus infinity, as floor division of the signed return
/// does. Requires `var_cpb2 < 2^120`; the capped return keeps every result below that bound.
pub fn ewma_update(var_cpb2: u128, prev_px: u64, px: u64, dt_s: u64) -> u128 {
    if prev_px == 0 || dt_s == 0 {
        return var_cpb2;
    }
    let move_cpb = u128::from(px.abs_diff(prev_px)) * 1_000_000;
    let r_cpb = if px >= prev_px {
        move_cpb / u128::from(prev_px)
    } else {
        move_cpb.div_ceil(u128::from(prev_px))
    }
    .min(R_CAP_CPB);
    let r2_per_min = r_cpb * r_cpb * 60 / u128::from(dt_s);
    (EWMA_NUM * var_cpb2 + (EWMA_DEN - EWMA_NUM) * r2_per_min) / EWMA_DEN
}

/// The variance, the sampled price and its package timestamp in milliseconds.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Sample {
    /// The variance after this sample.
    pub var_cpb2: u128,
    /// The sampled price.
    pub px: u64,
    /// The sampled price's package timestamp, in milliseconds; zero before the first sample.
    pub ms: u64,
}

/// Whether a price with package timestamp `ms` is due as a sample after the one taken at `last_ms`:
/// always for the first price, then once [`SAMPLE_MIN_GAP_MS`] has passed.
pub fn sample_due(last_ms: u64, ms: u64) -> bool {
    last_ms == 0 || ms >= last_ms.saturating_add(SAMPLE_MIN_GAP_MS)
}

/// The next sample once a price `px` with package timestamp `ms` arrives, or `None` when the price
/// is zero or no sample is due.
pub fn next_sample(last: Sample, px: u64, ms: u64) -> Option<Sample> {
    if px == 0 || !sample_due(last.ms, ms) {
        return None;
    }
    let var_cpb2 = ewma_update(last.var_cpb2, last.px, px, (ms - last.ms) / 1000);
    Some(Sample { var_cpb2, px, ms })
}

/// Seconds from `at_s` to `now_s`, zero when `at_s` is ahead.
pub fn age_s(now_s: u64, at_s: u64) -> u64 {
    now_s.saturating_sub(at_s)
}

/// Whole seconds from a package timestamp in milliseconds to `now_s`, zero when the package is ahead.
pub fn package_age_s(now_s: u64, ms: u64) -> u64 {
    now_s.saturating_mul(1000).saturating_sub(ms) / 1000
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;
    use serde_json::Value;

    const VECTORS: &str = include_str!("../testdata/quote-vectors.json");

    fn int(value: &Value) -> u128 {
        value.as_str().unwrap().parse().unwrap()
    }

    fn state(name: &str) -> State {
        match name {
            "OPEN" => State::Open,
            "CLOSED" => State::Closed,
            "DEGRADED" => State::Degraded,
            _ => State::Halted,
        }
    }

    #[test]
    fn quotes_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["quotes"].as_array().unwrap() {
            let (input, expected) = (&case["input"], &case["expected"]);
            let inputs = Inputs {
                live247_px: int(&input["live247_px"]) as u64,
                live247_age_s: int(&input["live247_age_s"]) as u64,
                cl_px: int(&input["cl_px"]) as u64,
                cl_age_s: int(&input["cl_age_s"]) as u64,
                cl_session_open: input["cl_session_open"].as_bool().unwrap(),
                var_cpb2: int(&input["var_cpb2"]),
                basis_bps: input.get("basis_bps").map_or(0, |v| int(v) as u64),
            };
            let quote = Quote {
                state: state(expected["state"].as_str().unwrap()),
                live: int(&expected["live"]) as u8,
                mid: int(&expected["mid"]) as u64,
                half_bps: int(&expected["half_bps"]) as u64,
                low: int(&expected["low"]) as u64,
                high: int(&expected["high"]),
            };
            assert_eq!(compute(&inputs), quote, "{}", case["name"]);
        }
    }

    #[test]
    fn variance_updates_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["ewma"].as_array().unwrap() {
            let var = ewma_update(
                int(&case["var"]),
                int(&case["prev_px"]) as u64,
                int(&case["px"]) as u64,
                int(&case["dt_s"]) as u64,
            );
            assert_eq!(var, int(&case["expected"]), "{}", case["name"]);
        }
    }

    #[test]
    fn the_first_price_is_the_first_sample() {
        let first = Sample {
            var_cpb2: 0,
            px: 0,
            ms: 0,
        };
        assert_eq!(
            next_sample(first, 22_000_000_000, 1_790_000_000_000),
            Some(Sample {
                var_cpb2: 0,
                px: 22_000_000_000,
                ms: 1_790_000_000_000
            })
        );
    }

    #[test]
    fn a_price_within_the_gap_is_not_a_sample() {
        let last = Sample {
            var_cpb2: 250_000,
            px: 22_000_000_000,
            ms: 1_790_000_000_000,
        };
        assert_eq!(next_sample(last, 22_300_000_000, 1_790_000_049_999), None);
    }

    #[test]
    fn a_price_after_the_gap_folds_into_the_variance() {
        let last = Sample {
            var_cpb2: 250_000,
            px: 22_000_000_000,
            ms: 1_790_000_000_000,
        };
        assert_eq!(
            next_sample(last, 22_300_000_000, 1_790_000_050_999),
            Some(Sample {
                var_cpb2: ewma_update(250_000, 22_000_000_000, 22_300_000_000, 50),
                px: 22_300_000_000,
                ms: 1_790_000_050_999,
            })
        );
    }

    #[test]
    fn a_long_gap_keeps_the_variance() {
        let last = Sample {
            var_cpb2: 250_000,
            px: 22_000_000_000,
            ms: 1_790_000_000_000,
        };
        let next = next_sample(last, 22_000_000_000, 1_790_003_600_000).unwrap();
        assert_eq!(
            next.var_cpb2,
            ewma_update(250_000, 22_000_000_000, 22_000_000_000, 3600)
        );
        assert!(next.var_cpb2 > 0);
    }

    #[test]
    fn a_zero_price_is_not_a_sample() {
        let last = Sample {
            var_cpb2: 250_000,
            px: 22_000_000_000,
            ms: 1_790_000_000_000,
        };
        assert_eq!(next_sample(last, 0, 1_790_000_060_000), None);
    }

    #[test]
    fn a_package_ahead_of_the_block_has_age_zero() {
        assert_eq!(package_age_s(1_790_000_000, 1_790_000_030_000), 0);
        assert_eq!(package_age_s(1_790_000_000, 1_789_999_879_500), 120);
        assert_eq!(package_age_s(1_790_000_000, 1_789_999_879_499), 120);
        assert_eq!(package_age_s(1_790_000_000, 1_789_999_879_000), 121);
        assert_eq!(age_s(1_790_000_000, 1_790_000_005), 0);
    }

    fn inputs() -> impl Strategy<Value = Inputs> {
        (
            any::<u64>(),
            0..300u64,
            any::<u64>(),
            0..200_000u64,
            any::<bool>(),
            any::<u128>(),
            0..2_000u64,
        )
            .prop_map(
                |(
                    live247_px,
                    live247_age_s,
                    cl_px,
                    cl_age_s,
                    cl_session_open,
                    var_cpb2,
                    basis_bps,
                )| {
                    Inputs {
                        live247_px,
                        live247_age_s,
                        cl_px,
                        cl_age_s,
                        cl_session_open,
                        var_cpb2,
                        basis_bps,
                    }
                },
            )
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, ..ProptestConfig::default() })]

        #[test]
        fn the_bounds_are_ordered(inputs in inputs()) {
            let quote = compute(&inputs);
            prop_assert!(u128::from(quote.low) <= u128::from(quote.mid));
            prop_assert!(u128::from(quote.mid) <= quote.high);
        }

        #[test]
        fn the_half_width_stays_between_floor_and_cap(inputs in inputs()) {
            let quote = compute(&inputs);
            if quote.state == State::Halted {
                prop_assert_eq!(quote, Quote { state: State::Halted, live: 0, mid: 0, half_bps: 0, low: 0, high: 0 });
            } else {
                prop_assert!(u128::from(quote.half_bps) >= FLOOR_BPS);
                prop_assert!(u128::from(quote.half_bps) <= MAX_HALF_BPS);
            }
        }

        #[test]
        fn a_live_24_7_leg_is_the_centre(inputs in inputs()) {
            let quote = compute(&inputs);
            if inputs.live247_px > 0 && inputs.live247_age_s <= LIVE_MAX_AGE_S {
                prop_assert_eq!(quote.mid, inputs.live247_px);
            }
        }

        #[test]
        fn a_sleeping_chainlink_leg_changes_nothing(inputs in inputs(), other_px in any::<u64>()) {
            let asleep = Inputs { cl_session_open: false, ..inputs };
            prop_assert_eq!(compute(&asleep), compute(&Inputs { cl_px: other_px, ..asleep }));
            let stale = Inputs { cl_age_s: CL_MAX_AGE_S + 1, ..inputs };
            prop_assert_eq!(compute(&stale), compute(&Inputs { cl_px: other_px, ..stale }));
        }

        #[test]
        fn extra_width_only_widens(inputs in inputs()) {
            let plain = compute(&Inputs { basis_bps: 0, ..inputs });
            let quote = compute(&inputs);
            prop_assert!(quote.half_bps >= plain.half_bps);
            prop_assert_eq!((quote.state, quote.live, quote.mid), (plain.state, plain.live, plain.mid));
        }

        #[test]
        fn the_variance_stays_below_its_bound(var in 0..(1u128 << 120), prev in any::<u64>(), px in any::<u64>(), dt in any::<u64>()) {
            prop_assert!(ewma_update(var, prev, px, dt) < 1u128 << 120);
        }
    }
}
