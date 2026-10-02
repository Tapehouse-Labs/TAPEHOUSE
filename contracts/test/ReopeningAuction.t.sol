// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {Basket} from "../src/Basket.sol";
import {BandFeed} from "../src/BandFeed.sol";
import {GapBackstop} from "../src/GapBackstop.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {ReopeningAuction} from "../src/ReopeningAuction.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {SealDouble} from "./doubles/SealDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

abstract contract ReopeningAuctionBase is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubToken internal weth;
    StubAggregator internal ethUsd;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    BandDouble internal band;
    SealDouble internal seal;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    ReopeningAuction internal auction;
    address internal owner = makeAddr("owner");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal dave = makeAddr("dave");
    address internal eve = makeAddr("eve");
    uint64 internal reopenMs;
    uint64 internal openMs;

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        weth = new StubToken(18);
        nvda = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(NVDA, 200e8, 202e8);
        _quote(SPY, 600e8, 606e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (NVDA, SPY);
        ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        MarginDouble engine = new MarginDouble(symbols, address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](2);
        (caps[0], caps[1]) = (type(uint256).max, type(uint256).max);
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            weth,
            owner,
            caps,
            type(uint256).max,
            type(uint256).max,
            0,
            10_00
        );
        liquidator = new Liquidator(accounts);
        seal = new SealDouble(address(band), NVDA);
        BandFeed[] memory feeds = new BandFeed[](1);
        feeds[0] = BandFeed(address(seal));
        bytes32[] memory sealed_ = new bytes32[](1);
        sealed_[0] = NVDA;
        auction = new ReopeningAuction(liquidator, sealed_, feeds);
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        liquidator.setAuction(address(auction));
        vm.stopPrank();
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), 1_000_000 * USDG);
        vault.deposit(1_000_000 * USDG, owner);
        vm.stopPrank();
    }

    function _emptied(address account) internal {
        _emptiedBut(account, 0);
    }

    function _emptiedBut(address account, uint256 left) internal {
        _quote(NVDA, 200e8, 202e8);
        _position(account, 10 * SHARE, 1_600 * USDG);
        _quote(NVDA, 100e8, 101e8);
        liquidator.start(account, CROSS);
        usdg.mint(dave, 10_000 * USDG);
        vm.startPrank(dave);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(account, CROSS, address(nvda), 10 * SHARE - left, type(uint256).max, dave);
        vm.stopPrank();
    }

    function _weekend(uint64 sealedLow) internal {
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        reopenMs = closes + 48 hours * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, reopenMs));
        accounts.accruePremium();
        if (sealedLow != 0) {
            seal.setSeal(
                reopenMs,
                SealDouble.Seal(2, 1, sealedLow + 1e8, 50, sealedLow, sealedLow + 2e8, uint64(vm.getBlockTimestamp()))
            );
        }
        skip(48 hours);
        openMs = reopenMs + 13.5 hours * 1000;
        band.setSession(BandDouble.Session(2, 3, 1, openMs, 0));
        accounts.accruePremium();
        _quote(NVDA, 190e8, 192e8);
    }

    function _toReveal() internal {
        vm.warp(openMs / 1000 - 30 minutes);
    }

    function _open() internal {
        vm.warp(openMs / 1000);
        band.setSession(BandDouble.Session(2, 1, 2, openMs + 6.5 hours * 1000, 0));
    }

    function _commit(address bidder, uint256 quantity, uint256 price, bytes32 salt)
        internal
        returns (bytes32 commitment)
    {
        commitment = keccak256(abi.encode(bidder, NVDA, openMs, quantity, price, salt));
        _fund(bidder, 10_100 * USDG);
        uint256 deposit = (quantity * price - 1) / 1e20 + 1;
        vm.prank(bidder);
        auction.commit(NVDA, commitment, deposit < 100 * USDG ? 100 * USDG : deposit);
    }

    function _reveal(address bidder, uint256 quantity, uint256 price, bytes32 salt) internal {
        vm.prank(bidder);
        auction.reveal(NVDA, openMs, quantity, price, salt);
    }

    function _fund(address who, uint256 amount) internal {
        usdg.mint(who, amount);
        vm.prank(who);
        usdg.approve(address(auction), type(uint256).max);
    }

    function _need(uint256 owed, uint256 floor) internal pure returns (uint256) {
        return (owed * 10_000 * 1e20 - 1) / (9_950 * floor) + 1;
    }

    function _position(address account, uint256 shares, uint256 debt) internal {
        _deposit(account, CROSS, address(nvda), shares);
        vm.prank(account);
        accounts.borrow(CROSS, debt, account, account);
    }

    function _quote(bytes32 symbol, uint64 low, uint128 high) internal {
        uint64 mid = uint64((uint256(low) + high) / 2);
        band.setQuote(symbol, BandDouble.Quote(low == 0 ? 0 : 3, 3, mid, 50, low, high));
    }

    function _deposit(address account, bytes32 position, address token, uint256 amount) internal {
        if (token == address(usdg)) usdg.mint(account, amount);
        else StubToken(token).mint(account, amount);
        vm.startPrank(account);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, account);
        vm.stopPrank();
    }
}

contract ReopeningAuctionTest is ReopeningAuctionBase {
    function test_TheAuctionTakesItsReadsFromTheLiquidatorAndItsFeeds() public {
        assertEq(address(auction.liquidator()), address(liquidator));
        assertEq(address(auction.accounts()), address(accounts));
        assertEq(address(auction.band()), address(band));
        assertEq(address(auction.usdg()), address(usdg));
        assertEq(address(auction.feeds(NVDA)), address(seal));
        assertEq(address(auction.feeds(SPY)), address(0));
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = SPY;
        BandFeed[] memory feeds = new BandFeed[](1);
        feeds[0] = BandFeed(address(seal));
        vm.expectRevert(ReopeningAuction.InvalidFeeds.selector);
        this.deploy(symbols, feeds);
        vm.expectRevert(ReopeningAuction.InvalidFeeds.selector);
        this.deploy(symbols, new BandFeed[](0));
        vm.expectRevert(ReopeningAuction.InvalidFeeds.selector);
        this.deploy(new bytes32[](0), feeds);
        symbols[0] = NVDA;
        feeds[0] = BandFeed(address(new SealDouble(address(new BandDouble()), NVDA)));
        vm.expectRevert(ReopeningAuction.InvalidFeeds.selector);
        this.deploy(symbols, feeds);
        assertEq(address(this.deploy(new bytes32[](0), new BandFeed[](0)).liquidator()), address(liquidator));
    }

    function test_AShortPositionIsEnrolledFromTheSealedBandAndHeldFromTheDutchAuction() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _position(bob, 10 * SHARE, 1_000 * USDG);
        _deposit(dave, SPY, address(spy), 1 * SHARE);
        vm.prank(dave);
        accounts.borrow(SPY, 480 * USDG, dave, dave);
        _position(eve, 10 * SHARE, 1_520 * USDG);
        _weekend(190e8);
        _quote(SPY, 590e8, 596e8);
        vm.expectEmit(address(auction));
        emit ReopeningAuction.RoundOpened(NVDA, openMs, reopenMs, 171e8, true);
        vm.expectEmit(address(liquidator));
        emit Liquidator.Held(alice, CROSS, uint64(openMs / 1000 + 1 hours));
        auction.enroll(alice, CROSS, NVDA);
        ReopeningAuction.Lot[] memory lots = auction.lots(NVDA, openMs);
        assertEq(lots.length, 1);
        assertEq(lots[0].account, alice);
        assertEq(lots[0].amount, _need(1_600 * USDG, 171e8));
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.sealMs, reopenMs);
        assertEq(r.floor, 171e8);
        assertEq(r.supply, lots[0].amount);
        assertTrue(auction.enrolled(openMs, alice, CROSS, NVDA));
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, alice, CROSS));
        auction.enroll(alice, CROSS, NVDA);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, bob, CROSS));
        auction.enroll(bob, CROSS, NVDA);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, carol, CROSS));
        auction.enroll(carol, CROSS, NVDA);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownAsset.selector, bytes32("TSLA")));
        auction.enroll(alice, CROSS, "TSLA");
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, dave, SPY));
        auction.enroll(dave, SPY, NVDA);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, eve, CROSS));
        auction.enroll(eve, CROSS, NVDA);
        _open();
        vm.expectRevert(
            abi.encodeWithSelector(Liquidator.PositionHeld.selector, alice, CROSS, uint64(openMs / 1000 + 1 hours))
        );
        liquidator.start(alice, CROSS);
        vm.expectRevert(
            abi.encodeWithSelector(Liquidator.PositionHeld.selector, alice, CROSS, uint64(openMs / 1000 + 1 hours))
        );
        liquidator.buy(alice, CROSS, address(nvda), 1, 1, alice);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.enroll(bob, CROSS, NVDA);
    }

    function test_ARoundWithNoSealTakesTheBandAtItsOpening() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(0);
        _quote(NVDA, 180e8, 182e8);
        vm.expectEmit(address(auction));
        emit ReopeningAuction.RoundOpened(NVDA, openMs, reopenMs, 162e8, false);
        auction.enroll(alice, CROSS, NVDA);
        _quote(NVDA, 170e8, 172e8);
        _commit(carol, 2 * SHARE, 185e8, "c");
        assertEq(auction.round(NVDA, openMs).floor, 162e8);
        _quote(SPY, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NoPrice.selector, SPY));
        auction.commit(SPY, bytes32(uint256(1)), 100 * USDG);
    }

    function test_BidsClearAtOneUniformPriceThatTheContractChecks() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _position(bob, 5 * SHARE, 800 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        auction.enroll(bob, CROSS, NVDA);
        uint256 supply = _need(1_600 * USDG, 171e8) + _need(800 * USDG, 171e8);
        _commit(carol, 8 * SHARE, 185e8, "c");
        _commit(dave, 5 * SHARE, 180e8, "d");
        _commit(eve, 5 * SHARE, 175e8, "e");
        _toReveal();
        _reveal(carol, 8 * SHARE, 185e8, "c");
        _reveal(dave, 5 * SHARE, 180e8, "d");
        _reveal(eve, 5 * SHARE, 175e8, "e");
        assertEq(usdg.balanceOf(carol), 100 * USDG + 10_000 * USDG - 1_480 * USDG);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.clear(NVDA, openMs, 175e8);
        _open();
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 180e8));
        auction.clear(NVDA, openMs, 180e8);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 174e8));
        auction.clear(NVDA, openMs, 174e8);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 171e8));
        auction.clear(NVDA, openMs, 171e8);
        uint256 reserve = accounts.reserve();
        auction.clear(NVDA, openMs, 175e8);
        vm.expectRevert(ReopeningAuction.WrongState.selector);
        auction.clear(NVDA, openMs, 175e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertTrue(r.cleared);
        assertEq(r.supply, supply);
        assertEq(r.price, 175e8);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.debt(bob, CROSS), 0);
        assertEq(
            r.sold,
            15 * SHARE - accounts.collateral(alice, CROSS, address(nvda))
                - accounts.collateral(bob, CROSS, address(nvda))
        );
        assertApproxEqAbs(accounts.reserve() - reserve, r.paid * 50 / 10_000, 2);
        assertEq(liquidator.heldUntil(alice, CROSS), 0);
        uint256 took = _claimed(0, carol, 1_480 * USDG, 8 * SHARE * r.sold / supply);
        took += _claimed(1, dave, 900 * USDG, 5 * SHARE * r.sold / supply);
        took += _claimed(2, eve, 875 * USDG, (supply - 13 * SHARE) * r.sold / supply);
        assertLe(took, r.sold);
        vm.expectRevert(ReopeningAuction.WrongState.selector);
        auction.claim(NVDA, openMs, 0);
        assertLe(usdg.balanceOf(address(auction)), 3);
    }

    function test_AnUndersubscribedRoundClearsAtItsLowestBidAndReleasesTheRest() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 2 * SHARE, 185e8, "c");
        _commit(dave, 1 * SHARE, 180e8, "d");
        _toReveal();
        _reveal(carol, 2 * SHARE, 185e8, "c");
        _reveal(dave, 1 * SHARE, 180e8, "d");
        _open();
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 171e8));
        auction.clear(NVDA, openMs, 171e8);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 185e8));
        auction.clear(NVDA, openMs, 185e8);
        auction.clear(NVDA, openMs, 180e8);
        uint256 before = usdg.balanceOf(carol);
        auction.claim(NVDA, openMs, 0);
        assertEq(nvda.balanceOf(carol), 2 * SHARE);
        assertEq(usdg.balanceOf(carol) - before, 10 * USDG);
        auction.claim(NVDA, openMs, 1);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 7 * SHARE);
        assertEq(liquidator.heldUntil(alice, CROSS), 0);
    }

    function test_ARevealMustMatchItsCommitmentItsBidderAndItsPhase() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        bytes32 commitment = _commit(carol, 2 * SHARE, 185e8, "c");
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.InvalidCommitment.selector, commitment));
        vm.prank(carol);
        auction.commit(NVDA, commitment, 370 * USDG);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.InvalidCommitment.selector, bytes32(uint256(2))));
        vm.prank(carol);
        auction.commit(NVDA, bytes32(uint256(2)), 99 * USDG);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 2 * SHARE, 185e8, "c");
        _toReveal();
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        vm.prank(dave);
        auction.commit(NVDA, bytes32(uint256(1)), 100 * USDG);
        _fund(dave, 10_000 * USDG);
        bytes32 copied = keccak256(abi.encode(dave, NVDA, openMs, 2 * SHARE, uint256(185e8), bytes32("c")));
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownCommitment.selector, copied));
        vm.prank(dave);
        auction.reveal(NVDA, openMs, 2 * SHARE, 185e8, "c");
        bytes32 wrong = keccak256(abi.encode(carol, NVDA, openMs, 2 * SHARE, uint256(185e8), bytes32("x")));
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownCommitment.selector, wrong));
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 2 * SHARE, 185e8, "x");
        _open();
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 2 * SHARE, 185e8, "c");
    }

    function test_BidsBelowTheFloorOrTheMinimumAreRefused() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 2 * SHARE, 170e8, "a");
        _commit(carol, 0.5e18, 180e8, "b");
        _commit(carol, 0.58e18, 10_000e8, "d");
        _toReveal();
        vm.expectRevert(ReopeningAuction.InvalidBid.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 2 * SHARE, 170e8, "a");
        vm.expectRevert(ReopeningAuction.InvalidBid.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 0.5e18, 180e8, "b");
        vm.expectRevert(ReopeningAuction.InvalidBid.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 0.58e18, 10_000e8, "d");
    }

    function test_AFullBookTakesAHigherBidInPlaceOfItsLowest() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        for (uint256 i; i < 64; ++i) {
            _commit(address(uint160(0x1000 + i)), _bookQuantity(i), _bookPrice(i), bytes32(i));
        }
        _commit(carol, 1 * SHARE, 172e8, "c");
        _commit(eve, 2 * SHARE, 172e8, "e");
        _commit(dave, 1 * SHARE, 173e8, "d");
        _toReveal();
        for (uint256 i; i < 64; ++i) {
            _reveal(address(uint160(0x1000 + i)), _bookQuantity(i), _bookPrice(i), bytes32(i));
        }
        uint256 pool = auction.round(NVDA, openMs).pool;
        uint256 deposits = auction.round(NVDA, openMs).deposits;
        uint256 refunded = usdg.balanceOf(carol);
        vm.expectEmit(address(auction));
        emit ReopeningAuction.Outbid(NVDA, openMs, carol, 1 * SHARE, 172e8);
        _reveal(carol, 1 * SHARE, 172e8, "c");
        assertEq(usdg.balanceOf(carol) - refunded, 172 * USDG);
        assertEq(auction.round(NVDA, openMs).pool, pool);
        assertEq(auction.round(NVDA, openMs).deposits, deposits - 172 * USDG);
        bytes32 outbid = keccak256(abi.encode(carol, NVDA, openMs, 1 * SHARE, uint256(172e8), bytes32("c")));
        (, uint128 left) = auction.commitments(NVDA, openMs, outbid);
        assertEq(left, 0);
        uint256 before = usdg.balanceOf(address(0x1005));
        vm.expectEmit(address(auction));
        emit ReopeningAuction.Evicted(NVDA, openMs, address(0x1005), 5, 1 * SHARE, 172e8);
        _reveal(eve, 2 * SHARE, 172e8, "e");
        assertEq(usdg.balanceOf(address(0x1005)) - before, 172 * USDG);
        before = usdg.balanceOf(address(0x1004));
        vm.expectEmit(address(auction));
        emit ReopeningAuction.Evicted(NVDA, openMs, address(0x1004), 4, 2 * SHARE, 172e8);
        _reveal(dave, 1 * SHARE, 173e8, "d");
        assertEq(usdg.balanceOf(address(0x1004)) - before, 344 * USDG);
        assertEq(auction.round(NVDA, openMs).pool, pool - 172 * USDG + 173 * USDG);
        ReopeningAuction.Bid[] memory bids = auction.bids(NVDA, openMs);
        assertEq(bids.length, 64);
        assertEq(bids[4].bidder, dave);
        assertEq(bids[5].bidder, eve);
    }

    function test_AnUnrevealedCommitmentLosesPartOfItsDeposit() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 2 * SHARE, 185e8, "c");
        bytes32 small = _commit(dave, 2 * SHARE, 185e8, "d");
        bytes32 large = keccak256(abi.encode(eve, NVDA, openMs, 10 * SHARE, uint256(185e8), bytes32("e")));
        _fund(eve, 5_000 * USDG);
        vm.prank(eve);
        auction.commit(NVDA, large, 5_000 * USDG);
        _toReveal();
        _reveal(carol, 2 * SHARE, 185e8, "c");
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.forfeit(NVDA, openMs, small);
        _open();
        uint256 reserve = accounts.reserve();
        uint256 before = usdg.balanceOf(dave);
        vm.expectEmit(address(auction));
        emit ReopeningAuction.Forfeited(NVDA, openMs, dave, 100 * USDG, 270 * USDG);
        auction.forfeit(NVDA, openMs, small);
        assertEq(usdg.balanceOf(dave) - before, 270 * USDG);
        auction.forfeit(NVDA, openMs, large);
        assertEq(accounts.reserve(), reserve + 100 * USDG + 500 * USDG);
        assertEq(usdg.balanceOf(eve), 4_500 * USDG);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownCommitment.selector, small));
        auction.forfeit(NVDA, openMs, small);
        assertEq(auction.round(NVDA, openMs).deposits, 0);
    }

    function test_ARoundNobodyClearsRefundsItsBidsAndReleasesItsLots() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 2 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 2 * SHARE, 185e8, "c");
        _open();
        vm.expectRevert(ReopeningAuction.WrongState.selector);
        auction.claim(NVDA, openMs, 0);
        skip(1 hours);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.clear(NVDA, openMs, 171e8);
        uint256 before = usdg.balanceOf(carol);
        auction.claim(NVDA, openMs, 0);
        assertEq(usdg.balanceOf(carol) - before, 370 * USDG);
        liquidator.start(alice, CROSS);
    }

    function test_APositionThatRecoveredBeforeTheClearingSellsNothing() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 12 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 12 * SHARE, 185e8, "c");
        _open();
        _quote(NVDA, 200e8, 202e8);
        auction.clear(NVDA, openMs, 185e8);
        uint256 sold = auction.round(NVDA, openMs).sold;
        assertEq(sold, 0);
        uint256 before = usdg.balanceOf(carol);
        auction.claim(NVDA, openMs, 0);
        assertEq(nvda.balanceOf(carol), 0);
        assertEq(usdg.balanceOf(carol) - before, 2_220 * USDG);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
    }

    function test_AMidweekHaltsEndOpensNoRoundAndTheDutchAuctionSellsAfterIt() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _quote(NVDA, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        _quote(NVDA, 190e8, 192e8);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.enroll(alice, CROSS, NVDA);
        liquidator.start(alice, CROSS);
    }

    function testFuzz_OnlyTheClearingPriceClearsAndNoBidderPaysAboveItsPrice(uint256 seed) public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _position(bob, 5 * SHARE, 800 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        auction.enroll(bob, CROSS, NVDA);
        uint256 supply = auction.round(NVDA, openMs).supply;
        uint256 n = bound(seed, 1, 6);
        uint256[] memory quantities = new uint256[](n);
        uint256[] memory prices = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            quantities[i] = bound(uint256(keccak256(abi.encode(seed, i, "q"))), 1 * SHARE, 10 * SHARE);
            prices[i] = 171e8 + bound(uint256(keccak256(abi.encode(seed, i, "p"))), 0, 5) * 3e8;
            _commit(address(uint160(0x2000 + i)), quantities[i], prices[i], bytes32(i));
        }
        _toReveal();
        for (uint256 i; i < n; ++i) {
            _reveal(address(uint160(0x2000 + i)), quantities[i], prices[i], bytes32(i));
        }
        _open();
        uint256 clearing = type(uint256).max;
        for (uint256 i; i < n; ++i) {
            if (prices[i] < clearing) clearing = prices[i];
        }
        for (uint256 step = 5; step != type(uint256).max; --step) {
            uint256 p = 171e8 + step * 3e8;
            uint256 demand;
            for (uint256 i; i < n; ++i) {
                if (prices[i] >= p) demand += quantities[i];
            }
            if (demand >= supply) {
                clearing = p;
                break;
            }
            if (step == 0) break;
        }
        for (uint256 step; step <= 5; ++step) {
            uint256 p = 171e8 + step * 3e8;
            if (p == clearing) continue;
            vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, p));
            auction.clear(NVDA, openMs, p);
        }
        auction.clear(NVDA, openMs, clearing);
        uint256 sold = auction.round(NVDA, openMs).sold;
        uint256 taken;
        for (uint256 i; i < n; ++i) {
            address bidder = address(uint160(0x2000 + i));
            uint256 before = usdg.balanceOf(bidder);
            auction.claim(NVDA, openMs, i);
            uint256 got = nvda.balanceOf(bidder);
            uint256 escrow = (quantities[i] * prices[i] - 1) / 1e20 + 1;
            assertLe(got, quantities[i]);
            assertLe(escrow - (usdg.balanceOf(bidder) - before), got * clearing / 1e20 + 3);
            if (prices[i] < clearing) assertEq(got, 0);
            taken += got;
        }
        assertLe(taken, sold);
        assertLe(sold - taken, 2 * n);
    }

    function test_ALotIsTheWholeHoldingWhenThatRepaysLess() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(150e8);
        _quote(NVDA, 150e8, 152e8);
        auction.enroll(alice, CROSS, NVDA);
        assertEq(auction.round(NVDA, openMs).floor, 135e8);
        assertEq(auction.lots(NVDA, openMs)[0].amount, 10 * SHARE);
    }

    function test_ALotCountsTheStockTokensOfThePositionsBaskets() public {
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = NVDA;
        uint256[] memory units = new uint256[](1);
        units[0] = 1e18;
        Basket basket = new Basket("Tapehouse Test Basket", "thTEST", IBand(address(band)), symbols, units, owner);
        vm.prank(owner);
        accounts.addBasket(basket);
        nvda.mint(alice, 7 * SHARE);
        units[0] = 7 * SHARE;
        vm.startPrank(alice);
        nvda.approve(address(basket), 7 * SHARE);
        basket.mint(7 * SHARE, alice, units);
        basket.transfer(bob, SHARE);
        basket.approve(address(accounts), 6 * SHARE);
        accounts.deposit(CROSS, address(basket), 6 * SHARE, alice);
        vm.stopPrank();
        units[0] = 6 * SHARE;
        vm.startPrank(bob);
        basket.approve(address(accounts), SHARE);
        accounts.deposit(CROSS, address(basket), SHARE, bob);
        vm.stopPrank();
        _position(alice, 4 * SHARE, 1_600 * USDG);
        _position(carol, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotEnrollable.selector, bob, CROSS));
        auction.enroll(bob, CROSS, NVDA);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unwrap(address(auction), alice, address(basket), 6 * SHARE, units);
        auction.enroll(alice, CROSS, NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 0);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        auction.enroll(carol, CROSS, NVDA);
        ReopeningAuction.Lot[] memory lots = auction.lots(NVDA, openMs);
        assertEq(lots[0].amount, _need(1_600 * USDG, 171e8));
        assertEq(lots[1].amount, lots[0].amount);
    }

    function test_ALotCountsWhatThePositionLentAndItsClearingTakesItFromTheVault() public {
        StockLendingVault lending =
            new StockLendingVault(IERC20(address(nvda)), owner, SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.startPrank(owner);
        lending.setDepositor(address(accounts));
        accounts.setLending(NVDA, lending);
        vm.stopPrank();
        _position(alice, 10 * SHARE, 1_200 * USDG);
        vm.prank(alice);
        accounts.lend(CROSS, address(nvda), 10 * SHARE, alice);
        _weekend(150e8);
        _quote(NVDA, 150e8, 152e8);
        auction.enroll(alice, CROSS, NVDA);
        uint256 lot = auction.lots(NVDA, openMs)[0].amount;
        assertEq(lot, _need(1_200 * USDG, 135e8));
        _commit(carol, lot, 140e8, "c");
        _toReveal();
        _reveal(carol, lot, 140e8, "c");
        _open();
        auction.clear(NVDA, openMs, 140e8);
        auction.claim(NVDA, openMs, 0);
        assertEq(nvda.balanceOf(carol) + accounts.lent(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.debt(alice, CROSS), 0);
    }

    function test_ALotIncludesThePremiumOwed() public {
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        uint256 owed = accounts.debt(alice, CROSS) + accounts.premium(alice, CROSS);
        assertGt(owed, accounts.debt(alice, CROSS));
        auction.enroll(alice, CROSS, NVDA);
        assertEq(auction.lots(NVDA, openMs)[0].amount, _need(owed, 171e8));
    }

    function test_ALargerLotTakesThePlaceOfTheSmallestAndReleasesIt() public {
        for (uint256 i; i < 33; ++i) {
            _position(address(uint160(0x4000 + i)), 1 * SHARE, 160 * USDG - i * 0.2e6);
        }
        _position(alice, 2 * SHARE, 320 * USDG);
        _position(bob, 1 * SHARE, 154 * USDG);
        _weekend(190e8);
        for (uint256 i; i < 32; ++i) {
            auction.enroll(address(uint160(0x4000 + i)), CROSS, NVDA);
        }
        uint256 supply = auction.round(NVDA, openMs).supply;
        vm.expectRevert(ReopeningAuction.TooManyLots.selector);
        auction.enroll(address(uint160(0x4000 + 32)), CROSS, NVDA);
        vm.expectEmit(address(auction));
        emit ReopeningAuction.LotEvicted(NVDA, openMs, address(0x401f), CROSS);
        auction.enroll(alice, CROSS, NVDA);
        ReopeningAuction.Lot[] memory lots = auction.lots(NVDA, openMs);
        assertEq(lots.length, 32);
        assertEq(lots[31].account, alice);
        assertEq(auction.round(NVDA, openMs).supply, supply - _need(153.8e6, 171e8) + _need(320 * USDG, 171e8));
        assertFalse(auction.enrolled(openMs, address(0x401f), CROSS, NVDA));
        assertEq(liquidator.heldUntil(address(0x401f), CROSS), 0);
        vm.expectRevert(ReopeningAuction.TooManyLots.selector);
        auction.enroll(bob, CROSS, NVDA);
        assertEq(liquidator.heldUntil(alice, CROSS), openMs / 1000 + 1 hours);
    }

    function test_APositionInTwoRoundsStaysHeldUntilBothAreDone() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(spy), 1 * SHARE);
        vm.prank(alice);
        accounts.borrow(CROSS, 2_050 * USDG, alice, alice);
        _weekend(190e8);
        _quote(SPY, 590e8, 596e8);
        auction.enroll(alice, CROSS, NVDA);
        auction.enroll(alice, CROSS, SPY);
        assertTrue(auction.enrolled(openMs, alice, CROSS, SPY));
        _commit(carol, 1 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 1 * SHARE, 185e8, "c");
        _open();
        auction.clear(NVDA, openMs, 185e8);
        assertEq(auction.round(NVDA, openMs).sold, 1 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(Liquidator.PositionHeld.selector, alice, CROSS, uint64(openMs / 1000 + 1 hours))
        );
        liquidator.start(alice, CROSS);
        auction.clear(SPY, openMs, auction.round(SPY, openMs).floor);
        assertEq(liquidator.heldUntil(alice, CROSS), 0);
        liquidator.start(alice, CROSS);
    }

    function test_TheLastRefundIsWhatTheRoundHoldsWhenAnEscrowCapsAPayment() public {
        uint24[8] memory debts = [uint24(155_232), 153_770, 156_662, 152_131, 154_374, 157_306, 157_000, 158_016];
        for (uint256 i; i < 8; ++i) {
            _position(address(uint160(0x7000 + i)), 1 * SHARE, uint256(debts[i]) * 1e3);
        }
        _weekend(190e8);
        for (uint256 i; i < 8; ++i) {
            auction.enroll(address(uint160(0x7000 + i)), CROSS, NVDA);
        }
        _commit(carol, 1 * SHARE, 180e8, "c");
        _commit(dave, 1 * SHARE, 200e8, "d");
        _toReveal();
        _reveal(carol, 1 * SHARE, 180e8, "c");
        _reveal(dave, 1 * SHARE, 200e8, "d");
        _open();
        auction.clear(NVDA, openMs, 180e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.paid, 360_000_006);
        uint256 before = usdg.balanceOf(carol);
        auction.claim(NVDA, openMs, 0);
        assertEq(usdg.balanceOf(carol), before);
        uint256 left = auction.round(NVDA, openMs).pool;
        before = usdg.balanceOf(dave);
        auction.claim(NVDA, openMs, 1);
        assertEq(usdg.balanceOf(dave) - before, left);
        assertEq(auction.round(NVDA, openMs).pool, 0);
        assertEq(usdg.balanceOf(address(auction)), 0);
    }

    function test_ARevealedCommitmentCannotBeRevealedAgain() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        bytes32 commitment = _commit(carol, 2 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 2 * SHARE, 185e8, "c");
        (, uint128 deposit) = auction.commitments(NVDA, openMs, commitment);
        assertEq(deposit, 0);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownCommitment.selector, commitment));
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 2 * SHARE, 185e8, "c");
    }

    function test_ARevealRefundsItsDepositAboveTheEscrowAndRefusesABidItDoesNotCover() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _fund(carol, 500 * USDG);
        bytes32 covered = keccak256(abi.encode(carol, NVDA, openMs, 2 * SHARE, uint256(185e8), bytes32("c")));
        bytes32 uncovered = keccak256(abi.encode(carol, NVDA, openMs, 3 * SHARE, uint256(185e8), bytes32("d")));
        vm.startPrank(carol);
        auction.commit(NVDA, covered, 400 * USDG);
        auction.commit(NVDA, uncovered, 100 * USDG);
        vm.stopPrank();
        assertEq(auction.round(NVDA, openMs).deposits, 500 * USDG);
        _toReveal();
        _reveal(carol, 2 * SHARE, 185e8, "c");
        assertEq(usdg.balanceOf(carol), 30 * USDG);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.deposits, 100 * USDG);
        assertEq(r.pool, 370 * USDG);
        vm.expectRevert(ReopeningAuction.InvalidBid.selector);
        vm.prank(carol);
        auction.reveal(NVDA, openMs, 3 * SHARE, 185e8, "d");
    }

    function test_ThePhaseIsTheNightBeforeARegularOpenAfterAClosure() public {
        uint64 nowMs = uint64(vm.getBlockTimestamp()) * 1000;
        uint64 change = nowMs + 2 hours * 1000;
        band.setSession(BandDouble.Session(2, 3, 1, change, 0));
        (uint64 open, bool revealing) = auction.phase();
        assertEq(open, change);
        assertFalse(revealing);
        band.setSession(BandDouble.Session(2, 3, 1, nowMs + 30 minutes * 1000, 0));
        (, revealing) = auction.phase();
        assertTrue(revealing);
        band.setSession(BandDouble.Session(2, 3, 1, change, nowMs + 1 hours * 1000));
        (open,) = auction.phase();
        assertEq(open, 0);
        band.setSession(BandDouble.Session(2, 2, 1, change, 0));
        (open,) = auction.phase();
        assertEq(open, 0);
        band.setSession(BandDouble.Session(1, 3, 1, change, 0));
        (open,) = auction.phase();
        assertEq(open, 0);
        band.setSession(BandDouble.Session(2, 3, 1, nowMs, 0));
        (open,) = auction.phase();
        assertEq(open, 0);
        band.setSession(BandDouble.Session(2, 3, 0, change, 0));
        (open,) = auction.phase();
        assertEq(open, 0);
    }

    function test_AfterAHolidayTheRoundsThatOpenedBeforeMidnightRunUntilTheOpen() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _quote(NVDA, 190e8, 192e8);
        uint64 midnight = uint64(vm.getBlockTimestamp() + 4 hours) * 1000;
        uint64 holidayOpen = midnight + 9.5 hours * 1000;
        band.setSession(BandDouble.Session(2, 3, 2, midnight, 0));
        (uint64 open,) = auction.phase();
        assertEq(open, holidayOpen);
        vm.warp(midnight / 1000);
        band.setSession(BandDouble.Session(2, 2, 0, 0, 0));
        (open,) = auction.phase();
        assertEq(open, 0);
        vm.warp(midnight / 1000 - 1 hours);
        band.setSession(BandDouble.Session(2, 3, 2, midnight, 0));
        auction.enroll(alice, CROSS, NVDA);
        assertEq(auction.lastOpenMs(), holidayOpen);
        vm.warp(midnight / 1000);
        band.setSession(BandDouble.Session(2, 2, 1, holidayOpen, 0));
        (open,) = auction.phase();
        assertEq(open, holidayOpen);
        _commit(carol, 2 * SHARE, 185e8, "c");
        vm.warp(holidayOpen / 1000);
        (open,) = auction.phase();
        assertEq(open, 0);
    }

    function test_AHaltedSealIsNotTheFloor() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(0);
        seal.setSeal(reopenMs, SealDouble.Seal(0, 0, 0, 0, 100e8, 0, uint64(vm.getBlockTimestamp())));
        vm.expectEmit(address(auction));
        emit ReopeningAuction.RoundOpened(NVDA, openMs, reopenMs, 171e8, false);
        auction.enroll(alice, CROSS, NVDA);
    }

    function test_TheClearingWindowNeedsTheRegularOpen() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        vm.warp(openMs / 1000);
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.clear(NVDA, openMs, 171e8);
        vm.warp(openMs / 1000 - 1);
        band.setSession(BandDouble.Session(2, 1, 2, openMs + 6.5 hours * 1000, 0));
        vm.expectRevert(ReopeningAuction.WrongPhase.selector);
        auction.clear(NVDA, openMs, 171e8);
    }

    function test_AnAssetWithoutAStockTokenHasNoRound() public {
        bytes32 tsla = "TSLA";
        band.setAsset(tsla, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(0)));
        _quote(tsla, 300e8, 303e8);
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (NVDA, tsla);
        MarginDouble engine = new MarginDouble(symbols, address(band), address(0));
        SupplyVault other = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        MarginAccounts bare = new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), other, weth, owner, new uint256[](2), 0, 0, 0, 10_00
        );
        ReopeningAuction lone = new ReopeningAuction(new Liquidator(bare), new bytes32[](0), new BandFeed[](0));
        _weekend(0);
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.UnknownAsset.selector, tsla));
        lone.commit(tsla, bytes32(uint256(1)), 100 * USDG);
    }

    function test_TheLargestRoundClearsInOneBlock() public {
        for (uint256 i; i < 32; ++i) {
            _position(address(uint160(0x5000 + i)), 1 * SHARE, 160 * USDG);
        }
        _weekend(190e8);
        for (uint256 i; i < 32; ++i) {
            auction.enroll(address(uint160(0x5000 + i)), CROSS, NVDA);
        }
        for (uint256 i; i < 64; ++i) {
            _commit(address(uint160(0x6000 + i)), 1 * SHARE, 172e8 + i * 1e7, bytes32(i));
        }
        _toReveal();
        for (uint256 i; i < 64; ++i) {
            _reveal(address(uint160(0x6000 + i)), 1 * SHARE, 172e8 + i * 1e7, bytes32(i));
        }
        _open();
        uint256 supply = auction.round(NVDA, openMs).supply;
        uint256 marginal = 64 - (supply - 1) / SHARE - 1;
        auction.clear(NVDA, openMs, 172e8 + marginal * 1e7);
        vm.snapshotGasLastCall("clearLargest");
        assertLt(vm.lastCallGas().gasTotalUsed, 32_000_000);
        assertGt(auction.round(NVDA, openMs).sold, supply * 9 / 10);
    }

    function test_GasOfEachCall() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _position(bob, 5 * SHARE, 800 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        vm.snapshotGasLastCall("enroll");
        auction.enroll(bob, CROSS, NVDA);
        _commit(carol, 8 * SHARE, 185e8, "c");
        vm.snapshotGasLastCall("commit");
        _commit(dave, 10 * SHARE, 175e8, "d");
        _toReveal();
        _reveal(carol, 8 * SHARE, 185e8, "c");
        vm.snapshotGasLastCall("reveal");
        _reveal(dave, 10 * SHARE, 175e8, "d");
        _open();
        auction.clear(NVDA, openMs, 175e8);
        vm.snapshotGasLastCall("clear");
        auction.claim(NVDA, openMs, 0);
        vm.snapshotGasLastCall("claim");
    }

    function _bookQuantity(uint256 i) internal pure returns (uint256) {
        return i == 4 ? 2 * SHARE : 1 * SHARE;
    }

    function _bookPrice(uint256 i) internal pure returns (uint256) {
        return i == 5 ? 172e8 : 172e8 + ((i + 60) % 64) * 1e8;
    }

    function deploy(bytes32[] memory symbols, BandFeed[] memory feeds) external returns (ReopeningAuction) {
        return new ReopeningAuction(liquidator, symbols, feeds);
    }

    function _claimed(uint256 index, address bidder, uint256 escrow, uint256 expected) internal returns (uint256 took) {
        uint256 before = usdg.balanceOf(bidder);
        auction.claim(NVDA, openMs, index);
        took = nvda.balanceOf(bidder);
        assertEq(took, expected);
        assertApproxEqAbs(escrow - (usdg.balanceOf(bidder) - before), took * 175e8 / 1e20, 2);
    }
}

contract ReopeningAuctionSettlementTest is ReopeningAuctionBase {
    function test_OnlyTheAccountsOwnerSetsTheAuctionOnceAndOnlyItHoldsSettlesAndCollects() public {
        Liquidator other = new Liquidator(accounts);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotAccountsOwner.selector, alice));
        vm.prank(alice);
        other.setAuction(alice);
        vm.expectRevert(Liquidator.InvalidAuction.selector);
        vm.prank(owner);
        other.setAuction(address(0));
        vm.expectEmit(address(other));
        emit Liquidator.AuctionSet(alice);
        vm.prank(owner);
        other.setAuction(alice);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.AuctionAlreadySet.selector, alice));
        vm.prank(owner);
        other.setAuction(bob);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotAuction.selector, bob));
        vm.prank(bob);
        liquidator.hold(alice, CROSS, 1);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotAuction.selector, bob));
        vm.prank(bob);
        liquidator.settle(alice, CROSS, address(nvda), 1, 1);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotAuction.selector, bob));
        vm.prank(bob);
        liquidator.collect(1);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _quote(NVDA, 190e8, 192e8);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotForSale.selector, address(usdg)));
        vm.prank(address(auction));
        liquidator.settle(alice, CROSS, address(usdg), 1, 1);
    }

    function test_AWriteOffSweepsHoldingsWorthLessThanAUsdgToItsCaller() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _deposit(alice, CROSS, address(usdg), 0.5e6);
        _quote(NVDA, 100e8, 101e8);
        liquidator.start(alice, CROSS);
        usdg.mint(carol, 10_000 * USDG);
        vm.startPrank(carol);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE - 1e15, type(uint256).max, carol);
        vm.stopPrank();
        weth.mint(dave, 1e14);
        vm.startPrank(dave);
        weth.approve(address(accounts), 1e14);
        accounts.deposit(CROSS, address(weth), 1e14, alice);
        vm.stopPrank();
        vm.prank(dave);
        liquidator.writeOff(alice, CROSS);
        assertEq(nvda.balanceOf(dave), 1e15);
        assertEq(weth.balanceOf(dave), 1e14);
        assertEq(usdg.balanceOf(dave), 0.5e6);
        assertEq(accounts.debt(alice, CROSS), 0);
        _quote(NVDA, 200e8, 202e8);
        _position(bob, 10 * SHARE, 1_600 * USDG);
        _quote(NVDA, 100e8, 101e8);
        liquidator.start(bob, CROSS);
        vm.prank(carol);
        liquidator.buy(bob, CROSS, address(nvda), 10 * SHARE - 0.01e18, type(uint256).max, carol);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(dave);
        liquidator.writeOff(bob, CROSS);
    }

    function test_AWriteOffLeavesAHoldingWorthAUsdgOrMore() public {
        _emptied(alice);
        weth.mint(alice, 1e15);
        vm.startPrank(alice);
        weth.approve(address(accounts), 1e15);
        accounts.deposit(CROSS, address(weth), 1e15, alice);
        vm.stopPrank();
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        liquidator.writeOff(alice, CROSS);
        _emptied(bob);
        _deposit(bob, CROSS, address(usdg), 2 * USDG);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        liquidator.writeOff(bob, CROSS);
        _emptiedBut(carol, 1e12);
        _quote(NVDA, 0, 0);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        liquidator.writeOff(carol, CROSS);
    }

    function test_ASettlementsRoundingExcessGoesToTheReserve() public {
        _deposit(alice, CROSS, address(nvda), 20 * SHARE);
        vm.prank(alice);
        accounts.borrow(CROSS, 150 * USDG, alice, alice);
        _weekend(9e8);
        _quote(NVDA, 9e8, 10e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 20 * SHARE, 10e8, "c");
        _toReveal();
        _reveal(carol, 20 * SHARE, 10e8, "c");
        _open();
        uint256 reserve = accounts.reserve();
        auction.clear(NVDA, openMs, 10e8);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(usdg.balanceOf(alice), 150 * USDG);
        assertEq(accounts.reserve() - reserve, 753_768 + 1);
    }

    function test_WhileTheMarketIsClosedAHeldPositionIsStillSold() public {
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        liquidator.start(alice, CROSS);
        usdg.mint(carol, 1_000 * USDG);
        vm.startPrank(carol);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 1 * SHARE, type(uint256).max, carol);
        vm.stopPrank();
        assertEq(nvda.balanceOf(carol), 1 * SHARE);
    }

    function test_AWriteOffTakesNothingFromAPositionWorthMoreThanItOwes() public {
        _deposit(alice, CROSS, address(usdg), 0.5e6);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(dave);
        liquidator.writeOff(alice, CROSS);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 0.5e6);
    }

    function test_AWriteOffKeepsWethWhileTheEthPriceIsStale() public {
        _emptied(alice);
        weth.mint(alice, 1e14);
        vm.startPrank(alice);
        weth.approve(address(accounts), 1e14);
        accounts.deposit(CROSS, address(weth), 1e14, alice);
        vm.stopPrank();
        ethUsd.setRound(2_000e8, block.timestamp - liquidator.MAX_FEED_AGE() - 1);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(eve);
        liquidator.writeOff(alice, CROSS);
        ethUsd.setRound(2_000e8, block.timestamp);
        vm.prank(eve);
        liquidator.writeOff(alice, CROSS);
        assertEq(weth.balanceOf(eve), 1e14);
    }

    function test_TheBackstopBuysWhatTheBidsLeaveAndSharesItAsGains() public {
        uint256[] memory limits = new uint256[](2);
        (limits[0], limits[1]) = (5_000 * USDG, 5_000 * USDG);
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 3_000 * USDG);
        _back(backstop, eve, 1_000 * USDG);
        vm.expectRevert(abi.encodeWithSelector(GapBackstop.NotAuction.selector, address(this)));
        backstop.buyRemainder(alice, CROSS, address(nvda), 1, 1);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 1 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 1 * SHARE, 185e8, "c");
        _open();
        auction.clear(NVDA, openMs, 185e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.sold, 1 * SHARE);
        assertGt(r.taken, 0);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(nvda.balanceOf(address(backstop)), r.taken);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 9 * SHARE - r.taken);
        uint256 cost = (r.taken * 185e8 - 1) / 1e20 + 1;
        assertEq(backstop.held(), 4_000 * USDG - cost);
        assertEq(usdg.allowance(address(backstop), address(liquidator)), 0);
        assertEq(backstop.gains(bob, address(nvda)), r.taken * 3 / 4);
        assertEq(backstop.gains(eve, address(nvda)), r.taken / 4);
        uint256 shares = backstop.balanceOf(eve);
        vm.prank(eve);
        assertTrue(backstop.transfer(bob, shares));
        assertEq(backstop.gains(eve, address(nvda)), r.taken / 4);
        assertEq(backstop.gains(bob, address(nvda)), r.taken * 3 / 4);
        vm.expectEmit(address(backstop));
        emit GapBackstop.GainsClaimed(eve, address(nvda), r.taken / 4);
        vm.prank(eve);
        backstop.claimGains(address(nvda));
        vm.prank(bob);
        backstop.claimGains(address(nvda));
        assertEq(nvda.balanceOf(eve), r.taken / 4);
        assertEq(nvda.balanceOf(bob), r.taken * 3 / 4);
        assertEq(backstop.gains(eve, address(nvda)), 0);
        vm.prank(eve);
        assertEq(backstop.claimGains(address(nvda)), 0);
    }

    function test_TheBackstopBuysNoMoreThanItsExposureLimit() public {
        GapBackstop backstop = _backstop(300 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        ReopeningAuction.Round memory r = _thinRound();
        uint256 cost = (r.taken * 185e8 - 1) / 1e20 + 1;
        assertEq(r.taken, 300 * USDG * 1e20 / 185e8);
        assertEq(backstop.covered(CROSS, backstop.closureMs()), cost);
        assertEq(backstop.held(), 4_000 * USDG - cost);
        assertGt(accounts.debt(alice, CROSS), 0);
        vm.prank(bob);
        assertApproxEqAbs(backstop.claimGains(address(nvda)), r.taken, 1);
    }

    function test_TheBackstopBuysNothingOnceTheBidsCureThePosition() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 5 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 5 * SHARE, 185e8, "c");
        _open();
        auction.clear(NVDA, openMs, 185e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertGt(r.sold, 0);
        assertEq(r.taken, 0);
        assertEq(backstop.held(), 4_000 * USDG);
        assertEq(backstop.gainsEpoch(), 0);
    }

    function test_TheBackstopPaysAtMostTheBandsLowEdge() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(170e8);
        _quote(NVDA, 170e8, 172e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(alice, 0.7e18, 300e8, "a");
        _toReveal();
        _reveal(alice, 0.7e18, 300e8, "a");
        _open();
        auction.clear(NVDA, openMs, 300e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.price, 300e8);
        assertGt(r.taken, 0);
        assertEq(4_000 * USDG - backstop.held(), (r.taken * 170e8 - 1) / 1e20 + 1);
        _quote(NVDA, 0, 0);
        assertEq(liquidator.heldUntil(alice, CROSS), 0);
    }

    function test_ABackstopWithNoRoomLeftBuysNothing() public {
        GapBackstop backstop = _backstop(0);
        _back(backstop, bob, 4_000 * USDG);
        ReopeningAuction.Round memory r = _thinRound();
        assertEq(r.sold, 1 * SHARE);
        assertEq(r.taken, 0);
        assertEq(backstop.held(), 4_000 * USDG);
        assertEq(backstop.gainsEpoch(), 0);
    }

    function test_ABackstopThatCannotBuyNeverBlocksTheClearing() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        nvda.blockAccount(address(backstop), true);
        ReopeningAuction.Round memory r = _thinRound();
        assertTrue(r.cleared);
        assertEq(r.sold, 1 * SHARE);
        assertEq(r.taken, 0);
        assertEq(backstop.held(), 4_000 * USDG);
    }

    function test_NoOneDepositsIntoTheBackstopWhileAClosureSettles() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        assertEq(backstop.maxDeposit(eve), 0);
        _open();
        assertEq(backstop.maxDeposit(eve), 0);
        usdg.mint(eve, 1_000 * USDG);
        vm.startPrank(eve);
        usdg.approve(address(backstop), 1_000 * USDG);
        vm.expectRevert(abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxDeposit.selector, eve, 1_000 * USDG, 0));
        backstop.deposit(1_000 * USDG, eve);
        vm.stopPrank();
        skip(1 days);
        assertEq(backstop.maxDeposit(eve), type(uint256).max);
    }

    function test_AFrozenBackstopBuysNothingAndTheRoundStillClears() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        usdg.freeze(address(backstop));
        ReopeningAuction.Round memory r = _thinRound();
        assertTrue(r.cleared);
        assertEq(r.sold, 1 * SHARE);
        assertEq(r.taken, 0);
    }

    function test_ABackstopWithoutSharesBuysNothing() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        _deposit(dave, SPY, address(spy), 10 * SHARE);
        vm.prank(dave);
        accounts.borrow(SPY, 1_000 * USDG, dave, dave);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        _weekend(0);
        _open();
        usdg.mint(dave, 100 * USDG);
        vm.startPrank(dave);
        usdg.approve(address(accounts), type(uint256).max);
        accounts.repay(SPY, type(uint256).max, dave);
        vm.stopPrank();
        vm.prank(owner);
        accounts.setPremiumRate(0);
        backstop.claim();
        vm.prank(bob);
        backstop.startCooldown();
        skip(7 days);
        uint256 shares = backstop.balanceOf(bob);
        vm.prank(bob);
        backstop.redeem(shares, bob, bob);
        assertEq(backstop.totalSupply(), 0);
        assertGt(backstop.held(), 0);
        assertEq(_thinRound().taken, 0);
    }

    function test_ARoundWithNoBidsClearsAtItsFloorAndTheBackstopBuysThere() public {
        GapBackstop backstop = _backstop(5_000 * USDG);
        _back(backstop, bob, 4_000 * USDG);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _open();
        vm.expectRevert(abi.encodeWithSelector(ReopeningAuction.NotClearingPrice.selector, 180e8));
        auction.clear(NVDA, openMs, 180e8);
        auction.clear(NVDA, openMs, 171e8);
        ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
        assertEq(r.sold, 0);
        assertGt(r.taken, 0);
        assertEq(backstop.held(), 4_000 * USDG - ((r.taken * 171e8 - 1) / 1e20 + 1));
    }

    function _backstop(uint256 limit) internal returns (GapBackstop backstop) {
        uint256[] memory limits = new uint256[](2);
        (limits[0], limits[1]) = (limit, limit);
        backstop = new GapBackstop(accounts, owner, limit, limits);
        vm.prank(owner);
        accounts.setBackstop(address(backstop));
    }

    function _thinRound() internal returns (ReopeningAuction.Round memory) {
        _quote(NVDA, 200e8, 202e8);
        _position(alice, 10 * SHARE, 1_600 * USDG);
        _weekend(190e8);
        auction.enroll(alice, CROSS, NVDA);
        _commit(carol, 1 * SHARE, 185e8, "c");
        _toReveal();
        _reveal(carol, 1 * SHARE, 185e8, "c");
        _open();
        auction.clear(NVDA, openMs, 185e8);
        return auction.round(NVDA, openMs);
    }

    function _back(GapBackstop backstop, address depositor, uint256 assets) internal {
        usdg.mint(depositor, assets);
        vm.startPrank(depositor);
        usdg.approve(address(backstop), assets);
        backstop.deposit(assets, depositor);
        vm.stopPrank();
    }
}
