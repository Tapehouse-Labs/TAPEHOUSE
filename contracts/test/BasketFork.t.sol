// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, console} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Basket} from "../src/Basket.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {IStockToken, IStockTokenRegistry} from "./conformance/Interfaces.sol";
import {MarginReference} from "./reference/MarginReference.sol";

contract BasketForkTest is Test {
    uint256 internal constant BLOCK = 75_093_578;
    string internal constant VECTORS = "../stylus/contracts/margin/testdata/requirement-vectors.json";
    string internal constant PARAMETERS = "../stylus/contracts/margin/parameters.json";
    bytes32 internal constant CROSS = bytes32(0);
    address internal constant PAUSER = 0xe7BCB188254Bc6eBBfF63014DfED4cD4A024F22A;
    address internal constant BLOCKER = 0x913cA87347391218e5De2C17c5A0AEba8B0b28fD;
    address internal constant ADMIN_BURNER = 0x957B6de6525C63349f7619743Ef1E0ad93cd74D4;
    address internal constant SANCTIONED = 0x910Cbd523D972eb0a6f4cAe4618aD62622b39DbF;

    string internal json;
    string internal parameters;
    string[] internal names;
    uint256[] internal prices;
    BandDouble internal band;
    MarginReference internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    Basket internal basket;
    IERC20 internal usdg;
    address[] internal tokens;
    uint256[] internal units;
    address internal alice = makeAddr("alice");
    address internal bob = makeAddr("bob");

    function setUp() public {
        vm.createSelectFork("robinhood", BLOCK);
        vm.pauseGasMetering();
        json = vm.readFile(VECTORS);
        parameters = vm.readFile(PARAMETERS);
        names = vm.parseJsonStringArray(json, ".symbols");
        prices = vm.parseJsonUintArray(json, ".prices");
        usdg = IERC20(vm.parseJsonAddress(json, ".usdg"));
        band = new BandDouble();
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        for (uint256 i; i < names.length; ++i) {
            bytes32 symbol = bytes32(bytes(names[i]));
            band.setAsset(symbol, BandDouble.Asset(address(0), bytes32(0), bytes32(0), _token(names[i])));
            _quote(symbol, uint64(prices[i]));
        }
        engine = _engine();
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        uint256[] memory caps = new uint256[](names.length);
        for (uint256 i; i < names.length; ++i) {
            caps[i] = type(uint256).max;
        }
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            IERC20(vm.parseJsonAddress(json, ".weth")),
            address(this),
            caps,
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        vault.setBorrower(address(accounts));
        liquidator = new Liquidator(accounts);
        accounts.setLiquidator(address(liquidator));
        deal(address(usdg), address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
        bytes32[] memory symbols = new bytes32[](5);
        units = new uint256[](5);
        for (uint256 i; i < 5; ++i) {
            symbols[i] = bytes32(bytes(names[i]));
            tokens.push(_token(names[i]));
        }
        (units[0], units[1], units[2], units[3], units[4]) = (4e17, 3e17, 3e17, 2e17, 3e17);
        basket = new Basket("Tapehouse Mega Cap Basket", "thMEGA5", IBand(address(band)), symbols, units, address(this));
        accounts.addBasket(basket);
    }

    function test_ABasketOfRealStockTokensMintsAndRedeemsInKind() public {
        _mint(alice, 50e18);
        for (uint256 i; i < 5; ++i) {
            assertEq(IERC20(tokens[i]).balanceOf(address(basket)), units[i] * 50);
            assertEq(IERC20(tokens[i]).balanceOf(alice), 0);
        }
        vm.prank(alice);
        uint256[] memory assets = basket.redeem(20e18, bob, alice);
        for (uint256 i; i < 5; ++i) {
            assertEq(assets[i], units[i] * 20);
            assertEq(IERC20(tokens[i]).balanceOf(bob), units[i] * 20);
            assertEq(IERC20(tokens[i]).balanceOf(address(basket)), units[i] * 30);
        }
    }

    function test_ABasketIsMarginedAsItsRealStockTokensWithTheirCorrelations() public {
        _mint(alice, 1_000e18);
        vm.startPrank(alice);
        basket.approve(address(accounts), 1_000e18);
        accounts.deposit(CROSS, address(basket), 1_000e18, alice);
        vm.stopPrank();
        int256[] memory quantities = new int256[](names.length);
        uint256 alone = 0;
        for (uint256 i; i < 5; ++i) {
            deal(tokens[i], bob, units[i] * 1_000);
            vm.startPrank(bob);
            IERC20(tokens[i]).approve(address(accounts), units[i] * 1_000);
            accounts.deposit(CROSS, tokens[i], units[i] * 1_000, bob);
            vm.stopPrank();
            quantities[i] = int256(units[i] * 1_000);
            int256[] memory one = new int256[](names.length);
            one[i] = quantities[i];
            (uint256 single,,) = engine.currentRequirement(one, prices);
            alone += single;
        }
        (int256 equity, uint256 requirement, uint8 missing, uint8 regime) = accounts.health(alice, CROSS);
        (int256 directEquity, uint256 directRequirement, uint8 directMissing, uint8 directRegime) =
            accounts.health(bob, CROSS);
        (uint256 expected,,) = engine.currentRequirement(quantities, prices);
        assertEq(requirement, expected);
        assertEq(requirement, directRequirement);
        assertEq(equity, directEquity);
        assertEq(missing, directMissing);
        assertEq(missing, 0);
        assertEq(regime, directRegime);
        assertLt(requirement, alone);
        console.log("basket: equity", uint256(equity) / 1e18);
        console.log("basket: requirement", requirement / 1e18, "its tokens one by one", alone / 1e18);
    }

    function test_TheIssuersControlsReachABasketAndItsHolders() public {
        IStockTokenRegistry registry = IStockTokenRegistry(IStockToken(tokens[0]).ACCESS_CONTROLLED_REGISTRY());
        _mint(alice, 10e18);
        vm.startPrank(alice);
        basket.approve(address(accounts), 5e18);
        accounts.deposit(CROSS, address(basket), 5e18, alice);
        accounts.borrow(CROSS, 100e6, alice, alice);
        vm.stopPrank();
        vm.prank(PAUSER);
        registry.pause();
        assertTrue(basket.paused());
        vm.startPrank(alice);
        vm.expectRevert(Basket.ComponentPaused.selector);
        basket.transfer(bob, 1e18);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketFrozen.selector, address(basket)));
        accounts.borrow(CROSS, 1e6, alice, alice);
        vm.stopPrank();
        vm.prank(PAUSER);
        registry.unpause();
        address[] memory blocked = new address[](1);
        blocked[0] = address(basket);
        vm.prank(BLOCKER);
        registry.blockAccounts(blocked);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(IStockToken.Blocked.selector, address(basket)));
        basket.redeem(1e18, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.BasketFrozen.selector, address(basket)));
        accounts.borrow(CROSS, 1e6, alice, alice);
        vm.stopPrank();
        vm.prank(BLOCKER);
        registry.unblockAccounts(blocked);
        assertTrue(basket.isBlocked(SANCTIONED));
        vm.expectRevert(abi.encodeWithSelector(Basket.Blocked.selector, SANCTIONED));
        vm.prank(alice);
        basket.transfer(SANCTIONED, 1e18);
        uint256 before = accounts.inBaskets(alice, CROSS)[0];
        vm.prank(ADMIN_BURNER);
        IStockToken(tokens[0]).adminBurn(address(basket), units[0] * 5);
        assertEq(accounts.inBaskets(alice, CROSS)[0], before / 2);
        assertEq(basket.previewRedeem(basket.balanceOf(alice))[0], units[0] * 5 / 2);
        vm.prank(alice);
        accounts.borrow(CROSS, 1e6, alice, alice);
    }

    function test_AShortPositionsBasketIsSoldAsItsRealStockTokensForRealUsdg() public {
        _mint(alice, 100e18);
        vm.startPrank(alice);
        basket.approve(address(accounts), 100e18);
        accounts.deposit(CROSS, address(basket), 100e18, alice);
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        accounts.borrow(CROSS, (uint256(equity) - requirement) / 1e12, alice, alice);
        vm.stopPrank();
        uint64 low = uint64(prices[0] * 95 / 100);
        _quote(bytes32(bytes(names[0])), low);
        (,, bool short,) = liquidator.shortfall(alice, CROSS);
        assertTrue(short);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.Unwrap(address(liquidator), alice, address(basket), 100e18, _times(units, 100));
        liquidator.start(alice, CROSS);
        assertEq(accounts.collateral(alice, CROSS, tokens[0]), units[0] * 100);
        address buyer = makeAddr("buyer");
        deal(address(usdg), buyer, 1_000_000e6);
        vm.startPrank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
        (uint256 bought, uint256 cost) = liquidator.buy(alice, CROSS, tokens[0], 10e18, type(uint256).max, buyer);
        vm.stopPrank();
        assertEq(bought, 10e18);
        assertEq(cost, (10e18 * (uint256(low) + low / 100) - 1) / 1e20 + 1);
        assertEq(IERC20(tokens[0]).balanceOf(buyer), 10e18);
        assertEq(accounts.collateral(alice, CROSS, tokens[0]), units[0] * 100 - 10e18);
    }

    function _mint(address account, uint256 shares) internal {
        uint256[] memory assets = basket.previewMint(shares);
        vm.startPrank(account);
        for (uint256 i; i < 5; ++i) {
            deal(tokens[i], account, assets[i]);
            IERC20(tokens[i]).approve(address(basket), assets[i]);
        }
        basket.mint(shares, account, assets);
        vm.stopPrank();
    }

    function _times(uint256[] memory values, uint256 factor) internal pure returns (uint256[] memory out) {
        out = new uint256[](values.length);
        for (uint256 i; i < values.length; ++i) {
            out[i] = values[i] * factor;
        }
    }

    function _quote(bytes32 symbol, uint64 low) internal {
        band.setQuote(symbol, BandDouble.Quote(3, 3, low, 50, low, low + low / 100));
    }

    function _token(string memory name) internal view returns (address) {
        return vm.parseJsonAddress(json, string.concat(".tokens.", name));
    }

    function _engine() internal returns (MarginReference) {
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](names.length);
        for (uint256 i; i < names.length; ++i) {
            uint256[] memory depth = vm.parseJsonUintArray(parameters, string.concat(".depth.", names[i], ".initial"));
            assets[i] = MarginReference.Asset(
                bytes32(bytes(names[i])),
                uint32(vm.parseJsonUint(parameters, string.concat(".volatility.", names[i], ".floor"))),
                uint32(vm.parseJsonUint(parameters, string.concat(".volatility.", names[i], ".initial"))),
                uint32(vm.parseJsonUint(parameters, string.concat(".weekendGap.", names[i], ".floor"))),
                uint32(vm.parseJsonUint(parameters, string.concat(".weekendGap.", names[i], ".initial"))),
                uint32(depth[0]),
                uint32(depth[1]),
                vm.parseJsonAddress(json, string.concat(".pools.", names[i], ".pool")),
                _token(names[i])
            );
        }
        uint256 pairs = names.length * (names.length - 1) / 2;
        uint16[] memory floors = new uint16[](pairs);
        uint16[] memory values = new uint16[](pairs);
        uint256 k;
        for (uint256 i; i < names.length; ++i) {
            for (uint256 j = i + 1; j < names.length; ++j) {
                string memory pair = string.concat(".correlation['", names[i], "/", names[j], "']");
                floors[k] = uint16(vm.parseJsonUint(parameters, string.concat(pair, ".floor")));
                values[k++] = uint16(vm.parseJsonUint(parameters, string.concat(pair, ".initial")));
            }
        }
        return new MarginReference(
            assets,
            floors,
            values,
            bytes32(bytes(vm.parseJsonString(parameters, ".market"))),
            address(usdg),
            vm.parseJsonAddress(json, ".weth"),
            vm.parseJsonAddress(json, ".ethUsd"),
            address(band),
            address(this)
        );
    }
}
