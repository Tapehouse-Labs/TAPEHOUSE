// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

import {IBand} from "../../src/interfaces/IBand.sol";

contract BandDouble is IBand {
    struct Quote {
        uint8 state;
        uint8 live;
        uint64 mid;
        uint64 halfBps;
        uint64 low;
        uint128 high;
    }

    struct Session {
        uint8 state;
        uint8 nyse;
        uint8 nyseNext;
        uint64 changeMs;
        uint64 boundaryMs;
    }

    struct Halt {
        bool signedHalt;
        uint64 until;
        uint64 issuedAt;
        bool oraclePaused;
    }

    struct Asset {
        address chainlinkFeed;
        bytes32 redstoneFeedId;
        bytes32 indexFeedId;
        address token;
    }

    mapping(bytes32 => Quote) public quotes;
    mapping(bytes32 => Halt) public halts;
    mapping(bytes32 => Asset) public assets;
    mapping(bytes32 => uint128) public variances;
    Session public currentSession;
    address public sequencer;
    bool public settled = true;

    function setQuote(bytes32 symbol, Quote calldata q) external {
        quotes[symbol] = q;
    }

    function setSession(Session calldata s) external {
        currentSession = s;
    }

    function setHalt(bytes32 symbol, Halt calldata h) external {
        halts[symbol] = h;
    }

    function setAsset(bytes32 symbol, Asset calldata a) external {
        assets[symbol] = a;
    }

    function setVariance(bytes32 feedId, uint128 v) external {
        variances[feedId] = v;
    }

    function setSequencer(address feed, bool settled_) external {
        sequencer = feed;
        settled = settled_;
    }

    function quote(bytes32 symbol) external view returns (uint8, uint8, uint64, uint64, uint64, uint128) {
        Quote memory q = quotes[symbol];
        return (q.state, q.live, q.mid, q.halfBps, q.low, q.high);
    }

    function session() external view returns (uint8, uint8, uint8, uint64, uint64) {
        Session memory s = currentSession;
        return (s.state, s.nyse, s.nyseNext, s.changeMs, s.boundaryMs);
    }

    function halt(bytes32 symbol) external view returns (bool, uint64, uint64, bool) {
        Halt memory h = halts[symbol];
        return (h.signedHalt, h.until, h.issuedAt, h.oraclePaused);
    }

    function asset(bytes32 symbol) external view returns (address, bytes32, bytes32, address) {
        Asset memory a = assets[symbol];
        return (a.chainlinkFeed, a.redstoneFeedId, a.indexFeedId, a.token);
    }

    function variance(bytes32 feedId) external view returns (uint128) {
        return variances[feedId];
    }

    function chainConfig() external view returns (address, bool) {
        return (sequencer, sequencer != address(0));
    }

    function sequencerSettled() external view returns (bool) {
        return settled;
    }
}
