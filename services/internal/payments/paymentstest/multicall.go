// SPDX-License-Identifier: MIT OR Apache-2.0

// Package paymentstest is for tests of payments on a dev node, as net/http/httptest is for tests of servers.
package paymentstest

import (
	"bytes"
	"context"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/x402-foundation/x402/go/v2/mechanisms/evm"
)

// Multicall3 is a dev node's RPC with the Multicall3 that Robinhood Chain has and a dev node lacks, which the x402
// facilitator reads through: it answers a call of tryAggregate by making each of its calls on the node.
type Multicall3 struct{ *ethclient.Client }

// CallContract makes the call on the node, or each call a tryAggregate batches.
func (m Multicall3) CallContract(ctx context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	if call.To == nil || *call.To != common.HexToAddress(evm.MULTICALL3Address) {
		return m.Client.CallContract(ctx, call, block)
	}
	parsed, err := abi.JSON(bytes.NewReader(evm.Multicall3TryAggregateABI))
	if err != nil {
		return nil, err
	}
	method := parsed.Methods["tryAggregate"]
	args, err := method.Inputs.Unpack(call.Data[4:])
	if err != nil {
		return nil, err
	}
	calls := reflect.ValueOf(args[1])
	results := make([]struct {
		Success    bool   `json:"success"`
		ReturnData []byte `json:"returnData"`
	}, calls.Len())
	for i := range results {
		target := calls.Index(i).FieldByName("Target").Interface().(common.Address)
		data := calls.Index(i).FieldByName("CallData").Bytes()
		out, err := m.Client.CallContract(ctx, ethereum.CallMsg{From: call.From, To: &target, Data: data}, block)
		results[i].Success, results[i].ReturnData = err == nil, out
	}
	return method.Outputs.Pack(results)
}
