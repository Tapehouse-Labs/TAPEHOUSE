// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stocktoken

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

// StockTokenMetaData contains all meta data concerning the StockToken contract.
var StockTokenMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ACCESS_CONTROLLED_REGISTRY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"adminBurn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOfUI\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"effectiveAt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"newUIMultiplier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oraclePaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseOracle\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupplyUI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uiMultiplier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpauseOracle\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMultiplier\",\"inputs\":[{\"name\":\"newMultiplier\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"effectiveAt_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OraclePaused\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleUnpaused\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferWithScaledUI\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"uiValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UIMultiplierUpdated\",\"inputs\":[{\"name\":\"oldMultiplier\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newMultiplier\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"effectiveAtTimestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Blocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"IsPaused\",\"inputs\":[]}]",
	ID:  "StockToken",
}

// StockToken is an auto generated Go binding around an Ethereum contract.
type StockToken struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *StockToken) GetABI() abi.ABI {
	return c.abi
}

// NewStockToken creates a new instance of StockToken.
func NewStockToken() *StockToken {
	parsed, err := StockTokenMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &StockToken{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *StockToken) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackACCESSCONTROLLEDREGISTRY is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50c09be3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ACCESS_CONTROLLED_REGISTRY() view returns(address)
func (stockToken *StockToken) PackACCESSCONTROLLEDREGISTRY() []byte {
	enc, err := stockToken.abi.Pack("ACCESS_CONTROLLED_REGISTRY")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackACCESSCONTROLLEDREGISTRY is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x50c09be3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ACCESS_CONTROLLED_REGISTRY() view returns(address)
func (stockToken *StockToken) TryPackACCESSCONTROLLEDREGISTRY() ([]byte, error) {
	return stockToken.abi.Pack("ACCESS_CONTROLLED_REGISTRY")
}

// UnpackACCESSCONTROLLEDREGISTRY is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x50c09be3.
//
// Solidity: function ACCESS_CONTROLLED_REGISTRY() view returns(address)
func (stockToken *StockToken) UnpackACCESSCONTROLLEDREGISTRY(data []byte) (common.Address, error) {
	out, err := stockToken.abi.Unpack("ACCESS_CONTROLLED_REGISTRY", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDOMAINSEPARATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3644e515.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (stockToken *StockToken) PackDOMAINSEPARATOR() []byte {
	enc, err := stockToken.abi.Pack("DOMAIN_SEPARATOR")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDOMAINSEPARATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3644e515.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (stockToken *StockToken) TryPackDOMAINSEPARATOR() ([]byte, error) {
	return stockToken.abi.Pack("DOMAIN_SEPARATOR")
}

// UnpackDOMAINSEPARATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (stockToken *StockToken) UnpackDOMAINSEPARATOR(data []byte) ([32]byte, error) {
	out, err := stockToken.abi.Unpack("DOMAIN_SEPARATOR", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAdminBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06dd0419.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function adminBurn(address from, uint256 amount) returns()
func (stockToken *StockToken) PackAdminBurn(from common.Address, amount *big.Int) []byte {
	enc, err := stockToken.abi.Pack("adminBurn", from, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAdminBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06dd0419.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function adminBurn(address from, uint256 amount) returns()
func (stockToken *StockToken) TryPackAdminBurn(from common.Address, amount *big.Int) ([]byte, error) {
	return stockToken.abi.Pack("adminBurn", from, amount)
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (stockToken *StockToken) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := stockToken.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (stockToken *StockToken) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return stockToken.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (stockToken *StockToken) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (stockToken *StockToken) PackApprove(spender common.Address, amount *big.Int) []byte {
	enc, err := stockToken.abi.Pack("approve", spender, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (stockToken *StockToken) TryPackApprove(spender common.Address, amount *big.Int) ([]byte, error) {
	return stockToken.abi.Pack("approve", spender, amount)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (stockToken *StockToken) UnpackApprove(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (stockToken *StockToken) PackBalanceOf(account common.Address) []byte {
	enc, err := stockToken.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (stockToken *StockToken) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return stockToken.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (stockToken *StockToken) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBalanceOfUI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x437a9958.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOfUI(address account) view returns(uint256)
func (stockToken *StockToken) PackBalanceOfUI(account common.Address) []byte {
	enc, err := stockToken.abi.Pack("balanceOfUI", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOfUI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x437a9958.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOfUI(address account) view returns(uint256)
func (stockToken *StockToken) TryPackBalanceOfUI(account common.Address) ([]byte, error) {
	return stockToken.abi.Pack("balanceOfUI", account)
}

// UnpackBalanceOfUI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x437a9958.
//
// Solidity: function balanceOfUI(address account) view returns(uint256)
func (stockToken *StockToken) UnpackBalanceOfUI(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("balanceOfUI", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (stockToken *StockToken) PackDecimals() []byte {
	enc, err := stockToken.abi.Pack("decimals")
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
func (stockToken *StockToken) TryPackDecimals() ([]byte, error) {
	return stockToken.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (stockToken *StockToken) UnpackDecimals(data []byte) (uint8, error) {
	out, err := stockToken.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackEffectiveAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97a4064f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function effectiveAt() view returns(uint256)
func (stockToken *StockToken) PackEffectiveAt() []byte {
	enc, err := stockToken.abi.Pack("effectiveAt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEffectiveAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97a4064f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function effectiveAt() view returns(uint256)
func (stockToken *StockToken) TryPackEffectiveAt() ([]byte, error) {
	return stockToken.abi.Pack("effectiveAt")
}

// UnpackEffectiveAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x97a4064f.
//
// Solidity: function effectiveAt() view returns(uint256)
func (stockToken *StockToken) UnpackEffectiveAt(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("effectiveAt", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (stockToken *StockToken) PackEip712Domain() []byte {
	enc, err := stockToken.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (stockToken *StockToken) TryPackEip712Domain() ([]byte, error) {
	return stockToken.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (stockToken *StockToken) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := stockToken.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (stockToken *StockToken) PackName() []byte {
	enc, err := stockToken.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (stockToken *StockToken) TryPackName() ([]byte, error) {
	return stockToken.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (stockToken *StockToken) UnpackName(data []byte) (string, error) {
	out, err := stockToken.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackNewUIMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc767007.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function newUIMultiplier() view returns(uint256)
func (stockToken *StockToken) PackNewUIMultiplier() []byte {
	enc, err := stockToken.abi.Pack("newUIMultiplier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNewUIMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdc767007.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function newUIMultiplier() view returns(uint256)
func (stockToken *StockToken) TryPackNewUIMultiplier() ([]byte, error) {
	return stockToken.abi.Pack("newUIMultiplier")
}

// UnpackNewUIMultiplier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdc767007.
//
// Solidity: function newUIMultiplier() view returns(uint256)
func (stockToken *StockToken) UnpackNewUIMultiplier(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("newUIMultiplier", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ecebe00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (stockToken *StockToken) PackNonces(owner common.Address) []byte {
	enc, err := stockToken.abi.Pack("nonces", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ecebe00.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (stockToken *StockToken) TryPackNonces(owner common.Address) ([]byte, error) {
	return stockToken.abi.Pack("nonces", owner)
}

// UnpackNonces is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (stockToken *StockToken) UnpackNonces(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("nonces", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOraclePaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7706ba52.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function oraclePaused() view returns(bool)
func (stockToken *StockToken) PackOraclePaused() []byte {
	enc, err := stockToken.abi.Pack("oraclePaused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOraclePaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7706ba52.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function oraclePaused() view returns(bool)
func (stockToken *StockToken) TryPackOraclePaused() ([]byte, error) {
	return stockToken.abi.Pack("oraclePaused")
}

// UnpackOraclePaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7706ba52.
//
// Solidity: function oraclePaused() view returns(bool)
func (stockToken *StockToken) UnpackOraclePaused(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("oraclePaused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8456cb59.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pause() returns()
func (stockToken *StockToken) PackPause() []byte {
	enc, err := stockToken.abi.Pack("pause")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8456cb59.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pause() returns()
func (stockToken *StockToken) TryPackPause() ([]byte, error) {
	return stockToken.abi.Pack("pause")
}

// PackPauseOracle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x253ea980.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseOracle() returns()
func (stockToken *StockToken) PackPauseOracle() []byte {
	enc, err := stockToken.abi.Pack("pauseOracle")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseOracle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x253ea980.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseOracle() returns()
func (stockToken *StockToken) TryPackPauseOracle() ([]byte, error) {
	return stockToken.abi.Pack("pauseOracle")
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paused() view returns(bool)
func (stockToken *StockToken) PackPaused() []byte {
	enc, err := stockToken.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function paused() view returns(bool)
func (stockToken *StockToken) TryPackPaused() ([]byte, error) {
	return stockToken.abi.Pack("paused")
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (stockToken *StockToken) UnpackPaused(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd505accf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (stockToken *StockToken) PackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) []byte {
	enc, err := stockToken.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd505accf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (stockToken *StockToken) TryPackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) ([]byte, error) {
	return stockToken.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (stockToken *StockToken) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := stockToken.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (stockToken *StockToken) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return stockToken.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (stockToken *StockToken) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (stockToken *StockToken) PackSymbol() []byte {
	enc, err := stockToken.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (stockToken *StockToken) TryPackSymbol() ([]byte, error) {
	return stockToken.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (stockToken *StockToken) UnpackSymbol(data []byte) (string, error) {
	out, err := stockToken.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTokenPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86c75e74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tokenPaused() view returns(bool)
func (stockToken *StockToken) PackTokenPaused() []byte {
	enc, err := stockToken.abi.Pack("tokenPaused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTokenPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x86c75e74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tokenPaused() view returns(bool)
func (stockToken *StockToken) TryPackTokenPaused() ([]byte, error) {
	return stockToken.abi.Pack("tokenPaused")
}

// UnpackTokenPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x86c75e74.
//
// Solidity: function tokenPaused() view returns(bool)
func (stockToken *StockToken) UnpackTokenPaused(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("tokenPaused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (stockToken *StockToken) PackTotalSupply() []byte {
	enc, err := stockToken.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (stockToken *StockToken) TryPackTotalSupply() ([]byte, error) {
	return stockToken.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (stockToken *StockToken) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalSupplyUI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9bea6429.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupplyUI() view returns(uint256)
func (stockToken *StockToken) PackTotalSupplyUI() []byte {
	enc, err := stockToken.abi.Pack("totalSupplyUI")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupplyUI is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9bea6429.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupplyUI() view returns(uint256)
func (stockToken *StockToken) TryPackTotalSupplyUI() ([]byte, error) {
	return stockToken.abi.Pack("totalSupplyUI")
}

// UnpackTotalSupplyUI is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9bea6429.
//
// Solidity: function totalSupplyUI() view returns(uint256)
func (stockToken *StockToken) UnpackTotalSupplyUI(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("totalSupplyUI", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (stockToken *StockToken) PackTransfer(to common.Address, amount *big.Int) []byte {
	enc, err := stockToken.abi.Pack("transfer", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (stockToken *StockToken) TryPackTransfer(to common.Address, amount *big.Int) ([]byte, error) {
	return stockToken.abi.Pack("transfer", to, amount)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (stockToken *StockToken) UnpackTransfer(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (stockToken *StockToken) PackTransferFrom(from common.Address, to common.Address, amount *big.Int) []byte {
	enc, err := stockToken.abi.Pack("transferFrom", from, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (stockToken *StockToken) TryPackTransferFrom(from common.Address, to common.Address, amount *big.Int) ([]byte, error) {
	return stockToken.abi.Pack("transferFrom", from, to, amount)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (stockToken *StockToken) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := stockToken.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackUiMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa60bf13d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function uiMultiplier() view returns(uint256)
func (stockToken *StockToken) PackUiMultiplier() []byte {
	enc, err := stockToken.abi.Pack("uiMultiplier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUiMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa60bf13d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function uiMultiplier() view returns(uint256)
func (stockToken *StockToken) TryPackUiMultiplier() ([]byte, error) {
	return stockToken.abi.Pack("uiMultiplier")
}

// UnpackUiMultiplier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa60bf13d.
//
// Solidity: function uiMultiplier() view returns(uint256)
func (stockToken *StockToken) UnpackUiMultiplier(data []byte) (*big.Int, error) {
	out, err := stockToken.abi.Unpack("uiMultiplier", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpause() returns()
func (stockToken *StockToken) PackUnpause() []byte {
	enc, err := stockToken.abi.Pack("unpause")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unpause() returns()
func (stockToken *StockToken) TryPackUnpause() ([]byte, error) {
	return stockToken.abi.Pack("unpause")
}

// PackUnpauseOracle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fab6865.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpauseOracle() returns()
func (stockToken *StockToken) PackUnpauseOracle() []byte {
	enc, err := stockToken.abi.Pack("unpauseOracle")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnpauseOracle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fab6865.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unpauseOracle() returns()
func (stockToken *StockToken) TryPackUnpauseOracle() ([]byte, error) {
	return stockToken.abi.Pack("unpauseOracle")
}

// PackUpdateMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbad60f18.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateMultiplier(uint256 newMultiplier, uint256 effectiveAt_) returns()
func (stockToken *StockToken) PackUpdateMultiplier(newMultiplier *big.Int, effectiveAt *big.Int) []byte {
	enc, err := stockToken.abi.Pack("updateMultiplier", newMultiplier, effectiveAt)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateMultiplier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbad60f18.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateMultiplier(uint256 newMultiplier, uint256 effectiveAt_) returns()
func (stockToken *StockToken) TryPackUpdateMultiplier(newMultiplier *big.Int, effectiveAt *big.Int) ([]byte, error) {
	return stockToken.abi.Pack("updateMultiplier", newMultiplier, effectiveAt)
}

// StockTokenApproval represents a Approval event raised by the StockToken contract.
type StockTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const StockTokenApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (StockTokenApproval) ContractEventName() string {
	return StockTokenApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (stockToken *StockToken) UnpackApprovalEvent(log *types.Log) (*StockTokenApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenApproval)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenOraclePaused represents a OraclePaused event raised by the StockToken contract.
type StockTokenOraclePaused struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StockTokenOraclePausedEventName = "OraclePaused"

// ContractEventName returns the user-defined event name.
func (StockTokenOraclePaused) ContractEventName() string {
	return StockTokenOraclePausedEventName
}

// UnpackOraclePausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OraclePaused()
func (stockToken *StockToken) UnpackOraclePausedEvent(log *types.Log) (*StockTokenOraclePaused, error) {
	event := "OraclePaused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenOraclePaused)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenOracleUnpaused represents a OracleUnpaused event raised by the StockToken contract.
type StockTokenOracleUnpaused struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StockTokenOracleUnpausedEventName = "OracleUnpaused"

// ContractEventName returns the user-defined event name.
func (StockTokenOracleUnpaused) ContractEventName() string {
	return StockTokenOracleUnpausedEventName
}

// UnpackOracleUnpausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OracleUnpaused()
func (stockToken *StockToken) UnpackOracleUnpausedEvent(log *types.Log) (*StockTokenOracleUnpaused, error) {
	event := "OracleUnpaused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenOracleUnpaused)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenPaused represents a Paused event raised by the StockToken contract.
type StockTokenPaused struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StockTokenPausedEventName = "Paused"

// ContractEventName returns the user-defined event name.
func (StockTokenPaused) ContractEventName() string {
	return StockTokenPausedEventName
}

// UnpackPausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Paused()
func (stockToken *StockToken) UnpackPausedEvent(log *types.Log) (*StockTokenPaused, error) {
	event := "Paused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenPaused)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenTransfer represents a Transfer event raised by the StockToken contract.
type StockTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const StockTokenTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (StockTokenTransfer) ContractEventName() string {
	return StockTokenTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (stockToken *StockToken) UnpackTransferEvent(log *types.Log) (*StockTokenTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenTransfer)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenTransferWithScaledUI represents a TransferWithScaledUI event raised by the StockToken contract.
type StockTokenTransferWithScaledUI struct {
	From    common.Address
	To      common.Address
	Value   *big.Int
	UiValue *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const StockTokenTransferWithScaledUIEventName = "TransferWithScaledUI"

// ContractEventName returns the user-defined event name.
func (StockTokenTransferWithScaledUI) ContractEventName() string {
	return StockTokenTransferWithScaledUIEventName
}

// UnpackTransferWithScaledUIEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TransferWithScaledUI(address indexed from, address indexed to, uint256 value, uint256 uiValue)
func (stockToken *StockToken) UnpackTransferWithScaledUIEvent(log *types.Log) (*StockTokenTransferWithScaledUI, error) {
	event := "TransferWithScaledUI"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenTransferWithScaledUI)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenUIMultiplierUpdated represents a UIMultiplierUpdated event raised by the StockToken contract.
type StockTokenUIMultiplierUpdated struct {
	OldMultiplier        *big.Int
	NewMultiplier        *big.Int
	EffectiveAtTimestamp *big.Int
	Raw                  *types.Log // Blockchain specific contextual infos
}

const StockTokenUIMultiplierUpdatedEventName = "UIMultiplierUpdated"

// ContractEventName returns the user-defined event name.
func (StockTokenUIMultiplierUpdated) ContractEventName() string {
	return StockTokenUIMultiplierUpdatedEventName
}

// UnpackUIMultiplierUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event UIMultiplierUpdated(uint256 oldMultiplier, uint256 newMultiplier, uint256 effectiveAtTimestamp)
func (stockToken *StockToken) UnpackUIMultiplierUpdatedEvent(log *types.Log) (*StockTokenUIMultiplierUpdated, error) {
	event := "UIMultiplierUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenUIMultiplierUpdated)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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

// StockTokenUnpaused represents a Unpaused event raised by the StockToken contract.
type StockTokenUnpaused struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const StockTokenUnpausedEventName = "Unpaused"

// ContractEventName returns the user-defined event name.
func (StockTokenUnpaused) ContractEventName() string {
	return StockTokenUnpausedEventName
}

// UnpackUnpausedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Unpaused()
func (stockToken *StockToken) UnpackUnpausedEvent(log *types.Log) (*StockTokenUnpaused, error) {
	event := "Unpaused"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockTokenUnpaused)
	if len(log.Data) > 0 {
		if err := stockToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockToken.abi.Events[event].Inputs {
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
func (stockToken *StockToken) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], stockToken.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return stockToken.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockToken.abi.Errors["Blocked"].ID.Bytes()[:4]) {
		return stockToken.UnpackBlockedError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockToken.abi.Errors["IsPaused"].ID.Bytes()[:4]) {
		return stockToken.UnpackIsPausedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// StockTokenAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the StockToken contract.
type StockTokenAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func StockTokenAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (stockToken *StockToken) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*StockTokenAccessControlUnauthorizedAccount, error) {
	out := new(StockTokenAccessControlUnauthorizedAccount)
	if err := stockToken.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockTokenBlocked represents a Blocked error raised by the StockToken contract.
type StockTokenBlocked struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Blocked(address account)
func StockTokenBlockedErrorID() common.Hash {
	return common.HexToHash("0x75e91ce73c1d3352d8dd3610443539cd33dfe13b1de8f8caae54ec26dd0dc9cb")
}

// UnpackBlockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Blocked(address account)
func (stockToken *StockToken) UnpackBlockedError(raw []byte) (*StockTokenBlocked, error) {
	out := new(StockTokenBlocked)
	if err := stockToken.abi.UnpackIntoInterface(out, "Blocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockTokenIsPaused represents a IsPaused error raised by the StockToken contract.
type StockTokenIsPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IsPaused()
func StockTokenIsPausedErrorID() common.Hash {
	return common.HexToHash("0x1309a5639731e6b91c4aa19095f23eee07602ba240a393e2b5ccbf71a8e2952f")
}

// UnpackIsPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IsPaused()
func (stockToken *StockToken) UnpackIsPausedError(raw []byte) (*StockTokenIsPaused, error) {
	out := new(StockTokenIsPaused)
	if err := stockToken.abi.UnpackIntoInterface(out, "IsPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}
