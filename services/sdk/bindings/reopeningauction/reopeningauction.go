// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package reopeningauction

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

// ReopeningAuctionBid is an auto generated low-level Go binding around an user-defined struct.
type ReopeningAuctionBid struct {
	Bidder   common.Address
	Quantity *big.Int
	Price    *big.Int
	Escrow   *big.Int
	Claimed  bool
}

// ReopeningAuctionLot is an auto generated low-level Go binding around an user-defined struct.
type ReopeningAuctionLot struct {
	Account  common.Address
	Position [32]byte
	Amount   *big.Int
}

// ReopeningAuctionRound is an auto generated low-level Go binding around an user-defined struct.
type ReopeningAuctionRound struct {
	SealMs   uint64
	Floor    *big.Int
	Supply   *big.Int
	Deposits *big.Int
	Pool     *big.Int
	Price    *big.Int
	Above    *big.Int
	AtPrice  *big.Int
	Sold     *big.Int
	Paid     *big.Int
	Taken    *big.Int
	Cleared  bool
}

// ReopeningAuctionMetaData contains all meta data concerning the ReopeningAuction contract.
var ReopeningAuctionMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"liquidator_\",\"type\":\"address\",\"internalType\":\"contractLiquidator\"},{\"name\":\"symbols\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"feeds_\",\"type\":\"address[]\",\"internalType\":\"contractBandFeed[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BOND\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CLEAR_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FORFEIT_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"HOLIDAY_OPEN_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BIDS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_LOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REOPEN_LEAD_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REVEAL_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMarginAccounts\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bids\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structReopeningAuction.Bid[]\",\"components\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quantity\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"price\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"escrow\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"claimed\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"clear\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"deposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deposit\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"enroll\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enrolled\",\"inputs\":[{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeds\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractBandFeed\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"forfeit\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastOpenMs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractLiquidator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lots\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structReopeningAuction.Lot[]\",\"components\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"phase\",\"inputs\":[],\"outputs\":[{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"revealing\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reveal\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"quantity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"round\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structReopeningAuction.Round\",\"components\":[{\"name\":\"sealMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"floor\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"supply\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deposits\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"pool\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"price\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"above\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"atPrice\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"sold\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"paid\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"taken\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"cleared\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usdg\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Cleared\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"sold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"taken\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Committed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"deposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Enrolled\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Evicted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"quantity\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Forfeited\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"forfeited\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"returned\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LotEvicted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Outbid\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quantity\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Revealed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"quantity\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoundOpened\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"openMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sealMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"floor\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fromSeal\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidBid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommitment\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidFeeds\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoPrice\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotClearingPrice\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotEnrollable\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintToInt\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TooManyLots\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAsset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnknownCommitment\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"WrongPhase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongState\",\"inputs\":[]}]",
	ID:  "ReopeningAuction",
}

// ReopeningAuction is an auto generated Go binding around an Ethereum contract.
type ReopeningAuction struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *ReopeningAuction) GetABI() abi.ABI {
	return c.abi
}

// NewReopeningAuction creates a new instance of ReopeningAuction.
func NewReopeningAuction() *ReopeningAuction {
	parsed, err := ReopeningAuctionMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &ReopeningAuction{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *ReopeningAuction) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address liquidator_, bytes32[] symbols, address[] feeds_) returns()
func (reopeningAuction *ReopeningAuction) PackConstructor(liquidator_ common.Address, symbols [][32]byte, feeds_ []common.Address) []byte {
	enc, err := reopeningAuction.abi.Pack("", liquidator_, symbols, feeds_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackBOND is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1c1d218.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function BOND() view returns(uint256)
func (reopeningAuction *ReopeningAuction) PackBOND() []byte {
	enc, err := reopeningAuction.abi.Pack("BOND")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBOND is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1c1d218.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function BOND() view returns(uint256)
func (reopeningAuction *ReopeningAuction) TryPackBOND() ([]byte, error) {
	return reopeningAuction.abi.Pack("BOND")
}

// UnpackBOND is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc1c1d218.
//
// Solidity: function BOND() view returns(uint256)
func (reopeningAuction *ReopeningAuction) UnpackBOND(data []byte) (*big.Int, error) {
	out, err := reopeningAuction.abi.Unpack("BOND", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCLEARMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4695f585.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CLEAR_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) PackCLEARMS() []byte {
	enc, err := reopeningAuction.abi.Pack("CLEAR_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCLEARMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4695f585.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CLEAR_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) TryPackCLEARMS() ([]byte, error) {
	return reopeningAuction.abi.Pack("CLEAR_MS")
}

// UnpackCLEARMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4695f585.
//
// Solidity: function CLEAR_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) UnpackCLEARMS(data []byte) (uint64, error) {
	out, err := reopeningAuction.abi.Unpack("CLEAR_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackFORFEITBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc183086b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function FORFEIT_BPS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) PackFORFEITBPS() []byte {
	enc, err := reopeningAuction.abi.Pack("FORFEIT_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFORFEITBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc183086b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function FORFEIT_BPS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) TryPackFORFEITBPS() ([]byte, error) {
	return reopeningAuction.abi.Pack("FORFEIT_BPS")
}

// UnpackFORFEITBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc183086b.
//
// Solidity: function FORFEIT_BPS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) UnpackFORFEITBPS(data []byte) (*big.Int, error) {
	out, err := reopeningAuction.abi.Unpack("FORFEIT_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHOLIDAYOPENMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1e28048.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function HOLIDAY_OPEN_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) PackHOLIDAYOPENMS() []byte {
	enc, err := reopeningAuction.abi.Pack("HOLIDAY_OPEN_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHOLIDAYOPENMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1e28048.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function HOLIDAY_OPEN_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) TryPackHOLIDAYOPENMS() ([]byte, error) {
	return reopeningAuction.abi.Pack("HOLIDAY_OPEN_MS")
}

// UnpackHOLIDAYOPENMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb1e28048.
//
// Solidity: function HOLIDAY_OPEN_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) UnpackHOLIDAYOPENMS(data []byte) (uint64, error) {
	out, err := reopeningAuction.abi.Unpack("HOLIDAY_OPEN_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMAXBIDS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5d5f8e1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_BIDS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) PackMAXBIDS() []byte {
	enc, err := reopeningAuction.abi.Pack("MAX_BIDS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXBIDS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd5d5f8e1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_BIDS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) TryPackMAXBIDS() ([]byte, error) {
	return reopeningAuction.abi.Pack("MAX_BIDS")
}

// UnpackMAXBIDS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd5d5f8e1.
//
// Solidity: function MAX_BIDS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) UnpackMAXBIDS(data []byte) (*big.Int, error) {
	out, err := reopeningAuction.abi.Unpack("MAX_BIDS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcec117ba.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_LOTS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) PackMAXLOTS() []byte {
	enc, err := reopeningAuction.abi.Pack("MAX_LOTS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcec117ba.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_LOTS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) TryPackMAXLOTS() ([]byte, error) {
	return reopeningAuction.abi.Pack("MAX_LOTS")
}

// UnpackMAXLOTS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcec117ba.
//
// Solidity: function MAX_LOTS() view returns(uint256)
func (reopeningAuction *ReopeningAuction) UnpackMAXLOTS(data []byte) (*big.Int, error) {
	out, err := reopeningAuction.abi.Unpack("MAX_LOTS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINBID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce69cd20.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_BID() view returns(uint256)
func (reopeningAuction *ReopeningAuction) PackMINBID() []byte {
	enc, err := reopeningAuction.abi.Pack("MIN_BID")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINBID is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce69cd20.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_BID() view returns(uint256)
func (reopeningAuction *ReopeningAuction) TryPackMINBID() ([]byte, error) {
	return reopeningAuction.abi.Pack("MIN_BID")
}

// UnpackMINBID is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce69cd20.
//
// Solidity: function MIN_BID() view returns(uint256)
func (reopeningAuction *ReopeningAuction) UnpackMINBID(data []byte) (*big.Int, error) {
	out, err := reopeningAuction.abi.Unpack("MIN_BID", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackREOPENLEADMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a02dabe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function REOPEN_LEAD_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) PackREOPENLEADMS() []byte {
	enc, err := reopeningAuction.abi.Pack("REOPEN_LEAD_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackREOPENLEADMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a02dabe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function REOPEN_LEAD_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) TryPackREOPENLEADMS() ([]byte, error) {
	return reopeningAuction.abi.Pack("REOPEN_LEAD_MS")
}

// UnpackREOPENLEADMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5a02dabe.
//
// Solidity: function REOPEN_LEAD_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) UnpackREOPENLEADMS(data []byte) (uint64, error) {
	out, err := reopeningAuction.abi.Unpack("REOPEN_LEAD_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackREVEALMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6647c07c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function REVEAL_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) PackREVEALMS() []byte {
	enc, err := reopeningAuction.abi.Pack("REVEAL_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackREVEALMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6647c07c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function REVEAL_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) TryPackREVEALMS() ([]byte, error) {
	return reopeningAuction.abi.Pack("REVEAL_MS")
}

// UnpackREVEALMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6647c07c.
//
// Solidity: function REVEAL_MS() view returns(uint64)
func (reopeningAuction *ReopeningAuction) UnpackREVEALMS(data []byte) (uint64, error) {
	out, err := reopeningAuction.abi.Unpack("REVEAL_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackAccounts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68cd03f6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function accounts() view returns(address)
func (reopeningAuction *ReopeningAuction) PackAccounts() []byte {
	enc, err := reopeningAuction.abi.Pack("accounts")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAccounts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68cd03f6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function accounts() view returns(address)
func (reopeningAuction *ReopeningAuction) TryPackAccounts() ([]byte, error) {
	return reopeningAuction.abi.Pack("accounts")
}

// UnpackAccounts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68cd03f6.
//
// Solidity: function accounts() view returns(address)
func (reopeningAuction *ReopeningAuction) UnpackAccounts(data []byte) (common.Address, error) {
	out, err := reopeningAuction.abi.Unpack("accounts", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (reopeningAuction *ReopeningAuction) PackBand() []byte {
	enc, err := reopeningAuction.abi.Pack("band")
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
func (reopeningAuction *ReopeningAuction) TryPackBand() ([]byte, error) {
	return reopeningAuction.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (reopeningAuction *ReopeningAuction) UnpackBand(data []byte) (common.Address, error) {
	out, err := reopeningAuction.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBids is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x76dd11e1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function bids(bytes32 symbol, uint64 openMs) view returns((address,uint128,uint128,uint128,bool)[])
func (reopeningAuction *ReopeningAuction) PackBids(symbol [32]byte, openMs uint64) []byte {
	enc, err := reopeningAuction.abi.Pack("bids", symbol, openMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBids is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x76dd11e1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function bids(bytes32 symbol, uint64 openMs) view returns((address,uint128,uint128,uint128,bool)[])
func (reopeningAuction *ReopeningAuction) TryPackBids(symbol [32]byte, openMs uint64) ([]byte, error) {
	return reopeningAuction.abi.Pack("bids", symbol, openMs)
}

// UnpackBids is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x76dd11e1.
//
// Solidity: function bids(bytes32 symbol, uint64 openMs) view returns((address,uint128,uint128,uint128,bool)[])
func (reopeningAuction *ReopeningAuction) UnpackBids(data []byte) ([]ReopeningAuctionBid, error) {
	out, err := reopeningAuction.abi.Unpack("bids", data)
	if err != nil {
		return *new([]ReopeningAuctionBid), err
	}
	out0 := *abi.ConvertType(out[0], new([]ReopeningAuctionBid)).(*[]ReopeningAuctionBid)
	return out0, nil
}

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfb167d5f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim(bytes32 symbol, uint64 openMs, uint256 index) returns()
func (reopeningAuction *ReopeningAuction) PackClaim(symbol [32]byte, openMs uint64, index *big.Int) []byte {
	enc, err := reopeningAuction.abi.Pack("claim", symbol, openMs, index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfb167d5f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim(bytes32 symbol, uint64 openMs, uint256 index) returns()
func (reopeningAuction *ReopeningAuction) TryPackClaim(symbol [32]byte, openMs uint64, index *big.Int) ([]byte, error) {
	return reopeningAuction.abi.Pack("claim", symbol, openMs, index)
}

// PackClear is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b053c98.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function clear(bytes32 symbol, uint64 openMs, uint256 price) returns()
func (reopeningAuction *ReopeningAuction) PackClear(symbol [32]byte, openMs uint64, price *big.Int) []byte {
	enc, err := reopeningAuction.abi.Pack("clear", symbol, openMs, price)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClear is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3b053c98.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function clear(bytes32 symbol, uint64 openMs, uint256 price) returns()
func (reopeningAuction *ReopeningAuction) TryPackClear(symbol [32]byte, openMs uint64, price *big.Int) ([]byte, error) {
	return reopeningAuction.abi.Pack("clear", symbol, openMs, price)
}

// PackCommit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a2e7598.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function commit(bytes32 symbol, bytes32 commitment, uint256 deposit) returns()
func (reopeningAuction *ReopeningAuction) PackCommit(symbol [32]byte, commitment [32]byte, deposit *big.Int) []byte {
	enc, err := reopeningAuction.abi.Pack("commit", symbol, commitment, deposit)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCommit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4a2e7598.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function commit(bytes32 symbol, bytes32 commitment, uint256 deposit) returns()
func (reopeningAuction *ReopeningAuction) TryPackCommit(symbol [32]byte, commitment [32]byte, deposit *big.Int) ([]byte, error) {
	return reopeningAuction.abi.Pack("commit", symbol, commitment, deposit)
}

// PackCommitments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd3418a49.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function commitments(bytes32 symbol, uint64 openMs, bytes32 commitment) view returns(address bidder, uint128 deposit)
func (reopeningAuction *ReopeningAuction) PackCommitments(symbol [32]byte, openMs uint64, commitment [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("commitments", symbol, openMs, commitment)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCommitments is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd3418a49.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function commitments(bytes32 symbol, uint64 openMs, bytes32 commitment) view returns(address bidder, uint128 deposit)
func (reopeningAuction *ReopeningAuction) TryPackCommitments(symbol [32]byte, openMs uint64, commitment [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("commitments", symbol, openMs, commitment)
}

// CommitmentsOutput serves as a container for the return parameters of contract
// method Commitments.
type CommitmentsOutput struct {
	Bidder  common.Address
	Deposit *big.Int
}

// UnpackCommitments is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd3418a49.
//
// Solidity: function commitments(bytes32 symbol, uint64 openMs, bytes32 commitment) view returns(address bidder, uint128 deposit)
func (reopeningAuction *ReopeningAuction) UnpackCommitments(data []byte) (CommitmentsOutput, error) {
	out, err := reopeningAuction.abi.Unpack("commitments", data)
	outstruct := new(CommitmentsOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Bidder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Deposit = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackEnroll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdee83a00.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function enroll(address account, bytes32 position, bytes32 symbol) returns()
func (reopeningAuction *ReopeningAuction) PackEnroll(account common.Address, position [32]byte, symbol [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("enroll", account, position, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEnroll is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdee83a00.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function enroll(address account, bytes32 position, bytes32 symbol) returns()
func (reopeningAuction *ReopeningAuction) TryPackEnroll(account common.Address, position [32]byte, symbol [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("enroll", account, position, symbol)
}

// PackEnrolled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9439ce76.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function enrolled(uint64 openMs, address account, bytes32 position, bytes32 symbol) view returns(bool)
func (reopeningAuction *ReopeningAuction) PackEnrolled(openMs uint64, account common.Address, position [32]byte, symbol [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("enrolled", openMs, account, position, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEnrolled is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9439ce76.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function enrolled(uint64 openMs, address account, bytes32 position, bytes32 symbol) view returns(bool)
func (reopeningAuction *ReopeningAuction) TryPackEnrolled(openMs uint64, account common.Address, position [32]byte, symbol [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("enrolled", openMs, account, position, symbol)
}

// UnpackEnrolled is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9439ce76.
//
// Solidity: function enrolled(uint64 openMs, address account, bytes32 position, bytes32 symbol) view returns(bool)
func (reopeningAuction *ReopeningAuction) UnpackEnrolled(data []byte) (bool, error) {
	out, err := reopeningAuction.abi.Unpack("enrolled", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackFeeds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe90f1a43.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeds(bytes32 symbol) view returns(address)
func (reopeningAuction *ReopeningAuction) PackFeeds(symbol [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("feeds", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeds is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe90f1a43.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeds(bytes32 symbol) view returns(address)
func (reopeningAuction *ReopeningAuction) TryPackFeeds(symbol [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("feeds", symbol)
}

// UnpackFeeds is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe90f1a43.
//
// Solidity: function feeds(bytes32 symbol) view returns(address)
func (reopeningAuction *ReopeningAuction) UnpackFeeds(data []byte) (common.Address, error) {
	out, err := reopeningAuction.abi.Unpack("feeds", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackForfeit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x703a9fcb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function forfeit(bytes32 symbol, uint64 openMs, bytes32 commitment) returns()
func (reopeningAuction *ReopeningAuction) PackForfeit(symbol [32]byte, openMs uint64, commitment [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("forfeit", symbol, openMs, commitment)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackForfeit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x703a9fcb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function forfeit(bytes32 symbol, uint64 openMs, bytes32 commitment) returns()
func (reopeningAuction *ReopeningAuction) TryPackForfeit(symbol [32]byte, openMs uint64, commitment [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("forfeit", symbol, openMs, commitment)
}

// PackLastOpenMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3712970.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastOpenMs() view returns(uint64)
func (reopeningAuction *ReopeningAuction) PackLastOpenMs() []byte {
	enc, err := reopeningAuction.abi.Pack("lastOpenMs")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastOpenMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3712970.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastOpenMs() view returns(uint64)
func (reopeningAuction *ReopeningAuction) TryPackLastOpenMs() ([]byte, error) {
	return reopeningAuction.abi.Pack("lastOpenMs")
}

// UnpackLastOpenMs is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3712970.
//
// Solidity: function lastOpenMs() view returns(uint64)
func (reopeningAuction *ReopeningAuction) UnpackLastOpenMs(data []byte) (uint64, error) {
	out, err := reopeningAuction.abi.Unpack("lastOpenMs", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackLiquidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4046ebae.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidator() view returns(address)
func (reopeningAuction *ReopeningAuction) PackLiquidator() []byte {
	enc, err := reopeningAuction.abi.Pack("liquidator")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLiquidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4046ebae.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function liquidator() view returns(address)
func (reopeningAuction *ReopeningAuction) TryPackLiquidator() ([]byte, error) {
	return reopeningAuction.abi.Pack("liquidator")
}

// UnpackLiquidator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4046ebae.
//
// Solidity: function liquidator() view returns(address)
func (reopeningAuction *ReopeningAuction) UnpackLiquidator(data []byte) (common.Address, error) {
	out, err := reopeningAuction.abi.Unpack("liquidator", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackLots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0b15caa.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lots(bytes32 symbol, uint64 openMs) view returns((address,bytes32,uint128)[])
func (reopeningAuction *ReopeningAuction) PackLots(symbol [32]byte, openMs uint64) []byte {
	enc, err := reopeningAuction.abi.Pack("lots", symbol, openMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLots is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc0b15caa.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lots(bytes32 symbol, uint64 openMs) view returns((address,bytes32,uint128)[])
func (reopeningAuction *ReopeningAuction) TryPackLots(symbol [32]byte, openMs uint64) ([]byte, error) {
	return reopeningAuction.abi.Pack("lots", symbol, openMs)
}

// UnpackLots is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc0b15caa.
//
// Solidity: function lots(bytes32 symbol, uint64 openMs) view returns((address,bytes32,uint128)[])
func (reopeningAuction *ReopeningAuction) UnpackLots(data []byte) ([]ReopeningAuctionLot, error) {
	out, err := reopeningAuction.abi.Unpack("lots", data)
	if err != nil {
		return *new([]ReopeningAuctionLot), err
	}
	out0 := *abi.ConvertType(out[0], new([]ReopeningAuctionLot)).(*[]ReopeningAuctionLot)
	return out0, nil
}

// PackPhase is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1c9fe6e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function phase() view returns(uint64 openMs, bool revealing)
func (reopeningAuction *ReopeningAuction) PackPhase() []byte {
	enc, err := reopeningAuction.abi.Pack("phase")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPhase is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb1c9fe6e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function phase() view returns(uint64 openMs, bool revealing)
func (reopeningAuction *ReopeningAuction) TryPackPhase() ([]byte, error) {
	return reopeningAuction.abi.Pack("phase")
}

// PhaseOutput serves as a container for the return parameters of contract
// method Phase.
type PhaseOutput struct {
	OpenMs    uint64
	Revealing bool
}

// UnpackPhase is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb1c9fe6e.
//
// Solidity: function phase() view returns(uint64 openMs, bool revealing)
func (reopeningAuction *ReopeningAuction) UnpackPhase(data []byte) (PhaseOutput, error) {
	out, err := reopeningAuction.abi.Unpack("phase", data)
	outstruct := new(PhaseOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.OpenMs = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Revealing = *abi.ConvertType(out[1], new(bool)).(*bool)
	return *outstruct, nil
}

// PackReveal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc2162553.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reveal(bytes32 symbol, uint64 openMs, uint256 quantity, uint256 price, bytes32 salt) returns()
func (reopeningAuction *ReopeningAuction) PackReveal(symbol [32]byte, openMs uint64, quantity *big.Int, price *big.Int, salt [32]byte) []byte {
	enc, err := reopeningAuction.abi.Pack("reveal", symbol, openMs, quantity, price, salt)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReveal is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc2162553.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reveal(bytes32 symbol, uint64 openMs, uint256 quantity, uint256 price, bytes32 salt) returns()
func (reopeningAuction *ReopeningAuction) TryPackReveal(symbol [32]byte, openMs uint64, quantity *big.Int, price *big.Int, salt [32]byte) ([]byte, error) {
	return reopeningAuction.abi.Pack("reveal", symbol, openMs, quantity, price, salt)
}

// PackRound is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11e01d67.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function round(bytes32 symbol, uint64 openMs) view returns((uint64,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,bool))
func (reopeningAuction *ReopeningAuction) PackRound(symbol [32]byte, openMs uint64) []byte {
	enc, err := reopeningAuction.abi.Pack("round", symbol, openMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRound is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11e01d67.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function round(bytes32 symbol, uint64 openMs) view returns((uint64,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,bool))
func (reopeningAuction *ReopeningAuction) TryPackRound(symbol [32]byte, openMs uint64) ([]byte, error) {
	return reopeningAuction.abi.Pack("round", symbol, openMs)
}

// UnpackRound is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11e01d67.
//
// Solidity: function round(bytes32 symbol, uint64 openMs) view returns((uint64,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,uint128,bool))
func (reopeningAuction *ReopeningAuction) UnpackRound(data []byte) (ReopeningAuctionRound, error) {
	out, err := reopeningAuction.abi.Unpack("round", data)
	if err != nil {
		return *new(ReopeningAuctionRound), err
	}
	out0 := *abi.ConvertType(out[0], new(ReopeningAuctionRound)).(*ReopeningAuctionRound)
	return out0, nil
}

// PackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function usdg() view returns(address)
func (reopeningAuction *ReopeningAuction) PackUsdg() []byte {
	enc, err := reopeningAuction.abi.Pack("usdg")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function usdg() view returns(address)
func (reopeningAuction *ReopeningAuction) TryPackUsdg() ([]byte, error) {
	return reopeningAuction.abi.Pack("usdg")
}

// UnpackUsdg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5b91b7b.
//
// Solidity: function usdg() view returns(address)
func (reopeningAuction *ReopeningAuction) UnpackUsdg(data []byte) (common.Address, error) {
	out, err := reopeningAuction.abi.Unpack("usdg", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// ReopeningAuctionClaimed represents a Claimed event raised by the ReopeningAuction contract.
type ReopeningAuctionClaimed struct {
	Symbol [32]byte
	OpenMs uint64
	Bidder common.Address
	Index  *big.Int
	Amount *big.Int
	Paid   *big.Int
	Refund *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionClaimedEventName = "Claimed"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionClaimed) ContractEventName() string {
	return ReopeningAuctionClaimedEventName
}

// UnpackClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Claimed(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 index, uint256 amount, uint256 paid, uint256 refund)
func (reopeningAuction *ReopeningAuction) UnpackClaimedEvent(log *types.Log) (*ReopeningAuctionClaimed, error) {
	event := "Claimed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionClaimed)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionCleared represents a Cleared event raised by the ReopeningAuction contract.
type ReopeningAuctionCleared struct {
	Symbol [32]byte
	OpenMs uint64
	Price  *big.Int
	Sold   *big.Int
	Paid   *big.Int
	Taken  *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionClearedEventName = "Cleared"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionCleared) ContractEventName() string {
	return ReopeningAuctionClearedEventName
}

// UnpackClearedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Cleared(bytes32 indexed symbol, uint64 indexed openMs, uint256 price, uint256 sold, uint256 paid, uint256 taken)
func (reopeningAuction *ReopeningAuction) UnpackClearedEvent(log *types.Log) (*ReopeningAuctionCleared, error) {
	event := "Cleared"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionCleared)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionCommitted represents a Committed event raised by the ReopeningAuction contract.
type ReopeningAuctionCommitted struct {
	Symbol     [32]byte
	OpenMs     uint64
	Bidder     common.Address
	Commitment [32]byte
	Deposit    *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionCommittedEventName = "Committed"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionCommitted) ContractEventName() string {
	return ReopeningAuctionCommittedEventName
}

// UnpackCommittedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Committed(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, bytes32 commitment, uint256 deposit)
func (reopeningAuction *ReopeningAuction) UnpackCommittedEvent(log *types.Log) (*ReopeningAuctionCommitted, error) {
	event := "Committed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionCommitted)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionEnrolled represents a Enrolled event raised by the ReopeningAuction contract.
type ReopeningAuctionEnrolled struct {
	Symbol   [32]byte
	OpenMs   uint64
	Account  common.Address
	Position [32]byte
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionEnrolledEventName = "Enrolled"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionEnrolled) ContractEventName() string {
	return ReopeningAuctionEnrolledEventName
}

// UnpackEnrolledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Enrolled(bytes32 indexed symbol, uint64 indexed openMs, address indexed account, bytes32 position, uint256 amount)
func (reopeningAuction *ReopeningAuction) UnpackEnrolledEvent(log *types.Log) (*ReopeningAuctionEnrolled, error) {
	event := "Enrolled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionEnrolled)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionEvicted represents a Evicted event raised by the ReopeningAuction contract.
type ReopeningAuctionEvicted struct {
	Symbol   [32]byte
	OpenMs   uint64
	Bidder   common.Address
	Index    *big.Int
	Quantity *big.Int
	Price    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionEvictedEventName = "Evicted"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionEvicted) ContractEventName() string {
	return ReopeningAuctionEvictedEventName
}

// UnpackEvictedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Evicted(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 index, uint256 quantity, uint256 price)
func (reopeningAuction *ReopeningAuction) UnpackEvictedEvent(log *types.Log) (*ReopeningAuctionEvicted, error) {
	event := "Evicted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionEvicted)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionForfeited represents a Forfeited event raised by the ReopeningAuction contract.
type ReopeningAuctionForfeited struct {
	Symbol    [32]byte
	OpenMs    uint64
	Bidder    common.Address
	Forfeited *big.Int
	Returned  *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionForfeitedEventName = "Forfeited"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionForfeited) ContractEventName() string {
	return ReopeningAuctionForfeitedEventName
}

// UnpackForfeitedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Forfeited(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 forfeited, uint256 returned)
func (reopeningAuction *ReopeningAuction) UnpackForfeitedEvent(log *types.Log) (*ReopeningAuctionForfeited, error) {
	event := "Forfeited"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionForfeited)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionLotEvicted represents a LotEvicted event raised by the ReopeningAuction contract.
type ReopeningAuctionLotEvicted struct {
	Symbol   [32]byte
	OpenMs   uint64
	Account  common.Address
	Position [32]byte
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionLotEvictedEventName = "LotEvicted"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionLotEvicted) ContractEventName() string {
	return ReopeningAuctionLotEvictedEventName
}

// UnpackLotEvictedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LotEvicted(bytes32 indexed symbol, uint64 indexed openMs, address indexed account, bytes32 position)
func (reopeningAuction *ReopeningAuction) UnpackLotEvictedEvent(log *types.Log) (*ReopeningAuctionLotEvicted, error) {
	event := "LotEvicted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionLotEvicted)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionOutbid represents a Outbid event raised by the ReopeningAuction contract.
type ReopeningAuctionOutbid struct {
	Symbol   [32]byte
	OpenMs   uint64
	Bidder   common.Address
	Quantity *big.Int
	Price    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionOutbidEventName = "Outbid"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionOutbid) ContractEventName() string {
	return ReopeningAuctionOutbidEventName
}

// UnpackOutbidEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Outbid(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 quantity, uint256 price)
func (reopeningAuction *ReopeningAuction) UnpackOutbidEvent(log *types.Log) (*ReopeningAuctionOutbid, error) {
	event := "Outbid"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionOutbid)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionRevealed represents a Revealed event raised by the ReopeningAuction contract.
type ReopeningAuctionRevealed struct {
	Symbol   [32]byte
	OpenMs   uint64
	Bidder   common.Address
	Index    *big.Int
	Quantity *big.Int
	Price    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionRevealedEventName = "Revealed"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionRevealed) ContractEventName() string {
	return ReopeningAuctionRevealedEventName
}

// UnpackRevealedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Revealed(bytes32 indexed symbol, uint64 indexed openMs, address indexed bidder, uint256 index, uint256 quantity, uint256 price)
func (reopeningAuction *ReopeningAuction) UnpackRevealedEvent(log *types.Log) (*ReopeningAuctionRevealed, error) {
	event := "Revealed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionRevealed)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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

// ReopeningAuctionRoundOpened represents a RoundOpened event raised by the ReopeningAuction contract.
type ReopeningAuctionRoundOpened struct {
	Symbol   [32]byte
	OpenMs   uint64
	SealMs   uint64
	Floor    *big.Int
	FromSeal bool
	Raw      *types.Log // Blockchain specific contextual infos
}

const ReopeningAuctionRoundOpenedEventName = "RoundOpened"

// ContractEventName returns the user-defined event name.
func (ReopeningAuctionRoundOpened) ContractEventName() string {
	return ReopeningAuctionRoundOpenedEventName
}

// UnpackRoundOpenedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoundOpened(bytes32 indexed symbol, uint64 indexed openMs, uint64 sealMs, uint256 floor, bool fromSeal)
func (reopeningAuction *ReopeningAuction) UnpackRoundOpenedEvent(log *types.Log) (*ReopeningAuctionRoundOpened, error) {
	event := "RoundOpened"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != reopeningAuction.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ReopeningAuctionRoundOpened)
	if len(log.Data) > 0 {
		if err := reopeningAuction.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range reopeningAuction.abi.Events[event].Inputs {
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
func (reopeningAuction *ReopeningAuction) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["InvalidBid"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackInvalidBidError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["InvalidCommitment"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackInvalidCommitmentError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["InvalidFeeds"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackInvalidFeedsError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["NoPrice"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackNoPriceError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["NotClearingPrice"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackNotClearingPriceError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["NotEnrollable"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackNotEnrollableError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["SafeCastOverflowedUintToInt"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackSafeCastOverflowedUintToIntError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["TooManyLots"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackTooManyLotsError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["UnknownAsset"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackUnknownAssetError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["UnknownCommitment"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackUnknownCommitmentError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["WrongPhase"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackWrongPhaseError(raw[4:])
	}
	if bytes.Equal(raw[:4], reopeningAuction.abi.Errors["WrongState"].ID.Bytes()[:4]) {
		return reopeningAuction.UnpackWrongStateError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ReopeningAuctionInvalidBid represents a InvalidBid error raised by the ReopeningAuction contract.
type ReopeningAuctionInvalidBid struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidBid()
func ReopeningAuctionInvalidBidErrorID() common.Hash {
	return common.HexToHash("0xc6388ef7fa5b3311415adf90a33e975a1094a30eda7bb67fa93d9fd841e87494")
}

// UnpackInvalidBidError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidBid()
func (reopeningAuction *ReopeningAuction) UnpackInvalidBidError(raw []byte) (*ReopeningAuctionInvalidBid, error) {
	out := new(ReopeningAuctionInvalidBid)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "InvalidBid", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionInvalidCommitment represents a InvalidCommitment error raised by the ReopeningAuction contract.
type ReopeningAuctionInvalidCommitment struct {
	Commitment [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidCommitment(bytes32 commitment)
func ReopeningAuctionInvalidCommitmentErrorID() common.Hash {
	return common.HexToHash("0x537fbfab3dd427d9352b2479b77a7abd292451f570ffa496d1639dc23deb1ca4")
}

// UnpackInvalidCommitmentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidCommitment(bytes32 commitment)
func (reopeningAuction *ReopeningAuction) UnpackInvalidCommitmentError(raw []byte) (*ReopeningAuctionInvalidCommitment, error) {
	out := new(ReopeningAuctionInvalidCommitment)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "InvalidCommitment", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionInvalidFeeds represents a InvalidFeeds error raised by the ReopeningAuction contract.
type ReopeningAuctionInvalidFeeds struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidFeeds()
func ReopeningAuctionInvalidFeedsErrorID() common.Hash {
	return common.HexToHash("0x51a7f8794d70c380776e1d22ad9a4c49f1eb2c3e37fb5db68c58b6f07cf04093")
}

// UnpackInvalidFeedsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidFeeds()
func (reopeningAuction *ReopeningAuction) UnpackInvalidFeedsError(raw []byte) (*ReopeningAuctionInvalidFeeds, error) {
	out := new(ReopeningAuctionInvalidFeeds)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "InvalidFeeds", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionNoPrice represents a NoPrice error raised by the ReopeningAuction contract.
type ReopeningAuctionNoPrice struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoPrice(bytes32 symbol)
func ReopeningAuctionNoPriceErrorID() common.Hash {
	return common.HexToHash("0xcaf0b5a1e2c3ea83b95920138d8acdfc304bce20894c2b008fae66f4a15a05b1")
}

// UnpackNoPriceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoPrice(bytes32 symbol)
func (reopeningAuction *ReopeningAuction) UnpackNoPriceError(raw []byte) (*ReopeningAuctionNoPrice, error) {
	out := new(ReopeningAuctionNoPrice)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "NoPrice", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionNotClearingPrice represents a NotClearingPrice error raised by the ReopeningAuction contract.
type ReopeningAuctionNotClearingPrice struct {
	Price *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotClearingPrice(uint256 price)
func ReopeningAuctionNotClearingPriceErrorID() common.Hash {
	return common.HexToHash("0x772bb04b9792059495240058a3913003cb23a3c0697e95ef43f26dceb83695bd")
}

// UnpackNotClearingPriceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotClearingPrice(uint256 price)
func (reopeningAuction *ReopeningAuction) UnpackNotClearingPriceError(raw []byte) (*ReopeningAuctionNotClearingPrice, error) {
	out := new(ReopeningAuctionNotClearingPrice)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "NotClearingPrice", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionNotEnrollable represents a NotEnrollable error raised by the ReopeningAuction contract.
type ReopeningAuctionNotEnrollable struct {
	Account  common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEnrollable(address account, bytes32 position)
func ReopeningAuctionNotEnrollableErrorID() common.Hash {
	return common.HexToHash("0x2e6739b12a16770d84a548f894852a10b5d2de879578886fadd7204e58ed3e84")
}

// UnpackNotEnrollableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEnrollable(address account, bytes32 position)
func (reopeningAuction *ReopeningAuction) UnpackNotEnrollableError(raw []byte) (*ReopeningAuctionNotEnrollable, error) {
	out := new(ReopeningAuctionNotEnrollable)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "NotEnrollable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the ReopeningAuction contract.
type ReopeningAuctionSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func ReopeningAuctionSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (reopeningAuction *ReopeningAuction) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*ReopeningAuctionSafeCastOverflowedUintDowncast, error) {
	out := new(ReopeningAuctionSafeCastOverflowedUintDowncast)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionSafeCastOverflowedUintToInt represents a SafeCastOverflowedUintToInt error raised by the ReopeningAuction contract.
type ReopeningAuctionSafeCastOverflowedUintToInt struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func ReopeningAuctionSafeCastOverflowedUintToIntErrorID() common.Hash {
	return common.HexToHash("0x24775e0629ae69d78c11bae050651b81820407f300ff750ff2be51e4ce75c37f")
}

// UnpackSafeCastOverflowedUintToIntError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func (reopeningAuction *ReopeningAuction) UnpackSafeCastOverflowedUintToIntError(raw []byte) (*ReopeningAuctionSafeCastOverflowedUintToInt, error) {
	out := new(ReopeningAuctionSafeCastOverflowedUintToInt)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintToInt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the ReopeningAuction contract.
type ReopeningAuctionSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func ReopeningAuctionSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (reopeningAuction *ReopeningAuction) UnpackSafeERC20FailedOperationError(raw []byte) (*ReopeningAuctionSafeERC20FailedOperation, error) {
	out := new(ReopeningAuctionSafeERC20FailedOperation)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionTooManyLots represents a TooManyLots error raised by the ReopeningAuction contract.
type ReopeningAuctionTooManyLots struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyLots()
func ReopeningAuctionTooManyLotsErrorID() common.Hash {
	return common.HexToHash("0x3b263eefe3b20b303c9acb0aec0fb2ea61ccf4ce3534f0fbb23ebd1f95b36d39")
}

// UnpackTooManyLotsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyLots()
func (reopeningAuction *ReopeningAuction) UnpackTooManyLotsError(raw []byte) (*ReopeningAuctionTooManyLots, error) {
	out := new(ReopeningAuctionTooManyLots)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "TooManyLots", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionUnknownAsset represents a UnknownAsset error raised by the ReopeningAuction contract.
type ReopeningAuctionUnknownAsset struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func ReopeningAuctionUnknownAssetErrorID() common.Hash {
	return common.HexToHash("0x1059be3ed9f1c1f061f09edf206f56e7c80b74fd01671cbec1b3788a1f49417e")
}

// UnpackUnknownAssetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func (reopeningAuction *ReopeningAuction) UnpackUnknownAssetError(raw []byte) (*ReopeningAuctionUnknownAsset, error) {
	out := new(ReopeningAuctionUnknownAsset)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "UnknownAsset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionUnknownCommitment represents a UnknownCommitment error raised by the ReopeningAuction contract.
type ReopeningAuctionUnknownCommitment struct {
	Commitment [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownCommitment(bytes32 commitment)
func ReopeningAuctionUnknownCommitmentErrorID() common.Hash {
	return common.HexToHash("0x3f3ff3fce0a88e75f46b6412b76abfc6938b8ce21562ec2d541ab026ecccfa99")
}

// UnpackUnknownCommitmentError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownCommitment(bytes32 commitment)
func (reopeningAuction *ReopeningAuction) UnpackUnknownCommitmentError(raw []byte) (*ReopeningAuctionUnknownCommitment, error) {
	out := new(ReopeningAuctionUnknownCommitment)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "UnknownCommitment", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionWrongPhase represents a WrongPhase error raised by the ReopeningAuction contract.
type ReopeningAuctionWrongPhase struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongPhase()
func ReopeningAuctionWrongPhaseErrorID() common.Hash {
	return common.HexToHash("0xe2586bcc3922aeced192925d9e0eea38bd274dc9ca6c99bad431310815062029")
}

// UnpackWrongPhaseError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongPhase()
func (reopeningAuction *ReopeningAuction) UnpackWrongPhaseError(raw []byte) (*ReopeningAuctionWrongPhase, error) {
	out := new(ReopeningAuctionWrongPhase)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "WrongPhase", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ReopeningAuctionWrongState represents a WrongState error raised by the ReopeningAuction contract.
type ReopeningAuctionWrongState struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongState()
func ReopeningAuctionWrongStateErrorID() common.Hash {
	return common.HexToHash("0xde4168ba7198f38303be592cc28507102d9b35ab9e6847344eb3979e3a60fa52")
}

// UnpackWrongStateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongState()
func (reopeningAuction *ReopeningAuction) UnpackWrongStateError(raw []byte) (*ReopeningAuctionWrongState, error) {
	out := new(ReopeningAuctionWrongState)
	if err := reopeningAuction.abi.UnpackIntoInterface(out, "WrongState", raw); err != nil {
		return nil, err
	}
	return out, nil
}
