// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity ^0.8.12;

import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {IUniswapV3SwapCallback} from "@uniswap/v3-core/contracts/interfaces/callback/IUniswapV3SwapCallback.sol";

interface IQuoterFactory {
    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address);
}

/// @title Uniswap v3 single-pool quoter
/// @notice The single-pool quotes of QuoterV2 over one Uniswap v3 factory, for a chain that has no QuoterV2. A quote
/// swaps in the pool and reverts inside the callback with what it saw, so it is read with `eth_call`, never sent.
/// The initialized ticks a swap crosses are not counted, and its gas is the call's own.
contract UniswapV3Quoter is IUniswapV3SwapCallback {
    struct QuoteExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint256 amountIn;
        uint24 fee;
        uint160 sqrtPriceLimitX96;
    }

    struct QuoteExactOutputSingleParams {
        address tokenIn;
        address tokenOut;
        uint256 amount;
        uint24 fee;
        uint160 sqrtPriceLimitX96;
    }

    /// @notice The factory has no pool of the two tokens at the fee.
    error NoPool();
    /// @notice The callback came from a contract other than the pool being quoted.
    error NotThePool(address caller);

    uint160 private constant MIN_SQRT_RATIO = 4295128739;
    uint160 private constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    /// @notice The Uniswap v3 factory whose pools the quoter quotes.
    address public immutable factory;

    constructor(address factory_) {
        factory = factory_;
    }

    /// @notice What `amountIn` of `tokenIn` would pay out in `tokenOut`, and the pool's price after it.
    function quoteExactInputSingle(QuoteExactInputSingleParams memory params)
        external
        returns (uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
    {
        uint256 gasBefore = gasleft();
        (amountOut, sqrtPriceX96After) =
            _quote(params.tokenIn, params.tokenOut, params.fee, int256(params.amountIn), params.sqrtPriceLimitX96);
        gasEstimate = gasBefore - gasleft();
        initializedTicksCrossed = 0;
    }

    /// @notice What `amount` of `tokenOut` would cost in `tokenIn`, and the pool's price after it.
    function quoteExactOutputSingle(QuoteExactOutputSingleParams memory params)
        external
        returns (uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
    {
        uint256 gasBefore = gasleft();
        (amountIn, sqrtPriceX96After) =
            _quote(params.tokenIn, params.tokenOut, params.fee, -int256(params.amount), params.sqrtPriceLimitX96);
        gasEstimate = gasBefore - gasleft();
        initializedTicksCrossed = 0;
    }

    /// @inheritdoc IUniswapV3SwapCallback
    function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes calldata data)
        external
        view
        override
    {
        address pool = abi.decode(data, (address));
        if (msg.sender != pool) revert NotThePool(msg.sender);
        (uint160 sqrtPriceX96After,,,,,,) = IUniswapV3Pool(pool).slot0();
        bool zeroForOne = amount0Delta > 0;
        uint256 amount = uint256(zeroForOne ? -amount1Delta : -amount0Delta);
        uint256 cost = uint256(zeroForOne ? amount0Delta : amount1Delta);
        // exact input reports what it paid out and exact output what it took in; the caller picks the one it asked for
        bytes memory reply = abi.encode(amount, cost, sqrtPriceX96After);
        assembly {
            revert(add(reply, 32), mload(reply))
        }
    }

    function _quote(address tokenIn, address tokenOut, uint24 fee, int256 amountSpecified, uint160 sqrtPriceLimitX96)
        private
        returns (uint256 amount, uint160 sqrtPriceX96After)
    {
        address pool = IQuoterFactory(factory).getPool(tokenIn, tokenOut, fee);
        if (pool == address(0)) revert NoPool();
        bool zeroForOne = tokenIn < tokenOut;
        try IUniswapV3Pool(pool)
            .swap(
                address(this),
                zeroForOne,
                amountSpecified,
                sqrtPriceLimitX96 == 0 ? (zeroForOne ? MIN_SQRT_RATIO + 1 : MAX_SQRT_RATIO - 1) : sqrtPriceLimitX96,
                abi.encode(pool)
            ) {
            revert NoPool();
        } catch (bytes memory reason) {
            if (reason.length != 96) {
                assembly {
                    revert(add(reason, 32), mload(reason))
                }
            }
            (uint256 out, uint256 cost, uint160 sqrtAfter) = abi.decode(reason, (uint256, uint256, uint160));
            return (amountSpecified > 0 ? out : cost, sqrtAfter);
        }
    }
}
