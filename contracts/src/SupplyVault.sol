// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC4626} from "@openzeppelin/contracts/token/ERC20/extensions/ERC4626.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IUSDG} from "./interfaces/IUSDG.sol";

/// @title Tapehouse USDG supply vault
/// @notice Lenders deposit USDG and receive shares; one borrower, set once by the owner, borrows the USDG the
/// vault holds and pays interest at a rate set by utilisation. Lenders bear, last, any debt the borrower writes
/// off.
/// @dev The vault counts its assets itself, as the USDG it holds plus what the borrower owes, so USDG sent to
/// it directly changes no share's value. The debt grows by a borrow index with 27 decimals, so interest
/// compounds at every accrual without rounding away. USDG's issuer can pause it or freeze the vault, and then
/// wipe a frozen address's balance: while USDG is paused, the vault frozen, or its balance below its count, every
/// `max*` view returns zero, and `sync` brings the vault's count down to a wiped balance. Deposits open once the
/// borrower is set.
contract SupplyVault is ERC4626, Ownable2Step {
    /// @notice A two-slope borrow rate over utilisation, in basis points a year: `base` at zero, `base + slope1` at
    /// `optimal`, and `base + slope1 + slope2` at full utilisation, linear in between.
    struct RateModel {
        uint16 optimal;
        uint32 base;
        uint32 slope1;
        uint32 slope2;
    }

    /// @notice The lowest optimal utilisation, in basis points.
    uint16 public constant MIN_OPTIMAL = 1_00;
    /// @notice The highest optimal utilisation, in basis points.
    uint16 public constant MAX_OPTIMAL = 99_00;
    /// @notice The highest borrow rate the model may reach, in basis points a year: 1,000%.
    uint32 public constant MAX_RATE = 1000_00;

    uint256 private constant WAD = 1e18;
    uint256 private constant RAY = 1e27;
    uint256 private constant WAD_PER_BPS = 1e14;
    uint256 private constant YEAR = 365 days;

    IUSDG private immutable _usdg;

    /// @notice The only address that may borrow and write debt off.
    address public borrower;
    /// @notice The USDG the vault holds for lenders and the borrower, by its own count.
    uint128 public idle;
    /// @notice The borrower's debt divided by the borrow index at which it was taken, in USDG at an index of one.
    uint128 public scaledDebt;
    /// @notice What one unit of debt taken at the start is owed now, with 27 decimals, as of the last accrual.
    uint128 public borrowIndex;
    /// @notice When interest last accrued.
    uint64 public lastAccrual;
    RateModel private _rateModel;

    /// @notice The borrower was set.
    event BorrowerSet(address indexed borrower);
    /// @notice The rate model was set.
    event RateModelSet(uint16 optimal, uint32 base, uint32 slope1, uint32 slope2);
    /// @notice Interest accrued on the borrower's debt, which the index now reaches.
    event InterestAccrued(uint256 interest, uint256 index);
    /// @notice The borrower borrowed `assets` to `receiver`.
    event Borrow(address indexed receiver, uint256 assets);
    /// @notice `payer` repaid `assets` of the borrower's debt.
    event Repay(address indexed payer, uint256 assets);
    /// @notice The borrower wrote `assets` of its debt off, a loss to lenders.
    event WriteOff(uint256 assets);
    /// @notice The vault's count of the USDG it holds fell to its balance.
    event Sync(uint256 idle);

    /// @notice `account` is not the borrower.
    error NotBorrower(address account);
    /// @notice The borrower is already set, to `current`.
    error BorrowerAlreadySet(address current);
    /// @notice The borrower cannot be the zero address.
    error InvalidBorrower();
    /// @notice The vault holds `available` USDG, less than `requested`.
    error InsufficientLiquidity(uint256 requested, uint256 available);
    /// @notice The rate model is outside its bounds.
    error InvalidRateModel();
    /// @notice The owner cannot renounce: the vault needs one to set its borrower and its rate model.
    error OwnershipCannotBeRenounced();

    /// @param usdg The USDG token.
    /// @param initialOwner The owner, who sets the borrower once and the rate model within its bounds.
    /// @param initialRateModel The first rate model.
    constructor(IUSDG usdg, address initialOwner, RateModel memory initialRateModel)
        ERC20("Tapehouse USDG Supply", "thUSDG")
        ERC4626(IERC20(address(usdg)))
        Ownable(initialOwner)
    {
        _usdg = usdg;
        lastAccrual = SafeCast.toUint64(block.timestamp);
        borrowIndex = SafeCast.toUint128(RAY);
        _setRateModel(initialRateModel);
    }

    /// @notice The rate model.
    function rateModel() external view returns (RateModel memory) {
        return _rateModel;
    }

    /// @notice Always reverts: the vault keeps an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    /// @notice Sets the only address that may borrow. It can be set once, and never to zero.
    function setBorrower(address newBorrower) external onlyOwner {
        if (borrower != address(0)) revert BorrowerAlreadySet(borrower);
        if (newBorrower == address(0)) revert InvalidBorrower();
        borrower = newBorrower;
        emit BorrowerSet(newBorrower);
    }

    /// @notice Sets the rate model, within its bounds, at once. Interest accrues first at the old rate.
    function setRateModel(RateModel calldata model) external onlyOwner {
        _accrue();
        _setRateModel(model);
    }

    /// @notice Lends `assets` of the vault's USDG to `receiver`. Only the borrower may.
    function borrow(uint256 assets, address receiver) external {
        if (msg.sender != borrower) revert NotBorrower(msg.sender);
        uint256 index = _accrue();
        if (assets > idle) revert InsufficientLiquidity(assets, idle);
        scaledDebt += SafeCast.toUint128(_divUp(assets * RAY, index));
        emit Borrow(receiver, assets);
        _transferOut(receiver, assets);
    }

    /// @notice Repays up to `assets` of the debt from the borrower, and returns what it repaid. Only the borrower may:
    /// it keeps each of its accounts' share of the debt, which a repayment made around it would not reduce.
    function repay(uint256 assets) external returns (uint256 repaid) {
        if (msg.sender != borrower) revert NotBorrower(msg.sender);
        repaid = _reduceDebt(assets);
        emit Repay(msg.sender, repaid);
        _transferIn(msg.sender, repaid);
    }

    /// @notice Writes up to `assets` of the borrower's debt off, and returns what it wrote off. Every share loses
    /// its part: lenders bear the loss. Only the borrower may.
    function writeOff(uint256 assets) external returns (uint256 written) {
        if (msg.sender != borrower) revert NotBorrower(msg.sender);
        written = _reduceDebt(assets);
        emit WriteOff(written);
    }

    /// @notice Brings the vault's count of the USDG it holds down to its balance, as after the issuer wipes a
    /// frozen vault. It never counts USDG sent to the vault directly. Anyone may call it.
    function sync() external {
        uint256 balance = _usdg.balanceOf(address(this));
        if (balance < idle) {
            _accrue();
            idle = SafeCast.toUint128(balance);
            emit Sync(balance);
        }
    }

    /// @notice Deposits `assets` with an EIP-2612 permit for them. A permit already used, as by a front-runner,
    /// does not stop the deposit if the allowance stands.
    function depositWithPermit(uint256 assets, address receiver, uint256 deadline, uint8 v, bytes32 r, bytes32 s)
        external
        returns (uint256 shares)
    {
        try _usdg.permit(msg.sender, address(this), assets, deadline, v, r, s) {} catch {}
        return deposit(assets, receiver);
    }

    /// @notice What the borrower owes now, interest included, rounded up.
    function debt() public view returns (uint256) {
        return _divUp(scaledDebt * _currentIndex(), RAY);
    }

    /// @notice The share of the vault's assets lent, with 18 decimals.
    function utilization() public view returns (uint256) {
        return _utilization(debt(), idle);
    }

    /// @notice The borrow rate now, a year, with 18 decimals.
    function borrowRate() public view returns (uint256) {
        return _borrowRate(utilization());
    }

    /// @notice The rate lenders earn now, a year, with 18 decimals: the borrow rate on the share lent.
    function supplyRate() external view returns (uint256) {
        uint256 u = utilization();
        return _borrowRate(u) * u / WAD;
    }

    /// @inheritdoc ERC4626
    function totalAssets() public view override returns (uint256) {
        return idle + debt();
    }

    /// @inheritdoc ERC4626
    /// @dev Zero until the borrower is set, and while USDG is paused, the vault frozen or its USDG wiped.
    function maxDeposit(address receiver) public view override returns (uint256) {
        return borrower == address(0) || _unavailable() ? 0 : super.maxDeposit(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev Zero until the borrower is set, and while USDG is paused, the vault frozen or its USDG wiped.
    function maxMint(address receiver) public view override returns (uint256) {
        return borrower == address(0) || _unavailable() ? 0 : super.maxMint(receiver);
    }

    /// @inheritdoc ERC4626
    /// @dev Capped by the USDG the vault holds.
    function maxRedeem(address owner) public view override returns (uint256) {
        if (_unavailable()) return 0;
        uint256 shares = super.maxRedeem(owner);
        uint256 liquid = _convertToShares(idle, Math.Rounding.Floor);
        return shares < liquid ? shares : liquid;
    }

    function _deposit(address caller, address receiver, uint256 assets, uint256 shares) internal override {
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

    function _decimalsOffset() internal pure override returns (uint8) {
        return 6;
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
        if (index > previous) {
            borrowIndex = SafeCast.toUint128(index);
            // forge-lint: disable-next-line(reentrancy-events)
            emit InterestAccrued(scaledDebt * (index - previous) / RAY, index);
        }
    }

    function _currentIndex() private view returns (uint256) {
        uint256 index = borrowIndex;
        uint256 lent = _divUp(scaledDebt * index, RAY);
        uint256 rate = _borrowRate(_utilization(lent, idle));
        return index + index * rate * (block.timestamp - lastAccrual) / (WAD * YEAR);
    }

    function _setRateModel(RateModel memory model) private {
        if (
            model.optimal < MIN_OPTIMAL || model.optimal > MAX_OPTIMAL || model.slope1 > model.slope2
                || uint256(model.base) + model.slope1 + model.slope2 > MAX_RATE
        ) revert InvalidRateModel();
        _rateModel = model;
        emit RateModelSet(model.optimal, model.base, model.slope1, model.slope2);
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
        RateModel memory m = _rateModel;
        uint256 optimal = m.optimal * WAD_PER_BPS;
        uint256 rate = m.base * WAD_PER_BPS;
        if (u <= optimal) return rate + m.slope1 * WAD_PER_BPS * u / optimal;
        return rate + m.slope1 * WAD_PER_BPS + m.slope2 * WAD_PER_BPS * (u - optimal) / (WAD - optimal);
    }
}
