// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
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

contract StockLendingVaultTest is Test {
    uint256 internal constant SHARE = 1e18;

    StubStockToken internal nvda;
    StockLendingVault internal lending;
    address internal owner = makeAddr("owner");
    address internal borrower = makeAddr("borrower");
    address internal other = makeAddr("other");

    function setUp() public {
        vm.warp(1_790_000_000);
        nvda = new StubStockToken(1e18);
        lending =
            new StockLendingVault(IERC20(address(nvda)), owner, SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.startPrank(owner);
        lending.setDepositor(address(this));
        lending.setBorrower(borrower);
        vm.stopPrank();
        nvda.approve(address(lending), type(uint256).max);
        vm.prank(borrower);
        nvda.approve(address(lending), type(uint256).max);
    }

    function test_TheVaultLendsOneStockTokenAndNamesItself() public {
        assertEq(lending.asset(), address(nvda));
        assertEq(lending.name(), "Tapehouse STOCK Lending");
        assertEq(lending.symbol(), "thlSTOCK");
        assertEq(lending.decimals(), 18);
        assertEq(lending.feeShare(), 10_00);
        assertEq(lending.MAX_FEE_SHARE(), 20_00);
        assertEq(lending.MAX_UTILIZATION(), 90_00);
        SupplyVault.RateModel memory m = lending.rateModel();
        assertEq(m.optimal, 80_00);
        assertEq(m.base, 25);
        assertEq(m.slope1, 1_00);
        assertEq(m.slope2, 50_00);
        vm.expectRevert(StockLendingVault.InvalidFeeShare.selector);
        this.deploy(SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 20_01);
        vm.expectRevert(StockLendingVault.InvalidRateModel.selector);
        this.deploy(SupplyVault.RateModel(99_01, 0, 0, 0), 0);
        vm.expectRevert(StockLendingVault.InvalidRateModel.selector);
        this.deploy(SupplyVault.RateModel(80_00, 0, 2_00, 1_00), 0);
        vm.expectRevert(StockLendingVault.InvalidRateModel.selector);
        this.deploy(SupplyVault.RateModel(80_00, 1, 0, 1000_00), 0);
        assertEq(this.deploy(SupplyVault.RateModel(1_00, 0, 0, 1000_00), 20_00).feeShare(), 20_00);
    }

    function test_OnlyTheOwnerSetsTheDepositorAndTheBorrowerOnceAndKeepsItsOwnership() public {
        StockLendingVault fresh = this.deploy(SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, other));
        vm.prank(other);
        fresh.setDepositor(other);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, other));
        vm.prank(other);
        fresh.setBorrower(other);
        vm.startPrank(owner);
        vm.expectRevert(StockLendingVault.InvalidAddress.selector);
        fresh.setDepositor(address(0));
        vm.expectRevert(StockLendingVault.InvalidAddress.selector);
        fresh.setBorrower(address(0));
        vm.expectEmit(address(fresh));
        emit StockLendingVault.DepositorSet(other);
        fresh.setDepositor(other);
        vm.expectEmit(address(fresh));
        emit StockLendingVault.BorrowerSet(borrower);
        fresh.setBorrower(borrower);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.AlreadySet.selector, other));
        fresh.setDepositor(borrower);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.AlreadySet.selector, borrower));
        fresh.setBorrower(other);
        vm.expectRevert(StockLendingVault.OwnershipCannotBeRenounced.selector);
        fresh.renounceOwnership();
        vm.stopPrank();
    }

    function test_OnlyTheDepositorDeposits() public {
        nvda.mint(other, 10 * SHARE);
        vm.startPrank(other);
        nvda.approve(address(lending), type(uint256).max);
        assertEq(lending.maxDeposit(other), 0);
        assertEq(lending.maxMint(other), 0);
        vm.expectRevert(abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxDeposit.selector, other, SHARE, 0));
        lending.deposit(SHARE, other);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotDepositor.selector, other));
        lending.deposit(SHARE, address(this));
        vm.stopPrank();
        assertEq(lending.maxDeposit(address(this)), type(uint256).max);
        _lend(100 * SHARE);
        assertEq(lending.balanceOf(address(this)), 100 * SHARE);
        assertEq(lending.idle(), 100 * SHARE);
        assertEq(lending.totalAssets(), 100 * SHARE);
    }

    function test_TheBorrowerTakesAtMostNinetyPercentAndPaysAKinkedFeeInTheToken() public {
        _lend(100 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotBorrower.selector, other));
        vm.prank(other);
        lending.borrow(SHARE, other);
        assertEq(lending.borrowable(), 90 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(StockLendingVault.InsufficientLiquidity.selector, 91 * SHARE, 90 * SHARE)
        );
        vm.prank(borrower);
        lending.borrow(91 * SHARE, borrower);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Borrow(borrower, 72 * SHARE);
        vm.prank(borrower);
        lending.borrow(72 * SHARE, borrower);
        assertEq(nvda.balanceOf(borrower), 72 * SHARE);
        assertEq(lending.utilization(), 0.72e18);
        assertEq(lending.borrowRate(), 0.0025e18 + 0.01e18 * 72 / 80);
        vm.prank(borrower);
        lending.borrow(18 * SHARE, borrower);
        assertEq(lending.borrowable(), 0);
        assertEq(lending.utilization(), 0.9e18);
        uint256 rate = 0.0025e18 + 0.01e18 + 0.5e18 / 2;
        assertEq(lending.borrowRate(), rate);
        assertEq(lending.supplyRate(), rate * 9 / 10 * 9 / 10);
        skip(365 days);
        uint256 debt = 90 * SHARE + 90 * rate;
        assertEq(lending.debt(), debt);
        uint256 fee = debt - 90 * SHARE;
        nvda.mint(address(this), SHARE);
        vm.expectEmit(address(lending));
        emit StockLendingVault.FeeAccrued(fee, 1e27 + rate * 1e9, _feeShares(fee / 10, 100 * SHARE, 10 * SHARE + debt));
        lending.deposit(SHARE, address(this));
        assertApproxEqAbs(lending.convertToAssets(lending.balanceOf(owner)), fee / 10, 2);
        assertApproxEqAbs(lending.convertToAssets(lending.balanceOf(address(this))), 101 * SHARE + fee * 9 / 10, 2);
        assertEq(lending.borrowable(), 0);
    }

    function test_TheBorrowerRepaysAndNoOneElse() public {
        _lend(100 * SHARE);
        vm.prank(borrower);
        lending.borrow(50 * SHARE, borrower);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotBorrower.selector, other));
        vm.prank(other);
        lending.repay(SHARE);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Repay(20 * SHARE);
        vm.prank(borrower);
        assertEq(lending.repay(20 * SHARE), 20 * SHARE);
        assertEq(lending.idle(), 70 * SHARE);
        assertEq(lending.debt(), 30 * SHARE);
        assertEq(lending.totalAssets(), 100 * SHARE);
        vm.prank(borrower);
        assertEq(lending.repay(100 * SHARE), 30 * SHARE);
        assertEq(lending.debt(), 0);
        assertEq(lending.scaledDebt(), 0);
    }

    function test_RedemptionsAreCappedByTheTokensTheVaultHolds() public {
        _lend(100 * SHARE);
        vm.prank(borrower);
        lending.borrow(90 * SHARE, borrower);
        assertEq(lending.maxWithdraw(address(this)), 10 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxWithdraw.selector, address(this), 11 * SHARE, 10 * SHARE)
        );
        lending.withdraw(11 * SHARE, address(this), address(this));
        lending.withdraw(10 * SHARE, address(this), address(this));
        assertEq(lending.idle(), 0);
        assertEq(nvda.balanceOf(address(this)), 10 * SHARE);
    }

    function test_ThePauseTheBlocklistAndABurnStopTheViewsAndABurnFallsOnThisVaultsLendersAlone() public {
        _lend(100 * SHARE);
        nvda.pause();
        assertEq(lending.maxDeposit(address(this)), 0);
        assertEq(lending.maxMint(address(this)), 0);
        assertEq(lending.maxRedeem(address(this)), 0);
        assertEq(lending.maxWithdraw(address(this)), 0);
        nvda.unpause();
        nvda.blockAccount(address(lending), true);
        assertEq(lending.maxDeposit(address(this)), 0);
        assertEq(lending.maxRedeem(address(this)), 0);
        nvda.blockAccount(address(lending), false);
        nvda.adminBurn(address(lending), 30 * SHARE);
        assertEq(lending.maxRedeem(address(this)), 0);
        nvda.mint(address(lending), 5 * SHARE);
        assertEq(lending.maxRedeem(address(this)), 0);
        nvda.adminBurn(address(lending), 5 * SHARE);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Sync(70 * SHARE);
        lending.sync();
        assertEq(lending.idle(), 70 * SHARE);
        assertEq(lending.totalAssets(), 70 * SHARE);
        assertEq(lending.convertToAssets(lending.balanceOf(address(this))), 70 * SHARE);
        nvda.mint(address(lending), 5 * SHARE);
        lending.sync();
        assertEq(lending.idle(), 70 * SHARE);
        assertEq(lending.maxWithdraw(address(this)), 70 * SHARE - 1);
        vm.mockCallRevert(address(nvda), abi.encodeWithSignature("paused()"), "");
        assertEq(lending.maxRedeem(address(this)), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(nvda), abi.encodeWithSignature("isBlocked(address)", address(lending)), "");
        assertEq(lending.maxRedeem(address(this)), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(nvda), abi.encodeWithSignature("balanceOf(address)", address(lending)), "");
        assertEq(lending.maxRedeem(address(this)), 0);
    }

    function test_TheOwnerSetsTheRateModelWithinItsBoundsAfterTheFeeAccrues() public {
        _lend(100 * SHARE);
        vm.prank(borrower);
        lending.borrow(80 * SHARE, borrower);
        skip(30 days);
        uint256 owed = lending.debt();
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, other));
        vm.prank(other);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        vm.expectRevert(StockLendingVault.InvalidRateModel.selector);
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(0, 0, 0, 0));
        vm.expectEmit(address(lending));
        emit StockLendingVault.RateModelSet(80_00, 0, 0, 0);
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 0, 0));
        assertApproxEqAbs(lending.borrowIndex() * 80 * SHARE / 1e27, owed, 1);
        assertEq(lending.lastAccrual(), block.timestamp);
        skip(30 days);
        assertEq(lending.debt(), owed);
    }

    function test_TheFeeAccruesBeforeEveryChangeAndEveryRoundingFavoursTheLenders() public {
        _lend(100 * SHARE);
        skip(30 days);
        _lend(SHARE);
        assertEq(lending.borrowIndex(), 1e27);
        vm.prank(borrower);
        lending.borrow(80 * SHARE, borrower);
        vm.recordLogs();
        _lend(SHARE);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        for (uint256 i; i < logs.length; ++i) {
            assertTrue(logs[i].topics[0] != StockLendingVault.FeeAccrued.selector);
        }
        skip(365 days);
        uint256 worth = lending.convertToAssets(lending.balanceOf(address(this)));
        uint256 owed = lending.debt();
        lending.withdraw(10 * SHARE, address(this), address(this));
        assertEq(lending.debt(), owed);
        assertApproxEqAbs(lending.convertToAssets(lending.balanceOf(address(this))), worth - 10 * SHARE, 1);
        vm.startPrank(borrower);
        lending.borrow(1, borrower);
        assertGt(lending.debt(), owed);
        owed = lending.debt();
        assertEq(lending.repay(1), 1);
        vm.stopPrank();
        assertEq(lending.debt(), owed);
    }

    function deploy(SupplyVault.RateModel memory model, uint16 feeShare) external returns (StockLendingVault) {
        return new StockLendingVault(IERC20(address(nvda)), owner, model, feeShare);
    }

    function _lend(uint256 amount) internal {
        nvda.mint(address(this), amount);
        lending.deposit(amount, address(this));
    }

    function _feeShares(uint256 owed, uint256 supply, uint256 assets) internal pure returns (uint256) {
        return owed * (supply + 1) / (assets - owed + 1);
    }
}

contract MarginAccountsLendingTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant TSLA = "TSLA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);
    bytes32 internal constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");
    bytes32 internal constant DOMAIN_TYPEHASH =
        keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");

    StubUsdg internal usdg;
    StubToken internal weth;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    BandDouble internal band;
    MarginDouble internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    StockLendingVault internal lending;
    address internal owner = makeAddr("owner");
    address internal borrower = makeAddr("borrower");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal buyer = makeAddr("buyer");

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        weth = new StubToken(18);
        nvda = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(NVDA, 200e8);
        _quote(SPY, 600e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        _deploy(type(uint256).max);
    }

    function test_TheOwnerSetsAnAssetsLendingVaultOnce() public {
        StockLendingVault fresh = _lendingFor(spy, address(accounts));
        StockLendingVault foreign = _lendingFor(spy, alice);
        StockLendingVault other = _lendingFor(nvda, address(accounts));
        StockLendingVault alien =
            new StockLendingVault(IERC20(address(spy)), alice, SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.prank(alice);
        alien.setDepositor(address(accounts));
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.setLending(SPY, fresh);
        vm.startPrank(owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, TSLA));
        accounts.setLending(TSLA, fresh);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, bytes32("AAPL")));
        accounts.setLending("AAPL", fresh);
        vm.expectRevert(MarginAccounts.InvalidLending.selector);
        accounts.setLending(SPY, other);
        vm.expectRevert(MarginAccounts.InvalidLending.selector);
        accounts.setLending(SPY, foreign);
        vm.expectRevert(MarginAccounts.InvalidLending.selector);
        accounts.setLending(SPY, alien);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.LendingSet(SPY, address(fresh));
        accounts.setLending(SPY, fresh);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.LendingAlreadySet.selector, SPY, address(fresh)));
        accounts.setLending(SPY, fresh);
        vm.stopPrank();
        assertEq(address(accounts.lending(SPY)), address(fresh));
        assertEq(address(accounts.lending(NVDA)), address(lending));
        assertEq(address(accounts.lending(TSLA)), address(0));
        assertEq(spy.allowance(address(accounts), address(fresh)), type(uint256).max);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnknownAsset.selector, bytes32("AAPL")));
        accounts.lending("AAPL");
    }

    function test_APositionLendsItsTokensAndStillCountsThemLessTheRecallHaircut() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Lend(alice, CROSS, NVDA, 4 * SHARE, 4 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 6 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 4 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(usdg)), 0);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(usdg)), 0);
        assertEq(nvda.balanceOf(address(accounts)), 6 * SHARE);
        assertEq(lending.idle(), 4 * SHARE);
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertEq(units, 6 * SHARE);
        assertEq(scale, 1e18);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertEq(equity, 1_960e18);
        assertEq(requirement, 400e18);
        assertEq(accounts.RECALL_HAIRCUT(), 5_00);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.lend(CROSS, address(nvda), SHARE, alice);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(usdg), CROSS));
        accounts.lend(CROSS, address(usdg), SHARE, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NoLending.selector, SPY));
        accounts.lend(SPY, address(spy), SHARE, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.lend(CROSS, address(nvda), 0, alice);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 7 * SHARE)
        );
        accounts.lend(CROSS, address(nvda), 7 * SHARE, alice);
        vm.stopPrank();
        nvda.blockAccount(alice, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        vm.prank(alice);
        accounts.lend(CROSS, address(nvda), SHARE, alice);
        nvda.blockAccount(alice, false);
        _lend(alice, CROSS, 6 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        (equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 1_900e18);
    }

    function test_APositionTakesItsTokensBackAsFarAsTheVaultHoldsThemWithTheFee() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        vm.prank(borrower);
        lending.borrow(3 * SHARE, borrower);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 7 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(weth), CROSS));
        accounts.unlend(CROSS, address(weth), SHARE, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NoLending.selector, SPY));
        accounts.unlend(SPY, address(spy), SHARE, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.unlend(CROSS, address(nvda), 0, alice);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 5 * SHARE)
        );
        accounts.unlend(CROSS, address(nvda), 5 * SHARE, alice);
        vm.expectRevert(
            abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxWithdraw.selector, address(accounts), 2 * SHARE, SHARE)
        );
        accounts.unlend(CROSS, address(nvda), 2 * SHARE, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unlend(alice, CROSS, NVDA, SHARE, SHARE);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
        vm.stopPrank();
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 7 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 3 * SHARE);
        assertEq(nvda.balanceOf(address(accounts)), 7 * SHARE);
        skip(365 days);
        uint256 lent = accounts.lent(alice, CROSS, address(nvda));
        assertGt(lent, 3 * SHARE);
        nvda.mint(borrower, 2 * SHARE);
        vm.prank(borrower);
        lending.repay(type(uint256).max);
        lent = accounts.lent(alice, CROSS, address(nvda));
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), lent, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 7 * SHARE + lent);
        assertGt(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertEq(units * scale / 1e18, nvda.balanceOf(address(accounts)));
    }

    function test_APositionInDebtLendsOnlyIfItStillMeetsItsRequirement() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_550 * USDG);
        _lend(alice, CROSS, 4 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientMargin.selector, int256(350e18), 400e18));
        vm.prank(alice);
        accounts.lend(CROSS, address(nvda), 6 * SHARE, alice);
        vm.prank(owner);
        accounts.setGuardian(owner);
        vm.prank(owner);
        accounts.setBorrowingPaused(true);
        vm.expectRevert(MarginAccounts.BorrowingIsPaused.selector);
        vm.prank(alice);
        accounts.lend(CROSS, address(nvda), SHARE, alice);
    }

    function test_TheCapCountsWhatTheAccountsLent() public {
        _deploy(10 * SHARE);
        _deposit(alice, CROSS, address(nvda), 8 * SHARE);
        _lend(alice, CROSS, 5 * SHARE);
        _deposit(bob, CROSS, address(nvda), 2 * SHARE);
        nvda.mint(bob, 1);
        vm.startPrank(bob);
        nvda.approve(address(accounts), 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 10 * SHARE));
        accounts.deposit(CROSS, address(nvda), 1, bob);
        vm.stopPrank();
    }

    function test_ABurnFromTheVaultFallsOnItsLendersAndABurnFromTheAccountsOnTheirHolders() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        nvda.adminBurn(address(lending), 2 * SHARE);
        lending.sync();
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 2 * SHARE);
        accounts.sync(NVDA);
        (, uint256 scale,) = accounts.holding(NVDA);
        assertEq(scale, 1e18);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 10 * SHARE);
        nvda.adminBurn(address(accounts), 1.6e18);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 5.4e18);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 9e18);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 2 * SHARE);
        uint256 back = lending.maxWithdraw(address(accounts));
        assertEq(back, 2 * SHARE - 1);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), back, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 5.4e18 + back - 1);
    }

    function test_ABurnOfAllTheAccountsHoldBlocksAReturnUntilTheHoldingIsCleared() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        nvda.adminBurn(address(accounts), 6 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 4 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetWrittenOff.selector, NVDA));
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 4 * SHARE, alice);
        vm.prank(address(liquidator));
        accounts.seize(CROSS, address(nvda), SHARE, alice, buyer);
        assertEq(nvda.balanceOf(buyer), SHARE);
        accounts.clear(alice, CROSS, NVDA);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 570e18);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 3 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 3 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
    }

    function test_TheLiquidatorSeizesLentTokensAsFarAsTheVaultHoldsThem() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        vm.prank(borrower);
        lending.borrow(3 * SHARE, borrower);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NotLiquidator.selector, alice));
        vm.prank(alice);
        accounts.seize(CROSS, address(nvda), SHARE, alice, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        vm.prank(address(liquidator));
        accounts.seize(CROSS, address(nvda), 0, alice, buyer);
        vm.prank(address(liquidator));
        accounts.seize(CROSS, address(nvda), 7 * SHARE, alice, buyer);
        assertEq(nvda.balanceOf(buyer), 7 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 3 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 0);
        vm.expectRevert(
            abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxWithdraw.selector, address(accounts), SHARE, 0)
        );
        vm.prank(address(liquidator));
        accounts.seize(CROSS, address(nvda), SHARE, alice, buyer);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(address(liquidator));
        accounts.writeOff(alice, CROSS);
    }

    function test_WhatAPositionCanSellIsItsOwnLoanAsFarAsTheVaultHoldsIt() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, CROSS, address(nvda), 50 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        _lend(bob, CROSS, 50 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.sellable(bob, CROSS, address(nvda)), 50 * SHARE);
        assertEq(liquidator.hourlyAllowance(alice, CROSS, address(nvda)), SHARE);
        vm.prank(borrower);
        lending.borrow(48 * SHARE, borrower);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.sellable(bob, CROSS, address(nvda)), 6 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(usdg)), 0);
    }

    function test_AWriteOffSweepsLentDustWorthLessThanAUsdg() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 1e15);
        _borrow(alice, CROSS, 1_500 * USDG);
        _quote(NVDA, 100e8);
        liquidator.start(alice, CROSS);
        usdg.mint(buyer, 10_000 * USDG);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE - 1e15, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 1e15);
        vm.prank(bob);
        liquidator.writeOff(alice, CROSS);
        assertEq(nvda.balanceOf(bob), 1e15);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.debt(alice, CROSS), 0);
    }

    function test_ALoanWorthNothingAtAWriteOffCountsAgainOnceItIsWorthSomething() public {
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _lend(bob, CROSS, 10 * SHARE);
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 1);
        nvda.adminBurn(address(lending), 5 * SHARE);
        lending.sync();
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        _borrow(alice, CROSS, 1_500 * USDG);
        _quote(NVDA, 100e8);
        liquidator.start(alice, CROSS);
        usdg.mint(buyer, 10_000 * USDG);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE - 1, type(uint256).max, buyer);
        vm.stopPrank();
        liquidator.writeOff(alice, CROSS);
        vm.startPrank(borrower);
        lending.borrow(4.5e18, borrower);
        skip(10 * 365 days);
        nvda.mint(borrower, 20 * SHARE);
        lending.repay(20 * SHARE);
        vm.stopPrank();
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 1);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 1, alice);
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 100);
    }

    function test_ATokenHeldOnlyAsALoanStillGatesThePositionWhenFrozen() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
        _lend(alice, CROSS, 10 * SHARE);
        nvda.pause();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        vm.prank(alice);
        accounts.borrow(CROSS, USDG, alice, alice);
        nvda.unpause();
        nvda.blockAccount(address(lending), true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        vm.prank(alice);
        accounts.borrow(CROSS, USDG, alice, alice);
        nvda.blockAccount(address(lending), false);
        nvda.adminBurn(address(lending), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        vm.prank(alice);
        accounts.borrow(CROSS, USDG, alice, alice);
        lending.sync();
        _borrow(alice, CROSS, USDG);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 9 * SHARE);
    }

    function test_SyncingAnAssetSyncsItsLendingVaultBeforeAnyoneIsJudged() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 10 * SHARE);
        nvda.adminBurn(address(lending), 4 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 10 * SHARE);
        accounts.sync(NVDA);
        assertEq(lending.idle(), 6 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 6 * SHARE);
        accounts.sync(SPY);
        assertEq(spy.balanceOf(address(accounts)), 0);
    }

    function test_TheLiquidatorSyncsTheVaultOfAPositionThatLentAllItHeld() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 10 * SHARE);
        _borrow(alice, CROSS, 1_500 * USDG);
        nvda.adminBurn(address(lending), 5 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 10 * SHARE);
        liquidator.start(alice, CROSS);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 5 * SHARE);
    }

    function test_ALoanInRawUnitsCarriesAMultiplierChangeThrough() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        nvda.updateMultiplier(2e18, block.timestamp);
        assertEq(nvda.uiMultiplier(), 2e18);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 6 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 4 * SHARE);
        assertEq(lending.idle(), 4 * SHARE);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 4 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertEq(units * scale / 1e18, 10 * SHARE);
    }

    function test_TheLiquidatorJudgesALentHoldingLessItsHaircutAndSellsWhatItCanReach() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_500 * USDG);
        _lend(alice, CROSS, 5 * SHARE);
        vm.prank(borrower);
        lending.borrow(4 * SHARE, borrower);
        _quote(NVDA, 180e8);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertEq(equity, 255e18);
        (int256 judged, uint256 required, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(judged, equity);
        assertEq(required, requirement);
        assertTrue(short);
        liquidator.start(alice, CROSS);
        usdg.mint(buyer, 10_000 * USDG);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 7 * SHARE, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(nvda.balanceOf(buyer), 6 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 4 * SHARE);
    }

    function test_TheLiquidationPriceCountsALentHoldingLessItsHaircut() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        uint256 held = accounts.liquidationPrice(alice, CROSS, NVDA, 0);
        _lend(alice, CROSS, 10 * SHARE);
        uint256 lent = accounts.liquidationPrice(alice, CROSS, NVDA, 0);
        assertApproxEqAbs(held, 125e8, uint256(200e8) / 2 ** 13);
        assertApproxEqAbs(lent, uint256(1_000e8) * 10 / 75, uint256(200e8) / 2 ** 13);
    }

    function test_ADepositWithAPermitReadsTheTokensOwnDomain() public {
        Vm.Wallet memory wallet = vm.createWallet("carol");
        usdg.mint(wallet.addr, 100 * USDG);
        (uint8 v, bytes32 r, bytes32 s) = _permit(wallet, 60 * USDG, block.timestamp);
        vm.prank(wallet.addr);
        accounts.depositWithPermit(CROSS, address(usdg), 60 * USDG, wallet.addr, block.timestamp, v, r, s);
        assertEq(accounts.collateral(wallet.addr, CROSS, address(usdg)), 60 * USDG);
        (v, r, s) = _permit(wallet, 40 * USDG, block.timestamp);
        usdg.permit(wallet.addr, address(accounts), 40 * USDG, block.timestamp, v, r, s);
        vm.prank(wallet.addr);
        accounts.depositWithPermit(CROSS, address(usdg), 40 * USDG, wallet.addr, block.timestamp, v, r, s);
        assertEq(accounts.collateral(wallet.addr, CROSS, address(usdg)), 100 * USDG);
    }

    function test_GasOfEachCall() public {
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        _lend(alice, CROSS, 4 * SHARE);
        vm.snapshotGasLastCall("lend");
        vm.prank(borrower);
        lending.borrow(3 * SHARE, borrower);
        vm.snapshotGasLastCall("borrow");
        skip(1 days);
        vm.prank(borrower);
        lending.repay(SHARE);
        vm.snapshotGasLastCall("repay");
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
        vm.snapshotGasLastCall("unlend");
        _borrow(alice, CROSS, 100 * USDG);
        vm.snapshotGasLastCall("loanAgainstALoan");
    }

    function _deploy(uint256 cap) internal {
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = (NVDA, TSLA, SPY);
        engine = new MarginDouble(
            symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD"))
        );
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (cap, type(uint256).max, type(uint256).max);
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
        lending = _lendingFor(nvda, address(accounts));
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        accounts.setLending(NVDA, lending);
        lending.setBorrower(borrower);
        vm.stopPrank();
        vm.prank(borrower);
        nvda.approve(address(lending), type(uint256).max);
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), type(uint256).max);
        vault.deposit(1_000_000 * USDG, owner);
        vm.stopPrank();
    }

    function _lendingFor(StubStockToken token, address depositor) internal returns (StockLendingVault fresh) {
        fresh =
            new StockLendingVault(IERC20(address(token)), owner, SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00);
        vm.prank(owner);
        fresh.setDepositor(depositor);
    }

    function _quote(bytes32 symbol, uint64 low) internal {
        band.setQuote(symbol, BandDouble.Quote(3, 3, low, 50, low, low + low / 100));
    }

    function _deposit(address account, bytes32 position, address token, uint256 amount) internal {
        if (token == address(usdg)) usdg.mint(account, amount);
        else StubToken(token).mint(account, amount);
        vm.startPrank(account);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, account);
        vm.stopPrank();
    }

    function _lend(address account, bytes32 position, uint256 amount) internal {
        vm.prank(account);
        accounts.lend(position, address(nvda), amount, account);
    }

    function _borrow(address account, bytes32 position, uint256 assets) internal {
        vm.prank(account);
        accounts.borrow(position, assets, account, account);
    }

    function _permit(Vm.Wallet memory wallet, uint256 value, uint256 deadline)
        internal
        view
        returns (uint8, bytes32, bytes32)
    {
        bytes32 domain = _domain();
        bytes32 structHash = keccak256(
            abi.encode(PERMIT_TYPEHASH, wallet.addr, address(accounts), value, usdg.nonces(wallet.addr), deadline)
        );
        return vm.sign(wallet, keccak256(abi.encodePacked("\x19\x01", domain, structHash)));
    }

    function _domain() internal view returns (bytes32) {
        (, string memory name, string memory version, uint256 chainId, address verifying,,) = usdg.eip712Domain();
        return
            keccak256(
                abi.encode(DOMAIN_TYPEHASH, keccak256(bytes(name)), keccak256(bytes(version)), chainId, verifying)
            );
    }
}
