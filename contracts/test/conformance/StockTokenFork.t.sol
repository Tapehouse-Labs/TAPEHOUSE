// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test, Vm} from "forge-std/Test.sol";
import {IStockToken, IStockTokenRegistry} from "./Interfaces.sol";

contract StockTokenForkTest is Test {
    struct RpcLog {
        address emitter;
        bytes32 blockHash;
        bytes blockNumber;
        bytes blockTimestamp;
        bytes data;
        bytes logIndex;
        bool removed;
        bytes32[] topics;
        bytes32 transactionHash;
        bytes transactionIndex;
    }

    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    string internal constant ROBINHOOD_BLOCK_HEX = "0x42aeec9";
    bytes32 internal constant BEACON_SLOT = 0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50;
    address internal constant STOCK_IMPLEMENTATION = 0xb35490d6f9163DE4F80d88dc75c3516eb64C5aE2;

    bytes32 internal constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 internal constant BEACON_UPGRADER_ROLE = keccak256("BEACON_UPGRADER_ROLE");
    bytes32 internal constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 internal constant TOKEN_PAUSER_ROLE = keccak256("TOKEN_PAUSER_ROLE");
    bytes32 internal constant BLOCKER_ROLE = keccak256("BLOCKER_ROLE");
    bytes32 internal constant ADMIN_BURNER_ROLE = keccak256("ADMIN_BURNER_ROLE");
    bytes32 internal constant MULTIPLIER_UPDATER_ROLE = keccak256("MULTIPLIER_UPDATER_ROLE");
    bytes32 internal constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 internal constant BURNER_ROLE = keccak256("BURNER_ROLE");
    bytes32 internal constant ORACLE_PAUSER_ROLE = keccak256("ORACLE_PAUSER_ROLE");
    bytes32 internal constant TOKEN_DEPLOYER_ROLE = keccak256("TOKEN_DEPLOYER_ROLE");
    bytes32 internal constant FACTORY_UPGRADER_ROLE = keccak256("FACTORY_UPGRADER_ROLE");
    bytes32 internal constant METADATA_UPDATER_ROLE = keccak256("METADATA_UPDATER_ROLE");

    address internal constant ADMIN = 0xD6f8378F8e440c65F8382F5f2728c78DfD55B66d;
    address internal constant BEACON_UPGRADER = 0xCd8C6182e7C6Ca3B5156D6a90a67719d7e2Be094;
    address internal constant PAUSER = 0xe7BCB188254Bc6eBBfF63014DfED4cD4A024F22A;
    address internal constant TOKEN_PAUSER = 0xFCcF56B674113d9C4eb0F9B3370930ceD9E6Ab23;
    address internal constant BLOCKER = 0x913cA87347391218e5De2C17c5A0AEba8B0b28fD;
    address internal constant ADMIN_BURNER = 0x957B6de6525C63349f7619743Ef1E0ad93cd74D4;
    address internal constant MULTIPLIER_UPDATER = 0x92905e8d0e2301BA143215B8D86D63fFD4188143;
    address internal constant MINTER = 0x2b94105fFf37630f98e1f24811daD588FC5C3A87;
    address internal constant BURNER = 0x6E40B50A40C1db42A85a0E8fe8FF7d9CbFc2D8C1;
    address internal constant ORACLE_PAUSER = 0x7369d100c00F28E45D779ac9d4b1c7afa61e4aBC;
    address internal constant TOKEN_DEPLOYER = 0x5516B3451d4d6C9f63353Fe7Bc9537477ECCE000;
    address internal constant FACTORY_UPGRADER = 0x697e774d60c1a3769f2eD0b919AAcf17be0ae553;
    address internal constant METADATA_UPDATER = 0xcba16C2b9048AF033c5b34E43dd1D47D1358524A;

    address internal constant SANCTIONED = 0x910Cbd523D972eb0a6f4cAe4618aD62622b39DbF;

    uint256 internal constant SPY_UPDATE_BLOCK = 65_779_981;
    bytes32 internal constant SPY_UPDATE_TX = 0x2fe45ab24d1b3fa87883f8b08daf29dae969c5b43d0afe9ccca299c11a641025;
    uint256 internal constant SPY_LAST_BLOCK_BEFORE_EFFECTIVE = 65_785_790;
    uint256 internal constant SPY_OLD_MULTIPLIER = 1e18;
    uint256 internal constant SPY_NEW_MULTIPLIER = 1_001_717_991_187_472_003;
    uint256 internal constant SPY_EFFECTIVE_AT = 1_789_690_233;

    bytes32 internal constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)");

    IStockTokenRegistry internal registry;
    IStockToken internal nvda;
    IStockToken internal spy;
    address internal weth;

    function setUp() public {
        string memory json = vm.readFile("../deployments/4663.json");
        registry = IStockTokenRegistry(vm.parseJsonAddress(json, ".stockTokens.Registry"));
        nvda = IStockToken(vm.parseJsonAddress(json, ".tokens.NVDA"));
        spy = IStockToken(vm.parseJsonAddress(json, ".tokens.SPY"));
        weth = vm.parseJsonAddress(json, ".tokens.WETH");
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
    }

    function test_EveryStockTokenIsABeaconProxyOfTheRegistry() public view {
        assertEq(address(uint160(uint256(vm.load(address(nvda), BEACON_SLOT)))), address(registry));
        assertEq(address(uint160(uint256(vm.load(address(spy), BEACON_SLOT)))), address(registry));
        assertEq(registry.implementation(), STOCK_IMPLEMENTATION);
        assertEq(nvda.ACCESS_CONTROLLED_REGISTRY(), address(registry));
        assertEq(spy.ACCESS_CONTROLLED_REGISTRY(), address(registry));
        assertEq(nvda.decimals(), 18);
    }

    function test_RolesLiveInTheRegistryWhichIsNotEnumerable() public view {
        (bool tokenHasRoles,) = address(nvda).staticcall(abi.encodeWithSignature("hasRole(bytes32,address)", 0, ADMIN));
        assertFalse(tokenHasRoles);
        (bool enumerable,) = address(registry).staticcall(abi.encodeWithSignature("getRoleMemberCount(bytes32)", 0));
        assertFalse(enumerable);
    }

    function test_IssuerRoleHoldersAtThePinnedBlock() public view {
        (bytes32[] memory roles, address[] memory holders) = _issuerRoles();
        for (uint256 i; i < roles.length; ++i) {
            assertTrue(registry.hasRole(roles[i], holders[i]), vm.toString(roles[i]));
            assertEq(registry.getRoleAdmin(roles[i]), DEFAULT_ADMIN_ROLE, vm.toString(roles[i]));
            assertEq(holders[i].code.length, 0, vm.toString(holders[i]));
        }
    }

    function test_RoleHistoryHasNoOtherHolder() public {
        assertEq(vm.parseUint(ROBINHOOD_BLOCK_HEX), ROBINHOOD_BLOCK);
        bytes32 granted = keccak256("RoleGranted(bytes32,address,address)");
        bytes32 revoked = keccak256("RoleRevoked(bytes32,address,address)");
        string memory filter = string.concat(
            '[{"fromBlock":"0x0","toBlock":"',
            ROBINHOOD_BLOCK_HEX,
            '","address":"',
            vm.toString(address(registry)),
            '","topics":[["',
            vm.toString(granted),
            '","',
            vm.toString(revoked),
            '"]]}]'
        );
        RpcLog[] memory logs = abi.decode(vm.rpc("robinhood-logs", "eth_getLogs", filter), (RpcLog[]));

        bytes32[] memory roles = new bytes32[](logs.length);
        address[] memory accounts = new address[](logs.length);
        uint256 live;
        for (uint256 i; i < logs.length; ++i) {
            bytes32 role = logs[i].topics[1];
            address account = address(uint160(uint256(logs[i].topics[2])));
            uint256 k;
            while (k < live && (roles[k] != role || accounts[k] != account)) ++k;
            if (logs[i].topics[0] == granted && k == live) {
                (roles[live], accounts[live]) = (role, account);
                ++live;
            } else if (logs[i].topics[0] == revoked && k < live) {
                --live;
                (roles[k], accounts[k]) = (roles[live], accounts[live]);
            }
        }

        (bytes32[] memory expectedRoles, address[] memory holders) = _issuerRoles();
        assertEq(live, expectedRoles.length);
        for (uint256 k; k < live; ++k) {
            assertTrue(registry.hasRole(roles[k], accounts[k]), vm.toString(accounts[k]));
            bool expected;
            for (uint256 j; j < expectedRoles.length; ++j) {
                if (expectedRoles[j] == roles[k] && holders[j] == accounts[k]) expected = true;
            }
            assertTrue(expected, vm.toString(accounts[k]));
        }
    }

    function test_TransferToAContractNeedsNoAllowlist() public {
        address holder = makeAddr("holder");
        deal(address(nvda), holder, 5e18);
        vm.expectEmit(address(nvda));
        emit IStockToken.TransferWithScaledUI(holder, weth, 1e18, nvda.uiMultiplier());
        vm.prank(holder);
        assertTrue(nvda.transfer(weth, 1e18));
        assertEq(nvda.balanceOf(holder), 4e18);
        assertEq(nvda.balanceOf(weth), 1e18);
    }

    function test_PermitFromAFreshKeySetsAllowanceAndBumpsNonce() public {
        Vm.Wallet memory owner = vm.createWallet("owner");
        address spender = makeAddr("spender");
        deal(address(nvda), owner.addr, 5e18);

        (, string memory name, string memory version, uint256 chainId, address verifying,,) = nvda.eip712Domain();
        assertEq(name, nvda.name());
        assertEq(name, unicode"NVIDIA • Robinhood Token");
        assertEq(version, "1");
        assertEq(chainId, 4663);
        assertEq(verifying, address(nvda));

        (uint8 v, bytes32 r, bytes32 s) = _signPermit(nvda, owner, spender, 2e18, block.timestamp);
        nvda.permit(owner.addr, spender, 2e18, block.timestamp, v, r, s);
        assertEq(nvda.allowance(owner.addr, spender), 2e18);
        assertEq(nvda.nonces(owner.addr), 1);

        vm.prank(spender);
        assertTrue(nvda.transferFrom(owner.addr, spender, 2e18));
        assertEq(nvda.balanceOf(spender), 2e18);
    }

    function test_GlobalPauseStopsEveryTokenButNotAdminBurn() public {
        Vm.Wallet memory owner = vm.createWallet("owner");
        deal(address(nvda), owner.addr, 5e18);
        vm.prank(owner.addr);
        nvda.approve(address(this), 1e18);
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(nvda, owner, address(this), 1e18, block.timestamp);

        vm.prank(PAUSER);
        registry.pause();
        assertTrue(nvda.paused());
        assertTrue(spy.paused());
        assertFalse(nvda.tokenPaused());
        _assertEveryPathReverts(nvda, owner.addr, abi.encodeWithSelector(IStockToken.IsPaused.selector), v, r, s);

        vm.prank(ADMIN_BURNER);
        nvda.adminBurn(owner.addr, 1e18);
        assertEq(nvda.balanceOf(owner.addr), 4e18);

        vm.prank(PAUSER);
        registry.unpause();
        vm.prank(owner.addr);
        assertTrue(nvda.transfer(address(this), 1e18));
    }

    function test_TokenPauseStopsOnlyThatToken() public {
        Vm.Wallet memory owner = vm.createWallet("owner");
        deal(address(nvda), owner.addr, 5e18);
        deal(address(spy), owner.addr, 5e18);
        vm.prank(owner.addr);
        nvda.approve(address(this), 1e18);
        (uint8 v, bytes32 r, bytes32 s) = _signPermit(nvda, owner, address(this), 1e18, block.timestamp);

        vm.prank(TOKEN_PAUSER);
        nvda.pause();
        assertTrue(nvda.tokenPaused());
        assertTrue(nvda.paused());
        assertFalse(registry.paused());
        assertFalse(spy.paused());
        _assertEveryPathReverts(nvda, owner.addr, abi.encodeWithSelector(IStockToken.IsPaused.selector), v, r, s);
        vm.expectRevert(IStockToken.IsPaused.selector);
        vm.prank(MULTIPLIER_UPDATER);
        nvda.updateMultiplier(2e18, block.timestamp + 1 days);

        vm.prank(owner.addr);
        assertTrue(spy.transfer(address(this), 1e18));
        vm.prank(ADMIN_BURNER);
        nvda.adminBurn(owner.addr, 1e18);
        assertEq(nvda.balanceOf(owner.addr), 4e18);

        vm.prank(TOKEN_PAUSER);
        nvda.unpause();
        vm.prank(owner.addr);
        assertTrue(nvda.transfer(address(this), 1e18));
    }

    function test_BlocklistCoversSenderRecipientAndCaller() public {
        assertTrue(registry.isBlocked(SANCTIONED));
        address holder = makeAddr("holder");
        address blocked = makeAddr("blocked");
        deal(address(nvda), holder, 5e18);
        deal(address(nvda), blocked, 5e18);
        vm.prank(holder);
        nvda.approve(blocked, 1e18);
        vm.prank(blocked);
        nvda.approve(holder, 1e18);

        address[] memory accounts = new address[](1);
        accounts[0] = blocked;
        vm.prank(BLOCKER);
        registry.blockAccounts(accounts);
        bytes memory err = abi.encodeWithSelector(IStockToken.Blocked.selector, blocked);

        vm.expectRevert(err);
        vm.prank(blocked);
        nvda.transfer(holder, 1e18);
        vm.expectRevert(err);
        vm.prank(holder);
        nvda.transfer(blocked, 1e18);
        vm.expectRevert(err);
        vm.prank(blocked);
        nvda.transferFrom(holder, address(this), 1e18);
        vm.expectRevert(err);
        vm.prank(holder);
        nvda.transferFrom(blocked, address(this), 1e18);
        vm.expectRevert(err);
        vm.prank(holder);
        nvda.approve(blocked, 1e18);

        vm.prank(ADMIN_BURNER);
        nvda.adminBurn(blocked, 5e18);
        assertEq(nvda.balanceOf(blocked), 0);
    }

    function test_AdminBurnTakesFromAnyAddressWithoutApproval() public {
        address holder = makeAddr("holder");
        deal(address(nvda), holder, 5e18, true);
        uint256 supply = nvda.totalSupply();
        assertEq(nvda.allowance(holder, ADMIN_BURNER), 0);

        vm.expectRevert(
            abi.encodeWithSelector(IStockToken.AccessControlUnauthorizedAccount.selector, ADMIN, ADMIN_BURNER_ROLE)
        );
        vm.prank(ADMIN);
        nvda.adminBurn(holder, 1e18);

        vm.prank(ADMIN_BURNER);
        nvda.adminBurn(holder, 2e18);
        assertEq(nvda.balanceOf(holder), 3e18);
        assertEq(nvda.totalSupply(), supply - 2e18);
    }

    function test_SpyMultiplierReplaysTheRealChange() public {
        vm.createSelectFork("robinhood", SPY_UPDATE_BLOCK - 1);
        assertEq(spy.uiMultiplier(), SPY_OLD_MULTIPLIER);
        assertEq(spy.effectiveAt(), 0);
        vm.expectEmit(address(spy));
        emit IStockToken.UIMultiplierUpdated(SPY_OLD_MULTIPLIER, SPY_NEW_MULTIPLIER, SPY_EFFECTIVE_AT);
        vm.transact(SPY_UPDATE_TX);

        vm.createSelectFork("robinhood", SPY_LAST_BLOCK_BEFORE_EFFECTIVE);
        assertEq(vm.getBlockTimestamp(), SPY_EFFECTIVE_AT - 1);
        assertEq(spy.uiMultiplier(), SPY_OLD_MULTIPLIER);
        assertEq(spy.newUIMultiplier(), SPY_NEW_MULTIPLIER);
        assertEq(spy.effectiveAt(), SPY_EFFECTIVE_AT);
        uint256 supply = spy.totalSupply();
        assertEq(spy.totalSupplyUI(), supply);

        vm.warp(SPY_EFFECTIVE_AT);
        assertEq(spy.uiMultiplier(), SPY_NEW_MULTIPLIER);
        assertEq(spy.totalSupply(), supply);
        assertEq(spy.totalSupplyUI(), supply * SPY_NEW_MULTIPLIER / 1e18);

        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
        assertEq(spy.uiMultiplier(), SPY_NEW_MULTIPLIER);
        assertEq(spy.effectiveAt(), SPY_EFFECTIVE_AT);
    }

    function test_SpyBalanceOfUIScalesAndRoundsDown() public {
        address holder = makeAddr("holder");
        deal(address(spy), holder, 292);
        assertEq(spy.balanceOfUI(holder), 292);
        deal(address(spy), holder, 1e18);
        assertEq(spy.balanceOfUI(holder), SPY_NEW_MULTIPLIER);
    }

    function test_ScheduledMultiplierIsSilentlyReplaceable() public {
        vm.createSelectFork("robinhood", SPY_UPDATE_BLOCK);
        assertEq(spy.newUIMultiplier(), SPY_NEW_MULTIPLIER);
        assertEq(spy.effectiveAt(), SPY_EFFECTIVE_AT);
        assertLt(block.timestamp, SPY_EFFECTIVE_AT);
        vm.recordLogs();
        vm.prank(MULTIPLIER_UPDATER);
        spy.updateMultiplier(2e18, SPY_EFFECTIVE_AT + 1);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        assertEq(logs.length, 1);
        assertEq(logs[0].topics[0], IStockToken.UIMultiplierUpdated.selector);
        assertEq(abi.decode(logs[0].data, (uint256)), SPY_OLD_MULTIPLIER);
        assertEq(spy.newUIMultiplier(), 2e18);
        assertEq(spy.effectiveAt(), SPY_EFFECTIVE_AT + 1);
    }

    function test_Erc8056InterfaceIdsWithoutConversion() public view {
        assertTrue(spy.supportsInterface(0xa60bf13d));
        assertTrue(spy.supportsInterface(0x4bd27648));
        assertTrue(spy.supportsInterface(0xd890fd71));
        assertFalse(spy.supportsInterface(0x57854fc3));
    }

    function _assertEveryPathReverts(IStockToken token, address owner, bytes memory err, uint8 v, bytes32 r, bytes32 s)
        internal
    {
        vm.expectRevert(err);
        vm.prank(owner);
        token.transfer(address(this), 1e18);
        vm.expectRevert(err);
        token.transferFrom(owner, address(this), 1e18);
        vm.expectRevert(err);
        vm.prank(owner);
        token.approve(address(this), 1e18);
        vm.expectRevert(err);
        token.permit(owner, address(this), 1e18, block.timestamp, v, r, s);
    }

    function _signPermit(IStockToken token, Vm.Wallet memory owner, address spender, uint256 value, uint256 deadline)
        internal
        view
        returns (uint8, bytes32, bytes32)
    {
        bytes32 structHash =
            keccak256(abi.encode(PERMIT_TYPEHASH, owner.addr, spender, value, token.nonces(owner.addr), deadline));
        return vm.sign(owner, keccak256(abi.encodePacked("\x19\x01", token.DOMAIN_SEPARATOR(), structHash)));
    }

    function _issuerRoles() internal pure returns (bytes32[] memory roles, address[] memory holders) {
        roles = new bytes32[](13);
        holders = new address[](13);
        (roles[0], holders[0]) = (DEFAULT_ADMIN_ROLE, ADMIN);
        (roles[1], holders[1]) = (BEACON_UPGRADER_ROLE, BEACON_UPGRADER);
        (roles[2], holders[2]) = (PAUSER_ROLE, PAUSER);
        (roles[3], holders[3]) = (TOKEN_PAUSER_ROLE, TOKEN_PAUSER);
        (roles[4], holders[4]) = (BLOCKER_ROLE, BLOCKER);
        (roles[5], holders[5]) = (ADMIN_BURNER_ROLE, ADMIN_BURNER);
        (roles[6], holders[6]) = (MULTIPLIER_UPDATER_ROLE, MULTIPLIER_UPDATER);
        (roles[7], holders[7]) = (MINTER_ROLE, MINTER);
        (roles[8], holders[8]) = (BURNER_ROLE, BURNER);
        (roles[9], holders[9]) = (ORACLE_PAUSER_ROLE, ORACLE_PAUSER);
        (roles[10], holders[10]) = (TOKEN_DEPLOYER_ROLE, TOKEN_DEPLOYER);
        (roles[11], holders[11]) = (FACTORY_UPGRADER_ROLE, FACTORY_UPGRADER);
        (roles[12], holders[12]) = (METADATA_UPDATER_ROLE, METADATA_UPDATER);
    }
}
