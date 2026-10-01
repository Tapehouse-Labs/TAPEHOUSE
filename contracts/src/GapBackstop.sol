// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IUSDG} from "./interfaces/IUSDG.sol";
import {Liquidator} from "./Liquidator.sol";
import {MarginAccounts} from "./MarginAccounts.sol";

/// @title Tapehouse gap backstop
/// @notice The junior tranche behind the margin accounts' lenders. Depositors lend it USDG and receive shares; it
/// takes the weekend premium the accounts set aside for it, and covers what a position emptied by liquidation still
/// owes, before any of it is written off to lenders. It covers at most a declared limit per closure for each kind of
/// position: the cross positions together, and each Stock Token's isolated positions. Within the same limits it buys,
/// at the reopening auction's clearing price, what the auction's bids leave, and shares the Stock Tokens it takes among
/// its shares as gains each holder claims. A change to a limit applies only to closures that begin after every
/// depositor could have left. A depositor leaves by starting a cooldown: its shares may be redeemed from a week later,
/// for six days, never while the market is closed nor until a day after a closure's reopening.
/// @dev The backstop counts the USDG it holds itself, so USDG sent to it directly changes no share's value.
contract GapBackstop is ERC4626, Ownable2Step {
    using SafeERC20 for IERC20;

    struct Cooldown {
        uint192 shares;
        uint64 startedAt;
    }

    struct Limit {
        uint128 current;
        uint128 next;
        uint64 fromMs;
    }

    /// @notice How long after a depositor starts its cooldown its shares may be redeemed, in seconds.
    uint256 public constant COOLDOWN = 7 days;
    /// @notice How long, once the cooldown ends, its shares may be redeemed, in seconds.
    uint256 public constant WITHDRAWAL_WINDOW = 6 days;
    /// @notice How long after a closure's reopening redemptions stay shut, so its shortfalls are covered first, in
    /// seconds.
    uint256 public constant SETTLEMENT = 1 days;
    /// @notice How far apart two recorded closes must be to count as two closures, in seconds: a closure longer than
    /// this is recorded again as it goes on.
    uint256 public constant MIN_CLOSURE_GAP = 72 hours;

    bytes32 private constant CROSS = bytes32(0);
    uint256 private constant GAINS_SCALE = 1e36;
    uint256 private constant TOKEN_TO_USDG = 1e20;
    uint256 private constant MS = 1000;
    uint256 private constant LONGEST_CLOSURE = 96 hours;

    /// @notice The accounts this vault backstops.
    MarginAccounts public immutable accounts;
    IUSDG private immutable _usdg;

    /// @notice The USDG the backstop holds, by its own count.
    uint256 public held;
    /// @notice The close, in milliseconds, of the closure the backstop last counted a cover against.
    uint64 public closureMs;

    /// @notice What the backstop has covered for `position` in the closure it counts from `closesMs`, in USDG.
    mapping(bytes32 position => mapping(uint64 closesMs => uint256)) public covered;
    /// @notice Each depositor's cooldown: the shares it may redeem and when it started.
    mapping(address owner => Cooldown) public cooldowns;
    /// @notice The exposure limit of `position` for closures before `fromMs`, and for those from it on, in USDG.
    mapping(bytes32 position => Limit) public exposureLimits;
    /// @notice The units of each Stock Token taken for one share so far, scaled by 1e36.
    mapping(address token => uint256) public gainsPerShare;
    /// @notice How many times the backstop has taken Stock Tokens.
    uint256 public gainsEpoch;
    address[] private _tokens;
    mapping(address owner => uint256) private _seenEpoch;
    mapping(address owner => mapping(address token => uint256)) private _seenGains;
    mapping(address owner => mapping(address token => uint256)) private _owedGains;

    /// @notice The exposure limit of `position` is `limit` USDG a closure, for closures that close from `fromMs` on.
    event ExposureLimitSet(bytes32 indexed position, uint256 limit, uint64 fromMs);
    /// @notice `amount` of premium was claimed for the shares.
    event PremiumAdded(uint256 amount);
    /// @notice The backstop repaid `paid` of what `account`'s emptied `position` owed, in the closure it counts from
    /// `closesMs`, and `written` was written off to lenders.
    event Covered(
        address indexed account, bytes32 indexed position, uint64 indexed closesMs, uint256 paid, uint256 written
    );
    /// @notice `owner` started a cooldown for `shares`, redeemable from `from`.
    event CooldownStarted(address indexed owner, uint256 shares, uint64 from);
    /// @notice The backstop's count of the USDG it holds fell to its balance.
    event Sync(uint256 held);
    /// @notice The backstop bought `amount` of `token` from `account`'s `position` for `cost`, what the reopening
    /// auction's bids left, in the closure it counts from `closesMs`.
    event RemainderBought(
        address indexed account,
        bytes32 indexed position,
        address indexed token,
        uint256 amount,
        uint256 cost,
        uint64 closesMs
    );
    /// @notice `owner` claimed `amount` of `token` it gained.
    event GainsClaimed(address indexed owner, address indexed token, uint256 amount);

    /// @notice `position` is neither the cross position nor an asset of the accounts.
    error InvalidPosition(bytes32 position);
    /// @notice The exposure limits do not match the accounts' assets.
    error InvalidExposureLimits();
    /// @notice The owner cannot renounce: the backstop needs one to set its exposure limits.
    error OwnershipCannotBeRenounced();
    /// @notice Only the liquidator's reopening auction may.
    error NotAuction(address caller);

    /// @param accounts_ The accounts, whose owner must then set this vault as their backstop.
    /// @param initialOwner The owner, who sets the exposure limits.
    /// @param crossLimit The cross positions' exposure limit, in USDG a closure.
    /// @param assetLimits Each asset's isolated positions' exposure limit, in the accounts' order of assets.
    constructor(MarginAccounts accounts_, address initialOwner, uint256 crossLimit, uint256[] memory assetLimits)
        ERC20("Tapehouse Gap Backstop", "thGAP")
        ERC4626(accounts_.usdg())
        Ownable(initialOwner)
    {
        accounts = accounts_;
        _usdg = IUSDG(address(accounts_.usdg()));
        (bytes32[] memory symbols, address[] memory tokens) = accounts_.stocks();
        if (assetLimits.length != symbols.length) revert InvalidExposureLimits();
        _initLimit(CROSS, crossLimit);
        for (uint256 i; i < symbols.length; ++i) {
            _initLimit(symbols[i], assetLimits[i]);
            _tokens.push(tokens[i]);
        }
        IERC20(asset()).forceApprove(address(accounts_), type(uint256).max);
    }

    /// @notice Always reverts: the backstop keeps an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    /// @notice Sets the exposure limit of `position`, `bytes32(0)` for the cross positions or an asset's symbol for
    /// its isolated positions, in USDG a closure. It applies to closures that close `COOLDOWN + WITHDRAWAL_WINDOW` or
    /// more from now, when every depositor could have left; a later change before then replaces it.
    function setExposureLimit(bytes32 position, uint256 limit) external onlyOwner {
        if (position != CROSS) {
            // slither-disable-next-line unused-return
            (bytes32[] memory symbols,) = accounts.stocks(); // forge-lint: disable-line(unused-return)
            uint256 i = 0;
            while (i < symbols.length && symbols[i] != position) ++i;
            if (i == symbols.length) revert InvalidPosition(position);
        }
        Limit storage l = exposureLimits[position];
        (l.current, l.next) = (SafeCast.toUint128(_limitAt(l, _closure())), SafeCast.toUint128(limit));
        l.fromMs = SafeCast.toUint64((block.timestamp + COOLDOWN + WITHDRAWAL_WINDOW) * MS);
        emit ExposureLimitSet(position, limit, l.fromMs);
    }

    /// @notice Claims the premium the accounts set aside for the backstop, for the shares. Anyone may call it.
    function claim() external returns (uint256 amount) {
        amount = accounts.claimPremium();
        held += amount;
        // forge-lint: disable-next-line(reentrancy-events)
        emit PremiumAdded(amount);
    }

    /// @notice Repays what `account`'s emptied `position` owes, within what is left of its exposure limit in the
    /// current closure and of the USDG the backstop holds, then has the liquidator write off the rest to lenders. The
    /// position's premium is forgiven. While the issuer freezes the backstop it repays nothing, so lenders bear it all.
    /// Reverts unless the position holds nothing, or only holdings a burn left worthless. Anyone may call it.
    function cover(address account, bytes32 position) external returns (uint256 paid, uint256 written) {
        // slither-disable-next-line unused-return
        accounts.accruePremium(); // forge-lint: disable-line(unused-return)
        uint64 key = _closure();
        if (key != closureMs) closureMs = key;
        uint256 spent = covered[position][key];
        if (!_frozen()) {
            paid = Math.min(accounts.debt(account, position), _room(_limitAt(exposureLimits[position], key), spent));
        }
        if (paid != 0) {
            covered[position][key] = spent + paid;
            held -= paid;
            // slither-disable-next-line unused-return
            accounts.repay(position, paid, account); // forge-lint: disable-line(unused-return)
        }
        written = Liquidator(accounts.liquidator()).writeOff(account, position);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Covered(account, position, key, paid, written);
    }

    /// @notice Buys up to `amount` of `token` from `account`'s `position` at `price`, in USD with 8 decimals per token,
    /// what the reopening auction's bids left, within what is left of the position's exposure limit in the current
    /// closure and of the USDG the backstop holds, and returns what it took. The tokens are shared among the shares as
    /// gains. The auction offers at most the band's low edge as `price`, and leaves the rest in the position when this
    /// reverts: while the issuer freezes the backstop or blocks it from the Stock Token, or once it has no shares. Only
    /// the liquidator's reopening auction may.
    // slither-disable-next-line reentrancy-no-eth
    function buyRemainder(address account, bytes32 position, address token, uint256 amount, uint256 price)
        external
        returns (uint256 taken)
    {
        Liquidator liquidator = Liquidator(accounts.liquidator());
        if (msg.sender != liquidator.auction()) revert NotAuction(msg.sender);
        uint64 key = _closure();
        if (key != closureMs) closureMs = key;
        uint256 spent = covered[position][key];
        uint256 budget = _room(_limitAt(exposureLimits[position], key), spent);
        uint256 units = Math.min(amount, budget * TOKEN_TO_USDG / price);
        // slither-disable-next-line incorrect-equality
        if (units == 0) return 0;
        IERC20(asset()).forceApprove(address(liquidator), budget);
        uint256 cost;
        // forge-lint: disable-next-line(reentrancy-no-eth)
        (taken, cost) = liquidator.settle(account, position, token, units, price);
        IERC20(asset()).forceApprove(address(liquidator), 0);
        // slither-disable-next-line incorrect-equality
        if (taken == 0) return 0;
        covered[position][key] = spent + cost;
        held -= cost;
        gainsPerShare[token] += taken * GAINS_SCALE / totalSupply();
        ++gainsEpoch;
        // forge-lint: disable-next-line(reentrancy-events)
        emit RemainderBought(account, position, token, taken, cost, key);
    }

    /// @notice Sends the caller the `token` it gained from the Stock Tokens the backstop bought.
    function claimGains(address token) external returns (uint256 amount) {
        _checkpoint(msg.sender);
        amount = _owedGains[msg.sender][token];
        _owedGains[msg.sender][token] = 0;
        emit GainsClaimed(msg.sender, token, amount);
        if (amount != 0) IERC20(token).safeTransfer(msg.sender, amount);
    }

    /// @notice The `token` that `owner` has gained and not claimed.
    function gains(address owner, address token) external view returns (uint256) {
        return
            _owedGains[owner][token] + balanceOf(owner) * (gainsPerShare[token] - _seenGains[owner][token])
                / GAINS_SCALE;
    }

    /// @notice Starts the caller's cooldown for all its shares: they may be redeemed from `COOLDOWN` later, for
    /// `WITHDRAWAL_WINDOW`, while redemptions are open. A new cooldown replaces the last. Shares sent away leave it.
    function startCooldown() external {
        uint256 shares = balanceOf(msg.sender);
        cooldowns[msg.sender] = Cooldown(SafeCast.toUint192(shares), SafeCast.toUint64(block.timestamp));
        emit CooldownStarted(msg.sender, shares, SafeCast.toUint64(block.timestamp + COOLDOWN));
    }

    /// @notice Brings the backstop's count of the USDG it holds down to its balance, as after the issuer wipes a
    /// frozen backstop: a loss to the shares. It never counts USDG sent to the backstop directly. Anyone may call it.
    function sync() external {
        uint256 balance = _usdg.balanceOf(address(this));
        if (balance < held) {
            held = balance;
            emit Sync(balance);
        }
    }

    /// @notice The exposure limit of `position` in the current closure, in USDG.
    function exposureLimit(bytes32 position) external view returns (uint256) {
        return _limitAt(exposureLimits[position], _closure());
    }

    /// @notice What the backstop may still cover for `position` in the current closure: the rest of its exposure limit,
    /// at most the USDG the backstop holds.
    function exposureLeft(bytes32 position) external view returns (uint256) {
        uint64 key = _closure();
        return _room(_limitAt(exposureLimits[position], key), covered[position][key]);
    }

    /// @inheritdoc ERC4626
    function totalAssets() public view override returns (uint256) {
        return held;
    }

    /// @inheritdoc ERC4626
    /// @dev Zero until the accounts name this vault their backstop, while USDG is paused, the backstop frozen or its
    /// USDG wiped, and from the moment the accounts record a coming close until a day after its reopening.
    function maxDeposit(address receiver) public view override returns (uint256) {
        return _closedToDeposits() ? 0 : super.maxDeposit(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev Zero until the accounts name this vault their backstop, while USDG is paused, the backstop frozen or its
    /// USDG wiped, and from the moment the accounts record a coming close until a day after its reopening.
    function maxMint(address receiver) public view override returns (uint256) {
        return _closedToDeposits() ? 0 : super.maxMint(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev The shares of `owner`'s cooldown, from its end for `WITHDRAWAL_WINDOW`; zero while the market is closed,
    /// once the accounts record a coming close until `SETTLEMENT` after its reopening, and while USDG is paused, the
    /// backstop frozen or its USDG wiped.
    function maxRedeem(address owner) public view override returns (uint256) {
        Cooldown memory c = cooldowns[owner];
        uint256 from = uint256(c.startedAt) + COOLDOWN;
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp < from || block.timestamp > from + WITHDRAWAL_WINDOW) return 0;
        if (_unavailable() || _closed() || _settling()) return 0;
        return c.shares;
    }

    function _withdraw(address caller, address receiver, address owner, uint256 assets, uint256 shares)
        internal
        override
    {
        cooldowns[owner].shares -= SafeCast.toUint192(shares);
        super._withdraw(caller, receiver, owner, assets, shares);
    }

    function _update(address from, address to, uint256 value) internal override {
        _checkpoint(from);
        _checkpoint(to);
        super._update(from, to, value);
        if (from == address(0)) return;
        uint256 balance = balanceOf(from);
        if (cooldowns[from].shares > balance) cooldowns[from].shares = SafeCast.toUint192(balance);
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

    function _checkpoint(address owner) private {
        uint256 epoch = gainsEpoch;
        if (owner == address(0) || _seenEpoch[owner] == epoch) return;
        _seenEpoch[owner] = epoch;
        uint256 balance = balanceOf(owner);
        uint256 count = _tokens.length;
        for (uint256 i; i < count; ++i) {
            address token = _tokens[i];
            uint256 perShare = gainsPerShare[token];
            uint256 seen = _seenGains[owner][token];
            if (perShare == seen) continue;
            _owedGains[owner][token] += balance * (perShare - seen) / GAINS_SCALE;
            _seenGains[owner][token] = perShare;
        }
    }

    function _initLimit(bytes32 position, uint256 limit) private {
        uint128 value = SafeCast.toUint128(limit);
        exposureLimits[position] = Limit(value, value, 0);
        emit ExposureLimitSet(position, limit, 0);
    }

    function _limitAt(Limit memory l, uint64 key) private pure returns (uint256) {
        return key >= l.fromMs ? l.next : l.current;
    }

    function _room(uint256 limit, uint256 spent) private view returns (uint256) {
        return limit > spent ? Math.min(limit - spent, held) : 0;
    }

    function _closure() private view returns (uint64) {
        uint64 startMs = accounts.closureStart();
        uint64 last = closureMs;
        return startMs >= last + MIN_CLOSURE_GAP * MS ? startMs : last;
    }

    function _closedToDeposits() private view returns (bool) {
        try accounts.backstop() returns (address backstop) {
            return backstop != address(this) || _unavailable() || _settling();
        } catch {
            return true;
        }
    }

    function _frozen() private view returns (bool) {
        try _usdg.isFrozen(address(this)) returns (bool frozen) {
            return frozen;
        } catch {
            return true;
        }
    }

    function _unavailable() private view returns (bool) {
        try _usdg.paused() returns (bool paused) {
            if (paused) return true;
        } catch {
            return true;
        }
        if (_frozen()) return true;
        try _usdg.balanceOf(address(this)) returns (uint256 balance) {
            return balance < held;
        } catch {
            return true;
        }
    }

    function _settling() private view returns (bool) {
        // slither-disable-next-line unused-return
        try accounts.closure() returns (uint64 closesMs, uint64 reopensMs, uint64) {
            uint256 nowMs = block.timestamp * MS;
            uint256 until = reopensMs > closesMs ? reopensMs + SETTLEMENT * MS : closesMs + LONGEST_CLOSURE * MS;
            // forge-lint: disable-next-line(block-timestamp)
            return nowMs < until;
        } catch {
            return true;
        }
    }

    function _closed() private view returns (bool) {
        // slither-disable-next-line unused-return
        try accounts.band().session() returns (uint8 state, uint8 nyse, uint8 nyseNext, uint64, uint64 boundaryMs) {
            // forge-lint: disable-next-line(block-timestamp)
            if (state != 2 || (boundaryMs != 0 && boundaryMs <= block.timestamp * MS)) return true;
            return nyse != 1 && (nyse != 2 || nyseNext != 1);
        } catch {
            return true;
        }
    }
}
