// SPDX-License-Identifier: MIT OR Apache-2.0

package payments

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	x402 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	nethttp "github.com/x402-foundation/x402/go/v2/http/nethttp"
	evmserver "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/server"
)

const maxTimeoutSeconds = 60

// Terms are what a request costs: Price of Asset, in its atomic units, paid to PayTo on Network and settled through
// the facilitator at Facilitator, to which Token is the bearer token. Name and Version are the token's EIP-712 domain, which the payer signs under. Patterns
// are the routes it applies to, a method and a path as the x402 middleware reads them, such as "GET /v1/status".
type Terms struct {
	Network     x402.Network
	Asset       common.Address
	Name        string
	Version     string
	PayTo       common.Address
	Price       *big.Int
	Facilitator string
	Token       string
	Patterns    []string
}

// bearer is the facilitator's token, which each of the client's requests presents.
type bearer string

func (b bearer) GetAuthHeaders(context.Context) (x402http.AuthHeaders, error) {
	header := map[string]string{"Authorization": "Bearer " + string(b)}
	return x402http.AuthHeaders{Verify: header, Settle: header, Supported: header}, nil
}

// Paywall returns the middleware that answers a request to one of the terms' routes with 402 Payment Required and
// what it costs, and serves it once a payment for it has been verified, settling the payment after its response.
func Paywall(t Terms) (func(http.Handler) http.Handler, error) {
	if t.Price == nil || t.Price.Sign() <= 0 {
		return nil, errors.New("the price of a request must be positive")
	}
	routes := x402http.RoutesConfig{}
	for _, pattern := range t.Patterns {
		routes[pattern] = x402http.RouteConfig{
			Accepts: x402http.PaymentOptions{{
				Scheme:  "exact",
				Network: t.Network,
				PayTo:   t.PayTo.Hex(),
				Price: map[string]any{
					"amount": t.Price.String(),
					"asset":  t.Asset.Hex(),
					"extra":  map[string]any{"name": t.Name, "version": t.Version},
				},
				MaxTimeoutSeconds: maxTimeoutSeconds,
			}},
			Description: "A read of Tapehouse's API, paid in USDG per request",
			MimeType:    "application/json",
		}
	}
	return nethttp.PaymentMiddlewareFromConfig(routes,
		nethttp.WithFacilitatorClient(x402http.NewHTTPFacilitatorClient(&x402http.FacilitatorConfig{URL: t.Facilitator,
			AuthProvider: bearer(t.Token)})),
		nethttp.WithScheme(t.Network, evmserver.NewExactEvmScheme()),
		nethttp.WithTimeout(settleTimeout+time.Minute)), nil
}
