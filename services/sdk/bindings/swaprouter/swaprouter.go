// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swaprouter

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

// IV3SwapRouterExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	AmountIn          *big.Int
	AmountOutMinimum  *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IV3SwapRouterExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IV3SwapRouterExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Fee               *big.Int
	Recipient         common.Address
	AmountOut         *big.Int
	AmountInMaximum   *big.Int
	SqrtPriceLimitX96 *big.Int
}

// SwapRouterMetaData contains all meta data concerning the SwapRouter contract.
var SwapRouterMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"exactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV3SwapRouter.ExactInputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOutMinimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"exactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV3SwapRouter.ExactOutputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountInMaximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"}]",
	ID:  "SwapRouter",
}

// SwapRouter is an auto generated Go binding around an Ethereum contract.
type SwapRouter struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *SwapRouter) GetABI() abi.ABI {
	return c.abi
}

// NewSwapRouter creates a new instance of SwapRouter.
func NewSwapRouter() *SwapRouter {
	parsed, err := SwapRouterMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &SwapRouter{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *SwapRouter) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackExactInputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04e45aaf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (swapRouter *SwapRouter) PackExactInputSingle(params IV3SwapRouterExactInputSingleParams) []byte {
	enc, err := swapRouter.abi.Pack("exactInputSingle", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExactInputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x04e45aaf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (swapRouter *SwapRouter) TryPackExactInputSingle(params IV3SwapRouterExactInputSingleParams) ([]byte, error) {
	return swapRouter.abi.Pack("exactInputSingle", params)
}

// UnpackExactInputSingle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x04e45aaf.
//
// Solidity: function exactInputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountOut)
func (swapRouter *SwapRouter) UnpackExactInputSingle(data []byte) (*big.Int, error) {
	out, err := swapRouter.abi.Unpack("exactInputSingle", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackExactOutputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5023b4df.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (swapRouter *SwapRouter) PackExactOutputSingle(params IV3SwapRouterExactOutputSingleParams) []byte {
	enc, err := swapRouter.abi.Pack("exactOutputSingle", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackExactOutputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5023b4df.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (swapRouter *SwapRouter) TryPackExactOutputSingle(params IV3SwapRouterExactOutputSingleParams) ([]byte, error) {
	return swapRouter.abi.Pack("exactOutputSingle", params)
}

// UnpackExactOutputSingle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5023b4df.
//
// Solidity: function exactOutputSingle((address,address,uint24,address,uint256,uint256,uint160) params) payable returns(uint256 amountIn)
func (swapRouter *SwapRouter) UnpackExactOutputSingle(data []byte) (*big.Int, error) {
	out, err := swapRouter.abi.Unpack("exactOutputSingle", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc45a0155.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function factory() view returns(address)
func (swapRouter *SwapRouter) PackFactory() []byte {
	enc, err := swapRouter.abi.Pack("factory")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc45a0155.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function factory() view returns(address)
func (swapRouter *SwapRouter) TryPackFactory() ([]byte, error) {
	return swapRouter.abi.Pack("factory")
}

// UnpackFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (swapRouter *SwapRouter) UnpackFactory(data []byte) (common.Address, error) {
	out, err := swapRouter.abi.Unpack("factory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}
