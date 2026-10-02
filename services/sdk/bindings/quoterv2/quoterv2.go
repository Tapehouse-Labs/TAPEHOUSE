// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package quoterv2

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

// IQuoterV2QuoteExactInputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactInputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	AmountIn          *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// IQuoterV2QuoteExactOutputSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IQuoterV2QuoteExactOutputSingleParams struct {
	TokenIn           common.Address
	TokenOut          common.Address
	Amount            *big.Int
	Fee               *big.Int
	SqrtPriceLimitX96 *big.Int
}

// QuoterV2MetaData contains all meta data concerning the QuoterV2 contract.
var QuoterV2MetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quoteExactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIQuoterV2.QuoteExactInputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIQuoterV2.QuoteExactOutputSingleParams\",\"components\":[{\"name\":\"tokenIn\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenOut\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"sqrtPriceX96After\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"initializedTicksCrossed\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]",
	ID:  "QuoterV2",
}

// QuoterV2 is an auto generated Go binding around an Ethereum contract.
type QuoterV2 struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *QuoterV2) GetABI() abi.ABI {
	return c.abi
}

// NewQuoterV2 creates a new instance of QuoterV2.
func NewQuoterV2() *QuoterV2 {
	parsed, err := QuoterV2MetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &QuoterV2{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *QuoterV2) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc45a0155.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function factory() view returns(address)
func (quoterV2 *QuoterV2) PackFactory() []byte {
	enc, err := quoterV2.abi.Pack("factory")
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
func (quoterV2 *QuoterV2) TryPackFactory() ([]byte, error) {
	return quoterV2.abi.Pack("factory")
}

// UnpackFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (quoterV2 *QuoterV2) UnpackFactory(data []byte) (common.Address, error) {
	out, err := quoterV2.abi.Unpack("factory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackQuoteExactInputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6a5026a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) PackQuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) []byte {
	enc, err := quoterV2.abi.Pack("quoteExactInputSingle", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQuoteExactInputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6a5026a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) TryPackQuoteExactInputSingle(params IQuoterV2QuoteExactInputSingleParams) ([]byte, error) {
	return quoterV2.abi.Pack("quoteExactInputSingle", params)
}

// QuoteExactInputSingleOutput serves as a container for the return parameters of contract
// method QuoteExactInputSingle.
type QuoteExactInputSingleOutput struct {
	AmountOut               *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}

// UnpackQuoteExactInputSingle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6a5026a.
//
// Solidity: function quoteExactInputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountOut, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) UnpackQuoteExactInputSingle(data []byte) (QuoteExactInputSingleOutput, error) {
	out, err := quoterV2.abi.Unpack("quoteExactInputSingle", data)
	outstruct := new(QuoteExactInputSingleOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.AmountOut = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.SqrtPriceX96After = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackQuoteExactOutputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd21704a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) PackQuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) []byte {
	enc, err := quoterV2.abi.Pack("quoteExactOutputSingle", params)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQuoteExactOutputSingle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd21704a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) TryPackQuoteExactOutputSingle(params IQuoterV2QuoteExactOutputSingleParams) ([]byte, error) {
	return quoterV2.abi.Pack("quoteExactOutputSingle", params)
}

// QuoteExactOutputSingleOutput serves as a container for the return parameters of contract
// method QuoteExactOutputSingle.
type QuoteExactOutputSingleOutput struct {
	AmountIn                *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}

// UnpackQuoteExactOutputSingle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd21704a.
//
// Solidity: function quoteExactOutputSingle((address,address,uint256,uint24,uint160) params) returns(uint256 amountIn, uint160 sqrtPriceX96After, uint32 initializedTicksCrossed, uint256 gasEstimate)
func (quoterV2 *QuoterV2) UnpackQuoteExactOutputSingle(data []byte) (QuoteExactOutputSingleOutput, error) {
	out, err := quoterV2.abi.Unpack("quoteExactOutputSingle", data)
	outstruct := new(QuoteExactOutputSingleOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.AmountIn = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.SqrtPriceX96After = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.InitializedTicksCrossed = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	outstruct.GasEstimate = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, nil
}
