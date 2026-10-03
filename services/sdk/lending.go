// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
)

// LendingVault reads one of the registry's stock lending vaults and packs its transactions.
type LendingVault struct {
	target
	c     *Client
	vault *stocklendingvault.StockLendingVault
}

// Queue is a lending vault's recall queue: the next ticket in turn, how many were issued, and how much of the queue
// the borrower has returned.
type Queue struct {
	Head     uint64
	Tickets  uint64
	Assigned *big.Int
}

// LendingVault returns the stock lending vault of asset, the registry's .tapehouse.StockLending.<asset>.
func (c *Client) LendingVault(asset string) *LendingVault {
	return &LendingVault{
		target: lookup(c.deployments.StockLending, asset, ".tapehouse.StockLending"),
		c:      c,
		vault:  stocklendingvault.NewStockLendingVault(),
	}
}

// Queue reads the vault's recall queue at opts' block, or at the latest block where opts names none.
func (v *LendingVault) Queue(opts *bind.CallOpts) (Queue, error) {
	if v.err != nil {
		return Queue{}, v.err
	}
	pinned, err := v.c.pin(opts)
	if err != nil {
		return Queue{}, err
	}
	head, err := read(v.c, pinned, v.target, v.vault.UnpackHead)(v.vault.TryPackHead())
	if err != nil {
		return Queue{}, err
	}
	tickets, err := read(v.c, pinned, v.target, v.vault.UnpackTickets)(v.vault.TryPackTickets())
	if err != nil {
		return Queue{}, err
	}
	assigned, err := read(v.c, pinned, v.target, v.vault.UnpackAssigned)(v.vault.TryPackAssigned())
	return Queue{head, tickets, assigned}, err
}

// Ticket reads ticket id: its place in the queue, what of it was taken or given up, and when its notice runs out.
func (v *LendingVault) Ticket(opts *bind.CallOpts, id uint64) (stocklendingvault.TicketOutput, error) {
	return read(v.c, opts, v.target, v.vault.UnpackTicket)(v.vault.TryPackTicket(new(big.Int).SetUint64(id)))
}

// Claimable reads what may be taken now for ticket id: what the vault holds for it if it is next in turn.
func (v *LendingVault) Claimable(opts *bind.CallOpts, id uint64) (*big.Int, error) {
	return read(v.c, opts, v.target, v.vault.UnpackClaimable)(v.vault.TryPackClaimable(new(big.Int).SetUint64(id)))
}

// Counted reads the tokens the vault counts as held: those no ticket waits for and those held for tickets. More than
// its balance after a burn, until the accounts' Sync of its asset.
func (v *LendingVault) Counted(opts *bind.CallOpts) (*big.Int, error) {
	if v.err != nil {
		return nil, v.err
	}
	pinned, err := v.c.pin(opts)
	if err != nil {
		return nil, err
	}
	idle, err := read(v.c, pinned, v.target, v.vault.UnpackIdle)(v.vault.TryPackIdle())
	if err != nil {
		return nil, err
	}
	locked, err := read(v.c, pinned, v.target, v.vault.UnpackLocked)(v.vault.TryPackLocked())
	if err != nil {
		return nil, err
	}
	return idle.Add(idle, locked), nil
}

// BuyIn has the borrower buy in and return the queue up to the end of ticket id, whose notice has run out. Anyone may.
func (v *LendingVault) BuyIn(id uint64) (Tx, error) {
	return v.tx(v.vault.TryPackBuyIn(new(big.Int).SetUint64(id)))
}
