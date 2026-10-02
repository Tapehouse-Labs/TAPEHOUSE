// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package band

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

// BandMetaData contains all meta data concerning the Band contract.
var BandMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"chainlinkFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"redstoneFeedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"indexFeedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chainConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"sequencerUptimeFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainlinkRegularHours\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"corporateAction\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"effectiveAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"multiplierBefore\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"multiplierAfter\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"halt\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"signedHalt\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"until\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"issuedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"oraclePaused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"haltSigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"price\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"packageTimestampMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"writtenAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quote\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"live\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"mid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"halfBps\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"low\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"high\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sequencerSettled\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"session\",\"inputs\":[],\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nyse\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nyseNext\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"changeMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"boundaryMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"variance\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"writePrices\",\"inputs\":[{\"name\":\"feedIds\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Anchored\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"chainlinkPrice\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"indexPrice\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"updatedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HaltSignerUpdated\",\"inputs\":[{\"name\":\"previousSigner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSigner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HaltWritten\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"halted\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"issuedAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"expiresAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MultiplierConfirmed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"effectiveAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MultiplierRecorded\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"multiplierBefore\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"multiplierAfter\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"effectiveAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceWritten\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"packageTimestampMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CalldataMustHaveValidPayload\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CalldataOverOrUnderFlow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DataTimestampCannotBeZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncompleteStatus\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectUnsignedMetadataSize\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientNumberOfUniqueSigners\",\"inputs\":[{\"name\":\"receivedSignersCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requiredSignersCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[{\"name\":\"signedHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"PackageNotNewer\",\"inputs\":[{\"name\":\"feedId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"storedTimestampMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"packageTimestampMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"SignerNotAuthorised\",\"inputs\":[{\"name\":\"receivedSigner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TimestampFromTooLongFuture\",\"inputs\":[{\"name\":\"receivedTimestampSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampIsTooOld\",\"inputs\":[{\"name\":\"receivedTimestampSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampsMustBeEqual\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooLargeValueByteSize\",\"inputs\":[{\"name\":\"valueByteSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "Band",
}

// Band is an auto generated Go binding around an Ethereum contract.
type Band struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Band) GetABI() abi.ABI {
	return c.abi
}

// NewBand creates a new instance of Band.
func NewBand() *Band {
	parsed, err := BandMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Band{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Band) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8539ccf4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function asset(bytes32 symbol) view returns(address chainlinkFeed, bytes32 redstoneFeedId, bytes32 indexFeedId, address token)
func (band *Band) PackAsset(symbol [32]byte) []byte {
	enc, err := band.abi.Pack("asset", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8539ccf4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function asset(bytes32 symbol) view returns(address chainlinkFeed, bytes32 redstoneFeedId, bytes32 indexFeedId, address token)
func (band *Band) TryPackAsset(symbol [32]byte) ([]byte, error) {
	return band.abi.Pack("asset", symbol)
}

// AssetOutput serves as a container for the return parameters of contract
// method Asset.
type AssetOutput struct {
	ChainlinkFeed  common.Address
	RedstoneFeedId [32]byte
	IndexFeedId    [32]byte
	Token          common.Address
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8539ccf4.
//
// Solidity: function asset(bytes32 symbol) view returns(address chainlinkFeed, bytes32 redstoneFeedId, bytes32 indexFeedId, address token)
func (band *Band) UnpackAsset(data []byte) (AssetOutput, error) {
	out, err := band.abi.Unpack("asset", data)
	outstruct := new(AssetOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ChainlinkFeed = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.RedstoneFeedId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.IndexFeedId = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Token = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	return *outstruct, nil
}

// PackChainConfig is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xebb43738.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function chainConfig() view returns(address sequencerUptimeFeed, bool chainlinkRegularHours)
func (band *Band) PackChainConfig() []byte {
	enc, err := band.abi.Pack("chainConfig")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackChainConfig is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xebb43738.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function chainConfig() view returns(address sequencerUptimeFeed, bool chainlinkRegularHours)
func (band *Band) TryPackChainConfig() ([]byte, error) {
	return band.abi.Pack("chainConfig")
}

// ChainConfigOutput serves as a container for the return parameters of contract
// method ChainConfig.
type ChainConfigOutput struct {
	SequencerUptimeFeed   common.Address
	ChainlinkRegularHours bool
}

// UnpackChainConfig is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xebb43738.
//
// Solidity: function chainConfig() view returns(address sequencerUptimeFeed, bool chainlinkRegularHours)
func (band *Band) UnpackChainConfig(data []byte) (ChainConfigOutput, error) {
	out, err := band.abi.Unpack("chainConfig", data)
	outstruct := new(ChainConfigOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.SequencerUptimeFeed = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ChainlinkRegularHours = *abi.ConvertType(out[1], new(bool)).(*bool)
	return *outstruct, nil
}

// PackCorporateAction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a6d4db7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function corporateAction(bytes32 symbol) view returns(uint8 status, uint64 effectiveAt, uint128 multiplierBefore, uint128 multiplierAfter)
func (band *Band) PackCorporateAction(symbol [32]byte) []byte {
	enc, err := band.abi.Pack("corporateAction", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCorporateAction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a6d4db7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function corporateAction(bytes32 symbol) view returns(uint8 status, uint64 effectiveAt, uint128 multiplierBefore, uint128 multiplierAfter)
func (band *Band) TryPackCorporateAction(symbol [32]byte) ([]byte, error) {
	return band.abi.Pack("corporateAction", symbol)
}

// CorporateActionOutput serves as a container for the return parameters of contract
// method CorporateAction.
type CorporateActionOutput struct {
	Status           uint8
	EffectiveAt      uint64
	MultiplierBefore *big.Int
	MultiplierAfter  *big.Int
}

// UnpackCorporateAction is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a6d4db7.
//
// Solidity: function corporateAction(bytes32 symbol) view returns(uint8 status, uint64 effectiveAt, uint128 multiplierBefore, uint128 multiplierAfter)
func (band *Band) UnpackCorporateAction(data []byte) (CorporateActionOutput, error) {
	out, err := band.abi.Unpack("corporateAction", data)
	outstruct := new(CorporateActionOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Status = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.EffectiveAt = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.MultiplierBefore = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.MultiplierAfter = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackHalt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x537e7736.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function halt(bytes32 symbol) view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (band *Band) PackHalt(symbol [32]byte) []byte {
	enc, err := band.abi.Pack("halt", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHalt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x537e7736.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function halt(bytes32 symbol) view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (band *Band) TryPackHalt(symbol [32]byte) ([]byte, error) {
	return band.abi.Pack("halt", symbol)
}

// HaltOutput serves as a container for the return parameters of contract
// method Halt.
type HaltOutput struct {
	SignedHalt   bool
	Until        uint64
	IssuedAt     uint64
	OraclePaused bool
}

// UnpackHalt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x537e7736.
//
// Solidity: function halt(bytes32 symbol) view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (band *Band) UnpackHalt(data []byte) (HaltOutput, error) {
	out, err := band.abi.Unpack("halt", data)
	outstruct := new(HaltOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.SignedHalt = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Until = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.IssuedAt = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.OraclePaused = *abi.ConvertType(out[3], new(bool)).(*bool)
	return *outstruct, nil
}

// PackHaltSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x829793d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function haltSigner() view returns(address)
func (band *Band) PackHaltSigner() []byte {
	enc, err := band.abi.Pack("haltSigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHaltSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x829793d4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function haltSigner() view returns(address)
func (band *Band) TryPackHaltSigner() ([]byte, error) {
	return band.abi.Pack("haltSigner")
}

// UnpackHaltSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x829793d4.
//
// Solidity: function haltSigner() view returns(address)
func (band *Band) UnpackHaltSigner(data []byte) (common.Address, error) {
	out, err := band.abi.Unpack("haltSigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x88f3543a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function price(bytes32 feedId) view returns(uint256 value, uint64 packageTimestampMs, uint64 writtenAt)
func (band *Band) PackPrice(feedId [32]byte) []byte {
	enc, err := band.abi.Pack("price", feedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x88f3543a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function price(bytes32 feedId) view returns(uint256 value, uint64 packageTimestampMs, uint64 writtenAt)
func (band *Band) TryPackPrice(feedId [32]byte) ([]byte, error) {
	return band.abi.Pack("price", feedId)
}

// PriceOutput serves as a container for the return parameters of contract
// method Price.
type PriceOutput struct {
	Value              *big.Int
	PackageTimestampMs uint64
	WrittenAt          uint64
}

// UnpackPrice is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x88f3543a.
//
// Solidity: function price(bytes32 feedId) view returns(uint256 value, uint64 packageTimestampMs, uint64 writtenAt)
func (band *Band) UnpackPrice(data []byte) (PriceOutput, error) {
	out, err := band.abi.Unpack("price", data)
	outstruct := new(PriceOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Value = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.PackageTimestampMs = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.WrittenAt = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x15f05b05.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function quote(bytes32 symbol) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high)
func (band *Band) PackQuote(symbol [32]byte) []byte {
	enc, err := band.abi.Pack("quote", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x15f05b05.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function quote(bytes32 symbol) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high)
func (band *Band) TryPackQuote(symbol [32]byte) ([]byte, error) {
	return band.abi.Pack("quote", symbol)
}

// QuoteOutput serves as a container for the return parameters of contract
// method Quote.
type QuoteOutput struct {
	State   uint8
	Live    uint8
	Mid     uint64
	HalfBps uint64
	Low     uint64
	High    *big.Int
}

// UnpackQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x15f05b05.
//
// Solidity: function quote(bytes32 symbol) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high)
func (band *Band) UnpackQuote(data []byte) (QuoteOutput, error) {
	out, err := band.abi.Unpack("quote", data)
	outstruct := new(QuoteOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.State = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Live = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Mid = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.HalfBps = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Low = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.High = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackSequencerSettled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf561479d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sequencerSettled() view returns(bool)
func (band *Band) PackSequencerSettled() []byte {
	enc, err := band.abi.Pack("sequencerSettled")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSequencerSettled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf561479d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sequencerSettled() view returns(bool)
func (band *Band) TryPackSequencerSettled() ([]byte, error) {
	return band.abi.Pack("sequencerSettled")
}

// UnpackSequencerSettled is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf561479d.
//
// Solidity: function sequencerSettled() view returns(bool)
func (band *Band) UnpackSequencerSettled(data []byte) (bool, error) {
	out, err := band.abi.Unpack("sequencerSettled", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSession is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e3568b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function session() view returns(uint8 state, uint8 nyse, uint8 nyseNext, uint64 changeMs, uint64 boundaryMs)
func (band *Band) PackSession() []byte {
	enc, err := band.abi.Pack("session")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSession is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e3568b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function session() view returns(uint8 state, uint8 nyse, uint8 nyseNext, uint64 changeMs, uint64 boundaryMs)
func (band *Band) TryPackSession() ([]byte, error) {
	return band.abi.Pack("session")
}

// SessionOutput serves as a container for the return parameters of contract
// method Session.
type SessionOutput struct {
	State      uint8
	Nyse       uint8
	NyseNext   uint8
	ChangeMs   uint64
	BoundaryMs uint64
}

// UnpackSession is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e3568b8.
//
// Solidity: function session() view returns(uint8 state, uint8 nyse, uint8 nyseNext, uint64 changeMs, uint64 boundaryMs)
func (band *Band) UnpackSession(data []byte) (SessionOutput, error) {
	out, err := band.abi.Unpack("session", data)
	outstruct := new(SessionOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.State = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Nyse = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.NyseNext = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.ChangeMs = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.BoundaryMs = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackVariance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdecf0015.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function variance(bytes32 feedId) view returns(uint128)
func (band *Band) PackVariance(feedId [32]byte) []byte {
	enc, err := band.abi.Pack("variance", feedId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVariance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdecf0015.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function variance(bytes32 feedId) view returns(uint128)
func (band *Band) TryPackVariance(feedId [32]byte) ([]byte, error) {
	return band.abi.Pack("variance", feedId)
}

// UnpackVariance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdecf0015.
//
// Solidity: function variance(bytes32 feedId) view returns(uint128)
func (band *Band) UnpackVariance(data []byte) (*big.Int, error) {
	out, err := band.abi.Unpack("variance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWritePrices is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4267488b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function writePrices(bytes32[] feedIds, bytes payload) returns(uint256[])
func (band *Band) PackWritePrices(feedIds [][32]byte, payload []byte) []byte {
	enc, err := band.abi.Pack("writePrices", feedIds, payload)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWritePrices is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4267488b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function writePrices(bytes32[] feedIds, bytes payload) returns(uint256[])
func (band *Band) TryPackWritePrices(feedIds [][32]byte, payload []byte) ([]byte, error) {
	return band.abi.Pack("writePrices", feedIds, payload)
}

// UnpackWritePrices is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4267488b.
//
// Solidity: function writePrices(bytes32[] feedIds, bytes payload) returns(uint256[])
func (band *Band) UnpackWritePrices(data []byte) ([]*big.Int, error) {
	out, err := band.abi.Unpack("writePrices", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// BandAnchored represents a Anchored event raised by the Band contract.
type BandAnchored struct {
	Symbol         [32]byte
	ChainlinkPrice uint64
	IndexPrice     uint64
	UpdatedAt      uint64
	Raw            *types.Log // Blockchain specific contextual infos
}

const BandAnchoredEventName = "Anchored"

// ContractEventName returns the user-defined event name.
func (BandAnchored) ContractEventName() string {
	return BandAnchoredEventName
}

// UnpackAnchoredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Anchored(bytes32 indexed symbol, uint64 chainlinkPrice, uint64 indexPrice, uint64 updatedAt)
func (band *Band) UnpackAnchoredEvent(log *types.Log) (*BandAnchored, error) {
	event := "Anchored"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandAnchored)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandHaltSignerUpdated represents a HaltSignerUpdated event raised by the Band contract.
type BandHaltSignerUpdated struct {
	PreviousSigner common.Address
	NewSigner      common.Address
	Raw            *types.Log // Blockchain specific contextual infos
}

const BandHaltSignerUpdatedEventName = "HaltSignerUpdated"

// ContractEventName returns the user-defined event name.
func (BandHaltSignerUpdated) ContractEventName() string {
	return BandHaltSignerUpdatedEventName
}

// UnpackHaltSignerUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event HaltSignerUpdated(address indexed previousSigner, address indexed newSigner)
func (band *Band) UnpackHaltSignerUpdatedEvent(log *types.Log) (*BandHaltSignerUpdated, error) {
	event := "HaltSignerUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandHaltSignerUpdated)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandHaltWritten represents a HaltWritten event raised by the Band contract.
type BandHaltWritten struct {
	Symbol    [32]byte
	Halted    bool
	IssuedAt  uint64
	ExpiresAt uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const BandHaltWrittenEventName = "HaltWritten"

// ContractEventName returns the user-defined event name.
func (BandHaltWritten) ContractEventName() string {
	return BandHaltWrittenEventName
}

// UnpackHaltWrittenEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event HaltWritten(bytes32 indexed symbol, bool halted, uint64 issuedAt, uint64 expiresAt)
func (band *Band) UnpackHaltWrittenEvent(log *types.Log) (*BandHaltWritten, error) {
	event := "HaltWritten"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandHaltWritten)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandMultiplierConfirmed represents a MultiplierConfirmed event raised by the Band contract.
type BandMultiplierConfirmed struct {
	Symbol      [32]byte
	EffectiveAt uint64
	Raw         *types.Log // Blockchain specific contextual infos
}

const BandMultiplierConfirmedEventName = "MultiplierConfirmed"

// ContractEventName returns the user-defined event name.
func (BandMultiplierConfirmed) ContractEventName() string {
	return BandMultiplierConfirmedEventName
}

// UnpackMultiplierConfirmedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MultiplierConfirmed(bytes32 indexed symbol, uint64 effectiveAt)
func (band *Band) UnpackMultiplierConfirmedEvent(log *types.Log) (*BandMultiplierConfirmed, error) {
	event := "MultiplierConfirmed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandMultiplierConfirmed)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandMultiplierRecorded represents a MultiplierRecorded event raised by the Band contract.
type BandMultiplierRecorded struct {
	Symbol           [32]byte
	MultiplierBefore *big.Int
	MultiplierAfter  *big.Int
	EffectiveAt      uint64
	Raw              *types.Log // Blockchain specific contextual infos
}

const BandMultiplierRecordedEventName = "MultiplierRecorded"

// ContractEventName returns the user-defined event name.
func (BandMultiplierRecorded) ContractEventName() string {
	return BandMultiplierRecordedEventName
}

// UnpackMultiplierRecordedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MultiplierRecorded(bytes32 indexed symbol, uint128 multiplierBefore, uint128 multiplierAfter, uint64 effectiveAt)
func (band *Band) UnpackMultiplierRecordedEvent(log *types.Log) (*BandMultiplierRecorded, error) {
	event := "MultiplierRecorded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandMultiplierRecorded)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Band contract.
type BandOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const BandOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (BandOwnershipTransferStarted) ContractEventName() string {
	return BandOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (band *Band) UnpackOwnershipTransferStartedEvent(log *types.Log) (*BandOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandOwnershipTransferred represents a OwnershipTransferred event raised by the Band contract.
type BandOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const BandOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (BandOwnershipTransferred) ContractEventName() string {
	return BandOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (band *Band) UnpackOwnershipTransferredEvent(log *types.Log) (*BandOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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

// BandPriceWritten represents a PriceWritten event raised by the Band contract.
type BandPriceWritten struct {
	FeedId             [32]byte
	Value              *big.Int
	PackageTimestampMs uint64
	Raw                *types.Log // Blockchain specific contextual infos
}

const BandPriceWrittenEventName = "PriceWritten"

// ContractEventName returns the user-defined event name.
func (BandPriceWritten) ContractEventName() string {
	return BandPriceWrittenEventName
}

// UnpackPriceWrittenEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PriceWritten(bytes32 indexed feedId, uint256 value, uint64 packageTimestampMs)
func (band *Band) UnpackPriceWrittenEvent(log *types.Log) (*BandPriceWritten, error) {
	event := "PriceWritten"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != band.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandPriceWritten)
	if len(log.Data) > 0 {
		if err := band.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range band.abi.Events[event].Inputs {
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
func (band *Band) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], band.abi.Errors["CalldataMustHaveValidPayload"].ID.Bytes()[:4]) {
		return band.UnpackCalldataMustHaveValidPayloadError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["CalldataOverOrUnderFlow"].ID.Bytes()[:4]) {
		return band.UnpackCalldataOverOrUnderFlowError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["DataTimestampCannotBeZero"].ID.Bytes()[:4]) {
		return band.UnpackDataTimestampCannotBeZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["IncompleteStatus"].ID.Bytes()[:4]) {
		return band.UnpackIncompleteStatusError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["IncorrectUnsignedMetadataSize"].ID.Bytes()[:4]) {
		return band.UnpackIncorrectUnsignedMetadataSizeError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["InsufficientNumberOfUniqueSigners"].ID.Bytes()[:4]) {
		return band.UnpackInsufficientNumberOfUniqueSignersError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return band.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["PackageNotNewer"].ID.Bytes()[:4]) {
		return band.UnpackPackageNotNewerError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["SignerNotAuthorised"].ID.Bytes()[:4]) {
		return band.UnpackSignerNotAuthorisedError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["TimestampFromTooLongFuture"].ID.Bytes()[:4]) {
		return band.UnpackTimestampFromTooLongFutureError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["TimestampIsTooOld"].ID.Bytes()[:4]) {
		return band.UnpackTimestampIsTooOldError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["TimestampsMustBeEqual"].ID.Bytes()[:4]) {
		return band.UnpackTimestampsMustBeEqualError(raw[4:])
	}
	if bytes.Equal(raw[:4], band.abi.Errors["TooLargeValueByteSize"].ID.Bytes()[:4]) {
		return band.UnpackTooLargeValueByteSizeError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// BandCalldataMustHaveValidPayload represents a CalldataMustHaveValidPayload error raised by the Band contract.
type BandCalldataMustHaveValidPayload struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CalldataMustHaveValidPayload()
func BandCalldataMustHaveValidPayloadErrorID() common.Hash {
	return common.HexToHash("0xe7764c9e8e0caa7ff5b20e6f3065a6816f3e1be01d10790878af38007164e2da")
}

// UnpackCalldataMustHaveValidPayloadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CalldataMustHaveValidPayload()
func (band *Band) UnpackCalldataMustHaveValidPayloadError(raw []byte) (*BandCalldataMustHaveValidPayload, error) {
	out := new(BandCalldataMustHaveValidPayload)
	if err := band.abi.UnpackIntoInterface(out, "CalldataMustHaveValidPayload", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandCalldataOverOrUnderFlow represents a CalldataOverOrUnderFlow error raised by the Band contract.
type BandCalldataOverOrUnderFlow struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CalldataOverOrUnderFlow()
func BandCalldataOverOrUnderFlowErrorID() common.Hash {
	return common.HexToHash("0x5796f78a52824208b9c59aa4ea6a172cbe542ecc3198660c040f96765907dbb4")
}

// UnpackCalldataOverOrUnderFlowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CalldataOverOrUnderFlow()
func (band *Band) UnpackCalldataOverOrUnderFlowError(raw []byte) (*BandCalldataOverOrUnderFlow, error) {
	out := new(BandCalldataOverOrUnderFlow)
	if err := band.abi.UnpackIntoInterface(out, "CalldataOverOrUnderFlow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandDataTimestampCannotBeZero represents a DataTimestampCannotBeZero error raised by the Band contract.
type BandDataTimestampCannotBeZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DataTimestampCannotBeZero()
func BandDataTimestampCannotBeZeroErrorID() common.Hash {
	return common.HexToHash("0xdfb25a79cad62b89b7e8ce6836ecfeed6784cbd493e31afd94a2eb87f7685eac")
}

// UnpackDataTimestampCannotBeZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DataTimestampCannotBeZero()
func (band *Band) UnpackDataTimestampCannotBeZeroError(raw []byte) (*BandDataTimestampCannotBeZero, error) {
	out := new(BandDataTimestampCannotBeZero)
	if err := band.abi.UnpackIntoInterface(out, "DataTimestampCannotBeZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandIncompleteStatus represents a IncompleteStatus error raised by the Band contract.
type BandIncompleteStatus struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncompleteStatus()
func BandIncompleteStatusErrorID() common.Hash {
	return common.HexToHash("0x140a6d4d7b83b1a1290dc1b4b883c3d23a83e19452c323a6a2a6e7eff607a82e")
}

// UnpackIncompleteStatusError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncompleteStatus()
func (band *Band) UnpackIncompleteStatusError(raw []byte) (*BandIncompleteStatus, error) {
	out := new(BandIncompleteStatus)
	if err := band.abi.UnpackIntoInterface(out, "IncompleteStatus", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandIncorrectUnsignedMetadataSize represents a IncorrectUnsignedMetadataSize error raised by the Band contract.
type BandIncorrectUnsignedMetadataSize struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectUnsignedMetadataSize()
func BandIncorrectUnsignedMetadataSizeErrorID() common.Hash {
	return common.HexToHash("0xc30a7bd79ae58b3376eee0583848bba17d7381a8ed2924118a3b588367f8de08")
}

// UnpackIncorrectUnsignedMetadataSizeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectUnsignedMetadataSize()
func (band *Band) UnpackIncorrectUnsignedMetadataSizeError(raw []byte) (*BandIncorrectUnsignedMetadataSize, error) {
	out := new(BandIncorrectUnsignedMetadataSize)
	if err := band.abi.UnpackIntoInterface(out, "IncorrectUnsignedMetadataSize", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandInsufficientNumberOfUniqueSigners represents a InsufficientNumberOfUniqueSigners error raised by the Band contract.
type BandInsufficientNumberOfUniqueSigners struct {
	ReceivedSignersCount *big.Int
	RequiredSignersCount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientNumberOfUniqueSigners(uint256 receivedSignersCount, uint256 requiredSignersCount)
func BandInsufficientNumberOfUniqueSignersErrorID() common.Hash {
	return common.HexToHash("0x2b13aef51436e37db0ebf10aa5df2b1f452aae20ad834f63b5953fc8fda60d0a")
}

// UnpackInsufficientNumberOfUniqueSignersError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientNumberOfUniqueSigners(uint256 receivedSignersCount, uint256 requiredSignersCount)
func (band *Band) UnpackInsufficientNumberOfUniqueSignersError(raw []byte) (*BandInsufficientNumberOfUniqueSigners, error) {
	out := new(BandInsufficientNumberOfUniqueSigners)
	if err := band.abi.UnpackIntoInterface(out, "InsufficientNumberOfUniqueSigners", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandInvalidSignature represents a InvalidSignature error raised by the Band contract.
type BandInvalidSignature struct {
	SignedHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature(bytes32 signedHash)
func BandInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x666b7cba1c1e18e9729cc96bc32dad18492fd8ee8fcb6327b7a2337ed37534f6")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature(bytes32 signedHash)
func (band *Band) UnpackInvalidSignatureError(raw []byte) (*BandInvalidSignature, error) {
	out := new(BandInvalidSignature)
	if err := band.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandPackageNotNewer represents a PackageNotNewer error raised by the Band contract.
type BandPackageNotNewer struct {
	FeedId             [32]byte
	StoredTimestampMs  uint64
	PackageTimestampMs uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PackageNotNewer(bytes32 feedId, uint64 storedTimestampMs, uint64 packageTimestampMs)
func BandPackageNotNewerErrorID() common.Hash {
	return common.HexToHash("0x5d1d5d095998f1a82d0e33e491d01caecc92478d5377998e6a3655aef1a66674")
}

// UnpackPackageNotNewerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PackageNotNewer(bytes32 feedId, uint64 storedTimestampMs, uint64 packageTimestampMs)
func (band *Band) UnpackPackageNotNewerError(raw []byte) (*BandPackageNotNewer, error) {
	out := new(BandPackageNotNewer)
	if err := band.abi.UnpackIntoInterface(out, "PackageNotNewer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandSignerNotAuthorised represents a SignerNotAuthorised error raised by the Band contract.
type BandSignerNotAuthorised struct {
	ReceivedSigner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignerNotAuthorised(address receivedSigner)
func BandSignerNotAuthorisedErrorID() common.Hash {
	return common.HexToHash("0xec459bc05285e5b716796daffa940a2f1e2daf552dcda6cc9bc70577f533fce3")
}

// UnpackSignerNotAuthorisedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignerNotAuthorised(address receivedSigner)
func (band *Band) UnpackSignerNotAuthorisedError(raw []byte) (*BandSignerNotAuthorised, error) {
	out := new(BandSignerNotAuthorised)
	if err := band.abi.UnpackIntoInterface(out, "SignerNotAuthorised", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandTimestampFromTooLongFuture represents a TimestampFromTooLongFuture error raised by the Band contract.
type BandTimestampFromTooLongFuture struct {
	ReceivedTimestampSeconds *big.Int
	BlockTimestamp           *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampFromTooLongFuture(uint256 receivedTimestampSeconds, uint256 blockTimestamp)
func BandTimestampFromTooLongFutureErrorID() common.Hash {
	return common.HexToHash("0xb6b0916d36f78340ee5c06f513af9e20984692a4bde4eaf4be072bfa6089de90")
}

// UnpackTimestampFromTooLongFutureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampFromTooLongFuture(uint256 receivedTimestampSeconds, uint256 blockTimestamp)
func (band *Band) UnpackTimestampFromTooLongFutureError(raw []byte) (*BandTimestampFromTooLongFuture, error) {
	out := new(BandTimestampFromTooLongFuture)
	if err := band.abi.UnpackIntoInterface(out, "TimestampFromTooLongFuture", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandTimestampIsTooOld represents a TimestampIsTooOld error raised by the Band contract.
type BandTimestampIsTooOld struct {
	ReceivedTimestampSeconds *big.Int
	BlockTimestamp           *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampIsTooOld(uint256 receivedTimestampSeconds, uint256 blockTimestamp)
func BandTimestampIsTooOldErrorID() common.Hash {
	return common.HexToHash("0x0321d0b5b1ea0b5b678f54bffda902de00e8e7dc2a7257325feb9740c0d8a6d6")
}

// UnpackTimestampIsTooOldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampIsTooOld(uint256 receivedTimestampSeconds, uint256 blockTimestamp)
func (band *Band) UnpackTimestampIsTooOldError(raw []byte) (*BandTimestampIsTooOld, error) {
	out := new(BandTimestampIsTooOld)
	if err := band.abi.UnpackIntoInterface(out, "TimestampIsTooOld", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandTimestampsMustBeEqual represents a TimestampsMustBeEqual error raised by the Band contract.
type BandTimestampsMustBeEqual struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimestampsMustBeEqual()
func BandTimestampsMustBeEqualErrorID() common.Hash {
	return common.HexToHash("0x4cbc4742aa642e552d3094ab8187d27efcff7abd1047ac67e374330b798fd830")
}

// UnpackTimestampsMustBeEqualError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimestampsMustBeEqual()
func (band *Band) UnpackTimestampsMustBeEqualError(raw []byte) (*BandTimestampsMustBeEqual, error) {
	out := new(BandTimestampsMustBeEqual)
	if err := band.abi.UnpackIntoInterface(out, "TimestampsMustBeEqual", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandTooLargeValueByteSize represents a TooLargeValueByteSize error raised by the Band contract.
type BandTooLargeValueByteSize struct {
	ValueByteSize *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooLargeValueByteSize(uint256 valueByteSize)
func BandTooLargeValueByteSizeErrorID() common.Hash {
	return common.HexToHash("0xc000fc42deb431e500942408b9894a5625c9fbc62cb538fa5315c3f468f79639")
}

// UnpackTooLargeValueByteSizeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooLargeValueByteSize(uint256 valueByteSize)
func (band *Band) UnpackTooLargeValueByteSizeError(raw []byte) (*BandTooLargeValueByteSize, error) {
	out := new(BandTooLargeValueByteSize)
	if err := band.abi.UnpackIntoInterface(out, "TooLargeValueByteSize", raw); err != nil {
		return nil, err
	}
	return out, nil
}
