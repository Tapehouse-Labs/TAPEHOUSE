// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
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

contract MarginAccountsHandler is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal immutable usdg;
    StubStockToken internal immutable nvda;
    SupplyVault internal immutable vault;
    MarginAccounts internal immutable accounts;
    BandDouble internal immutable band;
    StubToken internal immutable weth;
    address[3] public actors = [address(0xA1), address(0xA2), address(0xA3)];
    uint256 public claimed;
    uint256 public withdrawn;
    bool public sharePriceHeld = true;

    modifier rounded() {
        (uint256 debtBefore, uint256 sharesBefore) = (vault.debt(), accounts.totalDebtShares());
        _;
        (uint256 debtAfter, uint256 sharesAfter) = (vault.debt(), accounts.totalDebtShares());
        if ((debtAfter + 3) * (sharesBefore + 1e6) < (debtBefore + 1) * (sharesAfter + 1e6)) sharePriceHeld = false;
    }

    constructor(
        StubUsdg usdg_,
        StubStockToken nvda_,
        SupplyVault vault_,
        MarginAccounts accounts_,
        BandDouble band_,
        StubToken weth_
    ) {
        (usdg, nvda, vault, accounts, band, weth) = (usdg_, nvda_, vault_, accounts_, band_, weth_);
    }

    function deposit(uint256 actor, bool isolated, bool stock, uint256 amount) external {
        address account = actors[actor % 3];
        address token = stock ? address(nvda) : address(usdg);
        amount = bound(amount, 1, stock ? 100e18 : 10_000e6);
        if (stock) nvda.mint(account, amount);
        else usdg.mint(account, amount);
        vm.startPrank(account);
        StubToken(token).approve(address(accounts), amount);
        try accounts.deposit(isolated ? NVDA : bytes32(0), token, amount, account) {} catch {}
        vm.stopPrank();
    }

    function depositWeth(uint256 actor, bool isolated, uint256 amount) external {
        address account = actors[actor % 3];
        amount = bound(amount, 1, 10e18);
        weth.mint(address(this), amount);
        weth.approve(address(accounts), amount);
        accounts.deposit(isolated ? NVDA : bytes32(0), address(weth), amount, account);
    }

    function liquidate(uint256 actor, bool isolated) external rounded {
        address account = actors[actor % 3];
        bytes32 position = isolated ? NVDA : bytes32(0);
        address[3] memory tokens = [address(nvda), address(usdg), address(weth)];
        for (uint256 i; i < 3; ++i) {
            uint256 held = accounts.collateral(account, position, tokens[i]);
            if (held != 0) accounts.seize(position, tokens[i], held, account, address(this));
        }
        try accounts.writeOff(account, position) {} catch {}
    }

    function pause(bool paused) external {
        accounts.setBorrowingPaused(paused);
    }

    function withdraw(uint256 actor, bool isolated, bool stock, uint256 amount) external {
        address account = actors[actor % 3];
        bytes32 position = isolated ? NVDA : bytes32(0);
        address token = stock ? address(nvda) : address(usdg);
        uint256 held = accounts.collateral(account, position, token);
        if (held == 0) return;
        vm.prank(account);
        try accounts.withdraw(position, token, bound(amount, 1, held), account, account) {} catch {}
    }

    function borrow(uint256 actor, bool isolated, uint256 amount) external rounded {
        address account = actors[actor % 3];
        vm.prank(account);
        try accounts.borrow(isolated ? NVDA : bytes32(0), bound(amount, 1, 20_000e6), account, account) {} catch {}
    }

    function repay(uint256 actor, bool isolated, uint256 amount) external rounded {
        address account = actors[actor % 3];
        amount = bound(amount, 0, 30_000e6);
        usdg.mint(address(this), amount);
        usdg.approve(address(accounts), amount);
        accounts.repay(isolated ? NVDA : bytes32(0), amount, account);
    }

    function repayWithCollateral(uint256 actor, bool isolated, uint256 amount) external rounded {
        address account = actors[actor % 3];
        vm.prank(account);
        accounts.repayWithCollateral(isolated ? NVDA : bytes32(0), bound(amount, 0, 30_000e6), account);
    }

    function burn(uint256 fraction) external {
        uint256 balance = nvda.balanceOf(address(accounts));
        nvda.adminBurn(address(accounts), balance * bound(fraction, 0, 50) / 100);
        accounts.sync(NVDA);
    }

    function wait(uint256 seconds_) external {
        skip(bound(seconds_, 1, 30 days));
    }

    function session(uint256 kind, uint256 boundary) external {
        uint64 boundaryMs = uint64(block.timestamp + bound(boundary, 1, 3 days)) * 1000;
        kind %= 4;
        if (kind == 0) band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        else if (kind == 1) band.setSession(BandDouble.Session(2, 1, 3, 0, boundaryMs));
        else if (kind == 2) band.setSession(BandDouble.Session(1, 3, 1, 0, boundary % 2 == 0 ? boundaryMs : 0));
        else band.setSession(BandDouble.Session(0, 0, 0, 0, 0));
        accounts.accruePremium();
    }

    function claim() external {
        claimed += accounts.claimPremium();
    }

    function withdrawReserve(uint256 amount) external {
        amount = bound(amount, 0, accounts.reserve());
        accounts.withdrawReserve(address(this), amount);
        withdrawn += amount;
    }
}

contract MarginAccountsInvariantTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal usdg;
    StubStockToken internal nvda;
    SupplyVault internal vault;
    MarginAccounts internal accounts;
    MarginAccountsHandler internal handler;
    BandDouble internal band;

    function setUp() public {
        usdg = new StubUsdg();
        nvda = new StubStockToken(1e18);
        band = new BandDouble();
        band.setAsset(NVDA, BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(nvda)));
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 200e8, 50, 200e8, 202e8));
        bytes32[] memory symbols = new bytes32[](1);
        symbols[0] = NVDA;
        StubAggregator ethUsd = new StubAggregator(8, 2_000e8, block.timestamp, "ETH / USD");
        MarginDouble engine = new MarginDouble(symbols, address(band), address(ethUsd));
        vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 6_00, 40_00));
        StubToken weth = new StubToken(18);
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            weth,
            address(this),
            _uncapped(1),
            type(uint256).max,
            type(uint256).max,
            5_00,
            10_00
        );
        vault.setBorrower(address(accounts));
        usdg.mint(address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
        handler = new MarginAccountsHandler(usdg, nvda, vault, accounts, band, weth);
        accounts.setBackstop(address(handler));
        accounts.setLiquidator(address(handler));
        accounts.setGuardian(address(handler));
        accounts.transferOwnership(address(handler));
        vm.prank(address(handler));
        accounts.acceptOwnership();
        targetContract(address(handler));
    }

    function test_APositionOwesAtMostTheVaultsDebtAfterAFullRepaymentRoundsUp() public {
        handler.deposit(9649, false, true, 1e27);
        handler.borrow(141999345260031776733590614, false, 13917442368293244976735);
        handler.deposit(2, false, true, 245023170005346524);
        handler.deposit(1, false, false, 7893659450690472933937595126916925581757859686455475350098575369859880643);
        handler.deposit(
            type(uint256).max, false, true, 234194790134726769463474428395652224556668084584978870636928907936376688336
        );
        handler.borrow(4436, false, 707);
        handler.wait(31536000);
        handler.wait(112753810612168508597436786459684542341147490844132390787);
        handler.wait(2756702308);
        handler.wait(1e27);
        handler.wait(1e14);
        handler.wait(1e14);
        handler.repayWithCollateral(
            1300878292806012235712410258922363316353718979, false, 125099703957394107027084500998
        );
        handler.wait(1e27);
        handler.borrow(0, false, 100000);
        handler.borrow(31536000, false, 5880);
        handler.wait(86400);
        handler.repay(1, false, 1000311948889617140886002766);
        handler.repay(8906, false, 271);
        handler.borrow(type(uint256).max, false, 86460);
        handler.repay(
            6342232749332734634618387916968659883317419865749329917296396908504238, false, 5677438057177084407467517
        );
        address last = handler.actors(2);
        assertEq(accounts.totalDebtShares(), accounts.debtShares(last, bytes32(0)));
        assertGt(accounts.totalDebtShares(), vault.debt() * 1e6);
        assertEq(accounts.debt(last, bytes32(0)), vault.debt());
    }

    function invariant_TheSharesAddUp() public view {
        uint256 total;
        for (uint256 i; i < 3; ++i) {
            total += accounts.debtShares(handler.actors(i), bytes32(0)) + accounts.debtShares(handler.actors(i), NVDA);
        }
        assertEq(total, accounts.totalDebtShares());
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

    function invariant_TheSharePriceFallsOnlyByRounding() public view {
        assertTrue(handler.sharePriceHeld());
    }

    function invariant_TheAccountsHoldTheCollateralTheyCount() public view {
        uint256 stock;
        uint256 dollars;
        for (uint256 i; i < 3; ++i) {
            address a = handler.actors(i);
            stock += accounts.collateral(a, bytes32(0), address(nvda)) + accounts.collateral(a, NVDA, address(nvda));
            dollars += accounts.collateral(a, bytes32(0), address(usdg)) + accounts.collateral(a, NVDA, address(usdg));
        }
        (uint256 units, uint256 scale,) = accounts.holding(NVDA);
        assertLe(stock, units * scale / 1e18);
        assertGe(nvda.balanceOf(address(accounts)), units * scale / 1e18);
        assertEq(usdg.balanceOf(address(accounts)), dollars + accounts.reserve() + accounts.backstopPremium());
    }

    function invariant_TheReserveTakesATenthOfEachPremiumPaidAtMost() public view {
        assertLe((accounts.reserve() + handler.withdrawn()) * 9, accounts.backstopPremium() + handler.claimed());
    }

    function _uncapped(uint256 n) internal pure returns (uint256[] memory caps) {
        caps = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            caps[i] = type(uint256).max;
        }
    }
}
