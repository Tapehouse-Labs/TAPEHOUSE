// SPDX-License-Identifier: MIT OR Apache-2.0
//! Reads of Chainlink data feeds through `AggregatorV3Interface`.

use stylus_sdk::alloy_primitives::{Address, I256, U256};
use stylus_sdk::prelude::*;

/// Heartbeat of the Robinhood Chain equity feeds, in seconds.
pub const CL_HEARTBEAT_S: u64 = 86_400;

/// Deviation threshold of the Robinhood Chain equity feeds, in basis points.
pub const CL_DEV_BPS: u64 = 50;

/// Maximum age of a live Chainlink answer: the heartbeat plus 60 s, because heartbeat rounds land
/// up to 27 s late.
pub const CL_MAX_AGE_S: u64 = CL_HEARTBEAT_S + 60;

/// Decimals every configured feed must report.
pub const DECIMALS: u8 = 8;

/// How long the L2 sequencer must have been up before its chain's prices are trusted again, in
/// seconds, as Chainlink's L2 sequencer guidance recommends: prices count once it has been up for
/// longer.
pub const SEQUENCER_GRACE_S: u64 = 3_600;

sol_interface! {
    interface AggregatorV3Interface {
        function decimals() external view returns (uint8);
        function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80);
    }
}

/// Whether `feed` answers `decimals()` with [`DECIMALS`].
pub fn has_expected_decimals(host: &impl Host, feed: Address) -> bool {
    AggregatorV3Interface::new(feed).decimals(host, Call::new()) == Ok(DECIMALS)
}

/// The latest answer of `feed`, its `startedAt` and its `updatedAt` in seconds, or zeros when the feed
/// is unset, the call fails, or the answer is not positive. Never reverts and never judges staleness.
pub fn latest(host: &impl Host, feed: Address) -> (U256, u64, u64) {
    if feed == Address::ZERO {
        return (U256::ZERO, 0, 0);
    }
    match AggregatorV3Interface::new(feed).latest_round_data(host, Call::new()) {
        Ok((_, answer, started_at, updated_at, _)) if answer.is_positive() => {
            match (u64::try_from(started_at), u64::try_from(updated_at)) {
                (Ok(started_at), Ok(updated_at)) => (answer.into_raw(), started_at, updated_at),
                _ => (U256::ZERO, 0, 0),
            }
        }
        _ => (U256::ZERO, 0, 0),
    }
}

/// Whether an L2 sequencer-uptime round says the sequencer is up and has been for more than
/// [`SEQUENCER_GRACE_S`] at `now`, as in Chainlink's reference consumer: `answer` 0 means up, and
/// `started_at` is when that status began.
pub fn sequencer_settled(answer: I256, started_at: U256, now: u64) -> bool {
    answer.is_zero()
        && u64::try_from(started_at)
            .is_ok_and(|since| since != 0 && since <= now && now - since > SEQUENCER_GRACE_S)
}

/// Whether the sequencer-uptime `feed` says the sequencer is up and settled at `now`. True for an
/// unset feed; false when the call fails.
pub fn sequencer_up(host: &impl Host, feed: Address, now: u64) -> bool {
    if feed == Address::ZERO {
        return true;
    }
    AggregatorV3Interface::new(feed)
        .latest_round_data(host, Call::new())
        .is_ok_and(|(_, answer, started_at, _, _)| sequencer_settled(answer, started_at, now))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn a_sequencer_is_settled_an_hour_after_it_comes_up() {
        let up = I256::ZERO;
        let down = I256::ONE;
        let now = 1_790_000_000;
        let since = |s: u64| U256::from(s);
        assert!(sequencer_settled(
            up,
            since(now - SEQUENCER_GRACE_S - 1),
            now
        ));
        assert!(!sequencer_settled(up, since(now - SEQUENCER_GRACE_S), now));
        assert!(!sequencer_settled(down, since(now - 86_400), now));
        assert!(!sequencer_settled(up, U256::ZERO, now));
        assert!(!sequencer_settled(up, since(now + 1), now));
        assert!(!sequencer_settled(up, U256::MAX, now));
    }
}
