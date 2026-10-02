// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IBand} from "./interfaces/IBand.sol";
import {IStockToken, IStockTokenRegistry} from "./interfaces/IStockToken.sol";

/// @title Tapehouse basket
/// @notice A share of a fixed set of Stock Tokens, minted and redeemed in kind. A share is a claim on its part of
/// every token the basket holds, in raw units, so dividends and splits, which move only a token's ERC-8056 multiplier,
/// pass through, and a burn by the issuer falls on every holder in proportion. Minting pays each token's part rounded
/// up and redeeming takes it rounded down, so neither ever takes from the other holders; the first mint pays the
/// target. The owner proposes a new target, in raw units of each token per share, which takes effect after `NOTICE`.
/// Anyone may then move the basket toward it with `rebalance`, giving it tokens it holds too few of and taking those it
/// holds too many of, as long as what it gives, at each band's low edge, is worth at least what it takes, at each band's
/// high edge: at every price the band allows, the basket is worth no less after. A share moves only as its tokens
/// could: not while any of them is paused, nor to, from or by an address the issuer has blocklisted.
/// @dev Shares follow ERC-4626's rounding, with a list of assets in place of one.
// slither-disable-next-line missing-inheritance
contract Basket is ERC20, Ownable2Step {
    using SafeERC20 for IERC20;

    /// @notice The most tokens a basket holds: the margin engine's most assets.
    uint256 public constant MAX_COMPONENTS = 8;
    /// @notice How long a proposed target waits before it takes effect, so a holder who disagrees may redeem first.
    uint64 public constant NOTICE = 7 days;

    uint256 private constant WAD = 1e18;

    /// @notice The band that prices every token in a rebalance.
    IBand public immutable band;

    bytes32[] private _symbols;
    address[] private _tokens;
    address[] private _registries;
    uint256[] private _target;
    uint256[] private _pending;
    uint64 private _pendingAt;

    /// @notice `sender` paid `assets` of the tokens, in their order, for `shares` minted to `owner`.
    event Deposit(address indexed sender, address indexed owner, uint256[] assets, uint256 shares);
    /// @notice `sender` redeemed `shares` of `owner` for `assets` of the tokens, in their order, sent to `receiver`.
    event Withdraw(
        address indexed sender, address indexed receiver, address indexed owner, uint256[] assets, uint256 shares
    );
    /// @notice The owner proposed `units` of each token per share as the target, from `effectiveAt`, in seconds.
    event TargetProposed(uint256[] units, uint64 effectiveAt);
    /// @notice `caller` gave the basket `assetsIn` and took `assetsOut` to `receiver`, toward the target.
    event Rebalanced(address indexed caller, address indexed receiver, uint256[] assetsIn, uint256[] assetsOut);

    /// @notice `symbol` has no Stock Token in the band.
    error UnknownAsset(bytes32 symbol);
    /// @notice `symbol` is listed twice.
    error DuplicateComponent(bytes32 symbol);
    /// @notice A basket holds from one to `MAX_COMPONENTS` tokens.
    error InvalidComponents();
    /// @notice A list is not one entry per token.
    error LengthMismatch();
    /// @notice A target holds nothing of every token, or the first holds nothing of one of them.
    error InvalidTarget();
    /// @notice An amount of zero.
    error ZeroAmount();
    /// @notice Every token the basket holds has been burned: it takes no mint.
    error EmptyBasket();
    /// @notice Minting would take `amount` of `token`, more than the caller's `maxAmount`.
    error AboveMaximum(address token, uint256 amount, uint256 maxAmount);
    /// @notice A share moves only as its tokens could: one of them is paused.
    error ComponentPaused();
    /// @notice The issuer has blocklisted `account`, so it neither holds, sends nor moves shares.
    error Blocked(address account);
    /// @notice The L2 sequencer is down or was, within the hour.
    error SequencerNotSettled();
    /// @notice `symbol`'s band is halted.
    error AssetHalted(bytes32 symbol);
    /// @notice `symbol`'s Stock Token has a multiplier change pending.
    error CorporateActionPending(bytes32 symbol);
    /// @notice The rebalance would take `symbol` past its target, or away from it.
    error PastTarget(bytes32 symbol);
    /// @notice What the basket would take, `valueOut` at the bands' high edges, is worth more than what it would get,
    /// `valueIn` at their low edges; both in raw units times USD with 8 decimals.
    error ValueLost(uint256 valueIn, uint256 valueOut);
    /// @notice The owner cannot renounce: the basket keeps one to propose targets.
    error OwnershipCannotBeRenounced();

    /// @param name_ The share's name.
    /// @param symbol_ The share's symbol.
    /// @param band_ The band whose Stock Tokens the basket holds and whose edges price a rebalance.
    /// @param symbols The band's assets the basket holds, in the order of every list of amounts.
    /// @param units The first target: raw units of each token per share, none of them zero.
    /// @param initialOwner The owner, who proposes targets.
    constructor(
        string memory name_,
        string memory symbol_,
        IBand band_,
        bytes32[] memory symbols,
        uint256[] memory units,
        address initialOwner
    ) ERC20(name_, symbol_) Ownable(initialOwner) {
        uint256 n = symbols.length;
        if (n == 0 || n > MAX_COMPONENTS) revert InvalidComponents();
        if (units.length != n) revert LengthMismatch();
        band = band_;
        for (uint256 i; i < n; ++i) {
            // slither-disable-next-line unused-return,calls-loop
            (,,, address token) = band_.asset(symbols[i]); // forge-lint: disable-line(calls-loop, unused-return)
            // forge-lint: disable-start(require-revert-in-loop)
            if (token == address(0)) revert UnknownAsset(symbols[i]);
            for (uint256 j; j < i; ++j) {
                if (_tokens[j] == token) revert DuplicateComponent(symbols[i]);
            }
            if (units[i] == 0) revert InvalidTarget();
            // forge-lint: disable-end(require-revert-in-loop)
            _symbols.push(symbols[i]);
            _tokens.push(token);
            // forge-lint: disable-next-line(calls-loop)
            address registry = IStockToken(token).ACCESS_CONTROLLED_REGISTRY(); // slither-disable-line calls-loop
            bool known = false;
            uint256 registries = _registries.length;
            for (uint256 j; j < registries; ++j) {
                if (_registries[j] == registry) known = true;
            }
            if (!known) _registries.push(registry);
        }
        _target = units;
    }

    /// @notice The band's assets the basket holds and their Stock Tokens, in the order of every list of amounts.
    function components() external view returns (bytes32[] memory symbols, address[] memory tokens) {
        return (_symbols, _tokens);
    }

    /// @notice What the basket holds of each token, in raw units.
    function totalAssets() public view returns (uint256[] memory assets) {
        uint256 n = _tokens.length;
        assets = new uint256[](n);
        for (uint256 i; i < n; ++i) {
            // slither-disable-next-line calls-loop
            assets[i] = IERC20(_tokens[i]).balanceOf(address(this)); // forge-lint: disable-line(calls-loop)
        }
    }

    /// @notice The target in effect: raw units of each token per share.
    function target() public view returns (uint256[] memory) {
        // forge-lint: disable-next-line(block-timestamp)
        return _pendingAt != 0 && block.timestamp >= _pendingAt ? _pending : _target;
    }

    /// @notice The target proposed and not yet in effect, and when it takes effect, in seconds; empty and zero if none.
    function pendingTarget() external view returns (uint256[] memory units, uint64 effectiveAt) {
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= _pendingAt) return (new uint256[](0), 0);
        return (_pending, _pendingAt);
    }

    /// @notice What minting `shares` takes of each token: their part of what the basket holds, rounded up, or of the
    /// target while there are no shares.
    function previewMint(uint256 shares) public view returns (uint256[] memory assets) {
        uint256 supply = totalSupply();
        if (supply == 0) {
            uint256[] memory units = target();
            assets = new uint256[](units.length);
            for (uint256 i; i < units.length; ++i) {
                assets[i] = Math.mulDiv(shares, units[i], WAD, Math.Rounding.Ceil);
            }
            return assets;
        }
        assets = totalAssets();
        for (uint256 i; i < assets.length; ++i) {
            assets[i] = Math.mulDiv(shares, assets[i], supply, Math.Rounding.Ceil);
        }
    }

    /// @notice What redeeming `shares` gives of each token: their part of what the basket holds, rounded down.
    function previewRedeem(uint256 shares) public view returns (uint256[] memory assets) {
        uint256 supply = totalSupply();
        assets = totalAssets();
        for (uint256 i; i < assets.length; ++i) {
            assets[i] = supply == 0 ? 0 : Math.mulDiv(shares, assets[i], supply);
        }
    }

    /// @notice Whether a share cannot move now, because one of its tokens is paused.
    function paused() public view returns (bool) {
        uint256 n = _tokens.length;
        for (uint256 i; i < n; ++i) {
            // slither-disable-next-line calls-loop
            if (IStockToken(_tokens[i]).paused()) return true; // forge-lint: disable-line(calls-loop)
        }
        return false;
    }

    /// @notice Whether the issuer has blocklisted `account` in the registry of any of the basket's tokens.
    function isBlocked(address account) public view returns (bool) {
        uint256 n = _registries.length;
        for (uint256 i; i < n; ++i) {
            // slither-disable-next-line calls-loop
            if (IStockTokenRegistry(_registries[i]).isBlocked(account)) return true; // forge-lint: disable-line(calls-loop)
        }
        return false;
    }

    /// @notice Mints `shares` to `receiver` for what `previewMint` takes from the caller, each token at most its
    /// `maxAssets`, and returns what it took.
    function mint(uint256 shares, address receiver, uint256[] calldata maxAssets)
        external
        returns (uint256[] memory assets)
    {
        if (shares == 0) revert ZeroAmount();
        if (maxAssets.length != _tokens.length) revert LengthMismatch();
        assets = previewMint(shares);
        bool any = false;
        for (uint256 i; i < assets.length; ++i) {
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (assets[i] > maxAssets[i]) revert AboveMaximum(_tokens[i], assets[i], maxAssets[i]);
            if (assets[i] != 0) any = true;
        }
        if (!any) revert EmptyBasket();
        _mint(receiver, shares);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Deposit(msg.sender, receiver, assets, shares);
        for (uint256 i; i < assets.length; ++i) {
            // forge-lint: disable-next-line(calls-loop)
            if (assets[i] != 0) IERC20(_tokens[i]).safeTransferFrom(msg.sender, address(this), assets[i]);
        }
    }

    /// @notice Burns `shares` of `owner`, with its allowance if the caller is not `owner`, sends what `previewRedeem`
    /// gives to `receiver`, and returns it.
    function redeem(uint256 shares, address receiver, address owner) external returns (uint256[] memory assets) {
        if (shares == 0) revert ZeroAmount();
        if (msg.sender != owner) _spendAllowance(owner, msg.sender, shares);
        assets = previewRedeem(shares);
        _burn(owner, shares);
        // forge-lint: disable-next-line(reentrancy-events)
        emit Withdraw(msg.sender, receiver, owner, assets, shares);
        for (uint256 i; i < assets.length; ++i) {
            // forge-lint: disable-next-line(calls-loop)
            if (assets[i] != 0) IERC20(_tokens[i]).safeTransfer(receiver, assets[i]);
        }
    }

    /// @notice Proposes `units` of each token per share as the target, in effect `NOTICE` from now. A target already
    /// in effect stays until then; one proposed and not yet in effect is replaced. Only the owner may.
    function proposeTarget(uint256[] calldata units) external onlyOwner {
        if (units.length != _tokens.length) revert LengthMismatch();
        bool any = false;
        for (uint256 i; i < units.length; ++i) {
            if (units[i] != 0) any = true;
        }
        if (!any) revert InvalidTarget();
        // forge-lint: disable-next-line(block-timestamp)
        if (_pendingAt != 0 && block.timestamp >= _pendingAt) _target = _pending;
        _pending = units;
        _pendingAt = SafeCast.toUint64(block.timestamp + NOTICE);
        emit TargetProposed(units, _pendingAt);
    }

    /// @notice Takes `assetsIn` from the caller and sends `assetsOut` to `receiver`, each in the tokens' order, if every
    /// token moves toward the target in effect without passing it, none both in and out, and what comes in at its
    /// band's low edge is worth at least what goes out at its band's high edge. Each token that moves needs a band with
    /// a price and no multiplier change pending, and the L2 sequencer settled. Anyone may.
    function rebalance(uint256[] calldata assetsIn, uint256[] calldata assetsOut, address receiver) external {
        uint256 n = _tokens.length;
        if (assetsIn.length != n || assetsOut.length != n) revert LengthMismatch();
        if (!band.sequencerSettled()) revert SequencerNotSettled();
        uint256 supply = totalSupply();
        uint256[] memory units = target();
        uint256 valueIn = 0;
        uint256 valueOut = 0;
        bool any = false;
        for (uint256 i; i < n; ++i) {
            (uint256 amountIn, uint256 amountOut) = (assetsIn[i], assetsOut[i]);
            if (amountIn == 0 && amountOut == 0) continue;
            bytes32 symbol = _symbols[i];
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (amountIn != 0 && amountOut != 0) revert PastTarget(symbol);
            // slither-disable-next-line unused-return,calls-loop
            (uint8 state,,,, uint64 low, uint128 high) = band.quote(symbol); // forge-lint: disable-line(calls-loop, unused-return)
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (state == 0) revert AssetHalted(symbol);
            // slither-disable-next-line unused-return,calls-loop
            (uint8 status,,,) = band.corporateAction(symbol); // forge-lint: disable-line(calls-loop, unused-return)
            // forge-lint: disable-next-line(require-revert-in-loop)
            if (status != 0) revert CorporateActionPending(symbol);
            // forge-lint: disable-next-line(calls-loop)
            uint256 held = IERC20(_tokens[i]).balanceOf(address(this));
            uint256 goal = units[i] * supply;
            if (amountIn != 0) {
                // forge-lint: disable-next-line(require-revert-in-loop)
                if ((held + amountIn) * WAD > goal) revert PastTarget(symbol);
                valueIn += amountIn * low;
            } else {
                // forge-lint: disable-next-line(require-revert-in-loop)
                if (amountOut > held || (held - amountOut) * WAD < goal) revert PastTarget(symbol);
                valueOut += amountOut * high;
            }
            any = true;
        }
        if (!any) revert ZeroAmount();
        if (valueIn < valueOut) revert ValueLost(valueIn, valueOut);
        emit Rebalanced(msg.sender, receiver, assetsIn, assetsOut);
        for (uint256 i; i < n; ++i) {
            // forge-lint: disable-start(calls-loop)
            if (assetsIn[i] != 0) IERC20(_tokens[i]).safeTransferFrom(msg.sender, address(this), assetsIn[i]);
            if (assetsOut[i] != 0) IERC20(_tokens[i]).safeTransfer(receiver, assetsOut[i]);
            // forge-lint: disable-end(calls-loop)
        }
    }

    /// @notice Always reverts: the basket keeps an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    function _update(address from, address to, uint256 value) internal override {
        if (paused()) revert ComponentPaused();
        if (from != address(0) && isBlocked(from)) revert Blocked(from);
        if (to != address(0) && isBlocked(to)) revert Blocked(to);
        if (msg.sender != from && msg.sender != to && isBlocked(msg.sender)) revert Blocked(msg.sender);
        super._update(from, to, value);
    }
}
