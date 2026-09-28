// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

/// @title Tapehouse band
/// @notice The Stylus program that prices each Robinhood Stock Token as a band, as its Solidity readers
/// call it. Prices have 8 decimals.
interface IBand {
    /// @notice The band of `symbol`: its state (0 halted, 1 degraded, 2 closed, 3 open), how many legs are
    /// live, its centre, half-width in basis points and bounds. All zero while halted.
    function quote(bytes32 symbol)
        external
        view
        returns (uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high);

    /// @notice Chainlink's 24/5 session (0 not known, 1 closed, 2 open), NYSE's state and next state, when
    /// NYSE changes state, and the session's next boundary, in milliseconds.
    function session()
        external
        view
        returns (uint8 state, uint8 nyse, uint8 nyseNext, uint64 changeMs, uint64 boundaryMs);

    /// @notice The trading halt of `symbol`: whether a halt signed by Tapehouse's halt signer holds, until
    /// when, when its last message was issued, and whether the issuer has paused the token's oracle.
    function halt(bytes32 symbol)
        external
        view
        returns (bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused);

    /// @notice The Chainlink feed, RedStone feed ID, index feed ID and Stock Token of `symbol`.
    function asset(bytes32 symbol)
        external
        view
        returns (address chainlinkFeed, bytes32 redstoneFeedId, bytes32 indexFeedId, address token);

    /// @notice The EWMA variance of a 24/7 feed, in centi-basis-points squared per minute.
    function variance(bytes32 feedId) external view returns (uint128);

    /// @notice The L2 sequencer-uptime feed the band follows, and whether its Chainlink feeds follow NYSE
    /// regular hours.
    function chainConfig() external view returns (address sequencerUptimeFeed, bool chainlinkRegularHours);

    /// @notice Whether the L2 sequencer is up and has been for over an hour; true without a feed.
    function sequencerSettled() external view returns (bool);
}
