// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Committed is a bid the bot committed, or is about to: everything its reveal needs.
type Committed struct {
	Symbol                   string
	OpenMs                   uint64
	Quantity, Price, Deposit *big.Int
	Salt, Commitment         [32]byte
}

type record struct {
	Symbol     string `json:"symbol"`
	OpenMs     uint64 `json:"openMs"`
	Quantity   string `json:"quantity"`
	Price      string `json:"price"`
	Deposit    string `json:"deposit"`
	Salt       string `json:"salt"`
	Commitment string `json:"commitment"`
}

// State is the bot's committed bids on disk. A bid is written, and synced, before its commit is sent, so a restart
// reveals what a crash interrupted.
type State struct {
	path           string
	mu             sync.Mutex
	bids           []Committed
	afterTemporary func() error
}

// LoadState reads the state file at path; a missing file is an empty state.
func LoadState(path string) (*State, error) {
	state := &State{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return nil, err
	}
	var file struct {
		Bids []record `json:"bids"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("the state file %s is not valid: %w", path, err)
	}
	for _, r := range file.Bids {
		bid, err := r.committed()
		if err != nil {
			return nil, fmt.Errorf("the state file %s is not valid: %w", path, err)
		}
		state.bids = append(state.bids, bid)
	}
	return state, nil
}

// Bids returns the committed bids, oldest first.
func (s *State) Bids() []Committed {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.bids)
}

// Add records bid, one for each asset and round, and returns once it is on disk.
func (s *State) Add(bid Committed) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, held := range s.bids {
		if held.Symbol == bid.Symbol && held.OpenMs == bid.OpenMs {
			return fmt.Errorf("a bid in %s's round at %d is already held", bid.Symbol, bid.OpenMs)
		}
	}
	return s.write(append(slices.Clone(s.bids), bid))
}

// Remove forgets asset's round at openMs.
func (s *State) Remove(symbol string, openMs uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.write(slices.DeleteFunc(slices.Clone(s.bids), func(bid Committed) bool {
		return bid.Symbol == symbol && bid.OpenMs == openMs
	}))
}

// write puts bids in a temporary file, syncs it and renames it over the state file.
func (s *State) write(bids []Committed) error {
	records := make([]record, len(bids))
	for i, bid := range bids {
		records[i] = record{bid.Symbol, bid.OpenMs, bid.Quantity.String(), bid.Price.String(), bid.Deposit.String(),
			hexutil.Encode(bid.Salt[:]), hexutil.Encode(bid.Commitment[:])}
	}
	data, err := json.MarshalIndent(map[string]any{"bids": records}, "", "  ")
	if err != nil {
		return err
	}
	temporary := s.path + ".tmp"
	file, err := os.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err == nil {
		err = file.Sync()
	}
	if closed := file.Close(); err == nil {
		err = closed
	}
	if err == nil && s.afterTemporary != nil {
		err = s.afterTemporary()
	}
	if err == nil {
		err = os.Rename(temporary, s.path)
	}
	if err != nil {
		return err
	}
	if dir, err := os.Open(filepath.Dir(s.path)); err == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	s.bids = bids
	return nil
}

func (r record) committed() (Committed, error) {
	bid := Committed{Symbol: r.Symbol, OpenMs: r.OpenMs}
	var ok [3]bool
	bid.Quantity, ok[0] = new(big.Int).SetString(r.Quantity, 10)
	bid.Price, ok[1] = new(big.Int).SetString(r.Price, 10)
	bid.Deposit, ok[2] = new(big.Int).SetString(r.Deposit, 10)
	if !ok[0] || !ok[1] || !ok[2] {
		return Committed{}, fmt.Errorf("the bid in %s's round at %d has an amount that is not a number", r.Symbol, r.OpenMs)
	}
	for text, into := range map[string]*[32]byte{r.Salt: &bid.Salt, r.Commitment: &bid.Commitment} {
		raw, err := hexutil.Decode(text)
		if err != nil || len(raw) != 32 {
			return Committed{}, fmt.Errorf("the bid in %s's round at %d has a word that is not 32 bytes", r.Symbol, r.OpenMs)
		}
		*into = [32]byte(raw)
	}
	return bid, nil
}
