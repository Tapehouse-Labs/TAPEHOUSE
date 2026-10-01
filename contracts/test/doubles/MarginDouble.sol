// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IMargin} from "../../src/interfaces/IMargin.sol";

/// @notice A margin engine whose requirement is a set share of the gross value, with the missing bits and regime it
/// is given, plus, once a pool price is set, the value above it, as the engine's liquidity add-on discounts a price
/// above its pool's. Test code only.
contract MarginDouble is IMargin {
    bytes32[] internal _assets;
    address public immutable band;
    address public immutable ethUsdFeed;
    uint256 public rateBps = 20_00;
    uint256 public openRateBps = 16_00;
    uint8 public missing;
    uint8 public regime = 2;
    uint64 public poolPrice;

    constructor(bytes32[] memory assets_, address band_, address ethUsdFeed_) {
        _assets = assets_;
        band = band_;
        ethUsdFeed = ethUsdFeed_;
    }

    function set(uint256 rateBps_, uint8 missing_, uint8 regime_) external {
        (rateBps, missing, regime) = (rateBps_, missing_, regime_);
    }

    function setOpen(uint256 openRateBps_) external {
        openRateBps = openRateBps_;
    }

    function setPool(uint64 poolPrice_) external {
        poolPrice = poolPrice_;
    }

    function requirement(int256[] calldata quantities, uint256[] calldata prices, uint64, bool)
        external
        view
        returns (uint256, uint8)
    {
        return (_gross(quantities, prices) * openRateBps / 10_000 + _addOn(quantities, prices), missing);
    }

    function assets() external view returns (bytes32[] memory) {
        return _assets;
    }

    function currentRequirement(int256[] calldata quantities, uint256[] calldata prices)
        external
        view
        returns (uint256, uint8, uint8)
    {
        return (_gross(quantities, prices) * rateBps / 10_000 + _addOn(quantities, prices), missing, regime);
    }

    function _gross(int256[] calldata quantities, uint256[] calldata prices) internal view returns (uint256 gross) {
        require(quantities.length == _assets.length && prices.length == _assets.length, "length");
        for (uint256 i; i < quantities.length; ++i) {
            uint256 quantity = uint256(quantities[i] < 0 ? -quantities[i] : quantities[i]);
            gross += quantity * prices[i] / 1e8;
        }
    }

    function _addOn(int256[] calldata quantities, uint256[] calldata prices) internal view returns (uint256 addOn) {
        if (poolPrice == 0) return 0;
        for (uint256 i; i < quantities.length; ++i) {
            if (prices[i] <= poolPrice) continue;
            uint256 quantity = uint256(quantities[i] < 0 ? -quantities[i] : quantities[i]);
            addOn += quantity * (prices[i] - poolPrice) / 1e8;
        }
    }

    function weekendLeverage() external pure returns (uint32) {
        return 50_000;
    }
}
