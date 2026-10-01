// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {BandFeed} from "./BandFeed.sol";
import {GapBackstop} from "./GapBackstop.sol";
import {IBand} from "./interfaces/IBand.sol";
import {Liquidator} from "./Liquidator.sol";
import {MarginAccounts} from "./MarginAccounts.sol";

/// @title Tapehouse reopening auction
/// @notice Sells the collateral of positions that fall short over a weekend or holiday closure in one sealed-bid,
/// uniform-price batch per Stock Token at NYSE's regular open, and the gap backstop takes what the bids leave. From
/// the 24/5 session's reopen until half an hour before the regular open, anyone may enroll a position that falls short
/// at its bands' low edges, which the liquidator then holds out of its open-market Dutch auction, and bidders commit
/// to hidden bids with a deposit. In the last half hour they reveal them, escrowing what they bid. In the first hour of
/// regular trading, anyone submits the clearing price, which the contract checks against the revealed bids and the
/// supply; the liquidator sells each lot's share to the bids at it, and the backstop buys the rest at it, at most the
/// band's low edge, within its exposure limit; a backstop that cannot buy leaves the rest in the position. Bidders then claim their tokens and the rest of their escrow. The floor is 10% below the low edge of
/// the band sealed before the reopen.
/// @dev Each round is keyed by its Stock Token and the regular open, in milliseconds, that the band's session names.
contract ReopeningAuction {
    using SafeERC20 for IERC20;

    struct Round {
        uint64 sealMs;
        uint128 floor;
        uint128 supply;
        uint128 deposits;
        uint128 pool;
        uint128 price;
        uint128 above;
        uint128 atPrice;
        uint128 sold;
        uint128 paid;
        uint128 taken;
        bool cleared;
    }

    struct Lot {
        address account;
        bytes32 position;
        uint128 amount;
    }

    struct Bid {
        address bidder;
        uint128 quantity;
        uint128 price;
        uint128 escrow;
        bool claimed;
    }

    struct Commitment {
        address bidder;
        uint128 deposit;
    }

    /// @notice How long before the regular open bids are revealed, in milliseconds.
    uint64 public constant REVEAL_MS = 30 minutes * 1000;
    /// @notice How long after the regular open the clearing price may be submitted, in milliseconds.
    uint64 public constant CLEAR_MS = 1 hours * 1000;
    /// @notice How long before the regular open the 24/5 session reopens, and the band is sealed, in milliseconds.
    uint64 public constant REOPEN_LEAD_MS = 13.5 hours * 1000;
    /// @notice How long after the midnight that ends a holiday NYSE opens, in milliseconds.
    uint64 public constant HOLIDAY_OPEN_MS = 9.5 hours * 1000;
    /// @notice The least deposit a commitment carries, in USDG units.
    uint256 public constant BOND = 100e6;
    /// @notice The share of an unrevealed commitment's deposit sent to the reserve, in basis points, at least `BOND`.
    uint256 public constant FORFEIT_BPS = 10_00;
    /// @notice The least a revealed bid may be worth at the floor, in USDG units.
    uint256 public constant MIN_BID = 100e6;
    /// @notice The most lots a round sells: past it, a larger lot takes the smallest one's place.
    uint256 public constant MAX_LOTS = 32;
    /// @notice The most bids a round keeps: past it, a bid at a higher price, or the same price for more, takes the
    /// lowest one's place, and any other gets its deposit back.
    uint256 public constant MAX_BIDS = 64;

    uint256 private constant BPS = 10_000;
    uint256 private constant TOKEN_TO_USDG = 1e20;

    /// @notice The liquidator that sells each lot.
    Liquidator public immutable liquidator;
    /// @notice The accounts whose positions it sells.
    MarginAccounts public immutable accounts;
    /// @notice The band whose session sets the phases.
    IBand public immutable band;
    /// @notice The token bids are paid in.
    IERC20 public immutable usdg;

    /// @notice The feed that seals each Stock Token's band before the reopen; zero where there is none.
    mapping(bytes32 symbol => BandFeed) public feeds;
    /// @notice Each commitment of a round: its bidder and deposit.
    mapping(bytes32 symbol => mapping(uint64 openMs => mapping(bytes32 commitment => Commitment))) public commitments;
    /// @notice Whether a position has a lot in the round of `symbol` for the regular open `openMs`.
    mapping(
        uint64 openMs => mapping(address account => mapping(bytes32 position => mapping(bytes32 symbol => bool)))
    ) public enrolled;
    /// @notice The regular open of the last round opened, which names the phase after a holiday's midnight.
    uint64 public lastOpenMs;
    mapping(bytes32 symbol => mapping(uint64 openMs => Round)) private _rounds;
    mapping(bytes32 symbol => mapping(uint64 openMs => Lot[])) private _lots;
    mapping(bytes32 symbol => mapping(uint64 openMs => Bid[])) private _bids;
    mapping(uint64 openMs => mapping(address account => mapping(bytes32 position => uint256))) private _lotCount;

    /// @notice The round of `symbol` for the regular open `openMs` opened with its floor, from the band sealed at
    /// `sealMs` if `fromSeal`.
    event RoundOpened(bytes32 indexed symbol, uint64 indexed openMs, uint64 sealMs, uint256 floor, bool fromSeal);
    /// @notice `account`'s `position` is enrolled, `amount` of the Stock Token to sell.
    event Enrolled(
        bytes32 indexed symbol, uint64 indexed openMs, address indexed account, bytes32 position, uint256 amount
    );
    /// @notice The lot of `account`'s `position` gave its place to a larger one, and the position left the round.
    event LotEvicted(bytes32 indexed symbol, uint64 indexed openMs, address indexed account, bytes32 position);
    /// @notice `bidder` committed to a hidden bid with `deposit`.
    event Committed(
        bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, bytes32 commitment, uint256 deposit
    );
    /// @notice `bidder` revealed a bid for `quantity` at `price`, kept at `index`.
    event Revealed(
        bytes32 indexed symbol,
        uint64 indexed openMs,
        address indexed bidder,
        uint256 index,
        uint256 quantity,
        uint256 price
    );
    /// @notice `bidder`'s bid found no place in a full book, below or equal to its lowest bid, and its deposit went back.
    event Outbid(
        bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 quantity, uint256 price
    );
    /// @notice The bid at `index` gave its place to a higher one, and its escrow went back.
    event Evicted(
        bytes32 indexed symbol,
        uint64 indexed openMs,
        address indexed bidder,
        uint256 index,
        uint256 quantity,
        uint256 price
    );
    /// @notice The round cleared at `price`: `sold` of the Stock Token to the bids for `paid`, and `taken` to the
    /// backstop.
    event Cleared(
        bytes32 indexed symbol, uint64 indexed openMs, uint256 price, uint256 sold, uint256 paid, uint256 taken
    );
    /// @notice The bid at `index` took `amount` of the Stock Token for `paid`, and `refund` of its escrow went back.
    event Claimed(
        bytes32 indexed symbol,
        uint64 indexed openMs,
        address indexed bidder,
        uint256 index,
        uint256 amount,
        uint256 paid,
        uint256 refund
    );
    /// @notice An unrevealed commitment of `bidder` lost `forfeited` of its deposit to the reserve; `returned` went back.
    event Forfeited(
        bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 forfeited, uint256 returned
    );

    /// @notice The feeds do not match the symbols or the liquidator's band.
    error InvalidFeeds();
    /// @notice The session or the clock is not in the phase the call needs.
    error WrongPhase();
    /// @notice The accounts have no Stock Token `symbol`.
    error UnknownAsset(bytes32 symbol);
    /// @notice The band of `symbol` gives no price.
    error NoPrice(bytes32 symbol);
    /// @notice The position meets its requirement at its bands' low edges, holds none of the Stock Token, or has a lot
    /// in the round already.
    error NotEnrollable(address account, bytes32 position);
    /// @notice The round sells `MAX_LOTS` lots, none smaller than this one.
    error TooManyLots();
    /// @notice The commitment is already pending, or its deposit is below `BOND`.
    error InvalidCommitment(bytes32 commitment);
    /// @notice No pending commitment matches the bid.
    error UnknownCommitment(bytes32 commitment);
    /// @notice The bid is below the floor, worth less than `MIN_BID` at the floor, or more than its deposit.
    error InvalidBid();
    /// @notice `price` is not the round's clearing price.
    error NotClearingPrice(uint256 price);
    /// @notice The round is cleared already, or cannot be claimed yet.
    error WrongState();

    /// @param liquidator_ The liquidator, whose accounts' owner must then set this contract as its auction.
    /// @param symbols The Stock Tokens whose sealed bands `feeds_` give.
    /// @param feeds_ Each symbol's feed, over the liquidator's band.
    constructor(Liquidator liquidator_, bytes32[] memory symbols, BandFeed[] memory feeds_) {
        if (symbols.length != feeds_.length) revert InvalidFeeds();
        liquidator = liquidator_;
        accounts = liquidator_.accounts();
        band = liquidator_.band();
        usdg = liquidator_.usdg();
        for (uint256 i; i < symbols.length; ++i) {
            // forge-lint: disable-next-line(calls-loop, require-revert-in-loop)
            if (feeds_[i].symbol() != symbols[i] || address(feeds_[i].band()) != address(band)) revert InvalidFeeds();
            feeds[symbols[i]] = feeds_[i];
        }
        usdg.forceApprove(address(liquidator_), type(uint256).max);
    }

    /// @notice Enrolls `account`'s `position` in the round of `symbol` if it falls short at its bands' low edges, as the
    /// accounts' `health` has it, and holds it out of the open-market Dutch auction until the round's clearing window
    /// closes. Its lot is what repays all it owes at the floor, fee included, at most its holding; a full round takes it
    /// only in place of a smaller lot. Anyone may call it while bids are committed.
    function enroll(address account, bytes32 position, bytes32 symbol) external {
        uint64 openMs = _commitPhase();
        address token = _token(symbol);
        uint256 floor = _open(symbol, openMs).floor;
        if (enrolled[openMs][account][position][symbol]) revert NotEnrollable(account, position);
        uint128 amount = _lotSize(account, position, token, floor);
        enrolled[openMs][account][position][symbol] = true;
        ++_lotCount[openMs][account][position];
        _addLot(symbol, openMs, Lot(account, position, amount));
        // forge-lint: disable-next-line(reentrancy-events)
        emit Enrolled(symbol, openMs, account, position, amount);
        liquidator.hold(account, position, SafeCast.toUint64((openMs + CLEAR_MS) / 1000));
    }

    /// @notice Commits the caller to a hidden bid in the round of `symbol` with `deposit` of its USDG, at least `BOND`
    /// and at least what the bid will escrow, so the deposit may hide its size. `commitment` is
    /// `keccak256(abi.encode(bidder, symbol, openMs, quantity, price, salt))`, which binds the bid to its bidder.
    function commit(bytes32 symbol, bytes32 commitment, uint256 deposit) external {
        uint64 openMs = _commitPhase();
        _token(symbol);
        Round storage r = _open(symbol, openMs);
        Commitment storage c = commitments[symbol][openMs][commitment];
        if (c.deposit != 0 || deposit < BOND) revert InvalidCommitment(commitment);
        (c.bidder, c.deposit) = (msg.sender, SafeCast.toUint128(deposit));
        r.deposits += SafeCast.toUint128(deposit);
        emit Committed(symbol, openMs, msg.sender, commitment, deposit);
        usdg.safeTransferFrom(msg.sender, address(this), deposit);
    }

    /// @notice Reveals the caller's bid for `quantity` of `symbol`'s Stock Token, with 18 decimals, at `price` in USD
    /// with 8 decimals per token, in the round for the regular open `openMs`, in the half hour before it. The bid's
    /// escrow, what it is worth at its price, stays from the deposit and the rest goes back. A bid below the floor, worth
    /// less than `MIN_BID` at the floor or more than its deposit is refused; past `MAX_BIDS`, it takes the lowest bid's
    /// place if it is higher, whose escrow goes back, and otherwise its whole deposit goes back.
    function reveal(bytes32 symbol, uint64 openMs, uint256 quantity, uint256 price, bytes32 salt) external {
        Round storage r = _rounds[symbol][openMs];
        _checkRevealWindow(openMs);
        uint256 deposit =
            _takeCommitment(symbol, openMs, keccak256(abi.encode(msg.sender, symbol, openMs, quantity, price, salt)));
        uint256 escrow = Math.mulDiv(quantity, price, TOKEN_TO_USDG, Math.Rounding.Ceil);
        if (price < r.floor || Math.mulDiv(quantity, r.floor, TOKEN_TO_USDG) < MIN_BID || escrow > deposit) {
            revert InvalidBid();
        }
        (bool placed, uint256 index, Bid memory evicted) = _place(symbol, openMs, _bid(quantity, price, escrow));
        if (!placed) {
            // forge-lint: disable-next-line(reentrancy-events)
            emit Outbid(symbol, openMs, msg.sender, quantity, price);
            usdg.safeTransfer(msg.sender, deposit);
            return;
        }
        r.pool += SafeCast.toUint128(escrow);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Revealed(symbol, openMs, msg.sender, index, quantity, price);
        if (deposit != escrow) usdg.safeTransfer(msg.sender, deposit - escrow);
        if (evicted.escrow != 0) usdg.safeTransfer(evicted.bidder, evicted.escrow);
    }

    /// @notice Sends a share of an unrevealed commitment's deposit, `FORFEIT_BPS` of it and at least `BOND`, to the
    /// accounts' reserve, and the rest back to its bidder, once the round's bids are revealed. Anyone may call it.
    function forfeit(bytes32 symbol, uint64 openMs, bytes32 commitment) external {
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp * 1000 < openMs) revert WrongPhase();
        Commitment memory c = commitments[symbol][openMs][commitment];
        if (c.deposit == 0) revert UnknownCommitment(commitment);
        delete commitments[symbol][openMs][commitment];
        _rounds[symbol][openMs].deposits -= c.deposit;
        uint256 forfeited = Math.min(c.deposit, Math.max(BOND, uint256(c.deposit) * FORFEIT_BPS / BPS));
        // forge-lint: disable-next-line(reentrancy-events)
        emit Forfeited(symbol, openMs, c.bidder, forfeited, c.deposit - forfeited);
        liquidator.collect(forfeited);
        if (c.deposit != forfeited) usdg.safeTransfer(c.bidder, c.deposit - forfeited);
    }

    /// @notice Clears the round of `symbol` for the regular open `openMs` at `price`: each lot sells its share of what
    /// the bids take at that price, and the backstop buys the rest at it, at most the band's low edge. `price` is the
    /// round's only clearing price: the
    /// highest bid price at which the bids at or above it take the whole supply; when all of them together take less,
    /// the lowest bid price; with no bid, the floor. Anyone may call it in the first `CLEAR_MS` of regular trading.
    // slither-disable-next-line reentrancy-no-eth
    function clear(bytes32 symbol, uint64 openMs, uint256 price) external {
        _checkClearWindow(openMs);
        Round storage r = _rounds[symbol][openMs];
        if (r.cleared || r.floor == 0) revert WrongState();
        _settlePrice(symbol, openMs, price);
        _sell(symbol, openMs, price);
        r.pool -= r.paid;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Cleared(symbol, openMs, price, r.sold, r.paid, r.taken);
    }

    /// @notice Sends the bidder of the bid at `index` its share of what the bids bought and the rest of its escrow, once
    /// the round is cleared; once its clearing window has closed uncleared, all its escrow. Anyone may call it.
    function claim(bytes32 symbol, uint64 openMs, uint256 index) external {
        Round storage r = _rounds[symbol][openMs];
        Bid storage bid = _bids[symbol][openMs][index];
        // forge-lint: disable-next-line(block-timestamp)
        bool lapsed = !r.cleared && block.timestamp * 1000 >= uint256(openMs) + CLEAR_MS;
        if (bid.claimed || !(r.cleared || lapsed)) revert WrongState();
        bid.claimed = true;
        uint256 amount = r.cleared ? _fill(r, bid) : 0;
        uint256 paid = amount == 0 ? 0 : Math.min(Math.mulDiv(r.paid, amount, r.sold, Math.Rounding.Ceil), bid.escrow);
        uint256 refund = Math.min(bid.escrow - paid, r.pool);
        r.pool -= SafeCast.toUint128(refund);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Claimed(symbol, openMs, bid.bidder, index, amount, paid, refund);
        if (amount != 0) IERC20(_token(symbol)).safeTransfer(bid.bidder, amount);
        if (refund != 0) usdg.safeTransfer(bid.bidder, refund);
    }

    /// @notice The round of `symbol` for the regular open `openMs`.
    function round(bytes32 symbol, uint64 openMs) external view returns (Round memory) {
        return _rounds[symbol][openMs];
    }

    /// @notice The lots of a round.
    function lots(bytes32 symbol, uint64 openMs) external view returns (Lot[] memory) {
        return _lots[symbol][openMs];
    }

    /// @notice The revealed bids of a round.
    function bids(bytes32 symbol, uint64 openMs) external view returns (Bid[] memory) {
        return _bids[symbol][openMs];
    }

    /// @notice The regular open that bids are committed or revealed for, and whether they are being revealed; zero
    /// outside those phases. The session names it from the 24/5 reopen after a weekend; after a holiday, until the
    /// midnight that ends it, and from then on only if a round opened before that midnight.
    function phase() public view returns (uint64 openMs, bool revealing) {
        // slither-disable-next-line unused-return
        (uint8 state, uint8 nyse, uint8 nyseNext, uint64 changeMs, uint64 boundaryMs) = band.session();
        uint256 nowMs = block.timestamp * 1000;
        if (state != 2 || boundaryMs != 0) return (openMs, revealing);
        if (nyse == 3 && nyseNext == 1) openMs = changeMs;
        else if (nyse == 3 && nyseNext == 2) openMs = changeMs + HOLIDAY_OPEN_MS;
        else if (nyse == 2) openMs = lastOpenMs;
        // forge-lint: disable-next-line(block-timestamp)
        if (openMs <= nowMs) return (0, revealing);
        // forge-lint: disable-next-line(block-timestamp)
        revealing = nowMs + REVEAL_MS >= openMs;
    }

    function _open(bytes32 symbol, uint64 openMs) private returns (Round storage r) {
        r = _rounds[symbol][openMs];
        if (r.floor != 0) return r;
        uint64 sealMs = openMs - REOPEN_LEAD_MS;
        uint256 low = 0;
        bool fromSeal = false;
        BandFeed feed = feeds[symbol];
        if (address(feed) != address(0)) {
            // slither-disable-next-line unused-return
            (uint8 state,,,, uint64 sealedLow,,) = feed.seals(sealMs); // forge-lint: disable-line(unused-return)
            fromSeal = state != 0;
            if (fromSeal) low = sealedLow;
        }
        if (!fromSeal) {
            // slither-disable-next-line unused-return
            (,,,, uint64 currentLow,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
            low = currentLow;
        }
        uint256 floor = low * (BPS - liquidator.MAX_DISCOUNT_BPS()) / BPS;
        if (floor == 0) revert NoPrice(symbol);
        (r.sealMs, r.floor) = (sealMs, SafeCast.toUint128(floor));
        lastOpenMs = openMs;
        emit RoundOpened(symbol, openMs, sealMs, floor, fromSeal);
    }

    function _lotSize(address account, bytes32 position, address token, uint256 floor) private view returns (uint128) {
        uint256 holding = accounts.collateral(account, position, token);
        // slither-disable-next-line unused-return
        (int256 equity, uint256 requirement,,) = accounts.health(account, position); // forge-lint: disable-line(unused-return)
        if (holding == 0 || equity >= SafeCast.toInt256(requirement)) revert NotEnrollable(account, position);
        uint256 owed = accounts.debt(account, position) + accounts.premium(account, position);
        uint256 need = Math.mulDiv(owed, BPS * TOKEN_TO_USDG, (BPS - liquidator.FEE_BPS()) * floor, Math.Rounding.Ceil);
        return SafeCast.toUint128(Math.min(holding, need));
    }

    function _addLot(bytes32 symbol, uint64 openMs, Lot memory lot) private {
        Round storage r = _rounds[symbol][openMs];
        r.supply += lot.amount;
        if (_lots[symbol][openMs].length < MAX_LOTS) {
            _lots[symbol][openMs].push(lot);
            return;
        }
        uint256 smallest = _smallest(_lots[symbol][openMs]);
        Lot memory out = _lots[symbol][openMs][smallest];
        if (out.amount >= lot.amount) revert TooManyLots();
        r.supply -= out.amount;
        _lots[symbol][openMs][smallest] = lot;
        // forge-lint: disable-next-line(reentrancy-events)
        emit LotEvicted(symbol, openMs, out.account, out.position);
        _leave(symbol, openMs, out);
    }

    function _settlePrice(bytes32 symbol, uint64 openMs, uint256 price) private {
        Round storage r = _rounds[symbol][openMs];
        Bid[] storage all = _bids[symbol][openMs];
        (uint256 above, uint256 atPrice, bool below) = _demand(all, price);
        uint256 supply = r.supply;
        bool full = above < supply && above + atPrice >= supply;
        bool short = !below && atPrice != 0 && above + atPrice < supply;
        bool none = all.length == 0 && price == r.floor;
        if (!(full || short || none)) revert NotClearingPrice(price);
        r.cleared = true;
        (r.price, r.above, r.atPrice) =
        (SafeCast.toUint128(price), SafeCast.toUint128(above), SafeCast.toUint128(atPrice));
    }

    // slither-disable-next-line reentrancy-no-eth
    function _sell(bytes32 symbol, uint64 openMs, uint256 price) private {
        Round storage r = _rounds[symbol][openMs];
        Lot[] storage all = _lots[symbol][openMs];
        uint256 allocated = Math.min(r.supply, uint256(r.above) + r.atPrice);
        address token = _token(symbol);
        uint256 cap = _backstopPrice(symbol, price);
        for (uint256 i; i < all.length; ++i) {
            Lot memory lot = all[i];
            _leave(symbol, openMs, lot);
            uint256 share = Math.min(lot.amount * allocated / r.supply, (r.pool - r.paid) * TOKEN_TO_USDG / price);
            _sellLot(r, lot, token, share, price, cap);
        }
    }

    // slither-disable-next-line reentrancy-no-eth
    function _sellLot(Round storage r, Lot memory lot, address token, uint256 share, uint256 price, uint256 cap)
        private
    {
        if (share != 0) {
            // forge-lint: disable-next-line(calls-loop)
            (uint256 sold, uint256 paid) = liquidator.settle(lot.account, lot.position, token, share, price);
            r.sold += SafeCast.toUint128(sold);
            r.paid += SafeCast.toUint128(paid);
        }
        // forge-lint: disable-next-line(calls-loop)
        address backstop = accounts.backstop();
        if (backstop != address(0) && lot.amount > share) {
            // forge-lint: disable-next-line(calls-loop)
            try GapBackstop(backstop).buyRemainder(lot.account, lot.position, token, lot.amount - share, cap) returns (
                uint256 taken
            ) {
                r.taken += SafeCast.toUint128(taken);
            } catch {}
        }
    }

    function _backstopPrice(bytes32 symbol, uint256 price) private view returns (uint256) {
        // slither-disable-next-line unused-return
        (,,,, uint64 low,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        return Math.min(price, low);
    }

    function _leave(bytes32 symbol, uint64 openMs, Lot memory lot) private {
        enrolled[openMs][lot.account][lot.position][symbol] = false;
        if (--_lotCount[openMs][lot.account][lot.position] == 0) {
            // forge-lint: disable-next-line(calls-loop)
            liquidator.hold(lot.account, lot.position, 0);
        }
    }

    function _place(bytes32 symbol, uint64 openMs, Bid memory bid)
        private
        returns (bool placed, uint256 index, Bid memory evicted)
    {
        Bid[] storage all = _bids[symbol][openMs];
        index = all.length;
        if (index < MAX_BIDS) {
            all.push(bid);
            placed = true;
            return (placed, index, evicted);
        }
        index = _lowest(all);
        evicted = all[index];
        if (evicted.price > bid.price || (evicted.price == bid.price && evicted.quantity >= bid.quantity)) {
            return (placed, index, evicted);
        }
        placed = true;
        all[index] = bid;
        _rounds[symbol][openMs].pool -= evicted.escrow;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Evicted(symbol, openMs, evicted.bidder, index, evicted.quantity, evicted.price);
    }

    function _bid(uint256 quantity, uint256 price, uint256 escrow) private view returns (Bid memory) {
        return
            Bid(msg.sender, SafeCast.toUint128(quantity), SafeCast.toUint128(price), SafeCast.toUint128(escrow), false);
    }

    function _commitPhase() private view returns (uint64 openMs) {
        bool revealing;
        (openMs, revealing) = phase();
        require(openMs != 0 && !revealing, WrongPhase());
    }

    function _checkRevealWindow(uint64 openMs) private view {
        uint256 nowMs = block.timestamp * 1000;
        // forge-lint: disable-next-line(block-timestamp)
        if (nowMs + REVEAL_MS < openMs || nowMs >= openMs) revert WrongPhase();
    }

    function _takeCommitment(bytes32 symbol, uint64 openMs, bytes32 commitment) private returns (uint256 deposit) {
        deposit = commitments[symbol][openMs][commitment].deposit;
        if (deposit == 0) revert UnknownCommitment(commitment);
        delete commitments[symbol][openMs][commitment];
        _rounds[symbol][openMs].deposits -= SafeCast.toUint128(deposit);
    }

    function _checkClearWindow(uint64 openMs) private view {
        // slither-disable-next-line unused-return
        (, uint8 nyse,,,) = band.session(); // forge-lint: disable-line(unused-return)
        uint256 nowMs = block.timestamp * 1000;
        // forge-lint: disable-next-line(block-timestamp)
        if (nyse != 1 || nowMs < openMs || nowMs >= uint256(openMs) + CLEAR_MS) revert WrongPhase();
    }

    function _token(bytes32 symbol) private view returns (address) {
        (bytes32[] memory symbols, address[] memory tokens) = accounts.stocks();
        for (uint256 i; i < symbols.length; ++i) {
            if (symbols[i] == symbol && tokens[i] != address(0)) return tokens[i];
        }
        revert UnknownAsset(symbol);
    }

    function _demand(Bid[] storage all, uint256 price)
        private
        view
        returns (uint256 above, uint256 atPrice, bool below)
    {
        for (uint256 i; i < all.length; ++i) {
            uint256 p = all[i].price;
            if (p > price) above += all[i].quantity;
            else if (p == price) atPrice += all[i].quantity;
            else below = true;
        }
    }

    function _lowest(Bid[] storage all) private view returns (uint256 index) {
        for (uint256 i = 1; i < all.length; ++i) {
            Bid storage b = all[i];
            Bid storage low = all[index];
            if (b.price < low.price || (b.price == low.price && b.quantity < low.quantity)) index = i;
        }
    }

    function _smallest(Lot[] storage all) private view returns (uint256 index) {
        for (uint256 i = 1; i < all.length; ++i) {
            if (all[i].amount < all[index].amount) index = i;
        }
    }

    function _fill(Round memory r, Bid memory bid) private pure returns (uint256) {
        uint256 allocated = Math.min(r.supply, uint256(r.above) + r.atPrice);
        if (allocated == 0 || bid.price < r.price) return 0;
        if (bid.price > r.price) return Math.mulDiv(bid.quantity, r.sold, allocated);
        return Math.mulDiv(uint256(bid.quantity) * (allocated - r.above), r.sold, uint256(r.atPrice) * allocated);
    }
}
