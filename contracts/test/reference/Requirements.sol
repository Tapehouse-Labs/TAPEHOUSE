// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Scenarios, Set} from "./Scenarios.sol";

/// @notice What an asset's pool gives the liquidity add-on: absent, unread (its fee still counts), or read.
struct Pool {
    uint8 kind;
    uint256 price;
    uint256 selling;
    uint256 buying;
    uint32 fee;
}

/// @notice The worst losses of one pass over a scenario set.
struct Tally {
    int256[3] joint;
    int256[3] apart;
    int256[3][] alone;
    int256 stress;
    int256 gaps;
}

/// @notice The engine's portfolio requirement, as the Stylus program computes it: expected shortfall at 99%
/// with capped diversification, stress rows, reference floors, the liquidity add-on and the leverage floor.
library Requirements {
    uint8 internal constant ABSENT = 0;
    uint8 internal constant UNREAD = 1;
    uint8 internal constant READ = 2;
    int256 internal constant PPM = 1e6;
    int256 internal constant WAD = 1e18;
    int256 internal constant MAX_EXPOSURE = 10_000_000_000 * WAD;
    int256 internal constant STOCK_MOVE = 150_000;
    int256 internal constant INDEX_DOWN = 80_000;
    int256 internal constant INDEX_UP = 60_000;
    uint256 internal constant WEEKEND_LEVERAGE = 50_000;

    /// @notice Each position's value in USD with 18 decimals, rounded toward zero; `bad` is one plus the first
    /// asset whose value does not fit, or zero.
    function exposures(int256[] memory quantities, uint256[] memory prices)
        internal
        pure
        returns (int256[] memory values, uint256 bad)
    {
        values = new int256[](quantities.length);
        for (uint256 i; i < quantities.length; ++i) {
            if (prices[i] > uint256(type(int256).max)) return (values, i + 1);
            int256 price = int256(prices[i]);
            int256 product;
            unchecked {
                product = quantities[i] * price;
            }
            if (price != 0 && product / price != quantities[i]) return (values, i + 1);
            int256 value = product / 1e8;
            if (value > MAX_EXPOSURE || value < -MAX_EXPOSURE) return (values, i + 1);
            values[i] = value;
        }
    }

    /// @notice The loss of `values` in a scenario of `returns`.
    function loss(int256[] memory values, int256[] memory returns_) internal pure returns (int256 total) {
        unchecked {
            for (uint256 i; i < values.length; ++i) {
                total += -(values[i] * returns_[i]) / PPM;
            }
        }
    }

    /// @notice The scenario requirement over the open market and across a closure, from one pass over the set.
    function scenarioRequirements(Set memory set, int256[] memory values)
        internal
        pure
        returns (int256 open, int256 closed)
    {
        Tally memory tally = _tally(set, values);
        uint256 size = set.size;
        int256 dependent = _max(_shortfall(tally.joint, size), _shortfall(tally.apart, size));
        int256 standalone;
        for (uint256 i; i < values.length; ++i) {
            standalone += _shortfall(tally.alone[i], size);
        }
        int256 capped = dependent + _max(standalone - dependent, 0) / 5;
        open = _max(_max(_max(capped, tally.stress), _floor(values, set.parameters.market)), 0);
        closed = _max(open, tally.gaps);
    }

    /// @notice The requirement over the open market and across a closure, the leverage floor, each with the
    /// liquidity add-on, and a bit for every asset held without a read pool.
    function requirements(
        Set memory set,
        int256[] memory values,
        uint256[] memory prices,
        uint32[] memory depths,
        Pool[] memory pools
    ) internal pure returns (uint256 open, uint256 closed, uint256 floor, uint8 missing) {
        (int256 openScenario, int256 closedScenario) = scenarioRequirements(set, values);
        uint256 charge;
        (charge, missing) = liquidity(values, prices, set.parameters.gaps, depths, pools);
        open = uint256(openScenario) + charge;
        closed = uint256(closedScenario) + charge;
        floor = leverageFloor(values) + charge;
    }

    /// @notice The liquidity add-on of every position, on the side a liquidation takes.
    function liquidity(
        int256[] memory values,
        uint256[] memory prices,
        uint32[] memory gaps,
        uint32[] memory depths,
        Pool[] memory pools
    ) internal pure returns (uint256 total, uint8 missing) {
        for (uint256 i; i < values.length; ++i) {
            if (values[i] == 0) continue;
            if (pools[i].kind != READ) missing |= uint8(1 << i);
            if (pools[i].kind == ABSENT) continue;
            total += _positionAddon(values[i], prices[i], depths[values[i] > 0 ? 2 * i : 2 * i + 1], pools[i], gaps[i]);
        }
    }

    /// @notice What liquidating one position costs beyond its value: the fee, the pool's discount against the
    /// position up to the weekend gap, and the price impact.
    function addon(uint256 value, bool long, uint256 price, uint256 poolPrice, uint256 depth, uint32 fee, uint32 gap)
        internal
        pure
        returns (uint256)
    {
        uint256 adverse =
            long ? (price > poolPrice ? price - poolPrice : 0) : (poolPrice > price ? poolPrice - price : 0);
        uint256 cap = price * gap / uint256(PPM);
        if (adverse > cap) adverse = cap;
        uint256 discount = price == 0 ? 0 : value * adverse / price;
        uint256 charged = value * fee / uint256(PPM);
        uint256 impact;
        if (value <= depth) {
            impact = depth == 0 ? 0 : value * value / (20 * depth);
        } else {
            impact = depth / 20 + (value - depth);
        }
        return discount + charged + impact;
    }

    /// @notice The margin that holds the gross value of `values` to 5×, rounded up.
    function leverageFloor(int256[] memory values) internal pure returns (uint256) {
        uint256 gross;
        for (uint256 i; i < values.length; ++i) {
            gross += uint256(values[i] < 0 ? -values[i] : values[i]);
        }
        return (gross * 10_000 + WEEKEND_LEVERAGE - 1) / WEEKEND_LEVERAGE;
    }

    /// @notice One position's add-on. A read pool's own depth lowers the governed one by half at most: one
    /// second with no liquidity in range empties a harmonic mean.
    function _positionAddon(int256 exposure, uint256 price, uint32 governed, Pool memory pool, uint32 gap)
        internal
        pure
        returns (uint256)
    {
        bool long = exposure > 0;
        uint256 depth = uint256(governed) * uint256(WAD);
        if (pool.kind == READ) {
            uint256 own = long ? pool.selling : pool.buying;
            if (own < depth / 2) own = depth / 2;
            if (own < depth) depth = own;
        }
        uint256 value = uint256(long ? exposure : -exposure);
        return addon(value, long, price, pool.kind == READ ? pool.price : price, depth, pool.fee, gap);
    }

    function _tally(Set memory set, int256[] memory values) private pure returns (Tally memory tally) {
        unchecked {
            uint256 n = values.length;
            uint256 size = set.size;
            tally.joint = _worst();
            tally.apart = _worst();
            tally.alone = new int256[3][](n);
            for (uint256 i; i < n; ++i) {
                tally.alone[i] = _worst();
            }
            int256[] memory row = new int256[](n);
            int256[] memory draws = new int256[](n);
            uint256 rows = Scenarios.count(size, n);
            for (uint256 index; index < rows; ++index) {
                Scenarios.fill(set, index, row, draws);
                if (index < size) {
                    _push(tally.joint, loss(values, row));
                } else if (index < 2 * size) {
                    for (uint256 i; i < n; ++i) {
                        _push(tally.alone[i], -(values[i] * row[i]) / PPM);
                    }
                } else if (index < 3 * size) {
                    _push(tally.apart, loss(values, row));
                } else if (index < 3 * size + 2) {
                    tally.stress = _max(tally.stress, loss(values, row));
                } else {
                    tally.gaps = _max(tally.gaps, loss(values, row));
                }
            }
        }
    }

    function _floor(int256[] memory values, uint256 market) private pure returns (int256 floor) {
        for (uint256 i; i < values.length; ++i) {
            int256 adverse = market == i + 1 ? (values[i] > 0 ? INDEX_DOWN : INDEX_UP) : STOCK_MOVE;
            int256 value = (values[i] < 0 ? -values[i] : values[i]) * adverse / PPM;
            if (i == 0 || value > floor) floor = value;
        }
    }

    function _worst() private pure returns (int256[3] memory worst) {
        int256 least = type(int128).min;
        worst = [least, least, least];
    }

    function _push(int256[3] memory worst, int256 value) private pure {
        unchecked {
            for (uint256 k; k < 3; ++k) {
                if (value > worst[k]) (worst[k], value) = (value, worst[k]);
            }
        }
    }

    function _shortfall(int256[3] memory worst, uint256 size) private pure returns (int256) {
        uint256 whole = size / 100;
        int256 sum;
        for (uint256 k; k < whole; ++k) {
            sum += worst[k];
        }
        return (100 * sum + int256(size % 100) * worst[whole]) / int256(size);
    }

    function _max(int256 a, int256 b) private pure returns (int256) {
        return a > b ? a : b;
    }
}
