// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Basket} from "../src/Basket.sol";
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

contract MarginAccountsBasketTest is Test {
    uint256 internal constant USDG = 1e6;
    uint256 internal constant SHARE = 1e18;
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant TSLA = "TSLA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant CROSS = bytes32(0);

    StubUsdg internal usdg;
    StubToken internal weth;
    StubStockToken internal nvda;
    StubStockToken internal tsla;
    StubStockToken internal spy;
    StubAggregator internal ethUsd;
    BandDouble internal band;
    MarginDouble internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    Basket internal basket;
    address internal owner = makeAddr("owner");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal buyer = makeAddr("buyer");

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        weth = new StubToken(18);
        nvda = new StubStockToken(1e18);
        tsla = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(TSLA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(tsla)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(NVDA, 3, 200e8, 202e8);
        _quote(TSLA, 3, 300e8, 303e8);
        _quote(SPY, 3, 600e8, 606e8);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        engine = new MarginDouble(_symbols(NVDA, TSLA, SPY), address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        _deployAccounts(_list(type(uint256).max, type(uint256).max, type(uint256).max));
        basket = _basket(_symbols(NVDA, TSLA, SPY), _list(1e18, 5e17, 2e17), owner);
        vm.prank(owner);
        accounts.addBasket(basket);
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), 1_000_000 * USDG);
        vault.deposit(1_000_000 * USDG, owner);
        vm.stopPrank();
        usdg.mint(buyer, 1_000_000 * USDG);
        vm.prank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
    }

    function test_OnlyTheOwnerAddsABasketOfTheAccountsStockTokensInTheEnginesOrder() public {
        assertEq(accounts.MAX_BASKETS(), 8);
        Basket[] memory baskets = accounts.baskets();
        assertEq(baskets.length, 1);
        assertEq(address(baskets[0]), address(basket));
        Basket pair = _basket(_symbols2(NVDA, SPY), _list2(1e18, 1e17), owner);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        accounts.addBasket(pair);
        vm.startPrank(owner);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.BasketAdded(address(pair));
        accounts.addBasket(pair);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketAlreadyAdded.selector, address(pair)));
        accounts.addBasket(pair);
        BandDouble other = new BandDouble();
        other.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        Basket elsewhere = new Basket("B", "B", IBand(address(other)), _symbols1(NVDA), _list1(1e18), owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InvalidBasket.selector, address(elsewhere)));
        accounts.addBasket(elsewhere);
        Basket reversed = _basket(_symbols2(SPY, NVDA), _list2(1e17, 1e18), owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InvalidBasket.selector, address(reversed)));
        accounts.addBasket(reversed);
        StubStockToken aapl = new StubStockToken(1e18);
        band.setAsset("AAPL", BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(aapl)));
        Basket unknown = _basket(_symbols2(NVDA, "AAPL"), _list2(1e18, 1e18), owner);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InvalidBasket.selector, address(unknown)));
        accounts.addBasket(unknown);
        Basket foreign = _basket(_symbols1(NVDA), _list1(1e18), alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InvalidBasket.selector, address(foreign)));
        accounts.addBasket(foreign);
        for (uint256 i = 2; i < 8; ++i) {
            accounts.addBasket(_basket(_symbols1(TSLA), _list1(i), owner));
        }
        Basket ninth = _basket(_symbols1(TSLA), _list1(9), owner);
        vm.expectRevert(MarginAccounts.TooManyBaskets.selector);
        accounts.addBasket(ninth);
        vm.stopPrank();
        assertEq(accounts.baskets().length, 8);
    }

    function test_ABasketSitsInTheCrossPositionAlone() public {
        _fundBasket(alice, 10 * SHARE);
        vm.startPrank(alice);
        basket.approve(address(accounts), 10 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(basket), NVDA));
        accounts.deposit(NVDA, address(basket), SHARE, alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.deposit(CROSS, address(basket), 0, alice);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, bob, alice));
        vm.prank(bob);
        accounts.deposit(CROSS, address(basket), SHARE, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Deposit(alice, alice, CROSS, address(basket), 10 * SHARE);
        vm.prank(alice);
        accounts.deposit(CROSS, address(basket), 10 * SHARE, alice);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 10 * SHARE);
        assertEq(accounts.collateral(alice, NVDA, address(basket)), 0);
        assertEq(accounts.sellable(alice, CROSS, address(basket)), 10 * SHARE);
        assertEq(basket.balanceOf(address(accounts)), 10 * SHARE);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(basket), SHARE));
        accounts.withdraw(NVDA, address(basket), SHARE, alice, alice);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(basket), 10 * SHARE + 1)
        );
        accounts.withdraw(CROSS, address(basket), 10 * SHARE + 1, alice, alice);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Withdraw(alice, alice, CROSS, address(basket), 4 * SHARE, bob);
        accounts.withdraw(CROSS, address(basket), 4 * SHARE, alice, bob);
        accounts.withdraw(CROSS, address(basket), 6 * SHARE, alice, alice);
        vm.stopPrank();
        assertEq(basket.balanceOf(bob), 4 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 0);
        assertEq(accounts.inBaskets(alice, CROSS), _list(0, 0, 0));
        _deposit(alice, NVDA, address(nvda), SHARE);
        assertEq(accounts.collateral(alice, NVDA, address(nvda)), SHARE);
    }

    function test_ABasketIsMarginedAsTheStockTokensItRedeemsFor() public {
        _depositBasket(alice, 10 * SHARE);
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _deposit(bob, CROSS, address(tsla), 5 * SHARE);
        _deposit(bob, CROSS, address(spy), 2 * SHARE);
        assertEq(accounts.inBaskets(alice, CROSS), _list(10 * SHARE, 5 * SHARE, 2 * SHARE));
        assertEq(accounts.inBaskets(alice, NVDA), _list(0, 0, 0));
        assertEq(accounts.inBaskets(bob, CROSS), _list(0, 0, 0));
        (int256 equity, uint256 requirement, uint8 missing, uint8 regime) = accounts.health(alice, CROSS);
        assertEq(equity, 4_700e18);
        assertEq(requirement, 940e18);
        (int256 bobEquity, uint256 bobRequirement, uint8 bobMissing, uint8 bobRegime) = accounts.health(bob, CROSS);
        assertEq(equity, bobEquity);
        assertEq(requirement, bobRequirement);
        assertEq(missing, bobMissing);
        assertEq(regime, bobRegime);
        assertEq(accounts.leverage(alice, CROSS), accounts.leverage(bob, CROSS));
        assertEq(
            accounts.liquidationPrice(alice, CROSS, SPY, 3_000 * USDG),
            accounts.liquidationPrice(bob, CROSS, SPY, 3_000 * USDG)
        );
        _borrow(alice, CROSS, 3_760 * USDG);
        _borrow(bob, CROSS, 3_760 * USDG);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        vm.prank(alice);
        accounts.borrow(CROSS, 1, alice, alice);
        nvda.mint(address(basket), 1);
        _depositBasket(carol, 3);
        assertEq(accounts.inBaskets(carol, CROSS), _list(3, 1, 0));
        Basket pair = _basket(_symbols2(NVDA, SPY), _list2(2e18, 1e18), owner);
        vm.prank(owner);
        accounts.addBasket(pair);
        _mintAndDeposit(pair, alice, SHARE);
        assertEq(accounts.inBaskets(alice, CROSS), _list(12 * SHARE + 1, 5 * SHARE, 3 * SHARE));
    }

    function test_AnAssetSitsInOnePositionThroughABasketToo() public {
        _deposit(alice, NVDA, address(nvda), SHARE);
        _fundBasket(alice, SHARE);
        vm.startPrank(alice);
        basket.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetInOtherPosition.selector, NVDA, NVDA));
        accounts.deposit(CROSS, address(basket), SHARE, alice);
        vm.stopPrank();
        _depositBasket(bob, SHARE);
        nvda.mint(bob, SHARE);
        vm.startPrank(bob);
        nvda.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetInOtherPosition.selector, NVDA, CROSS));
        accounts.deposit(NVDA, address(nvda), SHARE, bob);
        accounts.deposit(CROSS, address(nvda), SHARE, bob);
        vm.stopPrank();
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), SHARE);
        Basket other = _basket(_symbols1(TSLA), _list1(1e18), owner);
        vm.prank(owner);
        accounts.addBasket(other);
        _mintAndDeposit(other, carol, SHARE);
        _deposit(carol, NVDA, address(nvda), SHARE);
        assertEq(accounts.collateral(carol, NVDA, address(nvda)), SHARE);
    }

    function test_TheCapsCountWhatTheAccountsHoldThroughBaskets() public {
        _deployAccounts(_list(15 * SHARE, type(uint256).max, type(uint256).max));
        vm.prank(owner);
        accounts.addBasket(basket);
        _depositBasket(alice, 10 * SHARE);
        _deposit(bob, CROSS, address(nvda), 5 * SHARE);
        nvda.mint(bob, 1);
        vm.startPrank(bob);
        nvda.approve(address(accounts), 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        accounts.deposit(CROSS, address(nvda), 1, bob);
        vm.stopPrank();
        _fundBasket(carol, SHARE);
        vm.startPrank(carol);
        basket.approve(address(accounts), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        accounts.deposit(CROSS, address(basket), SHARE, carol);
        vm.stopPrank();
    }

    function test_AGiftToABasketCountsAgainstTheCapsBeforeAnyNewRisk() public {
        _redeployWithCaps(_list(15 * SHARE, type(uint256).max, type(uint256).max));
        Basket other = _basket(_symbols1(TSLA), _list1(1e18), owner);
        vm.prank(owner);
        accounts.addBasket(other);
        _depositBasket(alice, 10 * SHARE);
        _mintAndDeposit(other, alice, SHARE);
        _mintAndDeposit(other, carol, SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        nvda.mint(address(basket), 100 * SHARE);
        assertEq(accounts.inBaskets(alice, CROSS)[0], 110 * SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        _borrow(alice, CROSS, 15_000 * USDG);
        _borrow(carol, CROSS, 100 * USDG);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        vm.prank(alice);
        accounts.withdraw(CROSS, address(basket), SHARE, alice, alice);
        usdg.mint(alice, 1_000 * USDG);
        vm.startPrank(alice);
        usdg.approve(address(accounts), 1_000 * USDG);
        accounts.repay(CROSS, 1_000 * USDG, alice);
        accounts.withdraw(CROSS, address(basket), 9 * SHARE, alice, alice);
        vm.stopPrank();
        assertEq(accounts.inBaskets(alice, CROSS)[0], 11 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
    }

    function test_ARebalanceCountsAgainstTheCapsBeforeAnyNewRisk() public {
        _redeployWithCaps(_list(15 * SHARE, type(uint256).max, type(uint256).max));
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 1_000 * USDG);
        vm.prank(owner);
        basket.proposeTarget(_list(2e18, 5e17, 1e16));
        skip(7 days);
        nvda.mint(bob, 6 * SHARE);
        vm.startPrank(bob);
        nvda.approve(address(basket), 6 * SHARE);
        basket.rebalance(_list(6 * SHARE, 0, 0), _list(0, 0, 18e17), bob);
        vm.stopPrank();
        assertEq(accounts.inBaskets(alice, CROSS), _list(16 * SHARE, 5 * SHARE, 2e17));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        _borrow(alice, CROSS, 1);
        usdg.mint(alice, 1_000 * USDG);
        vm.startPrank(alice);
        usdg.approve(address(accounts), 1_000 * USDG);
        accounts.repay(CROSS, 1_000 * USDG, alice);
        accounts.withdraw(CROSS, address(basket), SHARE, alice, alice);
        vm.stopPrank();
        _borrow(alice, CROSS, 1);
    }

    function test_SharesSentToTheAccountsWithoutADepositTakeNoRoomUnderACap() public {
        _redeployWithCaps(_list(15 * SHARE, type(uint256).max, type(uint256).max));
        _depositBasket(alice, 5 * SHARE);
        _fundBasket(bob, 10 * SHARE);
        vm.prank(bob);
        basket.transfer(address(accounts), 10 * SHARE);
        assertEq(basket.balanceOf(address(accounts)), 15 * SHARE);
        Basket other = _basket(_symbols1(TSLA), _list1(1e18), owner);
        vm.prank(owner);
        accounts.addBasket(other);
        _mintAndDeposit(other, bob, SHARE);
        nvda.mint(address(other), 100 * SHARE);
        nvda.mint(address(basket), 1);
        assertEq(accounts.inBaskets(alice, CROSS)[0], 5 * SHARE);
        _deposit(carol, CROSS, address(nvda), 10 * SHARE);
        assertEq(accounts.collateral(carol, CROSS, address(nvda)), 10 * SHARE);
        nvda.mint(carol, 1);
        vm.startPrank(carol);
        nvda.approve(address(accounts), 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        accounts.deposit(CROSS, address(nvda), 1, carol);
        vm.stopPrank();
        _borrow(alice, CROSS, 100 * USDG);
        nvda.mint(address(basket), 2);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetCapExceeded.selector, NVDA, 15 * SHARE));
        _borrow(alice, CROSS, 1);
    }

    function test_ABasketBacksALoanAndLeavesOnlyWhileThePositionMeetsItsRequirement() public {
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 3_000 * USDG);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.InsufficientMargin.selector, int256(290e18), 658e18));
        accounts.withdraw(CROSS, address(basket), 3 * SHARE, alice, alice);
        accounts.withdraw(CROSS, address(basket), SHARE, alice, alice);
        vm.stopPrank();
        _quote(SPY, 0, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetHalted.selector, SPY));
        _borrow(alice, CROSS, 1);
        _quote(SPY, 3, 600e8, 606e8);
        band.setCorporateAction(TSLA, BandDouble.CorporateAction(1, uint64(block.timestamp + 1 hours), 1e18, 2e18));
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.CorporateActionPending.selector, TSLA));
        _borrow(alice, CROSS, 1);
        band.setCorporateAction(TSLA, BandDouble.CorporateAction(0, 0, 0, 0));
        engine.set(20_00, 2, 2);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.LiquidityUnknown.selector, TSLA));
        _borrow(alice, CROSS, 1);
        engine.set(20_00, 0, 2);
        _borrow(alice, CROSS, 1);
    }

    function test_ABasketWhoseStockTokensTheIssuerFrozeBacksNoNewRisk() public {
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 100 * USDG);
        spy.pause();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketFrozen.selector, address(basket)));
        _borrow(alice, CROSS, 1);
        spy.unpause();
        nvda.blockAccount(address(accounts), true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketFrozen.selector, address(basket)));
        _borrow(alice, CROSS, 1);
        nvda.blockAccount(address(accounts), false);
        tsla.blockAccount(address(basket), true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketFrozen.selector, address(basket)));
        _borrow(alice, CROSS, 1);
        tsla.blockAccount(address(basket), false);
        spy.blockAccount(alice, true);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        _borrow(alice, CROSS, 1);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        vm.prank(alice);
        accounts.withdraw(CROSS, address(basket), SHARE, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        vm.prank(alice);
        accounts.deposit(CROSS, address(basket), SHARE, alice);
        usdg.mint(alice, 100 * USDG);
        vm.startPrank(alice);
        usdg.approve(address(accounts), 100 * USDG);
        accounts.repay(CROSS, 100 * USDG, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        accounts.withdraw(CROSS, address(basket), SHARE, alice, bob);
        vm.stopPrank();
        assertEq(accounts.debt(alice, CROSS), 0);
        spy.blockAccount(alice, false);
        _borrow(alice, CROSS, 1);
    }

    function test_UnwrappingABasketPutsItsStockTokensInThePosition() public {
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 3_000 * USDG);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unwrap(alice, alice, address(basket), 4 * SHARE, _list(4 * SHARE, 2 * SHARE, 8e17));
        vm.prank(alice);
        assertEq(accounts.unwrap(alice, address(basket), 4 * SHARE), _list(4 * SHARE, 2 * SHARE, 8e17));
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 6 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 4 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(tsla)), 2 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(spy)), 8e17);
        assertEq(accounts.inBaskets(alice, CROSS), _list(6 * SHARE, 3 * SHARE, 12e17));
        (uint256 units,,) = accounts.holding(NVDA);
        assertEq(units, 4 * SHARE);
        assertEq(nvda.balanceOf(address(accounts)), 4 * SHARE);
        (int256 equityAfter, uint256 requirementAfter,,) = accounts.health(alice, CROSS);
        assertEq(equityAfter, equity);
        assertEq(requirementAfter, requirement);
        vm.prank(alice);
        accounts.setAuthorization(bob, true);
        vm.prank(bob);
        accounts.unwrap(alice, address(basket), SHARE);
        vm.prank(address(liquidator));
        accounts.unwrap(alice, address(basket), SHARE);
        vm.startPrank(carol);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, carol, alice));
        accounts.unwrap(alice, address(basket), SHARE);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.UnsupportedToken.selector, address(nvda), CROSS));
        accounts.unwrap(alice, address(nvda), SHARE);
        vm.stopPrank();
        vm.startPrank(alice);
        vm.expectRevert(MarginAccounts.ZeroAmount.selector);
        accounts.unwrap(alice, address(basket), 0);
        vm.expectRevert(
            abi.encodeWithSelector(MarginAccounts.InsufficientCollateral.selector, address(basket), 4 * SHARE + 1)
        );
        accounts.unwrap(alice, address(basket), 4 * SHARE + 1);
        vm.stopPrank();
        _quote(NVDA, 3, 50e8, 51e8);
        (equity, requirement,,) = accounts.health(alice, CROSS);
        assertLt(equity, int256(requirement));
        vm.prank(carol);
        accounts.unwrap(alice, address(basket), 4 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 0);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(accounts.inBaskets(alice, CROSS), _list(0, 0, 0));
        assertEq(basket.totalSupply(), 0);
    }

    function test_ABurnFromTheBasketLowersEveryHoldersPositionAtOnce() public {
        _depositBasket(alice, 10 * SHARE);
        _depositBasket(bob, 10 * SHARE);
        nvda.adminBurn(address(basket), 10 * SHARE);
        assertEq(accounts.inBaskets(alice, CROSS), _list(5 * SHARE, 5 * SHARE, 2 * SHARE));
        (int256 equity,,,) = accounts.health(alice, CROSS);
        assertEq(equity, 3_700e18);
    }

    function test_AnUnwrapCountsABurnFromTheAccountsFirst() public {
        _deposit(bob, CROSS, address(nvda), 10 * SHARE);
        _depositBasket(alice, 10 * SHARE);
        nvda.adminBurn(address(accounts), 5 * SHARE);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.AssetWrittenDown(NVDA, 10 * SHARE, 5 * SHARE);
        vm.prank(alice);
        accounts.unwrap(alice, address(basket), 10 * SHARE);
        assertEq(accounts.collateral(bob, CROSS, address(nvda)), 5 * SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        assertEq(nvda.balanceOf(address(accounts)), 15 * SHARE);
    }

    function test_AnUnwrapWaitsWhileABurnHasLeftTheAccountsNoneOfAStockToken() public {
        _deposit(bob, CROSS, address(spy), SHARE);
        _depositBasket(alice, 10 * SHARE);
        spy.adminBurn(address(accounts), SHARE);
        vm.prank(alice);
        accounts.unwrap(alice, address(basket), 3);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 3);
        assertEq(accounts.collateral(alice, CROSS, address(spy)), 0);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetWrittenOff.selector, SPY));
        vm.prank(alice);
        accounts.unwrap(alice, address(basket), SHARE);
    }

    function test_APositionHoldingABasketIsNotWrittenOff() public {
        _depositBasket(alice, SHARE);
        _borrow(alice, CROSS, 10 * USDG);
        vm.expectRevert(MarginAccounts.PositionNotEmpty.selector);
        vm.prank(address(liquidator));
        accounts.writeOff(alice, CROSS);
    }

    function test_SharesWorthNothingInTheirTokensDoNotBlockAWriteOff() public {
        Basket dust = _basket(_symbols2(NVDA, SPY), _list2(5e17, 4e17), owner);
        GapBackstop backstop = new GapBackstop(accounts, owner, 1_000 * USDG, _list(0, 0, 0));
        vm.startPrank(owner);
        accounts.addBasket(dust);
        accounts.setBackstop(address(backstop));
        vm.stopPrank();
        _mintAndDeposit(dust, bob, SHARE);
        _deposit(carol, CROSS, address(weth), SHARE);
        _borrow(carol, CROSS, 1_000 * USDG);
        ethUsd.setRound(500e8, block.timestamp);
        liquidator.start(carol, CROSS);
        vm.prank(buyer);
        liquidator.buy(carol, CROSS, address(weth), SHARE, type(uint256).max, buyer);
        _mintAndDeposit(dust, carol, 1);
        assertEq(accounts.inBaskets(carol, CROSS), _list(0, 0, 0));
        uint256 owed = accounts.debt(carol, CROSS);
        assertGt(owed, 0);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unwrap(address(liquidator), carol, address(dust), 1, _list2(0, 0));
        (uint256 paid, uint256 written) = backstop.cover(carol, CROSS);
        assertEq(paid, 0);
        assertEq(written, owed);
        assertEq(accounts.collateral(carol, CROSS, address(dust)), 0);
        assertEq(accounts.debt(carol, CROSS), 0);
    }

    function test_ALiquidationBreaksABasketIntoItsStockTokensAndSellsThem() public {
        Basket other = _basket(_symbols1(TSLA), _list1(1e18), owner);
        vm.prank(owner);
        accounts.addBasket(other);
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 3_760 * USDG);
        _quote(NVDA, 3, 190e8, 192e8);
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 840e18);
        assertEq(requirement, 920e18);
        assertTrue(short);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unwrap(
            address(liquidator), alice, address(basket), 10 * SHARE, _list(10 * SHARE, 5 * SHARE, 2 * SHARE)
        );
        liquidator.start(alice, CROSS);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 0);
        assertEq(accounts.collateral(alice, CROSS, address(nvda)), 10 * SHARE);
        (equity, requirement, short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 840e18);
        assertTrue(short);
        vm.prank(buyer);
        (uint256 bought, uint256 cost) = liquidator.buy(alice, CROSS, address(nvda), 2 * SHARE, 1_000 * USDG, buyer);
        assertEq(bought, 2 * SHARE);
        assertEq(cost, 384 * USDG);
        assertEq(nvda.balanceOf(buyer), 2 * SHARE);
    }

    function test_AHealthyPositionHoldingABasketIsNotLiquidated() public {
        _depositBasket(alice, 10 * SHARE);
        _borrow(alice, CROSS, 3_760 * USDG);
        (int256 equity, uint256 requirement, bool short,) = liquidator.shortfall(alice, CROSS);
        assertEq(equity, 940e18);
        assertEq(requirement, 940e18);
        assertFalse(short);
        vm.expectRevert(abi.encodeWithSelector(Liquidator.NotLiquidatable.selector, alice, CROSS));
        liquidator.start(alice, CROSS);
        liquidator.stop(alice, CROSS);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Unauthorized.selector, carol, alice));
        vm.prank(carol);
        accounts.unwrap(alice, address(basket), SHARE);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 10 * SHARE);
    }

    function test_AFrozenBasketWaitsWhileTheRestOfAShortPositionIsSold() public {
        _depositBasket(alice, 10 * SHARE);
        _deposit(alice, CROSS, address(weth), SHARE);
        _borrow(alice, CROSS, 5_000 * USDG);
        _quote(NVDA, 3, 100e8, 101e8);
        spy.pause();
        (,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertTrue(short);
        liquidator.start(alice, CROSS);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 10 * SHARE);
        vm.prank(buyer);
        (uint256 bought,) = liquidator.buy(alice, CROSS, address(weth), SHARE, 10_000 * USDG, buyer);
        assertEq(bought, SHARE);
        spy.unpause();
        vm.prank(buyer);
        liquidator.buy(alice, CROSS, address(nvda), SHARE, 10_000 * USDG, buyer);
        assertEq(accounts.collateral(alice, CROSS, address(basket)), 0);
        assertEq(nvda.balanceOf(buyer), SHARE);
    }

    function test_GasOfEachBasketCall() public {
        _fundBasket(alice, 10 * SHARE);
        vm.startPrank(alice);
        basket.approve(address(accounts), 10 * SHARE);
        accounts.deposit(CROSS, address(basket), 10 * SHARE, alice);
        vm.snapshotGasLastCall("deposit basket");
        accounts.borrow(CROSS, 3_000 * USDG, alice, alice);
        vm.snapshotGasLastCall("borrow against a basket");
        accounts.unwrap(alice, address(basket), 4 * SHARE);
        vm.snapshotGasLastCall("unwrap");
        vm.stopPrank();
        _quote(NVDA, 3, 100e8, 101e8);
        liquidator.start(alice, CROSS);
        vm.snapshotGasLastCall("start, unwrapping");
    }

    function test_GasOfAStockTokenDepositWithABasketOfFiveAdded() public {
        StubStockToken aapl = new StubStockToken(1e18);
        StubStockToken msft = new StubStockToken(1e18);
        band.setAsset("AAPL", BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(aapl)));
        band.setAsset("MSFT", BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(msft)));
        _quote("AAPL", 3, 250e8, 252e8);
        _quote("MSFT", 3, 500e8, 505e8);
        bytes32[] memory symbols = new bytes32[](5);
        (symbols[0], symbols[1], symbols[2], symbols[3], symbols[4]) = (NVDA, TSLA, SPY, "AAPL", "MSFT");
        uint256[] memory units = new uint256[](5);
        (units[0], units[1], units[2], units[3], units[4]) = (4e17, 3e17, 2e17, 3e17, 2e17);
        uint256[] memory caps = new uint256[](5);
        for (uint256 i; i < 5; ++i) {
            caps[i] = type(uint256).max;
        }
        engine = new MarginDouble(symbols, address(band), address(ethUsd));
        _deployAccounts(caps);
        Basket five = _basket(symbols, units, owner);
        uint256 fresh = vm.snapshotState();
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        vm.snapshotGasLastCall("deposit, no basket added");
        vm.revertToState(fresh);
        vm.prank(owner);
        accounts.addBasket(five);
        _mintAndDeposit(five, bob, 10 * SHARE);
        _deposit(alice, CROSS, address(nvda), 10 * SHARE);
        vm.snapshotGasLastCall("deposit, a basket added");
    }

    function _redeployWithCaps(uint256[] memory caps) internal {
        vault = new SupplyVault(IUSDG(address(usdg)), owner, SupplyVault.RateModel(90_00, 0, 0, 0));
        _deployAccounts(caps);
        usdg.mint(owner, 1_000_000 * USDG);
        vm.startPrank(owner);
        usdg.approve(address(vault), 1_000_000 * USDG);
        vault.deposit(1_000_000 * USDG, owner);
        accounts.addBasket(basket);
        vm.stopPrank();
    }

    function _deployAccounts(uint256[] memory caps) internal {
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
        if (vault.borrower() == address(0)) vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        vm.stopPrank();
    }

    function _basket(bytes32[] memory symbols, uint256[] memory units, address basketOwner) internal returns (Basket) {
        return new Basket("Tapehouse Test Basket", "thTEST", IBand(address(band)), symbols, units, basketOwner);
    }

    function _fundBasket(address account, uint256 shares) internal {
        uint256[] memory assets = basket.previewMint(shares);
        StubStockToken[3] memory tokens = [nvda, tsla, spy];
        vm.startPrank(account);
        for (uint256 i; i < 3; ++i) {
            tokens[i].mint(account, assets[i]);
            tokens[i].approve(address(basket), assets[i]);
        }
        basket.mint(shares, account, assets);
        vm.stopPrank();
    }

    function _mintAndDeposit(Basket basket_, address account, uint256 shares) internal {
        uint256[] memory assets = basket_.previewMint(shares);
        (, address[] memory tokens) = basket_.components();
        vm.startPrank(account);
        for (uint256 i; i < tokens.length; ++i) {
            StubStockToken(tokens[i]).mint(account, assets[i]);
            StubStockToken(tokens[i]).approve(address(basket_), assets[i]);
        }
        basket_.mint(shares, account, assets);
        basket_.approve(address(accounts), shares);
        accounts.deposit(CROSS, address(basket_), shares, account);
        vm.stopPrank();
    }

    function _depositBasket(address account, uint256 shares) internal {
        _fundBasket(account, shares);
        vm.startPrank(account);
        basket.approve(address(accounts), shares);
        accounts.deposit(CROSS, address(basket), shares, account);
        vm.stopPrank();
    }

    function _deposit(address account, bytes32 position, address token, uint256 amount) internal {
        StubToken(token).mint(account, amount);
        vm.startPrank(account);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, account);
        vm.stopPrank();
    }

    function _borrow(address account, bytes32 position, uint256 assets) internal {
        vm.prank(account);
        accounts.borrow(position, assets, account, account);
    }

    function _quote(bytes32 symbol, uint8 state, uint64 low, uint128 high) internal {
        band.setQuote(symbol, BandDouble.Quote(state, 3, uint64((uint256(low) + high) / 2), 50, low, high));
    }

    function _list(uint256 a, uint256 b, uint256 c) internal pure returns (uint256[] memory values) {
        values = new uint256[](3);
        (values[0], values[1], values[2]) = (a, b, c);
    }

    function _list2(uint256 a, uint256 b) internal pure returns (uint256[] memory values) {
        values = new uint256[](2);
        (values[0], values[1]) = (a, b);
    }

    function _list1(uint256 a) internal pure returns (uint256[] memory values) {
        values = new uint256[](1);
        values[0] = a;
    }

    function _symbols(bytes32 a, bytes32 b, bytes32 c) internal pure returns (bytes32[] memory values) {
        values = new bytes32[](3);
        (values[0], values[1], values[2]) = (a, b, c);
    }

    function _symbols2(bytes32 a, bytes32 b) internal pure returns (bytes32[] memory values) {
        values = new bytes32[](2);
        (values[0], values[1]) = (a, b);
    }

    function _symbols1(bytes32 a) internal pure returns (bytes32[] memory values) {
        values = new bytes32[](1);
        values[0] = a;
    }
}
