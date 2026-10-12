// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity ^0.8.12;

import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {IUniswapV3SwapCallback} from "@uniswap/v3-core/contracts/interfaces/callback/IUniswapV3SwapCallback.sol";

interface IFactory {
    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address);
}

/// @title Uniswap v3 single-pool swap router
/// @notice The single-pool swaps of SwapRouter02's `IV3SwapRouter` over one Uniswap v3 factory, for a chain that has
/// no SwapRouter02. The caller pays from an allowance to this contract.
contract UniswapV3Router is IUniswapV3SwapCallback {
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

    struct Swap {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        int256 amount;
        uint160 limit;
    }

    /// @notice The factory has no pool of the two tokens at the fee.
    error NoPool();
    /// @notice The swap paid out less than the minimum.
    error TooLittleReceived(uint256 amountOut);
    /// @notice The swap cost more than the maximum.
    error TooMuchRequested(uint256 amountIn);
    /// @notice A swap that must fill completely, because it has no price limit, did not.
    error NotFilled(uint256 amountOut);
    /// @notice The callback came from a contract other than the pool being swapped through.
    error NotThePool(address caller);
    /// @notice A token did not move.
    error TransferFailed(address token);

    uint160 private constant MIN_SQRT_RATIO = 4295128739;
    uint160 private constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    /// @notice The Uniswap v3 factory whose pools the router swaps through.
    address public immutable factory;

    constructor(address factory_) {
        factory = factory_;
    }

    /// @notice Swaps `amountIn` of `tokenIn` for as much `tokenOut` as the pool pays, at least `amountOutMinimum`.
    function exactInputSingle(ExactInputSingleParams calldata params) external payable returns (uint256 amountOut) {
        (, amountOut) = _swap(
            Swap(
                params.tokenIn,
                params.tokenOut,
                params.fee,
                params.recipient,
                int256(params.amountIn),
                params.sqrtPriceLimitX96
            )
        );
        if (amountOut < params.amountOutMinimum) revert TooLittleReceived(amountOut);
    }

    /// @notice Swaps as little `tokenIn` as the pool takes, at most `amountInMaximum`, for `amountOut` of `tokenOut`.
    function exactOutputSingle(ExactOutputSingleParams calldata params) external payable returns (uint256 amountIn) {
        uint256 received;
        (amountIn, received) = _swap(
            Swap(
                params.tokenIn,
                params.tokenOut,
                params.fee,
                params.recipient,
                -int256(params.amountOut),
                params.sqrtPriceLimitX96
            )
        );
        if (params.sqrtPriceLimitX96 == 0 && received != params.amountOut) revert NotFilled(received);
        if (amountIn > params.amountInMaximum) revert TooMuchRequested(amountIn);
    }

    /// @inheritdoc IUniswapV3SwapCallback
    function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes calldata data) external override {
        (address payer, address tokenIn, address tokenOut, uint24 fee) =
            abi.decode(data, (address, address, address, uint24));
        if (msg.sender != IFactory(factory).getPool(tokenIn, tokenOut, fee)) revert NotThePool(msg.sender);
        uint256 owed = uint256(amount0Delta > 0 ? amount0Delta : amount1Delta);
        (bool ok, bytes memory result) = tokenIn.call(abi.encodeWithSelector(0x23b872dd, payer, msg.sender, owed));
        if (!ok || (result.length != 0 && !abi.decode(result, (bool)))) revert TransferFailed(tokenIn);
    }

    function _swap(Swap memory swap) private returns (uint256 amountIn, uint256 amountOut) {
        address pool = IFactory(factory).getPool(swap.tokenIn, swap.tokenOut, swap.fee);
        if (pool == address(0)) revert NoPool();
        bool zeroForOne = swap.tokenIn < swap.tokenOut;
        (int256 amount0, int256 amount1) = _poolSwap(pool, swap, zeroForOne);
        (amountIn, amountOut) =
            zeroForOne ? (uint256(amount0), uint256(-amount1)) : (uint256(amount1), uint256(-amount0));
    }

    function _poolSwap(address pool, Swap memory swap, bool zeroForOne) private returns (int256, int256) {
        bytes memory data = abi.encode(msg.sender, swap.tokenIn, swap.tokenOut, swap.fee);
        return IUniswapV3Pool(pool).swap(swap.recipient, zeroForOne, swap.amount, _limit(swap.limit, zeroForOne), data);
    }

    function _limit(uint160 limit, bool zeroForOne) private pure returns (uint160) {
        return limit == 0 ? (zeroForOne ? MIN_SQRT_RATIO + 1 : MAX_SQRT_RATIO - 1) : limit;
    }
}
