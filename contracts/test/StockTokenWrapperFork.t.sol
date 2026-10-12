// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {StockTokenWrapper, IWrappedStock} from "../src/StockTokenWrapper.sol";

contract StockTokenWrapperForkTest is Test {
    address internal constant TSLA = 0xC9f9c86933092BbbfFF3CCb4b105A4A94bf3Bd4E;
    address internal constant AMZN = 0x5884aD2f920c162CFBbACc88C9C51AA75eC09E02;
    uint256 internal constant ROBINHOOD_TESTNET_BLOCK = 133_737_048;

    function test_WrapsTheTestnetStockTokens() public {
        vm.createSelectFork("robinhood-testnet", ROBINHOOD_TESTNET_BLOCK);
        address[2] memory tokens = [TSLA, AMZN];
        for (uint256 i; i < tokens.length; ++i) {
            IWrappedStock stock = IWrappedStock(tokens[i]);
            StockTokenWrapper wrapper = new StockTokenWrapper(stock);
            assertEq(wrapper.symbol(), string.concat("w", stock.symbol()));
            assertEq(wrapper.uiMultiplier(), stock.uiMultiplier());
            assertEq(wrapper.newUIMultiplier(), stock.newUIMultiplier());
            assertEq(wrapper.effectiveAt(), stock.effectiveAt());
            assertEq(wrapper.oraclePaused(), stock.tokenPaused());
            assertEq(wrapper.paused(), stock.paused());
            assertEq(wrapper.ACCESS_CONTROLLED_REGISTRY(), stock.ACCESS_CONTROLLED_REGISTRY());

            address holder = makeAddr(stock.symbol());
            deal(address(stock), holder, 5e18);
            vm.startPrank(holder);
            stock.approve(address(wrapper), 5e18);
            wrapper.depositFor(holder, 5e18);
            assertEq(wrapper.balanceOf(holder), 5e18);
            wrapper.withdrawTo(holder, 2e18);
            vm.stopPrank();
            assertEq(stock.balanceOf(holder), 2e18);
            assertEq(wrapper.balanceOf(holder), 3e18);
        }
    }
}
