// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

contract StubAggregator {
    uint8 public immutable decimals;
    int256 internal answer;
    uint256 internal updatedAt;
    string public description;

    constructor(uint8 decimals_, int256 answer_, uint256 updatedAt_, string memory description_) {
        decimals = decimals_;
        answer = answer_;
        updatedAt = updatedAt_;
        description = description_;
    }

    function setRound(int256 answer_, uint256 updatedAt_) external {
        answer = answer_;
        updatedAt = updatedAt_;
    }

    function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80) {
        return (1, answer, updatedAt, updatedAt, 1);
    }
}
