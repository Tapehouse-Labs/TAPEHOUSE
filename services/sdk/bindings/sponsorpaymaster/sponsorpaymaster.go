// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sponsorpaymaster

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

// PackedUserOperation is an auto generated low-level Go binding around an user-defined struct.
type PackedUserOperation struct {
	Sender             common.Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

// SponsorPaymasterMetaData contains all meta data concerning the SponsorPaymaster contract.
var SponsorPaymasterMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"entryPoint_\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"maxCost\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FREE_OPERATIONS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"unstakeDelaySec\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"entryPoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEntryPoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freeOperationsLeft\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxCostPerOperation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operations\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOp\",\"inputs\":[{\"name\":\"mode\",\"type\":\"uint8\",\"internalType\":\"enumIPaymaster.PostOpMode\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actualGasCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"setMaxCost\",\"inputs\":[{\"name\":\"maxCost\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPaused\",\"inputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSigner\",\"inputs\":[{\"name\":\"signer_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"signer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePaymasterUserOp\",\"inputs\":[{\"name\":\"userOp\",\"type\":\"tuple\",\"internalType\":\"structPackedUserOperation\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initCode\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"accountGasLimits\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"preVerificationGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFees\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"paymasterAndData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"userOpHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"validationData\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawStake\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxCostSet\",\"inputs\":[{\"name\":\"maxCost\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PausedSet\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignerSet\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CostAboveLimit\",\"inputs\":[{\"name\":\"maxCost\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FreeOperationsUsed\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymasterUnauthorized\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SponsorshipPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]}]",
	ID:  "SponsorPaymaster",
}

// SponsorPaymaster is an auto generated Go binding around an Ethereum contract.
type SponsorPaymaster struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *SponsorPaymaster) GetABI() abi.ABI {
	return c.abi
}

// NewSponsorPaymaster creates a new instance of SponsorPaymaster.
func NewSponsorPaymaster() *SponsorPaymaster {
	parsed, err := SponsorPaymasterMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &SponsorPaymaster{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *SponsorPaymaster) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address entryPoint_, address owner_, address signer_, uint128 maxCost) returns()
func (sponsorPaymaster *SponsorPaymaster) PackConstructor(entryPoint_ common.Address, owner_ common.Address, signer_ common.Address, maxCost *big.Int) []byte {
	enc, err := sponsorPaymaster.abi.Pack("", entryPoint_, owner_, signer_, maxCost)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackFREEOPERATIONS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x227e98ee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FREE_OPERATIONS() view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) PackFREEOPERATIONS() []byte {
	enc, err := sponsorPaymaster.abi.Pack("FREE_OPERATIONS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFREEOPERATIONS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x227e98ee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FREE_OPERATIONS() view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) TryPackFREEOPERATIONS() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("FREE_OPERATIONS")
}

// UnpackFREEOPERATIONS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x227e98ee.
//
// Solidity: function FREE_OPERATIONS() view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) UnpackFREEOPERATIONS(data []byte) (*big.Int, error) {
	out, err := sponsorPaymaster.abi.Unpack("FREE_OPERATIONS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (sponsorPaymaster *SponsorPaymaster) PackAcceptOwnership() []byte {
	enc, err := sponsorPaymaster.abi.Pack("acceptOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptOwnership() returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackAcceptOwnership() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("acceptOwnership")
}

// PackAddStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0396cb60.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (sponsorPaymaster *SponsorPaymaster) PackAddStake(unstakeDelaySec uint32) []byte {
	enc, err := sponsorPaymaster.abi.Pack("addStake", unstakeDelaySec)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0396cb60.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackAddStake(unstakeDelaySec uint32) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("addStake", unstakeDelaySec)
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit() payable returns()
func (sponsorPaymaster *SponsorPaymaster) PackDeposit() []byte {
	enc, err := sponsorPaymaster.abi.Pack("deposit")
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
func (sponsorPaymaster *SponsorPaymaster) TryPackDeposit() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("deposit")
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (sponsorPaymaster *SponsorPaymaster) PackEip712Domain() []byte {
	enc, err := sponsorPaymaster.abi.Pack("eip712Domain")
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
func (sponsorPaymaster *SponsorPaymaster) TryPackEip712Domain() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("eip712Domain")
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
func (sponsorPaymaster *SponsorPaymaster) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := sponsorPaymaster.abi.Unpack("eip712Domain", data)
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

// PackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function entryPoint() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) PackEntryPoint() []byte {
	enc, err := sponsorPaymaster.abi.Pack("entryPoint")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEntryPoint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0d691fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function entryPoint() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) TryPackEntryPoint() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("entryPoint")
}

// UnpackEntryPoint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) UnpackEntryPoint(data []byte) (common.Address, error) {
	out, err := sponsorPaymaster.abi.Unpack("entryPoint", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackFreeOperationsLeft is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72f0c6b7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function freeOperationsLeft(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) PackFreeOperationsLeft(account common.Address) []byte {
	enc, err := sponsorPaymaster.abi.Pack("freeOperationsLeft", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFreeOperationsLeft is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72f0c6b7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function freeOperationsLeft(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) TryPackFreeOperationsLeft(account common.Address) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("freeOperationsLeft", account)
}

// UnpackFreeOperationsLeft is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x72f0c6b7.
//
// Solidity: function freeOperationsLeft(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) UnpackFreeOperationsLeft(data []byte) (*big.Int, error) {
	out, err := sponsorPaymaster.abi.Unpack("freeOperationsLeft", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxCostPerOperation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1a07dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxCostPerOperation() view returns(uint128)
func (sponsorPaymaster *SponsorPaymaster) PackMaxCostPerOperation() []byte {
	enc, err := sponsorPaymaster.abi.Pack("maxCostPerOperation")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxCostPerOperation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf1a07dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxCostPerOperation() view returns(uint128)
func (sponsorPaymaster *SponsorPaymaster) TryPackMaxCostPerOperation() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("maxCostPerOperation")
}

// UnpackMaxCostPerOperation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf1a07dc.
//
// Solidity: function maxCostPerOperation() view returns(uint128)
func (sponsorPaymaster *SponsorPaymaster) UnpackMaxCostPerOperation(data []byte) (*big.Int, error) {
	out, err := sponsorPaymaster.abi.Unpack("maxCostPerOperation", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOperations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb83c8902.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function operations(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) PackOperations(account common.Address) []byte {
	enc, err := sponsorPaymaster.abi.Pack("operations", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOperations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb83c8902.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function operations(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) TryPackOperations(account common.Address) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("operations", account)
}

// UnpackOperations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb83c8902.
//
// Solidity: function operations(address account) view returns(uint256)
func (sponsorPaymaster *SponsorPaymaster) UnpackOperations(data []byte) (*big.Int, error) {
	out, err := sponsorPaymaster.abi.Unpack("operations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) PackOwner() []byte {
	enc, err := sponsorPaymaster.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) TryPackOwner() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) UnpackOwner(data []byte) (common.Address, error) {
	out, err := sponsorPaymaster.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paused() view returns(bool)
func (sponsorPaymaster *SponsorPaymaster) PackPaused() []byte {
	enc, err := sponsorPaymaster.abi.Pack("paused")
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
func (sponsorPaymaster *SponsorPaymaster) TryPackPaused() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("paused")
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (sponsorPaymaster *SponsorPaymaster) UnpackPaused(data []byte) (bool, error) {
	out, err := sponsorPaymaster.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) PackPendingOwner() []byte {
	enc, err := sponsorPaymaster.abi.Pack("pendingOwner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingOwner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) TryPackPendingOwner() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := sponsorPaymaster.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPostOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c627b21.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (sponsorPaymaster *SponsorPaymaster) PackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) []byte {
	enc, err := sponsorPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPostOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7c627b21.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackPostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() pure returns()
func (sponsorPaymaster *SponsorPaymaster) PackRenounceOwnership() []byte {
	enc, err := sponsorPaymaster.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() pure returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackRenounceOwnership() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("renounceOwnership")
}

// PackSetMaxCost is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2bcbaf98.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setMaxCost(uint128 maxCost) returns()
func (sponsorPaymaster *SponsorPaymaster) PackSetMaxCost(maxCost *big.Int) []byte {
	enc, err := sponsorPaymaster.abi.Pack("setMaxCost", maxCost)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetMaxCost is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2bcbaf98.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setMaxCost(uint128 maxCost) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackSetMaxCost(maxCost *big.Int) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("setMaxCost", maxCost)
}

// PackSetPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16c38b3c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPaused(bool paused_) returns()
func (sponsorPaymaster *SponsorPaymaster) PackSetPaused(paused bool) []byte {
	enc, err := sponsorPaymaster.abi.Pack("setPaused", paused)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x16c38b3c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPaused(bool paused_) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackSetPaused(paused bool) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("setPaused", paused)
}

// PackSetSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c19e783.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setSigner(address signer_) returns()
func (sponsorPaymaster *SponsorPaymaster) PackSetSigner(signer common.Address) []byte {
	enc, err := sponsorPaymaster.abi.Pack("setSigner", signer)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6c19e783.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setSigner(address signer_) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackSetSigner(signer common.Address) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("setSigner", signer)
}

// PackSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x238ac933.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function signer() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) PackSigner() []byte {
	enc, err := sponsorPaymaster.abi.Pack("signer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x238ac933.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function signer() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) TryPackSigner() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("signer")
}

// UnpackSigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x238ac933.
//
// Solidity: function signer() view returns(address)
func (sponsorPaymaster *SponsorPaymaster) UnpackSigner(data []byte) (common.Address, error) {
	out, err := sponsorPaymaster.abi.Unpack("signer", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (sponsorPaymaster *SponsorPaymaster) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := sponsorPaymaster.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("transferOwnership", newOwner)
}

// PackUnlockStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9fe6bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unlockStake() returns()
func (sponsorPaymaster *SponsorPaymaster) PackUnlockStake() []byte {
	enc, err := sponsorPaymaster.abi.Pack("unlockStake")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnlockStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbb9fe6bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unlockStake() returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackUnlockStake() ([]byte, error) {
	return sponsorPaymaster.abi.Pack("unlockStake")
}

// PackValidatePaymasterUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b7512c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (sponsorPaymaster *SponsorPaymaster) PackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) []byte {
	enc, err := sponsorPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidatePaymasterUserOp is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b7512c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (sponsorPaymaster *SponsorPaymaster) TryPackValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOpOutput serves as a container for the return parameters of contract
// method ValidatePaymasterUserOp.
type ValidatePaymasterUserOpOutput struct {
	Context        []byte
	ValidationData *big.Int
}

// UnpackValidatePaymasterUserOp is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (sponsorPaymaster *SponsorPaymaster) UnpackValidatePaymasterUserOp(data []byte) (ValidatePaymasterUserOpOutput, error) {
	out, err := sponsorPaymaster.abi.Unpack("validatePaymasterUserOp", data)
	outstruct := new(ValidatePaymasterUserOpOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Context = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.ValidationData = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackWithdrawStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc23a5cea.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawStake(address to) returns()
func (sponsorPaymaster *SponsorPaymaster) PackWithdrawStake(to common.Address) []byte {
	enc, err := sponsorPaymaster.abi.Pack("withdrawStake", to)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawStake is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc23a5cea.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawStake(address to) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackWithdrawStake(to common.Address) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("withdrawStake", to)
}

// PackWithdrawTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205c2878.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawTo(address to, uint256 value) returns()
func (sponsorPaymaster *SponsorPaymaster) PackWithdrawTo(to common.Address, value *big.Int) []byte {
	enc, err := sponsorPaymaster.abi.Pack("withdrawTo", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x205c2878.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawTo(address to, uint256 value) returns()
func (sponsorPaymaster *SponsorPaymaster) TryPackWithdrawTo(to common.Address, value *big.Int) ([]byte, error) {
	return sponsorPaymaster.abi.Pack("withdrawTo", to, value)
}

// SponsorPaymasterEIP712DomainChanged represents a EIP712DomainChanged event raised by the SponsorPaymaster contract.
type SponsorPaymasterEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterEIP712DomainChanged) ContractEventName() string {
	return SponsorPaymasterEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (sponsorPaymaster *SponsorPaymaster) UnpackEIP712DomainChangedEvent(log *types.Log) (*SponsorPaymasterEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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

// SponsorPaymasterMaxCostSet represents a MaxCostSet event raised by the SponsorPaymaster contract.
type SponsorPaymasterMaxCostSet struct {
	MaxCost *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterMaxCostSetEventName = "MaxCostSet"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterMaxCostSet) ContractEventName() string {
	return SponsorPaymasterMaxCostSetEventName
}

// UnpackMaxCostSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event MaxCostSet(uint128 maxCost)
func (sponsorPaymaster *SponsorPaymaster) UnpackMaxCostSetEvent(log *types.Log) (*SponsorPaymasterMaxCostSet, error) {
	event := "MaxCostSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterMaxCostSet)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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

// SponsorPaymasterOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the SponsorPaymaster contract.
type SponsorPaymasterOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterOwnershipTransferStarted) ContractEventName() string {
	return SponsorPaymasterOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (sponsorPaymaster *SponsorPaymaster) UnpackOwnershipTransferStartedEvent(log *types.Log) (*SponsorPaymasterOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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

// SponsorPaymasterOwnershipTransferred represents a OwnershipTransferred event raised by the SponsorPaymaster contract.
type SponsorPaymasterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterOwnershipTransferred) ContractEventName() string {
	return SponsorPaymasterOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (sponsorPaymaster *SponsorPaymaster) UnpackOwnershipTransferredEvent(log *types.Log) (*SponsorPaymasterOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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

// SponsorPaymasterPausedSet represents a PausedSet event raised by the SponsorPaymaster contract.
type SponsorPaymasterPausedSet struct {
	Paused bool
	Raw    *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterPausedSetEventName = "PausedSet"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterPausedSet) ContractEventName() string {
	return SponsorPaymasterPausedSetEventName
}

// UnpackPausedSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PausedSet(bool paused)
func (sponsorPaymaster *SponsorPaymaster) UnpackPausedSetEvent(log *types.Log) (*SponsorPaymasterPausedSet, error) {
	event := "PausedSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterPausedSet)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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

// SponsorPaymasterSignerSet represents a SignerSet event raised by the SponsorPaymaster contract.
type SponsorPaymasterSignerSet struct {
	Signer common.Address
	Raw    *types.Log // Blockchain specific contextual infos
}

const SponsorPaymasterSignerSetEventName = "SignerSet"

// ContractEventName returns the user-defined event name.
func (SponsorPaymasterSignerSet) ContractEventName() string {
	return SponsorPaymasterSignerSetEventName
}

// UnpackSignerSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event SignerSet(address indexed signer)
func (sponsorPaymaster *SponsorPaymaster) UnpackSignerSetEvent(log *types.Log) (*SponsorPaymasterSignerSet, error) {
	event := "SignerSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != sponsorPaymaster.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(SponsorPaymasterSignerSet)
	if len(log.Data) > 0 {
		if err := sponsorPaymaster.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range sponsorPaymaster.abi.Events[event].Inputs {
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
func (sponsorPaymaster *SponsorPaymaster) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["CostAboveLimit"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackCostAboveLimitError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["FreeOperationsUsed"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackFreeOperationsUsedError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["InvalidShortString"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackInvalidShortStringError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["PaymasterUnauthorized"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackPaymasterUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["SponsorshipPaused"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackSponsorshipPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], sponsorPaymaster.abi.Errors["StringTooLong"].ID.Bytes()[:4]) {
		return sponsorPaymaster.UnpackStringTooLongError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// SponsorPaymasterCostAboveLimit represents a CostAboveLimit error raised by the SponsorPaymaster contract.
type SponsorPaymasterCostAboveLimit struct {
	MaxCost *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CostAboveLimit(uint256 maxCost)
func SponsorPaymasterCostAboveLimitErrorID() common.Hash {
	return common.HexToHash("0x72dc312a28bbe1d24db6439e22c40c6af41abf9b9610745b9dc643196b94e419")
}

// UnpackCostAboveLimitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CostAboveLimit(uint256 maxCost)
func (sponsorPaymaster *SponsorPaymaster) UnpackCostAboveLimitError(raw []byte) (*SponsorPaymasterCostAboveLimit, error) {
	out := new(SponsorPaymasterCostAboveLimit)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "CostAboveLimit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterFreeOperationsUsed represents a FreeOperationsUsed error raised by the SponsorPaymaster contract.
type SponsorPaymasterFreeOperationsUsed struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error FreeOperationsUsed(address account)
func SponsorPaymasterFreeOperationsUsedErrorID() common.Hash {
	return common.HexToHash("0xd106fef0649e9784701a750497739bcfeba00ad4adca8b3b7e6cc7a9a10702df")
}

// UnpackFreeOperationsUsedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error FreeOperationsUsed(address account)
func (sponsorPaymaster *SponsorPaymaster) UnpackFreeOperationsUsedError(raw []byte) (*SponsorPaymasterFreeOperationsUsed, error) {
	out := new(SponsorPaymasterFreeOperationsUsed)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "FreeOperationsUsed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterInvalidShortString represents a InvalidShortString error raised by the SponsorPaymaster contract.
type SponsorPaymasterInvalidShortString struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidShortString()
func SponsorPaymasterInvalidShortStringErrorID() common.Hash {
	return common.HexToHash("0xb3512b0c6163e5f0bafab72bb631b9d58cd7a731b082f910338aa21c83d5c274")
}

// UnpackInvalidShortStringError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidShortString()
func (sponsorPaymaster *SponsorPaymaster) UnpackInvalidShortStringError(raw []byte) (*SponsorPaymasterInvalidShortString, error) {
	out := new(SponsorPaymasterInvalidShortString)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "InvalidShortString", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the SponsorPaymaster contract.
type SponsorPaymasterOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func SponsorPaymasterOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (sponsorPaymaster *SponsorPaymaster) UnpackOwnableInvalidOwnerError(raw []byte) (*SponsorPaymasterOwnableInvalidOwner, error) {
	out := new(SponsorPaymasterOwnableInvalidOwner)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the SponsorPaymaster contract.
type SponsorPaymasterOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func SponsorPaymasterOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (sponsorPaymaster *SponsorPaymaster) UnpackOwnableUnauthorizedAccountError(raw []byte) (*SponsorPaymasterOwnableUnauthorizedAccount, error) {
	out := new(SponsorPaymasterOwnableUnauthorizedAccount)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the SponsorPaymaster contract.
type SponsorPaymasterOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func SponsorPaymasterOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (sponsorPaymaster *SponsorPaymaster) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*SponsorPaymasterOwnershipCannotBeRenounced, error) {
	out := new(SponsorPaymasterOwnershipCannotBeRenounced)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterPaymasterUnauthorized represents a PaymasterUnauthorized error raised by the SponsorPaymaster contract.
type SponsorPaymasterPaymasterUnauthorized struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymasterUnauthorized(address sender)
func SponsorPaymasterPaymasterUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0x327499a713dd0d80b45c032f882b325294f7d27e79628eeb33c5ffd95da1138d")
}

// UnpackPaymasterUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymasterUnauthorized(address sender)
func (sponsorPaymaster *SponsorPaymaster) UnpackPaymasterUnauthorizedError(raw []byte) (*SponsorPaymasterPaymasterUnauthorized, error) {
	out := new(SponsorPaymasterPaymasterUnauthorized)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "PaymasterUnauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterSponsorshipPaused represents a SponsorshipPaused error raised by the SponsorPaymaster contract.
type SponsorPaymasterSponsorshipPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SponsorshipPaused()
func SponsorPaymasterSponsorshipPausedErrorID() common.Hash {
	return common.HexToHash("0x36592babc99cc16ef01d8076f73a438d1f38e0e8c4e3684047b8d68be4072a70")
}

// UnpackSponsorshipPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SponsorshipPaused()
func (sponsorPaymaster *SponsorPaymaster) UnpackSponsorshipPausedError(raw []byte) (*SponsorPaymasterSponsorshipPaused, error) {
	out := new(SponsorPaymasterSponsorshipPaused)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "SponsorshipPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// SponsorPaymasterStringTooLong represents a StringTooLong error raised by the SponsorPaymaster contract.
type SponsorPaymasterStringTooLong struct {
	Str string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StringTooLong(string str)
func SponsorPaymasterStringTooLongErrorID() common.Hash {
	return common.HexToHash("0x305a27a93f8e33b7392df0a0f91d6fc63847395853c45991eec52dbf24d72381")
}

// UnpackStringTooLongError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StringTooLong(string str)
func (sponsorPaymaster *SponsorPaymaster) UnpackStringTooLongError(raw []byte) (*SponsorPaymasterStringTooLong, error) {
	out := new(SponsorPaymasterStringTooLong)
	if err := sponsorPaymaster.abi.UnpackIntoInterface(out, "StringTooLong", raw); err != nil {
		return nil, err
	}
	return out, nil
}
