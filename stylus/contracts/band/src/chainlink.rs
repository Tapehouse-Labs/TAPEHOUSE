//! Reads of Chainlink data feeds through `AggregatorV3Interface`.

use stylus_sdk::alloy_primitives::{Address, U256};
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

/// The latest answer of `feed` and its `updatedAt` in seconds, or `(0, 0)` when the feed is unset,
/// the call fails, or the answer is not positive. Never reverts and never judges staleness.
pub fn latest(host: &impl Host, feed: Address) -> (U256, u64) {
    if feed == Address::ZERO {
        return (U256::ZERO, 0);
    }
    match AggregatorV3Interface::new(feed).latest_round_data(host, Call::new()) {
        Ok((_, answer, _, updated_at, _)) if answer.is_positive() => {
            match u64::try_from(updated_at) {
                Ok(updated_at) => (answer.into_raw(), updated_at),
                Err(_) => (U256::ZERO, 0),
            }
        }
        _ => (U256::ZERO, 0),
    }
}
