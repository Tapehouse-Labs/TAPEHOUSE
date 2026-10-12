// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {GapCover} from "../src/GapCover.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract GapCoverTest is Test {
    uint256 internal constant USDG = 1e6;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant TSLA = "TSLA";
    uint32 internal constant NVDA_GAP = 118_601;
    uint32 internal constant SPY_GAP = 54_840;
    uint32 internal constant NVDA_FLOOR = 65_230;
    uint256 internal constant FRIDAY = 1_790_949_600;
    uint64 internal constant REGULAR_CLOSE_MS = 1_790_971_200_000;
    uint64 internal constant CLOSE_MS = 1_790_985_600_000;
    uint256 internal constant SATURDAY = 1_791_028_800;
    uint64 internal constant REOPEN_MS = 1_791_158_400_000;
    uint256 internal constant REOPEN = 1_791_158_400;
    uint64 internal constant REGULAR_OPEN_MS = 1_791_207_000_000;
    uint80 internal constant FIRST = uint80(1) << 64 | 1;
    uint256 internal constant NOTIONAL = 10_000 * USDG;
    uint256 internal constant PREMIUM = 1_547_591;
    uint256 internal constant DEPOSIT = 100_000 * USDG;

    StubUsdg internal usdg;
    BandDouble internal band;
    MarginDouble internal engine;
    StubAggregator internal nvdaFeed;
    StubAggregator internal spyFeed;
    GapCover internal cover;
    address internal writer = makeAddr("writer");
    address internal other = makeAddr("other");
    address internal buyer = makeAddr("buyer");
    address internal holder = makeAddr("holder");
    address internal keeper = makeAddr("keeper");

    function setUp() public {
        vm.warp(FRIDAY);
        usdg = new StubUsdg();
        band = new BandDouble();
        nvdaFeed = new StubAggregator(8, 200e8, FRIDAY - 9 days, "NVDA / USD");
        spyFeed = new StubAggregator(8, 600e8, FRIDAY - 9 days, "SPY / USD");
        band.setAsset(NVDA, BandDouble.Asset(address(nvdaFeed), "NVDA---24_7", bytes32(0), address(0)));
        band.setAsset(SPY, BandDouble.Asset(address(spyFeed), bytes32(0), "USA500.Y---24_7", address(0)));
        band.setAsset(TSLA, BandDouble.Asset(address(0), "TSLA---24_7", bytes32(0), address(0)));
        _quote(NVDA, 200e8, 3);
        _quote(SPY, 600e8, 3);
        _quote(TSLA, 400e8, 3);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS, CLOSE_MS));
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = (NVDA, SPY, TSLA);
        engine = new MarginDouble(symbols, address(band), address(0));
        engine.setWeekendGap(NVDA, NVDA_GAP);
        engine.setWeekendGap(SPY, SPY_GAP);
        engine.setWeekendGap(TSLA, 135_345);
        cover = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        _deposit(writer, DEPOSIT);
        usdg.mint(buyer, 1_000_000 * USDG);
        vm.prank(buyer);
        usdg.approve(address(cover), type(uint256).max);
    }

    function test_TheCoverTakesTheEnginesBandAndItsAssetsChainlinkFeeds() public view {
        assertEq(address(cover.engine()), address(engine));
        assertEq(address(cover.band()), address(band));
        assertFalse(cover.regularHours());
        assertEq(cover.asset(), address(usdg));
        assertEq(cover.name(), "Tapehouse Gap Cover");
        assertEq(cover.symbol(), "thCOVER");
        assertEq(cover.decimals(), 12);
        assertEq(cover.feed(NVDA), address(nvdaFeed));
        assertEq(cover.feed(SPY), address(spyFeed));
        assertEq(cover.feed(TSLA), address(0));
        assertEq(cover.held(), DEPOSIT);
        assertEq(cover.totalAssets(), DEPOSIT);
        assertEq(cover.capacity(), DEPOSIT);
    }

    function test_ThePremiumIsTheLayersExpectedPayoutUnderTheFittedTailWithItsLoading() public {
        assertEq(cover.minDeductible(NVDA), 120);
        assertEq(cover.minDeductible(SPY), 56);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        assertEq(cover.quote(SPY, NOTIONAL, 300, 1000), 581_143);
        assertEq(cover.quote(NVDA, 1_000 * USDG, 120, 10_000), 926_836);
        assertEq(cover.quote(NVDA, 1, 120, 121), 1);
        band.setSession(BandDouble.Session(2, 2, 1, REGULAR_CLOSE_MS - 1 days, 0));
        assertEq(cover.minDeductible(NVDA), 218);
        assertEq(cover.minDeductible(SPY), 101);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), 5_056_153);
        assertEq(cover.quote(SPY, NOTIONAL, 300, 1000), 1_904_149);
        assertEq(cover.quote(NVDA, 1_000 * USDG, 218, 10_000), 1_669_052);
    }

    function test_EveryQuoteMatchesTheReferenceToTheUnit() public {
        string memory json = vm.readFile("test/data/gap-cover-quotes.json");
        band.setSession(BandDouble.Session(2, 2, 1, REGULAR_CLOSE_MS - 1 days, 0));
        uint256 count;
        for (; vm.keyExistsJson(json, string.concat(".quotes[", vm.toString(count), "]")); ++count) {
            string[] memory q = vm.parseJsonStringArray(json, string.concat(".quotes[", vm.toString(count), "]"));
            engine.setWeekendGap(NVDA, uint32(vm.parseUint(q[0])));
            nvdaFeed.setRound(int256(vm.parseUint(q[1])), FRIDAY);
            _quote(NVDA, uint64(vm.parseUint(q[2])), 3);
            assertEq(cover.minDeductible(NVDA), vm.parseUint(q[6]), q[0]);
            assertEq(
                cover.quote(NVDA, vm.parseUint(q[3]), vm.parseUint(q[4]), vm.parseUint(q[5])), vm.parseUint(q[7]), q[3]
            );
        }
        assertEq(count, 256);
    }

    function testFuzz_APremiumIsAtMostTheReserveAndRisesWithTheLayerAndTheGap(
        uint32 gap,
        uint128 notional,
        uint16 deductible,
        uint16 limit,
        uint64 drop
    ) public {
        gap = uint32(bound(gap, 1_000, 1_000_000));
        engine.setWeekendGap(NVDA, gap);
        _quote(NVDA, uint64(200e8 - bound(drop, 0, 2e8)), 3);
        uint256 minimum = cover.minDeductible(NVDA);
        vm.assume(minimum < 9_999);
        notional = uint128(bound(notional, 1, 1e30));
        deductible = uint16(bound(deductible, minimum, 9_998));
        limit = uint16(bound(limit, deductible + 1, 9_999));
        uint256 premium = cover.quote(NVDA, notional, deductible, limit);
        assertLe(premium, (uint256(notional) * (limit - deductible) + 9_999) / 10_000);
        assertGe(cover.quote(NVDA, notional, deductible, limit + 1), premium);
        if (deductible > minimum) assertGe(cover.quote(NVDA, notional, deductible - 1, limit), premium);
        engine.setWeekendGap(NVDA, uint32(uint256(gap) * 3 / 2));
        if (cover.minDeductible(NVDA) <= deductible) assertGe(cover.quote(NVDA, notional, deductible, limit), premium);
    }

    function test_AReferenceAboveTheBandsCentrePricesTheLayerFromTheCentre() public {
        _quote(NVDA, 199e8, 3);
        assertEq(cover.minDeductible(NVDA), 170);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), 1_776_481);
        assertEq(cover.quote(NVDA, NOTIONAL, 170, 1170), 8_339_677);
        vm.expectRevert(abi.encodeWithSelector(GapCover.StaleReference.selector, 169, 170));
        _buy(NVDA, NOTIONAL, 169, 1500);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLayer.selector, 119, 1500, 120));
        _buy(NVDA, NOTIONAL, 119, 1500);
        (, uint256 premium) = _buy(NVDA, NOTIONAL, 500, 1500);
        assertEq(premium, 1_776_481);
        _quote(NVDA, 0, 0);
        assertEq(cover.minDeductible(NVDA), 120);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        _quote(NVDA, 201e8, 3);
        assertEq(cover.minDeductible(NVDA), 120);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        _quote(NVDA, 200e8, 3);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        nvdaFeed.setRound(1e12, FRIDAY);
        _quote(NVDA, 994_955_876_240, 3);
        assertEq(cover.minDeductible(NVDA), 170);
        assertEq(cover.quote(NVDA, NOTIONAL, 170, 1170), 8_372_390);
        _quote(NVDA, 200e8, 3);
        nvdaFeed.setRound(0, FRIDAY);
        assertEq(cover.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        nvdaFeed.setRound(-1, FRIDAY);
        vm.expectRevert(abi.encodeWithSelector(SafeCast.SafeCastOverflowedIntToUint.selector, -1));
        cover.minDeductible(NVDA);
    }

    function test_TheGapFollowsTheRealisedMoveOfTheTradingWeekBeforeTheClose() public {
        (GapCover target,) = _coverOverWeek([int256(200e8), 194e8, 206e8, 198e8, 210e8, 200e8, 204e8, 200e8]);
        _assertGap(target, 223_890, 111_945);
        assertEq(target.minDeductible(NVDA), 411);
        assertEq(target.quote(NVDA, NOTIONAL, 600, 1600), 13_850_923);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLayer.selector, 410, 1500, 411));
        target.quote(NVDA, NOTIONAL, 410, 1500);
        engine.setWeekendGap(NVDA, 407_073);
        assertEq(cover.quote(NVDA, NOTIONAL, 600, 1600), 13_850_923);
        engine.setWeekendGap(NVDA, 500_000);
        _assertGap(target, 275_000, 111_945);
        engine.setWeekendGap(NVDA, NVDA_GAP);
        vm.prank(buyer);
        (, uint256 premium) = target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        assertEq(premium, 13_850_923);
    }

    function test_ACalmWeekPricesBelowTheEnginesGapAndEightTimesItBoundsATurbulentOne() public {
        (GapCover calm,) = _coverOverWeek([int256(200e8), 201e8, 200e8, 199e8, 200e8, 201e8, 200e8, 200e8]);
        _assertGap(calm, NVDA_FLOOR, 12_237);
        assertEq(calm.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        (GapCover moderate,) = _coverOverWeek([int256(200e8), 197e8, 201e8, 199e8, 202e8, 200e8, 198e8, 200e8]);
        _assertGap(moderate, 70_752, 35_376);
        assertEq(moderate.minDeductible(NVDA), 130);
        assertEq(moderate.quote(NVDA, NOTIONAL, 500, 1500), 1_818_476);
        (GapCover turbulent,) = _coverOverWeek([int256(200e8), 150e8, 200e8, 120e8, 200e8, 100e8, 200e8, 200e8]);
        _assertGap(turbulent, 948_808, 1_424_097);
        assertEq(turbulent.minDeductible(NVDA), 1740);
        assertEq(turbulent.quote(NVDA, NOTIONAL, 1740, 2740), 48_566_316);
        _assertGap(cover, NVDA_FLOOR, 0);
    }

    function test_AWeekTheFeedCannotShowPricesAtTheCap() public {
        (GapCover gone,) = _coverOverWeek([int256(200e8), 201e8, 200e8, 199e8, 200e8, 201e8, 200e8, 0]);
        _assertGap(gone, 8 * NVDA_GAP, type(uint64).max);
        (GapCover negative,) = _coverOverWeek([int256(200e8), 201e8, -1, 199e8, 200e8, 201e8, 200e8, 200e8]);
        _assertGap(negative, 8 * NVDA_GAP, type(uint64).max);
        (GapCover zero,) = _coverOverWeek([int256(0), 201e8, 200e8, 199e8, 200e8, 201e8, 200e8, 200e8]);
        _assertGap(zero, 8 * NVDA_GAP, type(uint64).max);
        int256 huge = int256(2 ** 250);
        (GapCover wild,) = _coverOverWeek([huge, 1, huge, 1, huge, 1, huge, 1]);
        _assertGap(wild, 8 * NVDA_GAP, type(uint64).max);
        StubAggregator late = new StubAggregator(8, 200e8, CLOSE_MS / 1000 - 8 days, "NVDA / USD");
        _assertGap(_coverOver(address(late)), 8 * NVDA_GAP, type(uint64).max);
        StubAggregator young = new StubAggregator(8, 200e8, CLOSE_MS / 1000 - 8 days - 1, "NVDA / USD");
        _assertGap(_coverOver(address(young)), NVDA_FLOOR, 0);
        PhasedAggregator upgraded = new PhasedAggregator();
        upgraded.add(2, 200e8, FRIDAY - 3 days, FRIDAY - 3 days);
        GapCover blind = _coverOver(address(upgraded));
        _assertGap(blind, 8 * NVDA_GAP, type(uint64).max);
        assertEq(blind.minDeductible(NVDA), 1740);
    }

    function test_TheWeekIsFoundAcrossTheFeedsPhases() public {
        uint256 closeS = CLOSE_MS / 1000;
        int256[8] memory week = [int256(200e8), 194e8, 206e8, 198e8, 210e8, 200e8, 204e8, 200e8];
        PhasedAggregator phased = new PhasedAggregator();
        phased.add(1, 300e8, closeS - 20 days, closeS - 20 days);
        for (uint256 k = 8; k > 3; --k) {
            phased.add(1, 300e8, closeS - k * 1 days - 3 hours, closeS - k * 1 days - 3 hours);
            phased.add(1, week[k - 1], closeS - k * 1 days - 1 hours, closeS - k * 1 days - 1 hours);
        }
        phased.add(2, week[2], closeS - 3 days - 1 hours, closeS - 3 days - 1 hours);
        phased.add(1, 100e8, closeS - 3 days - 30 minutes, closeS - 3 days - 30 minutes);
        phased.add(2, week[1], closeS - 2 days - 1 hours, closeS - 2 days - 1 hours);
        phased.add(1, 100e8, closeS - 2 days - 30 minutes, closeS - 2 days - 30 minutes);
        phased.add(4, week[0], closeS - 1 days - 1 hours, closeS - 1 days - 1 hours);
        phased.add(4, 200e8, FRIDAY - 1 hours, FRIDAY - 1 hours);
        GapCover target = _coverOver(address(phased));
        _assertGap(target, 223_890, 111_945);
        vm.prank(buyer);
        (, uint256 premium) = target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        assertEq(premium, 13_850_923);
    }

    function test_TheWeekIsKeptOnlyOnceItsLastDayHasPassed() public {
        uint256 closeS = CLOSE_MS / 1000;
        StubAggregator feed = new StubAggregator(8, 200e8, closeS - 8 days - 1 hours, "NVDA / USD");
        int256[6] memory week = [int256(194e8), 206e8, 198e8, 210e8, 200e8, 204e8];
        for (uint256 k = 7; k > 1; --k) {
            feed.setRound(week[k - 2], closeS - k * 1 days - 1 hours);
        }
        GapCover target = _coverOver(address(feed));
        vm.warp(closeS - 1 days - 2 hours);
        _assertGap(target, 215_178, 107_589);
        vm.expectRevert(abi.encodeWithSelector(GapCover.TooEarlyToMeasure.selector, CLOSE_MS - 1 days * 1000));
        target.measure(NVDA);
        vm.prank(buyer);
        target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        vm.warp(closeS - 1 days - 1 hours);
        feed.setRound(200e8, closeS - 1 days - 1 hours);
        _assertGap(target, 223_890, 111_945);
        vm.warp(closeS - 1 days);
        vm.prank(buyer);
        target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        feed.setRound(150e8, closeS - 1 days - 1);
        _assertGap(target, 223_890, 111_945);
    }

    function test_AnyoneMayMeasureTheWeekBeforeTheFirstBuyer() public {
        (GapCover target, StubAggregator feed) =
            _coverOverWeek([int256(200e8), 194e8, 206e8, 198e8, 210e8, 200e8, 204e8, 200e8]);
        vm.warp(CLOSE_MS / 1000 - 1 days);
        vm.expectEmit();
        emit GapCover.Measured(NVDA, CLOSE_MS, 111_945);
        vm.prank(keeper);
        assertEq(target.measure(NVDA), 111_945);
        feed.setRound(150e8, CLOSE_MS / 1000 - 1 days - 1);
        _assertGap(target, 223_890, 111_945);
        vm.recordLogs();
        vm.prank(keeper);
        assertEq(target.measure(NVDA), 111_945);
        vm.prank(buyer);
        (, uint256 premium) = target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        assertEq(premium, 13_850_923);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        for (uint256 i; i < logs.length; ++i) {
            assertTrue(logs[i].topics[0] != GapCover.Measured.selector);
        }
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotCoverable.selector, TSLA));
        target.measure(TSLA);
        vm.warp(CLOSE_MS / 1000);
        vm.expectRevert(GapCover.SalesClosed.selector);
        target.measure(NVDA);
    }

    function test_AFeedReadThatRunsOutOfGasReverts() public {
        GapCover target = _coverOver(address(new GreedyAggregator()));
        vm.expectRevert(GapCover.InsufficientGas.selector);
        target.pricingGap(NVDA);
    }

    function test_NoCallerChoosesTheRoundsOfTheWeek() public {
        (GapCover target, StubAggregator feed) =
            _coverOverWeek([int256(200e8), 194e8, 206e8, 198e8, 210e8, 200e8, 204e8, 200e8]);
        feed.setRound(150e8, FRIDAY);
        _quote(NVDA, 150e8, 3);
        _assertGap(target, 223_890, 111_945);
        vm.prank(buyer);
        target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        vm.warp(CLOSE_MS / 1000 - 1 hours);
        feed.setRound(300e8, CLOSE_MS / 1000 - 1 hours);
        _quote(NVDA, 300e8, 3);
        _assertGap(target, 223_890, 111_945);
        uint256 premium = target.quote(NVDA, NOTIONAL, 600, 1600);
        vm.prank(other);
        usdg.approve(address(target), premium);
        usdg.mint(other, premium);
        vm.prank(other);
        (, uint256 paid) = target.buy(NVDA, NOTIONAL, 600, 1600, premium, other);
        assertEq(paid, 13_850_923);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS + 7 days * 1000, CLOSE_MS + 7 days * 1000));
        vm.warp(FRIDAY + 7 days);
        _assertGap(target, 8 * NVDA_GAP, 500_000);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        _assertGap(target, NVDA_GAP, 0);
    }

    function testFuzz_TheWeekIsTheLastRoundBeforeEachOfItsDays(uint256 seed, uint8 perDay, uint8 days_) public {
        uint256 closeS = CLOSE_MS / 1000;
        uint256 t = closeS - bound(days_, 1, 20) * 1 days;
        uint256 spacing = 2 days / bound(perDay, 1, 40);
        StubAggregator feed = new StubAggregator(8, 200e8, t, "NVDA / USD");
        uint256[] memory times = new uint256[](800);
        int256[] memory answers = new int256[](800);
        (times[0], answers[0]) = (t, 200e8);
        uint256 n = 1;
        for (; n < 800; ++n) {
            seed = uint256(keccak256(abi.encode(seed)));
            t += 1 + seed % spacing;
            if (t >= closeS - 1 days) break;
            answers[n] = int256(answers[n - 1] * int256(9_700 + seed % 600) / 10_000);
            times[n] = t;
            feed.setRound(answers[n], t);
        }
        feed.setRound(200e8, FRIDAY - 1 hours);
        uint256 sum;
        uint256 later;
        bool readable = true;
        for (uint256 day = 1; day <= 8; ++day) {
            uint256 price;
            for (uint256 i; i < n && times[i] < closeS - day * 1 days; ++i) {
                price = uint256(answers[i]);
            }
            if (price == 0) readable = false;
            if (day > 1 && readable) {
                uint256 change = (later > price ? later - price : price - later) * 1e6 / price;
                sum += change * change;
            }
            later = price;
        }
        _assertGap(
            _coverOver(address(feed)),
            readable ? Math.min(Math.max(NVDA_FLOOR, Math.sqrt(sum) * 2), 8 * NVDA_GAP) : 8 * NVDA_GAP,
            readable ? Math.sqrt(sum) : type(uint64).max
        );
    }

    function test_TheWeekReadsOnlyAnsweredRoundsStartedBeforeEachDay() public {
        uint256 closeS = CLOSE_MS / 1000;
        StubAggregator feed = new StubAggregator(8, 200e8, closeS - 9 days, "NVDA / USD");
        for (uint256 k = 8; k > 0; --k) {
            uint256 day = closeS - k * 1 days;
            for (uint256 h = 23; h > 1; --h) {
                feed.setRound(200e8, day - h * 1 hours);
            }
            feed.setRound(int256(100e8 + k * 10e8), day - 30 minutes, 0);
            feed.setRound(int256(300e8 + k * 10e8), day);
            for (uint256 h = 1; h < 3; ++h) {
                feed.setRound(200e8, day + h * 1 hours);
            }
        }
        feed.setRound(200e8, FRIDAY - 1 hours);
        _assertGap(_coverOver(address(feed)), NVDA_FLOOR, 0);
        StubAggregator pending = new StubAggregator(8, 200e8, closeS - 9 days, "NVDA / USD");
        pending.setRound(200e8, closeS - 2 days);
        pending.setRound(100e8, closeS - 1 days - 1 hours, 0);
        _assertGap(_coverOver(address(pending)), NVDA_FLOOR, 0);
    }

    function test_TheSeriesKeepsTheWeekItWasSoldOnThroughAPhaseChange() public {
        PhasedAggregator phased = new PhasedAggregator();
        phased.add(1, 200e8, FRIDAY - 9 days, FRIDAY - 9 days);
        GapCover target = _coverOver(address(phased));
        vm.prank(buyer);
        (, uint256 premium) = target.buy(NVDA, NOTIONAL, 500, 1500, type(uint256).max, holder);
        assertEq(premium, PREMIUM);
        phased.add(2, 200e8, FRIDAY, FRIDAY);
        _assertGap(target, NVDA_FLOOR, 0);
        assertEq(target.quote(NVDA, NOTIONAL, 500, 1500), PREMIUM);
        _assertGap(_coverOver(address(phased)), NVDA_FLOOR, 0);
    }

    function test_EveryWeeksMoveMatchesTheReferenceToTheUnit() public {
        string memory json = vm.readFile("test/data/gap-cover-weeks.json");
        uint256 count;
        for (; vm.keyExistsJson(json, string.concat(".weeks[", vm.toString(count), "]")); ++count) {
            string[] memory w = vm.parseJsonStringArray(json, string.concat(".weeks[", vm.toString(count), "]"));
            uint256 closeS = CLOSE_MS / 1000 + (count + 1) * 8 days;
            for (uint256 k = 8; k > 0; --k) {
                nvdaFeed.setRound(vm.parseInt(w[k]), closeS - k * 1 days - 1 hours);
            }
            nvdaFeed.setRound(200e8, closeS - 11 hours);
            engine.setWeekendGap(NVDA, uint32(vm.parseUint(w[0])));
            band.setSession(BandDouble.Session(2, 1, 3, uint64(closeS - 4 hours) * 1000, uint64(closeS) * 1000));
            vm.warp(closeS - 10 hours);
            (uint256 gap, uint256 move) = cover.pricingGap(NVDA);
            assertEq(move, vm.parseUint(w[9]), w[0]);
            assertEq(gap, vm.parseUint(w[10]), w[0]);
        }
        assertEq(count, 256);
    }

    function test_ACoverIsSoldOnlyFromTheCloseTheSessionAnnouncesUntilItAndReservesItsWholePayout() public {
        band.setSession(BandDouble.Session(2, 2, 1, REGULAR_CLOSE_MS - 1 days, 0));
        vm.expectRevert(GapCover.SalesClosed.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS, CLOSE_MS));
        vm.expectEmit();
        emit GapCover.CloseRecorded(CLOSE_MS, CLOSE_MS);
        vm.expectEmit();
        emit GapCover.Bought(1, holder, NVDA, CLOSE_MS, NOTIONAL, 500, 1500, PREMIUM);
        (uint256 id, uint256 premium) = _buy(NVDA, NOTIONAL, 500, 1500);
        assertEq(id, 1);
        assertEq(premium, PREMIUM);
        assertEq(cover.lastCloseMs(), CLOSE_MS);
        assertEq(cover.coverCount(), 1);
        (address h, uint64 closes, uint16 deductible, uint16 limit, bytes32 symbol, uint128 notional, uint128 paid) =
            cover.covers(1);
        assertEq(h, holder);
        assertEq(closes, CLOSE_MS);
        assertEq(deductible, 500);
        assertEq(limit, 1500);
        assertEq(symbol, NVDA);
        assertEq(notional, NOTIONAL);
        assertEq(paid, PREMIUM);
        assertEq(cover.reserved(), 1_000 * USDG);
        assertEq(cover.premiums(), PREMIUM);
        assertEq(cover.held(), DEPOSIT);
        assertEq(cover.capacity(), DEPOSIT - 1_000 * USDG);
        assertEq(cover.outstanding(), 1);
        assertEq(cover.outstandingIn(CLOSE_MS), 1);
        assertEq(usdg.balanceOf(address(cover)), DEPOSIT + PREMIUM);
        (uint128 covered,,,,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(covered, NOTIONAL);

        vm.warp(CLOSE_MS / 1000 - 1);
        band.setSession(BandDouble.Session(2, 3, 1, REOPEN_MS + 48_600_000, CLOSE_MS));
        vm.recordLogs();
        _buy(SPY, NOTIONAL, 300, 1000);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        for (uint256 i; i < logs.length; ++i) {
            assertTrue(logs[i].topics[0] != GapCover.CloseRecorded.selector);
        }
        vm.warp(CLOSE_MS / 1000);
        vm.expectRevert(GapCover.SalesClosed.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        vm.expectRevert(GapCover.SalesClosed.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
    }

    function test_SalesShowTheClosureOnSaleAndWhenTheyEnd() public {
        _assertSales(cover, CLOSE_MS, CLOSE_MS);
        band.setSession(BandDouble.Session(2, 2, 1, REGULAR_CLOSE_MS - 1 days, 0));
        _assertSales(cover, 0, 0);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS, CLOSE_MS));
        _sold();
        uint64 revised = CLOSE_MS - 1 hours * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS - 1 hours * 1000, revised));
        _assertSales(cover, CLOSE_MS, revised);
        band.setSequencer(makeAddr("sequencer"), true);
        GapCover regular = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        _assertSales(regular, revised, revised - 4 hours * 1000);
        vm.warp(revised / 1000 - 1);
        _assertSales(cover, CLOSE_MS, revised);
        vm.warp(revised / 1000);
        _assertSales(cover, 0, 0);
        _closed();
        _assertSales(cover, 0, 0);
    }

    function test_ACoverNeedsALayerFromTheTailACoverableAssetAndAHolder() public {
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLayer.selector, 119, 1500, 120));
        _buy(NVDA, NOTIONAL, 119, 1500);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLayer.selector, 500, 500, 120));
        _buy(NVDA, NOTIONAL, 500, 500);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLayer.selector, 500, 10_001, 120));
        _buy(NVDA, NOTIONAL, 500, 10_001);
        vm.expectRevert(GapCover.ZeroAmount.selector);
        _buy(NVDA, 0, 500, 1500);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotCoverable.selector, TSLA));
        _buy(TSLA, NOTIONAL, 500, 1500);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotCoverable.selector, bytes32("AAPL")));
        cover.quote("AAPL", NOTIONAL, 500, 1500);
        vm.expectRevert(GapCover.InvalidHolder.selector);
        vm.prank(buyer);
        cover.buy(NVDA, NOTIONAL, 500, 1500, PREMIUM, address(0));
        vm.expectRevert(abi.encodeWithSelector(GapCover.PremiumAboveLimit.selector, PREMIUM, PREMIUM - 1));
        vm.prank(buyer);
        cover.buy(NVDA, NOTIONAL, 500, 1500, PREMIUM - 1, holder);
        _buy(NVDA, NOTIONAL, 120, 10_000);
    }

    function test_NoCoverIsSoldWhileTheBandCannotVouchForTheAsset() public {
        _quote(NVDA, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(GapCover.AssetHalted.selector, NVDA));
        _buy(NVDA, NOTIONAL, 500, 1500);
        _quote(NVDA, 200e8, 1);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(1, uint64(FRIDAY + 1 hours), 1e18, 2e18));
        vm.expectRevert(abi.encodeWithSelector(GapCover.CorporateActionPending.selector, NVDA));
        _buy(NVDA, NOTIONAL, 500, 1500);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(0, uint64(FRIDAY - 1 days), 1e18, 2e18));
        band.setSequencer(makeAddr("sequencer"), false);
        vm.expectRevert(GapCover.SequencerNotSettled.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
        band.setSequencer(makeAddr("sequencer"), true);
        _buy(NVDA, NOTIONAL, 500, 1500);
    }

    function test_ACoverIsSoldOnlyWithinTheWritersFreeUsdg() public {
        uint256 notional = DEPOSIT * 10_000 / 1_000;
        _buy(NVDA, notional - 1, 500, 1500);
        assertEq(cover.capacity(), 0);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NoCapacity.selector, 1, 0));
        _buy(SPY, 1, 9_999, 10_000);
        assertEq(cover.reserved(), DEPOSIT);
    }

    function test_TheCoverRecordsTheCloseAndTheReopenFromTheSession() public {
        cover.record();
        assertEq(cover.lastCloseMs(), CLOSE_MS);
        assertEq(cover.reopenOf(CLOSE_MS), 0);
        vm.warp(CLOSE_MS / 1000);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        vm.expectEmit();
        emit GapCover.ReopenRecorded(CLOSE_MS, REOPEN_MS);
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS + 1 days, REOPEN_MS + 1 days * 1000));
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS + 1 days * 1000);
        band.setSession(BandDouble.Session(1, 3, 2, 0, 0));
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS + 1 days * 1000);
        vm.warp(REOPEN + 1 days + 1);
        band.setSession(BandDouble.Session(1, 2, 1, REGULAR_OPEN_MS + 2 days, REOPEN_MS + 2 days * 1000));
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS + 1 days * 1000);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        cover.record();
        band.setSession(BandDouble.Session(2, 2, 1, 0, 0));
        cover.record();
        assertEq(cover.lastCloseMs(), CLOSE_MS);
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS + 1 days * 1000);
    }

    function test_AReopenIsRecordedOnlyForAClosureThatHasBegun() public {
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        vm.recordLogs();
        cover.record();
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS, uint64(FRIDAY * 1000)));
        cover.record();
        assertEq(vm.getRecordedLogs().length, 0);
        assertEq(cover.lastCloseMs(), 0);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS, CLOSE_MS));
        cover.record();
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), 0);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, 0));
        vm.warp(SATURDAY);
        cover.record();
        assertEq(cover.reopenOf(CLOSE_MS), 0);
    }

    function test_ARevisedCloseKeepsTheSeriesSoldBeforeItAndSettlesItAgainstTheRevisedClose() public {
        uint256 first = _sold();
        uint64 revised = CLOSE_MS - 1 hours * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS - 1 hours * 1000, revised));
        vm.expectEmit();
        emit GapCover.CloseRecorded(CLOSE_MS, revised);
        vm.expectEmit();
        emit GapCover.Bought(2, holder, NVDA, CLOSE_MS, NOTIONAL, 500, 1500, PREMIUM);
        (uint256 second,) = _buy(NVDA, NOTIONAL, 500, 1500);
        assertEq(cover.lastCloseMs(), CLOSE_MS);
        assertEq(cover.outstandingIn(CLOSE_MS), 2);
        assertEq(cover.maxDeposit(other), type(uint256).max);
        nvdaFeed.setRound(199e8, revised / 1000 - 60, revised / 1000 - 48);
        nvdaFeed.setRound(205e8, revised / 1000 + 60, revised / 1000 + 72);
        vm.warp(revised / 1000);
        assertEq(cover.maxDeposit(other), 0);
        vm.expectRevert(GapCover.SalesClosed.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
        uint256 snapshot = vm.snapshotState();
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS + 7 days * 1000, CLOSE_MS + 7 days * 1000));
        vm.expectEmit();
        emit GapCover.CloseRecorded(CLOSE_MS + 7 days * 1000, CLOSE_MS + 7 days * 1000);
        cover.record();
        vm.revertToState(snapshot);
        _closed();
        assertEq(cover.reopenOf(CLOSE_MS), REOPEN_MS);
        _fill(NVDA, 181e8);
        vm.warp(REOPEN + 15 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST + 2));
        cover.settle(NVDA, CLOSE_MS, FIRST + 2, FIRST + 2);
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 2);
        _assertSettled(NVDA, 181e8, true);
        (, uint64 referencePrice,,,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(referencePrice, 199e8);
        (uint256 payout,) = cover.release(first);
        assertEq(payout, NOTIONAL * (18e8 - uint256(199e8) * 500 / 10_000) / 199e8);
        cover.release(second);
    }

    function test_TheFirstRoundAfterTheReopenInsideTheBandSettlesTheSeries() public {
        uint256 id = _sold();
        nvdaFeed.setRound(199e8, FRIDAY + 5 hours, FRIDAY + 5 hours + 12);
        nvdaFeed.setRound(201e8, CLOSE_MS / 1000 + 100, CLOSE_MS / 1000 + 112);
        _closed();
        _reopened(REOPEN + 5);
        _quote(NVDA, 180e8, 1);
        _observe(NVDA);
        nvdaFeed.setRound(181e8, REOPEN + 20, REOPEN + 32);
        for (uint256 k = 1; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            _quote(NVDA, uint64(181e8 + k * 1e8), 3);
            _observe(NVDA);
        }
        vm.warp(REOPEN + 15 minutes);
        vm.expectEmit();
        emit GapCover.Settled(NVDA, CLOSE_MS, FIRST + 1, 199e8, FIRST + 3, 181e8, false);
        vm.prank(keeper);
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 2);
        (uint128 notional, uint64 referencePrice, uint64 price, uint16 shift, uint8 status, bool flagged) =
            cover.series(NVDA, CLOSE_MS);
        assertEq(notional, NOTIONAL);
        assertEq(referencePrice, 199e8);
        assertEq(price, 181e8);
        assertEq(shift, 0);
        assertEq(status, 1);
        assertFalse(flagged);
        uint256 payout = NOTIONAL * (18e8 - uint256(199e8) * 500 / 10_000) / 199e8;
        vm.expectEmit();
        emit GapCover.Released(id, holder, payout, 0);
        vm.prank(keeper);
        cover.release(id);
        assertEq(cover.payouts(holder), payout);
        assertEq(cover.held(), DEPOSIT + PREMIUM - payout);
        _claimAll(holder, payout);
    }

    function test_ARoundOutsideItsBandOrTheWindowSettlesOnTheMedianBandCentreAndIsFlagged() public {
        _sold();
        _closed();
        nvdaFeed.setRound(170e8, REOPEN + 20, REOPEN + 32);
        uint64[15] memory mids =
            [uint64(181e8), 185e8, 183e8, 184e8, 182e8, 186e8, 188e8, 187e8, 189e8, 190e8, 0, 191e8, 0, 0, 192e8];
        _window(mids);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 186.5e8, true);
    }

    function test_ARoundAboveItsBandIsNotTaken() public {
        _sold();
        _closed();
        nvdaFeed.setRound(182e8 + 1, REOPEN + 20, REOPEN + 32);
        uint64[15] memory mids = [uint64(181e8), 181e8, 181e8, 181e8, 181e8, 181e8, 181e8, 181e8, 0, 0, 0, 0, 0, 0, 0];
        _window(mids);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 181e8, true);
    }

    function test_AReferenceWithoutAPositiveAnswerIsRefused() public {
        _sold();
        nvdaFeed.setRound(0, FRIDAY + 1 hours, FRIDAY + 1 hours);
        _closed();
        _fill(NVDA, 181e8);
        vm.warp(REOPEN + 15 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST + 1));
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 1);
    }

    function test_ARoundThatLandsAfterTheWindowSettlesOnTheMedianOfAnEvenCount() public {
        _sold();
        _closed();
        uint64[15] memory mids = [uint64(181e8), 185e8, 183e8, 184e8, 182e8, 186e8, 187e8, 180e8, 0, 0, 0, 0, 0, 0, 0];
        _window(mids);
        nvdaFeed.setRound(181e8, REOPEN + 14 minutes + 59, REOPEN + 15 minutes);
        vm.warp(REOPEN + 15 minutes);
        vm.expectEmit();
        emit GapCover.Settled(NVDA, CLOSE_MS, FIRST, 200e8, 0, 183.5e8, true);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
    }

    function test_ARoundInAnUnobservedSlotOrWithoutAPositiveAnswerIsNotTaken() public {
        _sold();
        _closed();
        uint64[15] memory mids = [uint64(181e8), 0, 183e8, 184e8, 185e8, 186e8, 187e8, 188e8, 189e8, 0, 0, 0, 0, 0, 0];
        _window(mids);
        nvdaFeed.setRound(181e8, REOPEN + 61, REOPEN + 70);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 185.5e8, true);
    }

    function test_ARoundStartedAfterTheReopenButStampedBeforeItIsNotTaken() public {
        _sold();
        _closed();
        _fill(NVDA, 181e8);
        nvdaFeed.setRound(181e8, REOPEN, REOPEN - 1);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 181e8, true);
    }

    function test_ANonPositiveFirstAnswerSettlesOnTheMedian() public {
        _sold();
        _closed();
        uint64[15] memory mids = [uint64(181e8), 182e8, 181e8, 182e8, 181e8, 182e8, 181e8, 182e8, 0, 0, 0, 0, 0, 0, 0];
        _window(mids);
        vm.warp(REOPEN + 15 minutes);
        for (int256 answer = 0; answer > -2; --answer) {
            uint256 snapshot = vm.snapshotState();
            nvdaFeed.setRound(answer, REOPEN + 10, REOPEN + 20);
            cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
            _assertSettled(NVDA, 181.5e8, true);
            vm.revertToState(snapshot);
        }
    }

    function test_WithoutARoundSinceTheReopenTheSeriesSettlesOnTheMedian() public {
        _sold();
        _closed();
        uint64[15] memory mids = [uint64(0), 0, 0, 181e8, 181e8, 181e8, 181e8, 181e8, 181e8, 181e8, 181e8, 0, 0, 0, 0];
        _window(mids);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 181e8, true);
    }

    function test_SettlementTakesOnlyTheLastRoundsBeforeTheCloseAndBeforeTheReopen() public {
        _sold();
        nvdaFeed.setRound(199e8, FRIDAY + 5 hours, FRIDAY + 5 hours + 12);
        nvdaFeed.setRound(201e8, CLOSE_MS / 1000, CLOSE_MS / 1000 + 12);
        _closed();
        uint64[15] memory mids = [uint64(181e8), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
        _window(mids);
        nvdaFeed.setRound(181e8, REOPEN + 1, REOPEN + 13);
        vm.warp(REOPEN + 15 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST));
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST + 2);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST + 2));
        cover.settle(NVDA, CLOSE_MS, FIRST + 2, FIRST + 2);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST + 9));
        cover.settle(NVDA, CLOSE_MS, FIRST + 9, FIRST + 2);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLastRound.selector, FIRST + 1));
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 1);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLastRound.selector, FIRST + 3));
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 3);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidLastRound.selector, FIRST + 9));
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 9);
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 2);
        _assertSettled(NVDA, 181e8, false);
    }

    function test_TheReferenceAndTheFirstRoundAreFoundAcrossAPhaseChange() public {
        _sold();
        PhasedAggregator phased = new PhasedAggregator();
        phased.add(1, 600e8, FRIDAY - 9 days, FRIDAY - 9 days);
        band.setAsset(SPY, BandDouble.Asset(address(phased), bytes32(0), "USA500.Y---24_7", address(0)));
        GapCover second = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        _deposit(other, DEPOSIT, second);
        vm.prank(buyer);
        usdg.approve(address(second), type(uint256).max);
        vm.prank(buyer);
        second.buy(SPY, NOTIONAL, 300, 1000, type(uint256).max, holder);
        phased.add(2, 590e8, REOPEN + 30, REOPEN + 42);
        _closed();
        second.record();
        _reopened(REOPEN);
        _quote(SPY, 590e8, 3);
        for (uint256 k; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 50);
            second.observe(SPY, CLOSE_MS);
        }
        vm.warp(REOPEN + 15 minutes);
        second.settle(SPY, CLOSE_MS, FIRST, FIRST);
        (, uint64 referencePrice, uint64 price,,, bool flagged) = second.series(SPY, CLOSE_MS);
        assertEq(referencePrice, 600e8);
        assertEq(price, 590e8);
        assertFalse(flagged);
    }

    function test_TheWindowIsTheFifteenSlotsFromTheReopen() public {
        _sold();
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotReopened.selector, CLOSE_MS));
        cover.observe(NVDA, CLOSE_MS);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotReopened.selector, CLOSE_MS));
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _closed();
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotReopened.selector, CLOSE_MS));
        cover.observe(NVDA, CLOSE_MS);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NothingCovered.selector, SPY, CLOSE_MS));
        cover.observe(SPY, CLOSE_MS);
        _reopened(REOPEN + 59);
        _quote(NVDA, 180e8, 1);
        vm.expectEmit();
        emit GapCover.Observed(NVDA, CLOSE_MS, 0, 180e8, 179e8, 181e8);
        assertTrue(cover.observe(NVDA, CLOSE_MS));
        assertFalse(cover.observe(NVDA, CLOSE_MS));
        vm.warp(REOPEN + 3 minutes);
        _quote(NVDA, 0, 0);
        vm.recordLogs();
        assertFalse(cover.observe(NVDA, CLOSE_MS));
        assertEq(vm.getRecordedLogs().length, 0);
        (,,, uint16 shift,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(shift, 0);
        _quote(NVDA, 180e8, 3);
        for (uint256 k = 4; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 59);
            assertTrue(cover.observe(NVDA, CLOSE_MS));
        }
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowOpen.selector, REOPEN_MS + 15 minutes * 1000));
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        vm.warp(REOPEN + 15 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowClosed.selector, REOPEN_MS + 15 minutes * 1000));
        cover.observe(NVDA, CLOSE_MS);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 180e8, true);
        vm.expectRevert(abi.encodeWithSelector(GapCover.SeriesClosed.selector, NVDA, CLOSE_MS));
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        vm.expectRevert(abi.encodeWithSelector(GapCover.SeriesClosed.selector, NVDA, CLOSE_MS));
        cover.observe(NVDA, CLOSE_MS);
    }

    function test_ALateFirstRecordCannotChooseTheWindowOrThePrice() public {
        uint256 id = _sold();
        nvdaFeed.setRound(200e8, FRIDAY + 5 hours, FRIDAY + 5 hours + 12);
        _closed();
        nvdaFeed.setRound(200e8, REOPEN + 20, REOPEN + 32);
        _reopened(REOPEN + 3 days);
        _quote(NVDA, 170e8, 3);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowClosed.selector, REOPEN_MS + 15 minutes * 1000));
        cover.observe(NVDA, CLOSE_MS);
        vm.expectEmit();
        emit GapCover.Voided(NVDA, CLOSE_MS);
        cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 1);
        (uint256 payout, uint256 refund) = cover.release(id);
        assertEq(payout, 0);
        assertEq(refund, PREMIUM);
    }

    function test_TheMedianSettlesOnlyOnAQuorumOfTheWindowsSlots() public {
        uint256 id = _sold();
        _closed();
        uint64[15] memory mids = [uint64(181e8), 0, 183e8, 0, 185e8, 0, 187e8, 0, 189e8, 0, 191e8, 0, 193e8, 0, 0];
        _window(mids);
        vm.warp(REOPEN + 15 minutes);
        uint256 snapshot = vm.snapshotState();
        vm.expectEmit();
        emit GapCover.Voided(NVDA, CLOSE_MS);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        (, uint256 refund) = cover.release(id);
        assertEq(refund, PREMIUM);
        vm.revertToState(snapshot);
        vm.warp(REOPEN + 14 minutes + 50);
        _quote(NVDA, 180e8, 3);
        _observe(NVDA);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 186e8, true);
    }

    function test_ABandNotLiveAtTheReopenMovesTheWindowBySlotsRecordedSo() public {
        _sold();
        _closed();
        _reopened(REOPEN + 10);
        _quote(NVDA, 0, 0);
        vm.expectEmit();
        emit GapCover.Observed(NVDA, CLOSE_MS, 0, 0, 0, 0);
        assertTrue(cover.observe(NVDA, CLOSE_MS));
        _quote(NVDA, 150e8, 3);
        assertFalse(cover.observe(NVDA, CLOSE_MS));
        nvdaFeed.setRound(150e8, REOPEN + 15, REOPEN + 27);
        vm.warp(REOPEN + 70);
        _quote(NVDA, 170e8, 2);
        assertTrue(cover.observe(NVDA, CLOSE_MS));
        vm.warp(REOPEN + 130);
        _quote(NVDA, 170e8, 3);
        band.setSequencer(makeAddr("sequencer"), false);
        assertTrue(cover.observe(NVDA, CLOSE_MS));
        band.setSequencer(makeAddr("sequencer"), true);
        (,,, uint16 shift,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(shift, 3);
        for (uint256 k = 3; k < 18; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            _observe(NVDA);
        }
        vm.warp(REOPEN + 18 minutes - 1);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowOpen.selector, REOPEN_MS + 18 minutes * 1000));
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        vm.warp(REOPEN + 18 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowClosed.selector, REOPEN_MS + 18 minutes * 1000));
        cover.observe(NVDA, CLOSE_MS);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 170e8, true);
    }

    function test_AHaltNoOneRecordsInTheWindowsFirstSlotMovesNothing() public {
        _sold();
        _closed();
        _reopened(REOPEN + 70);
        _quote(NVDA, 0, 0);
        assertFalse(cover.observe(NVDA, CLOSE_MS));
        _quote(NVDA, 175e8, 3);
        for (uint256 k = 7; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            _observe(NVDA);
        }
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 175e8, true);
        (,,, uint16 shift,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(shift, 0);
    }

    function test_ABandNotLiveWithinFifteenSlotsOfTheReopenVoidsTheSeries() public {
        uint256 id = _sold();
        _closed();
        _reopened(REOPEN);
        _quote(NVDA, 0, 0);
        for (uint256 k; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            assertTrue(cover.observe(NVDA, CLOSE_MS));
        }
        uint256 snapshot = vm.snapshotState();
        _quote(NVDA, 170e8, 3);
        for (uint256 k = 15; k < 30; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            assertTrue(cover.observe(NVDA, CLOSE_MS));
        }
        (,,, uint16 shift,,) = cover.series(NVDA, CLOSE_MS);
        assertEq(shift, 15);
        vm.warp(REOPEN + 30 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        _assertSettled(NVDA, 170e8, true);
        vm.revertToState(snapshot);
        vm.warp(REOPEN + 15 minutes + 5);
        assertTrue(cover.observe(NVDA, CLOSE_MS));
        vm.warp(REOPEN + 16 minutes - 1);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowOpen.selector, REOPEN_MS + 16 minutes * 1000));
        cover.settle(NVDA, CLOSE_MS, 0, 0);
        vm.warp(REOPEN + 16 minutes);
        _quote(NVDA, 170e8, 3);
        vm.expectRevert(abi.encodeWithSelector(GapCover.WindowClosed.selector, REOPEN_MS + 16 minutes * 1000));
        cover.observe(NVDA, CLOSE_MS);
        vm.expectEmit();
        emit GapCover.Voided(NVDA, CLOSE_MS);
        cover.settle(NVDA, CLOSE_MS, 0, 0);
        (uint256 payout, uint256 refund) = cover.release(id);
        assertEq(payout, 0);
        assertEq(refund, PREMIUM);
    }

    function test_AMaterialMultiplierStepBetweenTheReferenceAndTheWindowVoidsTheSeries() public {
        _sold();
        nvdaFeed.setRound(199e8, FRIDAY + 5 hours, FRIDAY + 5 hours + 12);
        _closed();
        _fill(NVDA, 181e8);
        vm.warp(REOPEN + 15 minutes);
        uint64 end = uint64(REOPEN + 15 minutes);
        uint64 saturday = uint64(SATURDAY);
        uint64[7] memory effective = [saturday, uint64(FRIDAY + 5 hours), end, end - 1, saturday, saturday, saturday];
        uint128[7] memory after_ = [uint128(2e18), 2e18, 2e18, 2e18, 1.004e18, 1.005e18, 0.995e18];
        uint128[7] memory before = [uint128(1e18), 1e18, 1e18, 1e18, 1e18, 1e18, 1e18];
        bool[7] memory voids = [true, false, false, true, false, true, true];
        for (uint256 i; i < 8; ++i) {
            uint256 snapshot = vm.snapshotState();
            if (i < 7) {
                band.setCorporateAction(NVDA, BandDouble.CorporateAction(0, effective[i], before[i], after_[i]));
            } else {
                band.setCorporateAction(NVDA, BandDouble.CorporateAction(0, saturday, 0, 2e18));
            }
            vm.recordLogs();
            cover.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 1);
            Vm.Log[] memory logs = vm.getRecordedLogs();
            assertEq(logs[0].topics[0], i == 7 || voids[i] ? GapCover.Voided.selector : GapCover.Settled.selector);
            vm.revertToState(snapshot);
        }
    }

    function test_ACoverPaysTheFallBeyondItsDeductibleUpToItsLimit() public {
        uint256[6] memory prices = [uint256(210e8), 200e8, 190e8, 185e8, 170e8, 100e8];
        uint256[6] memory paid = [uint256(0), 0, 0, 250 * USDG, 1_000 * USDG, 1_000 * USDG];
        for (uint256 i; i < prices.length; ++i) {
            uint256 snapshot = vm.snapshotState();
            uint256 id = _sold();
            _closed();
            _fill(NVDA, uint64(prices[i]));
            vm.warp(REOPEN + 15 minutes);
            cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
            (uint256 payout, uint256 refund) = cover.release(id);
            assertEq(payout, paid[i]);
            assertEq(refund, 0);
            assertEq(cover.held(), DEPOSIT + PREMIUM - paid[i]);
            assertEq(cover.reserved(), 0);
            assertEq(cover.premiums(), 0);
            assertEq(cover.outstanding(), 0);
            assertEq(cover.outstandingIn(CLOSE_MS), 0);
            assertEq(cover.owed(), paid[i]);
            vm.revertToState(snapshot);
        }
    }

    function test_ReleaseNeedsASettledSeriesAndHappensOnce() public {
        uint256 id = _sold();
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotSettled.selector, id));
        cover.release(id);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NoCover.selector, 7));
        cover.release(7);
        _closed();
        _fill(NVDA, 200e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        cover.release(id);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NoCover.selector, id));
        cover.release(id);
        (address h,,,,,,) = cover.covers(id);
        assertEq(h, address(0));
        vm.prank(other);
        assertEq(cover.claim(other), 0);
    }

    function test_ASeriesNoOneSettlesWithinAWeekOfItsCloseIsVoidAndRefundsItsPremiums() public {
        uint256 id = _sold();
        vm.expectRevert(abi.encodeWithSelector(GapCover.NothingCovered.selector, SPY, CLOSE_MS));
        cover.void(SPY, CLOSE_MS);
        vm.warp(CLOSE_MS / 1000 + 7 days - 1);
        vm.expectRevert(abi.encodeWithSelector(GapCover.TooEarlyToVoid.selector, CLOSE_MS + 7 days * 1000));
        cover.void(NVDA, CLOSE_MS);
        vm.warp(CLOSE_MS / 1000 + 7 days);
        vm.expectEmit();
        emit GapCover.Voided(NVDA, CLOSE_MS);
        cover.void(NVDA, CLOSE_MS);
        vm.expectRevert(abi.encodeWithSelector(GapCover.SeriesClosed.selector, NVDA, CLOSE_MS));
        cover.void(NVDA, CLOSE_MS);
        (uint256 payout, uint256 refund) = cover.release(id);
        assertEq(payout, 0);
        assertEq(refund, PREMIUM);
        assertEq(cover.held(), DEPOSIT);
        assertEq(cover.premiums(), 0);
        _claimAll(holder, PREMIUM);
    }

    function test_WritersJoinOnlyBeforeTheSalesEndAndLeaveOnlyOnceEveryCoverIsReleased() public {
        assertEq(cover.maxDeposit(other), type(uint256).max);
        assertEq(cover.maxMint(other), type(uint256).max);
        assertEq(cover.maxRedeem(writer), cover.balanceOf(writer));
        assertEq(cover.maxWithdraw(writer), DEPOSIT);
        uint256 id = _sold();
        assertEq(cover.maxDeposit(other), type(uint256).max);
        assertEq(cover.maxRedeem(writer), 0);
        assertEq(cover.maxWithdraw(writer), 0);
        _deposit(other, DEPOSIT);
        assertEq(cover.balanceOf(other), cover.balanceOf(writer));
        _closed();
        assertEq(cover.maxDeposit(other), 0);
        assertEq(cover.maxMint(other), 0);
        band.setSession(BandDouble.Session(2, 1, 3, REGULAR_CLOSE_MS + 7 days * 1000, CLOSE_MS + 7 days * 1000));
        assertEq(cover.maxDeposit(other), 0);
        _fill(NVDA, 200e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        assertEq(cover.maxRedeem(writer), 0);
        cover.release(id);
        assertEq(cover.maxRedeem(writer), cover.balanceOf(writer));
        assertEq(cover.maxDeposit(other), type(uint256).max);
        uint256 shares = cover.balanceOf(writer);
        vm.prank(writer);
        uint256 assets = cover.redeem(shares, writer, writer);
        assertEq(assets, DEPOSIT + PREMIUM / 2);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        assertEq(cover.maxDeposit(other), type(uint256).max);
    }

    function test_DepositsCloseWhenTheBandsSessionCannotBeRead() public {
        _sold();
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        assertEq(cover.maxDeposit(other), 0);
        vm.mockCallRevert(address(band), abi.encodeWithSelector(BandDouble.session.selector), "");
        assertEq(cover.maxDeposit(other), 0);
    }

    function test_TheIssuersControlsStopTheVaultAndTheSalesButNotTheSettlement() public {
        uint256 id = _sold();
        usdg.pause();
        assertEq(cover.maxDeposit(other), 0);
        assertEq(cover.maxWithdraw(writer), 0);
        vm.expectRevert(StubUsdg.ContractPaused.selector);
        _buy(NVDA, NOTIONAL, 500, 1500);
        _closed();
        _fill(NVDA, 170e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        cover.release(id);
        vm.expectRevert(StubUsdg.ContractPaused.selector);
        vm.prank(holder);
        cover.claim(holder);
        usdg.unpause();
        usdg.freeze(holder);
        vm.expectRevert(StubUsdg.AddressFrozen.selector);
        vm.prank(holder);
        cover.claim(holder);
        assertEq(cover.payouts(holder), 1_000 * USDG);
        vm.prank(holder);
        cover.claim(other);
        assertEq(usdg.balanceOf(other), 1_000 * USDG);
        usdg.freeze(address(cover));
        assertEq(cover.maxRedeem(writer), 0);
        assertEq(cover.maxDeposit(other), 0);
        usdg.wipeFrozenAddress(address(cover));
        cover.sync();
        assertEq(cover.held(), 0);
        assertEq(cover.totalAssets(), 0);
    }

    function test_AWipeFallsOnTheWritersAfterThePremiumsAndCreditsTheCoverHolds() public {
        uint256 id = _sold();
        usdg.freeze(address(cover));
        usdg.wipeFrozenAddress(address(cover));
        usdg.unfreeze(address(cover));
        usdg.mint(address(cover), PREMIUM + 300 * USDG);
        vm.expectEmit();
        emit GapCover.Sync(300 * USDG);
        cover.sync();
        assertEq(cover.held(), 300 * USDG);
        assertEq(cover.capacity(), 0);
        cover.sync();
        _closed();
        _fill(NVDA, 150e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        (uint256 payout,) = cover.release(id);
        assertEq(payout, 300 * USDG + PREMIUM);
        assertEq(cover.held(), 0);
        _claimAll(holder, payout);
    }

    function test_TheWritersCountIsWhatTheUsdgLeavesAfterThePremiumsAndTheCredits() public {
        uint256 id = _sold();
        (, uint256 spyPremium) = _buy(SPY, NOTIONAL, 300, 1000);
        _closed();
        _fill(NVDA, 150e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        cover.release(id);
        assertEq(cover.owed(), 1_000 * USDG);
        assertEq(cover.premiums(), spyPremium);
        vm.recordLogs();
        cover.sync();
        assertEq(vm.getRecordedLogs().length, 0);
        usdg.freeze(address(cover));
        usdg.wipeFrozenAddress(address(cover));
        usdg.unfreeze(address(cover));
        usdg.mint(address(cover), 50_000 * USDG);
        cover.sync();
        assertEq(cover.held(), 50_000 * USDG - spyPremium - 1_000 * USDG);
    }

    function test_RedemptionsWaitWhileTheUsdgFallsShortOfTheCredits() public {
        uint256 id = _sold();
        _closed();
        _fill(NVDA, 150e8);
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        cover.release(id);
        uint256 shares = cover.balanceOf(writer);
        assertEq(cover.maxRedeem(writer), shares);
        uint256 counted = cover.held() + cover.premiums() + cover.owed();
        vm.mockCall(
            address(usdg), abi.encodeWithSignature("balanceOf(address)", address(cover)), abi.encode(counted - 1)
        );
        assertEq(cover.maxRedeem(writer), 0);
    }

    function test_DirectTransfersAndFailingReadsChangeNoShare() public {
        usdg.mint(address(cover), 5 * USDG);
        assertEq(cover.totalAssets(), DEPOSIT);
        vm.mockCallRevert(address(usdg), abi.encodeWithSelector(IUSDG.paused.selector), "");
        assertEq(cover.maxDeposit(other), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSelector(IUSDG.isFrozen.selector), "");
        assertEq(cover.maxDeposit(other), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("balanceOf(address)", address(cover)), "");
        assertEq(cover.maxRedeem(writer), 0);
    }

    function test_OnARegularHoursChainSalesEndAtTheRegularCloseAndTheReopeningIsTheRegularOpen() public {
        band.setSequencer(makeAddr("sequencer"), true);
        GapCover regular = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        assertTrue(regular.regularHours());
        _deposit(other, DEPOSIT, regular);
        vm.prank(buyer);
        usdg.approve(address(regular), type(uint256).max);
        _quote(NVDA, 199e8, 3);
        assertEq(regular.minDeductible(NVDA), 170);
        _quote(NVDA, 200e8, 3);
        vm.prank(buyer);
        uint256 id;
        (id,) = regular.buy(NVDA, NOTIONAL, 500, 1500, type(uint256).max, holder);
        nvdaFeed.setRound(199e8, REGULAR_CLOSE_MS / 1000 - 60, REGULAR_CLOSE_MS / 1000 - 48);
        nvdaFeed.setRound(201e8, REGULAR_CLOSE_MS / 1000 + 60, REGULAR_CLOSE_MS / 1000 + 72);
        vm.warp(REGULAR_CLOSE_MS / 1000);
        assertEq(regular.maxDeposit(writer), 0);
        vm.expectRevert(GapCover.SalesClosed.selector);
        vm.prank(buyer);
        regular.buy(NVDA, NOTIONAL, 500, 1500, type(uint256).max, holder);
        vm.warp(SATURDAY);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        regular.record();
        assertEq(regular.reopenOf(CLOSE_MS), REGULAR_OPEN_MS);
        nvdaFeed.setRound(195e8, REOPEN + 3600, REOPEN + 3612);
        nvdaFeed.setRound(180e8, REGULAR_OPEN_MS / 1000 + 30, REGULAR_OPEN_MS / 1000 + 42);
        vm.warp(REGULAR_OPEN_MS / 1000 - 1);
        band.setSession(BandDouble.Session(2, 3, 1, REGULAR_OPEN_MS, 0));
        _quote(NVDA, 181e8, 3);
        vm.expectRevert(abi.encodeWithSelector(GapCover.NotReopened.selector, CLOSE_MS));
        regular.observe(NVDA, CLOSE_MS);
        for (uint256 k; k < 15; ++k) {
            vm.warp(REGULAR_OPEN_MS / 1000 + k * 60 + 50);
            regular.observe(NVDA, CLOSE_MS);
        }
        vm.warp(REGULAR_OPEN_MS / 1000 + 15 minutes);
        vm.expectRevert(abi.encodeWithSelector(GapCover.InvalidReference.selector, FIRST + 2));
        regular.settle(NVDA, CLOSE_MS, FIRST + 2, FIRST + 3);
        regular.settle(NVDA, CLOSE_MS, FIRST + 1, FIRST + 3);
        (, uint64 referencePrice, uint64 price,,, bool flagged) = regular.series(NVDA, CLOSE_MS);
        assertEq(referencePrice, 199e8);
        assertEq(price, 180e8);
        assertFalse(flagged);
        regular.release(id);
    }

    function test_GasOfBuyingObservingSettlingAndReleasing() public {
        vm.prank(buyer);
        cover.buy(NVDA, NOTIONAL, 500, 1500, PREMIUM, holder);
        vm.snapshotGasLastCall("buy");
        vm.prank(buyer);
        cover.buy(NVDA, NOTIONAL, 500, 1500, PREMIUM, holder);
        vm.snapshotGasLastCall("buy, the series' next cover");
        vm.prank(buyer);
        cover.buy(SPY, NOTIONAL, 300, 1000, type(uint256).max, holder);
        nvdaFeed.setRound(181e8, REOPEN + 20, REOPEN + 32);
        _closed();
        _reopened(REOPEN + 5);
        _quote(NVDA, 181e8, 3);
        _quote(SPY, 590e8, 3);
        cover.observe(NVDA, CLOSE_MS);
        vm.snapshotGasLastCall("observe, the window's first slot");
        cover.observe(SPY, CLOSE_MS);
        for (uint256 k = 1; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 5);
            cover.observe(NVDA, CLOSE_MS);
            cover.observe(SPY, CLOSE_MS);
        }
        vm.snapshotGasLastCall("observe, a later slot");
        vm.warp(REOPEN + 15 minutes);
        cover.settle(NVDA, CLOSE_MS, FIRST, FIRST);
        vm.snapshotGasLastCall("settle, on the first round");
        cover.settle(SPY, CLOSE_MS, FIRST, FIRST);
        vm.snapshotGasLastCall("settle, on the median of 15 slots");
        cover.release(1);
        vm.snapshotGasLastCall("release");
        vm.prank(holder);
        cover.claim(holder);
        vm.snapshotGasLastCall("claim");
    }

    function test_GasOfMeasuringAWeekOfRounds() public {
        StubAggregator busy = new StubAggregator(8, 200e8, FRIDAY - 91 days, "NVDA / USD");
        for (uint256 k = 1; k < 91 * 13; ++k) {
            busy.setRound(int256(200e8 + (k % 7) * 1e8), FRIDAY - 91 days + k * 1 days / 13);
        }
        busy.setRound(200e8, FRIDAY - 1 hours);
        GapCover target = _coverOver(address(busy));
        vm.prank(buyer);
        target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        vm.snapshotGasLastCall("buy, measuring a week among 1,183 rounds");
        vm.prank(buyer);
        target.buy(NVDA, NOTIONAL, 600, 1600, type(uint256).max, holder);
        vm.snapshotGasLastCall("buy, the week measured");
        target = _coverOver(address(busy));
        target.measure(NVDA);
        vm.snapshotGasLastCall("measure, a week among 1,183 rounds");
    }

    function _sold() internal returns (uint256 id) {
        (id,) = _buy(NVDA, NOTIONAL, 500, 1500);
    }

    function _buy(bytes32 symbol, uint256 notional, uint256 deductible, uint256 limit)
        internal
        returns (uint256, uint256)
    {
        vm.prank(buyer);
        return cover.buy(symbol, notional, deductible, limit, type(uint256).max, holder);
    }

    function _closed() internal {
        vm.warp(SATURDAY);
        band.setSession(BandDouble.Session(1, 3, 1, REGULAR_OPEN_MS, REOPEN_MS));
        cover.record();
    }

    function _reopened(uint256 time) internal {
        vm.warp(time);
        band.setSession(BandDouble.Session(2, 3, 1, REGULAR_OPEN_MS, 0));
    }

    function _window(uint64[15] memory mids) internal {
        _reopened(REOPEN);
        for (uint256 k; k < 15; ++k) {
            if (mids[k] == 0) continue;
            vm.warp(REOPEN + k * 60 + 50);
            _quote(NVDA, mids[k], 3);
            _observe(NVDA);
        }
    }

    function _fill(bytes32 symbol, uint64 mid) internal {
        _reopened(REOPEN);
        _quote(symbol, mid, 3);
        for (uint256 k; k < 15; ++k) {
            vm.warp(REOPEN + k * 60 + 50);
            _observe(symbol);
        }
    }

    function _observe(bytes32 symbol) internal {
        vm.prank(keeper);
        assertTrue(cover.observe(symbol, CLOSE_MS));
    }

    function _quote(bytes32 symbol, uint64 mid, uint8 state) internal {
        band.setQuote(
            symbol,
            BandDouble.Quote(state, state == 0 ? 0 : 2, mid, 50, mid == 0 ? 0 : mid - 1e8, mid == 0 ? 0 : mid + 1e8)
        );
    }

    function _coverOverWeek(int256[8] memory week) internal returns (GapCover target, StubAggregator feed) {
        feed = new StubAggregator(8, week[7], CLOSE_MS / 1000 - 8 days - 1 hours, "NVDA / USD");
        for (uint256 k = 7; k > 0; --k) {
            feed.setRound(week[k - 1], CLOSE_MS / 1000 - k * 1 days - 1 hours);
        }
        feed.setRound(200e8, FRIDAY - 1 hours);
        target = _coverOver(address(feed));
    }

    function _coverOver(address feed) internal returns (GapCover target) {
        band.setAsset(NVDA, BandDouble.Asset(feed, "NVDA---24_7", bytes32(0), address(0)));
        target = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        band.setAsset(NVDA, BandDouble.Asset(address(nvdaFeed), "NVDA---24_7", bytes32(0), address(0)));
        vm.prank(buyer);
        usdg.approve(address(target), type(uint256).max);
        _deposit(writer, DEPOSIT, target);
    }

    function _assertGap(GapCover target, uint256 gap, uint256 move) internal view {
        (uint256 priced, uint256 moved) = target.pricingGap(NVDA);
        assertEq(priced, gap);
        assertEq(moved, move);
    }

    function _assertSettled(bytes32 symbol, uint256 price, bool flagged) internal view {
        (,, uint64 settledAt,, uint8 status, bool isFlagged) = cover.series(symbol, CLOSE_MS);
        assertEq(settledAt, price);
        assertEq(status, 1);
        assertEq(isFlagged, flagged);
    }

    function _assertSales(GapCover target, uint64 closesMs, uint64 endsMs) internal view {
        (uint64 closes, uint64 ends) = target.sales();
        assertEq(closes, closesMs);
        assertEq(ends, endsMs);
    }

    function _claimAll(address who, uint256 amount) internal {
        uint256 before = usdg.balanceOf(who);
        vm.expectEmit();
        emit GapCover.Claimed(who, who, amount);
        vm.prank(who);
        assertEq(cover.claim(who), amount);
        assertEq(usdg.balanceOf(who), before + amount);
        assertEq(cover.payouts(who), 0);
        assertEq(cover.owed(), 0);
    }

    function _deposit(address who, uint256 amount) internal {
        _deposit(who, amount, cover);
    }

    function _deposit(address who, uint256 amount, GapCover target) internal {
        usdg.mint(who, amount);
        vm.startPrank(who);
        usdg.approve(address(target), amount);
        target.deposit(amount, who);
        vm.stopPrank();
    }
}

/// @notice A Chainlink proxy whose rounds live in phases, as after aggregator upgrades: each phase numbers its own
/// rounds from one, an earlier phase's aggregator may go on answering after the upgrade, and the latest round is the
/// last one added. A round it does not have reverts, as the proxy's does for a phase without an aggregator. Test code
/// only.
contract PhasedAggregator {
    struct Round {
        int256 answer;
        uint256 startedAt;
        uint256 updatedAt;
    }

    mapping(uint80 => Round) internal rounds;
    mapping(uint16 => uint64) internal counts;
    uint80 internal latest;

    function add(uint16 phase, int256 answer, uint256 startedAt, uint256 updatedAt) external {
        uint80 id = uint80(phase) << 64 | ++counts[phase];
        rounds[id] = Round(answer, startedAt, updatedAt);
        latest = id;
    }

    function getRoundData(uint80 roundId) external view returns (uint80, int256, uint256, uint256, uint80) {
        Round memory r = rounds[roundId];
        if (r.updatedAt == 0) revert("No data present");
        return (roundId, r.answer, r.startedAt, r.updatedAt, roundId);
    }

    function latestRoundData() external view returns (uint80, int256, uint256, uint256, uint80) {
        Round memory r = rounds[latest];
        return (latest, r.answer, r.startedAt, r.updatedAt, latest);
    }
}

/// @notice A Chainlink proxy whose `getRoundData` spends all the gas it is given. Test code only.
contract GreedyAggregator {
    function getRoundData(uint80) external pure returns (uint80, int256, uint256, uint256, uint80) {
        assembly {
            invalid()
        }
    }

    function latestRoundData() external pure returns (uint80, int256, uint256, uint256, uint80) {
        return (uint80(1) << 64 | 2, 200e8, 1, 1, uint80(1) << 64 | 2);
    }
}
