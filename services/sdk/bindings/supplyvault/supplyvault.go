// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package supplyvault

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

// SupplyVaultRateModel is an auto generated low-level Go binding around an user-defined struct.
type SupplyVaultRateModel struct {
	Optimal uint16
	Base    uint32
	Slope1  uint32
	Slope2  uint32
}

// SupplyVaultMetaData contains all meta data concerning the SupplyVault contract.
var SupplyVaultMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"usdg\",\"type\":\"address\",\"internalType\":\"contractIUSDG\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialRateModel\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_OPTIMAL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_RATE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_OPTIMAL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrowIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrowRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrower\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToAssets\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToShares\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"debt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"idle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastAccrual\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewDeposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rateModel\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"repaid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scaledDebt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBorrower\",\"inputs\":[{\"name\":\"newBorrower\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateModel\",\"inputs\":[{\"name\":\"model\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sync\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"utilization\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"writeOff\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"written\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Borrow\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BorrowerSet\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InterestAccrued\",\"inputs\":[{\"name\":\"interest\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateModelSet\",\"inputs\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Repay\",\"inputs\":[{\"name\":\"payer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sync\",\"inputs\":[{\"name\":\"idle\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteOff\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BorrowerAlreadySet\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLiquidity\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBorrower\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRateModel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotBorrower\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "SupplyVault",
}

// SupplyVault is an auto generated Go binding around an Ethereum contract.
type SupplyVault struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *SupplyVault) GetABI() abi.ABI {
	return c.abi
}

// NewSupplyVault creates a new instance of SupplyVault.
func NewSupplyVault() *SupplyVault {
	parsed, err := SupplyVaultMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &SupplyVault{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *SupplyVault) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address usdg, address initialOwner, (uint16,uint32,uint32,uint32) initialRateModel) returns()
func (supplyVault *SupplyVault) PackConstructor(usdg common.Address, initialOwner common.Address, initialRateModel SupplyVaultRateModel) []byte {
	enc, err := supplyVault.abi.Pack("", usdg, initialOwner, initialRateModel)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMAXOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8718699e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) PackMAXOPTIMAL() []byte {
	enc, err := supplyVault.abi.Pack("MAX_OPTIMAL")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8718699e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) TryPackMAXOPTIMAL() ([]byte, error) {
	return supplyVault.abi.Pack("MAX_OPTIMAL")
}

// UnpackMAXOPTIMAL is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8718699e.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) UnpackMAXOPTIMAL(data []byte) (uint16, error) {
	out, err := supplyVault.abi.Unpack("MAX_OPTIMAL", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackMAXRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc24dbebd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (supplyVault *SupplyVault) PackMAXRATE() []byte {
	enc, err := supplyVault.abi.Pack("MAX_RATE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc24dbebd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (supplyVault *SupplyVault) TryPackMAXRATE() ([]byte, error) {
	return supplyVault.abi.Pack("MAX_RATE")
}

// UnpackMAXRATE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc24dbebd.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (supplyVault *SupplyVault) UnpackMAXRATE(data []byte) (uint32, error) {
	out, err := supplyVault.abi.Unpack("MAX_RATE", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackMINOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a1bc24d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) PackMINOPTIMAL() []byte {
	enc, err := supplyVault.abi.Pack("MIN_OPTIMAL")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a1bc24d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) TryPackMINOPTIMAL() ([]byte, error) {
	return supplyVault.abi.Pack("MIN_OPTIMAL")
}

// UnpackMINOPTIMAL is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a1bc24d.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (supplyVault *SupplyVault) UnpackMINOPTIMAL(data []byte) (uint16, error) {
	out, err := supplyVault.abi.Unpack("MIN_OPTIMAL", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (supplyVault *SupplyVault) PackAcceptOwnership() []byte {
	enc, err := supplyVault.abi.Pack("acceptOwnership")
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
func (supplyVault *SupplyVault) TryPackAcceptOwnership() ([]byte, error) {
	return supplyVault.abi.Pack("acceptOwnership")
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (supplyVault *SupplyVault) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := supplyVault.abi.Pack("allowance", owner, spender)
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
func (supplyVault *SupplyVault) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (supplyVault *SupplyVault) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("allowance", data)
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
func (supplyVault *SupplyVault) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("approve", spender, value)
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
func (supplyVault *SupplyVault) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (supplyVault *SupplyVault) UnpackApprove(data []byte) (bool, error) {
	out, err := supplyVault.abi.Unpack("approve", data)
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
func (supplyVault *SupplyVault) PackAsset() []byte {
	enc, err := supplyVault.abi.Pack("asset")
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
func (supplyVault *SupplyVault) TryPackAsset() ([]byte, error) {
	return supplyVault.abi.Pack("asset")
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (supplyVault *SupplyVault) UnpackAsset(data []byte) (common.Address, error) {
	out, err := supplyVault.abi.Unpack("asset", data)
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
func (supplyVault *SupplyVault) PackBalanceOf(account common.Address) []byte {
	enc, err := supplyVault.abi.Pack("balanceOf", account)
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
func (supplyVault *SupplyVault) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (supplyVault *SupplyVault) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b3fd148.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrow(uint256 assets, address receiver) returns()
func (supplyVault *SupplyVault) PackBorrow(assets *big.Int, receiver common.Address) []byte {
	enc, err := supplyVault.abi.Pack("borrow", assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b3fd148.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrow(uint256 assets, address receiver) returns()
func (supplyVault *SupplyVault) TryPackBorrow(assets *big.Int, receiver common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("borrow", assets, receiver)
}

// PackBorrowIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa5af0fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowIndex() view returns(uint128)
func (supplyVault *SupplyVault) PackBorrowIndex() []byte {
	enc, err := supplyVault.abi.Pack("borrowIndex")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa5af0fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowIndex() view returns(uint128)
func (supplyVault *SupplyVault) TryPackBorrowIndex() ([]byte, error) {
	return supplyVault.abi.Pack("borrowIndex")
}

// UnpackBorrowIndex is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa5af0fd.
//
// Solidity: function borrowIndex() view returns(uint128)
func (supplyVault *SupplyVault) UnpackBorrowIndex(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("borrowIndex", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrowRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc914b437.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowRate() view returns(uint256)
func (supplyVault *SupplyVault) PackBorrowRate() []byte {
	enc, err := supplyVault.abi.Pack("borrowRate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc914b437.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowRate() view returns(uint256)
func (supplyVault *SupplyVault) TryPackBorrowRate() ([]byte, error) {
	return supplyVault.abi.Pack("borrowRate")
}

// UnpackBorrowRate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc914b437.
//
// Solidity: function borrowRate() view returns(uint256)
func (supplyVault *SupplyVault) UnpackBorrowRate(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("borrowRate", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7df1f1b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrower() view returns(address)
func (supplyVault *SupplyVault) PackBorrower() []byte {
	enc, err := supplyVault.abi.Pack("borrower")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7df1f1b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrower() view returns(address)
func (supplyVault *SupplyVault) TryPackBorrower() ([]byte, error) {
	return supplyVault.abi.Pack("borrower")
}

// UnpackBorrower is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7df1f1b9.
//
// Solidity: function borrower() view returns(address)
func (supplyVault *SupplyVault) UnpackBorrower(data []byte) (common.Address, error) {
	out, err := supplyVault.abi.Unpack("borrower", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (supplyVault *SupplyVault) PackConvertToAssets(shares *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("convertToAssets", shares)
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
func (supplyVault *SupplyVault) TryPackConvertToAssets(shares *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("convertToAssets", shares)
}

// UnpackConvertToAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (supplyVault *SupplyVault) UnpackConvertToAssets(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("convertToAssets", data)
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
func (supplyVault *SupplyVault) PackConvertToShares(assets *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("convertToShares", assets)
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
func (supplyVault *SupplyVault) TryPackConvertToShares(assets *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("convertToShares", assets)
}

// UnpackConvertToShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (supplyVault *SupplyVault) UnpackConvertToShares(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("convertToShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dca59c1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function debt() view returns(uint256)
func (supplyVault *SupplyVault) PackDebt() []byte {
	enc, err := supplyVault.abi.Pack("debt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dca59c1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function debt() view returns(uint256)
func (supplyVault *SupplyVault) TryPackDebt() ([]byte, error) {
	return supplyVault.abi.Pack("debt")
}

// UnpackDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dca59c1.
//
// Solidity: function debt() view returns(uint256)
func (supplyVault *SupplyVault) UnpackDebt(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("debt", data)
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
func (supplyVault *SupplyVault) PackDecimals() []byte {
	enc, err := supplyVault.abi.Pack("decimals")
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
func (supplyVault *SupplyVault) TryPackDecimals() ([]byte, error) {
	return supplyVault.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (supplyVault *SupplyVault) UnpackDecimals(data []byte) (uint8, error) {
	out, err := supplyVault.abi.Unpack("decimals", data)
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
func (supplyVault *SupplyVault) PackDeposit(assets *big.Int, receiver common.Address) []byte {
	enc, err := supplyVault.abi.Pack("deposit", assets, receiver)
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
func (supplyVault *SupplyVault) TryPackDeposit(assets *big.Int, receiver common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("deposit", assets, receiver)
}

// UnpackDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (supplyVault *SupplyVault) UnpackDeposit(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("deposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDepositWithPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50921b23.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositWithPermit(uint256 assets, address receiver, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 shares)
func (supplyVault *SupplyVault) PackDepositWithPermit(assets *big.Int, receiver common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) []byte {
	enc, err := supplyVault.abi.Pack("depositWithPermit", assets, receiver, deadline, v, r, s)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositWithPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50921b23.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositWithPermit(uint256 assets, address receiver, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 shares)
func (supplyVault *SupplyVault) TryPackDepositWithPermit(assets *big.Int, receiver common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) ([]byte, error) {
	return supplyVault.abi.Pack("depositWithPermit", assets, receiver, deadline, v, r, s)
}

// UnpackDepositWithPermit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x50921b23.
//
// Solidity: function depositWithPermit(uint256 assets, address receiver, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 shares)
func (supplyVault *SupplyVault) UnpackDepositWithPermit(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("depositWithPermit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackIdle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3192164f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function idle() view returns(uint128)
func (supplyVault *SupplyVault) PackIdle() []byte {
	enc, err := supplyVault.abi.Pack("idle")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIdle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3192164f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function idle() view returns(uint128)
func (supplyVault *SupplyVault) TryPackIdle() ([]byte, error) {
	return supplyVault.abi.Pack("idle")
}

// UnpackIdle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3192164f.
//
// Solidity: function idle() view returns(uint128)
func (supplyVault *SupplyVault) UnpackIdle(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("idle", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLastAccrual is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b3baab4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastAccrual() view returns(uint64)
func (supplyVault *SupplyVault) PackLastAccrual() []byte {
	enc, err := supplyVault.abi.Pack("lastAccrual")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastAccrual is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b3baab4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastAccrual() view returns(uint64)
func (supplyVault *SupplyVault) TryPackLastAccrual() ([]byte, error) {
	return supplyVault.abi.Pack("lastAccrual")
}

// UnpackLastAccrual is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7b3baab4.
//
// Solidity: function lastAccrual() view returns(uint64)
func (supplyVault *SupplyVault) UnpackLastAccrual(data []byte) (uint64, error) {
	out, err := supplyVault.abi.Unpack("lastAccrual", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (supplyVault *SupplyVault) PackMaxDeposit(receiver common.Address) []byte {
	enc, err := supplyVault.abi.Pack("maxDeposit", receiver)
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
func (supplyVault *SupplyVault) TryPackMaxDeposit(receiver common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("maxDeposit", receiver)
}

// UnpackMaxDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (supplyVault *SupplyVault) UnpackMaxDeposit(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("maxDeposit", data)
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
func (supplyVault *SupplyVault) PackMaxMint(receiver common.Address) []byte {
	enc, err := supplyVault.abi.Pack("maxMint", receiver)
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
func (supplyVault *SupplyVault) TryPackMaxMint(receiver common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("maxMint", receiver)
}

// UnpackMaxMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (supplyVault *SupplyVault) UnpackMaxMint(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("maxMint", data)
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
func (supplyVault *SupplyVault) PackMaxRedeem(owner common.Address) []byte {
	enc, err := supplyVault.abi.Pack("maxRedeem", owner)
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
func (supplyVault *SupplyVault) TryPackMaxRedeem(owner common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("maxRedeem", owner)
}

// UnpackMaxRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (supplyVault *SupplyVault) UnpackMaxRedeem(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("maxRedeem", data)
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
func (supplyVault *SupplyVault) PackMaxWithdraw(owner common.Address) []byte {
	enc, err := supplyVault.abi.Pack("maxWithdraw", owner)
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
func (supplyVault *SupplyVault) TryPackMaxWithdraw(owner common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("maxWithdraw", owner)
}

// UnpackMaxWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (supplyVault *SupplyVault) UnpackMaxWithdraw(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("maxWithdraw", data)
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
func (supplyVault *SupplyVault) PackMint(shares *big.Int, receiver common.Address) []byte {
	enc, err := supplyVault.abi.Pack("mint", shares, receiver)
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
func (supplyVault *SupplyVault) TryPackMint(shares *big.Int, receiver common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("mint", shares, receiver)
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (supplyVault *SupplyVault) UnpackMint(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("mint", data)
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
func (supplyVault *SupplyVault) PackName() []byte {
	enc, err := supplyVault.abi.Pack("name")
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
func (supplyVault *SupplyVault) TryPackName() ([]byte, error) {
	return supplyVault.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (supplyVault *SupplyVault) UnpackName(data []byte) (string, error) {
	out, err := supplyVault.abi.Unpack("name", data)
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
func (supplyVault *SupplyVault) PackOwner() []byte {
	enc, err := supplyVault.abi.Pack("owner")
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
func (supplyVault *SupplyVault) TryPackOwner() ([]byte, error) {
	return supplyVault.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (supplyVault *SupplyVault) UnpackOwner(data []byte) (common.Address, error) {
	out, err := supplyVault.abi.Unpack("owner", data)
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
func (supplyVault *SupplyVault) PackPendingOwner() []byte {
	enc, err := supplyVault.abi.Pack("pendingOwner")
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
func (supplyVault *SupplyVault) TryPackPendingOwner() ([]byte, error) {
	return supplyVault.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (supplyVault *SupplyVault) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := supplyVault.abi.Unpack("pendingOwner", data)
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
func (supplyVault *SupplyVault) PackPreviewDeposit(assets *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("previewDeposit", assets)
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
func (supplyVault *SupplyVault) TryPackPreviewDeposit(assets *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("previewDeposit", assets)
}

// UnpackPreviewDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (supplyVault *SupplyVault) UnpackPreviewDeposit(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("previewDeposit", data)
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
func (supplyVault *SupplyVault) PackPreviewMint(shares *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("previewMint", shares)
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
func (supplyVault *SupplyVault) TryPackPreviewMint(shares *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("previewMint", shares)
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (supplyVault *SupplyVault) UnpackPreviewMint(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("previewMint", data)
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
func (supplyVault *SupplyVault) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("previewRedeem", shares)
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
func (supplyVault *SupplyVault) TryPackPreviewRedeem(shares *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("previewRedeem", shares)
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (supplyVault *SupplyVault) UnpackPreviewRedeem(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("previewRedeem", data)
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
func (supplyVault *SupplyVault) PackPreviewWithdraw(assets *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("previewWithdraw", assets)
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
func (supplyVault *SupplyVault) TryPackPreviewWithdraw(assets *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("previewWithdraw", assets)
}

// UnpackPreviewWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (supplyVault *SupplyVault) UnpackPreviewWithdraw(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("previewWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1088459.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (supplyVault *SupplyVault) PackRateModel() []byte {
	enc, err := supplyVault.abi.Pack("rateModel")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1088459.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (supplyVault *SupplyVault) TryPackRateModel() ([]byte, error) {
	return supplyVault.abi.Pack("rateModel")
}

// UnpackRateModel is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1088459.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (supplyVault *SupplyVault) UnpackRateModel(data []byte) (SupplyVaultRateModel, error) {
	out, err := supplyVault.abi.Unpack("rateModel", data)
	if err != nil {
		return *new(SupplyVaultRateModel), err
	}
	out0 := *abi.ConvertType(out[0], new(SupplyVaultRateModel)).(*SupplyVaultRateModel)
	return out0, nil
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (supplyVault *SupplyVault) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := supplyVault.abi.Pack("redeem", shares, receiver, owner)
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
func (supplyVault *SupplyVault) TryPackRedeem(shares *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("redeem", shares, receiver, owner)
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (supplyVault *SupplyVault) UnpackRedeem(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("redeem", data)
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
func (supplyVault *SupplyVault) PackRenounceOwnership() []byte {
	enc, err := supplyVault.abi.Pack("renounceOwnership")
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
func (supplyVault *SupplyVault) TryPackRenounceOwnership() ([]byte, error) {
	return supplyVault.abi.Pack("renounceOwnership")
}

// PackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x371fd8e6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (supplyVault *SupplyVault) PackRepay(assets *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("repay", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x371fd8e6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (supplyVault *SupplyVault) TryPackRepay(assets *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("repay", assets)
}

// UnpackRepay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x371fd8e6.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (supplyVault *SupplyVault) UnpackRepay(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("repay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackScaledDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf7bac74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scaledDebt() view returns(uint128)
func (supplyVault *SupplyVault) PackScaledDebt() []byte {
	enc, err := supplyVault.abi.Pack("scaledDebt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScaledDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf7bac74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scaledDebt() view returns(uint128)
func (supplyVault *SupplyVault) TryPackScaledDebt() ([]byte, error) {
	return supplyVault.abi.Pack("scaledDebt")
}

// UnpackScaledDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf7bac74.
//
// Solidity: function scaledDebt() view returns(uint128)
func (supplyVault *SupplyVault) UnpackScaledDebt(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("scaledDebt", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc762d5f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBorrower(address newBorrower) returns()
func (supplyVault *SupplyVault) PackSetBorrower(newBorrower common.Address) []byte {
	enc, err := supplyVault.abi.Pack("setBorrower", newBorrower)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc762d5f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBorrower(address newBorrower) returns()
func (supplyVault *SupplyVault) TryPackSetBorrower(newBorrower common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("setBorrower", newBorrower)
}

// PackSetRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x043319d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setRateModel((uint16,uint32,uint32,uint32) model) returns()
func (supplyVault *SupplyVault) PackSetRateModel(model SupplyVaultRateModel) []byte {
	enc, err := supplyVault.abi.Pack("setRateModel", model)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x043319d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setRateModel((uint16,uint32,uint32,uint32) model) returns()
func (supplyVault *SupplyVault) TryPackSetRateModel(model SupplyVaultRateModel) ([]byte, error) {
	return supplyVault.abi.Pack("setRateModel", model)
}

// PackSupplyRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad2961a3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supplyRate() view returns(uint256)
func (supplyVault *SupplyVault) PackSupplyRate() []byte {
	enc, err := supplyVault.abi.Pack("supplyRate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupplyRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad2961a3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supplyRate() view returns(uint256)
func (supplyVault *SupplyVault) TryPackSupplyRate() ([]byte, error) {
	return supplyVault.abi.Pack("supplyRate")
}

// UnpackSupplyRate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad2961a3.
//
// Solidity: function supplyRate() view returns(uint256)
func (supplyVault *SupplyVault) UnpackSupplyRate(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("supplyRate", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (supplyVault *SupplyVault) PackSymbol() []byte {
	enc, err := supplyVault.abi.Pack("symbol")
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
func (supplyVault *SupplyVault) TryPackSymbol() ([]byte, error) {
	return supplyVault.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (supplyVault *SupplyVault) UnpackSymbol(data []byte) (string, error) {
	out, err := supplyVault.abi.Unpack("symbol", data)
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
func (supplyVault *SupplyVault) PackSync() []byte {
	enc, err := supplyVault.abi.Pack("sync")
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
func (supplyVault *SupplyVault) TryPackSync() ([]byte, error) {
	return supplyVault.abi.Pack("sync")
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAssets() view returns(uint256)
func (supplyVault *SupplyVault) PackTotalAssets() []byte {
	enc, err := supplyVault.abi.Pack("totalAssets")
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
func (supplyVault *SupplyVault) TryPackTotalAssets() ([]byte, error) {
	return supplyVault.abi.Pack("totalAssets")
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (supplyVault *SupplyVault) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("totalAssets", data)
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
func (supplyVault *SupplyVault) PackTotalSupply() []byte {
	enc, err := supplyVault.abi.Pack("totalSupply")
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
func (supplyVault *SupplyVault) TryPackTotalSupply() ([]byte, error) {
	return supplyVault.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (supplyVault *SupplyVault) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("totalSupply", data)
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
func (supplyVault *SupplyVault) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("transfer", to, value)
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
func (supplyVault *SupplyVault) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (supplyVault *SupplyVault) UnpackTransfer(data []byte) (bool, error) {
	out, err := supplyVault.abi.Unpack("transfer", data)
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
func (supplyVault *SupplyVault) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("transferFrom", from, to, value)
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
func (supplyVault *SupplyVault) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (supplyVault *SupplyVault) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := supplyVault.abi.Unpack("transferFrom", data)
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
func (supplyVault *SupplyVault) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := supplyVault.abi.Pack("transferOwnership", newOwner)
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
func (supplyVault *SupplyVault) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("transferOwnership", newOwner)
}

// PackUtilization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea21cd92.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function utilization() view returns(uint256)
func (supplyVault *SupplyVault) PackUtilization() []byte {
	enc, err := supplyVault.abi.Pack("utilization")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUtilization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea21cd92.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function utilization() view returns(uint256)
func (supplyVault *SupplyVault) TryPackUtilization() ([]byte, error) {
	return supplyVault.abi.Pack("utilization")
}

// UnpackUtilization is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xea21cd92.
//
// Solidity: function utilization() view returns(uint256)
func (supplyVault *SupplyVault) UnpackUtilization(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("utilization", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (supplyVault *SupplyVault) PackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := supplyVault.abi.Pack("withdraw", assets, receiver, owner)
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
func (supplyVault *SupplyVault) TryPackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return supplyVault.abi.Pack("withdraw", assets, receiver, owner)
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (supplyVault *SupplyVault) UnpackWithdraw(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("withdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b3e2291.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (supplyVault *SupplyVault) PackWriteOff(assets *big.Int) []byte {
	enc, err := supplyVault.abi.Pack("writeOff", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b3e2291.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (supplyVault *SupplyVault) TryPackWriteOff(assets *big.Int) ([]byte, error) {
	return supplyVault.abi.Pack("writeOff", assets)
}

// UnpackWriteOff is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9b3e2291.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (supplyVault *SupplyVault) UnpackWriteOff(data []byte) (*big.Int, error) {
	out, err := supplyVault.abi.Unpack("writeOff", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// SupplyVaultApproval represents a Approval event raised by the SupplyVault contract.
type SupplyVaultApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const SupplyVaultApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (SupplyVaultApproval) ContractEventName() string {
	return SupplyVaultApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (supplyVault *SupplyVault) UnpackApprovalEvent(log *types.Log) (*SupplyVaultApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultApproval)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultBorrow represents a Borrow event raised by the SupplyVault contract.
type SupplyVaultBorrow struct {
	Receiver common.Address
	Assets   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const SupplyVaultBorrowEventName = "Borrow"

// ContractEventName returns the user-defined event name.
func (SupplyVaultBorrow) ContractEventName() string {
	return SupplyVaultBorrowEventName
}

// UnpackBorrowEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Borrow(address indexed receiver, uint256 assets)
func (supplyVault *SupplyVault) UnpackBorrowEvent(log *types.Log) (*SupplyVaultBorrow, error) {
	event := "Borrow"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultBorrow)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultBorrowerSet represents a BorrowerSet event raised by the SupplyVault contract.
type SupplyVaultBorrowerSet struct {
	Borrower common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const SupplyVaultBorrowerSetEventName = "BorrowerSet"

// ContractEventName returns the user-defined event name.
func (SupplyVaultBorrowerSet) ContractEventName() string {
	return SupplyVaultBorrowerSetEventName
}

// UnpackBorrowerSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BorrowerSet(address indexed borrower)
func (supplyVault *SupplyVault) UnpackBorrowerSetEvent(log *types.Log) (*SupplyVaultBorrowerSet, error) {
	event := "BorrowerSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultBorrowerSet)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultDeposit represents a Deposit event raised by the SupplyVault contract.
type SupplyVaultDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const SupplyVaultDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (SupplyVaultDeposit) ContractEventName() string {
	return SupplyVaultDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (supplyVault *SupplyVault) UnpackDepositEvent(log *types.Log) (*SupplyVaultDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultDeposit)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultInterestAccrued represents a InterestAccrued event raised by the SupplyVault contract.
type SupplyVaultInterestAccrued struct {
	Interest *big.Int
	Index    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const SupplyVaultInterestAccruedEventName = "InterestAccrued"

// ContractEventName returns the user-defined event name.
func (SupplyVaultInterestAccrued) ContractEventName() string {
	return SupplyVaultInterestAccruedEventName
}

// UnpackInterestAccruedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event InterestAccrued(uint256 interest, uint256 index)
func (supplyVault *SupplyVault) UnpackInterestAccruedEvent(log *types.Log) (*SupplyVaultInterestAccrued, error) {
	event := "InterestAccrued"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultInterestAccrued)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the SupplyVault contract.
type SupplyVaultOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const SupplyVaultOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (SupplyVaultOwnershipTransferStarted) ContractEventName() string {
	return SupplyVaultOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (supplyVault *SupplyVault) UnpackOwnershipTransferStartedEvent(log *types.Log) (*SupplyVaultOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultOwnershipTransferred represents a OwnershipTransferred event raised by the SupplyVault contract.
type SupplyVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const SupplyVaultOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (SupplyVaultOwnershipTransferred) ContractEventName() string {
	return SupplyVaultOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (supplyVault *SupplyVault) UnpackOwnershipTransferredEvent(log *types.Log) (*SupplyVaultOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultRateModelSet represents a RateModelSet event raised by the SupplyVault contract.
type SupplyVaultRateModelSet struct {
	Optimal uint16
	Base    uint32
	Slope1  uint32
	Slope2  uint32
	Raw     *types.Log // Blockchain specific contextual infos
}

const SupplyVaultRateModelSetEventName = "RateModelSet"

// ContractEventName returns the user-defined event name.
func (SupplyVaultRateModelSet) ContractEventName() string {
	return SupplyVaultRateModelSetEventName
}

// UnpackRateModelSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RateModelSet(uint16 optimal, uint32 base, uint32 slope1, uint32 slope2)
func (supplyVault *SupplyVault) UnpackRateModelSetEvent(log *types.Log) (*SupplyVaultRateModelSet, error) {
	event := "RateModelSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultRateModelSet)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultRepay represents a Repay event raised by the SupplyVault contract.
type SupplyVaultRepay struct {
	Payer  common.Address
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const SupplyVaultRepayEventName = "Repay"

// ContractEventName returns the user-defined event name.
func (SupplyVaultRepay) ContractEventName() string {
	return SupplyVaultRepayEventName
}

// UnpackRepayEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Repay(address indexed payer, uint256 assets)
func (supplyVault *SupplyVault) UnpackRepayEvent(log *types.Log) (*SupplyVaultRepay, error) {
	event := "Repay"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultRepay)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultSync represents a Sync event raised by the SupplyVault contract.
type SupplyVaultSync struct {
	Idle *big.Int
	Raw  *types.Log // Blockchain specific contextual infos
}

const SupplyVaultSyncEventName = "Sync"

// ContractEventName returns the user-defined event name.
func (SupplyVaultSync) ContractEventName() string {
	return SupplyVaultSyncEventName
}

// UnpackSyncEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sync(uint256 idle)
func (supplyVault *SupplyVault) UnpackSyncEvent(log *types.Log) (*SupplyVaultSync, error) {
	event := "Sync"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultSync)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultTransfer represents a Transfer event raised by the SupplyVault contract.
type SupplyVaultTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const SupplyVaultTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (SupplyVaultTransfer) ContractEventName() string {
	return SupplyVaultTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (supplyVault *SupplyVault) UnpackTransferEvent(log *types.Log) (*SupplyVaultTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultTransfer)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultWithdraw represents a Withdraw event raised by the SupplyVault contract.
type SupplyVaultWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const SupplyVaultWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (SupplyVaultWithdraw) ContractEventName() string {
	return SupplyVaultWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (supplyVault *SupplyVault) UnpackWithdrawEvent(log *types.Log) (*SupplyVaultWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultWithdraw)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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

// SupplyVaultWriteOff represents a WriteOff event raised by the SupplyVault contract.
type SupplyVaultWriteOff struct {
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const SupplyVaultWriteOffEventName = "WriteOff"

// ContractEventName returns the user-defined event name.
func (SupplyVaultWriteOff) ContractEventName() string {
	return SupplyVaultWriteOffEventName
}

// UnpackWriteOffEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WriteOff(uint256 assets)
func (supplyVault *SupplyVault) UnpackWriteOffEvent(log *types.Log) (*SupplyVaultWriteOff, error) {
	event := "WriteOff"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != supplyVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SupplyVaultWriteOff)
	if len(log.Data) > 0 {
		if err := supplyVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range supplyVault.abi.Events[event].Inputs {
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
func (supplyVault *SupplyVault) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["BorrowerAlreadySet"].ID.Bytes()[:4]) {
		return supplyVault.UnpackBorrowerAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC4626ExceededMaxDeposit"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC4626ExceededMaxDepositError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC4626ExceededMaxMint"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC4626ExceededMaxMintError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC4626ExceededMaxRedeem"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC4626ExceededMaxRedeemError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["ERC4626ExceededMaxWithdraw"].ID.Bytes()[:4]) {
		return supplyVault.UnpackERC4626ExceededMaxWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["InsufficientLiquidity"].ID.Bytes()[:4]) {
		return supplyVault.UnpackInsufficientLiquidityError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["InvalidBorrower"].ID.Bytes()[:4]) {
		return supplyVault.UnpackInvalidBorrowerError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["InvalidRateModel"].ID.Bytes()[:4]) {
		return supplyVault.UnpackInvalidRateModelError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["NotBorrower"].ID.Bytes()[:4]) {
		return supplyVault.UnpackNotBorrowerError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return supplyVault.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return supplyVault.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return supplyVault.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return supplyVault.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], supplyVault.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return supplyVault.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// SupplyVaultBorrowerAlreadySet represents a BorrowerAlreadySet error raised by the SupplyVault contract.
type SupplyVaultBorrowerAlreadySet struct {
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BorrowerAlreadySet(address current)
func SupplyVaultBorrowerAlreadySetErrorID() common.Hash {
	return common.HexToHash("0x94cc6745052ebc1c849563f4ebf4a536422ff0e3f99bac1a79a177e6f3cc4618")
}

// UnpackBorrowerAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BorrowerAlreadySet(address current)
func (supplyVault *SupplyVault) UnpackBorrowerAlreadySetError(raw []byte) (*SupplyVaultBorrowerAlreadySet, error) {
	out := new(SupplyVaultBorrowerAlreadySet)
	if err := supplyVault.abi.UnpackIntoInterface(out, "BorrowerAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the SupplyVault contract.
type SupplyVaultERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func SupplyVaultERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (supplyVault *SupplyVault) UnpackERC20InsufficientAllowanceError(raw []byte) (*SupplyVaultERC20InsufficientAllowance, error) {
	out := new(SupplyVaultERC20InsufficientAllowance)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the SupplyVault contract.
type SupplyVaultERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func SupplyVaultERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (supplyVault *SupplyVault) UnpackERC20InsufficientBalanceError(raw []byte) (*SupplyVaultERC20InsufficientBalance, error) {
	out := new(SupplyVaultERC20InsufficientBalance)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InvalidApprover represents a ERC20InvalidApprover error raised by the SupplyVault contract.
type SupplyVaultERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func SupplyVaultERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (supplyVault *SupplyVault) UnpackERC20InvalidApproverError(raw []byte) (*SupplyVaultERC20InvalidApprover, error) {
	out := new(SupplyVaultERC20InvalidApprover)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the SupplyVault contract.
type SupplyVaultERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func SupplyVaultERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (supplyVault *SupplyVault) UnpackERC20InvalidReceiverError(raw []byte) (*SupplyVaultERC20InvalidReceiver, error) {
	out := new(SupplyVaultERC20InvalidReceiver)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InvalidSender represents a ERC20InvalidSender error raised by the SupplyVault contract.
type SupplyVaultERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func SupplyVaultERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (supplyVault *SupplyVault) UnpackERC20InvalidSenderError(raw []byte) (*SupplyVaultERC20InvalidSender, error) {
	out := new(SupplyVaultERC20InvalidSender)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC20InvalidSpender represents a ERC20InvalidSpender error raised by the SupplyVault contract.
type SupplyVaultERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func SupplyVaultERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (supplyVault *SupplyVault) UnpackERC20InvalidSpenderError(raw []byte) (*SupplyVaultERC20InvalidSpender, error) {
	out := new(SupplyVaultERC20InvalidSpender)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC4626ExceededMaxDeposit represents a ERC4626ExceededMaxDeposit error raised by the SupplyVault contract.
type SupplyVaultERC4626ExceededMaxDeposit struct {
	Receiver common.Address
	Assets   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func SupplyVaultERC4626ExceededMaxDepositErrorID() common.Hash {
	return common.HexToHash("0x79012fb2819fcdc1de669c08773dbcd6bdc757862642f75fb1c584dadf259dfe")
}

// UnpackERC4626ExceededMaxDepositError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func (supplyVault *SupplyVault) UnpackERC4626ExceededMaxDepositError(raw []byte) (*SupplyVaultERC4626ExceededMaxDeposit, error) {
	out := new(SupplyVaultERC4626ExceededMaxDeposit)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxDeposit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC4626ExceededMaxMint represents a ERC4626ExceededMaxMint error raised by the SupplyVault contract.
type SupplyVaultERC4626ExceededMaxMint struct {
	Receiver common.Address
	Shares   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func SupplyVaultERC4626ExceededMaxMintErrorID() common.Hash {
	return common.HexToHash("0x284ff667dc615a39438518c22e8955b9470327d9de8a4d7e21c926b260d65176")
}

// UnpackERC4626ExceededMaxMintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func (supplyVault *SupplyVault) UnpackERC4626ExceededMaxMintError(raw []byte) (*SupplyVaultERC4626ExceededMaxMint, error) {
	out := new(SupplyVaultERC4626ExceededMaxMint)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxMint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC4626ExceededMaxRedeem represents a ERC4626ExceededMaxRedeem error raised by the SupplyVault contract.
type SupplyVaultERC4626ExceededMaxRedeem struct {
	Owner  common.Address
	Shares *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func SupplyVaultERC4626ExceededMaxRedeemErrorID() common.Hash {
	return common.HexToHash("0xb94abeec0557d36b5f0bc8f115deec7b184dcbff94ac66d55e37c8f301e75269")
}

// UnpackERC4626ExceededMaxRedeemError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func (supplyVault *SupplyVault) UnpackERC4626ExceededMaxRedeemError(raw []byte) (*SupplyVaultERC4626ExceededMaxRedeem, error) {
	out := new(SupplyVaultERC4626ExceededMaxRedeem)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxRedeem", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultERC4626ExceededMaxWithdraw represents a ERC4626ExceededMaxWithdraw error raised by the SupplyVault contract.
type SupplyVaultERC4626ExceededMaxWithdraw struct {
	Owner  common.Address
	Assets *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func SupplyVaultERC4626ExceededMaxWithdrawErrorID() common.Hash {
	return common.HexToHash("0xfe9cceec2bd1f9b68641914cc354eacaeb1cc2169f5ba0639930f241e87142f0")
}

// UnpackERC4626ExceededMaxWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func (supplyVault *SupplyVault) UnpackERC4626ExceededMaxWithdrawError(raw []byte) (*SupplyVaultERC4626ExceededMaxWithdraw, error) {
	out := new(SupplyVaultERC4626ExceededMaxWithdraw)
	if err := supplyVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultInsufficientLiquidity represents a InsufficientLiquidity error raised by the SupplyVault contract.
type SupplyVaultInsufficientLiquidity struct {
	Requested *big.Int
	Available *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientLiquidity(uint256 requested, uint256 available)
func SupplyVaultInsufficientLiquidityErrorID() common.Hash {
	return common.HexToHash("0xa17e11d5aed9d2c4240142d28db52e7da7283f43ac33766dd68e8f3dc6613c0d")
}

// UnpackInsufficientLiquidityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientLiquidity(uint256 requested, uint256 available)
func (supplyVault *SupplyVault) UnpackInsufficientLiquidityError(raw []byte) (*SupplyVaultInsufficientLiquidity, error) {
	out := new(SupplyVaultInsufficientLiquidity)
	if err := supplyVault.abi.UnpackIntoInterface(out, "InsufficientLiquidity", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultInvalidBorrower represents a InvalidBorrower error raised by the SupplyVault contract.
type SupplyVaultInvalidBorrower struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidBorrower()
func SupplyVaultInvalidBorrowerErrorID() common.Hash {
	return common.HexToHash("0x6f5f81d74c5f06945f90dbf5f9bf8227c338474b5bab4855ffad6fc0f5217f34")
}

// UnpackInvalidBorrowerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidBorrower()
func (supplyVault *SupplyVault) UnpackInvalidBorrowerError(raw []byte) (*SupplyVaultInvalidBorrower, error) {
	out := new(SupplyVaultInvalidBorrower)
	if err := supplyVault.abi.UnpackIntoInterface(out, "InvalidBorrower", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultInvalidRateModel represents a InvalidRateModel error raised by the SupplyVault contract.
type SupplyVaultInvalidRateModel struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRateModel()
func SupplyVaultInvalidRateModelErrorID() common.Hash {
	return common.HexToHash("0xd2375d9d660f19094039027d0f6695eff69fd031e8b29f0e5dc1923e6f5a0d25")
}

// UnpackInvalidRateModelError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRateModel()
func (supplyVault *SupplyVault) UnpackInvalidRateModelError(raw []byte) (*SupplyVaultInvalidRateModel, error) {
	out := new(SupplyVaultInvalidRateModel)
	if err := supplyVault.abi.UnpackIntoInterface(out, "InvalidRateModel", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultNotBorrower represents a NotBorrower error raised by the SupplyVault contract.
type SupplyVaultNotBorrower struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotBorrower(address account)
func SupplyVaultNotBorrowerErrorID() common.Hash {
	return common.HexToHash("0x3d19c0eaa213120f897b12cdf0df62bf048ed2228f0ad44d3e89ce9724036d29")
}

// UnpackNotBorrowerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotBorrower(address account)
func (supplyVault *SupplyVault) UnpackNotBorrowerError(raw []byte) (*SupplyVaultNotBorrower, error) {
	out := new(SupplyVaultNotBorrower)
	if err := supplyVault.abi.UnpackIntoInterface(out, "NotBorrower", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the SupplyVault contract.
type SupplyVaultOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func SupplyVaultOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (supplyVault *SupplyVault) UnpackOwnableInvalidOwnerError(raw []byte) (*SupplyVaultOwnableInvalidOwner, error) {
	out := new(SupplyVaultOwnableInvalidOwner)
	if err := supplyVault.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the SupplyVault contract.
type SupplyVaultOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func SupplyVaultOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (supplyVault *SupplyVault) UnpackOwnableUnauthorizedAccountError(raw []byte) (*SupplyVaultOwnableUnauthorizedAccount, error) {
	out := new(SupplyVaultOwnableUnauthorizedAccount)
	if err := supplyVault.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the SupplyVault contract.
type SupplyVaultOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func SupplyVaultOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (supplyVault *SupplyVault) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*SupplyVaultOwnershipCannotBeRenounced, error) {
	out := new(SupplyVaultOwnershipCannotBeRenounced)
	if err := supplyVault.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the SupplyVault contract.
type SupplyVaultSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func SupplyVaultSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (supplyVault *SupplyVault) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*SupplyVaultSafeCastOverflowedUintDowncast, error) {
	out := new(SupplyVaultSafeCastOverflowedUintDowncast)
	if err := supplyVault.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SupplyVaultSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the SupplyVault contract.
type SupplyVaultSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func SupplyVaultSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (supplyVault *SupplyVault) UnpackSafeERC20FailedOperationError(raw []byte) (*SupplyVaultSafeERC20FailedOperation, error) {
	out := new(SupplyVaultSafeERC20FailedOperation)
	if err := supplyVault.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}
