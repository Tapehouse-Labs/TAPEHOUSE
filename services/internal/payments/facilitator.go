// SPDX-License-Identifier: MIT OR Apache-2.0

package payments

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	x402 "github.com/x402-foundation/x402/go/v2"
	"github.com/x402-foundation/x402/go/v2/mechanisms/evm"
	exact "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/facilitator"
)

const (
	verifyTimeout = 30 * time.Second
	settleTimeout = 2*time.Minute + 10*time.Second
	maxBody       = 64 << 10
)

const authorizationStateABI = `[{"type":"function","name":"authorizationState","stateMutability":"view","inputs":[{"name":"authorizer","type":"address"},{"name":"nonce","type":"bytes32"}],"outputs":[{"name":"","type":"bool"}]}]`

const (
	reasonAsset     = "invalid_exact_evm_asset"
	reasonNonceUsed = "invalid_exact_evm_nonce_already_used"
)

// Facilitator is the x402 facilitator's HTTP interface, as the protocol defines it: GET /supported, POST /verify and
// POST /settle. It verifies and settles the "exact" scheme on one chain, for one token, and the transactions it sends
// are the signer's, so it serves only the client that holds its token. USDG answers a transfer authorization that was
// used already with an event and no transfer, not a revert, so the facilitator refuses one it sees used before it
// sends a transaction for it.
type Facilitator struct {
	core   *x402.X402Facilitator
	signer evm.FacilitatorEvmSigner
	asset  common.Address
	token  []byte
	mux    *http.ServeMux
	mu     sync.Mutex
	held   map[string]bool
}

// NewFacilitator returns a facilitator of the exact scheme on network, settling through signer payments of asset, for
// the client that presents token as a bearer token.
func NewFacilitator(network x402.Network, signer evm.FacilitatorEvmSigner, asset common.Address, token string) *Facilitator {
	f := &Facilitator{core: x402.Newx402Facilitator(), signer: signer, asset: asset, token: []byte(token),
		mux: http.NewServeMux(), held: map[string]bool{}}
	f.core.Register([]x402.Network{network}, exact.NewExactEvmScheme(signer, nil))
	f.core.OnBeforeVerify(func(c x402.FacilitatorVerifyContext) (*x402.FacilitatorBeforeHookResult, error) {
		return f.screen(c.Ctx, c.Payload, c.Requirements)
	})
	f.core.OnBeforeSettle(func(c x402.FacilitatorSettleContext) (*x402.FacilitatorBeforeHookResult, error) {
		return f.screen(c.Ctx, c.Payload, c.Requirements)
	})
	f.mux.HandleFunc("GET /supported", f.supported)
	f.mux.HandleFunc("POST /verify", f.verify)
	f.mux.HandleFunc("POST /settle", f.settle)
	return f
}

// ServeHTTP serves the facilitator's three operations to a client with its token.
func (f *Facilitator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	presented, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || len(f.token) == 0 || subtle.ConstantTimeCompare([]byte(presented), f.token) != 1 {
		reply(w, http.StatusUnauthorized, map[string]string{"error": "a bearer token is required"})
		return
	}
	f.mux.ServeHTTP(w, r)
}

// screen refuses a payment of another token, and one whose authorization the token has used already.
func (f *Facilitator) screen(ctx context.Context, payload x402.PaymentPayloadView,
	requirements x402.PaymentRequirementsView) (*x402.FacilitatorBeforeHookResult, error) {
	if !strings.EqualFold(requirements.GetAsset(), f.asset.Hex()) {
		return &x402.FacilitatorBeforeHookResult{Abort: true, Reason: reasonAsset,
			Message: "this facilitator settles " + f.asset.Hex() + " only"}, nil
	}
	authorization, _ := payload.GetPayload()["authorization"].(map[string]any)
	from, _ := authorization["from"].(string)
	nonce, _ := authorization["nonce"].(string)
	if !common.IsHexAddress(from) || len(nonce) != 66 {
		return nil, nil
	}
	used, err := f.signer.ReadContract(ctx, f.asset.Hex(), []byte(authorizationStateABI), "authorizationState",
		common.HexToAddress(from), common.HexToHash(nonce))
	if err != nil {
		return nil, err
	}
	if spent, _ := used.(bool); spent {
		return &x402.FacilitatorBeforeHookResult{Abort: true, Reason: reasonNonceUsed}, nil
	}
	return nil, nil
}

type request struct {
	PaymentPayload      json.RawMessage `json:"paymentPayload"`
	PaymentRequirements json.RawMessage `json:"paymentRequirements"`
}

func (f *Facilitator) supported(w http.ResponseWriter, _ *http.Request) {
	reply(w, http.StatusOK, f.core.GetSupported())
}

func (f *Facilitator) verify(w http.ResponseWriter, r *http.Request) {
	body, ok := read(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), verifyTimeout)
	defer cancel()
	result, err := f.core.Verify(ctx, body.PaymentPayload, body.PaymentRequirements)
	var invalid *x402.VerifyError
	switch {
	case errors.As(err, &invalid):
		reply(w, http.StatusBadRequest, x402.VerifyResponse{InvalidReason: invalid.InvalidReason,
			InvalidMessage: invalid.InvalidMessage, Payer: invalid.Payer})
	case err != nil:
		reply(w, http.StatusInternalServerError, map[string]string{"error": "the payment could not be verified"})
	default:
		reply(w, http.StatusOK, result)
	}
}

func (f *Facilitator) settle(w http.ResponseWriter, r *http.Request) {
	body, ok := read(w, r)
	if !ok {
		return
	}
	if key := authorizationOf(body.PaymentPayload); key != "" {
		if !f.claim(key) {
			reply(w, http.StatusBadRequest, x402.SettleResponse{ErrorReason: reasonNonceUsed, Network: network(body)})
			return
		}
		defer f.release(key)
	}
	ctx, cancel := context.WithTimeout(r.Context(), settleTimeout)
	defer cancel()
	result, err := f.core.Settle(ctx, body.PaymentPayload, body.PaymentRequirements)
	var failed *x402.SettleError
	switch {
	case errors.As(err, &failed):
		reply(w, http.StatusBadRequest, x402.SettleResponse{ErrorReason: failed.ErrorReason,
			ErrorMessage: failed.ErrorMessage, Payer: failed.Payer, Transaction: failed.Transaction, Network: failed.Network})
	case err != nil:
		reply(w, http.StatusInternalServerError, map[string]string{"error": "the payment could not be settled"})
	default:
		reply(w, http.StatusOK, result)
	}
}

// claim holds an authorization while its settlement is in flight, so that the same one, sent a second time before
// the first is mined, is refused rather than settled again at the facilitator's cost.
func (f *Facilitator) claim(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.held[key] {
		return false
	}
	f.held[key] = true
	return true
}

func (f *Facilitator) release(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.held, key)
}

// authorizationOf names the authorization a payment carries by its payer and nonce, or is empty where it carries none.
func authorizationOf(payload json.RawMessage) string {
	var parsed struct {
		Payload struct {
			Authorization struct{ From, Nonce string }
		}
	}
	if json.Unmarshal(payload, &parsed) != nil {
		return ""
	}
	authorization := parsed.Payload.Authorization
	if authorization.From == "" || authorization.Nonce == "" {
		return ""
	}
	return strings.ToLower(authorization.From + ":" + authorization.Nonce)
}

func network(body request) x402.Network {
	var parsed struct{ Network string }
	_ = json.Unmarshal(body.PaymentRequirements, &parsed)
	return x402.Network(parsed.Network)
}

func read(w http.ResponseWriter, r *http.Request) (request, bool) {
	var body request
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBody)).Decode(&body); err != nil ||
		len(body.PaymentPayload) == 0 || len(body.PaymentRequirements) == 0 {
		reply(w, http.StatusBadRequest, map[string]string{"error": "a payment payload and its requirements are required"})
		return body, false
	}
	return body, true
}

func reply(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
