// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

contract StubStockToken {
    uint256 internal multiplier;
    uint256 public newUIMultiplier;
    uint256 public effectiveAt;

    constructor(uint256 multiplier_) {
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
}
