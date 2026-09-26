//! Chainlink's 24/5 session, derived from RedStone's signed New York market status.
//!
//! The session runs from 20:00 ET on the evening before each trading day to 20:00 ET on it. The
//! signed status says whether NYSE is in regular hours, a short close between trading days or a long
//! close over a weekend or holiday, and when that changes, in UTC milliseconds. Every boundary of the
//! session is a fixed distance from a regular close or from that change time, and none of those
//! distances spans a daylight-saving change, so no time zone or calendar is needed.

use stylus_sdk::alloy_primitives::{B256, U256, b256};

/// Feed ID of `NY_MARKET_CURRENT_STATUS`.
pub const CURRENT_STATUS: B256 =
    b256!("0x4e595f4d41524b45545f43555252454e545f5354415455530000000000000000");
/// Feed ID of `NY_MARKET_NEXT_STATUS`.
pub const NEXT_STATUS: B256 =
    b256!("0x4e595f4d41524b45545f4e4558545f5354415455530000000000000000000000");
/// Feed ID of `NY_MARKET_NEXT_CHANGE_TIME`.
pub const NEXT_CHANGE_TIME: B256 =
    b256!("0x4e595f4d41524b45545f4e4558545f4348414e47455f54494d45000000000000");

/// The three status feeds, which are written together or not at all.
pub const FEEDS: [B256; 3] = [CURRENT_STATUS, NEXT_STATUS, NEXT_CHANGE_TIME];

/// Scale of every signed value.
pub const SCALE: u64 = 100_000_000;
/// Post-market after a regular close (16:00 to 20:00 ET, or 13:00 to 17:00 ET on an early close), in
/// milliseconds.
pub const POST_MARKET_MS: u64 = 14_400_000;
/// From the reopen at 20:00 ET to the regular open at 09:30 ET, in milliseconds.
pub const REOPEN_BEFORE_OPEN_MS: u64 = 48_600_000;
/// From the reopen at 20:00 ET to midnight ET, where a short close follows a holiday, in milliseconds.
pub const REOPEN_BEFORE_MIDNIGHT_MS: u64 = 14_400_000;
/// Maximum age of the signed status, in milliseconds.
pub const STATUS_MAX_AGE_MS: u64 = 3_600_000;

/// NYSE's state, as RedStone signs it.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
#[repr(u8)]
pub enum Nyse {
    /// Regular hours.
    Regular = 1,
    /// A short close between trading days.
    ClosedShort = 2,
    /// A long close over a weekend or holiday.
    ClosedLong = 3,
}

impl Nyse {
    fn decode(value: U256) -> Option<Self> {
        match u64::try_from(value).ok()? {
            100_000_000 => Some(Self::Regular),
            101_000_000 => Some(Self::ClosedShort),
            102_000_000 => Some(Self::ClosedLong),
            _ => None,
        }
    }
}

/// The signed status.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Status {
    /// NYSE's state now.
    pub current: Nyse,
    /// NYSE's state after `change_ms`.
    pub next: Nyse,
    /// When NYSE changes state, in milliseconds.
    pub change_ms: u64,
    /// The package timestamp of all three values, in milliseconds.
    pub package_ms: u64,
}

impl Status {
    /// Decodes the stored current status, next status and change time with their package timestamps.
    /// `None` unless all three come from one package and hold known values.
    pub fn decode(values: [U256; 3], package_ms: [u64; 3]) -> Option<Self> {
        if package_ms[1] != package_ms[0] || package_ms[2] != package_ms[0] {
            return None;
        }
        let change = u128::try_from(values[2]).ok()? / u128::from(SCALE);
        let change_ms = u64::try_from(change).ok().filter(|ms| *ms > 0)?;
        Some(Self {
            current: Nyse::decode(values[0])?,
            next: Nyse::decode(values[1])?,
            change_ms,
            package_ms: package_ms[0],
        })
    }
}

/// The 24/5 session at one instant.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Session {
    /// Whether the session is open.
    pub open: bool,
    /// NYSE's state.
    pub nyse: Nyse,
    /// NYSE's next state, while its change is still ahead.
    pub nyse_next: Option<Nyse>,
    /// When NYSE changes state, in milliseconds; zero once that has passed.
    pub change_ms: u64,
    /// The session's next boundary: when an open session closes or a closed one reopens, in
    /// milliseconds; zero while not known.
    pub boundary_ms: u64,
}

/// The session at `now_ms`, from the stored status and `close_ms`, the last regular close the program
/// saw. `None` when there is no status or it is more than [`STATUS_MAX_AGE_MS`] old.
pub fn derive(status: Option<Status>, close_ms: u64, now_ms: u64) -> Option<Session> {
    let status = status.filter(|s| now_ms.saturating_sub(s.package_ms) <= STATUS_MAX_AGE_MS)?;
    let (nyse, nyse_next, close_ms, change_ms, reopens_ms) = if now_ms < status.change_ms {
        let lead = match status.next {
            Nyse::Regular => Some(REOPEN_BEFORE_OPEN_MS),
            Nyse::ClosedShort => Some(REOPEN_BEFORE_MIDNIGHT_MS),
            Nyse::ClosedLong => None,
        };
        let reopens = lead.map_or(0, |lead| status.change_ms.saturating_sub(lead));
        (
            status.current,
            Some(status.next),
            close_ms,
            status.change_ms,
            reopens,
        )
    } else {
        let close_ms = if status.current == Nyse::Regular {
            status.change_ms
        } else {
            close_ms
        };
        (status.next, None, close_ms, 0, 0)
    };
    let post_market_ends = close_ms.saturating_add(POST_MARKET_MS);
    let in_post_market = now_ms < post_market_ends;
    let open = match (nyse, nyse_next) {
        (Nyse::Regular, _) => true,
        (Nyse::ClosedShort, Some(Nyse::ClosedLong)) => in_post_market,
        (Nyse::ClosedShort, _) => true,
        (Nyse::ClosedLong, _) => in_post_market || (reopens_ms != 0 && now_ms >= reopens_ms),
    };
    let boundary_ms = match (open, nyse, nyse_next) {
        (false, Nyse::ClosedLong, _) => reopens_ms,
        (false, _, _) => 0,
        (true, Nyse::ClosedLong, _) | (true, Nyse::ClosedShort, Some(Nyse::ClosedLong))
            if in_post_market =>
        {
            post_market_ends
        }
        (true, Nyse::Regular, Some(Nyse::ClosedLong)) => change_ms.saturating_add(POST_MARKET_MS),
        _ => 0,
    };
    Some(Session {
        open,
        nyse,
        nyse_next,
        change_ms,
        boundary_ms,
    })
}

/// The last regular close after `status` is written: its change time when it was signed in regular
/// hours, otherwise `close_ms` unchanged.
pub fn next_close(status: Option<Status>, close_ms: u64) -> u64 {
    match status {
        Some(s) if s.current == Nyse::Regular => s.change_ms,
        _ => close_ms,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::Value;

    const VECTORS: &str = include_str!("../testdata/session-vectors.json");

    fn int(value: &Value) -> u64 {
        value.as_str().unwrap().parse().unwrap()
    }

    fn values(case: &Value) -> [U256; 3] {
        let v: Vec<U256> = case["values"]
            .as_array()
            .unwrap()
            .iter()
            .map(|v| v.as_str().unwrap().parse().unwrap())
            .collect();
        [v[0], v[1], v[2]]
    }

    #[test]
    fn sessions_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["sessions"].as_array().unwrap() {
            let package_ms: Vec<u64> = case["package_ms"]
                .as_array()
                .unwrap()
                .iter()
                .map(int)
                .collect();
            let status =
                Status::decode(values(case), [package_ms[0], package_ms[1], package_ms[2]]);
            let session = derive(status, int(&case["close_ms"]), int(&case["now_ms"]));
            let expected = &case["expected"];
            let name = &case["name"];
            assert_eq!(
                session.is_some(),
                expected["known"].as_bool().unwrap(),
                "{name}"
            );
            let Some(session) = session else { continue };
            assert_eq!(session.open, expected["open"].as_bool().unwrap(), "{name}");
            assert_eq!(session.nyse as u64, int(&expected["nyse"]), "{name}");
            assert_eq!(session.change_ms, int(&expected["change_ms"]), "{name}");
            assert_eq!(
                session.nyse_next.map_or(0, |n| n as u64),
                int(&expected["nyse_next"]),
                "{name}"
            );
            assert_eq!(session.boundary_ms, int(&expected["boundary_ms"]), "{name}");
        }
    }

    #[test]
    fn closes_match_the_reference() {
        let vectors: Value = serde_json::from_str(VECTORS).unwrap();
        for case in vectors["closes"].as_array().unwrap() {
            let status = Status::decode(values(case), [0; 3]);
            assert_eq!(
                next_close(status, int(&case["close_ms"])),
                int(&case["expected"]),
                "{}",
                case["name"]
            );
        }
    }

    #[test]
    fn no_status_is_no_session() {
        assert_eq!(derive(None, 1_790_000_000_000, 1_790_000_000_000), None);
        assert_eq!(next_close(None, 1_790_000_000_000), 1_790_000_000_000);
    }

    #[test]
    fn a_status_between_the_enum_values_is_unknown() {
        assert_eq!(
            Nyse::decode(U256::from(101_000_000u64)),
            Some(Nyse::ClosedShort)
        );
        assert_eq!(Nyse::decode(U256::from(101_000_001u64)), None);
        assert_eq!(Nyse::decode(U256::from(100_500_000u64)), None);
        assert_eq!(Nyse::decode(U256::ZERO), None);
        assert_eq!(Nyse::decode(U256::MAX), None);
    }
}
