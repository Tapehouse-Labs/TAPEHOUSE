// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IBand} from "./IBand.sol";

/// @title Tapehouse band, as a RedStone relayer calls it
/// @notice The band's reads, its price writes from signed RedStone data packages, and the errors of a payload it
/// refuses.
interface IBandPrices is IBand {
    /// @notice The payload does not end with RedStone's marker.
    error CalldataMustHaveValidPayload();
    /// @notice The payload is shorter than the sizes it declares.
    error CalldataOverOrUnderFlow();
    /// @notice The payload's unsigned metadata is longer than the payload.
    error IncorrectUnsignedMetadataSize();
    /// @notice A data package's signature does not recover.
    error InvalidSignature(bytes32 signedHash);
    /// @notice A data package was signed by a node outside RedStone's primary production signers.
    error SignerNotAuthorised(address receivedSigner);
    /// @notice A data point's value is over 32 bytes.
    error TooLargeValueByteSize(uint256 valueByteSize);
    /// @notice A data package has no timestamp.
    error DataTimestampCannotBeZero();
    /// @notice The payload's data packages carry different timestamps.
    error TimestampsMustBeEqual();
    /// @notice A feed has fewer unique signers in the payload than the band requires.
    error InsufficientNumberOfUniqueSigners(uint256 receivedSignersCount, uint256 requiredSignersCount);
    /// @notice The packages are timestamped more than a minute ahead of the block.
    error TimestampFromTooLongFuture(uint256 receivedTimestampSeconds, uint256 blockTimestamp);
    /// @notice The packages are more than three minutes older than the block.
    error TimestampIsTooOld(uint256 receivedTimestampSeconds, uint256 blockTimestamp);
    /// @notice The band already stores a package of `feedId` at least as new.
    error PackageNotNewer(bytes32 feedId, uint64 storedTimestampMs, uint64 packageTimestampMs);
    /// @notice The three market-status feeds are written together or not at all.
    error IncompleteStatus();

    /// @notice Verifies the signed RedStone data packages in `payload` for `feedIds`, stores each feed's median value
    /// and returns the values. Anyone may call it.
    function writePrices(bytes32[] memory feedIds, bytes calldata payload) external returns (uint256[] memory);

    /// @notice The stored value of `feedId`, its package's timestamp in milliseconds, and the block timestamp it was
    /// written at. All zero for a feed never written.
    function price(bytes32 feedId) external view returns (uint256 value, uint64 packageTimestampMs, uint64 writtenAt);
}
