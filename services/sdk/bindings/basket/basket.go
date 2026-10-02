// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package basket

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

// BasketMetaData contains all meta data concerning the Basket contract.
var BasketMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"band_\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"symbols\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"units\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_COMPONENTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NOTICE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"components\",\"inputs\":[],\"outputs\":[{\"name\":\"symbols\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isBlocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxAssets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingTarget\",\"inputs\":[],\"outputs\":[{\"name\":\"units\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"effectiveAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeTarget\",\"inputs\":[{\"name\":\"units\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rebalance\",\"inputs\":[{\"name\":\"assetsIn\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"assetsOut\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"target\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Rebalanced\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assetsIn\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"assetsOut\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TargetProposed\",\"inputs\":[{\"name\":\"units\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"effectiveAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AboveMaximum\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AssetHalted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Blocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ComponentPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CorporateActionPending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DuplicateComponent\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EmptyBasket\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidComponents\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTarget\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PastTarget\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAsset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ValueLost\",\"inputs\":[{\"name\":\"valueIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valueOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
	ID:  "Basket",
}

// Basket is an auto generated Go binding around an Ethereum contract.
type Basket struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Basket) GetABI() abi.ABI {
	return c.abi
}

// NewBasket creates a new instance of Basket.
func NewBasket() *Basket {
	parsed, err := BasketMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Basket{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Basket) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(string name_, string symbol_, address band_, bytes32[] symbols, uint256[] units, address initialOwner) returns()
func (basket *Basket) PackConstructor(name_ string, symbol_ string, band_ common.Address, symbols [][32]byte, units []*big.Int, initialOwner common.Address) []byte {
	enc, err := basket.abi.Pack("", name_, symbol_, band_, symbols, units, initialOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMAXCOMPONENTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0ab50a6a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_COMPONENTS() view returns(uint256)
func (basket *Basket) PackMAXCOMPONENTS() []byte {
	enc, err := basket.abi.Pack("MAX_COMPONENTS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXCOMPONENTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0ab50a6a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_COMPONENTS() view returns(uint256)
func (basket *Basket) TryPackMAXCOMPONENTS() ([]byte, error) {
	return basket.abi.Pack("MAX_COMPONENTS")
}

// UnpackMAXCOMPONENTS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0ab50a6a.
//
// Solidity: function MAX_COMPONENTS() view returns(uint256)
func (basket *Basket) UnpackMAXCOMPONENTS(data []byte) (*big.Int, error) {
	out, err := basket.abi.Unpack("MAX_COMPONENTS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackNOTICE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4bfee686.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function NOTICE() view returns(uint64)
func (basket *Basket) PackNOTICE() []byte {
	enc, err := basket.abi.Pack("NOTICE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNOTICE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4bfee686.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function NOTICE() view returns(uint64)
func (basket *Basket) TryPackNOTICE() ([]byte, error) {
	return basket.abi.Pack("NOTICE")
}

// UnpackNOTICE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4bfee686.
//
// Solidity: function NOTICE() view returns(uint64)
func (basket *Basket) UnpackNOTICE(data []byte) (uint64, error) {
	out, err := basket.abi.Unpack("NOTICE", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (basket *Basket) PackAcceptOwnership() []byte {
	enc, err := basket.abi.Pack("acceptOwnership")
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
func (basket *Basket) TryPackAcceptOwnership() ([]byte, error) {
	return basket.abi.Pack("acceptOwnership")
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (basket *Basket) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := basket.abi.Pack("allowance", owner, spender)
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
func (basket *Basket) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return basket.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (basket *Basket) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := basket.abi.Unpack("allowance", data)
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
func (basket *Basket) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := basket.abi.Pack("approve", spender, value)
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
func (basket *Basket) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return basket.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (basket *Basket) UnpackApprove(data []byte) (bool, error) {
	out, err := basket.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (basket *Basket) PackBalanceOf(account common.Address) []byte {
	enc, err := basket.abi.Pack("balanceOf", account)
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
func (basket *Basket) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return basket.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (basket *Basket) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := basket.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (basket *Basket) PackBand() []byte {
	enc, err := basket.abi.Pack("band")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function band() view returns(address)
func (basket *Basket) TryPackBand() ([]byte, error) {
	return basket.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (basket *Basket) UnpackBand(data []byte) (common.Address, error) {
	out, err := basket.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackComponents is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba62fbe4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function components() view returns(bytes32[] symbols, address[] tokens)
func (basket *Basket) PackComponents() []byte {
	enc, err := basket.abi.Pack("components")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackComponents is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba62fbe4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function components() view returns(bytes32[] symbols, address[] tokens)
func (basket *Basket) TryPackComponents() ([]byte, error) {
	return basket.abi.Pack("components")
}

// ComponentsOutput serves as a container for the return parameters of contract
// method Components.
type ComponentsOutput struct {
	Symbols [][32]byte
	Tokens  []common.Address
}

// UnpackComponents is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba62fbe4.
//
// Solidity: function components() view returns(bytes32[] symbols, address[] tokens)
func (basket *Basket) UnpackComponents(data []byte) (ComponentsOutput, error) {
	out, err := basket.abi.Unpack("components", data)
	outstruct := new(ComponentsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Symbols = *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	outstruct.Tokens = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	return *outstruct, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (basket *Basket) PackDecimals() []byte {
	enc, err := basket.abi.Pack("decimals")
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
func (basket *Basket) TryPackDecimals() ([]byte, error) {
	return basket.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (basket *Basket) UnpackDecimals(data []byte) (uint8, error) {
	out, err := basket.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackIsBlocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbac3951.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isBlocked(address account) view returns(bool)
func (basket *Basket) PackIsBlocked(account common.Address) []byte {
	enc, err := basket.abi.Pack("isBlocked", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsBlocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbac3951.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isBlocked(address account) view returns(bool)
func (basket *Basket) TryPackIsBlocked(account common.Address) ([]byte, error) {
	return basket.abi.Pack("isBlocked", account)
}

// UnpackIsBlocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfbac3951.
//
// Solidity: function isBlocked(address account) view returns(bool)
func (basket *Basket) UnpackIsBlocked(data []byte) (bool, error) {
	out, err := basket.abi.Unpack("isBlocked", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa15da4a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(uint256 shares, address receiver, uint256[] maxAssets) returns(uint256[] assets)
func (basket *Basket) PackMint(shares *big.Int, receiver common.Address, maxAssets []*big.Int) []byte {
	enc, err := basket.abi.Pack("mint", shares, receiver, maxAssets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa15da4a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(uint256 shares, address receiver, uint256[] maxAssets) returns(uint256[] assets)
func (basket *Basket) TryPackMint(shares *big.Int, receiver common.Address, maxAssets []*big.Int) ([]byte, error) {
	return basket.abi.Pack("mint", shares, receiver, maxAssets)
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa15da4a.
//
// Solidity: function mint(uint256 shares, address receiver, uint256[] maxAssets) returns(uint256[] assets)
func (basket *Basket) UnpackMint(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("mint", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (basket *Basket) PackName() []byte {
	enc, err := basket.abi.Pack("name")
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
func (basket *Basket) TryPackName() ([]byte, error) {
	return basket.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (basket *Basket) UnpackName(data []byte) (string, error) {
	out, err := basket.abi.Unpack("name", data)
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
func (basket *Basket) PackOwner() []byte {
	enc, err := basket.abi.Pack("owner")
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
func (basket *Basket) TryPackOwner() ([]byte, error) {
	return basket.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (basket *Basket) UnpackOwner(data []byte) (common.Address, error) {
	out, err := basket.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paused() view returns(bool)
func (basket *Basket) PackPaused() []byte {
	enc, err := basket.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function paused() view returns(bool)
func (basket *Basket) TryPackPaused() ([]byte, error) {
	return basket.abi.Pack("paused")
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (basket *Basket) UnpackPaused(data []byte) (bool, error) {
	out, err := basket.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (basket *Basket) PackPendingOwner() []byte {
	enc, err := basket.abi.Pack("pendingOwner")
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
func (basket *Basket) TryPackPendingOwner() ([]byte, error) {
	return basket.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (basket *Basket) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := basket.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58ea52a5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingTarget() view returns(uint256[] units, uint64 effectiveAt)
func (basket *Basket) PackPendingTarget() []byte {
	enc, err := basket.abi.Pack("pendingTarget")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58ea52a5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingTarget() view returns(uint256[] units, uint64 effectiveAt)
func (basket *Basket) TryPackPendingTarget() ([]byte, error) {
	return basket.abi.Pack("pendingTarget")
}

// PendingTargetOutput serves as a container for the return parameters of contract
// method PendingTarget.
type PendingTargetOutput struct {
	Units       []*big.Int
	EffectiveAt uint64
}

// UnpackPendingTarget is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x58ea52a5.
//
// Solidity: function pendingTarget() view returns(uint256[] units, uint64 effectiveAt)
func (basket *Basket) UnpackPendingTarget(data []byte) (PendingTargetOutput, error) {
	out, err := basket.abi.Unpack("pendingTarget", data)
	outstruct := new(PendingTargetOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Units = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.EffectiveAt = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) PackPreviewMint(shares *big.Int) []byte {
	enc, err := basket.abi.Pack("previewMint", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) TryPackPreviewMint(shares *big.Int) ([]byte, error) {
	return basket.abi.Pack("previewMint", shares)
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) UnpackPreviewMint(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("previewMint", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := basket.abi.Pack("previewRedeem", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) TryPackPreviewRedeem(shares *big.Int) ([]byte, error) {
	return basket.abi.Pack("previewRedeem", shares)
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256[] assets)
func (basket *Basket) UnpackPreviewRedeem(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("previewRedeem", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackProposeTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x512dead6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function proposeTarget(uint256[] units) returns()
func (basket *Basket) PackProposeTarget(units []*big.Int) []byte {
	enc, err := basket.abi.Pack("proposeTarget", units)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProposeTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x512dead6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function proposeTarget(uint256[] units) returns()
func (basket *Basket) TryPackProposeTarget(units []*big.Int) ([]byte, error) {
	return basket.abi.Pack("proposeTarget", units)
}

// PackRebalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb08302b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rebalance(uint256[] assetsIn, uint256[] assetsOut, address receiver) returns()
func (basket *Basket) PackRebalance(assetsIn []*big.Int, assetsOut []*big.Int, receiver common.Address) []byte {
	enc, err := basket.abi.Pack("rebalance", assetsIn, assetsOut, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRebalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdb08302b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rebalance(uint256[] assetsIn, uint256[] assetsOut, address receiver) returns()
func (basket *Basket) TryPackRebalance(assetsIn []*big.Int, assetsOut []*big.Int, receiver common.Address) ([]byte, error) {
	return basket.abi.Pack("rebalance", assetsIn, assetsOut, receiver)
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256[] assets)
func (basket *Basket) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := basket.abi.Pack("redeem", shares, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256[] assets)
func (basket *Basket) TryPackRedeem(shares *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return basket.abi.Pack("redeem", shares, receiver, owner)
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256[] assets)
func (basket *Basket) UnpackRedeem(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("redeem", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() pure returns()
func (basket *Basket) PackRenounceOwnership() []byte {
	enc, err := basket.abi.Pack("renounceOwnership")
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
func (basket *Basket) TryPackRenounceOwnership() ([]byte, error) {
	return basket.abi.Pack("renounceOwnership")
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (basket *Basket) PackSymbol() []byte {
	enc, err := basket.abi.Pack("symbol")
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
func (basket *Basket) TryPackSymbol() ([]byte, error) {
	return basket.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (basket *Basket) UnpackSymbol(data []byte) (string, error) {
	out, err := basket.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4b83992.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function target() view returns(uint256[])
func (basket *Basket) PackTarget() []byte {
	enc, err := basket.abi.Pack("target")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTarget is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4b83992.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function target() view returns(uint256[])
func (basket *Basket) TryPackTarget() ([]byte, error) {
	return basket.abi.Pack("target")
}

// UnpackTarget is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd4b83992.
//
// Solidity: function target() view returns(uint256[])
func (basket *Basket) UnpackTarget(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("target", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAssets() view returns(uint256[] assets)
func (basket *Basket) PackTotalAssets() []byte {
	enc, err := basket.abi.Pack("totalAssets")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalAssets() view returns(uint256[] assets)
func (basket *Basket) TryPackTotalAssets() ([]byte, error) {
	return basket.abi.Pack("totalAssets")
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256[] assets)
func (basket *Basket) UnpackTotalAssets(data []byte) ([]*big.Int, error) {
	out, err := basket.abi.Unpack("totalAssets", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (basket *Basket) PackTotalSupply() []byte {
	enc, err := basket.abi.Pack("totalSupply")
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
func (basket *Basket) TryPackTotalSupply() ([]byte, error) {
	return basket.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (basket *Basket) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := basket.abi.Unpack("totalSupply", data)
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
func (basket *Basket) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := basket.abi.Pack("transfer", to, value)
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
func (basket *Basket) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return basket.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (basket *Basket) UnpackTransfer(data []byte) (bool, error) {
	out, err := basket.abi.Unpack("transfer", data)
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
func (basket *Basket) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := basket.abi.Pack("transferFrom", from, to, value)
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
func (basket *Basket) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return basket.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (basket *Basket) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := basket.abi.Unpack("transferFrom", data)
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
func (basket *Basket) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := basket.abi.Pack("transferOwnership", newOwner)
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
func (basket *Basket) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return basket.abi.Pack("transferOwnership", newOwner)
}

// BasketApproval represents a Approval event raised by the Basket contract.
type BasketApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const BasketApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (BasketApproval) ContractEventName() string {
	return BasketApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (basket *Basket) UnpackApprovalEvent(log *types.Log) (*BasketApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketApproval)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketDeposit represents a Deposit event raised by the Basket contract.
type BasketDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets []*big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const BasketDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (BasketDeposit) ContractEventName() string {
	return BasketDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256[] assets, uint256 shares)
func (basket *Basket) UnpackDepositEvent(log *types.Log) (*BasketDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketDeposit)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Basket contract.
type BasketOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const BasketOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (BasketOwnershipTransferStarted) ContractEventName() string {
	return BasketOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (basket *Basket) UnpackOwnershipTransferStartedEvent(log *types.Log) (*BasketOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketOwnershipTransferred represents a OwnershipTransferred event raised by the Basket contract.
type BasketOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const BasketOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (BasketOwnershipTransferred) ContractEventName() string {
	return BasketOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (basket *Basket) UnpackOwnershipTransferredEvent(log *types.Log) (*BasketOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketRebalanced represents a Rebalanced event raised by the Basket contract.
type BasketRebalanced struct {
	Caller    common.Address
	Receiver  common.Address
	AssetsIn  []*big.Int
	AssetsOut []*big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const BasketRebalancedEventName = "Rebalanced"

// ContractEventName returns the user-defined event name.
func (BasketRebalanced) ContractEventName() string {
	return BasketRebalancedEventName
}

// UnpackRebalancedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Rebalanced(address indexed caller, address indexed receiver, uint256[] assetsIn, uint256[] assetsOut)
func (basket *Basket) UnpackRebalancedEvent(log *types.Log) (*BasketRebalanced, error) {
	event := "Rebalanced"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketRebalanced)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketTargetProposed represents a TargetProposed event raised by the Basket contract.
type BasketTargetProposed struct {
	Units       []*big.Int
	EffectiveAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const BasketTargetProposedEventName = "TargetProposed"

// ContractEventName returns the user-defined event name.
func (BasketTargetProposed) ContractEventName() string {
	return BasketTargetProposedEventName
}

// UnpackTargetProposedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TargetProposed(uint256[] units, uint64 effectiveAt)
func (basket *Basket) UnpackTargetProposedEvent(log *types.Log) (*BasketTargetProposed, error) {
	event := "TargetProposed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketTargetProposed)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketTransfer represents a Transfer event raised by the Basket contract.
type BasketTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const BasketTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (BasketTransfer) ContractEventName() string {
	return BasketTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (basket *Basket) UnpackTransferEvent(log *types.Log) (*BasketTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketTransfer)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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

// BasketWithdraw represents a Withdraw event raised by the Basket contract.
type BasketWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   []*big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const BasketWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (BasketWithdraw) ContractEventName() string {
	return BasketWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256[] assets, uint256 shares)
func (basket *Basket) UnpackWithdrawEvent(log *types.Log) (*BasketWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != basket.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BasketWithdraw)
	if len(log.Data) > 0 {
		if err := basket.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range basket.abi.Events[event].Inputs {
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
func (basket *Basket) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], basket.abi.Errors["AboveMaximum"].ID.Bytes()[:4]) {
		return basket.UnpackAboveMaximumError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["AssetHalted"].ID.Bytes()[:4]) {
		return basket.UnpackAssetHaltedError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["Blocked"].ID.Bytes()[:4]) {
		return basket.UnpackBlockedError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ComponentPaused"].ID.Bytes()[:4]) {
		return basket.UnpackComponentPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["CorporateActionPending"].ID.Bytes()[:4]) {
		return basket.UnpackCorporateActionPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["DuplicateComponent"].ID.Bytes()[:4]) {
		return basket.UnpackDuplicateComponentError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return basket.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["EmptyBasket"].ID.Bytes()[:4]) {
		return basket.UnpackEmptyBasketError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["InvalidComponents"].ID.Bytes()[:4]) {
		return basket.UnpackInvalidComponentsError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["InvalidTarget"].ID.Bytes()[:4]) {
		return basket.UnpackInvalidTargetError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["LengthMismatch"].ID.Bytes()[:4]) {
		return basket.UnpackLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return basket.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return basket.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return basket.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["PastTarget"].ID.Bytes()[:4]) {
		return basket.UnpackPastTargetError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return basket.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return basket.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return basket.UnpackSequencerNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["UnknownAsset"].ID.Bytes()[:4]) {
		return basket.UnpackUnknownAssetError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ValueLost"].ID.Bytes()[:4]) {
		return basket.UnpackValueLostError(raw[4:])
	}
	if bytes.Equal(raw[:4], basket.abi.Errors["ZeroAmount"].ID.Bytes()[:4]) {
		return basket.UnpackZeroAmountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// BasketAboveMaximum represents a AboveMaximum error raised by the Basket contract.
type BasketAboveMaximum struct {
	Token     common.Address
	Amount    *big.Int
	MaxAmount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AboveMaximum(address token, uint256 amount, uint256 maxAmount)
func BasketAboveMaximumErrorID() common.Hash {
	return common.HexToHash("0x99d9388df1ca391cfb672a323a32e86ca75ac7e67a8bf44ebe7d3f7c71b82399")
}

// UnpackAboveMaximumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AboveMaximum(address token, uint256 amount, uint256 maxAmount)
func (basket *Basket) UnpackAboveMaximumError(raw []byte) (*BasketAboveMaximum, error) {
	out := new(BasketAboveMaximum)
	if err := basket.abi.UnpackIntoInterface(out, "AboveMaximum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketAssetHalted represents a AssetHalted error raised by the Basket contract.
type BasketAssetHalted struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetHalted(bytes32 symbol)
func BasketAssetHaltedErrorID() common.Hash {
	return common.HexToHash("0x3ec29ca3785f47585a9c96baed33ddbd8b2d9e8dfdeffb4be96925aaf9ec6468")
}

// UnpackAssetHaltedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetHalted(bytes32 symbol)
func (basket *Basket) UnpackAssetHaltedError(raw []byte) (*BasketAssetHalted, error) {
	out := new(BasketAssetHalted)
	if err := basket.abi.UnpackIntoInterface(out, "AssetHalted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketBlocked represents a Blocked error raised by the Basket contract.
type BasketBlocked struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Blocked(address account)
func BasketBlockedErrorID() common.Hash {
	return common.HexToHash("0x75e91ce73c1d3352d8dd3610443539cd33dfe13b1de8f8caae54ec26dd0dc9cb")
}

// UnpackBlockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Blocked(address account)
func (basket *Basket) UnpackBlockedError(raw []byte) (*BasketBlocked, error) {
	out := new(BasketBlocked)
	if err := basket.abi.UnpackIntoInterface(out, "Blocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketComponentPaused represents a ComponentPaused error raised by the Basket contract.
type BasketComponentPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ComponentPaused()
func BasketComponentPausedErrorID() common.Hash {
	return common.HexToHash("0xe706ab8b72b289ec71c62608e0de2085e321001abfc30eb217c2b6d07356b16d")
}

// UnpackComponentPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ComponentPaused()
func (basket *Basket) UnpackComponentPausedError(raw []byte) (*BasketComponentPaused, error) {
	out := new(BasketComponentPaused)
	if err := basket.abi.UnpackIntoInterface(out, "ComponentPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketCorporateActionPending represents a CorporateActionPending error raised by the Basket contract.
type BasketCorporateActionPending struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func BasketCorporateActionPendingErrorID() common.Hash {
	return common.HexToHash("0x1d1f5fd3c689a2e2561007ef605d914d864ac52012393f43a887a76f9f48be7b")
}

// UnpackCorporateActionPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func (basket *Basket) UnpackCorporateActionPendingError(raw []byte) (*BasketCorporateActionPending, error) {
	out := new(BasketCorporateActionPending)
	if err := basket.abi.UnpackIntoInterface(out, "CorporateActionPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketDuplicateComponent represents a DuplicateComponent error raised by the Basket contract.
type BasketDuplicateComponent struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DuplicateComponent(bytes32 symbol)
func BasketDuplicateComponentErrorID() common.Hash {
	return common.HexToHash("0x11c290dee5f8204f515831cea5d94af8523f9aa9459ddae96e0f72f586069f90")
}

// UnpackDuplicateComponentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DuplicateComponent(bytes32 symbol)
func (basket *Basket) UnpackDuplicateComponentError(raw []byte) (*BasketDuplicateComponent, error) {
	out := new(BasketDuplicateComponent)
	if err := basket.abi.UnpackIntoInterface(out, "DuplicateComponent", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the Basket contract.
type BasketERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func BasketERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (basket *Basket) UnpackERC20InsufficientAllowanceError(raw []byte) (*BasketERC20InsufficientAllowance, error) {
	out := new(BasketERC20InsufficientAllowance)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the Basket contract.
type BasketERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func BasketERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (basket *Basket) UnpackERC20InsufficientBalanceError(raw []byte) (*BasketERC20InsufficientBalance, error) {
	out := new(BasketERC20InsufficientBalance)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InvalidApprover represents a ERC20InvalidApprover error raised by the Basket contract.
type BasketERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func BasketERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (basket *Basket) UnpackERC20InvalidApproverError(raw []byte) (*BasketERC20InvalidApprover, error) {
	out := new(BasketERC20InvalidApprover)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the Basket contract.
type BasketERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func BasketERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (basket *Basket) UnpackERC20InvalidReceiverError(raw []byte) (*BasketERC20InvalidReceiver, error) {
	out := new(BasketERC20InvalidReceiver)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InvalidSender represents a ERC20InvalidSender error raised by the Basket contract.
type BasketERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func BasketERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (basket *Basket) UnpackERC20InvalidSenderError(raw []byte) (*BasketERC20InvalidSender, error) {
	out := new(BasketERC20InvalidSender)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketERC20InvalidSpender represents a ERC20InvalidSpender error raised by the Basket contract.
type BasketERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func BasketERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (basket *Basket) UnpackERC20InvalidSpenderError(raw []byte) (*BasketERC20InvalidSpender, error) {
	out := new(BasketERC20InvalidSpender)
	if err := basket.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketEmptyBasket represents a EmptyBasket error raised by the Basket contract.
type BasketEmptyBasket struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBasket()
func BasketEmptyBasketErrorID() common.Hash {
	return common.HexToHash("0x38f6c555b1959441054760a767163efc008492283ea4df3c6bfafd5a1ee24af8")
}

// UnpackEmptyBasketError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBasket()
func (basket *Basket) UnpackEmptyBasketError(raw []byte) (*BasketEmptyBasket, error) {
	out := new(BasketEmptyBasket)
	if err := basket.abi.UnpackIntoInterface(out, "EmptyBasket", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketInvalidComponents represents a InvalidComponents error raised by the Basket contract.
type BasketInvalidComponents struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidComponents()
func BasketInvalidComponentsErrorID() common.Hash {
	return common.HexToHash("0xf04ced8b0eb227c029948fabfd6fd7b3115532000191783f7140efdd58d66d0a")
}

// UnpackInvalidComponentsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidComponents()
func (basket *Basket) UnpackInvalidComponentsError(raw []byte) (*BasketInvalidComponents, error) {
	out := new(BasketInvalidComponents)
	if err := basket.abi.UnpackIntoInterface(out, "InvalidComponents", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketInvalidTarget represents a InvalidTarget error raised by the Basket contract.
type BasketInvalidTarget struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidTarget()
func BasketInvalidTargetErrorID() common.Hash {
	return common.HexToHash("0x82d5d76a554b1226e0bd53ae5121b9d30262dd775be75ba7d4bd87ba04ede349")
}

// UnpackInvalidTargetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidTarget()
func (basket *Basket) UnpackInvalidTargetError(raw []byte) (*BasketInvalidTarget, error) {
	out := new(BasketInvalidTarget)
	if err := basket.abi.UnpackIntoInterface(out, "InvalidTarget", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketLengthMismatch represents a LengthMismatch error raised by the Basket contract.
type BasketLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LengthMismatch()
func BasketLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xff633a3803c58b9bc21e58efecee59f27e033cc0b1883fccb4969c76146fe60f")
}

// UnpackLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LengthMismatch()
func (basket *Basket) UnpackLengthMismatchError(raw []byte) (*BasketLengthMismatch, error) {
	out := new(BasketLengthMismatch)
	if err := basket.abi.UnpackIntoInterface(out, "LengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the Basket contract.
type BasketOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func BasketOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (basket *Basket) UnpackOwnableInvalidOwnerError(raw []byte) (*BasketOwnableInvalidOwner, error) {
	out := new(BasketOwnableInvalidOwner)
	if err := basket.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the Basket contract.
type BasketOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func BasketOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (basket *Basket) UnpackOwnableUnauthorizedAccountError(raw []byte) (*BasketOwnableUnauthorizedAccount, error) {
	out := new(BasketOwnableUnauthorizedAccount)
	if err := basket.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the Basket contract.
type BasketOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func BasketOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (basket *Basket) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*BasketOwnershipCannotBeRenounced, error) {
	out := new(BasketOwnershipCannotBeRenounced)
	if err := basket.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketPastTarget represents a PastTarget error raised by the Basket contract.
type BasketPastTarget struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PastTarget(bytes32 symbol)
func BasketPastTargetErrorID() common.Hash {
	return common.HexToHash("0xfde097685aaefa97ba2f9145323d3df9ca91a60a8ee4ae5382bfe81dcf7d5e23")
}

// UnpackPastTargetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PastTarget(bytes32 symbol)
func (basket *Basket) UnpackPastTargetError(raw []byte) (*BasketPastTarget, error) {
	out := new(BasketPastTarget)
	if err := basket.abi.UnpackIntoInterface(out, "PastTarget", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the Basket contract.
type BasketSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func BasketSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (basket *Basket) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*BasketSafeCastOverflowedUintDowncast, error) {
	out := new(BasketSafeCastOverflowedUintDowncast)
	if err := basket.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the Basket contract.
type BasketSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func BasketSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (basket *Basket) UnpackSafeERC20FailedOperationError(raw []byte) (*BasketSafeERC20FailedOperation, error) {
	out := new(BasketSafeERC20FailedOperation)
	if err := basket.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketSequencerNotSettled represents a SequencerNotSettled error raised by the Basket contract.
type BasketSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func BasketSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (basket *Basket) UnpackSequencerNotSettledError(raw []byte) (*BasketSequencerNotSettled, error) {
	out := new(BasketSequencerNotSettled)
	if err := basket.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketUnknownAsset represents a UnknownAsset error raised by the Basket contract.
type BasketUnknownAsset struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func BasketUnknownAssetErrorID() common.Hash {
	return common.HexToHash("0x1059be3ed9f1c1f061f09edf206f56e7c80b74fd01671cbec1b3788a1f49417e")
}

// UnpackUnknownAssetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func (basket *Basket) UnpackUnknownAssetError(raw []byte) (*BasketUnknownAsset, error) {
	out := new(BasketUnknownAsset)
	if err := basket.abi.UnpackIntoInterface(out, "UnknownAsset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketValueLost represents a ValueLost error raised by the Basket contract.
type BasketValueLost struct {
	ValueIn  *big.Int
	ValueOut *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ValueLost(uint256 valueIn, uint256 valueOut)
func BasketValueLostErrorID() common.Hash {
	return common.HexToHash("0x0ad891bb90ff8de73a9f77e9241befa2bc734cb037bb1c0746b46f91da59f4f2")
}

// UnpackValueLostError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ValueLost(uint256 valueIn, uint256 valueOut)
func (basket *Basket) UnpackValueLostError(raw []byte) (*BasketValueLost, error) {
	out := new(BasketValueLost)
	if err := basket.abi.UnpackIntoInterface(out, "ValueLost", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BasketZeroAmount represents a ZeroAmount error raised by the Basket contract.
type BasketZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAmount()
func BasketZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x1f2a2005cb66a8e145327e8814e243a0996aec9bfe1e15a495778b1236dbd485")
}

// UnpackZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAmount()
func (basket *Basket) UnpackZeroAmountError(raw []byte) (*BasketZeroAmount, error) {
	out := new(BasketZeroAmount)
	if err := basket.abi.UnpackIntoInterface(out, "ZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
