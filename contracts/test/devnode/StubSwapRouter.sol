// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IUniswapV3Factory, IV3SwapRouter} from "../../src/interfaces/IUniswapV3.sol";

interface IMintable {
    function mint(address to, uint256 value) external;
}

/// @notice A Uniswap v3 swap router, and its factory, that swaps a Stock Token and USDG at the price it is given, in
/// USD with 8 decimals a token, and reverts as the real router and pool do when a swap misses its limit or swaps
/// nothing. It mints what it pays out. Test code only.
contract StubSwapRouter is IV3SwapRouter, IUniswapV3Factory {
    using SafeERC20 for IERC20;

    address public immutable usdg;
    mapping(address token => uint256) public prices;

    constructor(address usdg_) {
        usdg = usdg_;
    }

    function setPrice(address token, uint256 price) external {
        prices[token] = price;
    }

    function factory() external view returns (address) {
        return address(this);
    }

    function getPool(address tokenA, address tokenB, uint24) external view returns (address) {
        return (tokenA == usdg ? prices[tokenB] : prices[tokenA]) == 0 ? address(0) : address(this);
    }

    function exactInputSingle(ExactInputSingleParams calldata params) external payable returns (uint256 amountOut) {
        require(params.amountIn != 0, "AS");
        amountOut = params.tokenIn == usdg
            ? params.amountIn * 1e20 / prices[params.tokenOut]
            : params.amountIn * prices[params.tokenIn] / 1e20;
        require(amountOut >= params.amountOutMinimum, "Too little received");
        _swap(params.tokenIn, params.amountIn, params.tokenOut, amountOut, params.recipient);
    }

    function exactOutputSingle(ExactOutputSingleParams calldata params) external payable returns (uint256 amountIn) {
        require(params.amountOut != 0, "AS");
        amountIn = params.tokenIn == usdg
            ? _divUp(params.amountOut * prices[params.tokenOut], 1e20)
            : _divUp(params.amountOut * 1e20, prices[params.tokenIn]);
        require(amountIn <= params.amountInMaximum, "Too much requested");
        _swap(params.tokenIn, amountIn, params.tokenOut, params.amountOut, params.recipient);
    }

    function _swap(address tokenIn, uint256 amountIn, address tokenOut, uint256 amountOut, address recipient) internal {
        IERC20(tokenIn).safeTransferFrom(msg.sender, address(this), amountIn);
        IMintable(tokenOut).mint(recipient, amountOut);
    }

    function _divUp(uint256 a, uint256 b) internal pure returns (uint256) {
        return a == 0 ? 0 : (a - 1) / b + 1;
    }
}
