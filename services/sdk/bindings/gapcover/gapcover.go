// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gapcover

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

// GapCoverMetaData contains all meta data concerning the GapCover contract.
var GapCoverMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"engine_\",\"type\":\"address\",\"internalType\":\"contractIMargin\"},{\"name\":\"usdg_\",\"type\":\"address\",\"internalType\":\"contractIUSDG\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"LOADING_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_DELAY_SLOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_GAP_MULTIPLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_GAP_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MOVE_TO_GAP_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"POST_MARKET_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"QUORUM_SLOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REOPEN_BEFORE_OPEN_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SETTLEMENT_DEADLINE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SLOT_MS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TAIL_PROBABILITY_PPM\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TAIL_SCALE_PPM\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TAIL_THRESHOLD_PPM\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WEEK_DAYS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WINDOW_SLOTS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"asset\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buy\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"notional\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deductibleBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxPremium\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"premium\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"capacity\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"convertToAssets\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"convertToShares\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"coverCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"covers\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deductibleBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"limitBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"notional\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"premium\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"engine\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMargin\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"held\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastCloseMs\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"measure\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"weekMove\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minDeductible\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"observe\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"outstanding\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"outstandingIn\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payouts\",\"inputs\":[{\"name\":\"holder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"premiums\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewDeposit\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewMint\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRedeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewWithdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pricingGap\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"gap\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"weekMove\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quote\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"notional\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deductibleBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"record\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"regularHours\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"payout\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reopenOf\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reserved\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sales\",\"inputs\":[],\"outputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endsMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"series\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"notional\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"referencePrice\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"price\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"shift\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"flagged\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"referenceRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"},{\"name\":\"lastRound\",\"type\":\"uint80\",\"internalType\":\"uint80\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sync\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalAssets\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"void\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Bought\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"holder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"notional\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"deductibleBps\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"limitBps\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"premium\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"holder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CloseRecorded\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"atMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Measured\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"weekMove\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Observed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"slot\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"mid\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"low\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"high\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Released\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"holder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"refund\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReopenRecorded\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"reopensMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Settled\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"referenceRound\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"},{\"name\":\"referencePrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"firstRound\",\"type\":\"uint80\",\"indexed\":false,\"internalType\":\"uint80\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"flagged\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sync\",\"inputs\":[{\"name\":\"held\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Voided\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AssetHalted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CorporateActionPending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxDeposit\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxMint\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxRedeem\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC4626ExceededMaxWithdraw\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGas\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHolder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLastRound\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}]},{\"type\":\"error\",\"name\":\"InvalidLayer\",\"inputs\":[{\"name\":\"deductibleBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limitBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidReference\",\"inputs\":[{\"name\":\"roundId\",\"type\":\"uint80\",\"internalType\":\"uint80\"}]},{\"type\":\"error\",\"name\":\"NoCapacity\",\"inputs\":[{\"name\":\"need\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"free\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NoCover\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotCoverable\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotReopened\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"NotSettled\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NothingCovered\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"PremiumAboveLimit\",\"inputs\":[{\"name\":\"premium\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxPremium\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedIntToUint\",\"inputs\":[{\"name\":\"value\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SalesClosed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SeriesClosed\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"StaleReference\",\"inputs\":[{\"name\":\"deductibleBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TooEarlyToMeasure\",\"inputs\":[{\"name\":\"fromMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TooEarlyToVoid\",\"inputs\":[{\"name\":\"fromMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WindowClosed\",\"inputs\":[{\"name\":\"endMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WindowOpen\",\"inputs\":[{\"name\":\"endMs\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
	ID:  "GapCover",
}

// GapCover is an auto generated Go binding around an Ethereum contract.
type GapCover struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *GapCover) GetABI() abi.ABI {
	return c.abi
}

// NewGapCover creates a new instance of GapCover.
func NewGapCover() *GapCover {
	parsed, err := GapCoverMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &GapCover{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *GapCover) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address engine_, address usdg_) returns()
func (gapCover *GapCover) PackConstructor(engine_ common.Address, usdg_ common.Address) []byte {
	enc, err := gapCover.abi.Pack("", engine_, usdg_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackLOADINGBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa6f0b654.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function LOADING_BPS() view returns(uint256)
func (gapCover *GapCover) PackLOADINGBPS() []byte {
	enc, err := gapCover.abi.Pack("LOADING_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLOADINGBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa6f0b654.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function LOADING_BPS() view returns(uint256)
func (gapCover *GapCover) TryPackLOADINGBPS() ([]byte, error) {
	return gapCover.abi.Pack("LOADING_BPS")
}

// UnpackLOADINGBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa6f0b654.
//
// Solidity: function LOADING_BPS() view returns(uint256)
func (gapCover *GapCover) UnpackLOADINGBPS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("LOADING_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXDELAYSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad7f9cee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_DELAY_SLOTS() view returns(uint256)
func (gapCover *GapCover) PackMAXDELAYSLOTS() []byte {
	enc, err := gapCover.abi.Pack("MAX_DELAY_SLOTS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXDELAYSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xad7f9cee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_DELAY_SLOTS() view returns(uint256)
func (gapCover *GapCover) TryPackMAXDELAYSLOTS() ([]byte, error) {
	return gapCover.abi.Pack("MAX_DELAY_SLOTS")
}

// UnpackMAXDELAYSLOTS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xad7f9cee.
//
// Solidity: function MAX_DELAY_SLOTS() view returns(uint256)
func (gapCover *GapCover) UnpackMAXDELAYSLOTS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("MAX_DELAY_SLOTS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXGAPMULTIPLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x89b7d291.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_GAP_MULTIPLE() view returns(uint256)
func (gapCover *GapCover) PackMAXGAPMULTIPLE() []byte {
	enc, err := gapCover.abi.Pack("MAX_GAP_MULTIPLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXGAPMULTIPLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x89b7d291.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_GAP_MULTIPLE() view returns(uint256)
func (gapCover *GapCover) TryPackMAXGAPMULTIPLE() ([]byte, error) {
	return gapCover.abi.Pack("MAX_GAP_MULTIPLE")
}

// UnpackMAXGAPMULTIPLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x89b7d291.
//
// Solidity: function MAX_GAP_MULTIPLE() view returns(uint256)
func (gapCover *GapCover) UnpackMAXGAPMULTIPLE(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("MAX_GAP_MULTIPLE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINGAPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x990fdeb9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) PackMINGAPBPS() []byte {
	enc, err := gapCover.abi.Pack("MIN_GAP_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINGAPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x990fdeb9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) TryPackMINGAPBPS() ([]byte, error) {
	return gapCover.abi.Pack("MIN_GAP_BPS")
}

// UnpackMINGAPBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x990fdeb9.
//
// Solidity: function MIN_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) UnpackMINGAPBPS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("MIN_GAP_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMOVETOGAPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3833a4b9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MOVE_TO_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) PackMOVETOGAPBPS() []byte {
	enc, err := gapCover.abi.Pack("MOVE_TO_GAP_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMOVETOGAPBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3833a4b9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MOVE_TO_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) TryPackMOVETOGAPBPS() ([]byte, error) {
	return gapCover.abi.Pack("MOVE_TO_GAP_BPS")
}

// UnpackMOVETOGAPBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3833a4b9.
//
// Solidity: function MOVE_TO_GAP_BPS() view returns(uint256)
func (gapCover *GapCover) UnpackMOVETOGAPBPS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("MOVE_TO_GAP_BPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPOSTMARKETMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0529cd8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function POST_MARKET_MS() view returns(uint64)
func (gapCover *GapCover) PackPOSTMARKETMS() []byte {
	enc, err := gapCover.abi.Pack("POST_MARKET_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPOSTMARKETMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb0529cd8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function POST_MARKET_MS() view returns(uint64)
func (gapCover *GapCover) TryPackPOSTMARKETMS() ([]byte, error) {
	return gapCover.abi.Pack("POST_MARKET_MS")
}

// UnpackPOSTMARKETMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb0529cd8.
//
// Solidity: function POST_MARKET_MS() view returns(uint64)
func (gapCover *GapCover) UnpackPOSTMARKETMS(data []byte) (uint64, error) {
	out, err := gapCover.abi.Unpack("POST_MARKET_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackQUORUMSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf7d749c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function QUORUM_SLOTS() view returns(uint256)
func (gapCover *GapCover) PackQUORUMSLOTS() []byte {
	enc, err := gapCover.abi.Pack("QUORUM_SLOTS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQUORUMSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf7d749c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function QUORUM_SLOTS() view returns(uint256)
func (gapCover *GapCover) TryPackQUORUMSLOTS() ([]byte, error) {
	return gapCover.abi.Pack("QUORUM_SLOTS")
}

// UnpackQUORUMSLOTS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf7d749c.
//
// Solidity: function QUORUM_SLOTS() view returns(uint256)
func (gapCover *GapCover) UnpackQUORUMSLOTS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("QUORUM_SLOTS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackREOPENBEFOREOPENMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27dec4c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function REOPEN_BEFORE_OPEN_MS() view returns(uint64)
func (gapCover *GapCover) PackREOPENBEFOREOPENMS() []byte {
	enc, err := gapCover.abi.Pack("REOPEN_BEFORE_OPEN_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackREOPENBEFOREOPENMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27dec4c3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function REOPEN_BEFORE_OPEN_MS() view returns(uint64)
func (gapCover *GapCover) TryPackREOPENBEFOREOPENMS() ([]byte, error) {
	return gapCover.abi.Pack("REOPEN_BEFORE_OPEN_MS")
}

// UnpackREOPENBEFOREOPENMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27dec4c3.
//
// Solidity: function REOPEN_BEFORE_OPEN_MS() view returns(uint64)
func (gapCover *GapCover) UnpackREOPENBEFOREOPENMS(data []byte) (uint64, error) {
	out, err := gapCover.abi.Unpack("REOPEN_BEFORE_OPEN_MS", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackSETTLEMENTDEADLINE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd687cd62.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function SETTLEMENT_DEADLINE() view returns(uint256)
func (gapCover *GapCover) PackSETTLEMENTDEADLINE() []byte {
	enc, err := gapCover.abi.Pack("SETTLEMENT_DEADLINE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSETTLEMENTDEADLINE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd687cd62.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function SETTLEMENT_DEADLINE() view returns(uint256)
func (gapCover *GapCover) TryPackSETTLEMENTDEADLINE() ([]byte, error) {
	return gapCover.abi.Pack("SETTLEMENT_DEADLINE")
}

// UnpackSETTLEMENTDEADLINE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd687cd62.
//
// Solidity: function SETTLEMENT_DEADLINE() view returns(uint256)
func (gapCover *GapCover) UnpackSETTLEMENTDEADLINE(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("SETTLEMENT_DEADLINE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSLOTMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e692f8c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function SLOT_MS() view returns(uint256)
func (gapCover *GapCover) PackSLOTMS() []byte {
	enc, err := gapCover.abi.Pack("SLOT_MS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSLOTMS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e692f8c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function SLOT_MS() view returns(uint256)
func (gapCover *GapCover) TryPackSLOTMS() ([]byte, error) {
	return gapCover.abi.Pack("SLOT_MS")
}

// UnpackSLOTMS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3e692f8c.
//
// Solidity: function SLOT_MS() view returns(uint256)
func (gapCover *GapCover) UnpackSLOTMS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("SLOT_MS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTAILPROBABILITYPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4193114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function TAIL_PROBABILITY_PPM() view returns(uint256)
func (gapCover *GapCover) PackTAILPROBABILITYPPM() []byte {
	enc, err := gapCover.abi.Pack("TAIL_PROBABILITY_PPM")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTAILPROBABILITYPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4193114.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function TAIL_PROBABILITY_PPM() view returns(uint256)
func (gapCover *GapCover) TryPackTAILPROBABILITYPPM() ([]byte, error) {
	return gapCover.abi.Pack("TAIL_PROBABILITY_PPM")
}

// UnpackTAILPROBABILITYPPM is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc4193114.
//
// Solidity: function TAIL_PROBABILITY_PPM() view returns(uint256)
func (gapCover *GapCover) UnpackTAILPROBABILITYPPM(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("TAIL_PROBABILITY_PPM", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTAILSCALEPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x092f83aa.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function TAIL_SCALE_PPM() view returns(uint256)
func (gapCover *GapCover) PackTAILSCALEPPM() []byte {
	enc, err := gapCover.abi.Pack("TAIL_SCALE_PPM")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTAILSCALEPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x092f83aa.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function TAIL_SCALE_PPM() view returns(uint256)
func (gapCover *GapCover) TryPackTAILSCALEPPM() ([]byte, error) {
	return gapCover.abi.Pack("TAIL_SCALE_PPM")
}

// UnpackTAILSCALEPPM is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x092f83aa.
//
// Solidity: function TAIL_SCALE_PPM() view returns(uint256)
func (gapCover *GapCover) UnpackTAILSCALEPPM(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("TAIL_SCALE_PPM", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTAILTHRESHOLDPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xac3581c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function TAIL_THRESHOLD_PPM() view returns(uint256)
func (gapCover *GapCover) PackTAILTHRESHOLDPPM() []byte {
	enc, err := gapCover.abi.Pack("TAIL_THRESHOLD_PPM")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTAILTHRESHOLDPPM is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xac3581c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function TAIL_THRESHOLD_PPM() view returns(uint256)
func (gapCover *GapCover) TryPackTAILTHRESHOLDPPM() ([]byte, error) {
	return gapCover.abi.Pack("TAIL_THRESHOLD_PPM")
}

// UnpackTAILTHRESHOLDPPM is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xac3581c8.
//
// Solidity: function TAIL_THRESHOLD_PPM() view returns(uint256)
func (gapCover *GapCover) UnpackTAILTHRESHOLDPPM(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("TAIL_THRESHOLD_PPM", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWEEKDAYS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48eac431.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function WEEK_DAYS() view returns(uint256)
func (gapCover *GapCover) PackWEEKDAYS() []byte {
	enc, err := gapCover.abi.Pack("WEEK_DAYS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWEEKDAYS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x48eac431.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function WEEK_DAYS() view returns(uint256)
func (gapCover *GapCover) TryPackWEEKDAYS() ([]byte, error) {
	return gapCover.abi.Pack("WEEK_DAYS")
}

// UnpackWEEKDAYS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x48eac431.
//
// Solidity: function WEEK_DAYS() view returns(uint256)
func (gapCover *GapCover) UnpackWEEKDAYS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("WEEK_DAYS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWINDOWSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1252350d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function WINDOW_SLOTS() view returns(uint256)
func (gapCover *GapCover) PackWINDOWSLOTS() []byte {
	enc, err := gapCover.abi.Pack("WINDOW_SLOTS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWINDOWSLOTS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1252350d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function WINDOW_SLOTS() view returns(uint256)
func (gapCover *GapCover) TryPackWINDOWSLOTS() ([]byte, error) {
	return gapCover.abi.Pack("WINDOW_SLOTS")
}

// UnpackWINDOWSLOTS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1252350d.
//
// Solidity: function WINDOW_SLOTS() view returns(uint256)
func (gapCover *GapCover) UnpackWINDOWSLOTS(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("WINDOW_SLOTS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (gapCover *GapCover) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := gapCover.abi.Pack("allowance", owner, spender)
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
func (gapCover *GapCover) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return gapCover.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (gapCover *GapCover) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("allowance", data)
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
func (gapCover *GapCover) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := gapCover.abi.Pack("approve", spender, value)
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
func (gapCover *GapCover) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (gapCover *GapCover) UnpackApprove(data []byte) (bool, error) {
	out, err := gapCover.abi.Unpack("approve", data)
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
func (gapCover *GapCover) PackAsset() []byte {
	enc, err := gapCover.abi.Pack("asset")
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
func (gapCover *GapCover) TryPackAsset() ([]byte, error) {
	return gapCover.abi.Pack("asset")
}

// UnpackAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x38d52e0f.
//
// Solidity: function asset() view returns(address)
func (gapCover *GapCover) UnpackAsset(data []byte) (common.Address, error) {
	out, err := gapCover.abi.Unpack("asset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (gapCover *GapCover) PackBalanceOf(account common.Address) []byte {
	enc, err := gapCover.abi.Pack("balanceOf", account)
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
func (gapCover *GapCover) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return gapCover.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (gapCover *GapCover) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBand is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x10ea891a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function band() view returns(address)
func (gapCover *GapCover) PackBand() []byte {
	enc, err := gapCover.abi.Pack("band")
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
func (gapCover *GapCover) TryPackBand() ([]byte, error) {
	return gapCover.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (gapCover *GapCover) UnpackBand(data []byte) (common.Address, error) {
	out, err := gapCover.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBuy is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x99813cd5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function buy(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps, uint256 maxPremium, address holder) returns(uint256 id, uint256 premium)
func (gapCover *GapCover) PackBuy(symbol [32]byte, notional *big.Int, deductibleBps *big.Int, limitBps *big.Int, maxPremium *big.Int, holder common.Address) []byte {
	enc, err := gapCover.abi.Pack("buy", symbol, notional, deductibleBps, limitBps, maxPremium, holder)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBuy is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x99813cd5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function buy(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps, uint256 maxPremium, address holder) returns(uint256 id, uint256 premium)
func (gapCover *GapCover) TryPackBuy(symbol [32]byte, notional *big.Int, deductibleBps *big.Int, limitBps *big.Int, maxPremium *big.Int, holder common.Address) ([]byte, error) {
	return gapCover.abi.Pack("buy", symbol, notional, deductibleBps, limitBps, maxPremium, holder)
}

// BuyOutput serves as a container for the return parameters of contract
// method Buy.
type BuyOutput struct {
	Id      *big.Int
	Premium *big.Int
}

// UnpackBuy is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x99813cd5.
//
// Solidity: function buy(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps, uint256 maxPremium, address holder) returns(uint256 id, uint256 premium)
func (gapCover *GapCover) UnpackBuy(data []byte) (BuyOutput, error) {
	out, err := gapCover.abi.Unpack("buy", data)
	outstruct := new(BuyOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Id = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Premium = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackCapacity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5cfc1a51.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function capacity() view returns(uint256)
func (gapCover *GapCover) PackCapacity() []byte {
	enc, err := gapCover.abi.Pack("capacity")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCapacity is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5cfc1a51.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function capacity() view returns(uint256)
func (gapCover *GapCover) TryPackCapacity() ([]byte, error) {
	return gapCover.abi.Pack("capacity")
}

// UnpackCapacity is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5cfc1a51.
//
// Solidity: function capacity() view returns(uint256)
func (gapCover *GapCover) UnpackCapacity(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("capacity", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e83409a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim(address receiver) returns(uint256 amount)
func (gapCover *GapCover) PackClaim(receiver common.Address) []byte {
	enc, err := gapCover.abi.Pack("claim", receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1e83409a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim(address receiver) returns(uint256 amount)
func (gapCover *GapCover) TryPackClaim(receiver common.Address) ([]byte, error) {
	return gapCover.abi.Pack("claim", receiver)
}

// UnpackClaim is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1e83409a.
//
// Solidity: function claim(address receiver) returns(uint256 amount)
func (gapCover *GapCover) UnpackClaim(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("claim", data)
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
func (gapCover *GapCover) PackConvertToAssets(shares *big.Int) []byte {
	enc, err := gapCover.abi.Pack("convertToAssets", shares)
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
func (gapCover *GapCover) TryPackConvertToAssets(shares *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("convertToAssets", shares)
}

// UnpackConvertToAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 shares) view returns(uint256)
func (gapCover *GapCover) UnpackConvertToAssets(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("convertToAssets", data)
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
func (gapCover *GapCover) PackConvertToShares(assets *big.Int) []byte {
	enc, err := gapCover.abi.Pack("convertToShares", assets)
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
func (gapCover *GapCover) TryPackConvertToShares(assets *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("convertToShares", assets)
}

// UnpackConvertToShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 assets) view returns(uint256)
func (gapCover *GapCover) UnpackConvertToShares(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("convertToShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCoverCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeb0b8f5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function coverCount() view returns(uint256)
func (gapCover *GapCover) PackCoverCount() []byte {
	enc, err := gapCover.abi.Pack("coverCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCoverCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfeb0b8f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function coverCount() view returns(uint256)
func (gapCover *GapCover) TryPackCoverCount() ([]byte, error) {
	return gapCover.abi.Pack("coverCount")
}

// UnpackCoverCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfeb0b8f5.
//
// Solidity: function coverCount() view returns(uint256)
func (gapCover *GapCover) UnpackCoverCount(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("coverCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCovers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6299df6c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function covers(uint256 id) view returns(address holder, uint64 closesMs, uint16 deductibleBps, uint16 limitBps, bytes32 symbol, uint128 notional, uint128 premium)
func (gapCover *GapCover) PackCovers(id *big.Int) []byte {
	enc, err := gapCover.abi.Pack("covers", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCovers is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6299df6c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function covers(uint256 id) view returns(address holder, uint64 closesMs, uint16 deductibleBps, uint16 limitBps, bytes32 symbol, uint128 notional, uint128 premium)
func (gapCover *GapCover) TryPackCovers(id *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("covers", id)
}

// CoversOutput serves as a container for the return parameters of contract
// method Covers.
type CoversOutput struct {
	Holder        common.Address
	ClosesMs      uint64
	DeductibleBps uint16
	LimitBps      uint16
	Symbol        [32]byte
	Notional      *big.Int
	Premium       *big.Int
}

// UnpackCovers is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6299df6c.
//
// Solidity: function covers(uint256 id) view returns(address holder, uint64 closesMs, uint16 deductibleBps, uint16 limitBps, bytes32 symbol, uint128 notional, uint128 premium)
func (gapCover *GapCover) UnpackCovers(data []byte) (CoversOutput, error) {
	out, err := gapCover.abi.Unpack("covers", data)
	outstruct := new(CoversOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Holder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ClosesMs = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.DeductibleBps = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.LimitBps = *abi.ConvertType(out[3], new(uint16)).(*uint16)
	outstruct.Symbol = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.Notional = abi.ConvertType(out[5], new(big.Int)).(*big.Int)
	outstruct.Premium = abi.ConvertType(out[6], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() view returns(uint8)
func (gapCover *GapCover) PackDecimals() []byte {
	enc, err := gapCover.abi.Pack("decimals")
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
func (gapCover *GapCover) TryPackDecimals() ([]byte, error) {
	return gapCover.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (gapCover *GapCover) UnpackDecimals(data []byte) (uint8, error) {
	out, err := gapCover.abi.Unpack("decimals", data)
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
func (gapCover *GapCover) PackDeposit(assets *big.Int, receiver common.Address) []byte {
	enc, err := gapCover.abi.Pack("deposit", assets, receiver)
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
func (gapCover *GapCover) TryPackDeposit(assets *big.Int, receiver common.Address) ([]byte, error) {
	return gapCover.abi.Pack("deposit", assets, receiver)
}

// UnpackDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6e553f65.
//
// Solidity: function deposit(uint256 assets, address receiver) returns(uint256)
func (gapCover *GapCover) UnpackDeposit(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("deposit", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function engine() view returns(address)
func (gapCover *GapCover) PackEngine() []byte {
	enc, err := gapCover.abi.Pack("engine")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function engine() view returns(address)
func (gapCover *GapCover) TryPackEngine() ([]byte, error) {
	return gapCover.abi.Pack("engine")
}

// UnpackEngine is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (gapCover *GapCover) UnpackEngine(data []byte) (common.Address, error) {
	out, err := gapCover.abi.Unpack("engine", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackFeed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2679377f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function feed(bytes32 symbol) view returns(address)
func (gapCover *GapCover) PackFeed(symbol [32]byte) []byte {
	enc, err := gapCover.abi.Pack("feed", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackFeed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2679377f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function feed(bytes32 symbol) view returns(address)
func (gapCover *GapCover) TryPackFeed(symbol [32]byte) ([]byte, error) {
	return gapCover.abi.Pack("feed", symbol)
}

// UnpackFeed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2679377f.
//
// Solidity: function feed(bytes32 symbol) view returns(address)
func (gapCover *GapCover) UnpackFeed(data []byte) (common.Address, error) {
	out, err := gapCover.abi.Unpack("feed", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackHeld is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa285aed7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function held() view returns(uint256)
func (gapCover *GapCover) PackHeld() []byte {
	enc, err := gapCover.abi.Pack("held")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHeld is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa285aed7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function held() view returns(uint256)
func (gapCover *GapCover) TryPackHeld() ([]byte, error) {
	return gapCover.abi.Pack("held")
}

// UnpackHeld is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa285aed7.
//
// Solidity: function held() view returns(uint256)
func (gapCover *GapCover) UnpackHeld(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("held", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLastCloseMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x83b14afd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lastCloseMs() view returns(uint64)
func (gapCover *GapCover) PackLastCloseMs() []byte {
	enc, err := gapCover.abi.Pack("lastCloseMs")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLastCloseMs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x83b14afd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lastCloseMs() view returns(uint64)
func (gapCover *GapCover) TryPackLastCloseMs() ([]byte, error) {
	return gapCover.abi.Pack("lastCloseMs")
}

// UnpackLastCloseMs is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x83b14afd.
//
// Solidity: function lastCloseMs() view returns(uint64)
func (gapCover *GapCover) UnpackLastCloseMs(data []byte) (uint64, error) {
	out, err := gapCover.abi.Unpack("lastCloseMs", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackMaxDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x402d267d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (gapCover *GapCover) PackMaxDeposit(receiver common.Address) []byte {
	enc, err := gapCover.abi.Pack("maxDeposit", receiver)
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
func (gapCover *GapCover) TryPackMaxDeposit(receiver common.Address) ([]byte, error) {
	return gapCover.abi.Pack("maxDeposit", receiver)
}

// UnpackMaxDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x402d267d.
//
// Solidity: function maxDeposit(address receiver) view returns(uint256)
func (gapCover *GapCover) UnpackMaxDeposit(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("maxDeposit", data)
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
func (gapCover *GapCover) PackMaxMint(receiver common.Address) []byte {
	enc, err := gapCover.abi.Pack("maxMint", receiver)
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
func (gapCover *GapCover) TryPackMaxMint(receiver common.Address) ([]byte, error) {
	return gapCover.abi.Pack("maxMint", receiver)
}

// UnpackMaxMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc63d75b6.
//
// Solidity: function maxMint(address receiver) view returns(uint256)
func (gapCover *GapCover) UnpackMaxMint(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("maxMint", data)
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
func (gapCover *GapCover) PackMaxRedeem(owner common.Address) []byte {
	enc, err := gapCover.abi.Pack("maxRedeem", owner)
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
func (gapCover *GapCover) TryPackMaxRedeem(owner common.Address) ([]byte, error) {
	return gapCover.abi.Pack("maxRedeem", owner)
}

// UnpackMaxRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905777e.
//
// Solidity: function maxRedeem(address owner) view returns(uint256)
func (gapCover *GapCover) UnpackMaxRedeem(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("maxRedeem", data)
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
func (gapCover *GapCover) PackMaxWithdraw(owner common.Address) []byte {
	enc, err := gapCover.abi.Pack("maxWithdraw", owner)
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
func (gapCover *GapCover) TryPackMaxWithdraw(owner common.Address) ([]byte, error) {
	return gapCover.abi.Pack("maxWithdraw", owner)
}

// UnpackMaxWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xce96cb77.
//
// Solidity: function maxWithdraw(address owner) view returns(uint256)
func (gapCover *GapCover) UnpackMaxWithdraw(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("maxWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMeasure is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x25ee3714.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function measure(bytes32 symbol) returns(uint256 weekMove)
func (gapCover *GapCover) PackMeasure(symbol [32]byte) []byte {
	enc, err := gapCover.abi.Pack("measure", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMeasure is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x25ee3714.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function measure(bytes32 symbol) returns(uint256 weekMove)
func (gapCover *GapCover) TryPackMeasure(symbol [32]byte) ([]byte, error) {
	return gapCover.abi.Pack("measure", symbol)
}

// UnpackMeasure is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x25ee3714.
//
// Solidity: function measure(bytes32 symbol) returns(uint256 weekMove)
func (gapCover *GapCover) UnpackMeasure(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("measure", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMinDeductible is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a36cbf0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function minDeductible(bytes32 symbol) view returns(uint256)
func (gapCover *GapCover) PackMinDeductible(symbol [32]byte) []byte {
	enc, err := gapCover.abi.Pack("minDeductible", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMinDeductible is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9a36cbf0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function minDeductible(bytes32 symbol) view returns(uint256)
func (gapCover *GapCover) TryPackMinDeductible(symbol [32]byte) ([]byte, error) {
	return gapCover.abi.Pack("minDeductible", symbol)
}

// UnpackMinDeductible is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9a36cbf0.
//
// Solidity: function minDeductible(bytes32 symbol) view returns(uint256)
func (gapCover *GapCover) UnpackMinDeductible(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("minDeductible", data)
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
func (gapCover *GapCover) PackMint(shares *big.Int, receiver common.Address) []byte {
	enc, err := gapCover.abi.Pack("mint", shares, receiver)
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
func (gapCover *GapCover) TryPackMint(shares *big.Int, receiver common.Address) ([]byte, error) {
	return gapCover.abi.Pack("mint", shares, receiver)
}

// UnpackMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x94bf804d.
//
// Solidity: function mint(uint256 shares, address receiver) returns(uint256)
func (gapCover *GapCover) UnpackMint(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("mint", data)
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
func (gapCover *GapCover) PackName() []byte {
	enc, err := gapCover.abi.Pack("name")
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
func (gapCover *GapCover) TryPackName() ([]byte, error) {
	return gapCover.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (gapCover *GapCover) UnpackName(data []byte) (string, error) {
	out, err := gapCover.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackObserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8050f41a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function observe(bytes32 symbol, uint64 closesMs) returns(bool)
func (gapCover *GapCover) PackObserve(symbol [32]byte, closesMs uint64) []byte {
	enc, err := gapCover.abi.Pack("observe", symbol, closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackObserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8050f41a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function observe(bytes32 symbol, uint64 closesMs) returns(bool)
func (gapCover *GapCover) TryPackObserve(symbol [32]byte, closesMs uint64) ([]byte, error) {
	return gapCover.abi.Pack("observe", symbol, closesMs)
}

// UnpackObserve is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8050f41a.
//
// Solidity: function observe(bytes32 symbol, uint64 closesMs) returns(bool)
func (gapCover *GapCover) UnpackObserve(data []byte) (bool, error) {
	out, err := gapCover.abi.Unpack("observe", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackOutstanding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe50d33e3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function outstanding() view returns(uint256)
func (gapCover *GapCover) PackOutstanding() []byte {
	enc, err := gapCover.abi.Pack("outstanding")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOutstanding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe50d33e3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function outstanding() view returns(uint256)
func (gapCover *GapCover) TryPackOutstanding() ([]byte, error) {
	return gapCover.abi.Pack("outstanding")
}

// UnpackOutstanding is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe50d33e3.
//
// Solidity: function outstanding() view returns(uint256)
func (gapCover *GapCover) UnpackOutstanding(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("outstanding", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOutstandingIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3bd909c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function outstandingIn(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) PackOutstandingIn(closesMs uint64) []byte {
	enc, err := gapCover.abi.Pack("outstandingIn", closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOutstandingIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3bd909c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function outstandingIn(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) TryPackOutstandingIn(closesMs uint64) ([]byte, error) {
	return gapCover.abi.Pack("outstandingIn", closesMs)
}

// UnpackOutstandingIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3bd909c.
//
// Solidity: function outstandingIn(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) UnpackOutstandingIn(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("outstandingIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackOwed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87d6779.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owed() view returns(uint256)
func (gapCover *GapCover) PackOwed() []byte {
	enc, err := gapCover.abi.Pack("owed")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc87d6779.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owed() view returns(uint256)
func (gapCover *GapCover) TryPackOwed() ([]byte, error) {
	return gapCover.abi.Pack("owed")
}

// UnpackOwed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc87d6779.
//
// Solidity: function owed() view returns(uint256)
func (gapCover *GapCover) UnpackOwed(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("owed", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPayouts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65bcfbe7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function payouts(address holder) view returns(uint256)
func (gapCover *GapCover) PackPayouts(holder common.Address) []byte {
	enc, err := gapCover.abi.Pack("payouts", holder)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPayouts is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65bcfbe7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function payouts(address holder) view returns(uint256)
func (gapCover *GapCover) TryPackPayouts(holder common.Address) ([]byte, error) {
	return gapCover.abi.Pack("payouts", holder)
}

// UnpackPayouts is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x65bcfbe7.
//
// Solidity: function payouts(address holder) view returns(uint256)
func (gapCover *GapCover) UnpackPayouts(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("payouts", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPremiums is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde6a6740.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function premiums() view returns(uint256)
func (gapCover *GapCover) PackPremiums() []byte {
	enc, err := gapCover.abi.Pack("premiums")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPremiums is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde6a6740.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function premiums() view returns(uint256)
func (gapCover *GapCover) TryPackPremiums() ([]byte, error) {
	return gapCover.abi.Pack("premiums")
}

// UnpackPremiums is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xde6a6740.
//
// Solidity: function premiums() view returns(uint256)
func (gapCover *GapCover) UnpackPremiums(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("premiums", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPreviewDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xef8b30f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (gapCover *GapCover) PackPreviewDeposit(assets *big.Int) []byte {
	enc, err := gapCover.abi.Pack("previewDeposit", assets)
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
func (gapCover *GapCover) TryPackPreviewDeposit(assets *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("previewDeposit", assets)
}

// UnpackPreviewDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 assets) view returns(uint256)
func (gapCover *GapCover) UnpackPreviewDeposit(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("previewDeposit", data)
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
func (gapCover *GapCover) PackPreviewMint(shares *big.Int) []byte {
	enc, err := gapCover.abi.Pack("previewMint", shares)
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
func (gapCover *GapCover) TryPackPreviewMint(shares *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("previewMint", shares)
}

// UnpackPreviewMint is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb3d7f6b9.
//
// Solidity: function previewMint(uint256 shares) view returns(uint256)
func (gapCover *GapCover) UnpackPreviewMint(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("previewMint", data)
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
func (gapCover *GapCover) PackPreviewRedeem(shares *big.Int) []byte {
	enc, err := gapCover.abi.Pack("previewRedeem", shares)
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
func (gapCover *GapCover) TryPackPreviewRedeem(shares *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("previewRedeem", shares)
}

// UnpackPreviewRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 shares) view returns(uint256)
func (gapCover *GapCover) UnpackPreviewRedeem(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("previewRedeem", data)
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
func (gapCover *GapCover) PackPreviewWithdraw(assets *big.Int) []byte {
	enc, err := gapCover.abi.Pack("previewWithdraw", assets)
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
func (gapCover *GapCover) TryPackPreviewWithdraw(assets *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("previewWithdraw", assets)
}

// UnpackPreviewWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 assets) view returns(uint256)
func (gapCover *GapCover) UnpackPreviewWithdraw(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("previewWithdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPricingGap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaeb8186a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pricingGap(bytes32 symbol) view returns(uint256 gap, uint256 weekMove)
func (gapCover *GapCover) PackPricingGap(symbol [32]byte) []byte {
	enc, err := gapCover.abi.Pack("pricingGap", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPricingGap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaeb8186a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pricingGap(bytes32 symbol) view returns(uint256 gap, uint256 weekMove)
func (gapCover *GapCover) TryPackPricingGap(symbol [32]byte) ([]byte, error) {
	return gapCover.abi.Pack("pricingGap", symbol)
}

// PricingGapOutput serves as a container for the return parameters of contract
// method PricingGap.
type PricingGapOutput struct {
	Gap      *big.Int
	WeekMove *big.Int
}

// UnpackPricingGap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaeb8186a.
//
// Solidity: function pricingGap(bytes32 symbol) view returns(uint256 gap, uint256 weekMove)
func (gapCover *GapCover) UnpackPricingGap(data []byte) (PricingGapOutput, error) {
	out, err := gapCover.abi.Unpack("pricingGap", data)
	outstruct := new(PricingGapOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Gap = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.WeekMove = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5ef7e23.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function quote(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps) view returns(uint256)
func (gapCover *GapCover) PackQuote(symbol [32]byte, notional *big.Int, deductibleBps *big.Int, limitBps *big.Int) []byte {
	enc, err := gapCover.abi.Pack("quote", symbol, notional, deductibleBps, limitBps)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb5ef7e23.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function quote(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps) view returns(uint256)
func (gapCover *GapCover) TryPackQuote(symbol [32]byte, notional *big.Int, deductibleBps *big.Int, limitBps *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("quote", symbol, notional, deductibleBps, limitBps)
}

// UnpackQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb5ef7e23.
//
// Solidity: function quote(bytes32 symbol, uint256 notional, uint256 deductibleBps, uint256 limitBps) view returns(uint256)
func (gapCover *GapCover) UnpackQuote(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("quote", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRecord is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x266cf109.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function record() returns()
func (gapCover *GapCover) PackRecord() []byte {
	enc, err := gapCover.abi.Pack("record")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecord is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x266cf109.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function record() returns()
func (gapCover *GapCover) TryPackRecord() ([]byte, error) {
	return gapCover.abi.Pack("record")
}

// PackRedeem is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xba087652.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (gapCover *GapCover) PackRedeem(shares *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := gapCover.abi.Pack("redeem", shares, receiver, owner)
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
func (gapCover *GapCover) TryPackRedeem(shares *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return gapCover.abi.Pack("redeem", shares, receiver, owner)
}

// UnpackRedeem is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xba087652.
//
// Solidity: function redeem(uint256 shares, address receiver, address owner) returns(uint256)
func (gapCover *GapCover) UnpackRedeem(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("redeem", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRegularHours is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ff094a5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function regularHours() view returns(bool)
func (gapCover *GapCover) PackRegularHours() []byte {
	enc, err := gapCover.abi.Pack("regularHours")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegularHours is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1ff094a5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function regularHours() view returns(bool)
func (gapCover *GapCover) TryPackRegularHours() ([]byte, error) {
	return gapCover.abi.Pack("regularHours")
}

// UnpackRegularHours is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1ff094a5.
//
// Solidity: function regularHours() view returns(bool)
func (gapCover *GapCover) UnpackRegularHours(data []byte) (bool, error) {
	out, err := gapCover.abi.Unpack("regularHours", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRelease is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37bdc99b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function release(uint256 id) returns(uint256 payout, uint256 refund)
func (gapCover *GapCover) PackRelease(id *big.Int) []byte {
	enc, err := gapCover.abi.Pack("release", id)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRelease is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37bdc99b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function release(uint256 id) returns(uint256 payout, uint256 refund)
func (gapCover *GapCover) TryPackRelease(id *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("release", id)
}

// ReleaseOutput serves as a container for the return parameters of contract
// method Release.
type ReleaseOutput struct {
	Payout *big.Int
	Refund *big.Int
}

// UnpackRelease is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x37bdc99b.
//
// Solidity: function release(uint256 id) returns(uint256 payout, uint256 refund)
func (gapCover *GapCover) UnpackRelease(data []byte) (ReleaseOutput, error) {
	out, err := gapCover.abi.Unpack("release", data)
	outstruct := new(ReleaseOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Payout = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Refund = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackReopenOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x777b4a57.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reopenOf(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) PackReopenOf(closesMs uint64) []byte {
	enc, err := gapCover.abi.Pack("reopenOf", closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReopenOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x777b4a57.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reopenOf(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) TryPackReopenOf(closesMs uint64) ([]byte, error) {
	return gapCover.abi.Pack("reopenOf", closesMs)
}

// UnpackReopenOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x777b4a57.
//
// Solidity: function reopenOf(uint64 closesMs) view returns(uint256)
func (gapCover *GapCover) UnpackReopenOf(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("reopenOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackReserved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe60d12c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reserved() view returns(uint256)
func (gapCover *GapCover) PackReserved() []byte {
	enc, err := gapCover.abi.Pack("reserved")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReserved is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe60d12c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reserved() view returns(uint256)
func (gapCover *GapCover) TryPackReserved() ([]byte, error) {
	return gapCover.abi.Pack("reserved")
}

// UnpackReserved is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe60d12c.
//
// Solidity: function reserved() view returns(uint256)
func (gapCover *GapCover) UnpackReserved(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("reserved", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSales is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaace52fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sales() view returns(uint64 closesMs, uint64 endsMs)
func (gapCover *GapCover) PackSales() []byte {
	enc, err := gapCover.abi.Pack("sales")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSales is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xaace52fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sales() view returns(uint64 closesMs, uint64 endsMs)
func (gapCover *GapCover) TryPackSales() ([]byte, error) {
	return gapCover.abi.Pack("sales")
}

// SalesOutput serves as a container for the return parameters of contract
// method Sales.
type SalesOutput struct {
	ClosesMs uint64
	EndsMs   uint64
}

// UnpackSales is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xaace52fe.
//
// Solidity: function sales() view returns(uint64 closesMs, uint64 endsMs)
func (gapCover *GapCover) UnpackSales(data []byte) (SalesOutput, error) {
	out, err := gapCover.abi.Unpack("sales", data)
	outstruct := new(SalesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ClosesMs = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.EndsMs = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackSeries is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64a3f057.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function series(bytes32 symbol, uint64 closesMs) view returns(uint128 notional, uint64 referencePrice, uint64 price, uint16 shift, uint8 status, bool flagged)
func (gapCover *GapCover) PackSeries(symbol [32]byte, closesMs uint64) []byte {
	enc, err := gapCover.abi.Pack("series", symbol, closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeries is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64a3f057.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function series(bytes32 symbol, uint64 closesMs) view returns(uint128 notional, uint64 referencePrice, uint64 price, uint16 shift, uint8 status, bool flagged)
func (gapCover *GapCover) TryPackSeries(symbol [32]byte, closesMs uint64) ([]byte, error) {
	return gapCover.abi.Pack("series", symbol, closesMs)
}

// SeriesOutput serves as a container for the return parameters of contract
// method Series.
type SeriesOutput struct {
	Notional       *big.Int
	ReferencePrice uint64
	Price          uint64
	Shift          uint16
	Status         uint8
	Flagged        bool
}

// UnpackSeries is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x64a3f057.
//
// Solidity: function series(bytes32 symbol, uint64 closesMs) view returns(uint128 notional, uint64 referencePrice, uint64 price, uint16 shift, uint8 status, bool flagged)
func (gapCover *GapCover) UnpackSeries(data []byte) (SeriesOutput, error) {
	out, err := gapCover.abi.Unpack("series", data)
	outstruct := new(SeriesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Notional = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.ReferencePrice = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Price = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Shift = *abi.ConvertType(out[3], new(uint16)).(*uint16)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Flagged = *abi.ConvertType(out[5], new(bool)).(*bool)
	return *outstruct, nil
}

// PackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47a54490.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function settle(bytes32 symbol, uint64 closesMs, uint80 referenceRound, uint80 lastRound) returns()
func (gapCover *GapCover) PackSettle(symbol [32]byte, closesMs uint64, referenceRound *big.Int, lastRound *big.Int) []byte {
	enc, err := gapCover.abi.Pack("settle", symbol, closesMs, referenceRound, lastRound)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47a54490.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function settle(bytes32 symbol, uint64 closesMs, uint80 referenceRound, uint80 lastRound) returns()
func (gapCover *GapCover) TryPackSettle(symbol [32]byte, closesMs uint64, referenceRound *big.Int, lastRound *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("settle", symbol, closesMs, referenceRound, lastRound)
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (gapCover *GapCover) PackSymbol() []byte {
	enc, err := gapCover.abi.Pack("symbol")
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
func (gapCover *GapCover) TryPackSymbol() ([]byte, error) {
	return gapCover.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (gapCover *GapCover) UnpackSymbol(data []byte) (string, error) {
	out, err := gapCover.abi.Unpack("symbol", data)
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
func (gapCover *GapCover) PackSync() []byte {
	enc, err := gapCover.abi.Pack("sync")
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
func (gapCover *GapCover) TryPackSync() ([]byte, error) {
	return gapCover.abi.Pack("sync")
}

// PackTotalAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01e1d114.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalAssets() view returns(uint256)
func (gapCover *GapCover) PackTotalAssets() []byte {
	enc, err := gapCover.abi.Pack("totalAssets")
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
func (gapCover *GapCover) TryPackTotalAssets() ([]byte, error) {
	return gapCover.abi.Pack("totalAssets")
}

// UnpackTotalAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (gapCover *GapCover) UnpackTotalAssets(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("totalAssets", data)
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
func (gapCover *GapCover) PackTotalSupply() []byte {
	enc, err := gapCover.abi.Pack("totalSupply")
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
func (gapCover *GapCover) TryPackTotalSupply() ([]byte, error) {
	return gapCover.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (gapCover *GapCover) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("totalSupply", data)
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
func (gapCover *GapCover) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := gapCover.abi.Pack("transfer", to, value)
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
func (gapCover *GapCover) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (gapCover *GapCover) UnpackTransfer(data []byte) (bool, error) {
	out, err := gapCover.abi.Unpack("transfer", data)
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
func (gapCover *GapCover) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := gapCover.abi.Pack("transferFrom", from, to, value)
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
func (gapCover *GapCover) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return gapCover.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (gapCover *GapCover) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := gapCover.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackVoid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf1d40ad8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function void(bytes32 symbol, uint64 closesMs) returns()
func (gapCover *GapCover) PackVoid(symbol [32]byte, closesMs uint64) []byte {
	enc, err := gapCover.abi.Pack("void", symbol, closesMs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVoid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf1d40ad8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function void(bytes32 symbol, uint64 closesMs) returns()
func (gapCover *GapCover) TryPackVoid(symbol [32]byte, closesMs uint64) ([]byte, error) {
	return gapCover.abi.Pack("void", symbol, closesMs)
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb460af94.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (gapCover *GapCover) PackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) []byte {
	enc, err := gapCover.abi.Pack("withdraw", assets, receiver, owner)
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
func (gapCover *GapCover) TryPackWithdraw(assets *big.Int, receiver common.Address, owner common.Address) ([]byte, error) {
	return gapCover.abi.Pack("withdraw", assets, receiver, owner)
}

// UnpackWithdraw is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb460af94.
//
// Solidity: function withdraw(uint256 assets, address receiver, address owner) returns(uint256)
func (gapCover *GapCover) UnpackWithdraw(data []byte) (*big.Int, error) {
	out, err := gapCover.abi.Unpack("withdraw", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// GapCoverApproval represents a Approval event raised by the GapCover contract.
type GapCoverApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const GapCoverApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (GapCoverApproval) ContractEventName() string {
	return GapCoverApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (gapCover *GapCover) UnpackApprovalEvent(log *types.Log) (*GapCoverApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverApproval)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverBought represents a Bought event raised by the GapCover contract.
type GapCoverBought struct {
	Id            *big.Int
	Holder        common.Address
	Symbol        [32]byte
	ClosesMs      uint64
	Notional      *big.Int
	DeductibleBps *big.Int
	LimitBps      *big.Int
	Premium       *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const GapCoverBoughtEventName = "Bought"

// ContractEventName returns the user-defined event name.
func (GapCoverBought) ContractEventName() string {
	return GapCoverBoughtEventName
}

// UnpackBoughtEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Bought(uint256 indexed id, address indexed holder, bytes32 indexed symbol, uint64 closesMs, uint256 notional, uint256 deductibleBps, uint256 limitBps, uint256 premium)
func (gapCover *GapCover) UnpackBoughtEvent(log *types.Log) (*GapCoverBought, error) {
	event := "Bought"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverBought)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverClaimed represents a Claimed event raised by the GapCover contract.
type GapCoverClaimed struct {
	Holder   common.Address
	Receiver common.Address
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverClaimedEventName = "Claimed"

// ContractEventName returns the user-defined event name.
func (GapCoverClaimed) ContractEventName() string {
	return GapCoverClaimedEventName
}

// UnpackClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Claimed(address indexed holder, address indexed receiver, uint256 amount)
func (gapCover *GapCover) UnpackClaimedEvent(log *types.Log) (*GapCoverClaimed, error) {
	event := "Claimed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverClaimed)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverCloseRecorded represents a CloseRecorded event raised by the GapCover contract.
type GapCoverCloseRecorded struct {
	ClosesMs uint64
	AtMs     uint64
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverCloseRecordedEventName = "CloseRecorded"

// ContractEventName returns the user-defined event name.
func (GapCoverCloseRecorded) ContractEventName() string {
	return GapCoverCloseRecordedEventName
}

// UnpackCloseRecordedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CloseRecorded(uint64 indexed closesMs, uint64 atMs)
func (gapCover *GapCover) UnpackCloseRecordedEvent(log *types.Log) (*GapCoverCloseRecorded, error) {
	event := "CloseRecorded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverCloseRecorded)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverDeposit represents a Deposit event raised by the GapCover contract.
type GapCoverDeposit struct {
	Sender common.Address
	Owner  common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapCoverDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (GapCoverDeposit) ContractEventName() string {
	return GapCoverDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed sender, address indexed owner, uint256 assets, uint256 shares)
func (gapCover *GapCover) UnpackDepositEvent(log *types.Log) (*GapCoverDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverDeposit)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverMeasured represents a Measured event raised by the GapCover contract.
type GapCoverMeasured struct {
	Symbol   [32]byte
	ClosesMs uint64
	WeekMove *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverMeasuredEventName = "Measured"

// ContractEventName returns the user-defined event name.
func (GapCoverMeasured) ContractEventName() string {
	return GapCoverMeasuredEventName
}

// UnpackMeasuredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Measured(bytes32 indexed symbol, uint64 indexed closesMs, uint256 weekMove)
func (gapCover *GapCover) UnpackMeasuredEvent(log *types.Log) (*GapCoverMeasured, error) {
	event := "Measured"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverMeasured)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverObserved represents a Observed event raised by the GapCover contract.
type GapCoverObserved struct {
	Symbol   [32]byte
	ClosesMs uint64
	Slot     *big.Int
	Mid      uint64
	Low      uint64
	High     *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverObservedEventName = "Observed"

// ContractEventName returns the user-defined event name.
func (GapCoverObserved) ContractEventName() string {
	return GapCoverObservedEventName
}

// UnpackObservedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Observed(bytes32 indexed symbol, uint64 indexed closesMs, uint256 slot, uint64 mid, uint64 low, uint128 high)
func (gapCover *GapCover) UnpackObservedEvent(log *types.Log) (*GapCoverObserved, error) {
	event := "Observed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverObserved)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverReleased represents a Released event raised by the GapCover contract.
type GapCoverReleased struct {
	Id     *big.Int
	Holder common.Address
	Payout *big.Int
	Refund *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const GapCoverReleasedEventName = "Released"

// ContractEventName returns the user-defined event name.
func (GapCoverReleased) ContractEventName() string {
	return GapCoverReleasedEventName
}

// UnpackReleasedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Released(uint256 indexed id, address indexed holder, uint256 payout, uint256 refund)
func (gapCover *GapCover) UnpackReleasedEvent(log *types.Log) (*GapCoverReleased, error) {
	event := "Released"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverReleased)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverReopenRecorded represents a ReopenRecorded event raised by the GapCover contract.
type GapCoverReopenRecorded struct {
	ClosesMs  uint64
	ReopensMs uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const GapCoverReopenRecordedEventName = "ReopenRecorded"

// ContractEventName returns the user-defined event name.
func (GapCoverReopenRecorded) ContractEventName() string {
	return GapCoverReopenRecordedEventName
}

// UnpackReopenRecordedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReopenRecorded(uint64 indexed closesMs, uint64 reopensMs)
func (gapCover *GapCover) UnpackReopenRecordedEvent(log *types.Log) (*GapCoverReopenRecorded, error) {
	event := "ReopenRecorded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverReopenRecorded)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverSettled represents a Settled event raised by the GapCover contract.
type GapCoverSettled struct {
	Symbol         [32]byte
	ClosesMs       uint64
	ReferenceRound *big.Int
	ReferencePrice *big.Int
	FirstRound     *big.Int
	Price          *big.Int
	Flagged        bool
	Raw            *types.Log // Blockchain specific contextual infos
}

const GapCoverSettledEventName = "Settled"

// ContractEventName returns the user-defined event name.
func (GapCoverSettled) ContractEventName() string {
	return GapCoverSettledEventName
}

// UnpackSettledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Settled(bytes32 indexed symbol, uint64 indexed closesMs, uint80 referenceRound, uint256 referencePrice, uint80 firstRound, uint256 price, bool flagged)
func (gapCover *GapCover) UnpackSettledEvent(log *types.Log) (*GapCoverSettled, error) {
	event := "Settled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverSettled)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverSync represents a Sync event raised by the GapCover contract.
type GapCoverSync struct {
	Held *big.Int
	Raw  *types.Log // Blockchain specific contextual infos
}

const GapCoverSyncEventName = "Sync"

// ContractEventName returns the user-defined event name.
func (GapCoverSync) ContractEventName() string {
	return GapCoverSyncEventName
}

// UnpackSyncEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Sync(uint256 held)
func (gapCover *GapCover) UnpackSyncEvent(log *types.Log) (*GapCoverSync, error) {
	event := "Sync"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverSync)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverTransfer represents a Transfer event raised by the GapCover contract.
type GapCoverTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const GapCoverTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (GapCoverTransfer) ContractEventName() string {
	return GapCoverTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (gapCover *GapCover) UnpackTransferEvent(log *types.Log) (*GapCoverTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverTransfer)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverVoided represents a Voided event raised by the GapCover contract.
type GapCoverVoided struct {
	Symbol   [32]byte
	ClosesMs uint64
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverVoidedEventName = "Voided"

// ContractEventName returns the user-defined event name.
func (GapCoverVoided) ContractEventName() string {
	return GapCoverVoidedEventName
}

// UnpackVoidedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Voided(bytes32 indexed symbol, uint64 indexed closesMs)
func (gapCover *GapCover) UnpackVoidedEvent(log *types.Log) (*GapCoverVoided, error) {
	event := "Voided"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverVoided)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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

// GapCoverWithdraw represents a Withdraw event raised by the GapCover contract.
type GapCoverWithdraw struct {
	Sender   common.Address
	Receiver common.Address
	Owner    common.Address
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const GapCoverWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (GapCoverWithdraw) ContractEventName() string {
	return GapCoverWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed sender, address indexed receiver, address indexed owner, uint256 assets, uint256 shares)
func (gapCover *GapCover) UnpackWithdrawEvent(log *types.Log) (*GapCoverWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != gapCover.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(GapCoverWithdraw)
	if len(log.Data) > 0 {
		if err := gapCover.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range gapCover.abi.Events[event].Inputs {
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
func (gapCover *GapCover) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], gapCover.abi.Errors["AssetHalted"].ID.Bytes()[:4]) {
		return gapCover.UnpackAssetHaltedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["CorporateActionPending"].ID.Bytes()[:4]) {
		return gapCover.UnpackCorporateActionPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC4626ExceededMaxDeposit"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC4626ExceededMaxDepositError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC4626ExceededMaxMint"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC4626ExceededMaxMintError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC4626ExceededMaxRedeem"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC4626ExceededMaxRedeemError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ERC4626ExceededMaxWithdraw"].ID.Bytes()[:4]) {
		return gapCover.UnpackERC4626ExceededMaxWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["InsufficientGas"].ID.Bytes()[:4]) {
		return gapCover.UnpackInsufficientGasError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["InvalidHolder"].ID.Bytes()[:4]) {
		return gapCover.UnpackInvalidHolderError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["InvalidLastRound"].ID.Bytes()[:4]) {
		return gapCover.UnpackInvalidLastRoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["InvalidLayer"].ID.Bytes()[:4]) {
		return gapCover.UnpackInvalidLayerError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["InvalidReference"].ID.Bytes()[:4]) {
		return gapCover.UnpackInvalidReferenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NoCapacity"].ID.Bytes()[:4]) {
		return gapCover.UnpackNoCapacityError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NoCover"].ID.Bytes()[:4]) {
		return gapCover.UnpackNoCoverError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NotCoverable"].ID.Bytes()[:4]) {
		return gapCover.UnpackNotCoverableError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NotReopened"].ID.Bytes()[:4]) {
		return gapCover.UnpackNotReopenedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NotSettled"].ID.Bytes()[:4]) {
		return gapCover.UnpackNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["NothingCovered"].ID.Bytes()[:4]) {
		return gapCover.UnpackNothingCoveredError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["PremiumAboveLimit"].ID.Bytes()[:4]) {
		return gapCover.UnpackPremiumAboveLimitError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SafeCastOverflowedIntToUint"].ID.Bytes()[:4]) {
		return gapCover.UnpackSafeCastOverflowedIntToUintError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return gapCover.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return gapCover.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SalesClosed"].ID.Bytes()[:4]) {
		return gapCover.UnpackSalesClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return gapCover.UnpackSequencerNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["SeriesClosed"].ID.Bytes()[:4]) {
		return gapCover.UnpackSeriesClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["StaleReference"].ID.Bytes()[:4]) {
		return gapCover.UnpackStaleReferenceError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["TooEarlyToMeasure"].ID.Bytes()[:4]) {
		return gapCover.UnpackTooEarlyToMeasureError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["TooEarlyToVoid"].ID.Bytes()[:4]) {
		return gapCover.UnpackTooEarlyToVoidError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["WindowClosed"].ID.Bytes()[:4]) {
		return gapCover.UnpackWindowClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["WindowOpen"].ID.Bytes()[:4]) {
		return gapCover.UnpackWindowOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], gapCover.abi.Errors["ZeroAmount"].ID.Bytes()[:4]) {
		return gapCover.UnpackZeroAmountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// GapCoverAssetHalted represents a AssetHalted error raised by the GapCover contract.
type GapCoverAssetHalted struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetHalted(bytes32 symbol)
func GapCoverAssetHaltedErrorID() common.Hash {
	return common.HexToHash("0x3ec29ca3785f47585a9c96baed33ddbd8b2d9e8dfdeffb4be96925aaf9ec6468")
}

// UnpackAssetHaltedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetHalted(bytes32 symbol)
func (gapCover *GapCover) UnpackAssetHaltedError(raw []byte) (*GapCoverAssetHalted, error) {
	out := new(GapCoverAssetHalted)
	if err := gapCover.abi.UnpackIntoInterface(out, "AssetHalted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverCorporateActionPending represents a CorporateActionPending error raised by the GapCover contract.
type GapCoverCorporateActionPending struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func GapCoverCorporateActionPendingErrorID() common.Hash {
	return common.HexToHash("0x1d1f5fd3c689a2e2561007ef605d914d864ac52012393f43a887a76f9f48be7b")
}

// UnpackCorporateActionPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func (gapCover *GapCover) UnpackCorporateActionPendingError(raw []byte) (*GapCoverCorporateActionPending, error) {
	out := new(GapCoverCorporateActionPending)
	if err := gapCover.abi.UnpackIntoInterface(out, "CorporateActionPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the GapCover contract.
type GapCoverERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func GapCoverERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (gapCover *GapCover) UnpackERC20InsufficientAllowanceError(raw []byte) (*GapCoverERC20InsufficientAllowance, error) {
	out := new(GapCoverERC20InsufficientAllowance)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the GapCover contract.
type GapCoverERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func GapCoverERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (gapCover *GapCover) UnpackERC20InsufficientBalanceError(raw []byte) (*GapCoverERC20InsufficientBalance, error) {
	out := new(GapCoverERC20InsufficientBalance)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InvalidApprover represents a ERC20InvalidApprover error raised by the GapCover contract.
type GapCoverERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func GapCoverERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (gapCover *GapCover) UnpackERC20InvalidApproverError(raw []byte) (*GapCoverERC20InvalidApprover, error) {
	out := new(GapCoverERC20InvalidApprover)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the GapCover contract.
type GapCoverERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func GapCoverERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (gapCover *GapCover) UnpackERC20InvalidReceiverError(raw []byte) (*GapCoverERC20InvalidReceiver, error) {
	out := new(GapCoverERC20InvalidReceiver)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InvalidSender represents a ERC20InvalidSender error raised by the GapCover contract.
type GapCoverERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func GapCoverERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (gapCover *GapCover) UnpackERC20InvalidSenderError(raw []byte) (*GapCoverERC20InvalidSender, error) {
	out := new(GapCoverERC20InvalidSender)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC20InvalidSpender represents a ERC20InvalidSpender error raised by the GapCover contract.
type GapCoverERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func GapCoverERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (gapCover *GapCover) UnpackERC20InvalidSpenderError(raw []byte) (*GapCoverERC20InvalidSpender, error) {
	out := new(GapCoverERC20InvalidSpender)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC4626ExceededMaxDeposit represents a ERC4626ExceededMaxDeposit error raised by the GapCover contract.
type GapCoverERC4626ExceededMaxDeposit struct {
	Receiver common.Address
	Assets   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func GapCoverERC4626ExceededMaxDepositErrorID() common.Hash {
	return common.HexToHash("0x79012fb2819fcdc1de669c08773dbcd6bdc757862642f75fb1c584dadf259dfe")
}

// UnpackERC4626ExceededMaxDepositError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxDeposit(address receiver, uint256 assets, uint256 max)
func (gapCover *GapCover) UnpackERC4626ExceededMaxDepositError(raw []byte) (*GapCoverERC4626ExceededMaxDeposit, error) {
	out := new(GapCoverERC4626ExceededMaxDeposit)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxDeposit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC4626ExceededMaxMint represents a ERC4626ExceededMaxMint error raised by the GapCover contract.
type GapCoverERC4626ExceededMaxMint struct {
	Receiver common.Address
	Shares   *big.Int
	Max      *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func GapCoverERC4626ExceededMaxMintErrorID() common.Hash {
	return common.HexToHash("0x284ff667dc615a39438518c22e8955b9470327d9de8a4d7e21c926b260d65176")
}

// UnpackERC4626ExceededMaxMintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxMint(address receiver, uint256 shares, uint256 max)
func (gapCover *GapCover) UnpackERC4626ExceededMaxMintError(raw []byte) (*GapCoverERC4626ExceededMaxMint, error) {
	out := new(GapCoverERC4626ExceededMaxMint)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxMint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC4626ExceededMaxRedeem represents a ERC4626ExceededMaxRedeem error raised by the GapCover contract.
type GapCoverERC4626ExceededMaxRedeem struct {
	Owner  common.Address
	Shares *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func GapCoverERC4626ExceededMaxRedeemErrorID() common.Hash {
	return common.HexToHash("0xb94abeec0557d36b5f0bc8f115deec7b184dcbff94ac66d55e37c8f301e75269")
}

// UnpackERC4626ExceededMaxRedeemError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxRedeem(address owner, uint256 shares, uint256 max)
func (gapCover *GapCover) UnpackERC4626ExceededMaxRedeemError(raw []byte) (*GapCoverERC4626ExceededMaxRedeem, error) {
	out := new(GapCoverERC4626ExceededMaxRedeem)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxRedeem", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverERC4626ExceededMaxWithdraw represents a ERC4626ExceededMaxWithdraw error raised by the GapCover contract.
type GapCoverERC4626ExceededMaxWithdraw struct {
	Owner  common.Address
	Assets *big.Int
	Max    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func GapCoverERC4626ExceededMaxWithdrawErrorID() common.Hash {
	return common.HexToHash("0xfe9cceec2bd1f9b68641914cc354eacaeb1cc2169f5ba0639930f241e87142f0")
}

// UnpackERC4626ExceededMaxWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC4626ExceededMaxWithdraw(address owner, uint256 assets, uint256 max)
func (gapCover *GapCover) UnpackERC4626ExceededMaxWithdrawError(raw []byte) (*GapCoverERC4626ExceededMaxWithdraw, error) {
	out := new(GapCoverERC4626ExceededMaxWithdraw)
	if err := gapCover.abi.UnpackIntoInterface(out, "ERC4626ExceededMaxWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverInsufficientGas represents a InsufficientGas error raised by the GapCover contract.
type GapCoverInsufficientGas struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientGas()
func GapCoverInsufficientGasErrorID() common.Hash {
	return common.HexToHash("0x1c26714c035362f785e44b6773a87881da86705e5698c55bb550907617bb1353")
}

// UnpackInsufficientGasError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientGas()
func (gapCover *GapCover) UnpackInsufficientGasError(raw []byte) (*GapCoverInsufficientGas, error) {
	out := new(GapCoverInsufficientGas)
	if err := gapCover.abi.UnpackIntoInterface(out, "InsufficientGas", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverInvalidHolder represents a InvalidHolder error raised by the GapCover contract.
type GapCoverInvalidHolder struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidHolder()
func GapCoverInvalidHolderErrorID() common.Hash {
	return common.HexToHash("0x4971ba2d20f35ff7821f0f08772569203e04eae3391179015fff219711a04a3d")
}

// UnpackInvalidHolderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidHolder()
func (gapCover *GapCover) UnpackInvalidHolderError(raw []byte) (*GapCoverInvalidHolder, error) {
	out := new(GapCoverInvalidHolder)
	if err := gapCover.abi.UnpackIntoInterface(out, "InvalidHolder", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverInvalidLastRound represents a InvalidLastRound error raised by the GapCover contract.
type GapCoverInvalidLastRound struct {
	RoundId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidLastRound(uint80 roundId)
func GapCoverInvalidLastRoundErrorID() common.Hash {
	return common.HexToHash("0xd8d0cdc972fb593130ea470ea52d243674471dfc38afa627ef0fbe1b8da75097")
}

// UnpackInvalidLastRoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidLastRound(uint80 roundId)
func (gapCover *GapCover) UnpackInvalidLastRoundError(raw []byte) (*GapCoverInvalidLastRound, error) {
	out := new(GapCoverInvalidLastRound)
	if err := gapCover.abi.UnpackIntoInterface(out, "InvalidLastRound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverInvalidLayer represents a InvalidLayer error raised by the GapCover contract.
type GapCoverInvalidLayer struct {
	DeductibleBps *big.Int
	LimitBps      *big.Int
	Minimum       *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidLayer(uint256 deductibleBps, uint256 limitBps, uint256 minimum)
func GapCoverInvalidLayerErrorID() common.Hash {
	return common.HexToHash("0x2a3c197f04d3f6eab0687333eca3e3d73680a4a98e10dc60b6d2c03880b2781d")
}

// UnpackInvalidLayerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidLayer(uint256 deductibleBps, uint256 limitBps, uint256 minimum)
func (gapCover *GapCover) UnpackInvalidLayerError(raw []byte) (*GapCoverInvalidLayer, error) {
	out := new(GapCoverInvalidLayer)
	if err := gapCover.abi.UnpackIntoInterface(out, "InvalidLayer", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverInvalidReference represents a InvalidReference error raised by the GapCover contract.
type GapCoverInvalidReference struct {
	RoundId *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidReference(uint80 roundId)
func GapCoverInvalidReferenceErrorID() common.Hash {
	return common.HexToHash("0x4b8e8dcd8b8bc8d73201e80c0a80fad0fea3af339d8746a0e8637feb5839b44e")
}

// UnpackInvalidReferenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidReference(uint80 roundId)
func (gapCover *GapCover) UnpackInvalidReferenceError(raw []byte) (*GapCoverInvalidReference, error) {
	out := new(GapCoverInvalidReference)
	if err := gapCover.abi.UnpackIntoInterface(out, "InvalidReference", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNoCapacity represents a NoCapacity error raised by the GapCover contract.
type GapCoverNoCapacity struct {
	Need *big.Int
	Free *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoCapacity(uint256 need, uint256 free)
func GapCoverNoCapacityErrorID() common.Hash {
	return common.HexToHash("0xa25dc6ad4ec1d10e93da632ceeffb5e51e471ad23ff3f7d71f10e3cd418cc3aa")
}

// UnpackNoCapacityError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoCapacity(uint256 need, uint256 free)
func (gapCover *GapCover) UnpackNoCapacityError(raw []byte) (*GapCoverNoCapacity, error) {
	out := new(GapCoverNoCapacity)
	if err := gapCover.abi.UnpackIntoInterface(out, "NoCapacity", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNoCover represents a NoCover error raised by the GapCover contract.
type GapCoverNoCover struct {
	Id *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoCover(uint256 id)
func GapCoverNoCoverErrorID() common.Hash {
	return common.HexToHash("0xa18df25f578c47a77ab602884913fe0266475536fadedbef46389c6e98deaffd")
}

// UnpackNoCoverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoCover(uint256 id)
func (gapCover *GapCover) UnpackNoCoverError(raw []byte) (*GapCoverNoCover, error) {
	out := new(GapCoverNoCover)
	if err := gapCover.abi.UnpackIntoInterface(out, "NoCover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNotCoverable represents a NotCoverable error raised by the GapCover contract.
type GapCoverNotCoverable struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotCoverable(bytes32 symbol)
func GapCoverNotCoverableErrorID() common.Hash {
	return common.HexToHash("0x804a602ddeedf7b57d7b3ea9425370b41e9eaa542765f6fe05b4ea3b8d55ea7d")
}

// UnpackNotCoverableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotCoverable(bytes32 symbol)
func (gapCover *GapCover) UnpackNotCoverableError(raw []byte) (*GapCoverNotCoverable, error) {
	out := new(GapCoverNotCoverable)
	if err := gapCover.abi.UnpackIntoInterface(out, "NotCoverable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNotReopened represents a NotReopened error raised by the GapCover contract.
type GapCoverNotReopened struct {
	ClosesMs uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotReopened(uint64 closesMs)
func GapCoverNotReopenedErrorID() common.Hash {
	return common.HexToHash("0x667a563cc3cbefeee2bf91c9ee5b203ed59390c9be28acbfd004845697dbd0e2")
}

// UnpackNotReopenedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotReopened(uint64 closesMs)
func (gapCover *GapCover) UnpackNotReopenedError(raw []byte) (*GapCoverNotReopened, error) {
	out := new(GapCoverNotReopened)
	if err := gapCover.abi.UnpackIntoInterface(out, "NotReopened", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNotSettled represents a NotSettled error raised by the GapCover contract.
type GapCoverNotSettled struct {
	Id *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotSettled(uint256 id)
func GapCoverNotSettledErrorID() common.Hash {
	return common.HexToHash("0x38aa087c50b23897973148853603af2b0fa453a2b022b6c3c04f3ed1264df33b")
}

// UnpackNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotSettled(uint256 id)
func (gapCover *GapCover) UnpackNotSettledError(raw []byte) (*GapCoverNotSettled, error) {
	out := new(GapCoverNotSettled)
	if err := gapCover.abi.UnpackIntoInterface(out, "NotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverNothingCovered represents a NothingCovered error raised by the GapCover contract.
type GapCoverNothingCovered struct {
	Symbol   [32]byte
	ClosesMs uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NothingCovered(bytes32 symbol, uint64 closesMs)
func GapCoverNothingCoveredErrorID() common.Hash {
	return common.HexToHash("0x110d7b111b74a59d102d9afa60bbf0d0c4efe3028219d1d78c0b9fe234dc68f9")
}

// UnpackNothingCoveredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NothingCovered(bytes32 symbol, uint64 closesMs)
func (gapCover *GapCover) UnpackNothingCoveredError(raw []byte) (*GapCoverNothingCovered, error) {
	out := new(GapCoverNothingCovered)
	if err := gapCover.abi.UnpackIntoInterface(out, "NothingCovered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverPremiumAboveLimit represents a PremiumAboveLimit error raised by the GapCover contract.
type GapCoverPremiumAboveLimit struct {
	Premium    *big.Int
	MaxPremium *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PremiumAboveLimit(uint256 premium, uint256 maxPremium)
func GapCoverPremiumAboveLimitErrorID() common.Hash {
	return common.HexToHash("0x8183f41e96355ccb30a97737b291165a5e26e6b027ed3520237e5af207b97704")
}

// UnpackPremiumAboveLimitError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PremiumAboveLimit(uint256 premium, uint256 maxPremium)
func (gapCover *GapCover) UnpackPremiumAboveLimitError(raw []byte) (*GapCoverPremiumAboveLimit, error) {
	out := new(GapCoverPremiumAboveLimit)
	if err := gapCover.abi.UnpackIntoInterface(out, "PremiumAboveLimit", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSafeCastOverflowedIntToUint represents a SafeCastOverflowedIntToUint error raised by the GapCover contract.
type GapCoverSafeCastOverflowedIntToUint struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func GapCoverSafeCastOverflowedIntToUintErrorID() common.Hash {
	return common.HexToHash("0xa8ce4432b175c373e5f41aba830358e5361584f628450fd436c066323ad91ac2")
}

// UnpackSafeCastOverflowedIntToUintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func (gapCover *GapCover) UnpackSafeCastOverflowedIntToUintError(raw []byte) (*GapCoverSafeCastOverflowedIntToUint, error) {
	out := new(GapCoverSafeCastOverflowedIntToUint)
	if err := gapCover.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntToUint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the GapCover contract.
type GapCoverSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func GapCoverSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (gapCover *GapCover) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*GapCoverSafeCastOverflowedUintDowncast, error) {
	out := new(GapCoverSafeCastOverflowedUintDowncast)
	if err := gapCover.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the GapCover contract.
type GapCoverSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func GapCoverSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (gapCover *GapCover) UnpackSafeERC20FailedOperationError(raw []byte) (*GapCoverSafeERC20FailedOperation, error) {
	out := new(GapCoverSafeERC20FailedOperation)
	if err := gapCover.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSalesClosed represents a SalesClosed error raised by the GapCover contract.
type GapCoverSalesClosed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SalesClosed()
func GapCoverSalesClosedErrorID() common.Hash {
	return common.HexToHash("0x0671dd5e3fea89b895c1249668b5d4536fb7bc1c1981d7b9c170a9bf041f3909")
}

// UnpackSalesClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SalesClosed()
func (gapCover *GapCover) UnpackSalesClosedError(raw []byte) (*GapCoverSalesClosed, error) {
	out := new(GapCoverSalesClosed)
	if err := gapCover.abi.UnpackIntoInterface(out, "SalesClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSequencerNotSettled represents a SequencerNotSettled error raised by the GapCover contract.
type GapCoverSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func GapCoverSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (gapCover *GapCover) UnpackSequencerNotSettledError(raw []byte) (*GapCoverSequencerNotSettled, error) {
	out := new(GapCoverSequencerNotSettled)
	if err := gapCover.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverSeriesClosed represents a SeriesClosed error raised by the GapCover contract.
type GapCoverSeriesClosed struct {
	Symbol   [32]byte
	ClosesMs uint64
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SeriesClosed(bytes32 symbol, uint64 closesMs)
func GapCoverSeriesClosedErrorID() common.Hash {
	return common.HexToHash("0xb21f36f988c7ba81db086e8be5c39e79a90cba9b0f5fb241a80ccaa54e776861")
}

// UnpackSeriesClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SeriesClosed(bytes32 symbol, uint64 closesMs)
func (gapCover *GapCover) UnpackSeriesClosedError(raw []byte) (*GapCoverSeriesClosed, error) {
	out := new(GapCoverSeriesClosed)
	if err := gapCover.abi.UnpackIntoInterface(out, "SeriesClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverStaleReference represents a StaleReference error raised by the GapCover contract.
type GapCoverStaleReference struct {
	DeductibleBps *big.Int
	Minimum       *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error StaleReference(uint256 deductibleBps, uint256 minimum)
func GapCoverStaleReferenceErrorID() common.Hash {
	return common.HexToHash("0x6a0186ad4e2d797895d8a4ab43c51f248da11c886df2b49527fed3abdf0e8508")
}

// UnpackStaleReferenceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error StaleReference(uint256 deductibleBps, uint256 minimum)
func (gapCover *GapCover) UnpackStaleReferenceError(raw []byte) (*GapCoverStaleReference, error) {
	out := new(GapCoverStaleReference)
	if err := gapCover.abi.UnpackIntoInterface(out, "StaleReference", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverTooEarlyToMeasure represents a TooEarlyToMeasure error raised by the GapCover contract.
type GapCoverTooEarlyToMeasure struct {
	FromMs *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooEarlyToMeasure(uint256 fromMs)
func GapCoverTooEarlyToMeasureErrorID() common.Hash {
	return common.HexToHash("0x79b5a9574d89737e90d47e18c82c7db5dcfaf099ecafb8cde1f2521d146fb382")
}

// UnpackTooEarlyToMeasureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooEarlyToMeasure(uint256 fromMs)
func (gapCover *GapCover) UnpackTooEarlyToMeasureError(raw []byte) (*GapCoverTooEarlyToMeasure, error) {
	out := new(GapCoverTooEarlyToMeasure)
	if err := gapCover.abi.UnpackIntoInterface(out, "TooEarlyToMeasure", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverTooEarlyToVoid represents a TooEarlyToVoid error raised by the GapCover contract.
type GapCoverTooEarlyToVoid struct {
	FromMs *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooEarlyToVoid(uint256 fromMs)
func GapCoverTooEarlyToVoidErrorID() common.Hash {
	return common.HexToHash("0xa77c7e12ea61160b1e817f746fb7cefd4ad6cbde3c1dbd36205c45fc6686bc38")
}

// UnpackTooEarlyToVoidError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooEarlyToVoid(uint256 fromMs)
func (gapCover *GapCover) UnpackTooEarlyToVoidError(raw []byte) (*GapCoverTooEarlyToVoid, error) {
	out := new(GapCoverTooEarlyToVoid)
	if err := gapCover.abi.UnpackIntoInterface(out, "TooEarlyToVoid", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverWindowClosed represents a WindowClosed error raised by the GapCover contract.
type GapCoverWindowClosed struct {
	EndMs *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WindowClosed(uint256 endMs)
func GapCoverWindowClosedErrorID() common.Hash {
	return common.HexToHash("0xdb57bfb806503bdc5f96f060a25aead2076b299c2f2f3d4e36eab8e8336149e8")
}

// UnpackWindowClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WindowClosed(uint256 endMs)
func (gapCover *GapCover) UnpackWindowClosedError(raw []byte) (*GapCoverWindowClosed, error) {
	out := new(GapCoverWindowClosed)
	if err := gapCover.abi.UnpackIntoInterface(out, "WindowClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverWindowOpen represents a WindowOpen error raised by the GapCover contract.
type GapCoverWindowOpen struct {
	EndMs *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WindowOpen(uint256 endMs)
func GapCoverWindowOpenErrorID() common.Hash {
	return common.HexToHash("0xb317256bb83d3a0e7c565067a7c67553969518b9977a3d2d318c3e3fb33ef1e2")
}

// UnpackWindowOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WindowOpen(uint256 endMs)
func (gapCover *GapCover) UnpackWindowOpenError(raw []byte) (*GapCoverWindowOpen, error) {
	out := new(GapCoverWindowOpen)
	if err := gapCover.abi.UnpackIntoInterface(out, "WindowOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// GapCoverZeroAmount represents a ZeroAmount error raised by the GapCover contract.
type GapCoverZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAmount()
func GapCoverZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x1f2a2005cb66a8e145327e8814e243a0996aec9bfe1e15a495778b1236dbd485")
}

// UnpackZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAmount()
func (gapCover *GapCover) UnpackZeroAmountError(raw []byte) (*GapCoverZeroAmount, error) {
	out := new(GapCoverZeroAmount)
	if err := gapCover.abi.UnpackIntoInterface(out, "ZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
