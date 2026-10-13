// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package payments_test

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/payments"
	"github.com/tapehouse/tapehouse/services/internal/payments/paymentstest"
	x402 "github.com/x402-foundation/x402/go/v2"
	x402http "github.com/x402-foundation/x402/go/v2/http"
	exactclient "github.com/x402-foundation/x402/go/v2/mechanisms/evm/exact/client"
	evmsigners "github.com/x402-foundation/x402/go/v2/signers/evm"
)

const (
	price      = 1_000
	mintABI    = `[{"type":"function","name":"mint","stateMutability":"nonpayable","inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"outputs":[]}]`
	balanceABI = `[{"type":"function","name":"balanceOf","stateMutability":"view","inputs":[{"name":"account","type":"address"}],"outputs":[{"name":"","type":"uint256"}]}]`
)

type world struct {
	signer   *payments.Signer
	network  x402.Network
	token    common.Address
	payTo    common.Address
	resource *httptest.Server
	chain    paymentstest.Multicall3
}

func (w *world) balance(t *testing.T, account common.Address) *big.Int {
	t.Helper()
	got, err := w.signer.GetBalance(t.Context(), account.Hex(), w.token.Hex())
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func (w *world) fund(t *testing.T, account common.Address, amount int64) {
	t.Helper()
	hash, err := w.signer.WriteContract(t.Context(), w.token.Hex(), []byte(mintABI), "mint", nil, account, big.NewInt(amount))
	if err != nil {
		t.Fatal(err)
	}
	if receipt, err := w.signer.WaitForTransactionReceipt(t.Context(), hash); err != nil || receipt.Status != 1 {
		t.Fatalf("minting: %v %v", receipt, err)
	}
}

func setup(t *testing.T) *world {
	t.Helper()
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
	token := common.HexToAddress(os.Getenv("TAPEHOUSE_USDG"))
	name, version, err := payments.TokenDomain(t.Context(), front, token)
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := signer.GetChainID(t.Context())
	network := x402.Network(fmt.Sprintf("eip155:%d", chainID))
	facilitator := httptest.NewServer(payments.NewFacilitator(network, signer, token, "the token"))
	t.Cleanup(facilitator.Close)

	payTo := common.BytesToAddress(crypto.Keccak256([]byte(t.Name())))
	wrap, err := payments.Paywall(payments.Terms{Network: network, Asset: token, Name: name, Version: version,
		PayTo: payTo, Price: big.NewInt(price), Facilitator: facilitator.URL, Token: "the token",
		Patterns: []string{"GET /v1/status"}})
	if err != nil {
		t.Fatal(err)
	}
	resource := httptest.NewServer(wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"read":true}`))
	})))
	t.Cleanup(resource.Close)
	return &world{signer: signer, network: network, token: token, payTo: payTo, resource: resource, chain: front}
}

// payer is a client of the resource that pays from its own key, and records the payment it signs.
type payer struct {
	address common.Address
	client  *http.Client
	header  chan string
}

type recorder struct {
	next   http.RoundTripper
	header chan string
}

func (r recorder) RoundTrip(request *http.Request) (*http.Response, error) {
	if signature := request.Header.Get("PAYMENT-SIGNATURE"); signature != "" {
		select {
		case r.header <- signature:
		default:
		}
	}
	return r.next.RoundTrip(request)
}

func (w *world) newPayer(t *testing.T) payer {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signer, err := evmsigners.NewClientSignerFromPrivateKey(fmt.Sprintf("0x%x", crypto.FromECDSA(key)))
	if err != nil {
		t.Fatal(err)
	}
	header := make(chan string, 1)
	paying := x402http.WrapHTTPClientWithPayment(&http.Client{Transport: recorder{http.DefaultTransport, header}, Timeout: time.Minute},
		x402http.Newx402HTTPClient(x402.Newx402Client(x402.WithSpendControls(x402.SpendControls{
			AllowedAssets: []x402.SpendControlAsset{{Network: w.network, Asset: w.token.Hex(), MaxAmountPerPayment: "10000"}}})).Register("eip155:*", exactclient.NewExactEvmScheme(signer, nil))))
	return payer{address: crypto.PubkeyToAddress(key.PublicKey), client: paying, header: header}
}

func TestAClientPaysInUSDGForARead(t *testing.T) {
	w := setup(t)
	buyer := w.newPayer(t)
	w.fund(t, buyer.address, 5*price)

	response, err := buyer.client.Get(w.resource.URL + "/v1/status")
	if err != nil {
		t.Fatal(err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK || response.Header.Get("PAYMENT-RESPONSE") == "" {
		t.Fatalf("a paid read: %d %v", response.StatusCode, response.Header)
	}
	if got := w.balance(t, w.payTo); got.Int64() != price {
		t.Fatalf("the payee holds %s", got)
	}
	if got := w.balance(t, buyer.address); got.Int64() != 4*price {
		t.Fatalf("the payer holds %s", got)
	}
}

func TestAPaymentCannotBeReplayedAndAnEmptyAccountCannotPay(t *testing.T) {
	w := setup(t)
	buyer := w.newPayer(t)
	w.fund(t, buyer.address, 5*price)
	response, err := buyer.client.Get(w.resource.URL + "/v1/status")
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("the payment: %v %v", response, err)
	}
	_ = response.Body.Close()
	signed := <-buyer.header

	for range 3 {
		request, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, w.resource.URL+"/v1/status", nil)
		request.Header.Set("PAYMENT-SIGNATURE", signed)
		replay, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		_ = replay.Body.Close()
		if replay.StatusCode != http.StatusPaymentRequired {
			t.Fatalf("a replayed payment: %d", replay.StatusCode)
		}
	}
	if got := w.balance(t, w.payTo); got.Int64() != price {
		t.Fatalf("the payee holds %s after replays", got)
	}

	broke := w.newPayer(t)
	response, err = broke.client.Get(w.resource.URL + "/v1/status")
	if err == nil && response.StatusCode == http.StatusOK {
		t.Fatal("an account without USDG paid")
	}
	if response != nil {
		_ = response.Body.Close()
	}
	if got := w.balance(t, w.payTo); got.Int64() != price {
		t.Fatalf("the payee holds %s after an empty account's attempt", got)
	}
}

func TestPaymentsFromManyClientsAtOnceAllSettle(t *testing.T) {
	w := setup(t)
	const buyers = 6
	var group sync.WaitGroup
	failures := make(chan error, buyers)
	for range buyers {
		buyer := w.newPayer(t)
		w.fund(t, buyer.address, price)
		group.Add(1)
		go func() {
			defer group.Done()
			response, err := buyer.client.Get(w.resource.URL + "/v1/status")
			if err != nil {
				failures <- err
				return
			}
			_ = response.Body.Close()
			if response.StatusCode != http.StatusOK {
				failures <- fmt.Errorf("status %d", response.StatusCode)
			}
		}()
	}
	group.Wait()
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	if got := w.balance(t, w.payTo); got.Int64() != buyers*price {
		t.Fatalf("the payee holds %s", got)
	}
}
