// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @notice A Chainlink aggregator that keeps every round it is given, numbered from one in phase one as a Chainlink
/// proxy numbers them, and answers a round it does not have with zeros, as the proxy does. Test code only.
contract StubAggregator {
    struct Round {
        int256 answer;
        uint64 startedAt;
        uint64 updatedAt;
    }

    uint80 internal constant PHASE = uint80(1) << 64;

    uint8 public immutable decimals;
    string public description;
    Round[] internal rounds;

    constructor(uint8 decimals_, int256 answer_, uint256 updatedAt_, string memory description_) {
        decimals = decimals_;
        description = description_;
        rounds.push(Round(answer_, uint64(updatedAt_), uint64(updatedAt_)));
    }

    function setRound(int256 answer_, uint256 updatedAt_) external {
        rounds.push(Round(answer_, uint64(updatedAt_), uint64(updatedAt_)));
    }

    function setRound(int256 answer_, uint256 startedAt_, uint256 updatedAt_) external {
        rounds.push(Round(answer_, uint64(startedAt_), uint64(updatedAt_)));
    }

    function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80) {
        return getRoundData(PHASE | uint80(rounds.length));
    }

    function getRoundData(uint80 roundId) public view returns (uint80, int256, uint256, uint256, uint80) {
        uint256 index = roundId >> 64 == 1 ? uint64(roundId) : 0;
        if (index == 0 || index > rounds.length) return (roundId, 0, 0, 0, roundId);
        Round memory r = rounds[index - 1];
        return (roundId, r.answer, r.startedAt, r.updatedAt, roundId);
    }
}
