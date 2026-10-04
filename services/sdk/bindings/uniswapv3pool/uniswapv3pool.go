// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswapv3pool

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

// UniswapV3PoolMetaData contains all meta data concerning the UniswapV3Pool contract.
var UniswapV3PoolMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collect\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"amount0Requested\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Requested\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collectProtocol\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount0Requested\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Requested\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeGrowthGlobal0X128\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeGrowthGlobal1X128\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"flash\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseObservationCardinalityNext\",\"inputs\":[{\"name\":\"observationCardinalityNext\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"sqrtPriceX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"liquidity\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxLiquidityPerTick\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"observations\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"blockTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"tickCumulative\",\"type\":\"int56\",\"internalType\":\"int56\"},{\"name\":\"secondsPerLiquidityCumulativeX128\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"initialized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"observe\",\"inputs\":[{\"name\":\"secondsAgos\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"outputs\":[{\"name\":\"tickCumulatives\",\"type\":\"int56[]\",\"internalType\":\"int56[]\"},{\"name\":\"secondsPerLiquidityCumulativeX128s\",\"type\":\"uint160[]\",\"internalType\":\"uint160[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positions\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"liquidity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"feeGrowthInside0LastX128\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeGrowthInside1LastX128\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokensOwed0\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"tokensOwed1\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFees\",\"inputs\":[],\"outputs\":[{\"name\":\"token0\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"token1\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setFeeProtocol\",\"inputs\":[{\"name\":\"feeProtocol0\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"feeProtocol1\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slot0\",\"inputs\":[],\"outputs\":[{\"name\":\"sqrtPriceX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"tick\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"observationIndex\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"observationCardinality\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"observationCardinalityNext\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"feeProtocol\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unlocked\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"snapshotCumulativesInside\",\"inputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"outputs\":[{\"name\":\"tickCumulativeInside\",\"type\":\"int56\",\"internalType\":\"int56\"},{\"name\":\"secondsPerLiquidityInsideX128\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"secondsInside\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swap\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"amountSpecified\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount0\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"amount1\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tickBitmap\",\"inputs\":[{\"name\":\"wordPosition\",\"type\":\"int16\",\"internalType\":\"int16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tickSpacing\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ticks\",\"inputs\":[{\"name\":\"tick\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"outputs\":[{\"name\":\"liquidityGross\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"liquidityNet\",\"type\":\"int128\",\"internalType\":\"int128\"},{\"name\":\"feeGrowthOutside0X128\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"feeGrowthOutside1X128\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tickCumulativeOutside\",\"type\":\"int56\",\"internalType\":\"int56\"},{\"name\":\"secondsPerLiquidityOutsideX128\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"secondsOutside\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"initialized\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"token0\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"token1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Burn\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"amount\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Collect\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"amount0\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"amount1\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollectProtocol\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"amount1\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Flash\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paid0\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paid1\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IncreaseObservationCardinalityNext\",\"inputs\":[{\"name\":\"observationCardinalityNextOld\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"observationCardinalityNextNew\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialize\",\"inputs\":[{\"name\":\"sqrtPriceX96\",\"type\":\"uint160\",\"indexed\":false,\"internalType\":\"uint160\"},{\"name\":\"tick\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Mint\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"indexed\":true,\"internalType\":\"int24\"},{\"name\":\"amount\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"amount0\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount1\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetFeeProtocol\",\"inputs\":[{\"name\":\"feeProtocol0Old\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"feeProtocol1Old\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"feeProtocol0New\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"feeProtocol1New\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Swap\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount0\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"amount1\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"sqrtPriceX96\",\"type\":\"uint160\",\"indexed\":false,\"internalType\":\"uint160\"},{\"name\":\"liquidity\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"tick\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AS\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"F0\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"F1\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IIA\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"L\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LOK\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"M0\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"M1\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TLM\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TLU\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TUM\",\"inputs\":[]}]",
	ID:  "UniswapV3Pool",
}

// UniswapV3Pool is an auto generated Go binding around an Ethereum contract.
type UniswapV3Pool struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *UniswapV3Pool) GetABI() abi.ABI {
	return c.abi
}

// NewUniswapV3Pool creates a new instance of UniswapV3Pool.
func NewUniswapV3Pool() *UniswapV3Pool {
	parsed, err := UniswapV3PoolMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &UniswapV3Pool{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *UniswapV3Pool) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa34123a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) PackBurn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("burn", tickLower, tickUpper, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa34123a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) TryPackBurn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("burn", tickLower, tickUpper, amount)
}

// BurnOutput serves as a container for the return parameters of contract
// method Burn.
type BurnOutput struct {
	Amount0 *big.Int
	Amount1 *big.Int
}

// UnpackBurn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackBurn(data []byte) (BurnOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("burn", data)
	outstruct := new(BurnOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Amount0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Amount1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCollect is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1eb3d8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) PackCollect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("collect", recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollect is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f1eb3d8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) TryPackCollect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("collect", recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// CollectOutput serves as a container for the return parameters of contract
// method Collect.
type CollectOutput struct {
	Amount0 *big.Int
	Amount1 *big.Int
}

// UnpackCollect is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackCollect(data []byte) (CollectOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("collect", data)
	outstruct := new(CollectOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Amount0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Amount1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCollectProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85b66729.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) PackCollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("collectProtocol", recipient, amount0Requested, amount1Requested)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollectProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85b66729.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) TryPackCollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("collectProtocol", recipient, amount0Requested, amount1Requested)
}

// CollectProtocolOutput serves as a container for the return parameters of contract
// method CollectProtocol.
type CollectProtocolOutput struct {
	Amount0 *big.Int
	Amount1 *big.Int
}

// UnpackCollectProtocol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackCollectProtocol(data []byte) (CollectProtocolOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("collectProtocol", data)
	outstruct := new(CollectProtocolOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Amount0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Amount1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackFactory is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc45a0155.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function factory() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) PackFactory() []byte {
	enc, err := uniswapV3Pool.abi.Pack("factory")
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
func (uniswapV3Pool *UniswapV3Pool) TryPackFactory() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("factory")
}

// UnpackFactory is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) UnpackFactory(data []byte) (common.Address, error) {
	out, err := uniswapV3Pool.abi.Unpack("factory", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xddca3f43.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function fee() view returns(uint24)
func (uniswapV3Pool *UniswapV3Pool) PackFee() []byte {
	enc, err := uniswapV3Pool.abi.Pack("fee")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xddca3f43.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function fee() view returns(uint24)
func (uniswapV3Pool *UniswapV3Pool) TryPackFee() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("fee")
}

// UnpackFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (uniswapV3Pool *UniswapV3Pool) UnpackFee(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("fee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFeeGrowthGlobal0X128 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3058399.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) PackFeeGrowthGlobal0X128() []byte {
	enc, err := uniswapV3Pool.abi.Pack("feeGrowthGlobal0X128")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeGrowthGlobal0X128 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3058399.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) TryPackFeeGrowthGlobal0X128() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("feeGrowthGlobal0X128")
}

// UnpackFeeGrowthGlobal0X128 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf3058399.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) UnpackFeeGrowthGlobal0X128(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("feeGrowthGlobal0X128", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFeeGrowthGlobal1X128 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x46141319.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) PackFeeGrowthGlobal1X128() []byte {
	enc, err := uniswapV3Pool.abi.Pack("feeGrowthGlobal1X128")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeGrowthGlobal1X128 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x46141319.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) TryPackFeeGrowthGlobal1X128() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("feeGrowthGlobal1X128")
}

// UnpackFeeGrowthGlobal1X128 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x46141319.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) UnpackFeeGrowthGlobal1X128(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("feeGrowthGlobal1X128", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackFlash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x490e6cbc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (uniswapV3Pool *UniswapV3Pool) PackFlash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) []byte {
	enc, err := uniswapV3Pool.abi.Pack("flash", recipient, amount0, amount1, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFlash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x490e6cbc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (uniswapV3Pool *UniswapV3Pool) TryPackFlash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("flash", recipient, amount0, amount1, data)
}

// PackIncreaseObservationCardinalityNext is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32148f67.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (uniswapV3Pool *UniswapV3Pool) PackIncreaseObservationCardinalityNext(observationCardinalityNext uint16) []byte {
	enc, err := uniswapV3Pool.abi.Pack("increaseObservationCardinalityNext", observationCardinalityNext)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIncreaseObservationCardinalityNext is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x32148f67.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (uniswapV3Pool *UniswapV3Pool) TryPackIncreaseObservationCardinalityNext(observationCardinalityNext uint16) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("increaseObservationCardinalityNext", observationCardinalityNext)
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf637731d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (uniswapV3Pool *UniswapV3Pool) PackInitialize(sqrtPriceX96 *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("initialize", sqrtPriceX96)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf637731d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (uniswapV3Pool *UniswapV3Pool) TryPackInitialize(sqrtPriceX96 *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("initialize", sqrtPriceX96)
}

// PackLiquidity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a686502.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidity() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) PackLiquidity() []byte {
	enc, err := uniswapV3Pool.abi.Pack("liquidity")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLiquidity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1a686502.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function liquidity() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) TryPackLiquidity() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("liquidity")
}

// UnpackLiquidity is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1a686502.
//
// Solidity: function liquidity() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) UnpackLiquidity(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("liquidity", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxLiquidityPerTick is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70cf754a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) PackMaxLiquidityPerTick() []byte {
	enc, err := uniswapV3Pool.abi.Pack("maxLiquidityPerTick")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxLiquidityPerTick is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70cf754a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) TryPackMaxLiquidityPerTick() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("maxLiquidityPerTick")
}

// UnpackMaxLiquidityPerTick is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (uniswapV3Pool *UniswapV3Pool) UnpackMaxLiquidityPerTick(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("maxLiquidityPerTick", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3c8a7d8d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) PackMint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) []byte {
	enc, err := uniswapV3Pool.abi.Pack("mint", recipient, tickLower, tickUpper, amount, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3c8a7d8d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) TryPackMint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("mint", recipient, tickLower, tickUpper, amount, data)
}

// MintOutput serves as a container for the return parameters of contract
// method Mint.
type MintOutput struct {
	Amount0 *big.Int
	Amount1 *big.Int
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackMint(data []byte) (MintOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("mint", data)
	outstruct := new(MintOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Amount0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Amount1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackObservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x252c09d7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function observations(uint256 index) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) PackObservations(index *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("observations", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackObservations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x252c09d7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function observations(uint256 index) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) TryPackObservations(index *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("observations", index)
}

// ObservationsOutput serves as a container for the return parameters of contract
// method Observations.
type ObservationsOutput struct {
	BlockTimestamp                    uint32
	TickCumulative                    *big.Int
	SecondsPerLiquidityCumulativeX128 *big.Int
	Initialized                       bool
}

// UnpackObservations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x252c09d7.
//
// Solidity: function observations(uint256 index) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) UnpackObservations(data []byte) (ObservationsOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("observations", data)
	outstruct := new(ObservationsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.BlockTimestamp = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.TickCumulative = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.SecondsPerLiquidityCumulativeX128 = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.Initialized = *abi.ConvertType(out[3], new(bool)).(*bool)
	return *outstruct, nil
}

// PackObserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x883bdbfd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (uniswapV3Pool *UniswapV3Pool) PackObserve(secondsAgos []uint32) []byte {
	enc, err := uniswapV3Pool.abi.Pack("observe", secondsAgos)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackObserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x883bdbfd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (uniswapV3Pool *UniswapV3Pool) TryPackObserve(secondsAgos []uint32) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("observe", secondsAgos)
}

// ObserveOutput serves as a container for the return parameters of contract
// method Observe.
type ObserveOutput struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}

// UnpackObserve is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (uniswapV3Pool *UniswapV3Pool) UnpackObserve(data []byte) (ObserveOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("observe", data)
	outstruct := new(ObserveOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.TickCumulatives = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.SecondsPerLiquidityCumulativeX128s = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackPositions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x514ea4bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function positions(bytes32 key) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (uniswapV3Pool *UniswapV3Pool) PackPositions(key [32]byte) []byte {
	enc, err := uniswapV3Pool.abi.Pack("positions", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPositions is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x514ea4bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function positions(bytes32 key) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (uniswapV3Pool *UniswapV3Pool) TryPackPositions(key [32]byte) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("positions", key)
}

// PositionsOutput serves as a container for the return parameters of contract
// method Positions.
type PositionsOutput struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	TokensOwed0              *big.Int
	TokensOwed1              *big.Int
}

// UnpackPositions is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x514ea4bf.
//
// Solidity: function positions(bytes32 key) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (uniswapV3Pool *UniswapV3Pool) UnpackPositions(data []byte) (PositionsOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("positions", data)
	outstruct := new(PositionsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Liquidity = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.FeeGrowthInside0LastX128 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.FeeGrowthInside1LastX128 = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.TokensOwed0 = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.TokensOwed1 = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackProtocolFees is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ad8b03b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (uniswapV3Pool *UniswapV3Pool) PackProtocolFees() []byte {
	enc, err := uniswapV3Pool.abi.Pack("protocolFees")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackProtocolFees is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ad8b03b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (uniswapV3Pool *UniswapV3Pool) TryPackProtocolFees() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("protocolFees")
}

// ProtocolFeesOutput serves as a container for the return parameters of contract
// method ProtocolFees.
type ProtocolFeesOutput struct {
	Token0 *big.Int
	Token1 *big.Int
}

// UnpackProtocolFees is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ad8b03b.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (uniswapV3Pool *UniswapV3Pool) UnpackProtocolFees(data []byte) (ProtocolFeesOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("protocolFees", data)
	outstruct := new(ProtocolFeesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Token0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Token1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackSetFeeProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8206a4d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (uniswapV3Pool *UniswapV3Pool) PackSetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) []byte {
	enc, err := uniswapV3Pool.abi.Pack("setFeeProtocol", feeProtocol0, feeProtocol1)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetFeeProtocol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8206a4d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (uniswapV3Pool *UniswapV3Pool) TryPackSetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("setFeeProtocol", feeProtocol0, feeProtocol1)
}

// PackSlot0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3850c7bd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (uniswapV3Pool *UniswapV3Pool) PackSlot0() []byte {
	enc, err := uniswapV3Pool.abi.Pack("slot0")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSlot0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3850c7bd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (uniswapV3Pool *UniswapV3Pool) TryPackSlot0() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("slot0")
}

// Slot0Output serves as a container for the return parameters of contract
// method Slot0.
type Slot0Output struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

// UnpackSlot0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3850c7bd.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (uniswapV3Pool *UniswapV3Pool) UnpackSlot0(data []byte) (Slot0Output, error) {
	out, err := uniswapV3Pool.abi.Unpack("slot0", data)
	outstruct := new(Slot0Output)
	if err != nil {
		return *outstruct, err
	}
	outstruct.SqrtPriceX96 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Tick = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.ObservationIndex = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.ObservationCardinality = *abi.ConvertType(out[3], new(uint16)).(*uint16)
	outstruct.ObservationCardinalityNext = *abi.ConvertType(out[4], new(uint16)).(*uint16)
	outstruct.FeeProtocol = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Unlocked = *abi.ConvertType(out[6], new(bool)).(*bool)
	return *outstruct, nil
}

// PackSnapshotCumulativesInside is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa38807f2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (uniswapV3Pool *UniswapV3Pool) PackSnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("snapshotCumulativesInside", tickLower, tickUpper)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSnapshotCumulativesInside is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa38807f2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (uniswapV3Pool *UniswapV3Pool) TryPackSnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("snapshotCumulativesInside", tickLower, tickUpper)
}

// SnapshotCumulativesInsideOutput serves as a container for the return parameters of contract
// method SnapshotCumulativesInside.
type SnapshotCumulativesInsideOutput struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}

// UnpackSnapshotCumulativesInside is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (uniswapV3Pool *UniswapV3Pool) UnpackSnapshotCumulativesInside(data []byte) (SnapshotCumulativesInsideOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("snapshotCumulativesInside", data)
	outstruct := new(SnapshotCumulativesInsideOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.TickCumulativeInside = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.SecondsPerLiquidityInsideX128 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.SecondsInside = *abi.ConvertType(out[2], new(uint32)).(*uint32)
	return *outstruct, nil
}

// PackSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x128acb08.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (uniswapV3Pool *UniswapV3Pool) PackSwap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) []byte {
	enc, err := uniswapV3Pool.abi.Pack("swap", recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSwap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x128acb08.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (uniswapV3Pool *UniswapV3Pool) TryPackSwap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("swap", recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// SwapOutput serves as a container for the return parameters of contract
// method Swap.
type SwapOutput struct {
	Amount0 *big.Int
	Amount1 *big.Int
}

// UnpackSwap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackSwap(data []byte) (SwapOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("swap", data)
	outstruct := new(SwapOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Amount0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Amount1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackTickBitmap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5339c296.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tickBitmap(int16 wordPosition) view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) PackTickBitmap(wordPosition int16) []byte {
	enc, err := uniswapV3Pool.abi.Pack("tickBitmap", wordPosition)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTickBitmap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5339c296.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tickBitmap(int16 wordPosition) view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) TryPackTickBitmap(wordPosition int16) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("tickBitmap", wordPosition)
}

// UnpackTickBitmap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5339c296.
//
// Solidity: function tickBitmap(int16 wordPosition) view returns(uint256)
func (uniswapV3Pool *UniswapV3Pool) UnpackTickBitmap(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("tickBitmap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTickSpacing is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0c93a7c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tickSpacing() view returns(int24)
func (uniswapV3Pool *UniswapV3Pool) PackTickSpacing() []byte {
	enc, err := uniswapV3Pool.abi.Pack("tickSpacing")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTickSpacing is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0c93a7c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tickSpacing() view returns(int24)
func (uniswapV3Pool *UniswapV3Pool) TryPackTickSpacing() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("tickSpacing")
}

// UnpackTickSpacing is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (uniswapV3Pool *UniswapV3Pool) UnpackTickSpacing(data []byte) (*big.Int, error) {
	out, err := uniswapV3Pool.abi.Unpack("tickSpacing", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTicks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf30dba93.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ticks(int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) PackTicks(tick *big.Int) []byte {
	enc, err := uniswapV3Pool.abi.Pack("ticks", tick)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTicks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf30dba93.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ticks(int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) TryPackTicks(tick *big.Int) ([]byte, error) {
	return uniswapV3Pool.abi.Pack("ticks", tick)
}

// TicksOutput serves as a container for the return parameters of contract
// method Ticks.
type TicksOutput struct {
	LiquidityGross                 *big.Int
	LiquidityNet                   *big.Int
	FeeGrowthOutside0X128          *big.Int
	FeeGrowthOutside1X128          *big.Int
	TickCumulativeOutside          *big.Int
	SecondsPerLiquidityOutsideX128 *big.Int
	SecondsOutside                 uint32
	Initialized                    bool
}

// UnpackTicks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf30dba93.
//
// Solidity: function ticks(int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (uniswapV3Pool *UniswapV3Pool) UnpackTicks(data []byte) (TicksOutput, error) {
	out, err := uniswapV3Pool.abi.Unpack("ticks", data)
	outstruct := new(TicksOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.LiquidityGross = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.LiquidityNet = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.FeeGrowthOutside0X128 = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.FeeGrowthOutside1X128 = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.TickCumulativeOutside = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	outstruct.SecondsPerLiquidityOutsideX128 = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	outstruct.SecondsOutside = *abi.ConvertType(out[6], new(uint32)).(*uint32)
	outstruct.Initialized = *abi.ConvertType(out[7], new(bool)).(*bool)
	return *outstruct, nil
}

// PackToken0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dfe1681.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function token0() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) PackToken0() []byte {
	enc, err := uniswapV3Pool.abi.Pack("token0")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackToken0 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dfe1681.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function token0() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) TryPackToken0() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("token0")
}

// UnpackToken0 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) UnpackToken0(data []byte) (common.Address, error) {
	out, err := uniswapV3Pool.abi.Unpack("token0", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackToken1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd21220a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function token1() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) PackToken1() []byte {
	enc, err := uniswapV3Pool.abi.Pack("token1")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackToken1 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd21220a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function token1() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) TryPackToken1() ([]byte, error) {
	return uniswapV3Pool.abi.Pack("token1")
}

// UnpackToken1 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (uniswapV3Pool *UniswapV3Pool) UnpackToken1(data []byte) (common.Address, error) {
	out, err := uniswapV3Pool.abi.Unpack("token1", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// UniswapV3PoolBurn represents a Burn event raised by the UniswapV3Pool contract.
type UniswapV3PoolBurn struct {
	Owner     common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount    *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolBurnEventName = "Burn"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolBurn) ContractEventName() string {
	return UniswapV3PoolBurnEventName
}

// UnpackBurnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Burn(address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackBurnEvent(log *types.Log) (*UniswapV3PoolBurn, error) {
	event := "Burn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolBurn)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolCollect represents a Collect event raised by the UniswapV3Pool contract.
type UniswapV3PoolCollect struct {
	Owner     common.Address
	Recipient common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolCollectEventName = "Collect"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolCollect) ContractEventName() string {
	return UniswapV3PoolCollectEventName
}

// UnpackCollectEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Collect(address indexed owner, address recipient, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackCollectEvent(log *types.Log) (*UniswapV3PoolCollect, error) {
	event := "Collect"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolCollect)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolCollectProtocol represents a CollectProtocol event raised by the UniswapV3Pool contract.
type UniswapV3PoolCollectProtocol struct {
	Sender    common.Address
	Recipient common.Address
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolCollectProtocolEventName = "CollectProtocol"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolCollectProtocol) ContractEventName() string {
	return UniswapV3PoolCollectProtocolEventName
}

// UnpackCollectProtocolEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackCollectProtocolEvent(log *types.Log) (*UniswapV3PoolCollectProtocol, error) {
	event := "CollectProtocol"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolCollectProtocol)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolFlash represents a Flash event raised by the UniswapV3Pool contract.
type UniswapV3PoolFlash struct {
	Sender    common.Address
	Recipient common.Address
	Amount0   *big.Int
	Amount1   *big.Int
	Paid0     *big.Int
	Paid1     *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolFlashEventName = "Flash"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolFlash) ContractEventName() string {
	return UniswapV3PoolFlashEventName
}

// UnpackFlashEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Flash(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1, uint256 paid0, uint256 paid1)
func (uniswapV3Pool *UniswapV3Pool) UnpackFlashEvent(log *types.Log) (*UniswapV3PoolFlash, error) {
	event := "Flash"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolFlash)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolIncreaseObservationCardinalityNext represents a IncreaseObservationCardinalityNext event raised by the UniswapV3Pool contract.
type UniswapV3PoolIncreaseObservationCardinalityNext struct {
	ObservationCardinalityNextOld uint16
	ObservationCardinalityNextNew uint16
	Raw                           *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolIncreaseObservationCardinalityNextEventName = "IncreaseObservationCardinalityNext"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolIncreaseObservationCardinalityNext) ContractEventName() string {
	return UniswapV3PoolIncreaseObservationCardinalityNextEventName
}

// UnpackIncreaseObservationCardinalityNextEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event IncreaseObservationCardinalityNext(uint16 observationCardinalityNextOld, uint16 observationCardinalityNextNew)
func (uniswapV3Pool *UniswapV3Pool) UnpackIncreaseObservationCardinalityNextEvent(log *types.Log) (*UniswapV3PoolIncreaseObservationCardinalityNext, error) {
	event := "IncreaseObservationCardinalityNext"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolIncreaseObservationCardinalityNext)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolInitialize represents a Initialize event raised by the UniswapV3Pool contract.
type UniswapV3PoolInitialize struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolInitializeEventName = "Initialize"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolInitialize) ContractEventName() string {
	return UniswapV3PoolInitializeEventName
}

// UnpackInitializeEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialize(uint160 sqrtPriceX96, int24 tick)
func (uniswapV3Pool *UniswapV3Pool) UnpackInitializeEvent(log *types.Log) (*UniswapV3PoolInitialize, error) {
	event := "Initialize"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolInitialize)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolMint represents a Mint event raised by the UniswapV3Pool contract.
type UniswapV3PoolMint struct {
	Sender    common.Address
	Owner     common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount    *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolMintEventName = "Mint"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolMint) ContractEventName() string {
	return UniswapV3PoolMintEventName
}

// UnpackMintEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Mint(address sender, address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (uniswapV3Pool *UniswapV3Pool) UnpackMintEvent(log *types.Log) (*UniswapV3PoolMint, error) {
	event := "Mint"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolMint)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolSetFeeProtocol represents a SetFeeProtocol event raised by the UniswapV3Pool contract.
type UniswapV3PoolSetFeeProtocol struct {
	FeeProtocol0Old uint8
	FeeProtocol1Old uint8
	FeeProtocol0New uint8
	FeeProtocol1New uint8
	Raw             *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolSetFeeProtocolEventName = "SetFeeProtocol"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolSetFeeProtocol) ContractEventName() string {
	return UniswapV3PoolSetFeeProtocolEventName
}

// UnpackSetFeeProtocolEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New)
func (uniswapV3Pool *UniswapV3Pool) UnpackSetFeeProtocolEvent(log *types.Log) (*UniswapV3PoolSetFeeProtocol, error) {
	event := "SetFeeProtocol"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolSetFeeProtocol)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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

// UniswapV3PoolSwap represents a Swap event raised by the UniswapV3Pool contract.
type UniswapV3PoolSwap struct {
	Sender       common.Address
	Recipient    common.Address
	Amount0      *big.Int
	Amount1      *big.Int
	SqrtPriceX96 *big.Int
	Liquidity    *big.Int
	Tick         *big.Int
	Raw          *types.Log // Blockchain specific contextual infos
}

const UniswapV3PoolSwapEventName = "Swap"

// ContractEventName returns the user-defined event name.
func (UniswapV3PoolSwap) ContractEventName() string {
	return UniswapV3PoolSwapEventName
}

// UnpackSwapEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, int256 amount0, int256 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick)
func (uniswapV3Pool *UniswapV3Pool) UnpackSwapEvent(log *types.Log) (*UniswapV3PoolSwap, error) {
	event := "Swap"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != uniswapV3Pool.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(UniswapV3PoolSwap)
	if len(log.Data) > 0 {
		if err := uniswapV3Pool.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range uniswapV3Pool.abi.Events[event].Inputs {
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
func (uniswapV3Pool *UniswapV3Pool) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["AI"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackAIError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["AS"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackASError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["F0"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackF0Error(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["F1"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackF1Error(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["IIA"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackIIAError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["L"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackLError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["LOK"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackLOKError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["M0"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackM0Error(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["M1"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackM1Error(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["TLM"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackTLMError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["TLU"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackTLUError(raw[4:])
	}
	if bytes.Equal(raw[:4], uniswapV3Pool.abi.Errors["TUM"].ID.Bytes()[:4]) {
		return uniswapV3Pool.UnpackTUMError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// UniswapV3PoolAI represents a AI error raised by the UniswapV3Pool contract.
type UniswapV3PoolAI struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AI()
func UniswapV3PoolAIErrorID() common.Hash {
	return common.HexToHash("0x9cc0b7f81b0afdb90924ab4b40f96fb372a95da401bcd28ba95b721201664d16")
}

// UnpackAIError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AI()
func (uniswapV3Pool *UniswapV3Pool) UnpackAIError(raw []byte) (*UniswapV3PoolAI, error) {
	out := new(UniswapV3PoolAI)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "AI", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolAS represents a AS error raised by the UniswapV3Pool contract.
type UniswapV3PoolAS struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AS()
func UniswapV3PoolASErrorID() common.Hash {
	return common.HexToHash("0x03fff01897c9401082548001aea637ebe0e887627eacab27a854c4a0712c8d07")
}

// UnpackASError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AS()
func (uniswapV3Pool *UniswapV3Pool) UnpackASError(raw []byte) (*UniswapV3PoolAS, error) {
	out := new(UniswapV3PoolAS)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "AS", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolF0 represents a F0 error raised by the UniswapV3Pool contract.
type UniswapV3PoolF0 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error F0()
func UniswapV3PoolF0ErrorID() common.Hash {
	return common.HexToHash("0xf704e899c3d1537aa14e9f29fa77d04c1320459daede40ceff62a92076fd730a")
}

// UnpackF0Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error F0()
func (uniswapV3Pool *UniswapV3Pool) UnpackF0Error(raw []byte) (*UniswapV3PoolF0, error) {
	out := new(UniswapV3PoolF0)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "F0", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolF1 represents a F1 error raised by the UniswapV3Pool contract.
type UniswapV3PoolF1 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error F1()
func UniswapV3PoolF1ErrorID() common.Hash {
	return common.HexToHash("0xe90c349375555573daf3f3020e969556f3b86bb96d9257cf5744e93c2b99583d")
}

// UnpackF1Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error F1()
func (uniswapV3Pool *UniswapV3Pool) UnpackF1Error(raw []byte) (*UniswapV3PoolF1, error) {
	out := new(UniswapV3PoolF1)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "F1", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolIIA represents a IIA error raised by the UniswapV3Pool contract.
type UniswapV3PoolIIA struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IIA()
func UniswapV3PoolIIAErrorID() common.Hash {
	return common.HexToHash("0xba0b951e20f417a790f16c6192db00c389f27764ad7e46eb05bc67d2da424f5a")
}

// UnpackIIAError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IIA()
func (uniswapV3Pool *UniswapV3Pool) UnpackIIAError(raw []byte) (*UniswapV3PoolIIA, error) {
	out := new(UniswapV3PoolIIA)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "IIA", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolL represents a L error raised by the UniswapV3Pool contract.
type UniswapV3PoolL struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error L()
func UniswapV3PoolLErrorID() common.Hash {
	return common.HexToHash("0x9f13f76d2663f02e97ae485fda00006d267bd7aeeee1b717496fe4c46722c096")
}

// UnpackLError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error L()
func (uniswapV3Pool *UniswapV3Pool) UnpackLError(raw []byte) (*UniswapV3PoolL, error) {
	out := new(UniswapV3PoolL)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "L", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolLOK represents a LOK error raised by the UniswapV3Pool contract.
type UniswapV3PoolLOK struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LOK()
func UniswapV3PoolLOKErrorID() common.Hash {
	return common.HexToHash("0xa1bf78866bf823d0e6d330c8edac6a3fea6e06b3c47d5e9c51f2248057440c27")
}

// UnpackLOKError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LOK()
func (uniswapV3Pool *UniswapV3Pool) UnpackLOKError(raw []byte) (*UniswapV3PoolLOK, error) {
	out := new(UniswapV3PoolLOK)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "LOK", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolM0 represents a M0 error raised by the UniswapV3Pool contract.
type UniswapV3PoolM0 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error M0()
func UniswapV3PoolM0ErrorID() common.Hash {
	return common.HexToHash("0x748800af3a112bc2a48eac181c2c1ed4d378534a24ea90abc411ef667927f771")
}

// UnpackM0Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error M0()
func (uniswapV3Pool *UniswapV3Pool) UnpackM0Error(raw []byte) (*UniswapV3PoolM0, error) {
	out := new(UniswapV3PoolM0)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "M0", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolM1 represents a M1 error raised by the UniswapV3Pool contract.
type UniswapV3PoolM1 struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error M1()
func UniswapV3PoolM1ErrorID() common.Hash {
	return common.HexToHash("0x20e5672e2db58b29042af8b09d9525ac50290a75db4062b7ec7068b6d35daa6a")
}

// UnpackM1Error is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error M1()
func (uniswapV3Pool *UniswapV3Pool) UnpackM1Error(raw []byte) (*UniswapV3PoolM1, error) {
	out := new(UniswapV3PoolM1)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "M1", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolTLM represents a TLM error raised by the UniswapV3Pool contract.
type UniswapV3PoolTLM struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TLM()
func UniswapV3PoolTLMErrorID() common.Hash {
	return common.HexToHash("0x9ad612e8af7b588417a500664ab71a113d2ac3a720ea05f93d3ab76db208ad3a")
}

// UnpackTLMError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TLM()
func (uniswapV3Pool *UniswapV3Pool) UnpackTLMError(raw []byte) (*UniswapV3PoolTLM, error) {
	out := new(UniswapV3PoolTLM)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "TLM", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolTLU represents a TLU error raised by the UniswapV3Pool contract.
type UniswapV3PoolTLU struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TLU()
func UniswapV3PoolTLUErrorID() common.Hash {
	return common.HexToHash("0x2fe0284ff5085ced1982e27ec4cb26013bc9957bc51f79b227bb3100c4930dc4")
}

// UnpackTLUError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TLU()
func (uniswapV3Pool *UniswapV3Pool) UnpackTLUError(raw []byte) (*UniswapV3PoolTLU, error) {
	out := new(UniswapV3PoolTLU)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "TLU", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// UniswapV3PoolTUM represents a TUM error raised by the UniswapV3Pool contract.
type UniswapV3PoolTUM struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TUM()
func UniswapV3PoolTUMErrorID() common.Hash {
	return common.HexToHash("0xd7b54ab1b19d28cb459c4498f085576a19cfafbfdc2e657807931167e4bb8623")
}

// UnpackTUMError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TUM()
func (uniswapV3Pool *UniswapV3Pool) UnpackTUMError(raw []byte) (*UniswapV3PoolTUM, error) {
	out := new(UniswapV3PoolTUM)
	if err := uniswapV3Pool.abi.UnpackIntoInterface(out, "TUM", raw); err != nil {
		return nil, err
	}
	return out, nil
}
