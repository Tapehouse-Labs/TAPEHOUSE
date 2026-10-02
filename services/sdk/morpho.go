// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
)

// NoPrice is why a Morpho oracle has no price.
type NoPrice string

const (
	// NoPriceStale is a band with no live leg, which a fresh writePrices ends.
	NoPriceStale NoPrice = "stale"
	// NoPriceHalted is a band that holds a signed halt, the issuer's pause or an unconfirmed multiplier step for the
	// asset.
	NoPriceHalted NoPrice = "halted"
	// NoPriceSequencerNotSettled is an L2 sequencer that is down or came back an hour ago or less.
	NoPriceSequencerNotSettled NoPrice = "sequencerNotSettled"
)

// OraclePrice is a Morpho oracle's answer: the price of one whole collateral token in the loan token, with Morpho's
// 36 + loan decimals - collateral decimals, or, where Price is nil, why it has none and the revert that said so.
// Never a price of zero.
type OraclePrice struct {
	Price   *big.Int
	NoPrice NoPrice
	Revert  *Revert
}

// MorphoOracle is a Morpho oracle: the band it reads, the symbol it prices, its markets' collateral and loan tokens,
// the factor on the band's 8-decimal price, and its owner, who may re-point it to another band.
type MorphoOracle struct {
	Address         common.Address
	Band            common.Address
	Symbol          [32]byte
	CollateralToken common.Address
	LoanToken       common.Address
	ScaleFactor     *big.Int
	Owner           common.Address
}

// MorphoOracles reads the band's Morpho Blue oracles, one per asset in the registry's .morphoOracles.
type MorphoOracles struct {
	c      *Client
	oracle *morphobandoracle.MorphoBandOracle
	band   *band.Band
}

// MorphoOracles returns the Morpho oracles of the registry's .morphoOracles.
func (c *Client) MorphoOracles() *MorphoOracles {
	return &MorphoOracles{c: c, oracle: morphobandoracle.NewMorphoBandOracle(), band: band.NewBand()}
}

// Price reads the price asset's Morpho oracle answers Morpho Blue. A revert with NoAnswer or SequencerNotSettled is no
// price, never zero; the oracle's halt and the band's corporate action, read at the same block, tell a stale band from
// a halt.
func (m *MorphoOracles) Price(opts *bind.CallOpts, asset string) (OraclePrice, error) {
	oracle := lookup(m.c.deployments.MorphoOracles, asset, ".morphoOracles")
	if oracle.err != nil {
		return OraclePrice{}, oracle.err
	}
	pinned, err := m.c.pin(opts)
	if err != nil {
		return OraclePrice{}, err
	}
	price, err := read(m.c, pinned, oracle, m.oracle.UnpackPrice)(m.oracle.TryPackPrice())
	var revert *Revert
	switch {
	case err == nil:
		return OraclePrice{Price: price}, nil
	case !errors.As(err, &revert):
		return OraclePrice{}, err
	case revert.Name == "SequencerNotSettled":
		return OraclePrice{NoPrice: NoPriceSequencerNotSettled, Revert: revert}, nil
	case revert.Name != "NoAnswer":
		return OraclePrice{}, err
	}
	halt, err := read(m.c, pinned, oracle, m.oracle.UnpackHalt)(m.oracle.TryPackHalt())
	if err != nil {
		return OraclePrice{}, err
	}
	bandAddress, err := read(m.c, pinned, oracle, m.oracle.UnpackBand)(m.oracle.TryPackBand())
	if err != nil {
		return OraclePrice{}, err
	}
	symbol, err := read(m.c, pinned, oracle, m.oracle.UnpackSymbol)(m.oracle.TryPackSymbol())
	if err != nil {
		return OraclePrice{}, err
	}
	step, err := read(m.c, pinned, target{address: bandAddress}, m.band.UnpackCorporateAction)(
		m.band.TryPackCorporateAction(symbol))
	if err != nil {
		return OraclePrice{}, err
	}
	reason := NoPriceStale
	if halt.SignedHalt || halt.OraclePaused || step.Status == 2 {
		reason = NoPriceHalted
	}
	return OraclePrice{NoPrice: reason, Revert: revert}, nil
}

// Halt reads asset's trading halt as its Morpho oracle reports it from the band: whether a halt signed by Tapehouse's
// halt signer holds, until when, when it was issued, and whether the issuer has paused the Stock Token's oracle.
func (m *MorphoOracles) Halt(opts *bind.CallOpts, asset string) (morphobandoracle.HaltOutput, error) {
	oracle := lookup(m.c.deployments.MorphoOracles, asset, ".morphoOracles")
	return read(m.c, opts, oracle, m.oracle.UnpackHalt)(m.oracle.TryPackHalt())
}

// Oracle reads asset's Morpho oracle, every field at opts' block, or at the latest block where opts names none.
func (m *MorphoOracles) Oracle(opts *bind.CallOpts, asset string) (MorphoOracle, error) {
	oracle := lookup(m.c.deployments.MorphoOracles, asset, ".morphoOracles")
	if oracle.err != nil {
		return MorphoOracle{}, oracle.err
	}
	pinned, err := m.c.pin(opts)
	if err != nil {
		return MorphoOracle{}, err
	}
	out := MorphoOracle{Address: oracle.address}
	if out.Band, err = read(m.c, pinned, oracle, m.oracle.UnpackBand)(m.oracle.TryPackBand()); err != nil {
		return MorphoOracle{}, err
	}
	if out.Symbol, err = read(m.c, pinned, oracle, m.oracle.UnpackSymbol)(m.oracle.TryPackSymbol()); err != nil {
		return MorphoOracle{}, err
	}
	if out.CollateralToken, err = read(m.c, pinned, oracle, m.oracle.UnpackCollateralToken)(
		m.oracle.TryPackCollateralToken()); err != nil {
		return MorphoOracle{}, err
	}
	if out.LoanToken, err = read(m.c, pinned, oracle, m.oracle.UnpackLoanToken)(m.oracle.TryPackLoanToken()); err != nil {
		return MorphoOracle{}, err
	}
	if out.ScaleFactor, err = read(m.c, pinned, oracle, m.oracle.UnpackScaleFactor)(m.oracle.TryPackScaleFactor()); err != nil {
		return MorphoOracle{}, err
	}
	if out.Owner, err = read(m.c, pinned, oracle, m.oracle.UnpackOwner)(m.oracle.TryPackOwner()); err != nil {
		return MorphoOracle{}, err
	}
	return out, nil
}
