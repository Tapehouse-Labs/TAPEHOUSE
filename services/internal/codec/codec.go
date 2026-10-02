// SPDX-License-Identifier: MIT OR Apache-2.0

// Package codec turns ABI values into JSON and query strings into ABI values. Every integer is a decimal string, as
// JSON numbers cannot carry 256 bits; an address is checksummed and bytes are 0x-prefixed lowercase hex.
package codec

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// Address is an address that marshals checksummed, as the API serves every address.
type Address common.Address

// MarshalText returns the address checksummed.
func (a Address) MarshalText() ([]byte, error) {
	return []byte(common.Address(a).Hex()), nil
}

// UnmarshalText reads a hex address.
func (a *Address) UnmarshalText(text []byte) error {
	return (*common.Address)(a).UnmarshalText(text)
}

// Value returns the JSON form of v, a value of type t as go-ethereum unpacks it.
func Value(t abi.Type, v any) any {
	value := reflect.ValueOf(v)
	switch t.T {
	case abi.IntTy, abi.UintTy:
		switch n := v.(type) {
		case *big.Int:
			return n.String()
		default:
			if value.CanInt() {
				return strconv.FormatInt(value.Int(), 10)
			}
			return strconv.FormatUint(value.Uint(), 10)
		}
	case abi.BoolTy:
		return v.(bool)
	case abi.StringTy:
		return v.(string)
	case abi.AddressTy:
		return v.(common.Address).Hex()
	case abi.BytesTy:
		return hexutil.Encode(v.([]byte))
	case abi.FixedBytesTy, abi.HashTy:
		raw := make([]byte, value.Len())
		reflect.Copy(reflect.ValueOf(raw), value)
		return hexutil.Encode(raw)
	case abi.SliceTy, abi.ArrayTy:
		out := make([]any, value.Len())
		for i := range out {
			out[i] = Value(*t.Elem, value.Index(i).Interface())
		}
		return out
	case abi.TupleTy:
		out := make(map[string]any, len(t.TupleElems))
		for i, elem := range t.TupleElems {
			out[name(t.TupleRawNames[i], i)] = Value(*elem, value.Field(i).Interface())
		}
		return out
	}
	return fmt.Sprint(v)
}

// Values returns the JSON object of values unpacked for arguments, keyed by each argument's name, or arg<i> for an
// unnamed one.
func Values(arguments abi.Arguments, values []any) map[string]any {
	out := make(map[string]any, len(values))
	for i, argument := range arguments {
		out[name(argument.Name, i)] = Value(argument.Type, values[i])
	}
	return out
}

// Event decodes log as event: its indexed arguments from the topics and the rest from the data, as a JSON object.
// An indexed argument of a dynamic type is its topic, the hash of its value.
func Event(event abi.Event, log types.Log) (map[string]any, error) {
	if len(log.Topics) == 0 || log.Topics[0] != event.ID {
		return nil, errors.New("the log is not this event")
	}
	var indexed abi.Arguments
	for _, input := range event.Inputs {
		if input.Indexed {
			indexed = append(indexed, input)
		}
	}
	if len(log.Topics) != len(indexed)+1 {
		return nil, fmt.Errorf("%s has %d indexed arguments, the log %d topics", event.Sig, len(indexed), len(log.Topics)-1)
	}
	data, err := event.Inputs.NonIndexed().UnpackValues(log.Data)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(event.Inputs))
	topic, field := 1, 0
	for i, input := range event.Inputs {
		if !input.Indexed {
			out[name(input.Name, i)] = Value(input.Type, data[field])
			field++
			continue
		}
		out[name(input.Name, i)], err = topicValue(input.Type, log.Topics[topic])
		if err != nil {
			return nil, err
		}
		topic++
	}
	return out, nil
}

func topicValue(t abi.Type, topic common.Hash) (any, error) {
	switch t.T {
	case abi.StringTy, abi.BytesTy, abi.SliceTy, abi.ArrayTy, abi.TupleTy:
		return topic.Hex(), nil
	}
	values, err := abi.Arguments{{Type: t}}.UnpackValues(topic.Bytes())
	if err != nil {
		return nil, err
	}
	return Value(t, values[0]), nil
}

// Parse returns the ABI value of type t that text names: a decimal or 0x-prefixed integer, true or false, an address
// (checksummed where it is in mixed case), 0x-prefixed hex bytes, or, for bytes32, also an asset's symbol or a feed ID
// as the contracts take it, right-padded. An array is its elements separated by commas.
func Parse(t abi.Type, text string) (any, error) {
	switch t.T {
	case abi.IntTy, abi.UintTy:
		n, ok := new(big.Int).SetString(text, 0)
		if !ok || strings.HasPrefix(text, "+") {
			return nil, fmt.Errorf("%q is not an integer", text)
		}
		if t.T == abi.UintTy && n.Sign() < 0 || n.BitLen() > t.Size {
			return nil, fmt.Errorf("%q does not fit %s", text, t)
		}
		return integer(t, n), nil
	case abi.BoolTy:
		b, err := strconv.ParseBool(text)
		if err != nil || (text != "true" && text != "false") {
			return nil, fmt.Errorf("%q is not true or false", text)
		}
		return b, nil
	case abi.StringTy:
		return text, nil
	case abi.AddressTy:
		address := common.HexToAddress(text)
		if !common.IsHexAddress(text) || (strings.ToLower(text) != text && address.Hex() != text) {
			return nil, fmt.Errorf("%q is not an address", text)
		}
		return address, nil
	case abi.BytesTy:
		return hexutil.Decode(text)
	case abi.FixedBytesTy:
		raw, err := hexutil.Decode(text)
		if err != nil && t.Size == 32 {
			word, err := sdk.ToBytes32(text)
			return word, err
		}
		if err != nil || len(raw) != t.Size {
			return nil, fmt.Errorf("%q is not %d bytes of hex", text, t.Size)
		}
		out := reflect.New(t.GetType()).Elem()
		reflect.Copy(out, reflect.ValueOf(raw))
		return out.Interface(), nil
	case abi.SliceTy:
		var parts []string
		if text != "" {
			parts = strings.Split(text, ",")
		}
		out := reflect.MakeSlice(t.GetType(), len(parts), len(parts))
		for i, part := range parts {
			elem, err := Parse(*t.Elem, part)
			if err != nil {
				return nil, err
			}
			out.Index(i).Set(reflect.ValueOf(elem))
		}
		return out.Interface(), nil
	}
	return nil, fmt.Errorf("arguments of type %s are not supported", t)
}

// integer converts n to the Go type go-ethereum packs t from: a sized integer up to 64 bits, *big.Int above.
func integer(t abi.Type, n *big.Int) any {
	out := reflect.New(t.GetType()).Elem()
	switch out.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.SetInt(n.Int64())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.SetUint(n.Uint64())
	default:
		return n
	}
	return out.Interface()
}

func name(raw string, i int) string {
	if raw == "" {
		return "arg" + strconv.Itoa(i)
	}
	return raw
}
