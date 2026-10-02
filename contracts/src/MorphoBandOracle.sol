// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IOracle} from "./interfaces/IOracle.sol";

/// @title Tapehouse band as a Morpho Blue oracle
/// @notice Prices one Robinhood Stock Token, a Morpho Blue market's collateral, in the market's loan token at the
/// low edge of the asset's Tapehouse band, open or closed, with the loan token taken at par with the US dollar. A
/// band with no price, halted or with no live leg, has no answer: the oracle reverts, which stops borrowing,
/// withdrawing collateral against debt and liquidating, and never reports a price of zero.
/// @dev A market's oracle is fixed when the market is created, so the oracle holds the band as governed state: its
/// owner re-points it to a redeployed band that names the same Stock Token, and no market has to move. The band's
/// price is the token's, its ERC-8056 multiplier already applied, per whole token with 8 decimals; the oracle
/// scales it to Morpho's 36 + loan decimals − collateral decimals and applies no multiplier of its own.
contract MorphoBandOracle is IOracle, Ownable2Step {
    /// @notice The decimals of the band's prices.
    uint8 public constant BAND_DECIMALS = 8;

    /// @notice The asset's symbol in the band.
    bytes32 public immutable symbol;
    /// @notice The Stock Token the band names for `symbol`: the market's collateral.
    address public immutable collateralToken;
    /// @notice The market's loan token.
    address public immutable loanToken;
    /// @notice What the band's price is multiplied by: 10^(36 + loan decimals − collateral decimals − 8).
    uint256 public immutable scaleFactor;

    /// @notice The band program the oracle reads.
    IBand public band;
    bool private followsSequencer;

    /// @notice The oracle reads `newBand`, in place of `previousBand`.
    event BandSet(IBand indexed previousBand, IBand indexed newBand);

    /// @notice The band of `symbol` has no price: it is halted, or has no live leg, most often because no one has
    /// written a 24/7 price for two minutes while the market is closed.
    error NoAnswer(bytes32 symbol);
    /// @notice The L2 sequencer is down or came back an hour ago or less.
    error SequencerNotSettled();
    /// @notice `band` names no Stock Token for `symbol`.
    error NoStockToken(IBand band, bytes32 symbol);
    /// @notice `band` names `token` for the asset, not the oracle's collateral.
    error AssetMismatch(IBand band, address token);
    /// @notice A price scaled by `scaleFactor` could overflow: the loan token has more than 29 decimals more than the
    /// collateral.
    error ScaleTooLarge(uint256 scaleFactor);
    /// @notice The owner cannot renounce: the oracle needs one to follow a band redeploy.
    error OwnershipCannotBeRenounced();

    /// @param band_ The band program.
    /// @param symbol_ The asset's symbol in the band, which must name its Stock Token.
    /// @param loanToken_ The market's loan token.
    /// @param initialOwner The owner, who re-points the oracle to a redeployed band.
    constructor(IBand band_, bytes32 symbol_, IERC20Metadata loanToken_, address initialOwner) Ownable(initialOwner) {
        // slither-disable-next-line unused-return
        (,,, address token) = band_.asset(symbol_); // forge-lint: disable-line(unused-return)
        if (token == address(0)) revert NoStockToken(band_, symbol_);
        symbol = symbol_;
        collateralToken = token;
        loanToken = address(loanToken_);
        uint256 scale = 10 ** (36 + uint256(loanToken_.decimals()) - IERC20Metadata(token).decimals() - BAND_DECIMALS);
        if (scale > type(uint256).max / type(uint64).max) revert ScaleTooLarge(scale);
        scaleFactor = scale;
        _setBand(band_);
    }

    /// @inheritdoc IOracle
    /// @dev The band's low edge in every state that has a price: open, closed and degraded, a lapsed signed halt
    /// included. Reverts with `NoAnswer` while the band is halted or has no live leg, and with `SequencerNotSettled`
    /// where the band follows a sequencer-uptime feed and the sequencer is not settled.
    function price() external view returns (uint256) {
        // slither-disable-next-line unused-return
        (uint8 state,,,, uint64 low,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if (state == 0 || low == 0) revert NoAnswer(symbol);
        if (followsSequencer && !band.sequencerSettled()) revert SequencerNotSettled();
        return low * scaleFactor;
    }

    /// @notice The asset's trading halt as the band holds it: whether a halt signed by Tapehouse's halt signer holds,
    /// the one input the band takes on Tapehouse's own signature, until when, when its last message was issued, and
    /// whether the issuer has paused the Stock Token's oracle.
    function halt() external view returns (bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused) {
        return band.halt(symbol);
    }

    /// @notice Points the oracle at `newBand`, which must name the same Stock Token for the asset. Owner only.
    function setBand(IBand newBand) external onlyOwner {
        // slither-disable-next-line unused-return
        (,,, address token) = newBand.asset(symbol); // forge-lint: disable-line(unused-return)
        if (token != collateralToken) revert AssetMismatch(newBand, token);
        _setBand(newBand);
    }

    /// @notice Always reverts: the oracle needs an owner to follow a band redeploy.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    function _setBand(IBand newBand) private {
        // slither-disable-next-line unused-return
        (address sequencer,) = newBand.chainConfig(); // forge-lint: disable-line(unused-return)
        emit BandSet(band, newBand);
        band = newBand;
        followsSequencer = sequencer != address(0);
    }
}
