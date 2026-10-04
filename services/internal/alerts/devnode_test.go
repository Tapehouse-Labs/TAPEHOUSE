// SPDX-License-Identifier: MIT OR Apache-2.0

//go:build devnode

package alerts

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/internal/keeper"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/internal/redstone"
	"github.com/tapehouse/tapehouse/services/sdk"
)

func TestTheAlertsRunOnTheDevNode(t *testing.T) {
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, os.Getenv("TAPEHOUSE_RPC_URL"))
	if err != nil {
		t.Fatal(err)
	}
	deployments, err := sdk.LoadDeployments(os.Getenv("TAPEHOUSE_DEPLOYMENTS"))
	if err != nil {
		t.Fatal(err)
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(os.Getenv("PRIVATE_KEY"), "0x"))
	if err != nil {
		t.Fatal("set PRIVATE_KEY to the dev node's funded key")
	}
	chainID := new(big.Int).SetUint64(deployments.ChainID)
	account := crypto.PubkeyToAddress(key.PublicKey)
	direct := sdk.NewClient(client, deployments)
	send := func(to common.Address, signature string, args ...any) {
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
		tx, err := direct.Send(bind.NewKeyedTransactor(key, chainID), sdk.Tx{To: to, Data: data})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if _, err := bind.WaitMined(ctx, client, tx.Hash()); err != nil {
			t.Fatal(err)
		}
	}
	act := func(tx sdk.Tx, err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		sent, err := direct.Send(bind.NewKeyedTransactor(key, chainID), tx)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := bind.WaitMined(ctx, client, sent.Hash()); err != nil {
			t.Fatal(err)
		}
	}

	// The band is live for two minutes after RedStone's packages are written, so each read follows a write.
	relay := redstone.New(redstone.Gateways(os.Getenv), &http.Client{Timeout: 20 * time.Second})
	feeds := append(append([]string{}, keeper.StatusFeeds...), "USA500.Y---24_7", "NVDA---24_7")
	refresh := func() {
		t.Helper()
		for attempt := 0; ; attempt++ {
			tx, err := direct.Band().WritePrices(ctx, relay, feeds)
			if err == nil {
				if err = direct.Simulate(&bind.CallOpts{Context: ctx, From: account}, tx); err == nil {
					act(tx, nil)
					return
				}
			}
			if attempt == 5 {
				t.Fatalf("the band's prices were not written: %v", err)
			}
			time.Sleep(11 * time.Second)
		}
	}

	// A position of ten SPY, near its limit.
	position := symbol("SPY")
	spy, accountsAddress := deployments.Tokens["SPY"], deployments.Tapehouse["MarginAccounts"]
	ten := new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18))
	send(spy, "mint(address,uint256)", account, ten)
	send(spy, "approve(address,uint256)", accountsAddress, ten)
	act(direct.Accounts().Deposit(position, spy, ten, account))
	refresh()
	reader := NewReader(direct)
	m, err := reader.Market(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	j, err := reader.Judge(ctx, m, account, position)
	if err != nil || !j.Valuation.Judged {
		t.Fatalf("the position is not judged: settled %v assets %+v pos %+v, %v", m.Settled, m.Assets, j.Position, err)
	}
	borrow := new(big.Int).Div(new(big.Int).Sub(j.Valuation.Equity, new(big.Int).Div(new(big.Int).Mul(j.Valuation.Requirement, big.NewInt(115)), big.NewInt(100))), big.NewInt(1e12))
	if borrow.Sign() <= 0 {
		t.Fatalf("no room to borrow: %+v", j.Valuation)
	}
	usdg, vault := deployments.Tokens["USDG"], deployments.Tapehouse["SupplyVault"]
	supply := big.NewInt(20_000_000_000)
	send(usdg, "mint(address,uint256)", account, supply)
	send(usdg, "approve(address,uint256)", vault, supply)
	send(vault, "deposit(uint256,address)", supply, account)
	refresh()
	act(direct.Accounts().Borrow(position, borrow, account, account))

	// The index names the position the account deposited into, as the indexer serves its Deposit events.
	index := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		events := []map[string]any{}
		if q := r.URL.Query(); q.Get("contract") == "tapehouse.MarginAccounts" && q.Get("event") == "Deposit" && q.Get("arg[account]") == account.Hex() {
			events = append(events, map[string]any{"block": 1, "logIndex": 0, "args": map[string]any{"position": hexutil.Encode(position[:])}})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"events": events, "next": nil})
	}))
	defer index.Close()
	indexURL := index.URL

	// The service over the real chain, and doubles of the SMTP server and the Bot API.
	smtpServer, bot := newSMTPDouble(t, true), newBotDouble(t)
	email, err := NewEmail(smtpServer.url("alerts", "s3cret"), "alerts@tapehouse.example")
	if err != nil {
		t.Fatal(err)
	}
	email.TLS = &tls.Config{RootCAs: smtpServer.pool}
	chat := NewTelegram("TOKEN", bot.server.URL)
	store, err := OpenStore(filepath.Join(t.TempDir(), "alerts.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))
	links := Links{BaseURL: "https://alerts.example", Secret: bytes.Repeat([]byte{9}, 32)}
	server := NewServer(store, deployments.ChainID, links, email, chat, log)
	front := httptest.NewServer(server.Handler())
	defer front.Close()
	subscribe := func(channel, destination string) map[string]any {
		issued := time.Now().Unix()
		body, _ := json.Marshal(map[string]any{"account": account, "channel": channel, "destination": destination, "issuedAt": issued,
			"signature": sign(t, key, SubscribeText(deployments.ChainID, account, channel, issued))})
		response, err := http.Post(front.URL+"/v1/subscriptions", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = response.Body.Close() }()
		var out map[string]any
		_ = json.NewDecoder(response.Body).Decode(&out)
		if response.StatusCode != http.StatusAccepted {
			t.Fatalf("%s subscription: %d %v", channel, response.StatusCode, out)
		}
		return out
	}
	subscribe(Email, "ada@example.org")
	smtpServer.mu.Lock()
	confirmation := link.FindStringSubmatch(smtpServer.data)
	smtpServer.mu.Unlock()
	if confirmation == nil {
		t.Fatal("no confirmation email")
	}
	if response, err := http.Get(front.URL + "/v1/confirm/" + confirmation[1]); err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("confirm: %v", err)
	}
	code := subscribe(Telegram, "")["code"].(string)
	server.Chat(ctx, "4242", "/start "+code)
	smtpServer.mu.Lock()
	smtpServer.data = ""
	smtpServer.mu.Unlock()
	bot.mu.Lock()
	bot.sent = nil
	bot.mu.Unlock()

	service := NewService(store, reader, keeper.NewIndex(indexURL, http.DefaultClient),
		map[string]Channel{Email: email, Telegram: chat}, links, server, log)
	refresh()
	if err := service.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	smtpServer.mu.Lock()
	mailed := smtpServer.data
	smtpServer.mu.Unlock()
	bot.mu.Lock()
	chatted := bot.sent
	bot.mu.Unlock()
	if !strings.Contains(mailed, "close to its requirement") || !strings.Contains(mailed, "/v1/unsubscribe/") {
		t.Fatalf("the email:\n%s\nlogs:\n%s", mailed, logs.String())
	}
	if len(chatted) != 1 || chatted[0]["chat_id"] != "4242" || !strings.Contains(chatted[0]["text"].(string), "close to its requirement") {
		t.Fatalf("telegram %v\nlogs:\n%s", chatted, logs.String())
	}
	if err := service.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	bot.mu.Lock()
	again := len(bot.sent)
	bot.mu.Unlock()
	if again != 1 {
		t.Fatalf("a second alert at the same ratio: %d messages", again)
	}
	for _, secret := range []string{"ada@example.org", "4242"} {
		if strings.Contains(logs.String(), secret) {
			t.Fatalf("a log line names a destination:\n%s", logs.String())
		}
	}

	// The valuation against the liquidator's and the requirement against the engine's, at one block.
	refresh()
	m, err = reader.Market(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	j, err = reader.Judge(ctx, m, account, position)
	if err != nil {
		t.Fatal(err)
	}
	if !j.Evaluated || !j.Valuation.Judged {
		t.Fatalf("the Go requirement is not the engine's at block %d: %+v", m.Block, j)
	}
	opts := &bind.CallOpts{Context: ctx, BlockNumber: new(big.Int).SetUint64(m.Block)}
	judged, err := direct.Liquidator().Shortfall(opts, account, position)
	if err != nil {
		t.Fatal(err)
	}
	if judged.Equity.Cmp(j.Valuation.Equity) != 0 || judged.Requirement.Cmp(j.Valuation.Requirement) != 0 ||
		judged.Short != j.Valuation.Short() || judged.Closed != m.Closed {
		t.Fatalf("the liquidator says %+v, Go says %+v closed %v", judged, j.Valuation, m.Closed)
	}
	asset := m.Assets[m.index("SPY")]
	rebuilt := closedBand(asset, asset.Mid)
	if m.Closed && (rebuilt.Low != asset.Low || rebuilt.High.Cmp(asset.High) != 0) {
		t.Fatalf("the band rebuilt at its centre is %d..%s, the chain's is %d..%s", rebuilt.Low, rebuilt.High, asset.Low, asset.High)
	}

	// The weekend price: short at both edges there, not a unit above it.
	p := j.Position
	weekend := m.WeekendPrice(p, "SPY")
	if weekend.Status != Priced {
		t.Fatalf("weekend price %+v", weekend)
	}
	short := func(centre uint64) bool {
		edges := m.centres(closedBand)
		edges.Low[m.index("SPY")], edges.High[m.index("SPY")] = m.edge(m.index("SPY"), closedBand, centre)
		return m.Assess(p, edges, true, margin.Regime{Kind: margin.Closed}).Short()
	}
	if !short(weekend.Centre) || (weekend.Centre < asset.Mid && short(weekend.Centre+1)) {
		t.Fatalf("the weekend price %d: short there %v, a unit above %v", weekend.Centre, short(weekend.Centre), short(weekend.Centre+1))
	}
	fmt.Printf("alerts: threshold alert by email and Telegram at block %d; weekend price %s against a centre of %s, the liquidator agrees on equity %s and requirement %s\n",
		m.Block, price(weekend.Centre), price(asset.Mid), usd(judged.Equity, 18), usd(judged.Requirement, 18))
}
