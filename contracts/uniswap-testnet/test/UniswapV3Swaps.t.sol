// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity ^0.8.12;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3Factory} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Factory.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {UniswapV3Seeder} from "../src/UniswapV3Seeder.sol";
import {UniswapV3Router} from "../src/UniswapV3Router.sol";
import {UniswapV3Quoter} from "../src/UniswapV3Quoter.sol";
import {Token} from "./UniswapV3Seeder.t.sol";

contract UniswapV3SwapsTest is Test {
    uint160 internal constant SQRT_PRICE = 79228162514264337593543950336;

    IUniswapV3Pool internal pool;
    UniswapV3Router internal router;
    UniswapV3Quoter internal quoter;
    Token internal token0;
    Token internal token1;
    address internal trader = makeAddr("trader");

    function setUp() public {
        IUniswapV3Factory factory = IUniswapV3Factory(deployCode("UniswapV3Factory.sol:UniswapV3Factory"));
        UniswapV3Seeder seeder = new UniswapV3Seeder();
        router = new UniswapV3Router(address(factory));
        quoter = new UniswapV3Quoter(address(factory));
        Token a = new Token();
        Token b = new Token();
        (token0, token1) = address(a) < address(b) ? (a, b) : (b, a);
        pool = IUniswapV3Pool(factory.createPool(address(token0), address(token1), 3000));
        pool.initialize(SQRT_PRICE);
        token0.mint(address(this), 1_000_000e18);
        token1.mint(address(this), 1_000_000e18);
        token0.approve(address(seeder), type(uint256).max);
        token1.approve(address(seeder), type(uint256).max);
        vm.warp(block.timestamp + 10);
        seeder.seed(pool, -6000, 6000, 1_000e18);
        token0.mint(trader, 100e18);
        token1.mint(trader, 100e18);
        vm.startPrank(trader);
        token0.approve(address(router), type(uint256).max);
        token1.approve(address(router), type(uint256).max);
        vm.stopPrank();
    }

    function test_ExactInputPaysTheQuoteInEitherDirection() public {
        (uint256 quoted,,,) = quoter.quoteExactInputSingle(
            UniswapV3Quoter.QuoteExactInputSingleParams(address(token0), address(token1), 1e18, 3000, 0)
        );
        vm.prank(trader);
        uint256 out = router.exactInputSingle(
            UniswapV3Router.ExactInputSingleParams(address(token0), address(token1), 3000, trader, 1e18, quoted, 0)
        );
        assertEq(out, quoted);
        assertEq(token1.balanceOf(trader), 100e18 + out);
        assertEq(token0.balanceOf(trader), 99e18);

        (uint256 back,,,) = quoter.quoteExactInputSingle(
            UniswapV3Quoter.QuoteExactInputSingleParams(address(token1), address(token0), 1e18, 3000, 0)
        );
        vm.prank(trader);
        uint256 got = router.exactInputSingle(
            UniswapV3Router.ExactInputSingleParams(address(token1), address(token0), 3000, trader, 1e18, back, 0)
        );
        assertEq(got, back);
    }

    function test_ExactOutputCostsTheQuote() public {
        (uint256 cost, uint160 sqrtAfter,,) = quoter.quoteExactOutputSingle(
            UniswapV3Quoter.QuoteExactOutputSingleParams(address(token0), address(token1), 1e18, 3000, 0)
        );
        assertGt(cost, 1e18);
        assertGt(uint256(SQRT_PRICE), uint256(sqrtAfter));
        vm.prank(trader);
        uint256 paid = router.exactOutputSingle(
            UniswapV3Router.ExactOutputSingleParams(address(token0), address(token1), 3000, trader, 1e18, cost, 0)
        );
        assertEq(paid, cost);
        assertEq(token1.balanceOf(trader), 101e18);
    }

    function test_SlippageBoundsRevert() public {
        vm.startPrank(trader);
        vm.expectRevert();
        router.exactInputSingle(
            UniswapV3Router.ExactInputSingleParams(address(token0), address(token1), 3000, trader, 1e18, 2e18, 0)
        );
        vm.expectRevert();
        router.exactOutputSingle(
            UniswapV3Router.ExactOutputSingleParams(address(token0), address(token1), 3000, trader, 1e18, 1e18, 0)
        );
        vm.stopPrank();
    }

    function test_AnUnknownPoolOrCallbackReverts() public {
        vm.prank(trader);
        vm.expectRevert(UniswapV3Router.NoPool.selector);
        router.exactInputSingle(
            UniswapV3Router.ExactInputSingleParams(address(token0), address(token1), 500, trader, 1e18, 0, 0)
        );
        vm.expectRevert(abi.encodeWithSelector(UniswapV3Router.NotThePool.selector, address(this)));
        router.uniswapV3SwapCallback(1, 0, abi.encode(trader, address(token0), address(token1), uint24(3000)));
        vm.expectRevert(abi.encodeWithSelector(UniswapV3Quoter.NotThePool.selector, address(this)));
        quoter.uniswapV3SwapCallback(1, 0, abi.encode(address(pool)));
    }
}
