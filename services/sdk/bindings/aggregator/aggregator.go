// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package aggregator

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

// AggregatorMetaData contains all meta data concerning the Aggregator contract.
var AggregatorMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"_roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"answeredInRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	ID:  "Aggregator",
}

// Aggregator is an auto generated Go binding around an Ethereum contract.
type Aggregator struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Aggregator) GetABI() abi.ABI {
	return c.abi
}

// NewAggregator creates a new instance of Aggregator.
func NewAggregator() *Aggregator {
	parsed, err := AggregatorMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Aggregator{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Aggregator) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (aggregator *Aggregator) PackDecimals() []byte {
	enc, err := aggregator.abi.Pack("decimals")
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
func (aggregator *Aggregator) TryPackDecimals() ([]byte, error) {
	return aggregator.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (aggregator *Aggregator) UnpackDecimals(data []byte) (uint8, error) {
	out, err := aggregator.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackDescription is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7284e416.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function description() view returns(string)
func (aggregator *Aggregator) PackDescription() []byte {
	enc, err := aggregator.abi.Pack("description")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDescription is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7284e416.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function description() view returns(string)
func (aggregator *Aggregator) TryPackDescription() ([]byte, error) {
	return aggregator.abi.Pack("description")
}

// UnpackDescription is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7284e416.
//
// Solidity: function description() view returns(string)
func (aggregator *Aggregator) UnpackDescription(data []byte) (string, error) {
	out, err := aggregator.abi.Unpack("description", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackGetRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a6fc8f5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) PackGetRoundData(roundId *big.Int) []byte {
	enc, err := aggregator.abi.Pack("getRoundData", roundId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a6fc8f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) TryPackGetRoundData(roundId *big.Int) ([]byte, error) {
	return aggregator.abi.Pack("getRoundData", roundId)
}

// GetRoundDataOutput serves as a container for the return parameters of contract
// method GetRoundData.
type GetRoundDataOutput struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// UnpackGetRoundData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 _roundId) view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) UnpackGetRoundData(data []byte) (GetRoundDataOutput, error) {
	out, err := aggregator.abi.Unpack("getRoundData", data)
	outstruct := new(GetRoundDataOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.RoundId = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Answer = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.StartedAt = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.UpdatedAt = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.AnsweredInRound = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackLatestRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeaf968c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) PackLatestRoundData() []byte {
	enc, err := aggregator.abi.Pack("latestRoundData")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLatestRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeaf968c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) TryPackLatestRoundData() ([]byte, error) {
	return aggregator.abi.Pack("latestRoundData")
}

// LatestRoundDataOutput serves as a container for the return parameters of contract
// method LatestRoundData.
type LatestRoundDataOutput struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// UnpackLatestRoundData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
func (aggregator *Aggregator) UnpackLatestRoundData(data []byte) (LatestRoundDataOutput, error) {
	out, err := aggregator.abi.Unpack("latestRoundData", data)
	outstruct := new(LatestRoundDataOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.RoundId = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Answer = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.StartedAt = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.UpdatedAt = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.AnsweredInRound = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function version() view returns(uint256)
func (aggregator *Aggregator) PackVersion() []byte {
	enc, err := aggregator.abi.Pack("version")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function version() view returns(uint256)
func (aggregator *Aggregator) TryPackVersion() ([]byte, error) {
	return aggregator.abi.Pack("version")
}

// UnpackVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (aggregator *Aggregator) UnpackVersion(data []byte) (*big.Int, error) {
	out, err := aggregator.abi.Unpack("version", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}
