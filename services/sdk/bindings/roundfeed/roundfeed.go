// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package roundfeed

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

// RoundFeedMetaData contains all meta data concerning the RoundFeed contract.
var RoundFeedMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"dataFeedId_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"description_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"stockToken_\",\"type\":\"address\",\"internalType\":\"contractIMultiplier\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dataFeedId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"answer\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"startedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stockToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMultiplier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"update\",\"inputs\":[],\"outputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AnswerUpdated\",\"inputs\":[{\"name\":\"current\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"},{\"name\":\"roundId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AnswerTooLarge\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientNumberOfUniqueSigners\",\"inputs\":[{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidPayload\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotNewer\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SignerNotAuthorised\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TimestampIsFromTheFuture\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampIsTooOld\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampsMustMatchAndBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValueTooWide\",\"inputs\":[{\"name\":\"bytesWide\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "RoundFeed",
}

// RoundFeed is an auto generated Go binding around an Ethereum contract.
type RoundFeed struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *RoundFeed) GetABI() abi.ABI {
	return c.abi
}

// NewRoundFeed creates a new instance of RoundFeed.
func NewRoundFeed() *RoundFeed {
	parsed, err := RoundFeedMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &RoundFeed{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *RoundFeed) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(bytes32 dataFeedId_, string description_, address stockToken_) returns()
func (roundFeed *RoundFeed) PackConstructor(dataFeedId_ [32]byte, description_ string, stockToken_ common.Address) []byte {
	enc, err := roundFeed.abi.Pack("", dataFeedId_, description_, stockToken_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDataFeedId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x370c826b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function dataFeedId() view returns(bytes32)
func (roundFeed *RoundFeed) PackDataFeedId() []byte {
	enc, err := roundFeed.abi.Pack("dataFeedId")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDataFeedId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x370c826b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function dataFeedId() view returns(bytes32)
func (roundFeed *RoundFeed) TryPackDataFeedId() ([]byte, error) {
	return roundFeed.abi.Pack("dataFeedId")
}

// UnpackDataFeedId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x370c826b.
//
// Solidity: function dataFeedId() view returns(bytes32)
func (roundFeed *RoundFeed) UnpackDataFeedId(data []byte) ([32]byte, error) {
	out, err := roundFeed.abi.Unpack("dataFeedId", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (roundFeed *RoundFeed) PackDecimals() []byte {
	enc, err := roundFeed.abi.Pack("decimals")
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
func (roundFeed *RoundFeed) TryPackDecimals() ([]byte, error) {
	return roundFeed.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (roundFeed *RoundFeed) UnpackDecimals(data []byte) (uint8, error) {
	out, err := roundFeed.abi.Unpack("decimals", data)
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
func (roundFeed *RoundFeed) PackDescription() []byte {
	enc, err := roundFeed.abi.Pack("description")
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
func (roundFeed *RoundFeed) TryPackDescription() ([]byte, error) {
	return roundFeed.abi.Pack("description")
}

// UnpackDescription is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7284e416.
//
// Solidity: function description() view returns(string)
func (roundFeed *RoundFeed) UnpackDescription(data []byte) (string, error) {
	out, err := roundFeed.abi.Unpack("description", data)
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
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256 answer, uint256 startedAt, uint256 updatedAt, uint80)
func (roundFeed *RoundFeed) PackGetRoundData(roundId *big.Int) []byte {
	enc, err := roundFeed.abi.Pack("getRoundData", roundId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a6fc8f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256 answer, uint256 startedAt, uint256 updatedAt, uint80)
func (roundFeed *RoundFeed) TryPackGetRoundData(roundId *big.Int) ([]byte, error) {
	return roundFeed.abi.Pack("getRoundData", roundId)
}

// GetRoundDataOutput serves as a container for the return parameters of contract
// method GetRoundData.
type GetRoundDataOutput struct {
	Arg0      *big.Int
	Answer    *big.Int
	StartedAt *big.Int
	UpdatedAt *big.Int
	Arg4      *big.Int
}

// UnpackGetRoundData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256 answer, uint256 startedAt, uint256 updatedAt, uint80)
func (roundFeed *RoundFeed) UnpackGetRoundData(data []byte) (GetRoundDataOutput, error) {
	out, err := roundFeed.abi.Unpack("getRoundData", data)
	outstruct := new(GetRoundDataOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Answer = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.StartedAt = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.UpdatedAt = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.Arg4 = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackLatestRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeaf968c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function latestRoundData() view returns(uint80, int256, uint256, uint256, uint80)
func (roundFeed *RoundFeed) PackLatestRoundData() []byte {
	enc, err := roundFeed.abi.Pack("latestRoundData")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLatestRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeaf968c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function latestRoundData() view returns(uint80, int256, uint256, uint256, uint80)
func (roundFeed *RoundFeed) TryPackLatestRoundData() ([]byte, error) {
	return roundFeed.abi.Pack("latestRoundData")
}

// LatestRoundDataOutput serves as a container for the return parameters of contract
// method LatestRoundData.
type LatestRoundDataOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
	Arg2 *big.Int
	Arg3 *big.Int
	Arg4 *big.Int
}

// UnpackLatestRoundData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfeaf968c.
//
// Solidity: function latestRoundData() view returns(uint80, int256, uint256, uint256, uint80)
func (roundFeed *RoundFeed) UnpackLatestRoundData(data []byte) (LatestRoundDataOutput, error) {
	out, err := roundFeed.abi.Unpack("latestRoundData", data)
	outstruct := new(LatestRoundDataOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Arg2 = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.Arg3 = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.Arg4 = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackStockToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbaa8e786.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stockToken() view returns(address)
func (roundFeed *RoundFeed) PackStockToken() []byte {
	enc, err := roundFeed.abi.Pack("stockToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStockToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbaa8e786.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stockToken() view returns(address)
func (roundFeed *RoundFeed) TryPackStockToken() ([]byte, error) {
	return roundFeed.abi.Pack("stockToken")
}

// UnpackStockToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbaa8e786.
//
// Solidity: function stockToken() view returns(address)
func (roundFeed *RoundFeed) UnpackStockToken(data []byte) (common.Address, error) {
	out, err := roundFeed.abi.Unpack("stockToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2e62045.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function update() returns(uint80 roundId)
func (roundFeed *RoundFeed) PackUpdate() []byte {
	enc, err := roundFeed.abi.Pack("update")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa2e62045.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function update() returns(uint80 roundId)
func (roundFeed *RoundFeed) TryPackUpdate() ([]byte, error) {
	return roundFeed.abi.Pack("update")
}

// UnpackUpdate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa2e62045.
//
// Solidity: function update() returns(uint80 roundId)
func (roundFeed *RoundFeed) UnpackUpdate(data []byte) (*big.Int, error) {
	out, err := roundFeed.abi.Unpack("update", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function version() view returns(uint256)
func (roundFeed *RoundFeed) PackVersion() []byte {
	enc, err := roundFeed.abi.Pack("version")
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
func (roundFeed *RoundFeed) TryPackVersion() ([]byte, error) {
	return roundFeed.abi.Pack("version")
}

// UnpackVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (roundFeed *RoundFeed) UnpackVersion(data []byte) (*big.Int, error) {
	out, err := roundFeed.abi.Unpack("version", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// RoundFeedAnswerUpdated represents a AnswerUpdated event raised by the RoundFeed contract.
type RoundFeedAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const RoundFeedAnswerUpdatedEventName = "AnswerUpdated"

// ContractEventName returns the user-defined event name.
func (RoundFeedAnswerUpdated) ContractEventName() string {
	return RoundFeedAnswerUpdatedEventName
}

// UnpackAnswerUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt)
func (roundFeed *RoundFeed) UnpackAnswerUpdatedEvent(log *types.Log) (*RoundFeedAnswerUpdated, error) {
	event := "AnswerUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != roundFeed.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(RoundFeedAnswerUpdated)
	if len(log.Data) > 0 {
		if err := roundFeed.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range roundFeed.abi.Events[event].Inputs {
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
func (roundFeed *RoundFeed) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["AnswerTooLarge"].ID.Bytes()[:4]) {
		return roundFeed.UnpackAnswerTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["InsufficientNumberOfUniqueSigners"].ID.Bytes()[:4]) {
		return roundFeed.UnpackInsufficientNumberOfUniqueSignersError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["InvalidPayload"].ID.Bytes()[:4]) {
		return roundFeed.UnpackInvalidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return roundFeed.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["NotNewer"].ID.Bytes()[:4]) {
		return roundFeed.UnpackNotNewerError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["SignerNotAuthorised"].ID.Bytes()[:4]) {
		return roundFeed.UnpackSignerNotAuthorisedError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["TimestampIsFromTheFuture"].ID.Bytes()[:4]) {
		return roundFeed.UnpackTimestampIsFromTheFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["TimestampIsTooOld"].ID.Bytes()[:4]) {
		return roundFeed.UnpackTimestampIsTooOldError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["TimestampsMustMatchAndBePositive"].ID.Bytes()[:4]) {
		return roundFeed.UnpackTimestampsMustMatchAndBePositiveError(raw[4:])
	}
	if bytes.Equal(raw[:4], roundFeed.abi.Errors["ValueTooWide"].ID.Bytes()[:4]) {
		return roundFeed.UnpackValueTooWideError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// RoundFeedAnswerTooLarge represents a AnswerTooLarge error raised by the RoundFeed contract.
type RoundFeedAnswerTooLarge struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AnswerTooLarge()
func RoundFeedAnswerTooLargeErrorID() common.Hash {
	return common.HexToHash("0x12318764b194413b52193d69fdb5c0e10afc79a49ca960bd665f474b6c79c3ea")
}

// UnpackAnswerTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AnswerTooLarge()
func (roundFeed *RoundFeed) UnpackAnswerTooLargeError(raw []byte) (*RoundFeedAnswerTooLarge, error) {
	out := new(RoundFeedAnswerTooLarge)
	if err := roundFeed.abi.UnpackIntoInterface(out, "AnswerTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedInsufficientNumberOfUniqueSigners represents a InsufficientNumberOfUniqueSigners error raised by the RoundFeed contract.
type RoundFeedInsufficientNumberOfUniqueSigners struct {
	Received *big.Int
	Required *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientNumberOfUniqueSigners(uint256 received, uint256 required)
func RoundFeedInsufficientNumberOfUniqueSignersErrorID() common.Hash {
	return common.HexToHash("0x2b13aef51436e37db0ebf10aa5df2b1f452aae20ad834f63b5953fc8fda60d0a")
}

// UnpackInsufficientNumberOfUniqueSignersError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientNumberOfUniqueSigners(uint256 received, uint256 required)
func (roundFeed *RoundFeed) UnpackInsufficientNumberOfUniqueSignersError(raw []byte) (*RoundFeedInsufficientNumberOfUniqueSigners, error) {
	out := new(RoundFeedInsufficientNumberOfUniqueSigners)
	if err := roundFeed.abi.UnpackIntoInterface(out, "InsufficientNumberOfUniqueSigners", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedInvalidPayload represents a InvalidPayload error raised by the RoundFeed contract.
type RoundFeedInvalidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPayload()
func RoundFeedInvalidPayloadErrorID() common.Hash {
	return common.HexToHash("0x7c6953f9f55fcfd713fe1d72e370e4757e17ae7187e1acdcc02b4b76b56a9a35")
}

// UnpackInvalidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPayload()
func (roundFeed *RoundFeed) UnpackInvalidPayloadError(raw []byte) (*RoundFeedInvalidPayload, error) {
	out := new(RoundFeedInvalidPayload)
	if err := roundFeed.abi.UnpackIntoInterface(out, "InvalidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedInvalidSignature represents a InvalidSignature error raised by the RoundFeed contract.
type RoundFeedInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func RoundFeedInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (roundFeed *RoundFeed) UnpackInvalidSignatureError(raw []byte) (*RoundFeedInvalidSignature, error) {
	out := new(RoundFeedInvalidSignature)
	if err := roundFeed.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedNotNewer represents a NotNewer error raised by the RoundFeed contract.
type RoundFeedNotNewer struct {
	Timestamp     *big.Int
	LastTimestamp *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotNewer(uint256 timestamp, uint256 lastTimestamp)
func RoundFeedNotNewerErrorID() common.Hash {
	return common.HexToHash("0xa3051cadb2230b21b1fe8d85416cf7f3cf3b97691188f1a64ef0b2a20b69f87a")
}

// UnpackNotNewerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotNewer(uint256 timestamp, uint256 lastTimestamp)
func (roundFeed *RoundFeed) UnpackNotNewerError(raw []byte) (*RoundFeedNotNewer, error) {
	out := new(RoundFeedNotNewer)
	if err := roundFeed.abi.UnpackIntoInterface(out, "NotNewer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedSignerNotAuthorised represents a SignerNotAuthorised error raised by the RoundFeed contract.
type RoundFeedSignerNotAuthorised struct {
	Signer common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerNotAuthorised(address signer)
func RoundFeedSignerNotAuthorisedErrorID() common.Hash {
	return common.HexToHash("0xec459bc05285e5b716796daffa940a2f1e2daf552dcda6cc9bc70577f533fce3")
}

// UnpackSignerNotAuthorisedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerNotAuthorised(address signer)
func (roundFeed *RoundFeed) UnpackSignerNotAuthorisedError(raw []byte) (*RoundFeedSignerNotAuthorised, error) {
	out := new(RoundFeedSignerNotAuthorised)
	if err := roundFeed.abi.UnpackIntoInterface(out, "SignerNotAuthorised", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedTimestampIsFromTheFuture represents a TimestampIsFromTheFuture error raised by the RoundFeed contract.
type RoundFeedTimestampIsFromTheFuture struct {
	Timestamp      *big.Int
	BlockTimestamp *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampIsFromTheFuture(uint256 timestamp, uint256 blockTimestamp)
func RoundFeedTimestampIsFromTheFutureErrorID() common.Hash {
	return common.HexToHash("0x1056e2ae08e930b01b5af21623f241b7f9d200bffe311de7af3fc00e7a0b4301")
}

// UnpackTimestampIsFromTheFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampIsFromTheFuture(uint256 timestamp, uint256 blockTimestamp)
func (roundFeed *RoundFeed) UnpackTimestampIsFromTheFutureError(raw []byte) (*RoundFeedTimestampIsFromTheFuture, error) {
	out := new(RoundFeedTimestampIsFromTheFuture)
	if err := roundFeed.abi.UnpackIntoInterface(out, "TimestampIsFromTheFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedTimestampIsTooOld represents a TimestampIsTooOld error raised by the RoundFeed contract.
type RoundFeedTimestampIsTooOld struct {
	Timestamp      *big.Int
	BlockTimestamp *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampIsTooOld(uint256 timestamp, uint256 blockTimestamp)
func RoundFeedTimestampIsTooOldErrorID() common.Hash {
	return common.HexToHash("0x0321d0b5b1ea0b5b678f54bffda902de00e8e7dc2a7257325feb9740c0d8a6d6")
}

// UnpackTimestampIsTooOldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampIsTooOld(uint256 timestamp, uint256 blockTimestamp)
func (roundFeed *RoundFeed) UnpackTimestampIsTooOldError(raw []byte) (*RoundFeedTimestampIsTooOld, error) {
	out := new(RoundFeedTimestampIsTooOld)
	if err := roundFeed.abi.UnpackIntoInterface(out, "TimestampIsTooOld", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedTimestampsMustMatchAndBePositive represents a TimestampsMustMatchAndBePositive error raised by the RoundFeed contract.
type RoundFeedTimestampsMustMatchAndBePositive struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampsMustMatchAndBePositive()
func RoundFeedTimestampsMustMatchAndBePositiveErrorID() common.Hash {
	return common.HexToHash("0x97cddea8f1b19aeb5035e8ab95d706d3c5420fb2de3785bddc15ee4050576360")
}

// UnpackTimestampsMustMatchAndBePositiveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampsMustMatchAndBePositive()
func (roundFeed *RoundFeed) UnpackTimestampsMustMatchAndBePositiveError(raw []byte) (*RoundFeedTimestampsMustMatchAndBePositive, error) {
	out := new(RoundFeedTimestampsMustMatchAndBePositive)
	if err := roundFeed.abi.UnpackIntoInterface(out, "TimestampsMustMatchAndBePositive", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// RoundFeedValueTooWide represents a ValueTooWide error raised by the RoundFeed contract.
type RoundFeedValueTooWide struct {
	BytesWide *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ValueTooWide(uint256 bytesWide)
func RoundFeedValueTooWideErrorID() common.Hash {
	return common.HexToHash("0x0af119b7488854d1f64d10bc11443520ceeabbb4feead48eb1e9e411f321ad14")
}

// UnpackValueTooWideError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ValueTooWide(uint256 bytesWide)
func (roundFeed *RoundFeed) UnpackValueTooWideError(raw []byte) (*RoundFeedValueTooWide, error) {
	out := new(RoundFeedValueTooWide)
	if err := roundFeed.abi.UnpackIntoInterface(out, "ValueTooWide", raw); err != nil {
		return nil, err
	}
	return out, nil
}
