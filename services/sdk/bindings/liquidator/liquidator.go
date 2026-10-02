// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package liquidator

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

// LiquidatorMetaData contains all meta data concerning the Liquidator contract.
var LiquidatorMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"accounts_\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CLOSED_AUCTION_LIFETIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CLOSED_DECAY_BPS_PER_MINUTE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CLOSED_HOURLY_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CLOSE_FACTOR_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DUST\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FULL_CLOSE_OWED\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"HORIZON\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_DISCOUNT_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_FEED_AGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OPEN_AUCTION_LIFETIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OPEN_DECAY_BPS_PER_MINUTE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WETH_LIQUIDATION_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"auction\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"auctions\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"startedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"closed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buy\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"bought\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collect\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"engine\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMargin\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethUsd\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"heldUntil\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hold\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"until\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hourlyAllowance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"left\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"price\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recall\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recallHaircut\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAuction\",\"inputs\":[{\"name\":\"newAuction\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unitPrice\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"sold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleCash\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"settled\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"shortfall\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"short\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"closed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"start\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stop\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usdg\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weth\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"writeOff\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AuctionSet\",\"inputs\":[{\"name\":\"auction\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuctionStarted\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closed\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuctionStopped\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Bought\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CashSettled\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Held\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"until\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Settled\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cost\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AuctionAlreadySet\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CannotJudge\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CostAboveLimit\",\"inputs\":[{\"name\":\"cost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAuction\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoAuction\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotAccountsOwner\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotAuction\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotBackstop\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotForSale\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotLiquidatable\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NothingToBuy\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToRecall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionHeld\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"until\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedIntToUint\",\"inputs\":[{\"name\":\"value\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintToInt\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StillShort\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "Liquidator",
}

// Liquidator is an auto generated Go binding around an Ethereum contract.
type Liquidator struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Liquidator) GetABI() abi.ABI {
	return c.abi
}

// NewLiquidator creates a new instance of Liquidator.
func NewLiquidator() *Liquidator {
	parsed, err := LiquidatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Liquidator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Liquidator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address accounts_) returns()
func (liquidator *Liquidator) PackConstructor(accounts_ common.Address) []byte {
	enc, err := liquidator.abi.Pack("", accounts_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCLOSEDAUCTIONLIFETIME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c827f3c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CLOSED_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) PackCLOSEDAUCTIONLIFETIME() []byte {
	enc, err := liquidator.abi.Pack("CLOSED_AUCTION_LIFETIME")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCLOSEDAUCTIONLIFETIME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c827f3c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CLOSED_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) TryPackCLOSEDAUCTIONLIFETIME() ([]byte, error) {
	return liquidator.abi.Pack("CLOSED_AUCTION_LIFETIME")
}

// UnpackCLOSEDAUCTIONLIFETIME is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c827f3c.
//
// Solidity: function CLOSED_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) UnpackCLOSEDAUCTIONLIFETIME(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("CLOSED_AUCTION_LIFETIME", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCLOSEDDECAYBPSPERMINUTE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0267c06b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CLOSED_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) PackCLOSEDDECAYBPSPERMINUTE() []byte {
	enc, err := liquidator.abi.Pack("CLOSED_DECAY_BPS_PER_MINUTE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCLOSEDDECAYBPSPERMINUTE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0267c06b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CLOSED_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) TryPackCLOSEDDECAYBPSPERMINUTE() ([]byte, error) {
	return liquidator.abi.Pack("CLOSED_DECAY_BPS_PER_MINUTE")
}

// UnpackCLOSEDDECAYBPSPERMINUTE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0267c06b.
//
// Solidity: function CLOSED_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) UnpackCLOSEDDECAYBPSPERMINUTE(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("CLOSED_DECAY_BPS_PER_MINUTE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCLOSEDHOURLYBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21cfccc0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CLOSED_HOURLY_BPS() view returns(uint256)
func (liquidator *Liquidator) PackCLOSEDHOURLYBPS() []byte {
	enc, err := liquidator.abi.Pack("CLOSED_HOURLY_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCLOSEDHOURLYBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21cfccc0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CLOSED_HOURLY_BPS() view returns(uint256)
func (liquidator *Liquidator) TryPackCLOSEDHOURLYBPS() ([]byte, error) {
	return liquidator.abi.Pack("CLOSED_HOURLY_BPS")
}

// UnpackCLOSEDHOURLYBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x21cfccc0.
//
// Solidity: function CLOSED_HOURLY_BPS() view returns(uint256)
func (liquidator *Liquidator) UnpackCLOSEDHOURLYBPS(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("CLOSED_HOURLY_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCLOSEFACTORBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48ddab64.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CLOSE_FACTOR_BPS() view returns(uint256)
func (liquidator *Liquidator) PackCLOSEFACTORBPS() []byte {
	enc, err := liquidator.abi.Pack("CLOSE_FACTOR_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCLOSEFACTORBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48ddab64.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CLOSE_FACTOR_BPS() view returns(uint256)
func (liquidator *Liquidator) TryPackCLOSEFACTORBPS() ([]byte, error) {
	return liquidator.abi.Pack("CLOSE_FACTOR_BPS")
}

// UnpackCLOSEFACTORBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x48ddab64.
//
// Solidity: function CLOSE_FACTOR_BPS() view returns(uint256)
func (liquidator *Liquidator) UnpackCLOSEFACTORBPS(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("CLOSE_FACTOR_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDUST is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e0cd799.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DUST() view returns(uint256)
func (liquidator *Liquidator) PackDUST() []byte {
	enc, err := liquidator.abi.Pack("DUST")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDUST is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e0cd799.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DUST() view returns(uint256)
func (liquidator *Liquidator) TryPackDUST() ([]byte, error) {
	return liquidator.abi.Pack("DUST")
}

// UnpackDUST is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4e0cd799.
//
// Solidity: function DUST() view returns(uint256)
func (liquidator *Liquidator) UnpackDUST(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("DUST", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFEEBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf333f2c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (liquidator *Liquidator) PackFEEBPS() []byte {
	enc, err := liquidator.abi.Pack("FEE_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFEEBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf333f2c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (liquidator *Liquidator) TryPackFEEBPS() ([]byte, error) {
	return liquidator.abi.Pack("FEE_BPS")
}

// UnpackFEEBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint256)
func (liquidator *Liquidator) UnpackFEEBPS(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("FEE_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFULLCLOSEOWED is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfc12838.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FULL_CLOSE_OWED() view returns(uint256)
func (liquidator *Liquidator) PackFULLCLOSEOWED() []byte {
	enc, err := liquidator.abi.Pack("FULL_CLOSE_OWED")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFULLCLOSEOWED is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfc12838.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FULL_CLOSE_OWED() view returns(uint256)
func (liquidator *Liquidator) TryPackFULLCLOSEOWED() ([]byte, error) {
	return liquidator.abi.Pack("FULL_CLOSE_OWED")
}

// UnpackFULLCLOSEOWED is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbfc12838.
//
// Solidity: function FULL_CLOSE_OWED() view returns(uint256)
func (liquidator *Liquidator) UnpackFULLCLOSEOWED(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("FULL_CLOSE_OWED", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHORIZON is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb444ef9d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function HORIZON() view returns(uint64)
func (liquidator *Liquidator) PackHORIZON() []byte {
	enc, err := liquidator.abi.Pack("HORIZON")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHORIZON is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb444ef9d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function HORIZON() view returns(uint64)
func (liquidator *Liquidator) TryPackHORIZON() ([]byte, error) {
	return liquidator.abi.Pack("HORIZON")
}

// UnpackHORIZON is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb444ef9d.
//
// Solidity: function HORIZON() view returns(uint64)
func (liquidator *Liquidator) UnpackHORIZON(data []byte) (uint64, error) {
	out, err := liquidator.abi.Unpack("HORIZON", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMAXDISCOUNTBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x390947cf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_DISCOUNT_BPS() view returns(uint256)
func (liquidator *Liquidator) PackMAXDISCOUNTBPS() []byte {
	enc, err := liquidator.abi.Pack("MAX_DISCOUNT_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXDISCOUNTBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x390947cf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_DISCOUNT_BPS() view returns(uint256)
func (liquidator *Liquidator) TryPackMAXDISCOUNTBPS() ([]byte, error) {
	return liquidator.abi.Pack("MAX_DISCOUNT_BPS")
}

// UnpackMAXDISCOUNTBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x390947cf.
//
// Solidity: function MAX_DISCOUNT_BPS() view returns(uint256)
func (liquidator *Liquidator) UnpackMAXDISCOUNTBPS(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("MAX_DISCOUNT_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXFEEDAGE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x696453b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (liquidator *Liquidator) PackMAXFEEDAGE() []byte {
	enc, err := liquidator.abi.Pack("MAX_FEED_AGE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXFEEDAGE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x696453b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (liquidator *Liquidator) TryPackMAXFEEDAGE() ([]byte, error) {
	return liquidator.abi.Pack("MAX_FEED_AGE")
}

// UnpackMAXFEEDAGE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x696453b8.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (liquidator *Liquidator) UnpackMAXFEEDAGE(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("MAX_FEED_AGE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOPENAUCTIONLIFETIME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2d6318d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function OPEN_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) PackOPENAUCTIONLIFETIME() []byte {
	enc, err := liquidator.abi.Pack("OPEN_AUCTION_LIFETIME")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOPENAUCTIONLIFETIME is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2d6318d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function OPEN_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) TryPackOPENAUCTIONLIFETIME() ([]byte, error) {
	return liquidator.abi.Pack("OPEN_AUCTION_LIFETIME")
}

// UnpackOPENAUCTIONLIFETIME is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe2d6318d.
//
// Solidity: function OPEN_AUCTION_LIFETIME() view returns(uint256)
func (liquidator *Liquidator) UnpackOPENAUCTIONLIFETIME(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("OPEN_AUCTION_LIFETIME", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOPENDECAYBPSPERMINUTE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd664e2ec.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function OPEN_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) PackOPENDECAYBPSPERMINUTE() []byte {
	enc, err := liquidator.abi.Pack("OPEN_DECAY_BPS_PER_MINUTE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOPENDECAYBPSPERMINUTE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd664e2ec.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function OPEN_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) TryPackOPENDECAYBPSPERMINUTE() ([]byte, error) {
	return liquidator.abi.Pack("OPEN_DECAY_BPS_PER_MINUTE")
}

// UnpackOPENDECAYBPSPERMINUTE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd664e2ec.
//
// Solidity: function OPEN_DECAY_BPS_PER_MINUTE() view returns(uint256)
func (liquidator *Liquidator) UnpackOPENDECAYBPSPERMINUTE(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("OPEN_DECAY_BPS_PER_MINUTE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWETHLIQUIDATIONBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d7dca41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function WETH_LIQUIDATION_BPS() view returns(uint256)
func (liquidator *Liquidator) PackWETHLIQUIDATIONBPS() []byte {
	enc, err := liquidator.abi.Pack("WETH_LIQUIDATION_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWETHLIQUIDATIONBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d7dca41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function WETH_LIQUIDATION_BPS() view returns(uint256)
func (liquidator *Liquidator) TryPackWETHLIQUIDATIONBPS() ([]byte, error) {
	return liquidator.abi.Pack("WETH_LIQUIDATION_BPS")
}

// UnpackWETHLIQUIDATIONBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0d7dca41.
//
// Solidity: function WETH_LIQUIDATION_BPS() view returns(uint256)
func (liquidator *Liquidator) UnpackWETHLIQUIDATIONBPS(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("WETH_LIQUIDATION_BPS", data)
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
func (liquidator *Liquidator) PackAccounts() []byte {
	enc, err := liquidator.abi.Pack("accounts")
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
func (liquidator *Liquidator) TryPackAccounts() ([]byte, error) {
	return liquidator.abi.Pack("accounts")
}

// UnpackAccounts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68cd03f6.
//
// Solidity: function accounts() view returns(address)
func (liquidator *Liquidator) UnpackAccounts(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("accounts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAuction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d9f6db5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function auction() view returns(address)
func (liquidator *Liquidator) PackAuction() []byte {
	enc, err := liquidator.abi.Pack("auction")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAuction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d9f6db5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function auction() view returns(address)
func (liquidator *Liquidator) TryPackAuction() ([]byte, error) {
	return liquidator.abi.Pack("auction")
}

// UnpackAuction is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d9f6db5.
//
// Solidity: function auction() view returns(address)
func (liquidator *Liquidator) UnpackAuction(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("auction", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAuctions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xee5bcb62.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function auctions(address account, bytes32 position) view returns(uint64 startedAt, bool closed)
func (liquidator *Liquidator) PackAuctions(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("auctions", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAuctions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xee5bcb62.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function auctions(address account, bytes32 position) view returns(uint64 startedAt, bool closed)
func (liquidator *Liquidator) TryPackAuctions(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("auctions", account, position)
}

// AuctionsOutput serves as a container for the return parameters of contract
// method Auctions.
type AuctionsOutput struct {
	StartedAt uint64
	Closed    bool
}

// UnpackAuctions is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xee5bcb62.
//
// Solidity: function auctions(address account, bytes32 position) view returns(uint64 startedAt, bool closed)
func (liquidator *Liquidator) UnpackAuctions(data []byte) (AuctionsOutput, error) {
	out, err := liquidator.abi.Unpack("auctions", data)
	outstruct := new(AuctionsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.StartedAt = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Closed = *abi.ConvertType(out[1], new(bool)).(*bool)
	return *outstruct, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (liquidator *Liquidator) PackBand() []byte {
	enc, err := liquidator.abi.Pack("band")
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
func (liquidator *Liquidator) TryPackBand() ([]byte, error) {
	return liquidator.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (liquidator *Liquidator) UnpackBand(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBuy is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57520fbc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buy(address account, bytes32 position, address token, uint256 amount, uint256 maxCost, address receiver) returns(uint256 bought, uint256 cost)
func (liquidator *Liquidator) PackBuy(account common.Address, position [32]byte, token common.Address, amount *big.Int, maxCost *big.Int, receiver common.Address) []byte {
	enc, err := liquidator.abi.Pack("buy", account, position, token, amount, maxCost, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuy is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57520fbc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buy(address account, bytes32 position, address token, uint256 amount, uint256 maxCost, address receiver) returns(uint256 bought, uint256 cost)
func (liquidator *Liquidator) TryPackBuy(account common.Address, position [32]byte, token common.Address, amount *big.Int, maxCost *big.Int, receiver common.Address) ([]byte, error) {
	return liquidator.abi.Pack("buy", account, position, token, amount, maxCost, receiver)
}

// BuyOutput serves as a container for the return parameters of contract
// method Buy.
type BuyOutput struct {
	Bought *big.Int
	Cost   *big.Int
}

// UnpackBuy is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x57520fbc.
//
// Solidity: function buy(address account, bytes32 position, address token, uint256 amount, uint256 maxCost, address receiver) returns(uint256 bought, uint256 cost)
func (liquidator *Liquidator) UnpackBuy(data []byte) (BuyOutput, error) {
	out, err := liquidator.abi.Unpack("buy", data)
	outstruct := new(BuyOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Bought = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Cost = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCollect is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce3f865f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collect(uint256 amount) returns()
func (liquidator *Liquidator) PackCollect(amount *big.Int) []byte {
	enc, err := liquidator.abi.Pack("collect", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollect is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce3f865f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collect(uint256 amount) returns()
func (liquidator *Liquidator) TryPackCollect(amount *big.Int) ([]byte, error) {
	return liquidator.abi.Pack("collect", amount)
}

// PackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function engine() view returns(address)
func (liquidator *Liquidator) PackEngine() []byte {
	enc, err := liquidator.abi.Pack("engine")
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
func (liquidator *Liquidator) TryPackEngine() ([]byte, error) {
	return liquidator.abi.Pack("engine")
}

// UnpackEngine is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (liquidator *Liquidator) UnpackEngine(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("engine", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEthUsd is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a960216.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ethUsd() view returns(address)
func (liquidator *Liquidator) PackEthUsd() []byte {
	enc, err := liquidator.abi.Pack("ethUsd")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEthUsd is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a960216.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ethUsd() view returns(address)
func (liquidator *Liquidator) TryPackEthUsd() ([]byte, error) {
	return liquidator.abi.Pack("ethUsd")
}

// UnpackEthUsd is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5a960216.
//
// Solidity: function ethUsd() view returns(address)
func (liquidator *Liquidator) UnpackEthUsd(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("ethUsd", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackHeldUntil is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf46a3e69.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function heldUntil(address account, bytes32 position) view returns(uint64)
func (liquidator *Liquidator) PackHeldUntil(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("heldUntil", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHeldUntil is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf46a3e69.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function heldUntil(address account, bytes32 position) view returns(uint64)
func (liquidator *Liquidator) TryPackHeldUntil(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("heldUntil", account, position)
}

// UnpackHeldUntil is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf46a3e69.
//
// Solidity: function heldUntil(address account, bytes32 position) view returns(uint64)
func (liquidator *Liquidator) UnpackHeldUntil(data []byte) (uint64, error) {
	out, err := liquidator.abi.Unpack("heldUntil", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackHold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f6b365b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hold(address account, bytes32 position, uint64 until) returns()
func (liquidator *Liquidator) PackHold(account common.Address, position [32]byte, until uint64) []byte {
	enc, err := liquidator.abi.Pack("hold", account, position, until)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f6b365b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hold(address account, bytes32 position, uint64 until) returns()
func (liquidator *Liquidator) TryPackHold(account common.Address, position [32]byte, until uint64) ([]byte, error) {
	return liquidator.abi.Pack("hold", account, position, until)
}

// PackHourlyAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04d6ef61.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hourlyAllowance(address account, bytes32 position, address token) view returns(uint256 left)
func (liquidator *Liquidator) PackHourlyAllowance(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := liquidator.abi.Pack("hourlyAllowance", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHourlyAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04d6ef61.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hourlyAllowance(address account, bytes32 position, address token) view returns(uint256 left)
func (liquidator *Liquidator) TryPackHourlyAllowance(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return liquidator.abi.Pack("hourlyAllowance", account, position, token)
}

// UnpackHourlyAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x04d6ef61.
//
// Solidity: function hourlyAllowance(address account, bytes32 position, address token) view returns(uint256 left)
func (liquidator *Liquidator) UnpackHourlyAllowance(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("hourlyAllowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b3349ae.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function price(address account, bytes32 position, address token) view returns(uint256)
func (liquidator *Liquidator) PackPrice(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := liquidator.abi.Pack("price", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b3349ae.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function price(address account, bytes32 position, address token) view returns(uint256)
func (liquidator *Liquidator) TryPackPrice(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return liquidator.abi.Pack("price", account, position, token)
}

// UnpackPrice is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b3349ae.
//
// Solidity: function price(address account, bytes32 position, address token) view returns(uint256)
func (liquidator *Liquidator) UnpackPrice(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("price", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f5737a0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recall(address account, bytes32 position, address token) returns(uint256 amount)
func (liquidator *Liquidator) PackRecall(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := liquidator.abi.Pack("recall", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f5737a0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recall(address account, bytes32 position, address token) returns(uint256 amount)
func (liquidator *Liquidator) TryPackRecall(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return liquidator.abi.Pack("recall", account, position, token)
}

// UnpackRecall is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4f5737a0.
//
// Solidity: function recall(address account, bytes32 position, address token) returns(uint256 amount)
func (liquidator *Liquidator) UnpackRecall(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("recall", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRecallHaircut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x057e81a1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recallHaircut() view returns(uint256)
func (liquidator *Liquidator) PackRecallHaircut() []byte {
	enc, err := liquidator.abi.Pack("recallHaircut")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecallHaircut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x057e81a1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recallHaircut() view returns(uint256)
func (liquidator *Liquidator) TryPackRecallHaircut() ([]byte, error) {
	return liquidator.abi.Pack("recallHaircut")
}

// UnpackRecallHaircut is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x057e81a1.
//
// Solidity: function recallHaircut() view returns(uint256)
func (liquidator *Liquidator) UnpackRecallHaircut(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("recallHaircut", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetAuction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8c6f579.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAuction(address newAuction) returns()
func (liquidator *Liquidator) PackSetAuction(newAuction common.Address) []byte {
	enc, err := liquidator.abi.Pack("setAuction", newAuction)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAuction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8c6f579.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAuction(address newAuction) returns()
func (liquidator *Liquidator) TryPackSetAuction(newAuction common.Address) ([]byte, error) {
	return liquidator.abi.Pack("setAuction", newAuction)
}

// PackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6fcb2e6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function settle(address account, bytes32 position, address token, uint256 amount, uint256 unitPrice) returns(uint256 sold, uint256 cost)
func (liquidator *Liquidator) PackSettle(account common.Address, position [32]byte, token common.Address, amount *big.Int, unitPrice *big.Int) []byte {
	enc, err := liquidator.abi.Pack("settle", account, position, token, amount, unitPrice)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6fcb2e6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function settle(address account, bytes32 position, address token, uint256 amount, uint256 unitPrice) returns(uint256 sold, uint256 cost)
func (liquidator *Liquidator) TryPackSettle(account common.Address, position [32]byte, token common.Address, amount *big.Int, unitPrice *big.Int) ([]byte, error) {
	return liquidator.abi.Pack("settle", account, position, token, amount, unitPrice)
}

// SettleOutput serves as a container for the return parameters of contract
// method Settle.
type SettleOutput struct {
	Sold *big.Int
	Cost *big.Int
}

// UnpackSettle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe6fcb2e6.
//
// Solidity: function settle(address account, bytes32 position, address token, uint256 amount, uint256 unitPrice) returns(uint256 sold, uint256 cost)
func (liquidator *Liquidator) UnpackSettle(data []byte) (SettleOutput, error) {
	out, err := liquidator.abi.Unpack("settle", data)
	outstruct := new(SettleOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Sold = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Cost = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackSettleCash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x990e9055.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function settleCash(address account, bytes32 position) returns(uint256 settled)
func (liquidator *Liquidator) PackSettleCash(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("settleCash", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSettleCash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x990e9055.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function settleCash(address account, bytes32 position) returns(uint256 settled)
func (liquidator *Liquidator) TryPackSettleCash(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("settleCash", account, position)
}

// UnpackSettleCash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x990e9055.
//
// Solidity: function settleCash(address account, bytes32 position) returns(uint256 settled)
func (liquidator *Liquidator) UnpackSettleCash(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("settleCash", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackShortfall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbe602b74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function shortfall(address account, bytes32 position) view returns(int256 equity, uint256 requirement, bool short, bool closed)
func (liquidator *Liquidator) PackShortfall(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("shortfall", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackShortfall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbe602b74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function shortfall(address account, bytes32 position) view returns(int256 equity, uint256 requirement, bool short, bool closed)
func (liquidator *Liquidator) TryPackShortfall(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("shortfall", account, position)
}

// ShortfallOutput serves as a container for the return parameters of contract
// method Shortfall.
type ShortfallOutput struct {
	Equity      *big.Int
	Requirement *big.Int
	Short       bool
	Closed      bool
}

// UnpackShortfall is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbe602b74.
//
// Solidity: function shortfall(address account, bytes32 position) view returns(int256 equity, uint256 requirement, bool short, bool closed)
func (liquidator *Liquidator) UnpackShortfall(data []byte) (ShortfallOutput, error) {
	out, err := liquidator.abi.Unpack("shortfall", data)
	outstruct := new(ShortfallOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Equity = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Requirement = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Short = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.Closed = *abi.ConvertType(out[3], new(bool)).(*bool)
	return *outstruct, nil
}

// PackStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa88f0663.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function start(address account, bytes32 position) returns()
func (liquidator *Liquidator) PackStart(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("start", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa88f0663.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function start(address account, bytes32 position) returns()
func (liquidator *Liquidator) TryPackStart(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("start", account, position)
}

// PackStop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4c05535c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stop(address account, bytes32 position) returns()
func (liquidator *Liquidator) PackStop(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("stop", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4c05535c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stop(address account, bytes32 position) returns()
func (liquidator *Liquidator) TryPackStop(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("stop", account, position)
}

// PackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function usdg() view returns(address)
func (liquidator *Liquidator) PackUsdg() []byte {
	enc, err := liquidator.abi.Pack("usdg")
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
func (liquidator *Liquidator) TryPackUsdg() ([]byte, error) {
	return liquidator.abi.Pack("usdg")
}

// UnpackUsdg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5b91b7b.
//
// Solidity: function usdg() view returns(address)
func (liquidator *Liquidator) UnpackUsdg(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("usdg", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWeth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fc8cef3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weth() view returns(address)
func (liquidator *Liquidator) PackWeth() []byte {
	enc, err := liquidator.abi.Pack("weth")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fc8cef3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weth() view returns(address)
func (liquidator *Liquidator) TryPackWeth() ([]byte, error) {
	return liquidator.abi.Pack("weth")
}

// UnpackWeth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (liquidator *Liquidator) UnpackWeth(data []byte) (common.Address, error) {
	out, err := liquidator.abi.Unpack("weth", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f3596c0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256)
func (liquidator *Liquidator) PackWriteOff(account common.Address, position [32]byte) []byte {
	enc, err := liquidator.abi.Pack("writeOff", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f3596c0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256)
func (liquidator *Liquidator) TryPackWriteOff(account common.Address, position [32]byte) ([]byte, error) {
	return liquidator.abi.Pack("writeOff", account, position)
}

// UnpackWriteOff is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1f3596c0.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256)
func (liquidator *Liquidator) UnpackWriteOff(data []byte) (*big.Int, error) {
	out, err := liquidator.abi.Unpack("writeOff", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// LiquidatorAuctionSet represents a AuctionSet event raised by the Liquidator contract.
type LiquidatorAuctionSet struct {
	Auction common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const LiquidatorAuctionSetEventName = "AuctionSet"

// ContractEventName returns the user-defined event name.
func (LiquidatorAuctionSet) ContractEventName() string {
	return LiquidatorAuctionSetEventName
}

// UnpackAuctionSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuctionSet(address indexed auction)
func (liquidator *Liquidator) UnpackAuctionSetEvent(log *types.Log) (*LiquidatorAuctionSet, error) {
	event := "AuctionSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorAuctionSet)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorAuctionStarted represents a AuctionStarted event raised by the Liquidator contract.
type LiquidatorAuctionStarted struct {
	Account  common.Address
	Position [32]byte
	Closed   bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorAuctionStartedEventName = "AuctionStarted"

// ContractEventName returns the user-defined event name.
func (LiquidatorAuctionStarted) ContractEventName() string {
	return LiquidatorAuctionStartedEventName
}

// UnpackAuctionStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuctionStarted(address indexed account, bytes32 indexed position, bool closed)
func (liquidator *Liquidator) UnpackAuctionStartedEvent(log *types.Log) (*LiquidatorAuctionStarted, error) {
	event := "AuctionStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorAuctionStarted)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorAuctionStopped represents a AuctionStopped event raised by the Liquidator contract.
type LiquidatorAuctionStopped struct {
	Account  common.Address
	Position [32]byte
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorAuctionStoppedEventName = "AuctionStopped"

// ContractEventName returns the user-defined event name.
func (LiquidatorAuctionStopped) ContractEventName() string {
	return LiquidatorAuctionStoppedEventName
}

// UnpackAuctionStoppedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuctionStopped(address indexed account, bytes32 indexed position)
func (liquidator *Liquidator) UnpackAuctionStoppedEvent(log *types.Log) (*LiquidatorAuctionStopped, error) {
	event := "AuctionStopped"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorAuctionStopped)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorBought represents a Bought event raised by the Liquidator contract.
type LiquidatorBought struct {
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Cost     *big.Int
	Fee      *big.Int
	Buyer    common.Address
	Receiver common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorBoughtEventName = "Bought"

// ContractEventName returns the user-defined event name.
func (LiquidatorBought) ContractEventName() string {
	return LiquidatorBoughtEventName
}

// UnpackBoughtEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Bought(address indexed account, bytes32 indexed position, address indexed token, uint256 amount, uint256 cost, uint256 fee, address buyer, address receiver)
func (liquidator *Liquidator) UnpackBoughtEvent(log *types.Log) (*LiquidatorBought, error) {
	event := "Bought"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorBought)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorCashSettled represents a CashSettled event raised by the Liquidator contract.
type LiquidatorCashSettled struct {
	Account  common.Address
	Position [32]byte
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorCashSettledEventName = "CashSettled"

// ContractEventName returns the user-defined event name.
func (LiquidatorCashSettled) ContractEventName() string {
	return LiquidatorCashSettledEventName
}

// UnpackCashSettledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CashSettled(address indexed account, bytes32 indexed position, uint256 amount)
func (liquidator *Liquidator) UnpackCashSettledEvent(log *types.Log) (*LiquidatorCashSettled, error) {
	event := "CashSettled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorCashSettled)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorHeld represents a Held event raised by the Liquidator contract.
type LiquidatorHeld struct {
	Account  common.Address
	Position [32]byte
	Until    uint64
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorHeldEventName = "Held"

// ContractEventName returns the user-defined event name.
func (LiquidatorHeld) ContractEventName() string {
	return LiquidatorHeldEventName
}

// UnpackHeldEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Held(address indexed account, bytes32 indexed position, uint64 until)
func (liquidator *Liquidator) UnpackHeldEvent(log *types.Log) (*LiquidatorHeld, error) {
	event := "Held"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorHeld)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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

// LiquidatorSettled represents a Settled event raised by the Liquidator contract.
type LiquidatorSettled struct {
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Price    *big.Int
	Cost     *big.Int
	Fee      *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const LiquidatorSettledEventName = "Settled"

// ContractEventName returns the user-defined event name.
func (LiquidatorSettled) ContractEventName() string {
	return LiquidatorSettledEventName
}

// UnpackSettledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Settled(address indexed account, bytes32 indexed position, address indexed token, uint256 amount, uint256 price, uint256 cost, uint256 fee)
func (liquidator *Liquidator) UnpackSettledEvent(log *types.Log) (*LiquidatorSettled, error) {
	event := "Settled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != liquidator.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(LiquidatorSettled)
	if len(log.Data) > 0 {
		if err := liquidator.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range liquidator.abi.Events[event].Inputs {
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
func (liquidator *Liquidator) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], liquidator.abi.Errors["AuctionAlreadySet"].ID.Bytes()[:4]) {
		return liquidator.UnpackAuctionAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["CannotJudge"].ID.Bytes()[:4]) {
		return liquidator.UnpackCannotJudgeError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["CostAboveLimit"].ID.Bytes()[:4]) {
		return liquidator.UnpackCostAboveLimitError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["InvalidAuction"].ID.Bytes()[:4]) {
		return liquidator.UnpackInvalidAuctionError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NoAuction"].ID.Bytes()[:4]) {
		return liquidator.UnpackNoAuctionError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NotAccountsOwner"].ID.Bytes()[:4]) {
		return liquidator.UnpackNotAccountsOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NotAuction"].ID.Bytes()[:4]) {
		return liquidator.UnpackNotAuctionError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NotBackstop"].ID.Bytes()[:4]) {
		return liquidator.UnpackNotBackstopError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NotForSale"].ID.Bytes()[:4]) {
		return liquidator.UnpackNotForSaleError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NotLiquidatable"].ID.Bytes()[:4]) {
		return liquidator.UnpackNotLiquidatableError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NothingToBuy"].ID.Bytes()[:4]) {
		return liquidator.UnpackNothingToBuyError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["NothingToRecall"].ID.Bytes()[:4]) {
		return liquidator.UnpackNothingToRecallError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["PositionHeld"].ID.Bytes()[:4]) {
		return liquidator.UnpackPositionHeldError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["SafeCastOverflowedIntToUint"].ID.Bytes()[:4]) {
		return liquidator.UnpackSafeCastOverflowedIntToUintError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return liquidator.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["SafeCastOverflowedUintToInt"].ID.Bytes()[:4]) {
		return liquidator.UnpackSafeCastOverflowedUintToIntError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return liquidator.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], liquidator.abi.Errors["StillShort"].ID.Bytes()[:4]) {
		return liquidator.UnpackStillShortError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// LiquidatorAuctionAlreadySet represents a AuctionAlreadySet error raised by the Liquidator contract.
type LiquidatorAuctionAlreadySet struct {
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AuctionAlreadySet(address current)
func LiquidatorAuctionAlreadySetErrorID() common.Hash {
	return common.HexToHash("0xd8333fda46d5bb44554a7362336ad352faf8e18929dd57760defb4feb04bfb71")
}

// UnpackAuctionAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AuctionAlreadySet(address current)
func (liquidator *Liquidator) UnpackAuctionAlreadySetError(raw []byte) (*LiquidatorAuctionAlreadySet, error) {
	out := new(LiquidatorAuctionAlreadySet)
	if err := liquidator.abi.UnpackIntoInterface(out, "AuctionAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorCannotJudge represents a CannotJudge error raised by the Liquidator contract.
type LiquidatorCannotJudge struct {
	Account  common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotJudge(address account, bytes32 position)
func LiquidatorCannotJudgeErrorID() common.Hash {
	return common.HexToHash("0x0f0b5c2aaf1f0f405bbb18057717f93b11e68554b26b498cd01776bdea65673f")
}

// UnpackCannotJudgeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotJudge(address account, bytes32 position)
func (liquidator *Liquidator) UnpackCannotJudgeError(raw []byte) (*LiquidatorCannotJudge, error) {
	out := new(LiquidatorCannotJudge)
	if err := liquidator.abi.UnpackIntoInterface(out, "CannotJudge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorCostAboveLimit represents a CostAboveLimit error raised by the Liquidator contract.
type LiquidatorCostAboveLimit struct {
	Cost    *big.Int
	MaxCost *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CostAboveLimit(uint256 cost, uint256 maxCost)
func LiquidatorCostAboveLimitErrorID() common.Hash {
	return common.HexToHash("0xc36742822815e1d5fd9c4b44dd3de018250452f22e0d9ac2a79c68108488273d")
}

// UnpackCostAboveLimitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CostAboveLimit(uint256 cost, uint256 maxCost)
func (liquidator *Liquidator) UnpackCostAboveLimitError(raw []byte) (*LiquidatorCostAboveLimit, error) {
	out := new(LiquidatorCostAboveLimit)
	if err := liquidator.abi.UnpackIntoInterface(out, "CostAboveLimit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorInvalidAuction represents a InvalidAuction error raised by the Liquidator contract.
type LiquidatorInvalidAuction struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidAuction()
func LiquidatorInvalidAuctionErrorID() common.Hash {
	return common.HexToHash("0x215621609d0db68990f9d9e1ae5d6827e5646f31530f785840378fd986ffa0ca")
}

// UnpackInvalidAuctionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidAuction()
func (liquidator *Liquidator) UnpackInvalidAuctionError(raw []byte) (*LiquidatorInvalidAuction, error) {
	out := new(LiquidatorInvalidAuction)
	if err := liquidator.abi.UnpackIntoInterface(out, "InvalidAuction", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNoAuction represents a NoAuction error raised by the Liquidator contract.
type LiquidatorNoAuction struct {
	Account  common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoAuction(address account, bytes32 position)
func LiquidatorNoAuctionErrorID() common.Hash {
	return common.HexToHash("0x60b14a33925ceb442da6240965d0f4027008fcbc5bef7f921106eeca36d7d3d0")
}

// UnpackNoAuctionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoAuction(address account, bytes32 position)
func (liquidator *Liquidator) UnpackNoAuctionError(raw []byte) (*LiquidatorNoAuction, error) {
	out := new(LiquidatorNoAuction)
	if err := liquidator.abi.UnpackIntoInterface(out, "NoAuction", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNotAccountsOwner represents a NotAccountsOwner error raised by the Liquidator contract.
type LiquidatorNotAccountsOwner struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAccountsOwner(address caller)
func LiquidatorNotAccountsOwnerErrorID() common.Hash {
	return common.HexToHash("0x56b052ac71074dda62d2ec89d4c6b974cec638b68d0e6c76081cb4320c2bd431")
}

// UnpackNotAccountsOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAccountsOwner(address caller)
func (liquidator *Liquidator) UnpackNotAccountsOwnerError(raw []byte) (*LiquidatorNotAccountsOwner, error) {
	out := new(LiquidatorNotAccountsOwner)
	if err := liquidator.abi.UnpackIntoInterface(out, "NotAccountsOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNotAuction represents a NotAuction error raised by the Liquidator contract.
type LiquidatorNotAuction struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAuction(address caller)
func LiquidatorNotAuctionErrorID() common.Hash {
	return common.HexToHash("0xf459dd89b92e732bb4d2d5f95e8b4b5405c659d120e94467be4b958332750795")
}

// UnpackNotAuctionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAuction(address caller)
func (liquidator *Liquidator) UnpackNotAuctionError(raw []byte) (*LiquidatorNotAuction, error) {
	out := new(LiquidatorNotAuction)
	if err := liquidator.abi.UnpackIntoInterface(out, "NotAuction", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNotBackstop represents a NotBackstop error raised by the Liquidator contract.
type LiquidatorNotBackstop struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotBackstop(address caller)
func LiquidatorNotBackstopErrorID() common.Hash {
	return common.HexToHash("0xae8055b28d17e6998cbd5fb24affe947d018d63d8ee19fcf0b6e18781022dd40")
}

// UnpackNotBackstopError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotBackstop(address caller)
func (liquidator *Liquidator) UnpackNotBackstopError(raw []byte) (*LiquidatorNotBackstop, error) {
	out := new(LiquidatorNotBackstop)
	if err := liquidator.abi.UnpackIntoInterface(out, "NotBackstop", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNotForSale represents a NotForSale error raised by the Liquidator contract.
type LiquidatorNotForSale struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotForSale(address token)
func LiquidatorNotForSaleErrorID() common.Hash {
	return common.HexToHash("0x8589eea1e707568c95ac1402f43544027ad06a47e5f92fbc7e246309eab6089e")
}

// UnpackNotForSaleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotForSale(address token)
func (liquidator *Liquidator) UnpackNotForSaleError(raw []byte) (*LiquidatorNotForSale, error) {
	out := new(LiquidatorNotForSale)
	if err := liquidator.abi.UnpackIntoInterface(out, "NotForSale", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNotLiquidatable represents a NotLiquidatable error raised by the Liquidator contract.
type LiquidatorNotLiquidatable struct {
	Account  common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotLiquidatable(address account, bytes32 position)
func LiquidatorNotLiquidatableErrorID() common.Hash {
	return common.HexToHash("0xfb45d351dab2353495a276e25efe1cd744c3372f234c249ba3ab36e3b1bcaaba")
}

// UnpackNotLiquidatableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotLiquidatable(address account, bytes32 position)
func (liquidator *Liquidator) UnpackNotLiquidatableError(raw []byte) (*LiquidatorNotLiquidatable, error) {
	out := new(LiquidatorNotLiquidatable)
	if err := liquidator.abi.UnpackIntoInterface(out, "NotLiquidatable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNothingToBuy represents a NothingToBuy error raised by the Liquidator contract.
type LiquidatorNothingToBuy struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NothingToBuy()
func LiquidatorNothingToBuyErrorID() common.Hash {
	return common.HexToHash("0xb2f536819389fd5ac4a1d6994bba2c706ffe7bfda10e8443370d8d6a398a78f3")
}

// UnpackNothingToBuyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NothingToBuy()
func (liquidator *Liquidator) UnpackNothingToBuyError(raw []byte) (*LiquidatorNothingToBuy, error) {
	out := new(LiquidatorNothingToBuy)
	if err := liquidator.abi.UnpackIntoInterface(out, "NothingToBuy", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorNothingToRecall represents a NothingToRecall error raised by the Liquidator contract.
type LiquidatorNothingToRecall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NothingToRecall()
func LiquidatorNothingToRecallErrorID() common.Hash {
	return common.HexToHash("0x6f0bd69e7685e2649e802243999976ec1e7d1dc7f78ee2b7ea77075b38cc5641")
}

// UnpackNothingToRecallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NothingToRecall()
func (liquidator *Liquidator) UnpackNothingToRecallError(raw []byte) (*LiquidatorNothingToRecall, error) {
	out := new(LiquidatorNothingToRecall)
	if err := liquidator.abi.UnpackIntoInterface(out, "NothingToRecall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorPositionHeld represents a PositionHeld error raised by the Liquidator contract.
type LiquidatorPositionHeld struct {
	Account  common.Address
	Position [32]byte
	Until    uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PositionHeld(address account, bytes32 position, uint64 until)
func LiquidatorPositionHeldErrorID() common.Hash {
	return common.HexToHash("0x7a8acc81833af9fd188f447ea4c99d56daf3f60f0c774c9945865d8cf5fcf9ea")
}

// UnpackPositionHeldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PositionHeld(address account, bytes32 position, uint64 until)
func (liquidator *Liquidator) UnpackPositionHeldError(raw []byte) (*LiquidatorPositionHeld, error) {
	out := new(LiquidatorPositionHeld)
	if err := liquidator.abi.UnpackIntoInterface(out, "PositionHeld", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorSafeCastOverflowedIntToUint represents a SafeCastOverflowedIntToUint error raised by the Liquidator contract.
type LiquidatorSafeCastOverflowedIntToUint struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func LiquidatorSafeCastOverflowedIntToUintErrorID() common.Hash {
	return common.HexToHash("0xa8ce4432b175c373e5f41aba830358e5361584f628450fd436c066323ad91ac2")
}

// UnpackSafeCastOverflowedIntToUintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func (liquidator *Liquidator) UnpackSafeCastOverflowedIntToUintError(raw []byte) (*LiquidatorSafeCastOverflowedIntToUint, error) {
	out := new(LiquidatorSafeCastOverflowedIntToUint)
	if err := liquidator.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntToUint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the Liquidator contract.
type LiquidatorSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func LiquidatorSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (liquidator *Liquidator) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*LiquidatorSafeCastOverflowedUintDowncast, error) {
	out := new(LiquidatorSafeCastOverflowedUintDowncast)
	if err := liquidator.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorSafeCastOverflowedUintToInt represents a SafeCastOverflowedUintToInt error raised by the Liquidator contract.
type LiquidatorSafeCastOverflowedUintToInt struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func LiquidatorSafeCastOverflowedUintToIntErrorID() common.Hash {
	return common.HexToHash("0x24775e0629ae69d78c11bae050651b81820407f300ff750ff2be51e4ce75c37f")
}

// UnpackSafeCastOverflowedUintToIntError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func (liquidator *Liquidator) UnpackSafeCastOverflowedUintToIntError(raw []byte) (*LiquidatorSafeCastOverflowedUintToInt, error) {
	out := new(LiquidatorSafeCastOverflowedUintToInt)
	if err := liquidator.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintToInt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the Liquidator contract.
type LiquidatorSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func LiquidatorSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (liquidator *Liquidator) UnpackSafeERC20FailedOperationError(raw []byte) (*LiquidatorSafeERC20FailedOperation, error) {
	out := new(LiquidatorSafeERC20FailedOperation)
	if err := liquidator.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// LiquidatorStillShort represents a StillShort error raised by the Liquidator contract.
type LiquidatorStillShort struct {
	Account  common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StillShort(address account, bytes32 position)
func LiquidatorStillShortErrorID() common.Hash {
	return common.HexToHash("0x630feaf7563d2b92c16831a1992ff860aa9d8a79abc658ae836c207c1b43b478")
}

// UnpackStillShortError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StillShort(address account, bytes32 position)
func (liquidator *Liquidator) UnpackStillShortError(raw []byte) (*LiquidatorStillShort, error) {
	out := new(LiquidatorStillShort)
	if err := liquidator.abi.UnpackIntoInterface(out, "StillShort", raw); err != nil {
		return nil, err
	}
	return out, nil
}
