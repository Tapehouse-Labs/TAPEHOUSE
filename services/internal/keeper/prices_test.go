// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type source struct{ asked [][]string }

func (s *source) Payload(_ context.Context, ids []string) ([]byte, error) {
	s.asked = append(s.asked, ids)
	return []byte(strings.Join(ids, ",")), nil
}

// pricedWorld is a world whose band prices NVDA on its 24/7 feed, TSLA with no Stock Token, and SPY on the S&P 500 index.
func pricedWorld(t *testing.T) *world {
	w := newWorld(t)
	for _, name := range []string{"NVDA", "TSLA", "SPY"} {
		w.asset(name)
	}
	w.read("Band", "asset", func(args []any) []any {
		feeds := map[string][2]string{"NVDA": {"NVDA---24_7", ""}, "TSLA": {"TSLA---24_7", ""}, "SPY": {"", "USA500.Y---24_7"}}
		name := text(args[0].([32]byte))
		f, ok := feeds[name]
		if !ok {
			return []any{common.Address{}, [32]byte{}, [32]byte{}, common.Address{}}
		}
		token := w.deployments.Tokens[name]
		if name == "TSLA" {
			token = common.Address{}
		}
		return []any{w.deployments.Chainlink[name+"_USD"], symbol(t, f[0]), symbol(t, f[1]), token}
	})
	w.read("Band", "price", func([]any) []any {
		return []any{n(101_000_000), uint64(w.time*1000 - 20_000), w.time - 15}
	})
	return w
}

func TestOneWriteCarriesTheStatusThenTheIndexThenEvery247FeedAndTheStatusIsOnePackage(t *testing.T) {
	w := pricedWorld(t)
	var written [][32]byte
	w.write("Band", "writePrices", func(args []any, mined bool) error {
		if mined {
			written = args[0].([][32]byte)
		}
		return nil
	})
	k := w.keeper(Config{})
	feeds := &source{}
	k.prices = statusGroup{feeds}
	if err := k.Pass(context.Background(), "prices"); err != nil {
		t.Fatal(err)
	}
	want := []string{"NY_MARKET_CURRENT_STATUS", "NY_MARKET_NEXT_STATUS", "NY_MARKET_NEXT_CHANGE_TIME", "USA500.Y---24_7",
		"NVDA---24_7", "TSLA---24_7"}
	got := make([]string, len(written))
	for i, id := range written {
		got[i] = text(id)
	}
	if !slices.Equal(got, want) || len(w.mined()) != 1 {
		t.Fatalf("wrote %v", got)
	}
	if !slices.Equal(feeds.asked[0], []string{"NY_MARKET_STATUS", "USA500.Y---24_7", "NVDA---24_7", "TSLA---24_7"}) {
		t.Fatalf("asked RedStone for %v", feeds.asked)
	}
}

func TestAWriteAFrontRunnerBeatTakesTheNextPackageAndAnOldStatusIsAlerted(t *testing.T) {
	w := pricedWorld(t)
	beaten := 2
	w.write("Band", "writePrices", func(_ []any, mined bool) error {
		if beaten > 0 && !mined {
			beaten--
			return w.fail("Band", "PackageNotNewer", symbol(t, "NVDA---24_7"), uint64(1790775000000), uint64(1790775000000))
		}
		return nil
	})
	k := w.keeper(Config{})
	k.prices = &source{}
	var slept []time.Duration
	k.sleep = func(_ context.Context, d time.Duration) error {
		slept = append(slept, d)
		return nil
	}
	if err := k.Pass(context.Background(), "prices"); err != nil {
		t.Fatal(err)
	}
	if len(w.mined()) != 1 || !slices.Equal(slept, []time.Duration{PackageInterval, PackageInterval}) {
		t.Fatalf("slept %v", slept)
	}
	if strings.Contains(w.logs.String(), "market status is old") {
		t.Fatal("a fresh status was alerted")
	}
	w.read("Band", "price", func([]any) []any {
		return []any{n(101_000_000), uint64(w.time*1000) - uint64(StatusAlertAge.Milliseconds()), w.time - 900}
	})
	beaten = PackageRetries + 1
	if err := k.Pass(context.Background(), "prices"); err == nil || !strings.Contains(err.Error(), "PackageNotNewer") {
		t.Fatalf("a write beaten every time: %v", err)
	}
	if !strings.Contains(w.logs.String(), `level=ERROR msg="the band's market status is old`) {
		t.Fatalf("logs: %s", w.logs.String())
	}
}

func TestTheMultiplierIsSyncedAtStartOnANewTermAndOnEachRoundWhileAStepWaits(t *testing.T) {
	w := pricedWorld(t)
	newMultiplier, effective, paused := n(1e18), n(0), false
	for _, name := range []string{"NVDA", "SPY"} {
		w.read("Token:"+name, "uiMultiplier", func([]any) []any { return []any{n(1e18)} })
		w.read("Token:"+name, "newUIMultiplier", func([]any) []any { return []any{newMultiplier} })
		w.read("Token:"+name, "effectiveAt", func([]any) []any { return []any{effective} })
		w.read("Token:"+name, "oraclePaused", func([]any) []any { return []any{paused} })
	}
	status, round := uint8(0), int64(7)
	w.read("Band", "corporateAction", func([]any) []any { return []any{status, uint64(0), n(0), n(0)} })
	w.read("Feed:NVDA", "latestRoundData", func([]any) []any { return []any{n(round), n(1), n(1), n(1), n(round)} })
	w.read("Feed:SPY", "latestRoundData", func([]any) []any { return []any{n(round), n(1), n(1), n(1), n(round)} })
	w.write("Band", "syncMultiplier", func([]any, bool) error { return nil })
	w.write("Band", "writePrices", func([]any, bool) error { return nil })
	k := w.keeper(Config{})
	k.prices = &source{}
	sync := func(want ...string) {
		t.Helper()
		if err := k.Pass(context.Background(), "multiplier"); err != nil {
			t.Fatal(err)
		}
		if got := w.mined(); fmt.Sprint(got) != fmt.Sprint(want) {
			t.Fatalf("mined %v, want %v", got, want)
		}
	}
	nvda, spy := fmt.Sprintf("Band.syncMultiplier[%v]", symbol(t, "NVDA")), fmt.Sprintf("Band.syncMultiplier[%v]", symbol(t, "SPY"))
	sync(nvda, spy)
	sync()
	newMultiplier, effective = n(4e18), n(int64(w.time)+600)
	sync(nvda, spy)
	sync()
	status = 2
	if err := k.Pass(context.Background(), "multiplier"); err != nil {
		t.Fatal(err)
	}
	if got := w.mined(); len(got) != 4 || !strings.HasPrefix(got[0], "Band.writePrices") || got[1] != nvda {
		t.Fatalf("a waiting step's first pass: %v", got)
	}
	sync()
	round = 8
	paused = true
	sync()
	paused = false
	if err := k.Pass(context.Background(), "multiplier"); err != nil {
		t.Fatal(err)
	}
	if got := w.mined(); len(got) != 4 || got[1] != nvda || got[3] != spy {
		t.Fatalf("after the pause ended: %v", got)
	}
}

func TestAHaltIsSignedRenewedAndLiftedAndOneRelayedFirstIsWritten(t *testing.T) {
	w := pricedWorld(t)
	halted := map[string]bool{"NVDA": true, "TSLA": false, "SPY": false}
	quotes := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		var out []map[string]any
		for name, h := range halted {
			out = append(out, map[string]any{"tokenSymbol": name, "isTradingHalt": h, "bid": "1"})
		}
		_ = json.NewEncoder(rw).Encode(map[string]any{"quotes": out})
	}))
	defer quotes.Close()
	type stored struct {
		halted          bool
		until, issuedAt uint64
	}
	store := map[string]stored{}
	relayed := false
	signer, _ := crypto.GenerateKey()
	w.read("Band", "halt", func(args []any) []any {
		s := store[text(args[0].([32]byte))]
		return []any{s.halted && s.until > w.time, s.until, s.issuedAt, false}
	})
	w.write("Band", "writeHalt", func(args []any, mined bool) error {
		asset, isHalted, issued, expires, signature := text(args[0].([32]byte)), args[1].(bool), args[2].(uint64), args[3].(uint64), args[4].([]byte)
		digest, err := HaltDigest(w.deployments.ChainID, w.deployments.Tapehouse["Band"], asset, isHalted, issued, expires)
		if err != nil {
			t.Fatal(err)
		}
		signature = append([]byte{}, signature...)
		signature[64] -= 27
		key, err := crypto.SigToPub(digest, signature)
		if err != nil || crypto.PubkeyToAddress(*key) != crypto.PubkeyToAddress(signer.PublicKey) {
			t.Fatalf("a halt signed by another key: %v", err)
		}
		if relayed {
			return w.fail("Band", "HaltNotNewer", args[0], issued, issued)
		}
		if mined {
			if isHalted {
				store[asset] = stored{true, expires, issued}
			} else {
				store[asset] = stored{false, 0, issued}
			}
		}
		return nil
	})
	haltSigner := crypto.PubkeyToAddress(signer.PublicKey)
	w.read("Band", "haltSigner", func([]any) []any { return []any{haltSigner} })
	k := w.keeper(Config{HaltsURL: quotes.URL})
	k.halts = &halts{key: signer, url: quotes.URL, client: http.DefaultClient}
	haltSigner = alice
	if err := k.Pass(context.Background(), "halt"); err == nil || !strings.Contains(err.Error(), "not the band's halt signer") {
		t.Fatalf("a key that is not the band's halt signer: %v", err)
	}
	haltSigner = crypto.PubkeyToAddress(signer.PublicKey)
	pass := func() []string {
		t.Helper()
		if err := k.Pass(context.Background(), "halt"); err != nil {
			t.Fatal(err)
		}
		return w.mined()
	}
	if got := pass(); len(got) != 1 || store["NVDA"] != (stored{true, w.time + 900, w.time}) {
		t.Fatalf("a halt: %v %v", got, store)
	}
	w.time += 299
	if got := pass(); len(got) != 0 {
		t.Fatalf("renewed early: %v", got)
	}
	w.time++
	if got := pass(); len(got) != 1 || store["NVDA"].issuedAt != w.time {
		t.Fatalf("not renewed after five minutes: %v %v", got, store)
	}
	halted["NVDA"] = false
	w.time++
	if got := pass(); len(got) != 1 || store["NVDA"] != (stored{false, 0, w.time}) {
		t.Fatalf("not lifted: %v %v", got, store)
	}
	if got := pass(); len(got) != 0 {
		t.Fatalf("lifted twice: %v", got)
	}
	store["SPY"] = stored{true, w.time - 1, w.time - 901}
	relayed = true
	w.time++
	if got := pass(); len(got) != 0 {
		t.Fatalf("a lift relayed first was sent again: %v", got)
	}
	delete(halted, "SPY")
	relayed = false
	if got := pass(); len(got) != 0 || !strings.Contains(w.logs.String(), "Robinhood quotes no such Stock Token") {
		t.Fatalf("an asset Robinhood does not quote: %v", got)
	}
}

func TestTheHaltSignatureIsTheBandsVectorsByteForByte(t *testing.T) {
	data, err := os.ReadFile("../../../stylus/contracts/band/testdata/halt-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors struct {
		Signed []struct {
			Name      string `json:"name"`
			ChainID   string `json:"chain_id"`
			Band      string `json:"band"`
			Symbol    string `json:"symbol"`
			Halted    bool   `json:"halted"`
			IssuedAt  string `json:"issued_at"`
			ExpiresAt string `json:"expires_at"`
			Hash      string `json:"hash"`
			Signature string `json:"signature"`
			Signer    string `json:"signer"`
		} `json:"signed"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	anvil, _ := crypto.HexToECDSA("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	checked := 0
	for _, v := range vectors.Signed {
		if common.HexToAddress(v.Signer) != crypto.PubkeyToAddress(anvil.PublicKey) {
			continue
		}
		checked++
		chainID, _ := strconv.ParseUint(v.ChainID, 10, 64)
		issued, _ := strconv.ParseUint(v.IssuedAt, 10, 64)
		expires, _ := strconv.ParseUint(v.ExpiresAt, 10, 64)
		asset := text([32]byte(common.FromHex(v.Symbol)))
		digest, err := HaltDigest(chainID, common.HexToAddress(v.Band), asset, v.Halted, issued, expires)
		if err != nil || hexutil.Encode(digest) != v.Hash {
			t.Errorf("%s: digest %x, %v", v.Name, digest, err)
		}
		signature, err := SignHalt(anvil, chainID, common.HexToAddress(v.Band), asset, v.Halted, issued, expires)
		if err != nil || hexutil.Encode(signature) != v.Signature {
			t.Errorf("%s: signature %x, %v", v.Name, signature, err)
		}
	}
	if checked < 2 {
		t.Fatalf("%d vectors signed by the dev node's halt signer", checked)
	}
}
