// SPDX-License-Identifier: MIT OR Apache-2.0

package band

import (
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"testing"
)

const (
	quoteVectors   = "../../../stylus/contracts/band/testdata/quote-vectors.json"
	sessionVectors = "../../../stylus/contracts/band/testdata/session-vectors.json"
)

func load(t *testing.T, path string, into any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, into); err != nil {
		t.Fatal(err)
	}
}

func u64(t *testing.T, s string) uint64 {
	t.Helper()
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func wide(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("not an integer: %q", s)
	}
	return v
}

func TestQuotesMatchTheReference(t *testing.T) {
	var v struct {
		Quotes []struct {
			Name  string
			Input struct {
				Live247Px     string `json:"live247_px"`
				Live247AgeS   string `json:"live247_age_s"`
				ClPx          string `json:"cl_px"`
				ClAgeS        string `json:"cl_age_s"`
				ClSessionOpen bool   `json:"cl_session_open"`
				VarCpb2       string `json:"var_cpb2"`
				BasisBps      string `json:"basis_bps"`
			}
			Expected struct {
				State   string
				Live    string
				Mid     string
				HalfBps string `json:"half_bps"`
				Low     string
				High    string
			}
		}
		Ewma []struct {
			Name     string
			Var      string
			PrevPx   string `json:"prev_px"`
			Px       string
			DtS      string `json:"dt_s"`
			Expected string
		}
		IndexLegs []struct {
			Name          string
			AnchorClPx    string `json:"anchor_cl_px"`
			AnchorIndexPx string `json:"anchor_index_px"`
			IndexPx       string `json:"index_px"`
			Expected      string
		} `json:"index_legs"`
		Anchors []struct {
			Name     string
			Anchor   map[string]string
			ClPx     string `json:"cl_px"`
			ClAtS    string `json:"cl_at_s"`
			IndexPx  string `json:"index_px"`
			IndexMs  string `json:"index_ms"`
			Expected map[string]string
		}
	}
	load(t, quoteVectors, &v)
	if len(v.Quotes) != 31 || len(v.Ewma) != 11 || len(v.IndexLegs) != 6 || len(v.Anchors) != 11 {
		t.Fatalf("vectors: %d quotes, %d variance updates, %d legs, %d anchors", len(v.Quotes), len(v.Ewma), len(v.IndexLegs), len(v.Anchors))
	}
	states := map[string]State{"HALTED": Halted, "DEGRADED": Degraded, "CLOSED": Closed, "OPEN": Open}
	for _, c := range v.Quotes {
		in := Inputs{
			Live247Px:     u64(t, c.Input.Live247Px),
			Live247AgeS:   u64(t, c.Input.Live247AgeS),
			ClPx:          u64(t, c.Input.ClPx),
			ClAgeS:        u64(t, c.Input.ClAgeS),
			ClSessionOpen: c.Input.ClSessionOpen,
			VarCpb2:       wide(t, c.Input.VarCpb2),
		}
		if c.Input.BasisBps != "" {
			in.BasisBps = u64(t, c.Input.BasisBps)
		}
		got := Compute(in)
		e := c.Expected
		if got.State != states[e.State] || uint64(got.Live) != u64(t, e.Live) || got.Mid != u64(t, e.Mid) ||
			got.HalfBps != u64(t, e.HalfBps) || got.Low != u64(t, e.Low) || got.High.Cmp(wide(t, e.High)) != 0 {
			t.Errorf("%s: got %+v, want %+v", c.Name, got, e)
		}
	}
	for _, c := range v.Ewma {
		got := EwmaUpdate(wide(t, c.Var), u64(t, c.PrevPx), u64(t, c.Px), u64(t, c.DtS))
		if got.Cmp(wide(t, c.Expected)) != 0 {
			t.Errorf("%s: got %v, want %s", c.Name, got, c.Expected)
		}
	}
	for _, c := range v.IndexLegs {
		anchor := Anchor{ClPx: u64(t, c.AnchorClPx), IndexPx: u64(t, c.AnchorIndexPx)}
		if got := LegPx(anchor, u64(t, c.IndexPx)); got != u64(t, c.Expected) {
			t.Errorf("%s: got %d, want %s", c.Name, got, c.Expected)
		}
	}
	anchor := func(m map[string]string) Anchor {
		return Anchor{ClPx: u64(t, m["cl_px"]), IndexPx: u64(t, m["index_px"]), AtS: u64(t, m["at_s"])}
	}
	for _, c := range v.Anchors {
		before := anchor(c.Anchor)
		after, ok := NextAnchor(before, u64(t, c.ClPx), u64(t, c.ClAtS), u64(t, c.IndexPx), u64(t, c.IndexMs))
		if !ok {
			after = before
		}
		if after != anchor(c.Expected) {
			t.Errorf("%s: got %+v, want %+v", c.Name, after, c.Expected)
		}
	}
}

func TestSamplesFollowTheFiftySecondRule(t *testing.T) {
	last := Sample{Var: big.NewInt(250_000), Px: 22_000_000_000, Ms: 1_790_000_000_000}
	if _, ok := NextSample(last, 22_300_000_000, 1_790_000_049_999); ok {
		t.Error("a price within 50 s of the last sample is a sample")
	}
	next, ok := NextSample(last, 22_300_000_000, 1_790_000_050_999)
	if !ok || next.Var.Cmp(EwmaUpdate(big.NewInt(250_000), 22_000_000_000, 22_300_000_000, 50)) != 0 {
		t.Errorf("a price 50 s after the last sample: %+v, %v", next, ok)
	}
	long, ok := NextSample(last, 22_000_000_000, 1_790_003_600_000)
	if !ok || long.Var.Sign() <= 0 {
		t.Error("a long gap reset the variance")
	}
	if _, ok := NextSample(last, 0, 1_790_000_060_000); ok {
		t.Error("a zero price is a sample")
	}
	first, ok := NextSample(Sample{Var: new(big.Int)}, 22_000_000_000, 1_790_000_000_000)
	if !ok || first.Var.Sign() != 0 || first.Px != 22_000_000_000 {
		t.Errorf("the first price: %+v", first)
	}
	if PackageAgeS(1_790_000_000, 1_790_000_030_000) != 0 || PackageAgeS(1_790_000_000, 1_789_999_879_000) != 121 {
		t.Error("package ages")
	}
}

func TestTokenPricesScaleTheSharePriceDown(t *testing.T) {
	if got := TokenPx(76_024_470_000, big.NewInt(1_001_717_991_187_472_003)); got != 76_155_079_369 {
		t.Errorf("SPY's step of 18 Sep 2026: got %d", got)
	}
	if got := TokenPx(1<<63, new(big.Int).Lsh(big.NewInt(1), 70)); got != 0 {
		t.Errorf("a price beyond u64: got %d", got)
	}
}

func TestSessionsMatchTheReference(t *testing.T) {
	var v struct {
		Sessions []struct {
			Name      string
			Values    []string
			PackageMs []string `json:"package_ms"`
			CloseMs   string   `json:"close_ms"`
			NowMs     string   `json:"now_ms"`
			Expected  struct {
				Known      bool
				Open       bool
				Nyse       string
				NyseNext   string `json:"nyse_next"`
				ChangeMs   string `json:"change_ms"`
				BoundaryMs string `json:"boundary_ms"`
			}
		}
		Closes []struct {
			Name     string
			Values   []string
			CloseMs  string `json:"close_ms"`
			Expected string
		}
	}
	load(t, sessionVectors, &v)
	if len(v.Sessions) != 38 || len(v.Closes) != 3 {
		t.Fatalf("vectors: %d sessions, %d closes", len(v.Sessions), len(v.Closes))
	}
	values := func(s []string) [3]*big.Int { return [3]*big.Int{wide(t, s[0]), wide(t, s[1]), wide(t, s[2])} }
	for _, c := range v.Sessions {
		status, known := DecodeStatus(values(c.Values), [3]uint64{u64(t, c.PackageMs[0]), u64(t, c.PackageMs[1]), u64(t, c.PackageMs[2])})
		s, ok := Derive(status, known, u64(t, c.CloseMs), u64(t, c.NowMs))
		e := c.Expected
		if ok != e.Known {
			t.Errorf("%s: known %v", c.Name, ok)
			continue
		}
		if ok && (s.Open != e.Open || uint64(s.Nyse) != u64(t, e.Nyse) || uint64(s.NyseNext) != u64(t, e.NyseNext) ||
			s.ChangeMs != u64(t, e.ChangeMs) || s.BoundaryMs != u64(t, e.BoundaryMs)) {
			t.Errorf("%s: got %+v, want %+v", c.Name, s, e)
		}
	}
	for _, c := range v.Closes {
		status, known := DecodeStatus(values(c.Values), [3]uint64{})
		if got := NextClose(status, known, u64(t, c.CloseMs)); got != u64(t, c.Expected) {
			t.Errorf("%s: got %d, want %s", c.Name, got, c.Expected)
		}
	}
}

func TestTheSequencerSettlesAnHourAfterItComesBack(t *testing.T) {
	cases := []struct {
		answer, startedAt, now uint64
		settled                bool
	}{
		{0, 1_000, 4_601, true},
		{0, 1_000, 4_600, false},
		{1, 1_000, 9_000, false},
		{0, 0, 9_000, false},
		{0, 9_001, 9_000, false},
	}
	for _, c := range cases {
		if got := SequencerSettled(c.answer, c.startedAt, c.now); got != c.settled {
			t.Errorf("%+v: got %v", c, got)
		}
	}
}
