// SPDX-License-Identifier: MIT OR Apache-2.0
//! The band's 24/5 session, and the regime the current requirement follows through it.

use stylus_sdk::alloy_primitives::{Address, U256};
use stylus_sdk::prelude::*;

/// The margin period of risk, in seconds: two days in every regime.
pub const HORIZON: u64 = 172_800;
/// How long before a weekend or holiday session close the requirement starts rising to the closed one, in
/// milliseconds: the last three hours of the regular session and the four-hour post-market.
pub const RAMP_MS: u64 = 25_200_000;
/// The most gross exposure a requirement allows per unit of margin across a closure, in basis points: 5×.
pub const WEEKEND_LEVERAGE: u32 = 50_000;

sol_interface! {
    interface IBand {
        function session() external view returns (uint8, uint8, uint8, uint64, uint64);
    }
}

/// The band's session: `(open, nyse, nyse_next, change_ms, boundary_ms)`, as `session()` returns it.
pub type Session = (u8, u8, u8, u64, u64);

/// Where the current requirement stands.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Regime {
    /// The band cannot say whether the session is open.
    Unknown,
    /// The session is closed.
    Closed,
    /// The session is open and no weekend or holiday close is within [`RAMP_MS`].
    Open,
    /// The session closes for a weekend or holiday within [`RAMP_MS`], `elapsed_ms` into the ramp.
    Closing { elapsed_ms: u64 },
}

impl Regime {
    /// The code `currentRequirement` reports: 0 unknown, 1 closed, 2 open, 3 closing.
    pub fn code(self) -> u8 {
        match self {
            Self::Unknown => 0,
            Self::Closed => 1,
            Self::Open => 2,
            Self::Closing { .. } => 3,
        }
    }
}

/// `band`'s session. `None` when the call fails.
pub fn read(host: &impl Host, band: Address) -> Option<Session> {
    IBand::new(band).session(host, Call::new()).ok()
}

/// Whether a call that failed with `after` gas left, of `before`, failed for want of gas: a call keeps back
/// a 64th of the gas it may forward, so a caller that starves the band's call keeps no more than that.
pub fn starved(before: u64, after: u64) -> bool {
    u128::from(after) * 64 <= u128::from(before)
}

/// The regime of `session` at `now_ms`. The band sets an open session's boundary only on the last trading
/// day before a weekend or holiday close, from the regular open, to the end of that day's post-market; a
/// session still open at or past it reads as closed.
pub fn regime(session: Option<Session>, now_ms: u64) -> Regime {
    match session {
        Some((1, ..)) => Regime::Closed,
        Some((2, .., boundary_ms)) if boundary_ms != 0 && boundary_ms <= now_ms => Regime::Closed,
        Some((2, .., boundary_ms))
            if boundary_ms != 0 && boundary_ms.saturating_sub(now_ms) < RAMP_MS =>
        {
            Regime::Closing {
                elapsed_ms: RAMP_MS - boundary_ms.saturating_sub(now_ms),
            }
        }
        Some((2, ..)) => Regime::Open,
        _ => Regime::Unknown,
    }
}

/// The margin that holds `exposures`' gross value to [`WEEKEND_LEVERAGE`], rounded up.
pub fn leverage_floor(exposures: &[i128]) -> U256 {
    (exposures
        .iter()
        .map(|e| U256::from(e.unsigned_abs()))
        .fold(U256::ZERO, |sum, e| sum + e)
        * U256::from(10_000))
    .div_ceil(U256::from(WEEKEND_LEVERAGE))
}

/// The current requirement from the requirement over the open market, `open`, across a closure, `closed`,
/// and the [`leverage_floor`], `floor`. The open market holds a 25% buffer, rounded up. A closure needs the
/// largest of the buffered requirement, the closed one and the floor, so the buffer covers the weekend first;
/// the ramp before the close rises to it in a straight line, and a reopen only lowers it.
pub fn current(open: U256, closed: U256, floor: U256, regime: Regime) -> U256 {
    let buffered = (open * U256::from(5)).div_ceil(U256::from(4));
    let across = buffered.max(closed).max(floor);
    match regime {
        Regime::Open => buffered,
        Regime::Closing { elapsed_ms } => {
            buffered + (across - buffered) * U256::from(elapsed_ms) / U256::from(RAMP_MS)
        }
        Regime::Closed | Regime::Unknown => across,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;

    const BAND_SESSIONS: &str = include_str!("../../band/testdata/session-vectors.json");
    const TIMELINE: &str = include_str!("../testdata/session-timeline.json");
    const DOLLAR: u128 = 1_000_000_000_000_000_000;

    #[test]
    fn a_band_call_starved_of_gas_is_told_apart_from_a_band_that_fails() {
        let before = 1_000_000;
        assert!(starved(before, before / 64));
        assert!(starved(before, 0));
        assert!(!starved(before, before / 64 + 1));
        assert!(!starved(before, before - 2_600));
        assert!(!starved(2_700, 100));
    }

    fn dollars(d: u128) -> U256 {
        U256::from(d * DOLLAR)
    }

    #[test]
    fn every_session_the_band_reports_has_its_regime() {
        let v: serde_json::Value = serde_json::from_str(BAND_SESSIONS).unwrap();
        let expected = |name: &str| match name {
            "Friday, regular hours"
            | "Friday, the close passed on the last regular package"
            | "Friday, post-market after the close"
            | "Friday, last second of post-market"
            | "synthetic: daylight saving ends, post-market"
            | "synthetic: daylight saving starts, post-market"
            | "synthetic: early close, end of the late session"
            | "synthetic: evening before a holiday signed as a short close, post-market" => 3,
            "Saturday, the 24/5 close"
            | "Sunday, a second before the reopen"
            | "Friday post-market without the Friday close"
            | "weekend on a band that never saw a close"
            | "Labor Day, a second before the reopen"
            | "long close followed by a long close"
            | "synthetic: daylight saving ends, a second before the reopen"
            | "synthetic: daylight saving starts, the 24/5 close"
            | "synthetic: Thanksgiving, a second before the reopen"
            | "synthetic: early close, after the late session"
            | "synthetic: evening before a holiday signed as a short close, after post-market"
            | "synthetic: that short close passed into the holiday" => 1,
            "status an hour and a second old"
            | "points from different packages"
            | "change time from another package"
            | "unknown status value"
            | "change time zero"
            | "change time beyond u64" => 0,
            _ => 2,
        };
        let mut seen = [0; 4];
        for case in v["sessions"].as_array().unwrap() {
            let e = &case["expected"];
            let int = |x: &serde_json::Value| x.as_str().unwrap().parse::<u64>().unwrap();
            let code = if !e["known"].as_bool().unwrap() {
                0
            } else if e["open"].as_bool().unwrap() {
                2
            } else {
                1
            };
            let session = (
                code,
                int(&e["nyse"]) as u8,
                int(&e["nyse_next"]) as u8,
                int(&e["change_ms"]),
                int(&e["boundary_ms"]),
            );
            let name = case["name"].as_str().unwrap();
            let got = regime((code != 0).then_some(session), int(&case["now_ms"]));
            assert_eq!(got.code(), expected(name), "{name}");
            seen[usize::from(got.code())] += 1;
        }
        assert_eq!(seen, [6, 12, 12, 8]);
    }

    #[test]
    fn the_ramp_rises_in_a_straight_line_to_the_requirement_across_the_closure() {
        let (open, closed) = (dollars(1_000), dollars(1_400));
        let at = |elapsed_ms| current(open, closed, U256::ZERO, Regime::Closing { elapsed_ms });
        assert_eq!(
            current(open, closed, U256::ZERO, Regime::Open),
            dollars(1_250)
        );
        assert_eq!(at(0), dollars(1_250));
        assert_eq!(at(RAMP_MS / 4), dollars(1_287) + dollars(1) / U256::from(2));
        assert_eq!(at(RAMP_MS), dollars(1_400));
        assert_eq!(
            current(open, closed, U256::ZERO, Regime::Closed),
            dollars(1_400)
        );
        assert_eq!(
            current(open, closed, U256::ZERO, Regime::Unknown),
            dollars(1_400)
        );
        let covered = dollars(1_100);
        assert_eq!(
            current(
                open,
                covered,
                U256::ZERO,
                Regime::Closing {
                    elapsed_ms: RAMP_MS / 2
                }
            ),
            dollars(1_250)
        );
        assert_eq!(
            current(open, covered, U256::ZERO, Regime::Closed),
            dollars(1_250)
        );
    }

    #[test]
    fn across_a_closure_gross_exposure_stays_within_five_times_the_requirement() {
        let exposures = [
            4_000 * DOLLAR as i128,
            -2_000 * DOLLAR as i128,
            4_000 * DOLLAR as i128,
        ];
        let floor = leverage_floor(&exposures);
        assert_eq!(floor, dollars(2_000));
        let (open, closed) = (dollars(1_000), dollars(1_400));
        assert_eq!(current(open, closed, floor, Regime::Open), dollars(1_250));
        assert_eq!(
            current(
                open,
                closed,
                floor,
                Regime::Closing {
                    elapsed_ms: RAMP_MS / 2
                }
            ),
            dollars(1_625)
        );
        assert_eq!(current(open, closed, floor, Regime::Closed), dollars(2_000));
        assert_eq!(
            current(open, closed, floor, Regime::Unknown),
            dollars(2_000)
        );
        assert_eq!(
            current(open, dollars(2_500), floor, Regime::Closed),
            dollars(2_500)
        );
        assert_eq!(leverage_floor(&[]), U256::ZERO);
        assert_eq!(leverage_floor(&[-1, 5]), U256::from(2));
    }

    #[test]
    fn a_month_of_real_sessions_raises_a_requirement_only_through_the_ramp() {
        let v: serde_json::Value = serde_json::from_str(TIMELINE).unwrap();
        let step = v["step_ms"].as_u64().unwrap();
        let rows: Vec<(u64, Regime)> = v["rows"]
            .as_array()
            .unwrap()
            .iter()
            .map(|r| {
                let r: Vec<u64> = r
                    .as_array()
                    .unwrap()
                    .iter()
                    .map(|x| x.as_u64().unwrap())
                    .collect();
                let session = (r[1] as u8, r[2] as u8, r[3] as u8, r[4], r[5]);
                (r[0], regime((r[1] != 0).then_some(session), r[0]))
            })
            .collect();
        let count = |code| rows.iter().filter(|(_, g)| g.code() == code).count();
        assert_eq!(
            (count(0), count(1), count(2), count(3)),
            (0, 1_056, 1_880, 135)
        );
        for closed in [dollars(1_100), dollars(1_400), dollars(3_000)] {
            let open = dollars(1_000);
            let buffered = (open * U256::from(5)).div_ceil(U256::from(4));
            let across = buffered.max(closed);
            let largest =
                (across - buffered) * U256::from(step) / U256::from(RAMP_MS) + U256::from(1);
            for pair in rows.windows(2) {
                let (before, after) = (
                    current(open, closed, U256::ZERO, pair[0].1),
                    current(open, closed, U256::ZERO, pair[1].1),
                );
                if after > before {
                    let closing = |r: Regime| matches!(r, Regime::Closing { .. });
                    assert!(
                        closing(pair[0].1) || closing(pair[1].1),
                        "{} to {}",
                        pair[0].0,
                        pair[1].0
                    );
                    assert!(after - before <= largest, "{} to {}", pair[0].0, pair[1].0);
                }
            }
        }
    }

    #[test]
    fn the_ramp_starts_three_hours_before_the_regular_close() {
        let close = 1_789_761_600_000u64;
        let at = |now_ms| regime(Some((2, 1, 3, close, close + 14_400_000)), now_ms);
        assert_eq!(at(close - 10_800_000), Regime::Open);
        assert_eq!(at(close - 10_799_999), Regime::Closing { elapsed_ms: 1 });
        assert_eq!(
            regime(Some((2, 3, 1, 0, close + 14_400_000)), close),
            Regime::Closing {
                elapsed_ms: 10_800_000
            }
        );
    }

    #[test]
    fn an_open_session_past_its_boundary_reads_as_closed() {
        let boundary = 1_789_776_000_000u64;
        let at = |now_ms| regime(Some((2, 3, 1, 0, boundary)), now_ms);
        assert_eq!(
            at(boundary - 1),
            Regime::Closing {
                elapsed_ms: RAMP_MS - 1
            }
        );
        assert_eq!(at(boundary), Regime::Closed);
        assert_eq!(at(boundary + 1), Regime::Closed);
    }

    #[test]
    fn the_buffer_rounds_up() {
        assert_eq!(
            current(U256::from(1), U256::ZERO, U256::ZERO, Regime::Open),
            U256::from(2)
        );
        assert_eq!(
            current(U256::from(4), U256::ZERO, U256::ZERO, Regime::Open),
            U256::from(5)
        );
        assert_eq!(
            current(U256::from(5), U256::ZERO, U256::ZERO, Regime::Closed),
            U256::from(7)
        );
    }

    #[test]
    fn a_holiday_eve_signed_as_a_short_close_starts_its_ramp_at_the_close() {
        let close = 1_789_761_600_000u64;
        assert_eq!(regime(Some((2, 1, 2, close, 0)), close - 1), Regime::Open);
        assert_eq!(
            regime(Some((2, 2, 3, 0, close + 14_400_000)), close),
            Regime::Closing {
                elapsed_ms: 10_800_000
            }
        );
    }

    proptest! {
        #[test]
        fn no_scheduled_change_of_regime_raises_the_requirement_at_once(
            open in 0u128..1_000_000_000_000_000_000_000_000_000,
            extra in 0u128..1_000_000_000_000_000_000_000_000_000,
            floor in 0u128..2_000_000_000_000_000_000_000_000_000,
            elapsed in 0u64..RAMP_MS,
        ) {
            let (open, closed, floor) = (U256::from(open), U256::from(open + extra), U256::from(floor));
            let value = |regime| current(open, closed, floor, regime);
            prop_assert_eq!(value(Regime::Closing { elapsed_ms: 0 }), value(Regime::Open));
            prop_assert_eq!(value(Regime::Closing { elapsed_ms: RAMP_MS }), value(Regime::Closed));
            let later = value(Regime::Closing { elapsed_ms: elapsed + 1 });
            let earlier = value(Regime::Closing { elapsed_ms: elapsed });
            prop_assert!(later >= earlier);
            prop_assert!(value(Regime::Open) <= value(Regime::Closed));
            prop_assert_eq!(value(Regime::Unknown), value(Regime::Closed));
        }
    }
}
