// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stocklendingvault

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

// SupplyVaultRateModel is an auto generated low-level Go binding around an user-defined struct.
type SupplyVaultRateModel struct {
	Optimal uint16
	Base    uint32
	Slope1  uint32
	Slope2  uint32
}

// StockLendingVaultMetaData contains all meta data concerning the StockLendingVault contract.
var StockLendingVaultMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialRateModel\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"feeShare_\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_FEE_SHARE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_OPTIMAL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_RATE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_UTILIZATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_OPTIMAL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NOTICE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"assigned\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrowIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrowRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrowable\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrower\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyIn\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimable\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToAssets\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToShares\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"debt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeShare\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"forfeit\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"head\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"idle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastAccrual\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"locked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewDeposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rateModel\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reachable\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recall\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reclaim\",\"inputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"repaid\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requested\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"scaledDebt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"served\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setBorrower\",\"inputs\":[{\"name\":\"newBorrower\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositor\",\"inputs\":[{\"name\":\"newDepositor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRateModel\",\"inputs\":[{\"name\":\"model\",\"type\":\"tuple\",\"internalType\":\"structSupplyVault.RateModel\",\"components\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sync\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ticket\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"start\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"taken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tickets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"utilization\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"writeOff\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"written\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Borrow\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BorrowerSet\",\"inputs\":[{\"name\":\"borrower\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BuyIn\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepositorSet\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeAccrued\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Forfeit\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RateModelSet\",\"inputs\":[{\"name\":\"optimal\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"base\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"slope1\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"slope2\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Recall\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"dueAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Repay\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sync\",\"inputs\":[{\"name\":\"idle\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Take\",\"inputs\":[{\"name\":\"ticket\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteOff\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BuyInFailed\",\"inputs\":[{\"name\":\"outstanding\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientLiquidity\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeeShare\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRateModel\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotBorrower\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotDepositor\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotDue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotHead\",\"inputs\":[{\"name\":\"head\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OutOfReach\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reachable\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Unavailable\",\"inputs\":[]}]",
	ID:  "StockLendingVault",
}

// StockLendingVault is an auto generated Go binding around an Ethereum contract.
type StockLendingVault struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *StockLendingVault) GetABI() abi.ABI {
	return c.abi
}

// NewStockLendingVault creates a new instance of StockLendingVault.
func NewStockLendingVault() *StockLendingVault {
	parsed, err := StockLendingVaultMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &StockLendingVault{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *StockLendingVault) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address token, address initialOwner, (uint16,uint32,uint32,uint32) initialRateModel, uint16 feeShare_) returns()
func (stockLendingVault *StockLendingVault) PackConstructor(token common.Address, initialOwner common.Address, initialRateModel SupplyVaultRateModel, feeShare_ uint16) []byte {
	enc, err := stockLendingVault.abi.Pack("", token, initialOwner, initialRateModel, feeShare_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackMAXFEESHARE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x848a8ce7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_FEE_SHARE() view returns(uint16)
func (stockLendingVault *StockLendingVault) PackMAXFEESHARE() []byte {
	enc, err := stockLendingVault.abi.Pack("MAX_FEE_SHARE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXFEESHARE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x848a8ce7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_FEE_SHARE() view returns(uint16)
func (stockLendingVault *StockLendingVault) TryPackMAXFEESHARE() ([]byte, error) {
	return stockLendingVault.abi.Pack("MAX_FEE_SHARE")
}

// UnpackMAXFEESHARE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x848a8ce7.
//
// Solidity: function MAX_FEE_SHARE() view returns(uint16)
func (stockLendingVault *StockLendingVault) UnpackMAXFEESHARE(data []byte) (uint16, error) {
	out, err := stockLendingVault.abi.Unpack("MAX_FEE_SHARE", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackMAXOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8718699e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) PackMAXOPTIMAL() []byte {
	enc, err := stockLendingVault.abi.Pack("MAX_OPTIMAL")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8718699e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) TryPackMAXOPTIMAL() ([]byte, error) {
	return stockLendingVault.abi.Pack("MAX_OPTIMAL")
}

// UnpackMAXOPTIMAL is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8718699e.
//
// Solidity: function MAX_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) UnpackMAXOPTIMAL(data []byte) (uint16, error) {
	out, err := stockLendingVault.abi.Unpack("MAX_OPTIMAL", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackMAXRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc24dbebd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (stockLendingVault *StockLendingVault) PackMAXRATE() []byte {
	enc, err := stockLendingVault.abi.Pack("MAX_RATE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc24dbebd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (stockLendingVault *StockLendingVault) TryPackMAXRATE() ([]byte, error) {
	return stockLendingVault.abi.Pack("MAX_RATE")
}

// UnpackMAXRATE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc24dbebd.
//
// Solidity: function MAX_RATE() view returns(uint32)
func (stockLendingVault *StockLendingVault) UnpackMAXRATE(data []byte) (uint32, error) {
	out, err := stockLendingVault.abi.Unpack("MAX_RATE", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackMAXUTILIZATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x124fea4f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_UTILIZATION() view returns(uint16)
func (stockLendingVault *StockLendingVault) PackMAXUTILIZATION() []byte {
	enc, err := stockLendingVault.abi.Pack("MAX_UTILIZATION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXUTILIZATION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x124fea4f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_UTILIZATION() view returns(uint16)
func (stockLendingVault *StockLendingVault) TryPackMAXUTILIZATION() ([]byte, error) {
	return stockLendingVault.abi.Pack("MAX_UTILIZATION")
}

// UnpackMAXUTILIZATION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x124fea4f.
//
// Solidity: function MAX_UTILIZATION() view returns(uint16)
func (stockLendingVault *StockLendingVault) UnpackMAXUTILIZATION(data []byte) (uint16, error) {
	out, err := stockLendingVault.abi.Unpack("MAX_UTILIZATION", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackMINOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a1bc24d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) PackMINOPTIMAL() []byte {
	enc, err := stockLendingVault.abi.Pack("MIN_OPTIMAL")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINOPTIMAL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2a1bc24d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) TryPackMINOPTIMAL() ([]byte, error) {
	return stockLendingVault.abi.Pack("MIN_OPTIMAL")
}

// UnpackMINOPTIMAL is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2a1bc24d.
//
// Solidity: function MIN_OPTIMAL() view returns(uint16)
func (stockLendingVault *StockLendingVault) UnpackMINOPTIMAL(data []byte) (uint16, error) {
	out, err := stockLendingVault.abi.Unpack("MIN_OPTIMAL", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackNOTICE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4bfee686.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function NOTICE() view returns(uint64)
func (stockLendingVault *StockLendingVault) PackNOTICE() []byte {
	enc, err := stockLendingVault.abi.Pack("NOTICE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNOTICE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4bfee686.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function NOTICE() view returns(uint64)
func (stockLendingVault *StockLendingVault) TryPackNOTICE() ([]byte, error) {
	return stockLendingVault.abi.Pack("NOTICE")
}

// UnpackNOTICE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4bfee686.
//
// Solidity: function NOTICE() view returns(uint64)
func (stockLendingVault *StockLendingVault) UnpackNOTICE(data []byte) (uint64, error) {
	out, err := stockLendingVault.abi.Unpack("NOTICE", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackAcceptOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ba5097.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptOwnership() returns()
func (stockLendingVault *StockLendingVault) PackAcceptOwnership() []byte {
	enc, err := stockLendingVault.abi.Pack("acceptOwnership")
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
func (stockLendingVault *StockLendingVault) TryPackAcceptOwnership() ([]byte, error) {
	return stockLendingVault.abi.Pack("acceptOwnership")
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("allowance", owner, spender)
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
func (stockLendingVault *StockLendingVault) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("allowance", data)
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
func (stockLendingVault *StockLendingVault) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("approve", spender, value)
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
func (stockLendingVault *StockLendingVault) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (stockLendingVault *StockLendingVault) UnpackApprove(data []byte) (bool, error) {
	out, err := stockLendingVault.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function asset() view returns(address)
func (stockLendingVault *StockLendingVault) PackAsset() []byte {
	enc, err := stockLendingVault.abi.Pack("asset")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x38d52e0f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function asset() view returns(address)
func (stockLendingVault *StockLendingVault) TryPackAsset() ([]byte, error) {
	return stockLendingVault.abi.Pack("asset")
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (stockLendingVault *StockLendingVault) UnpackAsset(data []byte) (common.Address, error) {
	out, err := stockLendingVault.abi.Unpack("asset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackAssigned is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xadb4d990.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function assigned() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackAssigned() []byte {
	enc, err := stockLendingVault.abi.Pack("assigned")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAssigned is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xadb4d990.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function assigned() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackAssigned() ([]byte, error) {
	return stockLendingVault.abi.Pack("assigned")
}

// UnpackAssigned is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xadb4d990.
//
// Solidity: function assigned() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackAssigned(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("assigned", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackBalanceOf(account common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("balanceOf", account)
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
func (stockLendingVault *StockLendingVault) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b3fd148.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrow(uint256 assets, address receiver) returns()
func (stockLendingVault *StockLendingVault) PackBorrow(assets *big.Int, receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("borrow", assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b3fd148.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrow(uint256 assets, address receiver) returns()
func (stockLendingVault *StockLendingVault) TryPackBorrow(assets *big.Int, receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("borrow", assets, receiver)
}

// PackBorrowIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa5af0fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowIndex() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackBorrowIndex() []byte {
	enc, err := stockLendingVault.abi.Pack("borrowIndex")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaa5af0fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowIndex() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackBorrowIndex() ([]byte, error) {
	return stockLendingVault.abi.Pack("borrowIndex")
}

// UnpackBorrowIndex is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaa5af0fd.
//
// Solidity: function borrowIndex() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackBorrowIndex(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("borrowIndex", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrowRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc914b437.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackBorrowRate() []byte {
	enc, err := stockLendingVault.abi.Pack("borrowRate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc914b437.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackBorrowRate() ([]byte, error) {
	return stockLendingVault.abi.Pack("borrowRate")
}

// UnpackBorrowRate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc914b437.
//
// Solidity: function borrowRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackBorrowRate(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("borrowRate", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrowable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe99a6ac7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowable() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackBorrowable() []byte {
	enc, err := stockLendingVault.abi.Pack("borrowable")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe99a6ac7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowable() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackBorrowable() ([]byte, error) {
	return stockLendingVault.abi.Pack("borrowable")
}

// UnpackBorrowable is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe99a6ac7.
//
// Solidity: function borrowable() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackBorrowable(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("borrowable", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7df1f1b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrower() view returns(address)
func (stockLendingVault *StockLendingVault) PackBorrower() []byte {
	enc, err := stockLendingVault.abi.Pack("borrower")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7df1f1b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrower() view returns(address)
func (stockLendingVault *StockLendingVault) TryPackBorrower() ([]byte, error) {
	return stockLendingVault.abi.Pack("borrower")
}

// UnpackBorrower is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7df1f1b9.
//
// Solidity: function borrower() view returns(address)
func (stockLendingVault *StockLendingVault) UnpackBorrower(data []byte) (common.Address, error) {
	out, err := stockLendingVault.abi.Unpack("borrower", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBuyIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb287d8ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buyIn(uint256 id) returns()
func (stockLendingVault *StockLendingVault) PackBuyIn(id *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("buyIn", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuyIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb287d8ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buyIn(uint256 id) returns()
func (stockLendingVault *StockLendingVault) TryPackBuyIn(id *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("buyIn", id)
}

// PackClaimable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd1d58b25.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimable(uint256 id) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackClaimable(id *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("claimable", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd1d58b25.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimable(uint256 id) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackClaimable(id *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("claimable", id)
}

// UnpackClaimable is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd1d58b25.
//
// Solidity: function claimable(uint256 id) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackClaimable(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("claimable", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackConvertToAssets(shares *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("convertToAssets", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackConvertToAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07a2d13a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackConvertToAssets(shares *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("convertToAssets", shares)
}

// UnpackConvertToAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackConvertToAssets(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("convertToAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackConvertToShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6e6f592.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackConvertToShares(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("convertToShares", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackConvertToShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc6e6f592.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackConvertToShares(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("convertToShares", assets)
}

// UnpackConvertToShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackConvertToShares(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("convertToShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dca59c1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function debt() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackDebt() []byte {
	enc, err := stockLendingVault.abi.Pack("debt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0dca59c1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function debt() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackDebt() ([]byte, error) {
	return stockLendingVault.abi.Pack("debt")
}

// UnpackDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0dca59c1.
//
// Solidity: function debt() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackDebt(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("debt", data)
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
func (stockLendingVault *StockLendingVault) PackDecimals() []byte {
	enc, err := stockLendingVault.abi.Pack("decimals")
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
func (stockLendingVault *StockLendingVault) TryPackDecimals() ([]byte, error) {
	return stockLendingVault.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (stockLendingVault *StockLendingVault) UnpackDecimals(data []byte) (uint8, error) {
	out, err := stockLendingVault.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e553f65.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) PackDeposit(assets *big.Int, receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("deposit", assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6e553f65.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackDeposit(assets *big.Int, receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("deposit", assets, receiver)
}

// UnpackDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackDeposit(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("deposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDepositor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7c4ff46.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositor() view returns(address)
func (stockLendingVault *StockLendingVault) PackDepositor() []byte {
	enc, err := stockLendingVault.abi.Pack("depositor")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7c4ff46.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositor() view returns(address)
func (stockLendingVault *StockLendingVault) TryPackDepositor() ([]byte, error) {
	return stockLendingVault.abi.Pack("depositor")
}

// UnpackDepositor is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7c4ff46.
//
// Solidity: function depositor() view returns(address)
func (stockLendingVault *StockLendingVault) UnpackDepositor(data []byte) (common.Address, error) {
	out, err := stockLendingVault.abi.Unpack("depositor", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackFeeShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9ade90e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feeShare() view returns(uint16)
func (stockLendingVault *StockLendingVault) PackFeeShare() []byte {
	enc, err := stockLendingVault.abi.Pack("feeShare")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeeShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9ade90e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feeShare() view returns(uint16)
func (stockLendingVault *StockLendingVault) TryPackFeeShare() ([]byte, error) {
	return stockLendingVault.abi.Pack("feeShare")
}

// UnpackFeeShare is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe9ade90e.
//
// Solidity: function feeShare() view returns(uint16)
func (stockLendingVault *StockLendingVault) UnpackFeeShare(data []byte) (uint16, error) {
	out, err := stockLendingVault.abi.Unpack("feeShare", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackForfeit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ed546ab.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function forfeit(uint256 id) returns()
func (stockLendingVault *StockLendingVault) PackForfeit(id *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("forfeit", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackForfeit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ed546ab.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function forfeit(uint256 id) returns()
func (stockLendingVault *StockLendingVault) TryPackForfeit(id *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("forfeit", id)
}

// PackHead is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f7dcfa3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function head() view returns(uint64)
func (stockLendingVault *StockLendingVault) PackHead() []byte {
	enc, err := stockLendingVault.abi.Pack("head")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHead is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f7dcfa3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function head() view returns(uint64)
func (stockLendingVault *StockLendingVault) TryPackHead() ([]byte, error) {
	return stockLendingVault.abi.Pack("head")
}

// UnpackHead is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8f7dcfa3.
//
// Solidity: function head() view returns(uint64)
func (stockLendingVault *StockLendingVault) UnpackHead(data []byte) (uint64, error) {
	out, err := stockLendingVault.abi.Unpack("head", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackIdle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3192164f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function idle() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackIdle() []byte {
	enc, err := stockLendingVault.abi.Pack("idle")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIdle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3192164f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function idle() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackIdle() ([]byte, error) {
	return stockLendingVault.abi.Pack("idle")
}

// UnpackIdle is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3192164f.
//
// Solidity: function idle() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackIdle(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("idle", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLastAccrual is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b3baab4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastAccrual() view returns(uint64)
func (stockLendingVault *StockLendingVault) PackLastAccrual() []byte {
	enc, err := stockLendingVault.abi.Pack("lastAccrual")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastAccrual is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7b3baab4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastAccrual() view returns(uint64)
func (stockLendingVault *StockLendingVault) TryPackLastAccrual() ([]byte, error) {
	return stockLendingVault.abi.Pack("lastAccrual")
}

// UnpackLastAccrual is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7b3baab4.
//
// Solidity: function lastAccrual() view returns(uint64)
func (stockLendingVault *StockLendingVault) UnpackLastAccrual(data []byte) (uint64, error) {
	out, err := stockLendingVault.abi.Unpack("lastAccrual", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackLocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf309012.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function locked() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackLocked() []byte {
	enc, err := stockLendingVault.abi.Pack("locked")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf309012.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function locked() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackLocked() ([]byte, error) {
	return stockLendingVault.abi.Pack("locked")
}

// UnpackLocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf309012.
//
// Solidity: function locked() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackLocked(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("locked", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackMaxDeposit(receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("maxDeposit", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackMaxDeposit(receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("maxDeposit", receiver)
}

// UnpackMaxDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackMaxDeposit(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("maxDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc63d75b6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackMaxMint(receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("maxMint", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc63d75b6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackMaxMint(receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("maxMint", receiver)
}

// UnpackMaxMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackMaxMint(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("maxMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905777e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackMaxRedeem(owner common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("maxRedeem", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905777e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackMaxRedeem(owner common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("maxRedeem", owner)
}

// UnpackMaxRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackMaxRedeem(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("maxRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMaxWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce96cb77.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackMaxWithdraw(owner common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("maxWithdraw", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMaxWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xce96cb77.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackMaxWithdraw(owner common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("maxWithdraw", owner)
}

// UnpackMaxWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackMaxWithdraw(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("maxWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94bf804d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) PackMint(shares *big.Int, receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("mint", shares, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x94bf804d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackMint(shares *big.Int, receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("mint", shares, receiver)
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackMint(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("mint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (stockLendingVault *StockLendingVault) PackName() []byte {
	enc, err := stockLendingVault.abi.Pack("name")
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
func (stockLendingVault *StockLendingVault) TryPackName() ([]byte, error) {
	return stockLendingVault.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (stockLendingVault *StockLendingVault) UnpackName(data []byte) (string, error) {
	out, err := stockLendingVault.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (stockLendingVault *StockLendingVault) PackOwner() []byte {
	enc, err := stockLendingVault.abi.Pack("owner")
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
func (stockLendingVault *StockLendingVault) TryPackOwner() ([]byte, error) {
	return stockLendingVault.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (stockLendingVault *StockLendingVault) UnpackOwner(data []byte) (common.Address, error) {
	out, err := stockLendingVault.abi.Unpack("owner", data)
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
func (stockLendingVault *StockLendingVault) PackPendingOwner() []byte {
	enc, err := stockLendingVault.abi.Pack("pendingOwner")
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
func (stockLendingVault *StockLendingVault) TryPackPendingOwner() ([]byte, error) {
	return stockLendingVault.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (stockLendingVault *StockLendingVault) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := stockLendingVault.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackPreviewDeposit(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("previewDeposit", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackPreviewDeposit(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("previewDeposit", assets)
}

// UnpackPreviewDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackPreviewDeposit(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("previewDeposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackPreviewMint(shares *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("previewMint", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb3d7f6b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackPreviewMint(shares *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("previewMint", shares)
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackPreviewMint(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("previewMint", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("previewRedeem", shares)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4cdad506.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackPreviewRedeem(shares *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("previewRedeem", shares)
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackPreviewRedeem(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("previewRedeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a28a477.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) PackPreviewWithdraw(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("previewWithdraw", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPreviewWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a28a477.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackPreviewWithdraw(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("previewWithdraw", assets)
}

// UnpackPreviewWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackPreviewWithdraw(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("previewWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1088459.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (stockLendingVault *StockLendingVault) PackRateModel() []byte {
	enc, err := stockLendingVault.abi.Pack("rateModel")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1088459.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (stockLendingVault *StockLendingVault) TryPackRateModel() ([]byte, error) {
	return stockLendingVault.abi.Pack("rateModel")
}

// UnpackRateModel is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1088459.
//
// Solidity: function rateModel() view returns((uint16,uint32,uint32,uint32))
func (stockLendingVault *StockLendingVault) UnpackRateModel(data []byte) (SupplyVaultRateModel, error) {
	out, err := stockLendingVault.abi.Unpack("rateModel", data)
	if err != nil {
		return *new(SupplyVaultRateModel), err
	}
	out0 := *abi.ConvertType(out[0], new(SupplyVaultRateModel)).(*SupplyVaultRateModel)
	return out0, nil
}

// PackReachable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58fde595.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reachable(uint256[] ids) view returns(uint256 assets)
func (stockLendingVault *StockLendingVault) PackReachable(ids []*big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("reachable", ids)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReachable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x58fde595.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reachable(uint256[] ids) view returns(uint256 assets)
func (stockLendingVault *StockLendingVault) TryPackReachable(ids []*big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("reachable", ids)
}

// UnpackReachable is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x58fde595.
//
// Solidity: function reachable(uint256[] ids) view returns(uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackReachable(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("reachable", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d32e793.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recall(uint256 assets) returns(uint256 id)
func (stockLendingVault *StockLendingVault) PackRecall(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("recall", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7d32e793.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recall(uint256 assets) returns(uint256 id)
func (stockLendingVault *StockLendingVault) TryPackRecall(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("recall", assets)
}

// UnpackRecall is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7d32e793.
//
// Solidity: function recall(uint256 assets) returns(uint256 id)
func (stockLendingVault *StockLendingVault) UnpackRecall(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("recall", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackReclaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6a9fdab.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reclaim(uint256[] ids, uint256 assets, address receiver) returns(uint256 shares)
func (stockLendingVault *StockLendingVault) PackReclaim(ids []*big.Int, assets *big.Int, receiver common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("reclaim", ids, assets, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReclaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6a9fdab.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reclaim(uint256[] ids, uint256 assets, address receiver) returns(uint256 shares)
func (stockLendingVault *StockLendingVault) TryPackReclaim(ids []*big.Int, assets *big.Int, receiver common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("reclaim", ids, assets, receiver)
}

// UnpackReclaim is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd6a9fdab.
//
// Solidity: function reclaim(uint256[] ids, uint256 assets, address receiver) returns(uint256 shares)
func (stockLendingVault *StockLendingVault) UnpackReclaim(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("reclaim", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("redeem", shares, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackRedeem(shares *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("redeem", shares, receiver, owner)
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackRedeem(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("redeem", data)
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
func (stockLendingVault *StockLendingVault) PackRenounceOwnership() []byte {
	enc, err := stockLendingVault.abi.Pack("renounceOwnership")
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
func (stockLendingVault *StockLendingVault) TryPackRenounceOwnership() ([]byte, error) {
	return stockLendingVault.abi.Pack("renounceOwnership")
}

// PackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x371fd8e6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (stockLendingVault *StockLendingVault) PackRepay(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("repay", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x371fd8e6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (stockLendingVault *StockLendingVault) TryPackRepay(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("repay", assets)
}

// UnpackRepay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x371fd8e6.
//
// Solidity: function repay(uint256 assets) returns(uint256 repaid)
func (stockLendingVault *StockLendingVault) UnpackRepay(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("repay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRequested is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fae9651.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requested() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackRequested() []byte {
	enc, err := stockLendingVault.abi.Pack("requested")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequested is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fae9651.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requested() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackRequested() ([]byte, error) {
	return stockLendingVault.abi.Pack("requested")
}

// UnpackRequested is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3fae9651.
//
// Solidity: function requested() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackRequested(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("requested", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackScaledDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf7bac74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function scaledDebt() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackScaledDebt() []byte {
	enc, err := stockLendingVault.abi.Pack("scaledDebt")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackScaledDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdf7bac74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function scaledDebt() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackScaledDebt() ([]byte, error) {
	return stockLendingVault.abi.Pack("scaledDebt")
}

// UnpackScaledDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdf7bac74.
//
// Solidity: function scaledDebt() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackScaledDebt(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("scaledDebt", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackServed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x671a308b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function served() view returns(uint128)
func (stockLendingVault *StockLendingVault) PackServed() []byte {
	enc, err := stockLendingVault.abi.Pack("served")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackServed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x671a308b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function served() view returns(uint128)
func (stockLendingVault *StockLendingVault) TryPackServed() ([]byte, error) {
	return stockLendingVault.abi.Pack("served")
}

// UnpackServed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x671a308b.
//
// Solidity: function served() view returns(uint128)
func (stockLendingVault *StockLendingVault) UnpackServed(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("served", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc762d5f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBorrower(address newBorrower) returns()
func (stockLendingVault *StockLendingVault) PackSetBorrower(newBorrower common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("setBorrower", newBorrower)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBorrower is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc762d5f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBorrower(address newBorrower) returns()
func (stockLendingVault *StockLendingVault) TryPackSetBorrower(newBorrower common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("setBorrower", newBorrower)
}

// PackSetDepositor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2c098b7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setDepositor(address newDepositor) returns()
func (stockLendingVault *StockLendingVault) PackSetDepositor(newDepositor common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("setDepositor", newDepositor)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetDepositor is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2c098b7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setDepositor(address newDepositor) returns()
func (stockLendingVault *StockLendingVault) TryPackSetDepositor(newDepositor common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("setDepositor", newDepositor)
}

// PackSetRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x043319d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setRateModel((uint16,uint32,uint32,uint32) model) returns()
func (stockLendingVault *StockLendingVault) PackSetRateModel(model SupplyVaultRateModel) []byte {
	enc, err := stockLendingVault.abi.Pack("setRateModel", model)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetRateModel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x043319d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setRateModel((uint16,uint32,uint32,uint32) model) returns()
func (stockLendingVault *StockLendingVault) TryPackSetRateModel(model SupplyVaultRateModel) ([]byte, error) {
	return stockLendingVault.abi.Pack("setRateModel", model)
}

// PackSupplyRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad2961a3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supplyRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackSupplyRate() []byte {
	enc, err := stockLendingVault.abi.Pack("supplyRate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupplyRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad2961a3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supplyRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackSupplyRate() ([]byte, error) {
	return stockLendingVault.abi.Pack("supplyRate")
}

// UnpackSupplyRate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad2961a3.
//
// Solidity: function supplyRate() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackSupplyRate(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("supplyRate", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (stockLendingVault *StockLendingVault) PackSymbol() []byte {
	enc, err := stockLendingVault.abi.Pack("symbol")
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
func (stockLendingVault *StockLendingVault) TryPackSymbol() ([]byte, error) {
	return stockLendingVault.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (stockLendingVault *StockLendingVault) UnpackSymbol(data []byte) (string, error) {
	out, err := stockLendingVault.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfff6cae9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sync() returns()
func (stockLendingVault *StockLendingVault) PackSync() []byte {
	enc, err := stockLendingVault.abi.Pack("sync")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfff6cae9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sync() returns()
func (stockLendingVault *StockLendingVault) TryPackSync() ([]byte, error) {
	return stockLendingVault.abi.Pack("sync")
}

// PackTicket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc129ff32.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ticket(uint256 id) view returns(uint256 start, uint256 end, uint256 taken, uint64 dueAt)
func (stockLendingVault *StockLendingVault) PackTicket(id *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("ticket", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTicket is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc129ff32.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ticket(uint256 id) view returns(uint256 start, uint256 end, uint256 taken, uint64 dueAt)
func (stockLendingVault *StockLendingVault) TryPackTicket(id *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("ticket", id)
}

// TicketOutput serves as a container for the return parameters of contract
// method Ticket.
type TicketOutput struct {
	Start *big.Int
	End   *big.Int
	Taken *big.Int
	DueAt uint64
}

// UnpackTicket is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc129ff32.
//
// Solidity: function ticket(uint256 id) view returns(uint256 start, uint256 end, uint256 taken, uint64 dueAt)
func (stockLendingVault *StockLendingVault) UnpackTicket(data []byte) (TicketOutput, error) {
	out, err := stockLendingVault.abi.Unpack("ticket", data)
	outstruct := new(TicketOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Start = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.End = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Taken = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	outstruct.DueAt = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackTickets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21858521.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function tickets() view returns(uint64)
func (stockLendingVault *StockLendingVault) PackTickets() []byte {
	enc, err := stockLendingVault.abi.Pack("tickets")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTickets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x21858521.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function tickets() view returns(uint64)
func (stockLendingVault *StockLendingVault) TryPackTickets() ([]byte, error) {
	return stockLendingVault.abi.Pack("tickets")
}

// UnpackTickets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x21858521.
//
// Solidity: function tickets() view returns(uint64)
func (stockLendingVault *StockLendingVault) UnpackTickets(data []byte) (uint64, error) {
	out, err := stockLendingVault.abi.Unpack("tickets", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAssets() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackTotalAssets() []byte {
	enc, err := stockLendingVault.abi.Pack("totalAssets")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalAssets() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackTotalAssets() ([]byte, error) {
	return stockLendingVault.abi.Pack("totalAssets")
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("totalAssets", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackTotalSupply() []byte {
	enc, err := stockLendingVault.abi.Pack("totalSupply")
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
func (stockLendingVault *StockLendingVault) TryPackTotalSupply() ([]byte, error) {
	return stockLendingVault.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("totalSupply", data)
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
func (stockLendingVault *StockLendingVault) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("transfer", to, value)
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
func (stockLendingVault *StockLendingVault) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (stockLendingVault *StockLendingVault) UnpackTransfer(data []byte) (bool, error) {
	out, err := stockLendingVault.abi.Unpack("transfer", data)
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
func (stockLendingVault *StockLendingVault) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("transferFrom", from, to, value)
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
func (stockLendingVault *StockLendingVault) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (stockLendingVault *StockLendingVault) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := stockLendingVault.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (stockLendingVault *StockLendingVault) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("transferOwnership", newOwner)
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
func (stockLendingVault *StockLendingVault) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("transferOwnership", newOwner)
}

// PackUtilization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea21cd92.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function utilization() view returns(uint256)
func (stockLendingVault *StockLendingVault) PackUtilization() []byte {
	enc, err := stockLendingVault.abi.Pack("utilization")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUtilization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea21cd92.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function utilization() view returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackUtilization() ([]byte, error) {
	return stockLendingVault.abi.Pack("utilization")
}

// UnpackUtilization is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xea21cd92.
//
// Solidity: function utilization() view returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackUtilization(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("utilization", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) PackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := stockLendingVault.abi.Pack("withdraw", assets, receiver, owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) TryPackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return stockLendingVault.abi.Pack("withdraw", assets, receiver, owner)
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (stockLendingVault *StockLendingVault) UnpackWithdraw(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("withdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b3e2291.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (stockLendingVault *StockLendingVault) PackWriteOff(assets *big.Int) []byte {
	enc, err := stockLendingVault.abi.Pack("writeOff", assets)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9b3e2291.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (stockLendingVault *StockLendingVault) TryPackWriteOff(assets *big.Int) ([]byte, error) {
	return stockLendingVault.abi.Pack("writeOff", assets)
}

// UnpackWriteOff is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9b3e2291.
//
// Solidity: function writeOff(uint256 assets) returns(uint256 written)
func (stockLendingVault *StockLendingVault) UnpackWriteOff(data []byte) (*big.Int, error) {
	out, err := stockLendingVault.abi.Unpack("writeOff", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// StockLendingVaultApproval represents a Approval event raised by the StockLendingVault contract.
type StockLendingVaultApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultApproval) ContractEventName() string {
	return StockLendingVaultApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (stockLendingVault *StockLendingVault) UnpackApprovalEvent(log *types.Log) (*StockLendingVaultApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultApproval)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultBorrow represents a Borrow event raised by the StockLendingVault contract.
type StockLendingVaultBorrow struct {
	Receiver common.Address
	Assets   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultBorrowEventName = "Borrow"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultBorrow) ContractEventName() string {
	return StockLendingVaultBorrowEventName
}

// UnpackBorrowEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Borrow(address indexed receiver, uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackBorrowEvent(log *types.Log) (*StockLendingVaultBorrow, error) {
	event := "Borrow"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultBorrow)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultBorrowerSet represents a BorrowerSet event raised by the StockLendingVault contract.
type StockLendingVaultBorrowerSet struct {
	Borrower common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultBorrowerSetEventName = "BorrowerSet"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultBorrowerSet) ContractEventName() string {
	return StockLendingVaultBorrowerSetEventName
}

// UnpackBorrowerSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BorrowerSet(address indexed borrower)
func (stockLendingVault *StockLendingVault) UnpackBorrowerSetEvent(log *types.Log) (*StockLendingVaultBorrowerSet, error) {
	event := "BorrowerSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultBorrowerSet)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultBuyIn represents a BuyIn event raised by the StockLendingVault contract.
type StockLendingVaultBuyIn struct {
	Ticket *big.Int
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultBuyInEventName = "BuyIn"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultBuyIn) ContractEventName() string {
	return StockLendingVaultBuyInEventName
}

// UnpackBuyInEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BuyIn(uint256 indexed ticket, uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackBuyInEvent(log *types.Log) (*StockLendingVaultBuyIn, error) {
	event := "BuyIn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultBuyIn)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultDeposit represents a Deposit event raised by the StockLendingVault contract.
type StockLendingVaultDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultDeposit) ContractEventName() string {
	return StockLendingVaultDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (stockLendingVault *StockLendingVault) UnpackDepositEvent(log *types.Log) (*StockLendingVaultDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultDeposit)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultDepositorSet represents a DepositorSet event raised by the StockLendingVault contract.
type StockLendingVaultDepositorSet struct {
	Depositor common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultDepositorSetEventName = "DepositorSet"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultDepositorSet) ContractEventName() string {
	return StockLendingVaultDepositorSetEventName
}

// UnpackDepositorSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DepositorSet(address indexed depositor)
func (stockLendingVault *StockLendingVault) UnpackDepositorSetEvent(log *types.Log) (*StockLendingVaultDepositorSet, error) {
	event := "DepositorSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultDepositorSet)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultFeeAccrued represents a FeeAccrued event raised by the StockLendingVault contract.
type StockLendingVaultFeeAccrued struct {
	Fee    *big.Int
	Index  *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultFeeAccruedEventName = "FeeAccrued"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultFeeAccrued) ContractEventName() string {
	return StockLendingVaultFeeAccruedEventName
}

// UnpackFeeAccruedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeAccrued(uint256 fee, uint256 index, uint256 shares)
func (stockLendingVault *StockLendingVault) UnpackFeeAccruedEvent(log *types.Log) (*StockLendingVaultFeeAccrued, error) {
	event := "FeeAccrued"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultFeeAccrued)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultForfeit represents a Forfeit event raised by the StockLendingVault contract.
type StockLendingVaultForfeit struct {
	Ticket *big.Int
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultForfeitEventName = "Forfeit"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultForfeit) ContractEventName() string {
	return StockLendingVaultForfeitEventName
}

// UnpackForfeitEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Forfeit(uint256 indexed ticket, uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackForfeitEvent(log *types.Log) (*StockLendingVaultForfeit, error) {
	event := "Forfeit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultForfeit)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the StockLendingVault contract.
type StockLendingVaultOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultOwnershipTransferStarted) ContractEventName() string {
	return StockLendingVaultOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (stockLendingVault *StockLendingVault) UnpackOwnershipTransferStartedEvent(log *types.Log) (*StockLendingVaultOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultOwnershipTransferred represents a OwnershipTransferred event raised by the StockLendingVault contract.
type StockLendingVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultOwnershipTransferred) ContractEventName() string {
	return StockLendingVaultOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (stockLendingVault *StockLendingVault) UnpackOwnershipTransferredEvent(log *types.Log) (*StockLendingVaultOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultRateModelSet represents a RateModelSet event raised by the StockLendingVault contract.
type StockLendingVaultRateModelSet struct {
	Optimal uint16
	Base    uint32
	Slope1  uint32
	Slope2  uint32
	Raw     *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultRateModelSetEventName = "RateModelSet"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultRateModelSet) ContractEventName() string {
	return StockLendingVaultRateModelSetEventName
}

// UnpackRateModelSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RateModelSet(uint16 optimal, uint32 base, uint32 slope1, uint32 slope2)
func (stockLendingVault *StockLendingVault) UnpackRateModelSetEvent(log *types.Log) (*StockLendingVaultRateModelSet, error) {
	event := "RateModelSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultRateModelSet)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultRecall represents a Recall event raised by the StockLendingVault contract.
type StockLendingVaultRecall struct {
	Ticket *big.Int
	Assets *big.Int
	DueAt  uint64
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultRecallEventName = "Recall"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultRecall) ContractEventName() string {
	return StockLendingVaultRecallEventName
}

// UnpackRecallEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Recall(uint256 indexed ticket, uint256 assets, uint64 dueAt)
func (stockLendingVault *StockLendingVault) UnpackRecallEvent(log *types.Log) (*StockLendingVaultRecall, error) {
	event := "Recall"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultRecall)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultRepay represents a Repay event raised by the StockLendingVault contract.
type StockLendingVaultRepay struct {
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultRepayEventName = "Repay"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultRepay) ContractEventName() string {
	return StockLendingVaultRepayEventName
}

// UnpackRepayEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Repay(uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackRepayEvent(log *types.Log) (*StockLendingVaultRepay, error) {
	event := "Repay"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultRepay)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultSync represents a Sync event raised by the StockLendingVault contract.
type StockLendingVaultSync struct {
	Idle *big.Int
	Raw  *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultSyncEventName = "Sync"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultSync) ContractEventName() string {
	return StockLendingVaultSyncEventName
}

// UnpackSyncEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sync(uint256 idle)
func (stockLendingVault *StockLendingVault) UnpackSyncEvent(log *types.Log) (*StockLendingVaultSync, error) {
	event := "Sync"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultSync)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultTake represents a Take event raised by the StockLendingVault contract.
type StockLendingVaultTake struct {
	Ticket *big.Int
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultTakeEventName = "Take"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultTake) ContractEventName() string {
	return StockLendingVaultTakeEventName
}

// UnpackTakeEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Take(uint256 indexed ticket, uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackTakeEvent(log *types.Log) (*StockLendingVaultTake, error) {
	event := "Take"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultTake)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultTransfer represents a Transfer event raised by the StockLendingVault contract.
type StockLendingVaultTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultTransfer) ContractEventName() string {
	return StockLendingVaultTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (stockLendingVault *StockLendingVault) UnpackTransferEvent(log *types.Log) (*StockLendingVaultTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultTransfer)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultWithdraw represents a Withdraw event raised by the StockLendingVault contract.
type StockLendingVaultWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultWithdraw) ContractEventName() string {
	return StockLendingVaultWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (stockLendingVault *StockLendingVault) UnpackWithdrawEvent(log *types.Log) (*StockLendingVaultWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultWithdraw)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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

// StockLendingVaultWriteOff represents a WriteOff event raised by the StockLendingVault contract.
type StockLendingVaultWriteOff struct {
	Assets *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const StockLendingVaultWriteOffEventName = "WriteOff"

// ContractEventName returns the user-defined event name.
func (StockLendingVaultWriteOff) ContractEventName() string {
	return StockLendingVaultWriteOffEventName
}

// UnpackWriteOffEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WriteOff(uint256 assets)
func (stockLendingVault *StockLendingVault) UnpackWriteOffEvent(log *types.Log) (*StockLendingVaultWriteOff, error) {
	event := "WriteOff"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != stockLendingVault.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(StockLendingVaultWriteOff)
	if len(log.Data) > 0 {
		if err := stockLendingVault.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range stockLendingVault.abi.Events[event].Inputs {
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
func (stockLendingVault *StockLendingVault) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["AlreadySet"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["BuyInFailed"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackBuyInFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC4626ExceededMaxDeposit"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC4626ExceededMaxDepositError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC4626ExceededMaxMint"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC4626ExceededMaxMintError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC4626ExceededMaxRedeem"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC4626ExceededMaxRedeemError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["ERC4626ExceededMaxWithdraw"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackERC4626ExceededMaxWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["InsufficientLiquidity"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackInsufficientLiquidityError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["InvalidAddress"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackInvalidAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["InvalidFeeShare"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackInvalidFeeShareError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["InvalidRateModel"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackInvalidRateModelError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["NotBorrower"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackNotBorrowerError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["NotDepositor"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackNotDepositorError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["NotDue"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackNotDueError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["NotHead"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackNotHeadError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["OutOfReach"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackOutOfReachError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], stockLendingVault.abi.Errors["Unavailable"].ID.Bytes()[:4]) {
		return stockLendingVault.UnpackUnavailableError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// StockLendingVaultAlreadySet represents a AlreadySet error raised by the StockLendingVault contract.
type StockLendingVaultAlreadySet struct {
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadySet(address current)
func StockLendingVaultAlreadySetErrorID() common.Hash {
	return common.HexToHash("0xe0ecad1a50cc9e7a2cf6f3ea22cabc7010728d195e4e6a5e82ac97e868fdf77e")
}

// UnpackAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadySet(address current)
func (stockLendingVault *StockLendingVault) UnpackAlreadySetError(raw []byte) (*StockLendingVaultAlreadySet, error) {
	out := new(StockLendingVaultAlreadySet)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "AlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultBuyInFailed represents a BuyInFailed error raised by the StockLendingVault contract.
type StockLendingVaultBuyInFailed struct {
	Outstanding *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BuyInFailed(uint256 outstanding)
func StockLendingVaultBuyInFailedErrorID() common.Hash {
	return common.HexToHash("0xc296d0b466019c174a92ee788d61fc4ecc06e9dfd91455bad49738059a3db79b")
}

// UnpackBuyInFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BuyInFailed(uint256 outstanding)
func (stockLendingVault *StockLendingVault) UnpackBuyInFailedError(raw []byte) (*StockLendingVaultBuyInFailed, error) {
	out := new(StockLendingVaultBuyInFailed)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "BuyInFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the StockLendingVault contract.
type StockLendingVaultERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func StockLendingVaultERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (stockLendingVault *StockLendingVault) UnpackERC20InsufficientAllowanceError(raw []byte) (*StockLendingVaultERC20InsufficientAllowance, error) {
	out := new(StockLendingVaultERC20InsufficientAllowance)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the StockLendingVault contract.
type StockLendingVaultERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func StockLendingVaultERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (stockLendingVault *StockLendingVault) UnpackERC20InsufficientBalanceError(raw []byte) (*StockLendingVaultERC20InsufficientBalance, error) {
	out := new(StockLendingVaultERC20InsufficientBalance)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InvalidApprover represents a ERC20InvalidApprover error raised by the StockLendingVault contract.
type StockLendingVaultERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func StockLendingVaultERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (stockLendingVault *StockLendingVault) UnpackERC20InvalidApproverError(raw []byte) (*StockLendingVaultERC20InvalidApprover, error) {
	out := new(StockLendingVaultERC20InvalidApprover)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the StockLendingVault contract.
type StockLendingVaultERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func StockLendingVaultERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (stockLendingVault *StockLendingVault) UnpackERC20InvalidReceiverError(raw []byte) (*StockLendingVaultERC20InvalidReceiver, error) {
	out := new(StockLendingVaultERC20InvalidReceiver)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InvalidSender represents a ERC20InvalidSender error raised by the StockLendingVault contract.
type StockLendingVaultERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func StockLendingVaultERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (stockLendingVault *StockLendingVault) UnpackERC20InvalidSenderError(raw []byte) (*StockLendingVaultERC20InvalidSender, error) {
	out := new(StockLendingVaultERC20InvalidSender)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC20InvalidSpender represents a ERC20InvalidSpender error raised by the StockLendingVault contract.
type StockLendingVaultERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func StockLendingVaultERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (stockLendingVault *StockLendingVault) UnpackERC20InvalidSpenderError(raw []byte) (*StockLendingVaultERC20InvalidSpender, error) {
	out := new(StockLendingVaultERC20InvalidSpender)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC4626ExceededMaxDeposit represents a ERC4626ExceededMaxDeposit error raised by the StockLendingVault contract.
type StockLendingVaultERC4626ExceededMaxDeposit struct {
	Receiver common.Address
	Assets   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func StockLendingVaultERC4626ExceededMaxDepositErrorID() common.Hash {
	return common.HexToHash("0x79012fb2819fcdc1de669c08773dbcd6bdc757862642f75fb1c584dadf259dfe")
}

// UnpackERC4626ExceededMaxDepositError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func (stockLendingVault *StockLendingVault) UnpackERC4626ExceededMaxDepositError(raw []byte) (*StockLendingVaultERC4626ExceededMaxDeposit, error) {
	out := new(StockLendingVaultERC4626ExceededMaxDeposit)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxDeposit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC4626ExceededMaxMint represents a ERC4626ExceededMaxMint error raised by the StockLendingVault contract.
type StockLendingVaultERC4626ExceededMaxMint struct {
	Receiver common.Address
	Shares   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func StockLendingVaultERC4626ExceededMaxMintErrorID() common.Hash {
	return common.HexToHash("0x284ff667dc615a39438518c22e8955b9470327d9de8a4d7e21c926b260d65176")
}

// UnpackERC4626ExceededMaxMintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func (stockLendingVault *StockLendingVault) UnpackERC4626ExceededMaxMintError(raw []byte) (*StockLendingVaultERC4626ExceededMaxMint, error) {
	out := new(StockLendingVaultERC4626ExceededMaxMint)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxMint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC4626ExceededMaxRedeem represents a ERC4626ExceededMaxRedeem error raised by the StockLendingVault contract.
type StockLendingVaultERC4626ExceededMaxRedeem struct {
	Owner  common.Address
	Shares *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func StockLendingVaultERC4626ExceededMaxRedeemErrorID() common.Hash {
	return common.HexToHash("0xb94abeec0557d36b5f0bc8f115deec7b184dcbff94ac66d55e37c8f301e75269")
}

// UnpackERC4626ExceededMaxRedeemError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func (stockLendingVault *StockLendingVault) UnpackERC4626ExceededMaxRedeemError(raw []byte) (*StockLendingVaultERC4626ExceededMaxRedeem, error) {
	out := new(StockLendingVaultERC4626ExceededMaxRedeem)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxRedeem", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultERC4626ExceededMaxWithdraw represents a ERC4626ExceededMaxWithdraw error raised by the StockLendingVault contract.
type StockLendingVaultERC4626ExceededMaxWithdraw struct {
	Owner  common.Address
	Assets *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func StockLendingVaultERC4626ExceededMaxWithdrawErrorID() common.Hash {
	return common.HexToHash("0xfe9cceec2bd1f9b68641914cc354eacaeb1cc2169f5ba0639930f241e87142f0")
}

// UnpackERC4626ExceededMaxWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func (stockLendingVault *StockLendingVault) UnpackERC4626ExceededMaxWithdrawError(raw []byte) (*StockLendingVaultERC4626ExceededMaxWithdraw, error) {
	out := new(StockLendingVaultERC4626ExceededMaxWithdraw)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultInsufficientLiquidity represents a InsufficientLiquidity error raised by the StockLendingVault contract.
type StockLendingVaultInsufficientLiquidity struct {
	Requested *big.Int
	Available *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientLiquidity(uint256 requested, uint256 available)
func StockLendingVaultInsufficientLiquidityErrorID() common.Hash {
	return common.HexToHash("0xa17e11d5aed9d2c4240142d28db52e7da7283f43ac33766dd68e8f3dc6613c0d")
}

// UnpackInsufficientLiquidityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientLiquidity(uint256 requested, uint256 available)
func (stockLendingVault *StockLendingVault) UnpackInsufficientLiquidityError(raw []byte) (*StockLendingVaultInsufficientLiquidity, error) {
	out := new(StockLendingVaultInsufficientLiquidity)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "InsufficientLiquidity", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultInvalidAddress represents a InvalidAddress error raised by the StockLendingVault contract.
type StockLendingVaultInvalidAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidAddress()
func StockLendingVaultInvalidAddressErrorID() common.Hash {
	return common.HexToHash("0xe6c4247b90bd06996a32d386bb770af9c0018dd1b0ebbb2df2c4499f1eda7b16")
}

// UnpackInvalidAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidAddress()
func (stockLendingVault *StockLendingVault) UnpackInvalidAddressError(raw []byte) (*StockLendingVaultInvalidAddress, error) {
	out := new(StockLendingVaultInvalidAddress)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "InvalidAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultInvalidFeeShare represents a InvalidFeeShare error raised by the StockLendingVault contract.
type StockLendingVaultInvalidFeeShare struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidFeeShare()
func StockLendingVaultInvalidFeeShareErrorID() common.Hash {
	return common.HexToHash("0xe8cdd5bd149b84048627b9039a9005cfe5a7bbd788bf3d365c4aaff77b4ccda1")
}

// UnpackInvalidFeeShareError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidFeeShare()
func (stockLendingVault *StockLendingVault) UnpackInvalidFeeShareError(raw []byte) (*StockLendingVaultInvalidFeeShare, error) {
	out := new(StockLendingVaultInvalidFeeShare)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "InvalidFeeShare", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultInvalidRateModel represents a InvalidRateModel error raised by the StockLendingVault contract.
type StockLendingVaultInvalidRateModel struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRateModel()
func StockLendingVaultInvalidRateModelErrorID() common.Hash {
	return common.HexToHash("0xd2375d9d660f19094039027d0f6695eff69fd031e8b29f0e5dc1923e6f5a0d25")
}

// UnpackInvalidRateModelError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRateModel()
func (stockLendingVault *StockLendingVault) UnpackInvalidRateModelError(raw []byte) (*StockLendingVaultInvalidRateModel, error) {
	out := new(StockLendingVaultInvalidRateModel)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "InvalidRateModel", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultNotBorrower represents a NotBorrower error raised by the StockLendingVault contract.
type StockLendingVaultNotBorrower struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotBorrower(address account)
func StockLendingVaultNotBorrowerErrorID() common.Hash {
	return common.HexToHash("0x3d19c0eaa213120f897b12cdf0df62bf048ed2228f0ad44d3e89ce9724036d29")
}

// UnpackNotBorrowerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotBorrower(address account)
func (stockLendingVault *StockLendingVault) UnpackNotBorrowerError(raw []byte) (*StockLendingVaultNotBorrower, error) {
	out := new(StockLendingVaultNotBorrower)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "NotBorrower", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultNotDepositor represents a NotDepositor error raised by the StockLendingVault contract.
type StockLendingVaultNotDepositor struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotDepositor(address account)
func StockLendingVaultNotDepositorErrorID() common.Hash {
	return common.HexToHash("0xc6cd8a0b10e8457312985cc892e795e9216d2349829347dbe8134cd5a6a0602e")
}

// UnpackNotDepositorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotDepositor(address account)
func (stockLendingVault *StockLendingVault) UnpackNotDepositorError(raw []byte) (*StockLendingVaultNotDepositor, error) {
	out := new(StockLendingVaultNotDepositor)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "NotDepositor", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultNotDue represents a NotDue error raised by the StockLendingVault contract.
type StockLendingVaultNotDue struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotDue()
func StockLendingVaultNotDueErrorID() common.Hash {
	return common.HexToHash("0x47a2375f61a795c75725939fdb7dae821fe40ef2dc1cdd65800568492a89fcc4")
}

// UnpackNotDueError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotDue()
func (stockLendingVault *StockLendingVault) UnpackNotDueError(raw []byte) (*StockLendingVaultNotDue, error) {
	out := new(StockLendingVaultNotDue)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "NotDue", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultNotHead represents a NotHead error raised by the StockLendingVault contract.
type StockLendingVaultNotHead struct {
	Head *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotHead(uint256 head)
func StockLendingVaultNotHeadErrorID() common.Hash {
	return common.HexToHash("0xa6e80ffd441b79901dd759c724970127cf68e862bb0b101ff77b5ee62ed5604f")
}

// UnpackNotHeadError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotHead(uint256 head)
func (stockLendingVault *StockLendingVault) UnpackNotHeadError(raw []byte) (*StockLendingVaultNotHead, error) {
	out := new(StockLendingVaultNotHead)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "NotHead", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultOutOfReach represents a OutOfReach error raised by the StockLendingVault contract.
type StockLendingVaultOutOfReach struct {
	Assets    *big.Int
	Reachable *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OutOfReach(uint256 assets, uint256 reachable)
func StockLendingVaultOutOfReachErrorID() common.Hash {
	return common.HexToHash("0x6433ea35ffd4fcf838d617aacafc90ea85bf99456c7efa209245c75e03dd68e8")
}

// UnpackOutOfReachError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OutOfReach(uint256 assets, uint256 reachable)
func (stockLendingVault *StockLendingVault) UnpackOutOfReachError(raw []byte) (*StockLendingVaultOutOfReach, error) {
	out := new(StockLendingVaultOutOfReach)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "OutOfReach", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the StockLendingVault contract.
type StockLendingVaultOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func StockLendingVaultOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (stockLendingVault *StockLendingVault) UnpackOwnableInvalidOwnerError(raw []byte) (*StockLendingVaultOwnableInvalidOwner, error) {
	out := new(StockLendingVaultOwnableInvalidOwner)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the StockLendingVault contract.
type StockLendingVaultOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func StockLendingVaultOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (stockLendingVault *StockLendingVault) UnpackOwnableUnauthorizedAccountError(raw []byte) (*StockLendingVaultOwnableUnauthorizedAccount, error) {
	out := new(StockLendingVaultOwnableUnauthorizedAccount)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the StockLendingVault contract.
type StockLendingVaultOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func StockLendingVaultOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (stockLendingVault *StockLendingVault) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*StockLendingVaultOwnershipCannotBeRenounced, error) {
	out := new(StockLendingVaultOwnershipCannotBeRenounced)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the StockLendingVault contract.
type StockLendingVaultSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func StockLendingVaultSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (stockLendingVault *StockLendingVault) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*StockLendingVaultSafeCastOverflowedUintDowncast, error) {
	out := new(StockLendingVaultSafeCastOverflowedUintDowncast)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the StockLendingVault contract.
type StockLendingVaultSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func StockLendingVaultSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (stockLendingVault *StockLendingVault) UnpackSafeERC20FailedOperationError(raw []byte) (*StockLendingVaultSafeERC20FailedOperation, error) {
	out := new(StockLendingVaultSafeERC20FailedOperation)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// StockLendingVaultUnavailable represents a Unavailable error raised by the StockLendingVault contract.
type StockLendingVaultUnavailable struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Unavailable()
func StockLendingVaultUnavailableErrorID() common.Hash {
	return common.HexToHash("0xa3b8915fe1a2de6a78cee859ded8d8d4b8e246f746433fbddef272cfe0c13c07")
}

// UnpackUnavailableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Unavailable()
func (stockLendingVault *StockLendingVault) UnpackUnavailableError(raw []byte) (*StockLendingVaultUnavailable, error) {
	out := new(StockLendingVaultUnavailable)
	if err := stockLendingVault.abi.UnpackIntoInterface(out, "Unavailable", raw); err != nil {
		return nil, err
	}
	return out, nil
}
