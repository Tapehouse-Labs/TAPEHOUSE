// SPDX-License-Identifier: MIT OR Apache-2.0

package sponsor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/sponsorpaymaster"
)

var (
	entryPoint     = common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032")
	factory        = common.HexToAddress("0x91E60e0613810449d098b0b5Ec8b51A0FE8c8985")
	implementation = common.HexToAddress("0x5555555555555555555555555555555555555555")
	paymaster      = common.HexToAddress("0x6666666666666666666666666666666666666666")
	accounts       = common.HexToAddress("0x7777777777777777777777777777777777777777")
	usdg           = common.HexToAddress("0x8888888888888888888888888888888888888888")
	owner          = common.HexToAddress("0x9999999999999999999999999999999999999999")
	stranger       = common.HexToAddress("0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa")
	signerKey, _   = crypto.HexToECDSA("4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318")
	now            = time.Unix(1_790_000_000, 0)
)

type fakeChain struct {
	signer     common.Address
	used       map[common.Address]int64
	code       map[common.Address][]byte
	slots      map[common.Address]common.Address
	baseFee    *big.Int
	l1Gas      uint64
	failHeader bool
}

func newChain() *fakeChain {
	return &fakeChain{signer: crypto.PubkeyToAddress(signerKey.PublicKey), used: map[common.Address]int64{},
		code: map[common.Address][]byte{}, slots: map[common.Address]common.Address{}, baseFee: big.NewInt(100_000_000),
		l1Gas: 1_000_000}
}

func counterfactual(owner common.Address, salt *big.Int) common.Address {
	return common.BytesToAddress(crypto.Keccak256(owner.Bytes(), word(salt))[12:])
}

func (c *fakeChain) CodeAt(_ context.Context, account common.Address, _ *big.Int) ([]byte, error) {
	return c.code[account], nil
}

func (c *fakeChain) StorageAt(_ context.Context, account common.Address, key common.Hash, _ *big.Int) ([]byte, error) {
	if key != implementationSlot {
		return nil, errors.New("unexpected slot")
	}
	return common.LeftPadBytes(c.slots[account].Bytes(), 32), nil
}

func (c *fakeChain) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	if c.failHeader {
		return nil, errors.New("the node is down")
	}
	return &types.Header{BaseFee: c.baseFee}, nil
}

func (c *fakeChain) CallContract(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	paymasterABI := sponsorpaymaster.NewSponsorPaymaster()
	switch *msg.To {
	case factory:
		method, err := factoryABI.MethodById(msg.Data[:4])
		if err != nil {
			return nil, err
		}
		if method.Name == "accountImplementation" {
			return method.Outputs.Pack(implementation)
		}
		args, _ := method.Inputs.Unpack(msg.Data[4:])
		return method.Outputs.Pack(counterfactual(args[0].(common.Address), args[1].(*big.Int)))
	case paymaster:
		switch {
		case bytes.Equal(msg.Data, paymasterABI.PackSigner()):
			return common.LeftPadBytes(c.signer.Bytes(), 32), nil
		case bytes.Equal(msg.Data, paymasterABI.PackEntryPoint()):
			return common.LeftPadBytes(entryPoint.Bytes(), 32), nil
		default:
			sender := common.BytesToAddress(msg.Data[4:])
			return word(big.NewInt(c.used[sender])), nil
		}
	case nodeInterface:
		return nodeInterfaceABI.Methods["gasEstimateL1Component"].Outputs.Pack(c.l1Gas, c.baseFee, big.NewInt(1))
	}
	return nil, errors.New("unexpected call")
}

func registry(t *testing.T) *sdk.Deployments {
	t.Helper()
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,
		"tokens":{"USDG":"` + usdg.Hex() + `"},
		"tapehouse":{"MarginAccounts":"` + accounts.Hex() + `","Owner":"` + owner.Hex() + `","SponsorPaymaster":"` +
		paymaster.Hex() + `"},
		"erc4337":{"EntryPoint":"` + entryPoint.Hex() + `","SimpleAccountFactory":"` + factory.Hex() + `"}}`))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func service(t *testing.T, chain *fakeChain) *Service {
	t.Helper()
	s, err := New(context.Background(), chain, registry(t), signerKey, DefaultLimits,
		[]netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return now }
	return s
}

func big64(v int64) *hexutil.Big { return (*hexutil.Big)(big.NewInt(v)) }

func execute(to common.Address, value int64, data []byte) []byte {
	out, _ := accountABI.Pack("execute", to, big.NewInt(value), data)
	return out
}

func approve(spender common.Address) []byte {
	out, _ := erc20ABI.Pack("approve", spender, big.NewInt(1))
	return out
}

// creation is an operation that creates owner's account with callData and the gas a bundler estimates.
func creation(callData []byte) *UserOperation {
	factoryData, _ := factoryABI.Pack("createAccount", owner, big.NewInt(0))
	f := factory
	return &UserOperation{Sender: counterfactual(owner, big.NewInt(0)), Nonce: big64(0), Factory: &f,
		FactoryData: factoryData, CallData: callData, CallGasLimit: big64(150_000), VerificationGasLimit: big64(400_000),
		PreVerificationGas: big64(1_200_000), MaxFeePerGas: big64(240_000_000), MaxPriorityFeePerGas: big64(0),
		PaymasterVerificationGasLimit: big64(PaymasterVerificationGas), PaymasterPostOpGasLimit: big64(0),
		Signature: make([]byte, 65)}
}

// existing is an operation of owner's account once it exists.
func existing(chain *fakeChain, callData []byte) *UserOperation {
	op := creation(callData)
	op.Factory, op.FactoryData = nil, nil
	chain.code[op.Sender] = []byte{0x60}
	chain.slots[op.Sender] = implementation
	return op
}

func refused(t *testing.T, err error, reason string) {
	t.Helper()
	var refusal *Refusal
	if !errors.As(err, &refusal) || !strings.Contains(refusal.Reason, reason) {
		t.Fatalf("got %v, want a refusal with %q", err, reason)
	}
}

func TestTheTypedDataHashIsThePaymasters(t *testing.T) {
	f := common.HexToAddress("0x2222222222222222222222222222222222222222")
	op := &UserOperation{Sender: common.HexToAddress("0x1111111111111111111111111111111111111111"), Nonce: big64(7),
		Factory: &f, FactoryData: hexutil.MustDecode("0x5fbfb9cf"), CallData: hexutil.MustDecode("0xb61d27f6"),
		CallGasLimit: big64(150_000), VerificationGasLimit: big64(300_000), PreVerificationGas: big64(60_000),
		MaxFeePerGas: big64(200_000_000), MaxPriorityFeePerGas: big64(1),
		PaymasterVerificationGasLimit: big64(75_000), PaymasterPostOpGasLimit: big64(0)}
	domain := domainSeparator(412346, common.HexToAddress("0xe3BC0969E9C4AfcB1cF345C5F5E089F3766a8924"))
	got := typedDataHash(domain, op, 0, 1_800_000_000)
	if want := common.HexToHash("0x38cb6da73f6865edcae8f3ad68862e499d3a0734804a442c9c089931a9875cd8"); got != want {
		t.Fatalf("hash %s, want the paymaster's %s", got, want)
	}
}

func TestTheServiceStartsOnlyAsThePaymastersSigner(t *testing.T) {
	chain := newChain()
	chain.signer = stranger
	if _, err := New(context.Background(), chain, registry(t), signerKey, DefaultLimits, nil, slog.Default()); err == nil ||
		!strings.Contains(err.Error(), "not the paymaster's signer") {
		t.Fatalf("a foreign signer: %v", err)
	}
	d := registry(t)
	delete(d.Tapehouse, "SponsorPaymaster")
	if _, err := New(context.Background(), newChain(), d, signerKey, DefaultLimits, nil, slog.Default()); err == nil ||
		err.Error() != "the registry has no .tapehouse.SponsorPaymaster" {
		t.Fatalf("no paymaster: %v", err)
	}
	if service(t, newChain()).Signer() != crypto.PubkeyToAddress(signerKey.PublicKey) {
		t.Fatal("signer")
	}
}

func TestAnAccountsCreationAloneIsSponsored(t *testing.T) {
	s := service(t, newChain())
	stub, err := s.StubData(context.Background(), creation([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	if stub["paymaster"] != paymaster || stub["paymasterVerificationGasLimit"] != hexutil.Uint64(PaymasterVerificationGas) ||
		stub["isFinal"] != false || len(stub["paymasterData"].(hexutil.Bytes)) != 12+65 {
		t.Fatalf("stub %v", stub)
	}
}

func TestCallsToTapehouseAndApprovalsOfItAreSponsored(t *testing.T) {
	chain := newChain()
	s := service(t, chain)
	batch, _ := accountABI.Pack("executeBatch", []common.Address{usdg, accounts}, []*big.Int{},
		[][]byte{approve(accounts), {0x01, 0x02, 0x03, 0x04}})
	for _, callData := range [][]byte{execute(accounts, 0, nil), batch} {
		if _, err := s.StubData(context.Background(), existing(chain, callData)); err != nil {
			t.Fatal(err)
		}
	}
	chain.used[counterfactual(owner, big.NewInt(0))] = 2
	if _, err := s.Data(context.Background(), creation(execute(accounts, 0, nil))); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Data(context.Background(), existing(chain, execute(accounts, 0, nil))); err != nil {
		t.Fatal(err)
	}
}

func TestEverythingElseIsRefused(t *testing.T) {
	chain := newChain()
	s := service(t, chain)
	ctx := context.Background()
	otherFactory := creation([]byte{})
	otherFactory.Factory = &stranger
	refused(t, s.check(ctx, otherFactory), "not by the registry's SimpleAccountFactory")
	notCreate := creation([]byte{})
	notCreate.FactoryData = approve(owner)
	refused(t, s.check(ctx, notCreate), "not createAccount")
	otherSender := creation([]byte{})
	otherSender.Sender = stranger
	refused(t, s.check(ctx, otherSender), "not "+stranger.Hex())
	upgraded := existing(chain, execute(accounts, 0, nil))
	chain.slots[upgraded.Sender] = stranger
	refused(t, s.check(ctx, upgraded), "is not a SimpleAccount v0.7")
	emptyBatch, _ := accountABI.Pack("executeBatch", []common.Address{}, []*big.Int{}, [][]byte{})
	shortValues, _ := accountABI.Pack("executeBatch", []common.Address{accounts, accounts}, []*big.Int{big.NewInt(0)},
		[][]byte{nil, nil})
	for callData, reason := range map[string]string{
		string(execute(stranger, 0, nil)):                         "is not to a Tapehouse contract",
		string(execute(usdg, 0, approve(stranger))):               "is not to a Tapehouse contract",
		string(execute(stranger, 0, approve(accounts))):           "is not to a Tapehouse contract",
		string(execute(usdg, 0, approve(accounts)[:36])):          "is not to a Tapehouse contract",
		string(execute(accounts, 1, nil)):                         "sends value",
		string(emptyBatch):                                        "arrays do not match",
		string(shortValues):                                       "arrays do not match",
		"\x34\xfc\xd5\xbe":                                        "not execute or executeBatch",
		"\xb6\x1d":                                                "not execute or executeBatch",
		string(append(execute(accounts, 0, nil)[:4], 0x01, 0x02)): "does not decode",
	} {
		refused(t, s.check(ctx, existing(chain, []byte(callData))), reason)
	}
	refused(t, s.check(ctx, existing(chain, []byte{})), "neither creates the account nor makes a call")
	chain.used[counterfactual(owner, big.NewInt(0))] = FreeOperations
	refused(t, s.check(ctx, existing(chain, execute(accounts, 0, nil))), "has used its 3 free operations")
	if _, err := s.StubData(ctx, creation([]byte{})); err != nil {
		t.Fatalf("a creation alone is free whatever the count: %v", err)
	}
	refused(t, s.check(ctx, &UserOperation{}), "no nonce or call data")
}

func TestTheSignatureIsWhatThePaymasterChecks(t *testing.T) {
	s := service(t, newChain())
	op := creation(execute(accounts, 0, nil))
	data, err := s.Data(context.Background(), op)
	if err != nil {
		t.Fatal(err)
	}
	raw := data["paymasterData"].(hexutil.Bytes)
	validUntil := new(big.Int).SetBytes(raw[6:12]).Uint64()
	if new(big.Int).SetBytes(raw[:6]).Sign() != 0 || validUntil != uint64(now.Add(10*time.Minute).Unix()) {
		t.Fatalf("window %x", raw[:12])
	}
	signature := append([]byte{}, raw[12:]...)
	signature[64] -= 27
	key, err := crypto.SigToPub(typedDataHash(s.domain, op, 0, validUntil).Bytes(), signature)
	if err != nil || crypto.PubkeyToAddress(*key) != s.Signer() {
		t.Fatalf("the signature does not recover the signer: %v", err)
	}
}

func TestGasBeyondTheOperationsCostIsRefused(t *testing.T) {
	chain := newChain()
	s := service(t, chain)
	ctx := context.Background()
	for reason, raise := range map[string]func(op *UserOperation){
		"gas is not complete":                   func(op *UserOperation) { op.MaxFeePerGas = nil },
		"paymaster's gas limits":                func(op *UserOperation) { op.PaymasterVerificationGasLimit = big64(200_000) },
		"post-op":                               func(op *UserOperation) { op.PaymasterPostOpGasLimit = big64(1) },
		"verification gas limit is above":       func(op *UserOperation) { op.VerificationGasLimit = big64(1_000_001) },
		"call gas limit is above":               func(op *UserOperation) { op.CallGasLimit = big64(2_000_001) },
		"priority fee is above":                 func(op *UserOperation) { op.MaxPriorityFeePerGas = big64(10_000_001) },
		"fee cap is above 3 times the base fee": func(op *UserOperation) { op.MaxFeePerGas = big64(300_000_001) },
		"pre-verification gas is above 1400000": func(op *UserOperation) { op.PreVerificationGas = big64(1_400_001) },
	} {
		op := creation(execute(accounts, 0, nil))
		raise(op)
		_, err := s.Data(ctx, op)
		if reason == "post-op" {
			reason = "paymaster's gas limits"
		}
		refused(t, err, reason)
	}
	chain.failHeader = true
	if _, err := s.Data(ctx, creation(execute(accounts, 0, nil))); err == nil || err.Error() != "the node is down" {
		t.Fatalf("a node down: %v", err)
	}
}

func rpc(t *testing.T, s *Service, remote, forwarded string, method string, params ...any) map[string]any {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	r.RemoteAddr = remote + ":5000"
	if forwarded != "" {
		r.Header.Set("X-Forwarded-For", forwarded)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	var out map[string]any
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func messageOf(out map[string]any) string {
	if failure, ok := out["error"].(map[string]any); ok {
		return failure["message"].(string)
	}
	return ""
}

func TestTheServiceAnswersERC7677OverJSONRPC(t *testing.T) {
	s := service(t, newChain())
	chainID := hexutil.Uint64(412346)
	stub := rpc(t, s, "203.0.113.1", "", "pm_getPaymasterStubData", creation([]byte{}), entryPoint, chainID, map[string]any{})
	if result := stub["result"].(map[string]any); result["paymaster"] != strings.ToLower(paymaster.Hex()) &&
		result["paymaster"] != paymaster.Hex() {
		t.Fatalf("stub %v", stub)
	}
	data := rpc(t, s, "203.0.113.1", "", "pm_getPaymasterData", creation(execute(accounts, 0, nil)), entryPoint, chainID, nil)
	if _, ok := data["result"].(map[string]any)["paymasterData"]; !ok {
		t.Fatalf("data %v", data)
	}
	for want, out := range map[string]map[string]any{
		"the method eth_chainId does not exist": rpc(t, s, "203.0.113.1", "", "eth_chainId"),
		"the params are":                        rpc(t, s, "203.0.113.1", "", "pm_getPaymasterData", 1),
		"the params do not parse":               rpc(t, s, "203.0.113.1", "", "pm_getPaymasterData", 1, 2, 3),
		"on chain 412346":                       rpc(t, s, "203.0.113.1", "", "pm_getPaymasterData", creation([]byte{}), entryPoint, hexutil.Uint64(1)),
		"is not to a Tapehouse contract":        rpc(t, s, "203.0.113.1", "", "pm_getPaymasterStubData", creation(execute(stranger, 0, nil)), entryPoint, chainID),
	} {
		if !strings.Contains(messageOf(out), want) {
			t.Fatalf("got %v, want %q", out, want)
		}
	}
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), "the request is not JSON") {
		t.Fatalf("not JSON: %s", w.Body)
	}
	w = httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodOptions, "/", nil))
	if w.Code != http.StatusNoContent || w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("preflight %d %v", w.Code, w.Header())
	}
	w = httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET %d", w.Code)
	}
}

func TestEachClientGetsItsDailySignatures(t *testing.T) {
	chain := newChain()
	s := service(t, chain)
	s.limits.Daily = 2
	chainID := hexutil.Uint64(412346)
	op := creation(execute(accounts, 0, nil))
	sign := func(remote, forwarded string) string {
		return messageOf(rpc(t, s, remote, forwarded, "pm_getPaymasterData", op, entryPoint, chainID))
	}
	for range 2 {
		if sign("203.0.113.1", "") != "" {
			t.Fatal("the first two")
		}
	}
	if got := sign("203.0.113.1", ""); got != "the daily 2 sponsorships of this client are used" {
		t.Fatalf("the third: %q", got)
	}
	if sign("10.0.0.1", "198.51.100.7, 10.0.0.2") != "" {
		t.Fatal("another client through a trusted proxy")
	}
	for i := range 3 {
		if refused := sign("203.0.113.9", "198.51.100.7") != ""; refused != (i == 2) {
			t.Fatal("an untrusted proxy's header is ignored")
		}
	}
	now = now.Add(25 * time.Hour)
	defer func() { now = now.Add(-25 * time.Hour) }()
	if sign("203.0.113.1", "") != "" {
		t.Fatal("a day later")
	}
	chain.failHeader = true
	if got := sign("192.0.2.1", ""); got != "the chain could not be read" {
		t.Fatalf("a node down: %q", got)
	}
}
