// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gapbackstop

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// GapBackstopMetaData contains all meta data concerning the GapBackstop contract.
var GapBackstopMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"accounts_\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"crossLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"assetLimits\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"COOLDOWN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_CLOSURE_GAP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SETTLEMENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WITHDRAWAL_WINDOW\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"accounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyRemainder\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"taken\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimGains\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"closureMs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToAssets\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToShares\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cooldowns\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint192\",\"internalType\":\"uint192\"},{\"name\":\"startedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cover\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"paid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"written\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"covered\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"exposureLeft\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposureLimit\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exposureLimits\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"current\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"next\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"fromMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gains\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gainsEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gainsPerShare\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"held\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewDeposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setExposureLimit\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startCooldown\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sync\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CooldownStarted\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Covered\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"paid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"written\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ExposureLimitSet\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"limit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fromMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GainsClaimed\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumAdded\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemainderBought\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sync\",\"inputs\":[{\"name\":\"held\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidExposureLimits\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPosition\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotAuction\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "GapBackstop",
}

// GapBackstop is an auto generated Go binding around an Ethereum contract.
type GapBackstop struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *GapBackstop) GetABI() abi.ABI {
	return c.abi
}

// NewGapBackstop creates a new instance of GapBackstop.
func NewGapBackstop() *GapBackstop {
	parsed, err := GapBackstopMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &GapBackstop{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *GapBackstop) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address accounts_, address initialOwner, uint256 crossLimit, uint256[] assetLimits) returns()
func (gapBackstop *GapBackstop) PackConstructor(accounts_ common.Address, initialOwner common.Address, crossLimit *big.Int, assetLimits []*big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("", accounts_, initialOwner, crossLimit, assetLimits)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCOOLDOWN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2724a4d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function COOLDOWN() view returns(uint256)
func (gapBackstop *GapBackstop) PackCOOLDOWN() []byte {
	enc, err := gapBackstop.abi.Pack("COOLDOWN")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCOOLDOWN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2724a4d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function COOLDOWN() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackCOOLDOWN() ([]byte, error) {
	return gapBackstop.abi.Pack("COOLDOWN")
}

// UnpackCOOLDOWN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa2724a4d.
//
// Solidity: function COOLDOWN() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackCOOLDOWN(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("COOLDOWN", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINCLOSUREGAP is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4fdf7233.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_CLOSURE_GAP() view returns(uint256)
func (gapBackstop *GapBackstop) PackMINCLOSUREGAP() []byte {
	enc, err := gapBackstop.abi.Pack("MIN_CLOSURE_GAP")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINCLOSUREGAP is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4fdf7233.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_CLOSURE_GAP() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackMINCLOSUREGAP() ([]byte, error) {
	return gapBackstop.abi.Pack("MIN_CLOSURE_GAP")
}

// UnpackMINCLOSUREGAP is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4fdf7233.
//
// Solidity: function MIN_CLOSURE_GAP() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackMINCLOSUREGAP(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("MIN_CLOSURE_GAP", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSETTLEMENT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef0fe245.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function SETTLEMENT() view returns(uint256)
func (gapBackstop *GapBackstop) PackSETTLEMENT() []byte {
	enc, err := gapBackstop.abi.Pack("SETTLEMENT")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSETTLEMENT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef0fe245.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function SETTLEMENT() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackSETTLEMENT() ([]byte, error) {
	return gapBackstop.abi.Pack("SETTLEMENT")
}

// UnpackSETTLEMENT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef0fe245.
//
// Solidity: function SETTLEMENT() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackSETTLEMENT(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("SETTLEMENT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWITHDRAWALWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56765c51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function WITHDRAWAL_WINDOW() view returns(uint256)
func (gapBackstop *GapBackstop) PackWITHDRAWALWINDOW() []byte {
	enc, err := gapBackstop.abi.Pack("WITHDRAWAL_WINDOW")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWITHDRAWALWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56765c51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function WITHDRAWAL_WINDOW() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackWITHDRAWALWINDOW() ([]byte, error) {
	return gapBackstop.abi.Pack("WITHDRAWAL_WINDOW")
}

// UnpackWITHDRAWALWINDOW is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56765c51.
//
// Solidity: function WITHDRAWAL_WINDOW() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackWITHDRAWALWINDOW(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("WITHDRAWAL_WINDOW", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (gapBackstop *GapBackstop) PackAcceptOwnership() []byte {
	enc, err := gapBackstop.abi.Pack("acceptOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptOwnership() returns()
func (gapBackstop *GapBackstop) TryPackAcceptOwnership() ([]byte, error) {
	return gapBackstop.abi.Pack("acceptOwnership")
}

// PackAccounts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68cd03f6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function accounts() view returns(address)
func (gapBackstop *GapBackstop) PackAccounts() []byte {
	enc, err := gapBackstop.abi.Pack("accounts")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAccounts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68cd03f6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function accounts() view returns(address)
func (gapBackstop *GapBackstop) TryPackAccounts() ([]byte, error) {
	return gapBackstop.abi.Pack("accounts")
}

// UnpackAccounts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68cd03f6.
//
// Solidity: function accounts() view returns(address)
func (gapBackstop *GapBackstop) UnpackAccounts(data []byte) (common.Address, error) {
	out, err := gapBackstop.abi.Unpack("accounts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (gapBackstop *GapBackstop) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) UnpackApprove(data []byte) (bool, error) {
	out, err := gapBackstop.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function asset() view returns(address)
func (gapBackstop *GapBackstop) PackAsset() []byte {
	enc, err := gapBackstop.abi.Pack("asset")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function asset() view returns(address)
func (gapBackstop *GapBackstop) TryPackAsset() ([]byte, error) {
	return gapBackstop.abi.Pack("asset")
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (gapBackstop *GapBackstop) UnpackAsset(data []byte) (common.Address, error) {
	out, err := gapBackstop.abi.Unpack("asset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (gapBackstop *GapBackstop) PackBalanceOf(account common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBuyRemainder is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53b9ad4b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buyRemainder(address account, bytes32 position, address token, uint256 amount, uint256 price) returns(uint256 taken)
func (gapBackstop *GapBackstop) PackBuyRemainder(account common.Address, position [32]byte, token common.Address, amount *big.Int, price *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("buyRemainder", account, position, token, amount, price)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuyRemainder is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x53b9ad4b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buyRemainder(address account, bytes32 position, address token, uint256 amount, uint256 price) returns(uint256 taken)
func (gapBackstop *GapBackstop) TryPackBuyRemainder(account common.Address, position [32]byte, token common.Address, amount *big.Int, price *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("buyRemainder", account, position, token, amount, price)
}

// UnpackBuyRemainder is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x53b9ad4b.
//
// Solidity: function buyRemainder(address account, bytes32 position, address token, uint256 amount, uint256 price) returns(uint256 taken)
func (gapBackstop *GapBackstop) UnpackBuyRemainder(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("buyRemainder", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e71d92d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim() returns(uint256 amount)
func (gapBackstop *GapBackstop) PackClaim() []byte {
	enc, err := gapBackstop.abi.Pack("claim")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e71d92d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim() returns(uint256 amount)
func (gapBackstop *GapBackstop) TryPackClaim() ([]byte, error) {
	return gapBackstop.abi.Pack("claim")
}

// UnpackClaim is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4e71d92d.
//
// Solidity: function claim() returns(uint256 amount)
func (gapBackstop *GapBackstop) UnpackClaim(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("claim", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClaimGains is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205778bd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimGains(address token) returns(uint256 amount)
func (gapBackstop *GapBackstop) PackClaimGains(token common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("claimGains", token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimGains is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205778bd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimGains(address token) returns(uint256 amount)
func (gapBackstop *GapBackstop) TryPackClaimGains(token common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("claimGains", token)
}

// UnpackClaimGains is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x205778bd.
//
// Solidity: function claimGains(address token) returns(uint256 amount)
func (gapBackstop *GapBackstop) UnpackClaimGains(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("claimGains", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClosureMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f249e34.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function closureMs() view returns(uint64)
func (gapBackstop *GapBackstop) PackClosureMs() []byte {
	enc, err := gapBackstop.abi.Pack("closureMs")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClosureMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0f249e34.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function closureMs() view returns(uint64)
func (gapBackstop *GapBackstop) TryPackClosureMs() ([]byte, error) {
	return gapBackstop.abi.Pack("closureMs")
}

// UnpackClosureMs is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0f249e34.
//
// Solidity: function closureMs() view returns(uint64)
func (gapBackstop *GapBackstop) UnpackClosureMs(data []byte) (uint64, error) {
	out, err := gapBackstop.abi.Unpack("closureMs", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) PackConvertToAssets(shares *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("convertToAssets", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackConvertToAssets(shares *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("convertToAssets", shares)
}

// UnpackConvertToAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackConvertToAssets(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("convertToAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackConvertToShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6e6f592.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) PackConvertToShares(assets *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("convertToShares", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackConvertToShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6e6f592.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackConvertToShares(assets *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("convertToShares", assets)
}

// UnpackConvertToShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackConvertToShares(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("convertToShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCooldowns is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01320fe2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cooldowns(address owner) view returns(uint192 shares, uint64 startedAt)
func (gapBackstop *GapBackstop) PackCooldowns(owner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("cooldowns", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCooldowns is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01320fe2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cooldowns(address owner) view returns(uint192 shares, uint64 startedAt)
func (gapBackstop *GapBackstop) TryPackCooldowns(owner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("cooldowns", owner)
}

// CooldownsOutput serves as a container for the return parameters of contract
// method Cooldowns.
type CooldownsOutput struct {
	Shares    *big.Int
	StartedAt uint64
}

// UnpackCooldowns is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01320fe2.
//
// Solidity: function cooldowns(address owner) view returns(uint192 shares, uint64 startedAt)
func (gapBackstop *GapBackstop) UnpackCooldowns(data []byte) (CooldownsOutput, error) {
	out, err := gapBackstop.abi.Unpack("cooldowns", data)
	outstruct := new(CooldownsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Shares = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.StartedAt = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackCover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f273805.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cover(address account, bytes32 position) returns(uint256 paid, uint256 written)
func (gapBackstop *GapBackstop) PackCover(account common.Address, position [32]byte) []byte {
	enc, err := gapBackstop.abi.Pack("cover", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5f273805.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cover(address account, bytes32 position) returns(uint256 paid, uint256 written)
func (gapBackstop *GapBackstop) TryPackCover(account common.Address, position [32]byte) ([]byte, error) {
	return gapBackstop.abi.Pack("cover", account, position)
}

// CoverOutput serves as a container for the return parameters of contract
// method Cover.
type CoverOutput struct {
	Paid    *big.Int
	Written *big.Int
}

// UnpackCover is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5f273805.
//
// Solidity: function cover(address account, bytes32 position) returns(uint256 paid, uint256 written)
func (gapBackstop *GapBackstop) UnpackCover(data []byte) (CoverOutput, error) {
	out, err := gapBackstop.abi.Unpack("cover", data)
	outstruct := new(CoverOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Paid = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Written = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCovered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd87408fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function covered(bytes32 position, uint64 closesMs) view returns(uint256)
func (gapBackstop *GapBackstop) PackCovered(position [32]byte, closesMs uint64) []byte {
	enc, err := gapBackstop.abi.Pack("covered", position, closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCovered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd87408fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function covered(bytes32 position, uint64 closesMs) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackCovered(position [32]byte, closesMs uint64) ([]byte, error) {
	return gapBackstop.abi.Pack("covered", position, closesMs)
}

// UnpackCovered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd87408fe.
//
// Solidity: function covered(bytes32 position, uint64 closesMs) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackCovered(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("covered", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (gapBackstop *GapBackstop) PackDecimals() []byte {
	enc, err := gapBackstop.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() view returns(uint8)
func (gapBackstop *GapBackstop) TryPackDecimals() ([]byte, error) {
	return gapBackstop.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (gapBackstop *GapBackstop) UnpackDecimals(data []byte) (uint8, error) {
	out, err := gapBackstop.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e553f65.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) PackDeposit(assets *big.Int, receiver common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("deposit", assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e553f65.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) TryPackDeposit(assets *big.Int, receiver common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("deposit", assets, receiver)
}

// UnpackDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) UnpackDeposit(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("deposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackExposureLeft is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6970fee7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function exposureLeft(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) PackExposureLeft(position [32]byte) []byte {
	enc, err := gapBackstop.abi.Pack("exposureLeft", position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExposureLeft is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6970fee7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function exposureLeft(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackExposureLeft(position [32]byte) ([]byte, error) {
	return gapBackstop.abi.Pack("exposureLeft", position)
}

// UnpackExposureLeft is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6970fee7.
//
// Solidity: function exposureLeft(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackExposureLeft(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("exposureLeft", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackExposureLimit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0581c502.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function exposureLimit(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) PackExposureLimit(position [32]byte) []byte {
	enc, err := gapBackstop.abi.Pack("exposureLimit", position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExposureLimit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0581c502.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function exposureLimit(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackExposureLimit(position [32]byte) ([]byte, error) {
	return gapBackstop.abi.Pack("exposureLimit", position)
}

// UnpackExposureLimit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0581c502.
//
// Solidity: function exposureLimit(bytes32 position) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackExposureLimit(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("exposureLimit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackExposureLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc8c5462.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function exposureLimits(bytes32 position) view returns(uint128 current, uint128 next, uint64 fromMs)
func (gapBackstop *GapBackstop) PackExposureLimits(position [32]byte) []byte {
	enc, err := gapBackstop.abi.Pack("exposureLimits", position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExposureLimits is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc8c5462.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function exposureLimits(bytes32 position) view returns(uint128 current, uint128 next, uint64 fromMs)
func (gapBackstop *GapBackstop) TryPackExposureLimits(position [32]byte) ([]byte, error) {
	return gapBackstop.abi.Pack("exposureLimits", position)
}

// ExposureLimitsOutput serves as a container for the return parameters of contract
// method ExposureLimits.
type ExposureLimitsOutput struct {
	Current *big.Int
	Next    *big.Int
	FromMs  uint64
}

// UnpackExposureLimits is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdc8c5462.
//
// Solidity: function exposureLimits(bytes32 position) view returns(uint128 current, uint128 next, uint64 fromMs)
func (gapBackstop *GapBackstop) UnpackExposureLimits(data []byte) (ExposureLimitsOutput, error) {
	out, err := gapBackstop.abi.Unpack("exposureLimits", data)
	outstruct := new(ExposureLimitsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Current = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Next = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.FromMs = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackGains is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf33f306c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function gains(address owner, address token) view returns(uint256)
func (gapBackstop *GapBackstop) PackGains(owner common.Address, token common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("gains", owner, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGains is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf33f306c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function gains(address owner, address token) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackGains(owner common.Address, token common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("gains", owner, token)
}

// UnpackGains is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf33f306c.
//
// Solidity: function gains(address owner, address token) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackGains(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("gains", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGainsEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x181929b6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function gainsEpoch() view returns(uint256)
func (gapBackstop *GapBackstop) PackGainsEpoch() []byte {
	enc, err := gapBackstop.abi.Pack("gainsEpoch")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGainsEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x181929b6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function gainsEpoch() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackGainsEpoch() ([]byte, error) {
	return gapBackstop.abi.Pack("gainsEpoch")
}

// UnpackGainsEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x181929b6.
//
// Solidity: function gainsEpoch() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackGainsEpoch(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("gainsEpoch", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGainsPerShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f45b76d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function gainsPerShare(address token) view returns(uint256)
func (gapBackstop *GapBackstop) PackGainsPerShare(token common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("gainsPerShare", token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGainsPerShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f45b76d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function gainsPerShare(address token) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackGainsPerShare(token common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("gainsPerShare", token)
}

// UnpackGainsPerShare is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f45b76d.
//
// Solidity: function gainsPerShare(address token) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackGainsPerShare(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("gainsPerShare", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHeld is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa285aed7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function held() view returns(uint256)
func (gapBackstop *GapBackstop) PackHeld() []byte {
	enc, err := gapBackstop.abi.Pack("held")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHeld is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa285aed7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function held() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackHeld() ([]byte, error) {
	return gapBackstop.abi.Pack("held")
}

// UnpackHeld is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa285aed7.
//
// Solidity: function held() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackHeld(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("held", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) PackMaxDeposit(receiver common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("maxDeposit", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackMaxDeposit(receiver common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("maxDeposit", receiver)
}

// UnpackMaxDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackMaxDeposit(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("maxDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc63d75b6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) PackMaxMint(receiver common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("maxMint", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc63d75b6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackMaxMint(receiver common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("maxMint", receiver)
}

// UnpackMaxMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackMaxMint(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("maxMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905777e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) PackMaxRedeem(owner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("maxRedeem", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905777e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackMaxRedeem(owner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("maxRedeem", owner)
}

// UnpackMaxRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackMaxRedeem(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("maxRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce96cb77.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) PackMaxWithdraw(owner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("maxWithdraw", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce96cb77.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackMaxWithdraw(owner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("maxWithdraw", owner)
}

// UnpackMaxWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackMaxWithdraw(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("maxWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94bf804d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) PackMint(shares *big.Int, receiver common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("mint", shares, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94bf804d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) TryPackMint(shares *big.Int, receiver common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("mint", shares, receiver)
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (gapBackstop *GapBackstop) UnpackMint(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("mint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (gapBackstop *GapBackstop) PackName() []byte {
	enc, err := gapBackstop.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (gapBackstop *GapBackstop) TryPackName() ([]byte, error) {
	return gapBackstop.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (gapBackstop *GapBackstop) UnpackName(data []byte) (string, error) {
	out, err := gapBackstop.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (gapBackstop *GapBackstop) PackOwner() []byte {
	enc, err := gapBackstop.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (gapBackstop *GapBackstop) TryPackOwner() ([]byte, error) {
	return gapBackstop.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (gapBackstop *GapBackstop) UnpackOwner(data []byte) (common.Address, error) {
	out, err := gapBackstop.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (gapBackstop *GapBackstop) PackPendingOwner() []byte {
	enc, err := gapBackstop.abi.Pack("pendingOwner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingOwner() view returns(address)
func (gapBackstop *GapBackstop) TryPackPendingOwner() ([]byte, error) {
	return gapBackstop.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (gapBackstop *GapBackstop) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := gapBackstop.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) PackPreviewDeposit(assets *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("previewDeposit", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackPreviewDeposit(assets *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("previewDeposit", assets)
}

// UnpackPreviewDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackPreviewDeposit(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("previewDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) PackPreviewMint(shares *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("previewMint", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackPreviewMint(shares *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("previewMint", shares)
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackPreviewMint(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("previewMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("previewRedeem", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackPreviewRedeem(shares *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("previewRedeem", shares)
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackPreviewRedeem(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("previewRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a28a477.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) PackPreviewWithdraw(assets *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("previewWithdraw", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a28a477.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) TryPackPreviewWithdraw(assets *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("previewWithdraw", assets)
}

// UnpackPreviewWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (gapBackstop *GapBackstop) UnpackPreviewWithdraw(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("previewWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("redeem", shares, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) TryPackRedeem(shares *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("redeem", shares, receiver, owner)
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) UnpackRedeem(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("redeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() pure returns()
func (gapBackstop *GapBackstop) PackRenounceOwnership() []byte {
	enc, err := gapBackstop.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() pure returns()
func (gapBackstop *GapBackstop) TryPackRenounceOwnership() ([]byte, error) {
	return gapBackstop.abi.Pack("renounceOwnership")
}

// PackSetExposureLimit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1077a2e5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setExposureLimit(bytes32 position, uint256 limit) returns()
func (gapBackstop *GapBackstop) PackSetExposureLimit(position [32]byte, limit *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("setExposureLimit", position, limit)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetExposureLimit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1077a2e5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setExposureLimit(bytes32 position, uint256 limit) returns()
func (gapBackstop *GapBackstop) TryPackSetExposureLimit(position [32]byte, limit *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("setExposureLimit", position, limit)
}

// PackStartCooldown is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x680ff458.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function startCooldown() returns()
func (gapBackstop *GapBackstop) PackStartCooldown() []byte {
	enc, err := gapBackstop.abi.Pack("startCooldown")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStartCooldown is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x680ff458.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function startCooldown() returns()
func (gapBackstop *GapBackstop) TryPackStartCooldown() ([]byte, error) {
	return gapBackstop.abi.Pack("startCooldown")
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (gapBackstop *GapBackstop) PackSymbol() []byte {
	enc, err := gapBackstop.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (gapBackstop *GapBackstop) TryPackSymbol() ([]byte, error) {
	return gapBackstop.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (gapBackstop *GapBackstop) UnpackSymbol(data []byte) (string, error) {
	out, err := gapBackstop.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfff6cae9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sync() returns()
func (gapBackstop *GapBackstop) PackSync() []byte {
	enc, err := gapBackstop.abi.Pack("sync")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfff6cae9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sync() returns()
func (gapBackstop *GapBackstop) TryPackSync() ([]byte, error) {
	return gapBackstop.abi.Pack("sync")
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAssets() view returns(uint256)
func (gapBackstop *GapBackstop) PackTotalAssets() []byte {
	enc, err := gapBackstop.abi.Pack("totalAssets")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalAssets() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackTotalAssets() ([]byte, error) {
	return gapBackstop.abi.Pack("totalAssets")
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("totalAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (gapBackstop *GapBackstop) PackTotalSupply() []byte {
	enc, err := gapBackstop.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (gapBackstop *GapBackstop) TryPackTotalSupply() ([]byte, error) {
	return gapBackstop.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (gapBackstop *GapBackstop) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) UnpackTransfer(data []byte) (bool, error) {
	out, err := gapBackstop.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := gapBackstop.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return gapBackstop.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (gapBackstop *GapBackstop) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := gapBackstop.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (gapBackstop *GapBackstop) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (gapBackstop *GapBackstop) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("transferOwnership", newOwner)
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) PackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := gapBackstop.abi.Pack("withdraw", assets, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) TryPackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return gapBackstop.abi.Pack("withdraw", assets, receiver, owner)
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (gapBackstop *GapBackstop) UnpackWithdraw(data []byte) (*big.Int, error) {
	out, err := gapBackstop.abi.Unpack("withdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// GapBackstopApproval represents a Approval event raised by the GapBackstop contract.
type GapBackstopApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const GapBackstopApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (GapBackstopApproval) ContractEventName() string {
	return GapBackstopApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (gapBackstop *GapBackstop) UnpackApprovalEvent(log *types.Log) (*GapBackstopApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopApproval)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopCooldownStarted represents a CooldownStarted event raised by the GapBackstop contract.
type GapBackstopCooldownStarted struct {
	Owner  common.Address
	Shares *big.Int
	From   uint64
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapBackstopCooldownStartedEventName = "CooldownStarted"

// ContractEventName returns the user-defined event name.
func (GapBackstopCooldownStarted) ContractEventName() string {
	return GapBackstopCooldownStartedEventName
}

// UnpackCooldownStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CooldownStarted(address indexed owner, uint256 shares, uint64 from)
func (gapBackstop *GapBackstop) UnpackCooldownStartedEvent(log *types.Log) (*GapBackstopCooldownStarted, error) {
	event := "CooldownStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopCooldownStarted)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopCovered represents a Covered event raised by the GapBackstop contract.
type GapBackstopCovered struct {
	Account  common.Address
	Position [32]byte
	ClosesMs uint64
	Paid     *big.Int
	Written  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapBackstopCoveredEventName = "Covered"

// ContractEventName returns the user-defined event name.
func (GapBackstopCovered) ContractEventName() string {
	return GapBackstopCoveredEventName
}

// UnpackCoveredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Covered(address indexed account, bytes32 indexed position, uint64 indexed closesMs, uint256 paid, uint256 written)
func (gapBackstop *GapBackstop) UnpackCoveredEvent(log *types.Log) (*GapBackstopCovered, error) {
	event := "Covered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopCovered)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopDeposit represents a Deposit event raised by the GapBackstop contract.
type GapBackstopDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapBackstopDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (GapBackstopDeposit) ContractEventName() string {
	return GapBackstopDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (gapBackstop *GapBackstop) UnpackDepositEvent(log *types.Log) (*GapBackstopDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopDeposit)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopExposureLimitSet represents a ExposureLimitSet event raised by the GapBackstop contract.
type GapBackstopExposureLimitSet struct {
	Position [32]byte
	Limit    *big.Int
	FromMs   uint64
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapBackstopExposureLimitSetEventName = "ExposureLimitSet"

// ContractEventName returns the user-defined event name.
func (GapBackstopExposureLimitSet) ContractEventName() string {
	return GapBackstopExposureLimitSetEventName
}

// UnpackExposureLimitSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ExposureLimitSet(bytes32 indexed position, uint256 limit, uint64 fromMs)
func (gapBackstop *GapBackstop) UnpackExposureLimitSetEvent(log *types.Log) (*GapBackstopExposureLimitSet, error) {
	event := "ExposureLimitSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopExposureLimitSet)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopGainsClaimed represents a GainsClaimed event raised by the GapBackstop contract.
type GapBackstopGainsClaimed struct {
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapBackstopGainsClaimedEventName = "GainsClaimed"

// ContractEventName returns the user-defined event name.
func (GapBackstopGainsClaimed) ContractEventName() string {
	return GapBackstopGainsClaimedEventName
}

// UnpackGainsClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GainsClaimed(address indexed owner, address indexed token, uint256 amount)
func (gapBackstop *GapBackstop) UnpackGainsClaimedEvent(log *types.Log) (*GapBackstopGainsClaimed, error) {
	event := "GainsClaimed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopGainsClaimed)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the GapBackstop contract.
type GapBackstopOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const GapBackstopOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (GapBackstopOwnershipTransferStarted) ContractEventName() string {
	return GapBackstopOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (gapBackstop *GapBackstop) UnpackOwnershipTransferStartedEvent(log *types.Log) (*GapBackstopOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopOwnershipTransferred represents a OwnershipTransferred event raised by the GapBackstop contract.
type GapBackstopOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const GapBackstopOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (GapBackstopOwnershipTransferred) ContractEventName() string {
	return GapBackstopOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (gapBackstop *GapBackstop) UnpackOwnershipTransferredEvent(log *types.Log) (*GapBackstopOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopPremiumAdded represents a PremiumAdded event raised by the GapBackstop contract.
type GapBackstopPremiumAdded struct {
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapBackstopPremiumAddedEventName = "PremiumAdded"

// ContractEventName returns the user-defined event name.
func (GapBackstopPremiumAdded) ContractEventName() string {
	return GapBackstopPremiumAddedEventName
}

// UnpackPremiumAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PremiumAdded(uint256 amount)
func (gapBackstop *GapBackstop) UnpackPremiumAddedEvent(log *types.Log) (*GapBackstopPremiumAdded, error) {
	event := "PremiumAdded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopPremiumAdded)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopRemainderBought represents a RemainderBought event raised by the GapBackstop contract.
type GapBackstopRemainderBought struct {
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Cost     *big.Int
	ClosesMs uint64
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapBackstopRemainderBoughtEventName = "RemainderBought"

// ContractEventName returns the user-defined event name.
func (GapBackstopRemainderBought) ContractEventName() string {
	return GapBackstopRemainderBoughtEventName
}

// UnpackRemainderBoughtEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RemainderBought(address indexed account, bytes32 indexed position, address indexed token, uint256 amount, uint256 cost, uint64 closesMs)
func (gapBackstop *GapBackstop) UnpackRemainderBoughtEvent(log *types.Log) (*GapBackstopRemainderBought, error) {
	event := "RemainderBought"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopRemainderBought)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopSync represents a Sync event raised by the GapBackstop contract.
type GapBackstopSync struct {
	Held *big.Int
	Raw  *types.Log // Blockchain specific contextual infos
}

const GapBackstopSyncEventName = "Sync"

// ContractEventName returns the user-defined event name.
func (GapBackstopSync) ContractEventName() string {
	return GapBackstopSyncEventName
}

// UnpackSyncEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sync(uint256 held)
func (gapBackstop *GapBackstop) UnpackSyncEvent(log *types.Log) (*GapBackstopSync, error) {
	event := "Sync"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopSync)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopTransfer represents a Transfer event raised by the GapBackstop contract.
type GapBackstopTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const GapBackstopTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (GapBackstopTransfer) ContractEventName() string {
	return GapBackstopTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (gapBackstop *GapBackstop) UnpackTransferEvent(log *types.Log) (*GapBackstopTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopTransfer)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// GapBackstopWithdraw represents a Withdraw event raised by the GapBackstop contract.
type GapBackstopWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapBackstopWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (GapBackstopWithdraw) ContractEventName() string {
	return GapBackstopWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (gapBackstop *GapBackstop) UnpackWithdrawEvent(log *types.Log) (*GapBackstopWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapBackstop.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapBackstopWithdraw)
	if len(log.Data) > 0 {
		if err := gapBackstop.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapBackstop.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (gapBackstop *GapBackstop) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC4626ExceededMaxDeposit"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC4626ExceededMaxDepositError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC4626ExceededMaxMint"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC4626ExceededMaxMintError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC4626ExceededMaxRedeem"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC4626ExceededMaxRedeemError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["ERC4626ExceededMaxWithdraw"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackERC4626ExceededMaxWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["InvalidExposureLimits"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackInvalidExposureLimitsError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["InvalidPosition"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackInvalidPositionError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["NotAuction"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackNotAuctionError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapBackstop.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return gapBackstop.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// GapBackstopERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the GapBackstop contract.
type GapBackstopERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func GapBackstopERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (gapBackstop *GapBackstop) UnpackERC20InsufficientAllowanceError(raw []byte) (*GapBackstopERC20InsufficientAllowance, error) {
	out := new(GapBackstopERC20InsufficientAllowance)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the GapBackstop contract.
type GapBackstopERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func GapBackstopERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (gapBackstop *GapBackstop) UnpackERC20InsufficientBalanceError(raw []byte) (*GapBackstopERC20InsufficientBalance, error) {
	out := new(GapBackstopERC20InsufficientBalance)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC20InvalidApprover represents a ERC20InvalidApprover error raised by the GapBackstop contract.
type GapBackstopERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func GapBackstopERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (gapBackstop *GapBackstop) UnpackERC20InvalidApproverError(raw []byte) (*GapBackstopERC20InvalidApprover, error) {
	out := new(GapBackstopERC20InvalidApprover)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the GapBackstop contract.
type GapBackstopERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func GapBackstopERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (gapBackstop *GapBackstop) UnpackERC20InvalidReceiverError(raw []byte) (*GapBackstopERC20InvalidReceiver, error) {
	out := new(GapBackstopERC20InvalidReceiver)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC20InvalidSender represents a ERC20InvalidSender error raised by the GapBackstop contract.
type GapBackstopERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func GapBackstopERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (gapBackstop *GapBackstop) UnpackERC20InvalidSenderError(raw []byte) (*GapBackstopERC20InvalidSender, error) {
	out := new(GapBackstopERC20InvalidSender)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC20InvalidSpender represents a ERC20InvalidSpender error raised by the GapBackstop contract.
type GapBackstopERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func GapBackstopERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (gapBackstop *GapBackstop) UnpackERC20InvalidSpenderError(raw []byte) (*GapBackstopERC20InvalidSpender, error) {
	out := new(GapBackstopERC20InvalidSpender)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC4626ExceededMaxDeposit represents a ERC4626ExceededMaxDeposit error raised by the GapBackstop contract.
type GapBackstopERC4626ExceededMaxDeposit struct {
	Receiver common.Address
	Assets   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func GapBackstopERC4626ExceededMaxDepositErrorID() common.Hash {
	return common.HexToHash("0x79012fb2819fcdc1de669c08773dbcd6bdc757862642f75fb1c584dadf259dfe")
}

// UnpackERC4626ExceededMaxDepositError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func (gapBackstop *GapBackstop) UnpackERC4626ExceededMaxDepositError(raw []byte) (*GapBackstopERC4626ExceededMaxDeposit, error) {
	out := new(GapBackstopERC4626ExceededMaxDeposit)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxDeposit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC4626ExceededMaxMint represents a ERC4626ExceededMaxMint error raised by the GapBackstop contract.
type GapBackstopERC4626ExceededMaxMint struct {
	Receiver common.Address
	Shares   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func GapBackstopERC4626ExceededMaxMintErrorID() common.Hash {
	return common.HexToHash("0x284ff667dc615a39438518c22e8955b9470327d9de8a4d7e21c926b260d65176")
}

// UnpackERC4626ExceededMaxMintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func (gapBackstop *GapBackstop) UnpackERC4626ExceededMaxMintError(raw []byte) (*GapBackstopERC4626ExceededMaxMint, error) {
	out := new(GapBackstopERC4626ExceededMaxMint)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxMint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC4626ExceededMaxRedeem represents a ERC4626ExceededMaxRedeem error raised by the GapBackstop contract.
type GapBackstopERC4626ExceededMaxRedeem struct {
	Owner  common.Address
	Shares *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func GapBackstopERC4626ExceededMaxRedeemErrorID() common.Hash {
	return common.HexToHash("0xb94abeec0557d36b5f0bc8f115deec7b184dcbff94ac66d55e37c8f301e75269")
}

// UnpackERC4626ExceededMaxRedeemError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func (gapBackstop *GapBackstop) UnpackERC4626ExceededMaxRedeemError(raw []byte) (*GapBackstopERC4626ExceededMaxRedeem, error) {
	out := new(GapBackstopERC4626ExceededMaxRedeem)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxRedeem", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopERC4626ExceededMaxWithdraw represents a ERC4626ExceededMaxWithdraw error raised by the GapBackstop contract.
type GapBackstopERC4626ExceededMaxWithdraw struct {
	Owner  common.Address
	Assets *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func GapBackstopERC4626ExceededMaxWithdrawErrorID() common.Hash {
	return common.HexToHash("0xfe9cceec2bd1f9b68641914cc354eacaeb1cc2169f5ba0639930f241e87142f0")
}

// UnpackERC4626ExceededMaxWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func (gapBackstop *GapBackstop) UnpackERC4626ExceededMaxWithdrawError(raw []byte) (*GapBackstopERC4626ExceededMaxWithdraw, error) {
	out := new(GapBackstopERC4626ExceededMaxWithdraw)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopInvalidExposureLimits represents a InvalidExposureLimits error raised by the GapBackstop contract.
type GapBackstopInvalidExposureLimits struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidExposureLimits()
func GapBackstopInvalidExposureLimitsErrorID() common.Hash {
	return common.HexToHash("0x25d62b8197ee1e494151d428aa2100b20045638dbccd374282d72c48bd0fa7dd")
}

// UnpackInvalidExposureLimitsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidExposureLimits()
func (gapBackstop *GapBackstop) UnpackInvalidExposureLimitsError(raw []byte) (*GapBackstopInvalidExposureLimits, error) {
	out := new(GapBackstopInvalidExposureLimits)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "InvalidExposureLimits", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopInvalidPosition represents a InvalidPosition error raised by the GapBackstop contract.
type GapBackstopInvalidPosition struct {
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPosition(bytes32 position)
func GapBackstopInvalidPositionErrorID() common.Hash {
	return common.HexToHash("0x8ea9158f98eb1f30e0e3d0c2079e9b5d3d09ce66836abe2f020eb75644263da9")
}

// UnpackInvalidPositionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPosition(bytes32 position)
func (gapBackstop *GapBackstop) UnpackInvalidPositionError(raw []byte) (*GapBackstopInvalidPosition, error) {
	out := new(GapBackstopInvalidPosition)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "InvalidPosition", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopNotAuction represents a NotAuction error raised by the GapBackstop contract.
type GapBackstopNotAuction struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAuction(address caller)
func GapBackstopNotAuctionErrorID() common.Hash {
	return common.HexToHash("0xf459dd89b92e732bb4d2d5f95e8b4b5405c659d120e94467be4b958332750795")
}

// UnpackNotAuctionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAuction(address caller)
func (gapBackstop *GapBackstop) UnpackNotAuctionError(raw []byte) (*GapBackstopNotAuction, error) {
	out := new(GapBackstopNotAuction)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "NotAuction", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the GapBackstop contract.
type GapBackstopOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func GapBackstopOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (gapBackstop *GapBackstop) UnpackOwnableInvalidOwnerError(raw []byte) (*GapBackstopOwnableInvalidOwner, error) {
	out := new(GapBackstopOwnableInvalidOwner)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the GapBackstop contract.
type GapBackstopOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func GapBackstopOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (gapBackstop *GapBackstop) UnpackOwnableUnauthorizedAccountError(raw []byte) (*GapBackstopOwnableUnauthorizedAccount, error) {
	out := new(GapBackstopOwnableUnauthorizedAccount)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the GapBackstop contract.
type GapBackstopOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func GapBackstopOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (gapBackstop *GapBackstop) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*GapBackstopOwnershipCannotBeRenounced, error) {
	out := new(GapBackstopOwnershipCannotBeRenounced)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the GapBackstop contract.
type GapBackstopSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func GapBackstopSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (gapBackstop *GapBackstop) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*GapBackstopSafeCastOverflowedUintDowncast, error) {
	out := new(GapBackstopSafeCastOverflowedUintDowncast)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapBackstopSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the GapBackstop contract.
type GapBackstopSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func GapBackstopSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (gapBackstop *GapBackstop) UnpackSafeERC20FailedOperationError(raw []byte) (*GapBackstopSafeERC20FailedOperation, error) {
	out := new(GapBackstopSafeERC20FailedOperation)
	if err := gapBackstop.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}
