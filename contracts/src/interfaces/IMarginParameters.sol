// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IMargin} from "./IMargin.sol";

/// @title Tapehouse margin engine, as its risk parameters are read
/// @notice The engine's reads, the risk parameters it generates its scenarios from with their floors and ceilings, and
/// the events that record every value they take.
interface IMarginParameters is IMargin {
    /// @notice The daily volatility of `symbol` is `value`, in centi-basis-points.
    event VolatilitySet(bytes32 indexed symbol, uint32 value);
    /// @notice The correlation of `symbol` and `other` is `value`, in basis points.
    event CorrelationSet(bytes32 indexed symbol, bytes32 indexed other, uint16 value);
    /// @notice The weekend gap of `symbol` is `value`, in centi-basis-points.
    event GapSet(bytes32 indexed symbol, uint32 value);
    /// @notice The USD a liquidation of `symbol` can sell, then buy, within a 10% move of its pool.
    event DepthSet(bytes32 indexed symbol, uint32 selling, uint32 buying);
    /// @notice The owner started a two-step transfer of ownership to `newOwner`.
    event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner);
    /// @notice `newOwner` accepted the engine's ownership.
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /// @notice The daily volatility of `symbol` and its floor, in centi-basis-points.
    function volatility(bytes32 symbol) external view returns (uint32 value, uint32 floor);

    /// @notice The correlation of `symbol` and `other` and its floor, in basis points; 10 000 for an asset with
    /// itself.
    function correlation(bytes32 symbol, bytes32 other) external view returns (uint16 value, uint16 floor);

    /// @notice The USD a liquidation of `symbol` can sell, then buy, within a 10% move of its pool, and the ceilings
    /// over both.
    function depth(bytes32 symbol)
        external
        view
        returns (uint32 selling, uint32 buying, uint32 sellingCeiling, uint32 buyingCeiling);

    /// @notice The Uniswap v3 pool `symbol` is liquidated in; zero for none.
    function pool(bytes32 symbol) external view returns (address);

    /// @notice The asset that stands for the market; zero for the equal-weighted portfolio of the assets.
    function market() external view returns (bytes32);

    /// @notice When the parameters were last updated; zero before the first update.
    function lastUpdate() external view returns (uint64);

    /// @notice The engine's owner, who may update the parameters.
    function owner() external view returns (address);

    /// @notice The address the owner offered ownership to; zero for none.
    function pendingOwner() external view returns (address);
}
