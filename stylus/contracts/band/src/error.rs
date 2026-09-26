//! Errors of the band program. RedStone errors keep the reference contract's names and selectors.

use stylus_sdk::alloy_sol_types::sol;
use stylus_sdk::prelude::*;

sol! {
    #[derive(Debug, PartialEq, Eq)]
    error CalldataMustHaveValidPayload();
    #[derive(Debug, PartialEq, Eq)]
    error CalldataOverOrUnderFlow();
    #[derive(Debug, PartialEq, Eq)]
    error IncorrectUnsignedMetadataSize();
    #[derive(Debug, PartialEq, Eq)]
    error InvalidSignature(bytes32 signedHash);
    #[derive(Debug, PartialEq, Eq)]
    error SignerNotAuthorised(address receivedSigner);
    #[derive(Debug, PartialEq, Eq)]
    error TooLargeValueByteSize(uint256 valueByteSize);
    #[derive(Debug, PartialEq, Eq)]
    error DataTimestampCannotBeZero();
    #[derive(Debug, PartialEq, Eq)]
    error TimestampsMustBeEqual();
    #[derive(Debug, PartialEq, Eq)]
    error InsufficientNumberOfUniqueSigners(uint256 receivedSignersCount, uint256 requiredSignersCount);
    #[derive(Debug, PartialEq, Eq)]
    error TimestampFromTooLongFuture(uint256 receivedTimestampSeconds, uint256 blockTimestamp);
    #[derive(Debug, PartialEq, Eq)]
    error TimestampIsTooOld(uint256 receivedTimestampSeconds, uint256 blockTimestamp);
    #[derive(Debug, PartialEq, Eq)]
    error PackageNotNewer(bytes32 feedId, uint64 storedTimestampMs, uint64 packageTimestampMs);
    #[derive(Debug, PartialEq, Eq)]
    error LengthMismatch();
    #[derive(Debug, PartialEq, Eq)]
    error ZeroSymbol();
    #[derive(Debug, PartialEq, Eq)]
    error NoLegs(bytes32 symbol);
    #[derive(Debug, PartialEq, Eq)]
    error DuplicateAsset(bytes32 symbol);
    #[derive(Debug, PartialEq, Eq)]
    error InvalidFeed(address feed);
}

#[derive(SolidityError, Debug, PartialEq, Eq)]
pub enum BandError {
    CalldataMustHaveValidPayload(CalldataMustHaveValidPayload),
    CalldataOverOrUnderFlow(CalldataOverOrUnderFlow),
    IncorrectUnsignedMetadataSize(IncorrectUnsignedMetadataSize),
    InvalidSignature(InvalidSignature),
    SignerNotAuthorised(SignerNotAuthorised),
    TooLargeValueByteSize(TooLargeValueByteSize),
    DataTimestampCannotBeZero(DataTimestampCannotBeZero),
    TimestampsMustBeEqual(TimestampsMustBeEqual),
    InsufficientNumberOfUniqueSigners(InsufficientNumberOfUniqueSigners),
    TimestampFromTooLongFuture(TimestampFromTooLongFuture),
    TimestampIsTooOld(TimestampIsTooOld),
    PackageNotNewer(PackageNotNewer),
    LengthMismatch(LengthMismatch),
    ZeroSymbol(ZeroSymbol),
    NoLegs(NoLegs),
    DuplicateAsset(DuplicateAsset),
    InvalidFeed(InvalidFeed),
}
