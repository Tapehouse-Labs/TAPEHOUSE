// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract SupplyVaultTest is Test {
    uint256 internal constant USDG = 1e6;
    bytes32 internal constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");

    StubUsdg internal usdg;
    SupplyVault internal vault;
    address internal owner = makeAddr("owner");
    address internal lender = makeAddr("lender");
    address internal borrower = makeAddr("borrower");

    function setUp() public {
        usdg = new StubUsdg();
        vault = new SupplyVault(IUSDG(address(usdg)), owner, _model(90_00, 0, 6_00, 40_00));
        vm.prank(owner);
        vault.setBorrower(borrower);
    }

    function test_SharesCarrySixMoreDecimalsThanUsdg() public {
        assertEq(vault.name(), "Tapehouse USDG Supply");
        assertEq(vault.symbol(), "thUSDG");
        assertEq(vault.decimals(), 12);
        assertEq(vault.asset(), address(usdg));
        uint256 shares = _deposit(lender, 1_000 * USDG);
        assertEq(shares, 1_000 * USDG * 1e6);
        vm.prank(lender);
        assertEq(vault.redeem(shares, lender, lender), 1_000 * USDG);
        assertEq(usdg.balanceOf(lender), 1_000 * USDG);
    }

    function test_UsdgSentToTheVaultChangesNoShare() public {
        address attacker = makeAddr("attacker");
        _deposit(attacker, 1);
        usdg.mint(attacker, 1_000_000 * USDG);
        vm.prank(attacker);
        usdg.transfer(address(vault), 1_000_000 * USDG);
        assertEq(vault.totalAssets(), 1);
        uint256 shares = _deposit(lender, 1_000 * USDG);
        assertEq(shares, 1_000 * USDG * 1e6);
        assertEq(vault.previewRedeem(shares), 1_000 * USDG);
    }

    function test_OnlyTheBorrowerBorrowsAndOnlyWhatTheVaultHolds() public {
        _deposit(lender, 1_000 * USDG);
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.NotBorrower.selector, lender));
        vm.prank(lender);
        vault.borrow(1, lender);
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.InsufficientLiquidity.selector, 1_001 * USDG, 1_000 * USDG));
        vm.prank(borrower);
        vault.borrow(1_001 * USDG, borrower);
        vm.expectEmit(address(vault));
        emit SupplyVault.Borrow(borrower, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(1_000 * USDG, borrower);
        assertEq(vault.idle(), 0);
        assertEq(vault.debt(), 1_000 * USDG);
        assertEq(usdg.balanceOf(borrower), 1_000 * USDG);
    }

    function test_InterestAccruesAtTheModelsRate() public {
        _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(900 * USDG, borrower);
        assertEq(vault.utilization(), 0.9e18);
        assertEq(vault.borrowRate(), 0.06e18);
        assertEq(vault.supplyRate(), 0.054e18);
        skip(365 days);
        assertEq(vault.debt(), 954 * USDG);
        assertEq(vault.totalAssets(), 1_054 * USDG);
        vm.expectEmit(address(vault));
        emit SupplyVault.InterestAccrued(54 * USDG, 1.06e27);
        vm.prank(borrower);
        vault.repay(0);
        assertEq(vault.borrowIndex(), 1.06e27);
        assertEq(vault.scaledDebt(), 900 * USDG);
    }

    function test_AccruingEverySecondErasesNoInterest() public {
        _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(400 * USDG, borrower);
        uint256 once = 400 * USDG + 400 * USDG * vault.borrowRate() * 1 hours / (1e18 * 365 days);
        for (uint256 i; i < 1 hours; ++i) {
            skip(1);
            vm.prank(borrower);
            vault.repay(0);
        }
        assertGe(vault.debt(), once);
        assertLe(vault.debt(), once + 2);
    }

    function test_ADepositAccruesAtTheRateBeforeIt() public {
        _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(900 * USDG, borrower);
        skip(365 days);
        _deposit(makeAddr("second"), 1_000 * USDG);
        assertEq(vault.borrowIndex(), 1.06e27);
        assertEq(vault.lastAccrual(), block.timestamp);
        assertEq(vault.debt(), 954 * USDG);
    }

    function test_TheBorrowRateFollowsItsTwoSlopes() public {
        _deposit(lender, 1_000 * USDG);
        assertEq(vault.borrowRate(), 0);
        _borrowTo(450);
        assertEq(vault.borrowRate(), 0.03e18);
        _borrowTo(900);
        assertEq(vault.borrowRate(), 0.06e18);
        _borrowTo(950);
        assertEq(vault.borrowRate(), 0.26e18);
        _borrowTo(1_000);
        assertEq(vault.borrowRate(), 0.46e18);
        assertEq(vault.supplyRate(), 0.46e18);
    }

    function test_TheRateModelStaysWithinItsBounds() public {
        vm.startPrank(owner);
        vm.expectRevert(SupplyVault.InvalidRateModel.selector);
        vault.setRateModel(_model(99, 0, 1, 1));
        vm.expectRevert(SupplyVault.InvalidRateModel.selector);
        vault.setRateModel(_model(99_01, 0, 1, 1));
        vm.expectRevert(SupplyVault.InvalidRateModel.selector);
        vault.setRateModel(_model(90_00, 0, 2, 1));
        vm.expectRevert(SupplyVault.InvalidRateModel.selector);
        vault.setRateModel(_model(90_00, 1, 500_00, 500_00));
        vault.setRateModel(_model(1_00, 0, 500_00, 500_00));
        vm.expectEmit(address(vault));
        emit SupplyVault.RateModelSet(99_00, 1_00, 2_00, 3_00);
        vault.setRateModel(_model(99_00, 1_00, 2_00, 3_00));
        vm.stopPrank();
        SupplyVault.RateModel memory model = vault.rateModel();
        assertEq(model.optimal, 99_00);
        assertEq(model.base, 1_00);
        assertEq(model.slope1, 2_00);
        assertEq(model.slope2, 3_00);
        vm.expectRevert(SupplyVault.InvalidRateModel.selector);
        new SupplyVault(IUSDG(address(usdg)), owner, _model(0, 0, 0, 0));
    }

    function test_ANewRateAppliesFromWhenItIsSet() public {
        _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(900 * USDG, borrower);
        skip(365 days / 2);
        vm.prank(owner);
        vault.setRateModel(_model(90_00, 0, 12_00, 40_00));
        assertEq(vault.borrowIndex(), 1.03e27);
        assertEq(vault.debt(), 927 * USDG);
        uint256 rate = vault.borrowRate();
        assertGt(rate, 0.12e18);
        skip(365 days / 2);
        uint256 index = 1.03e27 + 1.03e27 * rate * (365 days / 2) / (1e18 * 365 days);
        assertEq(vault.debt(), (900 * USDG * index + 1e27 - 1) / 1e27);
    }

    function test_OnlyTheOwnerSetsTheBorrowerAndTheRate() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, lender));
        vm.prank(lender);
        vault.setBorrower(lender);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, lender));
        vm.prank(lender);
        vault.setRateModel(_model(90_00, 0, 6_00, 40_00));
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.BorrowerAlreadySet.selector, borrower));
        vm.prank(owner);
        vault.setBorrower(owner);
        SupplyVault fresh = new SupplyVault(IUSDG(address(usdg)), owner, _model(90_00, 0, 6_00, 40_00));
        assertEq(fresh.maxDeposit(lender), 0);
        assertEq(fresh.maxMint(lender), 0);
        vm.expectRevert(SupplyVault.InvalidBorrower.selector);
        vm.prank(owner);
        fresh.setBorrower(address(0));
        vm.expectEmit(address(fresh));
        emit SupplyVault.BorrowerSet(borrower);
        vm.prank(owner);
        fresh.setBorrower(borrower);
        assertEq(fresh.maxDeposit(lender), type(uint256).max);
        vm.prank(owner);
        vault.transferOwnership(lender);
        assertEq(vault.owner(), owner);
        vm.prank(lender);
        vault.acceptOwnership();
        assertEq(vault.owner(), lender);
    }

    function test_OnlyTheBorrowerRepaysAndNoMoreThanTheDebt() public {
        _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(500 * USDG, borrower);
        address payer = makeAddr("payer");
        usdg.mint(payer, 800 * USDG);
        vm.startPrank(payer);
        usdg.approve(address(vault), 800 * USDG);
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.NotBorrower.selector, payer));
        vault.repay(1);
        vm.stopPrank();
        usdg.mint(borrower, 300 * USDG);
        vm.startPrank(borrower);
        usdg.approve(address(vault), 800 * USDG);
        vm.expectEmit(address(vault));
        emit SupplyVault.Repay(borrower, 500 * USDG);
        assertEq(vault.repay(800 * USDG), 500 * USDG);
        vm.stopPrank();
        assertEq(vault.debt(), 0);
        assertEq(vault.idle(), 1_000 * USDG);
        assertEq(usdg.balanceOf(borrower), 300 * USDG);
    }

    function test_WithdrawalsWaitForRepayment() public {
        uint256 shares = _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(700 * USDG, borrower);
        assertEq(vault.maxWithdraw(lender), 300 * USDG);
        assertEq(vault.maxRedeem(lender), 300 * USDG * 1e6);
        vm.expectRevert(
            abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxWithdraw.selector, lender, 301 * USDG, 300 * USDG)
        );
        vm.prank(lender);
        vault.withdraw(301 * USDG, lender, lender);
        vm.prank(lender);
        vault.withdraw(300 * USDG, lender, lender);
        assertEq(vault.balanceOf(lender), shares - 300 * USDG * 1e6);
    }

    function test_APausedUsdgClosesTheVault() public {
        _deposit(lender, 1_000 * USDG);
        usdg.mint(lender, 1 * USDG);
        usdg.pause();
        assertEq(vault.maxDeposit(lender), 0);
        assertEq(vault.maxMint(lender), 0);
        assertEq(vault.maxWithdraw(lender), 0);
        assertEq(vault.maxRedeem(lender), 0);
        vm.expectRevert(abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxWithdraw.selector, lender, 1, 0));
        vm.prank(lender);
        vault.withdraw(1, lender, lender);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        vm.prank(borrower);
        vault.borrow(1, borrower);
        usdg.unpause();
        assertEq(vault.maxWithdraw(lender), 1_000 * USDG);
        assertEq(vault.maxDeposit(lender), type(uint256).max);
    }

    function test_AFrozenVaultClosesAndAWipeLeavesLendersAClaimOnNothing() public {
        _deposit(lender, 1_000 * USDG);
        usdg.freeze(address(vault));
        assertEq(vault.maxDeposit(lender), 0);
        assertEq(vault.maxRedeem(lender), 0);
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        vm.prank(borrower);
        vault.borrow(1, borrower);
        usdg.wipeFrozenAddress(address(vault));
        usdg.unfreeze(address(vault));
        assertEq(vault.maxWithdraw(lender), 0);
        assertEq(vault.maxDeposit(lender), 0);
        assertEq(vault.maxMint(lender), 0);
        assertEq(vault.totalAssets(), 1_000 * USDG);
        address late = makeAddr("late");
        usdg.mint(late, 1_000 * USDG);
        vm.startPrank(late);
        usdg.approve(address(vault), 1_000 * USDG);
        vm.expectRevert(abi.encodeWithSelector(ERC4626.ERC4626ExceededMaxDeposit.selector, late, 1_000 * USDG, 0));
        vault.deposit(1_000 * USDG, late);
        vm.stopPrank();
        vm.expectEmit(address(vault));
        emit SupplyVault.Sync(0);
        vault.sync();
        assertEq(vault.idle(), 0);
        assertEq(vault.totalAssets(), 0);
        uint256 lateShares = _deposit(late, 1_000 * USDG);
        assertEq(vault.previewRedeem(lateShares), 1_000 * USDG);
        assertLe(vault.maxWithdraw(lender), 1);
        usdg.mint(address(vault), 5 * USDG);
        vault.sync();
        assertEq(vault.idle(), 1_000 * USDG);
    }

    function test_TheBorrowerWritesDebtOffAndLendersBearIt() public {
        uint256 shares = _deposit(lender, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(500 * USDG, borrower);
        vm.expectRevert(abi.encodeWithSelector(SupplyVault.NotBorrower.selector, lender));
        vm.prank(lender);
        vault.writeOff(1);
        vm.expectEmit(address(vault));
        emit SupplyVault.WriteOff(200 * USDG);
        vm.prank(borrower);
        assertEq(vault.writeOff(200 * USDG), 200 * USDG);
        assertEq(vault.debt(), 300 * USDG);
        assertEq(vault.totalAssets(), 800 * USDG);
        assertEq(vault.previewRedeem(shares), 800 * USDG);
        vm.prank(borrower);
        assertEq(vault.writeOff(1_000 * USDG), 300 * USDG);
        assertEq(vault.scaledDebt(), 0);
    }

    function test_AMaxViewNeverRevertsWhenUsdgDoes() public {
        _deposit(lender, 1_000 * USDG);
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("paused()"), "");
        assertEq(vault.maxDeposit(lender), 0);
        assertEq(vault.maxRedeem(lender), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("isFrozen(address)", address(vault)), "");
        assertEq(vault.maxMint(lender), 0);
        vm.clearMockedCalls();
        vm.mockCallRevert(address(usdg), abi.encodeWithSignature("balanceOf(address)", address(vault)), "");
        assertEq(vault.maxWithdraw(lender), 0);
    }

    function test_AFrozenLenderCannotReceive() public {
        _deposit(lender, 1_000 * USDG);
        usdg.freeze(lender);
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        vm.prank(lender);
        vault.withdraw(1 * USDG, lender, lender);
    }

    function test_DepositWithAPermitEvenIfItWasSpent() public {
        Vm.Wallet memory wallet = vm.createWallet("wallet");
        usdg.mint(wallet.addr, 1_000 * USDG);
        (uint8 v, bytes32 r, bytes32 s) = _permit(wallet, 400 * USDG, block.timestamp);
        vm.prank(wallet.addr);
        assertEq(vault.depositWithPermit(400 * USDG, wallet.addr, block.timestamp, v, r, s), 400 * USDG * 1e6);
        (v, r, s) = _permit(wallet, 600 * USDG, block.timestamp);
        usdg.permit(wallet.addr, address(vault), 600 * USDG, block.timestamp, v, r, s);
        vm.prank(wallet.addr);
        vault.depositWithPermit(600 * USDG, wallet.addr, block.timestamp, v, r, s);
        assertEq(vault.balanceOf(wallet.addr), 1_000 * USDG * 1e6);
        assertEq(usdg.balanceOf(wallet.addr), 0);
    }

    function test_GasOfEachCall() public {
        usdg.mint(lender, 1_000 * USDG);
        vm.startPrank(lender);
        usdg.approve(address(vault), 1_000 * USDG);
        vault.deposit(1_000 * USDG, lender);
        vm.snapshotGasLastCall("deposit");
        vm.stopPrank();
        vm.prank(borrower);
        vault.borrow(900 * USDG, borrower);
        vm.snapshotGasLastCall("borrow");
        skip(1 days);
        usdg.mint(borrower, 1 * USDG);
        vm.startPrank(borrower);
        usdg.approve(address(vault), 901 * USDG);
        vault.repay(901 * USDG);
        vm.snapshotGasLastCall("repay");
        vm.stopPrank();
        vm.prank(lender);
        vault.withdraw(500 * USDG, lender, lender);
        vm.snapshotGasLastCall("withdraw");
        Vm.Wallet memory wallet = vm.createWallet("wallet");
        usdg.mint(wallet.addr, 100 * USDG);
        (uint8 v, bytes32 r, bytes32 s) = _permit(wallet, 100 * USDG, block.timestamp);
        vm.prank(wallet.addr);
        vault.depositWithPermit(100 * USDG, wallet.addr, block.timestamp, v, r, s);
        vm.snapshotGasLastCall("depositWithPermit");
    }

    function testFuzz_ADepositRedeemedReturnsNoMore(uint256 assets, uint256 borrowed, uint32 elapsed) public {
        assets = bound(assets, 1, 1e15 * USDG);
        _deposit(lender, 1_000 * USDG);
        borrowed = bound(borrowed, 0, 1_000 * USDG);
        vm.prank(borrower);
        vault.borrow(borrowed, borrower);
        skip(elapsed);
        address other = makeAddr("other");
        uint256 shares = _deposit(other, assets);
        assertLe(vault.previewRedeem(shares), assets);
    }

    function testFuzz_TheBorrowRateNeverFallsAsUtilizationRises(uint256 a, uint256 b) public {
        _deposit(lender, 1_000_000 * USDG);
        a = bound(a, 0, 1_000_000 * USDG);
        b = bound(b, a, 1_000_000 * USDG);
        vm.prank(borrower);
        vault.borrow(a, borrower);
        uint256 low = vault.borrowRate();
        vm.prank(borrower);
        vault.borrow(b - a, borrower);
        assertGe(vault.borrowRate(), low);
        assertLe(vault.borrowRate(), 0.46e18);
    }

    function _deposit(address who, uint256 assets) internal returns (uint256 shares) {
        usdg.mint(who, assets);
        vm.startPrank(who);
        usdg.approve(address(vault), assets);
        shares = vault.deposit(assets, who);
        vm.stopPrank();
    }

    function _borrowTo(uint256 total) internal {
        uint256 assets = total * USDG - vault.debt();
        vm.prank(borrower);
        vault.borrow(assets, borrower);
    }

    function _model(uint16 optimal, uint32 base, uint32 slope1, uint32 slope2)
        internal
        pure
        returns (SupplyVault.RateModel memory)
    {
        return SupplyVault.RateModel(optimal, base, slope1, slope2);
    }

    function _permit(Vm.Wallet memory wallet, uint256 value, uint256 deadline)
        internal
        view
        returns (uint8, bytes32, bytes32)
    {
        bytes32 structHash = keccak256(
            abi.encode(PERMIT_TYPEHASH, wallet.addr, address(vault), value, usdg.nonces(wallet.addr), deadline)
        );
        return vm.sign(wallet, keccak256(abi.encodePacked("\x19\x01", usdg.DOMAIN_SEPARATOR(), structHash)));
    }

    function test_OwnershipCannotBeRenounced() public {
        vm.expectRevert(SupplyVault.OwnershipCannotBeRenounced.selector);
        vm.prank(owner);
        vault.renounceOwnership();
        assertEq(vault.owner(), owner);
    }
}
