// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {FullMath} from "@uniswap/v3-core/contracts/libraries/FullMath.sol";
import {TickMath} from "@uniswap/v3-core/contracts/libraries/TickMath.sol";
import {AggregatorV3Interface} from "../conformance/Interfaces.sol";

/// @notice Reads of Uniswap v3 pools and Chainlink feeds, and the pool arithmetic of the liquidity add-on, as
/// the Stylus program performs them: a failed read or a value that does not fit gives no answer.
library Pools {
    uint32 internal constant WINDOW = 1_800;
    uint256 internal constant MAX_FEED_AGE = 86_400 + 60;
    uint256 internal constant SELL_FACTOR = 51_316_701;
    uint256 internal constant BUY_FACTOR = 48_808_848;

    /// @notice The arithmetic-mean tick and harmonic-mean liquidity of `pool` over the last 30 minutes, its
    /// cumulatives wrapping as the pool's `int56` and `uint160` do. No answer for a mean tick outside Uniswap's range.
    function consult(address pool) internal view returns (bool ok, int24 tick, uint128 liquidity) {
        uint32[] memory secondsAgos = new uint32[](2);
        secondsAgos[0] = WINDOW;
        try IUniswapV3Pool(pool).observe(secondsAgos) returns (int56[] memory ticks, uint160[] memory seconds_) {
            if (ticks.length < 2 || seconds_.length < 2) return (false, 0, 0);
            int56 tickDelta;
            uint160 secondsDelta;
            unchecked {
                tickDelta = ticks[1] - ticks[0];
                secondsDelta = seconds_[1] - seconds_[0];
            }
            int256 delta = int256(tickDelta);
            int256 mean = delta / int256(uint256(WINDOW));
            if (delta < 0 && delta % int256(uint256(WINDOW)) != 0) mean -= 1;
            uint256 perLiquidity = uint256(secondsDelta) << 32;
            if (perLiquidity == 0) return (false, 0, 0);
            uint256 harmonic = uint256(WINDOW) * type(uint160).max / perLiquidity;
            if (mean < TickMath.MIN_TICK || mean > TickMath.MAX_TICK || harmonic > type(uint128).max) {
                return (false, 0, 0);
            }
            return (true, int24(mean), uint128(harmonic));
        } catch {
            return (false, 0, 0);
        }
    }

    /// @notice The latest answer of `feed`, if positive and at most a heartbeat and a minute old at `now`.
    function answer(address feed, uint256 now_) internal view returns (bool ok, uint256 value) {
        try AggregatorV3Interface(feed).latestRoundData() returns (
            uint80, int256 answer_, uint256, uint256 updatedAt, uint80
        ) {
            if (updatedAt > type(uint64).max) return (false, 0);
            if (answer_ > 0 && updatedAt <= now_ && now_ - updatedAt <= MAX_FEED_AGE) return (true, uint256(answer_));
            return (false, 0);
        } catch {
            return (false, 0);
        }
    }

    /// @notice From a pool's mean tick and liquidity: the asset's price in USD with 8 decimals, and the USD with
    /// 18 decimals the pool pays out for a 10% fall and takes in for a 10% rise.
    function terms(
        bool stockIsToken0,
        uint8 stockDecimals,
        uint8 quoteDecimals,
        int24 tick,
        uint128 liquidity,
        uint256 usd
    ) internal pure returns (bool ok, uint256 price, uint256 selling, uint256 buying) {
        uint256 sqrt = TickMath.getSqrtRatioAtTick(stockIsToken0 ? tick : -tick);
        uint256 quoteUnit = 10 ** quoteDecimals;
        uint256 perStock = FullMath.mulDiv(FullMath.mulDiv(sqrt, sqrt, 1 << 64), 10 ** stockDecimals, 1 << 128);
        (bool fits, uint256 product) = _mul(perStock, usd);
        if (!fits) return (false, 0, 0, 0);
        price = product / quoteUnit;
        uint256 moved = FullMath.mulDiv(liquidity, sqrt, 1 << 96);
        (fits, selling) = _value(moved, SELL_FACTOR, usd, quoteUnit);
        if (!fits) return (false, 0, 0, 0);
        (fits, buying) = _value(moved, BUY_FACTOR, usd, quoteUnit);
        if (!fits) return (false, 0, 0, 0);
        ok = true;
    }

    function _value(uint256 moved, uint256 factor, uint256 usd, uint256 quoteUnit)
        private
        pure
        returns (bool ok, uint256 value)
    {
        (bool fits, uint256 raw) = _mul(moved, factor);
        if (!fits) return (false, 0);
        (fits, value) = _mul(raw / 1e9, usd);
        if (!fits) return (false, 0);
        (fits, value) = _mul(value, 1e18 / quoteUnit);
        if (!fits) return (false, 0);
        return (true, value / 1e8);
    }

    function _mul(uint256 a, uint256 b) private pure returns (bool, uint256) {
        unchecked {
            if (a == 0) return (true, 0);
            uint256 c = a * b;
            return (c / a == b, c);
        }
    }
}
