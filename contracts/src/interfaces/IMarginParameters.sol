// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IMargin} from "./IMargin.sol";

/// @title Tapehouse margin engine, as its risk parameters are read
/// @notice The engine's reads, the risk parameters it generates its scenarios from with their floors and ceilings, the
/// call that replaces them, its errors, and the events that record every value they take.
interface IMarginParameters is IMargin {
    /// @notice Another update came less than a day after the last; the next may come at `nextUpdateAt`.
    error UpdateTooSoon(uint64 nextUpdateAt);
    /// @notice An array's length does not match the engine's assets.
    error LengthMismatch();
    /// @notice The engine holds no asset `symbol`.
    error UnknownAsset(bytes32 symbol);
    /// @notice The volatility of `symbol` is below its `floor` or above the maximum.
    error InvalidVolatility(bytes32 symbol, uint32 value, uint32 floor);
    /// @notice The correlation of `symbol` and `other` is below its `floor` or above 10 000.
    error InvalidCorrelation(bytes32 symbol, bytes32 other, uint16 value, uint16 floor);
    /// @notice The weekend gap of `symbol` is below its `floor` or above the maximum.
    error InvalidGap(bytes32 symbol, uint32 value, uint32 floor);
    /// @notice A depth of `symbol` is zero or above its `ceiling`.
    error InvalidDepth(bytes32 symbol, uint32 value, uint32 ceiling);
    /// @notice The volatility of `symbol` moved from `previous` to `value`, by more than ×1.5.
    error VolatilityStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    /// @notice The correlation of `symbol` and `other` moved from `previous` to `value`, by more than ×1.5.
    error CorrelationStepTooLarge(bytes32 symbol, bytes32 other, uint16 previous, uint16 value);
    /// @notice The weekend gap of `symbol` moved from `previous` to `value`, by more than ×1.5.
    error GapStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    /// @notice A depth of `symbol` rose from `previous` to `value`, by more than ×1.5.
    error DepthStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    /// @notice The correlation matrix is not positive definite.
    error NotPositiveDefinite();
    /// @notice `account` is not the engine's owner.
    error OwnableUnauthorizedAccount(address account);

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

    /// @notice Replaces every parameter, in the order of `assets()`: a volatility and a weekend gap for each asset,
    /// the correlations of the upper triangle row by row, and each asset's selling then buying depth. Each value
    /// stays at or above its floor, a depth above zero and at or below its ceiling, and each moves by at most ×1.5
    /// or ÷1.5, except that a depth may fall by any amount. The correlation matrix stays positive definite. Owner
    /// only, at most once a day.
    function setParameters(
        uint32[] calldata volatilities,
        uint16[] calldata correlations,
        uint32[] calldata gaps,
        uint32[] calldata depths
    ) external;

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
