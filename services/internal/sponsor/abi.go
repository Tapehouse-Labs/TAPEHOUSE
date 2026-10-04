// SPDX-License-Identifier: MIT OR Apache-2.0

package sponsor

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	factoryABI = mustABI(`[
		{"type":"function","name":"createAccount","stateMutability":"nonpayable","inputs":[{"name":"owner","type":"address"},{"name":"salt","type":"uint256"}],"outputs":[{"name":"","type":"address"}]},
		{"type":"function","name":"getAddress","stateMutability":"view","inputs":[{"name":"owner","type":"address"},{"name":"salt","type":"uint256"}],"outputs":[{"name":"","type":"address"}]},
		{"type":"function","name":"accountImplementation","stateMutability":"view","inputs":[],"outputs":[{"name":"","type":"address"}]}]`)
	accountABI = mustABI(`[
		{"type":"function","name":"execute","stateMutability":"nonpayable","inputs":[{"name":"dest","type":"address"},{"name":"value","type":"uint256"},{"name":"func","type":"bytes"}],"outputs":[]},
		{"type":"function","name":"executeBatch","stateMutability":"nonpayable","inputs":[{"name":"dest","type":"address[]"},{"name":"value","type":"uint256[]"},{"name":"func","type":"bytes[]"}],"outputs":[]}]`)
	erc20ABI = mustABI(`[
		{"type":"function","name":"approve","stateMutability":"nonpayable","inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
	entryPointABI = mustABI(`[
		{"type":"function","name":"handleOps","stateMutability":"nonpayable","inputs":[{"name":"ops","type":"tuple[]","components":[
			{"name":"sender","type":"address"},{"name":"nonce","type":"uint256"},{"name":"initCode","type":"bytes"},
			{"name":"callData","type":"bytes"},{"name":"accountGasLimits","type":"bytes32"},
			{"name":"preVerificationGas","type":"uint256"},{"name":"gasFees","type":"bytes32"},
			{"name":"paymasterAndData","type":"bytes"},{"name":"signature","type":"bytes"}]},
			{"name":"beneficiary","type":"address"}],"outputs":[]}]`)
	nodeInterfaceABI = mustABI(`[
		{"type":"function","name":"gasEstimateL1Component","stateMutability":"payable","inputs":[{"name":"to","type":"address"},{"name":"contractCreation","type":"bool"},{"name":"data","type":"bytes"}],"outputs":[{"name":"gasEstimateForL1","type":"uint64"},{"name":"baseFee","type":"uint256"},{"name":"l1BaseFeeEstimate","type":"uint256"}]}]`)

	approveSelector = [4]byte(erc20ABI.Methods["approve"].ID)

	// nodeInterface is Arbitrum's NodeInterface, which answers only eth_call.
	nodeInterface = common.HexToAddress("0x00000000000000000000000000000000000000C8")

	// implementationSlot is ERC-1967's implementation slot, where a SimpleAccount's proxy keeps its implementation.
	implementationSlot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

	requestTypeHash = crypto.Keccak256Hash([]byte("UserOperationRequest(address sender,uint256 nonce,bytes initCode," +
		"bytes callData,bytes32 accountGasLimits,uint256 preVerificationGas,bytes32 gasFees," +
		"uint256 paymasterVerificationGasLimit,uint256 paymasterPostOpGasLimit,uint48 validAfter,uint48 validUntil)"))
	domainTypeHash = crypto.Keccak256Hash([]byte(
		"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
)

func mustABI(definition string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		panic(err)
	}
	return parsed
}

type call struct {
	to    common.Address
	value *big.Int
	data  []byte
}

type packedUserOperation struct {
	Sender             common.Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

// decodeCalls reads the calls of a SimpleAccount's execute or executeBatch, and refuses any other function.
func decodeCalls(callData []byte) ([]call, error) {
	if len(callData) < 4 {
		return nil, refuse("the call data is not execute or executeBatch")
	}
	method, err := accountABI.MethodById(callData[:4])
	if err != nil {
		return nil, refuse("the account is called with %#x, not execute or executeBatch", callData[:4])
	}
	args, err := method.Inputs.Unpack(callData[4:])
	if err != nil {
		return nil, refuse("the call data of %s does not decode", method.Name)
	}
	if method.Name == "execute" {
		return []call{{args[0].(common.Address), args[1].(*big.Int), args[2].([]byte)}}, nil
	}
	targets, values, data := args[0].([]common.Address), args[1].([]*big.Int), args[2].([][]byte)
	if len(targets) == 0 || len(data) != len(targets) || (len(values) != 0 && len(values) != len(targets)) {
		return nil, refuse("the batch's arrays do not match")
	}
	calls := make([]call, len(targets))
	for i, target := range targets {
		value := new(big.Int)
		if len(values) != 0 {
			value = values[i]
		}
		calls[i] = call{target, value, data[i]}
	}
	return calls, nil
}

// unpack decodes the arguments of a call to the method name of contract.
func unpack(contract abi.ABI, name string, data []byte) ([]any, error) {
	method := contract.Methods[name]
	if len(data) < 4 || [4]byte(data[:4]) != [4]byte(method.ID) {
		return nil, refuse("the call is not %s", name)
	}
	return method.Inputs.Unpack(data[4:])
}

// address calls a view of contract that returns an address.
func (s *Service) address(ctx context.Context, to common.Address, contract abi.ABI, name string, args ...any) (
	common.Address, error) {
	data, err := contract.Pack(name, args...)
	if err != nil {
		return common.Address{}, err
	}
	out, err := s.call(ctx, to, data)
	if err != nil {
		return common.Address{}, err
	}
	values, err := contract.Unpack(name, out)
	if err != nil {
		return common.Address{}, err
	}
	return values[0].(common.Address), nil
}

func initCode(op *UserOperation) []byte {
	if op.Factory == nil {
		return []byte{}
	}
	return append(op.Factory.Bytes(), op.FactoryData...)
}

// halves packs two 128-bit values into a word, high then low, as EntryPoint v0.7 packs gas limits and fees.
func halves(high, low *big.Int) [32]byte {
	var word [32]byte
	copy(word[:16], math.U256Bytes(new(big.Int).Set(high))[16:])
	copy(word[16:], math.U256Bytes(new(big.Int).Set(low))[16:])
	return word
}

func pack(op *UserOperation, paymasterAndData []byte) packedUserOperation {
	return packedUserOperation{Sender: op.Sender, Nonce: op.Nonce.ToInt(), InitCode: initCode(op), CallData: op.CallData,
		AccountGasLimits:   halves(op.VerificationGasLimit.ToInt(), op.CallGasLimit.ToInt()),
		PreVerificationGas: op.PreVerificationGas.ToInt(),
		GasFees:            halves(op.MaxPriorityFeePerGas.ToInt(), op.MaxFeePerGas.ToInt()),
		PaymasterAndData:   paymasterAndData, Signature: op.Signature}
}

// paymasterData is the paymaster's data: its validity window, 6 bytes each, then the signature.
func paymasterData(validAfter, validUntil uint64, signature []byte) []byte {
	data := make([]byte, 12, 12+len(signature))
	copy(data[:6], new(big.Int).SetUint64(validAfter).FillBytes(make([]byte, 6)))
	copy(data[6:12], new(big.Int).SetUint64(validUntil).FillBytes(make([]byte, 6)))
	return append(data, signature...)
}

// paymasterAndData is the paymaster, its gas limits and its data, as an operation carries them.
func paymasterAndData(paymaster common.Address, data []byte) []byte {
	out := append(paymaster.Bytes(), common.LeftPadBytes(big.NewInt(PaymasterVerificationGas).Bytes(), 16)...)
	out = append(out, make([]byte, 16)...)
	return append(out, data...)
}

func word(value *big.Int) []byte {
	return math.U256Bytes(new(big.Int).Set(value))
}

func domainSeparator(chainID uint64, paymaster common.Address) common.Hash {
	return crypto.Keccak256Hash(domainTypeHash.Bytes(), crypto.Keccak256([]byte("Tapehouse Sponsor Paymaster")),
		crypto.Keccak256([]byte("1")), word(new(big.Int).SetUint64(chainID)), common.LeftPadBytes(paymaster.Bytes(), 32))
}

// typedDataHash is the EIP-712 hash the paymaster checks the signature against: OpenZeppelin's UserOperationRequest.
func typedDataHash(domain common.Hash, op *UserOperation, validAfter, validUntil uint64) common.Hash {
	limits, fees := halves(op.VerificationGasLimit.ToInt(), op.CallGasLimit.ToInt()),
		halves(op.MaxPriorityFeePerGas.ToInt(), op.MaxFeePerGas.ToInt())
	request := crypto.Keccak256Hash(requestTypeHash.Bytes(), common.LeftPadBytes(op.Sender.Bytes(), 32),
		word(op.Nonce.ToInt()), crypto.Keccak256(initCode(op)), crypto.Keccak256(op.CallData), limits[:],
		word(op.PreVerificationGas.ToInt()), fees[:], word(op.PaymasterVerificationGasLimit.ToInt()),
		word(op.PaymasterPostOpGasLimit.ToInt()), word(new(big.Int).SetUint64(validAfter)),
		word(new(big.Int).SetUint64(validUntil)))
	return crypto.Keccak256Hash([]byte{0x19, 0x01}, domain.Bytes(), request.Bytes())
}
