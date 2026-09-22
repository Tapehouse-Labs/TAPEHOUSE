// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {ArbGasInfo, ArbSys} from "./Interfaces.sol";

contract ArbitrumPrecompilesForkTest is Test {
    ArbSys internal constant ARB_SYS = ArbSys(address(0x64));
    address internal constant ARB_GAS_INFO = address(0x6c);
    uint256 internal constant PROBE_GAS = 100_000;

    function test_RobinhoodForkGaps() public {
        _assertForkGaps("robinhood", 69_922_505, 26_035_252, 1_790_106_999);
    }

    function test_ArbitrumForkGaps() public {
        _assertForkGaps("arbitrum", 507_888_520, 26_035_253, 1_790_106_999);
    }

    function _assertForkGaps(string memory network, uint256 l2Block, uint256 l1Block, uint256 timestamp) internal {
        vm.createSelectFork(network, l2Block);
        assertEq(ARB_SYS.arbBlockNumber(), l2Block, "arbBlockNumber");
        assertEq(block.number, l1Block, "block.number");
        assertEq(block.timestamp, timestamp, "block.timestamp");

        (bool ok, bytes memory data) =
            address(ARB_SYS).staticcall{gas: PROBE_GAS}(abi.encodeCall(ArbSys.arbOSVersion, ()));
        assertFalse(ok, "arbOSVersion");
        assertEq(data.length, 0, "arbOSVersion");

        (ok, data) = ARB_GAS_INFO.staticcall{gas: PROBE_GAS}(abi.encodeCall(ArbGasInfo.getPricesInWei, ()));
        assertFalse(ok, "getPricesInWei");
        assertEq(data.length, 0, "getPricesInWei");
        assertEq(ARB_GAS_INFO.code, hex"fe", "ArbGasInfo code");
    }
}
