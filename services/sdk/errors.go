// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/usdg"
)

// Revert is a call's revert, decoded: a custom error of Tapehouse's contracts, its baskets, the band, the Stock Tokens or USDG, or
// Solidity's Error(string) and Panic(uint256).
type Revert struct {
	Name string
	Args []any
	Data []byte
	err  error
}

func (r *Revert) Error() string {
	args := make([]string, len(r.Args))
	for i, arg := range r.Args {
		if word, ok := arg.([32]byte); ok {
			args[i] = hexutil.Encode(word[:])
		} else {
			args[i] = fmt.Sprint(arg)
		}
	}
	return fmt.Sprintf("%s(%s)", r.Name, strings.Join(args, ", "))
}

func (r *Revert) Unwrap() error {
	return r.err
}

var revertErrors = func() map[[4]byte]abi.Error {
	errs := map[[4]byte]abi.Error{}
	for _, metadata := range []*bind.MetaData{
		&shortpositions.ShortPositionsMetaData,
		&gapcover.GapCoverMetaData,
		&marginaccounts.MarginAccountsMetaData,
		&basket.BasketMetaData,
		&bandfeed.BandFeedMetaData,
		&morphobandoracle.MorphoBandOracleMetaData,
		&band.BandMetaData,
		&stocktoken.StockTokenMetaData,
		&usdg.UsdgMetaData,
	} {
		parsed, err := metadata.ParseABI()
		if err != nil {
			panic(err)
		}
		for _, e := range parsed.Errors {
			if _, ok := errs[[4]byte(e.ID[:4])]; !ok {
				errs[[4]byte(e.ID[:4])] = e
			}
		}
	}
	return errs
}()

var (
	errorSelector = []byte{0x08, 0xc3, 0x79, 0xa0}
	panicSelector = []byte{0x4e, 0x48, 0x7b, 0x71}
)

// DecodeRevert decodes the revert behind err, as returned by a call or a gas estimate over ethclient.
func DecodeRevert(err error) (*Revert, bool) {
	var revert *Revert
	if errors.As(err, &revert) {
		return revert, true
	}
	data, ok := ethclient.RevertErrorData(err)
	if !ok {
		return nil, false
	}
	revert, ok = DecodeRevertData(data)
	if ok {
		revert.err = err
	}
	return revert, ok
}

// DecodeRevertData decodes raw revert data.
func DecodeRevertData(data []byte) (*Revert, bool) {
	if len(data) < 4 {
		return nil, false
	}
	switch {
	case bytes.Equal(data[:4], errorSelector):
		reason, err := abi.UnpackRevert(data)
		return &Revert{Name: "Error", Args: []any{reason}, Data: data}, err == nil
	case bytes.Equal(data[:4], panicSelector):
		code, err := abi.Arguments{{Type: uint256}}.Unpack(data[4:])
		return &Revert{Name: "Panic", Args: code, Data: data}, err == nil
	}
	e, ok := revertErrors[[4]byte(data[:4])]
	if !ok {
		return nil, false
	}
	args, err := e.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, false
	}
	return &Revert{Name: e.Name, Args: args, Data: data}, true
}

var uint256, _ = abi.NewType("uint256", "", nil)

func decoded(err error) error {
	if revert, ok := DecodeRevert(err); ok {
		return revert
	}
	return err
}
