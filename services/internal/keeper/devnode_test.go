// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/app"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// sent records each transaction the keepers log, by action.
type sent struct {
	slog.Handler
	mu      sync.Mutex
	actions []action
}

type action struct{ name, tx string }

func (s *sent) Handle(ctx context.Context, record slog.Record) error {
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "tx" {
			s.mu.Lock()
			s.actions = append(s.actions, action{record.Message, attr.Value.String()})
			s.mu.Unlock()
		}
		return true
	})
	return s.Handler.Handle(ctx, record)
}

func (s *sent) take() []action {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.actions
	s.actions = nil
	return out
}

// l2Gas reads a transaction's gas less its L1 part, from Arbitrum's receipt.
func l2Gas(t *testing.T, client *ethclient.Client, tx string) uint64 {
	t.Helper()
	var receipt struct {
		GasUsed      hexutil.Uint64 `json:"gasUsed"`
		GasUsedForL1 hexutil.Uint64 `json:"gasUsedForL1"`
		Status       hexutil.Uint64 `json:"status"`
	}
	if err := client.Client().Call(&receipt, "eth_getTransactionReceipt", tx); err != nil || receipt.Status != 1 {
		t.Fatalf("receipt of %s: %v, status %d", tx, err, receipt.Status)
	}
	return uint64(receipt.GasUsed - receipt.GasUsedForL1)
}

func TestTheKeepersRunOnTheDevNode(t *testing.T) {
	ctx := context.Background()
	rpcURL, registry := os.Getenv("TAPEHOUSE_RPC_URL"), os.Getenv("TAPEHOUSE_DEPLOYMENTS")
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		t.Fatal(err)
	}
	deployments, err := sdk.LoadDeployments(registry)
	if err != nil {
		t.Fatal(err)
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		t.Fatal("set PRIVATE_KEY to the dev node's funded key")
	}
	chainID := new(big.Int).SetUint64(deployments.ChainID)
	direct := sdk.NewClient(client, deployments)
	sendAs := func(from *ecdsa.PrivateKey, to common.Address, signature string, args ...any) {
		t.Helper()
		parsed, err := abi.JSON(strings.NewReader(`[` + function(signature) + `]`))
		if err != nil {
			t.Fatal(err)
		}
		name, _, _ := strings.Cut(signature, "(")
		data, err := parsed.Pack(name, args...)
		if err != nil {
			t.Fatal(err)
		}
		opts := bind.NewKeyedTransactor(from, chainID)
		tx, err := direct.Send(opts, sdk.Tx{To: to, Data: data})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if _, err := bind.WaitMined(ctx, client, tx.Hash()); err != nil {
			t.Fatal(err)
		}
	}
	send := func(to common.Address, signature string, args ...any) {
		t.Helper()
		sendAs(key, to, signature, args...)
	}

	// The indexer, as cmd/indexer runs it, for the keepers that enumerate.
	indexer, cancel := context.WithCancel(ctx)
	defer cancel()
	ready, done := make(chan string, 1), make(chan error, 1)
	cfg, err := app.ConfigFromEnv(func(name string) string {
		switch name {
		case "TAPEHOUSE_DB":
			return filepath.Join(t.TempDir(), "index.db")
		case "TAPEHOUSE_LISTEN":
			return "127.0.0.1:0"
		case "TAPEHOUSE_RPC_RATE":
			return "1000"
		}
		return os.Getenv(name)
	})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		done <- app.Run(indexer, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), func(address net.Addr) {
			ready <- "http://" + address.String()
		})
	}()
	var indexURL string
	select {
	case indexURL = <-ready:
	case err := <-done:
		t.Fatalf("the indexer did not start: %v", err)
	}

	// Robinhood's halts, as the e2e decides them.
	var haltedMu sync.Mutex
	halted := map[string]bool{"NVDA": true, "TSLA": false, "SPY": false}
	quotes := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		haltedMu.Lock()
		defer haltedMu.Unlock()
		var out []map[string]any
		for name, h := range halted {
			out = append(out, map[string]any{"tokenSymbol": name, "isTradingHalt": h})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"quotes": out})
	}))
	defer quotes.Close()

	// The halt signer: the deployment's, or the one the band's suite rotates to, both Anvil's test keys.
	signer, err := direct.Band().HaltSigner(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var haltKey *ecdsa.PrivateKey
	for _, hex := range []string{"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"} {
		if candidate, _ := crypto.HexToECDSA(hex); crypto.PubkeyToAddress(candidate.PublicKey) == signer {
			haltKey = candidate
		}
	}
	if haltKey == nil {
		t.Fatalf("the band's halt signer %s is no test key", signer.Hex())
	}

	// The keepers, their first RPC endpoint dead, so every call fails over to the dev node.
	newKeeper := func(buy bool) (*Keeper, *sent, func()) {
		recorder := &sent{Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})}
		log := slog.New(recorder)
		chain, err := Dial(ctx, []string{"http://127.0.0.1:1", rpcURL}, 1000, log)
		if err != nil {
			t.Fatal(err)
		}
		k, err := New(Config{IndexerURL: indexURL, Gateways: redstone.Gateways(os.Getenv), HaltsURL: quotes.URL, Buy: buy},
			deployments, chain, key, haltKey, log)
		if err != nil {
			t.Fatal(err)
		}
		return k, recorder, chain.Close
	}
	k, recorder, closeChain := newKeeper(true)
	defer closeChain()
	gas := map[string]uint64{}
	pass := func(name string) []action {
		t.Helper()
		if err := k.Pass(ctx, name); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		actions := recorder.take()
		for _, a := range actions {
			gas[a.name] = l2Gas(t, client, a.tx)
		}
		return actions
	}

	// One write of every tracked feed, the status first; the band's session is then known.
	if got := pass("prices"); len(got) != 1 {
		t.Fatalf("prices: %v", got)
	}
	session, err := direct.Band().Session(&bind.CallOpts{})
	if err != nil || session.State == 0 {
		t.Fatalf("the session after the keeper's write: %+v, %v", session, err)
	}
	fmt.Printf("keeper: writePrices of %d feeds, %d L2 gas; session %d, NYSE %d\n", len(mustFeeds(t, k)), gas["writePrices"],
		session.State, session.Nyse)

	// The same package again is a front-runner's: the keeper takes the next one.
	if got := pass("prices"); len(got) != 1 {
		t.Fatalf("a second write: %v", got)
	}

	// The status as three packages costs more than as RedStone's one.
	time.Sleep(PackageInterval)
	grouped := gas["writePrices"]
	k.prices = redstone.New(redstone.Gateways(os.Getenv), &http.Client{Timeout: 10 * time.Second})
	if got := pass("prices"); len(got) != 1 {
		t.Fatalf("a write of three status packages: %v", got)
	}
	fmt.Printf("keeper: the status as NY_MARKET_STATUS %d L2 gas, as its three feeds %d\n", grouped, gas["writePrices"])
	k.prices = statusGroup{k.prices}

	// Every asset with a token synced once at start, then nothing.
	got := pass("multiplier")
	fmt.Printf("keeper: %d multiplier syncs at start, %v\n", len(got), gas)
	if len(got) == 0 {
		t.Fatal("no multiplier sync at start")
	}
	if got := pass("multiplier"); len(got) != 0 {
		t.Fatalf("synced again: %v", got)
	}

	// NVDA halted, then lifted once Robinhood shows it trading.
	if got := pass("halt"); len(got) != 1 {
		t.Fatalf("halt: %v", got)
	}
	stored, err := direct.Band().Halt(&bind.CallOpts{}, "NVDA")
	if err != nil || !stored.SignedHalt || stored.Until != stored.IssuedAt+900 {
		t.Fatalf("NVDA's halt %+v, %v", stored, err)
	}
	haltedMu.Lock()
	halted["NVDA"] = false
	haltedMu.Unlock()
	for {
		head, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if head.Time > stored.IssuedAt {
			break
		}
		time.Sleep(300 * time.Millisecond)
		send(deployments.Tokens["USDG"], "mint(address,uint256)", crypto.PubkeyToAddress(key.PublicKey), big.NewInt(1))
	}
	if got := pass("halt"); len(got) != 1 {
		t.Fatalf("lift: %v", got)
	}
	if stored, err := direct.Band().Halt(&bind.CallOpts{}, "NVDA"); err != nil || stored.SignedHalt || stored.Until != 0 {
		t.Fatalf("NVDA's lift %+v, %v", stored, err)
	}

	// A borrower's WETH position at 80% of its value; WETH falls 15%: the keeper starts its auction and buys at its
	// ask, then stops it once WETH is back.
	feed := deployments.Chainlink["ETH_USD"]
	round, err := direct.LatestRound(&bind.CallOpts{}, feed)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		head, _ := client.HeaderByNumber(ctx, nil)
		send(feed, "setRound(int256,uint256)", round.Answer, new(big.Int).SetUint64(head.Time))
	})
	borrowerKey, _ := crypto.GenerateKey()
	me := crypto.PubkeyToAddress(borrowerKey.PublicKey)
	position := sdk.Cross
	fund(t, client, key, chainID, me)
	weth, accounts := deployments.Tokens["WETH"], deployments.Tapehouse["MarginAccounts"]
	send(weth, "mint(address,uint256)", me, big.NewInt(1e18))
	sendAs(borrowerKey, weth, "approve(address,uint256)", accounts, big.NewInt(1e18))
	sendAs(borrowerKey, accounts, "deposit(bytes32,address,uint256,address)", position, weth, big.NewInt(1e18), me)
	sendAs(borrowerKey, accounts, "borrow(bytes32,uint256,address,address)", position,
		new(big.Int).Sub(new(big.Int).Div(new(big.Int).Mul(round.Answer, big.NewInt(8)), big.NewInt(1000)), big.NewInt(1)), me, me)
	reach(t, client, indexURL)
	if got := pass("liquidation"); len(got) != 0 {
		fmt.Printf("keeper: before WETH falls, %v\n", got)
	}
	head, _ := client.HeaderByNumber(ctx, nil)
	send(feed, "setRound(int256,uint256)", new(big.Int).Div(new(big.Int).Mul(round.Answer, big.NewInt(85)), big.NewInt(100)),
		new(big.Int).SetUint64(head.Time))
	send(deployments.Tokens["USDG"], "mint(address,uint256)", crypto.PubkeyToAddress(key.PublicKey), big.NewInt(10_000_000_000))
	reach(t, client, indexURL)
	judged, err := direct.Liquidator().Shortfall(&bind.CallOpts{}, me, position)
	if err != nil || !judged.Short {
		t.Fatalf("the WETH position after a 15%% fall: %+v, %v", judged, err)
	}
	ours := fmt.Sprintf("%s %s", me.Hex(), hexutil.Encode(position[:]))
	find := func(actions []action, prefix string) (action, bool) {
		for _, a := range actions {
			if strings.HasPrefix(a.name, prefix) && strings.HasSuffix(a.name, ours) {
				return a, true
			}
		}
		return action{}, false
	}
	got = pass("liquidation")
	started, ok := find(got, "start ")
	bought, ok2 := find(got, "buy WETH ")
	if !ok || !ok2 {
		t.Fatalf("liquidation: %v", got)
	}
	fmt.Printf("keeper: start %d, buy %d L2 gas\n", gas[started.name], gas[bought.name])

	// A keeper that crashes and restarts repeats nothing.
	restarted, restartedRecorder, closeRestarted := newKeeper(false)
	defer closeRestarted()
	if err := restarted.Pass(ctx, "liquidation"); err != nil {
		t.Fatal(err)
	}
	again := restartedRecorder.take()
	for _, prefix := range []string{"start ", "buy WETH "} {
		if a, repeated := find(again, prefix); repeated {
			t.Fatalf("a restarted keeper repeated %v", a)
		}
	}
	head, _ = client.HeaderByNumber(ctx, nil)
	send(feed, "setRound(int256,uint256)", round.Answer, new(big.Int).SetUint64(head.Time))
	if got := pass("liquidation"); !func() bool { _, ok := find(got, "stop "); return ok }() {
		t.Fatalf("a recovered position's auction: %v", got)
	}

	// Every other keeper passes against the dev node's state as the session stands.
	for _, name := range []string{"seal", "premium", "mark", "sync", "backstop", "reopening", "recall", "shorts",
		"basket", "gapcover"} {
		got := pass(name)
		names := make([]string, len(got))
		for i, a := range got {
			names[i] = fmt.Sprintf("%s %d", a.name, gas[a.name])
		}
		fmt.Printf("keeper: %s sent %v\n", name, names)
	}
	fmt.Println("PASS")
}

// fund sends to an ETH for its gas.
func fund(t *testing.T, client *ethclient.Client, from *ecdsa.PrivateKey, chainID *big.Int, to common.Address) {
	t.Helper()
	ctx := context.Background()
	nonce, err := client.PendingNonceAt(ctx, crypto.PubkeyToAddress(from.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{From: crypto.PubkeyToAddress(from.PublicKey), To: &to,
		Value: big.NewInt(1e18)})
	if err != nil {
		t.Fatal(err)
	}
	tx, err := types.SignNewTx(from, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{ChainID: chainID,
		Nonce: nonce, GasTipCap: big.NewInt(0), GasFeeCap: new(big.Int).Mul(head.BaseFee, big.NewInt(2)), Gas: gas,
		To: &to, Value: big.NewInt(1e18)})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SendTransaction(ctx, tx); err != nil {
		t.Fatal(err)
	}
	if _, err := bind.WaitMined(ctx, client, tx.Hash()); err != nil {
		t.Fatal(err)
	}
}

// function is the ABI of a function the keepers never call, from its signature.
func function(signature string) string {
	name, rest, _ := strings.Cut(signature, "(")
	var inputs []string
	for i, typ := range strings.Split(strings.TrimSuffix(rest, ")"), ",") {
		if typ != "" {
			inputs = append(inputs, fmt.Sprintf(`{"name":"a%d","type":"%s"}`, i, typ))
		}
	}
	return fmt.Sprintf(`{"type":"function","name":"%s","inputs":[%s],"outputs":[],"stateMutability":"nonpayable"}`,
		name, strings.Join(inputs, ","))
}

func mustFeeds(t *testing.T, k *Keeper) []string {
	feeds, err := k.priceFeeds(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return feeds
}

// reach waits for the indexer to reach the chain's head.
func reach(t *testing.T, client *ethclient.Client, base string) {
	t.Helper()
	head, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for range 600 {
		response, err := http.Get(base + "/v1/status")
		if err == nil {
			var status struct {
				Head struct {
					Number uint64 `json:"number"`
				} `json:"head"`
			}
			_ = json.NewDecoder(response.Body).Decode(&status)
			_ = response.Body.Close()
			if status.Head.Number >= head {
				return
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("the indexer did not reach block %d", head)
}
