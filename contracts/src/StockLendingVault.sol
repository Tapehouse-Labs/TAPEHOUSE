// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20, IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IStockToken, IStockTokenRegistry} from "./interfaces/IStockToken.sol";
import {SupplyVault} from "./SupplyVault.sol";

/// @title Tapehouse stock lending vault
/// @notice Lends one Stock Token, in its raw units, so dividends and splits, which move only the token's ERC-8056
/// multiplier, pass through every loan. Its one depositor, the margin accounts, set once, lends the tokens of the
/// positions that choose to; its one borrower, set once, borrows them open-term at a fee set by utilisation on a
/// two-slope curve, paid in the token. The owner takes `feeShare` of the fee as shares; lenders keep the rest. The
/// borrower may take at most `MAX_UTILIZATION` of the vault; what stays idle is lenders' to take back, first come,
/// first served.
/// @dev The vault counts its tokens itself, as those it holds plus what the borrower owes, so tokens sent to it
/// directly change no share's value. The debt grows by a borrow index with 27 decimals. The issuer can pause the
/// token, blocklist the vault and burn from it: while the token is paused, the vault blocklisted or its balance below
/// its count, every `max*` view returns zero, and `sync` brings the count down to a burnt balance, a loss to this
/// vault's lenders alone.
contract StockLendingVault is ERC4626, Ownable2Step {
    /// @notice The highest share of the fee the owner may take, in basis points.
    uint16 public constant MAX_FEE_SHARE = 20_00;
    /// @notice The most the borrower may take of the vault's tokens, in basis points.
    uint16 public constant MAX_UTILIZATION = 90_00;
    /// @notice The lowest optimal utilisation, in basis points.
    uint16 public constant MIN_OPTIMAL = 1_00;
    /// @notice The highest optimal utilisation, in basis points.
    uint16 public constant MAX_OPTIMAL = 99_00;
    /// @notice The highest fee the curve may reach, in basis points a year: 1,000%.
    uint32 public constant MAX_RATE = 1000_00;

    uint256 private constant BPS = 10_000;
    uint256 private constant WAD = 1e18;
    uint256 private constant RAY = 1e27;
    uint256 private constant WAD_PER_BPS = 1e14;
    uint256 private constant YEAR = 365 days;

    /// @notice The owner's share of the fee, in basis points.
    uint16 public immutable feeShare;
    IStockTokenRegistry private immutable _registry;

    /// @notice The only address that may deposit: the margin accounts.
    address public depositor;
    /// @notice The only address that may borrow.
    address public borrower;
    /// @notice The tokens the vault holds, by its own count.
    uint128 public idle;
    /// @notice The borrower's debt divided by the borrow index at which it was taken.
    uint128 public scaledDebt;
    /// @notice What one unit of debt taken at the start is owed now, with 27 decimals, as of the last accrual.
    uint128 public borrowIndex;
    /// @notice When the fee last accrued.
    uint64 public lastAccrual;
    SupplyVault.RateModel private _rateModel;

    /// @notice The depositor was set.
    event DepositorSet(address indexed depositor);
    /// @notice The borrower was set.
    event BorrowerSet(address indexed borrower);
    /// @notice The rate model was set.
    event RateModelSet(uint16 optimal, uint32 base, uint32 slope1, uint32 slope2);
    /// @notice A fee of `fee` tokens accrued on the borrower's debt, which the index now reaches; the owner took
    /// `shares` for its share of it.
    event FeeAccrued(uint256 fee, uint256 index, uint256 shares);
    /// @notice The borrower borrowed `assets` to `receiver`.
    event Borrow(address indexed receiver, uint256 assets);
    /// @notice The borrower repaid `assets` of its debt.
    event Repay(uint256 assets);
    /// @notice The vault's count of the tokens it holds fell to its balance.
    event Sync(uint256 idle);

    /// @notice `account` is not the depositor.
    error NotDepositor(address account);
    /// @notice `account` is not the borrower.
    error NotBorrower(address account);
    /// @notice The depositor or the borrower is already set, to `current`.
    error AlreadySet(address current);
    /// @notice The depositor or the borrower cannot be the zero address.
    error InvalidAddress();
    /// @notice Lending `requested` would take the borrower past `MAX_UTILIZATION`; it may borrow `available`.
    error InsufficientLiquidity(uint256 requested, uint256 available);
    /// @notice The rate model is outside its bounds.
    error InvalidRateModel();
    /// @notice The owner's share of the fee is above `MAX_FEE_SHARE`.
    error InvalidFeeShare();
    /// @notice The owner cannot renounce: the vault needs one to set its depositor, borrower and rate model.
    error OwnershipCannotBeRenounced();

    /// @param token The Stock Token lent.
    /// @param initialOwner The owner, who sets the depositor and the borrower once and the rate model within its
    /// bounds, and receives its share of the fee.
    /// @param initialRateModel The first rate model.
    /// @param feeShare_ The owner's share of the fee, in basis points, at most `MAX_FEE_SHARE`.
    constructor(IERC20 token, address initialOwner, SupplyVault.RateModel memory initialRateModel, uint16 feeShare_)
        ERC20(
            string.concat("Tapehouse ", IERC20Metadata(address(token)).symbol(), " Lending"),
            string.concat("thl", IERC20Metadata(address(token)).symbol())
        )
        ERC4626(token)
        Ownable(initialOwner)
    {
        if (feeShare_ > MAX_FEE_SHARE) revert InvalidFeeShare();
        feeShare = feeShare_;
        _registry = IStockTokenRegistry(IStockToken(address(token)).ACCESS_CONTROLLED_REGISTRY());
        lastAccrual = SafeCast.toUint64(block.timestamp);
        borrowIndex = SafeCast.toUint128(RAY);
        _setRateModel(initialRateModel);
    }

    /// @notice The rate model.
    function rateModel() external view returns (SupplyVault.RateModel memory) {
        return _rateModel;
    }

    /// @notice Always reverts: the vault keeps an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    /// @notice Sets the only address that may deposit. It can be set once, and never to zero.
    function setDepositor(address newDepositor) external onlyOwner {
        if (depositor != address(0)) revert AlreadySet(depositor);
        if (newDepositor == address(0)) revert InvalidAddress();
        depositor = newDepositor;
        emit DepositorSet(newDepositor);
    }

    /// @notice Sets the only address that may borrow. It can be set once, and never to zero.
    function setBorrower(address newBorrower) external onlyOwner {
        if (borrower != address(0)) revert AlreadySet(borrower);
        if (newBorrower == address(0)) revert InvalidAddress();
        borrower = newBorrower;
        emit BorrowerSet(newBorrower);
    }

    /// @notice Sets the rate model, within its bounds, at once. The fee accrues first at the old rate.
    function setRateModel(SupplyVault.RateModel calldata model) external onlyOwner {
        _accrue();
        _setRateModel(model);
    }

    /// @notice Lends `assets` of the vault's tokens to `receiver`, within `MAX_UTILIZATION`. Only the borrower may.
    function borrow(uint256 assets, address receiver) external {
        if (msg.sender != borrower) revert NotBorrower(msg.sender);
        uint256 index = _accrue();
        uint256 available = borrowable();
        if (assets > available) revert InsufficientLiquidity(assets, available);
        scaledDebt += SafeCast.toUint128(_divUp(assets * RAY, index));
        // forge-lint: disable-next-line(reentrancy-events)
        emit Borrow(receiver, assets);
        _transferOut(receiver, assets);
    }

    /// @notice Repays up to `assets` of the debt from the borrower, and returns what it repaid. Only the borrower may.
    function repay(uint256 assets) external returns (uint256 repaid) {
        if (msg.sender != borrower) revert NotBorrower(msg.sender);
        repaid = _reduceDebt(assets);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Repay(repaid);
        _transferIn(msg.sender, repaid);
    }

    /// @notice Brings the vault's count of the tokens it holds down to its balance, as after the issuer burns from
    /// the vault: a loss to this vault's lenders alone. It never counts tokens sent to the vault directly. Anyone may
    /// call it.
    function sync() external {
        uint256 balance = IERC20(asset()).balanceOf(address(this));
        if (balance < idle) {
            _accrue();
            idle = SafeCast.toUint128(balance);
            // forge-lint: disable-next-line(reentrancy-events)
            emit Sync(balance);
        }
    }

    /// @notice What the borrower owes now, the fee included, rounded up.
    function debt() public view returns (uint256) {
        return _divUp(scaledDebt * _currentIndex(), RAY);
    }

    /// @notice What the borrower may still borrow: the vault's tokens up to `MAX_UTILIZATION` of its assets.
    function borrowable() public view returns (uint256) {
        uint256 owed = debt();
        uint256 cap = (owed + idle) * MAX_UTILIZATION / BPS;
        return cap > owed ? cap - owed : 0;
    }

    /// @notice The share of the vault's assets lent, with 18 decimals.
    function utilization() public view returns (uint256) {
        return _utilization(debt(), idle);
    }

    /// @notice The fee the borrower pays now, a year, with 18 decimals.
    function borrowRate() public view returns (uint256) {
        return _borrowRate(utilization());
    }

    /// @notice The rate lenders earn now, a year, with 18 decimals: the fee on the share lent, less the owner's share.
    function supplyRate() external view returns (uint256) {
        uint256 u = utilization();
        return _borrowRate(u) * u * (BPS - feeShare) / (WAD * BPS);
    }

    /// @inheritdoc ERC4626
    function totalAssets() public view override returns (uint256) {
        return idle + debt();
    }

    /// @inheritdoc ERC4626
    /// @dev Zero but for the depositor, and while the token is paused, the vault blocklisted or its tokens burnt.
    function maxDeposit(address receiver) public view override returns (uint256) {
        return receiver != depositor || _unavailable() ? 0 : super.maxDeposit(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev Zero but for the depositor, and while the token is paused, the vault blocklisted or its tokens burnt.
    function maxMint(address receiver) public view override returns (uint256) {
        return receiver != depositor || _unavailable() ? 0 : super.maxMint(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev Capped by the tokens the vault holds; zero while the token is paused, the vault blocklisted or its
    /// tokens burnt.
    function maxRedeem(address owner) public view override returns (uint256) {
        if (_unavailable()) return 0;
        uint256 shares = super.maxRedeem(owner);
        uint256 liquid = _convertToShares(idle, Math.Rounding.Floor);
        return shares < liquid ? shares : liquid;
    }

    function _convertToShares(uint256 assets, Math.Rounding rounding) internal view override returns (uint256) {
        return
            Math.mulDiv(
                assets, totalSupply() + _pendingFeeShares() + 10 ** _decimalsOffset(), totalAssets() + 1, rounding
            );
    }

    function _convertToAssets(uint256 shares, Math.Rounding rounding) internal view override returns (uint256) {
        return
            Math.mulDiv(
                shares, totalAssets() + 1, totalSupply() + _pendingFeeShares() + 10 ** _decimalsOffset(), rounding
            );
    }

    function _deposit(address caller, address receiver, uint256 assets, uint256 shares) internal override {
        if (caller != depositor) revert NotDepositor(caller);
        _accrue();
        super._deposit(caller, receiver, assets, shares);
    }

    function _withdraw(address caller, address receiver, address owner, uint256 assets, uint256 shares)
        internal
        override
    {
        _accrue();
        super._withdraw(caller, receiver, owner, assets, shares);
    }

    function _transferIn(address from, uint256 assets) internal override {
        idle += SafeCast.toUint128(assets); // forge-lint: disable-line(missing-events-arithmetic)
        super._transferIn(from, assets);
    }

    function _transferOut(address to, uint256 assets) internal override {
        idle -= SafeCast.toUint128(assets); // forge-lint: disable-line(missing-events-arithmetic)
        super._transferOut(to, assets);
    }

    function _reduceDebt(uint256 assets) private returns (uint256 reduced) {
        uint256 index = _accrue();
        uint256 owed = _divUp(scaledDebt * index, RAY);
        if (assets >= owed) {
            scaledDebt = 0;
            return owed;
        }
        scaledDebt -= SafeCast.toUint128(assets * RAY / index);
        return assets;
    }

    function _accrue() private returns (uint256 index) {
        index = _currentIndex();
        lastAccrual = SafeCast.toUint64(block.timestamp);
        uint256 previous = borrowIndex;
        // forge-lint: disable-next-line(block-timestamp)
        if (index <= previous) return index;
        borrowIndex = SafeCast.toUint128(index);
        uint256 fee = scaledDebt * (index - previous) / RAY;
        uint256 owed = scaledDebt * (index - previous) * feeShare / (RAY * BPS);
        uint256 shares = _feeShares(owed);
        // forge-lint: disable-next-line(reentrancy-events)
        emit FeeAccrued(fee, index, shares);
        // forge-lint: disable-next-line(block-timestamp)
        if (shares != 0) _mint(owner(), shares);
    }

    function _pendingFeeShares() private view returns (uint256) {
        uint256 previous = borrowIndex;
        uint256 index = _currentIndex();
        // forge-lint: disable-next-line(block-timestamp)
        if (index <= previous) return 0;
        return _feeShares(scaledDebt * (index - previous) * feeShare / (RAY * BPS));
    }

    function _feeShares(uint256 owed) private view returns (uint256) {
        return Math.mulDiv(owed, totalSupply() + 10 ** _decimalsOffset(), totalAssets() - owed + 1);
    }

    function _currentIndex() private view returns (uint256) {
        uint256 index = borrowIndex;
        // slither-disable-next-line incorrect-equality
        if (scaledDebt == 0) return index;
        uint256 lent = _divUp(scaledDebt * index, RAY);
        uint256 rate = _borrowRate(_utilization(lent, idle));
        return index + index * rate * (block.timestamp - lastAccrual) / (WAD * YEAR);
    }

    function _setRateModel(SupplyVault.RateModel memory model) private {
        if (
            model.optimal < MIN_OPTIMAL || model.optimal > MAX_OPTIMAL || model.slope1 > model.slope2
                || uint256(model.base) + model.slope1 + model.slope2 > MAX_RATE
        ) revert InvalidRateModel();
        _rateModel = model;
        // forge-lint: disable-next-line(reentrancy-events)
        emit RateModelSet(model.optimal, model.base, model.slope1, model.slope2);
    }

    function _unavailable() private view returns (bool) {
        try IStockToken(asset()).paused() returns (bool paused) {
            if (paused) return true;
        } catch {
            return true;
        }
        try _registry.isBlocked(address(this)) returns (bool blocked) {
            if (blocked) return true;
        } catch {
            return true;
        }
        try IERC20(asset()).balanceOf(address(this)) returns (uint256 balance) {
            return balance < idle;
        } catch {
            return true;
        }
    }

    function _divUp(uint256 a, uint256 b) private pure returns (uint256) {
        return a > 0 ? (a - 1) / b + 1 : 0;
    }

    function _utilization(uint256 lent, uint256 held) private pure returns (uint256) {
        return lent > 0 ? lent * WAD / (lent + held) : 0;
    }

    function _borrowRate(uint256 u) private view returns (uint256) {
        SupplyVault.RateModel memory m = _rateModel;
        uint256 optimal = m.optimal * WAD_PER_BPS;
        uint256 rate = m.base * WAD_PER_BPS;
        if (u <= optimal) return rate + m.slope1 * WAD_PER_BPS * u / optimal;
        return rate + m.slope1 * WAD_PER_BPS + m.slope2 * WAD_PER_BPS * (u - optimal) / (WAD - optimal);
    }
}
