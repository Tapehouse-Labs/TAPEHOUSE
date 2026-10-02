// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @title The borrower of a Tapehouse stock lending vault
/// @notice What a stock lending vault's borrower must do when a recall's notice runs out.
interface IStockLendingBorrower {
    /// @notice Brings back `assets` of the vault's token, buying them if it must, and repays them through the vault's
    /// `repay` before returning. The vault calls it once a recall's notice has run out.
    function buyIn(uint256 assets) external;
}
