// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";
import {AggregatorV3Interface, IUniswapV3Factory} from "./conformance/Interfaces.sol";

contract RegistryForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    uint256 internal constant ROBINHOOD_TESTNET_BLOCK = 124_738_159;
    uint256 internal constant ARBITRUM_BLOCK = 507_888_520;
    bytes32 internal constant BEACON_SLOT = 0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50;
    bytes4 internal constant STYLUS_ROOT_PREFIX = 0xeff00200;

    function test_RobinhoodEntriesAreTheContractsTheyClaimToBe() public {
        string memory json = vm.readFile("../deployments/4663.json");
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
        assertEq(block.chainid, 4663);
        _assertTokens(json);

        _assertFeeds(json);

        address registry = vm.parseJsonAddress(json, ".stockTokens.Registry");
        string[] memory symbols = vm.parseJsonKeys(json, ".tokens");
        for (uint256 i; i < symbols.length; ++i) {
            if (keccak256(bytes(symbols[i])) == keccak256("USDG") || keccak256(bytes(symbols[i])) == keccak256("WETH"))
            {
                continue;
            }
            address token = vm.parseJsonAddress(json, string.concat(".tokens.", symbols[i]));
            assertEq(address(uint160(uint256(vm.load(token, BEACON_SLOT)))), registry, symbols[i]);
        }
        assertGt(vm.parseJsonAddress(json, ".morpho.Blue").code.length, 0, "morpho.Blue");
        IUniswapV3Factory factory = IUniswapV3Factory(vm.parseJsonAddress(json, ".uniswapV3.Factory"));
        assertEq(factory.feeAmountTickSpacing(500), 10, "uniswapV3.Factory");
    }

    function test_RobinhoodTestnetEntriesAreTheContractsTheyClaimToBe() public {
        string memory json = vm.readFile("../deployments/46630.json");
        vm.createSelectFork("robinhood-testnet", ROBINHOOD_TESTNET_BLOCK);
        assertEq(block.chainid, 46630);
        _assertTokens(json);
        assertEq(bytes4(vm.parseJsonAddress(json, ".tapehouse.Band").code), STYLUS_ROOT_PREFIX, "tapehouse.Band");
    }

    function test_ArbitrumEntriesAreTheContractsTheyClaimToBe() public {
        string memory json = vm.readFile("../deployments/42161.json");
        vm.createSelectFork("arbitrum", ARBITRUM_BLOCK);
        assertEq(block.chainid, 42161);
        _assertFeeds(json);
    }

    function _assertFeeds(string memory json) internal view {
        string[] memory feeds = vm.parseJsonKeys(json, ".chainlink");
        for (uint256 i; i < feeds.length; ++i) {
            AggregatorV3Interface feed =
                AggregatorV3Interface(vm.parseJsonAddress(json, string.concat(".chainlink.", feeds[i])));
            assertTrue(vm.contains(feed.description(), vm.replace(feeds[i], "_USD", " / USD")), feeds[i]);
            assertEq(feed.decimals(), 8, feeds[i]);
            (, int256 answer,,,) = feed.latestRoundData();
            assertGt(answer, 0, feeds[i]);
        }
    }

    function _assertTokens(string memory json) internal view {
        string[] memory symbols = vm.parseJsonKeys(json, ".tokens");
        for (uint256 i; i < symbols.length; ++i) {
            address token = vm.parseJsonAddress(json, string.concat(".tokens.", symbols[i]));
            assertGt(token.code.length, 0, symbols[i]);
            (bool ok, bytes memory symbol) = token.staticcall(abi.encodeCall(IERC20.symbol, ()));
            assertTrue(ok, symbols[i]);
            assertEq(abi.decode(symbol, (string)), symbols[i], symbols[i]);
        }
    }
}
