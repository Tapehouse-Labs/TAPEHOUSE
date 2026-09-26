//! Reads of Robinhood Stock Tokens through their ERC-8056 multiplier interface.

use stylus_sdk::alloy_primitives::Address;
use stylus_sdk::prelude::*;

sol_interface! {
    interface IStockToken {
        function uiMultiplier() external view returns (uint256);
        function newUIMultiplier() external view returns (uint256);
        function effectiveAt() external view returns (uint256);
    }
}

/// The token's scheduled multiplier and when it takes effect, in seconds. `None` when a call fails or
/// a value does not fit.
pub fn schedule(host: &impl Host, token: Address) -> Option<(u128, u64)> {
    let token = IStockToken::new(token);
    Some((
        u128::try_from(token.new_ui_multiplier(host, Call::new()).ok()?).ok()?,
        u64::try_from(token.effective_at(host, Call::new()).ok()?).ok()?,
    ))
}

/// The token's current multiplier. `None` when the call fails or the value does not fit.
pub fn multiplier(host: &impl Host, token: Address) -> Option<u128> {
    u128::try_from(
        IStockToken::new(token)
            .ui_multiplier(host, Call::new())
            .ok()?,
    )
    .ok()
}
