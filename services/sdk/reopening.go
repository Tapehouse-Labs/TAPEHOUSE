// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
)

// ReopeningAuction reads the reopening auction and packs its transactions.
type ReopeningAuction struct {
	target
	c       *Client
	auction *reopeningauction.ReopeningAuction
}

// ReopeningAuction returns the reopening auction of the registry's .tapehouse.ReopeningAuction.
func (c *Client) ReopeningAuction() *ReopeningAuction {
	return &ReopeningAuction{
		target:  lookup(c.deployments.Tapehouse, "ReopeningAuction", ".tapehouse"),
		c:       c,
		auction: reopeningauction.NewReopeningAuction(),
	}
}

// Phase reads the regular open, in milliseconds, that bids are committed or revealed for, and whether they are being
// revealed; zero outside those phases.
func (r *ReopeningAuction) Phase(opts *bind.CallOpts) (reopeningauction.PhaseOutput, error) {
	return read(r.c, opts, r.target, r.auction.UnpackPhase)(r.auction.TryPackPhase())
}

// LastOpenMs reads the regular open of the last round opened, in milliseconds.
func (r *ReopeningAuction) LastOpenMs(opts *bind.CallOpts) (uint64, error) {
	return read(r.c, opts, r.target, r.auction.UnpackLastOpenMs)(r.auction.TryPackLastOpenMs())
}

// Round reads asset's round for the regular open openMs.
func (r *ReopeningAuction) Round(opts *bind.CallOpts, asset string, openMs uint64) (reopeningauction.ReopeningAuctionRound, error) {
	to, symbol := r.of(asset)
	return read(r.c, opts, to, r.auction.UnpackRound)(r.auction.TryPackRound(symbol, openMs))
}

// Lots reads the lots of asset's round for the regular open openMs.
func (r *ReopeningAuction) Lots(opts *bind.CallOpts, asset string, openMs uint64) ([]reopeningauction.ReopeningAuctionLot, error) {
	to, symbol := r.of(asset)
	return read(r.c, opts, to, r.auction.UnpackLots)(r.auction.TryPackLots(symbol, openMs))
}

// Bids reads the revealed bids of asset's round for the regular open openMs.
func (r *ReopeningAuction) Bids(opts *bind.CallOpts, asset string, openMs uint64) ([]reopeningauction.ReopeningAuctionBid, error) {
	to, symbol := r.of(asset)
	return read(r.c, opts, to, r.auction.UnpackBids)(r.auction.TryPackBids(symbol, openMs))
}

// Commitment reads a commitment of asset's round: its bidder and its deposit, zero once revealed or forfeited.
func (r *ReopeningAuction) Commitment(opts *bind.CallOpts, asset string, openMs uint64, commitment [32]byte) (reopeningauction.CommitmentsOutput, error) {
	to, symbol := r.of(asset)
	return read(r.c, opts, to, r.auction.UnpackCommitments)(r.auction.TryPackCommitments(symbol, openMs, commitment))
}

// Enrolled reports whether account's position has a lot in asset's round for the regular open openMs.
func (r *ReopeningAuction) Enrolled(opts *bind.CallOpts, openMs uint64, account common.Address, position [32]byte, asset string) (bool, error) {
	to, symbol := r.of(asset)
	return read(r.c, opts, to, r.auction.UnpackEnrolled)(r.auction.TryPackEnrolled(openMs, account, position, symbol))
}

// Enroll enrolls account's position, short at its bands' low edges, in asset's round. Anyone may while bids are
// committed.
func (r *ReopeningAuction) Enroll(account common.Address, position [32]byte, asset string) (Tx, error) {
	to, symbol := r.of(asset)
	return to.tx(r.auction.TryPackEnroll(account, position, symbol))
}

// Clear clears asset's round for the regular open openMs at its clearing price, which ClearingPrice computes. Anyone
// may in the first hour of regular trading.
func (r *ReopeningAuction) Clear(asset string, openMs uint64, price *big.Int) (Tx, error) {
	to, symbol := r.of(asset)
	if err := present(price); err != nil {
		return Tx{}, err
	}
	return to.tx(r.auction.TryPackClear(symbol, openMs, price))
}

// Claim sends the bidder of the bid at index its tokens and the rest of its escrow once the round is cleared, or all
// its escrow once it lapsed. Anyone may.
func (r *ReopeningAuction) Claim(asset string, openMs uint64, index *big.Int) (Tx, error) {
	to, symbol := r.of(asset)
	if err := present(index); err != nil {
		return Tx{}, err
	}
	return to.tx(r.auction.TryPackClaim(symbol, openMs, index))
}

// Forfeit sends an unrevealed commitment's forfeit to the accounts' reserve and the rest of its deposit back, once the
// round's bids are revealed. Anyone may.
func (r *ReopeningAuction) Forfeit(asset string, openMs uint64, commitment [32]byte) (Tx, error) {
	to, symbol := r.of(asset)
	return to.tx(r.auction.TryPackForfeit(symbol, openMs, commitment))
}

// ClearingPrice is a round's only clearing price: the highest bid price at which the bids at or above it take the
// whole supply; when all of them together take less, the lowest bid price; with no bid, the floor. False when no price
// clears it, as for a round with bids and no supply.
func ClearingPrice(bids []reopeningauction.ReopeningAuctionBid, supply, floor *big.Int) (*big.Int, bool) {
	if len(bids) == 0 {
		return new(big.Int).Set(floor), true
	}
	if supply.Sign() == 0 {
		return nil, false
	}
	var lowest *big.Int
	var best *big.Int
	for _, candidate := range bids {
		if lowest == nil || candidate.Price.Cmp(lowest) < 0 {
			lowest = candidate.Price
		}
		above, at := new(big.Int), new(big.Int)
		for _, bid := range bids {
			switch bid.Price.Cmp(candidate.Price) {
			case 1:
				above.Add(above, bid.Quantity)
			case 0:
				at.Add(at, bid.Quantity)
			}
		}
		if above.Cmp(supply) < 0 && new(big.Int).Add(above, at).Cmp(supply) >= 0 &&
			(best == nil || candidate.Price.Cmp(best) > 0) {
			best = candidate.Price
		}
	}
	if best == nil {
		best = lowest
	}
	return new(big.Int).Set(best), true
}

func (r *ReopeningAuction) of(asset string) (target, [32]byte) {
	symbol, err := ToBytes32(asset)
	return r.with(err), symbol
}
