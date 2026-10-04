// SPDX-License-Identifier: MIT OR Apache-2.0

// Package sponsor is Tapehouse's sponsor service. It decides which ERC-4337 user operations the sponsor paymaster pays
// for and signs each one it sponsors, serving ERC-7677's pm_getPaymasterStubData and pm_getPaymasterData over JSON-RPC.
// It sponsors an operation of a SimpleAccount v0.7 made by the registry's factory when it only creates the account, or
// when every call is a user's call to a Tapehouse contract of the registry, approves one to spend a token of the
// registry, or permits or transfers a token of the registry to the account beside such a call, and sends no value;
// while the account has free operations left on the paymaster; at gas limits, fees and a pre-verification gas close to
// what the operation costs; and within a daily number of signatures for each client and for all of them.
package sponsor

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/program"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/sponsorpaymaster"
	"golang.org/x/time/rate"
)

const (
	// PaymasterVerificationGas is the paymaster's verification gas limit in every operation the service signs.
	PaymasterVerificationGas = 75_000
	// FreeOperations is how many operations the paymaster sponsors for each account, besides its creation.
	FreeOperations = 3
)

// Limits are what the service accepts of an operation's gas, and how many signatures it gives each client a day.
type Limits struct {
	// MaxVerificationGas bounds the account's verification gas limit.
	MaxVerificationGas uint64
	// MaxCreationCallGas bounds the call gas limit of an operation that only creates its account, and MaxCallGas that of
	// one that also creates it, whose calls cannot be estimated before. An existing account's call gas limit is bounded
	// by 160% of its calls' estimate, plus EstimateMargin.
	MaxCreationCallGas, MaxCallGas, EstimateMargin uint64
	// MaxPriorityFee bounds the priority fee, in wei per gas.
	MaxPriorityFee *big.Int
	// FeeMultiple bounds the fee cap to this many times the latest base fee.
	FeeMultiple int64
	// StaticPreVerificationGas is the pre-verification gas allowed beyond 130% of the operation's L1 gas, as
	// Arbitrum's NodeInterface estimates it.
	StaticPreVerificationGas uint64
	// Validity is how long a signature stays valid.
	Validity time.Duration
	// Daily is how many signatures each client gets a day, a client being an IPv4 address or an IPv6 /64, and
	// GlobalDaily how many all clients get together.
	Daily, GlobalDaily int
}

// DefaultLimits are the service's limits unless configured otherwise.
var DefaultLimits = Limits{MaxVerificationGas: 1_000_000, MaxCreationCallGas: 50_000, MaxCallGas: 2_000_000,
	EstimateMargin: 20_000, MaxPriorityFee: new(big.Int), FeeMultiple: 3, StaticPreVerificationGas: 60_000,
	Validity: 10 * time.Minute, Daily: 30, GlobalDaily: 5_000}

// Chain is what the service reads of the chain. *ethclient.Client implements it.
type Chain interface {
	CodeAt(ctx context.Context, account common.Address, block *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, block *big.Int) ([]byte, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// UserOperation is an EntryPoint v0.7 user operation as ERC-4337's JSON-RPC carries it.
type UserOperation struct {
	Sender                        common.Address  `json:"sender"`
	Nonce                         *hexutil.Big    `json:"nonce"`
	Factory                       *common.Address `json:"factory,omitempty"`
	FactoryData                   hexutil.Bytes   `json:"factoryData,omitempty"`
	CallData                      hexutil.Bytes   `json:"callData"`
	CallGasLimit                  *hexutil.Big    `json:"callGasLimit"`
	VerificationGasLimit          *hexutil.Big    `json:"verificationGasLimit"`
	PreVerificationGas            *hexutil.Big    `json:"preVerificationGas"`
	MaxFeePerGas                  *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          *hexutil.Big    `json:"maxPriorityFeePerGas"`
	PaymasterVerificationGasLimit *hexutil.Big    `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       *hexutil.Big    `json:"paymasterPostOpGasLimit,omitempty"`
	Signature                     hexutil.Bytes   `json:"signature"`
}

// Service decides and signs sponsorships. It serves ERC-7677 as an http.Handler.
type Service struct {
	chain          Chain
	key            *ecdsa.PrivateKey
	chainID        uint64
	entryPoint     common.Address
	factory        common.Address
	paymaster      common.Address
	implementation common.Address
	proxyCodeHash  common.Hash
	domain         common.Hash
	targets        map[common.Address]map[[4]byte]bool
	tokens         map[common.Address]bool
	backstop       common.Address
	limits         Limits
	proxies        []netip.Prefix
	log            *slog.Logger
	now            func() time.Time

	mu      sync.Mutex
	clients map[string]*client
	global  *rate.Limiter
	swept   time.Time
}

// requestTimeout bounds the chain reads a request makes.
const requestTimeout = 20 * time.Second

// maxClients bounds the clients the service remembers; beyond it, a new client waits for the next sweep.
const maxClients = 100_000

type client struct {
	limiter, stubs *rate.Limiter
	seen           time.Time
}

// New returns the service of the registry d, signing with key, which must be the paymaster's signer. It trusts
// X-Forwarded-For from proxies.
func New(ctx context.Context, chain Chain, d *sdk.Deployments, key *ecdsa.PrivateKey, limits Limits,
	proxies []netip.Prefix, log *slog.Logger) (*Service, error) {
	s := &Service{chain: chain, key: key, chainID: d.ChainID, limits: limits, proxies: proxies, log: log, now: time.Now,
		clients: map[string]*client{}, tokens: map[common.Address]bool{},
		global: rate.NewLimiter(rate.Every(24*time.Hour/time.Duration(limits.GlobalDaily)), limits.GlobalDaily)}
	for _, entry := range []struct {
		into  *common.Address
		group map[string]common.Address
		name  string
	}{
		{&s.entryPoint, d.ERC4337, ".erc4337.EntryPoint"},
		{&s.factory, d.ERC4337, ".erc4337.SimpleAccountFactory"},
		{&s.paymaster, d.Tapehouse, ".tapehouse.SponsorPaymaster"},
	} {
		address, ok := entry.group[entry.name[strings.LastIndex(entry.name, ".")+1:]]
		if !ok {
			return nil, fmt.Errorf("the registry has no %s", entry.name)
		}
		*entry.into = address
	}
	var err error
	if s.targets, err = sponsoredFunctions(d); err != nil {
		return nil, err
	}
	for _, address := range d.Tokens {
		s.tokens[address] = true
	}
	for _, address := range d.Baskets {
		s.tokens[address] = true
	}
	s.backstop = d.Tapehouse["GapBackstop"]
	if s.implementation, err = s.address(ctx, s.factory, factoryABI, "accountImplementation"); err != nil {
		return nil, fmt.Errorf("the factory's account implementation: %w", err)
	}
	if s.proxyCodeHash, err = s.accountCodeHash(ctx); err != nil {
		return nil, fmt.Errorf("the factory's account code: %w", err)
	}
	paymaster := sponsorpaymaster.NewSponsorPaymaster()
	signer, err := s.call(ctx, s.paymaster, paymaster.PackSigner())
	if err != nil {
		return nil, fmt.Errorf("the paymaster's signer: %w", err)
	}
	if got, _ := paymaster.UnpackSigner(signer); got != crypto.PubkeyToAddress(key.PublicKey) {
		return nil, fmt.Errorf("the key is not the paymaster's signer, %s", got.Hex())
	}
	entryPoint, err := s.call(ctx, s.paymaster, paymaster.PackEntryPoint())
	if err != nil {
		return nil, fmt.Errorf("the paymaster's EntryPoint: %w", err)
	}
	if got, _ := paymaster.UnpackEntryPoint(entryPoint); got != s.entryPoint {
		return nil, fmt.Errorf("the paymaster is on another EntryPoint, %s", got.Hex())
	}
	s.domain = domainSeparator(d.ChainID, s.paymaster)
	return s, nil
}

// Signer is the address the service signs with.
func (s *Service) Signer() common.Address {
	return crypto.PubkeyToAddress(s.key.PublicKey)
}

// Refusal is why the service refuses to sponsor an operation.
type Refusal struct{ Reason string }

func (r *Refusal) Error() string { return r.Reason }

func refuse(format string, args ...any) error {
	return &Refusal{Reason: fmt.Sprintf(format, args...)}
}

// StubData answers pm_getPaymasterStubData: the paymaster fields to estimate the operation's gas with, once the
// operation is one the service sponsors.
func (s *Service) StubData(ctx context.Context, op *UserOperation) (map[string]any, error) {
	if err := s.check(ctx, op); err != nil {
		return nil, err
	}
	return map[string]any{
		"paymaster":                     s.paymaster,
		"paymasterData":                 hexutil.Bytes(paymasterData(0, 0, make([]byte, 65))),
		"paymasterVerificationGasLimit": hexutil.Uint64(PaymasterVerificationGas),
		"paymasterPostOpGasLimit":       hexutil.Uint64(0),
		"sponsor":                       map[string]string{"name": "Tapehouse"},
		"isFinal":                       false,
	}, nil
}

// Data answers pm_getPaymasterData: the paymaster's signed data for the operation, once its gas is within the limits.
func (s *Service) Data(ctx context.Context, op *UserOperation) (map[string]any, error) {
	if err := s.check(ctx, op); err != nil {
		return nil, err
	}
	if err := s.checkGas(ctx, op); err != nil {
		return nil, err
	}
	validUntil := uint64(s.now().Add(s.limits.Validity).Unix())
	signature, err := crypto.Sign(typedDataHash(s.domain, op, 0, validUntil).Bytes(), s.key)
	if err != nil {
		return nil, err
	}
	signature[64] += 27
	return map[string]any{"paymaster": s.paymaster, "paymasterData": hexutil.Bytes(paymasterData(0, validUntil, signature))}, nil
}

// check refuses an operation the paymaster does not sponsor, whatever its gas.
func (s *Service) check(ctx context.Context, op *UserOperation) error {
	if op.Nonce == nil || op.CallData == nil {
		return refuse("the operation has no nonce or call data")
	}
	if err := s.checkAccount(ctx, op); err != nil {
		return err
	}
	creates := op.Factory != nil
	if len(op.CallData) == 0 {
		if !creates {
			return refuse("the operation neither creates the account nor makes a call")
		}
		return nil
	}
	calls, err := decodeCalls(op.CallData)
	if err != nil {
		return err
	}
	moves, tapehouse := false, false
	claimed := map[common.Address]bool{}
	var holder *common.Address
	for _, c := range calls {
		if s.tokens[c.to] && len(c.data) >= 4 && [4]byte(c.data[:4]) == transferSelector {
			if err := s.checkTransfer(ctx, op, c, claimed, &holder); err != nil {
				return err
			}
			moves = true
			continue
		}
		if err := s.checkCall(op.Sender, c); err != nil {
			return err
		}
		if len(c.data) >= 4 {
			selector := [4]byte(c.data[:4])
			moves = moves || (s.tokens[c.to] && (selector == permitSelector || selector == transferFromSelector))
			tapehouse = tapehouse || s.targets[c.to][selector]
			if c.to == s.backstop && selector == claimGainsSelector {
				if args, err := unpack(backstopABI, "claimGains", c.data); err == nil {
					claimed[args[0].(common.Address)] = true
				}
			}
		}
	}
	if moves && !tapehouse {
		return refuse("the operation moves tokens without a call to Tapehouse")
	}
	used, err := s.call(ctx, s.paymaster, sponsorpaymaster.NewSponsorPaymaster().PackOperations(op.Sender))
	if err != nil {
		return err
	}
	if new(big.Int).SetBytes(used).Cmp(big.NewInt(FreeOperations)) >= 0 {
		return refuse("the account %s has used its %d free operations", op.Sender.Hex(), FreeOperations)
	}
	return nil
}

// checkAccount refuses an account the factory does not make, or one whose implementation is not the factory's.
func (s *Service) checkAccount(ctx context.Context, op *UserOperation) error {
	if op.Factory != nil {
		if *op.Factory != s.factory {
			return refuse("the account is made by %s, not by the registry's SimpleAccountFactory", op.Factory.Hex())
		}
		args, err := unpack(factoryABI, "createAccount", op.FactoryData)
		if err != nil {
			return refuse("the factory data is not createAccount(owner, salt)")
		}
		sender, err := s.address(ctx, s.factory, factoryABI, "getAddress", args...)
		if err != nil {
			return err
		}
		if sender != op.Sender {
			return refuse("the factory makes %s for that owner and salt, not %s", sender.Hex(), op.Sender.Hex())
		}
		return nil
	}
	code, err := s.chain.CodeAt(ctx, op.Sender, nil)
	if err != nil {
		return err
	}
	slot, err := s.chain.StorageAt(ctx, op.Sender, implementationSlot, nil)
	if err != nil {
		return err
	}
	if crypto.Keccak256Hash(code) != s.proxyCodeHash || common.BytesToAddress(slot) != s.implementation {
		return refuse("the account %s is not a SimpleAccount v0.7 of the registry's factory", op.Sender.Hex())
	}
	return nil
}

// checkCall refuses a call of sender that is not a user's call to a Tapehouse contract, an approval of one to spend a
// registry token, nor a permit or a transfer of a registry token to sender, by which it takes what its owner holds.
func (s *Service) checkCall(sender common.Address, c call) error {
	if c.value.Sign() != 0 {
		return refuse("the call to %s sends value", c.to.Hex())
	}
	if len(c.data) >= 4 {
		selector := [4]byte(c.data[:4])
		if s.targets[c.to][selector] {
			return nil
		}
		if s.tokens[c.to] {
			switch selector {
			case approveSelector:
				if args, err := unpack(erc20ABI, "approve", c.data); err == nil && s.targets[args[0].(common.Address)] != nil {
					return nil
				}
			case permitSelector:
				if args, err := unpack(erc20ABI, "permit", c.data); err == nil && args[1].(common.Address) == sender {
					return nil
				}
			case transferFromSelector:
				if args, err := unpack(erc20ABI, "transferFrom", c.data); err == nil && args[1].(common.Address) == sender {
					return nil
				}
			}
		}
	}
	return refuse("the call to %s is not a user's call to a Tapehouse contract, nor an approval of one", c.to.Hex())
}

// checkTransfer refuses a transfer of a registry token unless it sends op's account owner a token the account claimed
// from the backstop's gains earlier in the same operation, once per claim. The owner is read once an operation, into
// holder.
func (s *Service) checkTransfer(ctx context.Context, op *UserOperation, c call, claimed map[common.Address]bool,
	holder **common.Address) error {
	args, err := unpack(erc20ABI, "transfer", c.data)
	if err != nil || !claimed[c.to] {
		return refuse("the operation sends a token it did not just claim")
	}
	delete(claimed, c.to)
	if *holder == nil {
		owner, err := s.owner(ctx, op)
		if err != nil {
			return err
		}
		*holder = &owner
	}
	if args[0].(common.Address) != **holder {
		return refuse("the operation sends a token to someone other than its owner")
	}
	return nil
}

// owner is the owner of op's account: the factory data's while the operation creates it, else the account's own.
func (s *Service) owner(ctx context.Context, op *UserOperation) (common.Address, error) {
	if op.Factory != nil {
		args, err := unpack(factoryABI, "createAccount", op.FactoryData)
		if err != nil {
			return common.Address{}, refuse("the factory data is not createAccount(owner, salt)")
		}
		return args[0].(common.Address), nil
	}
	return s.address(ctx, op.Sender, accountABI, "owner")
}

// accountCodeHash is the hash of the code of the factory's accounts, ERC-1967 proxies with no immutable: the code
// returned by a creation, simulated with eth_call, that has the factory create an account and returns its code.
func (s *Service) accountCodeHash(ctx context.Context) (common.Hash, error) {
	create, err := factoryABI.Pack("createAccount", common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
		new(big.Int))
	if err != nil {
		return common.Hash{}, err
	}
	probe := program.New().Mstore(create, 0).Call(nil, s.factory, 0, 0, len(create), 0, 32).Op(vm.POP).
		Push(0).Op(vm.MLOAD, vm.DUP1, vm.EXTCODESIZE, vm.DUP1).Push(0).Push(32).Op(vm.DUP5, vm.EXTCODECOPY).
		Push(32).Op(vm.RETURN)
	code, err := s.chain.CallContract(ctx, ethereum.CallMsg{Data: probe.Bytes()}, nil)
	if err != nil {
		return common.Hash{}, err
	}
	if len(code) == 0 {
		return common.Hash{}, errors.New("the factory created no account")
	}
	return crypto.Keccak256Hash(code), nil
}

// checkGas refuses gas limits, fees or a pre-verification gas beyond what the operation costs.
func (s *Service) checkGas(ctx context.Context, op *UserOperation) error {
	for _, field := range []*hexutil.Big{op.CallGasLimit, op.VerificationGasLimit, op.PreVerificationGas,
		op.MaxFeePerGas, op.MaxPriorityFeePerGas, op.PaymasterVerificationGasLimit, op.PaymasterPostOpGasLimit} {
		if field == nil {
			return refuse("the operation's gas is not complete")
		}
	}
	if op.PaymasterVerificationGasLimit.ToInt().Cmp(big.NewInt(PaymasterVerificationGas)) != 0 ||
		op.PaymasterPostOpGasLimit.ToInt().Sign() != 0 {
		return refuse("the paymaster's gas limits are not %d and 0", PaymasterVerificationGas)
	}
	if !op.VerificationGasLimit.ToInt().IsUint64() || op.VerificationGasLimit.ToInt().Uint64() > s.limits.MaxVerificationGas {
		return refuse("the verification gas limit is above %d", s.limits.MaxVerificationGas)
	}
	callGas, err := s.callGasLimit(ctx, op)
	if err != nil {
		return err
	}
	if !op.CallGasLimit.ToInt().IsUint64() || op.CallGasLimit.ToInt().Uint64() > callGas {
		return refuse("the call gas limit is above %d", callGas)
	}
	if op.MaxPriorityFeePerGas.ToInt().Cmp(s.limits.MaxPriorityFee) > 0 {
		return refuse("the priority fee is above %s wei", s.limits.MaxPriorityFee)
	}
	head, err := s.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	if head.BaseFee == nil {
		return errors.New("the chain reports no base fee")
	}
	if limit := new(big.Int).Mul(head.BaseFee, big.NewInt(s.limits.FeeMultiple)); op.MaxFeePerGas.ToInt().Cmp(limit) > 0 {
		return refuse("the fee cap is above %d times the base fee", s.limits.FeeMultiple)
	}
	l1, err := s.l1Gas(ctx, op)
	if err != nil {
		return err
	}
	limit := l1*13/10 + s.limits.StaticPreVerificationGas
	if !op.PreVerificationGas.ToInt().IsUint64() || op.PreVerificationGas.ToInt().Uint64() > limit {
		return refuse("the pre-verification gas is above %d", limit)
	}
	return nil
}

// callGasLimit is the most call gas the operation may ask for: for an existing account, 160% of its calls' L2 gas, sent
// by the EntryPoint as a transaction of their own, which Arbitrum's NodeInterface splits from the L1 gas the
// pre-verification gas pays, plus EstimateMargin. The estimate is net of storage refunds, which hide up to a quarter of
// the gas a call that clears storage, such as a short's full buy-back, must hold at its peak.
func (s *Service) callGasLimit(ctx context.Context, op *UserOperation) (uint64, error) {
	switch {
	case len(op.CallData) == 0:
		return s.limits.MaxCreationCallGas, nil
	case op.Factory != nil:
		return s.limits.MaxCallGas, nil
	}
	data, err := nodeInterfaceABI.Pack("gasEstimateComponents", op.Sender, false, []byte(op.CallData))
	if err != nil {
		return 0, err
	}
	out, err := s.chain.CallContract(ctx, ethereum.CallMsg{From: s.entryPoint, To: &nodeInterface, Data: data}, nil)
	if err != nil {
		return 0, refuse("the calls fail: %v", err)
	}
	values, err := nodeInterfaceABI.Unpack("gasEstimateComponents", out)
	if err != nil {
		return 0, err
	}
	estimate := values[0].(uint64) - values[1].(uint64)
	return min(estimate*8/5+s.limits.EstimateMargin, s.limits.MaxCallGas), nil
}

// l1Gas is the L1 gas of handleOps carrying the operation, signed by the account and the service, as Arbitrum's
// NodeInterface estimates it. A SimpleAccount's signature is 65 bytes, whatever the operation carries yet.
func (s *Service) l1Gas(ctx context.Context, op *UserOperation) (uint64, error) {
	signed := *op
	signed.Signature = make([]byte, 65)
	handleOps, err := entryPointABI.Pack("handleOps", []packedUserOperation{pack(&signed, paymasterAndData(s.paymaster,
		paymasterData(0, 0, make([]byte, 65))))}, common.Address{})
	if err != nil {
		return 0, err
	}
	data, err := nodeInterfaceABI.Pack("gasEstimateL1Component", s.entryPoint, false, handleOps)
	if err != nil {
		return 0, err
	}
	out, err := s.call(ctx, nodeInterface, data)
	if err != nil {
		return 0, fmt.Errorf("NodeInterface's L1 estimate: %w", err)
	}
	values, err := nodeInterfaceABI.Unpack("gasEstimateL1Component", out)
	if err != nil {
		return 0, err
	}
	return values[0].(uint64), nil
}

// allowStub spends one of the client's stubs, a second each with bursts of 20, or reports that none is left.
func (s *Service) allowStub(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.remember(address)
	return c != nil && c.stubs.AllowN(s.now(), 1)
}

// allow spends one of the client's daily signatures and one of all clients', or reports that none is left.
func (s *Service) allow(address string) bool {
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.remember(address)
	if c == nil || !c.limiter.AllowN(now, 1) {
		return false
	}
	if !s.global.AllowN(now, 1) {
		c.limiter.ReserveN(now, -1)
		return false
	}
	return true
}

// remember returns the client at address, forgetting those unseen for a day, or nil when it is new and the service
// remembers as many as it may. The caller holds s.mu.
func (s *Service) remember(address string) *client {
	now := s.now()
	if now.Sub(s.swept) > time.Minute {
		for key, c := range s.clients {
			if now.Sub(c.seen) > 24*time.Hour {
				delete(s.clients, key)
			}
		}
		s.swept = now
	}
	c, ok := s.clients[address]
	if !ok {
		if len(s.clients) >= maxClients {
			return nil
		}
		c = &client{limiter: rate.NewLimiter(rate.Every(24*time.Hour/time.Duration(s.limits.Daily)), s.limits.Daily),
			stubs: rate.NewLimiter(1, 20)}
		s.clients[address] = c
	}
	c.seen = now
	return c
}

type request struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      json.RawMessage   `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServeHTTP serves ERC-7677 over JSON-RPC 2.0, to any origin.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "POST a JSON-RPC request", http.StatusMethodNotAllowed)
		return
	}
	var req request
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&req); err != nil {
		reply(w, nil, nil, &rpcError{Code: -32700, Message: "the request is not JSON"})
		return
	}
	result, failure := s.serve(r, &req)
	reply(w, req.ID, result, failure)
}

func (s *Service) serve(r *http.Request, req *request) (any, *rpcError) {
	if req.Method != "pm_getPaymasterStubData" && req.Method != "pm_getPaymasterData" {
		return nil, &rpcError{Code: -32601, Message: fmt.Sprintf("the method %s does not exist", req.Method)}
	}
	if len(req.Params) < 3 {
		return nil, &rpcError{Code: -32602, Message: "the params are [userOperation, entryPoint, chainId, context]"}
	}
	var op UserOperation
	var entryPoint common.Address
	var chainID hexutil.Uint64
	if json.Unmarshal(req.Params[0], &op) != nil || json.Unmarshal(req.Params[1], &entryPoint) != nil ||
		json.Unmarshal(req.Params[2], &chainID) != nil {
		return nil, &rpcError{Code: -32602, Message: "the params do not parse"}
	}
	if entryPoint != s.entryPoint || uint64(chainID) != s.chainID {
		return nil, &rpcError{Code: -32602, Message: fmt.Sprintf("the service sponsors EntryPoint %s on chain %d",
			s.entryPoint.Hex(), s.chainID)}
	}
	var result map[string]any
	var err error
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()
	switch {
	case req.Method == "pm_getPaymasterStubData" && !s.allowStub(s.client(r)):
		return nil, &rpcError{Code: -32005, Message: "this client asks too often"}
	case req.Method == "pm_getPaymasterStubData":
		result, err = s.StubData(ctx, &op)
	case !s.allow(s.client(r)):
		return nil, &rpcError{Code: -32005, Message: "the day's sponsorships of this client, or of all clients, are used"}
	default:
		result, err = s.Data(ctx, &op)
	}
	var refusal *Refusal
	switch {
	case errors.As(err, &refusal):
		return nil, &rpcError{Code: -32000, Message: refusal.Reason}
	case err != nil:
		s.log.Error("a sponsorship could not be decided", "method", req.Method, "error", err)
		return nil, &rpcError{Code: -32603, Message: "the chain could not be read"}
	}
	if req.Method == "pm_getPaymasterData" {
		s.log.Info("sponsored", "sender", op.Sender.Hex(), "nonce", op.Nonce.String())
	}
	return result, nil
}

// client is the request's client: its remote address, or, through a trusted proxy, the last address in
// X-Forwarded-For that is not a trusted proxy; an IPv6 address stands for its /64.
func (s *Service) client(r *http.Request) string {
	remote, err := parseHop(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	if s.trusted(remote) {
		hops := strings.Split(strings.Join(r.Header.Values("X-Forwarded-For"), ","), ",")
		for i := len(hops) - 1; i >= 0; i-- {
			hop, err := parseHop(strings.TrimSpace(hops[i]))
			if err != nil {
				break
			}
			if !s.trusted(hop) {
				remote = hop
				break
			}
		}
	}
	if remote.Is6() {
		prefix, _ := remote.Prefix(64)
		return prefix.String()
	}
	return remote.String()
}

// parseHop reads an address, with or without its port.
func parseHop(text string) (netip.Addr, error) {
	if host, _, err := net.SplitHostPort(text); err == nil {
		text = host
	}
	address, err := netip.ParseAddr(text)
	return address.Unmap(), err
}

func (s *Service) trusted(address netip.Addr) bool {
	for _, prefix := range s.proxies {
		if prefix.Contains(address.Unmap()) {
			return true
		}
	}
	return false
}

func reply(w http.ResponseWriter, id json.RawMessage, result any, failure *rpcError) {
	if id == nil {
		id = json.RawMessage("null")
	}
	body := map[string]any{"jsonrpc": "2.0", "id": id}
	if failure != nil {
		body["error"] = failure
	} else {
		body["result"] = result
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func (s *Service) call(ctx context.Context, to common.Address, data []byte) ([]byte, error) {
	return s.chain.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
}
