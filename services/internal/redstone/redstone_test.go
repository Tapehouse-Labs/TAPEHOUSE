// SPDX-License-Identifier: MIT OR Apache-2.0

package redstone_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/sdk"
)

var _ sdk.PackageSource = (*redstone.Relay)(nil)

// latest is a keyed gateway's answer for NVDA---24_7, NY_MARKET_STATUS and NY_MARKET_CURRENT_STATUS, as RedStone
// signed it on 2 October 2026.
func latest(t *testing.T) map[string][]map[string]any {
	t.Helper()
	raw, err := os.ReadFile("testdata/latest.json")
	if err != nil {
		t.Fatal(err)
	}
	var out map[string][]map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// gateway serves body for the path a keyed gateway or a public one answers on, counting its requests.
type gateway struct {
	*httptest.Server
	requests atomic.Int64
	key      string
}

func serve(t *testing.T, status int, body any, key string) *gateway {
	t.Helper()
	g := &gateway{key: key}
	g.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.requests.Add(1)
		if key != "" && (r.Header.Get("x-api-key") != key ||
			r.URL.Path != "/v2/data-packages/latest-by-data-feeds/redstone-primary-prod") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if key == "" && r.URL.Path != "/data-packages/latest/redstone-primary-prod" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(g.Close)
	return g
}

func TestThePayloadIsRedStonesSerialisationOfVerifiedPackages(t *testing.T) {
	keyed := serve(t, http.StatusOK, latest(t), "key")
	relay := redstone.New([]redstone.Gateway{{URL: keyed.URL, Key: "key"}}, keyed.Client())
	signed, err := relay.Signed(context.Background(), []string{"NVDA---24_7", "NY_MARKET_STATUS"})
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/nvda-status.hex")
	if err != nil {
		t.Fatal(err)
	}
	if hexutil.Encode(signed.Payload) != strings.TrimSpace(string(want)) {
		t.Fatal("the payload differs from redstone-payload.py's")
	}
	if signed.TimestampMs != 1_790_949_660_000 || len(signed.Packages) != 10 {
		t.Fatalf("%d packages at %d", len(signed.Packages), signed.TimestampMs)
	}
	nvda, status := signed.Packages[0], signed.Packages[5]
	if common.Address(nvda.Signer) != redstone.Signers[0] || nvda.DataPoints[0].Value != "23729678801" {
		t.Fatalf("NVDA---24_7 %+v", nvda)
	}
	if status.DataPackageID != "NY_MARKET_STATUS" || len(status.DataPoints) != 3 ||
		status.DataPoints[1].DataFeedID != "NY_MARKET_NEXT_CHANGE_TIME" ||
		status.DataPoints[1].Value != "179097120000000000000" || status.DataPoints[2].Value != "102000000" {
		t.Fatalf("NY_MARKET_STATUS %+v", status)
	}
	payload, err := relay.Payload(context.Background(), []string{"NVDA---24_7", "NY_MARKET_STATUS"})
	if err != nil || hexutil.Encode(payload) != strings.TrimSpace(string(want)) {
		t.Fatalf("as a PackageSource: %v", err)
	}
	if keyed.requests.Load() != 1 {
		t.Fatalf("%d requests in MaxAge", keyed.requests.Load())
	}
}

func TestAPackageThatDoesNotRecoverToASignerIsNotServed(t *testing.T) {
	body := latest(t)
	for _, p := range body["NVDA---24_7"][:3] {
		p["dataPoints"].([]any)[0].(map[string]any)["value"] = 1
	}
	keyed := serve(t, http.StatusOK, body, "key")
	relay := redstone.New([]redstone.Gateway{{URL: keyed.URL, Key: "key"}}, keyed.Client())
	if _, err := relay.Signed(context.Background(), []string{"NVDA---24_7"}); err == nil ||
		!strings.Contains(err.Error(), "NVDA---24_7 has 2 of the 3 signatures the band requires") {
		t.Fatalf("tampered packages: %v", err)
	}
	signed, err := relay.Signed(context.Background(), []string{"NY_MARKET_CURRENT_STATUS"})
	if err != nil || len(signed.Packages) != 5 {
		t.Fatalf("an untouched package: %v", err)
	}
}

func TestPackagesOfDifferentTimesOrValuesTooPreciseAreRefused(t *testing.T) {
	body := latest(t)
	raw, err := os.ReadFile("testdata/earlier.json")
	if err != nil {
		t.Fatal(err)
	}
	var earlier map[string][]map[string]any
	if err := json.Unmarshal(raw, &earlier); err != nil {
		t.Fatal(err)
	}
	body["NY_MARKET_STATUS"] = earlier["NY_MARKET_STATUS"]
	body["NY_MARKET_CURRENT_STATUS"][0]["timestampMilliseconds"] = 1_790_949_650_000
	g := serve(t, http.StatusOK, body, "key")
	relay := redstone.New([]redstone.Gateway{{URL: g.URL, Key: "key"}}, g.Client())
	if _, err := relay.Signed(context.Background(), []string{"NVDA---24_7", "NY_MARKET_STATUS"}); err == nil ||
		!strings.Contains(err.Error(), "different timestamps") {
		t.Fatalf("two timestamps: %v", err)
	}
	signed, err := relay.Signed(context.Background(), []string{"NY_MARKET_CURRENT_STATUS"})
	if err != nil || len(signed.Packages) != 4 {
		t.Fatalf("a package whose signed timestamp was changed: %v", err)
	}
	body = latest(t)
	body["NVDA---24_7"][0]["dataPoints"].([]any)[0].(map[string]any)["value"] = 1.000000001
	g = serve(t, http.StatusOK, body, "key")
	relay = redstone.New([]redstone.Gateway{{URL: g.URL, Key: "key"}}, g.Client())
	signed, err = relay.Signed(context.Background(), []string{"NVDA---24_7"})
	if err != nil || len(signed.Packages) != 4 {
		t.Fatalf("a value with nine decimals: %v", err)
	}
	if _, err := relay.Signed(context.Background(), nil); err == nil {
		t.Fatal("no package named")
	}
}

func TestTheKeyedGatewaysComeFirstAndThePublicOnesAfter(t *testing.T) {
	refused := serve(t, http.StatusOK, latest(t), "right")
	public := serve(t, http.StatusOK, latest(t), "")
	down := serve(t, http.StatusServiceUnavailable, nil, "")
	relay := redstone.New([]redstone.Gateway{{URL: refused.URL, Key: "wrong"}, {URL: down.URL}, {URL: public.URL}},
		public.Client())
	signed, err := relay.Signed(context.Background(), []string{"NVDA---24_7"})
	if err != nil || len(signed.Packages) != 5 {
		t.Fatalf("fallback: %v", err)
	}
	if refused.requests.Load() != 1 || down.requests.Load() != 1 || public.requests.Load() != 1 {
		t.Fatalf("requests %d, %d, %d", refused.requests.Load(), down.requests.Load(), public.requests.Load())
	}
	if _, err := relay.Signed(context.Background(), []string{"TSLA---24_7"}); err == nil ||
		!strings.Contains(err.Error(), "no RedStone gateway answered") || strings.Contains(err.Error(), "wrong") {
		t.Fatalf("a package no gateway has: %v", err)
	}
	gateways := redstone.Gateways(func(name string) string {
		return map[string]string{"REDSTONE_BACKUP_API_KEY": "backup"}[name]
	})
	if len(gateways) != 3 || gateways[0].Key != "backup" || gateways[1].Key != "" ||
		gateways[2].URL != "https://oracle-gateway-2.a.redstone.finance" {
		t.Fatalf("gateways %+v", gateways)
	}
}
