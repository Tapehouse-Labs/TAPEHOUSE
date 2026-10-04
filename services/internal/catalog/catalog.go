// SPDX-License-Identifier: MIT OR Apache-2.0

// Package catalog lists the contracts a chain's registry names, each with the ABI of its Go SDK binding: the
// contracts the indexer follows and the API reads.
package catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapbackstop"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphoblue"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/sponsorpaymaster"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/supplyvault"
)

// Contract is a contract the registry names: its path in the registry without the leading dot, such as
// tapehouse.StockLending.NVDA, its address, the name and ABI of its binding, and the events the indexer keeps.
type Contract struct {
	Name    string
	Address common.Address
	Binding string
	ABI     *abi.ABI
	// Events are the events the indexer keeps, by topic: every event of the ABI for Tapehouse's own contracts, a
	// chosen few for an external contract it tracks, none for one it only reads.
	Events map[common.Hash]abi.Event
	// Tapehouse is set for Tapehouse's own contracts, whose earliest deployment is where indexing starts.
	Tapehouse bool
}

// StockTokenEvents are the events of a Stock Token the indexer keeps: the issuer's pause of the token and of its
// oracle, which halt the band, and its ERC-8056 multiplier updates.
var StockTokenEvents = []string{"OraclePaused", "OracleUnpaused", "Paused", "Unpaused", "UIMultiplierUpdated"}

// accounts are the registry's entries under .tapehouse that are keys, not contracts.
var accounts = []string{"HaltSigner", "Owner"}

// notStockTokens are the registry's tokens that are not Stock Tokens.
var notStockTokens = []string{"USDG", "WETH"}

var tapehouse = map[string]*bind.MetaData{
	"Band":             &band.BandMetaData,
	"Margin":           &margin.MarginMetaData,
	"SupplyVault":      &supplyvault.SupplyVaultMetaData,
	"MarginAccounts":   &marginaccounts.MarginAccountsMetaData,
	"Liquidator":       &liquidator.LiquidatorMetaData,
	"GapBackstop":      &gapbackstop.GapBackstopMetaData,
	"ReopeningAuction": &reopeningauction.ReopeningAuctionMetaData,
	"ShortPositions":   &shortpositions.ShortPositionsMetaData,
	"GapCover":         &gapcover.GapCoverMetaData,
	"SponsorPaymaster": &sponsorpaymaster.SponsorPaymasterMetaData,
}

var bindingNames = map[*bind.MetaData]string{
	&band.BandMetaData:                           "Band",
	&margin.MarginMetaData:                       "Margin",
	&supplyvault.SupplyVaultMetaData:             "SupplyVault",
	&marginaccounts.MarginAccountsMetaData:       "MarginAccounts",
	&liquidator.LiquidatorMetaData:               "Liquidator",
	&gapbackstop.GapBackstopMetaData:             "GapBackstop",
	&reopeningauction.ReopeningAuctionMetaData:   "ReopeningAuction",
	&shortpositions.ShortPositionsMetaData:       "ShortPositions",
	&gapcover.GapCoverMetaData:                   "GapCover",
	&sponsorpaymaster.SponsorPaymasterMetaData:   "SponsorPaymaster",
	&stocklendingvault.StockLendingVaultMetaData: "StockLendingVault",
	&basket.BasketMetaData:                       "Basket",
	&bandfeed.BandFeedMetaData:                   "BandFeed",
	&morphobandoracle.MorphoBandOracleMetaData:   "MorphoBandOracle",
	&morphoblue.MorphoBlueMetaData:               "MorphoBlue",
	&stocktoken.StockTokenMetaData:               "StockToken",
}

// Catalog is every contract of a registry, in a stable order.
type Catalog struct {
	ChainID   uint64
	Contracts []Contract
	byName    map[string]*Contract
	byAddress map[common.Address][]*Contract
}

// New lists the contracts d names: every contract under .tapehouse, its stock lending vaults and baskets, each
// .bandFeeds and .morphoOracles entry, and Morpho Blue, all events kept; and each Stock Token under .tokens, keeping
// StockTokenEvents. Morpho Blue is read, not indexed. An entry under .tapehouse with no binding is refused, so a
// contract added to the registry is never silently skipped.
func New(d *sdk.Deployments) (*Catalog, error) {
	c := &Catalog{ChainID: d.ChainID, byName: map[string]*Contract{}, byAddress: map[common.Address][]*Contract{}}
	for _, name := range sorted(d.Tapehouse) {
		if slices.Contains(accounts, name) {
			continue
		}
		metadata, ok := tapehouse[name]
		if !ok {
			return nil, fmt.Errorf("the indexer has no binding for .tapehouse.%s", name)
		}
		if err := c.add("tapehouse."+name, d.Tapehouse[name], metadata, nil, true); err != nil {
			return nil, err
		}
	}
	for _, group := range []struct {
		prefix    string
		entries   map[string]common.Address
		metadata  *bind.MetaData
		events    []string
		tapehouse bool
	}{
		{"tapehouse.StockLending.", d.StockLending, &stocklendingvault.StockLendingVaultMetaData, nil, true},
		{"tapehouse.Baskets.", d.Baskets, &basket.BasketMetaData, nil, true},
		{"bandFeeds.", d.BandFeeds, &bandfeed.BandFeedMetaData, nil, true},
		{"morphoOracles.", d.MorphoOracles, &morphobandoracle.MorphoBandOracleMetaData, nil, true},
		{"tokens.", stockTokens(d.Tokens), &stocktoken.StockTokenMetaData, StockTokenEvents, false},
	} {
		for _, name := range sorted(group.entries) {
			if err := c.add(group.prefix+name, group.entries[name], group.metadata, group.events, group.tapehouse); err != nil {
				return nil, err
			}
		}
	}
	if blue, ok := d.Morpho["Blue"]; ok {
		if err := c.add("morpho.Blue", blue, &morphoblue.MorphoBlueMetaData, []string{}, false); err != nil {
			return nil, err
		}
	}
	return c.index(), nil
}

func (c *Catalog) add(name string, address common.Address, metadata *bind.MetaData, events []string, tapehouse bool) error {
	parsed, err := metadata.ParseABI()
	if err != nil {
		return err
	}
	kept := map[common.Hash]abi.Event{}
	for _, event := range parsed.Events {
		if event.Anonymous {
			continue
		}
		if events == nil || slices.Contains(events, event.RawName) {
			kept[event.ID] = event
		}
	}
	c.Contracts = append(c.Contracts, Contract{name, address, bindingNames[metadata], parsed, kept, tapehouse})
	return nil
}

// index builds the catalog's lookups, once every contract is added.
func (c *Catalog) index() *Catalog {
	for i := range c.Contracts {
		contract := &c.Contracts[i]
		c.byName[contract.Name] = contract
		c.byAddress[contract.Address] = append(c.byAddress[contract.Address], contract)
	}
	return c
}

// Contract returns the contract the registry names name.
func (c *Catalog) Contract(name string) (*Contract, bool) {
	contract, ok := c.byName[name]
	return contract, ok
}

// At returns the contract at address that keeps the event with topic, as an address may appear under two names.
func (c *Catalog) At(address common.Address, topic common.Hash) (*Contract, abi.Event, bool) {
	for _, contract := range c.byAddress[address] {
		if event, ok := contract.Events[topic]; ok {
			return contract, event, true
		}
	}
	return nil, abi.Event{}, false
}

// Fingerprint identifies what the catalog indexes: each contract's name, address and kept events. An index built from
// another fingerprint is rebuilt.
func (c *Catalog) Fingerprint() string {
	var b strings.Builder
	fmt.Fprintf(&b, "chain %d\n", c.ChainID)
	for _, contract := range c.Contracts {
		fmt.Fprintf(&b, "%s %s %s", contract.Name, contract.Address.Hex(), contract.Binding)
		for _, topic := range contract.topics() {
			fmt.Fprintf(&b, " %s", contract.Events[topic].Sig)
		}
		b.WriteByte('\n')
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func (c *Contract) topics() []common.Hash {
	topics := make([]common.Hash, 0, len(c.Events))
	for topic := range c.Events {
		topics = append(topics, topic)
	}
	slices.SortFunc(topics, func(a, b common.Hash) int { return strings.Compare(a.Hex(), b.Hex()) })
	return topics
}

// Filters are the address sets and topics one eth_getLogs query each covers: the contracts whose every event is kept,
// with no topic filter, and the rest with the union of the events they keep.
func (c *Catalog) Filters() (all []common.Address, some []common.Address, topics []common.Hash) {
	seen := map[common.Hash]bool{}
	for _, contract := range c.Contracts {
		switch {
		case len(contract.Events) == 0:
		case len(contract.Events) == len(eventsOf(contract.ABI)):
			all = append(all, contract.Address)
		default:
			some = append(some, contract.Address)
			for _, topic := range contract.topics() {
				if !seen[topic] {
					seen[topic] = true
					topics = append(topics, topic)
				}
			}
		}
	}
	return all, some, topics
}

func eventsOf(parsed *abi.ABI) []abi.Event {
	var events []abi.Event
	for _, event := range parsed.Events {
		if !event.Anonymous {
			events = append(events, event)
		}
	}
	return events
}

// Tapehouse returns the addresses of Tapehouse's own contracts.
func (c *Catalog) Tapehouse() []common.Address {
	var out []common.Address
	for _, contract := range c.Contracts {
		if contract.Tapehouse {
			out = append(out, contract.Address)
		}
	}
	return out
}

func stockTokens(tokens map[string]common.Address) map[string]common.Address {
	out := map[string]common.Address{}
	for name, address := range tokens {
		if !slices.Contains(notStockTokens, name) {
			out[name] = address
		}
	}
	return out
}

func sorted[V any](group map[string]V) []string {
	names := make([]string, 0, len(group))
	for name := range group {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
