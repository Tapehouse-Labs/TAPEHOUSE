// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IStockLendingBorrower} from "../../src/interfaces/IStockLendingBorrower.sol";
import {StockLendingVault} from "../../src/StockLendingVault.sol";

interface IUniswapV3SwapPool {
    function token0() external view returns (address);
    function swap(
        address recipient,
        bool zeroForOne,
        int256 amountSpecified,
        uint160 sqrtPriceLimitX96,
        bytes calldata data
    ) external returns (int256 amount0, int256 amount1);
}

/// @notice A stock lending vault's borrower that, on a buy-in, buys exactly what it must return from a Uniswap v3
/// pool with the USDG it holds, and repays it. Test code only.
contract PoolBorrowerDouble is IStockLendingBorrower {
    using SafeERC20 for IERC20;

    uint160 internal constant MIN_SQRT_RATIO = 4295128739;
    uint160 internal constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

    StockLendingVault public immutable vault;
    IERC20 public immutable usdg;
    IUniswapV3SwapPool public immutable pool;
    uint256 public paid;

    constructor(StockLendingVault vault_, IERC20 usdg_, IUniswapV3SwapPool pool_) {
        (vault, usdg, pool) = (vault_, usdg_, pool_);
        IERC20(vault_.asset()).forceApprove(address(vault_), type(uint256).max);
    }

    function borrow(uint256 assets) external {
        vault.borrow(assets, address(this));
    }

    function buyIn(uint256 assets) external {
        require(msg.sender == address(vault));
        bool usdgIsToken0 = pool.token0() == address(usdg);
        pool.swap(
            address(this), usdgIsToken0, -int256(assets), usdgIsToken0 ? MIN_SQRT_RATIO + 1 : MAX_SQRT_RATIO - 1, ""
        );
        vault.repay(assets);
    }

    function uniswapV3SwapCallback(int256 amount0Delta, int256 amount1Delta, bytes calldata) external {
        require(msg.sender == address(pool));
        uint256 owed = uint256(amount0Delta > 0 ? amount0Delta : amount1Delta);
        paid += owed;
        usdg.safeTransfer(msg.sender, owed);
    }
}
