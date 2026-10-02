// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IQuoterV2} from "../../src/interfaces/IQuoterV2.sol";
import {StubSwapRouter} from "./StubSwapRouter.sol";

/// @notice Uniswap's QuoterV2 over a `StubSwapRouter`: it quotes the single-pool swaps the router makes, at the price
/// the router is given and with its rounding, and reverts as the real quoter does on a swap of nothing. Test code only.
contract StubQuoterV2 is IQuoterV2 {
    StubSwapRouter public immutable router;

    constructor(StubSwapRouter router_) {
        router = router_;
    }

    function factory() external view returns (address) {
        return address(router);
    }

    function quoteExactInputSingle(QuoteExactInputSingleParams memory params)
        external
        view
        returns (uint256, uint160, uint32, uint256)
    {
        require(params.amountIn != 0, "AS");
        uint256 amountOut = params.tokenIn == router.usdg()
            ? params.amountIn * 1e20 / router.prices(params.tokenOut)
            : params.amountIn * router.prices(params.tokenIn) / 1e20;
        return (amountOut, 0, 0, 0);
    }

    function quoteExactOutputSingle(QuoteExactOutputSingleParams memory params)
        external
        view
        returns (uint256, uint160, uint32, uint256)
    {
        require(params.amount != 0, "AS");
        uint256 amountIn = params.tokenIn == router.usdg()
            ? _divUp(params.amount * router.prices(params.tokenOut), 1e20)
            : _divUp(params.amount * 1e20, router.prices(params.tokenIn));
        return (amountIn, 0, 0, 0);
    }

    function _divUp(uint256 a, uint256 b) internal pure returns (uint256) {
        return a == 0 ? 0 : (a - 1) / b + 1;
    }
}
