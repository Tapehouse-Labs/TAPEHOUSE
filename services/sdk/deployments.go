// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// SharePriceChains are the chains whose Chainlink feeds price the share rather than the Stock Token.
var SharePriceChains = []uint64{42161}

// Cross is the margin accounts' cross position.
var Cross [32]byte

// Deployments is a chain's address registry, deployments/<chainId>.json.
type Deployments struct {
	ChainID      uint64
	Tokens       map[string]common.Address
	Chainlink    map[string]common.Address
	BandFeeds    map[string]common.Address
	Tapehouse    map[string]common.Address
	StockLending map[string]common.Address
	UniswapV3    map[string]common.Address
}

// TokenPriceFeed is a Chainlink feed that prices the Stock Token, as Robinhood Chain's do.
type TokenPriceFeed struct{ Address common.Address }

// SharePriceFeed is a Chainlink feed that prices the share, as Arbitrum One's do.
type SharePriceFeed struct{ Address common.Address }

// LoadDeployments reads a chain's registry from path.
func LoadDeployments(path string) (*Deployments, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseDeployments(data)
}

// ParseDeployments parses a chain's registry. An address in mixed case must carry its checksum.
func ParseDeployments(data []byte) (*Deployments, error) {
	var registry struct {
		ChainID   uint64                     `json:"chainId"`
		Tokens    map[string]string          `json:"tokens"`
		Chainlink map[string]string          `json:"chainlink"`
		BandFeeds map[string]string          `json:"bandFeeds"`
		Tapehouse map[string]json.RawMessage `json:"tapehouse"`
		UniswapV3 map[string]string          `json:"uniswapV3"`
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}
	if registry.ChainID == 0 {
		return nil, fmt.Errorf("the registry has no chainId")
	}
	d := &Deployments{ChainID: registry.ChainID, Tapehouse: map[string]common.Address{}}
	var err error
	for _, group := range []struct {
		into *map[string]common.Address
		from map[string]string
		path string
	}{
		{&d.Tokens, registry.Tokens, ".tokens"},
		{&d.Chainlink, registry.Chainlink, ".chainlink"},
		{&d.BandFeeds, registry.BandFeeds, ".bandFeeds"},
		{&d.UniswapV3, registry.UniswapV3, ".uniswapV3"},
	} {
		if *group.into, err = addresses(group.from, group.path); err != nil {
			return nil, err
		}
	}
	lending := map[string]string{}
	for name, raw := range registry.Tapehouse {
		target := new(string)
		if name == "StockLending" {
			if err := json.Unmarshal(raw, &lending); err != nil {
				return nil, fmt.Errorf(".tapehouse.StockLending: %w", err)
			}
			continue
		}
		if err := json.Unmarshal(raw, target); err != nil {
			return nil, fmt.Errorf(".tapehouse.%s is not an address", name)
		}
		if d.Tapehouse[name], err = address(*target, ".tapehouse."+name); err != nil {
			return nil, err
		}
	}
	if d.StockLending, err = addresses(lending, ".tapehouse.StockLending"); err != nil {
		return nil, err
	}
	return d, nil
}

// SharePrices reports whether the chain's Chainlink feeds price the share rather than the Stock Token.
func (d *Deployments) SharePrices() bool {
	return slices.Contains(SharePriceChains, d.ChainID)
}

// TokenPriceFeed returns the Chainlink feed name, refusing a chain whose feeds price the share.
func (d *Deployments) TokenPriceFeed(name string) (TokenPriceFeed, error) {
	feed, err := entry(d.Chainlink, name, ".chainlink")
	if err == nil && d.SharePrices() {
		err = fmt.Errorf(".chainlink.%s prices the share on chain %d", name, d.ChainID)
	}
	return TokenPriceFeed{feed}, err
}

// SharePriceFeed returns the Chainlink feed name, refusing a chain whose feeds price the Stock Token.
func (d *Deployments) SharePriceFeed(name string) (SharePriceFeed, error) {
	feed, err := entry(d.Chainlink, name, ".chainlink")
	if err == nil && !d.SharePrices() {
		err = fmt.Errorf(".chainlink.%s prices the Stock Token on chain %d", name, d.ChainID)
	}
	return SharePriceFeed{feed}, err
}

// ToBytes32 returns an asset's symbol, or a RedStone feed ID, as the contracts take it: its bytes, right-padded to 32.
// A name longer than 32 bytes is refused.
func ToBytes32(name string) ([32]byte, error) {
	var out [32]byte
	if len(name) > len(out) {
		return out, fmt.Errorf("%q is longer than 32 bytes", name)
	}
	copy(out[:], name)
	return out, nil
}

func entry(group map[string]common.Address, name, path string) (common.Address, error) {
	address, ok := group[name]
	if !ok {
		return common.Address{}, fmt.Errorf("the registry has no %s.%s", path, name)
	}
	return address, nil
}

func addresses(group map[string]string, path string) (map[string]common.Address, error) {
	out := make(map[string]common.Address, len(group))
	for name, value := range group {
		parsed, err := address(value, path+"."+name)
		if err != nil {
			return nil, err
		}
		out[name] = parsed
	}
	return out, nil
}

func address(value, path string) (common.Address, error) {
	parsed := common.HexToAddress(value)
	if !common.IsHexAddress(value) || (strings.ToLower(value) != value && parsed.Hex() != value) {
		return common.Address{}, fmt.Errorf("%s is not an address", path)
	}
	return parsed, nil
}
