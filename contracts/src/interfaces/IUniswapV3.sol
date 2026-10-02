// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @title Uniswap v3's SwapRouter02, as its Tapehouse callers use it
/// @notice The single-pool swaps of `IV3SwapRouter` in Uniswap's swap-router-contracts 1.1.0, and the factory it routes
/// through.
interface IV3SwapRouter {
    struct ExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 amountIn;
        uint256 amountOutMinimum;
        uint160 sqrtPriceLimitX96;
    }

    struct ExactOutputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 amountOut;
        uint256 amountInMaximum;
        uint160 sqrtPriceLimitX96;
    }

    /// @notice Swaps `amountIn` of one token for as much as possible of another, at least `amountOutMinimum`.
    function exactInputSingle(ExactInputSingleParams calldata params) external payable returns (uint256 amountOut);

    /// @notice Swaps as little as possible of one token, at most `amountInMaximum`, for `amountOut` of another.
    function exactOutputSingle(ExactOutputSingleParams calldata params) external payable returns (uint256 amountIn);

    /// @notice The Uniswap v3 factory whose pools the router swaps through.
    function factory() external view returns (address);
}

/// @title Uniswap v3's factory, as its Tapehouse callers use it
interface IUniswapV3Factory {
    /// @notice The pool of `tokenA` and `tokenB` at `fee`, in hundredths of a basis point; zero where there is none.
    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address pool);
}
