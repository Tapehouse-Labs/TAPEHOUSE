// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/store"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// gapCoverSeries serves a series of the gap cover with what rebuilds its pricing and its settlement: the series at the
// head, the week's move measured for it, the covers bought on it, every band its window recorded, its settlement with
// both rounds, its reference and its price, or its voiding.
func (s *Server) gapCoverSeries(w http.ResponseWriter, r *http.Request) {
	asset := r.PathValue("symbol")
	symbol, err := sdk.ToBytes32(asset)
	if err != nil {
		s.reply(w, nil, invalid("%v", err))
		return
	}
	closesMs, err := strconv.ParseUint(r.PathValue("closesMs"), 10, 64)
	if err != nil {
		s.reply(w, nil, invalid("closesMs must be milliseconds"))
		return
	}
	s.atHead(w, r, fmt.Sprintf("gapCoverSeries %s %d", asset, closesMs), func(ctx context.Context, block uint64) (map[string]any, error) {
		cover, err := s.contract("tapehouse.GapCover")
		if err != nil {
			return nil, err
		}
		series, err := s.sdk.GapCover().Series(at(ctx, block), asset, closesMs)
		if err != nil {
			return nil, err
		}
		key := map[string]any{"symbol": hexutil.Encode(symbol[:]), "closesMs": strconv.FormatUint(closesMs, 10)}
		out := map[string]any{"symbol": asset, "closesMs": closesMs, "series": map[string]any{
			"notional": series.Notional.String(), "status": [...]string{"open", "settled", "void"}[series.Status],
			"referencePrice": strconv.FormatUint(series.ReferencePrice, 10), "price": strconv.FormatUint(series.Price, 10),
			"flagged": series.Flagged, "reopenMs": series.ReopenMs.String(),
		}}
		for name, event := range map[string]string{"measured": "Measured", "bought": "Bought", "observed": "Observed",
			"settled": "Settled", "voided": "Voided"} {
			events, err := s.all(ctx, store.Query{Contract: cover.Name, Event: event, Args: key, To: block})
			if err != nil {
				return nil, err
			}
			out[name] = s.served(events)
		}
		return out, nil
	})
}

// all returns every event q selects, page by page.
func (s *Server) all(ctx context.Context, q store.Query) ([]store.Event, error) {
	var out []store.Event
	q.Limit = MaxPage
	for {
		page, err := s.Store.Events(ctx, q)
		if err != nil {
			return nil, err
		}
		out = append(out, page...)
		if len(page) < q.Limit {
			return out, nil
		}
		q.After = &store.Cursor{Block: page[len(page)-1].Block, LogIndex: page[len(page)-1].LogIndex}
	}
}
