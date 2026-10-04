// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

/// @notice The dev node's Stock Token: `StubStockToken`, an 18-decimal ERC-20 with an ERC-8056 multiplier, an oracle the
/// issuer can pause, and the issuer's pause, blocklist and admin burn, with the EIP-2612 permit Robinhood's Stock Tokens
/// carry, in the domain their `eip712Domain()` reads: the token's name and version 1. Anyone may mint and act as the
/// issuer. Test code only.
contract StubPermitStockToken is ERC20Permit {
    uint256 internal multiplier;
    uint256 public newUIMultiplier;
    uint256 public effectiveAt;
    bool public oraclePaused;
    bool public paused;
    mapping(address => bool) public isBlocked;

    error IsPaused();
    error Blocked(address account);

    constructor(uint256 multiplier_) ERC20("Stub Stock Token", "STOCK") ERC20Permit("Stub Stock Token") {
        multiplier = multiplier_;
        newUIMultiplier = multiplier_;
    }

    function uiMultiplier() public view returns (uint256) {
        return block.timestamp >= effectiveAt ? newUIMultiplier : multiplier;
    }

    function updateMultiplier(uint256 newMultiplier, uint256 effectiveAt_) external {
        require(effectiveAt_ >= block.timestamp);
        multiplier = uiMultiplier();
        newUIMultiplier = newMultiplier;
        effectiveAt = effectiveAt_;
    }

    function mint(address to, uint256 value) external {
        _mint(to, value);
    }

    function pause() external {
        paused = true;
    }

    function unpause() external {
        paused = false;
    }

    // forge-lint: disable-next-line(mixed-case-function)
    function ACCESS_CONTROLLED_REGISTRY() external view returns (address) {
        return address(this);
    }

    function blockAccount(address account, bool blocked) external {
        isBlocked[account] = blocked;
    }

    function adminBurn(address from, uint256 amount) external {
        super._update(from, address(0), amount);
    }

    function _update(address from, address to, uint256 value) internal override {
        if (paused) revert IsPaused();
        if (isBlocked[from]) revert Blocked(from);
        if (isBlocked[to]) revert Blocked(to);
        if (isBlocked[msg.sender]) revert Blocked(msg.sender);
        super._update(from, to, value);
    }

    function pauseOracle() external {
        oraclePaused = true;
    }

    function unpauseOracle() external {
        oraclePaused = false;
    }
}
