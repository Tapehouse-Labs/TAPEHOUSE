// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";

/// @notice Arbitrum's wrapped ether, Robinhood Chain's WETH: an ERC-20 with EIP-2612 permits, minted one for one
/// against the ether sent to it.
interface IWETH is IERC20, IERC20Permit {
    /// @notice Wraps the ether sent for the caller.
    function deposit() external payable;

    /// @notice Wraps the ether sent for `account`.
    function depositTo(address account) external payable;

    /// @notice Unwraps `amount` to the caller.
    function withdraw(uint256 amount) external;
}
