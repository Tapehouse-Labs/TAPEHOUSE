// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity ^0.8.12;

import {Test} from "forge-std/Test.sol";
import {IUniswapV3Factory} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Factory.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {UniswapV3Seeder} from "../src/UniswapV3Seeder.sol";

contract Token {
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        return true;
    }
}

contract UniswapV3SeederTest is Test {
    uint160 internal constant SQRT_PRICE = 79228162514264337593543950336;

    IUniswapV3Factory internal factory;
    UniswapV3Seeder internal seeder;
    IUniswapV3Pool internal pool;
    Token internal token0;
    Token internal token1;
    address internal lp = makeAddr("lp");

    function setUp() public {
        factory = IUniswapV3Factory(deployCode("UniswapV3Factory.sol:UniswapV3Factory"));
        seeder = new UniswapV3Seeder();
        Token a = new Token();
        Token b = new Token();
        (token0, token1) = address(a) < address(b) ? (a, b) : (b, a);
        pool = IUniswapV3Pool(factory.createPool(address(token0), address(token1), 3000));
        pool.initialize(SQRT_PRICE);
        pool.increaseObservationCardinalityNext(10);
        token0.mint(lp, 1_000e18);
        token1.mint(lp, 1_000e18);
        vm.startPrank(lp);
        token0.approve(address(seeder), type(uint256).max);
        token1.approve(address(seeder), type(uint256).max);
        vm.stopPrank();
    }

    function test_SeedsTheLiquidityForTheCallerWhoPaysIt() public {
        vm.warp(block.timestamp + 10);
        vm.prank(lp);
        (uint256 amount0, uint256 amount1) = seeder.seed(pool, -600, 600, 1e18);
        assertGt(amount0, 0);
        assertGt(amount1, 0);
        assertEq(token0.balanceOf(lp), 1_000e18 - amount0);
        assertEq(token1.balanceOf(lp), 1_000e18 - amount1);
        assertEq(pool.liquidity(), 1e18);
        (uint128 liquidity,,,,) = pool.positions(keccak256(abi.encodePacked(lp, int24(-600), int24(600))));
        assertEq(liquidity, 1e18);
        (,,, uint16 cardinality,,,) = pool.slot0();
        assertEq(cardinality, 10);
    }

    function test_OnlyThePoolBeingSeededMayCallBack() public {
        vm.expectRevert(abi.encodeWithSelector(UniswapV3Seeder.NotThePool.selector, address(this)));
        seeder.uniswapV3MintCallback(1, 1, "");
    }

    function test_AnUnfundedCallerCannotSeed() public {
        address broke = makeAddr("broke");
        vm.prank(broke);
        vm.expectRevert();
        seeder.seed(pool, -600, 600, 1e18);
    }
}
