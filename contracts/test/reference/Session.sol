// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

/// @notice The band's session, the regime of the current requirement and the requirement itself, as the Stylus
/// program derives them.
library Session {
    uint8 internal constant UNKNOWN = 0;
    uint8 internal constant CLOSED = 1;
    uint8 internal constant OPEN = 2;
    uint8 internal constant CLOSING = 3;
    uint64 internal constant HORIZON = 172_800;
    uint64 internal constant RAMP_MS = 25_200_000;

    /// @notice `band`'s `session()`: whether it answered, then `(open, nyse, nyseNext, changeMs, boundaryMs)`.
    function read(address band)
        internal
        view
        returns (bool ok, uint256 open, uint256 nyse, uint256 nyseNext, uint256 changeMs, uint256 boundaryMs)
    {
        (bool success, bytes memory data) = band.staticcall(abi.encodeWithSignature("session()"));
        if (!success || data.length < 160) return (false, 0, 0, 0, 0, 0);
        (open, nyse, nyseNext, changeMs, boundaryMs) = abi.decode(data, (uint256, uint256, uint256, uint256, uint256));
        ok = open <= type(uint8).max && nyse <= type(uint8).max && nyseNext <= type(uint8).max
            && changeMs <= type(uint64).max && boundaryMs <= type(uint64).max;
    }

    /// @notice The regime at `nowMs`, and how far into the ramp a closing session is.
    function regime(bool known, uint256 open, uint256 boundaryMs, uint256 nowMs)
        internal
        pure
        returns (uint8 code, uint256 elapsedMs)
    {
        if (!known) return (UNKNOWN, 0);
        if (open == 1) return (CLOSED, 0);
        if (open != 2) return (UNKNOWN, 0);
        if (boundaryMs != 0 && boundaryMs <= nowMs) return (CLOSED, 0);
        uint256 remaining = boundaryMs > nowMs ? boundaryMs - nowMs : 0;
        if (boundaryMs != 0 && remaining < RAMP_MS) return (CLOSING, RAMP_MS - remaining);
        return (OPEN, 0);
    }

    /// @notice The current requirement from the open and closed requirements and the leverage floor.
    function current(uint256 open, uint256 closed, uint256 floor, uint8 code, uint256 elapsedMs)
        internal
        pure
        returns (uint256)
    {
        uint256 buffered = (open * 5 + 3) / 4;
        uint256 across = buffered;
        if (closed > across) across = closed;
        if (floor > across) across = floor;
        if (code == OPEN) return buffered;
        if (code == CLOSING) return buffered + (across - buffered) * elapsedMs / RAMP_MS;
        return across;
    }
}
