// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// Address returns the auction's address.
func (r *ReopeningAuction) Address() (common.Address, error) {
	return r.address, r.err
}

// AuctionConstants are the auction's public constants: the least deposit and the least bid worth, in USDG with 6
// decimals, and the lengths of the reveal and clearing windows in milliseconds.
type AuctionConstants struct {
	Bond, MinBid      *big.Int
	RevealMs, ClearMs uint64
}

// Constants reads the auction's BOND, MIN_BID, REVEAL_MS and CLEAR_MS at one block, opts' or the latest.
func (r *ReopeningAuction) Constants(opts *bind.CallOpts) (AuctionConstants, error) {
	if r.err != nil {
		return AuctionConstants{}, r.err
	}
	pinned, err := r.c.pin(opts)
	if err != nil {
		return AuctionConstants{}, err
	}
	var out AuctionConstants
	if out.Bond, err = read(r.c, pinned, r.target, r.auction.UnpackBOND)(r.auction.TryPackBOND()); err != nil {
		return AuctionConstants{}, err
	}
	if out.MinBid, err = read(r.c, pinned, r.target, r.auction.UnpackMINBID)(r.auction.TryPackMINBID()); err != nil {
		return AuctionConstants{}, err
	}
	if out.RevealMs, err = read(r.c, pinned, r.target, r.auction.UnpackREVEALMS)(r.auction.TryPackREVEALMS()); err != nil {
		return AuctionConstants{}, err
	}
	if out.ClearMs, err = read(r.c, pinned, r.target, r.auction.UnpackCLEARMS)(r.auction.TryPackCLEARMS()); err != nil {
		return AuctionConstants{}, err
	}
	return out, nil
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

// Commit commits the sender to a hidden bid in asset's round with deposit of its USDG, at least BOND and at least
// the bid's escrow. Commitment computes the commitment.
func (r *ReopeningAuction) Commit(asset string, commitment [32]byte, deposit *big.Int) (Tx, error) {
	to, symbol := r.of(asset)
	if err := present(deposit); err != nil {
		return Tx{}, err
	}
	return to.tx(r.auction.TryPackCommit(symbol, commitment, deposit))
}

// Reveal reveals the sender's bid for quantity of asset's Stock Token, with 18 decimals, at price in USD with 8
// decimals per token, in the round for the regular open openMs, in the half hour before it.
func (r *ReopeningAuction) Reveal(asset string, openMs uint64, quantity, price *big.Int, salt [32]byte) (Tx, error) {
	to, symbol := r.of(asset)
	if err := present(quantity, price); err != nil {
		return Tx{}, err
	}
	return to.tx(r.auction.TryPackReveal(symbol, openMs, quantity, price, salt))
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

var commitmentArguments = func() abi.Arguments {
	var arguments abi.Arguments
	for _, name := range []string{"address", "bytes32", "uint64", "uint256", "uint256", "bytes32"} {
		kind, err := abi.NewType(name, "", nil)
		if err != nil {
			panic(err)
		}
		arguments = append(arguments, abi.Argument{Type: kind})
	}
	return arguments
}()

// Commitment is the hash a bid is committed under: keccak256(abi.encode(bidder, symbol, openMs, quantity, price,
// salt)), which binds the bid to its bidder.
func Commitment(bidder common.Address, asset string, openMs uint64, quantity, price *big.Int, salt [32]byte) ([32]byte, error) {
	symbol, err := ToBytes32(asset)
	if err != nil {
		return [32]byte{}, err
	}
	if err := present(quantity, price); err != nil {
		return [32]byte{}, err
	}
	packed, err := commitmentArguments.Pack(bidder, symbol, openMs, quantity, price, salt)
	if err != nil {
		return [32]byte{}, err
	}
	return crypto.Keccak256Hash(packed), nil
}
