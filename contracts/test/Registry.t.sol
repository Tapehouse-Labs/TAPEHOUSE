// SPDX-License-Identifier: UNLICENSED
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
                if (keccak256(bytes(groups[g])) == keccak256("chainId")) continue;
                string memory group = string.concat(".", groups[g]);
                string[] memory names = vm.parseJsonKeys(json, group);
                for (uint256 n; n < names.length; ++n) {
                    string memory key = string.concat(group, ".", names[n]);
                    address entry = vm.parseJsonAddress(json, key);
                    assertTrue(entry != address(0), key);
                    assertEq(vm.parseJsonString(json, key), vm.toString(entry), key);
                }
            }
            ++checked;
        }
        assertGt(checked, 0, "no registry files");
    }
}
