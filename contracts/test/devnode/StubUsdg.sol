// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

/// @notice USDG as the conformance tests pin it: six decimals, EIP-2612 permits under "Global Dollar" version 1,
/// a pause that stops transfers, approvals and permits with `ContractPaused()`, and frozen addresses that
/// neither send nor receive (`AddressFrozen()`) and whose balance the operator can wipe. Anyone may mint. Test code
/// only.
contract StubUsdg is ERC20Permit {
    address public immutable operator;
    bool public paused;
    mapping(address => bool) public isFrozen;

    error ContractPaused();
    error AddressFrozen();
    error AddressNotFrozen();
    error NotOperator();

    constructor() ERC20("Global Dollar", "USDG") ERC20Permit("Global Dollar") {
        operator = msg.sender;
    }

    modifier onlyOperator() {
        if (msg.sender != operator) revert NotOperator();
        _;
    }

    function decimals() public pure override returns (uint8) {
        return 6;
    }

    function mint(address to, uint256 value) external {
        _mint(to, value);
    }

    function pause() external onlyOperator {
        paused = true;
    }

    function unpause() external onlyOperator {
        paused = false;
    }

    function freeze(address account) external onlyOperator {
        isFrozen[account] = true;
    }

    function unfreeze(address account) external onlyOperator {
        isFrozen[account] = false;
    }

    function wipeFrozenAddress(address account) external onlyOperator {
        if (!isFrozen[account]) revert AddressNotFrozen();
        _burn(account, balanceOf(account));
    }

    function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        public
        override
    {
        if (paused) revert ContractPaused();
        super.permit(owner, spender, value, deadline, v, r, s);
    }

    function _update(address from, address to, uint256 value) internal override {
        bool minting = from == address(0);
        bool burning = to == address(0);
        if (paused && !minting && !burning) revert ContractPaused();
        if (!burning && isFrozen[to] || !minting && !burning && isFrozen[from]) revert AddressFrozen();
        super._update(from, to, value);
    }

    function _approve(address owner, address spender, uint256 value, bool emitEvent) internal override {
        if (paused) revert ContractPaused();
        super._approve(owner, spender, value, emitEvent);
    }
}
