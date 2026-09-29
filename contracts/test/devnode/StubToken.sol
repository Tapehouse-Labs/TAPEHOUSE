// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.37;

contract StubToken {
    uint8 public immutable decimals;

    constructor(uint8 decimals_) {
        decimals = decimals_;
    }
}
