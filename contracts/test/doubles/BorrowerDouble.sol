// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IStockLendingBorrower} from "../../src/interfaces/IStockLendingBorrower.sol";
import {StockLendingVault} from "../../src/StockLendingVault.sol";
import {StubStockToken} from "../devnode/StubStockToken.sol";

/// @notice A stock lending vault's borrower that borrows on request and, on a buy-in, mints what it must return and
/// repays it, or returns only `toReturn` of it. Test code only.
contract BorrowerDouble is IStockLendingBorrower {
    StockLendingVault public immutable vault;
    StubStockToken public immutable token;
    uint256 public toReturn = type(uint256).max;

    constructor(StockLendingVault vault_, StubStockToken token_) {
        (vault, token) = (vault_, token_);
        token_.approve(address(vault_), type(uint256).max);
    }

    function borrow(uint256 assets) external {
        vault.borrow(assets, address(this));
    }

    function repay(uint256 assets) external returns (uint256) {
        token.mint(address(this), assets);
        return vault.repay(assets);
    }

    function setToReturn(uint256 assets) external {
        toReturn = assets;
    }

    function buyIn(uint256 assets) external {
        require(msg.sender == address(vault));
        uint256 back = assets < toReturn ? assets : toReturn;
        token.mint(address(this), back);
        vault.repay(back);
    }
}
