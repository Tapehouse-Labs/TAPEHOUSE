// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Index reads the indexer's API for what the chain cannot enumerate: the positions with debt, the open shorts and the
// events that name recalls, deposits, commitments and covers. Every judgement is then read from the chain.
type Index struct {
	base   string
	client *http.Client
}

// NewIndex returns the index API at base, such as the indexer's TAPEHOUSE_INDEXER_URL, read over client.
func NewIndex(base string, client *http.Client) *Index {
	return &Index{base: strings.TrimRight(base, "/"), client: client}
}

// Position is a margin position: its account and position, and its equity and requirement as the indexer last read
// them, which order the keepers' work and decide nothing.
type Position struct {
	Account     common.Address
	Position    [32]byte
	Equity      *big.Int
	Requirement *big.Int
}

// Short is an open short: its account and asset, its USDG, negative in deficit, and its equity and requirement as
// the indexer last read them.
type Short struct {
	Account     common.Address
	Asset       string
	USDG        *big.Int
	Equity      *big.Int
	Requirement *big.Int
}

// Event is an indexed event's arguments, as the API serves them: integers as decimal strings, addresses checksummed and
// bytes as hex.
type Event struct {
	Block    uint64         `json:"block"`
	LogIndex uint           `json:"logIndex"`
	Args     map[string]any `json:"args"`
}

// Arg returns the event's argument name as text.
func (e Event) Arg(name string) string {
	return fmt.Sprint(e.Args[name])
}

type health struct {
	Equity      string `json:"equity"`
	Requirement string `json:"requirement"`
}

// Debts reads every position with debt.
func (i *Index) Debts(ctx context.Context) ([]Position, error) {
	var out struct {
		Positions []struct {
			Account  common.Address `json:"account"`
			Position string         `json:"position"`
			Health   health         `json:"health"`
		} `json:"positions"`
	}
	if err := i.get(ctx, "/v1/accounts/debts", nil, &out); err != nil {
		return nil, err
	}
	positions := make([]Position, len(out.Positions))
	for n, p := range out.Positions {
		position, err := word(p.Position)
		if err != nil {
			return nil, err
		}
		positions[n] = Position{p.Account, position, integer(p.Health.Equity), integer(p.Health.Requirement)}
	}
	return positions, nil
}

// Shorts reads every open short.
func (i *Index) Shorts(ctx context.Context) ([]Short, error) {
	var out struct {
		Shorts []struct {
			Account  common.Address `json:"account"`
			Symbol   string         `json:"symbol"`
			Position struct {
				USDG string `json:"usdgHeld"`
			} `json:"position"`
			Health health `json:"health"`
		} `json:"shorts"`
	}
	if err := i.get(ctx, "/v1/shorts", nil, &out); err != nil {
		return nil, err
	}
	shorts := make([]Short, len(out.Shorts))
	for n, s := range out.Shorts {
		shorts[n] = Short{s.Account, s.Symbol, integer(s.Position.USDG), integer(s.Health.Equity),
			integer(s.Health.Requirement)}
	}
	return shorts, nil
}

// Events reads every event of contract, by its registry name, named event and carrying args, oldest first.
func (i *Index) Events(ctx context.Context, contract, event string, args map[string]string) ([]Event, error) {
	query := url.Values{"contract": {contract}, "event": {event}, "limit": {"1000"}}
	for name, value := range args {
		query.Set("arg["+name+"]", value)
	}
	var events []Event
	for {
		var page struct {
			Events []Event `json:"events"`
			Next   *string `json:"next"`
		}
		if err := i.get(ctx, "/v1/events", query, &page); err != nil {
			return nil, err
		}
		events = append(events, page.Events...)
		if page.Next == nil {
			return events, nil
		}
		query.Set("after", *page.Next)
	}
}

func (i *Index) get(ctx context.Context, path string, query url.Values, out any) error {
	target := i.base + path
	if query != nil {
		target += "?" + query.Encode()
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	response, err := i.client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("the index answered %s to %s", response.Status, path)
	}
	return json.NewDecoder(response.Body).Decode(out)
}

// deepestFirst orders positions by how far short they fall, (requirement - equity) / requirement, deepest first.
func deepestFirst[T any](items []T, of func(T) (equity, requirement *big.Int)) {
	depth := func(item T) *big.Rat {
		equity, requirement := of(item)
		if requirement == nil || requirement.Sign() == 0 || equity == nil {
			return new(big.Rat)
		}
		return new(big.Rat).SetFrac(new(big.Int).Sub(requirement, equity), requirement)
	}
	sort.SliceStable(items, func(a, b int) bool { return depth(items[a]).Cmp(depth(items[b])) > 0 })
}

func integer(text string) *big.Int {
	value, ok := new(big.Int).SetString(text, 10)
	if !ok {
		return new(big.Int)
	}
	return value
}

func word(hex string) ([32]byte, error) {
	raw, err := hexutil.Decode(hex)
	if err != nil || len(raw) != 32 {
		return [32]byte{}, fmt.Errorf("%q is not a 32-byte word", hex)
	}
	return [32]byte(raw), nil
}
