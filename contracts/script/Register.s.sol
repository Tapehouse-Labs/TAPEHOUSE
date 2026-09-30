// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Script, VmSafe} from "forge-std/Script.sol";

/// @notice Records the USDG supply vault's latest successful deployment on this chain, read from Foundry's broadcast
/// files, newest first, in the registry at `REGISTRY` (`deployments/<chainId>.json` by default) under
/// `.tapehouse.SupplyVault`, checksummed.
contract Register is Script {
    error NoDeployment();

    function run() external returns (address vault) {
        VmSafe.BroadcastTxSummary[] memory txs =
            vm.getBroadcasts("SupplyVault", uint64(block.chainid), VmSafe.BroadcastTxType.Create);
        for (uint256 i; i < txs.length && vault == address(0); ++i) {
            if (txs[i].success) vault = txs[i].contractAddress;
        }
        if (vault == address(0) || vault.code.length == 0) revert NoDeployment();
        string memory path = vm.envOr("REGISTRY", string.concat("../deployments/", vm.toString(block.chainid), ".json"));
        vm.writeJson(vm.toString(vault), path, ".tapehouse.SupplyVault");
        vm.writeLine(path, "");
    }
}
