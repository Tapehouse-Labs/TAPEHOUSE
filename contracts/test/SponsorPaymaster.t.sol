// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {EntryPoint} from "account-abstraction/core/EntryPoint.sol";
import {IEntryPoint} from "account-abstraction/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/interfaces/PackedUserOperation.sol";
import {SimpleAccount} from "account-abstraction/samples/SimpleAccount.sol";
import {SimpleAccountFactory} from "account-abstraction/samples/SimpleAccountFactory.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Paymaster} from "@openzeppelin/contracts/account/paymaster/Paymaster.sol";
import {
    IEntryPoint as IOzEntryPoint,
    IPaymaster,
    PackedUserOperation as OzPackedUserOperation
} from "@openzeppelin/contracts/interfaces/IERC4337.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {SponsorPaymaster} from "../src/SponsorPaymaster.sol";
import {TargetDouble} from "./doubles/TargetDouble.sol";

contract SponsorPaymasterTest is Test {
    uint256 internal constant VERIFICATION_GAS = 500_000;
    uint256 internal constant CALL_GAS = 200_000;
    uint256 internal constant PAYMASTER_GAS = 100_000;
    uint256 internal constant PRE_VERIFICATION_GAS = 50_000;
    uint256 internal constant MAX_FEE = 1 gwei;
    uint256 internal constant REQUIRED = (VERIFICATION_GAS + CALL_GAS + PAYMASTER_GAS + PRE_VERIFICATION_GAS) * MAX_FEE;
    bytes32 internal constant REQUEST_TYPEHASH = keccak256(
        "UserOperationRequest(address sender,uint256 nonce,bytes initCode,bytes callData,bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees,uint256 paymasterVerificationGasLimit,uint256 paymasterPostOpGasLimit,uint48 validAfter,uint48 validUntil)"
    );

    EntryPoint internal entryPoint;
    SimpleAccountFactory internal factory;
    SponsorPaymaster internal paymaster;
    TargetDouble internal target;
    address internal owner = makeAddr("owner");
    address payable internal bundler = payable(makeAddr("bundler"));
    address internal signer;
    uint256 internal signerKey;
    address internal user;
    uint256 internal userKey;
    address internal account;

    event SignerSet(address indexed signer);
    event MaxCostSet(uint128 maxCost);
    event PausedSet(bool paused);

    function setUp() public {
        vm.warp(1_790_000_000);
        entryPoint = new EntryPoint();
        factory = new SimpleAccountFactory(entryPoint);
        (signer, signerKey) = makeAddrAndKey("signer");
        paymaster = new SponsorPaymaster(IOzEntryPoint(address(entryPoint)), owner, signer, uint128(REQUIRED));
        target = new TargetDouble();
        vm.deal(owner, 10 ether);
        vm.startPrank(owner);
        paymaster.deposit{value: 1 ether}();
        paymaster.addStake{value: 1 ether}(1 days);
        vm.stopPrank();
        (user, userKey) = makeAddrAndKey("user");
        account = factory.getAddress(user, 0);
    }

    function test_CreatesTheAccountForFree() public {
        _handle(_signed(_op("", true)));
        assertGt(account.code.length, 0);
        assertEq(paymaster.operations(account), 0);
        assertEq(paymaster.freeOperationsLeft(account), 3);
        assertLt(entryPoint.balanceOf(address(paymaster)), 1 ether);
    }

    function test_CountsAnOperationThatCreatesTheAccountAndCalls() public {
        _handle(_signed(_op(_touch(), true)));
        assertEq(target.calls(account), 1);
        assertEq(paymaster.operations(account), 1);
        assertEq(paymaster.freeOperationsLeft(account), 2);
    }

    function test_CountsAnOperationWhoseCallReverts() public {
        _handle(_signed(_op(_execute(abi.encodeCall(TargetDouble.refuse, ())), true)));
        assertEq(paymaster.operations(account), 1);
    }

    function test_SponsorsThreeOperationsAndRefusesTheFourth() public {
        _handle(_signed(_op("", true)));
        for (uint256 i; i < 3; ++i) {
            _handle(_signed(_op(_touch(), false)));
        }
        assertEq(target.calls(account), 3);
        assertEq(paymaster.freeOperationsLeft(account), 0);
        _refuse(
            _signed(_op(_touch(), false)),
            _reverted(abi.encodeWithSelector(SponsorPaymaster.FreeOperationsUsed.selector, account))
        );
    }

    function test_RefusesAnOperationWithoutTheServicesSignature() public {
        (, uint256 stranger) = makeAddrAndKey("stranger");
        _refuse(_signedBy(_op(_touch(), true), stranger, 0, 0), _failed("AA34 signature error"));
        PackedUserOperation memory op = _op(_touch(), true);
        op.paymasterAndData = abi.encodePacked(op.paymasterAndData, uint48(0), uint48(0), new bytes(65));
        op.signature = _accountSignature(op);
        _refuse(op, _failed("AA34 signature error"));
        op = _op(_touch(), true);
        _refuse(op, _failed("AA34 signature error"));
    }

    function test_RefusesAnOperationWhoseGasWasRaisedAfterSigning() public {
        vm.prank(owner);
        paymaster.setMaxCost(uint128(REQUIRED * 2));
        PackedUserOperation memory op = _signed(_op("", true));
        op.preVerificationGas += 1;
        op.signature = _accountSignature(op);
        _refuse(op, _failed("AA34 signature error"));
        op = _signed(_op("", true));
        op.gasFees = bytes32(MAX_FEE << 128 | MAX_FEE + 1);
        op.signature = _accountSignature(op);
        _refuse(op, _failed("AA34 signature error"));
    }

    function test_RefusesAnOperationOutsideItsWindow() public {
        _refuse(
            _signedBy(_op("", true), signerKey, 0, uint48(block.timestamp - 1)),
            _failed("AA32 paymaster expired or not due")
        );
        _refuse(
            _signedBy(_op("", true), signerKey, uint48(block.timestamp + 60), 0),
            _failed("AA32 paymaster expired or not due")
        );
        _handle(_signedBy(_op("", true), signerKey, uint48(block.timestamp - 60), uint48(block.timestamp + 60)));
        assertGt(account.code.length, 0);
    }

    function test_RefusesAnOperationAboveTheCostLimit() public {
        vm.prank(owner);
        paymaster.setMaxCost(uint128(REQUIRED - 1));
        _refuse(
            _signed(_op("", true)),
            _reverted(abi.encodeWithSelector(SponsorPaymaster.CostAboveLimit.selector, REQUIRED))
        );
    }

    function test_RefusesWhilePaused() public {
        vm.prank(owner);
        paymaster.setPaused(true);
        _refuse(_signed(_op("", true)), _reverted(abi.encodeWithSelector(SponsorPaymaster.SponsorshipPaused.selector)));
        vm.prank(owner);
        paymaster.setPaused(false);
        _handle(_signed(_op("", true)));
    }

    function test_OwnerRotatesTheSigner() public {
        (address next, uint256 nextKey) = makeAddrAndKey("next");
        vm.expectEmit(address(paymaster));
        emit SignerSet(next);
        vm.prank(owner);
        paymaster.setSigner(next);
        assertEq(paymaster.signer(), next);
        _refuse(_signed(_op("", true)), _failed("AA34 signature error"));
        _handle(_signedBy(_op("", true), nextKey, 0, 0));
    }

    function test_RefusesCallsFromAnyoneButTheEntryPoint() public {
        OzPackedUserOperation memory op = abi.decode(abi.encode(_op("", true)), (OzPackedUserOperation));
        vm.expectRevert(abi.encodeWithSelector(Paymaster.PaymasterUnauthorized.selector, address(this)));
        paymaster.validatePaymasterUserOp(op, bytes32(0), 0);
        vm.expectRevert(abi.encodeWithSelector(Paymaster.PaymasterUnauthorized.selector, address(this)));
        paymaster.postOp(IPaymaster.PostOpMode.opSucceeded, "", 0, 0);
    }

    function test_ReportsItsEntryPointAndSigner() public view {
        assertEq(address(paymaster.entryPoint()), address(entryPoint));
        assertEq(paymaster.signer(), signer);
    }

    function test_OwnerSetsTheCostLimitAndPause() public {
        vm.expectEmit(address(paymaster));
        emit MaxCostSet(7);
        vm.prank(owner);
        paymaster.setMaxCost(7);
        assertEq(paymaster.maxCostPerOperation(), 7);
        vm.expectEmit(address(paymaster));
        emit PausedSet(true);
        vm.prank(owner);
        paymaster.setPaused(true);
        assertTrue(paymaster.paused());
    }

    function test_AnyoneDepositsAndTheOwnerWithdraws() public {
        address payable treasury = payable(makeAddr("treasury"));
        vm.deal(address(this), 1 ether);
        paymaster.deposit{value: 1 ether}();
        assertEq(entryPoint.balanceOf(address(paymaster)), 2 ether);
        vm.prank(owner);
        paymaster.withdrawTo(treasury, 1.5 ether);
        assertEq(treasury.balance, 1.5 ether);
        assertEq(entryPoint.balanceOf(address(paymaster)), 0.5 ether);
    }

    function test_OwnerUnlocksAndWithdrawsTheStake() public {
        address payable treasury = payable(makeAddr("treasury"));
        vm.prank(owner);
        paymaster.unlockStake();
        vm.warp(block.timestamp + 1 days);
        vm.prank(owner);
        paymaster.withdrawStake(treasury);
        assertEq(treasury.balance, 1 ether);
    }

    function test_OnlyTheOwnerManages() public {
        address intruder = makeAddr("intruder");
        bytes memory notOwner = abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, intruder);
        vm.deal(intruder, 1 ether);
        vm.startPrank(intruder);
        vm.expectRevert(notOwner);
        paymaster.setSigner(intruder);
        vm.expectRevert(notOwner);
        paymaster.setMaxCost(1);
        vm.expectRevert(notOwner);
        paymaster.setPaused(true);
        vm.expectRevert(notOwner);
        paymaster.withdrawTo(payable(intruder), 1);
        vm.expectRevert(notOwner);
        paymaster.addStake{value: 1}(1 days);
        vm.expectRevert(notOwner);
        paymaster.unlockStake();
        vm.expectRevert(notOwner);
        paymaster.withdrawStake(payable(intruder));
        vm.stopPrank();
    }

    function test_KeepsAnOwner() public {
        vm.prank(owner);
        vm.expectRevert(SponsorPaymaster.OwnershipCannotBeRenounced.selector);
        paymaster.renounceOwnership();
    }

    function test_GasOfEachCall() public {
        uint256 state = vm.snapshotState();
        _validate(_signed(_op(_touch(), true)));
        vm.snapshotGasLastCall("validate first operation");
        vm.revertToStateAndDelete(state);
        _handle(_signed(_op(_touch(), true)));
        vm.snapshotGasLastCall("handleOps create and call");
        _validate(_signed(_op(_touch(), false)));
        vm.snapshotGasLastCall("validate later operation");
        _handle(_signed(_op(_touch(), false)));
        vm.snapshotGasLastCall("handleOps call");
    }

    function testFuzz_RefusesAnySignatureButTheServices(bytes calldata signature) public {
        PackedUserOperation memory op = _op("", true);
        op.paymasterAndData = abi.encodePacked(op.paymasterAndData, uint48(0), uint48(0), signature);
        op.signature = _accountSignature(op);
        _refuse(op, _failed("AA34 signature error"));
    }

    function _op(bytes memory callData, bool create) internal view returns (PackedUserOperation memory op) {
        op.sender = account;
        op.nonce = entryPoint.getNonce(account, 0);
        if (create) op.initCode = abi.encodePacked(address(factory), abi.encodeCall(factory.createAccount, (user, 0)));
        op.callData = callData;
        op.accountGasLimits = bytes32(VERIFICATION_GAS << 128 | CALL_GAS);
        op.preVerificationGas = PRE_VERIFICATION_GAS;
        op.gasFees = bytes32(MAX_FEE << 128 | MAX_FEE);
        op.paymasterAndData = abi.encodePacked(address(paymaster), uint128(PAYMASTER_GAS), uint128(0));
        op.signature = _accountSignature(op);
    }

    function _signed(PackedUserOperation memory op) internal view returns (PackedUserOperation memory) {
        return _signedBy(op, signerKey, 0, 0);
    }

    function _signedBy(PackedUserOperation memory op, uint256 key, uint48 validAfter, uint48 validUntil)
        internal
        view
        returns (PackedUserOperation memory)
    {
        bytes32 request = keccak256(
            abi.encode(
                REQUEST_TYPEHASH,
                op.sender,
                op.nonce,
                keccak256(op.initCode),
                keccak256(op.callData),
                op.accountGasLimits,
                op.preVerificationGas,
                op.gasFees,
                PAYMASTER_GAS,
                uint256(0),
                validAfter,
                validUntil
            )
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, MessageHashUtils.toTypedDataHash(_domainSeparator(), request));
        op.paymasterAndData = abi.encodePacked(op.paymasterAndData, validAfter, validUntil, r, s, v);
        op.signature = _accountSignature(op);
        return op;
    }

    function _domainSeparator() internal view returns (bytes32) {
        (, string memory name, string memory version, uint256 chainId, address verifyingContract,,) =
            paymaster.eip712Domain();
        return keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256(bytes(name)),
                keccak256(bytes(version)),
                chainId,
                verifyingContract
            )
        );
    }

    function _accountSignature(PackedUserOperation memory op) internal view returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) =
            vm.sign(userKey, MessageHashUtils.toEthSignedMessageHash(entryPoint.getUserOpHash(op)));
        return abi.encodePacked(r, s, v);
    }

    function _handle(PackedUserOperation memory op) internal {
        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = op;
        vm.fee(MAX_FEE);
        entryPoint.handleOps(ops, bundler);
    }

    function _refuse(PackedUserOperation memory op, bytes memory reason) internal {
        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0] = op;
        vm.fee(MAX_FEE);
        vm.expectRevert(reason);
        entryPoint.handleOps(ops, bundler);
    }

    function _validate(PackedUserOperation memory op) internal {
        OzPackedUserOperation memory oz = abi.decode(abi.encode(op), (OzPackedUserOperation));
        vm.prank(address(entryPoint));
        paymaster.validatePaymasterUserOp(oz, bytes32(0), REQUIRED);
    }

    function _reverted(bytes memory inner) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IEntryPoint.FailedOpWithRevert.selector, 0, "AA33 reverted", inner);
    }

    function _failed(string memory reason) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IEntryPoint.FailedOp.selector, 0, reason);
    }

    function _execute(bytes memory data) internal view returns (bytes memory) {
        return abi.encodeCall(SimpleAccount.execute, (address(target), 0, data));
    }

    function _touch() internal view returns (bytes memory) {
        return _execute(abi.encodeCall(TargetDouble.touch, ()));
    }
}
