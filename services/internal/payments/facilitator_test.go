// SPDX-License-Identifier: MIT OR Apache-2.0

package payments_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/internal/payments"
	x402 "github.com/x402-foundation/x402/go/v2"
)

const network = x402.Network("eip155:46630")

func facilitatorOf(t *testing.T, node *chain) *httptest.Server {
	t.Helper()
	signer, err := payments.NewSigner(t.Context(), node, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(payments.NewFacilitator(network, signer, usdg, "the token"))
	t.Cleanup(server.Close)
	return server
}

func requirements(asset common.Address) map[string]any {
	return map[string]any{"scheme": "exact", "network": string(network), "amount": "1000", "asset": asset.Hex(),
		"payTo": payTo.Hex(), "maxTimeoutSeconds": 60, "extra": map[string]any{"name": "Global Dollar", "version": "1"}}
}

func payment(asset common.Address, from common.Address) map[string]any {
	return map[string]any{"x402Version": 2, "accepted": requirements(asset), "payload": map[string]any{
		"signature": "0x" + string(bytes.Repeat([]byte("11"), 65)),
		"authorization": map[string]any{"from": from.Hex(), "to": payTo.Hex(), "value": "1000", "validAfter": "0",
			"validBefore": "4102444800", "nonce": common.Hash{7}.Hex()}}}
}

func post(t *testing.T, server *httptest.Server, path, token string, body any) (int, map[string]any) {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL+path, bytes.NewReader(raw))
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	var out map[string]any
	_ = json.NewDecoder(response.Body).Decode(&out)
	return response.StatusCode, out
}

func TestTheFacilitatorServesOnlyTheClientThatHoldsItsToken(t *testing.T) {
	server := facilitatorOf(t, &chain{})
	for _, token := range []string{"", "another token"} {
		request, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/supported", nil)
		if token != "" {
			request.Header.Set("Authorization", "Bearer "+token)
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("token %q: %d", token, response.StatusCode)
		}
	}
	request, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/supported", nil)
	request.Header.Set("Authorization", "Bearer the token")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	var supported struct {
		Kinds []struct{ Scheme, Network string }
	}
	if err := json.NewDecoder(response.Body).Decode(&supported); err != nil || response.StatusCode != http.StatusOK ||
		len(supported.Kinds) != 1 || supported.Kinds[0].Scheme != "exact" || supported.Kinds[0].Network != string(network) {
		t.Fatalf("/supported: %d %+v %v", response.StatusCode, supported, err)
	}
}

func TestTheFacilitatorRefusesAnotherTokenAndAnAuthorizationUsedBefore(t *testing.T) {
	payer := common.HexToAddress("0x2222222222222222222222222222222222222222")
	node := &chain{answers: map[string]func([]any) ([]any, error){
		"authorizationState": func([]any) ([]any, error) { return []any{false}, nil }}}
	server := facilitatorOf(t, node)
	body := func(asset common.Address) map[string]any {
		return map[string]any{"paymentPayload": payment(asset, payer), "paymentRequirements": requirements(asset)}
	}

	other := common.HexToAddress("0x3333333333333333333333333333333333333333")
	for _, path := range []string{"/verify", "/settle"} {
		_, answer := post(t, server, path, "the token", body(other))
		if answer["isValid"] == true || answer["success"] == true {
			t.Fatalf("%s of another token: %v", path, answer)
		}
		if answer["invalidReason"] != "invalid_exact_evm_asset" && answer["errorReason"] != "invalid_exact_evm_asset" {
			t.Fatalf("%s of another token: %v", path, answer)
		}
	}

	node.answers["authorizationState"] = func([]any) ([]any, error) { return []any{true}, nil }
	for _, path := range []string{"/verify", "/settle"} {
		_, answer := post(t, server, path, "the token", body(usdg))
		if answer["invalidReason"] != "invalid_exact_evm_nonce_already_used" &&
			answer["errorReason"] != "invalid_exact_evm_nonce_already_used" {
			t.Fatalf("%s of a used authorization: %v", path, answer)
		}
	}
	if len(node.sent) != 0 {
		t.Fatalf("%d transactions were sent for payments that were refused", len(node.sent))
	}
	if status, _ := post(t, server, "/verify", "the token", map[string]any{"paymentPayload": nil}); status != http.StatusBadRequest {
		t.Fatalf("a request without requirements: %d", status)
	}
}

func TestTheSameAuthorizationIsNotSettledTwiceAtOnce(t *testing.T) {
	payer := common.HexToAddress("0x2222222222222222222222222222222222222222")
	entered, release := make(chan struct{}), make(chan struct{})
	node := &chain{answers: map[string]func([]any) ([]any, error){
		"authorizationState": func([]any) ([]any, error) {
			select {
			case entered <- struct{}{}:
				<-release
			default:
			}
			return []any{false}, nil
		}}}
	server := facilitatorOf(t, node)
	body := map[string]any{"paymentPayload": payment(usdg, payer), "paymentRequirements": requirements(usdg)}

	first := make(chan struct{})
	go func() {
		defer close(first)
		post(t, server, "/settle", "the token", body)
	}()
	<-entered
	status, answer := post(t, server, "/settle", "the token", body)
	if status != http.StatusBadRequest || answer["errorReason"] != "invalid_exact_evm_nonce_already_used" {
		t.Fatalf("the same authorization while it settles: %d %v", status, answer)
	}
	close(release)
	<-first
}
