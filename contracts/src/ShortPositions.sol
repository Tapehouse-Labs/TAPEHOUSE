// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IMargin} from "./interfaces/IMargin.sol";
import {IStockLendingBorrower} from "./interfaces/IStockLendingBorrower.sol";
import {IStockToken, IStockTokenRegistry} from "./interfaces/IStockToken.sol";
import {IUniswapV3Factory, IV3SwapRouter} from "./interfaces/IUniswapV3.sol";
import {MarginAccounts} from "./MarginAccounts.sol";
import {StockLendingVault} from "./StockLendingVault.sol";

/// @title Tapehouse short positions
/// @notice Each address may hold, for any Stock Token, a short position in it alone: it borrows the token from the
/// margin accounts' lending vault for that asset, sells it for USDG through its Uniswap v3 pool, and keeps the USDG,
/// with the margin it adds, as its only collateral, margined by the engine at the band's high edge. A short sells at no
/// less than the band's low edge, or its centre once the asset fell `RESTRICTION_DROP_BPS` from its last close; a
/// partial buy-back pays no more than the high edge, and a liquidation or a buy-in no more than `MAX_PREMIUM_BPS` above
/// the centre. The token's fee accrues in the vault, in the token. When the vault's lenders recall what the shorts
/// borrowed and a notice runs out, every short of the token buys its part back in proportion. Anyone may buy back a
/// short that falls short and take a bonus; a short whose USDG cannot buy back what it owes leaves the rest to the
/// vault's lenders. USDG in these positions never reaches the supply vault, and nothing here borrows from it. An address
/// authorized in the margin accounts may sell, buy back, deposit and withdraw for the account here too.
/// @dev The contract is each lending vault's one borrower. Each short owes its share of the vault's debt, kept as
/// borrow shares with Morpho Blue's virtual offsets in its token's book. A buy-in's cost is charged to every short of
/// the book by a cost index per borrow share, paid when each short is next touched; each book counts its shorts' USDG
/// apart, so a buy-in or a short in deficit never spends another book's. Once the vault is owed nothing while borrow
/// shares remain, the next sale opens a new book, an epoch, and the earlier shorts owe nothing and hold only their USDG.
/// Authorization and the guardian's pause are the margin accounts'.
contract ShortPositions is IStockLendingBorrower {
    using SafeERC20 for IERC20;

    struct Short {
        uint128 usdg;
        uint128 shares;
        uint192 costIndex;
        uint64 epoch;
    }

    struct Book {
        uint128 shares;
        uint128 usdg;
        uint256 costIndex;
    }

    struct Mark {
        uint32 day;
        uint64 last;
        uint64 close;
        uint32 until;
    }

    /// @notice What the caller of a liquidation takes from the short's USDG, in basis points of what the buy-back cost.
    uint256 public constant LIQUIDATION_BONUS_BPS = 50;
    /// @notice The fall from an asset's last close, in basis points, after which a short sells at no less than the
    /// band's centre for the rest of that day and the next: Regulation SHO Rule 201's 10%.
    uint256 public constant RESTRICTION_DROP_BPS = 10_00;
    /// @notice The most a liquidation or a buy-in pays for a token above the band's centre, in basis points, and never
    /// more than the band's high edge.
    uint256 public constant MAX_PREMIUM_BPS = 1_00;

    uint256 private constant BPS = 10_000;
    uint256 private constant PRICE_UNIT = 1e8;
    uint256 private constant USDG_TO_USD = 1e12;
    uint256 private constant TOKEN_TO_USDG = 1e20;
    uint256 private constant VIRTUAL_SHARES = 1e6;
    uint256 private constant VIRTUAL_ASSETS = 1;
    uint256 private constant RAY = 1e27;
    uint256 private constant MAX_SHARES_PER_UNIT = 1e12;

    /// @notice The margin accounts whose lending vaults the shorts borrow from, whose authorizations they honour and
    /// whose guardian pauses them.
    MarginAccounts public immutable accounts;
    /// @notice The band that prices every Stock Token.
    IBand public immutable band;
    /// @notice The engine that margins every short.
    IMargin public immutable engine;
    /// @notice The token the shorts sell for and keep.
    IERC20 public immutable usdg;
    /// @notice Uniswap v3's swap router, through which every short sells and buys back.
    IV3SwapRouter public immutable router;
    /// @notice The most gross exposure a short may carry per unit of equity from the regular open on the last trading
    /// day before a closure until the reopening, in basis points: the engine's weekend leverage cap.
    uint32 public immutable weekendLeverage;

    bytes32[] private _symbols;
    address[] private _tokens;
    uint24[] private _fees;
    mapping(address token => uint256) private _indexOfToken;
    mapping(uint256 asset => uint256) private _epochs;
    mapping(uint256 asset => mapping(uint256 epoch => Book)) private _books;
    mapping(uint256 asset => Mark) private _marks;
    mapping(address account => mapping(uint256 asset => Short)) private _shorts;

    /// @notice `caller` added `amount` of USDG to `account`'s short of `symbol`.
    event Deposit(address indexed caller, address indexed account, bytes32 indexed symbol, uint256 amount);
    /// @notice `caller` took `amount` of USDG from `account`'s short of `symbol` to `receiver`.
    event Withdraw(
        address indexed caller, address indexed account, bytes32 indexed symbol, uint256 amount, address receiver
    );
    /// @notice `account` borrowed `amount` of `symbol` for `shares` borrow shares and sold it for `proceeds` of USDG.
    event Sell(address indexed account, bytes32 indexed symbol, uint256 amount, uint256 proceeds, uint256 shares);
    /// @notice `account` bought back `amount` of `symbol` for `cost` of USDG and repaid it, retiring `shares`.
    event Cover(address indexed account, bytes32 indexed symbol, uint256 amount, uint256 cost, uint256 shares);
    /// @notice `caller` bought back `amount` of `account`'s short of `symbol` for `cost`, taking `bonus` of its USDG,
    /// and wrote off `writtenOff` it could not buy back.
    event Liquidate(
        address indexed caller,
        address indexed account,
        bytes32 indexed symbol,
        uint256 amount,
        uint256 cost,
        uint256 bonus,
        uint256 writtenOff
    );
    /// @notice The shorts of `symbol` bought in `amount` for its lending vault, for `cost` of their USDG.
    event BuyIn(bytes32 indexed symbol, uint256 amount, uint256 cost);
    /// @notice `symbol` fell `RESTRICTION_DROP_BPS` from its last close: until the end of day `until`, counted in days
    /// since 1970, a short sells at no less than the band's centre.
    event Restricted(bytes32 indexed symbol, uint32 until);
    /// @notice The vault was owed nothing of `symbol` while borrow shares remained, so its shorts start a new book,
    /// `epoch`; the earlier shorts owe nothing.
    event NewEpoch(bytes32 indexed symbol, uint256 epoch);

    /// @notice The fees are not one per asset.
    error LengthMismatch();
    /// @notice The router's factory has no pool of `token` and USDG at `fee`.
    error InvalidPool(address token, uint24 fee);
    /// @notice `symbol` has no Stock Token, no pool to short through, or no lending vault.
    error NotShortable(bytes32 symbol);
    /// @notice `caller` may not act for `account`.
    error Unauthorized(address caller, address account);
    /// @notice An amount of zero.
    error ZeroAmount();
    /// @notice The guardian has paused new risk.
    error BorrowingIsPaused();
    /// @notice The issuer has blocklisted `account`.
    error Blocked(address account);
    /// @notice `symbol`'s band is halted.
    error AssetHalted(bytes32 symbol);
    /// @notice `symbol`'s Stock Token has a multiplier change pending.
    error CorporateActionPending(bytes32 symbol);
    /// @notice The band cannot tell whether the session is open.
    error SessionUnknown();
    /// @notice The engine could not read `symbol`'s liquidity.
    error LiquidityUnknown(bytes32 symbol);
    /// @notice The L2 sequencer is down or was, within the hour.
    error SequencerNotSettled();
    /// @notice The short's equity, in USD with 18 decimals, is below its requirement.
    error InsufficientMargin(int256 equity, uint256 requirement);
    /// @notice The short's gross exposure exceeds its equity times the weekend leverage cap, from the regular open on
    /// the last trading day before a closure until the reopening.
    error WeekendLeverageExceeded(uint256 gross, int256 equity);
    /// @notice The short holds less than `amount` of USDG.
    error InsufficientCollateral(uint256 amount);
    /// @notice The short owes more of the cost of buy-ins than it holds: only a liquidation closes it.
    error Insolvent();
    /// @notice Another short of the book owes more of the buy-ins than it holds, so the book cannot pay this out until
    /// that short is liquidated.
    error DeficitOpen();
    /// @notice The account has no short of the asset.
    error NoShort();
    /// @notice The short's equity is not below its requirement.
    error NotShortfall(int256 equity, uint256 requirement);
    /// @notice The pool asks more than a liquidation may pay to buy back what the short owes.
    error AboveLimit();
    /// @notice `caller` is not the asset's lending vault.
    error NotLending(address caller);

    /// @param accounts_ The margin accounts, whose band, engine, USDG, lending vaults, authorizations and guardian the
    /// shorts use.
    /// @param router_ Uniswap v3's SwapRouter02.
    /// @param fees The fee tier, in hundredths of a basis point, of each asset's pool with USDG, in the engine's order;
    /// zero where the asset may not be shorted.
    constructor(MarginAccounts accounts_, IV3SwapRouter router_, uint24[] memory fees) {
        accounts = accounts_;
        band = accounts_.band();
        engine = accounts_.engine();
        usdg = accounts_.usdg();
        router = router_;
        weekendLeverage = engine.weekendLeverage();
        (bytes32[] memory symbols, address[] memory tokens) = accounts_.stocks();
        if (fees.length != symbols.length) revert LengthMismatch();
        IUniswapV3Factory factory = IUniswapV3Factory(router_.factory());
        for (uint256 i; i < symbols.length; ++i) {
            if (fees[i] != 0) {
                // forge-lint: disable-start(calls-loop, require-revert-in-loop)
                // slither-disable-next-line calls-loop
                if (tokens[i] == address(0) || factory.getPool(tokens[i], address(usdg), fees[i]) == address(0)) {
                    revert InvalidPool(tokens[i], fees[i]);
                }
                // forge-lint: disable-end(calls-loop, require-revert-in-loop)
                _indexOfToken[tokens[i]] = i + 1;
                IERC20(tokens[i]).forceApprove(address(router_), type(uint256).max);
            }
        }
        (_symbols, _tokens, _fees) = (symbols, tokens, fees);
        usdg.forceApprove(address(router_), type(uint256).max);
    }

    /// @notice The fee tier of `symbol`'s pool with USDG, in hundredths of a basis point; zero where it may not be
    /// shorted.
    function fee(bytes32 symbol) external view returns (uint24) {
        return _fees[_index(symbol)];
    }

    /// @notice `account`'s short of `symbol`: its USDG once it pays its part of every buy-in, negative if it cannot,
    /// what it owes the lending vault, in raw units of the token, its borrow shares, and the epoch of its book; a short
    /// of an earlier epoch owes nothing.
    function position(address account, bytes32 symbol)
        external
        view
        returns (int256 usdgHeld, uint256 debt, uint256 shares, uint256 epoch_)
    {
        uint256 i = _index(symbol);
        Short memory s = _shorts[account][i];
        return (_settled(s, _books[i][s.epoch].costIndex), _debt(i, s), s.shares, s.epoch);
    }

    /// @notice `account`'s short of `symbol` as the engine sees it now: its equity and requirement at the band's high
    /// edge, in USD with 18 decimals, the engine's missing bits and the regime.
    function health(address account, bytes32 symbol)
        external
        view
        returns (int256 equity, uint256 requirement, uint8 missing, uint8 regime)
    {
        uint256 i = _index(symbol);
        // slither-disable-next-line unused-return
        (,,,,, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        return _health(account, i, high);
    }

    /// @notice The epoch of `symbol`'s current book: new shorts borrow in it.
    function epoch(bytes32 symbol) external view returns (uint256) {
        return _epochs[_index(symbol)];
    }

    /// @notice The borrow shares, the USDG and the cost index of the shorts of `symbol` in book `epoch_`.
    function book(bytes32 symbol, uint256 epoch_)
        external
        view
        returns (uint256 shares, uint256 usdgHeld, uint256 costIndex)
    {
        Book memory b = _books[_index(symbol)][epoch_];
        return (b.shares, b.usdg, b.costIndex);
    }

    /// @notice Whether a short of `symbol` must sell at no less than the band's centre now, after a
    /// `RESTRICTION_DROP_BPS` fall from the asset's last close, and the band's centre the shorts last recorded as of an
    /// earlier day: that close.
    function restriction(bytes32 symbol) external view returns (bool restricted, uint64 close) {
        uint256 i = _index(symbol);
        // slither-disable-next-line unused-return
        (,, uint64 mid,,,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        Mark memory m = _marked(_marks[i], mid);
        // forge-lint: disable-next-line(block-timestamp)
        return (_today() <= m.until, m.close);
    }

    /// @notice Records `symbol`'s band centre for the restriction after a `RESTRICTION_DROP_BPS` fall: the last one
    /// recorded on an earlier day is that asset's close. Every sale records one; anyone may call it.
    function mark(bytes32 symbol) external {
        uint256 i = _index(symbol);
        // slither-disable-next-line unused-return
        (uint8 state,, uint64 mid,,,) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if (state == 0) revert AssetHalted(symbol);
        _mark(i, mid);
    }

    /// @notice Adds `amount` of the caller's USDG to `account`'s short of `symbol`. Only the account or an address it
    /// authorized in the margin accounts may.
    function deposit(bytes32 symbol, uint256 amount, address account) external {
        if (amount == 0) revert ZeroAmount();
        _checkAuthorized(account);
        uint256 i = _index(symbol);
        Short storage s = _shorts[account][i];
        _enter(i, s);
        _settleOrRevert(i, s);
        s.usdg += SafeCast.toUint128(amount);
        if (s.shares != 0) _books[i][s.epoch].usdg += SafeCast.toUint128(amount);
        emit Deposit(msg.sender, account, symbol, amount);
        usdg.safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Sends `amount` of USDG from `account`'s short of `symbol` to `receiver`. A short that owes the vault must
    /// then pass its checks, and waits while the guardian has paused borrowing. Only the account or an address it
    /// authorized in the margin accounts may.
    function withdraw(bytes32 symbol, uint256 amount, address account, address receiver) external {
        if (amount == 0) revert ZeroAmount();
        _checkAuthorized(account);
        uint256 i = _index(symbol);
        Short storage s = _shorts[account][i];
        _enter(i, s);
        _settleOrRevert(i, s);
        if (amount > s.usdg) revert InsufficientCollateral(amount);
        s.usdg -= SafeCast.toUint128(amount);
        if (s.shares != 0) {
            if (accounts.borrowingPaused()) revert BorrowingIsPaused();
            _leave(_books[i][s.epoch], amount);
            _check(account, i);
        }
        emit Withdraw(msg.sender, account, symbol, amount, receiver);
        usdg.safeTransfer(receiver, amount);
    }

    /// @notice Borrows `amount` of `symbol` from its lending vault for `account`'s short and sells it through its pool
    /// for at least `minProceeds` of USDG and never below the band's low edge, or its centre while the asset is
    /// restricted after a `RESTRICTION_DROP_BPS` fall. The proceeds stay in the short, which must then pass its checks.
    /// Only the account or an address it authorized may, for an account the issuer has not blocklisted, while the
    /// guardian has not paused borrowing.
    // slither-disable-next-line reentrancy-no-eth
    function sell(bytes32 symbol, uint256 amount, uint256 minProceeds, address account)
        external
        returns (uint256 proceeds)
    {
        if (amount == 0) revert ZeroAmount();
        if (accounts.borrowingPaused()) revert BorrowingIsPaused();
        _checkAuthorized(account);
        uint256 i = _index(symbol);
        uint256 least = _divUp(amount * _floor(account, i), TOKEN_TO_USDG);
        StockLendingVault lending = _lending(i);
        uint256 owed = lending.debt();
        uint256 current = _epochs[i];
        if (owed == 0 && _books[i][current].shares != 0) {
            _epochs[i] = ++current;
            emit NewEpoch(symbol, current);
        }
        Short storage s = _shorts[account][i];
        _enter(i, s);
        _settleOrRevert(i, s);
        Book storage b = _books[i][current];
        uint256 shares = _toSharesUp(amount, owed, b.shares);
        if (s.shares == 0) b.usdg += s.usdg;
        s.shares += SafeCast.toUint128(shares);
        b.shares += SafeCast.toUint128(shares);
        lending.borrow(amount, address(this));
        proceeds = router.exactInputSingle(
            IV3SwapRouter.ExactInputSingleParams(
                _tokens[i], address(usdg), _fees[i], address(this), amount, least > minProceeds ? least : minProceeds, 0
            )
        );
        s.usdg += SafeCast.toUint128(proceeds);
        b.usdg += SafeCast.toUint128(proceeds);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Sell(account, symbol, amount, proceeds, shares);
        _check(account, i);
    }

    /// @notice Buys back `amount` of `symbol`, at most what `account`'s short owes, through its pool for at most
    /// `maxCost` of the short's USDG, and repays it to the lending vault. A buy-back of part of what it owes pays no more
    /// than the band's high edge, so it never leaves the short weaker; buying back all it owes, or so much that it would
    /// retire all the short's borrow shares, closes the short. Only the account or an address it authorized may.
    // slither-disable-next-line reentrancy-no-eth
    function cover(bytes32 symbol, uint256 amount, uint256 maxCost, address account) external returns (uint256 cost) {
        if (amount == 0) revert ZeroAmount();
        _checkAuthorized(account);
        uint256 i = _index(symbol);
        Short storage s = _shorts[account][i];
        _enter(i, s);
        if (s.shares == 0) revert NoShort();
        _settleOrRevert(i, s);
        uint256 owed = _debt(i, s);
        if (amount >= owed || _toSharesDown(amount, _lending(i).debt(), _books[i][s.epoch].shares) >= s.shares) {
            amount = owed;
        }
        uint256 limit = maxCost < s.usdg ? maxCost : s.usdg;
        if (amount < owed) {
            (, uint256 high) = _quote(i);
            uint256 atHigh = _divUp(amount * high, TOKEN_TO_USDG);
            if (atHigh < limit) limit = atHigh;
        }
        if (amount != 0) cost = _buy(i, amount, limit);
        uint256 shares = _repay(i, s, amount, amount == owed);
        _spend(i, s, cost);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Cover(account, symbol, amount, cost, shares);
    }

    /// @notice Buys back all `account`'s short of `symbol` owes once its equity is below its requirement, through its
    /// pool and at no more than `MAX_PREMIUM_BPS` above the band's centre, nor its high edge, and pays the caller
    /// `LIQUIDATION_BONUS_BPS` of the cost from the short's USDG. While the market is closed a short falls short only if
    /// it does so at the band's low edge too. A short whose USDG cannot buy back all it owes spends all of it and leaves
    /// the rest written off, a loss to the vault's lenders; what it owed of the shorts' buy-ins and could not pay is
    /// charged to the book's other shorts, and the lenders bear it too, as far as those still owe the vault. Anyone may
    /// call it.
    // slither-disable-next-line reentrancy-no-eth
    function liquidate(address account, bytes32 symbol) external {
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        uint256 i = _index(symbol);
        Short storage s = _shorts[account][i];
        if (s.shares == 0) revert NoShort();
        uint256 price = _shortfall(account, i);
        uint256 deficit = _settle(i, s);
        uint256 owed = _debt(i, s);
        (uint256 bought, uint256 cost, uint256 bonus) = _buyBack(i, owed, s.usdg, price);
        _repay(i, s, bought, bought == owed);
        _spend(i, s, cost);
        uint256 writtenOff = bought == owed && deficit == 0 ? 0 : _writeOff(i, s, owed - bought, deficit, price);
        if (bonus != 0) s.usdg -= SafeCast.toUint128(bonus);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Liquidate(msg.sender, account, symbol, bought, cost, bonus, writtenOff);
        if (bonus != 0) usdg.safeTransfer(msg.sender, bonus);
    }

    /// @notice Buys `assets` of the calling lending vault's token through its pool, at no more than `MAX_PREMIUM_BPS`
    /// above the band's centre, nor its high edge, and with no more than the USDG of the token's book, and repays them
    /// to the vault; every short of the book pays its part of the cost in proportion to its borrow shares. Where that
    /// would leave the vault owed a remnant too small for the book's borrow shares, it buys in all the vault is owed.
    /// Only the asset's lending vault may, while the L2 sequencer is settled.
    // slither-disable-next-line reentrancy-no-eth
    function buyIn(uint256 assets) external {
        uint256 i = _indexOfToken[StockLendingVault(msg.sender).asset()];
        if (i == 0 || address(accounts.lending(_symbols[--i])) != msg.sender) revert NotLending(msg.sender);
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        Book storage b = _books[i][_epochs[i]];
        uint256 owed = StockLendingVault(msg.sender).debt();
        if ((owed - assets) * MAX_SHARES_PER_UNIT < b.shares) assets = owed;
        (uint256 mid, uint256 high) = _quote(i);
        uint256 limit = _divUp(assets * _limitPrice(mid, high), TOKEN_TO_USDG);
        uint256 cost = _buy(i, assets, b.usdg < limit ? b.usdg : limit);
        b.usdg -= SafeCast.toUint128(cost);
        b.costIndex += _divUp(cost * RAY, b.shares);
        // forge-lint: disable-next-line(reentrancy-events)
        emit BuyIn(_symbols[i], assets, cost);
        _repayTo(StockLendingVault(msg.sender), i, assets);
    }

    /// @dev Checks the issuer has not blocklisted `account` or the caller, and returns the least a sale may take for a
    /// token: the band's low edge, or its centre while the asset is restricted, which the sale records.
    function _floor(address account, uint256 i) private returns (uint256) {
        IStockTokenRegistry registry = IStockTokenRegistry(IStockToken(_tokens[i]).ACCESS_CONTROLLED_REGISTRY());
        if (registry.isBlocked(account)) revert Blocked(account);
        if (msg.sender != account && registry.isBlocked(msg.sender)) revert Blocked(msg.sender);
        // slither-disable-next-line unused-return
        (,, uint64 mid,, uint64 low,) = band.quote(_symbols[i]); // forge-lint: disable-line(unused-return)
        return _mark(i, mid) ? mid : low;
    }

    /// @dev Reverts unless the short falls short: at the band's high edge while the market is open, and at both edges
    /// while it is not. Returns the most a liquidation may pay for a token.
    function _shortfall(address account, uint256 i) private view returns (uint256) {
        bytes32 symbol = _symbols[i];
        // slither-disable-next-line unused-return
        (uint8 state,, uint64 mid,, uint64 low, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if (state == 0) revert AssetHalted(symbol);
        (int256 equity, uint256 requirement,, uint8 regime) = _health(account, i, high);
        if (regime != 2) (equity, requirement,,) = _health(account, i, low);
        if (equity >= 0 && SafeCast.toUint256(equity) >= requirement) revert NotShortfall(equity, requirement);
        return _limitPrice(mid, high);
    }

    /// @dev Retires all the short's borrow shares. In the current book, writes off `rest` of what it owed, or all the vault
    /// is still owed once no short is left; where it owed `deficit` of the buy-ins it could not pay, charges that to the
    /// book's other shorts and writes off what it buys at `price`, so the lenders bear it as far as those shorts still owe
    /// the vault. It never leaves the vault owed less than a raw unit for every `MAX_SHARES_PER_UNIT` of the book's
    /// borrow shares: what that keeps falls on the book's other shorts. In an earlier book, which owes the vault nothing,
    /// the deficit falls on that book's other shorts.
    function _writeOff(uint256 i, Short storage s, uint256 rest, uint256 deficit, uint256 price)
        private
        returns (uint256 written)
    {
        Book storage b = _books[i][s.epoch];
        b.shares -= s.shares;
        s.shares = 0;
        if (deficit != 0 && b.shares != 0) b.costIndex += _divUp(deficit * RAY, b.shares);
        if (s.epoch != _epochs[i]) return 0;
        StockLendingVault lending = _lending(i);
        uint256 owed = lending.debt();
        written = rest + deficit * TOKEN_TO_USDG / price;
        if (b.shares == 0 || written >= owed) {
            written = owed;
        } else {
            uint256 keep = _divUp(b.shares, MAX_SHARES_PER_UNIT);
            if (owed - written < keep) written = owed > keep ? owed - keep : 0;
        }
        // slither-disable-next-line unused-return
        lending.writeOff(written); // forge-lint: disable-line(unused-return)
    }

    /// @dev Buys back `owed` with at most `held` and at no more than `price` a token, and the bonus that leaves; where
    /// the pool asks more than `held` for it and `held` is below what it costs at `price`, buys what all `held` can.
    function _buyBack(uint256 i, uint256 owed, uint256 held, uint256 price)
        private
        returns (uint256 bought, uint256 cost, uint256 bonus)
    {
        if (owed == 0) return (0, 0, 0);
        uint256 limit = _divUp(owed * price, TOKEN_TO_USDG);
        try router.exactOutputSingle(
            IV3SwapRouter.ExactOutputSingleParams(
                address(usdg), _tokens[i], _fees[i], address(this), owed, held < limit ? held : limit, 0
            )
        ) returns (
            uint256 paid
        ) {
            return (owed, paid, _min(paid * LIQUIDATION_BONUS_BPS / BPS, held - paid));
        } catch {
            if (held >= limit) revert AboveLimit();
            if (held != 0) bought = _buyWith(i, held, price);
            return (bought, held, 0);
        }
    }

    function _buy(uint256 i, uint256 amount, uint256 maxCost) private returns (uint256) {
        return router.exactOutputSingle(
            IV3SwapRouter.ExactOutputSingleParams(
                address(usdg), _tokens[i], _fees[i], address(this), amount, maxCost, 0
            )
        );
    }

    function _buyWith(uint256 i, uint256 held, uint256 price) private returns (uint256) {
        return router.exactInputSingle(
            IV3SwapRouter.ExactInputSingleParams(
                address(usdg), _tokens[i], _fees[i], address(this), held, held * TOKEN_TO_USDG / price, 0
            )
        );
    }

    /// @dev Repays `amount` of the short's debt to the vault and retires its borrow shares, never more than it holds: all
    /// of them once it repays all it owes.
    function _repay(uint256 i, Short storage s, uint256 amount, bool all) private returns (uint256 shares) {
        Book storage b = _books[i][s.epoch];
        shares = all ? s.shares : _min(_toSharesDown(amount, _lending(i).debt(), b.shares), s.shares);
        s.shares -= SafeCast.toUint128(shares);
        b.shares -= SafeCast.toUint128(shares);
        if (amount != 0) _repayTo(_lending(i), i, amount);
    }

    function _repayTo(StockLendingVault lending, uint256 i, uint256 amount) private {
        IERC20 token = IERC20(_tokens[i]);
        if (token.allowance(address(this), address(lending)) < amount) {
            token.forceApprove(address(lending), type(uint256).max);
        }
        // slither-disable-next-line unused-return
        lending.repay(amount); // forge-lint: disable-line(unused-return)
    }

    /// @dev Pays `cost` from the short's USDG; a short left with no borrow shares takes what is left out of its book.
    function _spend(uint256 i, Short storage s, uint256 cost) private {
        s.usdg -= SafeCast.toUint128(cost);
        _leave(_books[i][s.epoch], s.shares == 0 ? cost + s.usdg : cost);
    }

    function _leave(Book storage b, uint256 amount) private {
        uint256 held = b.usdg;
        if (amount > held) revert DeficitOpen();
        b.usdg = SafeCast.toUint128(held - amount);
    }

    /// @dev Moves a short of an earlier epoch, which owes nothing, out of its book with what it holds, into the current
    /// epoch.
    function _enter(uint256 i, Short storage s) private {
        uint256 current = _epochs[i];
        if (s.epoch == current) return;
        _settleOrRevert(i, s);
        if (s.shares != 0) {
            Book storage b = _books[i][s.epoch];
            _leave(b, s.usdg);
            b.shares -= s.shares;
            s.shares = 0;
        }
        s.epoch = SafeCast.toUint64(current);
        s.costIndex = SafeCast.toUint192(_books[i][current].costIndex);
    }

    function _settle(uint256 i, Short storage s) private returns (uint256 deficit) {
        uint256 index = _books[i][s.epoch].costIndex;
        uint256 paid = s.costIndex;
        if (index == paid) return 0;
        s.costIndex = SafeCast.toUint192(index);
        uint256 charge = _divUp(s.shares * (index - paid), RAY);
        uint256 held = s.usdg;
        if (charge <= held) {
            s.usdg = SafeCast.toUint128(held - charge);
            return 0;
        }
        s.usdg = 0;
        return charge - held;
    }

    function _settleOrRevert(uint256 i, Short storage s) private {
        if (_settle(i, s) != 0) revert Insolvent();
    }

    function _settled(Short memory s, uint256 index) private pure returns (int256) {
        return SafeCast.toInt256(s.usdg) - SafeCast.toInt256(_divUp(s.shares * (index - s.costIndex), RAY));
    }

    function _check(address account, uint256 i) private view {
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        uint256 high = _priced(_symbols[i]);
        (int256 equity, uint256 requirement, uint8 missing, uint8 regime) = _health(account, i, high);
        if (regime == 0) revert SessionUnknown();
        if (missing & (1 << i) != 0) revert LiquidityUnknown(_symbols[i]);
        if (equity < 0 || SafeCast.toUint256(equity) < requirement) revert InsufficientMargin(equity, requirement);
        uint256 exposure = _debt(i, _shorts[account][i]) * high;
        if (
            exposure * BPS > SafeCast.toUint256(equity) * weekendLeverage * PRICE_UNIT
                && (regime != 2 || _weekendAhead())
        ) {
            revert WeekendLeverageExceeded(exposure / PRICE_UNIT, equity);
        }
    }

    /// @dev The band's high edge for `symbol`, which must be neither halted nor waiting on a multiplier change.
    function _priced(bytes32 symbol) private view returns (uint256) {
        // slither-disable-next-line unused-return
        (uint8 state,,,,, uint128 high) = band.quote(symbol); // forge-lint: disable-line(unused-return)
        if (state == 0) revert AssetHalted(symbol);
        // slither-disable-next-line unused-return
        (uint8 status,,,) = band.corporateAction(symbol); // forge-lint: disable-line(unused-return)
        if (status != 0) revert CorporateActionPending(symbol);
        return high;
    }

    /// @dev The band's centre and high edge for asset `i`, which must not be halted.
    function _quote(uint256 i) private view returns (uint256 mid, uint256 high) {
        // slither-disable-next-line unused-return
        (uint8 state,, uint64 centre,,, uint128 top) = band.quote(_symbols[i]); // forge-lint: disable-line(unused-return)
        if (state == 0) revert AssetHalted(_symbols[i]);
        return (centre, top);
    }

    function _limitPrice(uint256 mid, uint256 high) private pure returns (uint256) {
        uint256 limit = mid * (BPS + MAX_PREMIUM_BPS) / BPS;
        return limit < high ? limit : high;
    }

    function _health(address account, uint256 i, uint256 price)
        private
        view
        returns (int256 equity, uint256 requirement, uint8 missing, uint8 regime)
    {
        Short memory s = _shorts[account][i];
        uint256 owed = _debt(i, s);
        int256[] memory quantities = new int256[](_symbols.length);
        uint256[] memory prices = new uint256[](_symbols.length);
        quantities[i] = -SafeCast.toInt256(owed);
        prices[i] = price;
        (requirement, missing, regime) = engine.currentRequirement(quantities, prices);
        equity = _settled(s, _books[i][s.epoch].costIndex) * SafeCast.toInt256(USDG_TO_USD)
            - SafeCast.toInt256(owed * price / PRICE_UNIT);
    }

    function _mark(uint256 i, uint64 mid) private returns (bool restricted) {
        Mark memory m = _marked(_marks[i], mid);
        if (m.until > _marks[i].until) emit Restricted(_symbols[i], m.until);
        _marks[i] = m;
        // forge-lint: disable-next-line(block-timestamp)
        return _today() <= m.until;
    }

    function _marked(Mark memory m, uint64 mid) private view returns (Mark memory) {
        if (mid == 0) return m;
        uint32 today = _today();
        // forge-lint: disable-next-line(block-timestamp)
        if (today > m.day) (m.day, m.close) = (today, m.last);
        m.last = mid;
        if (uint256(mid) * BPS <= uint256(m.close) * (BPS - RESTRICTION_DROP_BPS)) m.until = today + 1;
        return m;
    }

    function _today() private view returns (uint32) {
        return SafeCast.toUint32(block.timestamp / 1 days);
    }

    function _weekendAhead() private view returns (bool) {
        // slither-disable-next-line unused-return
        (uint8 state,,,, uint64 boundaryMs) = band.session(); // forge-lint: disable-line(unused-return)
        return state != 2 || boundaryMs != 0;
    }

    function _lending(uint256 i) private view returns (StockLendingVault lending) {
        lending = accounts.lending(_symbols[i]);
        if (address(lending) == address(0)) revert NotShortable(_symbols[i]);
    }

    function _index(bytes32 symbol) private view returns (uint256) {
        for (uint256 i; i < _symbols.length; ++i) {
            if (_symbols[i] == symbol) {
                if (_fees[i] == 0) break;
                return i;
            }
        }
        revert NotShortable(symbol);
    }

    /// @dev What a short owes the vault: its share of the debt in the current book, all of it as the book's last short,
    /// and nothing in an earlier book.
    function _debt(uint256 i, Short memory s) private view returns (uint256) {
        if (s.shares == 0 || s.epoch != _epochs[i]) return 0;
        uint256 total = _lending(i).debt();
        uint256 all = _books[i][s.epoch].shares;
        return s.shares == all ? total : _min(_toAssetsUp(s.shares, total, all), total);
    }

    function _checkAuthorized(address account) private view {
        if (msg.sender != account && !accounts.isAuthorized(account, msg.sender)) {
            revert Unauthorized(msg.sender, account);
        }
    }

    function _toSharesUp(uint256 assets, uint256 totalAssets, uint256 totalShares) private pure returns (uint256) {
        return _divUp(assets * (totalShares + VIRTUAL_SHARES), totalAssets + VIRTUAL_ASSETS);
    }

    function _toSharesDown(uint256 assets, uint256 totalAssets, uint256 totalShares) private pure returns (uint256) {
        return assets * (totalShares + VIRTUAL_SHARES) / (totalAssets + VIRTUAL_ASSETS);
    }

    function _toAssetsUp(uint256 shares, uint256 totalAssets, uint256 totalShares) private pure returns (uint256) {
        return _divUp(shares * (totalAssets + VIRTUAL_ASSETS), totalShares + VIRTUAL_SHARES);
    }

    function _divUp(uint256 a, uint256 b) private pure returns (uint256) {
        return a > 0 ? (a - 1) / b + 1 : 0;
    }

    function _min(uint256 a, uint256 b) private pure returns (uint256) {
        return a < b ? a : b;
    }
}
