// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract StockLendingHandler is Test {
    bytes32 internal constant CROSS = bytes32(0);
    StubStockToken internal immutable nvda;
    MarginAccounts internal immutable accounts;
    StockLendingVault internal immutable lending;
    address internal immutable borrower;
    address[2] internal holders = [makeAddr("alice"), makeAddr("bob")];
    bool public borrowedPastTheCap;
    uint256 public burnt;

    constructor(StubStockToken nvda_, MarginAccounts accounts_, StockLendingVault lending_, address borrower_) {
        (nvda, accounts, lending, borrower) = (nvda_, accounts_, lending_, borrower_);
    }

    function holder(uint256 i) external view returns (address) {
        return holders[i];
    }

    function deposit(uint256 who, uint256 amount) external {
        address account = holders[who % 2];
        amount = bound(amount, 1, 100e18);
        nvda.mint(account, amount);
        vm.startPrank(account);
        nvda.approve(address(accounts), amount);
        accounts.deposit(CROSS, address(nvda), amount, account);
        vm.stopPrank();
    }

    function withdraw(uint256 who, uint256 amount) external {
        address account = holders[who % 2];
        uint256 held = accounts.collateral(account, CROSS, address(nvda));
        if (held == 0) return;
        vm.prank(account);
        accounts.withdraw(CROSS, address(nvda), bound(amount, 1, held), account, account);
    }

    function lend(uint256 who, uint256 amount) external {
        address account = holders[who % 2];
        uint256 held = accounts.collateral(account, CROSS, address(nvda));
        if (held == 0) return;
        vm.prank(account);
        accounts.lend(CROSS, address(nvda), bound(amount, 1, held), account);
    }

    function unlend(uint256 who, uint256 amount) external {
        address account = holders[who % 2];
        uint256 reachable =
            accounts.sellable(account, CROSS, address(nvda)) - accounts.collateral(account, CROSS, address(nvda));
        if (reachable == 0) return;
        vm.prank(account);
        accounts.unlend(CROSS, address(nvda), bound(amount, 1, reachable), account);
    }

    function borrow(uint256 amount) external {
        uint256 available = lending.borrowable();
        if (available == 0) return;
        amount = bound(amount, 1, available);
        vm.prank(borrower);
        lending.borrow(amount, borrower);
        if (lending.debt() * 10_000 > lending.totalAssets() * 9_000 + 10_000) borrowedPastTheCap = true;
    }

    function repay(uint256 amount) external {
        uint256 owed = lending.debt();
        if (owed == 0) return;
        amount = bound(amount, 1, owed);
        nvda.mint(borrower, amount);
        vm.prank(borrower);
        lending.repay(amount);
    }

    function wait(uint256 time) external {
        skip(bound(time, 1, 30 days));
    }

    function burn(uint256 amount) external {
        uint256 idle = lending.idle();
        if (idle == 0) return;
        amount = bound(amount, 1, idle);
        nvda.adminBurn(address(lending), amount);
        lending.sync();
        burnt += amount;
    }
}

/// forge-config: default.invariant.depth = 64
contract StockLendingInvariantTest is Test {
    StubStockToken internal nvda;
    MarginAccounts internal accounts;
    StockLendingVault internal lending;
    StockLendingHandler internal handler;

    function setUp() public {
        vm.warp(1_790_000_000);
        StubUsdg usdg = new StubUsdg();
        nvda = new StubStockToken(1e18);
        BandDouble band = new BandDouble();
        band.setAsset("NVDA", BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setQuote("NVDA", BandDouble.Quote(3, 3, 200e8, 50, 200e8, 202e8));
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = "NVDA";
        MarginDouble engine =
            new MarginDouble(symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH")));
        SupplyVault vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](1);
        caps[0] = type(uint256).max;
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            new StubToken(18),
            address(this),
            caps,
            type(uint256).max,
            type(uint256).max,
            0,
            10_00
        );
        lending = new StockLendingVault(
            IERC20(address(nvda)), address(this), SupplyVault.RateModel(80_00, 25, 1_00, 50_00), 10_00
        );
        address borrower = makeAddr("borrower");
        lending.setDepositor(address(accounts));
        lending.setBorrower(borrower);
        accounts.setLending("NVDA", lending);
        vm.prank(borrower);
        nvda.approve(address(lending), type(uint256).max);
        handler = new StockLendingHandler(nvda, accounts, lending, borrower);
        targetContract(address(handler));
    }

    function invariant_LendingNeverReadsAsABurnOfTheAccountsHolding() public view {
        (uint256 units, uint256 scale,) = accounts.holding("NVDA");
        assertEq(scale, 1e18);
        assertGe(nvda.balanceOf(address(accounts)), units * scale / 1e18);
    }

    function invariant_TheVaultHoldsTheTokensItCounts() public view {
        assertGe(nvda.balanceOf(address(lending)), lending.idle());
        assertEq(lending.totalAssets(), lending.idle() + lending.debt());
    }

    function invariant_ThePositionsLendWhatTheAccountsHoldOfTheVault() public view {
        uint256 lent;
        for (uint256 i; i < 2; ++i) {
            lent += accounts.lent(handler.holder(i), bytes32(0), address(nvda));
        }
        uint256 held = lending.convertToAssets(lending.balanceOf(address(accounts)));
        assertLe(lent, held);
        assertGe(lent + 2, held);
    }

    function invariant_TheBorrowerNeverTakesMoreThanNinetyPercent() public view {
        assertFalse(handler.borrowedPastTheCap());
    }
}
