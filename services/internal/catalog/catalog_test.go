// SPDX-License-Identifier: MIT OR Apache-2.0

package catalog_test

import (
	"fmt"
	"math/big"
	"reflect"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/sdk"
)

func address(n int) string {
	return common.BigToAddress(new(big.Int).Lsh(common.Big1, uint(n))).Hex()
}

// devnode is a registry shaped as the dev node's: every group, the halt signer and owner among .tapehouse.
func devnode(t *testing.T, extra string) *sdk.Deployments {
	t.Helper()
	d, err := sdk.ParseDeployments(fmt.Appendf(nil, `{"chainId":412346,
		"tokens":{"NVDA":"%s","SPY":"%s","USDG":"%s","WETH":"%s"},
		"tapehouse":{"HaltSigner":"%s","Owner":"%s","Band":"%s","Margin":"%s","SupplyVault":"%s","MarginAccounts":"%s",
			"Liquidator":"%s","GapBackstop":"%s","ReopeningAuction":"%s","ShortPositions":"%s","GapCover":"%s",
			"StockLending":{"NVDA":"%s","SPY":"%s"},"Baskets":{"PAIR":"%s"}%s},
		"morphoOracles":{"NVDA":"%s"},"bandFeeds":{"NVDA":"%s","SPY":"%s"},"morpho":{"Blue":"%s"}}`,
		address(1), address(2), address(3), address(4), address(5), address(6), address(7), address(8), address(9),
		address(10), address(11), address(12), address(13), address(14), address(15), address(16), address(17),
		address(18), extra, address(19), address(20), address(21), address(22)))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestEveryContractOfTheRegistryIsListedWithItsBinding(t *testing.T) {
	c, err := catalog.New(devnode(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, contract := range c.Contracts {
		got = append(got, contract.Name+" "+contract.Binding)
	}
	want := []string{
		"tapehouse.Band Band", "tapehouse.GapBackstop GapBackstop", "tapehouse.GapCover GapCover",
		"tapehouse.Liquidator Liquidator", "tapehouse.Margin Margin", "tapehouse.MarginAccounts MarginAccounts",
		"tapehouse.ReopeningAuction ReopeningAuction", "tapehouse.ShortPositions ShortPositions",
		"tapehouse.SupplyVault SupplyVault", "tapehouse.StockLending.NVDA StockLendingVault",
		"tapehouse.StockLending.SPY StockLendingVault", "tapehouse.Baskets.PAIR Basket", "bandFeeds.NVDA BandFeed",
		"bandFeeds.SPY BandFeed", "morphoOracles.NVDA MorphoBandOracle", "tokens.NVDA StockToken",
		"tokens.SPY StockToken", "morpho.Blue MorphoBlue",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("contracts:\n%v\nwant\n%v", got, want)
	}
	if len(c.Tapehouse()) != 15 {
		t.Fatalf("Tapehouse's own contracts: %v", c.Tapehouse())
	}
	band, _ := c.Contract("tapehouse.Band")
	for _, name := range []string{"PriceWritten", "Anchored", "MultiplierRecorded", "MultiplierConfirmed", "HaltWritten",
		"HaltSignerUpdated", "OwnershipTransferStarted", "OwnershipTransferred"} {
		if _, _, ok := c.At(band.Address, band.ABI.Events[name].ID); !ok {
			t.Errorf("the band's %s is not kept", name)
		}
	}
	cover, _ := c.Contract("tapehouse.GapCover")
	for _, name := range []string{"Measured", "Bought", "Observed", "Settled", "Voided"} {
		if _, ok := cover.Events[cover.ABI.Events[name].ID]; !ok {
			t.Errorf("the gap cover's %s is not kept", name)
		}
	}
	margin, _ := c.Contract("tapehouse.Margin")
	if _, ok := margin.Events[margin.ABI.Events["CorrelationSet"].ID]; !ok || len(margin.Events) != 6 {
		t.Errorf("the engine keeps %d events", len(margin.Events))
	}
}

func TestAStockTokenKeepsItsHaltsAndMultiplierAndMorphoBlueNothing(t *testing.T) {
	c, err := catalog.New(devnode(t, ""))
	if err != nil {
		t.Fatal(err)
	}
	token, _ := c.Contract("tokens.NVDA")
	var kept []string
	for _, event := range token.Events {
		kept = append(kept, event.RawName)
	}
	slices.Sort(kept)
	if !reflect.DeepEqual(kept, []string{"OraclePaused", "OracleUnpaused", "Paused", "UIMultiplierUpdated", "Unpaused"}) {
		t.Fatalf("a Stock Token keeps %v", kept)
	}
	if _, _, ok := c.At(token.Address, token.ABI.Events["Transfer"].ID); ok {
		t.Fatal("a Stock Token's transfers are kept")
	}
	blue, _ := c.Contract("morpho.Blue")
	if len(blue.Events) != 0 || blue.Tapehouse {
		t.Fatalf("Morpho Blue: %+v", blue.Events)
	}
	all, some, topics := c.Filters()
	if len(all) != 15 || len(some) != 2 || len(topics) != 5 || slices.Contains(all, blue.Address) {
		t.Fatalf("filters: %d addresses with every event, %d with %d topics", len(all), len(some), len(topics))
	}
}

func TestAContractWithoutABindingIsRefusedAndTheFingerprintFollowsTheRegistry(t *testing.T) {
	if _, err := catalog.New(devnode(t, `,"Options":"`+address(30)+`"`)); err == nil ||
		err.Error() != "the indexer has no binding for .tapehouse.Options" {
		t.Fatalf("an unknown contract: %v", err)
	}
	first, _ := catalog.New(devnode(t, ""))
	again, _ := catalog.New(devnode(t, ""))
	if first.Fingerprint() != again.Fingerprint() {
		t.Fatal("the same registry has two fingerprints")
	}
	moved := devnode(t, "")
	moved.BandFeeds["SPY"] = common.HexToAddress(address(31))
	other, _ := catalog.New(moved)
	if other.Fingerprint() == first.Fingerprint() {
		t.Fatal("a moved feed kept the fingerprint")
	}
}

func TestTheRegistriesInTheRepositoryAreCatalogued(t *testing.T) {
	for chain, contracts := range map[string]int{"4663": 7, "46630": 12, "42161": 7} {
		d, err := sdk.LoadDeployments("../../../deployments/" + chain + ".json")
		if err != nil {
			t.Fatal(err)
		}
		c, err := catalog.New(d)
		if err != nil || len(c.Contracts) != contracts {
			t.Errorf("%s: %d contracts, %v", chain, len(c.Contracts), err)
		}
	}
}
