// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IMargin} from "./interfaces/IMargin.sol";
import {IUSDG} from "./interfaces/IUSDG.sol";

/// @title Tapehouse weekend gap cover
/// @notice Cover against a Stock Token's fall over a market closure, paid in USDG. A cover pays its notional times the
/// fall from the asset's last Chainlink round before the close to its price at the reopening, beyond a deductible and
/// up to a limit. It is priced from the margin engine's weekend gap of the asset, moved to follow the asset's realised
/// move over the week before the close on the same feed, and sold only before the close.
/// Writers deposit USDG and receive shares; every cover reserves its whole possible payout from their USDG until it is
/// released, and its premium joins their USDG then. The reopening price is the first Chainlink round that starts at
/// or after the reopen when it lands within the settlement window inside the band, and otherwise the median of the
/// band's centre over the window, which anyone records a minute at a time, once a quorum of its minutes holds it.
/// The window is the fifteen minutes from the reopen, moved only by minutes in which the band could not vouch for
/// the asset; a series that cannot settle so is void and refunds its premiums.
/// @dev The cover counts the USDG it holds itself, so USDG sent to it directly changes no share's value.
contract GapCover is ERC4626 {
    using SafeERC20 for IERC20;

    struct Series {
        uint128 notional;
        uint64 referencePrice;
        uint64 price;
        uint16 shift;
        uint8 status;
        bool flagged;
    }

    struct Cover {
        address holder;
        uint64 closesMs;
        uint16 deductibleBps;
        uint16 limitBps;
        bytes32 symbol;
        uint128 notional;
        uint128 premium;
    }

    struct Observation {
        uint64 mid;
        uint64 low;
        uint128 high;
    }

    /// @notice The share of the expected payout a premium adds, in basis points.
    uint256 public constant LOADING_BPS = 5_000;
    /// @notice The fall over the closure, in millionths of the engine's weekend gap, past which the fitted tail starts:
    /// the smallest deductible.
    uint256 public constant TAIL_THRESHOLD_PPM = 183_288;
    /// @notice The fitted tail's scale, in millionths of the engine's weekend gap.
    uint256 public constant TAIL_SCALE_PPM = 95_837;
    /// @notice How often a weekend falls past the tail's threshold, in millionths.
    uint256 public constant TAIL_PROBABILITY_PPM = 50_224;
    /// @notice The weekend gap the cover prices at is at least this many basis points of the realised move of the
    /// trading week before the close, as the asset's Chainlink feed shows it: twice that move.
    uint256 public constant MOVE_TO_GAP_BPS = 20_000;
    /// @notice The least weekend gap the cover prices at, in basis points of the engine's: after a calm week it prices
    /// at 55% of the engine's gap.
    uint256 public constant MIN_GAP_BPS = 5_500;
    /// @notice The most the cover raises the engine's weekend gap to, as a multiple of it.
    uint256 public constant MAX_GAP_MULTIPLE = 8;
    /// @notice How many daily changes of the feed's answers make the trading week before a close: from the answers at
    /// the eight days that end a day before it.
    uint256 public constant WEEK_DAYS = 7;
    /// @notice How many one-minute slots the settlement window holds.
    uint256 public constant WINDOW_SLOTS = 15;
    /// @notice How many of the window's slots must hold the band for its median centre to settle a series.
    uint256 public constant QUORUM_SLOTS = 8;
    /// @notice How many slots a band that cannot vouch for the asset may move the window by; a series whose band
    /// cannot vouch for it longer after the reopen is void.
    uint256 public constant MAX_DELAY_SLOTS = 15;
    /// @notice How long a slot of the settlement window is, in milliseconds.
    uint256 public constant SLOT_MS = 60_000;
    /// @notice How long after its close a series no one has settled may be voided, in seconds.
    uint256 public constant SETTLEMENT_DEADLINE = 7 days;
    /// @notice From a regular close to the close of Chainlink's 24/5 session, in milliseconds.
    uint64 public constant POST_MARKET_MS = 14_400_000;
    /// @notice From the reopen of Chainlink's 24/5 session to the regular open, in milliseconds.
    uint64 public constant REOPEN_BEFORE_OPEN_MS = 48_600_000;

    uint8 private constant SETTLED = 1;
    uint8 private constant VOID = 2;
    uint256 private constant BPS = 10_000;
    uint256 private constant PPM = 1_000_000;
    uint256 private constant MS = 1000;
    uint256 private constant BPS_TO_PICO = 1e8;
    uint256 private constant PICO = 1e12;
    uint256 private constant DEVIATION_BPS = 50;
    uint256 private constant DAY_MS = 86_400_000;
    uint256 private constant UNREADABLE = type(uint64).max;

    /// @notice The margin engine whose weekend gap prices the cover.
    IMargin public immutable engine;
    /// @notice The engine's band, whose session times the sales and whose quotes check the reopening price.
    IBand public immutable band;
    /// @notice Whether the band's Chainlink feeds follow NYSE's regular hours: the sales then end at the regular close
    /// and the reopening is the regular open, rather than the close and reopen of the 24/5 session.
    bool public immutable regularHours;
    IUSDG private immutable _usdg;

    /// @notice The USDG the writers' shares stand for, by the cover's own count, reserved or not.
    uint256 public held;
    /// @notice The USDG reserved for the covers not yet released.
    uint256 public reserved;
    /// @notice The premiums of the covers not yet released.
    uint256 public premiums;
    /// @notice The payouts and refunds credited to holders and not yet claimed.
    uint256 public owed;
    /// @notice How many covers are not yet released.
    uint256 public outstanding;
    /// @notice The close that keys the series of the coming or the last recorded closure, in milliseconds: the first
    /// close of the 24/5 session the cover recorded for it.
    uint64 public lastCloseMs;
    /// @notice How many covers sold over the closure keyed by `closesMs` are not yet released.
    mapping(uint64 closesMs => uint256) public outstandingIn;
    /// @notice The cover sold on `symbol` over the closure from `closesMs`, and how it settled.
    mapping(bytes32 symbol => mapping(uint64 closesMs => Series)) public series;
    /// @notice Each cover: its holder, closure, layer, asset, notional and premium. Zero once released.
    mapping(uint256 id => Cover) public covers;
    /// @notice How many covers have been bought.
    uint256 public coverCount;
    /// @notice The USDG credited to `holder` and not yet claimed.
    mapping(address holder => uint256) public payouts;
    mapping(bytes32 symbol => AggregatorV3Interface) private _feeds;
    mapping(uint64 closesMs => uint64) private _closes;
    mapping(uint64 closesMs => uint64) private _reopens;
    mapping(bytes32 symbol => mapping(uint64 closesMs => mapping(uint256 slot => Observation))) private _observations;
    mapping(bytes32 symbol => mapping(uint64 closesMs => uint256)) private _weekMoves;

    /// @notice The cover recorded that the closure keyed by `closesMs` closes at `atMs`: first, or as the session
    /// revised it.
    event CloseRecorded(uint64 indexed closesMs, uint64 atMs);
    /// @notice The cover recorded that the closure keyed by `closesMs` reopens at `reopensMs`.
    event ReopenRecorded(uint64 indexed closesMs, uint64 reopensMs);
    /// @notice The series of `symbol` over the closure from `closesMs` keeps `weekMove`, the realised move of the week
    /// before its close in millionths, for every cover it sells.
    event Measured(bytes32 indexed symbol, uint64 indexed closesMs, uint256 weekMove);
    /// @notice `holder` holds cover `id` on `symbol` over the closure from `closesMs`, paying the fall beyond
    /// `deductibleBps` up to `limitBps` on `notional`, bought for `premium`.
    event Bought(
        uint256 indexed id,
        address indexed holder,
        bytes32 indexed symbol,
        uint64 closesMs,
        uint256 notional,
        uint256 deductibleBps,
        uint256 limitBps,
        uint256 premium
    );
    /// @notice The band of `symbol` in `slot`, counted from the reopen, after the closure from `closesMs`; zeros when the
    /// band could not vouch for the asset in the window's first slot, which moves the window by that slot.
    event Observed(bytes32 indexed symbol, uint64 indexed closesMs, uint256 slot, uint64 mid, uint64 low, uint128 high);
    /// @notice The series of `symbol` over the closure from `closesMs` settled at `price` against `referencePrice`,
    /// the answer of `referenceRound`: at the answer of `firstRound`, or, when `flagged`, at the band's median centre.
    event Settled(
        bytes32 indexed symbol,
        uint64 indexed closesMs,
        uint80 referenceRound,
        uint256 referencePrice,
        uint80 firstRound,
        uint256 price,
        bool flagged
    );
    /// @notice The series of `symbol` over the closure from `closesMs` is void: it could not settle on a price the
    /// band vouched for, or no one settled it in time.
    event Voided(bytes32 indexed symbol, uint64 indexed closesMs);
    /// @notice Cover `id` was released: `payout` and `refund` were credited to `holder`.
    event Released(uint256 indexed id, address indexed holder, uint256 payout, uint256 refund);
    /// @notice `holder` claimed `amount` to `receiver`.
    event Claimed(address indexed holder, address indexed receiver, uint256 amount);
    /// @notice The cover's count of the writers' USDG fell to `held`.
    event Sync(uint256 held);

    /// @notice The engine has no Chainlink feed in the band for `symbol`.
    error NotCoverable(bytes32 symbol);
    /// @notice No cover is sold now: the band's session shows no close ahead, or the sales have ended.
    error SalesClosed();
    /// @notice The amount is zero.
    error ZeroAmount();
    /// @notice The holder is the zero address.
    error InvalidHolder();
    /// @notice The band is halted for `symbol`.
    error AssetHalted(bytes32 symbol);
    /// @notice A multiplier change of `symbol`'s Stock Token is scheduled or not yet confirmed.
    error CorporateActionPending(bytes32 symbol);
    /// @notice The L2 sequencer is down or came back an hour ago or less.
    error SequencerNotSettled();
    /// @notice The deductible must be at least `minimum` and below the limit, and the limit at most 10,000 bps.
    error InvalidLayer(uint256 deductibleBps, uint256 limitBps, uint256 minimum);
    /// @notice The feed's last answer stands above the band's centre, so the deductible, measured from the centre, must
    /// be at least `minimum`.
    error StaleReference(uint256 deductibleBps, uint256 minimum);
    /// @notice The premium is above the buyer's limit.
    error PremiumAboveLimit(uint256 premium, uint256 maxPremium);
    /// @notice The cover would reserve `need` of the writers' USDG, more than the `free` USDG.
    error NoCapacity(uint256 need, uint256 free);
    /// @notice No cover was sold on `symbol` over the closure from `closesMs`.
    error NothingCovered(bytes32 symbol, uint64 closesMs);
    /// @notice The series has settled or is void.
    error SeriesClosed(bytes32 symbol, uint64 closesMs);
    /// @notice The reopen after the closure from `closesMs` is not recorded or not yet reached.
    error NotReopened(uint64 closesMs);
    /// @notice The settlement window of the series ended at `endMs`.
    error WindowClosed(uint256 endMs);
    /// @notice The settlement window of the series ends at `endMs`.
    error WindowOpen(uint256 endMs);
    /// @notice `roundId` is not the feed's last round started before the close.
    error InvalidReference(uint80 roundId);
    /// @notice `roundId` is not the feed's last round started before the reopen.
    error InvalidLastRound(uint80 roundId);
    /// @notice A series may be voided only from `fromMs`.
    error TooEarlyToVoid(uint256 fromMs);
    /// @notice The week before the close may be measured only from `fromMs`, once its last day has passed.
    error TooEarlyToMeasure(uint256 fromMs);
    /// @notice A read of a Chainlink feed ran out of gas: the call needs more.
    error InsufficientGas();
    /// @notice Cover `id` does not exist or has been released.
    error NoCover(uint256 id);
    /// @notice The series of cover `id` has not settled.
    error NotSettled(uint256 id);

    /// @param engine_ The margin engine, whose band and assets the cover takes.
    /// @param usdg_ USDG, which the writers deposit and the covers are paid in.
    constructor(IMargin engine_, IUSDG usdg_) ERC20("Tapehouse Gap Cover", "thCOVER") ERC4626(usdg_) {
        engine = engine_;
        band = IBand(engine_.band());
        _usdg = usdg_;
        // slither-disable-next-line unused-return
        (, regularHours) = band.chainConfig(); // forge-lint: disable-line(unused-return)
        bytes32[] memory symbols = engine_.assets();
        for (uint256 i; i < symbols.length; ++i) {
            // slither-disable-next-line unused-return,calls-loop
            (address chainlinkFeed,,,) = band.asset(symbols[i]); // forge-lint: disable-line(calls-loop, unused-return)
            _feeds[symbols[i]] = AggregatorV3Interface(chainlinkFeed);
        }
    }

    /// @notice Records the session's next close while it is open, and the reopen of the last recorded closure while the
    /// session is closed. A closure keeps the close first recorded for it as its key when the session revises the
    /// close before it is reached. Anyone may call it, and buying calls it first: a series settles only once its reopen
    /// has been recorded while the session was closed.
    function record() external {
        _record();
    }

    /// @notice Buys cover on `symbol` for `holder` over the coming closure: it pays `notional` times the fall from the
    /// last Chainlink round before the close to the reopening price beyond `deductibleBps`, up to `limitBps`. The
    /// caller pays the premium, at most `maxPremium`. Sold while the band's session shows a close ahead, until the sales
    /// end, while the band vouches for the asset, at the weekend gap of `pricingGap`. Where the feed's last answer
    /// stands above the band's centre, the layer is priced as measured from the centre, and its deductible so measured
    /// must reach the fitted tail. The series' first cover sold once the week's last day has passed keeps the realised
    /// move of the week before its close for the series' later covers, unless `measure` kept it first.
    function buy(
        bytes32 symbol,
        uint256 notional,
        uint256 deductibleBps,
        uint256 limitBps,
        uint256 maxPremium,
        address holder
    ) external returns (uint256 id, uint256 premium) {
        uint64 closesMs = _salesClose();
        if (notional == 0) revert ZeroAmount();
        if (holder == address(0)) revert InvalidHolder();
        uint256 weekMove;
        (premium, weekMove) = _vouched(symbol, closesMs, notional, deductibleBps, limitBps);
        if (premium > maxPremium) revert PremiumAboveLimit(premium, maxPremium);
        uint256 need = _reserve(notional, deductibleBps, limitBps);
        uint256 free = capacity();
        if (need > free) revert NoCapacity(need, free);
        reserved += need;
        premiums += premium;
        ++outstanding;
        ++outstandingIn[closesMs];
        series[symbol][closesMs].notional += SafeCast.toUint128(notional);
        // forge-lint: disable-next-line(block-timestamp)
        if (_weekMoves[symbol][closesMs] == 0 && _nowMs() >= closesMs - DAY_MS) _keep(symbol, closesMs, weekMove);
        id = ++coverCount;
        covers[id] = Cover(
            holder,
            closesMs,
            SafeCast.toUint16(deductibleBps),
            SafeCast.toUint16(limitBps),
            symbol,
            SafeCast.toUint128(notional),
            SafeCast.toUint128(premium)
        );
        // forge-lint: disable-next-line(reentrancy-events)
        emit Bought(id, holder, symbol, closesMs, notional, deductibleBps, limitBps, premium);
        IERC20(asset()).safeTransferFrom(msg.sender, address(this), premium);
    }

    /// @notice Keeps the realised move of the week before the close on sale for `symbol`'s series, as its first cover
    /// would, so that a keeper rather than the first buyer pays for reading it from the feed; returns it. Only during
    /// the sales, once the week's last day, a day before the close, has passed. Anyone may call it.
    function measure(bytes32 symbol) external returns (uint256 weekMove) {
        uint64 closesMs = _salesClose();
        AggregatorV3Interface feed_ = _feeds[symbol];
        if (address(feed_) == address(0)) revert NotCoverable(symbol);
        uint256 fromMs = closesMs - DAY_MS;
        // forge-lint: disable-next-line(block-timestamp)
        if (_nowMs() < fromMs) revert TooEarlyToMeasure(fromMs);
        weekMove = _weekMoves[symbol][closesMs];
        if (weekMove != 0) return weekMove - 1;
        // slither-disable-next-line unused-return
        (uint80 latest,,,,) = feed_.latestRoundData(); // forge-lint: disable-line(unused-return)
        weekMove = _weekMove(feed_, latest, closesMs);
        _keep(symbol, closesMs, weekMove);
    }

    /// @notice Records the band of `symbol` in this minute's slot of the settlement window after the closure from
    /// `closesMs`, once the reopen is reached. The window holds the `WINDOW_SLOTS` one-minute slots from the reopen,
    /// moved by one slot each time the band cannot vouch for the asset (halted or closed, or the sequencer not
    /// settled) in its first slot, as recorded here; a band that cannot vouch for it for more than `MAX_DELAY_SLOTS`
    /// slots closes the window, and the series is void. Returns whether it recorded the slot, which it does once.
    /// Anyone may call it.
    function observe(bytes32 symbol, uint64 closesMs) external returns (bool) {
        Series storage s = _open(symbol, closesMs);
        uint256 reopenMs = reopenOf(closesMs);
        uint256 nowMs = _nowMs();
        // forge-lint: disable-next-line(block-timestamp)
        if (reopenMs == 0 || nowMs < reopenMs) revert NotReopened(closesMs);
        uint256 slot = (nowMs - reopenMs) / SLOT_MS;
        uint256 shift = s.shift;
        uint256 end = _end(shift);
        // forge-lint: disable-next-line(block-timestamp)
        if (slot >= end) revert WindowClosed(reopenMs + end * SLOT_MS);
        // forge-lint: disable-next-line(block-timestamp)
        if (slot < shift || _observations[symbol][closesMs][slot].mid != 0) return false;
        if (_recordBand(symbol, closesMs, slot)) return true;
        // forge-lint: disable-next-line(block-timestamp)
        if (slot != shift) return false;
        s.shift = SafeCast.toUint16(shift + 1);
        emit Observed(symbol, closesMs, slot, 0, 0, 0);
        return true;
    }

    /// @notice Settles the series of `symbol` over the closure from `closesMs` once its window has ended.
    /// `referenceRound` must be the feed's last round started before the sales ended, and `lastRound` its last round
    /// started before the reopen. The round after `lastRound` settles the series when it lands within the window,
    /// inside the band recorded in its slot; otherwise the median of the band centres recorded in the window does,
    /// and the series is flagged. The series is void instead when the window never opened, when fewer than
    /// `QUORUM_SLOTS` of its slots hold the band and the round does not settle it, or when a material multiplier
    /// change took effect after the reference round started and before the window's end. Anyone may call it.
    function settle(bytes32 symbol, uint64 closesMs, uint80 referenceRound, uint80 lastRound) external {
        Series storage s = _open(symbol, closesMs);
        uint256 reopenMs = reopenOf(closesMs);
        if (reopenMs == 0) revert NotReopened(closesMs);
        uint256 shift = s.shift;
        uint256 endMs = reopenMs + _end(shift) * SLOT_MS;
        // forge-lint: disable-next-line(block-timestamp)
        if (_nowMs() < endMs) revert WindowOpen(endMs);
        if (shift > MAX_DELAY_SLOTS) return _void(s, symbol, closesMs);
        AggregatorV3Interface feed_ = _feeds[symbol];
        (uint256 referencePrice, uint256 referenceStartedAt) =
            _reference(feed_, referenceRound, _salesEndMs(_closeOf(closesMs)));
        if (_stepped(symbol, referenceStartedAt, endMs)) return _void(s, symbol, closesMs);
        (uint80 firstRound, uint256 price) = _anchor(feed_, symbol, closesMs, reopenMs, lastRound);
        bool flagged = firstRound == 0;
        if (flagged) price = _median(symbol, closesMs, shift);
        if (price == 0) return _void(s, symbol, closesMs);
        s.referencePrice = SafeCast.toUint64(referencePrice);
        s.price = SafeCast.toUint64(price);
        s.flagged = flagged;
        s.status = SETTLED;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Settled(symbol, closesMs, referenceRound, referencePrice, firstRound, price, flagged);
    }

    /// @notice Voids the series of `symbol` over the closure from `closesMs` when no one has settled it within
    /// `SETTLEMENT_DEADLINE` of its close: its covers refund their premiums. Anyone may call it.
    function void(bytes32 symbol, uint64 closesMs) external {
        Series storage s = _open(symbol, closesMs);
        uint256 fromMs = uint256(closesMs) + SETTLEMENT_DEADLINE * MS;
        // forge-lint: disable-next-line(block-timestamp)
        if (_nowMs() < fromMs) revert TooEarlyToVoid(fromMs);
        _void(s, symbol, closesMs);
    }

    /// @notice Releases cover `id` once its series has settled or is void: credits its holder the payout, or the
    /// premium of a void series, frees its reserve, and adds its premium less its payout to the writers' USDG. Anyone
    /// may call it.
    function release(uint256 id) external returns (uint256 payout, uint256 refund) {
        Cover memory c = covers[id];
        if (c.holder == address(0)) revert NoCover(id);
        Series memory s = series[c.symbol][c.closesMs];
        if (s.status == 0) revert NotSettled(id);
        delete covers[id];
        reserved -= _reserve(c.notional, c.deductibleBps, c.limitBps);
        premiums -= c.premium;
        --outstanding;
        --outstandingIn[c.closesMs];
        if (s.status == VOID) {
            refund = c.premium;
        } else {
            uint256 pool = held + c.premium;
            payout = Math.min(_payout(c.notional, c.deductibleBps, c.limitBps, s.referencePrice, s.price), pool);
            held = pool - payout;
        }
        payouts[c.holder] += payout + refund;
        owed += payout + refund;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Released(id, c.holder, payout, refund);
    }

    /// @notice Sends `receiver` the USDG credited to the caller.
    function claim(address receiver) external returns (uint256 amount) {
        amount = payouts[msg.sender];
        payouts[msg.sender] = 0;
        owed -= amount;
        emit Claimed(msg.sender, receiver, amount);
        if (amount != 0) IERC20(asset()).safeTransfer(receiver, amount);
    }

    /// @notice Brings the writers' count down to what the cover's USDG leaves after the premiums and the credits it
    /// holds, as after the issuer wipes a frozen cover: a loss to the shares. Anyone may call it.
    function sync() external {
        uint256 balance = _usdg.balanceOf(address(this));
        uint256 others = premiums + owed;
        uint256 left = balance > others ? balance - others : 0;
        if (left < held) {
            held = left;
            emit Sync(left);
        }
    }

    /// @notice The premium of cover on `symbol` paying the fall beyond `deductibleBps` up to `limitBps` on `notional`,
    /// at the weekend gap of `pricingGap` and the band's centre now.
    function quote(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps)
        external
        view
        returns (uint256)
    {
        (uint64 closesMs,) = _onSale();
        (uint256 gap,, uint256 stale,) = _terms(symbol, closesMs);
        return _premium(gap, stale, notional, deductibleBps, limitBps);
    }

    /// @notice The smallest deductible of cover on `symbol` now, in basis points: the start of the fitted tail at the
    /// weekend gap of `pricingGap`, plus how far the feed's last answer stands above the band's centre.
    function minDeductible(bytes32 symbol) external view returns (uint256) {
        (uint64 closesMs,) = _onSale();
        (uint256 gap,, uint256 stale,) = _terms(symbol, closesMs);
        return Math.ceilDiv(gap * TAIL_THRESHOLD_PPM + stale, BPS_TO_PICO);
    }

    /// @notice The weekend gap the cover prices `symbol` at now, in millionths, and the realised move of the trading
    /// week before the close on sale that it follows. During the sales the gap is `MOVE_TO_GAP_BPS` of that move, at
    /// least `MIN_GAP_BPS` of the engine's and at most `MAX_GAP_MULTIPLE` times it. The move is the root of the sum of
    /// the squared daily changes, in millionths, of the asset's Chainlink answers at the eight days that end a day
    /// before the close: each the answer of the last round started before it, in the phase of the feed's latest round
    /// or, where that phase has none, in the latest phase before it that has. The week is fixed by the close and the
    /// feed alone, so no caller chooses its rounds, and the series keeps it once its last day has passed. A week the
    /// feed cannot show, a day without a round or a positive answer, reads as `type(uint64).max` and prices at the
    /// most. Outside the sales the gap is the engine's and the move zero.
    function pricingGap(bytes32 symbol) external view returns (uint256 gap, uint256 weekMove) {
        (uint64 closesMs,) = _onSale();
        (gap, weekMove,,) = _terms(symbol, closesMs);
    }

    /// @notice The closure on sale now: the close that keys its series and when its sales end, in milliseconds; zeros
    /// while no cover is sold.
    function sales() external view returns (uint64 closesMs, uint64 endsMs) {
        return _onSale();
    }

    /// @notice The Chainlink feed that settles cover on `symbol`; zero where the asset is not coverable.
    function feed(bytes32 symbol) external view returns (address) {
        return address(_feeds[symbol]);
    }

    /// @notice The writers' USDG not reserved for a cover.
    function capacity() public view returns (uint256) {
        return held > reserved ? held - reserved : 0;
    }

    /// @notice When the closure keyed by `closesMs` reopens, in milliseconds: the recorded reopen of the 24/5 session,
    /// or the regular open after it where the feeds follow regular hours. Zero while not recorded.
    function reopenOf(uint64 closesMs) public view returns (uint256) {
        uint64 reopenMs = _reopens[closesMs];
        if (reopenMs == 0) return 0;
        return regularHours ? uint256(reopenMs) + REOPEN_BEFORE_OPEN_MS : reopenMs;
    }

    /// @inheritdoc ERC4626
    function totalAssets() public view override returns (uint256) {
        return held;
    }

    /// @inheritdoc ERC4626
    /// @dev Zero while USDG is paused, the cover frozen or its USDG wiped, and while a cover is outstanding over a
    /// closure other than the coming one, or the sales for the coming one have ended.
    function maxDeposit(address receiver) public view override returns (uint256) {
        return _depositsOpen() ? super.maxDeposit(receiver) : 0;
    }

    /// @inheritdoc ERC4626
    /// @dev Zero while USDG is paused, the cover frozen or its USDG wiped, and while a cover is outstanding over a
    /// closure other than the coming one, or the sales for the coming one have ended.
    function maxMint(address receiver) public view override returns (uint256) {
        return _depositsOpen() ? super.maxMint(receiver) : 0;
    }

    /// @inheritdoc ERC4626
    /// @dev Zero while a cover is outstanding, and while USDG is paused, the cover frozen or its USDG wiped.
    function maxRedeem(address owner) public view override returns (uint256) {
        return outstanding != 0 || _unavailable() ? 0 : super.maxRedeem(owner);
    }

    function _transferIn(address from, uint256 assets) internal override {
        held += assets; // forge-lint: disable-line(missing-events-arithmetic)
        super._transferIn(from, assets);
    }

    function _transferOut(address to, uint256 assets) internal override {
        held -= assets; // forge-lint: disable-line(missing-events-arithmetic)
        super._transferOut(to, assets);
    }

    function _decimalsOffset() internal pure override returns (uint8) {
        return 6;
    }

    function _record() private returns (uint8 session, uint64 boundaryMs, uint64 closesMs) {
        // slither-disable-next-line unused-return
        (session,,,, boundaryMs) = band.session(); // forge-lint: disable-line(unused-return)
        uint64 nowMs = SafeCast.toUint64(_nowMs());
        closesMs = lastCloseMs;
        if (session == 2 && boundaryMs > nowMs) {
            uint64 key = _keyOf(boundaryMs, nowMs);
            if (key != closesMs) {
                lastCloseMs = key;
                emit CloseRecorded(key, boundaryMs);
            } else if (_closeOf(key) != boundaryMs) {
                _closes[key] = boundaryMs;
                emit CloseRecorded(key, boundaryMs);
            }
            closesMs = key;
        } else if (session == 1 && boundaryMs != 0 && closesMs != 0 && _closeOf(closesMs) <= nowMs) {
            uint64 reopenMs = _reopens[closesMs];
            if (boundaryMs != reopenMs && (reopenMs == 0 || nowMs < reopenMs)) {
                _reopens[closesMs] = boundaryMs;
                emit ReopenRecorded(closesMs, boundaryMs);
            }
        }
    }

    function _salesClose() private returns (uint64 closesMs) {
        (uint8 session, uint64 boundaryMs, uint64 key) = _record();
        // forge-lint: disable-next-line(block-timestamp)
        if (session != 2 || _nowMs() >= _salesEndMs(boundaryMs)) revert SalesClosed();
        return key;
    }

    /// @dev The premium of the layer over the closure keyed by `closesMs` while the band vouches for the asset, and
    /// the realised move of the week before its close.
    function _vouched(bytes32 symbol, uint64 closesMs, uint256 notional, uint256 deductibleBps, uint256 limitBps)
        private
        view
        returns (uint256, uint256)
    {
        (uint256 gap, uint256 weekMove, uint256 stale, uint8 state) = _terms(symbol, closesMs);
        if (state == 0) revert AssetHalted(symbol);
        // slither-disable-next-line unused-return
        (uint8 action,,,) = band.corporateAction(symbol); // forge-lint: disable-line(unused-return)
        if (action != 0) revert CorporateActionPending(symbol);
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        return (_premium(gap, stale, notional, deductibleBps, limitBps), weekMove);
    }

    /// @dev The weekend gap the cover prices at over the closure keyed by `closesMs`, none for zero, the realised move
    /// of the week before its close, how far the feed's last answer stands above the band's centre in picounits of
    /// the price, rounded up, and the band's state. A halted band has no centre, and nothing stands above it.
    function _terms(bytes32 symbol, uint64 closesMs)
        private
        view
        returns (uint256 gap, uint256 weekMove, uint256 stale, uint8 state)
    {
        gap = _gap(symbol);
        uint64 mid;
        // slither-disable-next-line unused-return
        (state,, mid,,,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        AggregatorV3Interface feed_ = _feeds[symbol];
        // slither-disable-next-line unused-return
        (uint80 latest, int256 answer,,,) = feed_.latestRoundData(); // forge-lint: disable-line(unused-return)
        uint256 last = SafeCast.toUint256(answer);
        if (mid != 0 && mid < last) stale = Math.mulDiv(last - mid, PICO, last, Math.Rounding.Ceil);
        if (closesMs == 0) return (gap, 0, stale, state);
        weekMove = _weekMoves[symbol][closesMs];
        weekMove = weekMove != 0 ? weekMove - 1 : _weekMove(feed_, latest, closesMs);
        gap = Math.min(Math.max(gap * MIN_GAP_BPS / BPS, weekMove * MOVE_TO_GAP_BPS / BPS), gap * MAX_GAP_MULTIPLE);
    }

    function _keep(bytes32 symbol, uint64 closesMs, uint256 weekMove) private {
        _weekMoves[symbol][closesMs] = weekMove + 1;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Measured(symbol, closesMs, weekMove);
    }

    /// @dev The realised move of the week before the close `closesMs`: the root of the sum of the squared daily
    /// changes of the feed's answers at the `WEEK_DAYS` + 1 days that end a day before the close, each in millionths of
    /// the earlier answer and at most `UNREADABLE`, found from round `latest` down through the feed's phases.
    /// `UNREADABLE` where a day has no round or no positive answer.
    function _weekMove(AggregatorV3Interface feed_, uint80 latest, uint64 closesMs) private view returns (uint256) {
        // forge-lint: disable-next-line(unsafe-typecast)
        uint256 below = uint64(latest) + 1;
        (uint256 phase, uint256 index, int256 answer) = _lastBefore(feed_, latest >> 64, below, closesMs - DAY_MS);
        if (answer <= 0) return UNREADABLE;
        uint256 later = SafeCast.toUint256(answer);
        uint256 sum = 0;
        for (uint256 day = 2; day <= WEEK_DAYS + 1; ++day) {
            (phase, index, answer) = _lastBefore(feed_, phase, index + 1, closesMs - day * DAY_MS);
            if (answer <= 0) return UNREADABLE;
            uint256 price = SafeCast.toUint256(answer);
            uint256 change = Math.saturatingMul(later > price ? later - price : price - later, PPM) / price;
            change = Math.min(change, UNREADABLE);
            sum += change * change;
            later = price;
        }
        return Math.min(Math.sqrt(sum), UNREADABLE);
    }

    /// @dev The phase and index of the feed's last round started before `atMs`, and its answer: among the rounds of
    /// `phase` below `below` or, where none started before it, of the latest earlier phase that has one; zeros where
    /// there is none. Gallops down from `below`, or up from an earlier phase's first round, and then bisects, so a day
    /// costs a few reads.
    function _lastBefore(AggregatorV3Interface feed_, uint256 phase, uint256 below, uint256 atMs)
        private
        view
        returns (uint256, uint256 index, int256 answer)
    {
        uint256 hi = below;
        uint256 step = 1;
        while (index == 0) {
            uint256 i = hi > step ? hi - step : 1;
            (bool before, int256 a) = _startedBefore(feed_, phase, i, atMs);
            if (before) (index, answer) = (i, a);
            else if (i != 1) (hi, step) = (i, step * 2);
            else if (phase < 2) return (0, 0, 0);
            else (index, answer, hi) = _climb(feed_, --phase, atMs);
        }
        while (hi - index > 1) {
            uint256 i = (index + hi) / 2;
            (bool before, int256 a) = _startedBefore(feed_, phase, i, atMs);
            if (before) (index, answer) = (i, a);
            else hi = i;
        }
        return (phase, index, answer);
    }

    /// @dev The last of the rounds 1, 2, 4, 8… of `phase` started before `atMs`, its answer and the next of them;
    /// a zero index where round 1 did not.
    function _climb(AggregatorV3Interface feed_, uint256 phase, uint256 atMs)
        private
        view
        returns (uint256 index, int256 answer, uint256 hi)
    {
        for (hi = 1;; hi *= 2) {
            (bool before, int256 a) = _startedBefore(feed_, phase, hi, atMs);
            if (!before) return (index, answer, hi);
            (index, answer) = (hi, a);
        }
    }

    function _startedBefore(AggregatorV3Interface feed_, uint256 phase, uint256 index, uint256 atMs)
        private
        view
        returns (bool, int256)
    {
        (bool exists, int256 answer, uint256 startedAt,) = _round(feed_, SafeCast.toUint80(phase << 64 | index));
        return (exists && startedAt * MS < atMs, answer);
    }

    /// @dev The closure on sale now, by its key, and when its sales end; zeros while no cover is sold.
    function _onSale() private view returns (uint64 closesMs, uint64 endsMs) {
        // slither-disable-next-line unused-return
        (uint8 session,,,, uint64 boundaryMs) = band.session(); // forge-lint: disable-line(unused-return)
        uint256 nowMs = _nowMs();
        uint256 end = _salesEndMs(boundaryMs);
        // forge-lint: disable-next-line(block-timestamp)
        if (session != 2 || nowMs >= end) return (0, 0);
        return (_keyOf(boundaryMs, nowMs), SafeCast.toUint64(end));
    }

    /// @dev The key of the closure whose close is `boundaryMs`: the recorded one while its close is still ahead.
    function _keyOf(uint64 boundaryMs, uint256 nowMs) private view returns (uint64) {
        uint64 closesMs = lastCloseMs;
        // forge-lint: disable-next-line(block-timestamp)
        return _closeOf(closesMs) > nowMs ? closesMs : boundaryMs;
    }

    /// @dev The close of the closure keyed by `closesMs`: the key itself, unless the session revised it.
    function _closeOf(uint64 closesMs) private view returns (uint64) {
        uint64 revised = _closes[closesMs];
        return revised == 0 ? closesMs : revised;
    }

    function _open(bytes32 symbol, uint64 closesMs) private view returns (Series storage s) {
        s = series[symbol][closesMs];
        if (s.notional == 0) revert NothingCovered(symbol, closesMs);
        if (s.status != 0) revert SeriesClosed(symbol, closesMs);
    }

    /// @dev Records the band of `symbol` in `slot` when it vouches for the asset: degraded or open, with the
    /// sequencer settled.
    function _recordBand(bytes32 symbol, uint64 closesMs, uint256 slot) private returns (bool) {
        // slither-disable-next-line unused-return
        (uint8 state,, uint64 mid,, uint64 low, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if ((state != 1 && state != 3) || !band.sequencerSettled()) return false;
        _observations[symbol][closesMs][slot] = Observation(mid, low, high);
        emit Observed(symbol, closesMs, slot, mid, low, high);
        return true;
    }

    function _void(Series storage s, bytes32 symbol, uint64 closesMs) private {
        s.status = VOID;
        // forge-lint: disable-next-line(reentrancy-events)
        emit Voided(symbol, closesMs);
    }

    /// @dev The window's end, in slots from the reopen: none opens once the band has failed for more than
    /// `MAX_DELAY_SLOTS` slots.
    function _end(uint256 shift) private pure returns (uint256) {
        return shift > MAX_DELAY_SLOTS ? shift : shift + WINDOW_SLOTS;
    }

    function _reference(AggregatorV3Interface feed_, uint80 roundId, uint256 closeMs)
        private
        view
        returns (uint256, uint256)
    {
        (, int256 answer, uint256 startedAt,) = _round(feed_, roundId);
        if (answer <= 0 || startedAt * MS >= closeMs) revert InvalidReference(roundId);
        (, bool more,, uint256 nextStartedAt,) = _next(feed_, roundId);
        if (more && nextStartedAt * MS < closeMs) revert InvalidReference(roundId);
        return (SafeCast.toUint256(answer), startedAt);
    }

    /// @dev Whether the band's last multiplier change of `symbol` is material and took effect after the reference
    /// round started and before the window's end: the reference and the reopening price would then be in different
    /// terms. Material as the band counts it: of unknown size, or at least Chainlink's deviation threshold.
    function _stepped(bytes32 symbol, uint256 fromS, uint256 untilMs) private view returns (bool) {
        // slither-disable-next-line unused-return
        (, uint64 effectiveAt, uint128 before, uint128 after_) = band.corporateAction(symbol); // forge-lint: disable-line(unused-return)
        // forge-lint: disable-next-line(block-timestamp)
        if (effectiveAt <= fromS || effectiveAt * MS >= untilMs) return false;
        uint256 change = after_ > before ? after_ - before : before - after_;
        return change * BPS >= uint256(before) * DEVIATION_BPS;
    }

    function _anchor(AggregatorV3Interface feed_, bytes32 symbol, uint64 closesMs, uint256 reopenMs, uint80 lastRound)
        private
        view
        returns (uint80 firstRound, uint256 price)
    {
        (bool exists,, uint256 startedAt,) = _round(feed_, lastRound);
        if (!exists || startedAt * MS >= reopenMs) revert InvalidLastRound(lastRound);
        (uint80 nextId, bool more, int256 answer, uint256 nextStartedAt, uint256 updatedAt) = _next(feed_, lastRound);
        if (!more) return (0, 0);
        if (nextStartedAt * MS < reopenMs) revert InvalidLastRound(lastRound);
        if (answer <= 0) return (0, 0);
        price = SafeCast.toUint256(answer);
        if (_inBand(symbol, closesMs, reopenMs, price, updatedAt * MS)) firstRound = nextId;
    }

    /// @dev Only the window's slots hold an observation, and a slot without one has bounds of zero.
    function _inBand(bytes32 symbol, uint64 closesMs, uint256 reopenMs, uint256 price, uint256 landedMs)
        private
        view
        returns (bool)
    {
        if (landedMs < reopenMs) return false;
        Observation memory o = _observations[symbol][closesMs][(landedMs - reopenMs) / SLOT_MS];
        return price >= o.low && price <= o.high;
    }

    /// @dev The median of the band centres recorded in the window, or zero when fewer than `QUORUM_SLOTS` were.
    function _median(bytes32 symbol, uint64 closesMs, uint256 shift) private view returns (uint256) {
        uint256[] memory mids = new uint256[](WINDOW_SLOTS);
        uint256 n = 0;
        for (uint256 slot = shift; slot < shift + WINDOW_SLOTS; ++slot) {
            uint256 mid = _observations[symbol][closesMs][slot].mid;
            if (mid == 0) continue;
            uint256 i = n++;
            while (i != 0 && mids[i - 1] > mid) {
                mids[i] = mids[i - 1];
                --i;
            }
            mids[i] = mid;
        }
        if (n < QUORUM_SLOTS) return 0;
        return n % 2 == 1 ? mids[n / 2] : (mids[n / 2 - 1] + mids[n / 2]) / 2;
    }

    function _next(AggregatorV3Interface feed_, uint80 roundId)
        private
        view
        returns (uint80 nextId, bool exists, int256 answer, uint256 startedAt, uint256 updatedAt)
    {
        nextId = roundId + 1;
        (exists, answer, startedAt, updatedAt) = _round(feed_, nextId);
        if (exists) return (nextId, exists, answer, startedAt, updatedAt);
        nextId = ((roundId >> 64) + 1) << 64 | 1;
        (exists, answer, startedAt, updatedAt) = _round(feed_, nextId);
    }

    function _round(AggregatorV3Interface feed_, uint80 roundId)
        private
        view
        returns (bool exists, int256 answer, uint256 startedAt, uint256 updatedAt)
    {
        uint256 budget = gasleft();
        // slither-disable-next-line unused-return,calls-loop
        try feed_.getRoundData(roundId) returns (uint80, int256 a, uint256 started, uint256 updated, uint80) { // forge-lint: disable-line(calls-loop)
            return (updated != 0, a, started, updated);
        } catch {
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (gasleft() < budget / 64) revert InsufficientGas();
        }
    }

    function _gap(bytes32 symbol) private view returns (uint32 gap) {
        if (address(_feeds[symbol]) == address(0)) revert NotCoverable(symbol);
        // slither-disable-next-line unused-return
        (gap,) = engine.weekendGap(symbol); // forge-lint: disable-line(unused-return)
    }

    /// @dev The expected payout of the layer under a generalized Pareto tail of shape one half, threshold and scale in
    /// proportion to the weekend gap, with the loading on top, rounded up. In picounits of the notional, with `b` the
    /// scale, `u` the threshold and `x`, `l` the deductible and limit, the tail's stop-loss transform is
    /// `4·p·b² / (2b + x − u)`, so the layer's expected payout is `4·p·b²·(l − x) / ((2b + x − u)(2b + l − u))`.
    /// The deductible and the limit are measured from the band's centre: both lowered by `stale`.
    function _premium(uint256 gap, uint256 stale, uint256 notional, uint256 deductibleBps, uint256 limitBps)
        private
        pure
        returns (uint256)
    {
        uint256 u = gap * TAIL_THRESHOLD_PPM;
        uint256 minimum = Math.ceilDiv(u, BPS_TO_PICO);
        if (deductibleBps < minimum || limitBps <= deductibleBps || limitBps > BPS) {
            revert InvalidLayer(deductibleBps, limitBps, minimum);
        }
        uint256 x = deductibleBps * BPS_TO_PICO;
        if (x < u + stale) revert StaleReference(deductibleBps, Math.ceilDiv(u + stale, BPS_TO_PICO));
        x -= stale;
        uint256 l = limitBps * BPS_TO_PICO - stale;
        uint256 b = gap * TAIL_SCALE_PPM;
        uint256 numerator = 4 * TAIL_PROBABILITY_PPM * b * b * (l - x) * (BPS + LOADING_BPS);
        uint256 denominator = PPM * PICO * BPS * (2 * b + x - u) * (2 * b + l - u);
        return Math.mulDiv(notional, numerator, denominator, Math.Rounding.Ceil);
    }

    function _reserve(uint256 notional, uint256 deductibleBps, uint256 limitBps) private pure returns (uint256) {
        return Math.mulDiv(notional, limitBps - deductibleBps, BPS, Math.Rounding.Ceil);
    }

    function _payout(uint256 notional, uint256 deductibleBps, uint256 limitBps, uint256 referencePrice, uint256 price)
        private
        pure
        returns (uint256)
    {
        if (price >= referencePrice) return 0;
        uint256 fall = Math.min((referencePrice - price) * BPS, referencePrice * limitBps);
        uint256 deductible = referencePrice * deductibleBps;
        if (fall <= deductible) return 0;
        return Math.mulDiv(notional, fall - deductible, referencePrice * BPS);
    }

    function _salesEndMs(uint256 closeMs) private view returns (uint256) {
        return regularHours ? closeMs - Math.min(closeMs, POST_MARKET_MS) : closeMs;
    }

    function _nowMs() private view returns (uint256) {
        return block.timestamp * MS;
    }

    function _depositsOpen() private view returns (bool) {
        if (_unavailable()) return false;
        uint256 count = outstanding;
        if (count == 0) return true;
        // slither-disable-next-line unused-return
        try band.session() returns (uint8 session, uint8, uint8, uint64, uint64 boundaryMs) {
            uint256 nowMs = _nowMs();
            // forge-lint: disable-next-line(block-timestamp)
            if (session != 2 || nowMs >= _salesEndMs(boundaryMs)) return false;
            // forge-lint: disable-next-line(block-timestamp)
            return count == outstandingIn[_keyOf(boundaryMs, nowMs)];
        } catch {
            return false;
        }
    }

    function _unavailable() private view returns (bool) {
        try _usdg.paused() returns (bool paused) {
            if (paused) return true;
        } catch {
            return true;
        }
        try _usdg.isFrozen(address(this)) returns (bool frozen) {
            if (frozen) return true;
        } catch {
            return true;
        }
        try _usdg.balanceOf(address(this)) returns (uint256 balance) {
            return balance < held + premiums + owed;
        } catch {
            return true;
        }
    }
}
