// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {BandFeed} from "../src/BandFeed.sol";
import {Liquidator} from "../src/Liquidator.sol";
import {MarginAccounts} from "../src/MarginAccounts.sol";
import {ReopeningAuction} from "../src/ReopeningAuction.sol";
import {SupplyVault} from "../src/SupplyVault.sol";
import {IBand} from "../src/interfaces/IBand.sol";
import {IMargin} from "../src/interfaces/IMargin.sol";
import {IUSDG} from "../src/interfaces/IUSDG.sol";
import {BandDouble} from "./doubles/BandDouble.sol";
import {MarginDouble} from "./doubles/MarginDouble.sol";
import {SealDouble} from "./doubles/SealDouble.sol";
import {StubAggregator} from "./devnode/StubAggregator.sol";
import {StubStockToken} from "./devnode/StubStockToken.sol";
import {StubToken} from "./devnode/StubToken.sol";
import {StubUsdg} from "./devnode/StubUsdg.sol";

contract ReopeningAuctionHandler is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal immutable usdg;
    StubStockToken internal immutable nvda;
    BandDouble internal immutable band;
    MarginAccounts internal immutable accounts;
    ReopeningAuction internal immutable auction;
    bool public paidAtMostTheirPrice = true;
    uint64[] public opens;
    uint256 internal nonce;

    constructor(
        StubUsdg usdg_,
        StubStockToken nvda_,
        BandDouble band_,
        MarginAccounts accounts_,
        ReopeningAuction auction_
    ) {
        (usdg, nvda, band, accounts, auction) = (usdg_, nvda_, band_, accounts_, auction_);
    }

    function openCount() external view returns (uint256) {
        return opens.length;
    }

    function round(uint256 seed, uint256 lots, uint256 bids, uint256 claims) external {
        uint64 openMs = _weekend();
        opens.push(openMs);
        lots = bound(lots, 1, 3);
        for (uint256 i; i < lots; ++i) {
            address account = address(uint160(uint256(keccak256(abi.encode(openMs, i)))));
            _position(account, bound(uint256(keccak256(abi.encode(seed, i, "d"))), 1_530e6, 1_600e6));
            try auction.enroll(account, bytes32(0), NVDA) {} catch {}
        }
        bids = bound(bids, 1, 5);
        uint256[] memory quantities = new uint256[](bids);
        uint256[] memory prices = new uint256[](bids);
        bytes32[] memory hashes = new bytes32[](bids);
        for (uint256 i; i < bids; ++i) {
            quantities[i] = bound(uint256(keccak256(abi.encode(seed, i, "q"))), 1e18, 12e18);
            prices[i] = 171e8 + bound(uint256(keccak256(abi.encode(seed, i, "p"))), 0, 4) * 4e8;
            address bidder = address(uint160(0x3000 + i));
            uint256 deposit = (quantities[i] * prices[i] - 1) / 1e20 + 1 + seed % 500e6;
            usdg.mint(bidder, deposit);
            hashes[i] = keccak256(abi.encode(bidder, NVDA, openMs, quantities[i], prices[i], bytes32(i)));
            vm.startPrank(bidder);
            usdg.approve(address(auction), type(uint256).max);
            auction.commit(NVDA, hashes[i], deposit);
            vm.stopPrank();
        }
        vm.warp(openMs / 1000 - 30 minutes);
        for (uint256 i; i < bids; ++i) {
            if (uint256(keccak256(abi.encode(seed, i, "r"))) % 5 == 0) continue;
            vm.prank(address(uint160(0x3000 + i)));
            try auction.reveal(NVDA, openMs, quantities[i], prices[i], bytes32(i)) {
                hashes[i] = bytes32(0);
            } catch {}
        }
        vm.warp(openMs / 1000);
        band.setSession(BandDouble.Session(2, 1, 2, openMs + 6.5 hours * 1000, 0));
        if (seed % 4 != 0) {
            for (uint256 step = 4;; --step) {
                try auction.clear(NVDA, openMs, 171e8 + step * 4e8) {
                    break;
                } catch {}
                if (step == 0) break;
            }
        } else {
            skip(1 hours);
        }
        for (uint256 i; i < bids; ++i) {
            if (hashes[i] != bytes32(0)) auction.forfeit(NVDA, openMs, hashes[i]);
        }
        ReopeningAuction.Bid[] memory revealed = auction.bids(NVDA, openMs);
        claims = bound(claims, 0, revealed.length);
        for (uint256 i; i < claims; ++i) {
            _claim(openMs, i, revealed[i]);
        }
        nonce++;
    }

    function _claim(uint64 openMs, uint256 index, ReopeningAuction.Bid memory bid) internal {
        uint256 cash = usdg.balanceOf(bid.bidder);
        uint256 stock = nvda.balanceOf(bid.bidder);
        try auction.claim(NVDA, openMs, index) {
            uint256 got = nvda.balanceOf(bid.bidder) - stock;
            uint256 paid = bid.escrow - (usdg.balanceOf(bid.bidder) - cash);
            uint256 price = auction.round(NVDA, openMs).price;
            if (paid > got * price / 1e20 + 4 || got > bid.quantity) paidAtMostTheirPrice = false;
        } catch {}
    }

    function _weekend() internal returns (uint64 openMs) {
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        uint64 closes = uint64(vm.getBlockTimestamp() + 1 hours) * 1000;
        band.setSession(BandDouble.Session(2, 1, 3, 0, closes));
        accounts.accruePremium();
        skip(1 hours);
        uint64 reopenMs = closes + 48 hours * 1000;
        band.setSession(BandDouble.Session(1, 3, 1, 0, reopenMs));
        accounts.accruePremium();
        skip(48 hours);
        openMs = reopenMs + 13.5 hours * 1000;
        band.setSession(BandDouble.Session(2, 3, 1, openMs, 0));
        accounts.accruePremium();
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 191e8, 50, 190e8, 192e8));
    }

    function _position(address account, uint256 debt) internal {
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 201e8, 50, 200e8, 202e8));
        nvda.mint(account, 10e18);
        vm.startPrank(account);
        nvda.approve(address(accounts), 10e18);
        accounts.deposit(bytes32(0), address(nvda), 10e18, account);
        try accounts.borrow(bytes32(0), debt, account, account) {} catch {}
        vm.stopPrank();
        band.setQuote(NVDA, BandDouble.Quote(3, 3, 191e8, 50, 190e8, 192e8));
    }
}

/// forge-config: default.invariant.depth = 32
contract ReopeningAuctionInvariantTest is Test {
    bytes32 internal constant NVDA = "NVDA";
    StubUsdg internal usdg;
    StubStockToken internal nvda;
    ReopeningAuction internal auction;
    ReopeningAuctionHandler internal handler;

    function setUp() public {
        vm.warp(1_790_000_000);
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
            0,
            10_00
        );
        Liquidator liquidator = new Liquidator(accounts);
        BandFeed[] memory feeds = new BandFeed[](0);
        auction = new ReopeningAuction(liquidator, new bytes32[](0), feeds);
        vault.setBorrower(address(accounts));
        accounts.setLiquidator(address(liquidator));
        liquidator.setAuction(address(auction));
        usdg.mint(address(this), 1_000_000e6);
        usdg.approve(address(vault), 1_000_000e6);
        vault.deposit(1_000_000e6, address(this));
        handler = new ReopeningAuctionHandler(usdg, nvda, band, accounts, auction);
        targetContract(address(handler));
    }

    function invariant_TheAuctionHoldsTheUsdgItOwes() public view {
        uint256 owed;
        for (uint256 i; i < handler.openCount(); ++i) {
            ReopeningAuction.Round memory r = auction.round(NVDA, handler.opens(i));
            owed += r.pool + r.deposits;
        }
        assertGe(usdg.balanceOf(address(auction)), owed);
    }

    function invariant_TheAuctionHoldsTheStockItOwes() public view {
        uint256 owed;
        for (uint256 i; i < handler.openCount(); ++i) {
            uint64 openMs = handler.opens(i);
            ReopeningAuction.Round memory r = auction.round(NVDA, openMs);
            ReopeningAuction.Bid[] memory bids = auction.bids(NVDA, openMs);
            for (uint256 j; j < bids.length; ++j) {
                if (bids[j].claimed || !r.cleared) continue;
                uint256 allocated = r.supply < uint256(r.above) + r.atPrice ? r.supply : uint256(r.above) + r.atPrice;
                if (bids[j].price < r.price || allocated == 0) continue;
                owed += bids[j].price > r.price
                    ? uint256(bids[j].quantity) * r.sold / allocated
                    : uint256(bids[j].quantity) * (allocated - r.above) * r.sold / (uint256(r.atPrice) * allocated);
            }
        }
        assertGe(nvda.balanceOf(address(auction)), owed);
    }

    function invariant_NoBidderPaysAboveItsPrice() public view {
        assertTrue(handler.paidAtMostTheirPrice());
    }

    function invariant_ARoundNeverSellsMoreThanItsSupply() public view {
        for (uint256 i; i < handler.openCount(); ++i) {
            ReopeningAuction.Round memory r = auction.round(NVDA, handler.opens(i));
            assertLe(r.sold, r.supply);
        }
    }
}
