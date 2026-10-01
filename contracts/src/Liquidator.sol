// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IMargin} from "./interfaces/IMargin.sol";
import {MarginAccounts} from "./MarginAccounts.sol";

/// @title Tapehouse liquidator
/// @notice Sells the collateral of a margin position that falls short, in a Dutch auction that follows the market.
/// The market is open while the 24/5 session is open and NYSE is in regular hours or between two trading days. Then
/// a position falls short when its equity at each band's low edge is below the engine's current requirement, and its
/// auction starts at each band's high edge and falls a basis point a second. Otherwise, from the regular close before
/// a weekend or holiday to the regular open after it, or while the session is not known, a position falls short only
/// if it does so at both edges of its bands, so a shortfall the band can explain waits for the reopening; its auction
/// falls a basis point every four seconds and sells at most a tenth of the position's holding of a Stock Token an
/// hour. An unknown session is judged against the open requirement with its buffer, never the closed one it steps up
/// to. WETH counts for 84% of its Chainlink price. Nothing is sold while a held asset is halted or its multiplier
/// change is not yet confirmed, while WETH is held and its price is missing or stale, or while the L2 sequencer is not
/// settled. A purchase repays at most half of what the
/// position owes, or all of it once that is 2,000 USDG or less. A buyer pays USDG: half a percent goes to the
/// accounts' reserve unless the position is worth less than it owes, and the rest repays the position, its debt first.
/// Anyone may start, stop, buy and settle a position's USDG against its debt; a position left with nothing is
/// written off by the accounts' backstop, or by anyone while they have none.
/// @dev Set as the accounts' liquidator, it is their only way into seizures and write-offs.
contract Liquidator {
    using SafeERC20 for IERC20;

    struct Auction {
        uint64 startedAt;
        bool closed;
    }

    struct Allowance {
        uint192 used;
        uint64 updatedAt;
    }

    struct Valuation {
        bytes32[] symbols;
        address[] tokens;
        int256[] quantities;
        uint256[] prices;
        int256 equity;
        uint256 requirement;
        uint256 owed;
        bool held;
        bool judged;
        bool insolvent;
    }

    /// @notice What WETH counts for when a position is judged, in basis points of its Chainlink value: Aave v3's
    /// liquidation threshold for WETH on Arbitrum.
    uint256 public constant WETH_LIQUIDATION_BPS = 84_00;
    /// @notice How fast an auction's price falls while the market is open, in basis points a minute.
    uint256 public constant OPEN_DECAY_BPS_PER_MINUTE = 60;
    /// @notice How fast an auction's price falls while the market is closed, in basis points a minute.
    uint256 public constant CLOSED_DECAY_BPS_PER_MINUTE = 15;
    /// @notice The lowest an auction's price falls, in basis points below the band's low edge or WETH's price.
    uint256 public constant MAX_DISCOUNT_BPS = 10_00;
    /// @notice How long an auction started while the market is open runs before it must start again, in seconds.
    uint256 public constant OPEN_AUCTION_LIFETIME = 1 hours;
    /// @notice How long an auction started while the market is closed runs before it must start again, in seconds.
    uint256 public constant CLOSED_AUCTION_LIFETIME = 4 hours;
    /// @notice The share of every purchase that goes to the accounts' reserve, in basis points.
    uint256 public constant FEE_BPS = 50;
    /// @notice The most of what a position owes one purchase repays, in basis points.
    uint256 public constant CLOSE_FACTOR_BPS = 50_00;
    /// @notice What a position owes, in USDG, at or below which one purchase may repay all of it.
    uint256 public constant FULL_CLOSE_OWED = 2_000e6;
    /// @notice The most of a position's holding of a Stock Token sold in an hour while the market is closed, in basis
    /// points.
    uint256 public constant CLOSED_HOURLY_BPS = 10_00;
    /// @notice The horizon the open requirement is measured over, in seconds: two days, as the engine's.
    uint64 public constant HORIZON = 172_800;
    /// @notice The oldest ETH/USD answer WETH is priced at, in seconds: its heartbeat plus a minute.
    uint256 public constant MAX_FEED_AGE = 86_400 + 60;

    uint256 private constant BPS = 10_000;
    uint256 private constant PRICE_UNIT = 1e8;
    uint256 private constant USDG_TO_USD = 1e12;
    uint256 private constant TOKEN_TO_USDG = 1e20;

    /// @notice The accounts whose positions this contract liquidates.
    MarginAccounts public immutable accounts;
    /// @notice The band that prices every Stock Token.
    IBand public immutable band;
    /// @notice The engine that margins every position.
    IMargin public immutable engine;
    /// @notice The token buyers pay in.
    IERC20 public immutable usdg;
    /// @notice Wrapped ETH.
    IERC20 public immutable weth;
    /// @notice The Chainlink feed that prices WETH; zero where there is none.
    AggregatorV3Interface public immutable ethUsd;

    /// @notice The auction of each position, if one runs.
    mapping(address account => mapping(bytes32 position => Auction)) public auctions;
    mapping(address account => mapping(bytes32 position => mapping(address token => Allowance))) private _allowances;

    /// @notice An auction of `account`'s `position` started, `closed` if the market was closed.
    event AuctionStarted(address indexed account, bytes32 indexed position, bool closed);
    /// @notice The auction of `account`'s `position` stopped: the position no longer falls short.
    event AuctionStopped(address indexed account, bytes32 indexed position);
    /// @notice `buyer` paid `cost` of USDG for `amount` of `token` from `account`'s `position`, `fee` of it to the
    /// reserve, and it went to `receiver`; any part of `cost` beyond what the position owed returned to `buyer`.
    event Bought(
        address indexed account,
        bytes32 indexed position,
        address indexed token,
        uint256 amount,
        uint256 cost,
        uint256 fee,
        address buyer,
        address receiver
    );
    /// @notice `amount` of `account`'s `position`'s USDG repaid its debt.
    event CashSettled(address indexed account, bytes32 indexed position, uint256 amount);

    /// @notice The position meets its requirement as this contract judges it, or cannot be judged now.
    error NotLiquidatable(address account, bytes32 position);
    /// @notice The position cannot be judged now: a held asset is halted or its multiplier change is not yet
    /// confirmed, or the sequencer is not settled.
    error CannotJudge(address account, bytes32 position);
    /// @notice The position has no auction running in the market's current state.
    error NoAuction(address account, bytes32 position);
    /// @notice The auction still runs on a position that falls short.
    error StillShort(address account, bytes32 position);
    /// @notice `token` is not sold in auctions: USDG settles through `settleCash`.
    error NotForSale(address token);
    /// @notice The purchase would cost `cost`, more than the buyer's `maxCost`.
    error CostAboveLimit(uint256 cost, uint256 maxCost);
    /// @notice Nothing is left to buy: the position owes nothing, holds none of the token, or the hour's share of
    /// it is spent.
    error NothingToBuy();
    /// @notice Only the accounts' backstop writes off a position once they have one.
    error NotBackstop(address caller);

    /// @param accounts_ The accounts, which must then set this contract as their liquidator.
    constructor(MarginAccounts accounts_) {
        accounts = accounts_;
        band = accounts_.band();
        engine = accounts_.engine();
        usdg = accounts_.usdg();
        weth = accounts_.weth();
        ethUsd = accounts_.ethUsd();
        usdg.forceApprove(address(accounts_), type(uint256).max);
    }

    /// @notice `account`'s `position` as this contract judges it now, at the holdings as last synced: its equity and
    /// requirement in USD with 18 decimals, at the edge that leaves it the larger surplus while the market is closed,
    /// whether it falls short, and whether the market is closed. A position that cannot be judged now reads as not
    /// short.
    function shortfall(address account, bytes32 position)
        external
        view
        returns (int256 equity, uint256 requirement, bool short, bool closed)
    {
        closed = _closed();
        Valuation memory v = _assess(account, position, closed);
        return (v.equity, v.requirement, _isShort(v), closed);
    }

    /// @notice The price an auction of `account`'s `position` asks now for `token`, in USD with 8 decimals per whole
    /// token; zero while no auction runs in the market's current state.
    function price(address account, bytes32 position, address token) public view returns (uint256) {
        Auction memory a = auctions[account][position];
        bool closed = _closed();
        if (!_live(a, closed)) return 0;
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        return _price(token, _symbol(token, symbols, tokens), closed, block.timestamp - a.startedAt);
    }

    /// @notice How much of `token` a purchase from `account`'s `position` may take now while the market is closed: a
    /// tenth of the position's holding an hour, refilled continuously.
    function hourlyAllowance(address account, bytes32 position, address token) external view returns (uint256 left) {
        (left,) = _allowance(account, position, token);
    }

    /// @notice Starts an auction of `account`'s `position` if it falls short, syncing its assets first. An auction
    /// already running in the market's current state runs on. Anyone may call it.
    function start(address account, bytes32 position) external {
        bool closed = _closed();
        if (!_isShort(_judge(account, position, closed))) revert NotLiquidatable(account, position);
        Auction storage a = auctions[account][position];
        if (_live(a, closed)) return;
        (a.startedAt, a.closed) = (SafeCast.toUint64(block.timestamp), closed);
        // forge-lint: disable-next-line(reentrancy-events)
        emit AuctionStarted(account, position, closed);
    }

    /// @notice Stops the auction of `account`'s `position` once it no longer falls short, so a later shortfall starts
    /// from the top again. Anyone may call it, the borrower included.
    function stop(address account, bytes32 position) external {
        Valuation memory v = _judge(account, position, _closed());
        if (!v.judged) revert CannotJudge(account, position);
        if (_isShort(v)) revert StillShort(account, position);
        delete auctions[account][position];
        // forge-lint: disable-next-line(reentrancy-events)
        emit AuctionStopped(account, position);
    }

    /// @notice Buys up to `amount` of `token` from `account`'s `position` at the auction's price, for at most `maxCost`
    /// of the caller's USDG, and sends it to `receiver`. The purchase stops at the collateral held, at what repays half
    /// of what the position owes, or all of it once that is 2,000 USDG or less, and, while the market is closed, at
    /// the position's hourly share of a Stock Token. Returns what it bought and cost.
    function buy(address account, bytes32 position, address token, uint256 amount, uint256 maxCost, address receiver)
        external
        returns (uint256 bought, uint256 cost)
    {
        bool closed = _closed();
        Valuation memory v = _judge(account, position, closed);
        if (!_isShort(v)) revert NotLiquidatable(account, position);
        (bought, cost) = _terms(account, position, token, amount, closed, v);
        if (cost > maxCost) revert CostAboveLimit(cost, maxCost);
        uint256 fee = cost * _feeBps(v) / BPS;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Bought(account, position, token, bought, cost, fee, msg.sender, receiver);
        _pay(account, position, cost, fee);
        accounts.seize(position, token, bought, account, receiver);
    }

    /// @notice Repays the debt of `account`'s `position`, then its premium, with the USDG it holds, if it falls short.
    /// It sells nothing, so it takes no fee and needs no auction: it settles a position whose Stock Tokens the issuer
    /// has frozen, as far as its cash goes. Anyone may call it.
    function settleCash(address account, bytes32 position) external returns (uint256 settled) {
        Valuation memory v = _judge(account, position, _closed());
        if (!_isShort(v)) revert NotLiquidatable(account, position);
        settled = _min(accounts.collateral(account, position, address(usdg)), v.owed);
        require(settled != 0, NothingToBuy());
        // forge-lint: disable-next-line(reentrancy-events)
        emit CashSettled(account, position, settled);
        accounts.seize(position, address(usdg), settled, account, address(this));
        // slither-disable-next-line unused-return
        accounts.repay(position, settled, account); // forge-lint: disable-line(unused-return)
    }

    /// @notice Writes off what `account`'s `position` still owes once it holds nothing, or only holdings a burn left
    /// worthless: lenders bear it, and its premium is forgiven. Once the accounts have a backstop, only it may call
    /// this, after covering what it can; until then, anyone may.
    function writeOff(address account, bytes32 position) external returns (uint256) {
        address backstop = accounts.backstop();
        require(backstop == address(0) || msg.sender == backstop, NotBackstop(msg.sender));
        delete auctions[account][position];
        return accounts.writeOff(account, position);
    }

    function _terms(address account, bytes32 position, address token, uint256 amount, bool closed, Valuation memory v)
        private
        returns (uint256 bought, uint256 cost)
    {
        Auction memory a = auctions[account][position];
        if (!_live(a, closed)) revert NoAuction(account, position);
        bytes32 symbol = _symbol(token, v.symbols, v.tokens);
        uint256 unit = _price(token, symbol, closed, block.timestamp - a.startedAt);
        bought = _min(amount, accounts.collateral(account, position, token));
        bought = _min(bought, _repayable(v) * TOKEN_TO_USDG / unit);
        if (closed && symbol != bytes32(0)) bought = _spend(account, position, token, bought);
        require(bought != 0, NothingToBuy());
        cost = _divUp(bought * unit, TOKEN_TO_USDG);
    }

    function _repayable(Valuation memory v) private pure returns (uint256) {
        uint256 owed = v.owed > FULL_CLOSE_OWED ? v.owed * CLOSE_FACTOR_BPS / BPS : v.owed;
        return _divUp(owed * BPS, BPS - _feeBps(v));
    }

    function _spend(address account, bytes32 position, address token, uint256 wanted) private returns (uint256 spent) {
        (uint256 left, uint256 used) = _allowance(account, position, token);
        spent = _min(wanted, left);
        _allowances[account][position][token] =
            Allowance(SafeCast.toUint192(used + spent), SafeCast.toUint64(block.timestamp));
    }

    function _pay(address account, bytes32 position, uint256 cost, uint256 fee) private {
        usdg.safeTransferFrom(msg.sender, address(this), cost);
        if (fee != 0) accounts.collectFee(fee);
        uint256 repaid = accounts.repay(position, cost - fee, account);
        if (repaid < cost - fee) usdg.safeTransfer(msg.sender, cost - fee - repaid);
    }

    function _judge(address account, bytes32 position, bool closed) private returns (Valuation memory) {
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        for (uint256 i; i < symbols.length; ++i) {
            // forge-lint: disable-next-line(calls-loop)
            if (tokens[i] != address(0) && accounts.collateral(account, position, tokens[i]) != 0) {
                accounts.sync(symbols[i]); // forge-lint: disable-line(calls-loop)
            }
        }
        return _assess(account, position, closed);
    }

    function _assess(address account, bytes32 position, bool closed) private view returns (Valuation memory v) {
        v = _value(account, position, false);
        if (!closed || !_isShort(v)) return v;
        Valuation memory atHigh = _value(account, position, true);
        if (_surplus(atHigh) > _surplus(v)) return atHigh;
    }

    function _isShort(Valuation memory v) private pure returns (bool) {
        return v.equity < SafeCast.toInt256(v.requirement);
    }

    function _surplus(Valuation memory v) private pure returns (int256) {
        return v.equity - SafeCast.toInt256(v.requirement);
    }

    function _feeBps(Valuation memory v) private pure returns (uint256) {
        return v.insolvent ? 0 : FEE_BPS;
    }

    function _value(address account, bytes32 position, bool atHigh) private view returns (Valuation memory v) {
        if (!band.sequencerSettled()) return v;
        (v.symbols, v.tokens) = accounts.stocks();
        v.quantities = new int256[](v.symbols.length);
        v.prices = new uint256[](v.symbols.length);
        uint256 gross = 0;
        for (uint256 i; i < v.symbols.length; ++i) {
            if (v.tokens[i] == address(0)) continue;
            // forge-lint: disable-next-line(calls-loop)
            uint256 quantity = accounts.collateral(account, position, v.tokens[i]);
            if (quantity == 0) continue;
            // forge-lint: disable-start(calls-loop, unused-return)
            // slither-disable-next-line unused-return,calls-loop
            (uint8 state,,,, uint64 low, uint128 high) = band.quote(v.symbols[i]);
            // slither-disable-next-line unused-return,calls-loop
            (uint8 status,,,) = band.corporateAction(v.symbols[i]);
            // forge-lint: disable-end(calls-loop, unused-return)
            if (state == 0 || status == 2) return v;
            uint256 unit = atHigh ? high : low;
            v.quantities[i] = SafeCast.toInt256(quantity);
            v.prices[i] = unit;
            v.held = true;
            gross += quantity * unit / PRICE_UNIT;
        }
        uint256 cash = accounts.collateral(account, position, address(usdg)) * USDG_TO_USD;
        uint256 eth = accounts.collateral(account, position, address(weth));
        if (eth != 0) {
            uint256 ethPrice = _ethPrice();
            if (ethPrice == 0) return v;
            eth = _wethWorth(eth, ethPrice);
        }
        v.owed = accounts.debt(account, position) + accounts.premium(account, position);
        v.insolvent = gross + cash + eth < v.owed * USDG_TO_USD;
        v.equity = SafeCast.toInt256(gross + cash + eth * WETH_LIQUIDATION_BPS / BPS)
            - SafeCast.toInt256(v.owed * USDG_TO_USD);
        if (v.held) v.requirement = _requirement(v);
        v.judged = true;
    }

    function _requirement(Valuation memory v) private view returns (uint256 requirement) {
        uint8 regime;
        // forge-lint: disable-start(unused-return)
        // slither-disable-next-line unused-return
        (requirement,, regime) = engine.currentRequirement(v.quantities, v.prices);
        if (regime == 0) {
            // slither-disable-next-line unused-return
            (uint256 open,) = engine.requirement(v.quantities, v.prices, HORIZON, false);
            requirement = _divUp(open * 5, 4);
        }
        // forge-lint: disable-end(unused-return)
    }

    function _price(address token, bytes32 symbol, bool closed, uint256 elapsed) private view returns (uint256) {
        (uint256 top, uint256 floor) = _range(token, symbol);
        uint256 decay = (closed ? CLOSED_DECAY_BPS_PER_MINUTE : OPEN_DECAY_BPS_PER_MINUTE) * elapsed / 1 minutes;
        uint256 unit = top * (BPS - decay) / BPS;
        return unit > floor ? unit : floor;
    }

    function _range(address token, bytes32 symbol) private view returns (uint256 top, uint256 floor) {
        if (symbol != bytes32(0)) {
            // slither-disable-next-line unused-return
            (uint8 state,,,, uint64 low, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
            require(state != 0, NotForSale(token));
            return (high, uint256(low) * (BPS - MAX_DISCOUNT_BPS) / BPS);
        }
        require(token == address(weth), NotForSale(token));
        top = _ethPrice();
        require(top != 0, NotForSale(token));
        floor = top * (BPS - MAX_DISCOUNT_BPS) / BPS;
    }

    function _allowance(address account, bytes32 position, address token)
        private
        view
        returns (uint256 left, uint256 used)
    {
        Allowance memory a = _allowances[account][position][token];
        uint256 held = accounts.collateral(account, position, token);
        uint256 cap = held * CLOSED_HOURLY_BPS / BPS;
        // forge-lint: disable-start(block-timestamp)
        uint256 refill = held * CLOSED_HOURLY_BPS * (block.timestamp - a.updatedAt) / (BPS * 1 hours);
        used = a.used > refill ? a.used - refill : 0;
        left = cap > used ? cap - used : 0;
        // forge-lint: disable-end(block-timestamp)
    }

    function _symbol(address token, bytes32[] memory symbols, address[] memory tokens) private pure returns (bytes32) {
        for (uint256 i; i < tokens.length; ++i) {
            if (tokens[i] == token && token != address(0)) return symbols[i];
        }
        return bytes32(0);
    }

    function _closed() private view returns (bool) {
        // forge-lint: disable-start(unused-return, block-timestamp)
        // slither-disable-next-line unused-return
        (uint8 state, uint8 nyse, uint8 nyseNext,, uint64 boundaryMs) = band.session();
        if (state != 2 || (boundaryMs != 0 && boundaryMs <= block.timestamp * 1000)) return true;
        // forge-lint: disable-end(unused-return, block-timestamp)
        return nyse != 1 && (nyse != 2 || nyseNext != 1);
    }

    function _live(Auction memory a, bool closed) private view returns (bool) {
        uint256 lifetime = closed ? CLOSED_AUCTION_LIFETIME : OPEN_AUCTION_LIFETIME;
        // forge-lint: disable-next-line(block-timestamp)
        return a.startedAt != 0 && a.closed == closed && block.timestamp < a.startedAt + lifetime;
    }

    function _wethWorth(uint256 amount, uint256 ethPrice) private pure returns (uint256) {
        return amount * ethPrice / PRICE_UNIT;
    }

    function _ethPrice() private view returns (uint256) {
        if (address(ethUsd) == address(0)) return 0;
        // slither-disable-next-line unused-return
        try ethUsd.latestRoundData() returns (uint80, int256 answer, uint256, uint256 updatedAt, uint80) {
            // forge-lint: disable-next-line(block-timestamp)
            if (answer <= 0 || updatedAt > block.timestamp || block.timestamp - updatedAt > MAX_FEED_AGE) return 0;
            return SafeCast.toUint256(answer);
        } catch {
            return 0;
        }
    }

    function _divUp(uint256 a, uint256 b) private pure returns (uint256) {
        return a > 0 ? (a - 1) / b + 1 : 0;
    }

    function _min(uint256 a, uint256 b) private pure returns (uint256) {
        return a < b ? a : b;
    }
}
