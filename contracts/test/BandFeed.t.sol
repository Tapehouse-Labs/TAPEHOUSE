// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {OracleLibrary} from "@uniswap/v3-periphery/contracts/libraries/OracleLibrary.sol";
import {BandFeed} from "../src/BandFeed.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {PoolDouble, TokenDouble} from "./doubles/PoolDouble.sol";

contract BandFeedTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant NVDA_24_7 = "NVDA---24_7";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant USA500 = "USA500.Y---24_7";
    uint64 internal constant REOPEN_MS = 1_790_632_800_000;
    int24 internal constant NVDA_TICK = -223_380;

    BandDouble internal band;
    TokenDouble internal nvda;
    TokenDouble internal usdg;

    function setUp() public {
        vm.warp(1_790_390_000);
        band = new BandDouble();
        nvda = new TokenDouble(18);
        usdg = new TokenDouble(6);
        band.setAsset(NVDA, BandDouble.Asset(address(0xC1), NVDA_24_7, bytes32(0), address(nvda)));
        band.setQuote(NVDA, BandDouble.Quote(3, 2, 18_000_000_000, 55, 17_901_000_000, 18_099_000_000));
    }

    function _feed(BandFeed.Side side) internal returns (BandFeed) {
        return new BandFeed(IBand(address(band)), NVDA, side, IUniswapV3Pool(address(0)), "NVDA / USD");
    }

    function _pooled(address token0, address token1) internal returns (BandFeed feed, PoolDouble pool) {
        pool = new PoolDouble(token0, token1);
        feed = new BandFeed(IBand(address(band)), NVDA, BandFeed.Side.Low, IUniswapV3Pool(address(pool)), "NVDA / USD");
    }

    function test_ConstructorRejectsAnUnknownAsset() public {
        vm.expectRevert(abi.encodeWithSelector(BandFeed.UnknownAsset.selector, bytes32("AAPL")));
        new BandFeed(IBand(address(band)), "AAPL", BandFeed.Side.Low, IUniswapV3Pool(address(0)), "AAPL / USD");
    }

    function test_ConstructorRejectsAPoolThatDoesNotTradeTheToken() public {
        PoolDouble other = new PoolDouble(address(usdg), address(0xBEEF));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.PoolWithoutToken.selector, address(other), address(nvda)));
        new BandFeed(IBand(address(band)), NVDA, BandFeed.Side.Low, IUniswapV3Pool(address(other)), "NVDA / USD");

        band.setAsset(SPY, BandDouble.Asset(address(0xC2), bytes32(0), USA500, address(0)));
        PoolDouble pool = new PoolDouble(address(usdg), address(nvda));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.PoolWithoutToken.selector, address(pool), address(0)));
        new BandFeed(IBand(address(band)), SPY, BandFeed.Side.Low, IUniswapV3Pool(address(pool)), "SPY / USD");
    }

    function test_ConstructorRejectsAnUninitializedPool() public {
        PoolDouble pool = new PoolDouble(address(usdg), address(nvda));
        pool.setCardinality(0);
        vm.expectRevert(abi.encodeWithSelector(BandFeed.PoolNotInitialized.selector, address(pool)));
        new BandFeed(IBand(address(band)), NVDA, BandFeed.Side.Low, IUniswapV3Pool(address(pool)), "NVDA / USD");
    }

    function test_MetadataHasChainlinksShape() public {
        BandFeed feed = _feed(BandFeed.Side.Low);
        assertEq(feed.decimals(), 8);
        assertEq(feed.version(), 1);
        assertEq(feed.description(), "NVDA / USD");
        assertEq(address(feed.band()), address(band));
        assertEq(feed.symbol(), NVDA);
        assertEq(uint8(feed.side()), uint8(BandFeed.Side.Low));
        assertEq(address(feed.pool()), address(0));
    }

    function test_EachSideAnswersItsBound() public {
        uint64[3] memory expected = [uint64(17_901_000_000), 18_000_000_000, 18_099_000_000];
        for (uint8 i; i < 3; ++i) {
            (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredIn) =
                _feed(BandFeed.Side(i)).latestRoundData();
            assertEq(answer, int256(uint256(expected[i])));
            assertEq(roundId, vm.getBlockTimestamp());
            assertEq(answeredIn, roundId);
            assertEq(startedAt, vm.getBlockTimestamp());
            assertEq(updatedAt, vm.getBlockTimestamp());
        }
    }

    function test_AHaltedBandHasNoAnswer() public {
        BandFeed feed = _feed(BandFeed.Side.Low);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.NoAnswer.selector, NVDA));
        feed.latestRoundData();
    }

    function test_AnUnsettledSequencerHasNoAnswer() public {
        band.setSequencer(address(0x5E9), false);
        BandFeed feed = _feed(BandFeed.Side.Low);
        vm.expectRevert(BandFeed.SequencerNotSettled.selector);
        feed.latestRoundData();
        band.setSequencer(address(0x5E9), true);
        (, int256 answer,,,) = feed.latestRoundData();
        assertEq(answer, 17_901_000_000);
    }

    function test_WithoutASequencerFeedTheSequencerIsNeverRead() public {
        BandFeed feed = _feed(BandFeed.Side.Low);
        band.setSequencer(address(0), false);
        (, int256 answer,,,) = feed.latestRoundData();
        assertEq(answer, 17_901_000_000);
        assertTrue(feed.latestBand().sequencerSettled);
    }

    function test_OnlyTheCurrentRoundHasAnAnswer() public {
        BandFeed feed = _feed(BandFeed.Side.Mid);
        (uint80 roundId, int256 answer,,,) = feed.getRoundData(uint80(vm.getBlockTimestamp()));
        assertEq(roundId, vm.getBlockTimestamp());
        assertEq(answer, 18_000_000_000);
        vm.expectRevert(abi.encodeWithSelector(BandFeed.NoRound.selector, uint80(vm.getBlockTimestamp() - 1)));
        feed.getRoundData(uint80(vm.getBlockTimestamp() - 1));
    }

    function test_LatestBandCarriesTheBandSessionAndHalts() public {
        band.setVariance(NVDA_24_7, 144);
        band.setSession(BandDouble.Session(1, 3, 1, REOPEN_MS + 48_600_000, REOPEN_MS));
        band.setHalt(NVDA, BandDouble.Halt(true, 1_790_390_600, 1_790_389_990, true));
        band.setSequencer(address(0x5E9), false);
        BandFeed.Band memory b = _feed(BandFeed.Side.Low).latestBand();
        assertEq(b.state, 3);
        assertEq(b.live, 2);
        assertEq(b.mid, 18_000_000_000);
        assertEq(b.halfBps, 55);
        assertEq(b.low, 17_901_000_000);
        assertEq(b.high, 18_099_000_000);
        assertEq(b.variance, 144);
        assertEq(b.session, 1);
        assertEq(b.nyse, 3);
        assertEq(b.nyseNext, 1);
        assertEq(b.nyseChangeMs, REOPEN_MS + 48_600_000);
        assertEq(b.sessionBoundaryMs, REOPEN_MS);
        assertTrue(b.signedHalt);
        assertEq(b.haltUntil, 1_790_390_600);
        assertEq(b.haltIssuedAt, 1_790_389_990);
        assertTrue(b.oraclePaused);
        assertFalse(b.sequencerSettled);
        assertFalse(b.twapValid);
        assertEq(b.twap, 0);
        assertEq(b.premiumBps, 0);
    }

    function test_AnIndexAssetReportsItsIndexVariance() public {
        band.setAsset(SPY, BandDouble.Asset(address(0xC2), bytes32(0), USA500, address(0)));
        band.setVariance(USA500, 81);
        BandFeed feed = new BandFeed(IBand(address(band)), SPY, BandFeed.Side.Low, IUniswapV3Pool(address(0)), "SPY");
        assertEq(feed.latestBand().variance, 81);
    }

    function test_TheTwapIsTheTokensPriceInTheQuoteToken() public {
        (BandFeed usdgFirst, PoolDouble pool) = _pooled(address(usdg), address(nvda));
        pool.setTwap(-NVDA_TICK, 1800);
        uint256 expected = OracleLibrary.getQuoteAtTick(-NVDA_TICK, 1e18, address(nvda), address(usdg)) * 1e8 / 1e6;
        BandFeed.Band memory b = usdgFirst.latestBand();
        assertTrue(b.twapValid);
        assertEq(b.twap, expected);
        assertEq(b.premiumBps, (int256(expected) - 18_000_000_000) * 10_000 / 18_000_000_000);

        (BandFeed tokenFirst, PoolDouble other) = _pooled(address(nvda), address(usdg));
        other.setTwap(NVDA_TICK, 3600);
        expected = OracleLibrary.getQuoteAtTick(NVDA_TICK, 1e18, address(nvda), address(usdg)) * 1e8 / 1e6;
        b = tokenFirst.latestBand();
        assertTrue(b.twapValid);
        assertEq(b.twap, expected);
    }

    function test_APoolWithoutThirtyMinutesOfHistoryIsFlagged() public {
        (BandFeed feed, PoolDouble pool) = _pooled(address(usdg), address(nvda));
        pool.setTwap(-NVDA_TICK, 1799);
        BandFeed.Band memory b = feed.latestBand();
        assertFalse(b.twapValid);
        assertEq(b.twap, 0);
        assertEq(b.premiumBps, 0);
    }

    function test_AHaltedBandHasNoPremium() public {
        (BandFeed feed, PoolDouble pool) = _pooled(address(usdg), address(nvda));
        pool.setTwap(-NVDA_TICK, 1800);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        BandFeed.Band memory b = feed.latestBand();
        assertTrue(b.twapValid);
        assertGt(b.twap, 0);
        assertEq(b.premiumBps, 0);
    }

    function test_TheBandIsSealedOnlyInTheWindowBeforeAReopen() public {
        BandFeed feed = _feed(BandFeed.Side.Low);
        uint64 nowMs = uint64(vm.getBlockTimestamp()) * 1000;

        band.setSession(BandDouble.Session(2, 1, 2, nowMs + 1000, nowMs + 1000));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.NotSealWindow.selector, uint8(2), nowMs + 1000));
        feed.seal();

        band.setSession(BandDouble.Session(1, 3, 1, nowMs + 600_001, nowMs + 600_001));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.NotSealWindow.selector, uint8(1), nowMs + 600_001));
        feed.seal();

        band.setSession(BandDouble.Session(1, 3, 1, nowMs, nowMs));
        vm.expectRevert(abi.encodeWithSelector(BandFeed.NotSealWindow.selector, uint8(1), nowMs));
        feed.seal();

        uint64 reopenMs = nowMs + 600_000;
        band.setSession(BandDouble.Session(1, 3, 1, reopenMs, reopenMs));
        vm.expectEmit(address(feed));
        emit BandFeed.Sealed(reopenMs, 3, 2, 18_000_000_000, 55, 17_901_000_000, 18_099_000_000);
        vm.startSnapshotGas("seal");
        feed.seal();
        vm.stopSnapshotGas();
        _assertSeal(feed, reopenMs, BandDouble.Quote(3, 2, 18_000_000_000, 55, 17_901_000_000, 18_099_000_000));

        skip(300);
        band.setQuote(NVDA, BandDouble.Quote(2, 1, 18_100_000_000, 60, 17_991_400_000, 18_208_600_000));
        feed.seal();
        _assertSeal(feed, reopenMs, BandDouble.Quote(2, 1, 18_100_000_000, 60, 17_991_400_000, 18_208_600_000));
    }

    function _assertSeal(BandFeed feed, uint64 reopenMs, BandDouble.Quote memory q) internal view {
        (uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high, uint64 sealedAt) =
            feed.seals(reopenMs);
        assertEq(state, q.state);
        assertEq(live, q.live);
        assertEq(mid, q.mid);
        assertEq(halfBps, q.halfBps);
        assertEq(low, q.low);
        assertEq(high, q.high);
        assertEq(sealedAt, vm.getBlockTimestamp());
    }

    function testFuzz_ASealStandsOnlyInsideTheWindow(int256 aheadMs) public {
        BandFeed feed = _feed(BandFeed.Side.Low);
        int256 window = int256(uint256(feed.SEAL_WINDOW_MS()));
        aheadMs = bound(aheadMs, -2 * window, 2 * window);
        uint64 reopenMs = uint64(uint256(int256(vm.getBlockTimestamp() * 1000) + aheadMs));
        band.setSession(BandDouble.Session(1, 3, 1, reopenMs, reopenMs));
        if (aheadMs <= 0 || aheadMs > window) {
            vm.expectRevert(abi.encodeWithSelector(BandFeed.NotSealWindow.selector, uint8(1), reopenMs));
        }
        feed.seal();
    }
}
