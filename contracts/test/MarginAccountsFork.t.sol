// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {IStockToken, IStockTokenRegistry} from "./conformance/Interfaces.sol";
import {MarginReference} from "./reference/MarginReference.sol";

contract MarginAccountsForkTest is Test {
    uint256 internal constant BLOCK = 75_093_578;
    string internal constant VECTORS = "../stylus/contracts/margin/testdata/requirement-vectors.json";
    string internal constant PARAMETERS = "../stylus/contracts/margin/parameters.json";
    bytes32 internal constant CROSS = bytes32(0);
    bytes32 internal constant NVDA = "NVDA";
    bytes32 internal constant SPY = "SPY";
    bytes32 internal constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 internal constant BLOCKER_ROLE = keccak256("BLOCKER_ROLE");
    bytes32 internal constant ADMIN_BURNER_ROLE = keccak256("ADMIN_BURNER_ROLE");
    address internal constant PAUSER = 0xe7BCB188254Bc6eBBfF63014DfED4cD4A024F22A;
    address internal constant BLOCKER = 0x913cA87347391218e5De2C17c5A0AEba8B0b28fD;
    address internal constant ADMIN_BURNER = 0x957B6de6525C63349f7619743Ef1E0ad93cd74D4;

    string internal json;
    string internal parameters;
    string[] internal names;
    uint256[] internal prices;
    BandDouble internal band;
    MarginReference internal engine;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    IERC20 internal usdg;
    address internal alice = makeAddr("alice");

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
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            IERC20(vm.parseJsonAddress(json, ".weth")),
            address(this),
            _uncapped(6),
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        vault.setBorrower(address(accounts));
        deal(address(usdg), address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
    }

    function test_APositionIsMarginedByTheEngineOverTheRealPools() public {
        _deposit(CROSS, "NVDA", 437e18);
        _deposit(CROSS, "TSLA", 279e18);
        _deposit(CROSS, "SPY", 130e18);
        int256[] memory quantities = new int256[](6);
        (quantities[0], quantities[1], quantities[5]) = (437e18, 279e18, 130e18);
        (uint256 expected, uint8 expectedMissing, uint8 expectedRegime) = engine.currentRequirement(quantities, prices);
        (int256 equity, uint256 requirement, uint8 missing, uint8 regime) = accounts.health(alice, CROSS);
        assertEq(requirement, expected);
        assertEq(missing, expectedMissing);
        assertEq(missing, 0);
        assertEq(regime, expectedRegime);
        assertEq(regime, 2);
        assertEq(equity, int256((437e18 * prices[0] + 279e18 * prices[1] + 130e18 * prices[5]) / 1e8));
        uint256 room = (uint256(equity) - requirement) / 1e12;
        vm.startPrank(alice);
        vm.expectPartialRevert(MarginAccounts.InsufficientMargin.selector);
        accounts.borrow(CROSS, room + 1, alice, alice);
        accounts.borrow(CROSS, room, alice, alice);
        vm.stopPrank();
        assertEq(usdg.balanceOf(alice), room);
    }

    function test_IsolatingAnAssetNeverFreesMargin() public {
        _deposit(CROSS, "NVDA", 437e18);
        _deposit(CROSS, "SPY", 130e18);
        (, uint256 together,,) = accounts.health(alice, CROSS);
        vm.startPrank(alice);
        accounts.withdraw(CROSS, _token("NVDA"), 437e18, alice, alice);
        IERC20(_token("NVDA")).approve(address(accounts), 437e18);
        accounts.deposit(NVDA, _token("NVDA"), 437e18, alice);
        vm.stopPrank();
        (, uint256 isolated,,) = accounts.health(alice, NVDA);
        (, uint256 rest,,) = accounts.health(alice, CROSS);
        assertGe(isolated + rest, together);
    }

    function test_TheLiquidationPriceHoldsAgainstTheEngine() public {
        _deposit(CROSS, "NVDA", 437e18);
        _deposit(CROSS, "SPY", 130e18);
        vm.prank(alice);
        accounts.borrow(CROSS, 100_000e6, alice, alice);
        uint256 price = accounts.liquidationPrice(alice, CROSS, NVDA, 0);
        assertGt(price, 0);
        assertLt(price, prices[0]);
        _quote(NVDA, uint64(price));
        (int256 equity, uint256 requirement,,) = accounts.health(alice, CROSS);
        assertGe(equity, int256(requirement));
        _quote(NVDA, uint64(price - prices[0] / 2 ** 14 - 1));
        (equity, requirement,,) = accounts.health(alice, CROSS);
        assertLt(equity, int256(requirement));
    }

    function test_TheIssuersPauseAndBlocklistStopStockTokenTransfersAndNewRisk() public {
        IStockTokenRegistry registry = IStockTokenRegistry(IStockToken(_token("NVDA")).ACCESS_CONTROLLED_REGISTRY());
        assertTrue(registry.hasRole(PAUSER_ROLE, PAUSER));
        assertTrue(registry.hasRole(BLOCKER_ROLE, BLOCKER));
        _deposit(CROSS, "NVDA", 437e18);
        _deposit(SPY, "SPY", 130e18);
        address liquidator = makeAddr("liquidator");
        accounts.setLiquidator(liquidator);
        vm.startPrank(alice);
        accounts.borrow(CROSS, 10_000e6, alice, alice);
        accounts.borrow(SPY, 10_000e6, alice, alice);
        vm.stopPrank();
        vm.prank(PAUSER);
        registry.pause();
        vm.expectRevert(IStockToken.IsPaused.selector);
        vm.prank(liquidator);
        accounts.seize(CROSS, _token("NVDA"), 1e18, alice, liquidator);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.withdraw(CROSS, _token("NVDA"), 1e18, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, SPY));
        accounts.withdraw(SPY, _token("SPY"), 1e18, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.borrow(CROSS, 1e6, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, SPY));
        accounts.borrow(SPY, 1e6, alice, alice);
        usdg.approve(address(accounts), 30_000e6);
        deal(address(usdg), alice, 30_000e6);
        accounts.repay(CROSS, 10_000e6, alice);
        vm.stopPrank();
        vm.prank(PAUSER);
        registry.unpause();
        deal(_token("NVDA"), alice, 1e18);
        vm.prank(alice);
        IERC20(_token("NVDA")).approve(address(accounts), 1e18);
        address[] memory blocked = new address[](1);
        blocked[0] = address(accounts);
        vm.prank(BLOCKER);
        registry.blockAccounts(blocked);
        vm.expectRevert(abi.encodeWithSelector(IStockToken.Blocked.selector, address(accounts)));
        vm.prank(alice);
        accounts.withdraw(CROSS, _token("NVDA"), 1e18, alice, alice);
        vm.startPrank(alice);
        vm.expectRevert(abi.encodeWithSelector(IStockToken.Blocked.selector, address(accounts)));
        accounts.deposit(CROSS, _token("NVDA"), 1e18, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, NVDA));
        accounts.borrow(CROSS, 1e6, alice, alice);
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.AssetFrozen.selector, SPY));
        accounts.borrow(SPY, 1e6, alice, alice);
        accounts.repay(SPY, 10_001e6, alice);
        vm.stopPrank();
        assertEq(accounts.debt(alice, SPY), 0);
        vm.prank(ADMIN_BURNER);
        IStockToken(_token("NVDA")).adminBurn(address(accounts), 37e18);
        accounts.sync(NVDA);
        assertApproxEqAbs(accounts.collateral(alice, CROSS, _token("NVDA")), 400e18, 100);
        blocked[0] = alice;
        address[] memory none = new address[](1);
        none[0] = address(accounts);
        vm.startPrank(BLOCKER);
        registry.unblockAccounts(none);
        registry.blockAccounts(blocked);
        vm.stopPrank();
        vm.expectRevert(abi.encodeWithSelector(MarginAccounts.Blocked.selector, alice));
        vm.prank(alice);
        accounts.withdraw(CROSS, _token("NVDA"), 1e18, alice, makeAddr("fresh"));
    }

    function test_AnAdminBurnFallsOnTheBurnedAssetsHoldersAlone() public {
        IStockTokenRegistry registry = IStockTokenRegistry(IStockToken(_token("NVDA")).ACCESS_CONTROLLED_REGISTRY());
        assertTrue(registry.hasRole(ADMIN_BURNER_ROLE, ADMIN_BURNER));
        address bob = makeAddr("bob");
        _deposit(CROSS, "NVDA", 300e18);
        _deposit(CROSS, "SPY", 100e18);
        deal(_token("NVDA"), bob, 100e18);
        vm.startPrank(bob);
        IERC20(_token("NVDA")).approve(address(accounts), 100e18);
        accounts.deposit(NVDA, _token("NVDA"), 100e18, bob);
        vm.stopPrank();
        vm.prank(PAUSER);
        registry.pause();
        vm.prank(ADMIN_BURNER);
        IStockToken(_token("NVDA")).adminBurn(address(accounts), 100e18);
        vm.expectEmit(address(accounts));
        emit MarginAccounts.AssetWrittenDown(NVDA, 400e18, 300e18);
        accounts.sync(NVDA);
        assertEq(accounts.collateral(alice, CROSS, _token("NVDA")), 225e18);
        assertEq(accounts.collateral(bob, NVDA, _token("NVDA")), 75e18);
        assertEq(accounts.collateral(alice, CROSS, _token("SPY")), 100e18);
    }

    function test_ThePremiumIsPaidInUsdgBesideTheVault() public {
        _deposit(CROSS, "NVDA", 437e18);
        vm.prank(alice);
        accounts.borrow(CROSS, 50_000e6, alice, alice);
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        uint256 debt = accounts.debt(alice, CROSS);
        uint256 premium = accounts.premium(alice, CROSS);
        assertApproxEqRel(premium, uint256(50_000e6) * 5_00 * 48 hours / (10_000 * 365 days), 1e15);
        uint256 idle = usdg.balanceOf(address(vault));
        deal(address(usdg), alice, debt + premium);
        vm.startPrank(alice);
        usdg.approve(address(accounts), debt + premium);
        assertEq(accounts.repay(CROSS, debt + premium, alice), debt + premium);
        vm.stopPrank();
        assertEq(usdg.balanceOf(address(vault)) - idle, debt);
        assertEq(usdg.balanceOf(address(accounts)), premium);
        address backstop = makeAddr("backstop");
        accounts.setBackstop(backstop);
        vm.prank(backstop);
        accounts.claimPremium();
        accounts.withdrawReserve(address(this), premium / 10);
        assertEq(usdg.balanceOf(backstop), premium - premium / 10);
        assertEq(usdg.balanceOf(address(accounts)), 0);
    }

    function _deposit(bytes32 position, string memory name, uint256 amount) internal {
        address token = _token(name);
        deal(token, alice, amount);
        vm.startPrank(alice);
        IERC20(token).approve(address(accounts), amount);
        accounts.deposit(position, token, amount, alice);
        vm.stopPrank();
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

    function _uncapped(uint256 n) internal pure returns (uint256[] memory caps) {
        caps = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            caps[i] = type(uint256).max;
        }
    }
}
