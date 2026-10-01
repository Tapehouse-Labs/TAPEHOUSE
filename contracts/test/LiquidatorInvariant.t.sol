// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
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

contract LiquidatorHandler is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal immutable usdg;
    StubStockToken internal immutable nvda;
    BandDouble internal immutable band;
    MarginAccounts internal immutable accounts;
    Liquidator internal immutable liquidator;
    address[3] public actors = [address(0xA1), address(0xA2), address(0xA3)];
    bool public paidAtMostTheAsk = true;

    constructor(
        StubUsdg usdg_,
        StubStockToken nvda_,
        BandDouble band_,
        MarginAccounts accounts_,
        Liquidator liquidator_
    ) {
        (usdg, nvda, band, accounts, liquidator) = (usdg_, nvda_, band_, accounts_, liquidator_);
        usdg.approve(address(liquidator_), type(uint256).max);
    }

    function deposit(uint256 actor, bool isolated, uint256 amount) external {
        address account = actors[actor % 3];
        amount = bound(amount, 1, 100e18);
        nvda.mint(account, amount);
        vm.startPrank(account);
        nvda.approve(address(accounts), amount);
        accounts.deposit(isolated ? NVDA : bytes32(0), address(nvda), amount, account);
        vm.stopPrank();
    }

    function depositCash(uint256 actor, bool isolated, uint256 amount) external {
        address account = actors[actor % 3];
        amount = bound(amount, 1, 10_000e6);
        usdg.mint(account, amount);
        vm.startPrank(account);
        usdg.approve(address(accounts), amount);
        accounts.deposit(isolated ? NVDA : bytes32(0), address(usdg), amount, account);
        vm.stopPrank();
    }

    function borrow(uint256 actor, bool isolated, uint256 amount) external {
        address account = actors[actor % 3];
        vm.prank(account);
        try accounts.borrow(isolated ? NVDA : bytes32(0), bound(amount, 1, 20_000e6), account, account) {} catch {}
    }

    function move(uint256 low, uint256 spread, uint256 state) external {
        low = bound(low, 1e8, 400e8);
        uint256 high = low + low * bound(spread, 0, 20_00) / 10_000;
        uint8 s = uint8(bound(state, 1, 3));
        band.setQuote(NVDA, BandDouble.Quote(s, 3, uint64((low + high) / 2), 50, uint64(low), uint128(high)));
    }

    function session(uint256 kind) external {
        kind %= 3;
        band.setSession(BandDouble.Session(uint8(kind == 0 ? 2 : kind == 1 ? 1 : 0), 1, 1, 0, 0));
    }

    function liquidate(uint256 actor, bool isolated, uint256 amount, uint256 wait) external {
        address account = actors[actor % 3];
        bytes32 position = isolated ? NVDA : bytes32(0);
        try liquidator.start(account, position) {} catch {}
        skip(bound(wait, 0, 2 hours));
        uint256 ask = liquidator.price(account, position, address(nvda));
        amount = bound(amount, 1, 100e18);
        usdg.mint(address(this), 1_000_000e6);
        uint256 before = usdg.balanceOf(address(this));
        try liquidator.buy(account, position, address(nvda), amount, type(uint256).max, address(this)) returns (
            uint256 bought, uint256
        ) {
            uint256 paid = before - usdg.balanceOf(address(this));
            if (paid > (bought * ask - 1) / 1e20 + 1) paidAtMostTheAsk = false;
        } catch {}
    }

    function settleCash(uint256 actor, bool isolated) external {
        try liquidator.settleCash(actors[actor % 3], isolated ? NVDA : bytes32(0)) {} catch {}
    }

    function writeOff(uint256 actor, bool isolated) external {
        try liquidator.writeOff(actors[actor % 3], isolated ? NVDA : bytes32(0)) {} catch {}
    }

    function stop(uint256 actor, bool isolated) external {
        try liquidator.stop(actors[actor % 3], isolated ? NVDA : bytes32(0)) {} catch {}
    }
}

contract LiquidatorInvariantTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal usdg;
    StubStockToken internal nvda;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    Liquidator internal liquidator;
    LiquidatorHandler internal handler;

    function setUp() public {
        usdg = new StubUsdg();
        nvda = new StubStockToken(1e18);
        BandDouble band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = NVDA;
        StubAggregator ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        MarginDouble engine = new MarginDouble(symbols, address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
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
            5_00,
            10_00
        );
        liquidator = new Liquidator(accounts);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        usdg.mint(address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
        handler = new LiquidatorHandler(usdg, nvda, band, accounts, liquidator);
        targetContract(address(handler));
    }

    function invariant_TheLiquidatorKeepsNothing() public view {
        assertEq(usdg.balanceOf(address(liquidator)), 0);
        assertEq(nvda.balanceOf(address(liquidator)), 0);
    }

    function invariant_TheAccountsHoldTheirUsdg() public view {
        uint256 cash;
        for (uint256 i; i < 3; ++i) {
            address a = handler.actors(i);
            cash += accounts.collateral(a, bytes32(0), address(usdg)) + accounts.collateral(a, NVDA, address(usdg));
        }
        assertEq(usdg.balanceOf(address(accounts)), cash + accounts.reserve() + accounts.backstopPremium());
    }

    function invariant_ThePositionsOweWhatTheVaultIsOwed() public view {
        uint256 owed;
        for (uint256 i; i < 3; ++i) {
            owed += accounts.debt(handler.actors(i), bytes32(0)) + accounts.debt(handler.actors(i), NVDA);
        }
        uint256 debt = vault.debt();
        uint256 virtualPart = (debt + 1) * 1e6 / (accounts.totalDebtShares() + 1e6) + 1;
        assertGe(owed + virtualPart, debt);
        assertLe(owed, debt + 6);
    }

    function invariant_TheAccountsHoldTheStockTheyCount() public view {
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertGe(nvda.balanceOf(address(accounts)), units * scale / 1e18);
    }

    function invariant_ABuyerPaysAtMostTheAsk() public view {
        assertTrue(handler.paidAtMostTheAsk());
    }
}
