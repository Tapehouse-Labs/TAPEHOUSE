// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {BandFeed} from "../src/BandFeed.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {BandDouble} from "./doubles/BandDouble.sol";

contract BandFeedForkTest is Test {
    uint256 internal constant ROBINHOOD_BLOCK = 69_922_505;
    uint256 internal constant NVDA_TWAP = 22_969_005_472;
    uint256 internal constant TSLA_TWAP = 37_876_115_313;

    function test_TwapsMatchTheIndependentComputationAtThePinnedBlock() public {
        string memory json = vm.readFile("../deployments/4663.json");
        vm.createSelectFork("robinhood", ROBINHOOD_BLOCK);
        _assertTwap(json, "NVDA", ".uniswapV3.NVDA_USDG_500", NVDA_TWAP);
        _assertTwap(json, "TSLA", ".uniswapV3.TSLA_USDG_3000", TSLA_TWAP);
    }

    function _assertTwap(string memory json, string memory name, string memory poolKey, uint256 expected) internal {
        BandDouble band = new BandDouble();
        bytes32 symbol = bytes32(bytes(name));
        address token = vm.parseJsonAddress(json, string.concat(".tokens.", name));
        band.setAsset(symbol, BandDouble.Asset(address(0), keccak256(bytes(name)), bytes32(0), token));
        band.setQuote(symbol, BandDouble.Quote(2, 1, uint64(expected), 50, uint64(expected), uint128(expected)));
        BandFeed feed = new BandFeed(
            IBand(address(band)),
            symbol,
            BandFeed.Side.Low,
            IUniswapV3Pool(vm.parseJsonAddress(json, poolKey)),
            string.concat(name, " / USD")
        );
        BandFeed.Band memory b = feed.latestBand();
        assertTrue(b.twapValid, name);
        assertApproxEqRel(b.twap, expected, 1e12, name);
        assertEq(b.premiumBps, 0, name);
    }
}
