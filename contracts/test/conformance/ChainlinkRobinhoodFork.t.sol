// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {AggregatorProxy, AggregatorV3Interface, IStockToken, IUniswapV3Pool} from "./Interfaces.sol";

contract ChainlinkRobinhoodForkTest is Test {
    uint256 internal constant WEEKDAY_BLOCK = 69_922_505;
    uint256 internal constant SUNDAY_BLOCK = 68_333_447;
    uint256 internal constant ARBITRUM_SUNDAY_BLOCK = 507_260_840;
    uint256 internal constant SUNDAY_TIMESTAMP = 1_789_947_000;
    uint256 internal constant SESSION_OPEN = 1_789_948_800;
    uint256 internal constant HEARTBEAT = 86_400;
    uint256 internal constant DEVIATION_BPS = 50;
    uint80 internal constant ARBITRUM_NVDA_FRIDAY_ROUND = (2 << 64) | 1598;

    string[6] internal equities = ["NVDA", "TSLA", "AAPL", "MSFT", "GOOGL", "SPY"];
    string[2] internal crypto = ["USDG", "ETH"];
    string internal robinhood = vm.readFile("../deployments/4663.json");
    string internal arbitrum = vm.readFile("../deployments/42161.json");

    function test_FeedsAreVersionSixInPhaseOneWithEightDecimals() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        string[] memory keys = vm.parseJsonKeys(robinhood, ".chainlink");
        for (uint256 i; i < keys.length; ++i) {
            AggregatorProxy feed =
                AggregatorProxy(vm.parseJsonAddress(robinhood, string.concat(".chainlink.", keys[i])));
            assertEq(feed.decimals(), 8, keys[i]);
            assertEq(feed.version(), 6, keys[i]);
            assertEq(feed.phaseId(), 1, keys[i]);
            (uint80 roundId,,,,) = feed.latestRoundData();
            assertEq(roundId >> 64, 1, keys[i]);
            assertGt(feed.aggregator().code.length, 0, keys[i]);
        }
    }

    function test_EveryFeedIsWithinHeartbeatOnAWeekday() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        assertEq(_dayOfWeek(block.timestamp), 2);
        string[] memory keys = vm.parseJsonKeys(robinhood, ".chainlink");
        for (uint256 i; i < keys.length; ++i) {
            (,,, uint256 updatedAt,) = _feed(robinhood, keys[i]).latestRoundData();
            assertLt(block.timestamp - updatedAt, HEARTBEAT + 1 minutes, keys[i]);
        }
    }

    function test_EquityFeedsHoldFridayPriceUntilTheSessionReopens() public {
        uint256 weekday = vm.createFork("robinhood", WEEKDAY_BLOCK);
        uint256 sunday = vm.createSelectFork("robinhood", SUNDAY_BLOCK);
        assertEq(block.timestamp, SUNDAY_TIMESTAMP);
        assertEq(_dayOfWeek(SUNDAY_TIMESTAMP), 0);
        for (uint256 i; i < equities.length; ++i) {
            string memory key = string.concat(equities[i], "_USD");
            AggregatorV3Interface feed = _feed(robinhood, key);
            (uint80 roundId,,, uint256 updatedAt,) = feed.latestRoundData();
            assertEq(_dayOfWeek(updatedAt), 5, key);
            assertGt(SUNDAY_TIMESTAMP - updatedAt, 48 hours, key);

            vm.selectFork(weekday);
            (,,, uint256 nextUpdatedAt,) = feed.getRoundData(roundId + 1);
            assertGe(nextUpdatedAt, SESSION_OPEN, key);
            assertLt(nextUpdatedAt, SESSION_OPEN + 1 minutes, key);
            vm.selectFork(sunday);
        }
        for (uint256 i; i < crypto.length; ++i) {
            string memory key = string.concat(crypto[i], "_USD");
            (,,, uint256 updatedAt,) = _feed(robinhood, key).latestRoundData();
            assertLt(SUNDAY_TIMESTAMP - updatedAt, HEARTBEAT, key);
        }
    }

    function test_RoundsAreEitherHeartbeatsOrDeviations() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        assertEq(_heartbeats("SPY_USD", 127, 138), 3);
        assertEq(_heartbeats("NVDA_USD", 1062, 1083), 0);
    }

    function test_PriceIsTheRawSharePriceTimesTheUiMultiplier() public {
        vm.createSelectFork("arbitrum", ARBITRUM_SUNDAY_BLOCK);
        (, int256 raw,, uint256 rawUpdatedAt,) = _feed(arbitrum, "NVDA_USD").getRoundData(ARBITRUM_NVDA_FRIDAY_ROUND);

        vm.createSelectFork("robinhood", SUNDAY_BLOCK);
        assertEq(_dayOfWeek(block.timestamp), 0);
        (, int256 answer,, uint256 updatedAt,) = _feed(robinhood, "NVDA_USD").latestRoundData();
        uint256 multiplier = IStockToken(vm.parseJsonAddress(robinhood, ".tokens.NVDA")).uiMultiplier();

        assertEq(multiplier, 1_000_775_159_164_630_595);
        assertApproxEqAbs(updatedAt, rawUpdatedAt, 5);
        assertEq(uint256(answer), uint256(raw) * multiplier / 1e18);
    }

    function test_PoolQuotesTheMultipliedPrice() public {
        vm.createSelectFork("robinhood", WEEKDAY_BLOCK);
        (, int256 answer,,,) = _feed(robinhood, "NVDA_USD").latestRoundData();
        IUniswapV3Pool pool = IUniswapV3Pool(vm.parseJsonAddress(robinhood, ".uniswapV3.NVDA_USDG_500"));
        assertEq(pool.token0(), vm.parseJsonAddress(robinhood, ".tokens.USDG"));

        (uint160 sqrtPriceX96, int24 tick,,,,,) = pool.slot0();
        uint256 spot = 1e20 * 2 ** 96 / (uint256(sqrtPriceX96) * sqrtPriceX96 >> 96);
        assertApproxEqRel(spot, uint256(answer), 0.001e18);

        uint32[] memory secondsAgos = new uint32[](2);
        secondsAgos[0] = 1800;
        (int56[] memory cumulatives,) = pool.observe(secondsAgos);
        int56 twapTick = (cumulatives[1] - cumulatives[0]) / 1800;
        assertApproxEqAbs(twapTick, tick, 25);
    }

    function _heartbeats(string memory key, uint64 first, uint64 last) internal view returns (uint256 heartbeats) {
        AggregatorProxy feed = AggregatorProxy(vm.parseJsonAddress(robinhood, string.concat(".chainlink.", key)));
        uint80 phase = uint80(feed.phaseId()) << 64;
        (, int256 previous,, uint256 previousAt,) = feed.getRoundData(phase | first);
        for (uint80 round = first + 1; round <= last; ++round) {
            (, int256 answer,, uint256 updatedAt,) = feed.getRoundData(phase | round);
            uint256 gap = updatedAt - previousAt;
            assertLe(gap, HEARTBEAT + 1 minutes, key);
            if (gap >= HEARTBEAT) {
                ++heartbeats;
            } else {
                uint256 move = uint256(answer > previous ? answer - previous : previous - answer);
                assertGe(move * 10_000, DEVIATION_BPS * uint256(previous), key);
            }
            (previous, previousAt) = (answer, updatedAt);
        }
    }

    function _feed(string memory json, string memory key) internal pure returns (AggregatorV3Interface) {
        return AggregatorV3Interface(vm.parseJsonAddress(json, string.concat(".chainlink.", key)));
    }

    function _dayOfWeek(uint256 timestamp) internal pure returns (uint256) {
        return (timestamp / 1 days + 4) % 7;
    }
}
