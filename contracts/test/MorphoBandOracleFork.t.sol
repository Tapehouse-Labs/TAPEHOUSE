// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {MorphoBandOracle} from "../src/MorphoBandOracle.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMorpho} from "./conformance/Interfaces.sol";
import {BandDouble} from "./doubles/BandDouble.sol";

contract MorphoBandOracleForkTest is Test {
    uint256 internal constant FRIDAY_BLOCK = 66_476_421;
    string internal constant VECTORS = "../stylus/contracts/band/testdata/quote-vectors.json";
    bytes32 internal constant NVDA = "NVDA";
    uint256 internal constant LLTV = 0.625e18;
    uint256 internal constant WAD = 1e18;
    uint256 internal constant ORACLE_PRICE_SCALE = 1e36;

    uint256[6] internal captures =
        [uint256(1_789_759_836), 1_789_765_236, 1_789_797_638, 1_789_945_210, 1_789_948_930, 1_789_959_611];
    string internal vectors;
    IMorpho internal morpho;
    IERC20 internal nvda;
    IERC20 internal usdg;
    BandDouble internal band;
    MorphoBandOracle internal oracle;
    IMorpho.MarketParams internal params;
    bytes32 internal id;
    address internal owner = makeAddr("owner");
    address internal lender = makeAddr("lender");

    function setUp() public {
        string memory json = vm.readFile("../deployments/4663.json");
        vm.createSelectFork("robinhood", FRIDAY_BLOCK);
        vectors = vm.readFile(VECTORS);
        morpho = IMorpho(vm.parseJsonAddress(json, ".morpho.Blue"));
        nvda = IERC20(vm.parseJsonAddress(json, ".tokens.NVDA"));
        usdg = IERC20(vm.parseJsonAddress(json, ".tokens.USDG"));
        band = new BandDouble();
        band.setAsset(
            NVDA,
            BandDouble.Asset(vm.parseJsonAddress(json, ".chainlink.NVDA_USD"), "NVDA---24_7", bytes32(0), address(nvda))
        );
        _replay(0);
        oracle = new MorphoBandOracle(IBand(address(band)), NVDA, IERC20Metadata(address(usdg)), owner);
        params = IMorpho.MarketParams(
            address(usdg), address(nvda), address(oracle), vm.parseJsonAddress(json, ".morpho.AdaptiveCurveIrm"), LLTV
        );
        morpho.createMarket(params);
        id = keccak256(abi.encode(params));
        deal(address(usdg), lender, 100_000e6);
        vm.startPrank(lender);
        usdg.approve(address(morpho), 100_000e6);
        morpho.supply(params, 100_000e6, 0, lender, "");
        vm.stopPrank();
    }

    function test_AMarketOnTheOracleStaysSafeAcrossARealWeekend() public {
        address[6] memory borrowers;
        for (uint256 i; i < 6; ++i) {
            (string memory name, uint64 low) = _replay(i);
            uint256 price = oracle.price();
            assertEq(price, uint256(low) * 1e16, name);
            morpho.accrueInterest(params);
            uint256 highest;
            uint256 liquidatable;
            for (uint256 j; j < i; ++j) {
                uint256 ltv = _assertSafe(borrowers[j], price);
                if (ltv > highest) highest = ltv;
                if (ltv > LLTV) ++liquidatable;
            }
            borrowers[i] = makeAddr(name);
            emit log_named_decimal_uint(
                string.concat(name, ", USDG lent against 10 NVDA"), _borrowTheLimit(borrowers[i], 10e18), 6
            );
            emit log_named_decimal_uint(string.concat(name, ", highest earlier loan-to-value, %"), highest / 1e12, 4);
            emit log_named_uint(string.concat(name, ", earlier loans past the 62.5% LLTV"), liquidatable);
            assertEq(liquidatable, i == 2 || i == 3 ? 1 : 0, name);
        }
    }

    function test_AHaltStopsBorrowingAndLiquidationWhileRepaymentsGoOn() public {
        _replay(2);
        address alice = makeAddr("alice");
        _borrowTheLimit(alice, 10e18);
        band.setQuote(NVDA, BandDouble.Quote(0, 0, 0, 0, 0, 0));
        uint64 issuedAt = uint64(vm.getBlockTimestamp());
        band.setHalt(NVDA, BandDouble.Halt(true, issuedAt + 15 minutes, issuedAt, false));
        _assertOnlyRepaymentsAndSuppliesGoOn(alice);

        skip(16 minutes);
        band.setQuote(NVDA, BandDouble.Quote(1, 1, 16_663_374_488, 55, 16_571_725_928, 16_755_023_048));
        (,, uint128 collateral) = morpho.position(id, alice);
        (uint256 seized, uint256 owed) = _liquidateAll(alice);
        assertEq(seized, owed * _incentive() / WAD * ORACLE_PRICE_SCALE / (uint256(16_571_725_928) * 1e16));
        (, uint128 left, uint128 kept) = morpho.position(id, alice);
        assertEq(left, 0);
        assertEq(kept, collateral - seized);
    }

    function test_ARepointedOracleKeepsTheMarketAndItsPositions() public {
        address alice = makeAddr("alice");
        uint256 debt = _borrowTheLimit(alice, 10e18);
        (, uint128 shares, uint128 collateral) = morpho.position(id, alice);
        BandDouble next = new BandDouble();
        next.setAsset(NVDA, BandDouble.Asset(address(0), "NVDA---24_7", bytes32(0), address(nvda)));
        next.setQuote(NVDA, BandDouble.Quote(3, 2, 22_282_370_662, 30, 22_215_523_550, 22_349_217_773));
        vm.prank(owner);
        oracle.setBand(IBand(address(next)));
        assertEq(morpho.idToMarketParams(id).oracle, address(oracle));
        (, uint128 sharesAfter, uint128 collateralAfter) = morpho.position(id, alice);
        assertEq(sharesAfter, shares);
        assertEq(collateralAfter, collateral);
        uint256 limit = uint256(collateral) * oracle.price() / ORACLE_PRICE_SCALE * LLTV / WAD;
        assertEq(oracle.price(), uint256(22_215_523_550) * 1e16);
        vm.startPrank(alice);
        vm.expectRevert(bytes("insufficient collateral"));
        morpho.borrow(params, limit - debt + 1, 0, alice, alice);
        morpho.borrow(params, limit - debt - 1, 0, alice, alice);
        vm.stopPrank();
    }

    function _replay(uint256 i) internal returns (string memory name, uint64 low) {
        string memory q = string.concat(".quotes[", vm.toString(i), "]");
        name = vm.parseJsonString(vectors, string.concat(q, ".name"));
        assertTrue(vm.contains(name, "NVDA: "), name);
        string memory state = vm.parseJsonString(vectors, string.concat(q, ".expected.state"));
        low = uint64(_field(q, "low"));
        vm.warp(captures[i]);
        band.setQuote(
            NVDA,
            BandDouble.Quote(
                keccak256(bytes(state)) == keccak256("OPEN")
                    ? 3
                    : keccak256(bytes(state)) == keccak256("CLOSED") ? 2 : 1,
                uint8(_field(q, "live")),
                uint64(_field(q, "mid")),
                uint64(_field(q, "half_bps")),
                low,
                uint128(_field(q, "high"))
            )
        );
    }

    function _field(string memory q, string memory key) internal view returns (uint256) {
        return vm.parseUint(vm.parseJsonString(vectors, string.concat(q, ".expected.", key)));
    }

    function _borrowTheLimit(address borrower, uint256 collateral) internal returns (uint256 limit) {
        deal(address(nvda), borrower, collateral);
        limit = collateral * oracle.price() / ORACLE_PRICE_SCALE * LLTV / WAD;
        vm.startPrank(borrower);
        nvda.approve(address(morpho), collateral);
        morpho.supplyCollateral(params, collateral, borrower, "");
        vm.expectRevert(bytes("insufficient collateral"));
        morpho.borrow(params, limit + 1, 0, borrower, borrower);
        morpho.borrow(params, limit - 1, 0, borrower, borrower);
        vm.stopPrank();
        return limit - 1;
    }

    function _assertOnlyRepaymentsAndSuppliesGoOn(address borrower) internal {
        bytes memory noAnswer = abi.encodeWithSelector(MorphoBandOracle.NoAnswer.selector, NVDA);
        deal(address(usdg), borrower, 100e6);
        deal(address(nvda), borrower, 1e18);
        vm.startPrank(borrower);
        vm.expectRevert(noAnswer);
        morpho.borrow(params, 1, 0, borrower, borrower);
        vm.expectRevert(noAnswer);
        morpho.withdrawCollateral(params, 1, borrower, borrower);
        usdg.approve(address(morpho), 100e6);
        morpho.repay(params, 100e6, 0, borrower, "");
        nvda.approve(address(morpho), 1e18);
        morpho.supplyCollateral(params, 1e18, borrower, "");
        vm.stopPrank();
        vm.expectRevert(noAnswer);
        morpho.liquidate(params, borrower, 1e18, 0, "");
        vm.prank(lender);
        morpho.withdraw(params, 1_000e6, 0, lender, lender);
    }

    function _liquidateAll(address borrower) internal returns (uint256 seized, uint256 owed) {
        (, uint128 shares,) = morpho.position(id, borrower);
        morpho.accrueInterest(params);
        (,, uint128 borrowed, uint128 borrowShares,,) = morpho.market(id);
        owed = uint256(shares) * (borrowed + 1) / (borrowShares + 1e6);
        uint256 repaid;
        (seized, repaid) = _liquidate(borrower, 0, shares);
        assertEq(repaid, owed + 1);
        assertEq(nvda.balanceOf(makeAddr("liquidator")), seized);
    }

    function _liquidate(address borrower, uint256 seize, uint256 shares) internal returns (uint256, uint256) {
        address liquidator = makeAddr("liquidator");
        deal(address(usdg), liquidator, 10_000e6);
        vm.startPrank(liquidator);
        usdg.approve(address(morpho), 10_000e6);
        (uint256 seized, uint256 repaid) = morpho.liquidate(params, borrower, seize, shares, "");
        vm.stopPrank();
        return (seized, repaid);
    }

    function _assertSafe(address borrower, uint256 price) internal returns (uint256 ltv) {
        (, uint128 shares, uint128 collateral) = morpho.position(id, borrower);
        (,, uint128 borrowed, uint128 borrowShares,,) = morpho.market(id);
        uint256 debt = (uint256(shares) * (borrowed + 1) + borrowShares + 1e6 - 1) / (borrowShares + 1e6);
        uint256 value = uint256(collateral) * price / ORACLE_PRICE_SCALE;
        ltv = debt * WAD / value;
        assertLt(ltv * _incentive(), WAD * WAD);
        if (value * LLTV / WAD >= debt) {
            vm.expectRevert(bytes("position is healthy"));
            morpho.liquidate(params, borrower, 1, 0, "");
        } else {
            uint256 snapshot = vm.snapshotState();
            (uint256 seized,) = _liquidate(borrower, 1e17, 0);
            assertEq(seized, 1e17);
            vm.revertToState(snapshot);
        }
    }

    function _incentive() internal pure returns (uint256) {
        return _min(1.15e18, WAD * WAD / (WAD - 0.3e18 * (WAD - LLTV) / WAD));
    }

    function _min(uint256 a, uint256 b) internal pure returns (uint256) {
        return a < b ? a : b;
    }
}
