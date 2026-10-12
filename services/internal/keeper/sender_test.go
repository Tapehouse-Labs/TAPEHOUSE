// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTheSenderKeepsItsOwnNonceWhereANodeAnswersWithAnUsedOne(t *testing.T) {
	w := newWorld(t)
	w.staleNonce = true
	w.write("Token:USDG", "approve", func([]any, bool) error { return nil })
	k := w.keeper(Config{})
	for i := range 3 {
		tx, err := k.client.Approve(w.deployments.Tokens["USDG"], common.Address{9}, big.NewInt(int64(i+1)))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := k.sender.Act(context.Background(), "approve", tx); err != nil {
			t.Fatalf("the transaction %d, with a node a block behind: %v", i, err)
		}
	}
	if got := w.mined(); len(got) != 3 {
		t.Fatalf("sent %v", got)
	}
}
