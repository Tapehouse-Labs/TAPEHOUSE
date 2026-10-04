// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

contract TargetDouble {
    error Refused();

    mapping(address caller => uint256) public calls;

    function touch() external {
        ++calls[msg.sender];
    }

    function refuse() external pure {
        revert Refused();
    }
}
