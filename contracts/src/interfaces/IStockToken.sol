// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @notice A Robinhood Stock Token, as the accounts read it: whether its transfers are paused, and the registry that
/// holds its issuer's blocklist.
interface IStockToken {
    /// @notice Whether the token's transfers are paused, by its own pause or its registry's.
    function paused() external view returns (bool);

    // slither-disable-start naming-convention
    /// @notice The registry that holds the issuer's roles, pause and blocklist.
    // forge-lint: disable-next-line(mixed-case-function)
    function ACCESS_CONTROLLED_REGISTRY() external view returns (address);
    // slither-disable-end naming-convention
}

/// @notice The Stock Tokens' registry, as the accounts read it.
interface IStockTokenRegistry {
    /// @notice Whether the issuer has blocklisted `account`.
    function isBlocked(address account) external view returns (bool);
}
