// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IMargin} from "./interfaces/IMargin.sol";
import {IStockToken, IStockTokenRegistry} from "./interfaces/IStockToken.sol";
import {IUSDG} from "./interfaces/IUSDG.sol";
import {SupplyVault} from "./SupplyVault.sol";

/// @title Tapehouse margin accounts
/// @notice Each address holds a cross position and, for any asset it chooses, an isolated position in that asset
/// alone. A position holds Stock Tokens, USDG and WETH as collateral and borrows USDG from the supply vault
/// against them, margined by the engine. Deposits and repayments stay open whatever the band, the session and the
/// guardian; borrowing, and withdrawing from a position in debt, pass the engine's current requirement and the band's
/// checks first, and back no new risk with a Stock Token the accounts could not seize. A position
/// without Stock Tokens needs only its equity, and a position holding USDG repays with it rather than borrow. An
/// account may authorize another address to act for it.
/// @dev Position `bytes32(0)` is the cross position; an asset's symbol names its isolated position, which holds
/// no other Stock Token. An asset sits in one position of an account at a time. Stock Tokens are held in raw
/// units and valued at the low edge of their band, which already carries the token's multiplier. Each
/// position's debt is its share of the vault's debt, kept as borrow shares with Morpho Blue's virtual offsets.
/// Each position's Stock Tokens are units of the accounts' holding of that asset, so a burn by the issuer falls on
/// that asset's holders in proportion, and a Stock Token moves through the accounts only for an account the issuer
/// has not blocklisted. Deposits, total debt and weekend debt are capped for good at deployment, and a guardian can
/// pause new borrowing. While the 24/5 session is closed, each position's debt also
/// carries a premium, owed beside it and paid after it, that funds the backstop; part of every premium paid goes to a
/// fee reserve.
contract MarginAccounts is Ownable2Step {
    using SafeERC20 for IERC20;

    struct Position {
        uint128 debtShares;
        uint128 usdg;
        uint128 weth;
        uint8 held;
        uint120 premium;
        uint256 premiumIndex;
    }

    struct Closure {
        uint64 closesMs;
        uint64 reopensMs;
        uint64 accruedMs;
        uint32 premiumRate;
    }

    struct Book {
        uint128 units;
        uint128 scale;
    }

    struct Valuation {
        int256[] quantities;
        uint256[] prices;
        int256 equity;
        uint256 gross;
        uint8 held;
    }

    /// @notice The cross position.
    bytes32 public constant CROSS = bytes32(0);
    /// @notice What WETH counts for as collateral, in basis points of its Chainlink value: Aave v3's loan-to-value
    /// for WETH on Arbitrum.
    uint256 public constant WETH_COLLATERAL_BPS = 80_00;
    /// @notice The oldest ETH/USD answer the accounts value WETH at, in seconds: its heartbeat plus a minute.
    uint256 public constant MAX_FEED_AGE = 86_400 + 60;
    /// @notice The halvings of the price range a liquidation price is searched over.
    uint256 public constant LIQUIDATION_PRICE_STEPS = 14;
    /// @notice The highest premium rate, in basis points a year.
    uint32 public constant MAX_PREMIUM_RATE = 100_00;

    uint256 private constant BPS = 10_000;
    uint256 private constant PRICE_UNIT = 1e8;
    uint256 private constant USDG_TO_USD = 1e12;
    uint256 private constant VIRTUAL_SHARES = 1e6;
    uint256 private constant VIRTUAL_ASSETS = 1;
    uint256 private constant WAD = 1e18;
    uint256 private constant RAY = 1e27;
    uint256 private constant YEAR_MS = 365 days * 1000;
    uint64 private constant MAX_CLOSURE_MS = 72 hours * 1000;

    /// @notice The band that prices every Stock Token.
    IBand public immutable band;
    /// @notice The engine that margins every position.
    IMargin public immutable engine;
    /// @notice The vault the positions borrow from.
    SupplyVault public immutable vault;
    /// @notice The token the positions borrow and repay.
    IERC20 public immutable usdg;
    /// @notice Wrapped ETH, taken as collateral.
    IERC20 public immutable weth;
    /// @notice The Chainlink feed that prices WETH; zero where the engine has none.
    AggregatorV3Interface public immutable ethUsd;
    /// @notice The most gross exposure a position may carry per unit of equity from the regular open on the last
    /// trading day before a closure until the reopening, in basis points: the engine's weekend leverage cap.
    uint32 public immutable weekendLeverage;
    /// @notice The most USDG the positions may owe together, interest included.
    uint256 public immutable debtCap;
    /// @notice The most USDG the positions may owe together for a new loan from the regular open on the last trading
    /// day before a closure until the reopening, or while the band cannot tell the session.
    uint256 public immutable weekendDebtCap;
    /// @notice The part of every premium paid that goes to the reserve, in basis points.
    uint16 public immutable reserveShare;

    /// @notice The only address that may seize collateral and write a position's debt off.
    address public liquidator;
    /// @notice The address that may pause new borrowing.
    address public guardian;
    /// @notice Whether new borrowing, and withdrawals from positions in debt, are paused.
    bool public borrowingPaused;
    /// @notice The borrow shares of every position.
    uint128 public totalDebtShares;
    /// @notice The premium each borrow share has owed since deployment, in USDG with 27 more decimals, as last
    /// accrued.
    uint128 public premiumIndex;
    /// @notice The fees kept, in USDG.
    uint128 public reserve;
    /// @notice The premium paid and set aside for the backstop, in USDG.
    uint128 public backstopPremium;
    /// @notice The only address that may claim the premium set aside.
    address public backstop;
    /// @notice Whether `authorized` may borrow, withdraw and repay with collateral for `account`, and deposit Stock
    /// Tokens and USDG into it.
    mapping(address account => mapping(address authorized => bool)) public isAuthorized;

    bytes32[] private _symbols;
    address[] private _tokens;
    mapping(bytes32 symbol => uint256) private _indexOfSymbol;
    mapping(address token => uint256) private _indexOfToken;
    uint256[] private _caps;
    Book[] private _books;
    address[] private _registries;
    mapping(address account => mapping(bytes32 position => Position)) private _positions;
    mapping(address account => mapping(bytes32 position => mapping(uint256 asset => uint256))) private _units;
    Closure private _closure;

    /// @notice `caller` deposited `amount` of `token` into `account`'s `position`.
    event Deposit(
        address indexed caller, address indexed account, bytes32 indexed position, address token, uint256 amount
    );
    /// @notice `caller` withdrew `amount` of `token` from `account`'s `position` to `receiver`.
    event Withdraw(
        address indexed caller,
        address indexed account,
        bytes32 indexed position,
        address token,
        uint256 amount,
        address receiver
    );
    /// @notice `caller` borrowed `assets` of USDG against `account`'s `position` to `receiver`, for `shares` borrow
    /// shares.
    event Borrow(
        address indexed caller,
        address indexed account,
        bytes32 indexed position,
        uint256 assets,
        uint256 shares,
        address receiver
    );
    /// @notice `caller` repaid `assets` of USDG of `account`'s `position`, retiring `shares` borrow shares.
    event Repay(
        address indexed caller, address indexed account, bytes32 indexed position, uint256 assets, uint256 shares
    );
    /// @notice The liquidator took `amount` of `token` from `account`'s `position` to `receiver`.
    event Seize(address indexed account, bytes32 indexed position, address token, uint256 amount, address receiver);
    /// @notice `account`'s `position` paid `premium` of USDG, `toReserve` of it to the reserve and the rest to the
    /// backstop.
    event PremiumPaid(address indexed account, bytes32 indexed position, uint256 premium, uint256 toReserve);
    /// @notice The liquidator wrote off `assets` of `account`'s `position`, retiring `shares` borrow shares, and
    /// `premium` it owed.
    event WriteOff(address indexed account, bytes32 indexed position, uint256 assets, uint256 shares, uint256 premium);
    /// @notice `account` allowed or stopped `authorized` acting for it.
    event AuthorizationSet(address indexed account, address indexed authorized, bool allowed);
    /// @notice The liquidator was set.
    event LiquidatorSet(address indexed liquidator);
    /// @notice The guardian was set.
    event GuardianSet(address indexed guardian);
    /// @notice The backstop was set.
    event BackstopSet(address indexed backstop);
    /// @notice The premium rate became `rate` basis points a year.
    event PremiumRateSet(uint32 rate);
    /// @notice Each borrow share has owed `index` of premium since deployment.
    event PremiumAccrued(uint256 index);
    /// @notice The closure the premium accrues over runs from `closesMs` until `reopensMs`, in milliseconds; a zero
    /// `reopensMs` is not yet known.
    event ClosureSet(uint64 closesMs, uint64 reopensMs);
    /// @notice The backstop claimed `amount` of premium.
    event PremiumClaimed(address indexed backstop, uint256 amount);
    /// @notice The owner took `amount` of the reserve to `receiver`.
    event ReserveWithdrawn(address indexed receiver, uint256 amount);
    /// @notice Anyone cleared `units` of `symbol` a burn left worth nothing from `account`'s `position`.
    event HoldingCleared(address indexed account, bytes32 indexed position, bytes32 indexed symbol, uint256 units);
    /// @notice The guardian paused or resumed new borrowing.
    event BorrowingPausedSet(bool paused);
    /// @notice The accounts held `balance` of `symbol`, less than the `counted` their positions held, and every
    /// position's holding of it fell in proportion.
    event AssetWrittenDown(bytes32 indexed symbol, uint256 counted, uint256 balance);

    /// @notice The engine's band is `engineBand`, not the accounts' band.
    error BandMismatch(address engineBand);
    /// @notice `position` is neither the cross position nor an asset with a Stock Token.
    error UnknownPosition(bytes32 position);
    /// @notice `symbol` is not one of the engine's assets.
    error UnknownAsset(bytes32 symbol);
    /// @notice `token` is not collateral here, or not in `position`.
    error UnsupportedToken(address token, bytes32 position);
    /// @notice The account holds `symbol` in `position` already.
    error AssetInOtherPosition(bytes32 symbol, bytes32 position);
    /// @notice The position holds less than `amount` of `token`.
    error InsufficientCollateral(address token, uint256 amount);
    /// @notice The position's equity, in USD with 18 decimals, is below its requirement.
    error InsufficientMargin(int256 equity, uint256 requirement);
    /// @notice The position's gross exposure exceeds its equity times the weekend leverage cap, from the regular open
    /// on the last trading day before a closure until the reopening.
    error WeekendLeverageExceeded(uint256 gross, int256 equity);
    /// @notice A position holding USDG does not borrow USDG: it repays with it, through `repayWithCollateral`.
    error BorrowingAgainstUsdg();
    /// @notice `symbol`'s band is halted.
    error AssetHalted(bytes32 symbol);
    /// @notice `symbol`'s Stock Token has a multiplier change pending.
    error CorporateActionPending(bytes32 symbol);
    /// @notice An amount of zero.
    error ZeroAmount();
    /// @notice The engine could not read `symbol`'s liquidity.
    error LiquidityUnknown(bytes32 symbol);
    /// @notice The band cannot tell whether the session is open.
    error SessionUnknown();
    /// @notice The L2 sequencer is down or was, within the hour.
    error SequencerNotSettled();
    /// @notice `caller` may not act for `account`.
    error Unauthorized(address caller, address account);
    /// @notice Only the liquidator may seize collateral and write debt off.
    error NotLiquidator(address caller);
    /// @notice The liquidator is already set, to `current`.
    error LiquidatorAlreadySet(address current);
    /// @notice The liquidator cannot be the zero address.
    error InvalidLiquidator();
    /// @notice Only the backstop may claim the premium.
    error NotBackstop(address caller);
    /// @notice The backstop is already set, to `current`.
    error BackstopAlreadySet(address current);
    /// @notice The backstop cannot be the zero address.
    error InvalidBackstop();
    /// @notice The premium rate is above `MAX_PREMIUM_RATE`, or the reserve's share above the whole.
    error InvalidPremium();
    /// @notice The reserve holds only `reserve`.
    error InsufficientReserve(uint256 reserve);
    /// @notice The reserve goes to an address.
    error InvalidReceiver();
    /// @notice The issuer has paused `symbol`'s Stock Token or blocklisted the accounts, so it could not be seized: it
    /// backs no new loan or withdrawal in debt.
    error AssetFrozen(bytes32 symbol);
    /// @notice The accounts' holding of `symbol` in the position is still worth a raw unit or more.
    error HoldingNotEmpty(bytes32 symbol);
    /// @notice The engine has `count` assets; the accounts take at most eight.
    error TooManyAssets(uint256 count);
    /// @notice The weekend debt cap is above the debt cap.
    error InvalidDebtCaps();
    /// @notice The owner cannot renounce: the liquidator, the backstop and the premium rate need one.
    error OwnershipCannotBeRenounced();
    /// @notice A position still holding collateral cannot be written off.
    error PositionNotEmpty();
    /// @notice The caps are not one per asset.
    error LengthMismatch();
    /// @notice The accounts would hold more than `cap` of `symbol`.
    error AssetCapExceeded(bytes32 symbol, uint256 cap);
    /// @notice The positions would owe `debt`, more than `cap`.
    error DebtCapExceeded(uint256 debt, uint256 cap);
    /// @notice The positions would owe `debt` across a closure, more than `cap`.
    error WeekendDebtCapExceeded(uint256 debt, uint256 cap);
    /// @notice The guardian has paused new borrowing.
    error BorrowingIsPaused();
    /// @notice Only the guardian may pause new borrowing.
    error NotGuardian(address caller);
    /// @notice Every unit of `symbol` the accounts held has been burned: it takes no deposit.
    error AssetWrittenOff(bytes32 symbol);
    /// @notice The issuer has blocklisted `account`, so its Stock Tokens do not move through the accounts either.
    error Blocked(address account);

    /// @param band_ The band, which must be the engine's.
    /// @param engine_ The margin engine; its assets are the Stock Tokens the positions take.
    /// @param vault_ The supply vault, whose asset is USDG and whose borrower these accounts become.
    /// @param weth_ Wrapped ETH.
    /// @param initialOwner The owner, who sets the liquidator once and the guardian.
    /// @param assetCaps The most of each asset's Stock Token the accounts may hold, in raw units, in the engine's
    /// order.
    /// @param debtCap_ The most USDG the positions may owe together.
    /// @param weekendDebtCap_ The most USDG the positions may owe together for a new loan before and across a closure.
    /// @param premiumRate_ The premium each position's debt carries while the session is closed, in basis points a
    /// year.
    /// @param reserveShare_ The part of every premium paid that goes to the reserve, in basis points.
    constructor(
        IBand band_,
        IMargin engine_,
        SupplyVault vault_,
        IERC20 weth_,
        address initialOwner,
        uint256[] memory assetCaps,
        uint256 debtCap_,
        uint256 weekendDebtCap_,
        uint32 premiumRate_,
        uint16 reserveShare_
    ) Ownable(initialOwner) {
        if (premiumRate_ > MAX_PREMIUM_RATE || reserveShare_ > BPS) revert InvalidPremium();
        if (weekendDebtCap_ > debtCap_) revert InvalidDebtCaps();
        address engineBand = engine_.band();
        if (engineBand != address(band_)) revert BandMismatch(engineBand);
        band = band_;
        engine = engine_;
        vault = vault_;
        usdg = IERC20(vault_.asset());
        weth = weth_;
        ethUsd = AggregatorV3Interface(engine_.ethUsdFeed());
        weekendLeverage = engine_.weekendLeverage();
        debtCap = debtCap_;
        weekendDebtCap = weekendDebtCap_;
        reserveShare = reserveShare_;
        _closure = Closure(0, 0, SafeCast.toUint64(block.timestamp * 1000), premiumRate_);
        bytes32[] memory symbols = engine_.assets();
        if (symbols.length > 8) revert TooManyAssets(symbols.length);
        if (assetCaps.length != symbols.length) revert LengthMismatch();
        _caps = assetCaps;
        for (uint256 i; i < symbols.length; ++i) {
            _books.push(Book(0, SafeCast.toUint128(WAD)));
            // slither-disable-next-line unused-return,calls-loop
            (,,, address token) = band_.asset(symbols[i]); // forge-lint: disable-line(calls-loop, unused-return)
            _symbols.push(symbols[i]);
            _tokens.push(token);
            // forge-lint: disable-start(calls-loop)
            // slither-disable-next-line calls-loop
            _registries.push(token == address(0) ? address(0) : IStockToken(token).ACCESS_CONTROLLED_REGISTRY());
            // forge-lint: disable-end(calls-loop)
            _indexOfSymbol[symbols[i]] = i + 1;
            if (token != address(0)) _indexOfToken[token] = i + 1;
        }
        usdg.forceApprove(address(vault_), type(uint256).max);
        emit PremiumRateSet(premiumRate_);
    }

    /// @notice The assets, in the engine's order, and their Stock Tokens; a zero token takes no deposit.
    function stocks() external view returns (bytes32[] memory symbols, address[] memory tokens) {
        return (_symbols, _tokens);
    }

    /// @notice What `account`'s `position` holds of `token`: a Stock Token at its holding's last synced scale.
    function collateral(address account, bytes32 position, address token) external view returns (uint256) {
        if (token == address(usdg)) return _positions[account][position].usdg;
        if (token == address(weth)) return _positions[account][position].weth;
        uint256 index = _indexOfToken[token];
        return index == 0 ? 0 : _quantity(_units[account][position][index - 1], _books[index - 1].scale);
    }

    /// @notice The accounts' holding of `symbol`: the units the positions hold, each worth `scale` / 10^18 of a
    /// token, and the most the accounts may hold, in raw units.
    function holding(bytes32 symbol) external view returns (uint256 units, uint256 scale, uint256 cap) {
        uint256 i = _indexOfSymbol[symbol];
        if (i == 0) revert UnknownAsset(symbol);
        Book memory b = _books[--i];
        return (b.units, b.scale, _caps[i]);
    }

    /// @notice What `account`'s `position` owes the vault, in USDG: what a full repayment takes, rounded up and at most
    /// the vault's debt. Its premium is apart, in `premium`.
    function debt(address account, bytes32 position) public view returns (uint256) {
        uint256 total = vault.debt();
        return _min(_toAssetsUp(_positions[account][position].debtShares, total, totalDebtShares), total);
    }

    /// @notice The premium `account`'s `position` owes, in USDG, rounded down.
    function premium(address account, bytes32 position) public view returns (uint256) {
        Position memory p = _positions[account][position];
        return p.debtShares == 0 ? p.premium : _premium(p, _premiumIndexNow());
    }

    /// @notice The premium each position's debt carries while the session is closed, in basis points a year.
    function premiumRate() external view returns (uint32) {
        return _closure.premiumRate;
    }

    /// @notice The closure the premium accrues over, from `closesMs` until `reopensMs`, and when it last accrued, in
    /// milliseconds. A zero `closesMs` is no closure yet, and a zero `reopensMs` one whose reopening is not yet known.
    function closure() external view returns (uint64 closesMs, uint64 reopensMs, uint64 accruedMs) {
        Closure memory c = _closure;
        return (c.closesMs, c.reopensMs, c.accruedMs);
    }

    /// @notice The borrow shares of `account`'s `position`.
    function debtShares(address account, bytes32 position) external view returns (uint256) {
        return _positions[account][position].debtShares;
    }

    /// @notice `account`'s `position` as the engine sees it now: its equity, net of its debt and premium, and its
    /// requirement, in USD with 18 decimals, the engine's bits for assets whose liquidity is missing, and the regime. A
    /// halted asset counts for nothing.
    function health(address account, bytes32 position)
        external
        view
        returns (int256 equity, uint256 requirement, uint8 missing, uint8 regime)
    {
        Valuation memory v = _value(account, position, false);
        (requirement, missing, regime) = engine.currentRequirement(v.quantities, v.prices);
        return (v.equity, requirement, missing, regime);
    }

    /// @notice The gross exposure of `account`'s `position` to the engine's assets over its equity, in basis
    /// points; the largest value for a position without equity. `weekendLeverage` caps it before a closure.
    function leverage(address account, bytes32 position) external view returns (uint256) {
        Valuation memory v = _value(account, position, false);
        if (v.equity <= 0) return type(uint256).max;
        return v.gross * BPS / SafeCast.toUint256(v.equity);
    }

    /// @notice The price of `symbol`, in USD with 8 decimals, at or above which `account`'s `position` meets its
    /// current requirement once it owes `borrowing` more USDG, every other price and every pool where they are:
    /// the price at its band's low edge where the position falls short already, zero where no price leaves it
    /// short, and zero for a halted asset, which has no price. Found to within a 2^14th of the current price.
    function liquidationPrice(address account, bytes32 position, bytes32 symbol, uint256 borrowing)
        external
        view
        returns (uint256)
    {
        uint256 i = _indexOfSymbol[symbol];
        if (i == 0) revert UnknownAsset(symbol);
        --i;
        Valuation memory v = _value(account, position, false);
        uint256 quantity = SafeCast.toUint256(v.quantities[i]);
        uint256 high = v.prices[i];
        if (quantity == 0) return 0;
        int256 rest = v.equity - SafeCast.toInt256(quantity * high / PRICE_UNIT + borrowing * USDG_TO_USD);
        if (!_meets(v, rest, i, quantity, high)) return high;
        if (_meets(v, rest, i, quantity, 0)) return 0;
        uint256 low = 0;
        for (uint256 step; step < LIQUIDATION_PRICE_STEPS; ++step) {
            uint256 mid = (low + high) / 2;
            if (_meets(v, rest, i, quantity, mid)) high = mid;
            else low = mid;
        }
        return high;
    }

    /// @notice Allows or stops `authorized` borrowing, withdrawing and repaying with collateral for the caller, and
    /// depositing Stock Tokens and USDG into its positions.
    function setAuthorization(address authorized, bool allowed) external {
        isAuthorized[msg.sender][authorized] = allowed;
        emit AuthorizationSet(msg.sender, authorized, allowed);
    }

    /// @notice Sets the only address that may seize collateral and write debt off. It can be set once, and never to
    /// zero.
    function setLiquidator(address newLiquidator) external onlyOwner {
        if (liquidator != address(0)) revert LiquidatorAlreadySet(liquidator);
        if (newLiquidator == address(0)) revert InvalidLiquidator();
        liquidator = newLiquidator;
        emit LiquidatorSet(newLiquidator);
    }

    /// @notice Sets the only address that may claim the premium set aside. It can be set once, and never to zero.
    function setBackstop(address newBackstop) external onlyOwner {
        if (backstop != address(0)) revert BackstopAlreadySet(backstop);
        if (newBackstop == address(0)) revert InvalidBackstop();
        backstop = newBackstop;
        emit BackstopSet(newBackstop);
    }

    /// @notice Sets the premium rate, in basis points a year, from now on: closed time up to now accrues at the old
    /// rate. Only while the band can tell the session, so a spell it cannot is never priced at the new rate.
    function setPremiumRate(uint32 rate) external onlyOwner {
        if (rate > MAX_PREMIUM_RATE) revert InvalidPremium();
        // slither-disable-next-line unused-return
        (uint8 state,,,,) = band.session(); // forge-lint: disable-line(unused-return)
        if (state == 0) revert SessionUnknown();
        _accrue();
        _closure.premiumRate = rate;
        emit PremiumRateSet(rate);
    }

    /// @notice Accrues the premium up to now and records the closure the band's session shows. Anyone may call it;
    /// every loan, repayment, write-off, withdrawal in debt and rate change calls it first. A closure accrues exactly
    /// from its close to its reopening once its close has been recorded while ahead and its reopening while closed.
    /// Otherwise it accrues only time seen closed: from when it is first seen closed, until its last accrual while
    /// closed, and over a spell the band cannot tell the session only once the same closure is seen closed after it.
    /// A closed session is taken as the recorded closure only when it reopens within 72 hours of that closure's close.
    function accruePremium() external returns (uint256) {
        return _accrue();
    }

    /// @notice Sends the backstop the premium set aside for it, and returns how much. Only the backstop may.
    function claimPremium() external returns (uint256 amount) {
        if (msg.sender != backstop) revert NotBackstop(msg.sender);
        amount = backstopPremium;
        backstopPremium = 0;
        emit PremiumClaimed(msg.sender, amount);
        usdg.safeTransfer(msg.sender, amount);
    }

    /// @notice Sends `amount` of the reserve to `receiver`. Only the owner may.
    function withdrawReserve(address receiver, uint256 amount) external onlyOwner {
        if (receiver == address(0)) revert InvalidReceiver();
        if (amount > reserve) revert InsufficientReserve(reserve);
        reserve -= SafeCast.toUint128(amount);
        emit ReserveWithdrawn(receiver, amount);
        usdg.safeTransfer(receiver, amount);
    }

    /// @notice Always reverts: the accounts keep an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    /// @notice Sets the address that may pause new borrowing; zero leaves none, and the owner may set another.
    // forge-lint: disable-next-line(missing-zero-check)
    function setGuardian(address newGuardian) external onlyOwner {
        // slither-disable-next-line missing-zero-check
        guardian = newGuardian;
        emit GuardianSet(newGuardian);
    }

    /// @notice Pauses or resumes new borrowing, and withdrawals from positions in debt. Deposits, repayments,
    /// seizures and write-offs stay open. Only the guardian may.
    function setBorrowingPaused(bool paused) external {
        if (msg.sender != guardian) revert NotGuardian(msg.sender);
        borrowingPaused = paused;
        emit BorrowingPausedSet(paused);
    }

    /// @notice Brings every position's holding of `symbol` down in proportion when the accounts hold less of it than
    /// their positions count, as after the issuer burns from the accounts. It never raises a holding. Anyone may call
    /// it; every deposit, withdrawal, loan and seizure touching the asset calls it first.
    function sync(bytes32 symbol) external {
        uint256 i = _indexOfSymbol[symbol];
        if (i == 0 || _tokens[i - 1] == address(0)) revert UnknownAsset(symbol);
        _sync(i - 1);
    }

    /// @notice Clears `account`'s holding of `symbol` in `position` once a burn has left it worth less than a raw unit,
    /// so it neither gates the position nor keeps the asset closed to deposits. Anyone may call it.
    function clear(address account, bytes32 position, bytes32 symbol) external {
        uint256 i = _indexOfSymbol[symbol];
        if (i == 0 || _tokens[i - 1] == address(0)) revert UnknownAsset(symbol);
        Book storage b = _sync(--i);
        uint256 units = _units[account][position][i];
        if (_quantity(units, b.scale) != 0) revert HoldingNotEmpty(symbol);
        b.units -= SafeCast.toUint128(units);
        _units[account][position][i] = 0;
        Position storage p = _positions[account][position];
        p.held &= ~uint8(1 << i); // forge-lint: disable-line(unsafe-typecast)
        emit HoldingCleared(account, position, symbol, units);
    }

    /// @notice Deposits `amount` of `token` from the caller into `account`'s `position`. A Stock Token goes to the
    /// cross position or its own isolated position, and reverts with `AssetInOtherPosition` if the account holds it in
    /// the other. Stock Tokens and USDG come only from the account or an address it authorized, and a Stock Token only
    /// within its cap and for an account and caller the issuer has not blocklisted; anyone may deposit WETH.
    function deposit(bytes32 position, address token, uint256 amount, address account) external {
        if (amount == 0) revert ZeroAmount();
        _checkPosition(position);
        Position storage p = _positions[account][position];
        if (token == address(usdg)) {
            _checkAuthorized(account);
            p.usdg += SafeCast.toUint128(amount);
        } else if (token == address(weth)) {
            p.weth += SafeCast.toUint128(amount);
        } else {
            uint256 index = _indexOfToken[token];
            if (index == 0) revert UnsupportedToken(token, position);
            _checkAuthorized(account);
            _checkNotBlocked(--index, account);
            bytes32 symbol = _symbols[index];
            if (position != CROSS && position != symbol) revert UnsupportedToken(token, position);
            bytes32 other = position == CROSS ? symbol : CROSS;
            Book storage b = _sync(index);
            if (_quantity(_units[account][other][index], b.scale) != 0) revert AssetInOtherPosition(symbol, other);
            if (b.scale == 0) revert AssetWrittenOff(symbol);
            uint256 units = amount * WAD / b.scale;
            b.units += SafeCast.toUint128(units);
            if (_quantity(b.units, b.scale) > _caps[index]) revert AssetCapExceeded(symbol, _caps[index]);
            _units[account][position][index] += units;
            p.held |= uint8(1 << index); // forge-lint: disable-line(unsafe-typecast)
        }
        emit Deposit(msg.sender, account, position, token, amount);
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
    }

    /// @notice Withdraws `amount` of `token` from `account`'s `position` to `receiver`. A Stock Token leaves only for
    /// an account and caller the issuer has not blocklisted. A position in debt must then pass its checks, and waits while
    /// the guardian has paused borrowing.
    function withdraw(bytes32 position, address token, uint256 amount, address account, address receiver) external {
        _checkAuthorized(account);
        uint256 index = _indexOfToken[token];
        if (index != 0) _checkNotBlocked(index - 1, account);
        Position storage p = _debit(account, position, token, amount);
        if (p.debtShares != 0 || p.premium != 0) {
            if (borrowingPaused) revert BorrowingIsPaused();
            _accrue();
            _syncHeld(account, position);
            _check(account, position);
        }
        emit Withdraw(msg.sender, account, position, token, amount, receiver);
        IERC20(token).safeTransfer(receiver, amount);
    }

    /// @notice Borrows `assets` of USDG from the vault against `account`'s `position`, to `receiver`. The position
    /// must hold no USDG and must then pass its checks.
    function borrow(bytes32 position, uint256 assets, address account, address receiver) external {
        if (assets == 0) revert ZeroAmount();
        if (borrowingPaused) revert BorrowingIsPaused();
        _checkAuthorized(account);
        _checkPosition(position);
        Position storage p = _positions[account][position];
        if (p.usdg != 0) revert BorrowingAgainstUsdg();
        if (IUSDG(address(usdg)).isFrozen(address(this))) revert IUSDG.AddressFrozen();
        _settle(p, _accrue());
        _syncHeld(account, position);
        uint256 total = vault.debt();
        if (total + assets > debtCap) revert DebtCapExceeded(total + assets, debtCap);
        if (total + assets > weekendDebtCap && _weekendAhead()) {
            revert WeekendDebtCapExceeded(total + assets, weekendDebtCap);
        }
        uint256 shares = _toSharesUp(assets, total, totalDebtShares);
        p.debtShares += SafeCast.toUint128(shares);
        totalDebtShares += SafeCast.toUint128(shares);
        emit Borrow(msg.sender, account, position, assets, shares, receiver);
        vault.borrow(assets, receiver);
        _check(account, position);
    }

    /// @notice Repays up to `assets` of `account`'s `position` debt with the caller's USDG, then its premium with what
    /// is left, and returns what it paid. Anyone may.
    function repay(bytes32 position, uint256 assets, address account) external returns (uint256) {
        (uint256 repaid, uint256 premiumPaid) = _burn(account, position, _positions[account][position], assets);
        usdg.safeTransferFrom(msg.sender, address(this), repaid + premiumPaid);
        // slither-disable-next-line unused-return
        if (repaid != 0) vault.repay(repaid); // forge-lint: disable-line(unused-return)
        return repaid + premiumPaid;
    }

    /// @notice Repays up to `assets` of `account`'s `position` debt with the USDG the position holds, then its
    /// premium with what is left, and returns what it paid. It lowers collateral and debt alike, so it needs no check
    /// and stays open while an asset is halted.
    function repayWithCollateral(bytes32 position, uint256 assets, address account) external returns (uint256) {
        _checkAuthorized(account);
        Position storage p = _positions[account][position];
        (uint256 repaid, uint256 premiumPaid) = _burn(account, position, p, _min(assets, p.usdg));
        p.usdg -= SafeCast.toUint128(repaid + premiumPaid);
        // slither-disable-next-line unused-return
        if (repaid != 0) vault.repay(repaid); // forge-lint: disable-line(unused-return)
        return repaid + premiumPaid;
    }

    /// @notice Takes `amount` of `token` from `account`'s `position` to `receiver`, whatever the position's health.
    /// Only the liquidator may.
    function seize(bytes32 position, address token, uint256 amount, address account, address receiver) external {
        if (msg.sender != liquidator) revert NotLiquidator(msg.sender);
        _debit(account, position, token, amount);
        emit Seize(account, position, token, amount, receiver);
        IERC20(token).safeTransfer(receiver, amount);
    }

    /// @notice Writes off what `account`'s `position` owes once it holds no collateral, or only holdings a burn has
    /// left worth nothing, retiring its borrow shares: lenders bear it, and every other position owes what it did. The
    /// premium it owes is forgiven. Only the liquidator may.
    function writeOff(address account, bytes32 position) external returns (uint256 written) {
        if (msg.sender != liquidator) revert NotLiquidator(msg.sender);
        Position storage p = _positions[account][position];
        if (p.usdg != 0 || p.weth != 0) revert PositionNotEmpty();
        uint8 held = p.held;
        for (uint256 i = 0; held >> i != 0; ++i) {
            if (held & (1 << i) == 0) continue;
            Book storage book = _sync(i);
            uint256 units = _units[account][position][i];
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (_quantity(units, book.scale) != 0) revert PositionNotEmpty();
            book.units -= SafeCast.toUint128(units);
            _units[account][position][i] = 0;
        }
        p.held = 0;
        _settle(p, _accrue());
        uint256 total = vault.debt();
        uint256 shares = p.debtShares;
        written = _min(_toAssetsUp(shares, total, totalDebtShares), total);
        p.debtShares = 0;
        totalDebtShares -= SafeCast.toUint128(shares);
        emit WriteOff(account, position, written, shares, p.premium);
        p.premium = 0;
        // slither-disable-next-line unused-return
        vault.writeOff(written); // forge-lint: disable-line(unused-return)
    }

    function _burn(address account, bytes32 position, Position storage p, uint256 assets)
        private
        returns (uint256 repaid, uint256 premiumPaid)
    {
        _settle(p, _accrue());
        uint256 total = vault.debt();
        uint256 owed = _min(_toAssetsUp(p.debtShares, total, totalDebtShares), total);
        uint256 shares;
        if (assets >= owed) {
            (repaid, shares) = (owed, p.debtShares);
        } else {
            (repaid, shares) = (assets, _toSharesDown(assets, total, totalDebtShares));
        }
        p.debtShares -= SafeCast.toUint128(shares);
        totalDebtShares -= SafeCast.toUint128(shares);
        if (repaid != 0) emit Repay(msg.sender, account, position, repaid, shares);
        premiumPaid = _min(assets - repaid, p.premium);
        if (premiumPaid != 0) {
            p.premium -= SafeCast.toUint120(premiumPaid);
            uint256 toReserve = premiumPaid * reserveShare / BPS;
            reserve += SafeCast.toUint128(toReserve);
            backstopPremium += SafeCast.toUint128(premiumPaid - toReserve);
            emit PremiumPaid(account, position, premiumPaid, toReserve);
        }
    }

    function _accrue() private returns (uint256 index) {
        Closure memory c = _closure;
        index = premiumIndex;
        (uint64 closesMs, uint64 reopensMs) = (c.closesMs, c.reopensMs);
        uint256 closedMs = _advance(c, SafeCast.toUint64(block.timestamp * 1000));
        if (closedMs != 0 && totalDebtShares != 0) {
            index += _premiumPerShare(closedMs, c.premiumRate);
            premiumIndex = SafeCast.toUint128(index);
            emit PremiumAccrued(index);
        }
        bool moved = c.closesMs != closesMs || c.reopensMs != reopensMs;
        if (moved) emit ClosureSet(c.closesMs, c.reopensMs);
        if (moved || closedMs != 0) _closure = c;
    }

    function _premiumIndexNow() private view returns (uint256 index) {
        Closure memory c = _closure;
        index = premiumIndex;
        uint64 nowMs = SafeCast.toUint64(block.timestamp * 1000);
        if (c.accruedMs != nowMs && totalDebtShares != 0) {
            uint256 closedMs = _advance(c, nowMs);
            if (closedMs != 0) index += _premiumPerShare(closedMs, c.premiumRate);
        }
    }

    function _advance(Closure memory c, uint64 nowMs) private view returns (uint256 closedMs) {
        // slither-disable-next-line unused-return
        (uint8 state,,,, uint64 boundaryMs) = band.session(); // forge-lint: disable-line(unused-return)
        if (state != 1 && state != 2) return 0;
        bool reopened = c.reopensMs != 0 && c.reopensMs <= nowMs;
        bool inside = c.closesMs != 0 && c.closesMs <= nowMs && !reopened;
        bool same = inside && state == 1 && (boundaryMs != 0 ? boundaryMs : nowMs) <= c.closesMs + MAX_CLOSURE_MS;
        if (inside && !same) c.reopensMs = c.accruedMs > c.closesMs ? c.accruedMs : c.closesMs;
        else if (same && boundaryMs != 0) c.reopensMs = boundaryMs;
        if (c.closesMs != 0) {
            uint64 from = c.accruedMs > c.closesMs ? c.accruedMs : c.closesMs;
            uint64 to = c.reopensMs != 0 && c.reopensMs <= nowMs ? c.reopensMs : nowMs;
            if (to > from) closedMs = to - from;
        }
        if (state == 2 && boundaryMs != 0) (c.closesMs, c.reopensMs) = (boundaryMs, 0);
        else if (state == 1 && !same) (c.closesMs, c.reopensMs) = (nowMs, boundaryMs);
        c.accruedMs = nowMs;
    }

    function _premiumPerShare(uint256 closedMs, uint256 rate) private view returns (uint256) {
        return
            (vault.debt() + VIRTUAL_ASSETS) * rate * closedMs * RAY
                / (BPS * YEAR_MS * (totalDebtShares + VIRTUAL_SHARES));
    }

    function _settle(Position storage p, uint256 index) private {
        uint256 paid = p.premiumIndex;
        if (index != paid) {
            p.premium += SafeCast.toUint120(p.debtShares * (index - paid) / RAY);
            p.premiumIndex = index;
        }
    }

    function _premium(Position memory p, uint256 index) private pure returns (uint256) {
        return p.premium + p.debtShares * (index - p.premiumIndex) / RAY;
    }

    function _debit(address account, bytes32 position, address token, uint256 amount)
        private
        returns (Position storage p)
    {
        if (amount == 0) revert ZeroAmount();
        p = _positions[account][position];
        if (token == address(usdg)) {
            if (amount > p.usdg) revert InsufficientCollateral(token, amount);
            p.usdg -= SafeCast.toUint128(amount);
        } else if (token == address(weth)) {
            if (amount > p.weth) revert InsufficientCollateral(token, amount);
            p.weth -= SafeCast.toUint128(amount);
        } else {
            uint256 index = _indexOfToken[token];
            if (index == 0) revert UnsupportedToken(token, position);
            Book storage b = _sync(--index);
            uint256 held = _units[account][position][index];
            uint256 units = b.scale == 0 ? type(uint256).max : _divUp(amount * WAD, b.scale);
            if (units > held) revert InsufficientCollateral(token, amount);
            if (_quantity(held - units, b.scale) == 0) units = held;
            _units[account][position][index] = held - units;
            b.units -= SafeCast.toUint128(units);
            if (held == units) p.held &= ~uint8(1 << index); // forge-lint: disable-line(unsafe-typecast)
        }
    }

    function _syncHeld(address account, bytes32 position) private {
        uint8 bits = _positions[account][position].held;
        for (uint256 i = 0; bits >> i != 0; ++i) {
            if (bits & (1 << i) == 0) continue;
            Book storage b = _sync(i);
            if (_quantity(_units[account][position][i], b.scale) == 0) continue;
            _checkNotBlocked(i, account);
            IStockTokenRegistry registry = IStockTokenRegistry(_registries[i]);
            // forge-lint: disable-start(calls-loop)
            // slither-disable-next-line calls-loop
            if (IStockToken(_tokens[i]).paused() || registry.isBlocked(address(this))) {
                // forge-lint: disable-end(calls-loop)
                revert AssetFrozen(_symbols[i]); // forge-lint: disable-line(require-revert-in-loop)
            }
        }
    }

    function _check(address account, bytes32 position) private view {
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        Valuation memory v = _value(account, position, true);
        if (v.equity < 0) revert InsufficientMargin(v.equity, 0);
        uint8 held = v.held;
        if (held == 0) return;
        (uint256 requirement, uint8 missing, uint8 regime) = engine.currentRequirement(v.quantities, v.prices);
        if (regime == 0) revert SessionUnknown();
        uint8 unread = missing & held;
        if (unread != 0) revert LiquidityUnknown(_symbols[_lowestBit(unread)]);
        uint256 equity = SafeCast.toUint256(v.equity);
        if (equity < requirement) revert InsufficientMargin(v.equity, requirement);
        if (v.gross * BPS > equity * weekendLeverage && (regime != 2 || _weekendAhead())) {
            revert WeekendLeverageExceeded(v.gross, v.equity);
        }
    }

    function _weekendAhead() private view returns (bool) {
        // slither-disable-next-line unused-return
        (uint8 state,,,, uint64 boundaryMs) = band.session(); // forge-lint: disable-line(unused-return)
        return state != 2 || boundaryMs != 0;
    }

    function _sync(uint256 i) private returns (Book storage b) {
        b = _books[i];
        uint256 units = b.units;
        if (units == 0) {
            if (b.scale != WAD) b.scale = SafeCast.toUint128(WAD);
            return b;
        }
        uint256 counted = _quantity(units, b.scale);
        // slither-disable-next-line calls-loop
        uint256 balance = IERC20(_tokens[i]).balanceOf(address(this)); // forge-lint: disable-line(calls-loop)
        if (balance < counted) {
            b.scale = SafeCast.toUint128(balance * WAD / units);
            emit AssetWrittenDown(_symbols[i], counted, balance); // forge-lint: disable-line(reentrancy-events)
        }
    }

    function _quantity(uint256 units, uint256 scale) private pure returns (uint256) {
        return units * scale / WAD;
    }

    function _value(address account, bytes32 position, bool strict) private view returns (Valuation memory v) {
        Position memory p = _positions[account][position];
        uint256 n = _symbols.length;
        v.quantities = new int256[](n);
        v.prices = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            if (p.held & (1 << i) == 0) continue;
            uint256 quantity = _quantity(_units[account][position][i], _books[i].scale);
            if (quantity == 0) continue;
            bytes32 symbol = _symbols[i];
            // slither-disable-next-line unused-return,calls-loop
            (uint8 state,,,, uint64 low,) = band.quote(symbol); // forge-lint: disable-line(calls-loop, unused-return)
            if (strict) {
                if (state == 0) revert AssetHalted(symbol); // forge-lint: disable-line(require-revert-in-loop)
                // slither-disable-next-line unused-return,calls-loop
                (uint8 status,,,) = band.corporateAction(symbol); // forge-lint: disable-line(calls-loop, unused-return)
                // forge-lint: disable-next-line(require-revert-in-loop)
                if (status != 0) revert CorporateActionPending(symbol);
            }
            uint256 value = quantity * low / PRICE_UNIT;
            v.quantities[i] = SafeCast.toInt256(quantity);
            v.prices[i] = low;
            v.held |= uint8(1 << i); // forge-lint: disable-line(unsafe-typecast)
            v.gross += value;
        }
        uint256 owed = p.debtShares == 0
            ? p.premium
            : _toAssetsUp(p.debtShares, vault.debt(), totalDebtShares)
                + _premium(p, strict ? premiumIndex : _premiumIndexNow());
        v.equity = SafeCast.toInt256(v.gross + p.usdg * USDG_TO_USD + _wethValue(p.weth))
            - SafeCast.toInt256(owed * USDG_TO_USD);
    }

    function _meets(Valuation memory v, int256 rest, uint256 i, uint256 quantity, uint256 price)
        private
        view
        returns (bool)
    {
        v.prices[i] = price;
        // forge-lint: disable-start(calls-loop, unused-return)
        // slither-disable-next-line unused-return,calls-loop
        (uint256 requirement,,) = engine.currentRequirement(v.quantities, v.prices);
        // forge-lint: disable-end(calls-loop, unused-return)
        int256 equity = rest + SafeCast.toInt256(quantity * price / PRICE_UNIT);
        return equity >= 0 && SafeCast.toUint256(equity) >= requirement;
    }

    function _wethValue(uint256 amount) private view returns (uint256) {
        if (amount == 0 || address(ethUsd) == address(0)) return 0;
        // slither-disable-next-line unused-return
        try ethUsd.latestRoundData() returns (uint80, int256 answer, uint256, uint256 updatedAt, uint80) {
            // forge-lint: disable-next-line(block-timestamp)
            if (answer <= 0 || updatedAt > block.timestamp || block.timestamp - updatedAt > MAX_FEED_AGE) return 0;
            return amount * SafeCast.toUint256(answer) * WETH_COLLATERAL_BPS / (PRICE_UNIT * BPS);
        } catch {
            return 0;
        }
    }

    function _checkNotBlocked(uint256 i, address account) private view {
        IStockTokenRegistry registry = IStockTokenRegistry(_registries[i]);
        // forge-lint: disable-start(calls-loop, require-revert-in-loop)
        // slither-disable-next-line calls-loop
        if (registry.isBlocked(account)) revert Blocked(account);
        // forge-lint: disable-end(calls-loop, require-revert-in-loop)
        // forge-lint: disable-start(calls-loop, require-revert-in-loop)
        // slither-disable-next-line calls-loop
        if (msg.sender != account && registry.isBlocked(msg.sender)) revert Blocked(msg.sender);
        // forge-lint: disable-end(calls-loop, require-revert-in-loop)
    }

    function _checkAuthorized(address account) private view {
        if (msg.sender != account && !isAuthorized[account][msg.sender]) revert Unauthorized(msg.sender, account);
    }

    function _checkPosition(bytes32 position) private view {
        if (position == CROSS) return;
        uint256 index = _indexOfSymbol[position];
        if (index == 0 || _tokens[index - 1] == address(0)) revert UnknownPosition(position);
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

    function _lowestBit(uint8 bits) private pure returns (uint256 i) {
        while (bits & (1 << i) == 0) ++i;
    }
}
