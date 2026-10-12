// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {GapCover} from "../src/GapCover.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

/// @notice Covers settled against the real Chainlink rounds around real reopens: on Robinhood Chain, with the band each
/// minute of the window as the band program draws it from the captured 24/7 prices and the real rounds, or from the
/// rounds alone where the 24/7 prices were not captured; on Arbitrum One, where SPY's band has its Chainlink leg alone.
/// A sale on Robinhood Chain prices its layer from the band's centre where the feed's last round stands above it.
contract GapCoverForkTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant ROBINHOOD_SEP_21 = 68_360_724;
    uint256 internal constant ROBINHOOD_SEP_8 = 57_254_869;
    uint256 internal constant ROBINHOOD_SEP_11 = 60_522_480;
    uint256 internal constant ARBITRUM_SEP_8 = 503_046_442;
    uint256 internal constant ARBITRUM_SEP_14 = 505_093_846;
    uint256 internal constant ARBITRUM_APR_3 = 448_719_146;
    uint64 internal constant POST_MARKET_MS = 14_400_000;
    uint64 internal constant REOPEN_BEFORE_OPEN_MS = 48_600_000;
    uint64 internal constant REOPEN_BEFORE_MIDNIGHT_MS = 14_400_000;

    string internal data;
    string internal parameters;
    BandDouble internal band;
    GapCover internal cover;
    IERC20 internal usdg;
    address internal writer = makeAddr("writer");
    address internal buyer = makeAddr("buyer");
    address internal keeper = makeAddr("keeper");
    mapping(bytes32 symbol => uint256) internal ids;

    function setUp() public {
        data = vm.readFile("test/data/gap-cover-reopens.json");
        parameters = vm.readFile("../stylus/contracts/margin/parameters.json");
    }

    function test_TheFirstRoundsAfterARealReopenSettleInsideTheRealBand() public {
        string[] memory symbols = _symbols("NVDA TSLA AAPL MSFT GOOGL");
        _robinhood(ROBINHOOD_SEP_21, ".robinhood21Sep", symbols);
        _sellAndClose(".robinhood21Sep", symbols, false, 1);
        for (uint256 k; k < 15; ++k) {
            for (uint256 i; i < symbols.length; ++i) {
                _observe(".robinhood21Sep", symbols[i], k);
            }
        }
        for (uint256 i; i < symbols.length; ++i) {
            _settleOnTheFirstRound(".robinhood21Sep", symbols[i]);
        }
    }

    function test_AMissedFirstMinuteSettlesOnTheMedianOfTheRealBandAndIsFlagged() public {
        string[] memory symbols = _symbols("NVDA");
        _robinhood(ROBINHOOD_SEP_21, ".robinhood21Sep", symbols);
        _sellAndClose(".robinhood21Sep", symbols, false, 1);
        for (uint256 k = 1; k < 15; ++k) {
            _observe(".robinhood21Sep", "NVDA", k);
        }
        vm.warp(_reopenMs(".robinhood21Sep") / 1000 + 16 minutes);
        _settle(".robinhood21Sep", "NVDA");
        (, uint64 referencePrice, uint64 price,,, bool flagged) = cover.series("NVDA", _closeMs(".robinhood21Sep"));
        assertEq(referencePrice, 222_447_298_49);
        assertEq(price, 222_832_048_65);
        assertTrue(flagged);
        emit log_named_decimal_uint("NVDA's median band centre, 21 September, minutes 1 to 14", price, 8);
    }

    function test_AKeeperThatStartsLateLeavesTooFewMinutesAndTheSeriesIsVoid() public {
        string[] memory symbols = _symbols("NVDA");
        _robinhood(ROBINHOOD_SEP_21, ".robinhood21Sep", symbols);
        _sellAndClose(".robinhood21Sep", symbols, false, 1);
        for (uint256 k = 8; k < 15; ++k) {
            _observe(".robinhood21Sep", "NVDA", k);
        }
        vm.warp(_reopenMs(".robinhood21Sep") / 1000 + 15 minutes);
        vm.expectRevert(
            abi.encodeWithSelector(GapCover.WindowClosed.selector, _reopenMs(".robinhood21Sep") + 15 minutes * 1000)
        );
        cover.observe("NVDA", _closeMs(".robinhood21Sep"));
        vm.expectEmit();
        emit GapCover.Voided("NVDA", _closeMs(".robinhood21Sep"));
        _settle(".robinhood21Sep", "NVDA");
        (uint256 payout, uint256 refund) = cover.release(ids["NVDA"]);
        assertEq(payout, 0);
        assertGt(refund, 0);
    }

    function test_AReferenceAboveTheBandsCentreAtTheSaleRaisesTheSmallestDeductible() public {
        string[] memory symbols = _symbols("NVDA");
        _robinhood(ROBINHOOD_SEP_11, ".robinhood11Sep", symbols);
        uint64 closeMs = _closeMs(".robinhood11Sep");
        assertEq(vm.getBlockTimestamp(), vm.parseJsonUint(data, ".robinhood11Sep.sale"));
        band.setSession(BandDouble.Session(2, 1, 3, closeMs - POST_MARKET_MS, closeMs));
        (uint80 latestRound, int256 latest,,,) = AggregatorV3Interface(cover.feed("NVDA")).latestRoundData();
        assertEq(latestRound, _round(".robinhood11Sep", "NVDA", ".latest"));
        assertEq(uint256(latest), _answer(".robinhood11Sep", "NVDA", ".latest"));
        uint256 centre = _answer(".robinhood11Sep", "NVDA", ".reference");
        band.setQuote("NVDA", BandDouble.Quote(3, 2, uint64(centre), 30, uint64(centre * 997 / 1000), uint128(centre)));
        uint256 stale = Math.mulDiv(uint256(latest) - centre, 1e12, uint256(latest), Math.Rounding.Ceil);
        (uint256 gap,) = cover.pricingGap("NVDA");
        uint256 tail = Math.ceilDiv(gap * cover.TAIL_THRESHOLD_PPM(), 1e8);
        uint256 minimum = cover.minDeductible("NVDA");
        assertEq(minimum, Math.ceilDiv(gap * cover.TAIL_THRESHOLD_PPM() + stale, 1e8));
        uint256 premium = cover.quote("NVDA", 10_000 * USDG, minimum, minimum + 1_000);
        vm.startPrank(buyer);
        usdg.approve(address(cover), type(uint256).max);
        vm.expectRevert(abi.encodeWithSelector(GapCover.StaleReference.selector, minimum - 1, minimum));
        cover.buy("NVDA", 10_000 * USDG, minimum - 1, minimum + 1_000, type(uint256).max, buyer);
        (, uint256 paid) = cover.buy("NVDA", 10_000 * USDG, minimum, minimum + 1_000, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(paid, premium);
        band.setQuote("NVDA", BandDouble.Quote(3, 2, uint64(uint256(latest)), 30, 0, uint128(uint256(latest))));
        assertEq(cover.minDeductible("NVDA"), tail);
        assertGt(premium, cover.quote("NVDA", 10_000 * USDG, minimum, minimum + 1_000));
        emit log_named_decimal_uint("NVDA's last round above the band's centre, 11 September, pct", stale / 1e6, 4);
        emit log_named_uint("NVDA's smallest deductible at the sale, bps", minimum);
    }

    function test_TheWeekBeforeARealCloseRaisesTheGapFromTheRealRounds() public {
        string[] memory symbols = _symbols("NVDA TSLA AAPL MSFT GOOGL SPY");
        _robinhood(ROBINHOOD_SEP_8, ".robinhood8Sep", symbols);
        uint64 closeMs = _closeMs(".robinhood8Sep");
        vm.warp(closeMs / 1000 - 6 hours);
        band.setSession(BandDouble.Session(2, 1, 3, closeMs - POST_MARKET_MS, closeMs));
        uint256[2][6] memory expected = [
            [uint256(120_352), 60_176],
            [uint256(165_148), 82_574],
            [uint256(63_964), 31_982],
            [uint256(74_572), 37_286],
            [uint256(58_696), 29_348],
            [uint256(30_250), 15_125]
        ];
        for (uint256 i; i < symbols.length; ++i) {
            (uint256 gap, uint256 move) = cover.pricingGap(bytes32(bytes(symbols[i])));
            assertEq(gap, expected[i][0], symbols[i]);
            assertEq(move, expected[i][1], symbols[i]);
            emit log_named_decimal_uint(string.concat(symbols[i], "'s week to 3 September, pct"), move, 4);
        }
    }

    /// @dev SPY's proxy moved to phase 2 at this block, on 3 April; its phase 2 aggregator had answered since 19 March.
    /// The clock is set back to the sales of Friday 20 March, so the week to then starts in phase 2 and goes on in
    /// phase 1.
    function test_SpysWeekOnArbitrumOneIsFoundAcrossItsAggregatorUpgrade() public {
        _arbitrum(ARBITRUM_APR_3, ".arbitrum8Sep", _symbols("SPY"));
        uint64 closeMs = 1_774_051_200_000;
        vm.warp(closeMs / 1000 - 6 hours);
        band.setSession(BandDouble.Session(2, 1, 3, closeMs - POST_MARKET_MS, closeMs));
        (uint80 latest,,,,) = AggregatorV3Interface(cover.feed("SPY")).latestRoundData();
        assertEq(latest >> 64, 2);
        (uint256 gap, uint256 move) = cover.pricingGap("SPY");
        assertEq(gap, 41_090);
        assertEq(move, 20_545);
        uint256 before = gasleft();
        vm.prank(keeper);
        assertEq(cover.measure("SPY"), 20_545);
        emit log_named_uint("SPY's week to 20 March across phases 1 and 2, gas to measure", before - gasleft());
        emit log_named_decimal_uint("SPY's week to 20 March across phases 1 and 2, pct", move, 4);
    }

    function test_ATuesdayReopenAfterAHolidaySettlesOnRobinhoodChain() public {
        string[] memory symbols = _symbols("NVDA TSLA AAPL MSFT GOOGL SPY");
        _robinhood(ROBINHOOD_SEP_8, ".robinhood8Sep", symbols);
        _sellAndClose(".robinhood8Sep", symbols, false, 2);
        assertEq(cover.reopenOf(_closeMs(".robinhood8Sep")), 1_788_825_600_000);
        for (uint256 k; k < 15; ++k) {
            for (uint256 i; i < symbols.length; ++i) {
                _observe(".robinhood8Sep", symbols[i], k);
            }
        }
        for (uint256 i; i < symbols.length; ++i) {
            _settleOnTheFirstRound(".robinhood8Sep", symbols[i]);
        }
    }

    function test_OnArbitrumOneTheReopenIsTheRegularOpenAndALateFirstRoundLeavesTheMedian() public {
        string[] memory symbols = _symbols("SPY");
        _arbitrum(ARBITRUM_SEP_8, ".arbitrum8Sep", symbols);
        _sellAndClose(".arbitrum8Sep", symbols, true, 2);
        assertEq(cover.reopenOf(_closeMs(".arbitrum8Sep")), 1_788_874_200_000);
        for (uint256 k; k < 15; ++k) {
            _observe(".arbitrum8Sep", "SPY", k);
        }
        (, uint256 answer, uint256 startedAt,) = _first(".arbitrum8Sep", "SPY");
        assertEq(startedAt, 1_788_876_622);
        vm.warp(_reopenMs(".arbitrum8Sep") / 1000 + 45 minutes);
        _settle(".arbitrum8Sep", "SPY");
        (, uint64 referencePrice, uint64 price,,, bool flagged) = cover.series("SPY", _closeMs(".arbitrum8Sep"));
        assertEq(referencePrice, _answer(".arbitrum8Sep", "SPY", ".reference"));
        assertEq(price, _answer(".arbitrum8Sep", "SPY", ".last"));
        assertTrue(flagged);
        assertTrue(answer != price);
    }

    function test_OnArbitrumOneAPromptFirstRoundSettlesInsideTheBandOfChainlinkAlone() public {
        string[] memory symbols = _symbols("SPY");
        _arbitrum(ARBITRUM_SEP_14, ".arbitrum14Sep", symbols);
        _sellAndClose(".arbitrum14Sep", symbols, true, 1);
        for (uint256 k; k < 15; ++k) {
            _observe(".arbitrum14Sep", "SPY", k);
        }
        _settleOnTheFirstRound(".arbitrum14Sep", "SPY");
    }

    function _robinhood(uint256 block_, string memory key, string[] memory symbols) internal {
        string memory registry = vm.readFile("../deployments/4663.json");
        vm.createSelectFork("robinhood", block_);
        usdg = IERC20(vm.parseJsonAddress(registry, ".tokens.USDG"));
        address morpho = vm.parseJsonAddress(registry, ".morpho.Blue");
        vm.startPrank(morpho);
        usdg.transfer(writer, 100_000 * USDG);
        usdg.transfer(buyer, 1_000 * USDG);
        vm.stopPrank();
        _deploy(key, symbols, address(0));
    }

    function _arbitrum(uint256 block_, string memory key, string[] memory symbols) internal {
        vm.createSelectFork("arbitrum", block_);
        StubUsdg stub = new StubUsdg();
        stub.mint(writer, 100_000 * USDG);
        stub.mint(buyer, 1_000 * USDG);
        usdg = IERC20(address(stub));
        string memory registry = vm.readFile("../deployments/42161.json");
        _deploy(key, symbols, vm.parseJsonAddress(registry, ".chainlinkSequencer.Uptime"));
    }

    /// @dev A band with a sequencer-uptime feed follows NYSE's regular hours, as on Arbitrum One.
    function _deploy(string memory key, string[] memory symbols, address sequencer) internal {
        band = new BandDouble();
        if (sequencer != address(0)) band.setSequencer(sequencer, true);
        bytes32[] memory assets = new bytes32[](symbols.length);
        for (uint256 i; i < symbols.length; ++i) {
            assets[i] = bytes32(bytes(symbols[i]));
            address feed = vm.parseJsonAddress(data, string.concat(key, ".assets.", symbols[i], ".feed"));
            band.setAsset(assets[i], BandDouble.Asset(feed, bytes32(0), bytes32(0), address(0)));
        }
        MarginDouble engine = new MarginDouble(assets, address(band), address(0));
        for (uint256 i; i < symbols.length; ++i) {
            engine.setWeekendGap(
                assets[i], uint32(vm.parseJsonUint(parameters, string.concat(".weekendGap.", symbols[i], ".initial")))
            );
        }
        cover = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        vm.startPrank(writer);
        usdg.approve(address(cover), 100_000 * USDG);
        cover.deposit(100_000 * USDG, writer);
        vm.stopPrank();
    }

    function _sellAndClose(string memory key, string[] memory symbols, bool regularHours, uint8 nyseNext) internal {
        uint64 closeMs = _closeMs(key);
        vm.warp(closeMs / 1000 - 6 hours);
        band.setSession(BandDouble.Session(2, 1, 3, closeMs - POST_MARKET_MS, closeMs));
        vm.startPrank(buyer);
        usdg.approve(address(cover), type(uint256).max);
        for (uint256 i; i < symbols.length; ++i) {
            bytes32 symbol = bytes32(bytes(symbols[i]));
            (, int256 latest,,,) = AggregatorV3Interface(cover.feed(symbol)).latestRoundData();
            uint64 mid = uint64(uint256(latest));
            band.setQuote(symbol, BandDouble.Quote(3, 2, mid, 30, mid * 997 / 1000, uint128(mid) * 1003 / 1000));
            uint256 deductible = cover.minDeductible(symbol);
            (ids[symbol],) = cover.buy(symbol, 10_000 * USDG, deductible, deductible + 1_000, type(uint256).max, buyer);
        }
        vm.stopPrank();
        uint64 sessionReopenMs = _reopenMs(key) - (regularHours ? REOPEN_BEFORE_OPEN_MS : 0);
        uint64 changeMs = sessionReopenMs + (nyseNext == 2 ? REOPEN_BEFORE_MIDNIGHT_MS : REOPEN_BEFORE_OPEN_MS);
        vm.warp(closeMs / 1000 + 1 days);
        band.setSession(BandDouble.Session(1, 3, nyseNext, changeMs, sessionReopenMs));
        cover.record();
        band.setSession(BandDouble.Session(2, 3, 1, sessionReopenMs + REOPEN_BEFORE_OPEN_MS, 0));
    }

    function _observe(string memory key, string memory symbol, uint256 k) internal {
        string memory q = string.concat(key, ".assets.", symbol, ".quotes[", vm.toString(k), "]");
        uint256[] memory quote = vm.parseJsonUintArray(data, q);
        vm.warp(quote[0]);
        band.setQuote(
            bytes32(bytes(symbol)),
            BandDouble.Quote(
                uint8(quote[1]),
                uint8(quote[2]),
                uint64(quote[3]),
                uint64(quote[4]),
                uint64(quote[5]),
                uint128(quote[6])
            )
        );
        vm.prank(keeper);
        assertTrue(cover.observe(bytes32(bytes(symbol)), _closeMs(key)));
    }

    function _settleOnTheFirstRound(string memory key, string memory symbol) internal {
        (uint80 firstRound, uint256 answer,,) = _first(key, symbol);
        vm.warp(_reopenMs(key) / 1000 + 15 minutes);
        uint256 referencePrice = _answer(key, symbol, ".reference");
        (,, uint16 deductible, uint16 limit,, uint128 notional,) = cover.covers(ids[bytes32(bytes(symbol))]);
        uint256 fall = answer < referencePrice ? (referencePrice - answer) * 10_000 : 0;
        uint256 owed = fall > referencePrice * deductible
            ? notional * (Math.min(fall, referencePrice * limit) - referencePrice * deductible)
                / (referencePrice * 10_000)
            : 0;
        vm.expectEmit();
        emit GapCover.Settled(
            bytes32(bytes(symbol)),
            _closeMs(key),
            _round(key, symbol, ".reference"),
            referencePrice,
            firstRound,
            answer,
            false
        );
        _settle(key, symbol);
        (uint256 payout,) = cover.release(ids[bytes32(bytes(symbol))]);
        int256 gapBps = (int256(answer) - int256(referencePrice)) * 10_000 / int256(referencePrice);
        emit log_named_int(string.concat(symbol, "'s gap at the first round, bps"), gapBps);
        assertEq(payout, owed);
        if (payout != 0) emit log_named_decimal_uint(string.concat(symbol, "'s cover pays, USDG"), payout, 6);
    }

    function _settle(string memory key, string memory symbol) internal {
        vm.prank(keeper);
        cover.settle(
            bytes32(bytes(symbol)), _closeMs(key), _round(key, symbol, ".reference"), _round(key, symbol, ".last")
        );
    }

    function _first(string memory key, string memory symbol)
        internal
        view
        returns (uint80 roundId, uint256 answer, uint256 startedAt, uint256 updatedAt)
    {
        string memory f = string.concat(key, ".assets.", symbol, ".first");
        roundId = uint80(vm.parseUint(vm.parseJsonString(data, string.concat(f, "[0]"))));
        answer = vm.parseJsonUint(data, string.concat(f, "[1]"));
        startedAt = vm.parseJsonUint(data, string.concat(f, "[2]"));
        updatedAt = vm.parseJsonUint(data, string.concat(f, "[3]"));
    }

    function _round(string memory key, string memory symbol, string memory which) internal view returns (uint80) {
        return uint80(vm.parseUint(vm.parseJsonString(data, string.concat(key, ".assets.", symbol, which, "[0]"))));
    }

    function _answer(string memory key, string memory symbol, string memory which) internal view returns (uint256) {
        return vm.parseJsonUint(data, string.concat(key, ".assets.", symbol, which, "[1]"));
    }

    function _closeMs(string memory key) internal view returns (uint64) {
        return uint64(vm.parseJsonUint(data, string.concat(key, ".close")) * 1000);
    }

    function _reopenMs(string memory key) internal view returns (uint64) {
        return uint64(vm.parseJsonUint(data, string.concat(key, ".reopen")) * 1000);
    }

    function _symbols(string memory list) internal pure returns (string[] memory) {
        return vm.split(list, " ");
    }
}
