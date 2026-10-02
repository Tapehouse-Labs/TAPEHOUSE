// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {IERC20Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";
import {Basket} from "../src/Basket.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";

contract BasketTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant TSLA = "TSLA";
    bytes32 internal constant SPY = "SPY";

    BandDouble internal band;
    StubStockToken internal nvda;
    StubStockToken internal tsla;
    StubStockToken internal spy;
    Basket internal basket;
    address internal owner = makeAddr("owner");
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");
    address internal carol = makeAddr("carol");
    address internal keeper = makeAddr("keeper");

    function setUp() public {
        vm.warp(1_790_000_000);
        nvda = new StubStockToken(1e18);
        tsla = new StubStockToken(1e18);
        spy = new StubStockToken(1e18);
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setAsset(TSLA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(tsla)));
        band.setAsset(SPY, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(spy)));
        _quote(NVDA, 3, 200e8, 202e8);
        _quote(TSLA, 3, 300e8, 303e8);
        _quote(SPY, 3, 600e8, 606e8);
        basket = _basket(_symbols3(NVDA, TSLA, SPY), _list(1e18, 5e17, 2e17));
    }

    function test_ABasketHoldsTheBandsStockTokensInItsOrder() public {
        (bytes32[] memory symbols, address[] memory tokens) = basket.components();
        assertEq(symbols.length, 3);
        assertEq(symbols[0], NVDA);
        assertEq(symbols[1], TSLA);
        assertEq(symbols[2], SPY);
        assertEq(tokens[0], address(nvda));
        assertEq(tokens[1], address(tsla));
        assertEq(tokens[2], address(spy));
        assertEq(basket.target(), _list(1e18, 5e17, 2e17));
        assertEq(basket.name(), "Tapehouse Test Basket");
        assertEq(basket.symbol(), "thTEST");
        assertEq(basket.decimals(), 18);
        assertEq(basket.owner(), owner);
        assertEq(address(basket.band()), address(band));
        assertEq(basket.NOTICE(), 7 days);
        assertEq(basket.MAX_COMPONENTS(), 8);
        assertEq(basket.totalAssets(), _list(0, 0, 0));
        vm.expectRevert(Basket.InvalidComponents.selector);
        this.deploy(new bytes32[](0), new uint256[](0));
        vm.expectRevert(Basket.InvalidComponents.selector);
        this.deploy(new bytes32[](9), new uint256[](9));
        bytes32[] memory eight = new bytes32[](8);
        uint256[] memory ones = new uint256[](8);
        for (uint256 i; i < 8; ++i) {
            eight[i] = bytes32(uint256(i + 1));
            ones[i] = 1;
            band.setAsset(
                eight[i], BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(new StubStockToken(1e18)))
            );
        }
        assertEq(_basket(eight, ones).target(), ones);
        vm.mockCall(address(tsla), abi.encodeWithSignature("ACCESS_CONTROLLED_REGISTRY()"), abi.encode(address(nvda)));
        Basket shared = _basket(_symbols3(NVDA, TSLA, SPY), _list(1, 1, 1));
        nvda.blockAccount(bob, true);
        tsla.blockAccount(carol, true);
        assertTrue(shared.isBlocked(bob));
        assertFalse(shared.isBlocked(carol));
        vm.expectRevert(Basket.LengthMismatch.selector);
        this.deploy(_symbols3(NVDA, TSLA, SPY), new uint256[](2));
        bytes32[] memory unknown = new bytes32[](1);
        unknown[0] = "AAPL";
        vm.expectRevert(abi.encodeWithSelector(Basket.UnknownAsset.selector, bytes32("AAPL")));
        this.deploy(unknown, new uint256[](1));
        vm.expectRevert(abi.encodeWithSelector(Basket.DuplicateComponent.selector, SPY));
        this.deploy(_symbols3(SPY, TSLA, SPY), _list(1, 1, 1));
        vm.expectRevert(Basket.InvalidTarget.selector);
        this.deploy(_symbols3(NVDA, TSLA, SPY), _list(1, 0, 1));
    }

    function test_TheFirstMintPaysTheTargetRoundedUp() public {
        assertEq(basket.previewMint(1), _list(1, 1, 1));
        uint256 shares = 3e18 + 1;
        uint256[] memory assets = basket.previewMint(shares);
        assertEq(assets, _list(3e18 + 1, 15e17 + 1, 6e17 + 1));
        _fund(alice, assets);
        vm.expectEmit(address(basket));
        emit Basket.Deposit(alice, bob, assets, shares);
        vm.prank(alice);
        assertEq(basket.mint(shares, bob, assets), assets);
        assertEq(basket.balanceOf(bob), shares);
        assertEq(basket.totalAssets(), assets);
        assertEq(nvda.balanceOf(alice), 0);
    }

    function test_MintingPaysEachTokensPartOfWhatTheBasketHoldsRoundedUp() public {
        _mint(alice, 3e18);
        nvda.mint(address(basket), 1);
        uint256[] memory assets = basket.previewMint(1e18);
        assertEq(assets, _list(1e18 + 1, 5e17, 2e17));
        _fund(bob, assets);
        uint256[] memory short = _list(1e18, 5e17, 2e17);
        vm.startPrank(bob);
        vm.expectRevert(abi.encodeWithSelector(Basket.AboveMaximum.selector, address(nvda), 1e18 + 1, 1e18));
        basket.mint(1e18, bob, short);
        vm.expectRevert(Basket.ZeroAmount.selector);
        basket.mint(0, bob, assets);
        vm.expectRevert(Basket.LengthMismatch.selector);
        basket.mint(1e18, bob, new uint256[](2));
        basket.mint(1e18, bob, assets);
        vm.stopPrank();
        assertEq(basket.totalAssets(), _list(4e18 + 2, 2e18, 8e17));
        assertEq(basket.totalSupply(), 4e18);
    }

    function test_RedeemingGivesEachTokensPartRoundedDown() public {
        _mint(alice, 3e18);
        nvda.mint(address(basket), 2);
        uint256[] memory assets = basket.previewRedeem(1e18);
        assertEq(assets, _list(1e18, 5e17, 2e17));
        vm.expectRevert(abi.encodeWithSelector(IERC20Errors.ERC20InsufficientAllowance.selector, bob, 0, 1e18));
        vm.prank(bob);
        basket.redeem(1e18, carol, alice);
        vm.prank(alice);
        basket.approve(bob, 1e18);
        vm.expectEmit(address(basket));
        emit Basket.Withdraw(bob, carol, alice, assets, 1e18);
        vm.prank(bob);
        assertEq(basket.redeem(1e18, carol, alice), assets);
        assertEq(basket.allowance(alice, bob), 0);
        assertEq(nvda.balanceOf(carol), 1e18);
        assertEq(tsla.balanceOf(carol), 5e17);
        assertEq(spy.balanceOf(carol), 2e17);
        vm.startPrank(alice);
        vm.expectRevert(Basket.ZeroAmount.selector);
        basket.redeem(0, alice, alice);
        assertEq(basket.redeem(2e18, alice, alice), _list(2e18 + 2, 1e18, 4e17));
        vm.stopPrank();
        assertEq(basket.previewRedeem(1e18), _list(0, 0, 0));
        assertEq(basket.previewMint(1e18), _list(1e18, 5e17, 2e17));
    }

    function testFuzz_MintsAndRedeemsNeverTakeFromTheOtherHolders(
        uint256 first,
        uint256 gift,
        uint256 burn,
        uint256 shares,
        bool redeeming
    ) public {
        first = bound(first, 1, 1e24);
        _mint(alice, first);
        nvda.mint(address(basket), bound(gift, 0, 1e18));
        tsla.adminBurn(address(basket), bound(burn, 0, tsla.balanceOf(address(basket))));
        uint256[] memory before = basket.totalAssets();
        uint256 supply = basket.totalSupply();
        uint256[] memory paid;
        if (redeeming) {
            shares = bound(shares, 1, first);
            vm.prank(alice);
            paid = basket.redeem(shares, alice, alice);
        } else {
            shares = bound(shares, 1, 1e24);
            paid = basket.previewMint(shares);
            if (_sum(paid) == 0) return;
            _mint(bob, shares);
        }
        uint256[] memory afterwards = basket.totalAssets();
        uint256 supplyAfter = basket.totalSupply();
        for (uint256 i; i < 3; ++i) {
            assertGe(afterwards[i] * supply, before[i] * supplyAfter);
        }
        if (!redeeming) {
            vm.prank(bob);
            uint256[] memory back = basket.redeem(shares, bob, bob);
            for (uint256 i; i < 3; ++i) {
                assertLe(back[i], paid[i]);
            }
        }
    }

    function test_AnAdminBurnFallsOnEveryHolderInProportion() public {
        _mint(alice, 3e18);
        _mint(bob, 1e18);
        nvda.adminBurn(address(basket), 2e18);
        assertEq(basket.previewRedeem(basket.balanceOf(alice)), _list(15e17, 15e17, 6e17));
        assertEq(basket.previewRedeem(basket.balanceOf(bob)), _list(5e17, 5e17, 2e17));
        assertEq(basket.previewMint(1e18), _list(5e17, 5e17, 2e17));
    }

    function test_ABasketWhoseTokensAreAllBurnedTakesNoMintUntilItsSharesAreGone() public {
        _mint(alice, 1e18);
        nvda.adminBurn(address(basket), 1e18);
        tsla.adminBurn(address(basket), 5e17);
        spy.adminBurn(address(basket), 2e17);
        assertEq(basket.previewMint(1e18), _list(0, 0, 0));
        vm.expectRevert(Basket.EmptyBasket.selector);
        vm.prank(bob);
        basket.mint(1e18, bob, _list(0, 0, 0));
        vm.prank(alice);
        assertEq(basket.redeem(1e18, alice, alice), _list(0, 0, 0));
        assertEq(basket.previewMint(1e18), _list(1e18, 5e17, 2e17));
    }

    function test_AShareMovesOnlyAsItsTokensCould() public {
        _mint(alice, 2e18);
        vm.prank(alice);
        basket.transfer(bob, 1e17);
        uint256[] memory assets = basket.previewMint(1e18);
        _fund(carol, assets);
        spy.pause();
        assertTrue(basket.paused());
        vm.startPrank(alice);
        vm.expectRevert(Basket.ComponentPaused.selector);
        basket.transfer(bob, 1);
        vm.expectRevert(Basket.ComponentPaused.selector);
        basket.redeem(1, alice, alice);
        vm.stopPrank();
        vm.expectRevert(Basket.ComponentPaused.selector);
        vm.prank(carol);
        basket.mint(1e18, carol, assets);
        spy.unpause();
        assertFalse(basket.paused());
        tsla.blockAccount(bob, true);
        assertTrue(basket.isBlocked(bob));
        assertFalse(basket.isBlocked(alice));
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, bob));
        vm.prank(bob);
        basket.transfer(alice, 1);
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, bob));
        vm.prank(alice);
        basket.transfer(bob, 1);
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, bob));
        vm.prank(bob);
        basket.redeem(1, bob, bob);
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, bob));
        vm.prank(carol);
        basket.mint(1e18, bob, assets);
        vm.prank(alice);
        basket.approve(bob, 1);
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, bob));
        vm.prank(bob);
        basket.transferFrom(alice, carol, 1);
        tsla.blockAccount(bob, false);
        vm.prank(bob);
        basket.transferFrom(alice, carol, 1);
        assertEq(basket.balanceOf(carol), 1);
    }

    function test_AMultiplierStepPassesThroughInRawUnits() public {
        _mint(alice, 2e18);
        nvda.updateMultiplier(2e18, block.timestamp + 1 hours);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(1, uint64(block.timestamp + 1 hours), 1e18, 2e18));
        assertEq(basket.previewRedeem(1e18), _list(1e18, 5e17, 2e17));
        skip(2 hours);
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(2, uint64(block.timestamp - 1 hours), 1e18, 2e18));
        assertEq(nvda.uiMultiplier(), 2e18);
        vm.prank(alice);
        assertEq(basket.redeem(1e18, alice, alice), _list(1e18, 5e17, 2e17));
        _mint(bob, 1e18);
        assertEq(basket.totalAssets(), _list(2e18, 1e18, 4e17));
    }

    function test_AProposedTargetTakesEffectAfterItsNotice() public {
        uint256[] memory next = _list(1e18, 0, 4e17);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        vm.prank(alice);
        basket.proposeTarget(next);
        vm.startPrank(owner);
        vm.expectRevert(Basket.LengthMismatch.selector);
        basket.proposeTarget(new uint256[](2));
        vm.expectRevert(Basket.InvalidTarget.selector);
        basket.proposeTarget(_list(0, 0, 0));
        vm.expectEmit(address(basket));
        emit Basket.TargetProposed(next, uint64(block.timestamp + 7 days));
        basket.proposeTarget(next);
        vm.stopPrank();
        (uint256[] memory pending, uint64 effectiveAt) = basket.pendingTarget();
        assertEq(pending, next);
        assertEq(effectiveAt, block.timestamp + 7 days);
        assertEq(basket.target(), _list(1e18, 5e17, 2e17));
        skip(7 days - 1);
        assertEq(basket.target(), _list(1e18, 5e17, 2e17));
        skip(1);
        assertEq(basket.target(), next);
        (pending, effectiveAt) = basket.pendingTarget();
        assertEq(pending.length, 0);
        assertEq(effectiveAt, 0);
        assertEq(basket.previewMint(1e18), next);
        vm.prank(owner);
        basket.proposeTarget(_list(1, 1, 1));
        assertEq(basket.target(), next);
        skip(1 days);
        vm.prank(owner);
        basket.proposeTarget(_list(2, 2, 2));
        skip(7 days - 1);
        assertEq(basket.target(), next);
        skip(1);
        assertEq(basket.target(), _list(2, 2, 2));
    }

    function test_TheOwnerCannotRenounce() public {
        vm.expectRevert(Basket.OwnershipCannotBeRenounced.selector);
        vm.prank(owner);
        basket.renounceOwnership();
    }

    function test_ARebalanceMovesTheBasketTowardItsTargetAtTheBandsEdges() public {
        _mint(alice, 10e18);
        _propose(_list(8e17, 5e17, 27e16));
        uint256[] memory assetsIn = _list(0, 0, 7e17);
        uint256[] memory assetsOut = _list(2e18, 0, 0);
        spy.mint(keeper, 7e17);
        vm.startPrank(keeper);
        spy.approve(address(basket), 7e17);
        vm.expectEmit(address(basket));
        emit Basket.Rebalanced(keeper, carol, assetsIn, assetsOut);
        basket.rebalance(assetsIn, assetsOut, carol);
        vm.stopPrank();
        assertEq(basket.totalAssets(), _list(8e18, 5e18, 27e17));
        assertEq(basket.previewRedeem(1e18), basket.target());
        assertEq(nvda.balanceOf(carol), 2e18);
        assertEq(spy.balanceOf(keeper), 0);
    }

    function test_ARebalanceNeverPassesItsTarget() public {
        _mint(alice, 10e18);
        spy.mint(keeper, 34e17 + 1);
        tsla.mint(keeper, 1);
        nvda.mint(keeper, 1);
        vm.startPrank(keeper);
        spy.approve(address(basket), type(uint256).max);
        tsla.approve(address(basket), type(uint256).max);
        nvda.approve(address(basket), type(uint256).max);
        vm.stopPrank();
        vm.prank(owner);
        basket.proposeTarget(_list(0, 5e17, 54e16));
        _expectPastTarget(SPY, _list(0, 0, 1), _list(0, 0, 0));
        skip(7 days);
        _expectPastTarget(SPY, _list(0, 0, 34e17 + 1), _list(0, 0, 0));
        _expectPastTarget(NVDA, _list(0, 0, 0), _list(10e18 + 1, 0, 0));
        _expectPastTarget(NVDA, _list(1, 0, 0), _list(0, 0, 0));
        _expectPastTarget(NVDA, _list(1, 0, 0), _list(1, 0, 0));
        _expectPastTarget(SPY, _list(0, 0, 1), _list(0, 0, 1));
        _expectPastTarget(TSLA, _list(0, 1, 0), _list(0, 0, 0));
        _expectPastTarget(TSLA, _list(0, 0, 0), _list(0, 1, 0));
        vm.prank(keeper);
        basket.rebalance(_list(0, 0, 34e17), _list(10e18, 0, 0), keeper);
        assertEq(basket.totalAssets(), _list(0, 5e18, 54e17));
    }

    function test_ARebalanceNeverGivesMoreThanItGets() public {
        _mint(alice, 10e18);
        _propose(_list(7e17, 5e17, 31e16));
        spy.mint(keeper, 101e16);
        vm.startPrank(keeper);
        spy.approve(address(basket), 101e16);
        vm.expectRevert(
            abi.encodeWithSelector(Basket.ValueLost.selector, (101e16 - 1) * uint256(600e8), 3e18 * uint256(202e8))
        );
        basket.rebalance(_list(0, 0, 101e16 - 1), _list(3e18, 0, 0), keeper);
        basket.rebalance(_list(0, 0, 101e16), _list(3e18, 0, 0), keeper);
        vm.stopPrank();
        assertEq(basket.totalAssets(), _list(7e18, 5e18, 301e16));
    }

    function test_ARebalanceWaitsForEveryBandItTrades() public {
        _mint(alice, 10e18);
        _propose(_list(8e17, 5e17, 27e16));
        spy.mint(keeper, 7e17);
        vm.prank(keeper);
        spy.approve(address(basket), 7e17);
        uint256[] memory assetsIn = _list(0, 0, 7e17);
        uint256[] memory assetsOut = _list(2e18, 0, 0);
        band.setSequencer(address(1), false);
        vm.expectRevert(Basket.SequencerNotSettled.selector);
        vm.prank(keeper);
        basket.rebalance(assetsIn, assetsOut, keeper);
        band.setSequencer(address(0), true);
        _quote(SPY, 0, 0, 0);
        vm.expectRevert(abi.encodeWithSelector(Basket.AssetHalted.selector, SPY));
        vm.prank(keeper);
        basket.rebalance(assetsIn, assetsOut, keeper);
        _quote(SPY, 1, 600e8, 606e8);
        for (uint8 status = 1; status <= 2; ++status) {
            band.setCorporateAction(NVDA, BandDouble.CorporateAction(status, uint64(block.timestamp), 1e18, 2e18));
            vm.expectRevert(abi.encodeWithSelector(Basket.CorporateActionPending.selector, NVDA));
            vm.prank(keeper);
            basket.rebalance(assetsIn, assetsOut, keeper);
        }
        band.setCorporateAction(NVDA, BandDouble.CorporateAction(0, 0, 0, 0));
        _quote(TSLA, 0, 0, 0);
        band.setCorporateAction(TSLA, BandDouble.CorporateAction(2, uint64(block.timestamp), 1e18, 2e18));
        vm.prank(keeper);
        basket.rebalance(assetsIn, assetsOut, keeper);
        assertEq(basket.totalAssets(), _list(8e18, 5e18, 27e17));
    }

    function test_ARebalanceOfNothingReverts() public {
        _mint(alice, 1e18);
        vm.startPrank(keeper);
        vm.expectRevert(Basket.ZeroAmount.selector);
        basket.rebalance(_list(0, 0, 0), _list(0, 0, 0), keeper);
        vm.expectRevert(Basket.LengthMismatch.selector);
        basket.rebalance(new uint256[](2), _list(0, 0, 0), keeper);
        vm.expectRevert(Basket.LengthMismatch.selector);
        basket.rebalance(_list(0, 0, 0), new uint256[](2), keeper);
        vm.stopPrank();
    }

    function testFuzz_ARebalanceNeverLowersTheBasketsWorthAtAnyPriceTheBandAllows(
        uint256[3] memory units,
        uint256[3] memory amounts,
        uint256[3] memory lows,
        uint256[3] memory widths,
        uint256[3] memory picks
    ) public {
        _mint(alice, 10e18);
        uint256[] memory next = new uint256[](3);
        for (uint256 i; i < 3; ++i) {
            next[i] = bound(units[i], 1, 3e18);
        }
        _propose(next);
        bytes32[3] memory symbols = [NVDA, TSLA, SPY];
        StubStockToken[3] memory tokens = [nvda, tsla, spy];
        uint256[] memory held = basket.totalAssets();
        uint256[] memory assetsIn = new uint256[](3);
        uint256[] memory assetsOut = new uint256[](3);
        uint256[] memory prices = new uint256[](3);
        for (uint256 i; i < 3; ++i) {
            uint64 low = uint64(bound(lows[i], 1e8, 1_000e8));
            uint128 high = uint128(low + bound(widths[i], 0, low / 5));
            _quote(symbols[i], 3, low, high);
            prices[i] = bound(picks[i], low, high);
            uint256 goal = next[i] * 10;
            if (goal > held[i]) {
                assetsIn[i] = bound(amounts[i], 0, goal - held[i]);
                tokens[i].mint(keeper, assetsIn[i]);
                vm.prank(keeper);
                tokens[i].approve(address(basket), assetsIn[i]);
            } else {
                assetsOut[i] = bound(amounts[i], 0, held[i] - goal);
            }
        }
        uint256 worth = _worth(held, prices);
        vm.prank(keeper);
        try basket.rebalance(assetsIn, assetsOut, keeper) {
            assertGe(_worth(basket.totalAssets(), prices), worth);
        } catch (bytes memory reason) {
            bytes4 selector = bytes4(reason);
            assertTrue(selector == Basket.ValueLost.selector || selector == Basket.ZeroAmount.selector);
        }
    }

    function test_GasOfEachCall() public {
        _mint(alice, 10e18);
        uint256[] memory assets = basket.previewMint(1e18);
        _fund(bob, assets);
        vm.prank(bob);
        basket.mint(1e18, bob, assets);
        vm.snapshotGasLastCall("mint");
        vm.prank(bob);
        basket.transfer(carol, 1e17);
        vm.snapshotGasLastCall("transfer");
        vm.prank(bob);
        basket.redeem(5e17, bob, bob);
        vm.snapshotGasLastCall("redeem");
        _propose(_list(8e17, 5e17, 27e16));
        spy.mint(keeper, 7e17);
        vm.startPrank(keeper);
        spy.approve(address(basket), 7e17);
        basket.rebalance(_list(0, 0, 7e17), _list(2e18, 0, 0), keeper);
        vm.snapshotGasLastCall("rebalance");
        vm.stopPrank();
    }

    function _expectPastTarget(bytes32 symbol, uint256[] memory assetsIn, uint256[] memory assetsOut) internal {
        vm.expectRevert(abi.encodeWithSelector(Basket.PastTarget.selector, symbol));
        vm.prank(keeper);
        basket.rebalance(assetsIn, assetsOut, keeper);
    }

    function _worth(uint256[] memory assets, uint256[] memory prices) internal pure returns (uint256 worth) {
        for (uint256 i; i < assets.length; ++i) {
            worth += assets[i] * prices[i];
        }
    }

    function _propose(uint256[] memory units) internal {
        vm.prank(owner);
        basket.proposeTarget(units);
        skip(7 days);
    }

    function _mint(address account, uint256 shares) internal {
        uint256[] memory assets = basket.previewMint(shares);
        _fund(account, assets);
        vm.prank(account);
        basket.mint(shares, account, assets);
    }

    function _fund(address account, uint256[] memory assets) internal {
        StubStockToken[3] memory tokens = [nvda, tsla, spy];
        vm.startPrank(account);
        for (uint256 i; i < 3; ++i) {
            tokens[i].mint(account, assets[i]);
            tokens[i].approve(address(basket), assets[i]);
        }
        vm.stopPrank();
    }

    function deploy(bytes32[] memory symbols, uint256[] memory units) external returns (Basket) {
        return _basket(symbols, units);
    }

    function _basket(bytes32[] memory symbols, uint256[] memory units) internal returns (Basket) {
        return new Basket("Tapehouse Test Basket", "thTEST", IBand(address(band)), symbols, units, owner);
    }

    function _quote(bytes32 symbol, uint8 state, uint64 low, uint128 high) internal {
        band.setQuote(symbol, BandDouble.Quote(state, 3, uint64((uint256(low) + high) / 2), 50, low, high));
    }

    function _sum(uint256[] memory values) internal pure returns (uint256 total) {
        for (uint256 i; i < values.length; ++i) {
            total += values[i];
        }
    }

    function _list(uint256 a, uint256 b, uint256 c) internal pure returns (uint256[] memory values) {
        values = new uint256[](3);
        (values[0], values[1], values[2]) = (a, b, c);
    }

    function _symbols3(bytes32 a, bytes32 b, bytes32 c) internal pure returns (bytes32[] memory values) {
        values = new bytes32[](3);
        (values[0], values[1], values[2]) = (a, b, c);
    }
}
