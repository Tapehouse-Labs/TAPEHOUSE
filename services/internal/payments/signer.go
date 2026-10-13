// SPDX-License-Identifier: MIT OR Apache-2.0

// Package payments settles x402 payments in USDG: a signer through which the x402 facilitator reads the chain and
// sends transfers, the facilitator's HTTP service, and the paywall that sells the API's reads one request at a time.
package payments

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/x402-foundation/x402/go/v2/mechanisms/evm"
)

// Chain is the node's RPC as the signer uses it; *ethclient.Client is one.
type Chain interface {
	ChainID(ctx context.Context) (*big.Int, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, block *big.Int) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
}

const (
	erc1271Magic = "1626ba7e"
	gasMargin    = 120
	gasCeiling   = 5_000_000
	nonceHold    = 2 * time.Minute
	receiptEvery = 250 * time.Millisecond
	receiptLimit = 2 * time.Minute
)

var domainFields = []apitypes.Type{{Name: "name", Type: "string"}, {Name: "version", Type: "string"},
	{Name: "chainId", Type: "uint256"}, {Name: "verifyingContract", Type: "address"}}

const erc1271ABI = `[{"type":"function","name":"isValidSignature","stateMutability":"view","inputs":[{"name":"hash","type":"bytes32"},{"name":"signature","type":"bytes"}],"outputs":[{"name":"magicValue","type":"bytes4"}]}]`

const erc20BalanceABI = `[{"type":"function","name":"balanceOf","stateMutability":"view","inputs":[{"name":"account","type":"address"}],"outputs":[{"name":"","type":"uint256"}]}]`

// Signer is the x402 facilitator's way onto the chain: it reads, verifies typed data, and sends the transactions of
// one key, one at a time so that their nonces follow. Its nonce is the greater of the chain's and the one after its
// last, since a node a block behind answers with a nonce already used; that one is dropped once two minutes pass
// without the chain reaching it, as a transaction that was evicted leaves a gap nothing fills.
type Signer struct {
	chain   Chain
	key     *ecdsa.PrivateKey
	address common.Address
	chainID *big.Int
	mu      sync.Mutex
	next    uint64
	sentAt  time.Time
}

var _ evm.FacilitatorEvmSigner = (*Signer)(nil)

// NewSigner returns a signer of key on the chain.
func NewSigner(ctx context.Context, chain Chain, key *ecdsa.PrivateKey) (*Signer, error) {
	chainID, err := chain.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("the chain's id: %w", err)
	}
	return &Signer{chain: chain, key: key, address: crypto.PubkeyToAddress(key.PublicKey), chainID: chainID}, nil
}

// Address is the account the signer sends from, which pays the gas of every settlement.
func (s *Signer) Address() common.Address {
	return s.address
}

// GetAddresses returns the account the signer sends from.
func (s *Signer) GetAddresses() []string {
	return []string{s.address.Hex()}
}

// GetChainID returns the chain's id.
func (s *Signer) GetChainID(context.Context) (*big.Int, error) {
	return new(big.Int).Set(s.chainID), nil
}

// GetCode returns the code at address, empty for an account without any.
func (s *Signer) GetCode(ctx context.Context, address string) ([]byte, error) {
	return s.chain.CodeAt(ctx, common.HexToAddress(address), nil)
}

// GetBalance returns address's balance of the ERC-20 token.
func (s *Signer) GetBalance(ctx context.Context, address, token string) (*big.Int, error) {
	out, err := s.ReadContract(ctx, token, []byte(erc20BalanceABI), "balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	balance, ok := out.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("balanceOf answered %T", out)
	}
	return balance, nil
}

// ReadContract calls a view of the contract at address and returns its output, or its outputs when it has several.
func (s *Signer) ReadContract(ctx context.Context, address string, abiJSON []byte, function string, args ...any) (any, error) {
	parsed, err := abi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}
	data, err := parsed.Pack(function, args...)
	if err != nil {
		return nil, err
	}
	to := common.HexToAddress(address)
	raw, err := s.chain.CallContract(ctx, ethereum.CallMsg{From: s.address, To: &to, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	out, err := parsed.Unpack(function, raw)
	if err != nil {
		return nil, err
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return out, nil
}

// VerifyTypedData reports whether signature is address's over the EIP-712 message: an ECDSA signature of the account
// itself, or, where the account is a contract, its answer to ERC-1271.
func (s *Signer) VerifyTypedData(ctx context.Context, address string, domain evm.TypedDataDomain,
	fields map[string][]evm.TypedDataField, primaryType string, message map[string]any, signature []byte) (bool, error) {
	typed := apitypes.TypedData{
		Types: apitypes.Types{"EIP712Domain": {{Name: "name", Type: "string"}, {Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"}, {Name: "verifyingContract", Type: "address"}}},
		PrimaryType: primaryType,
		Domain: apitypes.TypedDataDomain{Name: domain.Name, Version: domain.Version,
			ChainId: (*math.HexOrDecimal256)(domain.ChainID), VerifyingContract: domain.VerifyingContract},
		Message: message,
	}
	for name, list := range fields {
		if name == "EIP712Domain" {
			continue
		}
		for _, field := range list {
			typed.Types[name] = append(typed.Types[name], apitypes.Type{Name: field.Name, Type: field.Type})
		}
	}
	digest, _, err := apitypes.TypedDataAndHash(typed)
	if err != nil {
		return false, err
	}
	account := common.HexToAddress(address)
	if len(signature) == crypto.SignatureLength {
		rs := bytes.Clone(signature)
		if rs[64] >= 27 {
			rs[64] -= 27
		}
		if public, err := crypto.SigToPub(digest, rs); err == nil && crypto.PubkeyToAddress(*public) == account {
			return true, nil
		}
	}
	code, err := s.chain.CodeAt(ctx, account, nil)
	if err != nil || len(code) == 0 {
		return false, err
	}
	var hash [32]byte
	copy(hash[:], digest)
	out, err := s.ReadContract(ctx, address, []byte(erc1271ABI), "isValidSignature", hash, signature)
	if err != nil {
		return false, nil
	}
	magic, ok := out.([4]byte)
	return ok && hexutil.Encode(magic[:])[2:] == erc1271Magic, nil
}

// WriteContract sends a call of the contract at address, with dataSuffix appended to its calldata, and returns the
// transaction's hash once it is sent.
func (s *Signer) WriteContract(ctx context.Context, address string, abiJSON []byte, function string, dataSuffix []byte,
	args ...any) (string, error) {
	parsed, err := abi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		return "", err
	}
	data, err := parsed.Pack(function, args...)
	if err != nil {
		return "", err
	}
	return s.SendTransaction(ctx, address, append(data, dataSuffix...))
}

// SendTransaction sends data to the account at to and returns the transaction's hash.
func (s *Signer) SendTransaction(ctx context.Context, to string, data []byte) (string, error) {
	target := common.HexToAddress(to)
	call := ethereum.CallMsg{From: s.address, To: &target, Data: data}
	gas, err := s.chain.EstimateGas(ctx, call)
	if err != nil {
		return "", fmt.Errorf("the transfer would revert: %w", err)
	}
	if gas > gasCeiling {
		return "", fmt.Errorf("the transfer needs %d gas, over the %d a settlement may use", gas, gasCeiling)
	}
	head, err := s.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}
	if head.BaseFee == nil {
		return "", errors.New("the chain has no base fee")
	}
	tip, err := s.chain.SuggestGasTipCap(ctx)
	if err != nil {
		return "", err
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(head.BaseFee, big.NewInt(2)), tip)

	s.mu.Lock()
	defer s.mu.Unlock()
	pending, err := s.chain.PendingNonceAt(ctx, s.address)
	if err != nil {
		return "", err
	}
	if s.next > pending && time.Since(s.sentAt) > nonceHold {
		s.next = 0
	}
	nonce := max(pending, s.next)
	tx, err := types.SignNewTx(s.key, types.LatestSignerForChainID(s.chainID), &types.DynamicFeeTx{
		ChainID: s.chainID, Nonce: nonce, GasTipCap: tip, GasFeeCap: feeCap, Gas: gas * gasMargin / 100, To: &target,
		Data: data})
	if err != nil {
		return "", err
	}
	if err := s.chain.SendTransaction(ctx, tx); err != nil {
		return "", err
	}
	s.next, s.sentAt = nonce+1, time.Now()
	return tx.Hash().Hex(), nil
}

// WaitForTransactionReceipt waits, for two minutes at most, until the transaction is mined.
func (s *Signer) WaitForTransactionReceipt(ctx context.Context, txHash string) (*evm.TransactionReceipt, error) {
	ctx, cancel := context.WithTimeout(ctx, receiptLimit)
	defer cancel()
	tick := time.NewTicker(receiptEvery)
	defer tick.Stop()
	for {
		receipt, err := s.chain.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err == nil {
			logs := receipt.Logs
			if logs == nil {
				logs = []*types.Log{}
			}
			return &evm.TransactionReceipt{Status: receipt.Status, BlockNumber: receipt.BlockNumber.Uint64(),
				TxHash: receipt.TxHash.Hex(), Logs: logs}, nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-tick.C:
		}
	}
}
