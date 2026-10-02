// SPDX-License-Identifier: MIT OR Apache-2.0

package codec_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/internal/codec"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphoblue"
)

var (
	alice = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	bob   = common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	nvda  = [32]byte{'N', 'V', 'D', 'A'}
)

func parsed(t *testing.T, metadata *bind.MetaData) *abi.ABI {
	t.Helper()
	out, err := metadata.ParseABI()
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// emit builds the log a contract emits for event with args, in the order of the event's inputs.
func emit(t *testing.T, event abi.Event, args ...any) types.Log {
	t.Helper()
	var indexed [][]any
	var data []any
	for i, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, []any{args[i]})
		} else {
			data = append(data, args[i])
		}
	}
	topics, err := abi.MakeTopics(indexed...)
	if err != nil {
		t.Fatal(err)
	}
	packed, err := event.Inputs.NonIndexed().Pack(data...)
	if err != nil {
		t.Fatal(err)
	}
	log := types.Log{Topics: []common.Hash{event.ID}, Data: packed}
	for _, topic := range topics {
		log.Topics = append(log.Topics, topic[0])
	}
	return log
}

// asJSON is v marshalled and read back, as the API serves it.
func asJSON(t *testing.T, v any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestAnEventDecodesAsItsBindingsUnpackerDoes(t *testing.T) {
	bandABI := parsed(t, &band.BandMetaData)
	log := emit(t, bandABI.Events["HaltWritten"], nvda, true, uint64(1_790_000_000), uint64(1_790_000_900))
	typed, err := band.NewBand().UnpackHaltWrittenEvent(&log)
	if err != nil {
		t.Fatal(err)
	}
	got, err := codec.Event(bandABI.Events["HaltWritten"], log)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"symbol": "0x4e56444100000000000000000000000000000000000000000000000000000000",
		"halted": typed.Halted, "issuedAt": "1790000000", "expiresAt": "1790000900"}
	if !reflect.DeepEqual(asJSON(t, got), want) || typed.Symbol != nvda {
		t.Fatalf("HaltWritten: %v, want %v", got, want)
	}

	oracleABI := parsed(t, &morphobandoracle.MorphoBandOracleMetaData)
	log = emit(t, oracleABI.Events["BandSet"], alice, bob)
	set, err := morphobandoracle.NewMorphoBandOracle().UnpackBandSetEvent(&log)
	if err != nil {
		t.Fatal(err)
	}
	got, _ = codec.Event(oracleABI.Events["BandSet"], log)
	if got["previousBand"] != set.PreviousBand.Hex() || got["newBand"] != set.NewBand.Hex() {
		t.Fatalf("BandSet: %v, want %+v", got, set)
	}
}

func TestArraysAndWideIntegersDecodeAsDecimalStrings(t *testing.T) {
	basketABI := parsed(t, &basket.BasketMetaData)
	huge, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	log := emit(t, basketABI.Events["Rebalanced"], alice, bob, []*big.Int{big.NewInt(1), huge}, []*big.Int{})
	typed, err := basket.NewBasket().UnpackRebalancedEvent(&log)
	if err != nil {
		t.Fatal(err)
	}
	got, err := codec.Event(basketABI.Events["Rebalanced"], log)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"caller": alice.Hex(), "receiver": bob.Hex(), "assetsIn": []any{"1", huge.String()},
		"assetsOut": []any{}}
	if !reflect.DeepEqual(asJSON(t, got), want) || typed.AssetsIn[1].Cmp(huge) != 0 {
		t.Fatalf("Rebalanced: %v, want %v", got, want)
	}
}

func TestASealedBandDecodesEveryField(t *testing.T) {
	feedABI := parsed(t, &bandfeed.BandFeedMetaData)
	log := emit(t, feedABI.Events["Sealed"], uint64(1_790_200_800_000), uint8(2), uint8(1), uint64(22_846_110_842),
		uint64(150), uint64(22_503_419_179), big.NewInt(23_188_802_504))
	typed, err := bandfeed.NewBandFeed().UnpackSealedEvent(&log)
	if err != nil {
		t.Fatal(err)
	}
	got, err := codec.Event(feedABI.Events["Sealed"], log)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"reopenMs": "1790200800000", "state": "2", "live": "1", "mid": "22846110842",
		"halfBps": "150", "low": "22503419179", "high": typed.High.String()}
	if !reflect.DeepEqual(asJSON(t, got), want) {
		t.Fatalf("Sealed: %v, want %v", got, want)
	}
}

func TestATupleOutputDecodesByItsFieldNames(t *testing.T) {
	method := parsed(t, &morphoblue.MorphoBlueMetaData).Methods["idToMarketParams"]
	params := morphoblue.IMorphoMarketParams{LoanToken: alice, CollateralToken: bob, Oracle: alice, Irm: bob,
		Lltv: big.NewInt(625_000_000_000_000_000)}
	packed, err := method.Outputs.Pack(params)
	if err != nil {
		t.Fatal(err)
	}
	values, err := method.Outputs.UnpackValues(packed)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"arg0": map[string]any{"loanToken": alice.Hex(), "collateralToken": bob.Hex(),
		"oracle": alice.Hex(), "irm": bob.Hex(), "lltv": "625000000000000000"}}
	if got := asJSON(t, codec.Values(method.Outputs, values)); !reflect.DeepEqual(got, want) {
		t.Fatalf("idToMarketParams: %v, want %v", got, want)
	}
}

func TestALogOfAnotherEventOrShapeIsRefused(t *testing.T) {
	bandABI := parsed(t, &band.BandMetaData)
	log := emit(t, bandABI.Events["HaltWritten"], nvda, true, uint64(1), uint64(2))
	if _, err := codec.Event(bandABI.Events["PriceWritten"], log); err == nil {
		t.Fatal("a HaltWritten log decoded as PriceWritten")
	}
	log.Topics = log.Topics[:1]
	if _, err := codec.Event(bandABI.Events["HaltWritten"], log); err == nil {
		t.Fatal("a log without its indexed symbol decoded")
	}
}

func TestQueryValuesParseAsTheirABITypes(t *testing.T) {
	types := map[string]abi.Type{}
	for _, name := range []string{"uint8", "uint256", "int256", "bool", "address", "bytes32", "bytes", "uint256[]", "string", "bytes4"} {
		typ, err := abi.NewType(name, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		types[name] = typ
	}
	for _, c := range []struct {
		typ, text string
		want      any
	}{
		{"uint8", "255", uint8(255)},
		{"uint256", "0x10", big.NewInt(16)},
		{"int256", "-5", big.NewInt(-5)},
		{"bool", "true", true},
		{"address", "0x70997970c51812dc3a010c7d01b50e0d17dc79c8", alice},
		{"address", alice.Hex(), alice},
		{"bytes32", "NVDA", nvda},
		{"bytes32", "0x4e56444100000000000000000000000000000000000000000000000000000000", nvda},
		{"bytes", "0x0102", []byte{1, 2}},
		{"bytes4", "0x01020304", [4]byte{1, 2, 3, 4}},
		{"uint256[]", "1,2", []*big.Int{big.NewInt(1), big.NewInt(2)}},
		{"uint256[]", "", []*big.Int{}},
		{"string", "NVDA---24_7", "NVDA---24_7"},
	} {
		got, err := codec.Parse(types[c.typ], c.text)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s %q: %#v, %v; want %#v", c.typ, c.text, got, err, c.want)
		}
	}
	for _, c := range []struct{ typ, text string }{
		{"uint8", "256"},
		{"uint256", "-1"},
		{"uint256", "1.5"},
		{"uint256", "+1"},
		{"bool", "1"},
		{"address", "0x70997970C51812dc3a010c7d01b50e0d17dc79c8"},
		{"bytes32", "a symbol longer than thirty-two bytes"},
		{"bytes4", "0x0102"},
		{"bytes", "0102"},
	} {
		if got, err := codec.Parse(types[c.typ], c.text); err == nil {
			t.Errorf("%s %q parsed as %v", c.typ, c.text, got)
		}
	}
}
