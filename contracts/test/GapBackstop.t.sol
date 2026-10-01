// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {GapBackstop} from "../src/GapBackstop.sol";
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

contract GapBackstopTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    BandDouble internal band;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    GapBackstop internal backstop;
    address internal owner = makeAddr("owner");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal buyer = makeAddr("buyer");

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
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
        StubAggregator ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        MarginDouble engine = new MarginDouble(symbols, address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](2);
        (caps[0], caps[1]) = (type(uint256).max, type(uint256).max);
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            new StubToken(18),
            owner,
            caps,
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        liquidator = new Liquidator(accounts);
        uint256[] memory limits = new uint256[](2);
        (limits[0], limits[1]) = (500 * USDG, 300 * USDG);
        backstop = new GapBackstop(accounts, owner, 1_000 * USDG, limits);
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        accounts.setBackstop(address(backstop));
        vm.stopPrank();
        _lend(owner, 100_000 * USDG);
        usdg.mint(buyer, 1_000_000 * USDG);
        vm.prank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
    }

    function test_TheBackstopDeclaresItsExposureLimitsAndTakesDepositsOnceNamed() public {
        assertEq(address(backstop.accounts()), address(accounts));
        assertEq(backstop.asset(), address(usdg));
        assertEq(backstop.decimals(), 12);
        assertEq(backstop.exposureLimit(CROSS), 1_000 * USDG);
        assertEq(backstop.exposureLimit(NVDA), 500 * USDG);
        assertEq(backstop.exposureLimit(SPY), 300 * USDG);
        uint256[] memory limits = new uint256[](2);
        vm.expectEmit();
        emit GapBackstop.ExposureLimitSet(CROSS, 7, 0);
        vm.expectEmit();
        emit GapBackstop.ExposureLimitSet(NVDA, 0, 0);
        vm.expectEmit();
        emit GapBackstop.ExposureLimitSet(SPY, 0, 0);
        GapBackstop unnamed = new GapBackstop(accounts, owner, 7, limits);
        assertEq(unnamed.maxDeposit(alice), 0);
        vm.expectRevert(GapBackstop.InvalidExposureLimits.selector);
        this.deploy(7, new uint256[](1));
        assertEq(backstop.maxDeposit(alice), type(uint256).max);
        _back(alice, 1_000 * USDG);
        assertEq(backstop.balanceOf(alice), 1_000 * USDG * 1e6);
        assertEq(backstop.held(), 1_000 * USDG);
        assertEq(backstop.totalAssets(), 1_000 * USDG);
        usdg.mint(address(backstop), 5 * USDG);
        assertEq(backstop.totalAssets(), 1_000 * USDG);
    }

    function test_OnlyTheOwnerSetsAnExposureLimitForAKnownPosition() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        backstop.setExposureLimit(NVDA, 1);
        vm.expectRevert(abi.encodeWithSelector(GapBackstop.InvalidPosition.selector, bytes32("TSLA")));
        vm.prank(owner);
        backstop.setExposureLimit("TSLA", 1);
        uint64 fromMs = uint64(block.timestamp + 13 days) * 1000;
        vm.expectEmit(address(backstop));
        emit GapBackstop.ExposureLimitSet(SPY, 42 * USDG, fromMs);
        vm.prank(owner);
        backstop.setExposureLimit(SPY, 42 * USDG);
        vm.prank(owner);
        backstop.setExposureLimit(CROSS, 0);
        assertEq(backstop.exposureLimit(SPY), 300 * USDG);
        assertEq(backstop.exposureLimit(CROSS), 1_000 * USDG);
        (uint128 current, uint128 next, uint64 from) = backstop.exposureLimits(SPY);
        assertEq(current, 300 * USDG);
        assertEq(next, 42 * USDG);
        assertEq(from, fromMs);
        _closeAndReopen();
        assertEq(backstop.exposureLimit(SPY), 300 * USDG);
        skip(13 days);
        _closeAndReopen();
        assertEq(backstop.exposureLimit(SPY), 42 * USDG);
        assertEq(backstop.exposureLimit(CROSS), 0);
        vm.prank(owner);
        backstop.setExposureLimit(SPY, 7);
        (current, next,) = backstop.exposureLimits(SPY);
        assertEq(current, 42 * USDG);
        assertEq(next, 7);
        vm.expectRevert(GapBackstop.OwnershipCannotBeRenounced.selector);
        vm.prank(owner);
        backstop.renounceOwnership();
    }

    function test_ThePremiumClaimedGoesToTheShares() public {
        _back(alice, 1_000 * USDG);
        uint256 premium = _earnPremium();
        uint256 set = premium - premium / 10;
        assertEq(accounts.backstopPremium(), set);
        vm.expectEmit(address(backstop));
        emit GapBackstop.PremiumAdded(set);
        vm.prank(bob);
        assertEq(backstop.claim(), set);
        assertEq(accounts.backstopPremium(), 0);
        assertEq(backstop.claim(), 0);
        assertEq(backstop.held(), 1_000 * USDG + set);
        assertEq(backstop.totalAssets(), 1_000 * USDG + set);
        assertEq(usdg.balanceOf(address(backstop)), backstop.held());
    }

    function test_ACoverRepaysWhatAnEmptiedPositionOwesAndLendersLoseNothing() public {
        _back(alice, 1_000 * USDG);
        _empty(bob, CROSS);
        uint256 owed = accounts.debt(bob, CROSS);
        assertEq(owed, 590 * USDG);
        uint256 lenders = vault.totalAssets();
        vm.expectEmit(address(backstop));
        emit GapBackstop.Covered(bob, CROSS, 0, owed, 0);
        (uint256 paid, uint256 written) = backstop.cover(bob, CROSS);
        assertEq(paid, owed);
        assertEq(written, 0);
        assertEq(accounts.debt(bob, CROSS), 0);
        assertEq(accounts.debtShares(bob, CROSS), 0);
        assertEq(vault.totalAssets(), lenders);
        assertEq(backstop.totalAssets(), 1_000 * USDG - owed);
        assertEq(backstop.covered(CROSS, 0), owed);
        assertEq(backstop.exposureLeft(CROSS), 1_000 * USDG - owed);
        assertEq(usdg.balanceOf(address(backstop)), backstop.held());
    }

    function test_ACoverStopsAtTheExposureLimitAndLendersBearTheRest() public {
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        uint256 lenders = vault.totalAssets();
        (uint256 paid, uint256 written) = backstop.cover(bob, CROSS);
        (paid, written) = backstop.cover(carol, CROSS);
        assertEq(paid, 410 * USDG);
        assertEq(written, 180 * USDG);
        assertEq(vault.totalAssets(), lenders - 180 * USDG);
        assertEq(backstop.exposureLeft(CROSS), 0);
        _empty(alice, CROSS);
        (paid, written) = backstop.cover(alice, CROSS);
        assertEq(paid, 0);
        assertEq(written, 590 * USDG);
        _empty(bob, CROSS);
        _closeAndReopen();
        (uint64 closesMs,,) = accounts.closure();
        assertEq(backstop.exposureLeft(CROSS), 1_000 * USDG);
        (paid, written) = backstop.cover(bob, CROSS);
        assertEq(paid, 590 * USDG);
        assertEq(written, 0);
        assertEq(backstop.covered(CROSS, closesMs), 590 * USDG);
        assertEq(backstop.covered(CROSS, 0), 1_000 * USDG);
    }

    function test_AnIsolatedPositionDrawsOnItsAssetsLimitAndTheBackstopsUsdg() public {
        _back(alice, 400 * USDG);
        _empty(bob, NVDA);
        (uint256 paid, uint256 written) = backstop.cover(bob, NVDA);
        assertEq(paid, 400 * USDG);
        assertEq(written, 190 * USDG);
        assertEq(backstop.covered(NVDA, 0), 400 * USDG);
        assertEq(backstop.covered(CROSS, 0), 0);
        assertEq(backstop.totalAssets(), 0);
        assertEq(backstop.exposureLeft(NVDA), 0);
        assertEq(backstop.exposureLeft(CROSS), 0);
    }

    function test_ACoverForgivesThePremiumAndLeavesTheReserveAndThePremiumSetAside() public {
        _back(alice, 1_000 * USDG);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _borrow(bob, CROSS, 1_600 * USDG);
        _closeAndReopen();
        uint256 premium = accounts.premium(bob, CROSS);
        assertGt(premium, 0);
        _quote(NVDA, 100e8, 101e8);
        liquidator.start(bob, CROSS);
        vm.prank(buyer);
        liquidator.buy(bob, CROSS, address(nvda), 10 * SHARE, type(uint256).max, buyer);
        _quote(NVDA, 200e8, 202e8);
        uint256 owed = accounts.debt(bob, CROSS);
        uint256 reserve = accounts.reserve();
        uint256 set = accounts.backstopPremium();
        vm.expectEmit(address(accounts));
        emit MarginAccounts.WriteOff(bob, CROSS, 0, 0, premium);
        (uint256 paid, uint256 written) = backstop.cover(bob, CROSS);
        assertEq(paid, owed);
        assertEq(written, 0);
        assertEq(accounts.premium(bob, CROSS), 0);
        assertEq(accounts.reserve(), reserve);
        assertEq(accounts.backstopPremium(), set);
    }

    function test_AShortfallABurnLeavesIsCoveredBeforeLenders() public {
        _back(alice, 1_000 * USDG);
        _deposit(bob, NVDA, address(nvda), 10 * SHARE);
        _borrow(bob, NVDA, 400 * USDG);
        nvda.adminBurn(address(accounts), 10 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(bob, NVDA, address(nvda)), 0);
        uint256 lenders = vault.totalAssets();
        (uint256 paid, uint256 written) = backstop.cover(bob, NVDA);
        assertEq(paid, 400 * USDG);
        assertEq(written, 0);
        assertEq(vault.totalAssets(), lenders);
    }

    function test_AFrozenBackstopPaysNothingAndLendersBearTheShortfall() public {
        _back(alice, 1_000 * USDG);
        _empty(bob, CROSS);
        uint256 lenders = vault.totalAssets();
        usdg.freeze(address(backstop));
        (uint256 paid, uint256 written) = backstop.cover(bob, CROSS);
        assertEq(paid, 0);
        assertEq(written, 590 * USDG);
        assertEq(vault.totalAssets(), lenders - 590 * USDG);
        assertEq(backstop.held(), 1_000 * USDG);
    }

    function test_AClosureLongerThanThreeDaysIsCountedOnce() public {
        address dave = makeAddr("dave");
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        _empty(dave, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 96 hours * 1000));
        (uint256 paid,) = backstop.cover(bob, CROSS);
        assertEq(paid, 590 * USDG);
        skip(30 hours);
        (paid,) = backstop.cover(carol, CROSS);
        assertEq(paid, 410 * USDG);
        skip(30 hours);
        (uint64 closesMs,,) = accounts.closure();
        assertGt(closesMs, closes);
        uint256 written;
        (paid, written) = backstop.cover(dave, CROSS);
        assertEq(paid, 0);
        assertEq(written, 590 * USDG);
        assertEq(backstop.closureMs(), closes);
        assertEq(backstop.covered(CROSS, closes), 1_000 * USDG);
    }

    function test_AClosureRecordedAgainAfterThreeDaysIsCountedOnce() public {
        address dave = makeAddr("dave");
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        _empty(dave, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 96 hours * 1000));
        backstop.cover(bob, CROSS);
        backstop.cover(carol, CROSS);
        skip(73 hours);
        (uint256 paid, uint256 written) = backstop.cover(dave, CROSS);
        assertEq(paid, 0);
        assertEq(written, 590 * USDG);
        assertEq(backstop.closureMs(), closes);
    }

    function test_ClosuresLessThanThreeDaysApartShareALimit() public {
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 24 hours * 1000));
        (uint256 paid,) = backstop.cover(bob, CROSS);
        assertEq(paid, 590 * USDG);
        skip(24 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        skip(24 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 96 hours * 1000));
        (paid,) = backstop.cover(carol, CROSS);
        assertEq(paid, 410 * USDG);
        assertEq(backstop.closureMs(), closes);
    }

    function test_AClosureFirstRecordedWhileItRunsBeginsWhenItIsRecorded() public {
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        backstop.cover(bob, CROSS);
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        skip(5 days);
        uint64 seen = uint64(vm.getBlockTimestamp()) * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, seen + 48 hours * 1000));
        (uint256 paid,) = backstop.cover(carol, CROSS);
        assertEq(paid, 590 * USDG);
        assertEq(backstop.closureMs(), seen);
        assertEq(accounts.closureStart(), seen);
    }

    function test_AClosureWithNoKnownReopeningIsCountedOnce() public {
        address dave = makeAddr("dave");
        _back(alice, 10_000 * USDG);
        _empty(bob, CROSS);
        _empty(carol, CROSS);
        _empty(dave, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 3, 0, 0));
        backstop.cover(bob, CROSS);
        backstop.cover(carol, CROSS);
        for (uint256 i; i < 73; ++i) {
            skip(1 hours);
            accounts.accruePremium();
        }
        (uint256 paid, uint256 written) = backstop.cover(dave, CROSS);
        assertEq(paid, 0);
        assertEq(written, 590 * USDG);
        assertEq(backstop.closureMs(), closes);
    }

    function test_TheLastTradingDayBeforeAClosureKeepsTheLastOnesLock() public {
        _back(alice, 1_000 * USDG);
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        uint64 reopens = closes + 48 hours * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, reopens));
        accounts.accruePremium();
        skip(48 hours + 13.5 hours);
        band.setSession(BandDouble.Session(2, 1, 3, 0, reopens + 24 hours * 1000));
        accounts.accruePremium();
        assertEq(backstop.maxRedeem(alice), 0);
        skip(10.5 hours);
        assertEq(backstop.maxRedeem(alice), 0);
    }

    function test_ACoverCountsAgainstTheClosureItRecordsFirst() public {
        _back(alice, 1_000 * USDG);
        _empty(bob, CROSS);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        backstop.cover(bob, CROSS);
        assertEq(backstop.closureMs(), closes);
        assertEq(backstop.covered(CROSS, closes), 590 * USDG);
    }

    function test_AClosureWithNoReopeningRecordedShutsRedemptionsForFourDays() public {
        _back(alice, 1_000 * USDG);
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        skip(50 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        (uint64 closesMs, uint64 reopensMs,) = accounts.closure();
        assertEq(closesMs, closes);
        assertLt(reopensMs, closesMs);
        assertEq(backstop.maxRedeem(alice), 0);
        skip(46 hours - 1);
        assertEq(backstop.maxRedeem(alice), 0);
        skip(1);
        assertGt(backstop.maxRedeem(alice), 0);
    }

    function test_ACoverRevertsUntilThePositionIsEmpty() public {
        _back(alice, 1_000 * USDG);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _borrow(bob, CROSS, 1_600 * USDG);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        backstop.cover(bob, CROSS);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotBackstop.selector, alice));
        vm.prank(alice);
        liquidator.writeOff(bob, CROSS);
        (uint256 paid, uint256 written) = backstop.cover(carol, CROSS);
        assertEq(paid, 0);
        assertEq(written, 0);
    }

    function test_ARedemptionWaitsForACooldownAndKeepsToItsWindow() public {
        _back(alice, 1_000 * USDG);
        uint256 shares = backstop.balanceOf(alice);
        assertEq(backstop.maxRedeem(alice), 0);
        vm.expectRevert(abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxRedeem.selector, alice, 1, 0));
        vm.prank(alice);
        backstop.redeem(1, alice, alice);
        vm.expectEmit(address(backstop));
        emit GapBackstop.CooldownStarted(alice, shares, uint64(block.timestamp + 7 days));
        vm.prank(alice);
        backstop.startCooldown();
        _back(alice, 500 * USDG);
        skip(7 days - 1);
        assertEq(backstop.maxRedeem(alice), 0);
        skip(1);
        assertEq(backstop.maxRedeem(alice), shares);
        assertEq(backstop.maxWithdraw(alice), 1_000 * USDG);
        vm.prank(alice);
        backstop.withdraw(400 * USDG, alice, alice);
        assertEq(usdg.balanceOf(alice), 400 * USDG);
        assertEq(backstop.maxRedeem(alice), shares * 6 / 10);
        vm.prank(alice);
        IERC20(address(backstop)).transfer(bob, shares * 8 / 10);
        assertEq(backstop.maxRedeem(alice), shares * 3 / 10);
        skip(6 days);
        assertEq(backstop.maxRedeem(alice), shares * 3 / 10);
        skip(1);
        assertEq(backstop.maxRedeem(alice), 0);
        assertEq(backstop.maxRedeem(bob), 0);
    }

    function test_NoRedemptionWhileTheMarketIsClosed() public {
        _back(alice, 1_000 * USDG);
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        assertGt(backstop.maxRedeem(alice), 0);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        assertEq(backstop.maxRedeem(alice), 0);
        band.setSession(BandDouble.Session(2, 2, 3, 0, 0));
        assertEq(backstop.maxRedeem(alice), 0);
        band.setSession(BandDouble.Session(2, 1, 3, 0, uint64(block.timestamp) * 1000));
        assertEq(backstop.maxRedeem(alice), 0);
        band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        assertEq(backstop.maxRedeem(alice), 0);
        band.setSession(BandDouble.Session(2, 2, 1, 0, uint64(block.timestamp + 1) * 1000));
        assertGt(backstop.maxRedeem(alice), 0);
    }

    function test_NoRedemptionUntilADayAfterAClosuresReopening() public {
        _back(alice, 1_000 * USDG);
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        assertGt(backstop.maxRedeem(alice), 0);
        accounts.accruePremium();
        assertEq(backstop.maxRedeem(alice), 0);
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        assertEq(backstop.maxRedeem(alice), 0);
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
        assertEq(backstop.maxRedeem(alice), 0);
        skip(1 days - 1);
        assertEq(backstop.maxRedeem(alice), 0);
        skip(1);
        assertEq(backstop.maxRedeem(alice), backstop.balanceOf(alice));
        vm.mockCallRevert(address(accounts), abi.encodeWithSelector(MarginAccounts.closure.selector), "");
        assertEq(backstop.maxRedeem(alice), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(band), abi.encodeWithSelector(BandDouble.session.selector), "");
        assertEq(backstop.maxRedeem(alice), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(accounts), abi.encodeWithSignature("backstop()"), "");
        assertEq(backstop.maxDeposit(alice), 0);
    }

    function test_AFrozenBackstopClosesAndAWipeIsALossToItsShares() public {
        _back(alice, 1_000 * USDG);
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        usdg.pause();
        assertEq(backstop.maxDeposit(alice), 0);
        assertEq(backstop.maxMint(alice), 0);
        assertEq(backstop.maxRedeem(alice), 0);
        usdg.unpause();
        usdg.freeze(address(backstop));
        assertEq(backstop.maxDeposit(alice), 0);
        assertEq(backstop.maxRedeem(alice), 0);
        vm.recordLogs();
        backstop.sync();
        assertEq(vm.getRecordedLogs().length, 0);
        usdg.wipeFrozenAddress(address(backstop));
        usdg.unfreeze(address(backstop));
        assertEq(backstop.maxRedeem(alice), 0);
        vm.expectEmit(address(backstop));
        emit GapBackstop.Sync(0);
        backstop.sync();
        assertEq(backstop.held(), 0);
        assertEq(backstop.totalAssets(), 0);
        assertEq(backstop.maxRedeem(alice), backstop.balanceOf(alice));
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("paused()"), "");
        assertEq(backstop.maxDeposit(alice), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("isFrozen(address)", address(backstop)), "");
        assertEq(backstop.maxDeposit(alice), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("balanceOf(address)", address(backstop)), "");
        assertEq(backstop.maxDeposit(alice), 0);
    }

    function test_GasOfEachCall() public {
        _back(alice, 1_000 * USDG);
        vm.snapshotGasLastCall("deposit");
        _earnPremium();
        backstop.claim();
        vm.snapshotGasLastCall("claim");
        _empty(bob, CROSS);
        backstop.cover(bob, CROSS);
        vm.snapshotGasLastCall("cover");
        vm.prank(alice);
        backstop.startCooldown();
        skip(7 days);
        uint256 shares = backstop.balanceOf(alice);
        vm.prank(alice);
        backstop.redeem(shares, alice, alice);
        vm.snapshotGasLastCall("redeem");
    }

    function deploy(uint256 crossLimit, uint256[] memory limits) external returns (GapBackstop) {
        return new GapBackstop(accounts, owner, crossLimit, limits);
    }

    function _earnPremium() internal returns (uint256 premium) {
        _deposit(carol, SPY, address(spy), 10 * SHARE);
        _borrow(carol, SPY, 1_000 * USDG);
        vm.prank(owner);
        accounts.setPremiumRate(100_00);
        _closeAndReopen();
        premium = accounts.premium(carol, SPY);
        usdg.mint(carol, 1_100 * USDG);
        vm.startPrank(carol);
        usdg.approve(address(accounts), type(uint256).max);
        accounts.repay(SPY, type(uint256).max, carol);
        vm.stopPrank();
    }

    function _closeAndReopen() internal {
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
    }

    function _empty(address account, bytes32 position) internal {
        _quote(NVDA, 200e8, 202e8);
        _deposit(account, position, address(nvda), 10 * SHARE);
        _borrow(account, position, 1_600 * USDG);
        _quote(NVDA, 100e8, 101e8);
        liquidator.start(account, position);
        vm.prank(buyer);
        liquidator.buy(account, position, address(nvda), 10 * SHARE, type(uint256).max, buyer);
        _quote(NVDA, 200e8, 202e8);
    }

    function _back(address depositor, uint256 assets) internal {
        usdg.mint(depositor, assets);
        vm.startPrank(depositor);
        usdg.approve(address(backstop), assets);
        backstop.deposit(assets, depositor);
        vm.stopPrank();
    }

    function _lend(address lender, uint256 assets) internal {
        usdg.mint(lender, assets);
        vm.startPrank(lender);
        usdg.approve(address(vault), assets);
        vault.deposit(assets, lender);
        vm.stopPrank();
    }

    function _quote(bytes32 symbol, uint64 low, uint128 high) internal {
        uint64 mid = uint64((uint256(low) + high) / 2);
        band.setQuote(symbol, BandDouble.Quote(3, 3, mid, 50, low, high));
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
