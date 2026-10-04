// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package weth

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

// WethMetaData contains all meta data concerning the Weth contract.
var WethMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositTo\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	ID:  "Weth",
}

// Weth is an auto generated Go binding around an Ethereum contract.
type Weth struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Weth) GetABI() abi.ABI {
	return c.abi
}

// NewWeth creates a new instance of Weth.
func NewWeth() *Weth {
	parsed, err := WethMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Weth{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Weth) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDOMAINSEPARATOR is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3644e515.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (weth *Weth) PackDOMAINSEPARATOR() []byte {
	enc, err := weth.abi.Pack("DOMAIN_SEPARATOR")
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
func (weth *Weth) TryPackDOMAINSEPARATOR() ([]byte, error) {
	return weth.abi.Pack("DOMAIN_SEPARATOR")
}

// UnpackDOMAINSEPARATOR is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (weth *Weth) UnpackDOMAINSEPARATOR(data []byte) ([32]byte, error) {
	out, err := weth.abi.Unpack("DOMAIN_SEPARATOR", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (weth *Weth) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := weth.abi.Pack("allowance", owner, spender)
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
func (weth *Weth) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return weth.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (weth *Weth) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := weth.abi.Unpack("allowance", data)
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
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (weth *Weth) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := weth.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (weth *Weth) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return weth.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (weth *Weth) UnpackApprove(data []byte) (bool, error) {
	out, err := weth.abi.Unpack("approve", data)
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
func (weth *Weth) PackBalanceOf(account common.Address) []byte {
	enc, err := weth.abi.Pack("balanceOf", account)
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
func (weth *Weth) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return weth.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (weth *Weth) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := weth.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit() payable returns()
func (weth *Weth) PackDeposit() []byte {
	enc, err := weth.abi.Pack("deposit")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit() payable returns()
func (weth *Weth) TryPackDeposit() ([]byte, error) {
	return weth.abi.Pack("deposit")
}

// PackDepositTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb760faf9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositTo(address account) payable returns()
func (weth *Weth) PackDepositTo(account common.Address) []byte {
	enc, err := weth.abi.Pack("depositTo", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb760faf9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositTo(address account) payable returns()
func (weth *Weth) TryPackDepositTo(account common.Address) ([]byte, error) {
	return weth.abi.Pack("depositTo", account)
}

// PackNonces is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7ecebe00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (weth *Weth) PackNonces(owner common.Address) []byte {
	enc, err := weth.abi.Pack("nonces", owner)
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
func (weth *Weth) TryPackNonces(owner common.Address) ([]byte, error) {
	return weth.abi.Pack("nonces", owner)
}

// UnpackNonces is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (weth *Weth) UnpackNonces(data []byte) (*big.Int, error) {
	out, err := weth.abi.Unpack("nonces", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd505accf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (weth *Weth) PackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) []byte {
	enc, err := weth.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
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
func (weth *Weth) TryPackPermit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) ([]byte, error) {
	return weth.abi.Pack("permit", owner, spender, value, deadline, v, r, s)
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (weth *Weth) PackTotalSupply() []byte {
	enc, err := weth.abi.Pack("totalSupply")
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
func (weth *Weth) TryPackTotalSupply() ([]byte, error) {
	return weth.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (weth *Weth) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := weth.abi.Unpack("totalSupply", data)
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
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (weth *Weth) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := weth.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (weth *Weth) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return weth.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (weth *Weth) UnpackTransfer(data []byte) (bool, error) {
	out, err := weth.abi.Unpack("transfer", data)
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
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (weth *Weth) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := weth.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (weth *Weth) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return weth.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (weth *Weth) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := weth.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 amount) returns()
func (weth *Weth) PackWithdraw(amount *big.Int) []byte {
	enc, err := weth.abi.Pack("withdraw", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 amount) returns()
func (weth *Weth) TryPackWithdraw(amount *big.Int) ([]byte, error) {
	return weth.abi.Pack("withdraw", amount)
}

// WethApproval represents a Approval event raised by the Weth contract.
type WethApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const WethApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (WethApproval) ContractEventName() string {
	return WethApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (weth *Weth) UnpackApprovalEvent(log *types.Log) (*WethApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != weth.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WethApproval)
	if len(log.Data) > 0 {
		if err := weth.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range weth.abi.Events[event].Inputs {
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

// WethTransfer represents a Transfer event raised by the Weth contract.
type WethTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const WethTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (WethTransfer) ContractEventName() string {
	return WethTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (weth *Weth) UnpackTransferEvent(log *types.Log) (*WethTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != weth.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(WethTransfer)
	if len(log.Data) > 0 {
		if err := weth.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range weth.abi.Events[event].Inputs {
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
