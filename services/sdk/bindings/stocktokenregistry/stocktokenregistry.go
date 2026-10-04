// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stocktokenregistry

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

// StockTokenRegistryMetaData contains all meta data concerning the StockTokenRegistry contract.
var StockTokenRegistryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"isBlocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	ID:  "StockTokenRegistry",
}

// StockTokenRegistry is an auto generated Go binding around an Ethereum contract.
type StockTokenRegistry struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *StockTokenRegistry) GetABI() abi.ABI {
	return c.abi
}

// NewStockTokenRegistry creates a new instance of StockTokenRegistry.
func NewStockTokenRegistry() *StockTokenRegistry {
	parsed, err := StockTokenRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &StockTokenRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *StockTokenRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackIsBlocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbac3951.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isBlocked(address account) view returns(bool)
func (stockTokenRegistry *StockTokenRegistry) PackIsBlocked(account common.Address) []byte {
	enc, err := stockTokenRegistry.abi.Pack("isBlocked", account)
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
func (stockTokenRegistry *StockTokenRegistry) TryPackIsBlocked(account common.Address) ([]byte, error) {
	return stockTokenRegistry.abi.Pack("isBlocked", account)
}

// UnpackIsBlocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfbac3951.
//
// Solidity: function isBlocked(address account) view returns(bool)
func (stockTokenRegistry *StockTokenRegistry) UnpackIsBlocked(data []byte) (bool, error) {
	out, err := stockTokenRegistry.abi.Unpack("isBlocked", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}
