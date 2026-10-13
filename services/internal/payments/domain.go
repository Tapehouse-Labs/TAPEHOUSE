// SPDX-License-Identifier: MIT OR Apache-2.0

package payments

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const tokenDomainABI = `[
 {"type":"function","name":"name","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"string"}]},
 {"type":"function","name":"DOMAIN_SEPARATOR","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"bytes32"}]}]`

// TokenDomain reads the EIP-712 domain a payer signs a transfer authorization under: the token's name, and the
// version under which the domain separator is the one the token reports.
func TokenDomain(ctx context.Context, chain interface {
	ChainID(context.Context) (*big.Int, error)
	CallContract(context.Context, ethereum.CallMsg, *big.Int) ([]byte, error)
}, token common.Address) (name, version string, err error) {
	parsed, err := abi.JSON(strings.NewReader(tokenDomainABI))
	if err != nil {
		return "", "", err
	}
	read := func(function string) (any, error) {
		data, err := parsed.Pack(function)
		if err != nil {
			return nil, err
		}
		raw, err := chain.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
		if err != nil {
			return nil, err
		}
		out, err := parsed.Unpack(function, raw)
		if err != nil {
			return nil, err
		}
		return out[0], nil
	}
	answer, err := read("name")
	if err != nil {
		return "", "", fmt.Errorf("the token's name: %w", err)
	}
	name = answer.(string)
	answer, err = read("DOMAIN_SEPARATOR")
	if err != nil {
		return "", "", fmt.Errorf("the token's domain separator: %w", err)
	}
	reported := answer.([32]byte)
	chainID, err := chain.ChainID(ctx)
	if err != nil {
		return "", "", err
	}
	for _, version := range []string{"1", "2"} {
		domain := apitypes.TypedDataDomain{Name: name, Version: version, ChainId: (*math.HexOrDecimal256)(chainID),
			VerifyingContract: token.Hex()}
		typed := apitypes.TypedData{Types: apitypes.Types{"EIP712Domain": domainFields}, PrimaryType: "EIP712Domain",
			Domain: domain}
		hash, err := typed.HashStruct("EIP712Domain", domain.Map())
		if err != nil {
			return "", "", err
		}
		if [32]byte(hash) == reported {
			return name, version, nil
		}
	}
	return "", "", errors.New("the token's domain separator is not that of its name under version 1 or 2")
}
