// SPDX-License-Identifier: MIT OR Apache-2.0

package sponsor

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapbackstop"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/supplyvault"
)

// userFunctions are the functions a user calls, by contract: the service sponsors these, and approvals of their
// contracts, and no keeper's, owner's or bidder's call.
var userFunctions = []struct {
	group     func(d *sdk.Deployments) map[string]common.Address
	name      string
	meta      *bind.MetaData
	functions []string
}{
	{tapehouse, "MarginAccounts", &marginaccounts.MarginAccountsMetaData, []string{"deposit", "depositWithPermit",
		"withdraw", "borrow", "repay", "repayWithCollateral", "setAuthorization", "lend", "unlend", "recall", "settle",
		"unwrap"}},
	{tapehouse, "SupplyVault", &supplyvault.SupplyVaultMetaData, []string{"deposit", "mint", "withdraw", "redeem"}},
	{tapehouse, "ShortPositions", &shortpositions.ShortPositionsMetaData, []string{"sell", "cover", "deposit", "withdraw"}},
	{tapehouse, "GapCover", &gapcover.GapCoverMetaData, []string{"buy", "claim", "deposit", "mint", "withdraw", "redeem"}},
	{tapehouse, "GapBackstop", &gapbackstop.GapBackstopMetaData, []string{"deposit", "mint", "startCooldown", "withdraw",
		"redeem", "claimGains"}},
	{stockLending, "", &stocklendingvault.StockLendingVaultMetaData, []string{"deposit", "mint", "withdraw", "redeem"}},
	{baskets, "", &basket.BasketMetaData, []string{"mint", "redeem"}},
}

func tapehouse(d *sdk.Deployments) map[string]common.Address    { return d.Tapehouse }
func stockLending(d *sdk.Deployments) map[string]common.Address { return d.StockLending }
func baskets(d *sdk.Deployments) map[string]common.Address      { return d.Baskets }

// sponsoredFunctions maps each contract of the registry a user calls to the selectors of its user functions.
func sponsoredFunctions(d *sdk.Deployments) (map[common.Address]map[[4]byte]bool, error) {
	targets := map[common.Address]map[[4]byte]bool{}
	for _, contract := range userFunctions {
		parsed, err := contract.meta.ParseABI()
		if err != nil {
			return nil, err
		}
		selectors := map[[4]byte]bool{}
		for _, name := range contract.functions {
			found := false
			for _, method := range parsed.Methods {
				if method.RawName == name {
					selectors[[4]byte(method.ID)] = true
					found = true
				}
			}
			if !found {
				return nil, fmt.Errorf("the ABI has no function %s", name)
			}
		}
		for name, address := range contract.group(d) {
			if contract.name == "" || name == contract.name {
				targets[address] = selectors
			}
		}
	}
	return targets, nil
}
