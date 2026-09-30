// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IMargin} from "../../src/interfaces/IMargin.sol";

/// @notice A margin engine whose requirement is a set share of the gross value, with the missing bits and regime it
/// is given. Test code only.
contract MarginDouble is IMargin {
    bytes32[] internal _assets;
    address public immutable band;
    address public immutable ethUsdFeed;
    uint256 public rateBps = 20_00;
    uint8 public missing;
    uint8 public regime = 2;

    constructor(bytes32[] memory assets_, address band_, address ethUsdFeed_) {
        _assets = assets_;
        band = band_;
        ethUsdFeed = ethUsdFeed_;
    }

    function set(uint256 rateBps_, uint8 missing_, uint8 regime_) external {
        (rateBps, missing, regime) = (rateBps_, missing_, regime_);
    }

    function assets() external view returns (bytes32[] memory) {
        return _assets;
    }

    function currentRequirement(int256[] calldata quantities, uint256[] calldata prices)
        external
        view
        returns (uint256 requirement, uint8, uint8)
    {
        require(quantities.length == _assets.length && prices.length == _assets.length, "length");
        uint256 gross;
        for (uint256 i; i < quantities.length; ++i) {
            uint256 quantity = uint256(quantities[i] < 0 ? -quantities[i] : quantities[i]);
            gross += quantity * prices[i] / 1e8;
        }
        return (gross * rateBps / 10_000, missing, regime);
    }

    function weekendLeverage() external pure returns (uint32) {
        return 50_000;
    }
}
