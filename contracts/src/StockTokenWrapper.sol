// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC20Wrapper} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Wrapper.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";

/// @notice The Stock Token views the wrapper forwards, as Robinhood's testnet tokens expose them.
interface IWrappedStock is IERC20Metadata {
    /// @notice The ERC-8056 multiplier now.
    function uiMultiplier() external view returns (uint256);

    /// @notice The ERC-8056 multiplier after `effectiveAt`.
    function newUIMultiplier() external view returns (uint256);

    /// @notice When `newUIMultiplier` takes effect, in seconds; zero when no change is scheduled.
    function effectiveAt() external view returns (uint256);

    /// @notice Whether the issuer has paused the token's price.
    function tokenPaused() external view returns (bool);

    /// @notice Whether the token's transfers are paused.
    function paused() external view returns (bool);

    // slither-disable-start naming-convention
    /// @notice The registry that holds the issuer's roles, pause and blocklist.
    // forge-lint: disable-next-line(mixed-case-function)
    function ACCESS_CONTROLLED_REGISTRY() external view returns (address);
    // slither-disable-end naming-convention
}

/// @title Wrapped Stock Token
/// @notice A one-to-one ERC-20 over a Stock Token, with the views the band and the accounts read from a Stock
/// Token: the ERC-8056 multiplier, the issuer's pause and its registry, all those of the wrapped token, and
/// `oraclePaused`, which is the wrapped token's `tokenPaused`. Tokens that predate `oraclePaused` take part
/// in the band and the accounts through it. The wrapper holds the stock; whoever holds the wrapper owns it.
contract StockTokenWrapper is ERC20Wrapper, ERC20Permit {
    /// @notice The wrapped Stock Token, as the views read it.
    IWrappedStock public immutable stock;

    constructor(IWrappedStock stock_)
        ERC20(string.concat("Wrapped ", stock_.name()), string.concat("w", stock_.symbol()))
        ERC20Permit(string.concat("Wrapped ", stock_.name()))
        ERC20Wrapper(IERC20(address(stock_)))
    {
        stock = stock_;
    }

    /// @notice The wrapped token's ERC-8056 multiplier now.
    function uiMultiplier() external view returns (uint256) {
        return stock.uiMultiplier();
    }

    /// @notice The wrapped token's multiplier after `effectiveAt`.
    function newUIMultiplier() external view returns (uint256) {
        return stock.newUIMultiplier();
    }

    /// @notice When the wrapped token's next multiplier takes effect, in seconds; zero when none is scheduled.
    function effectiveAt() external view returns (uint256) {
        return stock.effectiveAt();
    }

    /// @notice Whether the issuer has paused the wrapped token's price.
    function oraclePaused() external view returns (bool) {
        return stock.tokenPaused();
    }

    /// @notice Whether the wrapped token's transfers are paused.
    function paused() external view returns (bool) {
        return stock.paused();
    }

    // slither-disable-start naming-convention
    /// @notice The wrapped token's registry, which holds the issuer's roles, pause and blocklist.
    // forge-lint: disable-next-line(mixed-case-function)
    function ACCESS_CONTROLLED_REGISTRY() external view returns (address) {
        return stock.ACCESS_CONTROLLED_REGISTRY();
    }

    // slither-disable-end naming-convention

    /// @inheritdoc ERC20Wrapper
    function decimals() public view override(ERC20, ERC20Wrapper) returns (uint8) {
        return super.decimals();
    }
}
