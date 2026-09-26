//! ERC-8056 multipliers. A Stock Token's price is its share's price times the token's multiplier, and
//! a multiplier change is a step at `effectiveAt`. The token does not expose the multiplier it replaced,
//! so the program records each change it sees.

use crate::chainlink::CL_DEV_BPS;

/// Scale of a multiplier: `SCALE` is 1.
pub const SCALE: u128 = 1_000_000_000_000_000_000;

/// A multiplier change as the program recorded it. All zero when none was seen.
#[derive(Clone, Copy, Debug, Default, PartialEq, Eq)]
pub struct Change {
    /// The multiplier before the step; zero when the step was seen only after it happened.
    pub before: u128,
    /// The multiplier after the step.
    pub after: u128,
    /// When the step happens, in seconds.
    pub effective_at: u64,
    /// Whether Chainlink has confirmed the new terms.
    pub confirmed: bool,
}

/// What a change means for the band now.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
#[repr(u8)]
pub enum Status {
    /// No change needs a window.
    None = 0,
    /// A material change is scheduled.
    Scheduled = 1,
    /// A material change happened and Chainlink has not confirmed the new terms: the band is halted.
    Unconfirmed = 2,
}

/// `floor(a × b / c)`, or `None` when `c` is zero or the result does not fit in `u64`.
pub fn mul_div(a: u64, b: u128, c: u128) -> Option<u64> {
    if c == 0 {
        return None;
    }
    let (b_hi, b_lo) = (b >> 64, b & u128::from(u64::MAX));
    let (x, y) = (u128::from(a) * b_lo, u128::from(a) * b_hi);
    let (lo, carry) = x.overflowing_add(y << 64);
    let hi = (y >> 64) + u128::from(carry);
    if hi >= c {
        return None;
    }
    let mut rem = hi;
    let mut quotient = 0u128;
    for bit in (0..128).rev() {
        let top = rem >> 127;
        rem = (rem << 1) | ((lo >> bit) & 1);
        quotient <<= 1;
        if top == 1 || rem >= c {
            rem = rem.wrapping_sub(c);
            quotient |= 1;
        }
    }
    u64::try_from(quotient).ok()
}

/// The token price for a share price, rounded down as the token does. Zero when it does not fit in
/// `u64`.
pub fn token_px(share_px: u64, multiplier: u128) -> u64 {
    mul_div(share_px, multiplier, SCALE).unwrap_or(0)
}

/// Whether a change needs a window: its size is not known, or it moves the price by at least
/// Chainlink's deviation threshold.
pub fn material(change: Change) -> bool {
    change.before == 0
        || change
            .after
            .abs_diff(change.before)
            .saturating_mul(u128::from(10_000 / CL_DEV_BPS))
            >= change.before
}

/// A Chainlink answer observed at `started_at` in the terms in force at `now`. A round observed before
/// the step is scaled by the change's ratio when the change is not material, because a reinvested
/// dividend leaves the share price where it was; after a material change, which may be a split that
/// also moves the share price, it is zero.
pub fn chainlink_px(answer: u64, started_at: u64, change: Change, now: u64) -> u64 {
    if change.effective_at == 0 || now < change.effective_at || started_at >= change.effective_at {
        return answer;
    }
    if material(change) {
        return 0;
    }
    mul_div(answer, change.after, change.before).unwrap_or(0)
}

/// The change's status at `now`.
pub fn status(change: Change, now: u64) -> Status {
    if change.effective_at == 0 || change.confirmed || !material(change) {
        Status::None
    } else if now < change.effective_at {
        Status::Scheduled
    } else {
        Status::Unconfirmed
    }
}

fn in_force(change: Change, now: u64) -> u128 {
    if change.effective_at == 0 || now >= change.effective_at {
        change.after
    } else {
        change.before
    }
}

/// The change in force given the token's `current` multiplier and its `new` one at `effective_at`.
/// A material change past its step and not confirmed holds until it is confirmed. Otherwise it is the
/// recorded change when the token reports it; a change of unknown size when its step has passed
/// unrecorded; a missed step of unknown size, counted from `now`, when the token's current multiplier
/// is not the one the recorded change implies; and a scheduled change of known size otherwise.
pub fn observe(recorded: Change, current: u128, new: u128, effective_at: u64, now: u64) -> Change {
    if status(recorded, now) == Status::Unconfirmed
        || (effective_at == recorded.effective_at && new == recorded.after)
    {
        return recorded;
    }
    let (before, after, effective_at) = if effective_at != 0 && effective_at <= now {
        (0, new, effective_at)
    } else if recorded.after != 0 && current != in_force(recorded, now) {
        (0, current, now)
    } else {
        (current, new, effective_at)
    };
    Change {
        before,
        after,
        effective_at,
        confirmed: false,
    }
}

/// The change to record, or `None` when it is the recorded one.
pub fn record(
    recorded: Change,
    current: u128,
    new: u128,
    effective_at: u64,
    now: u64,
) -> Option<Change> {
    let change = observe(recorded, current, new, effective_at, now);
    (change != recorded).then_some(change)
}

/// Whether a Chainlink round that started at `started_at` with price `cl_px` confirms the change's new
/// terms: it started at or after the step, lies inside the band `[low, high]` that the 24/7 leg draws
/// alone around `px_247`, and, when the change's size is known, is closer to `px_247` than to the
/// 24/7 price in the old terms.
pub fn confirms(
    change: Change,
    started_at: u64,
    cl_px: u64,
    px_247: u64,
    low: u64,
    high: u128,
) -> bool {
    if started_at < change.effective_at || cl_px == 0 || cl_px < low || u128::from(cl_px) > high {
        return false;
    }
    if change.before == 0 {
        return true;
    }
    let old = mul_div(px_247, change.before, change.after).unwrap_or(0);
    cl_px.abs_diff(px_247) < cl_px.abs_diff(old)
}

/// The token's current multiplier without asking the token, from its `new` multiplier and
/// `effective_at` and the recorded change; `None` when it has to be read.
pub fn current(change: Change, new: u128, effective_at: u64, now: u64) -> Option<u128> {
    if effective_at == 0 {
        return Some(new);
    }
    if change.effective_at != effective_at || change.after != new {
        return None;
    }
    if now >= effective_at {
        Some(change.after)
    } else {
        (change.before != 0).then_some(change.before)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;
    use serde_json::Value;
    use stylus_sdk::alloy_primitives::U256;

    const VECTORS: &str = include_str!("../testdata/multiplier-vectors.json");

    fn int(value: &Value) -> u128 {
        value.as_str().unwrap().parse().unwrap()
    }

    fn change(value: &Value) -> Change {
        Change {
            before: int(&value["before"]),
            after: int(&value["after"]),
            effective_at: int(&value["effective_at"]) as u64,
            confirmed: value["confirmed"].as_bool().unwrap(),
        }
    }

    fn cases(section: &str) -> Vec<Value> {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        vectors[section].as_array().unwrap().clone()
    }

    #[test]
    fn mul_div_matches_the_reference() {
        for case in cases("mul_div") {
            let expected = case["expected"].as_str().map(|v| v.parse::<u64>().unwrap());
            let got = mul_div(int(&case["a"]) as u64, int(&case["b"]), int(&case["c"]));
            assert_eq!(got, expected, "{}", case["name"]);
        }
    }

    #[test]
    fn materiality_matches_the_reference() {
        for case in cases("material") {
            let expected = case["expected"].as_bool().unwrap();
            assert_eq!(
                material(change(&case["change"])),
                expected,
                "{}",
                case["name"]
            );
        }
    }

    #[test]
    fn chainlink_prices_match_the_reference() {
        for case in cases("chainlink") {
            let got = chainlink_px(
                int(&case["answer"]) as u64,
                int(&case["started_at"]) as u64,
                change(&case["change"]),
                int(&case["now"]) as u64,
            );
            assert_eq!(u128::from(got), int(&case["expected"]), "{}", case["name"]);
        }
    }

    #[test]
    fn statuses_match_the_reference() {
        for case in cases("status") {
            let got = status(change(&case["change"]), int(&case["now"]) as u64);
            assert_eq!(got as u128, int(&case["expected"]), "{}", case["name"]);
        }
    }

    #[test]
    fn observed_changes_match_the_reference() {
        for case in cases("observe") {
            let got = observe(
                change(&case["recorded"]),
                int(&case["current"]),
                int(&case["new"]),
                int(&case["effective_at"]) as u64,
                int(&case["now"]) as u64,
            );
            assert_eq!(got, change(&case["expected"]), "{}", case["name"]);
        }
    }

    #[test]
    fn records_match_the_reference() {
        for case in cases("record") {
            let got = record(
                change(&case["recorded"]),
                int(&case["current"]),
                int(&case["new"]),
                int(&case["effective_at"]) as u64,
                int(&case["now"]) as u64,
            );
            let expected = (!case["expected"].is_null()).then(|| change(&case["expected"]));
            assert_eq!(got, expected, "{}", case["name"]);
        }
    }

    #[test]
    fn current_multipliers_match_the_reference() {
        for case in cases("current") {
            let got = current(
                change(&case["change"]),
                int(&case["new"]),
                int(&case["effective_at"]) as u64,
                int(&case["now"]) as u64,
            );
            let expected = case["expected"]
                .as_str()
                .map(|v| v.parse::<u128>().unwrap());
            assert_eq!(got, expected, "{}", case["name"]);
        }
    }

    #[test]
    fn confirmations_match_the_reference() {
        for case in cases("confirms") {
            let got = confirms(
                change(&case["change"]),
                int(&case["started_at"]) as u64,
                int(&case["cl_px"]) as u64,
                int(&case["px_247"]) as u64,
                int(&case["low"]) as u64,
                int(&case["high"]),
            );
            assert_eq!(got, case["expected"].as_bool().unwrap(), "{}", case["name"]);
        }
    }

    proptest! {
        #![proptest_config(ProptestConfig { failure_persistence: None, ..ProptestConfig::default() })]

        #[test]
        fn mul_div_is_the_exact_floor(a in any::<u64>(), b in any::<u128>(), c in any::<u128>()) {
            let exact = (c != 0).then(|| U256::from(a) * U256::from(b) / U256::from(c));
            let expected = exact.and_then(|q| u64::try_from(q).ok());
            prop_assert_eq!(mul_div(a, b, c), expected);
        }

        #[test]
        fn a_material_change_never_scales_an_old_round(answer in any::<u64>(), before in 1..u128::MAX, after in any::<u128>()) {
            let change = Change { before, after, effective_at: 100, confirmed: false };
            let got = chainlink_px(answer, 99, change, 100);
            prop_assert!(!material(change) || got == 0);
        }

        #[test]
        fn an_unconfirmed_change_holds_whatever_the_token_reports(
            after in 1..u128::MAX, current in any::<u128>(), new in any::<u128>(), effective_at in any::<u64>(), late in 0..u64::MAX / 2,
        ) {
            let change = Change { before: 0, after, effective_at: 100, confirmed: false };
            let now = 100 + late;
            prop_assert_eq!(observe(change, current, new, effective_at, now), change);
            prop_assert_eq!(record(change, current, new, effective_at, now), None);
        }
    }
}
