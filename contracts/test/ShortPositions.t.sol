// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {ShortPositions} from "../src/ShortPositions.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {IV3SwapRouter} from "../src/interfaces/IUniswapV3.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubSwapRouter} from "./devnode/StubSwapRouter.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract ShortPositionsTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant TSLA = "TSLA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    BandDouble internal band;
    MarginDouble internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    StockLendingVault internal lending;
    StockLendingVault internal spyLending;
    StubSwapRouter internal router;
    ShortPositions internal shorts;
    address internal owner = makeAddr("owner");
    address internal guardian = makeAddr("guardian");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal dave = makeAddr("dave");
    address internal keeper = makeAddr("keeper");

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        nvda = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = (NVDA, TSLA, SPY);
        engine = new MarginDouble(
            symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD"))
        );
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (type(uint256).max, type(uint256).max, type(uint256).max);
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            new StubToken(18),
            owner,
            caps,
            type(uint256).max,
            type(uint256).max,
            0,
            10_00
        );
        lending = _lendingFor(nvda);
        spyLending = _lendingFor(spy);
        router = new StubSwapRouter(address(usdg));
        router.setPrice(address(nvda), 200e8);
        router.setPrice(address(spy), 600e8);
        shorts = new ShortPositions(accounts, IV3SwapRouter(address(router)), _fees(500, 0, 500));
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setGuardian(guardian);
        accounts.setLending(NVDA, lending);
        accounts.setLending(SPY, spyLending);
        lending.setBorrower(address(shorts));
        spyLending.setBorrower(address(shorts));
        vm.stopPrank();
        _quote(NVDA, 200e8);
        _quote(SPY, 600e8);
        _lender(alice, nvda, 100 * SHARE);
    }

    function test_AShortBorrowsSellsAndKeepsTheProceedsWithItsMargin() public {
        uint256 supplied = vault.totalAssets();
        _margin(bob, NVDA, 1_000 * USDG);
        vm.prank(bob);
        uint256 proceeds = shorts.sell(NVDA, 10 * SHARE, 0, bob);
        assertEq(proceeds, 2_000 * USDG);
        (int256 held, uint256 debt, uint256 shares,) = shorts.position(bob, NVDA);
        assertEq(held, int256(3_000 * USDG));
        assertEq(debt, 10 * SHARE);
        assertEq(lending.debt(), 10 * SHARE);
        assertEq(nvda.balanceOf(address(shorts)), 0);
        (uint256 total, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(total, shares);
        assertEq(booked, 3_000 * USDG);
        assertEq(usdg.balanceOf(address(shorts)), 3_000 * USDG);
        assertEq(vault.totalAssets(), supplied);
        assertEq(vault.debt(), 0);
    }

    function test_AShortIsMarginedByTheEngineAsANegativeQuantityAtTheHighEdge() public {
        _margin(bob, NVDA, 1_000 * USDG);
        vm.prank(bob);
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        int256[] memory quantities = new int256[](3);
        uint256[] memory prices = new uint256[](3);
        (quantities[0], prices[0]) = (-10e18, 202e8);
        (uint256 expected,,) = engine.currentRequirement(quantities, prices);
        (int256 equity, uint256 requirement, uint8 missing, uint8 regime) = shorts.health(bob, NVDA);
        assertEq(requirement, expected);
        assertEq(requirement, 404e18);
        assertEq(equity, 3_000e18 - 2_020e18);
        assertEq(missing, 0);
        assertEq(regime, 2);
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.InsufficientMargin.selector, int256(940e18), 1_212e18));
        shorts.sell(NVDA, 20 * SHARE, 0, bob);
    }

    function test_AShortSellsAtNoLessThanTheBandsLowEdgeAndItsOwnLimit() public {
        _margin(bob, NVDA, 1_000 * USDG);
        router.setPrice(address(nvda), 197e8);
        vm.prank(bob);
        vm.expectRevert("Too little received");
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        router.setPrice(address(nvda), 199e8);
        vm.prank(bob);
        vm.expectRevert("Too little received");
        shorts.sell(NVDA, 10 * SHARE, 1_995 * USDG, bob);
        vm.prank(bob);
        assertEq(shorts.sell(NVDA, 10 * SHARE, 1_990 * USDG, bob), 1_990 * USDG);
    }

    function test_AfterATenPercentFallAShortSellsAtNoLessThanTheBandsCentreUntilTheNextDayEnds() public {
        _margin(bob, NVDA, 2_000 * USDG);
        shorts.mark(NVDA);
        skip(1 days);
        _quote(NVDA, 195e8);
        shorts.mark(NVDA);
        _quote(NVDA, 181e8);
        (bool restricted, uint64 close) = shorts.restriction(NVDA);
        assertFalse(restricted);
        assertEq(close, 200e8);
        _quote(NVDA, 180e8);
        (restricted,) = shorts.restriction(NVDA);
        assertTrue(restricted);
        router.setPrice(address(nvda), 179e8);
        uint32 today = uint32(vm.getBlockTimestamp() / 1 days);
        vm.expectEmit(address(shorts));
        emit ShortPositions.Restricted(NVDA, today + 1);
        shorts.mark(NVDA);
        vm.prank(bob);
        vm.expectRevert("Too little received");
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        _quote(NVDA, 200e8);
        router.setPrice(address(nvda), 199e8);
        vm.warp((uint256(today) + 2) * 1 days - 1);
        (restricted,) = shorts.restriction(NVDA);
        assertTrue(restricted);
        vm.prank(bob);
        vm.expectRevert("Too little received");
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        vm.warp((uint256(today) + 2) * 1 days);
        (restricted, close) = shorts.restriction(NVDA);
        assertFalse(restricted);
        assertEq(close, 180e8);
        vm.prank(bob);
        assertEq(shorts.sell(NVDA, 10 * SHARE, 0, bob), 1_990 * USDG);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        (restricted, close) = shorts.restriction(NVDA);
        assertFalse(restricted);
        assertEq(close, 180e8);
    }

    function test_ACoverBuysBackAndRepaysAndAFullCoverClosesTheShort() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        router.setPrice(address(nvda), 190e8);
        vm.prank(bob);
        assertEq(shorts.cover(NVDA, 4 * SHARE, type(uint256).max, bob), 760 * USDG);
        (int256 held, uint256 debt,,) = shorts.position(bob, NVDA);
        assertEq(held, int256(2_240 * USDG));
        assertEq(debt, 6 * SHARE);
        assertEq(lending.debt(), 6 * SHARE);
        vm.prank(bob);
        assertEq(shorts.cover(NVDA, type(uint256).max, type(uint256).max, bob), 1_140 * USDG);
        (held, debt,,) = shorts.position(bob, NVDA);
        assertEq(held, int256(1_100 * USDG));
        assertEq(debt, 0);
        assertEq(lending.debt(), 0);
        (uint256 total, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(total, 0);
        assertEq(booked, 0);
        vm.prank(guardian);
        accounts.setBorrowingPaused(true);
        vm.prank(bob);
        shorts.withdraw(NVDA, 1_100 * USDG, bob, bob);
        assertEq(usdg.balanceOf(bob), 1_100 * USDG);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.NoShort.selector);
        shorts.cover(NVDA, SHARE, type(uint256).max, bob);
    }

    function test_ACoverPaysNoMoreThanItsLimitOrTheShortsUsdg() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        router.setPrice(address(nvda), 210e8);
        vm.prank(bob);
        vm.expectRevert("Too much requested");
        shorts.cover(NVDA, 10 * SHARE, 2_099 * USDG, bob);
        router.setPrice(address(nvda), 301e8);
        vm.prank(bob);
        vm.expectRevert("Too much requested");
        shorts.cover(NVDA, 10 * SHARE, type(uint256).max, bob);
    }

    function test_TheFeeAccruesInTheTokenAndAFullCoverRepaysIt() public {
        _short(bob, 3_000 * USDG, 50 * SHARE);
        vm.warp(block.timestamp + 365 days);
        (, uint256 debt,,) = shorts.position(bob, NVDA);
        assertEq(debt, lending.debt());
        assertGt(debt, 50 * SHARE);
        vm.prank(bob);
        shorts.cover(NVDA, type(uint256).max, type(uint256).max, bob);
        assertEq(lending.debt(), 0);
        assertGt(accounts.lent(alice, CROSS, address(nvda)), 100 * SHARE);
    }

    function test_AWithdrawalFromAShortPassesItsChecksAndWaitsOnTheGuardian() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(bob);
        vm.expectRevert(
            abi.encodeWithSelector(ShortPositions.InsufficientMargin.selector, int256(403e18 - 1e12), 404e18)
        );
        shorts.withdraw(NVDA, 577 * USDG + 1, bob, bob);
        vm.prank(bob);
        shorts.withdraw(NVDA, 576 * USDG, bob, bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortfall.selector, int256(404e18), 404e18));
        shorts.liquidate(bob, NVDA);
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.InsufficientCollateral.selector, 3_000 * USDG));
        shorts.withdraw(NVDA, 3_000 * USDG, bob, bob);
        vm.prank(guardian);
        accounts.setBorrowingPaused(true);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.BorrowingIsPaused.selector);
        shorts.withdraw(NVDA, 1, bob, bob);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.BorrowingIsPaused.selector);
        shorts.sell(NVDA, SHARE, 0, bob);
    }

    function test_OnlyTheAccountOrAnAddressItAuthorizedInTheAccountsActsForIt() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        usdg.mint(carol, 100 * USDG);
        vm.startPrank(carol);
        usdg.approve(address(shorts), type(uint256).max);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Unauthorized.selector, carol, bob));
        shorts.deposit(NVDA, 100 * USDG, bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Unauthorized.selector, carol, bob));
        shorts.withdraw(NVDA, 1, bob, carol);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Unauthorized.selector, carol, bob));
        shorts.sell(NVDA, SHARE, 0, bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Unauthorized.selector, carol, bob));
        shorts.cover(NVDA, SHARE, type(uint256).max, bob);
        vm.stopPrank();
        vm.prank(bob);
        accounts.setAuthorization(carol, true);
        vm.startPrank(carol);
        shorts.deposit(NVDA, 100 * USDG, bob);
        (, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(booked, 3_100 * USDG);
        shorts.sell(NVDA, SHARE, 0, bob);
        shorts.cover(NVDA, SHARE, type(uint256).max, bob);
        shorts.withdraw(NVDA, 100 * USDG, bob, carol);
        vm.stopPrank();
        assertEq(usdg.balanceOf(carol), 100 * USDG);
    }

    function test_NoShortForAnAccountOrCallerTheIssuerBlocklisted() public {
        _margin(bob, NVDA, 1_000 * USDG);
        vm.prank(bob);
        accounts.setAuthorization(carol, true);
        nvda.blockAccount(bob, true);
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Blocked.selector, bob));
        shorts.sell(NVDA, SHARE, 0, bob);
        nvda.blockAccount(bob, false);
        nvda.blockAccount(carol, true);
        vm.prank(carol);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.Blocked.selector, carol));
        shorts.sell(NVDA, SHARE, 0, bob);
    }

    function test_NoShortWhileTheBandCannotVouchForTheAsset() public {
        _margin(bob, NVDA, 1_000 * USDG);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        shorts.sell(NVDA, SHARE, 0, bob);
        _quote(NVDA, 200e8);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(1, uint64(block.timestamp + 600), 1e18, 2e18));
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.CorporateActionPending.selector, NVDA));
        shorts.sell(NVDA, SHARE, 0, bob);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(0, 0, 0, 0));
        band.setSequencer(address(1), false);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.SequencerNotSettled.selector);
        shorts.sell(NVDA, SHARE, 0, bob);
        band.setSequencer(address(0), true);
        engine.set(20_00, 1, 2);
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.LiquidityUnknown.selector, NVDA));
        shorts.sell(NVDA, SHARE, 0, bob);
        engine.set(20_00, 0, 0);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.SessionUnknown.selector);
        shorts.sell(NVDA, SHARE, 0, bob);
    }

    function test_TheWeekendLeverageCapHoldsFromTheLastOpenBeforeAClosure() public {
        _margin(bob, NVDA, 300 * USDG);
        engine.set(1_00, 0, 2);
        vm.prank(bob);
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        engine.set(1_00, 0, 1);
        vm.prank(bob);
        vm.expectRevert(
            abi.encodeWithSelector(ShortPositions.WeekendLeverageExceeded.selector, 4_040e18, int256(260e18))
        );
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        engine.set(1_00, 0, 2);
        band.setSession(BandDouble.Session(2, 1, 2, 0, uint64(block.timestamp + 1 days) * 1000));
        vm.prank(bob);
        vm.expectRevert(
            abi.encodeWithSelector(ShortPositions.WeekendLeverageExceeded.selector, 4_040e18, int256(260e18))
        );
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
    }

    function test_EveryShortPaysItsPartOfABuyInInProportion() public {
        _short(bob, 3_000 * USDG, 60 * SHARE);
        _short(carol, 1_000 * USDG, 20 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 40 * SHARE, alice);
        assertEq(lending.requested(), 20 * SHARE);
        vm.warp(block.timestamp + 1 days);
        router.setPrice(address(nvda), 201e8);
        _quote(NVDA, 201e8);
        uint256 debt = lending.debt();
        (int256 bobBefore, uint256 bobDebt,,) = shorts.position(bob, NVDA);
        (int256 carolBefore, uint256 carolDebt,,) = shorts.position(carol, NVDA);
        vm.expectEmit(address(shorts));
        emit ShortPositions.BuyIn(NVDA, 20 * SHARE, 4_020 * USDG);
        lending.buyIn(0);
        assertEq(lending.debt(), debt - 20 * SHARE);
        (int256 bobAfter, uint256 bobLeft,,) = shorts.position(bob, NVDA);
        (int256 carolAfter, uint256 carolLeft,,) = shorts.position(carol, NVDA);
        assertApproxEqAbs(bobBefore - bobAfter, int256(3_015 * USDG), 1);
        assertApproxEqAbs(carolBefore - carolAfter, int256(1_005 * USDG), 1);
        assertApproxEqAbs(bobDebt - bobLeft, 15 * SHARE, 1e6);
        assertApproxEqAbs(carolDebt - carolLeft, 5 * SHARE, 1e6);
        (, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(booked, 3_000 * USDG + 12_000 * USDG + 1_000 * USDG + 4_000 * USDG - 4_020 * USDG);
        assertGe(booked, uint256(bobAfter + carolAfter));
        vm.prank(alice);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 40 * SHARE);
    }

    function test_ABuyInPaysNoMoreThanTheBandsHighEdgeNorTheTokensShortsUsdg() public {
        _short(bob, 5_000 * USDG, 90 * SHARE);
        _short(carol, 50_000 * USDG, 0);
        _spyShort(carol, 60_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 50 * SHARE, alice);
        vm.warp(block.timestamp + 1 days);
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 200e8, 50, 199e8, 201e8));
        router.setPrice(address(nvda), 2015e7);
        vm.expectRevert("Too much requested");
        lending.buyIn(0);
        _quote(NVDA, 600e8);
        router.setPrice(address(nvda), 600e8);
        (, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(booked, 23_000 * USDG);
        assertGt(usdg.balanceOf(address(shorts)), 24_000 * USDG);
        vm.expectRevert("Too much requested");
        lending.buyIn(0);
        router.setPrice(address(nvda), 202e8);
        lending.buyIn(0);
    }

    function test_OnlyTheAssetsLendingVaultBuysIn() public {
        vm.prank(alice);
        vm.expectRevert();
        shorts.buyIn(SHARE);
        StockLendingVault other = _lendingFor(nvda);
        vm.prank(address(other));
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotLending.selector, address(other)));
        shorts.buyIn(SHARE);
    }

    function test_ARepaymentGoesToTheRecallQueueBeforeTheShortsBorrowAgain() public {
        _short(bob, 5_000 * USDG, 90 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 30 * SHARE, alice);
        vm.prank(bob);
        shorts.cover(NVDA, 30 * SHARE, type(uint256).max, bob);
        assertEq(lending.assigned(), 20 * SHARE);
        assertEq(lending.idle(), 10 * SHARE);
        assertEq(lending.borrowable(), 10 * SHARE);
        vm.prank(bob);
        vm.expectRevert(
            abi.encodeWithSelector(StockLendingVault.InsufficientLiquidity.selector, 11 * SHARE, 10 * SHARE)
        );
        shorts.sell(NVDA, 11 * SHARE, 0, bob);
    }

    function test_AShortThatFallsShortIsBoughtBackAndTheCallerTakesABonus() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortfall.selector, int256(980e18), 404e18));
        shorts.liquidate(bob, NVDA);
        _quote(NVDA, 250e8);
        router.setPrice(address(nvda), 250e8);
        vm.expectEmit(address(shorts));
        emit ShortPositions.Liquidate(keeper, bob, NVDA, 10 * SHARE, 2_500 * USDG, 12_500_000, 0);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(usdg.balanceOf(keeper), 12_500_000);
        assertEq(lending.debt(), 0);
        (int256 held, uint256 debt, uint256 shares,) = shorts.position(bob, NVDA);
        assertEq(held, int256(487_500_000));
        assertEq(debt, 0);
        assertEq(shares, 0);
        vm.prank(bob);
        shorts.withdraw(NVDA, 487_500_000, bob, bob);
        assertEq(usdg.balanceOf(address(shorts)), 0);
    }

    function test_WhileTheMarketIsClosedAShortFallsShortOnlyAtBothEdges() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        engine.set(20_00, 0, 1);
        band.setQuote(NVDA, BandDouble.Quote(2, 3, 230e8, 13_00, 200e8, 260e8));
        router.setPrice(address(nvda), 265e8);
        (int256 equity, uint256 requirement,,) = shorts.health(bob, NVDA);
        assertLt(equity, int256(requirement));
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortfall.selector, int256(1_000e18), 400e18));
        shorts.liquidate(bob, NVDA);
        band.setQuote(NVDA, BandDouble.Quote(2, 3, 270e8, 4_00, 260e8, 280e8));
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(lending.debt(), 0);
    }

    function test_ALiquidationWaitsWhileThePoolAsksMoreThanTheBandsHighEdge() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _quote(NVDA, 250e8);
        router.setPrice(address(nvda), 253e8);
        vm.expectRevert(ShortPositions.AboveLimit.selector);
        shorts.liquidate(bob, NVDA);
    }

    function test_AShortThatCannotBuyBackAllItOwesLeavesTheRestToTheLenders() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        uint256 lentBefore = accounts.lent(alice, CROSS, address(nvda));
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        (int256 held, uint256 debt,,) = shorts.position(bob, NVDA);
        assertEq(held, 0);
        assertEq(debt, 0);
        assertEq(lending.debt(), 0);
        assertEq(usdg.balanceOf(keeper), 0);
        assertApproxEqAbs(lentBefore - accounts.lent(alice, CROSS, address(nvda)), 25 * SHARE / 10, 1);
        assertEq(nvda.balanceOf(address(lending)), 975 * SHARE / 10);
    }

    function test_AShortInDeficitAfterABuyInLeavesItToTheLendersNotTheOtherShorts() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _short(carol, 20_000 * USDG, 30 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        vm.warp(block.timestamp + 1 days);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        lending.buyIn(0);
        (int256 bobHeld, uint256 bobDebt,,) = shorts.position(bob, NVDA);
        assertLt(bobHeld, 0);
        (int256 carolHeld, uint256 carolDebt,,) = shorts.position(carol, NVDA);
        int256 carolWorth = carolHeld - int256(carolDebt * 404e8 / 1e20);
        uint256 assetsBefore = lending.totalAssets();
        vm.prank(bob);
        vm.expectRevert(ShortPositions.Insolvent.selector);
        shorts.cover(NVDA, SHARE, type(uint256).max, bob);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.Insolvent.selector);
        shorts.deposit(NVDA, 1, bob);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        (carolHeld, carolDebt,,) = shorts.position(carol, NVDA);
        assertApproxEqAbs(carolHeld - int256(carolDebt * 404e8 / 1e20), carolWorth, 2);
        uint256 lost = uint256(-bobHeld) * 1e20 / 404e8 + bobDebt;
        assertApproxEqAbs(assetsBefore - lending.totalAssets(), lost, 2);
        (, uint256 booked,) = shorts.book(NVDA, 0);
        assertGe(booked, uint256(carolHeld));
        uint256 held = accounts.collateral(alice, CROSS, address(nvda));
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), held + 35 * SHARE);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 0);
        vm.prank(carol);
        shorts.cover(NVDA, type(uint256).max, type(uint256).max, carol);
        assertEq(lending.debt(), 0);
        (carolHeld,,,) = shorts.position(carol, NVDA);
        vm.prank(carol);
        shorts.withdraw(NVDA, uint256(carolHeld), carol, carol);
        assertLe(usdg.balanceOf(address(shorts)), 1);
    }

    function test_ALaterShortOwesAtLeastWhatItBorrowed() public {
        _short(bob, 5_000 * USDG, 60 * SHARE);
        skip(97 days);
        _short(carol, 1_000 * USDG, 7 * SHARE);
        (, uint256 debt,,) = shorts.position(carol, NVDA);
        assertGe(debt, 7 * SHARE);
    }

    function test_TheLastShortToCloseRepaysAllTheVaultIsOwed() public {
        _short(bob, 100_000 * USDG, 90 * SHARE);
        skip(5 * 365 days);
        assertGt(lending.debt(), 180 * SHARE);
        vm.prank(bob);
        shorts.cover(NVDA, type(uint256).max, type(uint256).max, bob);
        assertEq(lending.debt(), 0);
        assertEq(lending.scaledDebt(), 0);
    }

    function test_TheBonusIsAtMostWhatTheShortHasLeft() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _quote(NVDA, 2999e7);
        router.setPrice(address(nvda), 2999e7);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(usdg.balanceOf(keeper), USDG);
        (int256 held,,,) = shorts.position(bob, NVDA);
        assertEq(held, 0);
    }

    function test_ADefaultBuysWhatItCanOnlyWithinTheBand() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 410e8);
        vm.expectRevert("Too little received");
        shorts.liquidate(bob, NVDA);
    }

    function test_ADeficitLeftOnceTheVaultIsOwedNothingFallsOnTheOtherShorts() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _short(carol, 20_000 * USDG, 30 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE, alice);
        skip(1 days);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        lending.buyIn(0);
        assertEq(lending.debt(), 0);
        (int256 bobHeld,,,) = shorts.position(bob, NVDA);
        assertLt(bobHeld, 0);
        (int256 carolHeld,,,) = shorts.position(carol, NVDA);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        (int256 carolAfter,,,) = shorts.position(carol, NVDA);
        assertApproxEqAbs(carolAfter, carolHeld + bobHeld, 1);
        (, uint256 booked,) = shorts.book(NVDA, 0);
        assertGe(booked, uint256(carolAfter));
        assertLe(booked, uint256(carolAfter) + 1);
    }

    function test_APartialCoverPaysNoMoreThanTheBandsHighEdge() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        router.setPrice(address(nvda), 203e8);
        vm.prank(bob);
        vm.expectRevert("Too much requested");
        shorts.cover(NVDA, 4 * SHARE, type(uint256).max, bob);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        shorts.cover(NVDA, 4 * SHARE, type(uint256).max, bob);
        vm.prank(bob);
        assertEq(shorts.cover(NVDA, 10 * SHARE, type(uint256).max, bob), 2_030 * USDG);
        assertEq(lending.debt(), 0);
    }

    function test_ALiquidationAndABuyInPayAtMostOnePercentAboveTheBandsCentre() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        skip(1 days);
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 200e8, 5_00, 190e8, 210e8));
        router.setPrice(address(nvda), 203e8);
        vm.expectRevert("Too much requested");
        lending.buyIn(0);
        router.setPrice(address(nvda), 202e8);
        lending.buyIn(0);
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 350e8, 5_00, 3325e7, 3675e7));
        router.setPrice(address(nvda), 354e8);
        vm.expectRevert(ShortPositions.AboveLimit.selector);
        shorts.liquidate(bob, NVDA);
        router.setPrice(address(nvda), 353e8);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(lending.debt(), 0);
    }

    function test_ABuyInThatWouldLeaveARemnantBuysInAllTheVaultIsOwed() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE - 1, alice);
        assertEq(lending.requested(), 10 * SHARE - 1);
        skip(1 days);
        vm.expectEmit(address(shorts));
        emit ShortPositions.BuyIn(NVDA, 10 * SHARE, 2_000 * USDG);
        lending.buyIn(0);
        assertEq(lending.debt(), 0);
    }

    function test_OnceTheVaultIsOwedNothingTheNextSaleOpensANewBook() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE, alice);
        skip(1 days);
        lending.buyIn(0);
        assertEq(lending.debt(), 0);
        _lender(carol, nvda, 50 * SHARE);
        _margin(dave, NVDA, 1_000 * USDG);
        vm.expectEmit(address(shorts));
        emit ShortPositions.NewEpoch(NVDA, 1);
        vm.prank(dave);
        shorts.sell(NVDA, 5 * SHARE, 0, dave);
        assertEq(shorts.epoch(NVDA), 1);
        (, uint256 daveDebt, uint256 daveShares, uint256 daveEpoch) = shorts.position(dave, NVDA);
        assertEq(daveDebt, 5 * SHARE);
        assertEq(daveShares, 5e24);
        assertEq(daveEpoch, 1);
        (int256 held, uint256 debt, uint256 shares, uint256 bobEpoch) = shorts.position(bob, NVDA);
        assertEq(held, int256(1_000 * USDG));
        assertEq(debt, 0);
        assertEq(shares, 10e24);
        assertEq(bobEpoch, 0);
        vm.prank(bob);
        vm.expectRevert(ShortPositions.NoShort.selector);
        shorts.cover(NVDA, SHARE, type(uint256).max, bob);
        vm.prank(bob);
        shorts.withdraw(NVDA, 1_000 * USDG, bob, bob);
        (uint256 oldShares, uint256 oldUsdg,) = shorts.book(NVDA, 0);
        assertEq(oldShares, 0);
        assertEq(oldUsdg, 0);
        (,, shares, bobEpoch) = shorts.position(bob, NVDA);
        assertEq(shares, 0);
        assertEq(bobEpoch, 1);
    }

    function test_AShortOfAnEarlierBookInDeficitIsChargedToThatBookAlone() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _short(carol, 20_000 * USDG, 30 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE, alice);
        skip(1 days);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        lending.buyIn(0);
        (int256 bobHeld,,,) = shorts.position(bob, NVDA);
        assertLt(bobHeld, 0);
        _lender(alice, nvda, 50 * SHARE);
        _margin(dave, NVDA, 5_000 * USDG);
        vm.prank(dave);
        shorts.sell(NVDA, 5 * SHARE, 0, dave);
        (int256 carolHeld,,,) = shorts.position(carol, NVDA);
        vm.prank(carol);
        vm.expectRevert(ShortPositions.DeficitOpen.selector);
        shorts.withdraw(NVDA, uint256(carolHeld), carol, carol);
        uint256 debt = lending.debt();
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(lending.debt(), debt);
        (int256 carolAfter,,,) = shorts.position(carol, NVDA);
        assertApproxEqAbs(carolAfter, carolHeld + bobHeld, 1);
        vm.prank(carol);
        shorts.withdraw(NVDA, uint256(carolAfter), carol, carol);
    }

    function test_ADeficitBeyondWhatTheOtherShortsOweFallsOnThem() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _short(carol, 20_000 * USDG, 30 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 99 * SHARE, alice);
        skip(1 days);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        lending.buyIn(0);
        (int256 bobHeld,,,) = shorts.position(bob, NVDA);
        (int256 carolHeld, uint256 carolDebt,,) = shorts.position(carol, NVDA);
        assertLt(bobHeld, 0);
        assertLt(carolDebt * 400e8 / 1e20, uint256(-bobHeld));
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        assertEq(lending.debt(), 0);
        (int256 carolAfter, uint256 carolLeft,,) = shorts.position(carol, NVDA);
        assertEq(carolLeft, 0);
        assertApproxEqAbs(carolAfter, carolHeld + bobHeld, 1);
    }

    function test_WhileTheIssuerBlocklistsTheShortsNothingIsBoughtInAndNoticesRun() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        nvda.blockAccount(address(shorts), true);
        _margin(carol, NVDA, 1_000 * USDG);
        vm.prank(carol);
        vm.expectRevert(abi.encodeWithSelector(StubStockToken.Blocked.selector, address(shorts)));
        shorts.sell(NVDA, SHARE, 0, carol);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        skip(1 days);
        vm.expectRevert(abi.encodeWithSelector(StubStockToken.Blocked.selector, address(shorts)));
        lending.buyIn(0);
        nvda.blockAccount(address(shorts), false);
        lending.buyIn(0);
        assertEq(lending.claimable(0), 5 * SHARE);
    }

    function test_NoBuyInWhileTheSequencerIsNotSettled() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        skip(1 days);
        band.setSequencer(address(1), false);
        vm.expectRevert(ShortPositions.SequencerNotSettled.selector);
        lending.buyIn(0);
    }

    function test_AWriteOffNeverLeavesTheVaultOwedTooLittleForTheBooksShares() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        _short(carol, 20_000 * USDG, 30 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        skip(1 days);
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        lending.buyIn(0);
        (int256 bobHeld,,,) = shorts.position(bob, NVDA);
        (, uint256 carolDebt, uint256 carolShares,) = shorts.position(carol, NVDA);
        uint64 price = uint64(uint256(-bobHeld) * 1e20 / carolDebt) + 6_000;
        band.setQuote(NVDA, BandDouble.Quote(3, 3, price, 1, price, price));
        assertGt(carolDebt - uint256(-bobHeld) * 1e20 / price, 0);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        uint256 keep = (carolShares - 1) / 1e12 + 1;
        assertEq(lending.debt(), keep);
        (, uint256 carolLeft,,) = shorts.position(carol, NVDA);
        assertEq(carolLeft, keep);
    }

    function test_ADeficitOfAWeiNeverForgivesTheBook() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 50_000 * USDG, 90 * SHARE);
        (uint256 shares,,) = shorts.book(NVDA, 0);
        uint256 target = shares / 1e12 + 270;
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE - target, alice);
        skip(1 days);
        lending.buyIn(0);
        assertEq(lending.debt(), target);
        _lender(dave, nvda, 700 * SHARE);
        _lender(alice, nvda, 400 * SHARE);
        _margin(bob, NVDA, 300_000 * USDG);
        vm.prank(bob);
        shorts.sell(NVDA, 800 * SHARE, 0, bob);
        _margin(carol, NVDA, 2);
        vm.prank(carol);
        shorts.sell(NVDA, 1e10, 0, carol);
        for (uint256 k; k < 10; ++k) {
            (int256 h,,,) = shorts.position(carol, NVDA);
            if (h == 0) break;
            vm.prank(carol);
            shorts.cover(NVDA, 1, type(uint256).max, carol);
        }
        uint256 free = lending.idle() + lending.claimable(0);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), free - 4e9, alice);
        address eve = makeAddr("eve");
        _lender(eve, nvda, 1e9);
        vm.prank(bob);
        shorts.sell(NVDA, 5e9, 0, bob);
        vm.prank(eve);
        accounts.recall(CROSS, address(nvda), 1e9, eve);
        skip(1 days);
        lending.buyIn(lending.tickets() - 1);
        (int256 carolHeld,,,) = shorts.position(carol, NVDA);
        assertLt(carolHeld, 0);
        (, uint256 bobDebt,,) = shorts.position(bob, NVDA);
        uint256 assets = lending.totalAssets();
        shorts.liquidate(carol, NVDA);
        (, uint256 bobAfter,,) = shorts.position(bob, NVDA);
        assertApproxEqAbs(bobAfter, bobDebt, 1e10);
        assertLe(assets - lending.totalAssets(), 2e10);
    }

    function test_ABuyBackOfDustThatWouldRetireAllTheSharesClosesTheShort() public {
        _short(bob, 1_000 * USDG, 1e15);
        vm.prank(bob);
        shorts.cover(NVDA, 1e15 - 2, type(uint256).max, bob);
        skip(3_650 days);
        for (uint256 k; k < 40; ++k) {
            (,, uint256 shares,) = shorts.position(bob, NVDA);
            if (shares == 0) break;
            vm.prank(bob);
            shorts.cover(NVDA, 1, type(uint256).max, bob);
        }
        (, uint256 debt, uint256 left,) = shorts.position(bob, NVDA);
        assertEq(left, 0);
        assertEq(debt, 0);
        assertEq(lending.debt(), 0);
    }

    function test_AFullCoverOfOneOfSeveralShortsRetiresAllItsShares() public {
        _short(bob, 3_000 * USDG, 30 * SHARE);
        _short(carol, 1_000 * USDG, 7 * SHARE);
        skip(97 days);
        vm.prank(bob);
        shorts.cover(NVDA, type(uint256).max, type(uint256).max, bob);
        (, uint256 debt, uint256 shares,) = shorts.position(bob, NVDA);
        assertEq(shares, 0);
        assertEq(debt, 0);
        (, uint256 carolDebt, uint256 carolShares,) = shorts.position(carol, NVDA);
        (uint256 total,,) = shorts.book(NVDA, 0);
        assertEq(total, carolShares);
        assertEq(carolDebt, lending.debt());
    }

    function test_AShortThatOwesNothingAfterABuyInClosesWithACover() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 100 * SHARE, alice);
        skip(1 days);
        lending.buyIn(0);
        vm.prank(bob);
        assertEq(shorts.cover(NVDA, SHARE, type(uint256).max, bob), 0);
        (int256 held,, uint256 shares,) = shorts.position(bob, NVDA);
        assertEq(shares, 0);
        (uint256 total, uint256 booked,) = shorts.book(NVDA, 0);
        assertEq(total, 0);
        assertEq(booked, 0);
        vm.prank(guardian);
        accounts.setBorrowingPaused(true);
        vm.prank(bob);
        shorts.withdraw(NVDA, uint256(held), bob, bob);
    }

    function test_TheConstructorTakesOnlyPoolsTheRoutersFactoryHas() public {
        vm.expectRevert(ShortPositions.LengthMismatch.selector);
        this.deploy(new uint24[](2));
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.InvalidPool.selector, address(0), uint24(3000)));
        this.deploy(_fees(500, 3000, 0));
        router.setPrice(address(spy), 0);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.InvalidPool.selector, address(spy), uint24(500)));
        this.deploy(_fees(500, 0, 500));
        ShortPositions nvdaOnly = this.deploy(_fees(500, 0, 0));
        assertEq(nvdaOnly.fee(NVDA), 500);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortable.selector, SPY));
        nvdaOnly.fee(SPY);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortable.selector, TSLA));
        shorts.fee(TSLA);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortable.selector, bytes32("AAPL")));
        shorts.fee("AAPL");
        assertEq(address(shorts.accounts()), address(accounts));
        assertEq(address(shorts.band()), address(band));
        assertEq(address(shorts.engine()), address(engine));
        assertEq(address(shorts.usdg()), address(usdg));
        assertEq(address(shorts.router()), address(router));
        assertEq(shorts.weekendLeverage(), 50_000);
    }

    function test_AnAssetWithoutALendingVaultCannotBeShorted() public {
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = (NVDA, TSLA, SPY);
        uint256[] memory caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (type(uint256).max, type(uint256).max, type(uint256).max);
        MarginAccounts bare = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0)),
            new StubToken(18),
            owner,
            caps,
            type(uint256).max,
            type(uint256).max,
            0,
            10_00
        );
        ShortPositions orphan = new ShortPositions(bare, IV3SwapRouter(address(router)), _fees(500, 0, 0));
        usdg.mint(bob, 1_000 * USDG);
        vm.startPrank(bob);
        usdg.approve(address(orphan), 1_000 * USDG);
        orphan.deposit(NVDA, 1_000 * USDG, bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.NotShortable.selector, NVDA));
        orphan.sell(NVDA, SHARE, 0, bob);
        vm.stopPrank();
    }

    function test_EveryEntryRefusesNothingAndWaitsWhileTheBandIsHalted() public {
        _short(bob, 1_000 * USDG, 10 * SHARE);
        vm.startPrank(bob);
        vm.expectRevert(ShortPositions.ZeroAmount.selector);
        shorts.deposit(NVDA, 0, bob);
        vm.expectRevert(ShortPositions.ZeroAmount.selector);
        shorts.withdraw(NVDA, 0, bob, bob);
        vm.expectRevert(ShortPositions.ZeroAmount.selector);
        shorts.sell(NVDA, 0, 0, bob);
        vm.expectRevert(ShortPositions.ZeroAmount.selector);
        shorts.cover(NVDA, 0, 0, bob);
        vm.stopPrank();
        vm.expectRevert(ShortPositions.NoShort.selector);
        shorts.liquidate(carol, NVDA);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 95 * SHARE, alice);
        skip(1 days);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        shorts.mark(NVDA);
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        shorts.withdraw(NVDA, USDG, bob, bob);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        shorts.liquidate(bob, NVDA);
        vm.expectRevert(abi.encodeWithSelector(ShortPositions.AssetHalted.selector, NVDA));
        lending.buyIn(0);
        _quote(NVDA, 300e8);
        band.setSequencer(address(1), false);
        vm.expectRevert(ShortPositions.SequencerNotSettled.selector);
        shorts.liquidate(bob, NVDA);
    }

    function test_TheShortsGas() public {
        _margin(bob, NVDA, 1_000 * USDG);
        _margin(carol, NVDA, 3_000 * USDG);
        vm.prank(carol);
        shorts.sell(NVDA, 10 * SHARE, 0, carol);
        vm.prank(bob);
        shorts.sell(NVDA, 10 * SHARE, 0, bob);
        vm.snapshotGasLastCall("sell");
        vm.prank(bob);
        shorts.cover(NVDA, 4 * SHARE, type(uint256).max, bob);
        vm.snapshotGasLastCall("cover");
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 90 * SHARE, alice);
        vm.warp(block.timestamp + 1 days);
        lending.buyIn(0);
        vm.snapshotGasLastCall("buyIn");
        _quote(NVDA, 400e8);
        router.setPrice(address(nvda), 400e8);
        vm.prank(keeper);
        shorts.liquidate(bob, NVDA);
        vm.snapshotGasLastCall("liquidate");
    }

    function deploy(uint24[] memory fees) external returns (ShortPositions) {
        return new ShortPositions(accounts, IV3SwapRouter(address(router)), fees);
    }

    function _lendingFor(StubStockToken token) internal returns (StockLendingVault fresh) {
        fresh =
            new StockLendingVault(IERC20(address(token)), owner, SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.prank(owner);
        fresh.setDepositor(address(accounts));
    }

    function _fees(uint24 a, uint24 b, uint24 c) internal pure returns (uint24[] memory fees) {
        fees = new uint24[](3);
        (fees[0], fees[1], fees[2]) = (a, b, c);
    }

    function _quote(bytes32 symbol, uint64 mid) internal {
        band.setQuote(symbol, BandDouble.Quote(3, 3, mid, 1_00, mid - mid / 100, mid + mid / 100));
    }

    function _lender(address account, StubStockToken token, uint256 amount) internal {
        token.mint(account, amount);
        vm.startPrank(account);
        token.approve(address(accounts), amount);
        accounts.deposit(CROSS, address(token), amount, account);
        accounts.lend(CROSS, address(token), amount, account);
        vm.stopPrank();
    }

    function _margin(address account, bytes32 symbol, uint256 amount) internal {
        usdg.mint(account, amount);
        vm.startPrank(account);
        usdg.approve(address(shorts), amount);
        shorts.deposit(symbol, amount, account);
        vm.stopPrank();
    }

    function _short(address account, uint256 margin, uint256 amount) internal {
        _margin(account, NVDA, margin);
        if (amount == 0) return;
        vm.prank(account);
        shorts.sell(NVDA, amount, 0, account);
    }

    function _spyShort(address account, uint256 margin, uint256 amount) internal {
        _lender(alice, spy, 100 * SHARE);
        _margin(account, SPY, margin);
        vm.prank(account);
        shorts.sell(SPY, amount, 0, account);
    }
}
