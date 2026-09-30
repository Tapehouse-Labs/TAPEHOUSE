// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

contract TokenDouble {
    uint8 public immutable decimals;

    constructor(uint8 decimals_) {
        decimals = decimals_;
    }
}

contract PoolDouble {
    address public immutable token0;
    address public immutable token1;
    int24 public meanTick;
    uint32 public history;
    uint16 public cardinality = 1;

    constructor(address token0_, address token1_) {
        token0 = token0_;
        token1 = token1_;
    }

    function setTwap(int24 meanTick_, uint32 history_) external {
        meanTick = meanTick_;
        history = history_;
    }

    function setCardinality(uint16 cardinality_) external {
        cardinality = cardinality_;
    }

    function slot0() external view returns (uint160, int24, uint16, uint16, uint16, uint8, bool) {
        return (0, 0, 0, cardinality, cardinality, 0, true);
    }

    function observations(uint256) external view returns (uint32, int56, uint160, bool) {
        return (uint32(block.timestamp) - history, 0, 0, true);
    }

    function observe(uint32[] calldata secondsAgos)
        external
        view
        returns (int56[] memory ticks, uint160[] memory liquidity)
    {
        ticks = new int56[](2);
        liquidity = new uint160[](2);
        ticks[1] = int56(meanTick) * int56(uint56(secondsAgos[0]));
        liquidity[1] = 1;
    }
}
