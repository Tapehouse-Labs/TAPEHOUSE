// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {GapBackstop} from "../src/GapBackstop.sol";
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

contract GapBackstopHandler is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal immutable usdg;
    StubStockToken internal immutable nvda;
    BandDouble internal immutable band;
    MarginAccounts internal immutable accounts;
    Liquidator internal immutable liquidator;
    GapBackstop internal immutable backstop;
    address[3] public actors = [address(0xA1), address(0xA2), address(0xA3)];
    bool public redeemedOnlyInWindow = true;
    bool public lendersLostOnlyPastTheLimit = true;
    bytes32[] public keys;
    uint64[] public closures;
    uint64 internal lastClose;

    constructor(
        StubUsdg usdg_,
        StubStockToken nvda_,
        BandDouble band_,
        MarginAccounts accounts_,
        Liquidator liquidator_,
        GapBackstop backstop_
    ) {
        (usdg, nvda, band, accounts, liquidator, backstop) = (usdg_, nvda_, band_, accounts_, liquidator_, backstop_);
        usdg.approve(address(liquidator_), type(uint256).max);
        keys.push(bytes32(0));
        keys.push(NVDA);
        closures.push(0);
    }

    function keyCount() external view returns (uint256) {
        return keys.length;
    }

    function closureCount() external view returns (uint256) {
        return closures.length;
    }

    function deposit(uint256 actor, uint256 assets) external {
        address depositor = actors[actor % 3];
        assets = bound(assets, 1, 1_000e6);
        usdg.mint(depositor, assets);
        vm.startPrank(depositor);
        usdg.approve(address(backstop), assets);
        backstop.deposit(assets, depositor);
        vm.stopPrank();
    }

    function startCooldown(uint256 actor) external {
        vm.prank(actors[actor % 3]);
        backstop.startCooldown();
    }

    function redeem(uint256 actor, uint256 shares) external {
        address depositor = actors[actor % 3];
        (uint192 cooling, uint64 startedAt) = backstop.cooldowns(depositor);
        shares = bound(shares, 1, backstop.balanceOf(depositor) + 1);
        vm.prank(depositor);
        try backstop.redeem(shares, depositor, depositor) {
            uint256 nowMs = vm.getBlockTimestamp() * 1000;
            bool inWindow = startedAt != 0 && vm.getBlockTimestamp() >= startedAt + 7 days
                && vm.getBlockTimestamp() <= startedAt + 13 days && shares <= cooling;
            (uint8 state, uint8 nyse, uint8 next,, uint64 boundaryMs) = band.session();
            bool open = state == 2 && (boundaryMs == 0 || boundaryMs > nowMs) && (nyse == 1 || (nyse == 2 && next == 1));
            (uint64 closesMs, uint64 reopensMs,) = accounts.closure();
            uint256 until = reopensMs > closesMs ? reopensMs + 1 days * 1000 : closesMs + 96 hours * 1000;
            bool settled = closesMs == 0 || closesMs > nowMs || nowMs >= until;
            if (!inWindow || !open || !settled) redeemedOnlyInWindow = false;
        } catch {}
    }

    function transfer(uint256 from, uint256 to, uint256 shares) external {
        address sender = actors[from % 3];
        shares = bound(shares, 0, backstop.balanceOf(sender));
        vm.prank(sender);
        IERC20(address(backstop)).transfer(actors[to % 3], shares);
    }

    function wait(uint256 time) external {
        skip(bound(time, 1 hours, 3 days));
    }

    function session(uint256 kind) external {
        kind %= 4;
        if (kind == 0) band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        else if (kind == 1) band.setSession(BandDouble.Session(2, 3, 1, 0, 0));
        else if (kind == 2) band.setSession(BandDouble.Session(1, 3, 1, 0, 0));
        else band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
    }

    function closure() external {
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 48 hours * 1000));
        accounts.accruePremium();
        skip(48 hours);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
    }

    function longClosure(uint256 crash) external {
        for (uint256 i; i < 3; ++i) {
            _empty(actors[i], bytes32(0), crash);
        }
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        band.setSession(BandDouble.Session(1, 3, 1, 0, closes + 96 hours * 1000));
        for (uint256 i; i < 3; ++i) {
            accounts.accruePremium();
            _cover(actors[i], bytes32(0));
            skip(30 hours);
        }
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        accounts.accruePremium();
    }

    function shortfall(uint256 actor, bool isolated, uint256 crash) external {
        address account = actors[actor % 3];
        bytes32 position = isolated ? NVDA : bytes32(0);
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        _empty(account, position, crash);
        _cover(account, position);
    }

    function _empty(address account, bytes32 position, uint256 crash) internal {
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
        nvda.mint(account, 10e18);
        vm.startPrank(account);
        nvda.approve(address(accounts), 10e18);
        try accounts.deposit(position, address(nvda), 10e18, account) {}
        catch {
            vm.stopPrank();
            return;
        }
        try accounts.borrow(position, 1_600e6, account, account) {} catch {}
        vm.stopPrank();
        uint64 low = uint64(bound(crash, 10e8, 150e8));
        band.setQuote(NVDA, BandDouble.Quote(3, 3, low, 50, low, low + low / 100));
        try liquidator.start(account, position) {} catch {}
        usdg.mint(address(this), 10_000e6);
        for (uint256 i; i < 4 && accounts.collateral(account, position, address(nvda)) != 0; ++i) {
            try liquidator.buy(account, position, address(nvda), 10e18, type(uint256).max, address(this)) {}
            catch {
                break;
            }
        }
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
    }

    function _cover(address account, bytes32 position) internal {
        try backstop.cover(account, position) returns (uint256, uint256 written) {
            if (written != 0 && backstop.exposureLeft(position) != 0) lendersLostOnlyPastTheLimit = false;
            uint64 key = backstop.closureMs();
            if (key != lastClose) {
                closures.push(key);
                lastClose = key;
            }
        } catch {}
    }
}

contract GapBackstopInvariantTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal usdg;
    GapBackstop internal backstop;
    GapBackstopHandler internal handler;

    function setUp() public {
        vm.warp(1_790_000_000);
        usdg = new StubUsdg();
        StubStockToken nvda = new StubStockToken(1e18);
        BandDouble band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = NVDA;
        StubAggregator ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        MarginDouble engine = new MarginDouble(symbols, address(band), address(ethUsd));
        SupplyVault vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 0, 0));
        uint256[] memory caps = new uint256[](1);
        caps[0] = type(uint256).max;
        MarginAccounts accounts = new MarginAccounts(
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
        Liquidator liquidator = new Liquidator(accounts);
        uint256[] memory limits = new uint256[](1);
        limits[0] = 700e6;
        backstop = new GapBackstop(accounts, address(this), 1_000e6, limits);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        accounts.setBackstop(address(backstop));
        usdg.mint(address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
        handler = new GapBackstopHandler(usdg, nvda, band, accounts, liquidator, backstop);
        targetContract(address(handler));
    }

    function invariant_TheBackstopHoldsTheUsdgItCounts() public view {
        assertEq(usdg.balanceOf(address(backstop)), backstop.held());
        assertEq(backstop.totalAssets(), backstop.held());
    }

    function invariant_NoClosureCoversMoreThanItsLimit() public view {
        for (uint256 k; k < handler.keyCount(); ++k) {
            bytes32 key = handler.keys(k);
            for (uint256 c; c < handler.closureCount(); ++c) {
                assertLe(backstop.covered(key, handler.closures(c)), backstop.exposureLimit(key));
            }
        }
    }

    function invariant_LendersLoseOnlyWhatTheBackstopCannotCover() public view {
        assertTrue(handler.lendersLostOnlyPastTheLimit());
    }

    function invariant_SharesLeaveOnlyInTheirWindowWhileTheMarketIsOpen() public view {
        assertTrue(handler.redeemedOnlyInWindow());
    }

    function invariant_NoCooldownExceedsItsOwnersShares() public view {
        for (uint256 i; i < 3; ++i) {
            address actor = handler.actors(i);
            (uint192 shares,) = backstop.cooldowns(actor);
            assertLe(shares, backstop.balanceOf(actor));
        }
    }
}
