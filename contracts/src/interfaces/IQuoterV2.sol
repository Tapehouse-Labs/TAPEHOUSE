// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @title Uniswap v3's QuoterV2, as Tapehouse's SDKs use it
/// @notice The single-pool quotes of Uniswap v3-periphery's `IQuoterV2` (1.1.1 and later), and the factory whose
/// pools it quotes. A quote simulates the swap and reverts inside the call, so it is read with `eth_call`, never sent.
interface IQuoterV2 {
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

    /// @notice What swapping `amountIn` of one token for another through the pool at `fee` would pay out, the pool's
    /// price after it, the initialized ticks it would cross, and an estimate of its gas.
    function quoteExactInputSingle(QuoteExactInputSingleParams memory params)
        external
        returns (uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate);

    /// @notice What buying `amount` of one token with another through the pool at `fee` would cost, the pool's price
    /// after it, the initialized ticks it would cross, and an estimate of its gas.
    function quoteExactOutputSingle(QuoteExactOutputSingleParams memory params)
        external
        returns (uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate);

    /// @notice The Uniswap v3 factory whose pools the quoter quotes.
    function factory() external view returns (address);
}
