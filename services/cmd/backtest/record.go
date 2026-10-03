// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bandreplay "github.com/tapehouse/tapehouse/services/internal/band"
	"github.com/tapehouse/tapehouse/services/internal/history"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
)

const logWindow = 10_000

func name(id [32]byte) string {
	return string(bytes.TrimRight(id[:], "\x00"))
}

func logs(ctx context.Context, c chain, address common.Address, topics []common.Hash, to uint64) ([]types.Log, error) {
	var out []types.Log
	for from := uint64(0); from <= to; from += logWindow {
		got, err := c.client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(min(from+logWindow-1, to)),
			Addresses: []common.Address{address},
			Topics:    [][]common.Hash{topics},
		})
		if err != nil {
			return nil, err
		}
		out = append(out, got...)
	}
	return out, nil
}

func runRecord(ctx context.Context, url, registry string) error {
	c, err := dial(ctx, url, registry)
	if err != nil {
		return err
	}
	address, ok := c.deployments.Tapehouse["Band"]
	if !ok {
		return fmt.Errorf("%s names no .tapehouse.Band", registry)
	}
	head, err := c.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	binding := band.NewBand()
	contract := bind.NewBoundContract(address, abi.ABI{}, c.client, c.client, c.client)
	events := binding.GetABI().Events
	topics := []common.Hash{events["PriceWritten"].ID, events["Anchored"].ID, events["MultiplierRecorded"].ID}
	found, err := logs(ctx, c, address, topics, head.Number.Uint64())
	if err != nil {
		return err
	}
	for i := range found {
		if m, err := binding.UnpackMultiplierRecordedEvent(&found[i]); err == nil && m.EffectiveAt != 0 {
			if head, err = c.client.HeaderByNumber(ctx, new(big.Int).SetUint64(found[i].BlockNumber-1)); err != nil {
				return err
			}
			fmt.Printf("record: checked at block %v, before the band recorded its first multiplier step\n", head.Number)
			break
		}
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: head.Number}
	config, err := bind.Call(contract, opts, binding.PackChainConfig(), binding.UnpackChainConfig)
	if err != nil {
		return err
	}
	type configured struct {
		asset band.AssetOutput
		quote bandreplay.Asset
	}
	var assets []configured
	var tracked []string
	for _, n := range names {
		symbol, err := sdk.ToBytes32(n)
		if err != nil {
			return err
		}
		a, err := bind.Call(contract, opts, binding.PackAsset(symbol), binding.UnpackAsset)
		if err != nil {
			return err
		}
		if a.ChainlinkFeed == (common.Address{}) && a.RedstoneFeedId == ([32]byte{}) {
			continue
		}
		q := bandreplay.Asset{Symbol: n, Feed: name(a.RedstoneFeedId), Index: name(a.IndexFeedId), RegularHours: config.ChainlinkRegularHours}
		assets = append(assets, configured{a, q})
		for _, f := range []string{q.Feed, q.Index} {
			if f != "" {
				tracked = append(tracked, f)
			}
		}
	}
	record := bandreplay.NewRecord(tracked...)
	var batch []bandreplay.Written
	var tx common.Hash
	flush := func() {
		if len(batch) > 0 {
			record.Write(batch...)
			batch = nil
		}
	}
	written, anchored := 0, 0
	for i := range found {
		l := &found[i]
		if l.BlockNumber > head.Number.Uint64() {
			break
		}
		if l.Topics[0] == events["MultiplierRecorded"].ID {
			continue
		}
		if l.TxHash != tx {
			flush()
			tx = l.TxHash
		}
		if p, err := binding.UnpackPriceWrittenEvent(l); err == nil {
			batch = append(batch, bandreplay.Written{Feed: name(p.FeedId), Value: p.Value, Ms: p.PackageTimestampMs})
			written++
			continue
		}
		a, err := binding.UnpackAnchoredEvent(l)
		if err != nil {
			return err
		}
		flush()
		record.Anchored(name(a.Symbol), bandreplay.Anchor{ClPx: a.ChainlinkPrice, IndexPx: a.IndexPrice, AtS: a.UpdatedAt})
		anchored++
	}
	flush()
	mismatches := 0
	check := func(ok bool, format string, args ...any) {
		status := "matches"
		if !ok {
			status, mismatches = "DIFFERS", mismatches+1
		}
		fmt.Printf("%s: %s\n", fmt.Sprintf(format, args...), status)
	}
	for _, f := range tracked {
		id, err := sdk.ToBytes32(f)
		if err != nil {
			return err
		}
		v, err := bind.Call(contract, opts, binding.PackVariance(id), binding.UnpackVariance)
		if err != nil {
			return err
		}
		check(v.Cmp(record.Samples[f].Var) == 0, "variance of %s, %v on chain", f, v)
	}
	session, err := bind.Call(contract, opts, binding.PackSession(), binding.UnpackSession)
	if err != nil {
		return err
	}
	s, known := record.Session(head.Time)
	code := uint8(0)
	if known {
		code = 1
		if s.Open {
			code = 2
		}
	}
	check(session.State == code && session.Nyse == uint8(s.Nyse) && session.NyseNext == uint8(s.NyseNext) &&
		session.ChangeMs == s.ChangeMs && session.BoundaryMs == s.BoundaryMs, "session %+v", session)
	settled, err := bind.Call(contract, opts, binding.PackSequencerSettled(), binding.UnpackSequencerSettled)
	if err != nil {
		return err
	}
	for _, a := range assets {
		symbol, _ := sdk.ToBytes32(a.quote.Symbol)
		halt, err := bind.Call(contract, opts, binding.PackHalt(symbol), binding.UnpackHalt)
		if err != nil {
			return err
		}
		action, err := bind.Call(contract, opts, binding.PackCorporateAction(symbol), binding.UnpackCorporateAction)
		if err != nil {
			return err
		}
		if halt.SignedHalt || halt.OraclePaused || halt.Until != 0 || action.Status == 2 {
			fmt.Printf("quote of %s: halted or pending a corporate action, not replayed\n", a.quote.Symbol)
			continue
		}
		cl := bandreplay.Chainlink{}
		if a.asset.ChainlinkFeed != (common.Address{}) {
			round, err := history.NewFeed(c.client, a.asset.ChainlinkFeed).Latest(opts)
			if err == nil && round.Answer.Sign() > 0 && round.Answer.IsUint64() {
				cl = bandreplay.Chainlink{Answer: round.Answer.Uint64(), UpdatedAt: round.UpdatedAt}
			}
		}
		multiplier, err := multiplierOf(opts, chain{c.client, &sdk.Deployments{Tokens: map[string]common.Address{a.quote.Symbol: a.asset.Token}}, c.limit}, a.quote.Symbol)
		if err != nil {
			return err
		}
		if a.asset.Token == (common.Address{}) {
			multiplier = one
		}
		got := record.Quote(a.quote, cl, multiplier, settled, head.Time)
		q, err := bind.Call(contract, opts, binding.PackQuote(symbol), binding.UnpackQuote)
		if err != nil {
			return err
		}
		check(q.State == uint8(got.State) && q.Live == got.Live && q.Mid == got.Mid && q.HalfBps == got.HalfBps &&
			q.Low == got.Low && q.High.Cmp(got.High) == 0, "quote of %s, %+v on chain", a.quote.Symbol, q)
	}
	fmt.Printf("record: %d prices and %d anchors written up to block %v\n", written, anchored, head.Number)
	if mismatches > 0 {
		return fmt.Errorf("%d of the band's views differ from its rebuilt record", mismatches)
	}
	return nil
}
