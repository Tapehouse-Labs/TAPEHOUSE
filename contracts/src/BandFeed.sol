// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {IUniswapV3Pool} from "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import {OracleLibrary} from "@uniswap/v3-periphery/contracts/libraries/OracleLibrary.sol";
import {IBand} from "./interfaces/IBand.sol";

/// @notice The ERC-20 decimals getter.
interface IERC20Decimals {
    function decimals() external view returns (uint8);
}

/// @title Tapehouse band feed
/// @notice One asset's Tapehouse band behind Chainlink's `AggregatorV3Interface`, its whole band with the
/// token's own Uniswap v3 market beside it, and the band sealed before each reopen. It answers one side of
/// the band, set at deployment; lending collateral reads the low side.
/// @dev The band is recomputed on every read, and its state already judges staleness, so a round is the
/// current block. A halted band, or an L2 sequencer that is down or back for an hour or less, has no
/// answer: the feed reverts rather than report a price of zero.
contract BandFeed is AggregatorV3Interface {
    /// @notice The side of the band this feed answers.
    enum Side {
        Low,
        Mid,
        High
    }

    /// @notice The whole band of the asset at one block.
    /// @param signedHalt Whether a halt signed by Tapehouse's halt signer holds: the one input the band
    /// takes on Tapehouse's own signature. `haltUntil` in the past and non-zero is a halt that lapsed
    /// without a lift.
    /// @param twap The token's 30-minute Uniswap v3 TWAP in its quote token (USDG on Robinhood Chain), with
    /// 8 decimals; `premiumBps` is its distance from the band's centre, the quote token taken at par. Both
    /// are zero and `twapValid` false without a pool or without 30 minutes of pool history.
    struct Band {
        uint8 state;
        uint8 live;
        uint64 mid;
        uint64 halfBps;
        uint64 low;
        uint128 high;
        uint128 variance;
        uint8 session;
        uint8 nyse;
        uint8 nyseNext;
        uint64 nyseChangeMs;
        uint64 sessionBoundaryMs;
        bool signedHalt;
        uint64 haltUntil;
        uint64 haltIssuedAt;
        bool oraclePaused;
        bool sequencerSettled;
        bool twapValid;
        uint256 twap;
        int256 premiumBps;
    }

    /// @notice The band sealed before a reopen, and when.
    struct Seal {
        uint8 state;
        uint8 live;
        uint64 mid;
        uint64 halfBps;
        uint64 low;
        uint128 high;
        uint64 sealedAt;
    }

    /// @notice The TWAP's window, in seconds.
    uint32 public constant TWAP_WINDOW = 1800;
    /// @notice How long before a reopen the band may be sealed, in milliseconds.
    uint64 public constant SEAL_WINDOW_MS = 600_000;
    /// @inheritdoc AggregatorV3Interface
    uint8 public constant override decimals = 8;
    /// @inheritdoc AggregatorV3Interface
    uint256 public constant override version = 1;

    /// @notice The band program.
    IBand public immutable band;
    /// @notice The asset's symbol in the band.
    bytes32 public immutable symbol;
    /// @notice The side of the band this feed answers.
    Side public immutable side;
    /// @notice The token's Uniswap v3 pool; zero for none.
    IUniswapV3Pool public immutable pool;

    bytes32 private immutable varianceFeedId;
    bool private immutable followsSequencer;
    address private immutable token;
    address private immutable quoteToken;
    uint128 private immutable tokenUnit;
    uint256 private immutable quoteUnit;

    /// @inheritdoc AggregatorV3Interface
    string public override description;

    /// @notice The band sealed before the reopen at `reopenMs`.
    mapping(uint64 reopenMs => Seal) public seals;

    /// @notice The band was sealed before the reopen at `reopenMs`.
    event Sealed(
        uint64 indexed reopenMs, uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high
    );

    /// @notice The band of `symbol` is halted.
    error NoAnswer(bytes32 symbol);
    /// @notice The L2 sequencer is down or came back an hour ago or less.
    error SequencerNotSettled();
    /// @notice Only the current round, the block's timestamp, has an answer.
    error NoRound(uint80 roundId);
    /// @notice A band is sealed only while the session is closed, within `SEAL_WINDOW_MS` of its reopen.
    error NotSealWindow(uint8 session, uint64 reopenMs);
    /// @notice The band does not configure `symbol`.
    error UnknownAsset(bytes32 symbol);
    /// @notice `pool` does not trade the asset's Stock Token.
    error PoolWithoutToken(address pool, address token);
    /// @notice `pool` has never been initialised, so it has no price history.
    error PoolNotInitialized(address pool);

    /// @param band_ The band program.
    /// @param symbol_ The asset's symbol in the band.
    /// @param side_ The side of the band to answer.
    /// @param pool_ An initialised Uniswap v3 pool of the asset's Stock Token, or zero.
    /// @param description_ What is priced: the Stock Token or, where the band has none, the share.
    constructor(IBand band_, bytes32 symbol_, Side side_, IUniswapV3Pool pool_, string memory description_) {
        (address feed, bytes32 feedId, bytes32 indexId, address token_) = band_.asset(symbol_);
        if (feed == address(0) && feedId == bytes32(0) && indexId == bytes32(0)) revert UnknownAsset(symbol_);
        // slither-disable-next-line unused-return
        (address sequencer,) = band_.chainConfig(); // forge-lint: disable-line(unused-return)
        band = band_;
        symbol = symbol_;
        side = side_;
        pool = pool_;
        description = description_;
        varianceFeedId = feedId != bytes32(0) ? feedId : indexId;
        followsSequencer = sequencer != address(0);
        if (address(pool_) != address(0)) {
            address token0 = pool_.token0();
            address token1 = pool_.token1();
            if (token0 != token_ && token1 != token_) {
                revert PoolWithoutToken(address(pool_), token_);
            }
            // slither-disable-next-line unused-return
            (,,, uint16 cardinality,,,) = pool_.slot0(); // forge-lint: disable-line(unused-return)
            if (cardinality == 0) revert PoolNotInitialized(address(pool_));
            address quote_ = token0 == token_ ? token1 : token0;
            token = token_;
            // slither-disable-next-line missing-zero-check
            quoteToken = quote_;
            tokenUnit = uint128(10 ** IERC20Decimals(token_).decimals());
            quoteUnit = 10 ** IERC20Decimals(quote_).decimals();
        }
    }

    /// @inheritdoc AggregatorV3Interface
    function latestRoundData() public view override returns (uint80, int256, uint256, uint256, uint80) {
        // slither-disable-next-line unused-return
        (uint8 state,, uint64 mid,, uint64 low, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if (state == 0) revert NoAnswer(symbol);
        if (followsSequencer && !band.sequencerSettled()) revert SequencerNotSettled();
        int256 answer =
            side == Side.Low ? int256(uint256(low)) : side == Side.Mid ? int256(uint256(mid)) : int256(uint256(high));
        // forge-lint: disable-next-line(unsafe-typecast)
        uint80 roundId = uint80(block.timestamp);
        return (roundId, answer, block.timestamp, block.timestamp, roundId);
    }

    /// @inheritdoc AggregatorV3Interface
    function getRoundData(uint80 roundId) external view override returns (uint80, int256, uint256, uint256, uint80) {
        // slither-disable-next-line timestamp
        if (roundId != block.timestamp) revert NoRound(roundId); // forge-lint: disable-line(block-timestamp)
        return latestRoundData();
    }

    /// @notice The whole band now, the token's market beside it, and both halt sources.
    function latestBand() external view returns (Band memory b) {
        (b.state, b.live, b.mid, b.halfBps, b.low, b.high) = band.quote(symbol);
        b.variance = band.variance(varianceFeedId);
        (b.session, b.nyse, b.nyseNext, b.nyseChangeMs, b.sessionBoundaryMs) = band.session();
        (b.signedHalt, b.haltUntil, b.haltIssuedAt, b.oraclePaused) = band.halt(symbol);
        b.sequencerSettled = !followsSequencer || band.sequencerSettled();
        (b.twapValid, b.twap) = twap();
        if (b.twapValid && b.mid != 0) {
            // forge-lint: disable-next-line(unsafe-typecast)
            b.premiumBps = (int256(b.twap) - int256(uint256(b.mid))) * 10_000 / int256(uint256(b.mid));
        }
    }

    /// @notice Seals the band before the session reopens. Anyone may call it while the session is closed and
    /// its reopen is at most `SEAL_WINDOW_MS` away; the latest call before the reopen stands. On a chain whose
    /// Chainlink feeds follow NYSE regular hours, the band sealed before the 24/5 reopen has its 24/7 leg alone.
    function seal() external {
        // slither-disable-next-line unused-return
        (uint8 session_,,,, uint64 reopenMs) = band.session(); // forge-lint: disable-line(unused-return)
        // forge-lint: disable-next-line(unsafe-typecast)
        uint64 nowMs = uint64(block.timestamp) * 1000;
        // slither-disable-next-line timestamp
        if (session_ != 1 || reopenMs <= nowMs || reopenMs - nowMs > SEAL_WINDOW_MS) {
            revert NotSealWindow(session_, reopenMs);
        }
        (uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high) = band.quote(symbol);
        // forge-lint: disable-next-line(unsafe-typecast)
        seals[reopenMs] = Seal(state, live, mid, halfBps, low, high, uint64(block.timestamp));
        emit Sealed(reopenMs, state, live, mid, halfBps, low, high);
    }

    function twap() private view returns (bool valid, uint256 price) {
        if (address(pool) == address(0) || OracleLibrary.getOldestObservationSecondsAgo(address(pool)) < TWAP_WINDOW) {
            return (valid, price);
        }
        // slither-disable-next-line unused-return
        (int24 tick,) = OracleLibrary.consult(address(pool), TWAP_WINDOW);
        price = OracleLibrary.getQuoteAtTick(tick, tokenUnit, token, quoteToken) * 10 ** decimals / quoteUnit;
        valid = true;
    }
}
