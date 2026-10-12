// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

/// @notice The multiplier of a Stock Token, as ERC-8056 reports it.
interface IMultiplier {
    /// @notice The multiplier now, with 18 decimals.
    function uiMultiplier() external view returns (uint256);
}

/// @title RedStone round feed
/// @notice A Chainlink-compatible price feed of one RedStone data feed, for a chain without a Chainlink feed of the
/// asset. Anyone may append a payload of RedStone's `redstone-primary-prod` data service to a call of `update`: the
/// feed checks its signatures and keeps its value as the next round, with the package's time as the round's start.
/// Its answers have 8 decimals, and its rounds are numbered from one in one phase, as Chainlink aggregators' are.
/// Where it is given a Stock Token, its answer is the token's price, the share's times the token's multiplier, as
/// the Chainlink feeds of Robinhood Chain's Stock Tokens are; without one, the share's.
/// @dev The payload is read as RedStone's contracts read it, from the end of the calldata: the packages, each of the
/// feed's values signed by one of the data service's five signers, and the value is the median of the first three
/// unique signers' values. The signers and the three are the data service's at deployment. A package may be at most
/// three minutes old and one minute ahead of the block.
contract RedStoneRoundFeed is AggregatorV3Interface {
    struct Round {
        uint128 answer;
        uint64 time;
    }

    /// @notice The payload does not end in RedStone's marker, or is cut short.
    error InvalidPayload();
    /// @notice A package is signed by an address that is not one of the data service's signers.
    error SignerNotAuthorised(address signer);
    /// @notice A signature is malformed, or its signer cannot be recovered.
    error InvalidSignature();
    /// @notice A data point's value is wider than 32 bytes.
    error ValueTooWide(uint256 bytesWide);
    /// @notice The packages do not all carry one timestamp, or one is zero.
    error TimestampsMustMatchAndBePositive();
    /// @notice The feed has fewer unique signers than the data service requires.
    error InsufficientNumberOfUniqueSigners(uint256 received, uint256 required);
    /// @notice The packages are older than three minutes.
    error TimestampIsTooOld(uint256 timestamp, uint256 blockTimestamp);
    /// @notice The packages are more than a minute ahead of the block.
    error TimestampIsFromTheFuture(uint256 timestamp, uint256 blockTimestamp);
    /// @notice The packages are not newer than the feed's last round.
    error NotNewer(uint256 timestamp, uint256 lastTimestamp);
    /// @notice The feed's answer is more than 128 bits wide.
    error AnswerTooLarge();

    /// @notice A round was recorded.
    event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);

    uint256 private constant SIGNERS_REQUIRED = 3;
    uint256 private constant MAX_AGE = 3 minutes;
    uint256 private constant MAX_AHEAD = 1 minutes;
    uint256 private constant PHASE = 1 << 64;
    uint256 private constant MARKER_SIZE = 9;
    uint256 private constant SIGNATURE_SIZE = 65;
    uint256 private constant HALF_ORDER = 0x7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0;
    // slither-disable-next-line too-many-digits
    bytes9 private constant MARKER = 0x000002ed57011e0000;

    address private constant SIGNER_0 = 0x8BB8F32Df04c8b654987DAaeD53D6B6091e3B774;
    address private constant SIGNER_1 = 0xdEB22f54738d54976C4c0fe5ce6d408E40d88499;
    address private constant SIGNER_2 = 0x51Ce04Be4b3E32572C4Ec9135221d0691Ba7d202;
    address private constant SIGNER_3 = 0xDD682daEC5A90dD295d14DA4b0bec9281017b5bE;
    address private constant SIGNER_4 = 0x9c5AE89C4Af6aA32cE58588DBaF90d18a855B6de;

    /// @notice The RedStone data feed this feed follows, as RedStone names it in bytes32.
    bytes32 public immutable dataFeedId;
    /// @notice The Stock Token whose multiplier scales the answer; zero for the share's price.
    IMultiplier public immutable stockToken;
    /// @notice The feed's number of decimals.
    uint8 public constant override decimals = 8;
    /// @notice The feed's version.
    uint256 public constant override version = 1;

    string private _description;
    uint64 private _latest;
    mapping(uint64 index => Round) private _rounds;

    constructor(bytes32 dataFeedId_, string memory description_, IMultiplier stockToken_) {
        dataFeedId = dataFeedId_;
        stockToken = stockToken_;
        _description = description_;
    }

    /// @inheritdoc AggregatorV3Interface
    function description() external view override returns (string memory) {
        return _description;
    }

    /// @notice Records the data feed's value in the payload appended to the call as the next round.
    /// @return roundId The new round.
    // slither-disable-next-line timestamp
    function update() external returns (uint80 roundId) {
        (uint256 answer, uint256 timestampMs) = _verify(msg.data[4:]);
        if (address(stockToken) != address(0)) answer = answer * stockToken.uiMultiplier() / 1e18;
        uint256 time = timestampMs / 1000;
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp > time + MAX_AGE) revert TimestampIsTooOld(time, block.timestamp);
        // forge-lint: disable-next-line(block-timestamp)
        if (time > block.timestamp + MAX_AHEAD) revert TimestampIsFromTheFuture(time, block.timestamp);
        uint64 last = _latest;
        if (time <= _rounds[last].time) revert NotNewer(time, _rounds[last].time);
        if (answer > type(uint128).max) revert AnswerTooLarge();
        // forge-lint: disable-next-line(unsafe-typecast)
        _rounds[++last] = Round(uint128(answer), uint64(time));
        _latest = last;
        // forge-lint: disable-next-line(unsafe-typecast)
        roundId = uint80(PHASE | last);
        // forge-lint: disable-next-line(unsafe-typecast)
        emit AnswerUpdated(int256(answer), roundId, time);
    }

    /// @inheritdoc AggregatorV3Interface
    function getRoundData(uint80 roundId)
        public
        view
        override
        returns (uint80, int256 answer, uint256 startedAt, uint256 updatedAt, uint80)
    {
        // forge-lint: disable-next-line(unsafe-typecast)
        uint64 index = uint64(roundId);
        if (roundId >> 64 != 1 || index == 0 || index > _latest) revert("No data present");
        Round memory round = _rounds[index];
        // forge-lint: disable-next-line(unsafe-typecast)
        return (roundId, int256(uint256(round.answer)), round.time, round.time, roundId);
    }

    /// @inheritdoc AggregatorV3Interface
    function latestRoundData() external view override returns (uint80, int256, uint256, uint256, uint80) {
        // forge-lint: disable-next-line(unsafe-typecast)
        return getRoundData(uint80(PHASE | _latest));
    }

    // forge-lint: disable-start(require-revert-in-loop)
    struct Collected {
        uint256[3] values;
        uint256 count;
        uint256 seen;
        uint256 timestampMs;
    }

    /// @dev The median of the first three unique signers' values of `dataFeedId` in `payload`, and the packages'
    /// timestamp in milliseconds.
    function _verify(bytes calldata payload) private view returns (uint256 median, uint256 timestampMs) {
        if (payload.length < MARKER_SIZE + 3 + 2 || bytes9(payload[payload.length - MARKER_SIZE:]) != MARKER) {
            revert InvalidPayload();
        }
        uint256 end = payload.length - MARKER_SIZE;
        uint256 metadata = _read(payload, end, 3);
        end -= 3;
        if (metadata > end || end - metadata < 2) revert InvalidPayload();
        end -= metadata;
        uint256 packages = _read(payload, end, 2);
        end -= 2;

        Collected memory found = Collected([uint256(0), 0, 0], 0, 0, 0);
        for (uint256 p; p < packages; ++p) {
            end = _collect(payload, end, found);
        }
        if (found.count < SIGNERS_REQUIRED) revert InsufficientNumberOfUniqueSigners(found.count, SIGNERS_REQUIRED);
        (uint256 a, uint256 b, uint256 c) = (found.values[0], found.values[1], found.values[2]);
        median = a > b ? (b > c ? b : (a > c ? c : a)) : (a > c ? a : (b > c ? c : b));
        timestampMs = found.timestampMs;
    }

    /// @dev Reads the package that ends at `end`, keeps the feed's value in it if its signer is new, and returns where
    /// the package before it ends.
    function _collect(bytes calldata payload, uint256 end, Collected memory found) private view returns (uint256) {
        if (end < SIGNATURE_SIZE + 13) revert InvalidPayload();
        uint256 signatureStart = end - SIGNATURE_SIZE;
        uint256 pointsCount = _read(payload, signatureStart, 3);
        uint256 valueSize = _read(payload, signatureStart - 3, 4);
        uint256 packageTime = _read(payload, signatureStart - 7, 6);
        uint256 pointSize = valueSize + 32;
        uint256 pointsStart = signatureStart - 13;
        if (pointsCount * pointSize > pointsStart) revert InvalidPayload();
        pointsStart -= pointsCount * pointSize;
        uint256 bit = 1 << _signerIndex(_recover(payload, pointsStart, signatureStart));
        if (valueSize > 32) revert ValueTooWide(valueSize);
        for (uint256 i = pointsCount; i > 0; --i) {
            uint256 offset = pointsStart + (i - 1) * pointSize;
            if (bytes32(payload[offset:offset + 32]) != dataFeedId) continue;
            if (found.seen & bit == 0 && found.count < SIGNERS_REQUIRED) {
                found.values[found.count++] = _read(payload, offset + pointSize, valueSize);
                found.seen |= bit;
            }
        }
        if (packageTime == 0 || (found.timestampMs != 0 && packageTime != found.timestampMs)) {
            revert TimestampsMustMatchAndBePositive();
        }
        found.timestampMs = packageTime;
        return pointsStart;
    }

    /// @dev The signer of the package whose signed bytes start at `start` and whose signature starts at `signatureStart`.
    function _recover(bytes calldata payload, uint256 start, uint256 signatureStart)
        private
        pure
        returns (address signer)
    {
        bytes32 hash = keccak256(payload[start:signatureStart]);
        bytes32 r = bytes32(payload[signatureStart:signatureStart + 32]);
        bytes32 s = bytes32(payload[signatureStart + 32:signatureStart + 64]);
        uint8 v = uint8(payload[signatureStart + 64]);
        if ((v != 27 && v != 28) || uint256(s) > HALF_ORDER) revert InvalidSignature();
        signer = ecrecover(hash, v, r, s);
        if (signer == address(0)) revert InvalidSignature();
    }

    function _signerIndex(address signer) private pure returns (uint256) {
        if (signer == SIGNER_0) return 0;
        if (signer == SIGNER_1) return 1;
        if (signer == SIGNER_2) return 2;
        if (signer == SIGNER_3) return 3;
        if (signer == SIGNER_4) return 4;
        revert SignerNotAuthorised(signer);
    }

    /// @dev The big-endian unsigned integer of `size` bytes that ends at `end`.
    function _read(bytes calldata payload, uint256 end, uint256 size) private pure returns (uint256 value) {
        if (size > end) revert InvalidPayload();
        for (uint256 i = end - size; i < end; ++i) {
            value = (value << 8) | uint8(payload[i]);
        }
    }
    // forge-lint: disable-end(require-revert-in-loop)
}
