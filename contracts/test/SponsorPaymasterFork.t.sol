// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IEntryPoint} from "account-abstraction/interfaces/IEntryPoint.sol";
import {PackedUserOperation} from "account-abstraction/interfaces/PackedUserOperation.sol";
import {SimpleAccount} from "account-abstraction/samples/SimpleAccount.sol";
import {SimpleAccountFactory} from "account-abstraction/samples/SimpleAccountFactory.sol";
import {IEntryPoint as IOzEntryPoint} from "@openzeppelin/contracts/interfaces/IERC4337.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MessageHashUtils} from "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import {SponsorPaymaster} from "../src/SponsorPaymaster.sol";
import {TargetDouble} from "./doubles/TargetDouble.sol";

contract SponsorPaymasterForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    uint256 internal constant ROBINHOOD_TESTNET_BLOCK = 127_372_363;
    uint256 internal constant PAYMASTER_GAS = 100_000;
    bytes32 internal constant REQUEST_TYPEHASH = keccak256(
        "UserOperationRequest(address sender,uint256 nonce,bytes initCode,bytes callData,bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees,uint256 paymasterVerificationGasLimit,uint256 paymasterPostOpGasLimit,uint48 validAfter,uint48 validUntil)"
    );

    function test_SponsorsANewAccountOnRobinhood() public {
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
        _sponsor(vm.readFile("../deployments/4663.json"));
    }

    function test_SponsorsANewAccountOnRobinhoodTestnet() public {
        vm.createSelectFork("robinhood-testnet", ROBINHOOD_TESTNET_BLOCK);
        _sponsor(vm.readFile("../deployments/46630.json"));
    }

    function _sponsor(string memory json) internal {
        IEntryPoint entryPoint = IEntryPoint(vm.parseJsonAddress(json, ".erc4337.EntryPoint"));
        SimpleAccountFactory factory = SimpleAccountFactory(vm.parseJsonAddress(json, ".erc4337.SimpleAccountFactory"));
        IERC20 usdg = IERC20(vm.parseJsonAddress(json, ".tokens.USDG"));
        TargetDouble target = new TargetDouble();
        (address signer, uint256 signerKey) = makeAddrAndKey("signer");
        SponsorPaymaster paymaster = _paymaster(entryPoint, signer);
        (address user, uint256 key) = makeAddrAndKey("user");
        address account = factory.getAddress(user, 0);
        address[] memory targets = new address[](2);
        (targets[0], targets[1]) = (address(usdg), address(target));
        bytes[] memory data = new bytes[](2);
        (data[0], data[1]) =
        (abi.encodeCall(IERC20.approve, (address(target), 1e6)), abi.encodeCall(TargetDouble.touch, ()));
        PackedUserOperation[] memory ops = new PackedUserOperation[](1);
        ops[0].sender = account;
        ops[0].initCode = abi.encodePacked(address(factory), abi.encodeCall(factory.createAccount, (user, 0)));
        ops[0].callData = abi.encodeCall(SimpleAccount.executeBatch, (targets, new uint256[](0), data));
        ops[0].accountGasLimits = bytes32(uint256(500_000) << 128 | 200_000);
        ops[0].preVerificationGas = 50_000;
        ops[0].gasFees = bytes32(uint256(block.basefee) << 128 | block.basefee * 2);
        _sign(entryPoint, paymaster, ops[0], signerKey, key);

        entryPoint.handleOps(ops, payable(makeAddr("bundler")));

        assertEq(address(factory.accountImplementation().entryPoint()), address(entryPoint));
        assertEq(account.balance, 0);
        assertEq(usdg.allowance(account, address(target)), 1e6);
        assertEq(target.calls(account), 1);
        assertEq(paymaster.freeOperationsLeft(account), 2);
        assertLt(entryPoint.balanceOf(address(paymaster)), 1 ether);
    }

    function _sign(
        IEntryPoint entryPoint,
        SponsorPaymaster paymaster,
        PackedUserOperation memory op,
        uint256 signerKey,
        uint256 key
    ) internal view {
        uint48 validUntil = uint48(block.timestamp + 600);
        op.paymasterAndData = abi.encodePacked(
            address(paymaster),
            uint128(PAYMASTER_GAS),
            uint128(0),
            uint48(0),
            validUntil,
            _sponsorSignature(paymaster, op, signerKey, validUntil)
        );
        (uint8 v, bytes32 r, bytes32 s) =
            vm.sign(key, MessageHashUtils.toEthSignedMessageHash(entryPoint.getUserOpHash(op)));
        op.signature = abi.encodePacked(r, s, v);
    }

    function _paymaster(IEntryPoint entryPoint, address signer) internal returns (SponsorPaymaster paymaster) {
        address owner = makeAddr("owner");
        paymaster = new SponsorPaymaster(IOzEntryPoint(address(entryPoint)), owner, signer, 0.01 ether);
        vm.deal(owner, 2 ether);
        vm.startPrank(owner);
        paymaster.deposit{value: 1 ether}();
        paymaster.addStake{value: 1 ether}(1 days);
        vm.stopPrank();
    }

    function _sponsorSignature(
        SponsorPaymaster paymaster,
        PackedUserOperation memory op,
        uint256 key,
        uint48 validUntil
    ) internal view returns (bytes memory) {
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
                uint48(0),
                validUntil
            )
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(key, MessageHashUtils.toTypedDataHash(_domain(paymaster), request));
        return abi.encodePacked(r, s, v);
    }

    function _domain(SponsorPaymaster paymaster) internal view returns (bytes32) {
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
}
