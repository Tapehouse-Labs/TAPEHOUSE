//! Errors of the margin program.

use ownable::{OwnableInvalidOwner, OwnableUnauthorizedAccount};
use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::prelude::*;

sol! {
    #[derive(Debug, PartialEq, Eq)]
    error AssetCount(uint256 count);
    #[derive(Debug, PartialEq, Eq)]
    error ZeroSymbol();
    #[derive(Debug, PartialEq, Eq)]
    error DuplicateAsset(bytes32 symbol);
    #[derive(Debug, PartialEq, Eq)]
    error LengthMismatch();
    #[derive(Debug, PartialEq, Eq)]
    error UnknownAsset(bytes32 symbol);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidVolatility(bytes32 symbol, uint32 value, uint32 floor);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidCorrelation(bytes32 symbol, bytes32 other, uint16 value, uint16 floor);
    #[derive(Debug, PartialEq, Eq)]
    error VolatilityStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    #[derive(Debug, PartialEq, Eq)]
    error CorrelationStepTooLarge(bytes32 symbol, bytes32 other, uint16 previous, uint16 value);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidGap(bytes32 symbol, uint32 value, uint32 floor);
    #[derive(Debug, PartialEq, Eq)]
    error GapStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    #[derive(Debug, PartialEq, Eq)]
    error NotPositiveDefinite();
    #[derive(Debug, PartialEq, Eq)]
    error UpdateTooSoon(uint64 nextUpdateAt);
    #[derive(Debug, PartialEq, Eq)]
    error UnsupportedScenarioSize(uint16 size);
    #[derive(Debug, PartialEq, Eq)]
    error ScenarioOutOfRange(uint16 index);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidDepth(bytes32 symbol, uint32 value, uint32 ceiling);
    #[derive(Debug, PartialEq, Eq)]
    error DepthStepTooLarge(bytes32 symbol, uint32 previous, uint32 value);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidPool(bytes32 symbol, address pool);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidFeed(address feed);
    #[derive(Debug, PartialEq, Eq)]
    error ExposureTooLarge(bytes32 symbol);
    #[derive(Debug, PartialEq, Eq)]
    error InsufficientGas();
}

/// The errors of the margin program.
#[derive(SolidityError, Debug, Clone, PartialEq, Eq)]
pub enum MarginError {
    AssetCount(AssetCount),
    ZeroSymbol(ZeroSymbol),
    DuplicateAsset(DuplicateAsset),
    LengthMismatch(LengthMismatch),
    UnknownAsset(UnknownAsset),
    InvalidVolatility(InvalidVolatility),
    InvalidCorrelation(InvalidCorrelation),
    VolatilityStepTooLarge(VolatilityStepTooLarge),
    CorrelationStepTooLarge(CorrelationStepTooLarge),
    InvalidGap(InvalidGap),
    GapStepTooLarge(GapStepTooLarge),
    NotPositiveDefinite(NotPositiveDefinite),
    UpdateTooSoon(UpdateTooSoon),
    UnsupportedScenarioSize(UnsupportedScenarioSize),
    ScenarioOutOfRange(ScenarioOutOfRange),
    InvalidDepth(InvalidDepth),
    DepthStepTooLarge(DepthStepTooLarge),
    InvalidPool(InvalidPool),
    InvalidFeed(InvalidFeed),
    ExposureTooLarge(ExposureTooLarge),
    InsufficientGas(InsufficientGas),
    OwnableUnauthorizedAccount(OwnableUnauthorizedAccount),
    OwnableInvalidOwner(OwnableInvalidOwner),
}
