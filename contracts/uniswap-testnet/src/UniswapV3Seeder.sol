// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity ^0.8.12;

import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {IUniswapV3MintCallback} from "@uniswap/v3-core/contracts/interfaces/callback/IUniswapV3MintCallback.sol";

/// @title Uniswap v3 liquidity seeder
/// @notice Adds liquidity to a pool for the caller, who owns the position and pays the pool's tokens from an
/// allowance to this contract. Uniswap v3 takes liquidity through a callback, which an account cannot answer.
contract UniswapV3Seeder is IUniswapV3MintCallback {
    /// @notice The callback came from a contract other than the pool being seeded.
    error NotThePool(address caller);
    /// @notice A token did not move.
    error TransferFailed(address token);

    address private pool;
    address private payer;

    /// @notice Mints `liquidity` between `lower` and `upper` for the caller, who pays the tokens it takes.
    /// @return amount0 The pool's token0 paid.
    /// @return amount1 The pool's token1 paid.
    function seed(IUniswapV3Pool pool_, int24 lower, int24 upper, uint128 liquidity)
        external
        returns (uint256 amount0, uint256 amount1)
    {
        pool = address(pool_);
        payer = msg.sender;
        (amount0, amount1) = pool_.mint(msg.sender, lower, upper, liquidity, "");
        delete pool;
        delete payer;
    }

    /// @inheritdoc IUniswapV3MintCallback
    function uniswapV3MintCallback(uint256 amount0Owed, uint256 amount1Owed, bytes calldata) external override {
        if (msg.sender != pool) revert NotThePool(msg.sender);
        if (amount0Owed > 0) _pay(IUniswapV3Pool(msg.sender).token0(), amount0Owed);
        if (amount1Owed > 0) _pay(IUniswapV3Pool(msg.sender).token1(), amount1Owed);
    }

    function _pay(address token, uint256 amount) private {
        (bool ok, bytes memory result) = token.call(abi.encodeWithSelector(0x23b872dd, payer, msg.sender, amount));
        if (!ok || (result.length != 0 && !abi.decode(result, (bool)))) revert TransferFailed(token);
    }
}
