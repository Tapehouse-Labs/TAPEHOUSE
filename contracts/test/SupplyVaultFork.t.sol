// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {IUSDG as IUsdgIssuer} from "./conformance/Interfaces.sol";

contract SupplyVaultForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    address internal constant OPERATOR = 0x3Af3e85f4f97De7AD0f000B724Fb77fE5ffc024B;
    bytes32 internal constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");

    IUsdgIssuer internal usdg;
    SupplyVault internal vault;
    address internal morpho;
    address internal lender = makeAddr("lender");
    address internal borrower = makeAddr("borrower");

    function setUp() public {
        string memory json = vm.readFile("../deployments/4663.json");
        usdg = IUsdgIssuer(vm.parseJsonAddress(json, ".tokens.USDG"));
        morpho = vm.parseJsonAddress(json, ".morpho.Blue");
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        vault.setBorrower(borrower);
        vm.prank(morpho);
        usdg.transfer(lender, 1_000e6);
        vm.startPrank(lender);
        usdg.approve(address(vault), 1_000e6);
        vault.deposit(1_000e6, lender);
        vm.stopPrank();
    }

    function test_LendersEarnTheBorrowersInterestInUsdg() public {
        vm.prank(borrower);
        vault.borrow(900e6, borrower);
        assertEq(usdg.balanceOf(borrower), 900e6);
        skip(365 days);
        vm.prank(morpho);
        usdg.transfer(borrower, 54e6);
        vm.startPrank(borrower);
        usdg.approve(address(vault), 954e6);
        assertEq(vault.repay(954e6), 954e6);
        vm.stopPrank();
        uint256 shares = vault.balanceOf(lender);
        vm.prank(lender);
        assertEq(vault.redeem(shares, lender, lender), 1_054e6 - 1);
    }

    function test_IssuerPauseClosesTheVaultWithUsdgsOwnError() public {
        vm.prank(OPERATOR);
        usdg.pause();
        assertEq(vault.maxDeposit(lender), 0);
        assertEq(vault.maxRedeem(lender), 0);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        vm.prank(borrower);
        vault.borrow(1, borrower);
    }

    function test_IssuerFreezeCanWipeWhatTheVaultHolds() public {
        vm.prank(OPERATOR);
        usdg.freeze(address(vault));
        assertEq(vault.maxWithdraw(lender), 0);
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        vm.prank(borrower);
        vault.borrow(1, borrower);
        vm.prank(OPERATOR);
        usdg.wipeFrozenAddress(address(vault));
        assertEq(usdg.balanceOf(address(vault)), 0);
        assertEq(vault.totalAssets(), 1_000e6);
        vault.sync();
        assertEq(vault.totalAssets(), 0);
    }

    function test_DepositWithAGlobalDollarPermit() public {
        Vm.Wallet memory wallet = vm.createWallet("wallet");
        vm.prank(morpho);
        usdg.transfer(wallet.addr, 500e6);
        bytes32 structHash = keccak256(
            abi.encode(PERMIT_TYPEHASH, wallet.addr, address(vault), 500e6, usdg.nonces(wallet.addr), block.timestamp)
        );
        (uint8 v, bytes32 r, bytes32 s) =
            vm.sign(wallet, keccak256(abi.encodePacked("\x19\x01", usdg.DOMAIN_SEPARATOR(), structHash)));
        vm.prank(wallet.addr);
        vault.depositWithPermit(500e6, wallet.addr, block.timestamp, v, r, s);
        assertEq(vault.balanceOf(wallet.addr), 500e6 * 1e6);
    }
}
