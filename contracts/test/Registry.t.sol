// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {VmSafe} from "forge-std/Vm.sol";

contract RegistryTest is Test {
    function test_EveryEntryIsChecksummedAndFilesAreNamedAfterTheirChain() public {
        vm.pauseGasMetering();
        VmSafe.DirEntry[] memory files = vm.readDir("../deployments");
        uint256 checked;
        for (uint256 i; i < files.length; ++i) {
            string memory path = files[i].path;
            if (!vm.contains(path, ".json")) continue;
            string memory json = vm.readFile(path);
            string memory chainId = vm.toString(vm.parseJsonUint(json, ".chainId"));
            assertTrue(vm.contains(path, string.concat("/", chainId, ".json")), path);
            string[] memory groups = vm.parseJsonKeys(json, "$");
            for (uint256 g; g < groups.length; ++g) {
                bytes32 group = keccak256(bytes(groups[g]));
                if (group != keccak256("chainId") && group != keccak256("chainlinkRegularHours")) {
                    _check(json, string.concat(".", groups[g]));
                }
            }
            ++checked;
        }
        assertGt(checked, 0, "no registry files");
    }

    function keysOf(string memory json, string memory key) external view returns (string[] memory) {
        return vm.parseJsonKeys(json, key);
    }

    /// @dev Every address under `group`, however deep the groups nest, is not zero and written as its checksum.
    function _check(string memory json, string memory group) internal view {
        string[] memory names = this.keysOf(json, group);
        for (uint256 n; n < names.length; ++n) {
            string memory key = string.concat(group, ".", names[n]);
            try this.keysOf(json, key) returns (string[] memory) {
                _check(json, key);
            } catch {
                address entry = vm.parseJsonAddress(json, key);
                assertTrue(entry != address(0), key);
                assertEq(vm.parseJsonString(json, key), vm.toString(entry), key);
            }
        }
    }
}
