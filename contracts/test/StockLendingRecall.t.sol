// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {BorrowerDouble} from "./doubles/BorrowerDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract StockLendingRecallVaultTest is Test {
    uint256 internal constant SHARE = 1e18;

    StubStockToken internal nvda;
    StockLendingVault internal lending;
    BorrowerDouble internal borrower;
    address internal owner = makeAddr("owner");
    address internal other = makeAddr("other");

    function setUp() public {
        vm.warp(1_790_000_000);
        nvda = new StubStockToken(1e18);
        lending = new StockLendingVault(IERC20(address(nvda)), owner, SupplyVault.RateModel(80_00, 0, 0, 0), 10_00);
        borrower = new BorrowerDouble(lending, nvda);
        vm.startPrank(owner);
        lending.setDepositor(address(this));
        lending.setBorrower(address(borrower));
        vm.stopPrank();
        nvda.approve(address(lending), type(uint256).max);
        nvda.mint(address(this), 100 * SHARE);
        lending.deposit(100 * SHARE, address(this));
        borrower.borrow(90 * SHARE);
    }

    function test_ARecallIsATicketTheVaultHoldsTokensForAndKeepsFromTheBorrower() public {
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotDepositor.selector, other));
        vm.prank(other);
        lending.recall(SHARE);
        uint64 due = uint64(vm.getBlockTimestamp() + 1 days);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Recall(0, 30 * SHARE, due);
        assertEq(lending.recall(30 * SHARE), 0);
        (uint256 start, uint256 end, uint256 taken, uint64 dueAt) = lending.ticket(0);
        assertEq(start, 0);
        assertEq(end, 30 * SHARE);
        assertEq(taken, 0);
        assertEq(dueAt, due);
        assertEq(lending.tickets(), 1);
        assertEq(lending.requested(), 30 * SHARE);
        assertEq(lending.claimable(0), 10 * SHARE);
        assertEq(lending.borrowable(), 0);
        assertEq(lending.maxWithdraw(address(this)), 0);
        assertEq(lending.utilization(), 0.9e18);
        uint256[] memory first = new uint256[](1);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.OutOfReach.selector, 11 * SHARE, 10 * SHARE));
        lending.reclaim(first, 11 * SHARE, address(this));
        lending.reclaim(first, 4 * SHARE, address(this));
        (,, taken,) = lending.ticket(0);
        assertEq(taken, 4 * SHARE);
        skip(12 hours);
        assertEq(lending.recall(20 * SHARE), 1);
        (start, end, taken, dueAt) = lending.ticket(1);
        assertEq(start, 30 * SHARE);
        assertEq(end, 50 * SHARE);
        assertEq(taken, 0);
        assertEq(dueAt, vm.getBlockTimestamp() + 1 days);
        assertEq(lending.claimable(1), 0);
    }

    function test_TicketsAreHeldForAndTakenInTheOrderTheyWereRecalled() public {
        lending.recall(30 * SHARE);
        lending.recall(20 * SHARE);
        borrower.repay(25 * SHARE);
        assertEq(lending.assigned(), 35 * SHARE);
        assertEq(lending.locked(), 35 * SHARE);
        assertEq(lending.idle(), 0);
        assertEq(lending.claimable(0), 30 * SHARE);
        assertEq(lending.claimable(1), 0);
        uint256[] memory later = new uint256[](1);
        later[0] = 1;
        assertEq(lending.reachable(later), 0);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.OutOfReach.selector, SHARE, 0));
        lending.reclaim(later, SHARE, address(this));
        uint256[] memory both = new uint256[](2);
        both[1] = 1;
        assertEq(lending.reachable(both), 35 * SHARE);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Take(0, 30 * SHARE);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Take(1, 2 * SHARE);
        lending.reclaim(both, 32 * SHARE, address(this));
        assertEq(lending.head(), 1);
        assertEq(lending.served(), 32 * SHARE);
        (,, uint256 taken,) = lending.ticket(0);
        assertEq(taken, 30 * SHARE);
        assertEq(lending.claimable(1), 3 * SHARE);
        assertEq(nvda.balanceOf(address(this)), 32 * SHARE);
        borrower.repay(15 * SHARE);
        assertEq(lending.claimable(1), 18 * SHARE);
        assertEq(lending.borrowable(), 0);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotDepositor.selector, other));
        vm.prank(other);
        lending.reclaim(later, SHARE, other);
    }

    function test_EachTicketHasItsOwnNoticeBeforeABuyIn() public {
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(0);
        lending.recall(30 * SHARE);
        skip(1 days - 1);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(0);
        lending.recall(50 * SHARE);
        skip(1);
        borrower.setToReturn(15 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.BuyInFailed.selector, 5 * SHARE));
        vm.prank(other);
        lending.buyIn(0);
        borrower.setToReturn(type(uint256).max);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(1);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(0, 20 * SHARE);
        vm.prank(other);
        lending.buyIn(0);
        assertEq(lending.claimable(0), 30 * SHARE);
        assertEq(lending.claimable(1), 0);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(0);
        skip(1 days);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(1, 50 * SHARE);
        lending.buyIn(1);
        uint256[] memory first = new uint256[](1);
        lending.reclaim(first, 30 * SHARE, address(this));
        assertEq(lending.claimable(1), 50 * SHARE);
        assertEq(lending.debt(), 20 * SHARE);
    }

    function test_ARecallThatNewLendingMeetsNeedsNoBuyIn() public {
        lending.recall(30 * SHARE);
        nvda.mint(address(this), 20 * SHARE);
        lending.deposit(20 * SHARE, address(this));
        assertEq(lending.claimable(0), 30 * SHARE);
        skip(1 days);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(0);
        nvda.mint(address(this), 5 * SHARE);
        lending.recall(40 * SHARE);
        lending.deposit(5 * SHARE, address(this));
        skip(1 days);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(1, 35 * SHARE);
        lending.buyIn(1);
    }

    function test_ABuyInAsksNoMoreThanTheBorrowerOwes() public {
        borrower.repay(90 * SHARE);
        lending.recall(30 * SHARE);
        assertEq(lending.borrowable(), 70 * SHARE);
        borrower.borrow(70 * SHARE);
        lending.recall(100 * SHARE);
        skip(1 days);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(1, 70 * SHARE);
        lending.buyIn(1);
        assertEq(lending.debt(), 0);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(1);
    }

    function test_ARecallFromAVaultFullyLentIsBoughtIn() public {
        lending.withdraw(10 * SHARE, address(this), address(this));
        assertEq(lending.utilization(), 1e18);
        lending.recall(30 * SHARE);
        assertEq(lending.claimable(0), 0);
        skip(1 days);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(0, 30 * SHARE);
        lending.buyIn(0);
        assertEq(lending.claimable(0), 30 * SHARE);
        assertEq(lending.debt(), 60 * SHARE);
    }

    function test_ABuyInForALaterTicketAsksForAllThatIsDueUpToIt() public {
        for (uint256 i; i < 64; ++i) {
            lending.recall(SHARE / 10);
        }
        lending.recall(30 * SHARE);
        skip(1 days);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(63);
        vm.expectEmit(address(lending));
        emit StockLendingVault.BuyIn(64, 26.4e18);
        lending.buyIn(64);
        vm.expectRevert(StockLendingVault.NotDue.selector);
        lending.buyIn(65);
    }

    function test_TheDepositorGivesUpWhatIsLeftOfTheTicketInTurn() public {
        lending.recall(30 * SHARE);
        lending.recall(5 * SHARE);
        uint256[] memory first = new uint256[](1);
        lending.reclaim(first, 4 * SHARE, address(this));
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotDepositor.selector, other));
        vm.prank(other);
        lending.forfeit(0);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotHead.selector, 0));
        lending.forfeit(1);
        vm.expectEmit(address(lending));
        emit StockLendingVault.Forfeit(0, 26 * SHARE);
        lending.forfeit(0);
        assertEq(lending.head(), 1);
        assertEq(lending.served(), 30 * SHARE);
        assertEq(lending.assigned(), 35 * SHARE);
        assertEq(lending.claimable(0), 0);
        assertEq(lending.claimable(1), 5 * SHARE);
        assertEq(lending.maxWithdraw(address(this)), SHARE);
        (,, uint256 taken,) = lending.ticket(0);
        assertEq(taken, 30 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(StockLendingVault.NotHead.selector, 1));
        lending.forfeit(2);
    }

    function test_TheFeeCountsTheTokensHeldForTicketsAsTheVaults() public {
        vm.prank(owner);
        lending.setRateModel(SupplyVault.RateModel(80_00, 0, 10_00, 10_00));
        lending.recall(30 * SHARE);
        assertEq(lending.locked(), 10 * SHARE);
        assertEq(lending.borrowRate(), 0.15e18);
        skip(365 days);
        assertEq(lending.debt(), 103.5e18);
    }

    function test_ABurnTakesTheFreeTokensFirstThenWhatIsHeldForTheLatestTickets() public {
        lending.recall(4 * SHARE);
        lending.recall(4 * SHARE);
        nvda.adminBurn(address(lending), 3 * SHARE);
        lending.sync();
        assertEq(lending.idle(), 0);
        assertEq(lending.locked(), 7 * SHARE);
        assertEq(lending.assigned(), 7 * SHARE);
        nvda.adminBurn(address(lending), 2 * SHARE);
        lending.sync();
        assertEq(lending.locked(), 5 * SHARE);
        assertEq(lending.assigned(), 5 * SHARE);
        assertEq(lending.claimable(0), 4 * SHARE);
        uint256[] memory both = new uint256[](2);
        both[1] = 1;
        assertEq(lending.reachable(both), 5 * SHARE);
    }

    function test_NoBuyInWhileTheTokenIsPausedTheVaultBlocklistedOrItsTokensBurnt() public {
        lending.recall(30 * SHARE);
        skip(1 days);
        uint256[] memory first = new uint256[](1);
        nvda.pause();
        assertEq(lending.reachable(first), 0);
        vm.expectRevert(StockLendingVault.Unavailable.selector);
        lending.buyIn(0);
        vm.expectRevert(StockLendingVault.Unavailable.selector);
        lending.reclaim(first, SHARE, address(this));
        vm.expectRevert(StubStockToken.IsPaused.selector);
        borrower.repay(SHARE);
        lending.recall(SHARE);
        assertEq(lending.tickets(), 2);
        nvda.unpause();
        nvda.blockAccount(address(lending), true);
        vm.expectRevert(StockLendingVault.Unavailable.selector);
        lending.buyIn(0);
        nvda.blockAccount(address(lending), false);
        nvda.blockAccount(address(borrower), true);
        vm.expectRevert(abi.encodeWithSelector(StubStockToken.Blocked.selector, address(borrower)));
        lending.buyIn(0);
        nvda.blockAccount(address(borrower), false);
        nvda.adminBurn(address(lending), SHARE);
        vm.expectRevert(StockLendingVault.Unavailable.selector);
        lending.buyIn(0);
        lending.sync();
        assertEq(lending.reachable(first), 9 * SHARE);
        lending.buyIn(0);
        assertEq(lending.claimable(0), 30 * SHARE);
    }
}

contract MarginAccountsRecallTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubStockToken internal nvda;
    StubStockToken internal spy;
    BandDouble internal band;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    StockLendingVault internal lending;
    BorrowerDouble internal borrower;
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
        band.setAsset("SPY", BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(200e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = (NVDA, "SPY");
        MarginDouble engine = new MarginDouble(
            symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD"))
        );
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
            0,
            10_00
        );
        liquidator = new Liquidator(accounts);
        lending = new StockLendingVault(IERC20(address(nvda)), owner, SupplyVault.RateModel(80_00, 0, 0, 0), 10_00);
        borrower = new BorrowerDouble(lending, nvda);
        vm.startPrank(owner);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        lending.setDepositor(address(accounts));
        lending.setBorrower(address(borrower));
        accounts.setLending(NVDA, lending);
        vm.stopPrank();
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), type(uint256).max);
        vault.deposit(1_000_000 * USDG, owner);
        vm.stopPrank();
        _lender(alice, 10 * SHARE);
        _lender(bob, 10 * SHARE);
        borrower.borrow(18 * SHARE);
    }

    function test_ARecallTakesTheFreeTokensAtOnceAndQueuesOnlyTheRest() public {
        borrower.repay(5 * SHARE);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Recall(alice, CROSS, NVDA, 0, 3 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 7 * SHARE);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 3 * SHARE);
        assertEq(lending.idle(), 0);
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, bob);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 0);
        assertEq(accounts.claim(bob, CROSS, address(nvda)), 2 * SHARE);
        borrower.repay(10 * SHARE);
        assertEq(lending.idle(), 5 * SHARE);
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 5 * SHARE, bob);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 5 * SHARE);
        assertEq(accounts.recalls(bob, CROSS, address(nvda)).length, 1);
    }

    function test_APositionsRecallIsATicketAndWhatComesBackIsItsBeforeLaterOnes() public {
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), SHARE, alice);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 11 * SHARE)
        );
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 11 * SHARE, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Recall(alice, CROSS, NVDA, 0, 4 * SHARE);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 6 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 2 * SHARE);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(nvda), 5 * SHARE)
        );
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 5 * SHARE, alice);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 4 * SHARE, alice);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 8 * SHARE);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 2);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 2 * SHARE);
        assertEq(accounts.sellable(bob, CROSS, address(nvda)), 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.OutOfReach.selector, NVDA, SHARE, 0));
        vm.prank(bob);
        accounts.unlend(CROSS, address(nvda), SHARE, bob);
        borrower.repay(5 * SHARE);
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 7 * SHARE, bob);
        assertEq(accounts.recalls(bob, CROSS, address(nvda))[0], 2);
        assertEq(accounts.sellable(bob, CROSS, address(nvda)), 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.OutOfReach.selector, NVDA, SHARE, 0));
        vm.prank(bob);
        accounts.unlend(CROSS, address(nvda), SHARE, bob);
        borrower.repay(3 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 10 * SHARE);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 8 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        borrower.repay(7 * SHARE);
        assertEq(accounts.sellable(bob, CROSS, address(nvda)), 7 * SHARE);
    }

    function test_AnyoneSettlesWhatCameBackIntoTheHolding() public {
        vm.expectRevert(MarginAccounts.NothingToSettle.selector);
        accounts.settle(alice, CROSS, address(nvda));
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 4 * SHARE, alice);
        borrower.repay(2 * SHARE);
        vm.prank(buyer);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 4 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 6 * SHARE);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.requested(), lending.served());
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, bob);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), SHARE, alice);
        borrower.repay(11 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 4 * SHARE);
        accounts.settle(bob, CROSS, address(nvda));
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 5 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(usdg), CROSS));
        accounts.settle(alice, CROSS, address(usdg));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NoLending.selector, bytes32("SPY")));
        accounts.settle(alice, CROSS, address(spy));
    }

    function test_SettlingTakesAPositionsTicketsInTurnInOneCall() public {
        vm.startPrank(alice);
        accounts.recall(CROSS, address(nvda), 4 * SHARE, alice);
        accounts.recall(CROSS, address(nvda), 3 * SHARE, alice);
        vm.stopPrank();
        borrower.repay(5 * SHARE);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 7 * SHARE);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.head(), 2);
    }

    function test_ALenderLeavesWithAllItLentAfterABuyIn() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        skip(1 days);
        lending.buyIn(0);
        vm.startPrank(alice);
        accounts.unlend(CROSS, address(nvda), 8 * SHARE, alice);
        accounts.withdraw(CROSS, address(nvda), 10 * SHARE, alice, alice);
        vm.stopPrank();
        assertEq(nvda.balanceOf(alice), 10 * SHARE);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        (uint256 units,,) = accounts.holding(NVDA);
        assertEq(units, 0);
        assertEq(nvda.balanceOf(address(accounts)), 0);
    }

    function test_TheLiquidatorRecallsWhatAPositionThatFallsShortLentAndSellsItOnceItIsBack() public {
        vm.prank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.recall(alice, CROSS, address(nvda));
        _quote(180e8);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Recall(alice, CROSS, NVDA, 0, 8 * SHARE);
        assertEq(liquidator.recall(alice, CROSS, address(nvda)), 10 * SHARE);
        vm.expectRevert(Liquidator.NothingToRecall.selector);
        liquidator.recall(alice, CROSS, address(nvda));
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 2 * SHARE);
        skip(1 days);
        lending.buyIn(0);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 10 * SHARE);
        liquidator.start(alice, CROSS);
        usdg.mint(buyer, 10_000 * USDG);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        (uint256 bought,) = liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(nvda.balanceOf(buyer), bought);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 10 * SHARE - bought);
        assertEq(lending.requested() - lending.served(), 10 * SHARE - bought);
    }

    function test_TheLiquidatorRecallsOnlyWhatThePositionHasNotRecalledYet() public {
        vm.startPrank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        accounts.recall(CROSS, address(nvda), 4 * SHARE, alice);
        vm.stopPrank();
        _quote(180e8);
        assertEq(liquidator.recall(alice, CROSS, address(nvda)), 6 * SHARE);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 8 * SHARE);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 2);
    }

    function test_ARecalledPositionIsWrittenOffOnceItsTokensAreBackAndSold() public {
        vm.prank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        _quote(100e8);
        liquidator.recall(alice, CROSS, address(nvda));
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        liquidator.writeOff(alice, CROSS);
        skip(1 days);
        lending.buyIn(0);
        liquidator.start(alice, CROSS);
        usdg.mint(buyer, 10_000 * USDG);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        liquidator.buy(alice, CROSS, address(nvda), 10 * SHARE, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        liquidator.writeOff(alice, CROSS);
        assertEq(accounts.debt(alice, CROSS), 0);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 0);
        assertEq(lending.requested(), lending.served());
    }

    function test_AfterAWriteOffAnyoneGivesUpTheRecallsOfALoanABurnEmptied() public {
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, bob);
        vm.startPrank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        accounts.recall(CROSS, address(nvda), 5 * SHARE, alice);
        vm.stopPrank();
        borrower.repay(18 * SHARE);
        nvda.adminBurn(address(lending), 18 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        liquidator.writeOff(alice, CROSS);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 1);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.requested(), lending.served());
        vm.expectRevert(MarginAccounts.NothingToSettle.selector);
        accounts.settle(alice, CROSS, address(nvda));
    }

    function test_AWorthlessPositionsRecallsAreGivenUpEachInItsTurn() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 5 * SHARE, alice);
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 5 * SHARE, bob);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 3 * SHARE, alice);
        borrower.repay(18 * SHARE);
        nvda.adminBurn(address(lending), 18 * SHARE);
        accounts.sync(NVDA);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        vm.expectRevert(MarginAccounts.NothingToSettle.selector);
        accounts.settle(bob, CROSS, address(nvda));
        accounts.settle(alice, CROSS, address(nvda));
        uint256[] memory left = accounts.recalls(alice, CROSS, address(nvda));
        assertEq(left.length, 1);
        assertEq(left[0], 2);
        accounts.settle(bob, CROSS, address(nvda));
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.requested(), lending.served());
    }

    function test_ARecallInRawUnitsIsUnmovedByAMultiplierChange() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        nvda.updateMultiplier(3e18, block.timestamp + 12 hours);
        skip(1 days);
        assertEq(nvda.uiMultiplier(), 3e18);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 8 * SHARE);
        lending.buyIn(0);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), 8 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
    }

    function test_ARecallLeftOnALoanABurnEmptiedIsGivenUp() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        borrower.repay(18 * SHARE);
        nvda.adminBurn(address(lending), 10 * SHARE);
        accounts.sync(NVDA);
        assertGt(accounts.lent(alice, CROSS, address(nvda)), 0);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 1);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.requested(), lending.served());
        assertApproxEqAbs(accounts.sellable(bob, CROSS, address(nvda)), accounts.lent(bob, CROSS, address(nvda)), 1);
    }

    function test_ASeizureThatTakesAllThatIsLeftLeavesTheRecallsRestToGiveUp() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        borrower.repay(18 * SHARE);
        nvda.adminBurn(address(lending), 10 * SHARE);
        accounts.sync(NVDA);
        uint256 all = accounts.sellable(alice, CROSS, address(nvda));
        vm.prank(address(liquidator));
        accounts.seize(CROSS, address(nvda), all, alice, buyer);
        assertEq(nvda.balanceOf(buyer), all);
        assertEq(accounts.lent(alice, CROSS, address(nvda)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 1);
        accounts.settle(alice, CROSS, address(nvda));
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 0);
        assertEq(lending.requested(), lending.served());
    }

    function test_NothingComesBackWhileTheTokenIsPausedTheVaultBlocklistedOrItsTokensBurnt() public {
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, bob);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, alice);
        borrower.repay(2 * SHARE);
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 2 * SHARE);
        nvda.pause();
        _nothingComesBack();
        nvda.unpause();
        nvda.blockAccount(address(lending), true);
        _nothingComesBack();
        nvda.blockAccount(address(lending), false);
        nvda.adminBurn(address(lending), SHARE / 2);
        _nothingComesBack();
        lending.sync();
        nvda.blockAccount(address(accounts), true);
        vm.expectRevert(abi.encodeWithSelector(StubStockToken.Blocked.selector, address(accounts)));
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
        nvda.blockAccount(address(accounts), false);
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), SHARE);
    }

    function test_AnOwnerKeepsAtMostThreeRecallsOpenAndTheLiquidatorAFourth() public {
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, bob);
        vm.startPrank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        for (uint256 i; i < 3; ++i) {
            accounts.recall(CROSS, address(nvda), 1e16, alice);
        }
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.TooManyRecalls.selector, NVDA));
        accounts.recall(CROSS, address(nvda), 1e16, alice);
        vm.stopPrank();
        _quote(180e8);
        assertEq(liquidator.recall(alice, CROSS, address(nvda)), 10 * SHARE - 3e16);
        assertEq(accounts.recalls(alice, CROSS, address(nvda)).length, 4);
        vm.expectRevert(Liquidator.NothingToRecall.selector);
        liquidator.recall(alice, CROSS, address(nvda));
    }

    function test_TheLiquidatorLeavesLessThanTheLeastRecallToTheOpenOnes() public {
        vm.startPrank(alice);
        accounts.borrow(CROSS, 1_500 * USDG, alice, alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE - 1e15, alice);
        vm.stopPrank();
        _quote(180e8);
        vm.expectRevert(Liquidator.NothingToRecall.selector);
        liquidator.recall(alice, CROSS, address(nvda));
    }

    function test_ARecallQueuesAtLeastAHundredthOfATokenUnlessItIsThePositionsWhole() public {
        vm.prank(bob);
        accounts.recall(CROSS, address(nvda), 2 * SHARE, bob);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.RecallTooSmall.selector, 1e16 - 1));
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 1e16 - 1, alice);
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 1e16, alice);
        assertEq(accounts.claim(alice, CROSS, address(nvda)), 1e16);
        _lender(carol, 1e15);
        vm.prank(carol);
        accounts.recall(CROSS, address(nvda), 1e15, carol);
        assertEq(accounts.claim(carol, CROSS, address(nvda)), 1e15);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.RecallTooSmall.selector, 1));
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 1, alice);
    }

    function test_ARecallNeedsALendingVaultAndAnAmount() public {
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(usdg), CROSS));
        accounts.recall(CROSS, address(usdg), SHARE, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.NoLending.selector, bytes32("SPY")));
        accounts.recall(CROSS, address(spy), SHARE, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.recall(CROSS, address(nvda), 0, alice);
        vm.stopPrank();
        assertEq(accounts.claim(alice, CROSS, address(usdg)), 0);
        assertEq(accounts.recalls(alice, CROSS, address(usdg)).length, 0);
    }

    function test_GasOfEachCall() public {
        vm.prank(alice);
        accounts.recall(CROSS, address(nvda), 10 * SHARE, alice);
        vm.snapshotGasLastCall("recall");
        skip(1 days);
        lending.buyIn(0);
        vm.snapshotGasLastCall("buyIn");
        accounts.settle(alice, CROSS, address(nvda));
        vm.snapshotGasLastCall("settle");
    }

    function _nothingComesBack() internal {
        assertEq(accounts.sellable(alice, CROSS, address(nvda)), 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.OutOfReach.selector, NVDA, SHARE, 0));
        vm.prank(alice);
        accounts.unlend(CROSS, address(nvda), SHARE, alice);
    }

    function _lender(address account, uint256 amount) internal {
        nvda.mint(account, amount);
        vm.startPrank(account);
        nvda.approve(address(accounts), amount);
        accounts.deposit(CROSS, address(nvda), amount, account);
        accounts.lend(CROSS, address(nvda), amount, account);
        vm.stopPrank();
    }

    function _quote(uint64 low) internal {
        band.setQuote(NVDA, BandDouble.Quote(3, 3, low, 50, low, low + low / 100));
    }
}
