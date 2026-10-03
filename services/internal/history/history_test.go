// SPDX-License-Identifier: MIT OR Apache-2.0

package history

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const chart = `{"chart":{"result":[{"meta":{"symbol":"NVDA"},"timestamp":[1757683800,1757943000],
"indicators":{"quote":[{"open":[176.0,null],"close":[177.82,174.88]}],"adjclose":[{"adjclose":[177.8,174.88]}]}}],"error":null}}`

func TestDailyBarsSkipDaysWithoutAPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v8/finance/chart/NVDA" || r.URL.Query().Get("period1") != "1262304000" || r.URL.Query().Get("interval") != "1d" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(chart))
	}))
	defer server.Close()
	s, err := Daily(context.Background(), server.Client(), server.URL, "NVDA", time.Unix(1262304000, 0), time.Unix(1790985600, 0))
	if err != nil {
		t.Fatal(err)
	}
	if s.Symbol != "NVDA" || len(s.Bars) != 1 {
		t.Fatalf("%+v", s)
	}
	b := s.Bars[0]
	if b.At != time.Date(2025, 9, 12, 13, 30, 0, 0, time.UTC) || b.Open != 176.0 || b.Close != 177.82 || b.AdjClose != 177.8 {
		t.Errorf("%+v", b)
	}
	if _, err := Daily(context.Background(), server.Client(), server.URL, "TSLA", time.Unix(1262304000, 0), time.Unix(1790985600, 0)); err == nil {
		t.Error("a failed request was accepted")
	}
}

const redstone = `{"success":true,"data":{"symbol":"NVDA---24_7","onChainUpdates":[
{"value":219.21,"timestamp":1787097659500},{"value":219.0755593,"timestamp":1787094059500}]}}`

func TestRedStoneHistoryIsInEightDecimalsAndAscending(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prices/pull/historical" || r.URL.Query().Get("symbol") != "NVDA---24_7" || r.URL.Query().Get("range") != "7" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(redstone))
	}))
	defer server.Close()
	points, err := RedStone(context.Background(), server.Client(), server.URL, "NVDA---24_7", 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(points) != 2 || points[0].Ms != 1787094059500 || points[0].Scaled().Int64() != 21907555930 || points[1].Scaled().Int64() != 21921000000 {
		t.Errorf("%+v", points)
	}
	if change := (Point{Value: 1788379200000}).Scaled().String(); change != "178837920000000000000" {
		t.Errorf("a change time: %s", change)
	}
	merged := Merge([]Point{{1, 10}, {3, 30}}, []Point{{2, 20}, {3, 31}})
	if len(merged) != 3 || merged[1] != (Point{2, 20}) || merged[2] != (Point{3, 30}) {
		t.Errorf("merge: %+v", merged)
	}
}

func TestTheRoundsAroundAnInstantAreFoundInTheLatestPhase(t *testing.T) {
	phase := new(big.Int).Lsh(big.NewInt(3), 64)
	id := func(n int64) *big.Int { return new(big.Int).Add(phase, big.NewInt(n)) }
	updated := []uint64{0, 100, 200, 200, 300, 400}
	reads := 0
	get := func(round *big.Int) (Round, bool, error) {
		reads++
		n := new(big.Int).Sub(round, phase).Int64()
		if n < 1 || n >= int64(len(updated)) {
			return Round{}, false, nil
		}
		return Round{ID: round, Answer: big.NewInt(n), UpdatedAt: updated[n]}, true, nil
	}
	latest := Round{ID: id(5), Answer: big.NewInt(5), UpdatedAt: 400}
	cases := []struct {
		at            uint64
		before, after int64
	}{
		{50, 0, 1}, {100, 1, 2}, {250, 3, 4}, {399, 4, 5}, {400, 5, 0}, {500, 5, 0},
	}
	for _, c := range cases {
		before, after, err := Around(get, latest, c.at)
		if err != nil {
			t.Fatal(err)
		}
		n := func(r Round) int64 {
			if r.ID == nil {
				return 0
			}
			return new(big.Int).Sub(r.ID, phase).Int64()
		}
		if n(before) != c.before || n(after) != c.after {
			t.Errorf("at %d: rounds %d and %d", c.at, n(before), n(after))
		}
	}
	if reads > 6*4 {
		t.Errorf("%d reads for six searches over five rounds", reads)
	}
}
