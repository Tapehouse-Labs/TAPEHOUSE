// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {ShortPositions} from "../src/ShortPositions.sol";
import {StockLendingVault} from "../src/StockLendingVault.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {IV3SwapRouter} from "../src/interfaces/IUniswapV3.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubSwapRouter} from "./devnode/StubSwapRouter.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract ShortPositionsHandler is Test {
    bytes32 internal constant CROSS = bytes32(0);
    StubUsdg internal immutable usdg;
    BandDouble internal immutable band;
    StubSwapRouter internal immutable router;
    MarginAccounts internal immutable accounts;
    ShortPositions internal immutable shorts;
    address internal immutable lender = makeAddr("lender");
    address[3] internal traders = [makeAddr("bob"), makeAddr("carol"), makeAddr("dave")];
    bytes32[2] public symbols = [bytes32("NVDA"), bytes32("SPY")];
    StubStockToken[2] internal tokens;
    StockLendingVault[2] internal vaults;
    bool public panicked;

    constructor(
        StubUsdg usdg_,
        BandDouble band_,
        StubSwapRouter router_,
        MarginAccounts accounts_,
        ShortPositions shorts_,
        StubStockToken[2] memory tokens_,
        StockLendingVault[2] memory vaults_
    ) {
        (usdg, band, router, accounts, shorts, tokens, vaults) =
        (usdg_, band_, router_, accounts_, shorts_, tokens_, vaults_);
        vm.startPrank(lender);
        for (uint256 t; t < 2; ++t) {
            tokens[t].mint(lender, 1_000e18);
            tokens[t].approve(address(accounts), type(uint256).max);
            accounts.deposit(CROSS, address(tokens[t]), 1_000e18, lender);
            accounts.lend(CROSS, address(tokens[t]), 1_000e18, lender);
        }
        vm.stopPrank();
        for (uint256 i; i < 3; ++i) {
            vm.prank(traders[i]);
            usdg.approve(address(shorts), type(uint256).max);
        }
    }

    function trader(uint256 i) external view returns (address) {
        return traders[i];
    }

    function deposit(uint256 who, uint256 which, uint256 amount) external {
        address account = traders[who % 3];
        bytes32 symbol = symbols[which % 2];
        (,, uint256 shares,) = shorts.position(account, symbol);
        if (shares != 0) return;
        amount = bound(amount, 1, 5_000e6);
        usdg.mint(account, amount);
        vm.prank(account);
        try shorts.deposit(symbol, amount, account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function withdraw(uint256 who, uint256 which, uint256 amount) external {
        address account = traders[who % 3];
        bytes32 symbol = symbols[which % 2];
        (int256 held,,,) = shorts.position(account, symbol);
        if (held <= 0) return;
        vm.prank(account);
        try shorts.withdraw(symbol, bound(amount, 1, uint256(held)), account, account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function sell(uint256 who, uint256 which, uint256 amount) external {
        address account = traders[who % 3];
        uint256 t = which % 2;
        uint256 available = vaults[t].borrowable();
        if (available == 0) return;
        vm.prank(account);
        try shorts.sell(symbols[t], bound(amount, 1, available), 0, account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function cover(uint256 who, uint256 which, uint256 amount) external {
        address account = traders[who % 3];
        bytes32 symbol = symbols[which % 2];
        (, uint256 debt,,) = shorts.position(account, symbol);
        if (debt == 0) return;
        vm.prank(account);
        try shorts.cover(symbol, bound(amount, 1, debt + 1), type(uint256).max, account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function move(uint256 which, uint256 percent, uint256 offset) external {
        uint256 t = which % 2;
        (,, uint64 now_,,,) = band.quote(symbols[t]);
        uint64 mid = uint64(bound(uint256(now_) * bound(percent, 70, 160) / 100, 50e8, 2_000e8));
        band.setQuote(symbols[t], BandDouble.Quote(3, 3, mid, 1_00, mid - mid / 100, mid + mid / 100));
        router.setPrice(address(tokens[t]), uint256(mid) * bound(offset, 97, 103) / 100);
    }

    function lend(uint256 which, uint256 amount) external {
        uint256 t = which % 2;
        amount = bound(amount, 1e18, 500e18);
        tokens[t].mint(lender, amount);
        vm.startPrank(lender);
        accounts.deposit(CROSS, address(tokens[t]), amount, lender);
        accounts.lend(CROSS, address(tokens[t]), amount, lender);
        vm.stopPrank();
    }

    function recall(uint256 which, uint256 amount) external {
        address token = address(tokens[which % 2]);
        uint256 lent = accounts.lent(lender, CROSS, token) - accounts.claim(lender, CROSS, token);
        if (lent < 1e16 || accounts.recalls(lender, CROSS, token).length >= 3) return;
        vm.prank(lender);
        try accounts.recall(CROSS, token, amount % 4 == 0 ? lent : bound(amount, 1e16, lent), lender) {} catch {}
    }

    function buyIn(uint256 which) external {
        uint256 t = which % 2;
        uint256 tickets = vaults[t].tickets();
        if (tickets == 0) return;
        vm.warp(block.timestamp + 1 days);
        try vaults[t].buyIn(tickets - 1) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
        try accounts.settle(lender, CROSS, address(tokens[t])) {} catch {}
    }

    function liquidate(uint256 who, uint256 which) external {
        try shorts.liquidate(traders[who % 3], symbols[which % 2]) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function wait(uint256 seconds_) external {
        vm.warp(block.timestamp + bound(seconds_, 1, 30 days));
    }
}

contract ShortPositionsInvariantTest is Test {
    StubUsdg internal usdg;
    ShortPositions internal shorts;
    StockLendingVault[2] internal vaults;
    SupplyVault internal vault;
    ShortPositionsHandler internal handler;
    uint256 internal supplied;

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        StubStockToken[2] memory tokens = [new StubStockToken(1e18), new StubStockToken(1e18)];
        BandDouble band = new BandDouble();
        bytes32[] memory symbols = new bytes32[](2);
        (symbols[0], symbols[1]) = ("NVDA", "SPY");
        StubSwapRouter router = new StubSwapRouter(address(usdg));
        for (uint256 t; t < 2; ++t) {
            uint64 mid = t == 0 ? 200e8 : 600e8;
            band.setAsset(symbols[t], BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(tokens[t])));
            band.setQuote(symbols[t], BandDouble.Quote(3, 3, mid, 1_00, mid - mid / 100, mid + mid / 100));
            router.setPrice(address(tokens[t]), mid);
        }
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        MarginDouble engine =
            new MarginDouble(symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH")));
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](2);
        (caps[0], caps[1]) = (type(uint256).max, type(uint256).max);
        MarginAccounts accounts = new MarginAccounts(
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
        vault.setBorrower(address(accounts));
        uint24[] memory fees = new uint24[](2);
        (fees[0], fees[1]) = (500, 500);
        shorts = new ShortPositions(accounts, IV3SwapRouter(address(router)), fees);
        for (uint256 t; t < 2; ++t) {
            vaults[t] = new StockLendingVault(
                IERC20(address(tokens[t])),
                address(this),
                t == 0 ? SupplyVault.RateModel(80_00, 25, 1_00, 50_00) : SupplyVault.RateModel(80_00, 0, 0, 0),
                10_00
            );
            vaults[t].setDepositor(address(accounts));
            vaults[t].setBorrower(address(shorts));
            accounts.setLending(symbols[t], vaults[t]);
        }
        handler = new ShortPositionsHandler(usdg, band, router, accounts, shorts, tokens, vaults);
        supplied = vault.totalAssets();
        targetContract(address(handler));
    }

    function invariant_TheShortsHoldTheUsdgTheyCount() public view {
        uint256 counted;
        for (uint256 t; t < 2; ++t) {
            bytes32 symbol = handler.symbols(t);
            for (uint256 e; e <= shorts.epoch(symbol); ++e) {
                (, uint256 booked,) = shorts.book(symbol, e);
                counted += booked;
            }
            for (uint256 i; i < 3; ++i) {
                (int256 held,, uint256 shares,) = shorts.position(handler.trader(i), symbol);
                if (shares == 0) counted += uint256(held);
            }
        }
        assertEq(usdg.balanceOf(address(shorts)), counted);
    }

    function invariant_EachBookCoversItsShortsAndCountsTheirShares() public view {
        for (uint256 t; t < 2; ++t) {
            bytes32 symbol = handler.symbols(t);
            for (uint256 e; e <= shorts.epoch(symbol); ++e) {
                (uint256 total, uint256 booked,) = shorts.book(symbol, e);
                int256 held;
                uint256 shares;
                for (uint256 i; i < 3; ++i) {
                    (int256 h,, uint256 s, uint256 epoch) = shorts.position(handler.trader(i), symbol);
                    if (epoch != e || s == 0) continue;
                    held += h;
                    shares += s;
                }
                assertGe(int256(booked), held);
                assertEq(total, shares);
            }
        }
    }

    function invariant_TheShortsOweAllTheVaultLent() public view {
        for (uint256 t; t < 2; ++t) {
            bytes32 symbol = handler.symbols(t);
            uint256 owed;
            for (uint256 i; i < 3; ++i) {
                (, uint256 debt,,) = shorts.position(handler.trader(i), symbol);
                owed += debt;
            }
            assertGe(owed + 3, vaults[t].debt());
            (uint256 total,,) = shorts.book(symbol, shorts.epoch(symbol));
            if (total == 0) assertEq(vaults[t].debt(), 0);
        }
    }

    function invariant_NoShortReachesTheSupplyVault() public view {
        assertEq(vault.totalAssets(), supplied);
        assertEq(vault.debt(), 0);
    }

    function invariant_NothingPanics() public view {
        assertFalse(handler.panicked());
    }
}
