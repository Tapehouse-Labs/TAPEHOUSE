// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @title Morpho Blue's oracle, as its markets call it
/// @notice The one call Morpho Blue makes of a market's oracle, to borrow, to withdraw collateral against debt and to
/// liquidate.
interface IOracle {
    /// @notice The price of one unit of the market's collateral token in units of its loan token, scaled by 1e36:
    /// the price of 10^(collateral decimals) units of collateral in 10^(loan decimals) units of the loan token, with
    /// 36 + loan decimals − collateral decimals decimals.
    function price() external view returns (uint256);
}
