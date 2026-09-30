// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @notice An ERC-20 with the decimals it is given, such as WETH's. Anyone may mint. Test code only.
contract StubToken is ERC20 {
    uint8 internal immutable _decimals;

    constructor(uint8 decimals_) ERC20("Stub Token", "STUB") {
        _decimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    function mint(address to, uint256 value) external {
        _mint(to, value);
    }
}
