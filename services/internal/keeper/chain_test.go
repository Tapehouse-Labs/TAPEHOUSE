// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// endpoint is an RPC endpoint double that answers eth_call as told, counting the calls it gets.
type endpoint struct {
	*httptest.Server
	calls  atomic.Int32
	answer atomic.Value
}

func newEndpoint(t *testing.T, answer string) *endpoint {
	e := &endpoint{}
	e.answer.Store(answer)
	e.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e.calls.Add(1)
		var request struct {
			ID json.RawMessage `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&request)
		switch answer := e.answer.Load().(string); answer {
		case "down":
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
		case "revert":
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(request.ID) +
				`,"error":{"code":3,"message":"execution reverted","data":"0x5d1d5d09` + string(make64()) + `"}}`))
		case "limited":
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(request.ID) +
				`,"error":{"code":-32005,"message":"limit exceeded"}}`))
		default:
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(request.ID) + `,"result":"` + answer + `"}`))
		}
	}))
	t.Cleanup(e.Close)
	return e
}

func make64() []byte {
	out := make([]byte, 192)
	for i := range out {
		out[i] = '0'
	}
	return out
}

func TestAFailedEndpointMovesTheCallsToTheNextAndARevertIsAnAnswer(t *testing.T) {
	first, second := newEndpoint(t, "down"), newEndpoint(t, "0x01")
	f, err := Dial(context.Background(), []string{first.URL, second.URL}, 1000, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	call := ethereum.CallMsg{To: &common.Address{1}, Data: []byte{1, 2, 3, 4}}
	out, err := f.CallContract(context.Background(), call, nil)
	if err != nil || len(out) != 1 || first.calls.Load() != 1 || second.calls.Load() != 1 {
		t.Fatalf("a dead first endpoint: %x, %v", out, err)
	}
	first.answer.Store("0x02")
	if out, err := f.CallContract(context.Background(), call, nil); err != nil || out[0] != 1 || first.calls.Load() != 1 {
		t.Fatalf("the endpoint that answered is not kept: %x, %v", out, err)
	}
	second.answer.Store("revert")
	_, err = f.CallContract(context.Background(), call, nil)
	if revert, ok := sdk.DecodeRevert(err); !ok || revert.Name != "PackageNotNewer" || first.calls.Load() != 1 {
		t.Fatalf("a revert retried on another endpoint: %v", err)
	}
	second.answer.Store("limited")
	if out, err := f.CallContract(context.Background(), call, nil); err != nil || out[0] != 2 || first.calls.Load() != 2 {
		t.Fatalf("a rate limit did not move the call: %x, %v", out, err)
	}
	first.answer.Store("down")
	second.answer.Store("down")
	if _, err := f.CallContract(context.Background(), call, nil); err == nil {
		t.Fatal("every endpoint down answered")
	}
}
