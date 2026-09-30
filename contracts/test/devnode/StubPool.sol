// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {TickMath} from "@uniswap/v3-core/contracts/libraries/TickMath.sol";

contract StubPool {
    address public immutable token0;
    address public immutable token1;
    uint24 public immutable fee;
    int24 internal immutable tick;
    uint160 internal immutable secondsPerLiquidity;
    bool internal young;

    constructor(address token0_, address token1_, uint24 fee_, int24 tick_, uint128 liquidity) {
        token0 = token0_;
        token1 = token1_;
        fee = fee_;
        tick = tick_;
        secondsPerLiquidity = uint160((uint256(1) << 128) / liquidity);
    }

    function setYoung(bool young_) external {
        young = young_;
    }

    function slot0() external view returns (uint160, int24, uint16, uint16, uint16, uint8, bool) {
        return (TickMath.getSqrtRatioAtTick(tick), tick, 0, 1801, 1801, 0, true);
    }

    function observations(uint256) external view returns (uint32, int56, uint160, bool) {
        uint256 time = young ? block.timestamp : block.timestamp - 1 days;
        return (uint32(time), int56(tick) * int56(uint56(time)), secondsPerLiquidity * uint160(time), true);
    }

    function observe(uint32[] calldata secondsAgos) external view returns (int56[] memory, uint160[] memory) {
        require(!young, "OLD");
        int56[] memory tickCumulatives = new int56[](secondsAgos.length);
        uint160[] memory secondsPerLiquidityCumulatives = new uint160[](secondsAgos.length);
        for (uint256 i; i < secondsAgos.length; ++i) {
            uint256 time = block.timestamp - secondsAgos[i];
            tickCumulatives[i] = int56(tick) * int56(uint56(time));
            secondsPerLiquidityCumulatives[i] = secondsPerLiquidity * uint160(time);
        }
        return (tickCumulatives, secondsPerLiquidityCumulatives);
    }
}
