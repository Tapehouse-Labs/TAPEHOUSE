// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IBand} from "./IBand.sol";

/// @title Tapehouse band, as a RedStone relayer and a keeper call it
/// @notice The band's reads, its events, its halt signer, its price writes from signed RedStone data packages, its
/// signed halts and multiplier syncs, and the errors of each.
interface IBandPrices is IBand {
    /// @notice `writePrices` stored `value`, the median of `feedId`'s signed package of `packageTimestampMs`.
    event PriceWritten(bytes32 indexed feedId, uint256 value, uint64 packageTimestampMs);
    /// @notice The 24/7 leg of `symbol`, priced from an index, was anchored to the Chainlink print of `updatedAt` and
    /// the index price written with it.
    event Anchored(bytes32 indexed symbol, uint64 chainlinkPrice, uint64 indexPrice, uint64 updatedAt);
    /// @notice The multiplier change of `symbol`'s Stock Token was recorded: the multipliers before and after it, and
    /// when it takes effect. The multiplier before is zero when the change's size is not known.
    event MultiplierRecorded(
        bytes32 indexed symbol, uint128 multiplierBefore, uint128 multiplierAfter, uint64 effectiveAt
    );
    /// @notice A Chainlink round confirmed the multiplier change of `symbol`'s Stock Token that took effect at
    /// `effectiveAt`.
    event MultiplierConfirmed(bytes32 indexed symbol, uint64 effectiveAt);
    /// @notice `writeHalt` stored a trading halt of `symbol` signed by Tapehouse's halt signer.
    event HaltWritten(bytes32 indexed symbol, bool halted, uint64 issuedAt, uint64 expiresAt);
    /// @notice The owner replaced the halt signer.
    event HaltSignerUpdated(address indexed previousSigner, address indexed newSigner);
    /// @notice The owner started a two-step transfer of ownership to `newOwner`.
    event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner);
    /// @notice `newOwner` accepted the band's ownership.
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

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
    /// @notice The band prices no asset `symbol`.
    error UnknownAsset(bytes32 symbol);
    /// @notice `symbol` has no Stock Token.
    error NoToken(bytes32 symbol);
    /// @notice The Stock Token's multipliers cannot be read.
    error InvalidToken(address token);
    /// @notice A halt message must be issued no later than now, expire after now and span at most an hour.
    error HaltOutsideWindow(uint64 issuedAt, uint64 expiresAt, uint64 blockTimestamp);
    /// @notice The band already stores a halt message of `symbol` issued at `storedIssuedAt`, at least as new.
    error HaltNotNewer(bytes32 symbol, uint64 storedIssuedAt, uint64 issuedAt);

    /// @notice Verifies the signed RedStone data packages in `payload` for `feedIds`, stores each feed's median value
    /// and returns the values. Anyone may call it.
    function writePrices(bytes32[] memory feedIds, bytes calldata payload) external returns (uint256[] memory);

    /// @notice The stored value of `feedId`, its package's timestamp in milliseconds, and the block timestamp it was
    /// written at. All zero for a feed never written.
    function price(bytes32 feedId) external view returns (uint256 value, uint64 packageTimestampMs, uint64 writtenAt);

    /// @notice Writes a trading halt of `symbol` signed by the halt signer as EIP-712
    /// `HaltState(bytes32 symbol,bool halted,uint64 issuedAt,uint64 expiresAt)` under the domain "Tapehouse Band",
    /// version "1", of this band: `halted` holds the asset halted until `expiresAt`, otherwise the halt is lifted. Anyone
    /// may call it.
    function writeHalt(bytes32 symbol, bool halted, uint64 issuedAt, uint64 expiresAt, bytes calldata signature)
        external;

    /// @notice Confirms a material multiplier change of `symbol`'s Stock Token past its step once a Chainlink round
    /// that started at or after it falls inside the band of the 24/7 leg alone, then records the token's latest change.
    /// Anyone may call it. Returns the status of the change in force, as `corporateAction` reports it.
    function syncMultiplier(bytes32 symbol) external returns (uint8);

    /// @notice The address whose signed trading halts the band accepts: the one input it takes on Tapehouse's own
    /// signature.
    function haltSigner() external view returns (address);
}
