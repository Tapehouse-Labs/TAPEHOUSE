// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract SupplyVaultHandler is Test {
    StubUsdg internal immutable usdg;
    SupplyVault internal immutable vault;
    address internal immutable borrower;
    address[3] internal lenders = [address(0x1001), address(0x1002), address(0x1003)];
    bool public shareValueHeld = true;

    modifier valued() {
        uint256 before = vault.convertToAssets(1e12);
        _;
        if (vault.convertToAssets(1e12) < before) shareValueHeld = false;
    }

    constructor(StubUsdg usdg_, SupplyVault vault_, address borrower_) {
        usdg = usdg_;
        vault = vault_;
        borrower = borrower_;
    }

    function deposit(uint256 who, uint256 assets) external valued {
        address lender = lenders[who % 3];
        assets = bound(assets, 1, 1e12 * 1e6);
        usdg.mint(lender, assets);
        vm.startPrank(lender);
        usdg.approve(address(vault), assets);
        vault.deposit(assets, lender);
        vm.stopPrank();
    }

    function redeem(uint256 who, uint256 shares) external valued {
        address lender = lenders[who % 3];
        shares = bound(shares, 0, vault.maxRedeem(lender));
        vm.prank(lender);
        vault.redeem(shares, lender, lender);
    }

    function borrow(uint256 assets) external valued {
        assets = bound(assets, 0, vault.idle());
        vm.prank(borrower);
        vault.borrow(assets, borrower);
    }

    function repay(uint256 assets) external valued {
        assets = bound(assets, 0, vault.debt());
        usdg.mint(borrower, assets);
        vm.startPrank(borrower);
        usdg.approve(address(vault), assets);
        vault.repay(assets);
        vm.stopPrank();
    }

    function donate(uint256 assets) external valued {
        assets = bound(assets, 0, 1e9 * 1e6);
        usdg.mint(address(vault), assets);
    }

    function wait(uint256 seconds_) external valued {
        skip(bound(seconds_, 0, 30 days));
    }
}

contract SupplyVaultInvariantTest is Test {
    StubUsdg internal usdg;
    SupplyVault internal vault;
    SupplyVaultHandler internal handler;

    function setUp() public {
        usdg = new StubUsdg();
        address borrower = makeAddr("borrower");
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        vault.setBorrower(borrower);
        handler = new SupplyVaultHandler(usdg, vault, borrower);
        targetContract(address(handler));
    }

    function invariant_TotalAssetsAreTheUsdgHeldPlusTheDebt() public view {
        assertEq(vault.totalAssets(), vault.idle() + vault.debt());
    }

    function invariant_TheVaultHoldsTheUsdgItCounts() public view {
        assertGe(usdg.balanceOf(address(vault)), vault.idle());
    }

    function invariant_AShareNeverLosesValue() public view {
        assertTrue(handler.shareValueHeld());
    }
}
