// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract LiquidatorTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubToken internal weth;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    StubAggregator internal ethUsd;
    BandDouble internal band;
    MarginDouble internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    address internal owner = makeAddr("owner");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal buyer = makeAddr("buyer");

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        weth = new StubToken(18);
        nvda = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(NVDA, 3, 200e8, 202e8);
        _quote(SPY, 3, 600e8, 606e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (NVDA, SPY);
        engine = new MarginDouble(symbols, address(band), address(ethUsd));
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
            5_00,
            10_00
        );
        liquidator = new Liquidator(accounts);
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        vm.stopPrank();
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), 1_000_000 * USDG);
        vault.deposit(1_000_000 * USDG, owner);
        vm.stopPrank();
        usdg.mint(buyer, 1_000_000 * USDG);
        vm.prank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
    }

    function test_TheLiquidatorTakesItsReadsFromTheAccounts() public view {
        assertEq(address(liquidator.accounts()), address(accounts));
        assertEq(address(liquidator.band()), address(band));
        assertEq(address(liquidator.engine()), address(engine));
        assertEq(address(liquidator.usdg()), address(usdg));
        assertEq(address(liquidator.weth()), address(weth));
        assertEq(address(liquidator.ethUsd()), address(ethUsd));
    }

    function test_AHealthyPositionIsNotLiquidated() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        (int256 equity, uint256 requirement, bool short, bool closed) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 400e18);
        assertEq(requirement, 400e18);
        assertFalse(short);
        assertFalse(closed);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 0);
    }

    function test_AnOpenMarketAuctionFallsFromTheHighEdgeABasisPointASecond() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 190e8, 192e8);
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 300e18);
        assertEq(requirement, 380e18);
        assertTrue(short);
        vm.expectEmit(address(liquidator));
        emit Liquidator.AuctionStarted(alice, CROSS, false);
        liquidator.start(alice, CROSS);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 192e8);
        skip(10 minutes);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 180_48000000);
        skip(40 minutes);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 171e8);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotForSale.selector, address(usdg)));
        liquidator.price(alice, CROSS, address(usdg));
    }

    function test_ABuyerPaysTheAuctionPriceAndItRepaysThePositionLessTheFee() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 190e8, 192e8);
        liquidator.start(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.CostAboveLimit.selector, 384 * USDG, 383 * USDG));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 2 * SHARE, 383 * USDG, buyer);
        vm.expectEmit(address(liquidator));
        emit Liquidator.Bought(alice, CROSS, address(nvda), 2 * SHARE, 384 * USDG, 1_920_000, buyer, bob);
        vm.prank(buyer);
        (uint256 bought, uint256 cost) = liquidator.buy(alice, CROSS, address(nvda), 2 * SHARE, 384 * USDG, bob);
        assertEq(bought, 2 * SHARE);
        assertEq(cost, 384 * USDG);
        assertEq(nvda.balanceOf(bob), 2 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 8 * SHARE);
        assertEq(accounts.debt(alice, CROSS), 1_600 * USDG - 382_080_000);
        assertEq(accounts.reserve(), 1_920_000);
        assertEq(usdg.balanceOf(address(liquidator)), 0);
        assertEq(usdg.balanceOf(buyer), 1_000_000 * USDG - 384 * USDG);
        vm.prank(buyer);
        (bought, cost) = liquidator.buy(alice, CROSS, address(nvda), 1, 1, bob);
        assertEq(bought, 1);
        assertEq(cost, 1);
    }

    function test_APurchaseStopsAtWhatRepaysThePositionInFull() public {
        _deposit(alice, CROSS, address(nvda), 20 * SHARE);
        _borrow(alice, CROSS, 150 * USDG);
        _quote(NVDA, 3, 9e8, 10e8);
        liquidator.start(alice, CROSS);
        vm.prank(buyer);
        (uint256 bought, uint256 cost) = liquidator.buy(alice, CROSS, address(nvda), 20 * SHARE, 1_000 * USDG, buyer);
        assertEq(cost, 150_753_769);
        assertEq(bought, 15_075_376_900_000_000_000);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.reserve(), 753_768);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 20 * SHARE - bought);
        assertEq(usdg.balanceOf(address(liquidator)), 0);
        assertEq(usdg.balanceOf(buyer), 1_000_000 * USDG - cost + 1);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 1_000 * USDG, buyer);
    }

    function test_WhileClosedAShortfallTheBandExplainsWaits() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        engine.set(20_00, 0, 1);
        _quote(NVDA, 2, 180e8, 200e8);
        (int256 equity,, bool short, bool closed) = liquidator.shortfall(alice, CROSS);
        assertTrue(closed);
        assertFalse(short);
        assertEq(equity, 400e18);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        _quote(NVDA, 2, 180e8, 195e8);
        vm.expectEmit(address(liquidator));
        emit Liquidator.AuctionStarted(alice, CROSS, true);
        liquidator.start(alice, CROSS);
        skip(20 minutes);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 189_15000000);
        vm.prank(buyer);
        (uint256 bought,) = liquidator.buy(alice, CROSS, address(nvda), 0.5e18, 1_000 * USDG, buyer);
        assertEq(bought, 0.5e18);
        assertEq(liquidator.hourlyAllowance(alice, CROSS, address(nvda)), 0.45e18);
        vm.prank(buyer);
        (bought,) = liquidator.buy(alice, CROSS, address(nvda), 5 * SHARE, 1_000 * USDG, buyer);
        assertEq(bought, 0.45e18);
        assertEq(liquidator.hourlyAllowance(alice, CROSS, address(nvda)), 0);
        vm.expectRevert(Liquidator.NothingToBuy.selector);
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 1_000 * USDG, buyer);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        assertEq(liquidator.hourlyAllowance(bob, CROSS, address(nvda)), 1 * SHARE);
        skip(30 minutes);
        assertEq(liquidator.hourlyAllowance(alice, CROSS, address(nvda)), 0.4075e18);
        skip(30 minutes);
        vm.prank(buyer);
        (bought,) = liquidator.buy(alice, CROSS, address(nvda), 5 * SHARE, 1_000 * USDG, buyer);
        assertEq(bought, 0.86e18);
    }

    function test_WhileClosedAPositionFallsShortOnlyAtBothEdges() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        engine.set(20_00, 0, 1);
        engine.setPool(200e8);
        _quote(NVDA, 2, 200e8, 210e8);
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 400e18);
        assertEq(requirement, 400e18);
        assertFalse(short);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        _quote(NVDA, 2, 195e8, 210e8);
        (equity, requirement, short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 500e18);
        assertEq(requirement, 520e18);
        assertTrue(short);
        _quote(NVDA, 2, 199e8, 230e8);
        (equity, requirement, short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 390e18);
        assertEq(requirement, 398e18);
        assertTrue(short);
        liquidator.start(alice, CROSS);
    }

    function test_TheMarketIsClosedFromTheCloseBeforeAWeekendToTheOpenAfterIt() public {
        uint64 nowMs = uint64(vm.getBlockTimestamp()) * 1000;
        assertFalse(_closedIn(BandDouble.Session(2, 1, 3, 0, nowMs + 1)));
        assertFalse(_closedIn(BandDouble.Session(2, 2, 1, 0, 0)));
        assertTrue(_closedIn(BandDouble.Session(2, 3, 1, 0, 0)));
        assertTrue(_closedIn(BandDouble.Session(2, 2, 3, 0, 0)));
        assertTrue(_closedIn(BandDouble.Session(2, 1, 3, 0, nowMs)));
        assertTrue(_closedIn(BandDouble.Session(1, 3, 1, 0, 0)));
        assertTrue(_closedIn(BandDouble.Session(0, 0, 0, 0, 0)));
    }

    function test_APurchaseRepaysAtMostHalfOfALargeDebt() public {
        _deposit(alice, CROSS, address(nvda), 20 * SHARE);
        _borrow(alice, CROSS, 3_000 * USDG);
        _quote(NVDA, 3, 180e8, 182e8);
        liquidator.start(alice, CROSS);
        vm.prank(buyer);
        (uint256 bought, uint256 cost) = liquidator.buy(alice, CROSS, address(nvda), 20 * SHARE, 10_000 * USDG, buyer);
        assertEq(bought, 8_283_174_115_384_615_384);
        assertEq(cost, 1_507_537_689);
        assertEq(accounts.debt(alice, CROSS), 3_000 * USDG - (cost - cost * 50 / 10_000));
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 10_000 * USDG, buyer);
    }

    function test_OnceTheAccountsHaveABackstopOnlyItWritesOff() public {
        address backstop = makeAddr("backstop");
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 50e8, 51e8);
        liquidator.start(alice, CROSS);
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE, 1_000 * USDG, buyer);
        vm.prank(owner);
        accounts.setBackstop(backstop);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotBackstop.selector, bob));
        vm.prank(bob);
        liquidator.writeOff(alice, CROSS);
        vm.prank(backstop);
        assertEq(liquidator.writeOff(alice, CROSS), 1_090 * USDG);
    }

    function test_AnUnknownSessionIsJudgedAgainstTheBufferedOpenRequirement() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        engine.set(40_00, 0, 0);
        (int256 equity, uint256 requirement, bool short, bool closed) = liquidator.shortfall(alice, CROSS);
        assertTrue(closed);
        assertEq(equity, 400e18);
        assertEq(requirement, 400e18);
        assertFalse(short);
        (, uint256 health,,) = accounts.health(alice, CROSS);
        assertEq(health, 800e18);
        _quote(NVDA, 1, 190e8, 199e8);
        (,, short,) = liquidator.shortfall(alice, CROSS);
        assertTrue(short);
        liquidator.start(alice, CROSS);
        (, bool auctionClosed) = liquidator.auctions(alice, CROSS);
        assertTrue(auctionClosed);
    }

    function test_WethIsJudgedAtEightyFourPercent() public {
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        ethUsd.setRound(1_905e8, vm.getBlockTimestamp());
        (int256 equity,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 0.2e18);
        assertFalse(short);
        ethUsd.setRound(1_900e8, vm.getBlockTimestamp());
        (equity,, short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, -4e18);
        assertTrue(short);
        liquidator.start(alice, CROSS);
        assertEq(liquidator.price(alice, CROSS, address(weth)), 1_900e8);
        vm.prank(buyer);
        (, uint256 cost) = liquidator.buy(alice, CROSS, address(weth), 0.1e18, 1_000 * USDG, buyer);
        assertEq(cost, 190 * USDG);
        assertEq(accounts.reserve(), 950_000);
        skip(50 minutes);
        ethUsd.setRound(1_900e8, block.timestamp);
        assertEq(liquidator.price(alice, CROSS, address(weth)), 1_710e8);
        ethUsd.setRound(1_900e8, block.timestamp - 1 days - 61);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotForSale.selector, address(weth)));
        liquidator.price(alice, CROSS, address(weth));
    }

    function test_NothingIsJudgedWhileTheSequencerIsUnsettledOrAHeldAssetCannotBePriced() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 150e8, 151e8);
        band.setSequencer(address(1), false);
        (,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertFalse(short);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.CannotJudge.selector, alice, CROSS));
        liquidator.stop(alice, CROSS);
        band.setSequencer(address(1), true);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(2, uint64(block.timestamp - 10 minutes), 1e18, 2e18));
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(1, uint64(block.timestamp + 10 minutes), 1e18, 2e18));
        liquidator.start(alice, CROSS);
        _quote(NVDA, 0, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 1_000 * USDG, buyer);
    }

    function test_AnIsolatedPositionIsLiquidatedAloneAndTheCrossIsUntouched() public {
        _deposit(alice, NVDA, address(nvda), 10 * SHARE);
        _borrow(alice, NVDA, 1_600 * USDG);
        _deposit(alice, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _quote(NVDA, 3, 150e8, 151e8);
        (int256 crossEquity,,,) = liquidator.shortfall(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        liquidator.start(alice, NVDA);
        vm.expectEmit(address(liquidator));
        emit Liquidator.Bought(alice, NVDA, address(nvda), 10 * SHARE, 1_510 * USDG, 0, buyer, buyer);
        vm.prank(buyer);
        liquidator.buy(alice, NVDA, address(nvda), 10 * SHARE, 10_000 * USDG, buyer);
        assertEq(accounts.collateral(alice, NVDA, address(nvda)), 0);
        assertEq(accounts.reserve(), 0);
        vm.expectRevert(Liquidator.NothingToBuy.selector);
        vm.prank(buyer);
        liquidator.buy(alice, NVDA, address(spy), 1, 10_000 * USDG, buyer);
        uint256 lent = vault.debt();
        uint256 written = liquidator.writeOff(alice, NVDA);
        assertEq(written, 90 * USDG);
        assertEq(vault.debt(), lent - written);
        assertEq(accounts.debt(alice, NVDA), 0);
        (int256 after_,,,) = liquidator.shortfall(alice, CROSS);
        assertEq(after_, crossEquity);
        assertEq(accounts.collateral(alice, CROSS, address(spy)), 10 * SHARE);
        assertEq(accounts.debt(alice, CROSS), 1_000 * USDG);
    }

    function test_CashSettlesAPositionWhoseStockTokenTheIssuerFroze() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _deposit(alice, CROSS, address(usdg), 500 * USDG);
        _quote(NVDA, 3, 100e8, 101e8);
        nvda.pause();
        liquidator.start(alice, CROSS);
        vm.expectRevert(StubStockToken.IsPaused.selector);
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1 * SHARE, 1_000 * USDG, buyer);
        vm.expectEmit(address(liquidator));
        emit Liquidator.CashSettled(alice, CROSS, 500 * USDG);
        assertEq(liquidator.settleCash(alice, CROSS), 500 * USDG);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 0);
        assertEq(accounts.debt(alice, CROSS), 1_100 * USDG);
        assertEq(accounts.reserve(), 0);
        vm.expectRevert(Liquidator.NothingToBuy.selector);
        liquidator.settleCash(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, bob, CROSS));
        liquidator.settleCash(bob, CROSS);
    }

    function test_AnAuctionRunsForItsLifetimeInOneMarketStateAndStopsWhenThePositionRecovers() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 190e8, 192e8);
        liquidator.start(alice, CROSS);
        (uint64 startedAt,) = liquidator.auctions(alice, CROSS);
        skip(30 minutes);
        liquidator.start(alice, CROSS);
        (uint64 again,) = liquidator.auctions(alice, CROSS);
        assertEq(again, startedAt);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.StillShort.selector, alice, CROSS));
        liquidator.stop(alice, CROSS);
        skip(30 minutes);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NoAuction.selector, alice, CROSS));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 1_000 * USDG, buyer);
        liquidator.start(alice, CROSS);
        (again,) = liquidator.auctions(alice, CROSS);
        assertEq(again, block.timestamp);
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        _quote(NVDA, 2, 150e8, 160e8);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NoAuction.selector, alice, CROSS));
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 1, 1_000 * USDG, buyer);
        liquidator.start(alice, CROSS);
        skip(4 hours - 1);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 135e8);
        skip(1);
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 0);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        _quote(NVDA, 3, 200e8, 202e8);
        vm.expectEmit(address(liquidator));
        emit Liquidator.AuctionStopped(alice, CROSS);
        liquidator.stop(alice, CROSS);
        (again,) = liquidator.auctions(alice, CROSS);
        assertEq(again, 0);
    }

    function test_ABurnIsCountedBeforeAPositionIsJudged() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        nvda.adminBurn(address(accounts), 1 * SHARE);
        (,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertFalse(short);
        liquidator.start(alice, CROSS);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 9 * SHARE);
    }

    function test_AWrittenOffPositionForgivesItsPremiumAndClearsItsAuction() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 50e8, 51e8);
        liquidator.start(alice, CROSS);
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE, 1_000 * USDG, buyer);
        uint256 debt = accounts.debt(alice, CROSS);
        uint256 shares = accounts.debtShares(alice, CROSS);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.WriteOff(alice, CROSS, debt, shares, 0);
        assertEq(liquidator.writeOff(alice, CROSS), debt);
        (uint64 startedAt,) = liquidator.auctions(alice, CROSS);
        assertEq(startedAt, 0);
        _deposit(bob, CROSS, address(nvda), 1 * SHARE);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(address(0xdead));
        liquidator.writeOff(bob, CROSS);
    }

    function test_AnAuctionBottomsOutAtItsFloorAndSellsNoHaltedAsset() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 3_000 * USDG);
        _quote(NVDA, 3, 100e8, 101e8);
        liquidator.start(alice, CROSS);
        skip(50 minutes);
        ethUsd.setRound(2_000e8, vm.getBlockTimestamp());
        assertEq(liquidator.price(alice, CROSS, address(nvda)), 90e8);
        assertEq(liquidator.price(alice, CROSS, address(weth)), 1_800e8);
        _quote(SPY, 0, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotForSale.selector, address(spy)));
        liquidator.price(alice, CROSS, address(spy));
    }

    function test_APositionHoldingWethIsNotJudgedWithoutAWorkingFeed() public {
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 2_000 * USDG);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _borrow(bob, CROSS, 1_600 * USDG);
        engine.set(150_00, 0, 2);
        vm.mockCallRevert(address(ethUsd), abi.encodeWithSelector(StubAggregator.latestRoundData.selector), "");
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 0);
        assertEq(requirement, 0);
        assertFalse(short);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.CannotJudge.selector, alice, CROSS));
        liquidator.stop(alice, CROSS);
        (,, short,) = liquidator.shortfall(bob, CROSS);
        assertTrue(short);
        vm.clearMockedCalls();
        skip(1 days + 61);
        (,, short,) = liquidator.shortfall(alice, CROSS);
        assertFalse(short);
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (NVDA, SPY);
        MarginDouble unpriced = new MarginDouble(symbols, address(band), address(0));
        SupplyVault other = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](2);
        MarginAccounts testnet = new MarginAccounts(
            IBand(address(band)), IMargin(address(unpriced)), other, weth, owner, caps, 0, 0, 5_00, 10_00
        );
        Liquidator bare = new Liquidator(testnet);
        assertEq(address(bare.ethUsd()), address(0));
        vm.startPrank(owner);
        testnet.setLiquidator(address(bare));
        vm.stopPrank();
        weth.mint(alice, 1 * SHARE);
        vm.startPrank(alice);
        weth.approve(address(testnet), 1 * SHARE);
        testnet.deposit(CROSS, address(weth), 1 * SHARE, alice);
        vm.stopPrank();
        (equity,, short,) = bare.shortfall(alice, CROSS);
        assertEq(equity, 0);
        assertFalse(short);
    }

    function test_CashSettlesNoMoreThanThePositionOwes() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _deposit(alice, CROSS, address(usdg), 1_700 * USDG);
        engine.set(150_00, 0, 2);
        assertEq(liquidator.settleCash(alice, CROSS), 1_600 * USDG);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 100 * USDG);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(usdg.balanceOf(address(liquidator)), 0);
    }

    function test_ThePositionsCashCountsForItsEquity() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _deposit(alice, CROSS, address(usdg), 500 * USDG);
        _quote(NVDA, 3, 190e8, 192e8);
        (int256 equity,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 800e18);
        assertFalse(short);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
    }

    function test_APositionShortOnItsPremiumAloneIsLiquidated() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        uint256 premium = accounts.premium(alice, CROSS);
        assertGt(premium, 0);
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 400e18 - int256(premium * 1e12));
        assertEq(requirement, 400e18);
        assertTrue(short);
        liquidator.start(alice, CROSS);
    }

    function test_OnlyTheLiquidatorCollectsFees() public {
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotLiquidator.selector, alice));
        vm.prank(alice);
        accounts.collectFee(1);
    }

    function test_GasOfEachCall() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(spy), 1 * SHARE);
        _borrow(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 3, 100e8, 101e8);
        skip(1 hours);
        liquidator.start(alice, CROSS);
        vm.snapshotGasLastCall("start");
        skip(1 minutes);
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), 2 * SHARE, 1_000 * USDG, buyer);
        vm.snapshotGasLastCall("buy");
    }

    function _closedIn(BandDouble.Session memory session) internal returns (bool closed) {
        band.setSession(session);
        (,,, closed) = liquidator.shortfall(alice, CROSS);
    }

    function _quote(bytes32 symbol, uint8 state, uint64 low, uint128 high) internal {
        uint64 mid = uint64((uint256(low) + high) / 2);
        band.setQuote(symbol, BandDouble.Quote(state, 3, mid, 50, low, high));
    }

    function _deposit(address account, bytes32 position, address token, uint256 amount) internal {
        if (token == address(usdg)) usdg.mint(account, amount);
        else StubToken(token).mint(account, amount);
        vm.startPrank(account);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, account);
        vm.stopPrank();
    }

    function _borrow(address account, bytes32 position, uint256 assets) internal {
        vm.prank(account);
        accounts.borrow(position, assets, account, account);
    }
}
