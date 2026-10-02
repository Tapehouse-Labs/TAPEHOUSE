// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/codec"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// Event is an indexed event as the API serves it: the store's, with the emitter's address and whether its block is
// final.
type Event struct {
	store.Event
	Address   codec.Address `json:"address"`
	Finalized bool          `json:"finalized"`
}

func (s *Server) served(events []store.Event) []Event {
	out := make([]Event, len(events))
	finalized := s.Index.Finalized()
	for i, event := range events {
		out[i] = Event{Event: event, Finalized: event.Block <= finalized}
		if contract, ok := s.Catalog.Contract(event.Contract); ok {
			out[i].Address = codec.Address(contract.Address)
		}
	}
	return out
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	head, err := s.head(r.Context())
	s.reply(w, map[string]any{
		"chainId":   s.Catalog.ChainID,
		"start":     s.Index.Start(),
		"head":      head,
		"finalized": s.Index.Finalized(),
	}, err)
}

func (s *Server) contracts(w http.ResponseWriter, _ *http.Request) {
	type served struct {
		Name    string        `json:"name"`
		Address codec.Address `json:"address"`
		Binding string        `json:"binding"`
		Events  []string      `json:"events"`
	}
	out := make([]served, len(s.Catalog.Contracts))
	for i, contract := range s.Catalog.Contracts {
		out[i] = served{contract.Name, codec.Address(contract.Address), contract.Binding, []string{}}
		for _, event := range contract.Events {
			out[i].Events = append(out[i].Events, event.Sig)
		}
		slices.Sort(out[i].Events)
	}
	s.reply(w, map[string]any{"contracts": out}, nil)
}

// MaxPage is the most events one page of /v1/events holds.
const MaxPage = 1000

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	page, err := s.eventQuery(r)
	if err != nil {
		s.reply(w, nil, err)
		return
	}
	events, err := s.Store.Events(r.Context(), page)
	var next *string
	if err == nil && len(events) == page.Limit {
		cursor := store.Cursor{Block: events[len(events)-1].Block, LogIndex: events[len(events)-1].LogIndex}.String()
		next = &cursor
	}
	s.reply(w, map[string]any{"events": s.served(events), "next": next}, err)
}

// eventQuery reads /v1/events' parameters: contract, event, from, to, after, limit and arg[<name>], an argument the
// event must carry, which needs the contract and the event named so it is read as the event's ABI types it.
func (s *Server) eventQuery(r *http.Request) (store.Query, error) {
	query := r.URL.Query()
	q := store.Query{Contract: query.Get("contract"), Event: query.Get("event"), Args: map[string]any{}, Limit: 100}
	var err error
	if q.From, err = uintParam(r, "from"); err != nil {
		return q, err
	}
	if q.To, err = uintParam(r, "to"); err != nil {
		return q, err
	}
	if after := query.Get("after"); after != "" {
		cursor, err := store.ParseCursor(after)
		if err != nil {
			return q, invalid("after: %v", err)
		}
		q.After = &cursor
	}
	if limit := query.Get("limit"); limit != "" {
		if q.Limit, err = strconv.Atoi(limit); err != nil || q.Limit < 1 || q.Limit > MaxPage {
			return q, invalid("limit must be from 1 to %d", MaxPage)
		}
	}
	for key, values := range query {
		switch name, ok := argName(key); {
		case ok:
			value, err := s.argFilter(q.Contract, q.Event, name, values[0])
			if err != nil {
				return q, err
			}
			q.Args[name] = value
		case !slices.Contains([]string{"contract", "event", "from", "to", "after", "limit"}, key):
			return q, invalid("unknown parameter %s", key)
		}
	}
	if q.Contract != "" {
		if _, ok := s.Catalog.Contract(q.Contract); !ok {
			return q, missing("the registry names no contract %s", q.Contract)
		}
	}
	return q, nil
}

func (s *Server) argFilter(contractName, eventName, name, text string) (any, error) {
	contract, ok := s.Catalog.Contract(contractName)
	if !ok || eventName == "" {
		return nil, invalid("arg[%s] needs a contract and an event", name)
	}
	for _, event := range contract.Events {
		if event.RawName != eventName {
			continue
		}
		for i, input := range event.Inputs {
			if input.Name != name && "arg"+strconv.Itoa(i) != name {
				continue
			}
			if input.Indexed && (input.Type.T == abi.StringTy || input.Type.T == abi.BytesTy) {
				return text, nil
			}
			value, err := codec.Parse(input.Type, text)
			if err != nil {
				return nil, invalid("arg[%s]: %v", name, err)
			}
			switch json := codec.Value(input.Type, value).(type) {
			case bool:
				return map[bool]int{false: 0, true: 1}[json], nil
			case string:
				return json, nil
			}
			return nil, invalid("arg[%s] of type %s cannot be filtered on", name, input.Type)
		}
		return nil, invalid("%s has no argument %s", event.Sig, name)
	}
	return nil, missing("%s keeps no event %s", contractName, eventName)
}

// argName is the name in a parameter arg[<name>], OpenAPI's deepObject style.
func argName(key string) (string, bool) {
	inner, ok := strings.CutPrefix(key, "arg[")
	if !ok || !strings.HasSuffix(inner, "]") || len(inner) < 2 {
		return "", false
	}
	return strings.TrimSuffix(inner, "]"), true
}

// view serves a view of a contract the registry names, at the index's head or at the block named: the inputs by name
// as query parameters, the outputs by name, or the revert by name.
func (s *Server) view(w http.ResponseWriter, r *http.Request) {
	contract, ok := s.Catalog.Contract(r.PathValue("contract"))
	if !ok {
		s.reply(w, nil, missing("the registry names no contract %s", r.PathValue("contract")))
		return
	}
	method, ok := contract.ABI.Methods[r.PathValue("function")]
	if !ok || !method.IsConstant() {
		s.reply(w, nil, missing("%s has no view %s", contract.Name, r.PathValue("function")))
		return
	}
	head, err := s.head(r.Context())
	if err != nil {
		s.reply(w, nil, err)
		return
	}
	block, err := uintParam(r, "block")
	switch {
	case err != nil:
		s.reply(w, nil, err)
		return
	case block > head.Number:
		s.reply(w, nil, invalid("block %d is ahead of the index's head, %d", block, head.Number))
		return
	case block == 0:
		block = head.Number
	}
	query := r.URL.Query()
	args := make([]any, len(method.Inputs))
	for i, input := range method.Inputs {
		name := input.Name
		if name == "" {
			name = "arg" + strconv.Itoa(i)
		}
		text, ok := query[name]
		if !ok {
			s.reply(w, nil, invalid("%s needs %s", method.Sig, name))
			return
		}
		if args[i], err = codec.Parse(input.Type, text[0]); err != nil {
			s.reply(w, nil, invalid("%s: %v", name, err))
			return
		}
		query.Del(name)
	}
	query.Del("block")
	for key := range query {
		s.reply(w, nil, invalid("unknown parameter %s", key))
		return
	}
	key := fmt.Sprintf("view %s %s %d %v", contract.Name, method.Name, block, args)
	out, err := s.cache.get(key, func() (any, error) {
		outputs, revert, err := s.call(r.Context(), contract, method.Name, block, args...)
		return map[string]any{"block": block, "outputs": outputs, "revert": revert}, err
	})
	s.reply(w, out, err)
}

// Revert is a view's revert: its error's name and arguments where an SDK binding knows it, and its data.
type Revert struct {
	Name string        `json:"name,omitempty"`
	Args []any         `json:"args,omitempty"`
	Data hexutil.Bytes `json:"data"`
}

func revertOf(err error) (*Revert, bool) {
	if decoded, ok := sdk.DecodeRevert(err); ok {
		args := make([]any, len(decoded.Args))
		for i, arg := range decoded.Args {
			args[i] = plain(arg)
		}
		return &Revert{decoded.Name, args, decoded.Data}, true
	}
	if data, ok := ethclient.RevertErrorData(err); ok {
		return &Revert{Data: data}, true
	}
	return nil, false
}

// plain is v in the API's JSON conventions, for a value whose ABI type is not at hand.
func plain(v any) any {
	switch v := v.(type) {
	case *big.Int:
		return v.String()
	case common.Address:
		return v.Hex()
	case string, bool:
		return v
	case []byte:
		return hexutil.Encode(v)
	}
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Array:
		raw := make([]byte, value.Len())
		reflect.Copy(reflect.ValueOf(raw), value)
		return hexutil.Encode(raw)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	}
	return fmt.Sprint(v)
}

func valuesOf(arguments abi.Arguments, values []any) map[string]any {
	return codec.Values(arguments, values)
}

func (s *Server) redstonePayload(w http.ResponseWriter, r *http.Request) {
	feeds := r.URL.Query().Get("feeds")
	if feeds == "" {
		s.reply(w, nil, invalid("feeds must name the data packages, separated by commas"))
		return
	}
	signed, err := s.Relay.Signed(r.Context(), strings.Split(feeds, ","))
	if err != nil {
		s.Log.Warn("the RedStone relay failed", "error", err)
		s.fail(w, http.StatusServiceUnavailable, "no RedStone gateway served every package signed")
		return
	}
	s.reply(w, signed, nil)
}
