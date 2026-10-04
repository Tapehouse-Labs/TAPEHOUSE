// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package margin

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

// MarginMetaData contains all meta data concerning the Margin contract.
var MarginMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"assets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"correlation\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"other\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"floor\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentRequirement\",\"inputs\":[{\"name\":\"quantities\",\"type\":\"int256[]\",\"internalType\":\"int256[]\"},{\"name\":\"prices\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"margin\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"missing\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"regime\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depth\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"selling\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"buying\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"sellingCeiling\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"buyingCeiling\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethUsdFeed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"market\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pool\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requirement\",\"inputs\":[{\"name\":\"quantities\",\"type\":\"int256[]\",\"internalType\":\"int256[]\"},{\"name\":\"prices\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"horizon\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"spansClosure\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"margin\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"missing\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setParameters\",\"inputs\":[{\"name\":\"volatilities\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"correlations\",\"type\":\"uint16[]\",\"internalType\":\"uint16[]\"},{\"name\":\"gaps\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"depths\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"volatility\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"floor\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weekendGap\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"floor\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weekendLeverage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CorrelationSet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"other\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepthSet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"selling\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"buying\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GapSet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VolatilitySet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CorrelationStepTooLarge\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"other\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"previous\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"value\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"DepthStepTooLarge\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"previous\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"GapStepTooLarge\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"previous\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidCorrelation\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"other\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"floor\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidDepth\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"ceiling\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidGap\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"floor\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"InvalidVolatility\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"floor\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPositiveDefinite\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnknownAsset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UpdateTooSoon\",\"inputs\":[{\"name\":\"nextUpdateAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"VolatilityStepTooLarge\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"previous\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"value\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]",
	ID:  "Margin",
}

// Margin is an auto generated Go binding around an Ethereum contract.
type Margin struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Margin) GetABI() abi.ABI {
	return c.abi
}

// NewMargin creates a new instance of Margin.
func NewMargin() *Margin {
	parsed, err := MarginMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Margin{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Margin) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71a97305.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function assets() view returns(bytes32[])
func (margin *Margin) PackAssets() []byte {
	enc, err := margin.abi.Pack("assets")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71a97305.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function assets() view returns(bytes32[])
func (margin *Margin) TryPackAssets() ([]byte, error) {
	return margin.abi.Pack("assets")
}

// UnpackAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x71a97305.
//
// Solidity: function assets() view returns(bytes32[])
func (margin *Margin) UnpackAssets(data []byte) ([][32]byte, error) {
	out, err := margin.abi.Unpack("assets", data)
	if err != nil {
		return *new([][32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (margin *Margin) PackBand() []byte {
	enc, err := margin.abi.Pack("band")
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
func (margin *Margin) TryPackBand() ([]byte, error) {
	return margin.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (margin *Margin) UnpackBand(data []byte) (common.Address, error) {
	out, err := margin.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCorrelation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d244f45.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function correlation(bytes32 symbol, bytes32 other) view returns(uint16 value, uint16 floor)
func (margin *Margin) PackCorrelation(symbol [32]byte, other [32]byte) []byte {
	enc, err := margin.abi.Pack("correlation", symbol, other)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCorrelation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d244f45.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function correlation(bytes32 symbol, bytes32 other) view returns(uint16 value, uint16 floor)
func (margin *Margin) TryPackCorrelation(symbol [32]byte, other [32]byte) ([]byte, error) {
	return margin.abi.Pack("correlation", symbol, other)
}

// CorrelationOutput serves as a container for the return parameters of contract
// method Correlation.
type CorrelationOutput struct {
	Value uint16
	Floor uint16
}

// UnpackCorrelation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1d244f45.
//
// Solidity: function correlation(bytes32 symbol, bytes32 other) view returns(uint16 value, uint16 floor)
func (margin *Margin) UnpackCorrelation(data []byte) (CorrelationOutput, error) {
	out, err := margin.abi.Unpack("correlation", data)
	outstruct := new(CorrelationOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = *abi.ConvertType(out[0], new(uint16)).(*uint16)
	outstruct.Floor = *abi.ConvertType(out[1], new(uint16)).(*uint16)
	return *outstruct, nil
}

// PackCurrentRequirement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x995bff62.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function currentRequirement(int256[] quantities, uint256[] prices) view returns(uint256 margin, uint8 missing, uint8 regime)
func (margin *Margin) PackCurrentRequirement(quantities []*big.Int, prices []*big.Int) []byte {
	enc, err := margin.abi.Pack("currentRequirement", quantities, prices)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCurrentRequirement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x995bff62.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function currentRequirement(int256[] quantities, uint256[] prices) view returns(uint256 margin, uint8 missing, uint8 regime)
func (margin *Margin) TryPackCurrentRequirement(quantities []*big.Int, prices []*big.Int) ([]byte, error) {
	return margin.abi.Pack("currentRequirement", quantities, prices)
}

// CurrentRequirementOutput serves as a container for the return parameters of contract
// method CurrentRequirement.
type CurrentRequirementOutput struct {
	Margin  *big.Int
	Missing uint8
	Regime  uint8
}

// UnpackCurrentRequirement is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x995bff62.
//
// Solidity: function currentRequirement(int256[] quantities, uint256[] prices) view returns(uint256 margin, uint8 missing, uint8 regime)
func (margin *Margin) UnpackCurrentRequirement(data []byte) (CurrentRequirementOutput, error) {
	out, err := margin.abi.Unpack("currentRequirement", data)
	outstruct := new(CurrentRequirementOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Margin = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Missing = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Regime = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9e77a298.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depth(bytes32 symbol) view returns(uint32 selling, uint32 buying, uint32 sellingCeiling, uint32 buyingCeiling)
func (margin *Margin) PackDepth(symbol [32]byte) []byte {
	enc, err := margin.abi.Pack("depth", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9e77a298.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depth(bytes32 symbol) view returns(uint32 selling, uint32 buying, uint32 sellingCeiling, uint32 buyingCeiling)
func (margin *Margin) TryPackDepth(symbol [32]byte) ([]byte, error) {
	return margin.abi.Pack("depth", symbol)
}

// DepthOutput serves as a container for the return parameters of contract
// method Depth.
type DepthOutput struct {
	Selling        uint32
	Buying         uint32
	SellingCeiling uint32
	BuyingCeiling  uint32
}

// UnpackDepth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9e77a298.
//
// Solidity: function depth(bytes32 symbol) view returns(uint32 selling, uint32 buying, uint32 sellingCeiling, uint32 buyingCeiling)
func (margin *Margin) UnpackDepth(data []byte) (DepthOutput, error) {
	out, err := margin.abi.Unpack("depth", data)
	outstruct := new(DepthOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Selling = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Buying = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.SellingCeiling = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.BuyingCeiling = *abi.ConvertType(out[3], new(uint32)).(*uint32)
	return *outstruct, nil
}

// PackEthUsdFeed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8009b7bd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ethUsdFeed() view returns(address)
func (margin *Margin) PackEthUsdFeed() []byte {
	enc, err := margin.abi.Pack("ethUsdFeed")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEthUsdFeed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8009b7bd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ethUsdFeed() view returns(address)
func (margin *Margin) TryPackEthUsdFeed() ([]byte, error) {
	return margin.abi.Pack("ethUsdFeed")
}

// UnpackEthUsdFeed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8009b7bd.
//
// Solidity: function ethUsdFeed() view returns(address)
func (margin *Margin) UnpackEthUsdFeed(data []byte) (common.Address, error) {
	out, err := margin.abi.Unpack("ethUsdFeed", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackLastUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0463711.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastUpdate() view returns(uint64)
func (margin *Margin) PackLastUpdate() []byte {
	enc, err := margin.abi.Pack("lastUpdate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0463711.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastUpdate() view returns(uint64)
func (margin *Margin) TryPackLastUpdate() ([]byte, error) {
	return margin.abi.Pack("lastUpdate")
}

// UnpackLastUpdate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc0463711.
//
// Solidity: function lastUpdate() view returns(uint64)
func (margin *Margin) UnpackLastUpdate(data []byte) (uint64, error) {
	out, err := margin.abi.Unpack("lastUpdate", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80f55605.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function market() view returns(bytes32)
func (margin *Margin) PackMarket() []byte {
	enc, err := margin.abi.Pack("market")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x80f55605.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function market() view returns(bytes32)
func (margin *Margin) TryPackMarket() ([]byte, error) {
	return margin.abi.Pack("market")
}

// UnpackMarket is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x80f55605.
//
// Solidity: function market() view returns(bytes32)
func (margin *Margin) UnpackMarket(data []byte) ([32]byte, error) {
	out, err := margin.abi.Unpack("market", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (margin *Margin) PackOwner() []byte {
	enc, err := margin.abi.Pack("owner")
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
func (margin *Margin) TryPackOwner() ([]byte, error) {
	return margin.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (margin *Margin) UnpackOwner(data []byte) (common.Address, error) {
	out, err := margin.abi.Unpack("owner", data)
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
func (margin *Margin) PackPendingOwner() []byte {
	enc, err := margin.abi.Pack("pendingOwner")
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
func (margin *Margin) TryPackPendingOwner() ([]byte, error) {
	return margin.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (margin *Margin) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := margin.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b7a4480.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pool(bytes32 symbol) view returns(address)
func (margin *Margin) PackPool(symbol [32]byte) []byte {
	enc, err := margin.abi.Pack("pool", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b7a4480.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pool(bytes32 symbol) view returns(address)
func (margin *Margin) TryPackPool(symbol [32]byte) ([]byte, error) {
	return margin.abi.Pack("pool", symbol)
}

// UnpackPool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b7a4480.
//
// Solidity: function pool(bytes32 symbol) view returns(address)
func (margin *Margin) UnpackPool(data []byte) (common.Address, error) {
	out, err := margin.abi.Unpack("pool", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRequirement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb490ff16.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requirement(int256[] quantities, uint256[] prices, uint64 horizon, bool spansClosure) view returns(uint256 margin, uint8 missing)
func (margin *Margin) PackRequirement(quantities []*big.Int, prices []*big.Int, horizon uint64, spansClosure bool) []byte {
	enc, err := margin.abi.Pack("requirement", quantities, prices, horizon, spansClosure)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequirement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb490ff16.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requirement(int256[] quantities, uint256[] prices, uint64 horizon, bool spansClosure) view returns(uint256 margin, uint8 missing)
func (margin *Margin) TryPackRequirement(quantities []*big.Int, prices []*big.Int, horizon uint64, spansClosure bool) ([]byte, error) {
	return margin.abi.Pack("requirement", quantities, prices, horizon, spansClosure)
}

// RequirementOutput serves as a container for the return parameters of contract
// method Requirement.
type RequirementOutput struct {
	Margin  *big.Int
	Missing uint8
}

// UnpackRequirement is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb490ff16.
//
// Solidity: function requirement(int256[] quantities, uint256[] prices, uint64 horizon, bool spansClosure) view returns(uint256 margin, uint8 missing)
func (margin *Margin) UnpackRequirement(data []byte) (RequirementOutput, error) {
	out, err := margin.abi.Unpack("requirement", data)
	outstruct := new(RequirementOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Margin = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Missing = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackSetParameters is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x332403cd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setParameters(uint32[] volatilities, uint16[] correlations, uint32[] gaps, uint32[] depths) returns()
func (margin *Margin) PackSetParameters(volatilities []uint32, correlations []uint16, gaps []uint32, depths []uint32) []byte {
	enc, err := margin.abi.Pack("setParameters", volatilities, correlations, gaps, depths)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetParameters is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x332403cd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setParameters(uint32[] volatilities, uint16[] correlations, uint32[] gaps, uint32[] depths) returns()
func (margin *Margin) TryPackSetParameters(volatilities []uint32, correlations []uint16, gaps []uint32, depths []uint32) ([]byte, error) {
	return margin.abi.Pack("setParameters", volatilities, correlations, gaps, depths)
}

// PackVolatility is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f205fca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function volatility(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) PackVolatility(symbol [32]byte) []byte {
	enc, err := margin.abi.Pack("volatility", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVolatility is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6f205fca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function volatility(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) TryPackVolatility(symbol [32]byte) ([]byte, error) {
	return margin.abi.Pack("volatility", symbol)
}

// VolatilityOutput serves as a container for the return parameters of contract
// method Volatility.
type VolatilityOutput struct {
	Value uint32
	Floor uint32
}

// UnpackVolatility is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6f205fca.
//
// Solidity: function volatility(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) UnpackVolatility(data []byte) (VolatilityOutput, error) {
	out, err := margin.abi.Unpack("volatility", data)
	outstruct := new(VolatilityOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Floor = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	return *outstruct, nil
}

// PackWeekendGap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc6ec89f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weekendGap(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) PackWeekendGap(symbol [32]byte) []byte {
	enc, err := margin.abi.Pack("weekendGap", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeekendGap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc6ec89f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weekendGap(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) TryPackWeekendGap(symbol [32]byte) ([]byte, error) {
	return margin.abi.Pack("weekendGap", symbol)
}

// WeekendGapOutput serves as a container for the return parameters of contract
// method WeekendGap.
type WeekendGapOutput struct {
	Value uint32
	Floor uint32
}

// UnpackWeekendGap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdc6ec89f.
//
// Solidity: function weekendGap(bytes32 symbol) view returns(uint32 value, uint32 floor)
func (margin *Margin) UnpackWeekendGap(data []byte) (WeekendGapOutput, error) {
	out, err := margin.abi.Unpack("weekendGap", data)
	outstruct := new(WeekendGapOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Floor = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	return *outstruct, nil
}

// PackWeekendLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb659ad89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (margin *Margin) PackWeekendLeverage() []byte {
	enc, err := margin.abi.Pack("weekendLeverage")
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
func (margin *Margin) TryPackWeekendLeverage() ([]byte, error) {
	return margin.abi.Pack("weekendLeverage")
}

// UnpackWeekendLeverage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb659ad89.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (margin *Margin) UnpackWeekendLeverage(data []byte) (uint32, error) {
	out, err := margin.abi.Unpack("weekendLeverage", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// MarginCorrelationSet represents a CorrelationSet event raised by the Margin contract.
type MarginCorrelationSet struct {
	Symbol [32]byte
	Other  [32]byte
	Value  uint16
	Raw    *types.Log // Blockchain specific contextual infos
}

const MarginCorrelationSetEventName = "CorrelationSet"

// ContractEventName returns the user-defined event name.
func (MarginCorrelationSet) ContractEventName() string {
	return MarginCorrelationSetEventName
}

// UnpackCorrelationSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CorrelationSet(bytes32 indexed symbol, bytes32 indexed other, uint16 value)
func (margin *Margin) UnpackCorrelationSetEvent(log *types.Log) (*MarginCorrelationSet, error) {
	event := "CorrelationSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginCorrelationSet)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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

// MarginDepthSet represents a DepthSet event raised by the Margin contract.
type MarginDepthSet struct {
	Symbol  [32]byte
	Selling uint32
	Buying  uint32
	Raw     *types.Log // Blockchain specific contextual infos
}

const MarginDepthSetEventName = "DepthSet"

// ContractEventName returns the user-defined event name.
func (MarginDepthSet) ContractEventName() string {
	return MarginDepthSetEventName
}

// UnpackDepthSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DepthSet(bytes32 indexed symbol, uint32 selling, uint32 buying)
func (margin *Margin) UnpackDepthSetEvent(log *types.Log) (*MarginDepthSet, error) {
	event := "DepthSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginDepthSet)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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

// MarginGapSet represents a GapSet event raised by the Margin contract.
type MarginGapSet struct {
	Symbol [32]byte
	Value  uint32
	Raw    *types.Log // Blockchain specific contextual infos
}

const MarginGapSetEventName = "GapSet"

// ContractEventName returns the user-defined event name.
func (MarginGapSet) ContractEventName() string {
	return MarginGapSetEventName
}

// UnpackGapSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GapSet(bytes32 indexed symbol, uint32 value)
func (margin *Margin) UnpackGapSetEvent(log *types.Log) (*MarginGapSet, error) {
	event := "GapSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginGapSet)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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

// MarginOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Margin contract.
type MarginOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MarginOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (MarginOwnershipTransferStarted) ContractEventName() string {
	return MarginOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (margin *Margin) UnpackOwnershipTransferStartedEvent(log *types.Log) (*MarginOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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

// MarginOwnershipTransferred represents a OwnershipTransferred event raised by the Margin contract.
type MarginOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MarginOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MarginOwnershipTransferred) ContractEventName() string {
	return MarginOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (margin *Margin) UnpackOwnershipTransferredEvent(log *types.Log) (*MarginOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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

// MarginVolatilitySet represents a VolatilitySet event raised by the Margin contract.
type MarginVolatilitySet struct {
	Symbol [32]byte
	Value  uint32
	Raw    *types.Log // Blockchain specific contextual infos
}

const MarginVolatilitySetEventName = "VolatilitySet"

// ContractEventName returns the user-defined event name.
func (MarginVolatilitySet) ContractEventName() string {
	return MarginVolatilitySetEventName
}

// UnpackVolatilitySetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VolatilitySet(bytes32 indexed symbol, uint32 value)
func (margin *Margin) UnpackVolatilitySetEvent(log *types.Log) (*MarginVolatilitySet, error) {
	event := "VolatilitySet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != margin.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginVolatilitySet)
	if len(log.Data) > 0 {
		if err := margin.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range margin.abi.Events[event].Inputs {
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
func (margin *Margin) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], margin.abi.Errors["CorrelationStepTooLarge"].ID.Bytes()[:4]) {
		return margin.UnpackCorrelationStepTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["DepthStepTooLarge"].ID.Bytes()[:4]) {
		return margin.UnpackDepthStepTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["GapStepTooLarge"].ID.Bytes()[:4]) {
		return margin.UnpackGapStepTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["InvalidCorrelation"].ID.Bytes()[:4]) {
		return margin.UnpackInvalidCorrelationError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["InvalidDepth"].ID.Bytes()[:4]) {
		return margin.UnpackInvalidDepthError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["InvalidGap"].ID.Bytes()[:4]) {
		return margin.UnpackInvalidGapError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["InvalidVolatility"].ID.Bytes()[:4]) {
		return margin.UnpackInvalidVolatilityError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["LengthMismatch"].ID.Bytes()[:4]) {
		return margin.UnpackLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["NotPositiveDefinite"].ID.Bytes()[:4]) {
		return margin.UnpackNotPositiveDefiniteError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return margin.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["UnknownAsset"].ID.Bytes()[:4]) {
		return margin.UnpackUnknownAssetError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["UpdateTooSoon"].ID.Bytes()[:4]) {
		return margin.UnpackUpdateTooSoonError(raw[4:])
	}
	if bytes.Equal(raw[:4], margin.abi.Errors["VolatilityStepTooLarge"].ID.Bytes()[:4]) {
		return margin.UnpackVolatilityStepTooLargeError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MarginCorrelationStepTooLarge represents a CorrelationStepTooLarge error raised by the Margin contract.
type MarginCorrelationStepTooLarge struct {
	Symbol   [32]byte
	Other    [32]byte
	Previous uint16
	Value    uint16
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CorrelationStepTooLarge(bytes32 symbol, bytes32 other, uint16 previous, uint16 value)
func MarginCorrelationStepTooLargeErrorID() common.Hash {
	return common.HexToHash("0xa94553bf2d20232e14257f0357cb3fabd1dc8058595cd5880d85d26025600959")
}

// UnpackCorrelationStepTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CorrelationStepTooLarge(bytes32 symbol, bytes32 other, uint16 previous, uint16 value)
func (margin *Margin) UnpackCorrelationStepTooLargeError(raw []byte) (*MarginCorrelationStepTooLarge, error) {
	out := new(MarginCorrelationStepTooLarge)
	if err := margin.abi.UnpackIntoInterface(out, "CorrelationStepTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginDepthStepTooLarge represents a DepthStepTooLarge error raised by the Margin contract.
type MarginDepthStepTooLarge struct {
	Symbol   [32]byte
	Previous uint32
	Value    uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepthStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func MarginDepthStepTooLargeErrorID() common.Hash {
	return common.HexToHash("0x393e5a737a0b014eb025529ebf646c70ab68ff951ad87b380281464500979553")
}

// UnpackDepthStepTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepthStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func (margin *Margin) UnpackDepthStepTooLargeError(raw []byte) (*MarginDepthStepTooLarge, error) {
	out := new(MarginDepthStepTooLarge)
	if err := margin.abi.UnpackIntoInterface(out, "DepthStepTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginGapStepTooLarge represents a GapStepTooLarge error raised by the Margin contract.
type MarginGapStepTooLarge struct {
	Symbol   [32]byte
	Previous uint32
	Value    uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error GapStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func MarginGapStepTooLargeErrorID() common.Hash {
	return common.HexToHash("0x4025bda0577ba1c10b1e9e8cda3c5932c41be72ce75524ed44db7e3754d9c338")
}

// UnpackGapStepTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error GapStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func (margin *Margin) UnpackGapStepTooLargeError(raw []byte) (*MarginGapStepTooLarge, error) {
	out := new(MarginGapStepTooLarge)
	if err := margin.abi.UnpackIntoInterface(out, "GapStepTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginInvalidCorrelation represents a InvalidCorrelation error raised by the Margin contract.
type MarginInvalidCorrelation struct {
	Symbol [32]byte
	Other  [32]byte
	Value  uint16
	Floor  uint16
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidCorrelation(bytes32 symbol, bytes32 other, uint16 value, uint16 floor)
func MarginInvalidCorrelationErrorID() common.Hash {
	return common.HexToHash("0x2793e7a38b41f4e235ba876e439ef42912bd52d6535de25a746e674d607a7f13")
}

// UnpackInvalidCorrelationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidCorrelation(bytes32 symbol, bytes32 other, uint16 value, uint16 floor)
func (margin *Margin) UnpackInvalidCorrelationError(raw []byte) (*MarginInvalidCorrelation, error) {
	out := new(MarginInvalidCorrelation)
	if err := margin.abi.UnpackIntoInterface(out, "InvalidCorrelation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginInvalidDepth represents a InvalidDepth error raised by the Margin contract.
type MarginInvalidDepth struct {
	Symbol  [32]byte
	Value   uint32
	Ceiling uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidDepth(bytes32 symbol, uint32 value, uint32 ceiling)
func MarginInvalidDepthErrorID() common.Hash {
	return common.HexToHash("0xb0d313fafa07d738d8c40f14fc13d6cacd381ae9c6f44fe13eb3c4f7a9d50798")
}

// UnpackInvalidDepthError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidDepth(bytes32 symbol, uint32 value, uint32 ceiling)
func (margin *Margin) UnpackInvalidDepthError(raw []byte) (*MarginInvalidDepth, error) {
	out := new(MarginInvalidDepth)
	if err := margin.abi.UnpackIntoInterface(out, "InvalidDepth", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginInvalidGap represents a InvalidGap error raised by the Margin contract.
type MarginInvalidGap struct {
	Symbol [32]byte
	Value  uint32
	Floor  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidGap(bytes32 symbol, uint32 value, uint32 floor)
func MarginInvalidGapErrorID() common.Hash {
	return common.HexToHash("0x3f8f34a39014bc849f003199531c41c98797c005c12dea027575ae15a68797e8")
}

// UnpackInvalidGapError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidGap(bytes32 symbol, uint32 value, uint32 floor)
func (margin *Margin) UnpackInvalidGapError(raw []byte) (*MarginInvalidGap, error) {
	out := new(MarginInvalidGap)
	if err := margin.abi.UnpackIntoInterface(out, "InvalidGap", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginInvalidVolatility represents a InvalidVolatility error raised by the Margin contract.
type MarginInvalidVolatility struct {
	Symbol [32]byte
	Value  uint32
	Floor  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidVolatility(bytes32 symbol, uint32 value, uint32 floor)
func MarginInvalidVolatilityErrorID() common.Hash {
	return common.HexToHash("0xd271f5a0fc3a046c1c6ddb3fab0f263b153ad9d4a5ec5f25acb3bfeaa67e3769")
}

// UnpackInvalidVolatilityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidVolatility(bytes32 symbol, uint32 value, uint32 floor)
func (margin *Margin) UnpackInvalidVolatilityError(raw []byte) (*MarginInvalidVolatility, error) {
	out := new(MarginInvalidVolatility)
	if err := margin.abi.UnpackIntoInterface(out, "InvalidVolatility", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginLengthMismatch represents a LengthMismatch error raised by the Margin contract.
type MarginLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LengthMismatch()
func MarginLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xff633a3803c58b9bc21e58efecee59f27e033cc0b1883fccb4969c76146fe60f")
}

// UnpackLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LengthMismatch()
func (margin *Margin) UnpackLengthMismatchError(raw []byte) (*MarginLengthMismatch, error) {
	out := new(MarginLengthMismatch)
	if err := margin.abi.UnpackIntoInterface(out, "LengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginNotPositiveDefinite represents a NotPositiveDefinite error raised by the Margin contract.
type MarginNotPositiveDefinite struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotPositiveDefinite()
func MarginNotPositiveDefiniteErrorID() common.Hash {
	return common.HexToHash("0x2637f6156652a16e5c1b9cbfa2476bf319413deb18f6064f6d972e9bb8f8c163")
}

// UnpackNotPositiveDefiniteError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotPositiveDefinite()
func (margin *Margin) UnpackNotPositiveDefiniteError(raw []byte) (*MarginNotPositiveDefinite, error) {
	out := new(MarginNotPositiveDefinite)
	if err := margin.abi.UnpackIntoInterface(out, "NotPositiveDefinite", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the Margin contract.
type MarginOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MarginOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (margin *Margin) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MarginOwnableUnauthorizedAccount, error) {
	out := new(MarginOwnableUnauthorizedAccount)
	if err := margin.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginUnknownAsset represents a UnknownAsset error raised by the Margin contract.
type MarginUnknownAsset struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func MarginUnknownAssetErrorID() common.Hash {
	return common.HexToHash("0x1059be3ed9f1c1f061f09edf206f56e7c80b74fd01671cbec1b3788a1f49417e")
}

// UnpackUnknownAssetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func (margin *Margin) UnpackUnknownAssetError(raw []byte) (*MarginUnknownAsset, error) {
	out := new(MarginUnknownAsset)
	if err := margin.abi.UnpackIntoInterface(out, "UnknownAsset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginUpdateTooSoon represents a UpdateTooSoon error raised by the Margin contract.
type MarginUpdateTooSoon struct {
	NextUpdateAt uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UpdateTooSoon(uint64 nextUpdateAt)
func MarginUpdateTooSoonErrorID() common.Hash {
	return common.HexToHash("0x9ec53c104f5bae60d4b964d7b3a05382856feed221816cc69bf052042f8a27bf")
}

// UnpackUpdateTooSoonError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UpdateTooSoon(uint64 nextUpdateAt)
func (margin *Margin) UnpackUpdateTooSoonError(raw []byte) (*MarginUpdateTooSoon, error) {
	out := new(MarginUpdateTooSoon)
	if err := margin.abi.UnpackIntoInterface(out, "UpdateTooSoon", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginVolatilityStepTooLarge represents a VolatilityStepTooLarge error raised by the Margin contract.
type MarginVolatilityStepTooLarge struct {
	Symbol   [32]byte
	Previous uint32
	Value    uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VolatilityStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func MarginVolatilityStepTooLargeErrorID() common.Hash {
	return common.HexToHash("0x7b13cffaf373ffa0aa2e300be8863abfb4ab9f6f2ea904e8087766f866c32147")
}

// UnpackVolatilityStepTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VolatilityStepTooLarge(bytes32 symbol, uint32 previous, uint32 value)
func (margin *Margin) UnpackVolatilityStepTooLargeError(raw []byte) (*MarginVolatilityStepTooLarge, error) {
	out := new(MarginVolatilityStepTooLarge)
	if err := margin.abi.UnpackIntoInterface(out, "VolatilityStepTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}
