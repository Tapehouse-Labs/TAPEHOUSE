// SPDX-License-Identifier: MIT OR Apache-2.0

package band

import (
	"math/big"
	"testing"
)

func TestTheRecordFoldsWrittenPricesAsTheBandDoes(t *testing.T) {
	r := NewRecord("NVDA---24_7", "USA500.Y---24_7")
	const friday = 1_789_761_600_000
	r.Write(Written{"NVDA---24_7", big.NewInt(22_000_000_000), friday})
	r.Write(Written{"NVDA---24_7", big.NewInt(22_100_000_000), friday + 49_000})
	r.Write(Written{"NVDA---24_7", big.NewInt(22_200_000_000), friday + 60_000})
	want := EwmaUpdate(new(big.Int), 22_000_000_000, 22_200_000_000, 60)
	if s := r.Samples["NVDA---24_7"]; s.Var.Cmp(want) != 0 || s.Px != 22_200_000_000 || s.Ms != friday+60_000 {
		t.Errorf("the 49 s price was a sample, or the 60 s one was not: %+v", s)
	}
	if r.Prices["NVDA---24_7"].Value.Int64() != 22_200_000_000 {
		t.Error("the last price is not stored")
	}
	status := func(current, next int64, change, at uint64) []Written {
		return []Written{
			{CurrentStatus, big.NewInt(current), at},
			{NextStatus, big.NewInt(next), at},
			{NextChangeTime, new(big.Int).Mul(new(big.Int).SetUint64(change), big.NewInt(statusScale)), at},
		}
	}
	r.Write(status(100_000_000, 102_000_000, friday, friday-60_000)...)
	if r.CloseMs != friday {
		t.Errorf("the regular close was not recorded: %d", r.CloseMs)
	}
	saturday := uint64(friday + 16*3_600_000)
	monday := uint64(friday + 3*86_400_000 - 14_400_000 + 48_600_000)
	r.Write(status(102_000_000, 100_000_000, monday, saturday)...)
	if r.CloseMs != friday {
		t.Errorf("a status signed outside regular hours moved the close: %d", r.CloseMs)
	}
	nowS := saturday/1000 + 30
	one := big.NewInt(1_000_000_000_000_000_000)
	nvda := Asset{Symbol: "NVDA", Feed: "NVDA---24_7"}
	if q := r.Quote(nvda, Chainlink{21_900_000_000, friday / 1000}, one, true, nowS); q.State != Halted {
		t.Errorf("a stale 24/7 leg over the weekend: %+v", q)
	}
	r.Write(Written{"NVDA---24_7", big.NewInt(22_300_000_000), (nowS - 10) * 1000})
	q := r.Quote(nvda, Chainlink{21_900_000_000, friday / 1000}, one, true, nowS)
	if q.State != Closed || q.Live != 1 || q.Mid != 22_300_000_000 {
		t.Errorf("the weekend band: %+v", q)
	}
	if q := r.Quote(nvda, Chainlink{21_900_000_000, friday / 1000}, one, false, nowS); q.State != Degraded {
		t.Errorf("a sequencer down did not degrade the band: %+v", q)
	}
	r.Write(Written{"USA500.Y---24_7", big.NewInt(774_135_000_000), (nowS - 5) * 1000})
	r.Reanchor(Asset{Symbol: "SPY", Index: "USA500.Y---24_7", RegularHours: true}, Chainlink{77_232_802_713, nowS - 5}, nowS)
	if _, ok := r.Anchors["SPY"]; ok {
		t.Error("a 42161 band anchored SPY outside regular hours")
	}
	r.Reanchor(Asset{Symbol: "SPY", Index: "USA500.Y---24_7"}, Chainlink{77_232_802_713, nowS - 5}, nowS)
	if a := r.Anchors["SPY"]; a != (Anchor{77_232_802_713, 774_135_000_000, nowS - 5}) {
		t.Errorf("SPY's anchor: %+v", a)
	}
	r.Anchored("SPY", Anchor{ClPx: 77_232_802_713, IndexPx: 774_263_746_783, AtS: friday / 1000})
	spy := r.Quote(Asset{Symbol: "SPY", Index: "USA500.Y---24_7"}, Chainlink{77_232_802_713, friday / 1000}, one, true, nowS)
	if spy.Mid != 77_219_960_222 || spy.HalfBps != FloorBps+SingleSourceBps+IndexBasisBps {
		t.Errorf("SPY's index leg: %+v", spy)
	}
	night := monday + 13*3_600_000
	r.Write(status(101_000_000, 100_000_000, night+11*3_600_000, night)...)
	r.Write(Written{"NVDA---24_7", big.NewInt(22_300_000_000), night})
	cl := Chainlink{22_310_000_000, night/1000 - 60}
	if q := r.Quote(Asset{Symbol: "NVDA", Feed: "NVDA---24_7", RegularHours: true}, cl, one, true, night/1000); q.Live != 1 || q.State != Degraded {
		t.Errorf("a Chainlink round outside regular hours counted: %+v", q)
	}
	if q := r.Quote(nvda, cl, one, true, night/1000); q.Live != 2 || q.State != Open {
		t.Errorf("a Chainlink round in the overnight session did not count: %+v", q)
	}
}
