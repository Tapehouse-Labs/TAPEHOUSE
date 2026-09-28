// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "forge-std/interfaces/IERC20.sol";
import {AggregatorV3Interface, IUniswapV3Factory} from "./conformance/Interfaces.sol";
import {BandFeed} from "../src/BandFeed.sol";

contract RegistryForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    uint256 internal constant ROBINHOOD_TESTNET_BLOCK = 125_558_861;
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
        address band = vm.parseJsonAddress(json, ".tapehouse.Band");
        assertEq(
            _bandAddress(band, "haltSigner()", ROBINHOOD_TESTNET_BLOCK),
            vm.parseJsonAddress(json, ".tapehouse.HaltSigner"),
            "tapehouse.HaltSigner"
        );
        assertEq(
            _bandAddress(band, "owner()", ROBINHOOD_TESTNET_BLOCK),
            vm.parseJsonAddress(json, ".tapehouse.Owner"),
            "tapehouse.Owner"
        );
        _assertBandFeeds(json);
    }

    function test_ArbitrumEntriesAreTheContractsTheyClaimToBe() public {
        string memory json = vm.readFile("../deployments/42161.json");
        vm.createSelectFork("arbitrum", ARBITRUM_BLOCK);
        assertEq(block.chainid, 42161);
        _assertFeeds(json);
        AggregatorV3Interface sequencer = AggregatorV3Interface(vm.parseJsonAddress(json, ".chainlinkSequencer.Uptime"));
        assertEq(sequencer.description(), "L2 Sequencer Uptime Status Feed");
        assertEq(sequencer.decimals(), 0);
        (, int256 status, uint256 startedAt,,) = sequencer.latestRoundData();
        assertEq(status, 0);
        assertGt(startedAt, 0);
        _assertBandFeeds(json);
    }

    function _assertBandFeeds(string memory json) internal view {
        if (!vm.keyExistsJson(json, ".bandFeeds")) return;
        address band = vm.parseJsonAddress(json, ".tapehouse.Band");
        string[] memory assets = vm.parseJsonKeys(json, ".bandFeeds");
        for (uint256 i; i < assets.length; ++i) {
            BandFeed feed = BandFeed(vm.parseJsonAddress(json, string.concat(".bandFeeds.", assets[i])));
            assertEq(address(feed.band()), band, assets[i]);
            assertEq(feed.symbol(), bytes32(bytes(assets[i])), assets[i]);
            assertEq(feed.decimals(), 8, assets[i]);
        }
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

    function _bandAddress(address band, string memory signature, uint256 blockNumber) internal returns (address) {
        bytes memory result = vm.rpc(
            "eth_call",
            string.concat(
                '[{"to":"',
                vm.toString(band),
                '","data":"',
                vm.toString(abi.encodeWithSignature(signature)),
                '"},"',
                _quantity(blockNumber),
                '"]'
            )
        );
        return abi.decode(result, (address));
    }

    function _quantity(uint256 value) internal pure returns (string memory) {
        bytes memory digits = "0123456789abcdef";
        bytes memory out = new bytes(64);
        uint256 length;
        do {
            out[63 - length++] = digits[value & 15];
            value >>= 4;
        } while (value != 0);
        bytes memory quantity = new bytes(length + 2);
        quantity[0] = "0";
        quantity[1] = "x";
        for (uint256 i; i < length; ++i) {
            quantity[i + 2] = out[64 - length + i];
        }
        return string(quantity);
    }
}
