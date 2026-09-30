// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3PoolDerivedState} from "@uniswap/v3-core/contracts/interfaces/pool/IUniswapV3PoolDerivedState.sol";
import {MarginReference} from "./MarginReference.sol";

contract MarginReferenceForkTest is Test {
    uint256 internal constant BLOCK = 75_093_578;
    string internal constant VECTORS = "../stylus/contracts/margin/testdata/requirement-vectors.json";
    string internal constant PARAMETERS = "../stylus/contracts/margin/parameters.json";

    string internal json;
    string internal parameters;
    string[] internal names;

    function test_RequirementsAreTheReferenceRequirementsOnRobinhoodChain() public {
        vm.createSelectFork("robinhood", BLOCK);
        vm.pauseGasMetering();
        json = vm.readFile(VECTORS);
        parameters = vm.readFile(PARAMETERS);
        names = vm.parseJsonStringArray(json, ".symbols");
        assertEq(block.timestamp, vm.parseJsonUint(json, ".timestamp"));
        MarginReference pooled = _deploy(true);
        MarginReference unpooled = _deploy(false);
        string[] memory labels = abi.decode(vm.parseJson(json, ".vectors[*].label"), (string[]));
        assertEq(labels.length, 120);
        for (uint256 k; k < labels.length; ++k) {
            string memory path = string.concat(".vectors[", vm.toString(k), "]");
            bytes32 pools = keccak256(bytes(vm.parseJsonString(json, string.concat(path, ".pools"))));
            if (pools == keccak256("unread")) _unread();
            (uint256 requirement, uint8 missing) = (pools == keccak256("absent") ? unpooled : pooled)
            .requirement(
                _quantities(path),
                vm.parseJsonUintArray(json, ".prices"),
                uint64(vm.parseJsonUint(json, string.concat(path, ".horizon"))),
                vm.parseJsonBool(json, string.concat(path, ".spansClosure"))
            );
            vm.clearMockedCalls();
            assertEq(
                requirement, vm.parseUint(vm.parseJsonString(json, string.concat(path, ".requirement"))), labels[k]
            );
            assertEq(missing, vm.parseJsonUint(json, string.concat(path, ".missing")), labels[k]);
        }
    }

    function _deploy(bool withPools) internal returns (MarginReference) {
        MarginReference.Asset[] memory assets = new MarginReference.Asset[](names.length);
        for (uint256 i; i < names.length; ++i) {
            assets[i] = _asset(names[i], withPools);
        }
        (uint16[] memory floors, uint16[] memory correlations) = _correlations();
        return new MarginReference(
            assets,
            floors,
            correlations,
            bytes32(bytes(vm.parseJsonString(parameters, ".market"))),
            vm.parseJsonAddress(json, ".usdg"),
            vm.parseJsonAddress(json, ".weth"),
            vm.parseJsonAddress(json, ".ethUsd"),
            address(0),
            address(this)
        );
    }

    function _asset(string memory name, bool withPools) internal view returns (MarginReference.Asset memory) {
        uint256[] memory depth = vm.parseJsonUintArray(parameters, string.concat(".depth.", name, ".initial"));
        return MarginReference.Asset(
            bytes32(bytes(name)),
            uint32(vm.parseJsonUint(parameters, string.concat(".volatility.", name, ".floor"))),
            uint32(vm.parseJsonUint(parameters, string.concat(".volatility.", name, ".initial"))),
            uint32(vm.parseJsonUint(parameters, string.concat(".weekendGap.", name, ".floor"))),
            uint32(vm.parseJsonUint(parameters, string.concat(".weekendGap.", name, ".initial"))),
            uint32(depth[0]),
            uint32(depth[1]),
            withPools ? vm.parseJsonAddress(json, string.concat(".pools.", name, ".pool")) : address(0),
            vm.parseJsonAddress(json, string.concat(".tokens.", name))
        );
    }

    function _correlations() internal view returns (uint16[] memory floors, uint16[] memory values) {
        uint256 pairs = names.length * (names.length - 1) / 2;
        floors = new uint16[](pairs);
        values = new uint16[](pairs);
        uint256 k;
        for (uint256 i; i < names.length; ++i) {
            for (uint256 j = i + 1; j < names.length; ++j) {
                string memory pair = string.concat(".correlation['", names[i], "/", names[j], "']");
                floors[k] = uint16(vm.parseJsonUint(parameters, string.concat(pair, ".floor")));
                values[k++] = uint16(vm.parseJsonUint(parameters, string.concat(pair, ".initial")));
            }
        }
    }

    function _unread() internal {
        for (uint256 i; i < names.length; ++i) {
            vm.mockCallRevert(
                vm.parseJsonAddress(json, string.concat(".pools.", names[i], ".pool")),
                abi.encodeWithSelector(IUniswapV3PoolDerivedState.observe.selector),
                ""
            );
        }
    }

    function _quantities(string memory path) internal view returns (int256[] memory quantities) {
        string[] memory values = vm.parseJsonStringArray(json, string.concat(path, ".quantities"));
        quantities = new int256[](values.length);
        for (uint256 i; i < values.length; ++i) {
            quantities[i] = vm.parseInt(values[i]);
        }
    }
}
