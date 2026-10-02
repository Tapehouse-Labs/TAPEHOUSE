// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package morphoblue

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

// IMorphoMarketParams is an auto generated low-level Go binding around an user-defined struct.
type IMorphoMarketParams struct {
	LoanToken       common.Address
	CollateralToken common.Address
	Oracle          common.Address
	Irm             common.Address
	Lltv            *big.Int
}

// MorphoBlueMetaData contains all meta data concerning the MorphoBlue contract.
var MorphoBlueMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"accrueInterest\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createMarket\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"idToMarketParams\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isIrmEnabled\",\"inputs\":[{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isLltvEnabled\",\"inputs\":[{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidate\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"borrower\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"seizedAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"repaidShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"market\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"totalSupplyAssets\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalSupplyShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalBorrowAssets\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"totalBorrowShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdate\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"fee\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"position\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"supplyShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"borrowShares\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"collateral\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supply\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyCollateral\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"marketParams\",\"type\":\"tuple\",\"internalType\":\"structIMorpho.MarketParams\",\"components\":[{\"name\":\"loanToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"irm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lltv\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	ID:  "MorphoBlue",
}

// MorphoBlue is an auto generated Go binding around an Ethereum contract.
type MorphoBlue struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *MorphoBlue) GetABI() abi.ABI {
	return c.abi
}

// NewMorphoBlue creates a new instance of MorphoBlue.
func NewMorphoBlue() *MorphoBlue {
	parsed, err := MorphoBlueMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MorphoBlue{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MorphoBlue) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAccrueInterest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x151c1ade.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function accrueInterest((address,address,address,address,uint256) marketParams) returns()
func (morphoBlue *MorphoBlue) PackAccrueInterest(marketParams IMorphoMarketParams) []byte {
	enc, err := morphoBlue.abi.Pack("accrueInterest", marketParams)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAccrueInterest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x151c1ade.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function accrueInterest((address,address,address,address,uint256) marketParams) returns()
func (morphoBlue *MorphoBlue) TryPackAccrueInterest(marketParams IMorphoMarketParams) ([]byte, error) {
	return morphoBlue.abi.Pack("accrueInterest", marketParams)
}

// PackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50d8cd4b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrow((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) PackBorrow(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, receiver common.Address) []byte {
	enc, err := morphoBlue.abi.Pack("borrow", marketParams, assets, shares, onBehalf, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50d8cd4b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrow((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) TryPackBorrow(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, receiver common.Address) ([]byte, error) {
	return morphoBlue.abi.Pack("borrow", marketParams, assets, shares, onBehalf, receiver)
}

// BorrowOutput serves as a container for the return parameters of contract
// method Borrow.
type BorrowOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackBorrow is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x50d8cd4b.
//
// Solidity: function borrow((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) UnpackBorrow(data []byte) (BorrowOutput, error) {
	out, err := morphoBlue.abi.Unpack("borrow", data)
	outstruct := new(BorrowOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCreateMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c1358a2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function createMarket((address,address,address,address,uint256) marketParams) returns()
func (morphoBlue *MorphoBlue) PackCreateMarket(marketParams IMorphoMarketParams) []byte {
	enc, err := morphoBlue.abi.Pack("createMarket", marketParams)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCreateMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8c1358a2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function createMarket((address,address,address,address,uint256) marketParams) returns()
func (morphoBlue *MorphoBlue) TryPackCreateMarket(marketParams IMorphoMarketParams) ([]byte, error) {
	return morphoBlue.abi.Pack("createMarket", marketParams)
}

// PackIdToMarketParams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c3c9157.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function idToMarketParams(bytes32 id) view returns((address,address,address,address,uint256))
func (morphoBlue *MorphoBlue) PackIdToMarketParams(id [32]byte) []byte {
	enc, err := morphoBlue.abi.Pack("idToMarketParams", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIdToMarketParams is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2c3c9157.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function idToMarketParams(bytes32 id) view returns((address,address,address,address,uint256))
func (morphoBlue *MorphoBlue) TryPackIdToMarketParams(id [32]byte) ([]byte, error) {
	return morphoBlue.abi.Pack("idToMarketParams", id)
}

// UnpackIdToMarketParams is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2c3c9157.
//
// Solidity: function idToMarketParams(bytes32 id) view returns((address,address,address,address,uint256))
func (morphoBlue *MorphoBlue) UnpackIdToMarketParams(data []byte) (IMorphoMarketParams, error) {
	out, err := morphoBlue.abi.Unpack("idToMarketParams", data)
	if err != nil {
		return *new(IMorphoMarketParams), err
	}
	out0 := *abi.ConvertType(out[0], new(IMorphoMarketParams)).(*IMorphoMarketParams)
	return out0, nil
}

// PackIsIrmEnabled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2b863ce.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isIrmEnabled(address irm) view returns(bool)
func (morphoBlue *MorphoBlue) PackIsIrmEnabled(irm common.Address) []byte {
	enc, err := morphoBlue.abi.Pack("isIrmEnabled", irm)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsIrmEnabled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2b863ce.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isIrmEnabled(address irm) view returns(bool)
func (morphoBlue *MorphoBlue) TryPackIsIrmEnabled(irm common.Address) ([]byte, error) {
	return morphoBlue.abi.Pack("isIrmEnabled", irm)
}

// UnpackIsIrmEnabled is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf2b863ce.
//
// Solidity: function isIrmEnabled(address irm) view returns(bool)
func (morphoBlue *MorphoBlue) UnpackIsIrmEnabled(data []byte) (bool, error) {
	out, err := morphoBlue.abi.Unpack("isIrmEnabled", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsLltvEnabled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb485f3b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isLltvEnabled(uint256 lltv) view returns(bool)
func (morphoBlue *MorphoBlue) PackIsLltvEnabled(lltv *big.Int) []byte {
	enc, err := morphoBlue.abi.Pack("isLltvEnabled", lltv)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsLltvEnabled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb485f3b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isLltvEnabled(uint256 lltv) view returns(bool)
func (morphoBlue *MorphoBlue) TryPackIsLltvEnabled(lltv *big.Int) ([]byte, error) {
	return morphoBlue.abi.Pack("isLltvEnabled", lltv)
}

// UnpackIsLltvEnabled is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb485f3b8.
//
// Solidity: function isLltvEnabled(uint256 lltv) view returns(bool)
func (morphoBlue *MorphoBlue) UnpackIsLltvEnabled(data []byte) (bool, error) {
	out, err := morphoBlue.abi.Unpack("isLltvEnabled", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackLiquidate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8eabcb8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidate((address,address,address,address,uint256) marketParams, address borrower, uint256 seizedAssets, uint256 repaidShares, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) PackLiquidate(marketParams IMorphoMarketParams, borrower common.Address, seizedAssets *big.Int, repaidShares *big.Int, data []byte) []byte {
	enc, err := morphoBlue.abi.Pack("liquidate", marketParams, borrower, seizedAssets, repaidShares, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLiquidate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8eabcb8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function liquidate((address,address,address,address,uint256) marketParams, address borrower, uint256 seizedAssets, uint256 repaidShares, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) TryPackLiquidate(marketParams IMorphoMarketParams, borrower common.Address, seizedAssets *big.Int, repaidShares *big.Int, data []byte) ([]byte, error) {
	return morphoBlue.abi.Pack("liquidate", marketParams, borrower, seizedAssets, repaidShares, data)
}

// LiquidateOutput serves as a container for the return parameters of contract
// method Liquidate.
type LiquidateOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackLiquidate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd8eabcb8.
//
// Solidity: function liquidate((address,address,address,address,uint256) marketParams, address borrower, uint256 seizedAssets, uint256 repaidShares, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) UnpackLiquidate(data []byte) (LiquidateOutput, error) {
	out, err := morphoBlue.abi.Unpack("liquidate", data)
	outstruct := new(LiquidateOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c60e39a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function market(bytes32 id) view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint128 lastUpdate, uint128 fee)
func (morphoBlue *MorphoBlue) PackMarket(id [32]byte) []byte {
	enc, err := morphoBlue.abi.Pack("market", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMarket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c60e39a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function market(bytes32 id) view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint128 lastUpdate, uint128 fee)
func (morphoBlue *MorphoBlue) TryPackMarket(id [32]byte) ([]byte, error) {
	return morphoBlue.abi.Pack("market", id)
}

// MarketOutput serves as a container for the return parameters of contract
// method Market.
type MarketOutput struct {
	TotalSupplyAssets *big.Int
	TotalSupplyShares *big.Int
	TotalBorrowAssets *big.Int
	TotalBorrowShares *big.Int
	LastUpdate        *big.Int
	Fee               *big.Int
}

// UnpackMarket is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c60e39a.
//
// Solidity: function market(bytes32 id) view returns(uint128 totalSupplyAssets, uint128 totalSupplyShares, uint128 totalBorrowAssets, uint128 totalBorrowShares, uint128 lastUpdate, uint128 fee)
func (morphoBlue *MorphoBlue) UnpackMarket(data []byte) (MarketOutput, error) {
	out, err := morphoBlue.abi.Unpack("market", data)
	outstruct := new(MarketOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.TotalSupplyAssets = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.TotalSupplyShares = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.TotalBorrowAssets = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.TotalBorrowShares = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.LastUpdate = abi.ConvertType(out[4], new(big.Int)).(*big.Int)
	outstruct.Fee = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackPosition is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93c52062.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function position(bytes32 id, address user) view returns(uint256 supplyShares, uint128 borrowShares, uint128 collateral)
func (morphoBlue *MorphoBlue) PackPosition(id [32]byte, user common.Address) []byte {
	enc, err := morphoBlue.abi.Pack("position", id, user)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPosition is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93c52062.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function position(bytes32 id, address user) view returns(uint256 supplyShares, uint128 borrowShares, uint128 collateral)
func (morphoBlue *MorphoBlue) TryPackPosition(id [32]byte, user common.Address) ([]byte, error) {
	return morphoBlue.abi.Pack("position", id, user)
}

// PositionOutput serves as a container for the return parameters of contract
// method Position.
type PositionOutput struct {
	SupplyShares *big.Int
	BorrowShares *big.Int
	Collateral   *big.Int
}

// UnpackPosition is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x93c52062.
//
// Solidity: function position(bytes32 id, address user) view returns(uint256 supplyShares, uint128 borrowShares, uint128 collateral)
func (morphoBlue *MorphoBlue) UnpackPosition(data []byte) (PositionOutput, error) {
	out, err := morphoBlue.abi.Unpack("position", data)
	outstruct := new(PositionOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.SupplyShares = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.BorrowShares = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Collateral = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x20b76e81.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function repay((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) PackRepay(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, data []byte) []byte {
	enc, err := morphoBlue.abi.Pack("repay", marketParams, assets, shares, onBehalf, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x20b76e81.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function repay((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) TryPackRepay(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, data []byte) ([]byte, error) {
	return morphoBlue.abi.Pack("repay", marketParams, assets, shares, onBehalf, data)
}

// RepayOutput serves as a container for the return parameters of contract
// method Repay.
type RepayOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackRepay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x20b76e81.
//
// Solidity: function repay((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) UnpackRepay(data []byte) (RepayOutput, error) {
	out, err := morphoBlue.abi.Unpack("repay", data)
	outstruct := new(RepayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa99aad89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supply((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) PackSupply(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, data []byte) []byte {
	enc, err := morphoBlue.abi.Pack("supply", marketParams, assets, shares, onBehalf, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa99aad89.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supply((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) TryPackSupply(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, data []byte) ([]byte, error) {
	return morphoBlue.abi.Pack("supply", marketParams, assets, shares, onBehalf, data)
}

// SupplyOutput serves as a container for the return parameters of contract
// method Supply.
type SupplyOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa99aad89.
//
// Solidity: function supply((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, bytes data) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) UnpackSupply(data []byte) (SupplyOutput, error) {
	out, err := morphoBlue.abi.Unpack("supply", data)
	outstruct := new(SupplyOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackSupplyCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x238d6579.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supplyCollateral((address,address,address,address,uint256) marketParams, uint256 assets, address onBehalf, bytes data) returns()
func (morphoBlue *MorphoBlue) PackSupplyCollateral(marketParams IMorphoMarketParams, assets *big.Int, onBehalf common.Address, data []byte) []byte {
	enc, err := morphoBlue.abi.Pack("supplyCollateral", marketParams, assets, onBehalf, data)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupplyCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x238d6579.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supplyCollateral((address,address,address,address,uint256) marketParams, uint256 assets, address onBehalf, bytes data) returns()
func (morphoBlue *MorphoBlue) TryPackSupplyCollateral(marketParams IMorphoMarketParams, assets *big.Int, onBehalf common.Address, data []byte) ([]byte, error) {
	return morphoBlue.abi.Pack("supplyCollateral", marketParams, assets, onBehalf, data)
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c2bea49.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) PackWithdraw(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, receiver common.Address) []byte {
	enc, err := morphoBlue.abi.Pack("withdraw", marketParams, assets, shares, onBehalf, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c2bea49.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) TryPackWithdraw(marketParams IMorphoMarketParams, assets *big.Int, shares *big.Int, onBehalf common.Address, receiver common.Address) ([]byte, error) {
	return morphoBlue.abi.Pack("withdraw", marketParams, assets, shares, onBehalf, receiver)
}

// WithdrawOutput serves as a container for the return parameters of contract
// method Withdraw.
type WithdrawOutput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c2bea49.
//
// Solidity: function withdraw((address,address,address,address,uint256) marketParams, uint256 assets, uint256 shares, address onBehalf, address receiver) returns(uint256, uint256)
func (morphoBlue *MorphoBlue) UnpackWithdraw(data []byte) (WithdrawOutput, error) {
	out, err := morphoBlue.abi.Unpack("withdraw", data)
	outstruct := new(WithdrawOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Arg1 = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackWithdrawCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8720316d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawCollateral((address,address,address,address,uint256) marketParams, uint256 assets, address onBehalf, address receiver) returns()
func (morphoBlue *MorphoBlue) PackWithdrawCollateral(marketParams IMorphoMarketParams, assets *big.Int, onBehalf common.Address, receiver common.Address) []byte {
	enc, err := morphoBlue.abi.Pack("withdrawCollateral", marketParams, assets, onBehalf, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8720316d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawCollateral((address,address,address,address,uint256) marketParams, uint256 assets, address onBehalf, address receiver) returns()
func (morphoBlue *MorphoBlue) TryPackWithdrawCollateral(marketParams IMorphoMarketParams, assets *big.Int, onBehalf common.Address, receiver common.Address) ([]byte, error) {
	return morphoBlue.abi.Pack("withdrawCollateral", marketParams, assets, onBehalf, receiver)
}
