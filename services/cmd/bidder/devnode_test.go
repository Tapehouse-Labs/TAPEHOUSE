// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package main

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
	"github.com/tapehouse/tapehouse/services/sdk"
)

type action struct{ name, tx string }

// txLog records each transaction the bidder logs, by action.
type txLog struct {
	slog.Handler
	mu      sync.Mutex
	actions []action
}

func (s *txLog) Handle(ctx context.Context, record slog.Record) error {
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

func (s *txLog) take() []action {
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

func TestTheBidderBuysAndBidsOnTheDevNode(t *testing.T) {
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
		name, rest, _ := strings.Cut(signature, "(")
		var inputs []string
		for i, typ := range strings.Split(strings.TrimSuffix(rest, ")"), ",") {
			inputs = append(inputs, fmt.Sprintf(`{"name":"a%d","type":"%s"}`, i, typ))
		}
		parsed, err := abi.JSON(strings.NewReader(fmt.Sprintf(`[{"type":"function","name":"%s","inputs":[%s],"outputs":[],"stateMutability":"nonpayable"}]`,
			name, strings.Join(inputs, ","))))
		if err != nil {
			t.Fatal(err)
		}
		data, err := parsed.Pack(name, args...)
		if err != nil {
			t.Fatal(err)
		}
		tx, err := direct.Send(bind.NewKeyedTransactor(from, chainID), sdk.Tx{To: to, Data: data})
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

	botKey, _ := crypto.GenerateKey()
	bot := crypto.PubkeyToAddress(botKey.PublicKey)
	fund(t, client, key, chainID, bot)
	usdg, nvda, weth := deployments.Tokens["USDG"], deployments.Tokens["NVDA"], deployments.Tokens["WETH"]
	send(usdg, "mint(address,uint256)", bot, big.NewInt(10_000_000_000))
	recorder := &txLog{Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})}
	bidder := &Bidder{Client: direct, Chain: client, Opts: bind.NewKeyedTransactor(botKey, chainID), Budget: big.NewInt(1000e6),
		DiscountBps: 0, CommitLead: 15 * time.Minute, Assets: []string{"NVDA"}, Log: slog.New(recorder), IndexerURL: indexURL}
	bidder.State, err = LoadState(filepath.Join(t.TempDir(), "bidder.json"))
	if err != nil {
		t.Fatal(err)
	}

	var head *types.Header
	band, err := direct.Band().Quote(&bind.CallOpts{}, "NVDA")
	if err != nil {
		t.Fatal(err)
	}
	if band.State == 0 {
		fmt.Println("bidder: NVDA's band is halted, so its Dutch purchase is not run now")
	} else {
		// A borrower's isolated NVDA position: a little NVDA and WETH borrowed to 80% of its value; WETH falls 10% on the
		// dev node's ETH/USD stub, which makes the position short, and its auction is started.
		feed := deployments.Chainlink["ETH_USD"]
		round, err := direct.LatestRound(&bind.CallOpts{}, feed)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			head, _ = client.HeaderByNumber(ctx, nil)
			send(feed, "setRound(int256,uint256)", round.Answer, new(big.Int).SetUint64(head.Time))
		})
		head, _ = client.HeaderByNumber(ctx, nil)
		send(feed, "setRound(int256,uint256)", round.Answer, new(big.Int).SetUint64(head.Time))
		borrowerKey, _ := crypto.GenerateKey()
		borrower := crypto.PubkeyToAddress(borrowerKey.PublicKey)
		position, err := sdk.ToBytes32("NVDA")
		if err != nil {
			t.Fatal(err)
		}
		fund(t, client, key, chainID, borrower)
		accounts, vault := deployments.Tapehouse["MarginAccounts"], deployments.Tapehouse["SupplyVault"]
		send(usdg, "mint(address,uint256)", crypto.PubkeyToAddress(key.PublicKey), big.NewInt(10_000_000_000))
		send(usdg, "approve(address,uint256)", vault, big.NewInt(10_000_000_000))
		send(vault, "deposit(uint256,address)", big.NewInt(10_000_000_000), crypto.PubkeyToAddress(key.PublicKey))
		small := big.NewInt(1e16)
		send(weth, "mint(address,uint256)", borrower, big.NewInt(1e18))
		send(nvda, "mint(address,uint256)", borrower, small)
		sendAs(borrowerKey, weth, "approve(address,uint256)", accounts, big.NewInt(1e18))
		sendAs(borrowerKey, nvda, "approve(address,uint256)", accounts, small)
		sendAs(borrowerKey, accounts, "deposit(bytes32,address,uint256,address)", position, weth, big.NewInt(1e18), borrower)
		sendAs(borrowerKey, accounts, "deposit(bytes32,address,uint256,address)", position, nvda, small, borrower)
		sendAs(borrowerKey, accounts, "borrow(bytes32,uint256,address,address)", position,
			new(big.Int).Sub(new(big.Int).Div(new(big.Int).Mul(round.Answer, big.NewInt(8)), big.NewInt(1000)), big.NewInt(1)), borrower, borrower)
		head, _ = client.HeaderByNumber(ctx, nil)
		send(feed, "setRound(int256,uint256)", new(big.Int).Div(new(big.Int).Mul(round.Answer, big.NewInt(90)), big.NewInt(100)),
			new(big.Int).SetUint64(head.Time))
		judged, err := direct.Liquidator().Shortfall(&bind.CallOpts{}, borrower, position)
		if err != nil || !judged.Short {
			t.Fatalf("the NVDA position after a 10%% fall: %+v, %v", judged, err)
		}
		start, err := direct.Liquidator().Start(borrower, position)
		if err != nil {
			t.Fatal(err)
		}
		started, err := direct.Send(bind.NewKeyedTransactor(key, chainID), start)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := bind.WaitMined(ctx, client, started.Hash()); err != nil {
			t.Fatal(err)
		}
		reach(t, client, indexURL)

		// The bidder buys NVDA once the ask has fallen to the band's low edge.
		gas := map[string]uint64{}
		var cost, bought *big.Int
		before, _ := direct.BalanceOf(&bind.CallOpts{}, usdg, bot)
		deadline := time.Now().Add(15 * time.Minute)
		for time.Now().Before(deadline) {
			head, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err := bidder.Dutch(ctx, head); err != nil {
				t.Fatalf("the Dutch loop: %v", err)
			}
			for _, a := range recorder.take() {
				gas[a.name] = l2Gas(t, client, a.tx)
			}
			held, _ := direct.BalanceOf(&bind.CallOpts{}, nvda, bot)
			if held.Sign() > 0 {
				after, _ := direct.BalanceOf(&bind.CallOpts{}, usdg, bot)
				bought, cost = held, new(big.Int).Sub(before, after)
				break
			}
			time.Sleep(5 * time.Second)
		}
		if bought == nil {
			quote, _ := direct.Band().Quote(&bind.CallOpts{}, "NVDA")
			ask, _ := direct.Liquidator().Price(&bind.CallOpts{}, borrower, position, nvda)
			t.Fatalf("the bidder bought nothing in 15 minutes: ask %v, band %+v", ask, quote)
		}
		quote, err := direct.Band().Quote(&bind.CallOpts{}, "NVDA")
		if err != nil {
			t.Fatal(err)
		}
		limit := Escrow(bought, new(big.Int).SetUint64(quote.Low))
		if cost.Cmp(limit) > 0 {
			t.Fatalf("the bidder paid %v for %v NVDA, above its price's %v", cost, bought, limit)
		}
		for name, used := range gas {
			fmt.Printf("bidder: %s, %d L2 gas\n", name, used)
		}
		fmt.Printf("bidder: bought %v NVDA for %v USDG at a low edge of %d\n", bought, cost, quote.Low)

	}

	// The reopening auction, in whatever phase the dev node's session is in.
	phase, err := direct.ReopeningAuction().Phase(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	nonce, _ := client.PendingNonceAt(ctx, bot)
	head, _ = client.HeaderByNumber(ctx, nil)
	if err := bidder.Reopening(ctx, head); err != nil {
		t.Fatalf("the reopening loop: %v", err)
	}
	after, _ := client.PendingNonceAt(ctx, bot)
	switch {
	case phase.OpenMs == 0:
		if after != nonce || len(bidder.State.Bids()) != 0 {
			t.Fatalf("the reopening auction is outside its phases and the bidder sent %d transactions", after-nonce)
		}
		fmt.Println("bidder: the reopening auction is outside its phases, nothing sent")
	case !phase.Revealing:
		fmt.Printf("bidder: bids are committed for %d; %d held, %d transactions sent\n", phase.OpenMs, len(bidder.State.Bids()), after-nonce)
		if len(bidder.State.Bids()) != 0 && after == nonce {
			t.Fatal("a bid was planned and none was committed")
		}
	default:
		fmt.Printf("bidder: bids are revealed for %d, %d held\n", phase.OpenMs, len(bidder.State.Bids()))
	}
	fmt.Println("bidder: PASS")
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
	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{From: crypto.PubkeyToAddress(from.PublicKey), To: &to, Value: big.NewInt(1e18)})
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
