// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity 0.8.37;

import {Ownable, Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";
import {Paymaster} from "@openzeppelin/contracts/account/paymaster/Paymaster.sol";
import {PaymasterSigner} from "@openzeppelin/contracts/account/paymaster/extensions/PaymasterSigner.sol";
import {IEntryPoint, PackedUserOperation} from "@openzeppelin/contracts/interfaces/IERC4337.sol";
import {EIP712} from "@openzeppelin/contracts/utils/cryptography/EIP712.sol";
import {SignerECDSA} from "@openzeppelin/contracts/utils/cryptography/signers/SignerECDSA.sol";

/// @title Tapehouse sponsor paymaster
/// @notice Pays the gas of a smart account's first operations with Tapehouse, so a new user needs no ether to start.
/// Tapehouse's sponsor service decides what it pays for and signs each operation it sponsors, gas fields and validity
/// window included; the paymaster pays only for operations so signed. Creating an account is free; each account then
/// has `FREE_OPERATIONS` sponsored operations. No operation may cost it more than `maxCostPerOperation`.
/// @dev The signature is over OpenZeppelin's `UserOperationRequest` (EIP-712), in `paymasterData` after the window:
/// `validAfter` and `validUntil`, 6 bytes each, then the signature. The paymaster reads its own storage during
/// validation, so it must be staked.
contract SponsorPaymaster is PaymasterSigner, SignerECDSA, Ownable2Step {
    /// @notice How many operations each account has sponsored, besides its creation.
    uint256 public constant FREE_OPERATIONS = 3;

    IEntryPoint private immutable _entryPoint;

    /// @notice The most an operation may cost the paymaster, in wei.
    uint128 public maxCostPerOperation;

    /// @notice Whether sponsorship is paused.
    bool public paused;

    /// @notice How many sponsored operations `account` has used.
    mapping(address account => uint256) public operations;

    /// @notice Emitted when the sponsor service's signing address changes.
    event SignerSet(address indexed signer);

    /// @notice Emitted when the cost limit per operation changes.
    event MaxCostSet(uint128 maxCost);

    /// @notice Emitted when sponsorship is paused or resumed.
    event PausedSet(bool paused);

    /// @notice `account` has used its free operations.
    error FreeOperationsUsed(address account);

    /// @notice The operation may cost more than `maxCostPerOperation`.
    error CostAboveLimit(uint256 maxCost);

    /// @notice Sponsorship is paused.
    error SponsorshipPaused();

    /// @notice The paymaster always keeps an owner.
    error OwnershipCannotBeRenounced();

    constructor(IEntryPoint entryPoint_, address owner_, address signer_, uint128 maxCost)
        EIP712("Tapehouse Sponsor Paymaster", "1")
        SignerECDSA(signer_)
        Ownable(owner_)
    {
        _entryPoint = entryPoint_;
        maxCostPerOperation = maxCost;
        emit SignerSet(signer_);
        emit MaxCostSet(maxCost);
    }

    /// @inheritdoc Paymaster
    function entryPoint() public view override returns (IEntryPoint) {
        return _entryPoint;
    }

    /// @notice How many sponsored operations `account` has left.
    function freeOperationsLeft(address account) external view returns (uint256) {
        uint256 used = operations[account];
        return used < FREE_OPERATIONS ? FREE_OPERATIONS - used : 0;
    }

    /// @notice Sets the sponsor service's signing address.
    function setSigner(address signer_) external onlyOwner {
        _setSigner(signer_);
        emit SignerSet(signer_);
    }

    /// @notice Sets the most an operation may cost the paymaster, in wei.
    function setMaxCost(uint128 maxCost) external onlyOwner {
        maxCostPerOperation = maxCost;
        emit MaxCostSet(maxCost);
    }

    /// @notice Pauses or resumes sponsorship.
    function setPaused(bool paused_) external onlyOwner {
        paused = paused_;
        emit PausedSet(paused_);
    }

    /// @notice Adds the ether sent to the paymaster's deposit in the EntryPoint.
    function deposit() external payable {
        _deposit(msg.value);
    }

    /// @notice Withdraws `value` wei of the paymaster's deposit to `to`.
    function withdrawTo(address payable to, uint256 value) external onlyOwner {
        _withdraw(to, value);
    }

    /// @notice Stakes the ether sent, unlockable `unstakeDelaySec` seconds after `unlockStake`.
    function addStake(uint32 unstakeDelaySec) external payable onlyOwner {
        _addStake(msg.value, unstakeDelaySec);
    }

    /// @notice Starts the delay after which the stake may be withdrawn.
    function unlockStake() external onlyOwner {
        _unlockStake();
    }

    /// @notice Withdraws the unlocked stake to `to`.
    function withdrawStake(address payable to) external onlyOwner {
        _withdrawStake(to);
    }

    /// @notice Always reverts: the paymaster keeps an owner.
    function renounceOwnership() public pure override {
        revert OwnershipCannotBeRenounced();
    }

    /// @dev Counts the operation unless it only creates the account, then checks the sponsor service's signature.
    function _validatePaymasterUserOp(PackedUserOperation calldata userOp, bytes32 userOpHash, uint256 requiredPreFund)
        internal
        override
        returns (bytes memory context, uint256 validationData)
    {
        if (paused) revert SponsorshipPaused();
        if (requiredPreFund > maxCostPerOperation) revert CostAboveLimit(requiredPreFund);
        if (userOp.callData.length != 0 || userOp.initCode.length == 0) {
            address sender = userOp.sender;
            uint256 used = operations[sender];
            if (used >= FREE_OPERATIONS) revert FreeOperationsUsed(sender);
            operations[sender] = used + 1;
        }
        return super._validatePaymasterUserOp(userOp, userOpHash, requiredPreFund);
    }
}
