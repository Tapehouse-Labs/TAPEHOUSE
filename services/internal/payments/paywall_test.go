// SPDX-License-Identifier: MIT OR Apache-2.0

package payments_test

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/payments"
)

// facilitator is a double of an x402 facilitator that approves what it is told to and counts what it settles.
type facilitator struct {
	settles atomic.Int64
	succeed atomic.Bool
}

func (f *facilitator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer the token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/supported":
		_, _ = w.Write([]byte(`{"kinds":[{"x402Version":2,"scheme":"exact","network":"eip155:46630"}],"extensions":[],"signers":{}}`))
	case "/verify":
		_, _ = w.Write([]byte(`{"isValid":true,"payer":"0x2222222222222222222222222222222222222222"}`))
	case "/settle":
		f.settles.Add(1)
		if !f.succeed.Load() {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"success":false,"errorReason":"transaction_failed","transaction":"","network":"eip155:46630"}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true,"payer":"0x2222222222222222222222222222222222222222",` +
			`"transaction":"0x0000000000000000000000000000000000000000000000000000000000000abc","network":"eip155:46630"}`))
	}
}

func TestAPaywalledRouteAsksForUSDGAndIsServedOnceItsPaymentSettles(t *testing.T) {
	approving := &facilitator{}
	remote := httptest.NewServer(approving)
	t.Cleanup(remote.Close)
	wrap, err := payments.Paywall(payments.Terms{Network: network, Asset: usdg, Name: "Global Dollar", Version: "1",
		PayTo: payTo, Price: big.NewInt(1000), Facilitator: remote.URL, Token: "the token",
		Patterns: []string{"GET /v1/status", "GET /v1/contracts/[contract]/views/[function]"}})
	if err != nil {
		t.Fatal(err)
	}
	served := httptest.NewServer(wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"secret":true}`))
	})))
	t.Cleanup(served.Close)

	get := func(path string, header ...string) (*http.Response, string) {
		request, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, served.URL+path, nil)
		for i := 0; i < len(header); i += 2 {
			request.Header.Set(header[i], header[i+1])
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = response.Body.Close() }()
		var body [64]byte
		n, _ := response.Body.Read(body[:])
		return response, string(body[:n])
	}

	response, _ := get("/v1/status")
	if response.StatusCode != http.StatusPaymentRequired {
		t.Fatalf("an unpaid request: %d", response.StatusCode)
	}
	raw, err := base64.StdEncoding.DecodeString(response.Header.Get("PAYMENT-REQUIRED"))
	if err != nil {
		t.Fatal(err)
	}
	var required struct {
		X402Version int `json:"x402Version"`
		Accepts     []map[string]any
	}
	if err := json.Unmarshal(raw, &required); err != nil || required.X402Version != 2 || len(required.Accepts) != 1 {
		t.Fatalf("the requirements: %s %v", raw, err)
	}
	accepted := required.Accepts[0]
	extra, _ := accepted["extra"].(map[string]any)
	if accepted["scheme"] != "exact" || accepted["network"] != string(network) || accepted["amount"] != "1000" ||
		accepted["asset"] != usdg.Hex() || accepted["payTo"] != payTo.Hex() || extra["name"] != "Global Dollar" ||
		extra["version"] != "1" {
		t.Fatalf("what is asked: %v", accepted)
	}
	if response, _ = get("/v1/contracts/Band/views/halt"); response.StatusCode != http.StatusPaymentRequired {
		t.Fatalf("a route with parameters, unpaid: %d", response.StatusCode)
	}
	signature := func() string {
		raw, _ := json.Marshal(map[string]any{"x402Version": 2, "accepted": accepted, "payload": map[string]any{
			"signature": "0x" + strings.Repeat("11", 65),
			"authorization": map[string]any{"from": common.Address{0x22}.Hex(), "to": payTo.Hex(), "value": "1000",
				"validAfter": "0", "validBefore": "4102444800", "nonce": common.Hash{7}.Hex()}}})
		return base64.StdEncoding.EncodeToString(raw)
	}

	response, body := get("/v1/status", "PAYMENT-SIGNATURE", signature())
	if response.StatusCode != http.StatusPaymentRequired || strings.Contains(body, "secret") || approving.settles.Load() != 1 {
		t.Fatalf("a payment that does not settle: %d %q, %d settled", response.StatusCode, body, approving.settles.Load())
	}
	approving.succeed.Store(true)
	response, body = get("/v1/status", "PAYMENT-SIGNATURE", signature())
	if response.StatusCode != http.StatusOK || body != `{"secret":true}` || response.Header.Get("PAYMENT-RESPONSE") == "" {
		t.Fatalf("a payment that settles: %d %q %v", response.StatusCode, body, response.Header)
	}
	settlement, err := base64.StdEncoding.DecodeString(response.Header.Get("PAYMENT-RESPONSE"))
	if err != nil || !strings.Contains(string(settlement), `"success":true`) {
		t.Fatalf("the settlement: %s %v", settlement, err)
	}
}

func TestAPaywallWithoutAPriceIsRefused(t *testing.T) {
	if _, err := payments.Paywall(payments.Terms{Network: network, Asset: usdg, Price: big.NewInt(0)}); err == nil {
		t.Fatal("a free paywall was built")
	}
}
