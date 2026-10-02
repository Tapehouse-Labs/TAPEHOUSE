// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package marginaccounts

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

// MarginAccountsMetaData contains all meta data concerning the MarginAccounts contract.
var MarginAccountsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"band_\",\"type\":\"address\",\"internalType\":\"contractIBand\"},{\"name\":\"engine_\",\"type\":\"address\",\"internalType\":\"contractIMargin\"},{\"name\":\"vault_\",\"type\":\"address\",\"internalType\":\"contractSupplyVault\"},{\"name\":\"weth_\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"assetCaps\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"debtCap_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"weekendDebtCap_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"premiumRate_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"reserveShare_\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CROSS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"LIQUIDATION_PRICE_STEPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_FEED_AGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_PREMIUM_RATE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_RECALLS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_RECALL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RECALL_HAIRCUT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WETH_COLLATERAL_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"accruePremium\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"backstop\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"backstopPremium\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"band\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBand\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"borrow\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"borrowingPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimPremium\",\"inputs\":[],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"clear\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"closure\",\"inputs\":[],\"outputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"reopensMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"accruedMs\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"closureStart\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collateral\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"collectFee\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"debt\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"debtCap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"debtShares\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"engine\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMargin\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ethUsd\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAggregatorV3Interface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"health\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"missing\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"regime\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"holding\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"units\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"scale\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAuthorized\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"authorized\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lend\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractStockLendingVault\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lent\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"leverage\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidationPrice\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"borrowing\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"liquidator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"premium\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"premiumIndex\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"premiumRate\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recall\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recalls\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"repay\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"repayWithCollateral\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reserveShare\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"seize\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sellable\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAuthorization\",\"inputs\":[{\"name\":\"authorized\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBackstop\",\"inputs\":[{\"name\":\"newBackstop\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBorrowingPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"newLending\",\"type\":\"address\",\"internalType\":\"contractStockLendingVault\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLiquidator\",\"inputs\":[{\"name\":\"newLiquidator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPremiumRate\",\"inputs\":[{\"name\":\"rate\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stocks\",\"inputs\":[],\"outputs\":[{\"name\":\"symbols\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sync\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalDebtShares\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlend\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usdg\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSupplyVault\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weekendDebtCap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weekendLeverage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"weth\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawReserve\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"writeOff\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"written\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AssetWrittenDown\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"counted\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"balance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuthorizationSet\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"authorized\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"allowed\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BackstopSet\",\"inputs\":[{\"name\":\"backstop\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Borrow\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BorrowingPausedSet\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClosureSet\",\"inputs\":[{\"name\":\"closesMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"reopensMs\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeCollected\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianSet\",\"inputs\":[{\"name\":\"guardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HoldingCleared\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"units\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Lend\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LendingSet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"lending\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LiquidatorSet\",\"inputs\":[{\"name\":\"liquidator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumAccrued\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumClaimed\",\"inputs\":[{\"name\":\"backstop\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumPaid\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"premium\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toReserve\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PremiumRateSet\",\"inputs\":[{\"name\":\"rate\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Recall\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"ticket\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Repay\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReserveWithdrawn\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Seize\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unlend\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"symbol\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WriteOff\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"assets\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"premium\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AssetCapExceeded\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AssetFrozen\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AssetHalted\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AssetInOtherPosition\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AssetWrittenOff\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"BackstopAlreadySet\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BandMismatch\",\"inputs\":[{\"name\":\"engineBand\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Blocked\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BorrowingAgainstUsdg\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BorrowingIsPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CorporateActionPending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DebtCapExceeded\",\"inputs\":[{\"name\":\"debt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"HoldingNotEmpty\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientCollateral\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientMargin\",\"inputs\":[{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"requirement\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientReserve\",\"inputs\":[{\"name\":\"reserve\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBackstop\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDebtCaps\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLiquidator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPremium\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceiver\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingAlreadySet\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LiquidatorAlreadySet\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LiquidityUnknown\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NoLending\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"NotBackstop\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotLiquidator\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NothingToSettle\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OutOfReach\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reachable\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnershipCannotBeRenounced\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionNotEmpty\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RecallTooSmall\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedIntToUint\",\"inputs\":[{\"name\":\"value\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintToInt\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SequencerNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SessionUnknown\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManyAssets\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TooManyRecalls\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnknownAsset\",\"inputs\":[{\"name\":\"symbol\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnknownPosition\",\"inputs\":[{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnsupportedToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"WeekendDebtCapExceeded\",\"inputs\":[{\"name\":\"debt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WeekendLeverageExceeded\",\"inputs\":[{\"name\":\"gross\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"equity\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
	ID:  "MarginAccounts",
}

// MarginAccounts is an auto generated Go binding around an Ethereum contract.
type MarginAccounts struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *MarginAccounts) GetABI() abi.ABI {
	return c.abi
}

// NewMarginAccounts creates a new instance of MarginAccounts.
func NewMarginAccounts() *MarginAccounts {
	parsed, err := MarginAccountsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &MarginAccounts{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *MarginAccounts) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address band_, address engine_, address vault_, address weth_, address initialOwner, uint256[] assetCaps, uint256 debtCap_, uint256 weekendDebtCap_, uint32 premiumRate_, uint16 reserveShare_) returns()
func (marginAccounts *MarginAccounts) PackConstructor(band_ common.Address, engine_ common.Address, vault_ common.Address, weth_ common.Address, initialOwner common.Address, assetCaps []*big.Int, debtCap_ *big.Int, weekendDebtCap_ *big.Int, premiumRate_ uint32, reserveShare_ uint16) []byte {
	enc, err := marginAccounts.abi.Pack("", band_, engine_, vault_, weth_, initialOwner, assetCaps, debtCap_, weekendDebtCap_, premiumRate_, reserveShare_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCROSS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x39e03dc4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function CROSS() view returns(bytes32)
func (marginAccounts *MarginAccounts) PackCROSS() []byte {
	enc, err := marginAccounts.abi.Pack("CROSS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCROSS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x39e03dc4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function CROSS() view returns(bytes32)
func (marginAccounts *MarginAccounts) TryPackCROSS() ([]byte, error) {
	return marginAccounts.abi.Pack("CROSS")
}

// UnpackCROSS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x39e03dc4.
//
// Solidity: function CROSS() view returns(bytes32)
func (marginAccounts *MarginAccounts) UnpackCROSS(data []byte) ([32]byte, error) {
	out, err := marginAccounts.abi.Unpack("CROSS", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackLIQUIDATIONPRICESTEPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8d635001.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function LIQUIDATION_PRICE_STEPS() view returns(uint256)
func (marginAccounts *MarginAccounts) PackLIQUIDATIONPRICESTEPS() []byte {
	enc, err := marginAccounts.abi.Pack("LIQUIDATION_PRICE_STEPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLIQUIDATIONPRICESTEPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8d635001.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function LIQUIDATION_PRICE_STEPS() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackLIQUIDATIONPRICESTEPS() ([]byte, error) {
	return marginAccounts.abi.Pack("LIQUIDATION_PRICE_STEPS")
}

// UnpackLIQUIDATIONPRICESTEPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8d635001.
//
// Solidity: function LIQUIDATION_PRICE_STEPS() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackLIQUIDATIONPRICESTEPS(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("LIQUIDATION_PRICE_STEPS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXFEEDAGE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x696453b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (marginAccounts *MarginAccounts) PackMAXFEEDAGE() []byte {
	enc, err := marginAccounts.abi.Pack("MAX_FEED_AGE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXFEEDAGE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x696453b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackMAXFEEDAGE() ([]byte, error) {
	return marginAccounts.abi.Pack("MAX_FEED_AGE")
}

// UnpackMAXFEEDAGE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x696453b8.
//
// Solidity: function MAX_FEED_AGE() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackMAXFEEDAGE(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("MAX_FEED_AGE", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMAXPREMIUMRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e8b948a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_PREMIUM_RATE() view returns(uint32)
func (marginAccounts *MarginAccounts) PackMAXPREMIUMRATE() []byte {
	enc, err := marginAccounts.abi.Pack("MAX_PREMIUM_RATE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXPREMIUMRATE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e8b948a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_PREMIUM_RATE() view returns(uint32)
func (marginAccounts *MarginAccounts) TryPackMAXPREMIUMRATE() ([]byte, error) {
	return marginAccounts.abi.Pack("MAX_PREMIUM_RATE")
}

// UnpackMAXPREMIUMRATE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e8b948a.
//
// Solidity: function MAX_PREMIUM_RATE() view returns(uint32)
func (marginAccounts *MarginAccounts) UnpackMAXPREMIUMRATE(data []byte) (uint32, error) {
	out, err := marginAccounts.abi.Unpack("MAX_PREMIUM_RATE", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackMAXRECALLS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1224b3e5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MAX_RECALLS() view returns(uint256)
func (marginAccounts *MarginAccounts) PackMAXRECALLS() []byte {
	enc, err := marginAccounts.abi.Pack("MAX_RECALLS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMAXRECALLS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1224b3e5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MAX_RECALLS() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackMAXRECALLS() ([]byte, error) {
	return marginAccounts.abi.Pack("MAX_RECALLS")
}

// UnpackMAXRECALLS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1224b3e5.
//
// Solidity: function MAX_RECALLS() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackMAXRECALLS(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("MAX_RECALLS", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackMINRECALL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x473206c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function MIN_RECALL() view returns(uint256)
func (marginAccounts *MarginAccounts) PackMINRECALL() []byte {
	enc, err := marginAccounts.abi.Pack("MIN_RECALL")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMINRECALL is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x473206c6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function MIN_RECALL() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackMINRECALL() ([]byte, error) {
	return marginAccounts.abi.Pack("MIN_RECALL")
}

// UnpackMINRECALL is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x473206c6.
//
// Solidity: function MIN_RECALL() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackMINRECALL(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("MIN_RECALL", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRECALLHAIRCUT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b220af3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function RECALL_HAIRCUT() view returns(uint256)
func (marginAccounts *MarginAccounts) PackRECALLHAIRCUT() []byte {
	enc, err := marginAccounts.abi.Pack("RECALL_HAIRCUT")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRECALLHAIRCUT is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4b220af3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function RECALL_HAIRCUT() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackRECALLHAIRCUT() ([]byte, error) {
	return marginAccounts.abi.Pack("RECALL_HAIRCUT")
}

// UnpackRECALLHAIRCUT is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4b220af3.
//
// Solidity: function RECALL_HAIRCUT() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackRECALLHAIRCUT(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("RECALL_HAIRCUT", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWETHCOLLATERALBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa44ff063.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function WETH_COLLATERAL_BPS() view returns(uint256)
func (marginAccounts *MarginAccounts) PackWETHCOLLATERALBPS() []byte {
	enc, err := marginAccounts.abi.Pack("WETH_COLLATERAL_BPS")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWETHCOLLATERALBPS is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa44ff063.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function WETH_COLLATERAL_BPS() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackWETHCOLLATERALBPS() ([]byte, error) {
	return marginAccounts.abi.Pack("WETH_COLLATERAL_BPS")
}

// UnpackWETHCOLLATERALBPS is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa44ff063.
//
// Solidity: function WETH_COLLATERAL_BPS() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackWETHCOLLATERALBPS(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("WETH_COLLATERAL_BPS", data)
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
func (marginAccounts *MarginAccounts) PackAcceptOwnership() []byte {
	enc, err := marginAccounts.abi.Pack("acceptOwnership")
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
func (marginAccounts *MarginAccounts) TryPackAcceptOwnership() ([]byte, error) {
	return marginAccounts.abi.Pack("acceptOwnership")
}

// PackAccruePremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40f3356e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function accruePremium() returns(uint256)
func (marginAccounts *MarginAccounts) PackAccruePremium() []byte {
	enc, err := marginAccounts.abi.Pack("accruePremium")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAccruePremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40f3356e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function accruePremium() returns(uint256)
func (marginAccounts *MarginAccounts) TryPackAccruePremium() ([]byte, error) {
	return marginAccounts.abi.Pack("accruePremium")
}

// UnpackAccruePremium is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x40f3356e.
//
// Solidity: function accruePremium() returns(uint256)
func (marginAccounts *MarginAccounts) UnpackAccruePremium(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("accruePremium", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBackstop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7dea1817.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function backstop() view returns(address)
func (marginAccounts *MarginAccounts) PackBackstop() []byte {
	enc, err := marginAccounts.abi.Pack("backstop")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBackstop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7dea1817.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function backstop() view returns(address)
func (marginAccounts *MarginAccounts) TryPackBackstop() ([]byte, error) {
	return marginAccounts.abi.Pack("backstop")
}

// UnpackBackstop is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7dea1817.
//
// Solidity: function backstop() view returns(address)
func (marginAccounts *MarginAccounts) UnpackBackstop(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("backstop", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBackstopPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x878cda36.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function backstopPremium() view returns(uint128)
func (marginAccounts *MarginAccounts) PackBackstopPremium() []byte {
	enc, err := marginAccounts.abi.Pack("backstopPremium")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBackstopPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x878cda36.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function backstopPremium() view returns(uint128)
func (marginAccounts *MarginAccounts) TryPackBackstopPremium() ([]byte, error) {
	return marginAccounts.abi.Pack("backstopPremium")
}

// UnpackBackstopPremium is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x878cda36.
//
// Solidity: function backstopPremium() view returns(uint128)
func (marginAccounts *MarginAccounts) UnpackBackstopPremium(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("backstopPremium", data)
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
func (marginAccounts *MarginAccounts) PackBand() []byte {
	enc, err := marginAccounts.abi.Pack("band")
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
func (marginAccounts *MarginAccounts) TryPackBand() ([]byte, error) {
	return marginAccounts.abi.Pack("band")
}

// UnpackBand is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x10ea891a.
//
// Solidity: function band() view returns(address)
func (marginAccounts *MarginAccounts) UnpackBand(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("band", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5bfe19bc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrow(bytes32 position, uint256 assets, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) PackBorrow(position [32]byte, assets *big.Int, account common.Address, receiver common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("borrow", position, assets, account, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrow is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5bfe19bc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrow(bytes32 position, uint256 assets, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) TryPackBorrow(position [32]byte, assets *big.Int, account common.Address, receiver common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("borrow", position, assets, account, receiver)
}

// PackBorrowingPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xded7abc6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function borrowingPaused() view returns(bool)
func (marginAccounts *MarginAccounts) PackBorrowingPaused() []byte {
	enc, err := marginAccounts.abi.Pack("borrowingPaused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBorrowingPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xded7abc6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function borrowingPaused() view returns(bool)
func (marginAccounts *MarginAccounts) TryPackBorrowingPaused() ([]byte, error) {
	return marginAccounts.abi.Pack("borrowingPaused")
}

// UnpackBorrowingPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xded7abc6.
//
// Solidity: function borrowingPaused() view returns(bool)
func (marginAccounts *MarginAccounts) UnpackBorrowingPaused(data []byte) (bool, error) {
	out, err := marginAccounts.abi.Unpack("borrowingPaused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa5eb206e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) PackClaim(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("claim", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa5eb206e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackClaim(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("claim", account, position, token)
}

// UnpackClaim is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa5eb206e.
//
// Solidity: function claim(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackClaim(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("claim", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClaimPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x529663db.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimPremium() returns(uint256 amount)
func (marginAccounts *MarginAccounts) PackClaimPremium() []byte {
	enc, err := marginAccounts.abi.Pack("claimPremium")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x529663db.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimPremium() returns(uint256 amount)
func (marginAccounts *MarginAccounts) TryPackClaimPremium() ([]byte, error) {
	return marginAccounts.abi.Pack("claimPremium")
}

// UnpackClaimPremium is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x529663db.
//
// Solidity: function claimPremium() returns(uint256 amount)
func (marginAccounts *MarginAccounts) UnpackClaimPremium(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("claimPremium", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClear is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea6865a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function clear(address account, bytes32 position, bytes32 symbol) returns()
func (marginAccounts *MarginAccounts) PackClear(account common.Address, position [32]byte, symbol [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("clear", account, position, symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClear is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xea6865a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function clear(address account, bytes32 position, bytes32 symbol) returns()
func (marginAccounts *MarginAccounts) TryPackClear(account common.Address, position [32]byte, symbol [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("clear", account, position, symbol)
}

// PackClosure is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc528c1c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function closure() view returns(uint64 closesMs, uint64 reopensMs, uint64 accruedMs)
func (marginAccounts *MarginAccounts) PackClosure() []byte {
	enc, err := marginAccounts.abi.Pack("closure")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClosure is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc528c1c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function closure() view returns(uint64 closesMs, uint64 reopensMs, uint64 accruedMs)
func (marginAccounts *MarginAccounts) TryPackClosure() ([]byte, error) {
	return marginAccounts.abi.Pack("closure")
}

// ClosureOutput serves as a container for the return parameters of contract
// method Closure.
type ClosureOutput struct {
	ClosesMs  uint64
	ReopensMs uint64
	AccruedMs uint64
}

// UnpackClosure is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfc528c1c.
//
// Solidity: function closure() view returns(uint64 closesMs, uint64 reopensMs, uint64 accruedMs)
func (marginAccounts *MarginAccounts) UnpackClosure(data []byte) (ClosureOutput, error) {
	out, err := marginAccounts.abi.Unpack("closure", data)
	outstruct := new(ClosureOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.ClosesMs = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.ReopensMs = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.AccruedMs = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackClosureStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6d52049.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function closureStart() view returns(uint64)
func (marginAccounts *MarginAccounts) PackClosureStart() []byte {
	enc, err := marginAccounts.abi.Pack("closureStart")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClosureStart is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb6d52049.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function closureStart() view returns(uint64)
func (marginAccounts *MarginAccounts) TryPackClosureStart() ([]byte, error) {
	return marginAccounts.abi.Pack("closureStart")
}

// UnpackClosureStart is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb6d52049.
//
// Solidity: function closureStart() view returns(uint64)
func (marginAccounts *MarginAccounts) UnpackClosureStart(data []byte) (uint64, error) {
	out, err := marginAccounts.abi.Unpack("closureStart", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, nil
}

// PackCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb4b1e6d8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collateral(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) PackCollateral(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("collateral", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb4b1e6d8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collateral(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackCollateral(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("collateral", account, position, token)
}

// UnpackCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb4b1e6d8.
//
// Solidity: function collateral(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackCollateral(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("collateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCollectFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa969ff0a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function collectFee(uint256 amount) returns()
func (marginAccounts *MarginAccounts) PackCollectFee(amount *big.Int) []byte {
	enc, err := marginAccounts.abi.Pack("collectFee", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCollectFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa969ff0a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function collectFee(uint256 amount) returns()
func (marginAccounts *MarginAccounts) TryPackCollectFee(amount *big.Int) ([]byte, error) {
	return marginAccounts.abi.Pack("collectFee", amount)
}

// PackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1fbe11f9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function debt(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) PackDebt(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("debt", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDebt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1fbe11f9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function debt(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackDebt(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("debt", account, position)
}

// UnpackDebt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1fbe11f9.
//
// Solidity: function debt(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackDebt(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("debt", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDebtCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31486c06.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function debtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) PackDebtCap() []byte {
	enc, err := marginAccounts.abi.Pack("debtCap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDebtCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x31486c06.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function debtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackDebtCap() ([]byte, error) {
	return marginAccounts.abi.Pack("debtCap")
}

// UnpackDebtCap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x31486c06.
//
// Solidity: function debtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackDebtCap(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("debtCap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDebtShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f7156aa.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function debtShares(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) PackDebtShares(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("debtShares", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDebtShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f7156aa.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function debtShares(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackDebtShares(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("debtShares", account, position)
}

// UnpackDebtShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2f7156aa.
//
// Solidity: function debtShares(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackDebtShares(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("debtShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c8e09fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) PackDeposit(position [32]byte, token common.Address, amount *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("deposit", position, token, amount, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c8e09fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) TryPackDeposit(position [32]byte, token common.Address, amount *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("deposit", position, token, amount, account)
}

// PackDepositWithPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d7da0f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositWithPermit(bytes32 position, address token, uint256 amount, address account, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (marginAccounts *MarginAccounts) PackDepositWithPermit(position [32]byte, token common.Address, amount *big.Int, account common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("depositWithPermit", position, token, amount, account, deadline, v, r, s)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositWithPermit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d7da0f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositWithPermit(bytes32 position, address token, uint256 amount, address account, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (marginAccounts *MarginAccounts) TryPackDepositWithPermit(position [32]byte, token common.Address, amount *big.Int, account common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("depositWithPermit", position, token, amount, account, deadline, v, r, s)
}

// PackEngine is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc9d4623f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function engine() view returns(address)
func (marginAccounts *MarginAccounts) PackEngine() []byte {
	enc, err := marginAccounts.abi.Pack("engine")
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
func (marginAccounts *MarginAccounts) TryPackEngine() ([]byte, error) {
	return marginAccounts.abi.Pack("engine")
}

// UnpackEngine is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (marginAccounts *MarginAccounts) UnpackEngine(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("engine", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackEthUsd is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a960216.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ethUsd() view returns(address)
func (marginAccounts *MarginAccounts) PackEthUsd() []byte {
	enc, err := marginAccounts.abi.Pack("ethUsd")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEthUsd is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a960216.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ethUsd() view returns(address)
func (marginAccounts *MarginAccounts) TryPackEthUsd() ([]byte, error) {
	return marginAccounts.abi.Pack("ethUsd")
}

// UnpackEthUsd is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5a960216.
//
// Solidity: function ethUsd() view returns(address)
func (marginAccounts *MarginAccounts) UnpackEthUsd(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("ethUsd", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGuardian is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x452a9320.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function guardian() view returns(address)
func (marginAccounts *MarginAccounts) PackGuardian() []byte {
	enc, err := marginAccounts.abi.Pack("guardian")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGuardian is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x452a9320.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function guardian() view returns(address)
func (marginAccounts *MarginAccounts) TryPackGuardian() ([]byte, error) {
	return marginAccounts.abi.Pack("guardian")
}

// UnpackGuardian is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (marginAccounts *MarginAccounts) UnpackGuardian(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("guardian", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackHealth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42d6fb7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function health(address account, bytes32 position) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (marginAccounts *MarginAccounts) PackHealth(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("health", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHealth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42d6fb7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function health(address account, bytes32 position) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (marginAccounts *MarginAccounts) TryPackHealth(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("health", account, position)
}

// HealthOutput serves as a container for the return parameters of contract
// method Health.
type HealthOutput struct {
	Equity      *big.Int
	Requirement *big.Int
	Missing     uint8
	Regime      uint8
}

// UnpackHealth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42d6fb7b.
//
// Solidity: function health(address account, bytes32 position) view returns(int256 equity, uint256 requirement, uint8 missing, uint8 regime)
func (marginAccounts *MarginAccounts) UnpackHealth(data []byte) (HealthOutput, error) {
	out, err := marginAccounts.abi.Unpack("health", data)
	outstruct := new(HealthOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Equity = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Requirement = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Missing = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.Regime = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackHolding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa12f9b1a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function holding(bytes32 symbol) view returns(uint256 units, uint256 scale, uint256 cap)
func (marginAccounts *MarginAccounts) PackHolding(symbol [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("holding", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHolding is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa12f9b1a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function holding(bytes32 symbol) view returns(uint256 units, uint256 scale, uint256 cap)
func (marginAccounts *MarginAccounts) TryPackHolding(symbol [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("holding", symbol)
}

// HoldingOutput serves as a container for the return parameters of contract
// method Holding.
type HoldingOutput struct {
	Units *big.Int
	Scale *big.Int
	Cap   *big.Int
}

// UnpackHolding is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa12f9b1a.
//
// Solidity: function holding(bytes32 symbol) view returns(uint256 units, uint256 scale, uint256 cap)
func (marginAccounts *MarginAccounts) UnpackHolding(data []byte) (HoldingOutput, error) {
	out, err := marginAccounts.abi.Unpack("holding", data)
	outstruct := new(HoldingOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Units = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Scale = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	outstruct.Cap = abi.ConvertType(out[2], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackIsAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65e4ad9e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isAuthorized(address account, address authorized) view returns(bool)
func (marginAccounts *MarginAccounts) PackIsAuthorized(account common.Address, authorized common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("isAuthorized", account, authorized)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsAuthorized is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x65e4ad9e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isAuthorized(address account, address authorized) view returns(bool)
func (marginAccounts *MarginAccounts) TryPackIsAuthorized(account common.Address, authorized common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("isAuthorized", account, authorized)
}

// UnpackIsAuthorized is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x65e4ad9e.
//
// Solidity: function isAuthorized(address account, address authorized) view returns(bool)
func (marginAccounts *MarginAccounts) UnpackIsAuthorized(data []byte) (bool, error) {
	out, err := marginAccounts.abi.Unpack("isAuthorized", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackLend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc53b4c75.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lend(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) PackLend(position [32]byte, token common.Address, amount *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("lend", position, token, amount, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc53b4c75.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lend(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) TryPackLend(position [32]byte, token common.Address, amount *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("lend", position, token, amount, account)
}

// PackLending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfbc3e80.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lending(bytes32 symbol) view returns(address)
func (marginAccounts *MarginAccounts) PackLending(symbol [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("lending", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbfbc3e80.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lending(bytes32 symbol) view returns(address)
func (marginAccounts *MarginAccounts) TryPackLending(symbol [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("lending", symbol)
}

// UnpackLending is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbfbc3e80.
//
// Solidity: function lending(bytes32 symbol) view returns(address)
func (marginAccounts *MarginAccounts) UnpackLending(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("lending", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackLent is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdfcd917.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lent(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) PackLent(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("lent", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLent is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbdfcd917.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lent(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackLent(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("lent", account, position, token)
}

// UnpackLent is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbdfcd917.
//
// Solidity: function lent(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackLent(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("lent", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb300ac07.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function leverage(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) PackLeverage(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("leverage", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb300ac07.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function leverage(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackLeverage(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("leverage", account, position)
}

// UnpackLeverage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb300ac07.
//
// Solidity: function leverage(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackLeverage(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("leverage", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLiquidationPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x824fe15a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidationPrice(address account, bytes32 position, bytes32 symbol, uint256 borrowing) view returns(uint256)
func (marginAccounts *MarginAccounts) PackLiquidationPrice(account common.Address, position [32]byte, symbol [32]byte, borrowing *big.Int) []byte {
	enc, err := marginAccounts.abi.Pack("liquidationPrice", account, position, symbol, borrowing)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLiquidationPrice is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x824fe15a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function liquidationPrice(address account, bytes32 position, bytes32 symbol, uint256 borrowing) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackLiquidationPrice(account common.Address, position [32]byte, symbol [32]byte, borrowing *big.Int) ([]byte, error) {
	return marginAccounts.abi.Pack("liquidationPrice", account, position, symbol, borrowing)
}

// UnpackLiquidationPrice is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x824fe15a.
//
// Solidity: function liquidationPrice(address account, bytes32 position, bytes32 symbol, uint256 borrowing) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackLiquidationPrice(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("liquidationPrice", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackLiquidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4046ebae.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function liquidator() view returns(address)
func (marginAccounts *MarginAccounts) PackLiquidator() []byte {
	enc, err := marginAccounts.abi.Pack("liquidator")
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
func (marginAccounts *MarginAccounts) TryPackLiquidator() ([]byte, error) {
	return marginAccounts.abi.Pack("liquidator")
}

// UnpackLiquidator is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4046ebae.
//
// Solidity: function liquidator() view returns(address)
func (marginAccounts *MarginAccounts) UnpackLiquidator(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("liquidator", data)
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
func (marginAccounts *MarginAccounts) PackOwner() []byte {
	enc, err := marginAccounts.abi.Pack("owner")
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
func (marginAccounts *MarginAccounts) TryPackOwner() ([]byte, error) {
	return marginAccounts.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (marginAccounts *MarginAccounts) UnpackOwner(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("owner", data)
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
func (marginAccounts *MarginAccounts) PackPendingOwner() []byte {
	enc, err := marginAccounts.abi.Pack("pendingOwner")
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
func (marginAccounts *MarginAccounts) TryPackPendingOwner() ([]byte, error) {
	return marginAccounts.abi.Pack("pendingOwner")
}

// UnpackPendingOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (marginAccounts *MarginAccounts) UnpackPendingOwner(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("pendingOwner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x353f1fd1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function premium(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) PackPremium(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("premium", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPremium is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x353f1fd1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function premium(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackPremium(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("premium", account, position)
}

// UnpackPremium is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x353f1fd1.
//
// Solidity: function premium(address account, bytes32 position) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackPremium(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("premium", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPremiumIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xab98e17a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function premiumIndex() view returns(uint128)
func (marginAccounts *MarginAccounts) PackPremiumIndex() []byte {
	enc, err := marginAccounts.abi.Pack("premiumIndex")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPremiumIndex is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xab98e17a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function premiumIndex() view returns(uint128)
func (marginAccounts *MarginAccounts) TryPackPremiumIndex() ([]byte, error) {
	return marginAccounts.abi.Pack("premiumIndex")
}

// UnpackPremiumIndex is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xab98e17a.
//
// Solidity: function premiumIndex() view returns(uint128)
func (marginAccounts *MarginAccounts) UnpackPremiumIndex(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("premiumIndex", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPremiumRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d0040d3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function premiumRate() view returns(uint32)
func (marginAccounts *MarginAccounts) PackPremiumRate() []byte {
	enc, err := marginAccounts.abi.Pack("premiumRate")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPremiumRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5d0040d3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function premiumRate() view returns(uint32)
func (marginAccounts *MarginAccounts) TryPackPremiumRate() ([]byte, error) {
	return marginAccounts.abi.Pack("premiumRate")
}

// UnpackPremiumRate is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5d0040d3.
//
// Solidity: function premiumRate() view returns(uint32)
func (marginAccounts *MarginAccounts) UnpackPremiumRate(data []byte) (uint32, error) {
	out, err := marginAccounts.abi.Unpack("premiumRate", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x064c3d70.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recall(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) PackRecall(position [32]byte, token common.Address, amount *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("recall", position, token, amount, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecall is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x064c3d70.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recall(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) TryPackRecall(position [32]byte, token common.Address, amount *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("recall", position, token, amount, account)
}

// PackRecalls is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a1ebef7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function recalls(address account, bytes32 position, address token) view returns(uint256[])
func (marginAccounts *MarginAccounts) PackRecalls(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("recalls", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRecalls is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a1ebef7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function recalls(address account, bytes32 position, address token) view returns(uint256[])
func (marginAccounts *MarginAccounts) TryPackRecalls(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("recalls", account, position, token)
}

// UnpackRecalls is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a1ebef7.
//
// Solidity: function recalls(address account, bytes32 position, address token) view returns(uint256[])
func (marginAccounts *MarginAccounts) UnpackRecalls(data []byte) ([]*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("recalls", data)
	if err != nil {
		return *new([]*big.Int), err
	}
	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() pure returns()
func (marginAccounts *MarginAccounts) PackRenounceOwnership() []byte {
	enc, err := marginAccounts.abi.Pack("renounceOwnership")
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
func (marginAccounts *MarginAccounts) TryPackRenounceOwnership() ([]byte, error) {
	return marginAccounts.abi.Pack("renounceOwnership")
}

// PackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xff268125.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function repay(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) PackRepay(position [32]byte, assets *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("repay", position, assets, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRepay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xff268125.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function repay(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) TryPackRepay(position [32]byte, assets *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("repay", position, assets, account)
}

// UnpackRepay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xff268125.
//
// Solidity: function repay(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) UnpackRepay(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("repay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRepayWithCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe327cb85.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function repayWithCollateral(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) PackRepayWithCollateral(position [32]byte, assets *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("repayWithCollateral", position, assets, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRepayWithCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe327cb85.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function repayWithCollateral(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) TryPackRepayWithCollateral(position [32]byte, assets *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("repayWithCollateral", position, assets, account)
}

// UnpackRepayWithCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe327cb85.
//
// Solidity: function repayWithCollateral(bytes32 position, uint256 assets, address account) returns(uint256)
func (marginAccounts *MarginAccounts) UnpackRepayWithCollateral(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("repayWithCollateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackReserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcd3293de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reserve() view returns(uint128)
func (marginAccounts *MarginAccounts) PackReserve() []byte {
	enc, err := marginAccounts.abi.Pack("reserve")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcd3293de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reserve() view returns(uint128)
func (marginAccounts *MarginAccounts) TryPackReserve() ([]byte, error) {
	return marginAccounts.abi.Pack("reserve")
}

// UnpackReserve is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcd3293de.
//
// Solidity: function reserve() view returns(uint128)
func (marginAccounts *MarginAccounts) UnpackReserve(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("reserve", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackReserveShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe7cb3d67.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function reserveShare() view returns(uint16)
func (marginAccounts *MarginAccounts) PackReserveShare() []byte {
	enc, err := marginAccounts.abi.Pack("reserveShare")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReserveShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe7cb3d67.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function reserveShare() view returns(uint16)
func (marginAccounts *MarginAccounts) TryPackReserveShare() ([]byte, error) {
	return marginAccounts.abi.Pack("reserveShare")
}

// UnpackReserveShare is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe7cb3d67.
//
// Solidity: function reserveShare() view returns(uint16)
func (marginAccounts *MarginAccounts) UnpackReserveShare(data []byte) (uint16, error) {
	out, err := marginAccounts.abi.Unpack("reserveShare", data)
	if err != nil {
		return *new(uint16), err
	}
	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)
	return out0, nil
}

// PackSeize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b208ccb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function seize(bytes32 position, address token, uint256 amount, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) PackSeize(position [32]byte, token common.Address, amount *big.Int, account common.Address, receiver common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("seize", position, token, amount, account, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSeize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b208ccb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function seize(bytes32 position, address token, uint256 amount, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) TryPackSeize(position [32]byte, token common.Address, amount *big.Int, account common.Address, receiver common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("seize", position, token, amount, account, receiver)
}

// PackSellable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d26c683.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sellable(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) PackSellable(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("sellable", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSellable is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d26c683.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sellable(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackSellable(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("sellable", account, position, token)
}

// UnpackSellable is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0d26c683.
//
// Solidity: function sellable(address account, bytes32 position, address token) view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackSellable(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("sellable", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetAuthorization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeecea000.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAuthorization(address authorized, bool allowed) returns()
func (marginAccounts *MarginAccounts) PackSetAuthorization(authorized common.Address, allowed bool) []byte {
	enc, err := marginAccounts.abi.Pack("setAuthorization", authorized, allowed)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAuthorization is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xeecea000.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAuthorization(address authorized, bool allowed) returns()
func (marginAccounts *MarginAccounts) TryPackSetAuthorization(authorized common.Address, allowed bool) ([]byte, error) {
	return marginAccounts.abi.Pack("setAuthorization", authorized, allowed)
}

// PackSetBackstop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x916c2b87.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBackstop(address newBackstop) returns()
func (marginAccounts *MarginAccounts) PackSetBackstop(newBackstop common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("setBackstop", newBackstop)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBackstop is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x916c2b87.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBackstop(address newBackstop) returns()
func (marginAccounts *MarginAccounts) TryPackSetBackstop(newBackstop common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("setBackstop", newBackstop)
}

// PackSetBorrowingPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1820a5f4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setBorrowingPaused(bool paused) returns()
func (marginAccounts *MarginAccounts) PackSetBorrowingPaused(paused bool) []byte {
	enc, err := marginAccounts.abi.Pack("setBorrowingPaused", paused)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetBorrowingPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1820a5f4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setBorrowingPaused(bool paused) returns()
func (marginAccounts *MarginAccounts) TryPackSetBorrowingPaused(paused bool) ([]byte, error) {
	return marginAccounts.abi.Pack("setBorrowingPaused", paused)
}

// PackSetGuardian is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a0dac4a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setGuardian(address newGuardian) returns()
func (marginAccounts *MarginAccounts) PackSetGuardian(newGuardian common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("setGuardian", newGuardian)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetGuardian is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a0dac4a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setGuardian(address newGuardian) returns()
func (marginAccounts *MarginAccounts) TryPackSetGuardian(newGuardian common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("setGuardian", newGuardian)
}

// PackSetLending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3656f23.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setLending(bytes32 symbol, address newLending) returns()
func (marginAccounts *MarginAccounts) PackSetLending(symbol [32]byte, newLending common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("setLending", symbol, newLending)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetLending is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa3656f23.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setLending(bytes32 symbol, address newLending) returns()
func (marginAccounts *MarginAccounts) TryPackSetLending(symbol [32]byte, newLending common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("setLending", symbol, newLending)
}

// PackSetLiquidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01c76f81.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setLiquidator(address newLiquidator) returns()
func (marginAccounts *MarginAccounts) PackSetLiquidator(newLiquidator common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("setLiquidator", newLiquidator)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetLiquidator is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01c76f81.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setLiquidator(address newLiquidator) returns()
func (marginAccounts *MarginAccounts) TryPackSetLiquidator(newLiquidator common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("setLiquidator", newLiquidator)
}

// PackSetPremiumRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdea0b5f8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPremiumRate(uint32 rate) returns()
func (marginAccounts *MarginAccounts) PackSetPremiumRate(rate uint32) []byte {
	enc, err := marginAccounts.abi.Pack("setPremiumRate", rate)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPremiumRate is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdea0b5f8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPremiumRate(uint32 rate) returns()
func (marginAccounts *MarginAccounts) TryPackSetPremiumRate(rate uint32) ([]byte, error) {
	return marginAccounts.abi.Pack("setPremiumRate", rate)
}

// PackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55cf15e1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function settle(address account, bytes32 position, address token) returns()
func (marginAccounts *MarginAccounts) PackSettle(account common.Address, position [32]byte, token common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("settle", account, position, token)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55cf15e1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function settle(address account, bytes32 position, address token) returns()
func (marginAccounts *MarginAccounts) TryPackSettle(account common.Address, position [32]byte, token common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("settle", account, position, token)
}

// PackStocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdcc37192.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function stocks() view returns(bytes32[] symbols, address[] tokens)
func (marginAccounts *MarginAccounts) PackStocks() []byte {
	enc, err := marginAccounts.abi.Pack("stocks")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackStocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdcc37192.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function stocks() view returns(bytes32[] symbols, address[] tokens)
func (marginAccounts *MarginAccounts) TryPackStocks() ([]byte, error) {
	return marginAccounts.abi.Pack("stocks")
}

// StocksOutput serves as a container for the return parameters of contract
// method Stocks.
type StocksOutput struct {
	Symbols [][32]byte
	Tokens  []common.Address
}

// UnpackStocks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdcc37192.
//
// Solidity: function stocks() view returns(bytes32[] symbols, address[] tokens)
func (marginAccounts *MarginAccounts) UnpackStocks(data []byte) (StocksOutput, error) {
	out, err := marginAccounts.abi.Unpack("stocks", data)
	outstruct := new(StocksOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Symbols = *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	outstruct.Tokens = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	return *outstruct, nil
}

// PackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xae4b2551.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function sync(bytes32 symbol) returns()
func (marginAccounts *MarginAccounts) PackSync(symbol [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("sync", symbol)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSync is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xae4b2551.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function sync(bytes32 symbol) returns()
func (marginAccounts *MarginAccounts) TryPackSync(symbol [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("sync", symbol)
}

// PackTotalDebtShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b859e41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalDebtShares() view returns(uint128)
func (marginAccounts *MarginAccounts) PackTotalDebtShares() []byte {
	enc, err := marginAccounts.abi.Pack("totalDebtShares")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalDebtShares is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b859e41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalDebtShares() view returns(uint128)
func (marginAccounts *MarginAccounts) TryPackTotalDebtShares() ([]byte, error) {
	return marginAccounts.abi.Pack("totalDebtShares")
}

// UnpackTotalDebtShares is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b859e41.
//
// Solidity: function totalDebtShares() view returns(uint128)
func (marginAccounts *MarginAccounts) UnpackTotalDebtShares(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("totalDebtShares", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (marginAccounts *MarginAccounts) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("transferOwnership", newOwner)
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
func (marginAccounts *MarginAccounts) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("transferOwnership", newOwner)
}

// PackUnlend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfebd126f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unlend(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) PackUnlend(position [32]byte, token common.Address, amount *big.Int, account common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("unlend", position, token, amount, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnlend is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfebd126f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unlend(bytes32 position, address token, uint256 amount, address account) returns()
func (marginAccounts *MarginAccounts) TryPackUnlend(position [32]byte, token common.Address, amount *big.Int, account common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("unlend", position, token, amount, account)
}

// PackUsdg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf5b91b7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function usdg() view returns(address)
func (marginAccounts *MarginAccounts) PackUsdg() []byte {
	enc, err := marginAccounts.abi.Pack("usdg")
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
func (marginAccounts *MarginAccounts) TryPackUsdg() ([]byte, error) {
	return marginAccounts.abi.Pack("usdg")
}

// UnpackUsdg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf5b91b7b.
//
// Solidity: function usdg() view returns(address)
func (marginAccounts *MarginAccounts) UnpackUsdg(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("usdg", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbfa77cf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function vault() view returns(address)
func (marginAccounts *MarginAccounts) PackVault() []byte {
	enc, err := marginAccounts.abi.Pack("vault")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfbfa77cf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function vault() view returns(address)
func (marginAccounts *MarginAccounts) TryPackVault() ([]byte, error) {
	return marginAccounts.abi.Pack("vault")
}

// UnpackVault is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (marginAccounts *MarginAccounts) UnpackVault(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("vault", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWeekendDebtCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa504f558.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weekendDebtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) PackWeekendDebtCap() []byte {
	enc, err := marginAccounts.abi.Pack("weekendDebtCap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeekendDebtCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa504f558.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weekendDebtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) TryPackWeekendDebtCap() ([]byte, error) {
	return marginAccounts.abi.Pack("weekendDebtCap")
}

// UnpackWeekendDebtCap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa504f558.
//
// Solidity: function weekendDebtCap() view returns(uint256)
func (marginAccounts *MarginAccounts) UnpackWeekendDebtCap(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("weekendDebtCap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackWeekendLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb659ad89.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (marginAccounts *MarginAccounts) PackWeekendLeverage() []byte {
	enc, err := marginAccounts.abi.Pack("weekendLeverage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeekendLeverage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb659ad89.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (marginAccounts *MarginAccounts) TryPackWeekendLeverage() ([]byte, error) {
	return marginAccounts.abi.Pack("weekendLeverage")
}

// UnpackWeekendLeverage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb659ad89.
//
// Solidity: function weekendLeverage() view returns(uint32)
func (marginAccounts *MarginAccounts) UnpackWeekendLeverage(data []byte) (uint32, error) {
	out, err := marginAccounts.abi.Unpack("weekendLeverage", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackWeth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fc8cef3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function weth() view returns(address)
func (marginAccounts *MarginAccounts) PackWeth() []byte {
	enc, err := marginAccounts.abi.Pack("weth")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWeth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3fc8cef3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function weth() view returns(address)
func (marginAccounts *MarginAccounts) TryPackWeth() ([]byte, error) {
	return marginAccounts.abi.Pack("weth")
}

// UnpackWeth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (marginAccounts *MarginAccounts) UnpackWeth(data []byte) (common.Address, error) {
	out, err := marginAccounts.abi.Unpack("weth", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x672911d9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(bytes32 position, address token, uint256 amount, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) PackWithdraw(position [32]byte, token common.Address, amount *big.Int, account common.Address, receiver common.Address) []byte {
	enc, err := marginAccounts.abi.Pack("withdraw", position, token, amount, account, receiver)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x672911d9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(bytes32 position, address token, uint256 amount, address account, address receiver) returns()
func (marginAccounts *MarginAccounts) TryPackWithdraw(position [32]byte, token common.Address, amount *big.Int, account common.Address, receiver common.Address) ([]byte, error) {
	return marginAccounts.abi.Pack("withdraw", position, token, amount, account, receiver)
}

// PackWithdrawReserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c58ce14.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawReserve(address receiver, uint256 amount) returns()
func (marginAccounts *MarginAccounts) PackWithdrawReserve(receiver common.Address, amount *big.Int) []byte {
	enc, err := marginAccounts.abi.Pack("withdrawReserve", receiver, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawReserve is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c58ce14.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawReserve(address receiver, uint256 amount) returns()
func (marginAccounts *MarginAccounts) TryPackWithdrawReserve(receiver common.Address, amount *big.Int) ([]byte, error) {
	return marginAccounts.abi.Pack("withdrawReserve", receiver, amount)
}

// PackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f3596c0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256 written)
func (marginAccounts *MarginAccounts) PackWriteOff(account common.Address, position [32]byte) []byte {
	enc, err := marginAccounts.abi.Pack("writeOff", account, position)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWriteOff is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1f3596c0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256 written)
func (marginAccounts *MarginAccounts) TryPackWriteOff(account common.Address, position [32]byte) ([]byte, error) {
	return marginAccounts.abi.Pack("writeOff", account, position)
}

// UnpackWriteOff is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1f3596c0.
//
// Solidity: function writeOff(address account, bytes32 position) returns(uint256 written)
func (marginAccounts *MarginAccounts) UnpackWriteOff(data []byte) (*big.Int, error) {
	out, err := marginAccounts.abi.Unpack("writeOff", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// MarginAccountsAssetWrittenDown represents a AssetWrittenDown event raised by the MarginAccounts contract.
type MarginAccountsAssetWrittenDown struct {
	Symbol  [32]byte
	Counted *big.Int
	Balance *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const MarginAccountsAssetWrittenDownEventName = "AssetWrittenDown"

// ContractEventName returns the user-defined event name.
func (MarginAccountsAssetWrittenDown) ContractEventName() string {
	return MarginAccountsAssetWrittenDownEventName
}

// UnpackAssetWrittenDownEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AssetWrittenDown(bytes32 indexed symbol, uint256 counted, uint256 balance)
func (marginAccounts *MarginAccounts) UnpackAssetWrittenDownEvent(log *types.Log) (*MarginAccountsAssetWrittenDown, error) {
	event := "AssetWrittenDown"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsAssetWrittenDown)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsAuthorizationSet represents a AuthorizationSet event raised by the MarginAccounts contract.
type MarginAccountsAuthorizationSet struct {
	Account    common.Address
	Authorized common.Address
	Allowed    bool
	Raw        *types.Log // Blockchain specific contextual infos
}

const MarginAccountsAuthorizationSetEventName = "AuthorizationSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsAuthorizationSet) ContractEventName() string {
	return MarginAccountsAuthorizationSetEventName
}

// UnpackAuthorizationSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AuthorizationSet(address indexed account, address indexed authorized, bool allowed)
func (marginAccounts *MarginAccounts) UnpackAuthorizationSetEvent(log *types.Log) (*MarginAccountsAuthorizationSet, error) {
	event := "AuthorizationSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsAuthorizationSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsBackstopSet represents a BackstopSet event raised by the MarginAccounts contract.
type MarginAccountsBackstopSet struct {
	Backstop common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsBackstopSetEventName = "BackstopSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsBackstopSet) ContractEventName() string {
	return MarginAccountsBackstopSetEventName
}

// UnpackBackstopSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BackstopSet(address indexed backstop)
func (marginAccounts *MarginAccounts) UnpackBackstopSetEvent(log *types.Log) (*MarginAccountsBackstopSet, error) {
	event := "BackstopSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsBackstopSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsBorrow represents a Borrow event raised by the MarginAccounts contract.
type MarginAccountsBorrow struct {
	Caller   common.Address
	Account  common.Address
	Position [32]byte
	Assets   *big.Int
	Shares   *big.Int
	Receiver common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsBorrowEventName = "Borrow"

// ContractEventName returns the user-defined event name.
func (MarginAccountsBorrow) ContractEventName() string {
	return MarginAccountsBorrowEventName
}

// UnpackBorrowEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Borrow(address indexed caller, address indexed account, bytes32 indexed position, uint256 assets, uint256 shares, address receiver)
func (marginAccounts *MarginAccounts) UnpackBorrowEvent(log *types.Log) (*MarginAccountsBorrow, error) {
	event := "Borrow"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsBorrow)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsBorrowingPausedSet represents a BorrowingPausedSet event raised by the MarginAccounts contract.
type MarginAccountsBorrowingPausedSet struct {
	Paused bool
	Raw    *types.Log // Blockchain specific contextual infos
}

const MarginAccountsBorrowingPausedSetEventName = "BorrowingPausedSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsBorrowingPausedSet) ContractEventName() string {
	return MarginAccountsBorrowingPausedSetEventName
}

// UnpackBorrowingPausedSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BorrowingPausedSet(bool paused)
func (marginAccounts *MarginAccounts) UnpackBorrowingPausedSetEvent(log *types.Log) (*MarginAccountsBorrowingPausedSet, error) {
	event := "BorrowingPausedSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsBorrowingPausedSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsClosureSet represents a ClosureSet event raised by the MarginAccounts contract.
type MarginAccountsClosureSet struct {
	ClosesMs  uint64
	ReopensMs uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const MarginAccountsClosureSetEventName = "ClosureSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsClosureSet) ContractEventName() string {
	return MarginAccountsClosureSetEventName
}

// UnpackClosureSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ClosureSet(uint64 closesMs, uint64 reopensMs)
func (marginAccounts *MarginAccounts) UnpackClosureSetEvent(log *types.Log) (*MarginAccountsClosureSet, error) {
	event := "ClosureSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsClosureSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsDeposit represents a Deposit event raised by the MarginAccounts contract.
type MarginAccountsDeposit struct {
	Caller   common.Address
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsDepositEventName = "Deposit"

// ContractEventName returns the user-defined event name.
func (MarginAccountsDeposit) ContractEventName() string {
	return MarginAccountsDepositEventName
}

// UnpackDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Deposit(address indexed caller, address indexed account, bytes32 indexed position, address token, uint256 amount)
func (marginAccounts *MarginAccounts) UnpackDepositEvent(log *types.Log) (*MarginAccountsDeposit, error) {
	event := "Deposit"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsDeposit)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsFeeCollected represents a FeeCollected event raised by the MarginAccounts contract.
type MarginAccountsFeeCollected struct {
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const MarginAccountsFeeCollectedEventName = "FeeCollected"

// ContractEventName returns the user-defined event name.
func (MarginAccountsFeeCollected) ContractEventName() string {
	return MarginAccountsFeeCollectedEventName
}

// UnpackFeeCollectedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event FeeCollected(uint256 amount)
func (marginAccounts *MarginAccounts) UnpackFeeCollectedEvent(log *types.Log) (*MarginAccountsFeeCollected, error) {
	event := "FeeCollected"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsFeeCollected)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsGuardianSet represents a GuardianSet event raised by the MarginAccounts contract.
type MarginAccountsGuardianSet struct {
	Guardian common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsGuardianSetEventName = "GuardianSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsGuardianSet) ContractEventName() string {
	return MarginAccountsGuardianSetEventName
}

// UnpackGuardianSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GuardianSet(address indexed guardian)
func (marginAccounts *MarginAccounts) UnpackGuardianSetEvent(log *types.Log) (*MarginAccountsGuardianSet, error) {
	event := "GuardianSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsGuardianSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsHoldingCleared represents a HoldingCleared event raised by the MarginAccounts contract.
type MarginAccountsHoldingCleared struct {
	Account  common.Address
	Position [32]byte
	Symbol   [32]byte
	Units    *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsHoldingClearedEventName = "HoldingCleared"

// ContractEventName returns the user-defined event name.
func (MarginAccountsHoldingCleared) ContractEventName() string {
	return MarginAccountsHoldingClearedEventName
}

// UnpackHoldingClearedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event HoldingCleared(address indexed account, bytes32 indexed position, bytes32 indexed symbol, uint256 units)
func (marginAccounts *MarginAccounts) UnpackHoldingClearedEvent(log *types.Log) (*MarginAccountsHoldingCleared, error) {
	event := "HoldingCleared"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsHoldingCleared)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsLend represents a Lend event raised by the MarginAccounts contract.
type MarginAccountsLend struct {
	Account  common.Address
	Position [32]byte
	Symbol   [32]byte
	Amount   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsLendEventName = "Lend"

// ContractEventName returns the user-defined event name.
func (MarginAccountsLend) ContractEventName() string {
	return MarginAccountsLendEventName
}

// UnpackLendEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Lend(address indexed account, bytes32 indexed position, bytes32 indexed symbol, uint256 amount, uint256 shares)
func (marginAccounts *MarginAccounts) UnpackLendEvent(log *types.Log) (*MarginAccountsLend, error) {
	event := "Lend"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsLend)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsLendingSet represents a LendingSet event raised by the MarginAccounts contract.
type MarginAccountsLendingSet struct {
	Symbol  [32]byte
	Lending common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const MarginAccountsLendingSetEventName = "LendingSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsLendingSet) ContractEventName() string {
	return MarginAccountsLendingSetEventName
}

// UnpackLendingSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LendingSet(bytes32 indexed symbol, address indexed lending)
func (marginAccounts *MarginAccounts) UnpackLendingSetEvent(log *types.Log) (*MarginAccountsLendingSet, error) {
	event := "LendingSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsLendingSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsLiquidatorSet represents a LiquidatorSet event raised by the MarginAccounts contract.
type MarginAccountsLiquidatorSet struct {
	Liquidator common.Address
	Raw        *types.Log // Blockchain specific contextual infos
}

const MarginAccountsLiquidatorSetEventName = "LiquidatorSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsLiquidatorSet) ContractEventName() string {
	return MarginAccountsLiquidatorSetEventName
}

// UnpackLiquidatorSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event LiquidatorSet(address indexed liquidator)
func (marginAccounts *MarginAccounts) UnpackLiquidatorSetEvent(log *types.Log) (*MarginAccountsLiquidatorSet, error) {
	event := "LiquidatorSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsLiquidatorSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the MarginAccounts contract.
type MarginAccountsOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MarginAccountsOwnershipTransferStartedEventName = "OwnershipTransferStarted"

// ContractEventName returns the user-defined event name.
func (MarginAccountsOwnershipTransferStarted) ContractEventName() string {
	return MarginAccountsOwnershipTransferStartedEventName
}

// UnpackOwnershipTransferStartedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (marginAccounts *MarginAccounts) UnpackOwnershipTransferStartedEvent(log *types.Log) (*MarginAccountsOwnershipTransferStarted, error) {
	event := "OwnershipTransferStarted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsOwnershipTransferStarted)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsOwnershipTransferred represents a OwnershipTransferred event raised by the MarginAccounts contract.
type MarginAccountsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const MarginAccountsOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (MarginAccountsOwnershipTransferred) ContractEventName() string {
	return MarginAccountsOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (marginAccounts *MarginAccounts) UnpackOwnershipTransferredEvent(log *types.Log) (*MarginAccountsOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsPremiumAccrued represents a PremiumAccrued event raised by the MarginAccounts contract.
type MarginAccountsPremiumAccrued struct {
	Index *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const MarginAccountsPremiumAccruedEventName = "PremiumAccrued"

// ContractEventName returns the user-defined event name.
func (MarginAccountsPremiumAccrued) ContractEventName() string {
	return MarginAccountsPremiumAccruedEventName
}

// UnpackPremiumAccruedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PremiumAccrued(uint256 index)
func (marginAccounts *MarginAccounts) UnpackPremiumAccruedEvent(log *types.Log) (*MarginAccountsPremiumAccrued, error) {
	event := "PremiumAccrued"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsPremiumAccrued)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsPremiumClaimed represents a PremiumClaimed event raised by the MarginAccounts contract.
type MarginAccountsPremiumClaimed struct {
	Backstop common.Address
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsPremiumClaimedEventName = "PremiumClaimed"

// ContractEventName returns the user-defined event name.
func (MarginAccountsPremiumClaimed) ContractEventName() string {
	return MarginAccountsPremiumClaimedEventName
}

// UnpackPremiumClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PremiumClaimed(address indexed backstop, uint256 amount)
func (marginAccounts *MarginAccounts) UnpackPremiumClaimedEvent(log *types.Log) (*MarginAccountsPremiumClaimed, error) {
	event := "PremiumClaimed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsPremiumClaimed)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsPremiumPaid represents a PremiumPaid event raised by the MarginAccounts contract.
type MarginAccountsPremiumPaid struct {
	Account   common.Address
	Position  [32]byte
	Premium   *big.Int
	ToReserve *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const MarginAccountsPremiumPaidEventName = "PremiumPaid"

// ContractEventName returns the user-defined event name.
func (MarginAccountsPremiumPaid) ContractEventName() string {
	return MarginAccountsPremiumPaidEventName
}

// UnpackPremiumPaidEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PremiumPaid(address indexed account, bytes32 indexed position, uint256 premium, uint256 toReserve)
func (marginAccounts *MarginAccounts) UnpackPremiumPaidEvent(log *types.Log) (*MarginAccountsPremiumPaid, error) {
	event := "PremiumPaid"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsPremiumPaid)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsPremiumRateSet represents a PremiumRateSet event raised by the MarginAccounts contract.
type MarginAccountsPremiumRateSet struct {
	Rate uint32
	Raw  *types.Log // Blockchain specific contextual infos
}

const MarginAccountsPremiumRateSetEventName = "PremiumRateSet"

// ContractEventName returns the user-defined event name.
func (MarginAccountsPremiumRateSet) ContractEventName() string {
	return MarginAccountsPremiumRateSetEventName
}

// UnpackPremiumRateSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PremiumRateSet(uint32 rate)
func (marginAccounts *MarginAccounts) UnpackPremiumRateSetEvent(log *types.Log) (*MarginAccountsPremiumRateSet, error) {
	event := "PremiumRateSet"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsPremiumRateSet)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsRecall represents a Recall event raised by the MarginAccounts contract.
type MarginAccountsRecall struct {
	Account  common.Address
	Position [32]byte
	Symbol   [32]byte
	Ticket   *big.Int
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsRecallEventName = "Recall"

// ContractEventName returns the user-defined event name.
func (MarginAccountsRecall) ContractEventName() string {
	return MarginAccountsRecallEventName
}

// UnpackRecallEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Recall(address indexed account, bytes32 indexed position, bytes32 indexed symbol, uint256 ticket, uint256 amount)
func (marginAccounts *MarginAccounts) UnpackRecallEvent(log *types.Log) (*MarginAccountsRecall, error) {
	event := "Recall"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsRecall)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsRepay represents a Repay event raised by the MarginAccounts contract.
type MarginAccountsRepay struct {
	Caller   common.Address
	Account  common.Address
	Position [32]byte
	Assets   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsRepayEventName = "Repay"

// ContractEventName returns the user-defined event name.
func (MarginAccountsRepay) ContractEventName() string {
	return MarginAccountsRepayEventName
}

// UnpackRepayEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Repay(address indexed caller, address indexed account, bytes32 indexed position, uint256 assets, uint256 shares)
func (marginAccounts *MarginAccounts) UnpackRepayEvent(log *types.Log) (*MarginAccountsRepay, error) {
	event := "Repay"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsRepay)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsReserveWithdrawn represents a ReserveWithdrawn event raised by the MarginAccounts contract.
type MarginAccountsReserveWithdrawn struct {
	Receiver common.Address
	Amount   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsReserveWithdrawnEventName = "ReserveWithdrawn"

// ContractEventName returns the user-defined event name.
func (MarginAccountsReserveWithdrawn) ContractEventName() string {
	return MarginAccountsReserveWithdrawnEventName
}

// UnpackReserveWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ReserveWithdrawn(address indexed receiver, uint256 amount)
func (marginAccounts *MarginAccounts) UnpackReserveWithdrawnEvent(log *types.Log) (*MarginAccountsReserveWithdrawn, error) {
	event := "ReserveWithdrawn"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsReserveWithdrawn)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsSeize represents a Seize event raised by the MarginAccounts contract.
type MarginAccountsSeize struct {
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Receiver common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsSeizeEventName = "Seize"

// ContractEventName returns the user-defined event name.
func (MarginAccountsSeize) ContractEventName() string {
	return MarginAccountsSeizeEventName
}

// UnpackSeizeEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Seize(address indexed account, bytes32 indexed position, address token, uint256 amount, address receiver)
func (marginAccounts *MarginAccounts) UnpackSeizeEvent(log *types.Log) (*MarginAccountsSeize, error) {
	event := "Seize"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsSeize)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsUnlend represents a Unlend event raised by the MarginAccounts contract.
type MarginAccountsUnlend struct {
	Account  common.Address
	Position [32]byte
	Symbol   [32]byte
	Amount   *big.Int
	Shares   *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsUnlendEventName = "Unlend"

// ContractEventName returns the user-defined event name.
func (MarginAccountsUnlend) ContractEventName() string {
	return MarginAccountsUnlendEventName
}

// UnpackUnlendEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Unlend(address indexed account, bytes32 indexed position, bytes32 indexed symbol, uint256 amount, uint256 shares)
func (marginAccounts *MarginAccounts) UnpackUnlendEvent(log *types.Log) (*MarginAccountsUnlend, error) {
	event := "Unlend"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsUnlend)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsWithdraw represents a Withdraw event raised by the MarginAccounts contract.
type MarginAccountsWithdraw struct {
	Caller   common.Address
	Account  common.Address
	Position [32]byte
	Token    common.Address
	Amount   *big.Int
	Receiver common.Address
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsWithdrawEventName = "Withdraw"

// ContractEventName returns the user-defined event name.
func (MarginAccountsWithdraw) ContractEventName() string {
	return MarginAccountsWithdrawEventName
}

// UnpackWithdrawEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdraw(address indexed caller, address indexed account, bytes32 indexed position, address token, uint256 amount, address receiver)
func (marginAccounts *MarginAccounts) UnpackWithdrawEvent(log *types.Log) (*MarginAccountsWithdraw, error) {
	event := "Withdraw"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsWithdraw)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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

// MarginAccountsWriteOff represents a WriteOff event raised by the MarginAccounts contract.
type MarginAccountsWriteOff struct {
	Account  common.Address
	Position [32]byte
	Assets   *big.Int
	Shares   *big.Int
	Premium  *big.Int
	Raw      *types.Log // Blockchain specific contextual infos
}

const MarginAccountsWriteOffEventName = "WriteOff"

// ContractEventName returns the user-defined event name.
func (MarginAccountsWriteOff) ContractEventName() string {
	return MarginAccountsWriteOffEventName
}

// UnpackWriteOffEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WriteOff(address indexed account, bytes32 indexed position, uint256 assets, uint256 shares, uint256 premium)
func (marginAccounts *MarginAccounts) UnpackWriteOffEvent(log *types.Log) (*MarginAccountsWriteOff, error) {
	event := "WriteOff"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != marginAccounts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(MarginAccountsWriteOff)
	if len(log.Data) > 0 {
		if err := marginAccounts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range marginAccounts.abi.Events[event].Inputs {
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
func (marginAccounts *MarginAccounts) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AddressFrozen"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAddressFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AssetCapExceeded"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAssetCapExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AssetFrozen"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAssetFrozenError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AssetHalted"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAssetHaltedError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AssetInOtherPosition"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAssetInOtherPositionError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["AssetWrittenOff"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackAssetWrittenOffError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["BackstopAlreadySet"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackBackstopAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["BandMismatch"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackBandMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["Blocked"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackBlockedError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["BorrowingAgainstUsdg"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackBorrowingAgainstUsdgError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["BorrowingIsPaused"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackBorrowingIsPausedError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["CorporateActionPending"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackCorporateActionPendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["DebtCapExceeded"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackDebtCapExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["HoldingNotEmpty"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackHoldingNotEmptyError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InsufficientCollateral"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInsufficientCollateralError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InsufficientMargin"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInsufficientMarginError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InsufficientReserve"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInsufficientReserveError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidBackstop"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidBackstopError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidDebtCaps"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidDebtCapsError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidLending"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidLendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidLiquidator"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidLiquidatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidPremium"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidPremiumError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["InvalidReceiver"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackInvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["LendingAlreadySet"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackLendingAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["LengthMismatch"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["LiquidatorAlreadySet"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackLiquidatorAlreadySetError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["LiquidityUnknown"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackLiquidityUnknownError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["NoLending"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackNoLendingError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["NotBackstop"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackNotBackstopError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["NotGuardian"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackNotGuardianError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["NotLiquidator"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackNotLiquidatorError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["NothingToSettle"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackNothingToSettleError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["OutOfReach"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackOutOfReachError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["OwnershipCannotBeRenounced"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackOwnershipCannotBeRenouncedError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["PositionNotEmpty"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackPositionNotEmptyError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["RecallTooSmall"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackRecallTooSmallError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SafeCastOverflowedIntToUint"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSafeCastOverflowedIntToUintError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SafeCastOverflowedUintToInt"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSafeCastOverflowedUintToIntError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SequencerNotSettled"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSequencerNotSettledError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["SessionUnknown"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackSessionUnknownError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["TooManyAssets"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackTooManyAssetsError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["TooManyRecalls"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackTooManyRecallsError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["Unauthorized"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackUnauthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["UnknownAsset"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackUnknownAssetError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["UnknownPosition"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackUnknownPositionError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["UnsupportedToken"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackUnsupportedTokenError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["WeekendDebtCapExceeded"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackWeekendDebtCapExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["WeekendLeverageExceeded"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackWeekendLeverageExceededError(raw[4:])
	}
	if bytes.Equal(raw[:4], marginAccounts.abi.Errors["ZeroAmount"].ID.Bytes()[:4]) {
		return marginAccounts.UnpackZeroAmountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// MarginAccountsAddressFrozen represents a AddressFrozen error raised by the MarginAccounts contract.
type MarginAccountsAddressFrozen struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressFrozen()
func MarginAccountsAddressFrozenErrorID() common.Hash {
	return common.HexToHash("0x1fd1cc4450f02aaade081fde79025f8828a7fb3ef46f6a3230a9c3914a5eb5ba")
}

// UnpackAddressFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressFrozen()
func (marginAccounts *MarginAccounts) UnpackAddressFrozenError(raw []byte) (*MarginAccountsAddressFrozen, error) {
	out := new(MarginAccountsAddressFrozen)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AddressFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsAssetCapExceeded represents a AssetCapExceeded error raised by the MarginAccounts contract.
type MarginAccountsAssetCapExceeded struct {
	Symbol [32]byte
	Cap    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetCapExceeded(bytes32 symbol, uint256 cap)
func MarginAccountsAssetCapExceededErrorID() common.Hash {
	return common.HexToHash("0xd26e4926582d061de714f1f5b1d76bff2734d0abd3c23ad2ea36029e4fd611ae")
}

// UnpackAssetCapExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetCapExceeded(bytes32 symbol, uint256 cap)
func (marginAccounts *MarginAccounts) UnpackAssetCapExceededError(raw []byte) (*MarginAccountsAssetCapExceeded, error) {
	out := new(MarginAccountsAssetCapExceeded)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AssetCapExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsAssetFrozen represents a AssetFrozen error raised by the MarginAccounts contract.
type MarginAccountsAssetFrozen struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetFrozen(bytes32 symbol)
func MarginAccountsAssetFrozenErrorID() common.Hash {
	return common.HexToHash("0xa9701556c60014385ace67d751b32db546a5f29063875844c48138ce118c1fff")
}

// UnpackAssetFrozenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetFrozen(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackAssetFrozenError(raw []byte) (*MarginAccountsAssetFrozen, error) {
	out := new(MarginAccountsAssetFrozen)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AssetFrozen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsAssetHalted represents a AssetHalted error raised by the MarginAccounts contract.
type MarginAccountsAssetHalted struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetHalted(bytes32 symbol)
func MarginAccountsAssetHaltedErrorID() common.Hash {
	return common.HexToHash("0x3ec29ca3785f47585a9c96baed33ddbd8b2d9e8dfdeffb4be96925aaf9ec6468")
}

// UnpackAssetHaltedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetHalted(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackAssetHaltedError(raw []byte) (*MarginAccountsAssetHalted, error) {
	out := new(MarginAccountsAssetHalted)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AssetHalted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsAssetInOtherPosition represents a AssetInOtherPosition error raised by the MarginAccounts contract.
type MarginAccountsAssetInOtherPosition struct {
	Symbol   [32]byte
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetInOtherPosition(bytes32 symbol, bytes32 position)
func MarginAccountsAssetInOtherPositionErrorID() common.Hash {
	return common.HexToHash("0x9603ebc4a6f1ea14122ca5c3662354410c77913f3c292bd318336a4a13b93795")
}

// UnpackAssetInOtherPositionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetInOtherPosition(bytes32 symbol, bytes32 position)
func (marginAccounts *MarginAccounts) UnpackAssetInOtherPositionError(raw []byte) (*MarginAccountsAssetInOtherPosition, error) {
	out := new(MarginAccountsAssetInOtherPosition)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AssetInOtherPosition", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsAssetWrittenOff represents a AssetWrittenOff error raised by the MarginAccounts contract.
type MarginAccountsAssetWrittenOff struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetWrittenOff(bytes32 symbol)
func MarginAccountsAssetWrittenOffErrorID() common.Hash {
	return common.HexToHash("0xc1e303a118fcd97571f7a4f2dae7541f42c680862e61d4521866da8994932629")
}

// UnpackAssetWrittenOffError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetWrittenOff(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackAssetWrittenOffError(raw []byte) (*MarginAccountsAssetWrittenOff, error) {
	out := new(MarginAccountsAssetWrittenOff)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "AssetWrittenOff", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsBackstopAlreadySet represents a BackstopAlreadySet error raised by the MarginAccounts contract.
type MarginAccountsBackstopAlreadySet struct {
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BackstopAlreadySet(address current)
func MarginAccountsBackstopAlreadySetErrorID() common.Hash {
	return common.HexToHash("0x5227892ec9afa323b580f66e4ae39b70a5c0a8c6d13b287d040e3178f45e58a1")
}

// UnpackBackstopAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BackstopAlreadySet(address current)
func (marginAccounts *MarginAccounts) UnpackBackstopAlreadySetError(raw []byte) (*MarginAccountsBackstopAlreadySet, error) {
	out := new(MarginAccountsBackstopAlreadySet)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "BackstopAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsBandMismatch represents a BandMismatch error raised by the MarginAccounts contract.
type MarginAccountsBandMismatch struct {
	EngineBand common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BandMismatch(address engineBand)
func MarginAccountsBandMismatchErrorID() common.Hash {
	return common.HexToHash("0x90e8bf8c67f0d0ce892d5748cf7cbe6da2f310e91938b0112662b9324c7abc3f")
}

// UnpackBandMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BandMismatch(address engineBand)
func (marginAccounts *MarginAccounts) UnpackBandMismatchError(raw []byte) (*MarginAccountsBandMismatch, error) {
	out := new(MarginAccountsBandMismatch)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "BandMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsBlocked represents a Blocked error raised by the MarginAccounts contract.
type MarginAccountsBlocked struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Blocked(address account)
func MarginAccountsBlockedErrorID() common.Hash {
	return common.HexToHash("0x75e91ce73c1d3352d8dd3610443539cd33dfe13b1de8f8caae54ec26dd0dc9cb")
}

// UnpackBlockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Blocked(address account)
func (marginAccounts *MarginAccounts) UnpackBlockedError(raw []byte) (*MarginAccountsBlocked, error) {
	out := new(MarginAccountsBlocked)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "Blocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsBorrowingAgainstUsdg represents a BorrowingAgainstUsdg error raised by the MarginAccounts contract.
type MarginAccountsBorrowingAgainstUsdg struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BorrowingAgainstUsdg()
func MarginAccountsBorrowingAgainstUsdgErrorID() common.Hash {
	return common.HexToHash("0x2f71a51ed8369bc9ae9df4f85c9cbe2320f3bd961e772e27256c35d27732ee30")
}

// UnpackBorrowingAgainstUsdgError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BorrowingAgainstUsdg()
func (marginAccounts *MarginAccounts) UnpackBorrowingAgainstUsdgError(raw []byte) (*MarginAccountsBorrowingAgainstUsdg, error) {
	out := new(MarginAccountsBorrowingAgainstUsdg)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "BorrowingAgainstUsdg", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsBorrowingIsPaused represents a BorrowingIsPaused error raised by the MarginAccounts contract.
type MarginAccountsBorrowingIsPaused struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BorrowingIsPaused()
func MarginAccountsBorrowingIsPausedErrorID() common.Hash {
	return common.HexToHash("0xaa864f37b700d896a72eafa272084f861aa9962378d6a62f8229dc64289fbc6a")
}

// UnpackBorrowingIsPausedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BorrowingIsPaused()
func (marginAccounts *MarginAccounts) UnpackBorrowingIsPausedError(raw []byte) (*MarginAccountsBorrowingIsPaused, error) {
	out := new(MarginAccountsBorrowingIsPaused)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "BorrowingIsPaused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsCorporateActionPending represents a CorporateActionPending error raised by the MarginAccounts contract.
type MarginAccountsCorporateActionPending struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func MarginAccountsCorporateActionPendingErrorID() common.Hash {
	return common.HexToHash("0x1d1f5fd3c689a2e2561007ef605d914d864ac52012393f43a887a76f9f48be7b")
}

// UnpackCorporateActionPendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CorporateActionPending(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackCorporateActionPendingError(raw []byte) (*MarginAccountsCorporateActionPending, error) {
	out := new(MarginAccountsCorporateActionPending)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "CorporateActionPending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsDebtCapExceeded represents a DebtCapExceeded error raised by the MarginAccounts contract.
type MarginAccountsDebtCapExceeded struct {
	Debt *big.Int
	Cap  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DebtCapExceeded(uint256 debt, uint256 cap)
func MarginAccountsDebtCapExceededErrorID() common.Hash {
	return common.HexToHash("0x42640c463bb188e8e4fd788d76edfd30e1189f8ae2b0a7a0a295e459b29970e8")
}

// UnpackDebtCapExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DebtCapExceeded(uint256 debt, uint256 cap)
func (marginAccounts *MarginAccounts) UnpackDebtCapExceededError(raw []byte) (*MarginAccountsDebtCapExceeded, error) {
	out := new(MarginAccountsDebtCapExceeded)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "DebtCapExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsHoldingNotEmpty represents a HoldingNotEmpty error raised by the MarginAccounts contract.
type MarginAccountsHoldingNotEmpty struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HoldingNotEmpty(bytes32 symbol)
func MarginAccountsHoldingNotEmptyErrorID() common.Hash {
	return common.HexToHash("0x1a1af020b6ed60664edf4b3e4968cddb32cc18ad449116b1c600006efa7bcdcd")
}

// UnpackHoldingNotEmptyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HoldingNotEmpty(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackHoldingNotEmptyError(raw []byte) (*MarginAccountsHoldingNotEmpty, error) {
	out := new(MarginAccountsHoldingNotEmpty)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "HoldingNotEmpty", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInsufficientCollateral represents a InsufficientCollateral error raised by the MarginAccounts contract.
type MarginAccountsInsufficientCollateral struct {
	Token  common.Address
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientCollateral(address token, uint256 amount)
func MarginAccountsInsufficientCollateralErrorID() common.Hash {
	return common.HexToHash("0x8564b7e511caa44438f2e47da8d489246aa57b29a865b4eafeb174867af3dfbe")
}

// UnpackInsufficientCollateralError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientCollateral(address token, uint256 amount)
func (marginAccounts *MarginAccounts) UnpackInsufficientCollateralError(raw []byte) (*MarginAccountsInsufficientCollateral, error) {
	out := new(MarginAccountsInsufficientCollateral)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InsufficientCollateral", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInsufficientMargin represents a InsufficientMargin error raised by the MarginAccounts contract.
type MarginAccountsInsufficientMargin struct {
	Equity      *big.Int
	Requirement *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientMargin(int256 equity, uint256 requirement)
func MarginAccountsInsufficientMarginErrorID() common.Hash {
	return common.HexToHash("0x1e4cbb75332dfc383c902bc80a5a1b1f64daea3dcc8e970b2d2b447b2b7640ea")
}

// UnpackInsufficientMarginError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientMargin(int256 equity, uint256 requirement)
func (marginAccounts *MarginAccounts) UnpackInsufficientMarginError(raw []byte) (*MarginAccountsInsufficientMargin, error) {
	out := new(MarginAccountsInsufficientMargin)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InsufficientMargin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInsufficientReserve represents a InsufficientReserve error raised by the MarginAccounts contract.
type MarginAccountsInsufficientReserve struct {
	Reserve *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientReserve(uint256 reserve)
func MarginAccountsInsufficientReserveErrorID() common.Hash {
	return common.HexToHash("0xd50317fadb9c4c16d543966c0cc754df3ed1a59a85a3be08c3fa56510189e078")
}

// UnpackInsufficientReserveError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientReserve(uint256 reserve)
func (marginAccounts *MarginAccounts) UnpackInsufficientReserveError(raw []byte) (*MarginAccountsInsufficientReserve, error) {
	out := new(MarginAccountsInsufficientReserve)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InsufficientReserve", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidBackstop represents a InvalidBackstop error raised by the MarginAccounts contract.
type MarginAccountsInvalidBackstop struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidBackstop()
func MarginAccountsInvalidBackstopErrorID() common.Hash {
	return common.HexToHash("0x15676e523b6197d055d8b5d5ab74e12e4a0bdb596c7c543ddd248b8403ab10e3")
}

// UnpackInvalidBackstopError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidBackstop()
func (marginAccounts *MarginAccounts) UnpackInvalidBackstopError(raw []byte) (*MarginAccountsInvalidBackstop, error) {
	out := new(MarginAccountsInvalidBackstop)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidBackstop", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidDebtCaps represents a InvalidDebtCaps error raised by the MarginAccounts contract.
type MarginAccountsInvalidDebtCaps struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidDebtCaps()
func MarginAccountsInvalidDebtCapsErrorID() common.Hash {
	return common.HexToHash("0xbd569264286df35726e56618241dcd038c513ed6edb018f9e6ac69f8ea03796d")
}

// UnpackInvalidDebtCapsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidDebtCaps()
func (marginAccounts *MarginAccounts) UnpackInvalidDebtCapsError(raw []byte) (*MarginAccountsInvalidDebtCaps, error) {
	out := new(MarginAccountsInvalidDebtCaps)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidDebtCaps", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidLending represents a InvalidLending error raised by the MarginAccounts contract.
type MarginAccountsInvalidLending struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidLending()
func MarginAccountsInvalidLendingErrorID() common.Hash {
	return common.HexToHash("0x006b3c69b08f0e6c455f21d96f8b54c3b5339573cfadf5deba962952af8e4fc1")
}

// UnpackInvalidLendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidLending()
func (marginAccounts *MarginAccounts) UnpackInvalidLendingError(raw []byte) (*MarginAccountsInvalidLending, error) {
	out := new(MarginAccountsInvalidLending)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidLending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidLiquidator represents a InvalidLiquidator error raised by the MarginAccounts contract.
type MarginAccountsInvalidLiquidator struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidLiquidator()
func MarginAccountsInvalidLiquidatorErrorID() common.Hash {
	return common.HexToHash("0x87e9041367999129b4aa3de21d092d2b3a17ceb9752c6ee137c5420ec65cdfaf")
}

// UnpackInvalidLiquidatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidLiquidator()
func (marginAccounts *MarginAccounts) UnpackInvalidLiquidatorError(raw []byte) (*MarginAccountsInvalidLiquidator, error) {
	out := new(MarginAccountsInvalidLiquidator)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidLiquidator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidPremium represents a InvalidPremium error raised by the MarginAccounts contract.
type MarginAccountsInvalidPremium struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPremium()
func MarginAccountsInvalidPremiumErrorID() common.Hash {
	return common.HexToHash("0x842a7f1a7934b68b3335a3c37c4edbdf0c9491124e29111f80d785815d701034")
}

// UnpackInvalidPremiumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPremium()
func (marginAccounts *MarginAccounts) UnpackInvalidPremiumError(raw []byte) (*MarginAccountsInvalidPremium, error) {
	out := new(MarginAccountsInvalidPremium)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidPremium", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsInvalidReceiver represents a InvalidReceiver error raised by the MarginAccounts contract.
type MarginAccountsInvalidReceiver struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidReceiver()
func MarginAccountsInvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0x1e4ec46ba639431522ed9b9ce6d5e4d79b6e4cb2c689f14af956db4edf067a7d")
}

// UnpackInvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidReceiver()
func (marginAccounts *MarginAccounts) UnpackInvalidReceiverError(raw []byte) (*MarginAccountsInvalidReceiver, error) {
	out := new(MarginAccountsInvalidReceiver)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsLendingAlreadySet represents a LendingAlreadySet error raised by the MarginAccounts contract.
type MarginAccountsLendingAlreadySet struct {
	Symbol  [32]byte
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingAlreadySet(bytes32 symbol, address current)
func MarginAccountsLendingAlreadySetErrorID() common.Hash {
	return common.HexToHash("0x09dc263d61dbe11703996d69c7a02e8058ca5fa3bd2365d9257760b1f70a3ef6")
}

// UnpackLendingAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingAlreadySet(bytes32 symbol, address current)
func (marginAccounts *MarginAccounts) UnpackLendingAlreadySetError(raw []byte) (*MarginAccountsLendingAlreadySet, error) {
	out := new(MarginAccountsLendingAlreadySet)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "LendingAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsLengthMismatch represents a LengthMismatch error raised by the MarginAccounts contract.
type MarginAccountsLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LengthMismatch()
func MarginAccountsLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xff633a3803c58b9bc21e58efecee59f27e033cc0b1883fccb4969c76146fe60f")
}

// UnpackLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LengthMismatch()
func (marginAccounts *MarginAccounts) UnpackLengthMismatchError(raw []byte) (*MarginAccountsLengthMismatch, error) {
	out := new(MarginAccountsLengthMismatch)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "LengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsLiquidatorAlreadySet represents a LiquidatorAlreadySet error raised by the MarginAccounts contract.
type MarginAccountsLiquidatorAlreadySet struct {
	Current common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LiquidatorAlreadySet(address current)
func MarginAccountsLiquidatorAlreadySetErrorID() common.Hash {
	return common.HexToHash("0x0b7867fa5b88925561dbe5c988222c527c888cf13e383deb5eacbf8406c3470f")
}

// UnpackLiquidatorAlreadySetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LiquidatorAlreadySet(address current)
func (marginAccounts *MarginAccounts) UnpackLiquidatorAlreadySetError(raw []byte) (*MarginAccountsLiquidatorAlreadySet, error) {
	out := new(MarginAccountsLiquidatorAlreadySet)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "LiquidatorAlreadySet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsLiquidityUnknown represents a LiquidityUnknown error raised by the MarginAccounts contract.
type MarginAccountsLiquidityUnknown struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LiquidityUnknown(bytes32 symbol)
func MarginAccountsLiquidityUnknownErrorID() common.Hash {
	return common.HexToHash("0x1e89df1a02046ee0b3974a0ccb9d6a45cd1e47b611031710fdf5be384e85f153")
}

// UnpackLiquidityUnknownError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LiquidityUnknown(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackLiquidityUnknownError(raw []byte) (*MarginAccountsLiquidityUnknown, error) {
	out := new(MarginAccountsLiquidityUnknown)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "LiquidityUnknown", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsNoLending represents a NoLending error raised by the MarginAccounts contract.
type MarginAccountsNoLending struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoLending(bytes32 symbol)
func MarginAccountsNoLendingErrorID() common.Hash {
	return common.HexToHash("0xc915b4a5a597bdf2cb77203d45160a7dba65cd6242f9613549c84f621f676d04")
}

// UnpackNoLendingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoLending(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackNoLendingError(raw []byte) (*MarginAccountsNoLending, error) {
	out := new(MarginAccountsNoLending)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "NoLending", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsNotBackstop represents a NotBackstop error raised by the MarginAccounts contract.
type MarginAccountsNotBackstop struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotBackstop(address caller)
func MarginAccountsNotBackstopErrorID() common.Hash {
	return common.HexToHash("0xae8055b28d17e6998cbd5fb24affe947d018d63d8ee19fcf0b6e18781022dd40")
}

// UnpackNotBackstopError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotBackstop(address caller)
func (marginAccounts *MarginAccounts) UnpackNotBackstopError(raw []byte) (*MarginAccountsNotBackstop, error) {
	out := new(MarginAccountsNotBackstop)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "NotBackstop", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsNotGuardian represents a NotGuardian error raised by the MarginAccounts contract.
type MarginAccountsNotGuardian struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotGuardian(address caller)
func MarginAccountsNotGuardianErrorID() common.Hash {
	return common.HexToHash("0xa252c151cd1ee155ad6a596fd962f839a8b963fd05f5e863eca095de2e0b08a4")
}

// UnpackNotGuardianError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotGuardian(address caller)
func (marginAccounts *MarginAccounts) UnpackNotGuardianError(raw []byte) (*MarginAccountsNotGuardian, error) {
	out := new(MarginAccountsNotGuardian)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "NotGuardian", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsNotLiquidator represents a NotLiquidator error raised by the MarginAccounts contract.
type MarginAccountsNotLiquidator struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotLiquidator(address caller)
func MarginAccountsNotLiquidatorErrorID() common.Hash {
	return common.HexToHash("0xc2a9bf37d4428dc601ac6dd467dbb3cdcb8f405155a3671c168317f0e557142e")
}

// UnpackNotLiquidatorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotLiquidator(address caller)
func (marginAccounts *MarginAccounts) UnpackNotLiquidatorError(raw []byte) (*MarginAccountsNotLiquidator, error) {
	out := new(MarginAccountsNotLiquidator)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "NotLiquidator", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsNothingToSettle represents a NothingToSettle error raised by the MarginAccounts contract.
type MarginAccountsNothingToSettle struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NothingToSettle()
func MarginAccountsNothingToSettleErrorID() common.Hash {
	return common.HexToHash("0x618104e821a31c94c175b9bf52827d750880c5f9e96152367da18c1fe40d3ee0")
}

// UnpackNothingToSettleError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NothingToSettle()
func (marginAccounts *MarginAccounts) UnpackNothingToSettleError(raw []byte) (*MarginAccountsNothingToSettle, error) {
	out := new(MarginAccountsNothingToSettle)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "NothingToSettle", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsOutOfReach represents a OutOfReach error raised by the MarginAccounts contract.
type MarginAccountsOutOfReach struct {
	Symbol    [32]byte
	Amount    *big.Int
	Reachable *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OutOfReach(bytes32 symbol, uint256 amount, uint256 reachable)
func MarginAccountsOutOfReachErrorID() common.Hash {
	return common.HexToHash("0x5add5dac287f1c082ea8bfcabffd260df037368d999a725758cf30fd49c68d27")
}

// UnpackOutOfReachError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OutOfReach(bytes32 symbol, uint256 amount, uint256 reachable)
func (marginAccounts *MarginAccounts) UnpackOutOfReachError(raw []byte) (*MarginAccountsOutOfReach, error) {
	out := new(MarginAccountsOutOfReach)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "OutOfReach", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the MarginAccounts contract.
type MarginAccountsOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func MarginAccountsOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (marginAccounts *MarginAccounts) UnpackOwnableInvalidOwnerError(raw []byte) (*MarginAccountsOwnableInvalidOwner, error) {
	out := new(MarginAccountsOwnableInvalidOwner)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the MarginAccounts contract.
type MarginAccountsOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func MarginAccountsOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (marginAccounts *MarginAccounts) UnpackOwnableUnauthorizedAccountError(raw []byte) (*MarginAccountsOwnableUnauthorizedAccount, error) {
	out := new(MarginAccountsOwnableUnauthorizedAccount)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsOwnershipCannotBeRenounced represents a OwnershipCannotBeRenounced error raised by the MarginAccounts contract.
type MarginAccountsOwnershipCannotBeRenounced struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnershipCannotBeRenounced()
func MarginAccountsOwnershipCannotBeRenouncedErrorID() common.Hash {
	return common.HexToHash("0x2fab92ca4da7e80162387e93e02720bc4b98838c093385cac8fecf71a9b7de25")
}

// UnpackOwnershipCannotBeRenouncedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnershipCannotBeRenounced()
func (marginAccounts *MarginAccounts) UnpackOwnershipCannotBeRenouncedError(raw []byte) (*MarginAccountsOwnershipCannotBeRenounced, error) {
	out := new(MarginAccountsOwnershipCannotBeRenounced)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "OwnershipCannotBeRenounced", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsPositionNotEmpty represents a PositionNotEmpty error raised by the MarginAccounts contract.
type MarginAccountsPositionNotEmpty struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PositionNotEmpty()
func MarginAccountsPositionNotEmptyErrorID() common.Hash {
	return common.HexToHash("0x1acb203e88c4114a90e62ce2cfe3f792e8e517ea70f6d08da4f088fedebf0c35")
}

// UnpackPositionNotEmptyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PositionNotEmpty()
func (marginAccounts *MarginAccounts) UnpackPositionNotEmptyError(raw []byte) (*MarginAccountsPositionNotEmpty, error) {
	out := new(MarginAccountsPositionNotEmpty)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "PositionNotEmpty", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsRecallTooSmall represents a RecallTooSmall error raised by the MarginAccounts contract.
type MarginAccountsRecallTooSmall struct {
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error RecallTooSmall(uint256 amount)
func MarginAccountsRecallTooSmallErrorID() common.Hash {
	return common.HexToHash("0x6136f79f44ff957ae710e86ae17267ddd179a0b4e1accf66d53a9294d9cf7600")
}

// UnpackRecallTooSmallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error RecallTooSmall(uint256 amount)
func (marginAccounts *MarginAccounts) UnpackRecallTooSmallError(raw []byte) (*MarginAccountsRecallTooSmall, error) {
	out := new(MarginAccountsRecallTooSmall)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "RecallTooSmall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSafeCastOverflowedIntToUint represents a SafeCastOverflowedIntToUint error raised by the MarginAccounts contract.
type MarginAccountsSafeCastOverflowedIntToUint struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func MarginAccountsSafeCastOverflowedIntToUintErrorID() common.Hash {
	return common.HexToHash("0xa8ce4432b175c373e5f41aba830358e5361584f628450fd436c066323ad91ac2")
}

// UnpackSafeCastOverflowedIntToUintError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedIntToUint(int256 value)
func (marginAccounts *MarginAccounts) UnpackSafeCastOverflowedIntToUintError(raw []byte) (*MarginAccountsSafeCastOverflowedIntToUint, error) {
	out := new(MarginAccountsSafeCastOverflowedIntToUint)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SafeCastOverflowedIntToUint", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the MarginAccounts contract.
type MarginAccountsSafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func MarginAccountsSafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (marginAccounts *MarginAccounts) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*MarginAccountsSafeCastOverflowedUintDowncast, error) {
	out := new(MarginAccountsSafeCastOverflowedUintDowncast)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSafeCastOverflowedUintToInt represents a SafeCastOverflowedUintToInt error raised by the MarginAccounts contract.
type MarginAccountsSafeCastOverflowedUintToInt struct {
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func MarginAccountsSafeCastOverflowedUintToIntErrorID() common.Hash {
	return common.HexToHash("0x24775e0629ae69d78c11bae050651b81820407f300ff750ff2be51e4ce75c37f")
}

// UnpackSafeCastOverflowedUintToIntError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintToInt(uint256 value)
func (marginAccounts *MarginAccounts) UnpackSafeCastOverflowedUintToIntError(raw []byte) (*MarginAccountsSafeCastOverflowedUintToInt, error) {
	out := new(MarginAccountsSafeCastOverflowedUintToInt)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintToInt", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the MarginAccounts contract.
type MarginAccountsSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func MarginAccountsSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (marginAccounts *MarginAccounts) UnpackSafeERC20FailedOperationError(raw []byte) (*MarginAccountsSafeERC20FailedOperation, error) {
	out := new(MarginAccountsSafeERC20FailedOperation)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSequencerNotSettled represents a SequencerNotSettled error raised by the MarginAccounts contract.
type MarginAccountsSequencerNotSettled struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SequencerNotSettled()
func MarginAccountsSequencerNotSettledErrorID() common.Hash {
	return common.HexToHash("0xc6b5066d9e2ddc0279201222f9bed84f15770e598b3f5e7315b589326806202e")
}

// UnpackSequencerNotSettledError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SequencerNotSettled()
func (marginAccounts *MarginAccounts) UnpackSequencerNotSettledError(raw []byte) (*MarginAccountsSequencerNotSettled, error) {
	out := new(MarginAccountsSequencerNotSettled)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SequencerNotSettled", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsSessionUnknown represents a SessionUnknown error raised by the MarginAccounts contract.
type MarginAccountsSessionUnknown struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SessionUnknown()
func MarginAccountsSessionUnknownErrorID() common.Hash {
	return common.HexToHash("0xbeb78049726639f27e32542d3c6bcfd536c907b423b81f72607e067b8770b849")
}

// UnpackSessionUnknownError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SessionUnknown()
func (marginAccounts *MarginAccounts) UnpackSessionUnknownError(raw []byte) (*MarginAccountsSessionUnknown, error) {
	out := new(MarginAccountsSessionUnknown)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "SessionUnknown", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsTooManyAssets represents a TooManyAssets error raised by the MarginAccounts contract.
type MarginAccountsTooManyAssets struct {
	Count *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyAssets(uint256 count)
func MarginAccountsTooManyAssetsErrorID() common.Hash {
	return common.HexToHash("0x04628a1cd07f76a5fb8a39505083fb033bbfe16fbbad9948409565818d9c277f")
}

// UnpackTooManyAssetsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyAssets(uint256 count)
func (marginAccounts *MarginAccounts) UnpackTooManyAssetsError(raw []byte) (*MarginAccountsTooManyAssets, error) {
	out := new(MarginAccountsTooManyAssets)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "TooManyAssets", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsTooManyRecalls represents a TooManyRecalls error raised by the MarginAccounts contract.
type MarginAccountsTooManyRecalls struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TooManyRecalls(bytes32 symbol)
func MarginAccountsTooManyRecallsErrorID() common.Hash {
	return common.HexToHash("0x502686c4a8f4e23cf3de69678e23aa19c66e7adc57cfd8c179d1ee590f4b3b47")
}

// UnpackTooManyRecallsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TooManyRecalls(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackTooManyRecallsError(raw []byte) (*MarginAccountsTooManyRecalls, error) {
	out := new(MarginAccountsTooManyRecalls)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "TooManyRecalls", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsUnauthorized represents a Unauthorized error raised by the MarginAccounts contract.
type MarginAccountsUnauthorized struct {
	Caller  common.Address
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Unauthorized(address caller, address account)
func MarginAccountsUnauthorizedErrorID() common.Hash {
	return common.HexToHash("0x295a81c15f8ed30f90d2dbdefc77f34f2603dc4ea3df178ac206913987b62f4b")
}

// UnpackUnauthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Unauthorized(address caller, address account)
func (marginAccounts *MarginAccounts) UnpackUnauthorizedError(raw []byte) (*MarginAccountsUnauthorized, error) {
	out := new(MarginAccountsUnauthorized)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "Unauthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsUnknownAsset represents a UnknownAsset error raised by the MarginAccounts contract.
type MarginAccountsUnknownAsset struct {
	Symbol [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func MarginAccountsUnknownAssetErrorID() common.Hash {
	return common.HexToHash("0x1059be3ed9f1c1f061f09edf206f56e7c80b74fd01671cbec1b3788a1f49417e")
}

// UnpackUnknownAssetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownAsset(bytes32 symbol)
func (marginAccounts *MarginAccounts) UnpackUnknownAssetError(raw []byte) (*MarginAccountsUnknownAsset, error) {
	out := new(MarginAccountsUnknownAsset)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "UnknownAsset", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsUnknownPosition represents a UnknownPosition error raised by the MarginAccounts contract.
type MarginAccountsUnknownPosition struct {
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownPosition(bytes32 position)
func MarginAccountsUnknownPositionErrorID() common.Hash {
	return common.HexToHash("0xf47ffff133cc31b081f55ac85e5e6f3db0426380ee9eeb83b5ed8560a6e168db")
}

// UnpackUnknownPositionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownPosition(bytes32 position)
func (marginAccounts *MarginAccounts) UnpackUnknownPositionError(raw []byte) (*MarginAccountsUnknownPosition, error) {
	out := new(MarginAccountsUnknownPosition)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "UnknownPosition", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsUnsupportedToken represents a UnsupportedToken error raised by the MarginAccounts contract.
type MarginAccountsUnsupportedToken struct {
	Token    common.Address
	Position [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnsupportedToken(address token, bytes32 position)
func MarginAccountsUnsupportedTokenErrorID() common.Hash {
	return common.HexToHash("0x1c12b532afe32bea30a94024e6ab8a6d4197b8984adbbba1d7a111b22a73d58c")
}

// UnpackUnsupportedTokenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnsupportedToken(address token, bytes32 position)
func (marginAccounts *MarginAccounts) UnpackUnsupportedTokenError(raw []byte) (*MarginAccountsUnsupportedToken, error) {
	out := new(MarginAccountsUnsupportedToken)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "UnsupportedToken", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsWeekendDebtCapExceeded represents a WeekendDebtCapExceeded error raised by the MarginAccounts contract.
type MarginAccountsWeekendDebtCapExceeded struct {
	Debt *big.Int
	Cap  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WeekendDebtCapExceeded(uint256 debt, uint256 cap)
func MarginAccountsWeekendDebtCapExceededErrorID() common.Hash {
	return common.HexToHash("0xfecedbab478ffc20763e6c90ea91a8e33d5b128a91c80e5c4f2f1bb222a047a4")
}

// UnpackWeekendDebtCapExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WeekendDebtCapExceeded(uint256 debt, uint256 cap)
func (marginAccounts *MarginAccounts) UnpackWeekendDebtCapExceededError(raw []byte) (*MarginAccountsWeekendDebtCapExceeded, error) {
	out := new(MarginAccountsWeekendDebtCapExceeded)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "WeekendDebtCapExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsWeekendLeverageExceeded represents a WeekendLeverageExceeded error raised by the MarginAccounts contract.
type MarginAccountsWeekendLeverageExceeded struct {
	Gross  *big.Int
	Equity *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WeekendLeverageExceeded(uint256 gross, int256 equity)
func MarginAccountsWeekendLeverageExceededErrorID() common.Hash {
	return common.HexToHash("0xc31b6ccbfd61becd2cc8cc7f5970ce1ff409b28eb055fd1fa3f3e5d8cb4d0479")
}

// UnpackWeekendLeverageExceededError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WeekendLeverageExceeded(uint256 gross, int256 equity)
func (marginAccounts *MarginAccounts) UnpackWeekendLeverageExceededError(raw []byte) (*MarginAccountsWeekendLeverageExceeded, error) {
	out := new(MarginAccountsWeekendLeverageExceeded)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "WeekendLeverageExceeded", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// MarginAccountsZeroAmount represents a ZeroAmount error raised by the MarginAccounts contract.
type MarginAccountsZeroAmount struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAmount()
func MarginAccountsZeroAmountErrorID() common.Hash {
	return common.HexToHash("0x1f2a2005cb66a8e145327e8814e243a0996aec9bfe1e15a495778b1236dbd485")
}

// UnpackZeroAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAmount()
func (marginAccounts *MarginAccounts) UnpackZeroAmountError(raw []byte) (*MarginAccountsZeroAmount, error) {
	out := new(MarginAccountsZeroAmount)
	if err := marginAccounts.abi.UnpackIntoInterface(out, "ZeroAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
