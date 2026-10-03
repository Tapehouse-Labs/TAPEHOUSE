// SPDX-License-Identifier: MIT OR Apache-2.0

package history

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/aggregator"
)

// Round is a Chainlink round: its proxy round ID, answer, and when it started and was updated, in seconds.
type Round struct {
	ID        *big.Int
	Answer    *big.Int
	StartedAt uint64
	UpdatedAt uint64
}

// Feed reads a Chainlink feed's rounds through AggregatorV3Interface.
type Feed struct {
	contract *bind.BoundContract
	binding  *aggregator.Aggregator
}

// NewFeed reads the feed at address over backend. The address comes from the SDK's TokenPriceFeed or
// SharePriceFeed, or the registry's sequencer-uptime feed.
func NewFeed(backend bind.ContractBackend, address common.Address) *Feed {
	return &Feed{bind.NewBoundContract(address, abi.ABI{}, backend, backend, backend), aggregator.NewAggregator()}
}

// Latest reads the feed's latest round.
func (f *Feed) Latest(opts *bind.CallOpts) (Round, error) {
	out, err := bind.Call(f.contract, opts, f.binding.PackLatestRoundData(), f.binding.UnpackLatestRoundData)
	if err != nil {
		return Round{}, err
	}
	return Round{out.RoundId, out.Answer, out.StartedAt.Uint64(), out.UpdatedAt.Uint64()}, nil
}

// Get reads round id, and reports whether it was ever updated.
func (f *Feed) Get(opts *bind.CallOpts, id *big.Int) (Round, bool, error) {
	out, err := bind.Call(f.contract, opts, f.binding.PackGetRoundData(id), f.binding.UnpackGetRoundData)
	if err != nil || out.UpdatedAt.Sign() == 0 {
		return Round{}, false, err
	}
	return Round{out.RoundId, out.Answer, out.StartedAt.Uint64(), out.UpdatedAt.Uint64()}, true, nil
}

// Around finds, within latest's phase, the last round updated at or before at and the first after it, by binary
// search over round IDs; a zero Round where there is none.
func Around(get func(*big.Int) (Round, bool, error), latest Round, at uint64) (before, after Round, err error) {
	phase := new(big.Int).Rsh(latest.ID, 64)
	phase.Lsh(phase, 64)
	if latest.UpdatedAt <= at {
		return latest, Round{}, nil
	}
	read := func(n int64) (Round, error) {
		r, _, err := get(new(big.Int).Add(phase, big.NewInt(n)))
		return r, err
	}
	seen := map[int64]Round{new(big.Int).Sub(latest.ID, phase).Int64(): latest}
	lo, hi := int64(1), new(big.Int).Sub(latest.ID, phase).Int64()
	for lo < hi {
		mid := lo + (hi-lo)/2
		r, err := read(mid)
		if err != nil {
			return Round{}, Round{}, err
		}
		seen[mid] = r
		if r.UpdatedAt > at {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	after, ok := seen[lo]
	if !ok {
		if after, err = read(lo); err != nil {
			return Round{}, Round{}, err
		}
	}
	if lo > 1 {
		if before, ok = seen[lo-1]; !ok {
			before, err = read(lo - 1)
		}
	}
	return before, after, err
}
