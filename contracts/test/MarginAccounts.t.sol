// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
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

contract MarginAccountsTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant TSLA = "TSLA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);
    uint256 internal constant STEP = uint256(200e8) / 2 ** 14;
    uint256 internal constant YEAR_MS = 365 days * 1000;
    string internal constant TIMELINE = "../stylus/contracts/margin/testdata/session-timeline.json";

    StubUsdg internal usdg;
    StubToken internal weth;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    StubAggregator internal ethUsd;
    BandDouble internal band;
    MarginDouble internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    address internal owner = makeAddr("owner");
    address internal lender = makeAddr("lender");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");

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
        _quote(NVDA, 3, 200e8);
        _quote(TSLA, 3, 300e8);
        _quote(SPY, 3, 600e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = (NVDA, TSLA, SPY);
        engine = new MarginDouble(symbols, address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            weth,
            owner,
            _uncapped(3),
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        vm.prank(owner);
        vault.setBorrower(address(accounts));
        usdg.mint(lender, 100_000 * USDG);
        vm.startPrank(lender);
        usdg.approve(address(vault), 100_000 * USDG);
        vault.deposit(100_000 * USDG, lender);
        vm.stopPrank();
    }

    function test_TheAccountsTakeTheEnginesAssetsAndTheirTokens() public {
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        assertEq(symbols.length, 3);
        assertEq(symbols[0], NVDA);
        assertEq(symbols[1], TSLA);
        assertEq(symbols[2], SPY);
        assertEq(tokens[0], address(nvda));
        assertEq(tokens[1], address(0));
        assertEq(tokens[2], address(spy));
        assertEq(address(accounts.usdg()), address(usdg));
        assertEq(address(accounts.ethUsd()), address(ethUsd));
        assertEq(accounts.weekendLeverage(), 50_000);
        MarginDouble other = new MarginDouble(symbols, address(1), address(ethUsd));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BandMismatch.selector, address(1)));
        new MarginAccounts(
            IBand(address(band)),
            IMargin(address(other)),
            vault,
            weth,
            owner,
            _uncapped(3),
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
    }

    function test_TheCapsAreOnePerAsset() public {
        vm.expectRevert(MarginAccounts.LengthMismatch.selector);
        new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), vault, weth, owner, _uncapped(2), 0, 0, 5_00, 10_00
        );
    }

    function test_AnAssetSitsInTheCrossPositionOrItsOwnIsolatedOne() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        nvda.mint(alice, 1);
        vm.startPrank(alice);
        nvda.approve(address(accounts), 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetInOtherPosition.selector, NVDA, CROSS));
        accounts.deposit(NVDA, address(nvda), 1, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(nvda), SPY));
        accounts.deposit(SPY, address(nvda), 1, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownPosition.selector, TSLA));
        accounts.deposit(TSLA, address(usdg), 1, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(ethUsd), CROSS));
        accounts.deposit(CROSS, address(ethUsd), 1, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.deposit(CROSS, address(nvda), 0, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(0), CROSS));
        accounts.deposit(CROSS, address(0), 1, alice);
        vm.stopPrank();
        _deposit(alice, SPY, address(spy), 2 * SHARE);
        _deposit(alice, SPY, address(usdg), 100 * USDG);
        assertEq(accounts.collateral(alice, SPY, address(spy)), 2 * SHARE);
        assertEq(accounts.collateral(alice, SPY, address(usdg)), 100 * USDG);
        _withdraw(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, NVDA, address(nvda), 1 * SHARE);
        assertEq(accounts.collateral(alice, NVDA, address(nvda)), 1 * SHARE);
    }

    function test_OnlyAnAuthorizedAddressActsForAnAccount() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        nvda.mint(bob, SHARE);
        usdg.mint(bob, 100 * USDG);
        vm.startPrank(bob);
        nvda.approve(address(accounts), SHARE);
        usdg.approve(address(accounts), 100 * USDG);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        accounts.deposit(NVDA, address(nvda), SHARE, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        accounts.borrow(CROSS, 1, alice, bob);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        accounts.withdraw(CROSS, address(nvda), 1, alice, bob);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        accounts.deposit(SPY, address(usdg), 100 * USDG, alice);
        vm.stopPrank();
        vm.expectEmit(address(accounts));
        emit MarginAccounts.AuthorizationSet(alice, bob, true);
        vm.prank(alice);
        accounts.setAuthorization(bob, true);
        assertTrue(accounts.isAuthorized(alice, bob));
        vm.startPrank(bob);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Deposit(bob, alice, CROSS, address(nvda), SHARE);
        accounts.deposit(CROSS, address(nvda), SHARE, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Borrow(bob, alice, CROSS, 100 * USDG, 100 * USDG * 1e6, bob);
        accounts.borrow(CROSS, 100 * USDG, alice, bob);
        accounts.withdraw(CROSS, address(nvda), SHARE, alice, bob);
        accounts.deposit(SPY, address(usdg), 100 * USDG, alice);
        vm.stopPrank();
        assertEq(usdg.balanceOf(bob), 100 * USDG);
        assertEq(nvda.balanceOf(bob), SHARE);
        assertEq(accounts.debt(alice, CROSS), 100 * USDG);
        assertEq(accounts.collateral(alice, SPY, address(usdg)), 100 * USDG);
        vm.prank(alice);
        accounts.setAuthorization(bob, false);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.borrow(CROSS, 1, alice, bob);
    }

    function test_APositionBorrowsUpToItsRequirement() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        (int256 equity, uint256 requirement,, uint8 regime) = accounts.health(alice, CROSS);
        assertEq(equity, 2_000e18);
        assertEq(requirement, 400e18);
        assertEq(regime, 2);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Borrow(alice, alice, CROSS, 1_600 * USDG, 1_600 * USDG * 1e6, alice);
        _borrow(alice, CROSS, 1_600 * USDG);
        assertEq(usdg.balanceOf(alice), 1_600 * USDG);
        assertEq(accounts.debt(alice, CROSS), 1_600 * USDG);
        assertEq(vault.debt(), 1_600 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(alice, CROSS, 1);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        _borrow(alice, CROSS, 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownPosition.selector, TSLA));
        _borrow(alice, TSLA, 1);
    }

    function test_AnyoneRepaysAnAccountAndNoMoreThanItOwes() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        usdg.mint(bob, 2_000 * USDG);
        vm.startPrank(bob);
        usdg.approve(address(accounts), 2_000 * USDG);
        assertEq(accounts.repay(CROSS, 400 * USDG, alice), 400 * USDG);
        assertEq(accounts.debt(alice, CROSS), 600 * USDG);
        assertEq(accounts.repay(CROSS, 2_000 * USDG, alice), 600 * USDG);
        vm.stopPrank();
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.debtShares(alice, CROSS), 0);
        assertEq(accounts.totalDebtShares(), 0);
        assertEq(vault.debt(), 0);
        assertEq(usdg.balanceOf(bob), 1_000 * USDG);
        _deposit(bob, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, CROSS, 700 * USDG);
        _borrow(bob, CROSS, 1_300 * USDG);
        skip(365 days);
        uint256 owed = accounts.debt(alice, CROSS);
        _repay(alice, CROSS, owed);
        assertEq(accounts.debtShares(alice, CROSS), 0);
        assertEq(accounts.debt(alice, CROSS), 0);
    }

    function test_DebtsGrowWithTheVaultsInterestProRata() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _borrow(bob, CROSS, 3_000 * USDG);
        skip(365 days);
        uint256 total = vault.debt();
        assertGt(total, 4_000 * USDG);
        uint256 a = accounts.debt(alice, CROSS);
        uint256 b = accounts.debt(bob, CROSS);
        assertApproxEqAbs(a + b, total, 2);
        assertGe(a + b, total);
        assertApproxEqAbs(a * 3, b, 3);
        usdg.mint(lender, total);
        vm.startPrank(lender);
        usdg.approve(address(vault), total);
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.NotBorrower.selector, lender));
        vault.repay(total);
        vm.stopPrank();
    }

    function test_AWithdrawalFromAPositionInDebtPassesItsChecks() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(weth)), 1 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(ethUsd)), 0);
        _borrow(alice, CROSS, 1_600 * USDG);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(usdg), 1));
        _withdraw(alice, CROSS, address(usdg), 1);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 11 * SHARE)
        );
        _withdraw(alice, CROSS, address(nvda), 11 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(weth), 2 * SHARE)
        );
        _withdraw(alice, CROSS, address(weth), 2 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(ethUsd), CROSS));
        _withdraw(alice, CROSS, address(ethUsd), 1);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        _withdraw(alice, CROSS, address(weth), 0);
        _withdraw(alice, CROSS, address(weth), 1 * SHARE);
        assertEq(weth.balanceOf(alice), 1 * SHARE);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _withdraw(alice, CROSS, address(nvda), 1);
        _repay(alice, CROSS, 1_600 * USDG);
        _quote(NVDA, 0, 0);
        _withdraw(alice, CROSS, address(nvda), 5 * SHARE);
        _withdraw(alice, CROSS, address(nvda), 5 * SHARE);
        assertEq(nvda.balanceOf(alice), 10 * SHARE);
    }

    function test_ACorporateActionClosesBorrowingAndWithdrawalsButNotRepaymentsOrDeposits() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(1, uint64(block.timestamp + 10 minutes), 1e18, 2e18));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.CorporateActionPending.selector, NVDA));
        _borrow(alice, CROSS, 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.CorporateActionPending.selector, NVDA));
        _withdraw(alice, CROSS, address(weth), 1);
        _deposit(alice, CROSS, address(nvda), 1 * SHARE);
        _repay(alice, CROSS, 100 * USDG);
        assertEq(accounts.debt(alice, CROSS), 900 * USDG);
        _deposit(bob, CROSS, address(spy), 1 * SHARE);
        _borrow(bob, CROSS, 100 * USDG);
    }

    function test_AHaltedAssetIsRepayOnlyAndCountsForNothing() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _quote(NVDA, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetHalted.selector, NVDA));
        _borrow(alice, CROSS, 1);
        _deposit(alice, CROSS, address(usdg), 500 * USDG);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertEq(equity, -500e18);
        assertEq(requirement, 0);
        vm.expectRevert(MarginAccounts.BorrowingAgainstUsdg.selector);
        _borrow(alice, CROSS, 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetHalted.selector, NVDA));
        _withdraw(alice, CROSS, address(usdg), 1);
        _repay(alice, CROSS, 1_000 * USDG);
        assertEq(accounts.debt(alice, CROSS), 0);
        _deposit(bob, CROSS, address(nvda), 1 * SHARE);
        _deposit(bob, CROSS, address(weth), 1 * SHARE);
        _withdraw(bob, CROSS, address(nvda), 1 * SHARE);
        _borrow(bob, CROSS, 100 * USDG);
    }

    function test_APositionRepaysWithItsUsdgWhileAnAssetIsHalted() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(alice, CROSS, address(usdg), 500 * USDG);
        _quote(NVDA, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.repayWithCollateral(CROSS, 500 * USDG, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Repay(alice, alice, CROSS, 500 * USDG, 500 * USDG * 1e6);
        vm.prank(alice);
        assertEq(accounts.repayWithCollateral(CROSS, 2_000 * USDG, alice), 500 * USDG);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 0);
        assertEq(accounts.debt(alice, CROSS), 500 * USDG);
        assertEq(vault.debt(), 500 * USDG);
        assertEq(usdg.balanceOf(address(accounts)), 0);
        _deposit(alice, CROSS, address(usdg), 800 * USDG);
        vm.prank(alice);
        assertEq(accounts.repayWithCollateral(CROSS, 800 * USDG, alice), 500 * USDG);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 300 * USDG);
        assertEq(accounts.debtShares(alice, CROSS), 0);
    }

    function test_ADegradedAssetBacksDebtAtItsBandsLowEdge() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _quote(NVDA, 1, 150e8);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertEq(equity, 1_500e18);
        assertEq(requirement, 300e18);
        _borrow(alice, CROSS, 1_200 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(alice, CROSS, 1);
    }

    function test_NoRiskIsAddedWhileTheSessionLiquidityOrSequencerIsUnknown() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        engine.set(20_00, 0, 0);
        vm.expectRevert(MarginAccounts.SessionUnknown.selector);
        _borrow(alice, CROSS, 1);
        _deposit(bob, SPY, address(weth), 1 * SHARE);
        _borrow(bob, SPY, 1_600 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(bob, SPY, 1);
        engine.set(20_00, 1, 2);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.LiquidityUnknown.selector, NVDA));
        _borrow(alice, CROSS, 1);
        _deposit(bob, CROSS, address(spy), 1 * SHARE);
        _borrow(bob, CROSS, 100 * USDG);
        engine.set(20_00, 4, 2);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.LiquidityUnknown.selector, SPY));
        _borrow(bob, CROSS, 1);
        engine.set(20_00, 0, 2);
        band.setSequencer(address(1), false);
        vm.expectRevert(MarginAccounts.SequencerNotSettled.selector);
        _borrow(alice, CROSS, 1);
        band.setSequencer(address(1), true);
        _borrow(alice, CROSS, 1);
    }

    function test_TheWeekendCapHoldsThroughTheRampAndTheClosure() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        engine.set(5_00, 0, 2);
        _borrow(alice, CROSS, 1_700 * USDG);
        assertEq(accounts.leverage(alice, CROSS), 66_666);
        band.setSession(BandDouble.Session(2, 1, 3, 0, uint64(block.timestamp + 10 hours) * 1000));
        vm.expectPartialRevert(MarginAccounts.WeekendLeverageExceeded.selector);
        _borrow(alice, CROSS, 1);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        _borrow(alice, CROSS, 1);
        engine.set(5_00, 0, 3);
        vm.expectPartialRevert(MarginAccounts.WeekendLeverageExceeded.selector);
        _borrow(alice, CROSS, 1);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _borrow(bob, CROSS, 1_600 * USDG);
        engine.set(5_00, 0, 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.WeekendLeverageExceeded.selector, 2_000e18, 399e18));
        _borrow(bob, CROSS, 1 * USDG);
        assertEq(accounts.leverage(bob, CROSS), 50_000);
    }

    function test_AnIsolatedPositionBearsItsOwnLoss() public {
        _deposit(alice, NVDA, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, NVDA, 1_500 * USDG);
        _borrow(alice, CROSS, 1_000 * USDG);
        (int256 crossEquity, uint256 crossRequirement,,) = accounts.health(alice, CROSS);
        assertEq(crossEquity, 5_000e18);
        assertEq(crossRequirement, 1_200e18);
        _quote(NVDA, 3, 100e8);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, NVDA);
        assertEq(equity, -500e18);
        assertEq(requirement, 200e18);
        (int256 after_, uint256 afterRequirement,,) = accounts.health(alice, CROSS);
        assertEq(after_, crossEquity);
        assertEq(afterRequirement, crossRequirement);
        _borrow(alice, CROSS, 100 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(alice, NVDA, 1);
    }

    function test_TheLiquidationPriceIsWhereEquityMeetsTheRequirement() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.liquidationPrice(alice, CROSS, NVDA, 0), 0);
        uint256 preview = accounts.liquidationPrice(alice, CROSS, NVDA, 1_000 * USDG);
        _borrow(alice, CROSS, 1_000 * USDG);
        uint256 price = accounts.liquidationPrice(alice, CROSS, NVDA, 0);
        assertEq(price, preview);
        assertGe(price, 125e8);
        assertLe(price, 125e8 + STEP + 1);
        assertEq(accounts.liquidationPrice(alice, CROSS, SPY, 0), 0);
        _quote(NVDA, 3, 120e8);
        vm.expectCall(address(engine), abi.encodeWithSelector(MarginDouble.currentRequirement.selector), 1);
        assertEq(accounts.liquidationPrice(alice, CROSS, NVDA, 0), 120e8);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, bytes32("AMD")));
        accounts.liquidationPrice(alice, CROSS, "AMD", 0);
    }

    function test_WethCountsForEightyPercentOfAFreshChainlinkPrice() public {
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 1_600e18);
        _borrow(alice, CROSS, 1_600 * USDG);
        ethUsd.setRound(2_000e8, block.timestamp - 86_461);
        (equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, -1_600e18);
        ethUsd.setRound(0, block.timestamp);
        (equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, -1_600e18);
        vm.mockCallRevert(address(ethUsd), abi.encodeWithSignature("latestRoundData()"), "");
        (equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, -1_600e18);
        (bytes32[] memory symbols,) = accounts.stocks();
        MarginDouble unpriced = new MarginDouble(symbols, address(band), address(0));
        MarginAccounts testnet = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(unpriced)),
            vault,
            weth,
            owner,
            _uncapped(3),
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        weth.mint(alice, 1 * SHARE);
        vm.startPrank(alice);
        weth.approve(address(testnet), 1 * SHARE);
        testnet.deposit(CROSS, address(weth), 1 * SHARE, alice);
        vm.stopPrank();
        (equity,,,) = testnet.health(alice, CROSS);
        assertEq(equity, 0);
    }

    function test_LeverageIsGrossExposureOverEquity() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.leverage(alice, CROSS), 10_000);
        _borrow(alice, CROSS, 1_000 * USDG);
        assertEq(accounts.leverage(alice, CROSS), 20_000);
        _deposit(alice, CROSS, address(usdg), 1_000 * USDG);
        assertEq(accounts.leverage(alice, CROSS), 10_000);
        _quote(NVDA, 3, 0);
        assertEq(accounts.leverage(alice, CROSS), type(uint256).max);
    }

    function test_OnlyTheLiquidatorSeizesAndWritesOffAndNoOtherDebtMoves() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _borrow(bob, CROSS, 2_000 * USDG);
        address liquidator = makeAddr("liquidator");
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.setLiquidator(liquidator);
        vm.expectRevert(MarginAccounts.InvalidLiquidator.selector);
        vm.prank(owner);
        accounts.setLiquidator(address(0));
        vm.expectEmit(address(accounts));
        emit MarginAccounts.LiquidatorSet(liquidator);
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.LiquidatorAlreadySet.selector, liquidator));
        vm.prank(owner);
        accounts.setLiquidator(owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotLiquidator.selector, alice));
        vm.prank(alice);
        accounts.writeOff(alice, CROSS);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(liquidator);
        accounts.writeOff(alice, CROSS);
        _deposit(bob, SPY, address(usdg), 1 * USDG);
        _deposit(bob, NVDA, address(weth), 1);
        vm.startPrank(liquidator);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        accounts.writeOff(bob, SPY);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        accounts.writeOff(bob, NVDA);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotLiquidator.selector, alice));
        vm.prank(alice);
        accounts.seize(CROSS, address(nvda), 10 * SHARE, alice, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Seize(alice, CROSS, address(nvda), 10 * SHARE, liquidator);
        vm.prank(liquidator);
        accounts.seize(CROSS, address(nvda), 10 * SHARE, alice, liquidator);
        assertEq(nvda.balanceOf(liquidator), 10 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        uint256 shares = accounts.debtShares(alice, CROSS);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.WriteOff(alice, CROSS, 1_000 * USDG, shares, 0);
        vm.prank(liquidator);
        assertEq(accounts.writeOff(alice, CROSS), 1_000 * USDG);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.debt(bob, CROSS), 2_000 * USDG);
        assertEq(vault.debt(), 2_000 * USDG);
        assertEq(vault.totalAssets(), 99_000 * USDG);
    }

    function test_UsdgsIssuerErrorsReachTheCaller() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
        vm.prank(alice);
        usdg.approve(address(accounts), 100 * USDG);
        usdg.pause();
        vm.expectRevert(IUSDG.ContractPaused.selector);
        _borrow(alice, CROSS, 1);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        vm.prank(alice);
        accounts.repay(CROSS, 1, alice);
        usdg.unpause();
        usdg.freeze(address(vault));
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        _borrow(alice, CROSS, 1);
    }

    function test_GasOfEachCall() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        vm.snapshotGasLastCall("deposit");
        skip(1 hours);
        _borrow(alice, CROSS, 1_000 * USDG);
        vm.snapshotGasLastCall("borrow");
        skip(1 hours);
        _repay(alice, CROSS, 500 * USDG);
        vm.snapshotGasLastCall("repay");
        skip(1 hours);
        _withdraw(alice, CROSS, address(nvda), 1 * SHARE);
        vm.snapshotGasLastCall("withdraw");
    }

    function testFuzz_ThePositionMeetsItsRequirementAtTheLiquidationPriceAndNotBelow(uint256 borrowed, uint256 rate)
        public
    {
        borrowed = bound(borrowed, 1 * USDG, 1_900 * USDG);
        rate = bound(rate, 1_00, 50_00);
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        engine.set(rate, 0, 2);
        vm.assume(borrowed * 1e12 * 10_000 <= 2_000e18 * (10_000 - rate));
        _borrow(alice, CROSS, borrowed);
        uint256 price = accounts.liquidationPrice(alice, CROSS, NVDA, 0);
        _quote(NVDA, 3, uint64(price));
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertGe(equity, int256(requirement));
        if (price > STEP + 1) {
            _quote(NVDA, 3, uint64(price - STEP - 1));
            (equity, requirement,,) = accounts.health(alice, CROSS);
            assertLt(equity, int256(requirement));
        }
    }

    function test_APausedOrBlockedTokenStopsItsTransfersAndNewRiskAgainstIt() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(bob, CROSS, address(spy), 10 * SHARE);
        address liquidator = makeAddr("liquidator");
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        nvda.mint(alice, SHARE);
        nvda.pause();
        vm.startPrank(alice);
        nvda.approve(address(accounts), SHARE);
        vm.expectRevert(StubStockToken.IsPaused.selector);
        accounts.deposit(CROSS, address(nvda), SHARE, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.withdraw(CROSS, address(nvda), SHARE, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.withdraw(CROSS, address(weth), 1 * SHARE, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.borrow(CROSS, 100 * USDG, alice, alice);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertGt(equity, 0);
        vm.stopPrank();
        vm.startPrank(liquidator);
        vm.expectRevert(StubStockToken.IsPaused.selector);
        accounts.seize(CROSS, address(nvda), SHARE, alice, liquidator);
        vm.stopPrank();
        _borrow(bob, CROSS, 1_000 * USDG);
        _withdraw(bob, CROSS, address(spy), 1 * SHARE);
        nvda.unpause();
        nvda.blockAccount(address(accounts), true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        _withdraw(alice, CROSS, address(nvda), 1);
        vm.expectRevert(abi.encodeWithSelector(StubStockToken.Blocked.selector, address(accounts)));
        vm.prank(liquidator);
        accounts.seize(CROSS, address(nvda), SHARE, alice, liquidator);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        _borrow(alice, CROSS, 1);
        _borrow(bob, CROSS, 100 * USDG);
        _repay(alice, CROSS, 1_000 * USDG);
        assertEq(accounts.debt(alice, CROSS), 0);
    }

    function test_ABurnFromTheAccountsFallsOnThatAssetsHoldersAlone() public {
        address carol = makeAddr("carol");
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, NVDA, address(nvda), 30 * SHARE);
        _deposit(carol, CROSS, address(spy), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        nvda.adminBurn(address(accounts), 20 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.AssetWrittenDown(NVDA, 40 * SHARE, 20 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 5 * SHARE);
        assertEq(accounts.collateral(bob, NVDA, address(nvda)), 15 * SHARE);
        assertEq(accounts.collateral(carol, CROSS, address(spy)), 10 * SHARE);
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertEq(units, 40 * SHARE);
        assertEq(scale, 0.5e18);
        address dave = makeAddr("dave");
        _deposit(dave, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.collateral(dave, CROSS, address(nvda)), 10 * SHARE);
        nvda.mint(address(accounts), 5 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(bob, NVDA, address(nvda)), 15 * SHARE);
        _withdraw(bob, NVDA, address(nvda), 15 * SHARE);
        assertEq(nvda.balanceOf(bob), 15 * SHARE);
        assertEq(accounts.collateral(bob, NVDA, address(nvda)), 0);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(alice, CROSS, 1);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 5 * SHARE);
        accounts.sync(NVDA);
        (, scale,) = accounts.holding(NVDA);
        assertEq(scale, uint256(10 * SHARE) * 1e18 / (30 * SHARE));
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE * scale / 1e18);
        assertEq(accounts.collateral(dave, CROSS, address(nvda)), 20 * SHARE * scale / 1e18);
        assertLe(
            accounts.collateral(alice, CROSS, address(nvda)) + accounts.collateral(dave, CROSS, address(nvda)),
            nvda.balanceOf(address(accounts))
        );
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, TSLA));
        accounts.sync(TSLA);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, bytes32("AMD")));
        accounts.holding("AMD");
    }

    function test_EveryCallTouchingABurnedAssetCountsTheBurnFirst() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(bob, NVDA, address(nvda), 10 * SHARE);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _borrow(alice, CROSS, 1_500 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        _withdraw(alice, CROSS, address(weth), 1 * SHARE);
        address carol = makeAddr("carol");
        _deposit(carol, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.collateral(carol, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.collateral(bob, NVDA, address(nvda)), 5 * SHARE);
    }

    function test_AWithdrawalAfterABurnTakesNoMoreThanItsShare() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, NVDA, address(nvda), 30 * SHARE);
        nvda.adminBurn(address(accounts), 20 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 30 * SHARE)
        );
        _withdraw(bob, NVDA, address(nvda), 30 * SHARE);
        _withdraw(bob, NVDA, address(nvda), 15 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 5 * SHARE);
    }

    function test_APositionThatWithdrewItsStocksNeedsNoSession() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _withdraw(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        engine.set(20_00, 0, 0);
        _borrow(alice, CROSS, 100 * USDG);
    }

    function test_AHoldingWorthNothingAfterAWithdrawalIsCleared() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _deposit(bob, NVDA, address(nvda), 20 * SHARE);
        nvda.adminBurn(address(accounts), 20 * SHARE);
        accounts.sync(NVDA);
        (, uint256 scale,) = accounts.holding(NVDA);
        _withdraw(alice, CROSS, address(nvda), (10 * SHARE - 2) * scale / 1e18);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        _quote(NVDA, 0, 0);
        _borrow(alice, CROSS, 100 * USDG);
        (uint256 units,,) = accounts.holding(NVDA);
        assertEq(units, 20 * SHARE);
    }

    function test_AnAssetBurnedEntirelyTakesNoDepositUntilItsHoldingsAreWrittenOff() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        address liquidator = makeAddr("liquidator");
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        nvda.mint(bob, SHARE);
        vm.startPrank(bob);
        nvda.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetWrittenOff.selector, NVDA));
        accounts.deposit(CROSS, address(nvda), SHARE, bob);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 1));
        _withdraw(alice, CROSS, address(nvda), 1);
        vm.prank(liquidator);
        assertEq(accounts.writeOff(alice, CROSS), 1_000 * USDG);
        (uint256 units,,) = accounts.holding(NVDA);
        assertEq(units, 0);
        assertEq(vault.debt(), 0);
        _deposit(bob, CROSS, address(nvda), 1 * SHARE);
        (, uint256 scale,) = accounts.holding(NVDA);
        assertEq(scale, 1e18);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 1 * SHARE);
    }

    function test_AHoldingABurnLeftWorthNothingIsWrittenOffWithItsDebt() public {
        _deposit(alice, CROSS, address(nvda), 1);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        nvda.adminBurn(address(accounts), 1 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        _quote(NVDA, 0, 0);
        _borrow(alice, CROSS, 100 * USDG);
        address liquidator = makeAddr("liquidator");
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        vm.startPrank(liquidator);
        accounts.seize(CROSS, address(weth), 1 * SHARE, alice, liquidator);
        assertEq(accounts.writeOff(alice, CROSS), 1_100 * USDG);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        accounts.writeOff(bob, CROSS);
        vm.stopPrank();
        (uint256 units,,) = accounts.holding(NVDA);
        assertEq(units, 10 * SHARE);
    }

    function test_AHoldingABurnLeftWorthNothingDoesNotBlockIsolatingTheAsset() public {
        _deposit(alice, CROSS, address(nvda), 1);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        nvda.adminBurn(address(accounts), 1 * SHARE);
        accounts.sync(NVDA);
        _deposit(alice, NVDA, address(nvda), 1 * SHARE);
        assertApproxEqAbs(accounts.collateral(alice, NVDA, address(nvda)), 1 * SHARE, 1);
    }

    function test_ABlocklistedAccountMovesNoStockTokenThroughTheAccounts() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
        nvda.blockAccount(alice, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        _withdraw(alice, CROSS, address(nvda), 1 * SHARE);
        vm.prank(alice);
        accounts.setAuthorization(bob, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        vm.prank(bob);
        accounts.withdraw(CROSS, address(nvda), 1 * SHARE, alice, bob);
        nvda.mint(bob, SHARE);
        vm.startPrank(bob);
        nvda.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        accounts.deposit(CROSS, address(nvda), SHARE, alice);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        _withdraw(alice, CROSS, address(weth), 1 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        _borrow(alice, CROSS, 1 * USDG);
        _repay(alice, CROSS, 100 * USDG);
        _withdraw(alice, CROSS, address(weth), 1 * SHARE);
        nvda.blockAccount(alice, false);
        nvda.blockAccount(bob, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, bob));
        vm.prank(bob);
        accounts.withdraw(CROSS, address(nvda), 1 * SHARE, alice, alice);
        address liquidator = makeAddr("liquidator");
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        nvda.blockAccount(alice, true);
        vm.prank(liquidator);
        accounts.seize(CROSS, address(nvda), 1 * SHARE, alice, liquidator);
        assertEq(nvda.balanceOf(liquidator), 1 * SHARE);
    }

    function test_DepositsDebtAndWeekendDebtStopAtTheirCaps() public {
        uint256[] memory caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (20 * SHARE, 0, 0);
        _deploy(caps, 3_000 * USDG, 1_000 * USDG);
        assertEq(accounts.debtCap(), 3_000 * USDG);
        assertEq(accounts.weekendDebtCap(), 1_000 * USDG);
        (,, uint256 cap) = accounts.holding(NVDA);
        assertEq(cap, 20 * SHARE);
        _deposit(alice, CROSS, address(nvda), 15 * SHARE);
        nvda.mint(bob, 6 * SHARE);
        vm.startPrank(bob);
        nvda.approve(address(accounts), 6 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 20 * SHARE));
        accounts.deposit(NVDA, address(nvda), 6 * SHARE, bob);
        accounts.deposit(NVDA, address(nvda), 5 * SHARE, bob);
        vm.stopPrank();
        spy.mint(bob, SHARE);
        vm.startPrank(bob);
        spy.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, SPY, 0));
        accounts.deposit(CROSS, address(spy), SHARE, bob);
        vm.stopPrank();
        _borrow(alice, CROSS, 1_500 * USDG);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.DebtCapExceeded.selector, 3_001 * USDG, 3_000 * USDG));
        _borrow(bob, NVDA, 1_501 * USDG);
        _borrow(bob, NVDA, 500 * USDG);
        band.setSession(BandDouble.Session(2, 1, 3, 0, uint64(block.timestamp + 10 hours) * 1000));
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.WeekendDebtCapExceeded.selector, 2_001 * USDG, 1_000 * USDG)
        );
        _borrow(bob, NVDA, 1 * USDG);
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.WeekendDebtCapExceeded.selector, 2_001 * USDG, 1_000 * USDG)
        );
        _borrow(bob, NVDA, 1 * USDG);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.WeekendDebtCapExceeded.selector, 2_001 * USDG, 1_000 * USDG)
        );
        _borrow(bob, NVDA, 1 * USDG);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        _borrow(bob, NVDA, 1 * USDG);
    }

    function test_TheGuardianPausesBorrowingAndWithdrawalsInDebtOnly() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        _deposit(bob, CROSS, address(spy), 1 * SHARE);
        address guardian = makeAddr("guardian");
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.setGuardian(guardian);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.GuardianSet(guardian);
        vm.prank(owner);
        accounts.setGuardian(guardian);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotGuardian.selector, owner));
        vm.prank(owner);
        accounts.setBorrowingPaused(true);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.BorrowingPausedSet(true);
        vm.prank(guardian);
        accounts.setBorrowingPaused(true);
        assertTrue(accounts.borrowingPaused());
        vm.expectRevert(MarginAccounts.BorrowingIsPaused.selector);
        _borrow(bob, CROSS, 1);
        vm.expectRevert(MarginAccounts.BorrowingIsPaused.selector);
        _withdraw(alice, CROSS, address(weth), 1);
        _withdraw(bob, CROSS, address(spy), 1 * SHARE);
        _deposit(alice, CROSS, address(nvda), 1 * SHARE);
        _repay(alice, CROSS, 1_000 * USDG);
        _withdraw(alice, CROSS, address(weth), 1 * SHARE);
        vm.prank(guardian);
        accounts.setBorrowingPaused(false);
        _borrow(alice, CROSS, 100 * USDG);
    }

    function test_ThePremiumAccruesExactlyOverEachClosureOfAMonthOfSessions() public {
        vm.pauseGasMetering();
        uint256[][] memory rows = abi.decode(vm.parseJson(vm.readFile(TIMELINE), ".rows"), (uint256[][]));
        uint256 shift = (vm.getBlockTimestamp() + 1 hours) * 1000 - rows[0][0];
        _borrowWithoutInterest(alice, 10_000 * USDG);
        uint256 state = vm.snapshotState();
        uint256 everyQuarterHour = _replay(rows, shift, 1);
        vm.revertToState(state);
        uint256 everyFiveHours = _replay(rows, shift, 20);
        uint256 closedMs = (4 * 48 hours + 72 hours) * 1000;
        assertApproxEqAbs(everyQuarterHour, _premiumOver(10_000 * USDG, closedMs), 1);
        assertApproxEqAbs(everyFiveHours, everyQuarterHour, 1);
    }

    function test_AClosureMissedAtItsBoundariesAccruesOnlyTimeSeenClosed() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        skip(10 hours);
        uint64 sighting = uint64(vm.getBlockTimestamp()) * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        vm.expectEmit(address(accounts));
        emit MarginAccounts.ClosureSet(sighting, 0);
        accounts.accruePremium();
        skip(20 hours);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        accounts.accruePremium();
        (uint64 closesMs, uint64 reopensMs, uint64 accruedMs) = accounts.closure();
        assertEq(closesMs, sighting);
        assertEq(reopensMs, 0);
        assertEq(accruedMs, sighting);
        assertEq(accounts.premium(alice, CROSS), 0);
        skip(10 hours);
        uint64 reopens = uint64(vm.getBlockTimestamp() + 8 hours) * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, reopens));
        vm.expectEmit(address(accounts));
        emit MarginAccounts.ClosureSet(sighting, reopens);
        accounts.accruePremium();
        assertApproxEqAbs(accounts.premium(alice, CROSS), _premiumOver(10_000 * USDG, 30 hours * 1000), 1);
        skip(13 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        uint256 premium = accounts.premium(alice, CROSS);
        assertApproxEqAbs(premium, _premiumOver(10_000 * USDG, 38 hours * 1000), 1);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        skip(10 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        (, reopensMs,) = accounts.closure();
        assertEq(reopensMs, closes);
        assertEq(accounts.premium(alice, CROSS), premium);
        skip(7 days);
        sighting = uint64(vm.getBlockTimestamp()) * 1000;
        band.setSession(BandDouble.Session(1, 3, 3, 0, 0));
        accounts.accruePremium();
        skip(5 hours);
        accounts.accruePremium();
        skip(10 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        vm.expectEmit(address(accounts));
        emit MarginAccounts.ClosureSet(sighting, sighting + 5 hours * 1000);
        accounts.accruePremium();
        assertApproxEqAbs(accounts.premium(alice, CROSS), premium + _premiumOver(10_000 * USDG, 5 hours * 1000), 1);
    }

    function test_ALaterClosureOrAnUnknownSessionNeverChargesOpenTime() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(7 days + 1 hours);
        uint64 sighting = uint64(vm.getBlockTimestamp()) * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, sighting + 40 hours * 1000));
        accounts.accruePremium();
        assertEq(accounts.premium(alice, CROSS), 0);
        (uint64 closesMs,,) = accounts.closure();
        assertEq(closesMs, sighting);
        skip(41 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        uint256 premium = accounts.premium(alice, CROSS);
        assertApproxEqAbs(premium, _premiumOver(10_000 * USDG, 40 hours * 1000), 1);
        closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        skip(78 hours);
        accounts.accruePremium();
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        assertEq(accounts.premium(alice, CROSS), premium);
        closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(7 days + 1 hours);
        band.setSession(BandDouble.Session(1, 3, 3, 0, 0));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        assertEq(accounts.premium(alice, CROSS), premium);
    }

    function test_ARepaymentPaysTheDebtBeforeThePremiumAndSplitsThePremium() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        assertApproxEqAbs(premium, _premiumOver(10_000 * USDG, 48 hours * 1000), 1);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 200_000e18 - int256((10_000 * USDG + premium) * 1e12));
        _repay(alice, CROSS, 4_000 * USDG);
        assertEq(accounts.debt(alice, CROSS), 6_000 * USDG);
        assertEq(accounts.premium(alice, CROSS), premium);
        usdg.mint(bob, 7_000 * USDG);
        vm.startPrank(bob);
        usdg.approve(address(accounts), 7_000 * USDG);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.PremiumPaid(alice, CROSS, premium, premium / 10);
        assertEq(accounts.repay(CROSS, 7_000 * USDG, alice), 6_000 * USDG + premium);
        vm.stopPrank();
        assertEq(accounts.premium(alice, CROSS), 0);
        assertEq(accounts.reserve(), premium / 10);
        assertEq(accounts.backstopPremium(), premium - premium / 10);
        assertEq(usdg.balanceOf(address(accounts)), premium);
        assertEq(vault.debt(), 0);
        assertEq(vault.totalAssets(), 100_000 * USDG);
        assertEq(usdg.balanceOf(bob), 1_000 * USDG - premium);
    }

    function test_APositionPaysItsPremiumFromItsUsdg() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        _deposit(alice, CROSS, address(usdg), 10_000 * USDG + premium + 1);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.PremiumPaid(alice, CROSS, premium, premium / 10);
        vm.prank(alice);
        assertEq(accounts.repayWithCollateral(CROSS, type(uint256).max, alice), 10_000 * USDG + premium);
        assertEq(accounts.collateral(alice, CROSS, address(usdg)), 1);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.premium(alice, CROSS), 0);
        assertEq(usdg.balanceOf(address(accounts)), premium + 1);
    }

    function test_ThePremiumIsOwedBesideTheDebtButNotCountedAgainstTheCaps() public {
        _deploy(_uncapped(3), 10_000 * USDG, 10_000 * USDG);
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        assertGt(premium, 0);
        _repay(alice, CROSS, accounts.debt(alice, CROSS));
        assertEq(accounts.debtShares(alice, CROSS), 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientMargin.selector, -int256(premium * 1e12), 0));
        _withdraw(alice, CROSS, address(nvda), 1_000 * SHARE);
        _deposit(bob, CROSS, address(nvda), 1_000 * SHARE);
        _borrow(bob, CROSS, 10_000 * USDG);
        assertEq(vault.debt(), 10_000 * USDG);
        _repay(alice, CROSS, premium);
        _withdraw(alice, CROSS, address(nvda), 1_000 * SHARE);
    }

    function test_ALoanOrRepaymentSettlesThePremiumOwedSoFar() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _deposit(bob, CROSS, address(spy), 100 * SHARE);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        skip(24 hours);
        _borrow(bob, CROSS, 10_000 * USDG);
        assertEq(accounts.premium(bob, CROSS), 0);
        _borrow(alice, CROSS, 10_000 * USDG);
        uint256 half = _premiumOver(10_000 * USDG, 24 hours * 1000);
        assertApproxEqAbs(accounts.premium(alice, CROSS), half, 1);
        skip(24 hours);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        assertApproxEqAbs(accounts.premium(alice, CROSS), half + 2 * half, 2);
        assertApproxEqAbs(accounts.premium(bob, CROSS), half, 1);
        skip(24 hours);
        assertApproxEqAbs(accounts.premium(bob, CROSS), half, 1);
        _repay(alice, CROSS, 1);
        assertApproxEqAbs(accounts.premium(alice, CROSS), half + 2 * half, 2);
    }

    function test_AWithdrawalCountsThePremiumAccruedSinceTheLastCall() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 72 hours * 1000));
        accounts.accruePremium();
        skip(72 hours);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        uint256 premium = accounts.premium(alice, CROSS);
        assertApproxEqAbs(premium, 20 * _premiumOver(10_000 * USDG, 72 hours * 1000), 1);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 200_000e18 - int256((10_000 * USDG + premium) * 1e12));
        vm.expectRevert(
            abi.encodeWithSelector(
                MarginAccounts.InsufficientMargin.selector,
                12_600e18 - int256((10_000 * USDG + premium) * 1e12),
                2_520e18
            )
        );
        _withdraw(alice, CROSS, address(nvda), 937 * SHARE);
        _withdraw(alice, CROSS, address(nvda), 936 * SHARE);
    }

    function test_TheOwnerSetsThePremiumRateWithinItsBoundFromNowOn() public {
        assertEq(accounts.premiumRate(), 5_00);
        assertEq(accounts.reserveShare(), 10_00);
        _borrowWithoutInterest(alice, 10_000 * USDG);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        skip(24 hours);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.setPremiumRate(1);
        vm.expectRevert(MarginAccounts.InvalidPremium.selector);
        vm.prank(owner);
        accounts.setPremiumRate(100_01);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.PremiumRateSet(100_00);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        assertEq(accounts.premiumRate(), 100_00);
        (uint64 closesMs, uint64 reopensMs, uint64 accruedMs) = accounts.closure();
        assertEq(closesMs, closes);
        assertEq(reopensMs, closes + 48 hours * 1000);
        assertEq(accruedMs, closes + 24 hours * 1000);
        skip(24 hours);
        uint256 half = _premiumOver(10_000 * USDG, 24 hours * 1000);
        assertApproxEqAbs(accounts.premium(alice, CROSS), half + 20 * half, 2);
    }

    function test_ThePremiumRateIsBoundedAtDeployment() public {
        vm.expectRevert(MarginAccounts.InvalidPremium.selector);
        new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), vault, weth, owner, _uncapped(3), 0, 0, 100_01, 0
        );
    }

    function test_TheReserveShareIsAtMostTheWhole() public {
        vm.expectRevert(MarginAccounts.InvalidPremium.selector);
        new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), vault, weth, owner, _uncapped(3), 0, 0, 0, 100_01
        );
    }

    function test_OnlyTheBackstopClaimsThePremiumSetAside() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        _repay(alice, CROSS, 10_000 * USDG + premium);
        address backstop = makeAddr("backstop");
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotBackstop.selector, alice));
        vm.prank(alice);
        accounts.claimPremium();
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.setBackstop(backstop);
        vm.expectRevert(MarginAccounts.InvalidBackstop.selector);
        vm.prank(owner);
        accounts.setBackstop(address(0));
        vm.expectEmit(address(accounts));
        emit MarginAccounts.BackstopSet(backstop);
        vm.prank(owner);
        accounts.setBackstop(backstop);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BackstopAlreadySet.selector, backstop));
        vm.prank(owner);
        accounts.setBackstop(owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotBackstop.selector, owner));
        vm.prank(owner);
        accounts.claimPremium();
        vm.expectEmit(address(accounts));
        emit MarginAccounts.PremiumClaimed(backstop, premium - premium / 10);
        vm.prank(backstop);
        assertEq(accounts.claimPremium(), premium - premium / 10);
        assertEq(usdg.balanceOf(backstop), premium - premium / 10);
        assertEq(accounts.backstopPremium(), 0);
        assertEq(usdg.balanceOf(address(accounts)), premium / 10);
    }

    function test_OnlyTheOwnerWithdrawsTheReserve() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        _repay(alice, CROSS, 10_000 * USDG + premium);
        uint256 reserve = premium / 10;
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.withdrawReserve(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientReserve.selector, reserve));
        vm.prank(owner);
        accounts.withdrawReserve(owner, reserve + 1);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.ReserveWithdrawn(bob, reserve);
        vm.prank(owner);
        accounts.withdrawReserve(bob, reserve);
        assertEq(usdg.balanceOf(bob), reserve);
        assertEq(accounts.reserve(), 0);
        assertEq(usdg.balanceOf(address(accounts)), premium - reserve);
    }

    function test_AWriteOffForgivesThePremium() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        uint256 premium = accounts.premium(alice, CROSS);
        address liquidator = makeAddr("liquidator");
        vm.prank(owner);
        accounts.setLiquidator(liquidator);
        vm.startPrank(liquidator);
        accounts.seize(CROSS, address(nvda), 1_000 * SHARE, alice, liquidator);
        uint256 shares = accounts.debtShares(alice, CROSS);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.WriteOff(alice, CROSS, 10_000 * USDG, shares, premium);
        assertEq(accounts.writeOff(alice, CROSS), 10_000 * USDG);
        vm.stopPrank();
        assertEq(accounts.premium(alice, CROSS), 0);
        assertEq(accounts.reserve() + accounts.backstopPremium(), 0);
    }

    function test_GasOfThePremium() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        vm.snapshotGasLastCall("accruePremium");
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        _repay(alice, CROSS, 11_000 * USDG);
        vm.snapshotGasLastCall("repayWithPremium");
    }

    function test_FrozenAccountsLendNothing() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        usdg.freeze(address(accounts));
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        _borrow(alice, CROSS, 100 * USDG);
        usdg.unfreeze(address(accounts));
        _borrow(alice, CROSS, 100 * USDG);
    }

    function test_AnyoneClearsAHoldingABurnLeftWorthNothing() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _deposit(bob, CROSS, address(spy), 1 * SHARE);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        accounts.sync(NVDA);
        nvda.mint(bob, SHARE);
        vm.startPrank(bob);
        nvda.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetWrittenOff.selector, NVDA));
        accounts.deposit(NVDA, address(nvda), SHARE, bob);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.HoldingNotEmpty.selector, SPY));
        accounts.clear(bob, CROSS, SPY);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, TSLA));
        accounts.clear(alice, CROSS, TSLA);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.HoldingCleared(alice, CROSS, NVDA, 10 * SHARE);
        accounts.clear(alice, CROSS, NVDA);
        accounts.deposit(NVDA, address(nvda), SHARE, bob);
        vm.stopPrank();
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertEq(units, SHARE);
        assertEq(scale, 1e18);
        assertEq(accounts.collateral(alice, CROSS, address(weth)), 1 * SHARE);
    }

    function test_AHoldingABurnLeftWorthNothingGatesNothing() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        accounts.sync(NVDA);
        engine.set(20_00, 0, 0);
        nvda.blockAccount(alice, true);
        _borrow(alice, CROSS, 1_000 * USDG);
        _withdraw(alice, CROSS, address(weth), SHARE / 10);
    }

    function test_AFeedAnswerFromTheFutureCountsForNothing() public {
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        ethUsd.setRound(2_000e8, block.timestamp + 1);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 0);
    }

    function test_TheAccountsTakeAtMostEightAssets() public {
        bytes32[] memory symbols = new bytes32[](9);
        for (uint256 i; i < 9; ++i) {
            symbols[i] = bytes32(i + 1);
        }
        MarginDouble nine = new MarginDouble(symbols, address(band), address(ethUsd));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.TooManyAssets.selector, 9));
        new MarginAccounts(
            IBand(address(band)), IMargin(address(nine)), vault, weth, owner, _uncapped(9), 0, 0, 5_00, 10_00
        );
    }

    function test_TheWeekendDebtCapIsAtMostTheDebtCap() public {
        vm.expectRevert(MarginAccounts.InvalidDebtCaps.selector);
        new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), vault, weth, owner, _uncapped(3), 1, 2, 5_00, 10_00
        );
    }

    function test_TheInitialPremiumRateIsAnnounced() public {
        vm.expectEmit();
        emit MarginAccounts.PremiumRateSet(5_00);
        new MarginAccounts(
            IBand(address(band)), IMargin(address(engine)), vault, weth, owner, _uncapped(3), 0, 0, 5_00, 10_00
        );
    }

    function test_OwnershipCannotBeRenounced() public {
        vm.expectRevert(MarginAccounts.OwnershipCannotBeRenounced.selector);
        vm.prank(owner);
        accounts.renounceOwnership();
        assertEq(accounts.owner(), owner);
    }

    function test_ThePremiumRateChangesOnlyWhileTheSessionIsKnown() public {
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        vm.expectRevert(MarginAccounts.SessionUnknown.selector);
        vm.prank(owner);
        accounts.setPremiumRate(1);
        band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        vm.prank(owner);
        accounts.setPremiumRate(1);
    }

    function test_TheReserveGoesNowhereButToAnAddress() public {
        vm.expectRevert(MarginAccounts.InvalidReceiver.selector);
        vm.prank(owner);
        accounts.withdrawReserve(address(0), 0);
    }

    function test_PayingOnlyPremiumLeavesTheVaultAlone() public {
        _borrowWithoutInterest(alice, 10_000 * USDG);
        _closeFor(48 hours);
        _repay(alice, CROSS, 10_000 * USDG);
        uint256 premium = accounts.premium(alice, CROSS);
        vm.expectCall(address(vault), abi.encodeCall(SupplyVault.repay, (0)), 0);
        vm.recordLogs();
        _repay(alice, CROSS, premium);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        for (uint256 i; i < logs.length; ++i) {
            assertTrue(logs[i].emitter != address(accounts) || logs[i].topics[0] != MarginAccounts.Repay.selector);
        }
        assertEq(accounts.premium(alice, CROSS), 0);
    }

    function test_AGuardianRemovedLeavesTheOwnerAbleToResume() public {
        address guardian = makeAddr("guardian");
        vm.prank(owner);
        accounts.setGuardian(guardian);
        vm.prank(guardian);
        accounts.setBorrowingPaused(true);
        vm.prank(owner);
        accounts.setGuardian(address(0));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotGuardian.selector, guardian));
        vm.prank(guardian);
        accounts.setBorrowingPaused(false);
        vm.prank(owner);
        accounts.setGuardian(owner);
        vm.prank(owner);
        accounts.setBorrowingPaused(false);
        assertFalse(accounts.borrowingPaused());
    }

    function test_TheLiquidationPriceOfAnAssetNotHeldIsZero() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.liquidationPrice(alice, CROSS, SPY, 0), 0);
    }

    function test_ABlocklistedAccountInDebtWithdrawsNothing() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(alice, CROSS, address(weth), 1 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
        nvda.blockAccount(alice, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        _withdraw(alice, CROSS, address(weth), 1);
    }

    function test_RoundingMovesTheSharePriceByAtMostAUnitACall() public {
        _deposit(alice, CROSS, address(nvda), 1_000 * SHARE);
        _deposit(bob, CROSS, address(spy), 10 * SHARE);
        _borrow(bob, CROSS, 1);
        for (uint256 cycle; cycle < 100; ++cycle) {
            _borrow(alice, CROSS, 10_000 * USDG);
            skip(1 days);
            _repay(alice, CROSS, accounts.debt(alice, CROSS));
            assertEq(accounts.debtShares(alice, CROSS), 0);
            assertLe(accounts.debt(bob, CROSS), vault.debt());
        }
        assertLe(accounts.totalDebtShares(), (vault.debt() + 1) * 1e6 + 100 * 2e6);
    }

    function _quote(bytes32 symbol, uint8 state, uint64 low) internal {
        band.setQuote(symbol, BandDouble.Quote(state, state == 0 ? 0 : 3, low, 50, low, low + low / 100));
    }

    function _deposit(address account, bytes32 position, address token, uint256 amount) internal {
        if (token == address(usdg)) usdg.mint(account, amount);
        else StubToken(token).mint(account, amount);
        vm.startPrank(account);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, account);
        vm.stopPrank();
    }

    function _withdraw(address account, bytes32 position, address token, uint256 amount) internal {
        vm.prank(account);
        accounts.withdraw(position, token, amount, account, account);
    }

    function _borrow(address account, bytes32 position, uint256 assets) internal {
        vm.prank(account);
        accounts.borrow(position, assets, account, account);
    }

    function _repay(address account, bytes32 position, uint256 assets) internal {
        usdg.mint(account, assets);
        vm.startPrank(account);
        usdg.approve(address(accounts), assets);
        accounts.repay(position, assets, account);
        vm.stopPrank();
    }

    function _deploy(uint256[] memory caps, uint256 debtCap, uint256 weekendDebtCap) internal {
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            weth,
            owner,
            caps,
            debtCap,
            weekendDebtCap,
            5_00,
            10_00
        );
        vm.prank(owner);
        vault.setBorrower(address(accounts));
        usdg.mint(lender, 100_000 * USDG);
        vm.startPrank(lender);
        usdg.approve(address(vault), 100_000 * USDG);
        vault.deposit(100_000 * USDG, lender);
        vm.stopPrank();
    }

    function _borrowWithoutInterest(address account, uint256 assets) internal {
        vm.prank(owner);
        vault.setRateModel(SupplyVault.RateModel(90_00, 0, 0, 0));
        _deposit(account, CROSS, address(nvda), 1_000 * SHARE);
        _borrow(account, CROSS, assets);
    }

    function _closeFor(uint256 closed) internal {
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + uint64(closed) * 1000));
        accounts.accruePremium();
        skip(closed);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        accounts.accruePremium();
    }

    function _replay(uint256[][] memory rows, uint256 shift, uint256 stride) internal returns (uint256) {
        for (uint256 k; k < rows.length; k += stride) {
            uint256[] memory row = rows[k];
            vm.warp((row[0] + shift) / 1000);
            band.setSession(
                BandDouble.Session(
                    uint8(row[1]),
                    uint8(row[2]),
                    uint8(row[3]),
                    uint64(row[4] == 0 ? 0 : row[4] + shift),
                    uint64(row[5] == 0 ? 0 : row[5] + shift)
                )
            );
            accounts.accruePremium();
        }
        return accounts.premium(alice, CROSS);
    }

    function _premiumOver(uint256 debt, uint256 closedMs) internal pure returns (uint256) {
        return debt * 5_00 * closedMs / (10_000 * YEAR_MS);
    }

    function _uncapped(uint256 n) internal pure returns (uint256[] memory caps) {
        caps = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            caps[i] = type(uint256).max;
        }
    }
}
