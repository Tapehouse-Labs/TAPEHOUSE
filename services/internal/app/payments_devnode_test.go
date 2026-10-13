// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package app_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/api"
	"github.com/tapehouse/tapehouse/services/internal/payments"
	"github.com/tapehouse/tapehouse/services/internal/payments/paymentstest"
	"github.com/tapehouse/tapehouse/services/sdk"
	x402 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	exactclient "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/client"
	evmsigners "github.com/x402-foundation/x402/go/v2/signers/evm"
)

const mintABI = `[{"type":"function","name":"mint","stateMutability":"nonpayable","inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"outputs":[]}]`

func TestTheIndexerSellsTheReadsBeyondThePublicAllowanceInUSDG(t *testing.T) {
	deployments, err := sdk.LoadDeployments(os.Getenv("TAPEHOUSE_DEPLOYMENTS"))
	if err != nil {
		t.Fatal(err)
	}
	dialed, err := ethclient.DialContext(t.Context(), os.Getenv("TAPEHOUSE_RPC_URL"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dialed.Close)
	front := paymentstest.Multicall3{Client: dialed}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		t.Fatal(err)
	}
	signer, err := payments.NewSigner(t.Context(), front, key)
	if err != nil {
		t.Fatal(err)
	}
	usdg := deployments.Tokens["USDG"]
	network := x402.Network(fmt.Sprintf("eip155:%d", deployments.ChainID))
	facilitator := httptest.NewServer(payments.NewFacilitator(network, signer, usdg, "the facilitator's token"))
	t.Cleanup(facilitator.Close)
	payeeKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	payTo := crypto.PubkeyToAddress(payeeKey.PublicKey)

	s := serveWith(t, filepath.Join(t.TempDir(), "index.db"), "", map[string]string{
		"TAPEHOUSE_X402_FACILITATOR": facilitator.URL, "TAPEHOUSE_X402_FACILITATOR_TOKEN": "the facilitator's token",
		"TAPEHOUSE_X402_PAY_TO": payTo.Hex(), "TAPEHOUSE_X402_PRICE": "1000"})

	buyerKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	hash, err := signer.WriteContract(t.Context(), usdg.Hex(), []byte(mintABI), "mint", nil,
		crypto.PubkeyToAddress(buyerKey.PublicKey), big.NewInt(10_000))
	if err != nil {
		t.Fatal(err)
	}
	if receipt, err := signer.WaitForTransactionReceipt(t.Context(), hash); err != nil || receipt.Status != 1 {
		t.Fatalf("funding the buyer: %v %v", receipt, err)
	}

	for status := 0; status != http.StatusOK; time.Sleep(250 * time.Millisecond) {
		if status = s.get(t, "/v1/status", "", nil); status != http.StatusOK && status != http.StatusServiceUnavailable {
			t.Fatalf("waiting for the index: %d", status)
		}
	}
	var unpaid *http.Response
	free := 0
	for range 5 * api.PublicLimit.Burst {
		request, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, s.base+"/v1/status", nil)
		if unpaid, err = http.DefaultClient.Do(request); err != nil {
			t.Fatal(err)
		}
		_ = unpaid.Body.Close()
		if unpaid.StatusCode == http.StatusPaymentRequired {
			break
		}
		free++
	}
	if free < api.PublicLimit.Burst-2 || unpaid.StatusCode != http.StatusPaymentRequired ||
		unpaid.Header.Get("PAYMENT-REQUIRED") == "" || unpaid.Header.Get("X-RateLimit-Tier") != "pay-per-request" {
		t.Fatalf("a read beyond the allowance, after %d free: %d %v", free, unpaid.StatusCode, unpaid.Header)
	}
	if status := s.get(t, "/v1/openapi.yaml", "", nil); status != http.StatusTooManyRequests {
		t.Fatalf("the specification, which is not sold: %d", status)
	}

	buyer, err := evmsigners.NewClientSignerFromPrivateKey(fmt.Sprintf("0x%x", crypto.FromECDSA(buyerKey)))
	if err != nil {
		t.Fatal(err)
	}
	paying := x402http.WrapHTTPClientWithPayment(&http.Client{Timeout: time.Minute},
		x402http.Newx402HTTPClient(x402.Newx402Client(x402.WithSpendControls(x402.SpendControls{
			AllowedAssets: []x402.SpendControlAsset{{Network: network, Asset: usdg.Hex(), MaxAmountPerPayment: "10000"}}})).
			Register("eip155:*", exactclient.NewExactEvmScheme(buyer, nil))))
	response, err := paying.Get(s.base + "/v1/status")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	var read status
	if err := json.NewDecoder(response.Body).Decode(&read); err != nil || response.StatusCode != http.StatusOK ||
		read.Head.Number == 0 {
		t.Fatalf("a paid read: %d %+v %v", response.StatusCode, read, err)
	}
	if balance, err := signer.GetBalance(t.Context(), payTo.Hex(), usdg.Hex()); err != nil || balance.Int64() != 1000 {
		t.Fatalf("the payee holds %v: %v", balance, err)
	}
}
