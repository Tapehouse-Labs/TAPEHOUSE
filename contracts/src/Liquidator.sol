// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IMargin} from "./interfaces/IMargin.sol";
import {Basket} from "./Basket.sol";
import {MarginAccounts} from "./MarginAccounts.sol";

/// @title Tapehouse liquidator
/// @notice Sells the collateral of a margin position that falls short, in a Dutch auction that follows the market. The
/// market is open while the 24/5 session is open and NYSE is in regular hours or between two trading days. Then a
/// position falls short when its equity at each band's low edge is below the engine's current requirement, and its
/// auction starts at each band's high edge and falls a basis point a second. Otherwise, from the regular close before a
/// weekend or holiday to the regular open after it, or while the session is not known, a position falls short only if
/// it does so at both edges of its bands, so a shortfall the band can explain waits for the reopening; its auction
/// falls a basis point every four seconds and sells at most a tenth of the position's holding of a Stock Token an hour.
/// An unknown session is judged against the open requirement with its buffer, never the closed one it steps up to. WETH
/// counts for 84% of its Chainlink price. Nothing is sold while a held asset is halted or its multiplier change is not
/// yet confirmed, while WETH is held and its price is missing or stale, or while the L2 sequencer is not settled. A
/// purchase repays at most half of what the position owes, or all of it once that is 2,000 USDG or less. A buyer pays
/// USDG: half a percent goes to the accounts' reserve unless the position is worth less than it owes, and the rest
/// repays the position, its debt first. Anyone may start, stop, buy and settle a position's USDG against its debt; a
/// position left with nothing, or worth less than it owes with holdings worth less than a USDG each, is written off by
/// the accounts' backstop, or by anyone while they have none. The reopening auction, set once by the accounts' owner,
/// holds the positions it will sell out of the Dutch auction while the market is open, and sells them, to its bids and
/// the backstop, at its clearing price. A basket in a position counts as the Stock Tokens it redeems for, and once the
/// position falls short its baskets are redeemed into those Stock Tokens, which are then sold as any other; a basket
/// whose tokens the issuer has frozen stays, and the rest is sold.
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
    /// @notice What a holding must be worth, in USDG, to keep a position from being written off: less is swept to the
    /// caller of `writeOff`.
    uint256 public constant DUST = 1e6;

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
    /// @notice What a lent Stock Token counts for less in a position's equity, in basis points: the accounts' own.
    uint256 public immutable recallHaircut;

    /// @notice The auction of each position, if one runs.
    mapping(address account => mapping(bytes32 position => Auction)) public auctions;
    mapping(address account => mapping(bytes32 position => mapping(address token => Allowance))) private _allowances;
    /// @notice The reopening auction, set once by the accounts' owner.
    address public auction;
    /// @notice Until when, in seconds, the reopening auction holds each position out of the Dutch auction.
    mapping(address account => mapping(bytes32 position => uint64)) public heldUntil;

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
    /// @notice The reopening auction is `auction`.
    event AuctionSet(address indexed auction);
    /// @notice The reopening auction holds `account`'s `position` out of the Dutch auction until `until`.
    event Held(address indexed account, bytes32 indexed position, uint64 until);
    /// @notice The reopening auction sold `amount` of `token` from `account`'s `position` at `price`, for `cost` of USDG,
    /// `fee` of it to the reserve.
    event Settled(
        address indexed account,
        bytes32 indexed position,
        address indexed token,
        uint256 amount,
        uint256 price,
        uint256 cost,
        uint256 fee
    );

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
    /// @notice The position has lent none of the token that it has not recalled already, or, with a recall open, less
    /// than the accounts' least recall.
    error NothingToRecall();
    /// @notice Only the accounts' backstop writes off a position once they have one.
    error NotBackstop(address caller);
    /// @notice Only the reopening auction may.
    error NotAuction(address caller);
    /// @notice Only the accounts' owner sets the reopening auction.
    error NotAccountsOwner(address caller);
    /// @notice The reopening auction is already set, to `current`.
    error AuctionAlreadySet(address current);
    /// @notice The reopening auction cannot be the zero address.
    error InvalidAuction();
    /// @notice The reopening auction holds the position until `until`.
    error PositionHeld(address account, bytes32 position, uint64 until);

    /// @param accounts_ The accounts, which must then set this contract as their liquidator.
    constructor(MarginAccounts accounts_) {
        accounts = accounts_;
        band = accounts_.band();
        engine = accounts_.engine();
        usdg = accounts_.usdg();
        weth = accounts_.weth();
        ethUsd = accounts_.ethUsd();
        recallHaircut = accounts_.RECALL_HAIRCUT();
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
        _checkNotHeld(account, position, closed);
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
        _checkNotHeld(account, position, closed);
        Valuation memory v = _judge(account, position, closed);
        if (!_isShort(v)) revert NotLiquidatable(account, position);
        (bought, cost) = _terms(account, position, token, amount, closed, v);
        if (cost > maxCost) revert CostAboveLimit(cost, maxCost);
        uint256 fee = cost * _feeBps(v) / BPS;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Bought(account, position, token, bought, cost, fee, msg.sender, receiver);
        _pay(account, position, cost, fee, msg.sender);
        accounts.seize(position, token, bought, account, receiver);
    }

    /// @notice Sets the reopening auction. Only the accounts' owner may, once, and never to zero.
    function setAuction(address newAuction) external {
        if (msg.sender != accounts.owner()) revert NotAccountsOwner(msg.sender);
        if (auction != address(0)) revert AuctionAlreadySet(auction);
        if (newAuction == address(0)) revert InvalidAuction();
        auction = newAuction;
        emit AuctionSet(newAuction);
    }

    /// @notice Holds `account`'s `position` out of the Dutch auction while the market is open, until `until`, while the
    /// reopening auction sells it. Only the reopening auction may.
    function hold(address account, bytes32 position, uint64 until) external {
        if (msg.sender != auction) revert NotAuction(msg.sender);
        heldUntil[account][position] = until;
        emit Held(account, position, until);
    }

    /// @notice Sells up to `amount` of `token` from `account`'s `position` to the caller at `unitPrice`, in USD with 8
    /// decimals per token, if the position still falls short, paid with the caller's USDG. The sale stops at the
    /// collateral held and at what repays all the position owes, fee included; half a percent goes to the reserve unless
    /// the position is worth less than it owes, the rest repays the position, and any rounding excess goes to the
    /// reserve. Returns what it sold and cost; nothing if the position no longer falls short. Only the reopening auction
    /// and the accounts' backstop, for what the auction's bids leave, may.
    function settle(address account, bytes32 position, address token, uint256 amount, uint256 unitPrice)
        external
        returns (uint256 sold, uint256 cost)
    {
        if (msg.sender != auction && msg.sender != accounts.backstop()) revert NotAuction(msg.sender);
        Valuation memory v = _judge(account, position, _closed());
        if (!_isShort(v)) return (0, 0);
        require(_symbol(token, v.symbols, v.tokens) != bytes32(0), NotForSale(token));
        uint256 feeBps = _feeBps(v);
        sold = _min(amount, accounts.sellable(account, position, token));
        sold = _min(sold, _divUp(v.owed * BPS, BPS - feeBps) * TOKEN_TO_USDG / unitPrice);
        if (sold != 0) {
            cost = _divUp(sold * unitPrice, TOKEN_TO_USDG);
            uint256 fee = cost * feeBps / BPS;
            // forge-lint: disable-next-line(reentrancy-events)
            emit Settled(account, position, token, sold, unitPrice, cost, fee);
            _pay(account, position, cost, fee, address(0));
            accounts.seize(position, token, sold, account, msg.sender);
        }
    }

    /// @notice Adds `amount` of the reopening auction's USDG to the accounts' reserve: the bonds its bidders forfeit.
    /// Only the reopening auction may.
    function collect(uint256 amount) external {
        if (msg.sender != auction) revert NotAuction(msg.sender);
        usdg.safeTransferFrom(msg.sender, address(this), amount);
        accounts.collectFee(amount);
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

    /// @notice Recalls, for `account`'s `position` once it falls short, all it lent of `token` that it has not
    /// recalled yet, so that what its lending vault cannot return now comes back within the vault's notice, or by a
    /// buy-in, for the auction to sell. Anyone may call it, and it returns what it recalled.
    function recall(address account, bytes32 position, address token) external returns (uint256 amount) {
        if (!_isShort(_judge(account, position, _closed()))) revert NotLiquidatable(account, position);
        uint256 lent = accounts.lent(account, position, token);
        uint256 claimed = accounts.claim(account, position, token);
        if (lent <= claimed || (claimed != 0 && lent - claimed < accounts.MIN_RECALL())) revert NothingToRecall();
        amount = lent - claimed;
        accounts.recall(position, token, amount, account);
    }

    /// @notice Writes off what `account`'s `position` still owes once it holds nothing, or only holdings a burn left
    /// worthless: lenders bear it, and its premium is forgiven. If the accounts refuse it and the position is worth less
    /// than it owes, holdings worth less than `DUST` each, a Stock Token's with what its lending vault can return of
    /// what the position lent, are swept to the caller first. Once the accounts have a backstop, only it may call this,
    /// after covering what it can; until then, anyone may.
    function writeOff(address account, bytes32 position) external returns (uint256) {
        address backstop = accounts.backstop();
        require(backstop == address(0) || msg.sender == backstop, NotBackstop(msg.sender));
        delete auctions[account][position];
        try accounts.writeOff(account, position) returns (uint256 written) {
            return written;
        } catch {
            Valuation memory v = _judge(account, position, _closed());
            if (v.insolvent) _sweep(account, position);
            return accounts.writeOff(account, position);
        }
    }

    function _terms(address account, bytes32 position, address token, uint256 amount, bool closed, Valuation memory v)
        private
        returns (uint256 bought, uint256 cost)
    {
        Auction memory a = auctions[account][position];
        if (!_live(a, closed)) revert NoAuction(account, position);
        bytes32 symbol = _symbol(token, v.symbols, v.tokens);
        uint256 unit = _price(token, symbol, closed, block.timestamp - a.startedAt);
        bought = _min(amount, accounts.sellable(account, position, token));
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

    function _pay(address account, bytes32 position, uint256 cost, uint256 fee, address excessTo) private {
        usdg.safeTransferFrom(msg.sender, address(this), cost);
        if (fee != 0) accounts.collectFee(fee);
        uint256 repaid = accounts.repay(position, cost - fee, account);
        if (repaid < cost - fee) {
            if (excessTo == address(0)) accounts.collectFee(cost - fee - repaid);
            else usdg.safeTransfer(excessTo, cost - fee - repaid);
        }
    }

    function _sweep(address account, bytes32 position) private {
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        for (uint256 i; i < symbols.length; ++i) {
            if (tokens[i] == address(0)) continue;
            // forge-lint: disable-next-line(calls-loop)
            uint256 held = accounts.sellable(account, position, tokens[i]);
            if (held == 0) continue;
            // slither-disable-next-line unused-return,calls-loop
            (,,,, uint64 low,) = band.quote(symbols[i]); // forge-lint: disable-line(unused-return, calls-loop)
            if (held * low / TOKEN_TO_USDG < DUST) _seizeTo(account, position, tokens[i], held);
        }
        uint256 eth = accounts.collateral(account, position, address(weth));
        if (eth != 0 && eth * _ethPrice() / TOKEN_TO_USDG < DUST) {
            _seizeTo(account, position, address(weth), eth);
        }
        uint256 cash = accounts.collateral(account, position, address(usdg));
        if (cash != 0 && cash < DUST) _seizeTo(account, position, address(usdg), cash);
    }

    function _seizeTo(address account, bytes32 position, address token, uint256 amount) private {
        accounts.seize(position, token, amount, account, msg.sender); // forge-lint: disable-line(calls-loop)
    }

    function _checkNotHeld(address account, bytes32 position, bool closed) private view {
        uint64 until = heldUntil[account][position];
        // forge-lint: disable-next-line(block-timestamp)
        if (!closed && block.timestamp < until) revert PositionHeld(account, position, until);
    }

    function _judge(address account, bytes32 position, bool closed) private returns (Valuation memory) {
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        for (uint256 i; i < symbols.length; ++i) {
            // forge-lint: disable-start(calls-loop)
            if (
                tokens[i] != address(0)
                    && (accounts.collateral(account, position, tokens[i]) != 0
                        || accounts.lent(account, position, tokens[i]) != 0)
            ) {
                accounts.sync(symbols[i]);
            }
            // forge-lint: disable-end(calls-loop)
        }
        Valuation memory v = _assess(account, position, closed);
        if (_isShort(v)) _unwrap(account, position);
        return v;
    }

    /// @dev Redeems every basket share the position holds into its Stock Tokens, which count the same for its
    /// valuation, whatever they come to: shares worth nothing in each token would otherwise hold up its write-off.
    function _unwrap(address account, bytes32 position) private {
        Basket[] memory baskets = accounts.baskets();
        for (uint256 j; j < baskets.length; ++j) {
            // forge-lint: disable-next-line(calls-loop)
            uint256 shares = accounts.collateral(account, position, address(baskets[j]));
            if (shares == 0) continue;
            // slither-disable-start unused-return,calls-loop
            // forge-lint: disable-next-line(calls-loop)
            try accounts.unwrap(account, address(baskets[j]), shares) returns (uint256[] memory) {} catch {}
            // slither-disable-end unused-return,calls-loop
        }
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
        uint256[] memory basketed = accounts.inBaskets(account, position);
        v.quantities = new int256[](v.symbols.length);
        v.prices = new uint256[](v.symbols.length);
        uint256 gross = 0;
        uint256 haircut = 0;
        for (uint256 i; i < v.symbols.length; ++i) {
            if (v.tokens[i] == address(0)) continue;
            // forge-lint: disable-start(calls-loop)
            uint256 lent = accounts.lent(account, position, v.tokens[i]);
            uint256 quantity = accounts.collateral(account, position, v.tokens[i]) + lent + basketed[i];
            // forge-lint: disable-end(calls-loop)
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
            haircut += lent * unit * recallHaircut / (PRICE_UNIT * BPS);
        }
        gross -= haircut;
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
        uint256 held = accounts.collateral(account, position, token) + accounts.lent(account, position, token);
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
