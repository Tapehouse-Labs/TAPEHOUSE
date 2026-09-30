// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Script} from "forge-std/Script.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";

/// @notice Deploys the USDG supply vault for the chain's USDG and owner, from the registry at `REGISTRY`
/// (`deployments/<chainId>.json` by default). Registering the address is `Register`'s job, after the broadcast.
contract DeploySupplyVault is Script {
    function run() external returns (SupplyVault vault) {
        string memory json = vm.readFile(registry());
        IUSDG usdg = IUSDG(vm.parseJsonAddress(json, ".tokens.USDG"));
        address owner = vm.parseJsonAddress(json, ".tapehouse.Owner");
        vm.startBroadcast();
        vault = new SupplyVault(usdg, owner, SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        vm.stopBroadcast();
    }

    /// @notice The registry the vault is configured from.
    function registry() public view returns (string memory) {
        return vm.envOr("REGISTRY", string.concat("../deployments/", vm.toString(block.chainid), ".json"));
    }
}
