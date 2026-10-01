// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

/// @notice A band feed that answers only its band, its symbol and the seals it is given. Test code only.
contract SealDouble {
    struct Seal {
        uint8 state;
        uint8 live;
        uint64 mid;
        uint64 halfBps;
        uint64 low;
        uint128 high;
        uint64 sealedAt;
    }

    address public immutable band;
    bytes32 public immutable symbol;
    mapping(uint64 reopenMs => Seal) public seals;

    constructor(address band_, bytes32 symbol_) {
        (band, symbol) = (band_, symbol_);
    }

    function setSeal(uint64 reopenMs, Seal calldata s) external {
        seals[reopenMs] = s;
    }
}
