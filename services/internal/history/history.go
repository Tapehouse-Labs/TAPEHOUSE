// SPDX-License-Identifier: MIT OR Apache-2.0

// Package history fetches the public record the backtest runs on: daily bars of the launch assets from Yahoo
// Finance's chart API, the RedStone 24/7 and market-status history from RedStone's data explorer, and Chainlink
// rounds from the chain.
package history

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"
)

// Public endpoints the backtest reads.
const (
	YahooBase    = "https://query1.finance.yahoo.com"
	RedStoneBase = "https://o40uhl5zq9.execute-api.us-east-1.amazonaws.com"
)

// Bar is one regular session: when it opened, its open, its close, and its close adjusted for splits and dividends.
type Bar struct {
	At       time.Time
	Open     float64
	Close    float64
	AdjClose float64
}

// Series is an asset's daily bars, ascending.
type Series struct {
	Symbol string
	Bars   []Bar
}

func get(ctx context.Context, client *http.Client, address string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", address, res.Status)
	}
	return io.ReadAll(res.Body)
}

// Daily fetches symbol's daily bars from from to to, with dividends and splits, from Yahoo Finance's chart API at
// base. A day without an open or a close is left out.
func Daily(ctx context.Context, client *http.Client, base, symbol string, from, to time.Time) (Series, error) {
	query := url.Values{
		"period1":  {strconv.FormatInt(from.Unix(), 10)},
		"period2":  {strconv.FormatInt(to.Unix(), 10)},
		"interval": {"1d"},
		"events":   {"div,splits"},
	}
	data, err := get(ctx, client, base+"/v8/finance/chart/"+url.PathEscape(symbol)+"?"+query.Encode())
	if err != nil {
		return Series{}, err
	}
	var chart struct {
		Chart struct {
			Result []struct {
				Meta       struct{ Symbol string }
				Timestamp  []int64
				Indicators struct {
					Quote    []struct{ Open, Close []*float64 }
					Adjclose []struct{ Adjclose []*float64 }
				}
			}
		}
	}
	if err := json.Unmarshal(data, &chart); err != nil {
		return Series{}, err
	}
	if len(chart.Chart.Result) != 1 || len(chart.Chart.Result[0].Indicators.Quote) != 1 || len(chart.Chart.Result[0].Indicators.Adjclose) != 1 {
		return Series{}, fmt.Errorf("%s: no daily series", symbol)
	}
	r := chart.Chart.Result[0]
	q, adj := r.Indicators.Quote[0], r.Indicators.Adjclose[0].Adjclose
	s := Series{Symbol: r.Meta.Symbol}
	for i, ts := range r.Timestamp {
		if i >= len(q.Open) || i >= len(q.Close) || i >= len(adj) || q.Open[i] == nil || q.Close[i] == nil || adj[i] == nil {
			continue
		}
		s.Bars = append(s.Bars, Bar{At: time.Unix(ts, 0).UTC(), Open: *q.Open[i], Close: *q.Close[i], AdjClose: *adj[i]})
	}
	return s, nil
}

// Point is a signed value at its package timestamp in milliseconds.
type Point struct {
	Ms    uint64
	Value float64
}

// Scaled is the value as the band stores it, with 8 decimals, rounded to the nearest integer.
func (p Point) Scaled() *big.Int {
	v := new(big.Float).SetPrec(128).SetFloat64(p.Value)
	v.Mul(v, big.NewFloat(1e8)).Add(v, big.NewFloat(0.5))
	out, _ := v.Int(nil)
	return out
}

// RedStone fetches the last days (1, 7 or 30) of symbol's signed values from RedStone's data explorer at base,
// ascending, at 60 s, 600 s and 3,600 s apart.
func RedStone(ctx context.Context, client *http.Client, base, symbol string, days int) ([]Point, error) {
	query := url.Values{"symbol": {symbol}, "range": {strconv.Itoa(days)}}
	data, err := get(ctx, client, base+"/prices/pull/historical?"+query.Encode())
	if err != nil {
		return nil, err
	}
	var history struct {
		Data struct {
			OnChainUpdates []struct {
				Value     float64
				Timestamp float64
			}
		}
	}
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}
	points := make([]Point, 0, len(history.Data.OnChainUpdates))
	for _, u := range history.Data.OnChainUpdates {
		points = append(points, Point{Ms: uint64(u.Timestamp), Value: u.Value})
	}
	slices.SortFunc(points, func(a, b Point) int { return cmp.Compare(a.Ms, b.Ms) })
	return points, nil
}

// Merge joins histories into one, ascending, the earlier history winning where two share a timestamp.
func Merge(histories ...[]Point) []Point {
	var out []Point
	for _, h := range histories {
		out = append(out, h...)
	}
	slices.SortStableFunc(out, func(a, b Point) int { return cmp.Compare(a.Ms, b.Ms) })
	return slices.CompactFunc(out, func(a, b Point) bool { return a.Ms == b.Ms })
}
