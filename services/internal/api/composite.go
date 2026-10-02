// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/catalog"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// SignedInput labels the halt: the one value Tapehouse signs.
const SignedInput = "A halt signed by Tapehouse's halt signer is the only value Tapehouse signs; every price the band " +
	"uses is signed by RedStone or Chainlink."

// atHead serves what read returns at the index's head, read once per head.
func (s *Server) atHead(w http.ResponseWriter, r *http.Request, name string, read func(ctx context.Context, block uint64) (map[string]any, error)) {
	head, err := s.head(r.Context())
	if err != nil {
		s.reply(w, nil, err)
		return
	}
	out, err := s.cache.get(fmt.Sprintf("%s %d", name, head.Number), func() (any, error) {
		out, err := read(r.Context(), head.Number)
		if out != nil {
			out["block"] = head.Number
		}
		return out, err
	})
	s.reply(w, out, err)
}

// contract returns the contract the registry names, or why it names none.
func (s *Server) contract(name string) (*catalog.Contract, error) {
	contract, ok := s.Catalog.Contract(name)
	if !ok {
		return nil, missing("the registry of chain %d names no %s", s.Catalog.ChainID, name)
	}
	return contract, nil
}

// views reads each of a contract's views at block, refusing a revert.
func (s *Server) views(ctx context.Context, contract *catalog.Contract, block uint64, reads map[string][]any) (map[string]map[string]any, error) {
	out := make(map[string]map[string]any, len(reads))
	for _, function := range slices.Sorted(maps.Keys(reads)) {
		outputs, revert, err := s.call(ctx, contract, function, block, reads[function]...)
		if err != nil {
			return nil, err
		}
		if revert != nil {
			return nil, fmt.Errorf("%s.%s reverted with %s", contract.Name, function, revert.Name)
		}
		out[function] = outputs
	}
	return out, nil
}

// debts serves each position of the margin accounts with debt, for the liquidation keeper: every position that ever
// borrowed, with its debt and health at the head, where it still owes.
func (s *Server) debts(w http.ResponseWriter, r *http.Request) {
	s.atHead(w, r, "debts", func(ctx context.Context, block uint64) (map[string]any, error) {
		accounts, err := s.contract("tapehouse.MarginAccounts")
		if err != nil {
			return nil, err
		}
		pairs, err := s.Store.Pairs(ctx, accounts.Name, "Borrow", "account", "position")
		if err != nil {
			return nil, err
		}
		positions := []map[string]any{}
		for _, pair := range pairs {
			account, position := common.HexToAddress(pair[0]), common.HexToHash(pair[1])
			read, err := s.views(ctx, accounts, block, map[string][]any{
				"debt": {account, [32]byte(position)}, "health": {account, [32]byte(position)},
			})
			if err != nil {
				return nil, err
			}
			if read["debt"]["arg0"] == "0" {
				continue
			}
			positions = append(positions, map[string]any{"account": pair[0], "position": pair[1],
				"debt": read["debt"]["arg0"], "health": read["health"]})
		}
		return map[string]any{"positions": positions}, nil
	})
}

// shorts serves each open short: every account and Stock Token that ever sold short, with its position and health at
// the head, where it still holds borrow shares.
func (s *Server) shorts(w http.ResponseWriter, r *http.Request) {
	s.atHead(w, r, "shorts", func(ctx context.Context, block uint64) (map[string]any, error) {
		shorts, err := s.contract("tapehouse.ShortPositions")
		if err != nil {
			return nil, err
		}
		pairs, err := s.Store.Pairs(ctx, shorts.Name, "Sell", "account", "symbol")
		if err != nil {
			return nil, err
		}
		open := []map[string]any{}
		for _, pair := range pairs {
			account, symbol := common.HexToAddress(pair[0]), common.HexToHash(pair[1])
			read, err := s.views(ctx, shorts, block, map[string][]any{
				"position": {account, [32]byte(symbol)}, "health": {account, [32]byte(symbol)},
			})
			if err != nil {
				return nil, err
			}
			if read["position"]["shares"] == "0" {
				continue
			}
			open = append(open, map[string]any{"account": pair[0], "symbol": symbolText(symbol),
				"position": read["position"], "health": read["health"]})
		}
		return map[string]any{"shorts": open}, nil
	})
}

// halts serves the band's halt of each asset the registry names, labelled as Tapehouse's one signed input, and the
// signer the band accepts.
func (s *Server) halts(w http.ResponseWriter, r *http.Request) {
	s.atHead(w, r, "halts", func(ctx context.Context, block uint64) (map[string]any, error) {
		band, err := s.contract("tapehouse.Band")
		if err != nil {
			return nil, err
		}
		signer, err := s.views(ctx, band, block, map[string][]any{"haltSigner": nil})
		if err != nil {
			return nil, err
		}
		assets := map[string]bool{}
		for name := range s.Deployments.BandFeeds {
			assets[name] = true
		}
		for _, contract := range s.Catalog.Contracts {
			if name, ok := strings.CutPrefix(contract.Name, "tokens."); ok {
				assets[name] = true
			}
		}
		halts := []map[string]any{}
		for _, asset := range slices.Sorted(maps.Keys(assets)) {
			symbol, err := sdk.ToBytes32(asset)
			if err != nil {
				return nil, err
			}
			read, err := s.views(ctx, band, block, map[string][]any{"halt": {symbol}})
			if err != nil {
				return nil, err
			}
			halt := read["halt"]
			halt["symbol"] = asset
			halts = append(halts, halt)
		}
		return map[string]any{"signedInput": SignedInput, "haltSigner": signer["haltSigner"]["arg0"],
			"halts": halts}, nil
	})
}

// morphoOracles serves each Morpho oracle: its band, symbol, tokens, scale and owner, and its price, or why it has
// none, never zero, as the Go SDK reads them at one block.
func (s *Server) morphoOracles(w http.ResponseWriter, r *http.Request) {
	s.atHead(w, r, "morphoOracles", func(ctx context.Context, block uint64) (map[string]any, error) {
		oracles := []map[string]any{}
		for _, asset := range slices.Sorted(maps.Keys(s.Deployments.MorphoOracles)) {
			oracle, err := s.sdk.MorphoOracles().Oracle(at(ctx, block), asset)
			if err != nil {
				return nil, err
			}
			price, err := s.sdk.MorphoOracles().Price(at(ctx, block), asset)
			if err != nil {
				return nil, err
			}
			served := map[string]any{"asset": asset, "address": oracle.Address.Hex(), "band": oracle.Band.Hex(),
				"symbol": symbolText(oracle.Symbol), "collateralToken": oracle.CollateralToken.Hex(),
				"loanToken": oracle.LoanToken.Hex(), "scaleFactor": oracle.ScaleFactor.String(),
				"owner": oracle.Owner.Hex(), "price": nil, "noPrice": nil}
			if price.Price != nil {
				served["price"] = price.Price.String()
			} else {
				served["noPrice"] = string(price.NoPrice)
			}
			oracles = append(oracles, served)
		}
		return map[string]any{"oracles": oracles}, nil
	})
}

// morphoMarkets serves each Morpho Blue market the registry records, by id: its parameters and its state.
func (s *Server) morphoMarkets(w http.ResponseWriter, r *http.Request) {
	s.atHead(w, r, "morphoMarkets", func(ctx context.Context, block uint64) (map[string]any, error) {
		markets := []map[string]any{}
		if len(s.Deployments.MorphoMarkets) == 0 {
			return map[string]any{"markets": markets}, nil
		}
		blue, err := s.contract("morpho.Blue")
		if err != nil {
			return nil, err
		}
		for _, name := range slices.Sorted(maps.Keys(s.Deployments.MorphoMarkets)) {
			id := s.Deployments.MorphoMarkets[name]
			read, err := s.views(ctx, blue, block, map[string][]any{"idToMarketParams": {[32]byte(id)}, "market": {[32]byte(id)}})
			if err != nil {
				return nil, err
			}
			markets = append(markets, map[string]any{"name": name, "id": id, "params": read["idToMarketParams"]["arg0"],
				"market": read["market"]})
		}
		return map[string]any{"markets": markets}, nil
	})
}

// symbolText is a symbol as text where its bytes are printable, as hex otherwise.
func symbolText(symbol [32]byte) string {
	text := strings.TrimRight(string(symbol[:]), "\x00")
	for _, c := range text {
		if c < 0x20 || c > 0x7e {
			return hexutil.Encode(symbol[:])
		}
	}
	return text
}
