// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/bandfeed"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/quoterv2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
)

var (
	alice = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	bob   = common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
)

func registry(t *testing.T, chainID string) *sdk.Deployments {
	t.Helper()
	d, err := sdk.LoadDeployments("../../deployments/" + chainID + ".json")
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestChainlinkFeedsPriceTheTokenOnRobinhoodChainAndTheShareOnArbitrumOne(t *testing.T) {
	robinhood, arbitrum := registry(t, "4663"), registry(t, "42161")
	token, err := robinhood.TokenPriceFeed("NVDA_USD")
	if err != nil || token.Address != common.HexToAddress("0x379EC4f7C378F34a1B47E4F3cbeBCbAC3E8E9F15") {
		t.Fatalf("token feed %v, %v", token, err)
	}
	share, err := arbitrum.SharePriceFeed("NVDA_USD")
	if err != nil || share.Address != common.HexToAddress("0x4881A4418b5F2460B21d6F08CD5aA0678a7f262F") {
		t.Fatalf("share feed %v, %v", share, err)
	}
	if _, err := robinhood.SharePriceFeed("NVDA_USD"); err == nil || err.Error() != ".chainlink.NVDA_USD prices the Stock Token on chain 4663" {
		t.Fatalf("a token feed read as a share feed: %v", err)
	}
	if _, err := arbitrum.TokenPriceFeed("NVDA_USD"); err == nil || err.Error() != ".chainlink.NVDA_USD prices the share on chain 42161" {
		t.Fatalf("a share feed read as a token feed: %v", err)
	}
}

func TestEveryGroupIsReadAndAMissingOneIsEmpty(t *testing.T) {
	testnet, robinhood := registry(t, "46630"), registry(t, "4663")
	if testnet.ChainID != 46630 {
		t.Fatalf("46630: %+v", testnet)
	}
	data, err := os.ReadFile("../../deployments/46630.json")
	if err != nil {
		t.Fatal(err)
	}
	var raw struct {
		Tapehouse map[string]any    `json:"tapehouse"`
		BandFeeds map[string]string `json:"bandFeeds"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	for name, value := range raw.Tapehouse {
		if address, ok := value.(string); ok && common.HexToAddress(address).Hex() != address {
			t.Errorf(".tapehouse.%s is not checksummed: %s", name, address)
		}
	}
	for name, address := range raw.BandFeeds {
		if common.HexToAddress(address).Hex() != address {
			t.Errorf(".bandFeeds.%s is not checksummed: %s", name, address)
		}
	}
	if robinhood.UniswapV3["QuoterV2"] != common.HexToAddress("0x33e885eD0Ec9bF04EcfB19341582aADCb4c8A9E7") {
		t.Fatalf("4663 quoter: %v", robinhood.UniswapV3["QuoterV2"])
	}
}

func TestAMixedCaseAddressMustCarryItsChecksum(t *testing.T) {
	lower, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"Band":"0xa70118d3324d90532e7d2854627b13cace305641","StockLending":{"SPY":"` + bob.Hex() + `"}}}`))
	if err != nil || lower.Tapehouse["Band"] != common.HexToAddress("0xa70118d3324D90532E7D2854627b13CacE305641") {
		t.Fatalf("lowercase: %v, %v", lower, err)
	}
	if lower.StockLending["SPY"] != bob || len(lower.Tapehouse) != 1 {
		t.Fatalf("stock lending: %+v", lower)
	}
	for registry, message := range map[string]string{
		`{"chainId":1,"tokens":{"USDG":"0xa70118D3324D90532E7D2854627b13CacE305641"}}`: ".tokens.USDG is not an address",
		`{"tokens":{}}`:                        "the registry has no chainId",
		`{"chainId":1,"tapehouse":{"Band":1}}`: ".tapehouse.Band is not an address",
	} {
		if _, err := sdk.ParseDeployments([]byte(registry)); err == nil || err.Error() != message {
			t.Errorf("%s: %v, want %q", registry, err, message)
		}
	}
}

func TestASaleAndACoverCarryTheLimitsAQuoteGives(t *testing.T) {
	d := registry(t, "46630")
	d.Tapehouse["ShortPositions"] = bob
	shorts := sdk.NewClient(nil, d).Shorts()
	tx, err := shorts.Sell("SPY", big.NewInt(1e18), big.NewInt(770), alice)
	if err != nil || tx.To != bob {
		t.Fatalf("sell: %v, %v", tx, err)
	}
	args := unpack(t, &shortpositions.ShortPositionsMetaData, "sell", tx.Data)
	if args[0].([32]byte) != bytes32(t, "SPY") || args[1].(*big.Int).Cmp(big.NewInt(1e18)) != 0 ||
		args[2].(*big.Int).Cmp(big.NewInt(770)) != 0 || args[3].(common.Address) != alice {
		t.Fatalf("sell args %v", args)
	}
	tx, _ = shorts.Cover("SPY", big.NewInt(1), big.NewInt(2), alice)
	unpack(t, &shortpositions.ShortPositionsMetaData, "cover", tx.Data)
}

func TestAnAuthorizationLetsOneAddressActForTheSender(t *testing.T) {
	d := registry(t, "46630")
	tx, err := sdk.NewClient(nil, d).Accounts().SetAuthorization(bob, true)
	if err != nil || tx.To != d.Tapehouse["MarginAccounts"] {
		t.Fatalf("setAuthorization: %v, %v", tx, err)
	}
	if args := unpack(t, &marginaccounts.MarginAccountsMetaData, "setAuthorization", tx.Data); args[0] != bob || args[1] != true {
		t.Fatalf("setAuthorization args %v", args)
	}
	tx, _ = sdk.NewClient(nil, d).Accounts().Repay(sdk.Cross, big.NewInt(5), alice)
	if args := unpack(t, &marginaccounts.MarginAccountsMetaData, "repay", tx.Data); args[0] != sdk.Cross || args[2] != alice {
		t.Fatalf("repay args %v", args)
	}
}

func TestAContractMissingFromTheRegistryIsNamed(t *testing.T) {
	empty, err := sdk.ParseDeployments([]byte(`{"chainId":1}`))
	if err != nil {
		t.Fatal(err)
	}
	client := sdk.NewClient(nil, empty)
	if _, err := client.Shorts().Mark("SPY"); err == nil || err.Error() != "the registry has no .tapehouse.ShortPositions" {
		t.Fatalf("shorts: %v", err)
	}
	if _, err := client.Band().Seal("NVDA"); err == nil || err.Error() != "the registry has no .bandFeeds.NVDA" {
		t.Fatalf("feed: %v", err)
	}
}

func TestAnInvalidArgumentIsAnErrorNeverAPanic(t *testing.T) {
	d := registry(t, "46630")
	d.Tapehouse["ShortPositions"] = bob
	shorts := sdk.NewClient(nil, d).Shorts()
	if _, err := shorts.Sell("SPY", nil, big.NewInt(1), alice); err == nil {
		t.Error("a nil amount packed")
	}
	if _, err := sdk.NewClient(nil, d).Accounts().Deposit(sdk.Cross, bob, big.NewInt(-1), alice); err == nil {
		t.Error("a negative amount packed")
	}
	long := strings.Repeat("A", 33)
	if _, err := sdk.ToBytes32(long); err == nil || err.Error() != `"`+long+`" is longer than 32 bytes` {
		t.Errorf("a 33-byte name: %v", err)
	}
	if _, err := shorts.Mark(long); err == nil {
		t.Error("a 33-byte asset packed")
	}
}

type fileSource struct {
	path  string
	asked [][]string
}

func (s *fileSource) Payload(_ context.Context, feedIDs []string) ([]byte, error) {
	s.asked = append(s.asked, feedIDs)
	text, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	return hexutil.Decode(strings.TrimSpace(string(text)))
}

func TestSignedRedStonePackagesComeFromTheSourceTheIntegratorPlugsIn(t *testing.T) {
	d := registry(t, "46630")
	source := &fileSource{path: "../../stylus/contracts/band/testdata/nvda-24_7.hex"}
	tx, err := sdk.NewClient(nil, d).Band().WritePrices(context.Background(), source, []string{"NVDA---24_7"})
	if err != nil || tx.To != d.Tapehouse["Band"] || len(source.asked) != 1 || source.asked[0][0] != "NVDA---24_7" {
		t.Fatalf("writePrices: %v, %v, %v", tx, err, source.asked)
	}
	payload, _ := source.Payload(context.Background(), nil)
	args := unpack(t, &band.BandMetaData, "writePrices", tx.Data)
	if ids := args[0].([][32]byte); len(ids) != 1 || ids[0] != bytes32(t, "NVDA---24_7") || string(args[1].([]byte)) != string(payload) {
		t.Fatalf("writePrices args %v", args)
	}
}

func TestTheRevertsOfTapehouseTheBandTheIssuerAndTheRouterAreDecoded(t *testing.T) {
	for _, c := range []struct {
		metadata *bind.MetaData
		name     string
		args     []any
		want     string
	}{
		{&shortpositions.ShortPositionsMetaData, "Unauthorized", []any{bob, alice}, "Unauthorized(" + bob.Hex() + ", " + alice.Hex() + ")"},
		{&bandfeed.BandFeedMetaData, "NotSealWindow", []any{uint8(2), uint64(1790200000000)}, "NotSealWindow(2, 1790200000000)"},
		{&stocktoken.StockTokenMetaData, "IsPaused", nil, "IsPaused()"},
		{&band.BandMetaData, "TimestampIsTooOld", []any{big.NewInt(1790119750), big.NewInt(1790300000)}, "TimestampIsTooOld(1790119750, 1790300000)"},
		{&marginaccounts.MarginAccountsMetaData, "DebtCapExceeded", []any{big.NewInt(2), big.NewInt(1)}, "DebtCapExceeded(2, 1)"},
	} {
		revert, ok := sdk.DecodeRevertData(encodeError(t, c.metadata, c.name, c.args...))
		if !ok || revert.Error() != c.want {
			t.Errorf("%s: %v, want %s", c.name, revert, c.want)
		}
	}
	reason, _ := abi.Arguments{{Type: mustType("string")}}.Pack("Too little received")
	revert, ok := sdk.DecodeRevertData(append([]byte{0x08, 0xc3, 0x79, 0xa0}, reason...))
	if !ok || revert.Error() != "Error(Too little received)" {
		t.Errorf("router: %v", revert)
	}
	if _, ok := sdk.DecodeRevertData([]byte{0xde, 0xad, 0xbe, 0xef}); ok {
		t.Error("an unknown selector decoded")
	}
}

type rpcError struct{ data string }

func (e rpcError) Error() string          { return "execution reverted" }
func (e rpcError) ErrorCode() int         { return 3 }
func (e rpcError) ErrorData() interface{} { return e.data }

func TestTheRevertBehindAFailedCallIsDecoded(t *testing.T) {
	data := encodeError(t, &marginaccounts.MarginAccountsMetaData, "Unauthorized", bob, alice)
	revert, ok := sdk.DecodeRevert(rpcError{hexutil.Encode(data)})
	if !ok || revert.Name != "Unauthorized" || revert.Args[0] != bob || revert.Args[1] != alice {
		t.Fatalf("%v", revert)
	}
	if !errors.As(revert, new(rpcError)) {
		t.Fatal("the revert does not wrap the call's error")
	}
	if _, ok := sdk.DecodeRevert(errors.New("not a revert")); ok {
		t.Fatal("decoded a plain error")
	}
}

type quoter struct {
	bind.ContractBackend
	t       *testing.T
	shorts  common.Address
	tokenIn map[string]common.Address
	asked   []string
}

func (q *quoter) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	if *call.To == q.shorts {
		return abi.Arguments{{Type: mustType("uint24")}}.Pack(big.NewInt(500))
	}
	parsed, _ := quoterv2.QuoterV2MetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		q.t.Fatal(err)
	}
	q.asked = append(q.asked, method.Name)
	args, _ := method.Inputs.Unpack(call.Data[4:])
	var fee *big.Int
	var tokenIn common.Address
	if method.Name == "quoteExactInputSingle" {
		params := abi.ConvertType(args[0], new(quoterv2.IQuoterV2QuoteExactInputSingleParams)).(*quoterv2.IQuoterV2QuoteExactInputSingleParams)
		fee, tokenIn = params.Fee, params.TokenIn
	} else {
		params := abi.ConvertType(args[0], new(quoterv2.IQuoterV2QuoteExactOutputSingleParams)).(*quoterv2.IQuoterV2QuoteExactOutputSingleParams)
		fee, tokenIn = params.Fee, params.TokenIn
	}
	if fee.Int64() != 500 || tokenIn != q.tokenIn[method.Name] {
		q.t.Fatalf("%s from %v at %v", method.Name, tokenIn, fee)
	}
	return method.Outputs.Pack(big.NewInt(77_232_802), new(big.Int), uint32(0), new(big.Int))
}

func TestASaleTakesItsQuoteLessTheSlippageAndACoverPaysItsQuotePlusIt(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["ShortPositions"] = bob
	backend := &quoter{t: t, shorts: bob, tokenIn: map[string]common.Address{
		"quoteExactInputSingle": d.Tokens["SPY"], "quoteExactOutputSingle": d.Tokens["USDG"],
	}}
	shorts := sdk.NewClient(backend, d).Shorts()
	sale, err := shorts.QuoteSale(&bind.CallOpts{}, "SPY", big.NewInt(1e17), 50)
	if err != nil || sale.Proceeds.Int64() != 77_232_802 || sale.MinProceeds.Int64() != 76_846_637 {
		t.Fatalf("sale %+v, %v", sale, err)
	}
	cover, err := shorts.QuoteCover(&bind.CallOpts{}, "SPY", big.NewInt(1e17), 50)
	if err != nil || cover.Cost.Int64() != 77_232_802 || cover.MaxCost.Int64() != 77_618_967 {
		t.Fatalf("cover %+v, %v", cover, err)
	}
	if strings.Join(backend.asked, ",") != "quoteExactInputSingle,quoteExactOutputSingle" {
		t.Fatalf("asked %v", backend.asked)
	}
}

func TestASlippageOutside0To10000IsRefused(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["ShortPositions"] = bob
	shorts := sdk.NewClient(nil, d).Shorts()
	if _, err := shorts.QuoteSale(&bind.CallOpts{}, "SPY", big.NewInt(1), 10_001); err == nil || err.Error() != "slippageBps 10001 is outside 0 to 10000" {
		t.Errorf("sale: %v", err)
	}
	if _, err := shorts.QuoteCover(&bind.CallOpts{}, "SPY", big.NewInt(1), -1); err == nil || err.Error() != "slippageBps -1 is outside 0 to 10000" {
		t.Errorf("cover: %v", err)
	}
}

type failingQuoter struct {
	bind.ContractBackend
	shorts common.Address
}

func (q failingQuoter) CallContract(_ context.Context, call ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	if *call.To == q.shorts {
		return abi.Arguments{{Type: mustType("uint24")}}.Pack(big.NewInt(0))
	}
	return nil, rpcError{"0x"}
}

func TestAQuoteThatFailsIsAnErrorNeverAZeroLimit(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["ShortPositions"] = bob
	shorts := sdk.NewClient(failingQuoter{shorts: bob}, d).Shorts()
	if sale, err := shorts.QuoteSale(&bind.CallOpts{}, "SPY", big.NewInt(1), 50); err == nil || err.Error() != "execution reverted" {
		t.Errorf("sale %+v, %v", sale, err)
	}
	if cover, err := shorts.QuoteCover(&bind.CallOpts{}, "SPY", big.NewInt(1), 50); err == nil {
		t.Errorf("cover %+v", cover)
	}
}

type ledger struct {
	bind.ContractBackend
	t      *testing.T
	blocks []*big.Int
}

func (l *ledger) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(7)}, nil
}

func (l *ledger) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	l.blocks = append(l.blocks, block)
	parsed, _ := marginaccounts.MarginAccountsMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		l.t.Fatal(err)
	}
	args, _ := method.Inputs.Unpack(call.Data[4:])
	if args[0] != alice || args[1] != sdk.Cross {
		l.t.Fatalf("%s of %v", method.Name, args)
	}
	answers := map[string]*big.Int{"debt": big.NewInt(100), "premium": big.NewInt(3), "collateral": big.NewInt(2e18), "leverage": big.NewInt(25_000)}
	if method.Name == "liquidationPrice" {
		if args[2] != bytes32(l.t, "SPY") || args[3].(*big.Int).Int64() != 5 {
			l.t.Fatalf("liquidationPrice of %v", args)
		}
		return method.Outputs.Pack(big.NewInt(69_412_000_000))
	}
	return method.Outputs.Pack(answers[method.Name])
}

func TestARepaymentReadsTheDebtAndThePremiumAtOneBlock(t *testing.T) {
	backend := &ledger{t: t}
	accounts := sdk.NewClient(backend, registry(t, "46630")).Accounts()
	repayment, err := accounts.Repayment(&bind.CallOpts{}, alice, sdk.Cross)
	if err != nil || repayment.Debt.Int64() != 100 || repayment.Premium.Int64() != 3 || repayment.Assets.Int64() != 103 {
		t.Fatalf("repayment %+v, %v", repayment, err)
	}
	if len(backend.blocks) != 2 || backend.blocks[0] == nil || backend.blocks[0].Int64() != 7 || backend.blocks[1].Int64() != 7 {
		t.Fatalf("read at %v", backend.blocks)
	}
}

func TestAPositionsCollateralLeverageAndLiquidationPriceAreRead(t *testing.T) {
	accounts := sdk.NewClient(&ledger{t: t}, registry(t, "46630")).Accounts()
	collateral, err := accounts.Collateral(&bind.CallOpts{}, alice, sdk.Cross, bob)
	if err != nil || collateral.Cmp(big.NewInt(2e18)) != 0 {
		t.Fatalf("collateral %v, %v", collateral, err)
	}
	leverage, err := accounts.Leverage(&bind.CallOpts{}, alice, sdk.Cross)
	if err != nil || leverage.Int64() != 25_000 {
		t.Fatalf("leverage %v, %v", leverage, err)
	}
	price, err := accounts.LiquidationPrice(&bind.CallOpts{}, alice, sdk.Cross, "SPY", big.NewInt(5))
	if err != nil || price.Int64() != 69_412_000_000 {
		t.Fatalf("liquidation price %v, %v", price, err)
	}
}

func unpack(t *testing.T, metadata *bind.MetaData, method string, data []byte) []any {
	t.Helper()
	parsed, err := metadata.ParseABI()
	if err != nil {
		t.Fatal(err)
	}
	if string(data[:4]) != string(parsed.Methods[method].ID) {
		t.Fatalf("%x is not %s", data[:4], method)
	}
	args, err := parsed.Methods[method].Inputs.Unpack(data[4:])
	if err != nil {
		t.Fatal(err)
	}
	return args
}

func encodeError(t *testing.T, metadata *bind.MetaData, name string, args ...any) []byte {
	t.Helper()
	parsed, err := metadata.ParseABI()
	if err != nil {
		t.Fatal(err)
	}
	e := parsed.Errors[name]
	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		t.Fatal(err)
	}
	return append(e.ID[:4:4], packed...)
}

func bytes32(t *testing.T, name string) [32]byte {
	t.Helper()
	word, err := sdk.ToBytes32(name)
	if err != nil {
		t.Fatal(err)
	}
	return word
}

func mustType(name string) abi.Type {
	typ, err := abi.NewType(name, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}
