// SPDX-License-Identifier: MIT OR Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// auctionKey names a liquidation auction: the position of an account.
type auctionKey struct {
	Account  common.Address
	Position [32]byte
}

// liveAuctions reads the auctions started since block from the indexer's public API, oldest first, each once; whether
// each still runs is for the chain to say.
func liveAuctions(ctx context.Context, indexerURL, key string, from uint64) ([]auctionKey, error) {
	query := url.Values{"contract": {"tapehouse.Liquidator"}, "event": {"AuctionStarted"},
		"from": {strconv.FormatUint(from, 10)}, "limit": {"1000"}}
	var started []auctionKey
	seen := map[auctionKey]bool{}
	for {
		var page struct {
			Events []struct {
				Args struct {
					Account  common.Address `json:"account"`
					Position string         `json:"position"`
				} `json:"args"`
			} `json:"events"`
			Next *string `json:"next"`
		}
		if err := get(ctx, indexerURL+"/v1/events?"+query.Encode(), key, &page); err != nil {
			return nil, err
		}
		for _, event := range page.Events {
			raw, err := hexutil.Decode(event.Args.Position)
			if err != nil || len(raw) != 32 {
				return nil, fmt.Errorf("%q is not a 32-byte position", event.Args.Position)
			}
			auction := auctionKey{event.Args.Account, [32]byte(raw)}
			if !seen[auction] {
				seen[auction] = true
				started = append(started, auction)
			}
		}
		if page.Next == nil {
			return started, nil
		}
		query.Set("after", *page.Next)
	}
}

func get(ctx context.Context, target, key string, out any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	if key != "" {
		request.Header.Set("X-API-Key", key)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("the indexer answered %s to /v1/events", response.Status)
	}
	return json.NewDecoder(response.Body).Decode(out)
}
