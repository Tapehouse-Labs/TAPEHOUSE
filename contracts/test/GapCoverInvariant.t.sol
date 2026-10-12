// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {GapCover} from "../src/GapCover.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

/// @notice Drives the cover through weeks of sales, closures, reopens, observations, settlements, voids, releases,
/// claims, deposits and redemptions over two assets, and records any panic.
contract GapCoverHandler is Test {
    struct Seen {
        address holder;
        uint8 status;
        uint128 premium;
        uint256 reserve;
        uint256 expected;
    }

    uint256 internal constant USDG = 1e6;
    uint256 internal constant FRIDAY = 1_790_949_600;
    uint256 internal constant WEEK = 7 days;
    uint80 internal constant PHASE = uint80(1) << 64;

    GapCover public immutable cover;
    StubUsdg public immutable usdg;
    BandDouble public immutable band;
    bytes32[2] internal symbols = [bytes32("NVDA"), bytes32("SPY")];
    StubAggregator[2] internal feeds;
    uint256[2] internal rounds;
    uint256[2] internal prices = [uint256(200e8), 600e8];
    mapping(uint256 asset => mapping(uint256 round => uint256)) internal started;
    address[3] internal actors = [makeAddr("first"), makeAddr("second"), makeAddr("third")];

    uint8 public phase;
    uint256 public week;
    uint256 internal minute;
    uint256[] internal open;

    bool public panicked;
    bool public overpaid;
    uint256 public deposited;
    uint256 public withdrawn;
    uint256 public premiumsReleased;
    uint256 public paidOut;
    uint256 public reservedOutstanding;

    constructor(GapCover cover_, StubUsdg usdg_, BandDouble band_, StubAggregator nvda, StubAggregator spy) {
        (cover, usdg, band) = (cover_, usdg_, band_);
        feeds = [nvda, spy];
        rounds = [uint256(1), 1];
        (started[0][1], started[1][1]) = (vm.getBlockTimestamp() - 1 hours, vm.getBlockTimestamp() - 1 hours);
    }

    function deposit(uint256 actorSeed, uint256 amount) external {
        address actor = actors[actorSeed % 3];
        if (cover.maxDeposit(actor) == 0) return;
        amount = bound(amount, 1, 50_000 * USDG);
        usdg.mint(actor, amount);
        vm.startPrank(actor);
        usdg.approve(address(cover), amount);
        try cover.deposit(amount, actor) {
            deposited += amount;
        } catch (bytes memory reason) {
            _check(reason);
        }
        vm.stopPrank();
    }

    function redeem(uint256 actorSeed, uint256 share) external {
        address actor = actors[actorSeed % 3];
        uint256 shares = cover.maxRedeem(actor) * bound(share, 1, 100) / 100;
        if (shares == 0) return;
        vm.prank(actor);
        try cover.redeem(shares, actor, actor) returns (uint256 assets) {
            withdrawn += assets;
        } catch (bytes memory reason) {
            _check(reason);
        }
    }

    function buy(uint256 actorSeed, uint256 assetSeed, uint256 notional, uint256 deductible, uint256 limit) external {
        if (phase != 0) return;
        uint256 i = assetSeed % 2;
        address actor = actors[actorSeed % 3];
        deductible = bound(deductible, cover.minDeductible(symbols[i]), 9_999);
        limit = bound(limit, deductible + 1, 10_000);
        uint256 room = cover.capacity() * 10_000 / (limit - deductible);
        if (room == 0) return;
        notional = bound(notional, 1, room);
        usdg.mint(actor, 1_000_000 * USDG);
        vm.startPrank(actor);
        usdg.approve(address(cover), type(uint256).max);
        try cover.buy(symbols[i], notional, deductible, limit, type(uint256).max, actor) returns (uint256 id, uint256) {
            open.push(id);
            reservedOutstanding += (notional * (limit - deductible) + 9_999) / 10_000;
        } catch (bytes memory reason) {
            _check(reason);
        }
        vm.stopPrank();
    }

    function measure(uint256 assetSeed) external {
        if (phase != 0) return;
        _call(abi.encodeCall(GapCover.measure, (symbols[assetSeed % 2])));
    }

    function move(uint256 seed) external {
        if (phase != 0) return;
        for (uint256 i; i < 2; ++i) {
            prices[i] = prices[i] * bound(seed >> (8 * i), 95, 105) / 100;
            _round(i, prices[i], _friday() + 3 hours + bound(seed >> (16 + i), 0, 2 hours));
        }
    }

    function close(uint256 seed) external {
        if (phase != 0 || seed % 4 != 0) return;
        vm.warp(_close() / 1000 + 1 hours);
        if (seed % 2 == 0) _round(seed % 4 / 2, prices[seed % 4 / 2], _close() / 1000 + 90);
        band.setSession(BandDouble.Session(1, 3, 1, _reopen() + 48_600_000, _reopen()));
        _call(abi.encodeCall(GapCover.record, ()));
        phase = 1;
        minute = 0;
    }

    function observe(uint256 seed, uint256 move_, uint256 state) external {
        for (uint256 n = 1 + seed % 5; n != 0 && phase == 1 && minute <= 20; --n) {
            _observe(seed, move_, state);
            (seed, move_, state) = (seed >> 3, move_ >> 16, state >> 16);
        }
    }

    function settle(uint256 seed) external {
        if (phase != 1 || minute < 16) return;
        vm.warp(vm.getBlockTimestamp() + 40 minutes + seed % 7 days);
        for (uint256 i; i < 2; ++i) {
            (uint128 notional,,,, uint8 status,) = cover.series(symbols[i], _close());
            if (notional == 0 || status != 0) continue;
            if (seed % 5 == i) {
                if (vm.getBlockTimestamp() * 1000 < _close() + 7 days * 1000) vm.warp(_close() / 1000 + 7 days);
                _call(abi.encodeCall(GapCover.void, (symbols[i], _close())));
                continue;
            }
            _call(abi.encodeCall(GapCover.settle, (symbols[i], _close(), _last(i, _close()), _last(i, _reopen()))));
        }
        phase = 2;
    }

    function release(uint256 seed) external {
        if (open.length == 0) return;
        uint256 k = seed % open.length;
        Seen memory v = _seen(open[k]);
        if (v.status == 0) return;
        uint256 owedBefore = cover.payouts(v.holder);
        try cover.release(open[k]) returns (uint256 payout, uint256 refund) {
            if (payout > v.reserve || (v.status == 1 && payout != v.expected) || (v.status == 2 && refund != v.premium))
            {
                overpaid = true;
            }
            if (cover.payouts(v.holder) != owedBefore + payout + refund) overpaid = true;
            reservedOutstanding -= v.reserve;
            if (v.status == 1) premiumsReleased += v.premium;
            paidOut += payout;
            open[k] = open[open.length - 1];
            open.pop();
        } catch (bytes memory reason) {
            _check(reason);
        }
    }

    function claim(uint256 actorSeed) external {
        address actor = actors[actorSeed % 3];
        vm.prank(actor);
        try cover.claim(actor) {}
        catch (bytes memory reason) {
            _check(reason);
        }
    }

    function nextWeek() external {
        if (phase != 2) return;
        while (_friday() <= vm.getBlockTimestamp()) ++week;
        vm.warp(_friday());
        band.setSession(BandDouble.Session(2, 1, 3, uint64(_close() - 14_400_000), _close()));
        for (uint256 i; i < 2; ++i) {
            band.setQuote(
                symbols[i],
                BandDouble.Quote(
                    3, 2, uint64(prices[i]), 50, uint64(prices[i] * 995 / 1000), uint128(prices[i] * 1005 / 1000)
                )
            );
        }
        phase = 0;
    }

    function openCount() external view returns (uint256) {
        return open.length;
    }

    function _seen(uint256 id) internal view returns (Seen memory v) {
        uint64 closesMs;
        uint16 deductible;
        uint16 limit;
        bytes32 symbol;
        uint128 notional;
        (v.holder, closesMs, deductible, limit, symbol, notional, v.premium) = cover.covers(id);
        v.reserve = (uint256(notional) * (limit - deductible) + 9_999) / 10_000;
        (, uint64 referencePrice, uint64 price,, uint8 status,) = cover.series(symbol, closesMs);
        v.status = status;
        v.expected = _payout(notional, deductible, limit, referencePrice, price);
    }

    function _observe(uint256 seed, uint256 move_, uint256 state) internal {
        vm.warp(_reopen() / 1000 + minute * 60 + bound(seed, 0, 59));
        band.setSession(BandDouble.Session(2, 3, 1, _reopen() + 48_600_000, 0));
        for (uint256 i; i < 2; ++i) {
            uint256 mid = prices[i] * bound(move_ >> (8 * i), 60, 120) / 100;
            uint8 s = uint8(bound(state >> (8 * i), 0, 3));
            band.setQuote(
                symbols[i],
                BandDouble.Quote(s, 1, uint64(mid), 50, uint64(mid * 995 / 1000), uint128(mid * 1005 / 1000))
            );
            if (minute == 0 && seed % 3 != 0) _round(i, mid, vm.getBlockTimestamp());
            (uint128 notional,,,, uint8 status,) = cover.series(symbols[i], _close());
            if (notional != 0 && status == 0) _call(abi.encodeCall(GapCover.observe, (symbols[i], _close())));
        }
        ++minute;
    }

    function _round(uint256 i, uint256 price, uint256 startedAt) internal {
        if (startedAt < started[i][rounds[i]]) return;
        feeds[i].setRound(int256(price), startedAt, startedAt + 12);
        started[i][++rounds[i]] = startedAt;
    }

    function _last(uint256 i, uint256 beforeMs) internal view returns (uint80) {
        uint256 r = rounds[i];
        while (r > 1 && started[i][r] * 1000 >= beforeMs) --r;
        return PHASE | uint80(r);
    }

    function _payout(uint256 notional, uint256 deductible, uint256 limit, uint256 reference_, uint256 price)
        internal
        pure
        returns (uint256)
    {
        if (price >= reference_) return 0;
        uint256 fall = (reference_ - price) * 10_000;
        if (fall > reference_ * limit) fall = reference_ * limit;
        if (fall <= reference_ * deductible) return 0;
        return notional * (fall - reference_ * deductible) / (reference_ * 10_000);
    }

    function _call(bytes memory data) internal {
        (bool ok, bytes memory reason) = address(cover).call(data);
        if (!ok) _check(reason);
    }

    function _check(bytes memory reason) internal {
        if (reason.length >= 4 && bytes4(reason) == bytes4(keccak256("Panic(uint256)"))) panicked = true;
    }

    function _friday() internal view returns (uint256) {
        return FRIDAY + week * WEEK;
    }

    function _close() internal view returns (uint64) {
        return uint64((FRIDAY + 10 hours + week * WEEK) * 1000);
    }

    function _reopen() internal view returns (uint64) {
        return _close() + 2 days * 1000;
    }
}

contract GapCoverInvariantTest is Test {
    GapCover internal cover;
    StubUsdg internal usdg;
    GapCoverHandler internal handler;

    function setUp() public {
        vm.warp(1_790_949_600);
        usdg = new StubUsdg();
        BandDouble band = new BandDouble();
        StubAggregator nvda = new StubAggregator(8, 200e8, vm.getBlockTimestamp() - 9 days, "NVDA / USD");
        StubAggregator spy = new StubAggregator(8, 600e8, vm.getBlockTimestamp() - 9 days, "SPY / USD");
        band.setAsset("NVDA", BandDouble.Asset(address(nvda), "NVDA---24_7", bytes32(0), address(0)));
        band.setAsset("SPY", BandDouble.Asset(address(spy), bytes32(0), "USA500.Y---24_7", address(0)));
        band.setQuote("NVDA", BandDouble.Quote(3, 2, 200e8, 50, 199e8, 201e8));
        band.setQuote("SPY", BandDouble.Quote(3, 2, 600e8, 50, 597e8, 603e8));
        band.setSession(BandDouble.Session(2, 1, 3, 1_790_971_200_000, 1_790_985_600_000));
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (bytes32("NVDA"), bytes32("SPY"));
        MarginDouble engine = new MarginDouble(symbols, address(band), address(0));
        engine.setWeekendGap("NVDA", 118_601);
        engine.setWeekendGap("SPY", 54_840);
        cover = new GapCover(IMargin(address(engine)), IUSDG(address(usdg)));
        handler = new GapCoverHandler(cover, usdg, band, nvda, spy);
        targetContract(address(handler));
    }

    /// The cover holds exactly the writers' USDG, the premiums of the covers not yet released and the credits not yet
    /// claimed.
    function invariant_TheCoverHoldsEveryCount() public view {
        assertEq(usdg.balanceOf(address(cover)), cover.held() + cover.premiums() + cover.owed());
    }

    /// Every outstanding cover's whole payout is reserved from the writers' USDG.
    function invariant_EveryCoverIsFullyReserved() public view {
        assertEq(cover.reserved(), handler.reservedOutstanding());
        assertLe(cover.reserved(), cover.held());
        assertEq(cover.outstanding(), handler.openCount());
    }

    /// The writers' USDG is what they put in and left, plus the premiums released, less the payouts.
    function invariant_WritersEarnThePremiumsLessThePayouts() public view {
        assertEq(
            cover.held() + handler.withdrawn() + handler.paidOut(), handler.deposited() + handler.premiumsReleased()
        );
    }

    /// No release paid more than its cover reserved or other than the fall beyond its deductible up to its limit.
    function invariant_NoReleasePaysOtherThanItsLayer() public view {
        assertFalse(handler.overpaid());
    }

    /// Nothing the handler called panicked.
    function invariant_NothingPanics() public view {
        assertFalse(handler.panicked());
    }
}
