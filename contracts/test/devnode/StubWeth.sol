// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {StubToken} from "./StubToken.sol";

/// @notice The dev node's WETH: an 18-decimal `StubToken` that wraps the ether sent to it one for one, for the caller
/// or for another account, as Arbitrum's WETH does. Test code only.
contract StubWeth is StubToken(18) {
    function deposit() external payable {
        _mint(msg.sender, msg.value);
    }

    function depositTo(address account) external payable {
        _mint(account, msg.value);
    }
}
