// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";

/// @notice Paxos's Global Dollar (USDG): an ERC-20 with EIP-2612 permits that one address can pause, and whose
/// asset-protection role can freeze an address, which then neither sends nor receives, and wipe its balance.
interface IUSDG is IERC20, IERC20Permit {
    /// @notice Transfers, approvals and permits revert while the token is paused.
    error ContractPaused();
    /// @notice A frozen address neither sends nor receives.
    error AddressFrozen();

    /// @notice Whether transfers, approvals and permits are paused.
    function paused() external view returns (bool);

    /// @notice Whether `account` is frozen.
    function isFrozen(address account) external view returns (bool);
}
