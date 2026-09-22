// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {AggregatorProxy} from "./Interfaces.sol";

contract ChainlinkArbitrumForkTest is Test {
    uint256 internal constant WEEKDAY_BLOCK = 507_888_520;
    uint256 internal constant SUNDAY_BLOCK = 507_260_840;
    uint256 internal constant SUNDAY_TIMESTAMP = 1_789_947_000;
    uint256 internal constant NYSE_OPEN = 1_789_997_400;
    uint256 internal constant HEARTBEAT = 86_400;

    string internal arbitrum = vm.readFile("../deployments/42161.json");

    function test_FeedsHaveEightDecimals() public {
        vm.createSelectFork("arbitrum", WEEKDAY_BLOCK);
        string[] memory keys = vm.parseJsonKeys(arbitrum, ".chainlink");
        for (uint256 i; i < keys.length; ++i) {
            AggregatorProxy feed = _feed(keys[i]);
            assertEq(feed.decimals(), 8, keys[i]);
            (uint80 roundId,,,,) = feed.latestRoundData();
            assertEq(roundId >> 64, feed.phaseId(), keys[i]);
        }
    }

    function test_WeekendHeartbeatsRepeatTheClose() public {
        uint256 weekday = vm.createFork("arbitrum", WEEKDAY_BLOCK);
        uint256 sunday = vm.createSelectFork("arbitrum", SUNDAY_BLOCK);
        assertEq(block.timestamp, SUNDAY_TIMESTAMP);
        string[] memory keys = vm.parseJsonKeys(arbitrum, ".chainlink");
        for (uint256 i; i < keys.length; ++i) {
            AggregatorProxy feed = _feed(keys[i]);
            (uint80 roundId, int256 answer,, uint256 updatedAt,) = feed.latestRoundData();
            assertLt(SUNDAY_TIMESTAMP - updatedAt, HEARTBEAT, keys[i]);
            (, int256 saturday,, uint256 saturdayAt,) = feed.getRoundData(roundId - 1);
            assertEq(answer, saturday, keys[i]);
            assertApproxEqAbs(updatedAt - saturdayAt, HEARTBEAT, 1 minutes, keys[i]);

            vm.selectFork(weekday);
            (,,, uint256 mondayAt,) = feed.getRoundData(roundId + 1);
            assertGe(mondayAt, NYSE_OPEN, keys[i]);
            vm.selectFork(sunday);
        }
    }

    function _feed(string memory key) internal view returns (AggregatorProxy) {
        return AggregatorProxy(vm.parseJsonAddress(arbitrum, string.concat(".chainlink.", key)));
    }
}
