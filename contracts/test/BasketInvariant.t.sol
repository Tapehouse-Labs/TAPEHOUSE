// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Test} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Basket} from "../src/Basket.sol";
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

contract BasketHandler is Test {
    bytes32 internal constant CROSS = bytes32(0);
    uint256 internal constant WAD = 1e18;
    StubUsdg internal immutable usdg;
    BandDouble internal immutable band;
    Basket internal immutable basket;
    MarginAccounts internal immutable accounts;
    Liquidator internal immutable liquidator;
    address internal immutable owner;
    address internal immutable keeper = makeAddr("keeper");
    address internal immutable buyer = makeAddr("buyer");
    address[3] internal actors = [makeAddr("alice"), makeAddr("bob"), makeAddr("carol")];
    bytes32[3] internal symbols = [bytes32("NVDA"), bytes32("TSLA"), bytes32("SPY")];
    StubStockToken[3] internal tokens;
    bool public panicked;
    bool public shareFell;
    bool public valueLost;
    bool public strayed;
    bool public pastCap;

    constructor(
        StubUsdg usdg_,
        BandDouble band_,
        Basket basket_,
        MarginAccounts accounts_,
        Liquidator liquidator_,
        StubStockToken[3] memory tokens_,
        address owner_
    ) {
        (usdg, band, basket, accounts, liquidator, tokens, owner) =
        (usdg_, band_, basket_, accounts_, liquidator_, tokens_, owner_);
        for (uint256 i; i < 3; ++i) {
            vm.startPrank(actors[i]);
            basket.approve(address(accounts), type(uint256).max);
            usdg.approve(address(accounts), type(uint256).max);
            for (uint256 t; t < 3; ++t) {
                tokens[t].approve(address(basket), type(uint256).max);
            }
            vm.stopPrank();
        }
        vm.startPrank(keeper);
        for (uint256 t; t < 3; ++t) {
            tokens[t].approve(address(basket), type(uint256).max);
        }
        vm.stopPrank();
        usdg.mint(buyer, 1e15);
        vm.prank(buyer);
        usdg.approve(address(liquidator), type(uint256).max);
    }

    function actor(uint256 i) external view returns (address) {
        return actors[i];
    }

    function mint(uint256 who, uint256 shares) external {
        address account = actors[who % 3];
        shares = bound(shares, 1, 20e18);
        uint256[] memory assets = basket.previewMint(shares);
        for (uint256 t; t < 3; ++t) {
            tokens[t].mint(account, assets[t]);
        }
        (uint256[] memory held, uint256 supply) = (basket.totalAssets(), basket.totalSupply());
        vm.prank(account);
        try basket.mint(shares, account, assets) {
            _checkShare(held, supply);
        } catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function redeem(uint256 who, uint256 shares) external {
        address account = actors[who % 3];
        uint256 balance = basket.balanceOf(account);
        if (balance == 0) return;
        (uint256[] memory held, uint256 supply) = (basket.totalAssets(), basket.totalSupply());
        vm.prank(account);
        try basket.redeem(bound(shares, 1, balance), account, account) {
            _checkShare(held, supply);
        } catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function transfer(uint256 from, uint256 to, uint256 shares) external {
        address account = actors[from % 3];
        uint256 balance = basket.balanceOf(account);
        if (balance == 0) return;
        vm.prank(account);
        try basket.transfer(actors[to % 3], bound(shares, 1, balance)) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function deposit(uint256 who, uint256 shares) external {
        address account = actors[who % 3];
        uint256 balance = basket.balanceOf(account);
        if (balance == 0) return;
        vm.prank(account);
        try accounts.deposit(CROSS, address(basket), bound(shares, 1, balance), account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function withdraw(uint256 who, uint256 shares) external {
        address account = actors[who % 3];
        uint256 held = accounts.collateral(account, CROSS, address(basket));
        if (held == 0) return;
        vm.prank(account);
        try accounts.withdraw(CROSS, address(basket), bound(shares, 1, held), account, account) {
            if (accounts.debt(account, CROSS) != 0) _checkCaps(account);
        } catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function borrow(uint256 who, uint256 part) external {
        address account = actors[who % 3];
        (int256 equity, uint256 requirement,,) = accounts.health(account, CROSS);
        if (equity <= int256(requirement) + 1e12) return;
        uint256 room = (uint256(equity) - requirement) / 1e12;
        vm.prank(account);
        try accounts.borrow(CROSS, room * bound(part, 50, 100) / 100, account, account) {
            _checkCaps(account);
        } catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function repay(uint256 who, uint256 part) external {
        address account = actors[who % 3];
        uint256 owed = (accounts.debt(account, CROSS) + accounts.premium(account, CROSS)) * bound(part, 1, 100) / 100;
        if (owed == 0) return;
        usdg.mint(account, owed);
        vm.prank(account);
        try accounts.repay(CROSS, owed, account) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function unwrap(uint256 who, uint256 shares) external {
        address account = actors[who % 3];
        uint256 held = accounts.collateral(account, CROSS, address(basket));
        if (held == 0) return;
        vm.prank(account);
        try accounts.unwrap(account, address(basket), bound(shares, 1, held)) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function move(uint256 which, uint256 percent) external {
        bytes32 symbol = symbols[which % 3];
        (,, uint64 mid,,,) = band.quote(symbol);
        mid = uint64(bound(uint256(mid) * bound(percent, 70, 140) / 100, 20e8, 2_000e8));
        band.setQuote(symbol, BandDouble.Quote(3, 3, mid, 1_00, mid - mid / 100, mid + mid / 100));
    }

    function liquidate(uint256 who, uint256 which, uint256 amount) external {
        address account = actors[who % 3];
        try liquidator.start(account, CROSS) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
        address token = address(tokens[which % 3]);
        uint256 held = accounts.collateral(account, CROSS, token);
        for (uint256 t; held == 0 && t < 3; ++t) {
            token = address(tokens[t]);
            held = accounts.collateral(account, CROSS, token);
        }
        if (held == 0) return;
        vm.prank(buyer);
        try liquidator.buy(account, CROSS, token, bound(amount, 1, held), type(uint256).max, buyer) {}
        catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function burn(uint256 which, uint256 part) external {
        StubStockToken token = tokens[which % 3];
        uint256 held = token.balanceOf(address(basket));
        token.adminBurn(address(basket), held * bound(part, 0, 5) / 100);
    }

    function gift(uint256 which, uint256 amount) external {
        tokens[which % 3].mint(address(basket), bound(amount, 0, 1e18));
    }

    function propose(uint256 a, uint256 b, uint256 c) external {
        uint256[] memory units = new uint256[](3);
        (units[0], units[1], units[2]) = (bound(a, 0, 2e18), bound(b, 0, 1e18), bound(c, 1, 5e17));
        vm.prank(owner);
        basket.proposeTarget(units);
        vm.warp(block.timestamp + 7 days);
    }

    function rebalance(uint256 seed) external {
        uint256 supply = basket.totalSupply();
        if (supply == 0) return;
        uint256[] memory units = basket.target();
        uint256[] memory held = basket.totalAssets();
        (uint256[] memory lows, uint256[] memory highs) = _edges();
        (uint256[] memory assetsIn, uint256[] memory assetsOut) = _trade(seed, units, held, supply, lows, highs);
        uint256 beforeLow = _worth(held, lows);
        uint256 beforeHigh = _worth(held, highs);
        vm.prank(keeper);
        try basket.rebalance(assetsIn, assetsOut, keeper) {
            uint256[] memory afterwards = basket.totalAssets();
            if (_worth(afterwards, lows) < beforeLow || _worth(afterwards, highs) < beforeHigh) valueLost = true;
            for (uint256 t; t < 3; ++t) {
                if (_distance(afterwards[t], units[t], supply) > _distance(held[t], units[t], supply)) strayed = true;
            }
        } catch Panic(uint256) {
            panicked = true;
        } catch {}
    }

    function _edges() internal view returns (uint256[] memory lows, uint256[] memory highs) {
        (lows, highs) = (new uint256[](3), new uint256[](3));
        for (uint256 t; t < 3; ++t) {
            (,,,, uint64 low, uint128 high) = band.quote(symbols[t]);
            (lows[t], highs[t]) = (low, high);
        }
    }

    function _trade(
        uint256 seed,
        uint256[] memory units,
        uint256[] memory held,
        uint256 supply,
        uint256[] memory lows,
        uint256[] memory highs
    ) internal returns (uint256[] memory assetsIn, uint256[] memory assetsOut) {
        (assetsIn, assetsOut) = (new uint256[](3), new uint256[](3));
        uint256 worthIn;
        uint256 worthOut;
        for (uint256 t; t < 3; ++t) {
            uint256 goal = units[t] * supply / WAD;
            uint256 part = uint256(keccak256(abi.encode(seed, t))) % 101;
            if (goal > held[t]) {
                assetsIn[t] = (goal - held[t]) * part / 100;
                tokens[t].mint(keeper, assetsIn[t]);
                worthIn += assetsIn[t] * lows[t];
            } else if (held[t] * WAD > units[t] * supply) {
                assetsOut[t] = (held[t] - _ceilGoal(units[t], supply)) * part / 100;
                worthOut += assetsOut[t] * highs[t];
            }
        }
        if (worthOut > worthIn) {
            for (uint256 t; t < 3; ++t) {
                assetsOut[t] = assetsOut[t] * worthIn / worthOut;
            }
        }
    }

    function wait(uint256 seconds_) external {
        vm.warp(block.timestamp + bound(seconds_, 1, 3 days));
    }

    function _checkCaps(address account) internal {
        if (accounts.collateral(account, CROSS, address(basket)) == 0) return;
        uint256 credited;
        for (uint256 i; i < 3; ++i) {
            credited += accounts.collateral(actors[i], CROSS, address(basket));
        }
        uint256[] memory inBaskets = basket.previewRedeem(credited);
        for (uint256 t; t < 3; ++t) {
            (uint256 units, uint256 scale, uint256 cap) = accounts.holding(symbols[t]);
            if (units * scale / WAD + inBaskets[t] > cap) pastCap = true;
        }
    }

    function _checkShare(uint256[] memory held, uint256 supply) internal {
        if (supply == 0) return;
        uint256[] memory afterwards = basket.totalAssets();
        uint256 supplyAfter = basket.totalSupply();
        for (uint256 t; t < 3; ++t) {
            if (afterwards[t] * supply < held[t] * supplyAfter) shareFell = true;
        }
    }

    function _worth(uint256[] memory assets, uint256[] memory prices) internal pure returns (uint256 worth) {
        for (uint256 t; t < assets.length; ++t) {
            worth += assets[t] * prices[t];
        }
    }

    function _ceilGoal(uint256 units, uint256 supply) internal pure returns (uint256) {
        return (units * supply + WAD - 1) / WAD;
    }

    function _distance(uint256 held, uint256 units, uint256 supply) internal pure returns (uint256) {
        uint256 have = held * WAD;
        uint256 want = units * supply;
        return have > want ? have - want : want - have;
    }
}

abstract contract BasketInvariantBase is Test {
    bytes32 internal constant CROSS = bytes32(0);
    Basket internal basket;
    MarginAccounts internal accounts;
    BasketHandler internal handler;
    StubStockToken[3] internal tokens;

    function setUp() public {
        vm.warp(1_790_000_000);
        StubUsdg usdg = new StubUsdg();
        tokens = [new StubStockToken(1e18), new StubStockToken(1e18), new StubStockToken(1e18)];
        BandDouble band = new BandDouble();
        bytes32[] memory symbols = new bytes32[](3);
        (symbols[0], symbols[1], symbols[2]) = ("NVDA", "TSLA", "SPY");
        uint64[3] memory mids = [uint64(200e8), 300e8, 600e8];
        for (uint256 t; t < 3; ++t) {
            band.setAsset(symbols[t], BandDouble.Asset(address(0), bytes32(0), bytes32(0), address(tokens[t])));
            band.setQuote(
                symbols[t], BandDouble.Quote(3, 3, mids[t], 1_00, mids[t] - mids[t] / 100, mids[t] + mids[t] / 100)
            );
        }
        band.setSession(BandDouble.Session(2, 1, 1, 0, 0));
        MarginDouble engine =
            new MarginDouble(symbols, address(band), address(new StubAggregator(8, 2_000e8, block.timestamp, "ETH")));
        SupplyVault vault = new SupplyVault(IUSDG(address(usdg)), address(this), SupplyVault.RateModel(90_00, 0, 0, 0));
        accounts = new MarginAccounts(
            IBand(address(band)),
            IMargin(address(engine)),
            vault,
            new StubToken(18),
            address(this),
            _caps(),
            type(uint256).max,
            type(uint256).max,
            0,
            10_00
        );
        vault.setBorrower(address(accounts));
        Liquidator liquidator = new Liquidator(accounts);
        accounts.setLiquidator(address(liquidator));
        uint256[] memory units = new uint256[](3);
        (units[0], units[1], units[2]) = (1e18, 5e17, 2e17);
        basket = new Basket("Tapehouse Test Basket", "thTEST", IBand(address(band)), symbols, units, address(this));
        accounts.addBasket(basket);
        usdg.mint(address(this), 1e12);
        usdg.approve(address(vault), 1e12);
        vault.deposit(1e12, address(this));
        handler = new BasketHandler(usdg, band, basket, accounts, liquidator, tokens, address(this));
        targetContract(address(handler));
    }

    function _caps() internal pure virtual returns (uint256[] memory);

    function invariant_NothingPanics() public view {
        assertFalse(handler.panicked());
    }

    function invariant_MintsAndRedeemsNeverTakeFromTheOtherHolders() public view {
        assertFalse(handler.shareFell());
    }

    function invariant_ARebalanceNeverLowersTheBasketsWorthAtTheBandsEdges() public view {
        assertFalse(handler.valueLost());
    }

    function invariant_ARebalanceNeverMovesATokenAwayFromItsTarget() public view {
        assertFalse(handler.strayed());
    }

    function invariant_NoNewRiskWhileTheAccountsHoldAStockTokenPastItsCap() public view {
        assertFalse(handler.pastCap());
    }

    function invariant_SharesNeverClaimMoreThanTheBasketHolds() public view {
        uint256[] memory held = basket.totalAssets();
        uint256[] memory claimed = basket.previewRedeem(basket.balanceOf(address(accounts)));
        for (uint256 i; i < 3; ++i) {
            uint256[] memory own = basket.previewRedeem(basket.balanceOf(handler.actor(i)));
            for (uint256 t; t < 3; ++t) {
                claimed[t] += own[t];
            }
        }
        for (uint256 t; t < 3; ++t) {
            assertLe(claimed[t], held[t]);
        }
    }

    function invariant_TheAccountsHoldTheSharesAndStockTokensTheirPositionsCount() public view {
        uint256 shares;
        uint256[3] memory counted;
        for (uint256 i; i < 3; ++i) {
            shares += accounts.collateral(handler.actor(i), CROSS, address(basket));
            for (uint256 t; t < 3; ++t) {
                counted[t] += accounts.collateral(handler.actor(i), CROSS, address(tokens[t]));
            }
        }
        assertEq(basket.balanceOf(address(accounts)), shares);
        for (uint256 t; t < 3; ++t) {
            assertLe(counted[t], tokens[t].balanceOf(address(accounts)));
        }
    }
}

contract BasketInvariantTest is BasketInvariantBase {
    function _caps() internal pure override returns (uint256[] memory caps) {
        caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (type(uint256).max, type(uint256).max, type(uint256).max);
    }
}

contract CappedBasketInvariantTest is BasketInvariantBase {
    function _caps() internal pure override returns (uint256[] memory caps) {
        caps = new uint256[](3);
        (caps[0], caps[1], caps[2]) = (40e18, 20e18, 8e18);
    }
}
