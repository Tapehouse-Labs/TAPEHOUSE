// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bandfeed

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

// BandFeedBand is an auto generated low-level Go binding around an user-defined struct.
type BandFeedBand struct {
	State             uint8
	Live              uint8
	Mid               uint64
	HalfBps           uint64
	Low               uint64
	High              *big.Int
	Variance          *big.Int
	Session           uint8
	Nyse              uint8
	NyseNext          uint8
	NyseChangeMs      uint64
	SessionBoundaryMs uint64
	SignedHalt        bool
	HaltUntil         uint64
	HaltIssuedAt      uint64
	OraclePaused      bool
	SequencerSettled  bool
	TwapValid         bool
	Twap              *big.Int
	PremiumBps        *big.Int
}

// BandFeedMetaData contains all meta data concerning the BandFeed contract.
var BandFeedMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"band_\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"symbol_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"side_\",\"type\":\"uint8\",\"internalType\":\"enumBandFeed.Side\"},{\"name\":\"pool_\",\"type\":\"address\",\"internalType\":\"contractIUniswapV3Pool\"},{\"name\":\"description_\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"SEAL_WINDOW_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TWAP_WINDOW\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoundData\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBand\",\"inputs\":[],\"outputs\":[{\"name\":\"b\",\"type\":\"tuple\",\"internalType\":\"structBandFeed.Band\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"live\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"mid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"halfBps\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"low\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"high\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"variance\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"session\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nyse\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nyseNext\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"nyseChangeMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sessionBoundaryMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"signedHalt\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"haltUntil\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"haltIssuedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"oraclePaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"sequencerSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"twapValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"twap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"premiumBps\",\"type\":\"int256\",\"internalType\":\"int256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestRoundData\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV3Pool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"seal\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"seals\",\"inputs\":[{\"name\":\"reopenMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"live\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"mid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"halfBps\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"low\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"high\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"sealedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"side\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumBandFeed.Side\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Sealed\",\"inputs\":[{\"name\":\"reopenMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"state\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"live\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"mid\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"halfBps\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"low\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"high\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"NoAnswer\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}]},{\"type\":\"error\",\"name\":\"NotSealWindow\",\"inputs\":[{\"name\":\"session\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"reopenMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PoolNotInitialized\",\"inputs\":[{\"name\":\"pool\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PoolWithoutToken\",\"inputs\":[{\"name\":\"pool\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"T\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAsset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "BandFeed",
}

// BandFeed is an auto generated Go binding around an Ethereum contract.
type BandFeed struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *BandFeed) GetABI() abi.ABI {
	return c.abi
}

// NewBandFeed creates a new instance of BandFeed.
func NewBandFeed() *BandFeed {
	parsed, err := BandFeedMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &BandFeed{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *BandFeed) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address band_, bytes32 symbol_, uint8 side_, address pool_, string description_) returns()
func (bandFeed *BandFeed) PackConstructor(band_ common.Address, symbol_ [32]byte, side_ uint8, pool_ common.Address, description_ string) []byte {
	enc, err := bandFeed.abi.Pack("", band_, symbol_, side_, pool_, description_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSEALWINDOWMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x001c3fa7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function SEAL_WINDOW_MS() view returns(uint64)
func (bandFeed *BandFeed) PackSEALWINDOWMS() []byte {
	enc, err := bandFeed.abi.Pack("SEAL_WINDOW_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSEALWINDOWMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x001c3fa7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function SEAL_WINDOW_MS() view returns(uint64)
func (bandFeed *BandFeed) TryPackSEALWINDOWMS() ([]byte, error) {
	return bandFeed.abi.Pack("SEAL_WINDOW_MS")
}

// UnpackSEALWINDOWMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x001c3fa7.
//
// Solidity: function SEAL_WINDOW_MS() view returns(uint64)
func (bandFeed *BandFeed) UnpackSEALWINDOWMS(data []byte) (uint64, error) {
	out, err := bandFeed.abi.Unpack("SEAL_WINDOW_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackTWAPWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2c1d1d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function TWAP_WINDOW() view returns(uint32)
func (bandFeed *BandFeed) PackTWAPWINDOW() []byte {
	enc, err := bandFeed.abi.Pack("TWAP_WINDOW")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTWAPWINDOW is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe2c1d1d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function TWAP_WINDOW() view returns(uint32)
func (bandFeed *BandFeed) TryPackTWAPWINDOW() ([]byte, error) {
	return bandFeed.abi.Pack("TWAP_WINDOW")
}

// UnpackTWAPWINDOW is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe2c1d1d1.
//
// Solidity: function TWAP_WINDOW() view returns(uint32)
func (bandFeed *BandFeed) UnpackTWAPWINDOW(data []byte) (uint32, error) {
	out, err := bandFeed.abi.Unpack("TWAP_WINDOW", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (bandFeed *BandFeed) PackBand() []byte {
	enc, err := bandFeed.abi.Pack("band")
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
func (bandFeed *BandFeed) TryPackBand() ([]byte, error) {
	return bandFeed.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (bandFeed *BandFeed) UnpackBand(data []byte) (common.Address, error) {
	out, err := bandFeed.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (bandFeed *BandFeed) PackDecimals() []byte {
	enc, err := bandFeed.abi.Pack("decimals")
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
func (bandFeed *BandFeed) TryPackDecimals() ([]byte, error) {
	return bandFeed.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (bandFeed *BandFeed) UnpackDecimals(data []byte) (uint8, error) {
	out, err := bandFeed.abi.Unpack("decimals", data)
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
func (bandFeed *BandFeed) PackDescription() []byte {
	enc, err := bandFeed.abi.Pack("description")
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
func (bandFeed *BandFeed) TryPackDescription() ([]byte, error) {
	return bandFeed.abi.Pack("description")
}

// UnpackDescription is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7284e416.
//
// Solidity: function description() view returns(string)
func (bandFeed *BandFeed) UnpackDescription(data []byte) (string, error) {
	out, err := bandFeed.abi.Unpack("description", data)
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
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256, uint256, uint256, uint80)
func (bandFeed *BandFeed) PackGetRoundData(roundId *big.Int) []byte {
	enc, err := bandFeed.abi.Pack("getRoundData", roundId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a6fc8f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256, uint256, uint256, uint80)
func (bandFeed *BandFeed) TryPackGetRoundData(roundId *big.Int) ([]byte, error) {
	return bandFeed.abi.Pack("getRoundData", roundId)
}

// GetRoundDataOutput serves as a container for the return parameters of contract
// method GetRoundData.
type GetRoundDataOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
	Arg2 *big.Int
	Arg3 *big.Int
	Arg4 *big.Int
}

// UnpackGetRoundData is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a6fc8f5.
//
// Solidity: function getRoundData(uint80 roundId) view returns(uint80, int256, uint256, uint256, uint80)
func (bandFeed *BandFeed) UnpackGetRoundData(data []byte) (GetRoundDataOutput, error) {
	out, err := bandFeed.abi.Unpack("getRoundData", data)
	outstruct := new(GetRoundDataOutput)
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

// PackLatestBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa1fec7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function latestBand() view returns((uint8,uint8,uint64,uint64,uint64,uint128,uint128,uint8,uint8,uint8,uint64,uint64,bool,uint64,uint64,bool,bool,bool,uint256,int256) b)
func (bandFeed *BandFeed) PackLatestBand() []byte {
	enc, err := bandFeed.abi.Pack("latestBand")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLatestBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa1fec7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function latestBand() view returns((uint8,uint8,uint64,uint64,uint64,uint128,uint128,uint8,uint8,uint8,uint64,uint64,bool,uint64,uint64,bool,bool,bool,uint256,int256) b)
func (bandFeed *BandFeed) TryPackLatestBand() ([]byte, error) {
	return bandFeed.abi.Pack("latestBand")
}

// UnpackLatestBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa1fec7e.
//
// Solidity: function latestBand() view returns((uint8,uint8,uint64,uint64,uint64,uint128,uint128,uint8,uint8,uint8,uint64,uint64,bool,uint64,uint64,bool,bool,bool,uint256,int256) b)
func (bandFeed *BandFeed) UnpackLatestBand(data []byte) (BandFeedBand, error) {
	out, err := bandFeed.abi.Unpack("latestBand", data)
	if err != nil {
		return *new(BandFeedBand), err
	}
	out0 := *abi.ConvertType(out[0], new(BandFeedBand)).(*BandFeedBand)
	return out0, nil
}

// PackLatestRoundData is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeaf968c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function latestRoundData() view returns(uint80, int256, uint256, uint256, uint80)
func (bandFeed *BandFeed) PackLatestRoundData() []byte {
	enc, err := bandFeed.abi.Pack("latestRoundData")
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
func (bandFeed *BandFeed) TryPackLatestRoundData() ([]byte, error) {
	return bandFeed.abi.Pack("latestRoundData")
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
func (bandFeed *BandFeed) UnpackLatestRoundData(data []byte) (LatestRoundDataOutput, error) {
	out, err := bandFeed.abi.Unpack("latestRoundData", data)
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

// PackPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16f0115b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pool() view returns(address)
func (bandFeed *BandFeed) PackPool() []byte {
	enc, err := bandFeed.abi.Pack("pool")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPool is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16f0115b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pool() view returns(address)
func (bandFeed *BandFeed) TryPackPool() ([]byte, error) {
	return bandFeed.abi.Pack("pool")
}

// UnpackPool is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x16f0115b.
//
// Solidity: function pool() view returns(address)
func (bandFeed *BandFeed) UnpackPool(data []byte) (common.Address, error) {
	out, err := bandFeed.abi.Unpack("pool", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fb27b85.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function seal() returns()
func (bandFeed *BandFeed) PackSeal() []byte {
	enc, err := bandFeed.abi.Pack("seal")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fb27b85.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function seal() returns()
func (bandFeed *BandFeed) TryPackSeal() ([]byte, error) {
	return bandFeed.abi.Pack("seal")
}

// PackSeals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d75962b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function seals(uint64 reopenMs) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high, uint64 sealedAt)
func (bandFeed *BandFeed) PackSeals(reopenMs uint64) []byte {
	enc, err := bandFeed.abi.Pack("seals", reopenMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d75962b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function seals(uint64 reopenMs) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high, uint64 sealedAt)
func (bandFeed *BandFeed) TryPackSeals(reopenMs uint64) ([]byte, error) {
	return bandFeed.abi.Pack("seals", reopenMs)
}

// SealsOutput serves as a container for the return parameters of contract
// method Seals.
type SealsOutput struct {
	State    uint8
	Live     uint8
	Mid      uint64
	HalfBps  uint64
	Low      uint64
	High     *big.Int
	SealedAt uint64
}

// UnpackSeals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d75962b.
//
// Solidity: function seals(uint64 reopenMs) view returns(uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high, uint64 sealedAt)
func (bandFeed *BandFeed) UnpackSeals(data []byte) (SealsOutput, error) {
	out, err := bandFeed.abi.Unpack("seals", data)
	outstruct := new(SealsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.State = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Live = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Mid = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.HalfBps = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Low = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.High = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	outstruct.SealedAt = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackSide is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8475c028.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function side() view returns(uint8)
func (bandFeed *BandFeed) PackSide() []byte {
	enc, err := bandFeed.abi.Pack("side")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSide is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8475c028.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function side() view returns(uint8)
func (bandFeed *BandFeed) TryPackSide() ([]byte, error) {
	return bandFeed.abi.Pack("side")
}

// UnpackSide is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8475c028.
//
// Solidity: function side() view returns(uint8)
func (bandFeed *BandFeed) UnpackSide(data []byte) (uint8, error) {
	out, err := bandFeed.abi.Unpack("side", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(bytes32)
func (bandFeed *BandFeed) PackSymbol() []byte {
	enc, err := bandFeed.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(bytes32)
func (bandFeed *BandFeed) TryPackSymbol() ([]byte, error) {
	return bandFeed.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(bytes32)
func (bandFeed *BandFeed) UnpackSymbol(data []byte) ([32]byte, error) {
	out, err := bandFeed.abi.Unpack("symbol", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackVersion is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x54fd4d50.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function version() view returns(uint256)
func (bandFeed *BandFeed) PackVersion() []byte {
	enc, err := bandFeed.abi.Pack("version")
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
func (bandFeed *BandFeed) TryPackVersion() ([]byte, error) {
	return bandFeed.abi.Pack("version")
}

// UnpackVersion is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x54fd4d50.
//
// Solidity: function version() view returns(uint256)
func (bandFeed *BandFeed) UnpackVersion(data []byte) (*big.Int, error) {
	out, err := bandFeed.abi.Unpack("version", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// BandFeedSealed represents a Sealed event raised by the BandFeed contract.
type BandFeedSealed struct {
	ReopenMs uint64
	State    uint8
	Live     uint8
	Mid      uint64
	HalfBps  uint64
	Low      uint64
	High     *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const BandFeedSealedEventName = "Sealed"

// ContractEventName returns the user-defined event name.
func (BandFeedSealed) ContractEventName() string {
	return BandFeedSealedEventName
}

// UnpackSealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sealed(uint64 indexed reopenMs, uint8 state, uint8 live, uint64 mid, uint64 halfBps, uint64 low, uint128 high)
func (bandFeed *BandFeed) UnpackSealedEvent(log *types.Log) (*BandFeedSealed, error) {
	event := "Sealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != bandFeed.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(BandFeedSealed)
	if len(log.Data) > 0 {
		if err := bandFeed.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range bandFeed.abi.Events[event].Inputs {
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
func (bandFeed *BandFeed) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["NoAnswer"].ID.Bytes()[:4]) {
		return bandFeed.UnpackNoAnswerError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["NoRound"].ID.Bytes()[:4]) {
		return bandFeed.UnpackNoRoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["NotSealWindow"].ID.Bytes()[:4]) {
		return bandFeed.UnpackNotSealWindowError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["PoolNotInitialized"].ID.Bytes()[:4]) {
		return bandFeed.UnpackPoolNotInitializedError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["PoolWithoutToken"].ID.Bytes()[:4]) {
		return bandFeed.UnpackPoolWithoutTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return bandFeed.UnpackSequencerNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["T"].ID.Bytes()[:4]) {
		return bandFeed.UnpackTError(raw[4:])
	}
	if bytes.Equal(raw[:4], bandFeed.abi.Errors["UnknownAsset"].ID.Bytes()[:4]) {
		return bandFeed.UnpackUnknownAssetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// BandFeedNoAnswer represents a NoAnswer error raised by the BandFeed contract.
type BandFeedNoAnswer struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoAnswer(bytes32 symbol)
func BandFeedNoAnswerErrorID() common.Hash {
	return common.HexToHash("0x39c5542dd45701eaba9384a6d4e59bd719eacb189e293835099e854560d87f0d")
}

// UnpackNoAnswerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoAnswer(bytes32 symbol)
func (bandFeed *BandFeed) UnpackNoAnswerError(raw []byte) (*BandFeedNoAnswer, error) {
	out := new(BandFeedNoAnswer)
	if err := bandFeed.abi.UnpackIntoInterface(out, "NoAnswer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedNoRound represents a NoRound error raised by the BandFeed contract.
type BandFeedNoRound struct {
	RoundId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoRound(uint80 roundId)
func BandFeedNoRoundErrorID() common.Hash {
	return common.HexToHash("0x2d24dabd7d25ff47b715533d828b281c1eb5b01a9edc5b1eeddec75a80baa0e5")
}

// UnpackNoRoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoRound(uint80 roundId)
func (bandFeed *BandFeed) UnpackNoRoundError(raw []byte) (*BandFeedNoRound, error) {
	out := new(BandFeedNoRound)
	if err := bandFeed.abi.UnpackIntoInterface(out, "NoRound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedNotSealWindow represents a NotSealWindow error raised by the BandFeed contract.
type BandFeedNotSealWindow struct {
	Session  uint8
	ReopenMs uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSealWindow(uint8 session, uint64 reopenMs)
func BandFeedNotSealWindowErrorID() common.Hash {
	return common.HexToHash("0x7c193929cb0bc6eff14453c7d972e4053ec403be3fff1f71d12d4cfd6f7f3822")
}

// UnpackNotSealWindowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSealWindow(uint8 session, uint64 reopenMs)
func (bandFeed *BandFeed) UnpackNotSealWindowError(raw []byte) (*BandFeedNotSealWindow, error) {
	out := new(BandFeedNotSealWindow)
	if err := bandFeed.abi.UnpackIntoInterface(out, "NotSealWindow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedPoolNotInitialized represents a PoolNotInitialized error raised by the BandFeed contract.
type BandFeedPoolNotInitialized struct {
	Pool common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PoolNotInitialized(address pool)
func BandFeedPoolNotInitializedErrorID() common.Hash {
	return common.HexToHash("0x4bdace133b6d2057f22364f72aba381b746b942afe76ad0de27f78d76b996c42")
}

// UnpackPoolNotInitializedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PoolNotInitialized(address pool)
func (bandFeed *BandFeed) UnpackPoolNotInitializedError(raw []byte) (*BandFeedPoolNotInitialized, error) {
	out := new(BandFeedPoolNotInitialized)
	if err := bandFeed.abi.UnpackIntoInterface(out, "PoolNotInitialized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedPoolWithoutToken represents a PoolWithoutToken error raised by the BandFeed contract.
type BandFeedPoolWithoutToken struct {
	Pool  common.Address
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PoolWithoutToken(address pool, address token)
func BandFeedPoolWithoutTokenErrorID() common.Hash {
	return common.HexToHash("0xefd789abed3d3330ec99bd1612c2b6e0985fbc38c4bde797883ec971a05b5735")
}

// UnpackPoolWithoutTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PoolWithoutToken(address pool, address token)
func (bandFeed *BandFeed) UnpackPoolWithoutTokenError(raw []byte) (*BandFeedPoolWithoutToken, error) {
	out := new(BandFeedPoolWithoutToken)
	if err := bandFeed.abi.UnpackIntoInterface(out, "PoolWithoutToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedSequencerNotSettled represents a SequencerNotSettled error raised by the BandFeed contract.
type BandFeedSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func BandFeedSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (bandFeed *BandFeed) UnpackSequencerNotSettledError(raw []byte) (*BandFeedSequencerNotSettled, error) {
	out := new(BandFeedSequencerNotSettled)
	if err := bandFeed.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedT represents a T error raised by the BandFeed contract.
type BandFeedT struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error T()
func BandFeedTErrorID() common.Hash {
	return common.HexToHash("0x2bc80f3a55f2773020ab09d1e81403c4d8281a546b7c92d5fa07b84f2ffb15ea")
}

// UnpackTError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error T()
func (bandFeed *BandFeed) UnpackTError(raw []byte) (*BandFeedT, error) {
	out := new(BandFeedT)
	if err := bandFeed.abi.UnpackIntoInterface(out, "T", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// BandFeedUnknownAsset represents a UnknownAsset error raised by the BandFeed contract.
type BandFeedUnknownAsset struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func BandFeedUnknownAssetErrorID() common.Hash {
	return common.HexToHash("0x1059be3ed9f1c1f061f09edf206f56e7c80b74fd01671cbec1b3788a1f49417e")
}

// UnpackUnknownAssetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func (bandFeed *BandFeed) UnpackUnknownAssetError(raw []byte) (*BandFeedUnknownAsset, error) {
	out := new(BandFeedUnknownAsset)
	if err := bandFeed.abi.UnpackIntoInterface(out, "UnknownAsset", raw); err != nil {
		return nil, err
	}
	return out, nil
}
