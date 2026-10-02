// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package shortpositions

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

// ShortPositionsMetaData contains all meta data concerning the ShortPositions contract.
var ShortPositionsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"accounts_\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"},{\"name\":\"router_\",\"type\":\"address\",\"internalType\":\"contractIV3SwapRouter\"},{\"name\":\"fees\",\"type\":\"uint24[]\",\"internalType\":\"uint24[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"LIQUIDATION_BONUS_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_PREMIUM_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RESTRICTION_DROP_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"book\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"usdgHeld\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"costIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyIn\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cover\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"cost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"engine\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMargin\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epoch\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"health\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"missing\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"regime\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidate\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mark\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"position\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"usdgHeld\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"debt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"epoch_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"restriction\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"restricted\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"close\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIV3SwapRouter\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sell\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minProceeds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"proceeds\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usdg\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weekendLeverage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BuyIn\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Cover\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Liquidate\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"bonus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"writtenOff\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewEpoch\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Restricted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"until\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sell\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"proceeds\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AboveLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AssetHalted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Blocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BorrowingIsPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CorporateActionPending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DeficitOpen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Insolvent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientCollateral\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientMargin\",\"inputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidPool\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"}]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LiquidityUnknown\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotLending\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotShortable\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotShortfall\",\"inputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedIntToUint\",\"inputs\":[{\"name\":\"value\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintToInt\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SessionUnknown\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WeekendLeverageExceeded\",\"inputs\":[{\"name\":\"gross\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
	ID:  "ShortPositions",
}

// ShortPositions is an auto generated Go binding around an Ethereum contract.
type ShortPositions struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *ShortPositions) GetABI() abi.ABI {
	return c.abi
}

// NewShortPositions creates a new instance of ShortPositions.
func NewShortPositions() *ShortPositions {
	parsed, err := ShortPositionsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ShortPositions{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ShortPositions) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address accounts_, address router_, uint24[] fees) returns()
func (shortPositions *ShortPositions) PackConstructor(accounts_ common.Address, router_ common.Address, fees []*big.Int) []byte {
	enc, err := shortPositions.abi.Pack("", accounts_, router_, fees)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackLIQUIDATIONBONUSBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecd1dae6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function LIQUIDATION_BONUS_BPS() view returns(uint256)
func (shortPositions *ShortPositions) PackLIQUIDATIONBONUSBPS() []byte {
	enc, err := shortPositions.abi.Pack("LIQUIDATION_BONUS_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLIQUIDATIONBONUSBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecd1dae6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function LIQUIDATION_BONUS_BPS() view returns(uint256)
func (shortPositions *ShortPositions) TryPackLIQUIDATIONBONUSBPS() ([]byte, error) {
	return shortPositions.abi.Pack("LIQUIDATION_BONUS_BPS")
}

// UnpackLIQUIDATIONBONUSBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecd1dae6.
//
// Solidity: function LIQUIDATION_BONUS_BPS() view returns(uint256)
func (shortPositions *ShortPositions) UnpackLIQUIDATIONBONUSBPS(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("LIQUIDATION_BONUS_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXPREMIUMBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x711b7a0c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_PREMIUM_BPS() view returns(uint256)
func (shortPositions *ShortPositions) PackMAXPREMIUMBPS() []byte {
	enc, err := shortPositions.abi.Pack("MAX_PREMIUM_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXPREMIUMBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x711b7a0c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_PREMIUM_BPS() view returns(uint256)
func (shortPositions *ShortPositions) TryPackMAXPREMIUMBPS() ([]byte, error) {
	return shortPositions.abi.Pack("MAX_PREMIUM_BPS")
}

// UnpackMAXPREMIUMBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x711b7a0c.
//
// Solidity: function MAX_PREMIUM_BPS() view returns(uint256)
func (shortPositions *ShortPositions) UnpackMAXPREMIUMBPS(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("MAX_PREMIUM_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRESTRICTIONDROPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcb0b3bc2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function RESTRICTION_DROP_BPS() view returns(uint256)
func (shortPositions *ShortPositions) PackRESTRICTIONDROPBPS() []byte {
	enc, err := shortPositions.abi.Pack("RESTRICTION_DROP_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRESTRICTIONDROPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcb0b3bc2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function RESTRICTION_DROP_BPS() view returns(uint256)
func (shortPositions *ShortPositions) TryPackRESTRICTIONDROPBPS() ([]byte, error) {
	return shortPositions.abi.Pack("RESTRICTION_DROP_BPS")
}

// UnpackRESTRICTIONDROPBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcb0b3bc2.
//
// Solidity: function RESTRICTION_DROP_BPS() view returns(uint256)
func (shortPositions *ShortPositions) UnpackRESTRICTIONDROPBPS(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("RESTRICTION_DROP_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAccounts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68cd03f6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function accounts() view returns(address)
func (shortPositions *ShortPositions) PackAccounts() []byte {
	enc, err := shortPositions.abi.Pack("accounts")
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
func (shortPositions *ShortPositions) TryPackAccounts() ([]byte, error) {
	return shortPositions.abi.Pack("accounts")
}

// UnpackAccounts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68cd03f6.
//
// Solidity: function accounts() view returns(address)
func (shortPositions *ShortPositions) UnpackAccounts(data []byte) (common.Address, error) {
	out, err := shortPositions.abi.Unpack("accounts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (shortPositions *ShortPositions) PackBand() []byte {
	enc, err := shortPositions.abi.Pack("band")
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
func (shortPositions *ShortPositions) TryPackBand() ([]byte, error) {
	return shortPositions.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (shortPositions *ShortPositions) UnpackBand(data []byte) (common.Address, error) {
	out, err := shortPositions.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBook is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56847cbd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function book(bytes32 symbol, uint256 epoch_) view returns(uint256 shares, uint256 usdgHeld, uint256 costIndex)
func (shortPositions *ShortPositions) PackBook(symbol [32]byte, epoch *big.Int) []byte {
	enc, err := shortPositions.abi.Pack("book", symbol, epoch)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBook is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56847cbd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function book(bytes32 symbol, uint256 epoch_) view returns(uint256 shares, uint256 usdgHeld, uint256 costIndex)
func (shortPositions *ShortPositions) TryPackBook(symbol [32]byte, epoch *big.Int) ([]byte, error) {
	return shortPositions.abi.Pack("book", symbol, epoch)
}

// BookOutput serves as a container for the return parameters of contract
// method Book.
type BookOutput struct {
	Shares    *big.Int
	UsdgHeld  *big.Int
	CostIndex *big.Int
}

// UnpackBook is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56847cbd.
//
// Solidity: function book(bytes32 symbol, uint256 epoch_) view returns(uint256 shares, uint256 usdgHeld, uint256 costIndex)
func (shortPositions *ShortPositions) UnpackBook(data []byte) (BookOutput, error) {
	out, err := shortPositions.abi.Unpack("book", data)
	outstruct := new(BookOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Shares = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.UsdgHeld = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.CostIndex = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackBuyIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb287d8ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buyIn(uint256 assets) returns()
func (shortPositions *ShortPositions) PackBuyIn(assets *big.Int) []byte {
	enc, err := shortPositions.abi.Pack("buyIn", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuyIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb287d8ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buyIn(uint256 assets) returns()
func (shortPositions *ShortPositions) TryPackBuyIn(assets *big.Int) ([]byte, error) {
	return shortPositions.abi.Pack("buyIn", assets)
}

// PackCover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5048818.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cover(bytes32 symbol, uint256 amount, uint256 maxCost, address account) returns(uint256 cost)
func (shortPositions *ShortPositions) PackCover(symbol [32]byte, amount *big.Int, maxCost *big.Int, account common.Address) []byte {
	enc, err := shortPositions.abi.Pack("cover", symbol, amount, maxCost, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCover is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5048818.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cover(bytes32 symbol, uint256 amount, uint256 maxCost, address account) returns(uint256 cost)
func (shortPositions *ShortPositions) TryPackCover(symbol [32]byte, amount *big.Int, maxCost *big.Int, account common.Address) ([]byte, error) {
	return shortPositions.abi.Pack("cover", symbol, amount, maxCost, account)
}

// UnpackCover is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb5048818.
//
// Solidity: function cover(bytes32 symbol, uint256 amount, uint256 maxCost, address account) returns(uint256 cost)
func (shortPositions *ShortPositions) UnpackCover(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("cover", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6ebf181b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(bytes32 symbol, uint256 amount, address account) returns()
func (shortPositions *ShortPositions) PackDeposit(symbol [32]byte, amount *big.Int, account common.Address) []byte {
	enc, err := shortPositions.abi.Pack("deposit", symbol, amount, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6ebf181b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(bytes32 symbol, uint256 amount, address account) returns()
func (shortPositions *ShortPositions) TryPackDeposit(symbol [32]byte, amount *big.Int, account common.Address) ([]byte, error) {
	return shortPositions.abi.Pack("deposit", symbol, amount, account)
}

// PackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function engine() view returns(address)
func (shortPositions *ShortPositions) PackEngine() []byte {
	enc, err := shortPositions.abi.Pack("engine")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function engine() view returns(address)
func (shortPositions *ShortPositions) TryPackEngine() ([]byte, error) {
	return shortPositions.abi.Pack("engine")
}

// UnpackEngine is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (shortPositions *ShortPositions) UnpackEngine(data []byte) (common.Address, error) {
	out, err := shortPositions.abi.Unpack("engine", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa814a90.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function epoch(bytes32 symbol) view returns(uint256)
func (shortPositions *ShortPositions) PackEpoch(symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("epoch", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEpoch is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa814a90.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function epoch(bytes32 symbol) view returns(uint256)
func (shortPositions *ShortPositions) TryPackEpoch(symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("epoch", symbol)
}

// UnpackEpoch is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfa814a90.
//
// Solidity: function epoch(bytes32 symbol) view returns(uint256)
func (shortPositions *ShortPositions) UnpackEpoch(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("epoch", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27cdab06.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function fee(bytes32 symbol) view returns(uint24)
func (shortPositions *ShortPositions) PackFee(symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("fee", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27cdab06.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function fee(bytes32 symbol) view returns(uint24)
func (shortPositions *ShortPositions) TryPackFee(symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("fee", symbol)
}

// UnpackFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27cdab06.
//
// Solidity: function fee(bytes32 symbol) view returns(uint24)
func (shortPositions *ShortPositions) UnpackFee(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("fee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHealth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42d6fb7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function health(address account, bytes32 symbol) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (shortPositions *ShortPositions) PackHealth(account common.Address, symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("health", account, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHealth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42d6fb7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function health(address account, bytes32 symbol) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (shortPositions *ShortPositions) TryPackHealth(account common.Address, symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("health", account, symbol)
}

// HealthOutput serves as a container for the return parameters of contract
// method Health.
type HealthOutput struct {
	Equity      *big.Int
	Requirement *big.Int
	Missing     uint8
	Regime      uint8
}

// UnpackHealth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42d6fb7b.
//
// Solidity: function health(address account, bytes32 symbol) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (shortPositions *ShortPositions) UnpackHealth(data []byte) (HealthOutput, error) {
	out, err := shortPositions.abi.Unpack("health", data)
	outstruct := new(HealthOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Equity = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Requirement = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Missing = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.Regime = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackLiquidate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3af8c7a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidate(address account, bytes32 symbol) returns()
func (shortPositions *ShortPositions) PackLiquidate(account common.Address, symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("liquidate", account, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLiquidate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3af8c7a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function liquidate(address account, bytes32 symbol) returns()
func (shortPositions *ShortPositions) TryPackLiquidate(account common.Address, symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("liquidate", account, symbol)
}

// PackMark is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xee02438e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mark(bytes32 symbol) returns()
func (shortPositions *ShortPositions) PackMark(symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("mark", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMark is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xee02438e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mark(bytes32 symbol) returns()
func (shortPositions *ShortPositions) TryPackMark(symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("mark", symbol)
}

// PackPosition is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefe10643.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function position(address account, bytes32 symbol) view returns(int256 usdgHeld, uint256 debt, uint256 shares, uint256 epoch_)
func (shortPositions *ShortPositions) PackPosition(account common.Address, symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("position", account, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPosition is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefe10643.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function position(address account, bytes32 symbol) view returns(int256 usdgHeld, uint256 debt, uint256 shares, uint256 epoch_)
func (shortPositions *ShortPositions) TryPackPosition(account common.Address, symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("position", account, symbol)
}

// PositionOutput serves as a container for the return parameters of contract
// method Position.
type PositionOutput struct {
	UsdgHeld *big.Int
	Debt     *big.Int
	Shares   *big.Int
	Epoch    *big.Int
}

// UnpackPosition is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefe10643.
//
// Solidity: function position(address account, bytes32 symbol) view returns(int256 usdgHeld, uint256 debt, uint256 shares, uint256 epoch_)
func (shortPositions *ShortPositions) UnpackPosition(data []byte) (PositionOutput, error) {
	out, err := shortPositions.abi.Unpack("position", data)
	outstruct := new(PositionOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.UsdgHeld = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Debt = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Shares = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.Epoch = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRestriction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4d866688.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function restriction(bytes32 symbol) view returns(bool restricted, uint64 close)
func (shortPositions *ShortPositions) PackRestriction(symbol [32]byte) []byte {
	enc, err := shortPositions.abi.Pack("restriction", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRestriction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4d866688.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function restriction(bytes32 symbol) view returns(bool restricted, uint64 close)
func (shortPositions *ShortPositions) TryPackRestriction(symbol [32]byte) ([]byte, error) {
	return shortPositions.abi.Pack("restriction", symbol)
}

// RestrictionOutput serves as a container for the return parameters of contract
// method Restriction.
type RestrictionOutput struct {
	Restricted bool
	Close      uint64
}

// UnpackRestriction is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4d866688.
//
// Solidity: function restriction(bytes32 symbol) view returns(bool restricted, uint64 close)
func (shortPositions *ShortPositions) UnpackRestriction(data []byte) (RestrictionOutput, error) {
	out, err := shortPositions.abi.Unpack("restriction", data)
	outstruct := new(RestrictionOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Restricted = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Close = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackRouter is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf887ea40.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function router() view returns(address)
func (shortPositions *ShortPositions) PackRouter() []byte {
	enc, err := shortPositions.abi.Pack("router")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRouter is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf887ea40.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function router() view returns(address)
func (shortPositions *ShortPositions) TryPackRouter() ([]byte, error) {
	return shortPositions.abi.Pack("router")
}

// UnpackRouter is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (shortPositions *ShortPositions) UnpackRouter(data []byte) (common.Address, error) {
	out, err := shortPositions.abi.Unpack("router", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackSell is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x433af24b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sell(bytes32 symbol, uint256 amount, uint256 minProceeds, address account) returns(uint256 proceeds)
func (shortPositions *ShortPositions) PackSell(symbol [32]byte, amount *big.Int, minProceeds *big.Int, account common.Address) []byte {
	enc, err := shortPositions.abi.Pack("sell", symbol, amount, minProceeds, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSell is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x433af24b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sell(bytes32 symbol, uint256 amount, uint256 minProceeds, address account) returns(uint256 proceeds)
func (shortPositions *ShortPositions) TryPackSell(symbol [32]byte, amount *big.Int, minProceeds *big.Int, account common.Address) ([]byte, error) {
	return shortPositions.abi.Pack("sell", symbol, amount, minProceeds, account)
}

// UnpackSell is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x433af24b.
//
// Solidity: function sell(bytes32 symbol, uint256 amount, uint256 minProceeds, address account) returns(uint256 proceeds)
func (shortPositions *ShortPositions) UnpackSell(data []byte) (*big.Int, error) {
	out, err := shortPositions.abi.Unpack("sell", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function usdg() view returns(address)
func (shortPositions *ShortPositions) PackUsdg() []byte {
	enc, err := shortPositions.abi.Pack("usdg")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function usdg() view returns(address)
func (shortPositions *ShortPositions) TryPackUsdg() ([]byte, error) {
	return shortPositions.abi.Pack("usdg")
}

// UnpackUsdg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5b91b7b.
//
// Solidity: function usdg() view returns(address)
func (shortPositions *ShortPositions) UnpackUsdg(data []byte) (common.Address, error) {
	out, err := shortPositions.abi.Unpack("usdg", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWeekendLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb659ad89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (shortPositions *ShortPositions) PackWeekendLeverage() []byte {
	enc, err := shortPositions.abi.Pack("weekendLeverage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeekendLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb659ad89.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (shortPositions *ShortPositions) TryPackWeekendLeverage() ([]byte, error) {
	return shortPositions.abi.Pack("weekendLeverage")
}

// UnpackWeekendLeverage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb659ad89.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (shortPositions *ShortPositions) UnpackWeekendLeverage(data []byte) (uint32, error) {
	out, err := shortPositions.abi.Unpack("weekendLeverage", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1e6ba5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(bytes32 symbol, uint256 amount, address account, address receiver) returns()
func (shortPositions *ShortPositions) PackWithdraw(symbol [32]byte, amount *big.Int, account common.Address, receiver common.Address) []byte {
	enc, err := shortPositions.abi.Pack("withdraw", symbol, amount, account, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1e6ba5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(bytes32 symbol, uint256 amount, address account, address receiver) returns()
func (shortPositions *ShortPositions) TryPackWithdraw(symbol [32]byte, amount *big.Int, account common.Address, receiver common.Address) ([]byte, error) {
	return shortPositions.abi.Pack("withdraw", symbol, amount, account, receiver)
}

// ShortPositionsBuyIn represents a BuyIn event raised by the ShortPositions contract.
type ShortPositionsBuyIn struct {
	Symbol [32]byte
	Amount *big.Int
	Cost   *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ShortPositionsBuyInEventName = "BuyIn"

// ContractEventName returns the user-defined event name.
func (ShortPositionsBuyIn) ContractEventName() string {
	return ShortPositionsBuyInEventName
}

// UnpackBuyInEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BuyIn(bytes32 indexed symbol, uint256 amount, uint256 cost)
func (shortPositions *ShortPositions) UnpackBuyInEvent(log *types.Log) (*ShortPositionsBuyIn, error) {
	event := "BuyIn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsBuyIn)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsCover represents a Cover event raised by the ShortPositions contract.
type ShortPositionsCover struct {
	Account common.Address
	Symbol  [32]byte
	Amount  *big.Int
	Cost    *big.Int
	Shares  *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ShortPositionsCoverEventName = "Cover"

// ContractEventName returns the user-defined event name.
func (ShortPositionsCover) ContractEventName() string {
	return ShortPositionsCoverEventName
}

// UnpackCoverEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Cover(address indexed account, bytes32 indexed symbol, uint256 amount, uint256 cost, uint256 shares)
func (shortPositions *ShortPositions) UnpackCoverEvent(log *types.Log) (*ShortPositionsCover, error) {
	event := "Cover"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsCover)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsDeposit represents a Deposit event raised by the ShortPositions contract.
type ShortPositionsDeposit struct {
	Caller  common.Address
	Account common.Address
	Symbol  [32]byte
	Amount  *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const ShortPositionsDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (ShortPositionsDeposit) ContractEventName() string {
	return ShortPositionsDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed caller, address indexed account, bytes32 indexed symbol, uint256 amount)
func (shortPositions *ShortPositions) UnpackDepositEvent(log *types.Log) (*ShortPositionsDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsDeposit)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsLiquidate represents a Liquidate event raised by the ShortPositions contract.
type ShortPositionsLiquidate struct {
	Caller     common.Address
	Account    common.Address
	Symbol     [32]byte
	Amount     *big.Int
	Cost       *big.Int
	Bonus      *big.Int
	WrittenOff *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const ShortPositionsLiquidateEventName = "Liquidate"

// ContractEventName returns the user-defined event name.
func (ShortPositionsLiquidate) ContractEventName() string {
	return ShortPositionsLiquidateEventName
}

// UnpackLiquidateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Liquidate(address indexed caller, address indexed account, bytes32 indexed symbol, uint256 amount, uint256 cost, uint256 bonus, uint256 writtenOff)
func (shortPositions *ShortPositions) UnpackLiquidateEvent(log *types.Log) (*ShortPositionsLiquidate, error) {
	event := "Liquidate"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsLiquidate)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsNewEpoch represents a NewEpoch event raised by the ShortPositions contract.
type ShortPositionsNewEpoch struct {
	Symbol [32]byte
	Epoch  *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ShortPositionsNewEpochEventName = "NewEpoch"

// ContractEventName returns the user-defined event name.
func (ShortPositionsNewEpoch) ContractEventName() string {
	return ShortPositionsNewEpochEventName
}

// UnpackNewEpochEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event NewEpoch(bytes32 indexed symbol, uint256 epoch)
func (shortPositions *ShortPositions) UnpackNewEpochEvent(log *types.Log) (*ShortPositionsNewEpoch, error) {
	event := "NewEpoch"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsNewEpoch)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsRestricted represents a Restricted event raised by the ShortPositions contract.
type ShortPositionsRestricted struct {
	Symbol [32]byte
	Until  uint32
	Raw    *types.Log // Blockchain specific contextual infos
}

const ShortPositionsRestrictedEventName = "Restricted"

// ContractEventName returns the user-defined event name.
func (ShortPositionsRestricted) ContractEventName() string {
	return ShortPositionsRestrictedEventName
}

// UnpackRestrictedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Restricted(bytes32 indexed symbol, uint32 until)
func (shortPositions *ShortPositions) UnpackRestrictedEvent(log *types.Log) (*ShortPositionsRestricted, error) {
	event := "Restricted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsRestricted)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsSell represents a Sell event raised by the ShortPositions contract.
type ShortPositionsSell struct {
	Account  common.Address
	Symbol   [32]byte
	Amount   *big.Int
	Proceeds *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ShortPositionsSellEventName = "Sell"

// ContractEventName returns the user-defined event name.
func (ShortPositionsSell) ContractEventName() string {
	return ShortPositionsSellEventName
}

// UnpackSellEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sell(address indexed account, bytes32 indexed symbol, uint256 amount, uint256 proceeds, uint256 shares)
func (shortPositions *ShortPositions) UnpackSellEvent(log *types.Log) (*ShortPositionsSell, error) {
	event := "Sell"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsSell)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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

// ShortPositionsWithdraw represents a Withdraw event raised by the ShortPositions contract.
type ShortPositionsWithdraw struct {
	Caller   common.Address
	Account  common.Address
	Symbol   [32]byte
	Amount   *big.Int
	Receiver common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const ShortPositionsWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (ShortPositionsWithdraw) ContractEventName() string {
	return ShortPositionsWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed caller, address indexed account, bytes32 indexed symbol, uint256 amount, address receiver)
func (shortPositions *ShortPositions) UnpackWithdrawEvent(log *types.Log) (*ShortPositionsWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != shortPositions.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ShortPositionsWithdraw)
	if len(log.Data) > 0 {
		if err := shortPositions.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range shortPositions.abi.Events[event].Inputs {
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
func (shortPositions *ShortPositions) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["AboveLimit"].ID.Bytes()[:4]) {
		return shortPositions.UnpackAboveLimitError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["AssetHalted"].ID.Bytes()[:4]) {
		return shortPositions.UnpackAssetHaltedError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["Blocked"].ID.Bytes()[:4]) {
		return shortPositions.UnpackBlockedError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["BorrowingIsPaused"].ID.Bytes()[:4]) {
		return shortPositions.UnpackBorrowingIsPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["CorporateActionPending"].ID.Bytes()[:4]) {
		return shortPositions.UnpackCorporateActionPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["DeficitOpen"].ID.Bytes()[:4]) {
		return shortPositions.UnpackDeficitOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["Insolvent"].ID.Bytes()[:4]) {
		return shortPositions.UnpackInsolventError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["InsufficientCollateral"].ID.Bytes()[:4]) {
		return shortPositions.UnpackInsufficientCollateralError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["InsufficientMargin"].ID.Bytes()[:4]) {
		return shortPositions.UnpackInsufficientMarginError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["InvalidPool"].ID.Bytes()[:4]) {
		return shortPositions.UnpackInvalidPoolError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["LengthMismatch"].ID.Bytes()[:4]) {
		return shortPositions.UnpackLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["LiquidityUnknown"].ID.Bytes()[:4]) {
		return shortPositions.UnpackLiquidityUnknownError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["NoShort"].ID.Bytes()[:4]) {
		return shortPositions.UnpackNoShortError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["NotLending"].ID.Bytes()[:4]) {
		return shortPositions.UnpackNotLendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["NotShortable"].ID.Bytes()[:4]) {
		return shortPositions.UnpackNotShortableError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["NotShortfall"].ID.Bytes()[:4]) {
		return shortPositions.UnpackNotShortfallError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SafeCastOverflowedIntToUint"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSafeCastOverflowedIntToUintError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SafeCastOverflowedUintToInt"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSafeCastOverflowedUintToIntError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSequencerNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["SessionUnknown"].ID.Bytes()[:4]) {
		return shortPositions.UnpackSessionUnknownError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["Unauthorized"].ID.Bytes()[:4]) {
		return shortPositions.UnpackUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["WeekendLeverageExceeded"].ID.Bytes()[:4]) {
		return shortPositions.UnpackWeekendLeverageExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], shortPositions.abi.Errors["ZeroAmount"].ID.Bytes()[:4]) {
		return shortPositions.UnpackZeroAmountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ShortPositionsAboveLimit represents a AboveLimit error raised by the ShortPositions contract.
type ShortPositionsAboveLimit struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AboveLimit()
func ShortPositionsAboveLimitErrorID() common.Hash {
	return common.HexToHash("0x0b9f18c808720c1563bb012343d71654417f3c7fd13a919e54436b75dfc942ff")
}

// UnpackAboveLimitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AboveLimit()
func (shortPositions *ShortPositions) UnpackAboveLimitError(raw []byte) (*ShortPositionsAboveLimit, error) {
	out := new(ShortPositionsAboveLimit)
	if err := shortPositions.abi.UnpackIntoInterface(out, "AboveLimit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsAssetHalted represents a AssetHalted error raised by the ShortPositions contract.
type ShortPositionsAssetHalted struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetHalted(bytes32 symbol)
func ShortPositionsAssetHaltedErrorID() common.Hash {
	return common.HexToHash("0x3ec29ca3785f47585a9c96baed33ddbd8b2d9e8dfdeffb4be96925aaf9ec6468")
}

// UnpackAssetHaltedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetHalted(bytes32 symbol)
func (shortPositions *ShortPositions) UnpackAssetHaltedError(raw []byte) (*ShortPositionsAssetHalted, error) {
	out := new(ShortPositionsAssetHalted)
	if err := shortPositions.abi.UnpackIntoInterface(out, "AssetHalted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsBlocked represents a Blocked error raised by the ShortPositions contract.
type ShortPositionsBlocked struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Blocked(address account)
func ShortPositionsBlockedErrorID() common.Hash {
	return common.HexToHash("0x75e91ce73c1d3352d8dd3610443539cd33dfe13b1de8f8caae54ec26dd0dc9cb")
}

// UnpackBlockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Blocked(address account)
func (shortPositions *ShortPositions) UnpackBlockedError(raw []byte) (*ShortPositionsBlocked, error) {
	out := new(ShortPositionsBlocked)
	if err := shortPositions.abi.UnpackIntoInterface(out, "Blocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsBorrowingIsPaused represents a BorrowingIsPaused error raised by the ShortPositions contract.
type ShortPositionsBorrowingIsPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BorrowingIsPaused()
func ShortPositionsBorrowingIsPausedErrorID() common.Hash {
	return common.HexToHash("0xaa864f37b700d896a72eafa272084f861aa9962378d6a62f8229dc64289fbc6a")
}

// UnpackBorrowingIsPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BorrowingIsPaused()
func (shortPositions *ShortPositions) UnpackBorrowingIsPausedError(raw []byte) (*ShortPositionsBorrowingIsPaused, error) {
	out := new(ShortPositionsBorrowingIsPaused)
	if err := shortPositions.abi.UnpackIntoInterface(out, "BorrowingIsPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsCorporateActionPending represents a CorporateActionPending error raised by the ShortPositions contract.
type ShortPositionsCorporateActionPending struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func ShortPositionsCorporateActionPendingErrorID() common.Hash {
	return common.HexToHash("0x1d1f5fd3c689a2e2561007ef605d914d864ac52012393f43a887a76f9f48be7b")
}

// UnpackCorporateActionPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func (shortPositions *ShortPositions) UnpackCorporateActionPendingError(raw []byte) (*ShortPositionsCorporateActionPending, error) {
	out := new(ShortPositionsCorporateActionPending)
	if err := shortPositions.abi.UnpackIntoInterface(out, "CorporateActionPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsDeficitOpen represents a DeficitOpen error raised by the ShortPositions contract.
type ShortPositionsDeficitOpen struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DeficitOpen()
func ShortPositionsDeficitOpenErrorID() common.Hash {
	return common.HexToHash("0x4ab6ae5d7df45670f2a7befd73acf72971ee7436e3a372bf2c6701a18a44c44b")
}

// UnpackDeficitOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DeficitOpen()
func (shortPositions *ShortPositions) UnpackDeficitOpenError(raw []byte) (*ShortPositionsDeficitOpen, error) {
	out := new(ShortPositionsDeficitOpen)
	if err := shortPositions.abi.UnpackIntoInterface(out, "DeficitOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsInsolvent represents a Insolvent error raised by the ShortPositions contract.
type ShortPositionsInsolvent struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Insolvent()
func ShortPositionsInsolventErrorID() common.Hash {
	return common.HexToHash("0xfc220038dd9537de7655560b6afd06dcea02819c57f7ce0a390d00abb0854ca4")
}

// UnpackInsolventError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Insolvent()
func (shortPositions *ShortPositions) UnpackInsolventError(raw []byte) (*ShortPositionsInsolvent, error) {
	out := new(ShortPositionsInsolvent)
	if err := shortPositions.abi.UnpackIntoInterface(out, "Insolvent", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsInsufficientCollateral represents a InsufficientCollateral error raised by the ShortPositions contract.
type ShortPositionsInsufficientCollateral struct {
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientCollateral(uint256 amount)
func ShortPositionsInsufficientCollateralErrorID() common.Hash {
	return common.HexToHash("0x2b3bc98557cfd9172d79165d096e2f1f67c00ee208ef0ac5df06614762a558fb")
}

// UnpackInsufficientCollateralError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientCollateral(uint256 amount)
func (shortPositions *ShortPositions) UnpackInsufficientCollateralError(raw []byte) (*ShortPositionsInsufficientCollateral, error) {
	out := new(ShortPositionsInsufficientCollateral)
	if err := shortPositions.abi.UnpackIntoInterface(out, "InsufficientCollateral", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsInsufficientMargin represents a InsufficientMargin error raised by the ShortPositions contract.
type ShortPositionsInsufficientMargin struct {
	Equity      *big.Int
	Requirement *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientMargin(int256 equity, uint256 requirement)
func ShortPositionsInsufficientMarginErrorID() common.Hash {
	return common.HexToHash("0x1e4cbb75332dfc383c902bc80a5a1b1f64daea3dcc8e970b2d2b447b2b7640ea")
}

// UnpackInsufficientMarginError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientMargin(int256 equity, uint256 requirement)
func (shortPositions *ShortPositions) UnpackInsufficientMarginError(raw []byte) (*ShortPositionsInsufficientMargin, error) {
	out := new(ShortPositionsInsufficientMargin)
	if err := shortPositions.abi.UnpackIntoInterface(out, "InsufficientMargin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsInvalidPool represents a InvalidPool error raised by the ShortPositions contract.
type ShortPositionsInvalidPool struct {
	Token common.Address
	Fee   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPool(address token, uint24 fee)
func ShortPositionsInvalidPoolErrorID() common.Hash {
	return common.HexToHash("0x221f6255948534b899aacd356721c6aeaf9ab194e4028f46124e81d38bf6a4cf")
}

// UnpackInvalidPoolError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPool(address token, uint24 fee)
func (shortPositions *ShortPositions) UnpackInvalidPoolError(raw []byte) (*ShortPositionsInvalidPool, error) {
	out := new(ShortPositionsInvalidPool)
	if err := shortPositions.abi.UnpackIntoInterface(out, "InvalidPool", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsLengthMismatch represents a LengthMismatch error raised by the ShortPositions contract.
type ShortPositionsLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LengthMismatch()
func ShortPositionsLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xff633a3803c58b9bc21e58efecee59f27e033cc0b1883fccb4969c76146fe60f")
}

// UnpackLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LengthMismatch()
func (shortPositions *ShortPositions) UnpackLengthMismatchError(raw []byte) (*ShortPositionsLengthMismatch, error) {
	out := new(ShortPositionsLengthMismatch)
	if err := shortPositions.abi.UnpackIntoInterface(out, "LengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsLiquidityUnknown represents a LiquidityUnknown error raised by the ShortPositions contract.
type ShortPositionsLiquidityUnknown struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LiquidityUnknown(bytes32 symbol)
func ShortPositionsLiquidityUnknownErrorID() common.Hash {
	return common.HexToHash("0x1e89df1a02046ee0b3974a0ccb9d6a45cd1e47b611031710fdf5be384e85f153")
}

// UnpackLiquidityUnknownError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LiquidityUnknown(bytes32 symbol)
func (shortPositions *ShortPositions) UnpackLiquidityUnknownError(raw []byte) (*ShortPositionsLiquidityUnknown, error) {
	out := new(ShortPositionsLiquidityUnknown)
	if err := shortPositions.abi.UnpackIntoInterface(out, "LiquidityUnknown", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsNoShort represents a NoShort error raised by the ShortPositions contract.
type ShortPositionsNoShort struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoShort()
func ShortPositionsNoShortErrorID() common.Hash {
	return common.HexToHash("0x13cd0df2fd8cf2143cebdb40a5387e7ec8ed269fc2854e6674946bd90c55a720")
}

// UnpackNoShortError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoShort()
func (shortPositions *ShortPositions) UnpackNoShortError(raw []byte) (*ShortPositionsNoShort, error) {
	out := new(ShortPositionsNoShort)
	if err := shortPositions.abi.UnpackIntoInterface(out, "NoShort", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsNotLending represents a NotLending error raised by the ShortPositions contract.
type ShortPositionsNotLending struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotLending(address caller)
func ShortPositionsNotLendingErrorID() common.Hash {
	return common.HexToHash("0x5390acfb124bea68b81348300c84aa23f4a80c1e5b44ad0142e9bec0db024cbd")
}

// UnpackNotLendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotLending(address caller)
func (shortPositions *ShortPositions) UnpackNotLendingError(raw []byte) (*ShortPositionsNotLending, error) {
	out := new(ShortPositionsNotLending)
	if err := shortPositions.abi.UnpackIntoInterface(out, "NotLending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsNotShortable represents a NotShortable error raised by the ShortPositions contract.
type ShortPositionsNotShortable struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotShortable(bytes32 symbol)
func ShortPositionsNotShortableErrorID() common.Hash {
	return common.HexToHash("0x66f454b882ac7da427de662929462ef53cae8d6a31e274ea25192e11f1e041ae")
}

// UnpackNotShortableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotShortable(bytes32 symbol)
func (shortPositions *ShortPositions) UnpackNotShortableError(raw []byte) (*ShortPositionsNotShortable, error) {
	out := new(ShortPositionsNotShortable)
	if err := shortPositions.abi.UnpackIntoInterface(out, "NotShortable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsNotShortfall represents a NotShortfall error raised by the ShortPositions contract.
type ShortPositionsNotShortfall struct {
	Equity      *big.Int
	Requirement *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotShortfall(int256 equity, uint256 requirement)
func ShortPositionsNotShortfallErrorID() common.Hash {
	return common.HexToHash("0x77ef61c8473dcc489476a72548df3c8516a92dfc059ed98dd0d451cdb5135d27")
}

// UnpackNotShortfallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotShortfall(int256 equity, uint256 requirement)
func (shortPositions *ShortPositions) UnpackNotShortfallError(raw []byte) (*ShortPositionsNotShortfall, error) {
	out := new(ShortPositionsNotShortfall)
	if err := shortPositions.abi.UnpackIntoInterface(out, "NotShortfall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSafeCastOverflowedIntToUint represents a SafeCastOverflowedIntToUint error raised by the ShortPositions contract.
type ShortPositionsSafeCastOverflowedIntToUint struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func ShortPositionsSafeCastOverflowedIntToUintErrorID() common.Hash {
	return common.HexToHash("0xa8ce4432b175c373e5f41aba830358e5361584f628450fd436c066323ad91ac2")
}

// UnpackSafeCastOverflowedIntToUintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func (shortPositions *ShortPositions) UnpackSafeCastOverflowedIntToUintError(raw []byte) (*ShortPositionsSafeCastOverflowedIntToUint, error) {
	out := new(ShortPositionsSafeCastOverflowedIntToUint)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntToUint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the ShortPositions contract.
type ShortPositionsSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func ShortPositionsSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (shortPositions *ShortPositions) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*ShortPositionsSafeCastOverflowedUintDowncast, error) {
	out := new(ShortPositionsSafeCastOverflowedUintDowncast)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSafeCastOverflowedUintToInt represents a SafeCastOverflowedUintToInt error raised by the ShortPositions contract.
type ShortPositionsSafeCastOverflowedUintToInt struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func ShortPositionsSafeCastOverflowedUintToIntErrorID() common.Hash {
	return common.HexToHash("0x24775e0629ae69d78c11bae050651b81820407f300ff750ff2be51e4ce75c37f")
}

// UnpackSafeCastOverflowedUintToIntError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func (shortPositions *ShortPositions) UnpackSafeCastOverflowedUintToIntError(raw []byte) (*ShortPositionsSafeCastOverflowedUintToInt, error) {
	out := new(ShortPositionsSafeCastOverflowedUintToInt)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintToInt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the ShortPositions contract.
type ShortPositionsSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func ShortPositionsSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (shortPositions *ShortPositions) UnpackSafeERC20FailedOperationError(raw []byte) (*ShortPositionsSafeERC20FailedOperation, error) {
	out := new(ShortPositionsSafeERC20FailedOperation)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSequencerNotSettled represents a SequencerNotSettled error raised by the ShortPositions contract.
type ShortPositionsSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func ShortPositionsSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (shortPositions *ShortPositions) UnpackSequencerNotSettledError(raw []byte) (*ShortPositionsSequencerNotSettled, error) {
	out := new(ShortPositionsSequencerNotSettled)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsSessionUnknown represents a SessionUnknown error raised by the ShortPositions contract.
type ShortPositionsSessionUnknown struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SessionUnknown()
func ShortPositionsSessionUnknownErrorID() common.Hash {
	return common.HexToHash("0xbeb78049726639f27e32542d3c6bcfd536c907b423b81f72607e067b8770b849")
}

// UnpackSessionUnknownError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SessionUnknown()
func (shortPositions *ShortPositions) UnpackSessionUnknownError(raw []byte) (*ShortPositionsSessionUnknown, error) {
	out := new(ShortPositionsSessionUnknown)
	if err := shortPositions.abi.UnpackIntoInterface(out, "SessionUnknown", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsUnauthorized represents a Unauthorized error raised by the ShortPositions contract.
type ShortPositionsUnauthorized struct {
	Caller  common.Address
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Unauthorized(address caller, address account)
func ShortPositionsUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0x295a81c15f8ed30f90d2dbdefc77f34f2603dc4ea3df178ac206913987b62f4b")
}

// UnpackUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Unauthorized(address caller, address account)
func (shortPositions *ShortPositions) UnpackUnauthorizedError(raw []byte) (*ShortPositionsUnauthorized, error) {
	out := new(ShortPositionsUnauthorized)
	if err := shortPositions.abi.UnpackIntoInterface(out, "Unauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsWeekendLeverageExceeded represents a WeekendLeverageExceeded error raised by the ShortPositions contract.
type ShortPositionsWeekendLeverageExceeded struct {
	Gross  *big.Int
	Equity *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WeekendLeverageExceeded(uint256 gross, int256 equity)
func ShortPositionsWeekendLeverageExceededErrorID() common.Hash {
	return common.HexToHash("0xc31b6ccbfd61becd2cc8cc7f5970ce1ff409b28eb055fd1fa3f3e5d8cb4d0479")
}

// UnpackWeekendLeverageExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WeekendLeverageExceeded(uint256 gross, int256 equity)
func (shortPositions *ShortPositions) UnpackWeekendLeverageExceededError(raw []byte) (*ShortPositionsWeekendLeverageExceeded, error) {
	out := new(ShortPositionsWeekendLeverageExceeded)
	if err := shortPositions.abi.UnpackIntoInterface(out, "WeekendLeverageExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ShortPositionsZeroAmount represents a ZeroAmount error raised by the ShortPositions contract.
type ShortPositionsZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAmount()
func ShortPositionsZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x1f2a2005cb66a8e145327e8814e243a0996aec9bfe1e15a495778b1236dbd485")
}

// UnpackZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAmount()
func (shortPositions *ShortPositions) UnpackZeroAmountError(raw []byte) (*ShortPositionsZeroAmount, error) {
	out := new(ShortPositionsZeroAmount)
	if err := shortPositions.abi.UnpackIntoInterface(out, "ZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
