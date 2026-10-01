// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @title Tapehouse margin engine
/// @notice The Stylus program that margins a portfolio of Stock Tokens, as its Solidity readers call it.
interface IMargin {
    /// @notice The assets, in the order of every quantity and price list.
    function assets() external view returns (bytes32[] memory);

    /// @notice The margin a portfolio needs over `horizon` seconds, in USD with 18 decimals, without the open market's
    /// buffer, and the missing bits of `currentRequirement`. `spansClosure` counts the weekend-gap scenarios.
    function requirement(int256[] calldata quantities, uint256[] calldata prices, uint64 horizon, bool spansClosure)
        external
        view
        returns (uint256 margin, uint8 missing);

    /// @notice The margin a portfolio needs now, in USD with 18 decimals, a bit for every asset whose liquidity
    /// input is missing, and the regime the band's session puts it in: 0 unknown, 1 closed, 2 open, 3 closing.
    /// `quantities` are signed token amounts with 18 decimals and `prices` USD with 8 decimals, in the order of
    /// `assets()`.
    function currentRequirement(int256[] calldata quantities, uint256[] calldata prices)
        external
        view
        returns (uint256 margin, uint8 missing, uint8 regime);

    /// @notice The most gross exposure the current requirement allows per unit of margin across a closure, in
    /// basis points.
    function weekendLeverage() external view returns (uint32);

    /// @notice The Chainlink feed that prices WETH.
    function ethUsdFeed() external view returns (address);

    /// @notice The band whose session sets the current requirement's regime.
    function band() external view returns (address);
}
