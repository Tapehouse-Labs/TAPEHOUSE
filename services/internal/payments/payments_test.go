// SPDX-License-Identifier: MIT OR Apache-2.0

package payments_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/tapehouse/tapehouse/services/internal/payments"
	"github.com/x402-foundation/x402/go/v2/mechanisms/evm"
)

var (
	usdg    = common.HexToAddress("0x7E955252E15c84f5768B83c41a71F9eba181802F")
	payTo   = common.HexToAddress("0x8A2631c8226D7EA5612798e6E8C34F4DC703d0ac")
	chainID = big.NewInt(46630)
)

// chain is a double of the node: it answers each call by the function it names and holds the code of some accounts.
type chain struct {
	answers map[string]func(args []any) ([]any, error)
	code    map[common.Address][]byte
	sent    []*types.Transaction
	nonce   uint64
}

func (c *chain) ChainID(context.Context) (*big.Int, error) { return chainID, nil }

func (c *chain) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	for name, answer := range c.answers {
		method, ok := tokenABI.Methods[name]
		if !ok || len(call.Data) < 4 || string(call.Data[:4]) != string(method.ID) {
			continue
		}
		args, err := method.Inputs.Unpack(call.Data[4:])
		if err != nil {
			return nil, err
		}
		out, err := answer(args)
		if err != nil {
			return nil, err
		}
		return method.Outputs.Pack(out...)
	}
	return nil, ethereum.NotFound
}

func (c *chain) CodeAt(_ context.Context, account common.Address, _ *big.Int) ([]byte, error) {
	return c.code[account], nil
}

func (c *chain) PendingNonceAt(context.Context, common.Address) (uint64, error) { return c.nonce, nil }

func (c *chain) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{BaseFee: big.NewInt(100)}, nil
}

func (c *chain) SuggestGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(1), nil }

func (c *chain) EstimateGas(context.Context, ethereum.CallMsg) (uint64, error) { return 100_000, nil }

func (c *chain) SendTransaction(_ context.Context, tx *types.Transaction) error {
	c.sent = append(c.sent, tx)
	return nil
}

func (c *chain) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	return nil, ethereum.NotFound
}

var tokenABI = func() abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(`[
	 {"type":"function","name":"name","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"string"}]},
	 {"type":"function","name":"DOMAIN_SEPARATOR","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"bytes32"}]},
	 {"type":"function","name":"balanceOf","stateMutability":"view","inputs":[{"name":"account","type":"address"}],"outputs":[{"name":"","type":"uint256"}]},
	 {"type":"function","name":"authorizationState","stateMutability":"view","inputs":[{"name":"authorizer","type":"address"},{"name":"nonce","type":"bytes32"}],"outputs":[{"name":"","type":"bool"}]},
	 {"type":"function","name":"isValidSignature","stateMutability":"view","inputs":[{"name":"hash","type":"bytes32"},{"name":"signature","type":"bytes"}],"outputs":[{"name":"magicValue","type":"bytes4"}]}]`))
	if err != nil {
		panic(err)
	}
	return parsed
}()

func domainSeparator(name, version string) [32]byte {
	domain := apitypes.TypedDataDomain{Name: name, Version: version, ChainId: (*math.HexOrDecimal256)(chainID),
		VerifyingContract: usdg.Hex()}
	typed := apitypes.TypedData{Types: apitypes.Types{"EIP712Domain": {{Name: "name", Type: "string"},
		{Name: "version", Type: "string"}, {Name: "chainId", Type: "uint256"}, {Name: "verifyingContract", Type: "address"}}},
		PrimaryType: "EIP712Domain", Domain: domain}
	hash, err := typed.HashStruct("EIP712Domain", domain.Map())
	if err != nil {
		panic(err)
	}
	return [32]byte(hash)
}

func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return key
}

var authorization = map[string][]evm.TypedDataField{"TransferWithAuthorization": {{Name: "from", Type: "address"},
	{Name: "to", Type: "address"}, {Name: "value", Type: "uint256"}, {Name: "validAfter", Type: "uint256"},
	{Name: "validBefore", Type: "uint256"}, {Name: "nonce", Type: "bytes32"}}}

func message(from common.Address) map[string]any {
	return map[string]any{"from": from.Hex(), "to": payTo.Hex(), "value": "1000", "validAfter": "0",
		"validBefore": "4102444800", "nonce": common.Hash{1}.Hex()}
}

func domain() evm.TypedDataDomain {
	return evm.TypedDataDomain{Name: "Global Dollar", Version: "1", ChainID: chainID, VerifyingContract: usdg.Hex()}
}

func sign(t *testing.T, key *ecdsa.PrivateKey, message map[string]any) []byte {
	t.Helper()
	typed := apitypes.TypedData{Types: apitypes.Types{"EIP712Domain": {{Name: "name", Type: "string"},
		{Name: "version", Type: "string"}, {Name: "chainId", Type: "uint256"}, {Name: "verifyingContract", Type: "address"}},
		"TransferWithAuthorization": {{Name: "from", Type: "address"}, {Name: "to", Type: "address"},
			{Name: "value", Type: "uint256"}, {Name: "validAfter", Type: "uint256"}, {Name: "validBefore", Type: "uint256"},
			{Name: "nonce", Type: "bytes32"}}},
		PrimaryType: "TransferWithAuthorization",
		Domain: apitypes.TypedDataDomain{Name: "Global Dollar", Version: "1", ChainId: (*math.HexOrDecimal256)(chainID),
			VerifyingContract: usdg.Hex()},
		Message: message}
	digest, _, err := apitypes.TypedDataAndHash(typed)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := crypto.Sign(digest, key)
	if err != nil {
		t.Fatal(err)
	}
	signature[64] += 27
	return signature
}

func TestAnAccountsSignatureOfATransferAuthorizationVerifies(t *testing.T) {
	signer, err := payments.NewSigner(t.Context(), &chain{}, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	payer := newKey(t)
	from := crypto.PubkeyToAddress(payer.PublicKey)
	signature := sign(t, payer, message(from))
	if ok, err := signer.VerifyTypedData(t.Context(), from.Hex(), domain(), authorization, "TransferWithAuthorization",
		message(from), signature); err != nil || !ok {
		t.Fatalf("a signature of the payer's: %v %v", ok, err)
	}
	tampered := message(from)
	tampered["value"] = "1001"
	if ok, _ := signer.VerifyTypedData(t.Context(), from.Hex(), domain(), authorization, "TransferWithAuthorization",
		tampered, signature); ok {
		t.Fatal("a signature of another amount verified")
	}
	other := newKey(t)
	if ok, _ := signer.VerifyTypedData(t.Context(), crypto.PubkeyToAddress(other.PublicKey).Hex(), domain(), authorization,
		"TransferWithAuthorization", message(from), signature); ok {
		t.Fatal("a signature verified for another account")
	}
}

func TestAContractWalletSignsByERC1271(t *testing.T) {
	wallet := common.HexToAddress("0x1111111111111111111111111111111111111111")
	node := &chain{code: map[common.Address][]byte{wallet: {0x60}}, answers: map[string]func([]any) ([]any, error){}}
	signer, err := payments.NewSigner(t.Context(), node, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	magic := [4]byte{0x16, 0x26, 0xba, 0x7e}
	node.answers["isValidSignature"] = func([]any) ([]any, error) { return []any{magic}, nil }
	if ok, err := signer.VerifyTypedData(t.Context(), wallet.Hex(), domain(), authorization, "TransferWithAuthorization",
		message(wallet), []byte("a wallet's own signature, longer than sixty-five bytes, as a smart account makes it..........")); err != nil || !ok {
		t.Fatalf("a contract that vouches: %v %v", ok, err)
	}
	node.answers["isValidSignature"] = func([]any) ([]any, error) { return []any{[4]byte{}}, nil }
	if ok, _ := signer.VerifyTypedData(t.Context(), wallet.Hex(), domain(), authorization, "TransferWithAuthorization",
		message(wallet), []byte("another")); ok {
		t.Fatal("a contract that does not vouch was believed")
	}
}

func TestTheTokensDomainIsReadAndCheckedAgainstItsSeparator(t *testing.T) {
	node := &chain{answers: map[string]func([]any) ([]any, error){
		"name":             func([]any) ([]any, error) { return []any{"Global Dollar"}, nil },
		"DOMAIN_SEPARATOR": func([]any) ([]any, error) { return []any{domainSeparator("Global Dollar", "1")}, nil },
	}}
	name, version, err := payments.TokenDomain(t.Context(), node, usdg)
	if err != nil || name != "Global Dollar" || version != "1" {
		t.Fatalf("%q %q %v", name, version, err)
	}
	node.answers["DOMAIN_SEPARATOR"] = func([]any) ([]any, error) { return []any{domainSeparator("Global Dollar", "2")}, nil }
	if _, version, err = payments.TokenDomain(t.Context(), node, usdg); err != nil || version != "2" {
		t.Fatalf("version 2: %q %v", version, err)
	}
	node.answers["DOMAIN_SEPARATOR"] = func([]any) ([]any, error) { return []any{domainSeparator("Another", "1")}, nil }
	if _, _, err = payments.TokenDomain(t.Context(), node, usdg); err == nil {
		t.Fatal("a separator of another name was taken")
	}
}

func TestTheSignerSendsEachTransactionWithTheNextNonceAndAMarginOnItsGas(t *testing.T) {
	node := &chain{nonce: 7}
	signer, err := payments.NewSigner(t.Context(), node, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := signer.SendTransaction(t.Context(), usdg.Hex(), []byte{1, 2, 3}); err != nil {
			t.Fatal(err)
		}
	}
	if len(node.sent) != 2 || node.sent[0].Nonce() != 7 || node.sent[1].Nonce() != 8 {
		t.Fatalf("nonces: %v", node.sent)
	}
	if node.sent[0].Gas() != 120_000 || node.sent[0].Type() != types.DynamicFeeTxType || *node.sent[0].To() != usdg {
		t.Fatalf("a transaction: gas %d type %d", node.sent[0].Gas(), node.sent[0].Type())
	}
}
