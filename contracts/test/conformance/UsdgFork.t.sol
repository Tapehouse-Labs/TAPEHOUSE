// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {IUSDG, ITimelockController} from "./Interfaces.sol";

contract UsdgForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    bytes32 internal constant IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    address internal constant USDG_IMPLEMENTATION = 0x68184C449E1a8f34fA18d289737129FD27B66f8F;
    address internal constant TIMELOCK = 0xcFA0388f5ddf905FdC08c45c716C15Dc10A14C6F;
    address internal constant OPERATOR = 0x3Af3e85f4f97De7AD0f000B724Fb77fE5ffc024B;

    bytes32 internal constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 internal constant PAUSE_ROLE = keccak256("PAUSE_ROLE");
    bytes32 internal constant ASSET_PROTECTION_ROLE = keccak256("ASSET_PROTECTION_ROLE");
    bytes32 internal constant PROPOSER_ROLE = keccak256("PROPOSER_ROLE");
    bytes32 internal constant EXECUTOR_ROLE = keccak256("EXECUTOR_ROLE");
    bytes32 internal constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");

    IUSDG internal usdg;
    address internal morpho;

    function setUp() public {
        string memory json = vm.readFile("../deployments/4663.json");
        usdg = IUSDG(vm.parseJsonAddress(json, ".tokens.USDG"));
        morpho = vm.parseJsonAddress(json, ".morpho.Blue");
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
    }

    function test_UsdgIsUupsBehindATwentyFourHourTimelock() public {
        assertEq(usdg.name(), "Global Dollar");
        assertEq(usdg.symbol(), "USDG");
        assertEq(usdg.decimals(), 6);
        assertEq(address(uint160(uint256(vm.load(address(usdg), IMPLEMENTATION_SLOT)))), USDG_IMPLEMENTATION);
        assertEq(usdg.owner(), TIMELOCK);
        assertEq(usdg.defaultAdmin(), TIMELOCK);
        assertTrue(usdg.hasRole(DEFAULT_ADMIN_ROLE, TIMELOCK));
        assertEq(usdg.defaultAdminDelay(), 3 hours);
        assertEq(ITimelockController(TIMELOCK).getMinDelay(), 24 hours);
        assertTrue(ITimelockController(TIMELOCK).hasRole(DEFAULT_ADMIN_ROLE, TIMELOCK));
        assertTrue(ITimelockController(TIMELOCK).hasRole(PROPOSER_ROLE, OPERATOR));
        assertTrue(ITimelockController(TIMELOCK).hasRole(EXECUTOR_ROLE, OPERATOR));

        vm.expectRevert(
            bytes(
                "AccessControl: account 0x3af3e85f4f97de7ad0f000b724fb77fe5ffc024b is missing role 0x0000000000000000000000000000000000000000000000000000000000000000"
            )
        );
        vm.prank(OPERATOR);
        usdg.upgradeToAndCall(address(usdg), "");
    }

    function test_UsdgPermitUsesGlobalDollarVersionOne() public {
        bytes32 expected = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
                keccak256("Global Dollar"),
                keccak256("1"),
                block.chainid,
                address(usdg)
            )
        );
        assertEq(usdg.DOMAIN_SEPARATOR(), expected);
        (bool hasDomain, bytes memory err) = address(usdg).staticcall(abi.encodeWithSignature("eip712Domain()"));
        assertFalse(hasDomain);
        assertEq(err, abi.encodeWithSelector(IUSDG.FacetNotFound.selector));

        Vm.Wallet memory owner = vm.createWallet("owner");
        address spender = makeAddr("spender");
        vm.prank(morpho);
        usdg.transfer(owner.addr, 1_000e6);
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(owner, spender, 400e6, block.timestamp);
        usdg.permit(owner.addr, spender, 400e6, block.timestamp, v, r, s);
        assertEq(usdg.allowance(owner.addr, spender), 400e6);
        assertEq(usdg.nonces(owner.addr), 1);

        vm.prank(spender);
        usdg.transferFrom(owner.addr, spender, 400e6);
        assertEq(usdg.balanceOf(spender), 400e6);
    }

    function test_UsdgPauseStopsTransfersApprovalsAndPermits() public {
        assertTrue(usdg.hasRole(PAUSE_ROLE, OPERATOR));
        Vm.Wallet memory owner = vm.createWallet("owner");
        vm.prank(morpho);
        usdg.transfer(owner.addr, 1_000e6);
        vm.prank(owner.addr);
        usdg.approve(address(this), 1e6);
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(owner, address(this), 1e6, block.timestamp);

        vm.prank(OPERATOR);
        usdg.pause();
        assertTrue(usdg.paused());

        vm.expectRevert(IUSDG.ContractPaused.selector);
        vm.prank(owner.addr);
        usdg.transfer(address(this), 1e6);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        usdg.transferFrom(owner.addr, address(this), 1e6);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        vm.prank(owner.addr);
        usdg.approve(address(this), 1e6);
        vm.expectRevert(IUSDG.ContractPaused.selector);
        usdg.permit(owner.addr, address(this), 1e6, block.timestamp, v, r, s);

        vm.prank(OPERATOR);
        usdg.unpause();
        vm.prank(owner.addr);
        usdg.transfer(address(this), 1e6);
    }

    function test_UsdgFrozenAddressCannotSendOrReceive() public {
        assertTrue(usdg.hasRole(ASSET_PROTECTION_ROLE, OPERATOR));
        address holder = makeAddr("holder");
        vm.prank(morpho);
        usdg.transfer(holder, 1_000e6);
        vm.prank(holder);
        usdg.approve(address(this), 1_000e6);

        vm.prank(OPERATOR);
        usdg.freeze(holder);
        assertTrue(usdg.isFrozen(holder));

        vm.expectRevert(IUSDG.AddressFrozen.selector);
        vm.prank(holder);
        usdg.transfer(address(this), 1e6);
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        usdg.transferFrom(holder, address(this), 1e6);
        vm.expectRevert(IUSDG.AddressFrozen.selector);
        vm.prank(morpho);
        usdg.transfer(holder, 1e6);
    }

    function _signPermit(Vm.Wallet memory owner, address spender, uint256 value, uint256 deadline)
        internal
        view
        returns (uint8, bytes32, bytes32)
    {
        bytes32 structHash =
            keccak256(abi.encode(PERMIT_TYPEHASH, owner.addr, spender, value, usdg.nonces(owner.addr), deadline));
        return vm.sign(owner, keccak256(abi.encodePacked("\x19\x01", usdg.DOMAIN_SEPARATOR(), structHash)));
    }
}
