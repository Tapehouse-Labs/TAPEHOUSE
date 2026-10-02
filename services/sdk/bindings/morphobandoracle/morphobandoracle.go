// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package morphobandoracle

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

// MorphoBandOracleMetaData contains all meta data concerning the MorphoBandOracle contract.
var MorphoBandOracleMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"band_\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"symbol_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"loanToken_\",\"type\":\"address\",\"internalType\":\"contractIERC20Metadata\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BAND_DECIMALS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collateralToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"halt\",\"inputs\":[],\"outputs\":[{\"name\":\"signedHalt\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"until\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"issuedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"oraclePaused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"loanToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"price\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"scaleFactor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBand\",\"inputs\":[{\"name\":\"newBand\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BandSet\",\"inputs\":[{\"name\":\"previousBand\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIBand\"},{\"name\":\"newBand\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIBand\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AssetMismatch\",\"inputs\":[{\"name\":\"band\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoAnswer\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoStockToken\",\"inputs\":[{\"name\":\"band\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ScaleTooLarge\",\"inputs\":[{\"name\":\"scaleFactor\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]}]",
	ID:  "MorphoBandOracle",
}

// MorphoBandOracle is an auto generated Go binding around an Ethereum contract.
type MorphoBandOracle struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *MorphoBandOracle) GetABI() abi.ABI {
	return c.abi
}

// NewMorphoBandOracle creates a new instance of MorphoBandOracle.
func NewMorphoBandOracle() *MorphoBandOracle {
	parsed, err := MorphoBandOracleMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MorphoBandOracle{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MorphoBandOracle) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address band_, bytes32 symbol_, address loanToken_, address initialOwner) returns()
func (morphoBandOracle *MorphoBandOracle) PackConstructor(band_ common.Address, symbol_ [32]byte, loanToken_ common.Address, initialOwner common.Address) []byte {
	enc, err := morphoBandOracle.abi.Pack("", band_, symbol_, loanToken_, initialOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBANDDECIMALS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd7a0dfa.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function BAND_DECIMALS() view returns(uint8)
func (morphoBandOracle *MorphoBandOracle) PackBANDDECIMALS() []byte {
	enc, err := morphoBandOracle.abi.Pack("BAND_DECIMALS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBANDDECIMALS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd7a0dfa.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function BAND_DECIMALS() view returns(uint8)
func (morphoBandOracle *MorphoBandOracle) TryPackBANDDECIMALS() ([]byte, error) {
	return morphoBandOracle.abi.Pack("BAND_DECIMALS")
}

// UnpackBANDDECIMALS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd7a0dfa.
//
// Solidity: function BAND_DECIMALS() view returns(uint8)
func (morphoBandOracle *MorphoBandOracle) UnpackBANDDECIMALS(data []byte) (uint8, error) {
	out, err := morphoBandOracle.abi.Unpack("BAND_DECIMALS", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (morphoBandOracle *MorphoBandOracle) PackAcceptOwnership() []byte {
	enc, err := morphoBandOracle.abi.Pack("acceptOwnership")
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
func (morphoBandOracle *MorphoBandOracle) TryPackAcceptOwnership() ([]byte, error) {
	return morphoBandOracle.abi.Pack("acceptOwnership")
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (morphoBandOracle *MorphoBandOracle) PackBand() []byte {
	enc, err := morphoBandOracle.abi.Pack("band")
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
func (morphoBandOracle *MorphoBandOracle) TryPackBand() ([]byte, error) {
	return morphoBandOracle.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (morphoBandOracle *MorphoBandOracle) UnpackBand(data []byte) (common.Address, error) {
	out, err := morphoBandOracle.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCollateralToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2016bd4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collateralToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) PackCollateralToken() []byte {
	enc, err := morphoBandOracle.abi.Pack("collateralToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollateralToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb2016bd4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collateralToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) TryPackCollateralToken() ([]byte, error) {
	return morphoBandOracle.abi.Pack("collateralToken")
}

// UnpackCollateralToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb2016bd4.
//
// Solidity: function collateralToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) UnpackCollateralToken(data []byte) (common.Address, error) {
	out, err := morphoBandOracle.abi.Unpack("collateralToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackHalt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ed7ca5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function halt() view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (morphoBandOracle *MorphoBandOracle) PackHalt() []byte {
	enc, err := morphoBandOracle.abi.Pack("halt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHalt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5ed7ca5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function halt() view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (morphoBandOracle *MorphoBandOracle) TryPackHalt() ([]byte, error) {
	return morphoBandOracle.abi.Pack("halt")
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
// from invoking the contract method with ID 0x5ed7ca5b.
//
// Solidity: function halt() view returns(bool signedHalt, uint64 until, uint64 issuedAt, bool oraclePaused)
func (morphoBandOracle *MorphoBandOracle) UnpackHalt(data []byte) (HaltOutput, error) {
	out, err := morphoBandOracle.abi.Unpack("halt", data)
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

// PackLoanToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06d37817.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function loanToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) PackLoanToken() []byte {
	enc, err := morphoBandOracle.abi.Pack("loanToken")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLoanToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06d37817.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function loanToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) TryPackLoanToken() ([]byte, error) {
	return morphoBandOracle.abi.Pack("loanToken")
}

// UnpackLoanToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06d37817.
//
// Solidity: function loanToken() view returns(address)
func (morphoBandOracle *MorphoBandOracle) UnpackLoanToken(data []byte) (common.Address, error) {
	out, err := morphoBandOracle.abi.Unpack("loanToken", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (morphoBandOracle *MorphoBandOracle) PackOwner() []byte {
	enc, err := morphoBandOracle.abi.Pack("owner")
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
func (morphoBandOracle *MorphoBandOracle) TryPackOwner() ([]byte, error) {
	return morphoBandOracle.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (morphoBandOracle *MorphoBandOracle) UnpackOwner(data []byte) (common.Address, error) {
	out, err := morphoBandOracle.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPendingOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe30c3978.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingOwner() view returns(address)
func (morphoBandOracle *MorphoBandOracle) PackPendingOwner() []byte {
	enc, err := morphoBandOracle.abi.Pack("pendingOwner")
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
func (morphoBandOracle *MorphoBandOracle) TryPackPendingOwner() ([]byte, error) {
	return morphoBandOracle.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (morphoBandOracle *MorphoBandOracle) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := morphoBandOracle.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa035b1fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function price() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) PackPrice() []byte {
	enc, err := morphoBandOracle.abi.Pack("price")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa035b1fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function price() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) TryPackPrice() ([]byte, error) {
	return morphoBandOracle.abi.Pack("price")
}

// UnpackPrice is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa035b1fe.
//
// Solidity: function price() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) UnpackPrice(data []byte) (*big.Int, error) {
	out, err := morphoBandOracle.abi.Unpack("price", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() pure returns()
func (morphoBandOracle *MorphoBandOracle) PackRenounceOwnership() []byte {
	enc, err := morphoBandOracle.abi.Pack("renounceOwnership")
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
func (morphoBandOracle *MorphoBandOracle) TryPackRenounceOwnership() ([]byte, error) {
	return morphoBandOracle.abi.Pack("renounceOwnership")
}

// PackScaleFactor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x683dd191.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scaleFactor() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) PackScaleFactor() []byte {
	enc, err := morphoBandOracle.abi.Pack("scaleFactor")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScaleFactor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x683dd191.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scaleFactor() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) TryPackScaleFactor() ([]byte, error) {
	return morphoBandOracle.abi.Pack("scaleFactor")
}

// UnpackScaleFactor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x683dd191.
//
// Solidity: function scaleFactor() view returns(uint256)
func (morphoBandOracle *MorphoBandOracle) UnpackScaleFactor(data []byte) (*big.Int, error) {
	out, err := morphoBandOracle.abi.Unpack("scaleFactor", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57fc67c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBand(address newBand) returns()
func (morphoBandOracle *MorphoBandOracle) PackSetBand(newBand common.Address) []byte {
	enc, err := morphoBandOracle.abi.Pack("setBand", newBand)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x57fc67c3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBand(address newBand) returns()
func (morphoBandOracle *MorphoBandOracle) TryPackSetBand(newBand common.Address) ([]byte, error) {
	return morphoBandOracle.abi.Pack("setBand", newBand)
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(bytes32)
func (morphoBandOracle *MorphoBandOracle) PackSymbol() []byte {
	enc, err := morphoBandOracle.abi.Pack("symbol")
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
func (morphoBandOracle *MorphoBandOracle) TryPackSymbol() ([]byte, error) {
	return morphoBandOracle.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(bytes32)
func (morphoBandOracle *MorphoBandOracle) UnpackSymbol(data []byte) ([32]byte, error) {
	out, err := morphoBandOracle.abi.Unpack("symbol", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (morphoBandOracle *MorphoBandOracle) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := morphoBandOracle.abi.Pack("transferOwnership", newOwner)
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
func (morphoBandOracle *MorphoBandOracle) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return morphoBandOracle.abi.Pack("transferOwnership", newOwner)
}

// MorphoBandOracleBandSet represents a BandSet event raised by the MorphoBandOracle contract.
type MorphoBandOracleBandSet struct {
	PreviousBand common.Address
	NewBand      common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const MorphoBandOracleBandSetEventName = "BandSet"

// ContractEventName returns the user-defined event name.
func (MorphoBandOracleBandSet) ContractEventName() string {
	return MorphoBandOracleBandSetEventName
}

// UnpackBandSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BandSet(address indexed previousBand, address indexed newBand)
func (morphoBandOracle *MorphoBandOracle) UnpackBandSetEvent(log *types.Log) (*MorphoBandOracleBandSet, error) {
	event := "BandSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != morphoBandOracle.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MorphoBandOracleBandSet)
	if len(log.Data) > 0 {
		if err := morphoBandOracle.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range morphoBandOracle.abi.Events[event].Inputs {
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

// MorphoBandOracleOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MorphoBandOracle contract.
type MorphoBandOracleOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MorphoBandOracleOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (MorphoBandOracleOwnershipTransferStarted) ContractEventName() string {
	return MorphoBandOracleOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (morphoBandOracle *MorphoBandOracle) UnpackOwnershipTransferStartedEvent(log *types.Log) (*MorphoBandOracleOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != morphoBandOracle.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MorphoBandOracleOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := morphoBandOracle.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range morphoBandOracle.abi.Events[event].Inputs {
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

// MorphoBandOracleOwnershipTransferred represents a OwnershipTransferred event raised by the MorphoBandOracle contract.
type MorphoBandOracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MorphoBandOracleOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MorphoBandOracleOwnershipTransferred) ContractEventName() string {
	return MorphoBandOracleOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (morphoBandOracle *MorphoBandOracle) UnpackOwnershipTransferredEvent(log *types.Log) (*MorphoBandOracleOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != morphoBandOracle.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MorphoBandOracleOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := morphoBandOracle.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range morphoBandOracle.abi.Events[event].Inputs {
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
func (morphoBandOracle *MorphoBandOracle) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["AssetMismatch"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackAssetMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["NoAnswer"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackNoAnswerError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["NoStockToken"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackNoStockTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["ScaleTooLarge"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackScaleTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], morphoBandOracle.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return morphoBandOracle.UnpackSequencerNotSettledError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MorphoBandOracleAssetMismatch represents a AssetMismatch error raised by the MorphoBandOracle contract.
type MorphoBandOracleAssetMismatch struct {
	Band  common.Address
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetMismatch(address band, address token)
func MorphoBandOracleAssetMismatchErrorID() common.Hash {
	return common.HexToHash("0x4e83a9b9b56317c86a1fe49825fbcaae2f7d2e01e50a25ee7fc1a426953fda50")
}

// UnpackAssetMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetMismatch(address band, address token)
func (morphoBandOracle *MorphoBandOracle) UnpackAssetMismatchError(raw []byte) (*MorphoBandOracleAssetMismatch, error) {
	out := new(MorphoBandOracleAssetMismatch)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "AssetMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleNoAnswer represents a NoAnswer error raised by the MorphoBandOracle contract.
type MorphoBandOracleNoAnswer struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoAnswer(bytes32 symbol)
func MorphoBandOracleNoAnswerErrorID() common.Hash {
	return common.HexToHash("0x39c5542dd45701eaba9384a6d4e59bd719eacb189e293835099e854560d87f0d")
}

// UnpackNoAnswerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoAnswer(bytes32 symbol)
func (morphoBandOracle *MorphoBandOracle) UnpackNoAnswerError(raw []byte) (*MorphoBandOracleNoAnswer, error) {
	out := new(MorphoBandOracleNoAnswer)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "NoAnswer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleNoStockToken represents a NoStockToken error raised by the MorphoBandOracle contract.
type MorphoBandOracleNoStockToken struct {
	Band   common.Address
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoStockToken(address band, bytes32 symbol)
func MorphoBandOracleNoStockTokenErrorID() common.Hash {
	return common.HexToHash("0xbbcae75dac9307441cd007d5846766c46c76ed50a627dc9ea48c8c6f387d8efb")
}

// UnpackNoStockTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoStockToken(address band, bytes32 symbol)
func (morphoBandOracle *MorphoBandOracle) UnpackNoStockTokenError(raw []byte) (*MorphoBandOracleNoStockToken, error) {
	out := new(MorphoBandOracleNoStockToken)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "NoStockToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MorphoBandOracle contract.
type MorphoBandOracleOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MorphoBandOracleOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (morphoBandOracle *MorphoBandOracle) UnpackOwnableInvalidOwnerError(raw []byte) (*MorphoBandOracleOwnableInvalidOwner, error) {
	out := new(MorphoBandOracleOwnableInvalidOwner)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MorphoBandOracle contract.
type MorphoBandOracleOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MorphoBandOracleOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (morphoBandOracle *MorphoBandOracle) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MorphoBandOracleOwnableUnauthorizedAccount, error) {
	out := new(MorphoBandOracleOwnableUnauthorizedAccount)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the MorphoBandOracle contract.
type MorphoBandOracleOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func MorphoBandOracleOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (morphoBandOracle *MorphoBandOracle) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*MorphoBandOracleOwnershipCannotBeRenounced, error) {
	out := new(MorphoBandOracleOwnershipCannotBeRenounced)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleScaleTooLarge represents a ScaleTooLarge error raised by the MorphoBandOracle contract.
type MorphoBandOracleScaleTooLarge struct {
	ScaleFactor *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ScaleTooLarge(uint256 scaleFactor)
func MorphoBandOracleScaleTooLargeErrorID() common.Hash {
	return common.HexToHash("0x02f1ca2f5fa1afde9a1d53a7dfabdaaaee945222f8e1ac93879b960d54b38bd5")
}

// UnpackScaleTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ScaleTooLarge(uint256 scaleFactor)
func (morphoBandOracle *MorphoBandOracle) UnpackScaleTooLargeError(raw []byte) (*MorphoBandOracleScaleTooLarge, error) {
	out := new(MorphoBandOracleScaleTooLarge)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "ScaleTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MorphoBandOracleSequencerNotSettled represents a SequencerNotSettled error raised by the MorphoBandOracle contract.
type MorphoBandOracleSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func MorphoBandOracleSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (morphoBandOracle *MorphoBandOracle) UnpackSequencerNotSettledError(raw []byte) (*MorphoBandOracleSequencerNotSettled, error) {
	out := new(MorphoBandOracleSequencerNotSettled)
	if err := morphoBandOracle.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}
