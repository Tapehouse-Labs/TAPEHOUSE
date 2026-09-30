// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {StubAggregator} from "../devnode/StubAggregator.sol";
import {StubPool} from "../devnode/StubPool.sol";
import {StubToken} from "../devnode/StubToken.sol";
import {MarginReference} from "./MarginReference.sol";
import {Pools} from "./Pools.sol";
import {Pool, Requirements} from "./Requirements.sol";
import {Session} from "./Session.sol";

contract ObservationsDouble {
    int56[] internal ticks;
    uint160[] internal secondsPerLiquidity;

    function set(int56 tickThen, int56 tickNow, uint160 secondsThen, uint160 secondsNow) external {
        ticks = [tickThen, tickNow];
        secondsPerLiquidity = [secondsThen, secondsNow];
    }

    function observe(uint32[] calldata) external view returns (int56[] memory, uint160[] memory) {
        return (ticks, secondsPerLiquidity);
    }
}

contract GasBurner {
    function session() external pure returns (uint8, uint8, uint8, uint64, uint64) {
        while (true) {}
        return (0, 0, 0, 0, 0);
    }
}

contract MarginReferenceTest is Test {
    struct Digest {
        string config;
        bytes32 digest;
        uint256 horizon;
        uint256 size;
    }

    string internal constant SCENARIOS = "../stylus/contracts/margin/testdata/scenario-vectors.json";
    string internal constant TIMELINE = "../stylus/contracts/margin/testdata/session-timeline.json";
    string internal constant SESSIONS = "../stylus/contracts/band/testdata/session-vectors.json";
    uint256 internal constant DOLLAR = 1e18;

    function test_ScenarioDigestsAreTheReferenceDigests() public {
        vm.pauseGasMetering();
        string memory json = vm.readFile(SCENARIOS);
        Digest[] memory digests = abi.decode(vm.parseJson(json, ".digests"), (Digest[]));
        assertEq(digests.length, 144);
        MarginReference engine;
        string memory deployed;
        for (uint256 k; k < digests.length; ++k) {
            if (keccak256(bytes(digests[k].config)) != keccak256(bytes(deployed))) {
                engine = _deploy(json, digests[k].config);
                deployed = digests[k].config;
            }
            assertEq(
                engine.scenarioDigest(uint16(digests[k].size), uint64(digests[k].horizon)),
                digests[k].digest,
                string.concat(
                    digests[k].config, " ", vm.toString(digests[k].size), " ", vm.toString(digests[k].horizon)
                )
            );
        }
    }

    function test_LaunchScenariosAreTheReferenceRows() public {
        string memory json = vm.readFile(SCENARIOS);
        MarginReference engine = _deploy(json, "launch");
        uint64 horizon = uint64(vm.parseJsonUint(json, ".launchRows.horizon"));
        uint256 rows = abi.decode(vm.parseJson(json, ".launchRows.rows[*].index"), (uint256[])).length;
        assertEq(rows, 14);
        for (uint256 k; k < rows; ++k) {
            string memory path = string.concat(".launchRows.rows[", vm.toString(k), "]");
            uint256 index = vm.parseJsonUint(json, string.concat(path, ".index"));
            int256[] memory expected = vm.parseJsonIntArray(json, string.concat(path, ".returns"));
            int32[] memory row = engine.scenario(uint16(index), horizon);
            for (uint256 i; i < row.length; ++i) {
                assertEq(int256(row[i]), expected[i], vm.toString(index));
            }
        }
    }

    function test_ScenariosOutsideTheSetAreRefused() public {
        MarginReference engine = _deploy(vm.readFile(SCENARIOS), "dev node");
        vm.expectRevert(abi.encodeWithSelector(MarginReference.ScenarioOutOfRange.selector, uint16(778)));
        engine.scenario(778, 172_800);
        vm.expectRevert(abi.encodeWithSelector(MarginReference.UnsupportedScenarioSize.selector, uint16(100)));
        engine.scenarioDigest(100, 172_800);
    }

    function test_EverySessionTheBandReportsHasItsRegime() public view {
        string memory json = vm.readFile(SESSIONS);
        string[] memory names = abi.decode(vm.parseJson(json, ".sessions[*].name"), (string[]));
        assertEq(names.length, 38);
        uint256[4] memory seen;
        for (uint256 k; k < names.length; ++k) {
            string memory path = string.concat(".sessions[", vm.toString(k), "]");
            bool known = vm.parseJsonBool(json, string.concat(path, ".expected.known"));
            uint256 open = known ? (vm.parseJsonBool(json, string.concat(path, ".expected.open")) ? 2 : 1) : 0;
            uint256 boundary = vm.parseUint(vm.parseJsonString(json, string.concat(path, ".expected.boundary_ms")));
            uint256 nowMs = vm.parseUint(vm.parseJsonString(json, string.concat(path, ".now_ms")));
            (uint8 code,) = Session.regime(open != 0, open, boundary, nowMs);
            assertEq(code, _expectedRegime(names[k]), names[k]);
            ++seen[code];
        }
        assertEq(seen[0], 6);
        assertEq(seen[1], 12);
        assertEq(seen[2], 12);
        assertEq(seen[3], 8);
    }

    function test_AMonthOfRealSessionsRaisesARequirementOnlyThroughTheRamp() public view {
        string memory json = vm.readFile(TIMELINE);
        uint256[][] memory rows = abi.decode(vm.parseJson(json, ".rows"), (uint256[][]));
        uint256 step = vm.parseJsonUint(json, ".step_ms");
        uint8[] memory codes = new uint8[](rows.length);
        uint256[] memory elapsed = new uint256[](rows.length);
        uint256[4] memory seen;
        for (uint256 k; k < rows.length; ++k) {
            (codes[k], elapsed[k]) = Session.regime(rows[k][1] != 0, rows[k][1], rows[k][5], rows[k][0]);
            ++seen[codes[k]];
        }
        assertEq(seen[0], 0);
        assertEq(seen[1], 1_056);
        assertEq(seen[2], 1_880);
        assertEq(seen[3], 135);
        uint256[3] memory closed = [uint256(1_100 * DOLLAR), 1_400 * DOLLAR, 3_000 * DOLLAR];
        for (uint256 c; c < 3; ++c) {
            _assertRisesOnlyThroughTheRamp(rows, codes, elapsed, closed[c], step);
        }
    }

    function _assertRisesOnlyThroughTheRamp(
        uint256[][] memory rows,
        uint8[] memory codes,
        uint256[] memory elapsed,
        uint256 closed,
        uint256 step
    ) internal pure {
        uint256 open = 1_000 * DOLLAR;
        uint256 buffered = open * 5 / 4;
        uint256 largest = ((buffered > closed ? buffered : closed) - buffered) * step / Session.RAMP_MS + 1;
        for (uint256 k = 1; k < rows.length; ++k) {
            uint256 before = Session.current(open, closed, 0, codes[k - 1], elapsed[k - 1]);
            uint256 later = Session.current(open, closed, 0, codes[k], elapsed[k]);
            if (later > before) {
                assertTrue(codes[k - 1] == Session.CLOSING || codes[k] == Session.CLOSING, vm.toString(rows[k][0]));
                assertLe(later - before, largest, vm.toString(rows[k][0]));
            }
        }
    }

    function test_TheRampRisesInAStraightLineToTheRequirementAcrossTheClosure() public pure {
        (uint256 open, uint256 closed) = (1_000 * DOLLAR, 1_400 * DOLLAR);
        assertEq(Session.current(open, closed, 0, Session.OPEN, 0), 1_250 * DOLLAR);
        assertEq(Session.current(open, closed, 0, Session.CLOSING, 0), 1_250 * DOLLAR);
        assertEq(Session.current(open, closed, 0, Session.CLOSING, Session.RAMP_MS / 4), 1_287 * DOLLAR + DOLLAR / 2);
        assertEq(Session.current(open, closed, 0, Session.CLOSING, Session.RAMP_MS), 1_400 * DOLLAR);
        assertEq(Session.current(open, closed, 0, Session.CLOSED, 0), 1_400 * DOLLAR);
        assertEq(Session.current(open, closed, 0, Session.UNKNOWN, 0), 1_400 * DOLLAR);
        assertEq(Session.current(open, 1_100 * DOLLAR, 0, Session.CLOSED, 0), 1_250 * DOLLAR);
    }

    function test_TheRampStartsThreeHoursBeforeTheRegularClose() public pure {
        uint256 close = 1_789_761_600_000;
        (uint8 code,) = Session.regime(true, 2, close + 14_400_000, close - 10_800_000);
        assertEq(code, Session.OPEN);
        uint256 elapsed;
        (code, elapsed) = Session.regime(true, 2, close + 14_400_000, close - 10_799_999);
        assertEq(code, Session.CLOSING);
        assertEq(elapsed, 1);
        (code, elapsed) = Session.regime(true, 2, close + 14_400_000, close);
        assertEq(elapsed, 10_800_000);
    }

    function test_AcrossAClosureGrossExposureStaysWithinFiveTimesTheRequirement() public pure {
        int256[] memory values = new int256[](3);
        values[0] = int256(4_000 * DOLLAR);
        values[1] = -int256(2_000 * DOLLAR);
        values[2] = int256(4_000 * DOLLAR);
        uint256 floor = Requirements.leverageFloor(values);
        assertEq(floor, 2_000 * DOLLAR);
        (uint256 open, uint256 closed) = (1_000 * DOLLAR, 1_400 * DOLLAR);
        assertEq(Session.current(open, closed, floor, Session.OPEN, 0), 1_250 * DOLLAR);
        assertEq(Session.current(open, closed, floor, Session.CLOSING, Session.RAMP_MS / 2), 1_625 * DOLLAR);
        assertEq(Session.current(open, closed, floor, Session.CLOSED, 0), 2_000 * DOLLAR);
        assertEq(Session.current(open, 2_500 * DOLLAR, floor, Session.CLOSED, 0), 2_500 * DOLLAR);
        int256[] memory small = new int256[](2);
        small[0] = -1;
        small[1] = 5;
        assertEq(Requirements.leverageFloor(small), 2);
    }

    function test_NoReturnFallsBelowMinusOneHundredPercent() public {
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](1);
        assets[0] = MarginReference.Asset("X", 1_000_000, 1_000_000, 1, 1, 1, 1, address(0), address(0));
        uint16[] memory none = new uint16[](0);
        MarginReference engine = new MarginReference(
            assets, none, none, bytes32(0), address(0), address(0), address(0), address(0), address(this)
        );
        assertEq(engine.scenario(0, 864_000)[0], -1_000_000);
        assertGt(engine.scenario(255, 864_000)[0], 0);
    }

    function test_TheMeanTickRoundsTowardNegativeInfinity() public {
        ObservationsDouble pool = new ObservationsDouble();
        uint160 perLiquidity = 1 << 20;
        pool.set(0, -1801, 0, perLiquidity);
        (bool ok, int24 tick, uint128 liquidity) = Pools.consult(address(pool));
        assertTrue(ok);
        assertEq(tick, -2);
        assertEq(liquidity, uint256(1_800) * type(uint160).max / (uint256(perLiquidity) << 32));
        pool.set(0, -3600, 0, perLiquidity);
        (, tick,) = Pools.consult(address(pool));
        assertEq(tick, -2);
        pool.set(0, 1801, 0, perLiquidity);
        (, tick,) = Pools.consult(address(pool));
        assertEq(tick, 1);
    }

    function test_CumulativesWrapAsThePoolsDo() public {
        ObservationsDouble pool = new ObservationsDouble();
        uint160 span = uint160(1 << 100) * 1_800;
        pool.set(0, 1_800 * 42, 0, span);
        (, int24 tick, uint128 liquidity) = Pools.consult(address(pool));
        uint160 then = type(uint160).max - 5;
        uint160 now_;
        unchecked {
            now_ = then + span;
        }
        pool.set(0, 1_800 * 42, then, now_);
        (bool ok, int24 wrappedTick, uint128 wrappedLiquidity) = Pools.consult(address(pool));
        assertTrue(ok);
        assertEq(wrappedTick, tick);
        assertEq(wrappedLiquidity, liquidity);
        int56 top = type(int56).max;
        int56 wrapped;
        unchecked {
            wrapped = top + 1_800 * 42;
        }
        pool.set(top, wrapped, 0, span);
        (ok, wrappedTick, wrappedLiquidity) = Pools.consult(address(pool));
        assertTrue(ok);
        assertEq(wrappedTick, tick);
    }

    function test_AMeanTickOutsideUniswapsRangeIsNoAnswer() public {
        ObservationsDouble pool = new ObservationsDouble();
        uint160 span = uint160(1 << 100) * 1_800;
        pool.set(0, 887_272 * 1_800, 0, span);
        (bool ok,,) = Pools.consult(address(pool));
        assertTrue(ok);
        pool.set(0, 887_273 * 1_800, 0, span);
        (ok,,) = Pools.consult(address(pool));
        assertFalse(ok);
        pool.set(0, -887_273 * 1_800, 0, span);
        (ok,,) = Pools.consult(address(pool));
        assertFalse(ok);
    }

    function test_APoolEmptiedForAMomentLowersTheDepthByHalfAtMost() public pure {
        Pool memory emptied = Pool(Requirements.READ, 100e8, 0, 1_799e18, 0);
        uint256 value = 100_000e18;
        assertEq(
            Requirements._positionAddon(int256(value), 100e8, 1_000_000, emptied, 50_000),
            Requirements.addon(value, true, 100e8, 100e8, 500_000e18, 0, 50_000)
        );
        assertEq(
            Requirements._positionAddon(-int256(value), 100e8, 2_000_000, emptied, 50_000),
            Requirements.addon(value, false, 100e8, 100e8, 1_000_000e18, 0, 50_000)
        );
    }

    function test_TheBufferRoundsUpAndAPassedBoundaryReadsAsClosed() public pure {
        assertEq(Session.current(1, 0, 0, Session.OPEN, 0), 2);
        assertEq(Session.current(4, 0, 0, Session.OPEN, 0), 5);
        assertEq(Session.current(5, 0, 0, Session.CLOSED, 0), 7);
        uint256 boundary = 1_789_776_000_000;
        (uint8 code, uint256 elapsed) = Session.regime(true, 2, boundary, boundary - 1);
        assertEq(code, Session.CLOSING);
        assertEq(elapsed, Session.RAMP_MS - 1);
        (code,) = Session.regime(true, 2, boundary, boundary);
        assertEq(code, Session.CLOSED);
        (code,) = Session.regime(true, 2, boundary, boundary + 1);
        assertEq(code, Session.CLOSED);
    }

    function test_ABandCallStarvedOfGasReverts() public {
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](1);
        assets[0] = MarginReference.Asset("X", 10_000, 10_000, 50_000, 50_000, 1, 1, address(0), address(0));
        uint16[] memory none = new uint16[](0);
        MarginReference engine = new MarginReference(
            assets, none, none, bytes32(0), address(0), address(0), address(0), address(new GasBurner()), address(this)
        );
        int256[] memory quantities = new int256[](1);
        quantities[0] = 1e18;
        uint256[] memory prices = new uint256[](1);
        prices[0] = 100e8;
        vm.expectRevert(MarginReference.InsufficientGas.selector);
        engine.currentRequirement{gas: 30_000_000}(quantities, prices);
        MarginReference silent = new MarginReference(
            assets, none, none, bytes32(0), address(0), address(0), address(0), address(this), address(this)
        );
        (,, uint8 code) = silent.currentRequirement(quantities, prices);
        assertEq(code, Session.UNKNOWN);
    }

    function test_AWethPoolNeedsALiveEtherPrice() public {
        vm.warp(1_790_000_000);
        StubToken spy = new StubToken(18);
        StubToken weth = new StubToken(18);
        StubPool pool = new StubPool(address(spy), address(weth), 500, -12_513, 16_029_297_629_534_329_325_587);
        StubAggregator feed = new StubAggregator(8, 268_330_550_000, block.timestamp, "ETH / USD");
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](1);
        assets[0] =
            MarginReference.Asset("SPY", 11_335, 11_335, 54_840, 54_840, 320_861, 877_772, address(pool), address(spy));
        uint16[] memory none = new uint16[](0);
        MarginReference engine = new MarginReference(
            assets,
            none,
            none,
            "SPY",
            address(new StubToken(6)),
            address(weth),
            address(feed),
            address(0),
            address(this)
        );
        int256[] memory quantities = new int256[](1);
        quantities[0] = -130e18;
        uint256[] memory prices = new uint256[](1);
        prices[0] = 76_779_000_000;
        (uint256 live, uint8 missing) = engine.requirement(quantities, prices, 172_800, false);
        assertEq(missing, 0);
        feed.setRound(268_330_550_000, block.timestamp - 86_461);
        (uint256 stale, uint8 bit) = engine.requirement(quantities, prices, 172_800, false);
        assertEq(bit, 1);
        assertLe(stale, live);
        feed.setRound(0, block.timestamp);
        (, bit) = engine.requirement(quantities, prices, 172_800, false);
        assertEq(bit, 1);
        feed.setRound(268_330_550_000, block.timestamp + 1);
        (, bit) = engine.requirement(quantities, prices, 172_800, false);
        assertEq(bit, 1);
    }

    function testFuzz_NoScheduledChangeOfRegimeRaisesTheRequirementAtOnce(
        uint96 open,
        uint96 extra,
        uint128 floor,
        uint64 elapsed
    ) public pure {
        elapsed = uint64(bound(elapsed, 0, Session.RAMP_MS - 1));
        uint256 closed = uint256(open) + extra;
        assertEq(
            Session.current(open, closed, floor, Session.CLOSING, 0),
            Session.current(open, closed, floor, Session.OPEN, 0)
        );
        assertEq(
            Session.current(open, closed, floor, Session.CLOSING, Session.RAMP_MS),
            Session.current(open, closed, floor, Session.CLOSED, 0)
        );
        assertGe(
            Session.current(open, closed, floor, Session.CLOSING, elapsed + 1),
            Session.current(open, closed, floor, Session.CLOSING, elapsed)
        );
        assertEq(
            Session.current(open, closed, floor, Session.UNKNOWN, 0),
            Session.current(open, closed, floor, Session.CLOSED, 0)
        );
    }

    function _deploy(string memory json, string memory config) internal returns (MarginReference) {
        string memory path = string.concat(".configs['", config, "']");
        string[] memory names = vm.parseJsonStringArray(json, string.concat(path, ".names"));
        uint256[] memory volatilities = vm.parseJsonUintArray(json, string.concat(path, ".volatilities"));
        uint256[] memory gaps = vm.parseJsonUintArray(json, string.concat(path, ".gaps"));
        uint256[] memory pairs = vm.parseJsonUintArray(json, string.concat(path, ".correlations"));
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](names.length);
        for (uint256 i; i < names.length; ++i) {
            assets[i] = MarginReference.Asset(
                bytes32(bytes(names[i])),
                uint32(volatilities[i]),
                uint32(volatilities[i]),
                uint32(gaps[i]),
                uint32(gaps[i]),
                1,
                1,
                address(0),
                address(0)
            );
        }
        uint16[] memory correlations = new uint16[](pairs.length);
        for (uint256 k; k < pairs.length; ++k) {
            correlations[k] = uint16(pairs[k]);
        }
        bytes memory market = vm.parseJson(json, string.concat(path, ".market"));
        bytes32 symbol = market.length == 0 ? bytes32(0) : bytes32(bytes(abi.decode(market, (string))));
        return new MarginReference(
            assets, correlations, correlations, symbol, address(0), address(0), address(0), address(0), address(this)
        );
    }

    function _expectedRegime(string memory name) internal pure returns (uint8) {
        string[8] memory closing = [
            "Friday, regular hours",
            "Friday, the close passed on the last regular package",
            "Friday, post-market after the close",
            "Friday, last second of post-market",
            "synthetic: daylight saving ends, post-market",
            "synthetic: daylight saving starts, post-market",
            "synthetic: early close, end of the late session",
            "synthetic: evening before a holiday signed as a short close, post-market"
        ];
        string[12] memory closed = [
            "Saturday, the 24/5 close",
            "Sunday, a second before the reopen",
            "Friday post-market without the Friday close",
            "weekend on a band that never saw a close",
            "Labor Day, a second before the reopen",
            "long close followed by a long close",
            "synthetic: daylight saving ends, a second before the reopen",
            "synthetic: daylight saving starts, the 24/5 close",
            "synthetic: Thanksgiving, a second before the reopen",
            "synthetic: early close, after the late session",
            "synthetic: evening before a holiday signed as a short close, after post-market",
            "synthetic: that short close passed into the holiday"
        ];
        string[6] memory unknown = [
            "status an hour and a second old",
            "points from different packages",
            "change time from another package",
            "unknown status value",
            "change time zero",
            "change time beyond u64"
        ];
        bytes32 id = keccak256(bytes(name));
        for (uint256 k; k < closing.length; ++k) {
            if (id == keccak256(bytes(closing[k]))) return Session.CLOSING;
        }
        for (uint256 k; k < closed.length; ++k) {
            if (id == keccak256(bytes(closed[k]))) return Session.CLOSED;
        }
        for (uint256 k; k < unknown.length; ++k) {
            if (id == keccak256(bytes(unknown[k]))) return Session.UNKNOWN;
        }
        return Session.OPEN;
    }
}
