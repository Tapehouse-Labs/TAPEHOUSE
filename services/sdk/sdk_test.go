// SPDX-License-Identifier: MIT OR Apache-2.0

package sdk_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapbackstop"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/gapcover"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/morphobandoracle"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/quoterv2"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/reopeningauction"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocklendingvault"
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
	if arbitrum.ChainlinkSequencer["Uptime"] != common.HexToAddress("0xFdB631F5EE196F0ed6FAa767959853A9F217697D") || len(robinhood.ChainlinkSequencer) != 0 {
		t.Fatalf("sequencer-uptime feeds: %v, %v", arbitrum.ChainlinkSequencer, robinhood.ChainlinkSequencer)
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

func TestTheBasketsAreReadFromTapehouseBaskets(t *testing.T) {
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"tapehouse":{"MarginAccounts":"` + alice.Hex() + `","Baskets":{"PAIR":"` + strings.ToLower(bob.Hex()) + `"}}}`))
	if err != nil || d.Baskets["PAIR"] != bob || len(d.Tapehouse) != 1 || d.Tapehouse["MarginAccounts"] != alice {
		t.Fatalf("baskets: %+v, %v", d, err)
	}
	if testnet := registry(t, "46630"); len(testnet.Baskets) != 0 {
		t.Fatalf("46630 baskets: %v", testnet.Baskets)
	}
	if _, err := sdk.ParseDeployments([]byte(`{"chainId":1,"tapehouse":{"Baskets":{"PAIR":1}}}`)); err == nil {
		t.Fatal("a basket that is not an address parsed")
	}
}

func TestABasketMintsAndRedeemsInKindAndTheAccountsUnwrapOne(t *testing.T) {
	d := registry(t, "46630")
	d.Baskets = map[string]common.Address{"PAIR": bob}
	client := sdk.NewClient(nil, d)
	tx, err := client.Basket("PAIR").Mint(big.NewInt(1e18), alice, []*big.Int{big.NewInt(3), big.NewInt(4)})
	if err != nil || tx.To != bob {
		t.Fatalf("mint: %v, %v", tx, err)
	}
	args := unpack(t, &basket.BasketMetaData, "mint", tx.Data)
	if args[0].(*big.Int).Cmp(big.NewInt(1e18)) != 0 || args[1] != alice || args[2].([]*big.Int)[1].Int64() != 4 {
		t.Fatalf("mint args %v", args)
	}
	tx, _ = client.Basket("PAIR").Redeem(big.NewInt(5), bob, alice)
	if args := unpack(t, &basket.BasketMetaData, "redeem", tx.Data); args[1] != bob || args[2] != alice {
		t.Fatalf("redeem args %v", args)
	}
	tx, err = client.Accounts().Unwrap(alice, "PAIR", big.NewInt(5))
	if err != nil || tx.To != d.Tapehouse["MarginAccounts"] {
		t.Fatalf("unwrap: %v, %v", tx, err)
	}
	if args := unpack(t, &marginaccounts.MarginAccountsMetaData, "unwrap", tx.Data); args[0] != alice || args[1] != bob || args[2].(*big.Int).Int64() != 5 {
		t.Fatalf("unwrap args %v", args)
	}
	if _, err := client.Basket("NONE").Redeem(big.NewInt(1), bob, alice); err == nil || err.Error() != "the registry has no .tapehouse.Baskets.NONE" {
		t.Fatalf("missing basket: %v", err)
	}
	if _, err := client.Basket("PAIR").Mint(nil, alice, nil); err == nil {
		t.Fatal("a nil amount packed")
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
		{&marginaccounts.MarginAccountsMetaData, "BasketFrozen", []any{bob}, "BasketFrozen(" + bob.Hex() + ")"},
		{&basket.BasketMetaData, "PastTarget", []any{bytes32(t, "SPY")}, "PastTarget(0x5350590000000000000000000000000000000000000000000000000000000000)"},
		{&liquidator.LiquidatorMetaData, "NotLiquidatable", []any{bob, sdk.Cross}, "NotLiquidatable(" + bob.Hex() + ", 0x0000000000000000000000000000000000000000000000000000000000000000)"},
		{&stocklendingvault.StockLendingVaultMetaData, "InsufficientLiquidity", []any{big.NewInt(2), big.NewInt(1)}, "InsufficientLiquidity(2, 1)"},
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

func TestMorphoBlueAndTheBandsMorphoOraclesAreRead(t *testing.T) {
	robinhood := registry(t, "4663")
	if robinhood.Morpho["Blue"] != common.HexToAddress("0x9D53d5E3bd5E8d4Cbfa6DB1ca238AEA02E651010") ||
		robinhood.Morpho["AdaptiveCurveIrm"] != common.HexToAddress("0x2BD3d5965B26B51814AC95127B2b80dD6CcC0fa1") {
		t.Fatalf("4663 morpho: %v", robinhood.Morpho)
	}
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"morphoOracles":{"NVDA":"` + bob.Hex() + `"}}`))
	if err != nil || d.MorphoOracles["NVDA"] != bob {
		t.Fatalf("morphoOracles: %v, %v", d, err)
	}
	if _, err := sdk.NewClient(nil, registry(t, "46630")).MorphoOracles().Price(&bind.CallOpts{}, "NVDA"); err == nil ||
		err.Error() != "the registry has no .morphoOracles.NVDA" {
		t.Fatalf("a missing oracle: %v", err)
	}
}

func TestMorphoBluesMarketsAreReadFromMarketsAs32ByteIDs(t *testing.T) {
	id := "0x3a85e619751152991742810df6ec69ce473daef99e28a64ab2340d7b7ccfee49"
	d, err := sdk.ParseDeployments([]byte(`{"chainId":412346,"morpho":{"Blue":"` + alice.Hex() +
		`","Markets":{"NVDA_USDG":"0x` + strings.ToUpper(id[2:]) + `"}}}`))
	if err != nil || len(d.Morpho) != 1 || d.Morpho["Blue"] != alice || d.MorphoMarkets["NVDA_USDG"] != common.HexToHash(id) {
		t.Fatalf("morpho: %v, %v", d, err)
	}
	if markets := registry(t, "4663").MorphoMarkets; len(markets) != 0 {
		t.Fatalf("4663 markets: %v", markets)
	}
	for _, wrong := range []string{`"` + alice.Hex() + `"`, `"` + id + `00"`, `"` + id[:65] + `"`} {
		_, err := sdk.ParseDeployments([]byte(`{"chainId":1,"morpho":{"Markets":{"NVDA_USDG":` + wrong + `}}}`))
		if err == nil || err.Error() != ".morpho.Markets.NVDA_USDG is not a 32-byte id" {
			t.Fatalf("%s: %v", wrong, err)
		}
	}
}

var bandAddress = common.HexToAddress("0xa70118d3324D90532E7D2854627b13CacE305641")

type oracleChain struct {
	bind.ContractBackend
	t      *testing.T
	price  []byte
	revert []byte
	halt   []any
	step   uint8
	blocks []*big.Int
}

func (o *oracleChain) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(7)}, nil
}

func (o *oracleChain) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	o.blocks = append(o.blocks, block)
	if *call.To == bandAddress {
		parsed, _ := band.BandMetaData.ParseABI()
		method, err := parsed.MethodById(call.Data)
		if err != nil || method.Name != "corporateAction" {
			o.t.Fatalf("band %v, %v", method, err)
		}
		return method.Outputs.Pack(o.step, uint64(0), new(big.Int), new(big.Int))
	}
	parsed, _ := morphobandoracle.MorphoBandOracleMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil || *call.To != bob {
		o.t.Fatalf("oracle %v at %v, %v", method, call.To, err)
	}
	switch method.Name {
	case "price":
		if o.revert != nil {
			return nil, rpcError{hexutil.Encode(o.revert)}
		}
		return o.price, nil
	case "halt":
		return method.Outputs.Pack(o.halt...)
	case "band":
		return method.Outputs.Pack(bandAddress)
	case "symbol":
		return method.Outputs.Pack(bytes32(o.t, "NVDA"))
	case "collateralToken", "owner":
		return method.Outputs.Pack(alice)
	case "loanToken":
		return method.Outputs.Pack(bob)
	case "scaleFactor":
		return method.Outputs.Pack(big.NewInt(1e16))
	}
	o.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func oracles(t *testing.T, chain *oracleChain) *sdk.MorphoOracles {
	d := registry(t, "46630")
	d.MorphoOracles["NVDA"] = bob
	chain.t = t
	if chain.halt == nil {
		chain.halt = []any{false, uint64(0), uint64(0), false}
	}
	return sdk.NewClient(chain, d).MorphoOracles()
}

func TestAnOracleAnswersItsPriceAtOneBlock(t *testing.T) {
	price, _ := abi.Arguments{{Type: uint256Type}}.Pack(new(big.Int).Mul(big.NewInt(765_517_754_420), big.NewInt(1e15)))
	chain := &oracleChain{price: price}
	answer, err := oracles(t, chain).Price(&bind.CallOpts{}, "NVDA")
	if err != nil || answer.Price == nil || answer.Price.String() != "765517754420000000000000000" || answer.NoPrice != "" {
		t.Fatalf("price %+v, %v", answer, err)
	}
	if len(chain.blocks) != 1 || chain.blocks[0].Int64() != 7 {
		t.Fatalf("read at %v", chain.blocks)
	}
}

func TestNoAnswerWithNoHaltPauseOrStepIsAStaleBandNeverAPriceOfZero(t *testing.T) {
	data := encodeError(t, &morphobandoracle.MorphoBandOracleMetaData, "NoAnswer", bytes32(t, "NVDA"))
	chain := &oracleChain{revert: data}
	answer, err := oracles(t, chain).Price(&bind.CallOpts{}, "NVDA")
	if err != nil || answer.Price != nil || answer.NoPrice != sdk.NoPriceStale || answer.Revert.Name != "NoAnswer" {
		t.Fatalf("stale %+v, %v", answer, err)
	}
	for _, block := range chain.blocks {
		if block.Int64() != 7 {
			t.Fatalf("read at %v", chain.blocks)
		}
	}
	if len(chain.blocks) != 5 {
		t.Fatalf("%d reads", len(chain.blocks))
	}
}

func TestNoAnswerUnderASignedHaltThePauseOrAnUnconfirmedStepIsAHalt(t *testing.T) {
	data := encodeError(t, &morphobandoracle.MorphoBandOracleMetaData, "NoAnswer", bytes32(t, "NVDA"))
	for _, chain := range []*oracleChain{
		{revert: data, halt: []any{true, uint64(1_790_200_000), uint64(1_790_196_400), false}},
		{revert: data, halt: []any{false, uint64(0), uint64(0), true}},
		{revert: data, step: 2},
	} {
		answer, err := oracles(t, chain).Price(&bind.CallOpts{}, "NVDA")
		if err != nil || answer.Price != nil || answer.NoPrice != sdk.NoPriceHalted {
			t.Fatalf("halted %+v, %v", answer, err)
		}
	}
}

func TestSequencerNotSettledIsNoPriceAndAnyOtherRevertAnError(t *testing.T) {
	data := encodeError(t, &morphobandoracle.MorphoBandOracleMetaData, "SequencerNotSettled")
	answer, err := oracles(t, &oracleChain{revert: data}).Price(&bind.CallOpts{}, "NVDA")
	if err != nil || answer.Price != nil || answer.NoPrice != sdk.NoPriceSequencerNotSettled || answer.Revert.Name != "SequencerNotSettled" {
		t.Fatalf("sequencer %+v, %v", answer, err)
	}
	if answer, err := oracles(t, &oracleChain{revert: []byte{}}).Price(&bind.CallOpts{}, "NVDA"); err == nil {
		t.Fatalf("an empty revert answered %+v", answer)
	}
}

func TestAnOraclesBandSymbolTokensScaleOwnerAndHaltAreRead(t *testing.T) {
	chain := &oracleChain{halt: []any{true, uint64(1_790_200_000), uint64(1_790_196_400), false}}
	oracle, err := oracles(t, chain).Oracle(&bind.CallOpts{}, "NVDA")
	if err != nil || oracle.Address != bob || oracle.Band != bandAddress || oracle.Symbol != bytes32(t, "NVDA") ||
		oracle.CollateralToken != alice || oracle.LoanToken != bob || oracle.ScaleFactor.Int64() != 1e16 || oracle.Owner != alice {
		t.Fatalf("oracle %+v, %v", oracle, err)
	}
	if len(chain.blocks) != 6 || chain.blocks[5].Int64() != 7 {
		t.Fatalf("read at %v", chain.blocks)
	}
	halt, err := oracles(t, chain).Halt(&bind.CallOpts{}, "NVDA")
	if err != nil || !halt.SignedHalt || halt.Until != 1_790_200_000 || halt.IssuedAt != 1_790_196_400 || halt.OraclePaused {
		t.Fatalf("halt %+v, %v", halt, err)
	}
}

func TestARepointsAssetMismatchAndBandSetAreDecoded(t *testing.T) {
	revert, ok := sdk.DecodeRevertData(encodeError(t, &morphobandoracle.MorphoBandOracleMetaData, "AssetMismatch", bandAddress, bob))
	if !ok || revert.Error() != "AssetMismatch("+bandAddress.Hex()+", "+bob.Hex()+")" {
		t.Fatalf("%v", revert)
	}
	parsed, _ := morphobandoracle.MorphoBandOracleMetaData.ParseABI()
	set, err := morphobandoracle.NewMorphoBandOracle().UnpackBandSetEvent(&types.Log{
		Address: bob,
		Topics:  []common.Hash{parsed.Events["BandSet"].ID, common.BytesToHash(alice[:]), common.BytesToHash(bandAddress[:])},
	})
	if err != nil || set.PreviousBand != alice || set.NewBand != bandAddress {
		t.Fatalf("BandSet %+v, %v", set, err)
	}
}

type shelf struct {
	bind.ContractBackend
	t      *testing.T
	basket common.Address
	blocks []*big.Int
}

func (s *shelf) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(9)}, nil
}

func (s *shelf) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	s.blocks = append(s.blocks, block)
	metadata := &marginaccounts.MarginAccountsMetaData
	if *call.To == s.basket {
		metadata = &basket.BasketMetaData
	}
	parsed, _ := metadata.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		s.t.Fatal(err)
	}
	args, _ := method.Inputs.Unpack(call.Data[4:])
	symbols := [][32]byte{bytes32(s.t, "NVDA"), bytes32(s.t, "SPY")}
	switch method.Name {
	case "stocks", "components":
		return method.Outputs.Pack(symbols, []common.Address{alice, bob})
	case "inBaskets":
		if args[0] != alice || args[1] != sdk.Cross {
			s.t.Fatalf("inBaskets of %v", args)
		}
		return method.Outputs.Pack([]*big.Int{big.NewInt(10), big.NewInt(5)})
	case "pendingTarget":
		return method.Outputs.Pack([]*big.Int{big.NewInt(1), big.NewInt(2)}, uint64(1_790_604_800))
	case "target":
		return method.Outputs.Pack([]*big.Int{big.NewInt(7), big.NewInt(8)})
	}
	if args[0].(*big.Int).Int64() != 3 {
		s.t.Fatalf("%s of %v", method.Name, args)
	}
	if method.Name == "previewMint" {
		return method.Outputs.Pack([]*big.Int{big.NewInt(4), big.NewInt(2)})
	}
	return method.Outputs.Pack([]*big.Int{big.NewInt(3), big.NewInt(1)})
}

func TestABasketsComponentsPreviewsAndTargetsAreRead(t *testing.T) {
	d := registry(t, "46630")
	d.Baskets = map[string]common.Address{"PAIR": bob}
	backend := &shelf{t: t, basket: bob}
	pair := sdk.NewClient(backend, d).Basket("PAIR")
	components, err := pair.Components(&bind.CallOpts{})
	if err != nil || strings.Join(components.Assets, ",") != "NVDA,SPY" || components.Tokens[1] != bob {
		t.Fatalf("components %+v, %v", components, err)
	}
	for name, read := range map[string]func(*bind.CallOpts, *big.Int) ([]*big.Int, error){
		"4,2": pair.PreviewMint,
		"3,1": pair.PreviewRedeem,
	} {
		amounts, err := read(&bind.CallOpts{}, big.NewInt(3))
		if err != nil || amounts[0].String()+","+amounts[1].String() != name {
			t.Errorf("preview %v, %v, want %s", amounts, err, name)
		}
	}
	units, err := pair.Target(&bind.CallOpts{})
	if err != nil || units[1].Int64() != 8 {
		t.Fatalf("target %v, %v", units, err)
	}
	pending, err := pair.PendingTarget(&bind.CallOpts{BlockNumber: big.NewInt(4)})
	if err != nil || pending.Units[0].Int64() != 1 || pending.EffectiveAt != 1_790_604_800 {
		t.Fatalf("pending %+v, %v", pending, err)
	}
	if last := backend.blocks[len(backend.blocks)-1]; last.Int64() != 4 {
		t.Fatalf("read at %v", last)
	}
}

func TestWhatAPositionHoldsThroughItsBasketsIsReadByAssetAtOneBlock(t *testing.T) {
	backend := &shelf{t: t}
	held, err := sdk.NewClient(backend, registry(t, "46630")).Accounts().InBaskets(&bind.CallOpts{}, alice, sdk.Cross)
	if err != nil || len(held) != 2 || held["NVDA"].Int64() != 10 || held["SPY"].Int64() != 5 {
		t.Fatalf("inBaskets %v, %v", held, err)
	}
	if len(backend.blocks) != 2 || backend.blocks[0].Int64() != 9 || backend.blocks[1].Int64() != 9 {
		t.Fatalf("read at %v", backend.blocks)
	}
}

func coverLayer() sdk.Layer {
	return sdk.Layer{Notional: big.NewInt(10_000_000_000), DeductibleBps: big.NewInt(218), LimitBps: big.NewInt(1_218)}
}

func TestAGapCoverPurchaseCarriesItsLayerItsPremiumLimitAndItsHolder(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["GapCover"] = bob
	cover := sdk.NewClient(nil, d).GapCover()
	tx, err := cover.Buy("NVDA", coverLayer(), big.NewInt(13_908_962), alice)
	if err != nil || tx.To != bob {
		t.Fatalf("buy: %v, %v", tx, err)
	}
	args := unpack(t, &gapcover.GapCoverMetaData, "buy", tx.Data)
	if args[0].([32]byte) != bytes32(t, "NVDA") || args[1].(*big.Int).Int64() != 10_000_000_000 ||
		args[2].(*big.Int).Int64() != 218 || args[3].(*big.Int).Int64() != 1_218 ||
		args[4].(*big.Int).Int64() != 13_908_962 || args[5].(common.Address) != alice {
		t.Fatalf("buy args %v", args)
	}
	reference, last := new(big.Int).Lsh(big.NewInt(1), 64), new(big.Int).Lsh(big.NewInt(1), 64)
	reference.Or(reference, big.NewInt(7))
	last.Or(last, big.NewInt(8))
	tx, _ = cover.Settle("SPY", 1_790_985_600_000, reference, last)
	args = unpack(t, &gapcover.GapCoverMetaData, "settle", tx.Data)
	if args[0].([32]byte) != bytes32(t, "SPY") || args[1].(uint64) != 1_790_985_600_000 ||
		args[2].(*big.Int).Cmp(reference) != 0 || args[3].(*big.Int).Cmp(last) != 0 {
		t.Fatalf("settle args %v", args)
	}
	for method, packed := range map[string]func() (sdk.Tx, error){
		"observe": func() (sdk.Tx, error) { return cover.Observe("SPY", 1) },
		"release": func() (sdk.Tx, error) { return cover.Release(big.NewInt(3)) },
		"claim":   func() (sdk.Tx, error) { return cover.Claim(alice) },
		"record":  cover.Record,
		"measure": func() (sdk.Tx, error) { return cover.Measure("SPY") },
		"deposit": func() (sdk.Tx, error) { return cover.Deposit(big.NewInt(5), alice) },
		"redeem":  func() (sdk.Tx, error) { return cover.Redeem(big.NewInt(5), alice, alice) },
	} {
		tx, err := packed()
		if err != nil || tx.To != bob {
			t.Fatalf("%s: %v, %v", method, tx, err)
		}
		unpack(t, &gapcover.GapCoverMetaData, method, tx.Data)
	}
	if _, err := cover.Buy("NVDA", sdk.Layer{}, big.NewInt(1), alice); err == nil {
		t.Error("a nil layer packed")
	}
	empty, _ := sdk.ParseDeployments([]byte(`{"chainId":1}`))
	if _, err := sdk.NewClient(nil, empty).GapCover().Record(); err == nil || err.Error() != "the registry has no .tapehouse.GapCover" {
		t.Fatalf("missing: %v", err)
	}
}

type coverBackend struct {
	bind.ContractBackend
	t      *testing.T
	status uint8
	blocks []*big.Int
}

func (c *coverBackend) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(7)}, nil
}

func (c *coverBackend) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	c.blocks = append(c.blocks, block)
	parsed, _ := gapcover.GapCoverMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil || *call.To != bob {
		c.t.Fatalf("%v to %v", err, call.To)
	}
	switch method.Name {
	case "quote":
		return method.Outputs.Pack(big.NewInt(13_908_962))
	case "pricingGap":
		return method.Outputs.Pack(big.NewInt(279_862), big.NewInt(111_945))
	case "series":
		return method.Outputs.Pack(big.NewInt(10_000_000_000), uint64(22_244_729_849), uint64(22_280_257_368), uint16(0), c.status, false)
	case "reopenOf":
		return method.Outputs.Pack(big.NewInt(1_789_948_800_000))
	case "sales":
		if c.status == 0 {
			return method.Outputs.Pack(uint64(1_790_985_600_000), uint64(1_790_985_600_000))
		}
		return method.Outputs.Pack(uint64(0), uint64(0))
	case "covers":
		return method.Outputs.Pack(common.Address{}, uint64(0), uint16(0), uint16(0), [32]byte{}, new(big.Int), new(big.Int))
	}
	c.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func TestAGapCoverQuoteReadsThePremiumAndGivesTheWritersUsdgItReserves(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["GapCover"] = bob
	cover := sdk.NewClient(&coverBackend{t: t}, d).GapCover()
	quote, err := cover.Quote(&bind.CallOpts{}, "NVDA", coverLayer())
	if err != nil || quote.Premium.Int64() != 13_908_962 || quote.Reserve.Int64() != 1_000_000_000 {
		t.Fatalf("quote %+v, %v", quote, err)
	}
	priced, err := cover.PricingGap(&bind.CallOpts{}, "NVDA")
	if err != nil || priced.Gap.Int64() != 279_862 || priced.WeekMove.Int64() != 111_945 {
		t.Fatalf("pricing gap %+v, %v", priced, err)
	}
}

func TestTheSalesReadTheClosureOnSaleAndWhenTheyEnd(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["GapCover"] = bob
	backend := &coverBackend{t: t}
	cover := sdk.NewClient(backend, d).GapCover()
	sales, open, err := cover.Sales(&bind.CallOpts{})
	if err != nil || !open || sales.ClosesMs != 1_790_985_600_000 || sales.EndsMs != 1_790_985_600_000 {
		t.Fatalf("sales %+v, %v, %v", sales, open, err)
	}
	backend.status = 1
	if _, open, err := cover.Sales(&bind.CallOpts{}); open || err != nil {
		t.Fatalf("closed sales read as open: %v", err)
	}
}

func TestASeriesReadsWhereItStandsAndItsSettlementAtOneBlock(t *testing.T) {
	d := registry(t, "4663")
	d.Tapehouse["GapCover"] = bob
	backend := &coverBackend{t: t, status: 1}
	cover := sdk.NewClient(backend, d).GapCover()
	series, err := cover.Series(&bind.CallOpts{}, "NVDA", 1_789_776_000_000)
	if err != nil || series.Status != sdk.SeriesSettled || series.ReferencePrice != 22_244_729_849 ||
		series.Price != 22_280_257_368 || series.Flagged || series.ReopenMs.Int64() != 1_789_948_800_000 ||
		series.Notional.Int64() != 10_000_000_000 {
		t.Fatalf("series %+v, %v", series, err)
	}
	if len(backend.blocks) != 2 || backend.blocks[0].Int64() != 7 || backend.blocks[1].Int64() != 7 {
		t.Fatalf("read at %v", backend.blocks)
	}
	backend.status = 2
	if series, _ := cover.Series(&bind.CallOpts{}, "NVDA", 1); series.Status != sdk.SeriesVoid {
		t.Fatalf("void read as %v", series.Status)
	}
	backend.status = 3
	if _, err := cover.Series(&bind.CallOpts{}, "NVDA", 1); err == nil || err.Error() != "unknown series status 3" {
		t.Fatalf("status 3: %v", err)
	}
	if _, held, err := cover.Cover(&bind.CallOpts{}, big.NewInt(1)); held || err != nil {
		t.Fatalf("a released cover read as held: %v", err)
	}
}

func TestACoverPaysTheFallBeyondItsDeductibleUpToItsLimit(t *testing.T) {
	layer := sdk.Layer{Notional: big.NewInt(10_000_000_000), DeductibleBps: big.NewInt(500), LimitBps: big.NewInt(1_500)}
	for price, paid := range map[uint64]int64{210e8: 0, 200e8: 0, 190e8: 0, 185e8: 250_000_000, 100e8: 1_000_000_000} {
		if got := sdk.Payout(layer, 200e8, price); got.Int64() != paid {
			t.Errorf("at %d: %v, want %d", price, got, paid)
		}
	}
	revert, ok := sdk.DecodeRevertData(encodeError(t, &gapcover.GapCoverMetaData, "NoCapacity", big.NewInt(2), big.NewInt(1)))
	if !ok || revert.Error() != "NoCapacity(2, 1)" {
		t.Fatalf("NoCapacity: %v", revert)
	}
	if revert, ok := sdk.DecodeRevertData(encodeError(t, &gapcover.GapCoverMetaData, "SalesClosed")); !ok || revert.Name != "SalesClosed" {
		t.Fatalf("SalesClosed: %v", revert)
	}
	stale, ok := sdk.DecodeRevertData(encodeError(t, &gapcover.GapCoverMetaData, "StaleReference", big.NewInt(267), big.NewInt(268)))
	if !ok || stale.Error() != "StaleReference(267, 268)" {
		t.Fatalf("StaleReference: %v", stale)
	}
	early, ok := sdk.DecodeRevertData(encodeError(t, &gapcover.GapCoverMetaData, "TooEarlyToMeasure", big.NewInt(1_790_899_200_000)))
	if !ok || early.Error() != "TooEarlyToMeasure(1790899200000)" {
		t.Fatalf("TooEarlyToMeasure: %v", early)
	}
}

func TestTheKeepersCallsGoToTheRegistrysContracts(t *testing.T) {
	d := registry(t, "46630")
	d.StockLending = map[string]common.Address{"SPY": bob}
	d.Baskets = map[string]common.Address{"PAIR": alice}
	d.BandFeeds = map[string]common.Address{"SPY": alice}
	client := sdk.NewClient(nil, d)
	spy := bytes32(t, "SPY")
	for _, c := range []struct {
		pack     func() (sdk.Tx, error)
		to       common.Address
		metadata *bind.MetaData
		method   string
		args     []any
	}{
		{func() (sdk.Tx, error) { return client.Liquidator().Start(alice, sdk.Cross) }, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "start", []any{alice, sdk.Cross}},
		{func() (sdk.Tx, error) { return client.Liquidator().Stop(alice, spy) }, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "stop", []any{alice, spy}},
		{func() (sdk.Tx, error) {
			return client.Liquidator().Buy(alice, sdk.Cross, bob, big.NewInt(5), big.NewInt(7), alice)
		}, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "buy", []any{alice, sdk.Cross, bob, big.NewInt(5), big.NewInt(7), alice}},
		{func() (sdk.Tx, error) { return client.Liquidator().SettleCash(alice, sdk.Cross) }, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "settleCash", []any{alice, sdk.Cross}},
		{func() (sdk.Tx, error) { return client.Liquidator().WriteOff(alice, sdk.Cross) }, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "writeOff", []any{alice, sdk.Cross}},
		{func() (sdk.Tx, error) { return client.Liquidator().Recall(alice, sdk.Cross, bob) }, d.Tapehouse["Liquidator"], &liquidator.LiquidatorMetaData, "recall", []any{alice, sdk.Cross, bob}},
		{func() (sdk.Tx, error) { return client.Backstop().Claim() }, d.Tapehouse["GapBackstop"], &gapbackstop.GapBackstopMetaData, "claim", []any{}},
		{func() (sdk.Tx, error) { return client.Backstop().Cover(alice, sdk.Cross) }, d.Tapehouse["GapBackstop"], &gapbackstop.GapBackstopMetaData, "cover", []any{alice, sdk.Cross}},
		{func() (sdk.Tx, error) { return client.ReopeningAuction().Enroll(alice, sdk.Cross, "SPY") }, d.Tapehouse["ReopeningAuction"], &reopeningauction.ReopeningAuctionMetaData, "enroll", []any{alice, sdk.Cross, spy}},
		{func() (sdk.Tx, error) { return client.ReopeningAuction().Clear("SPY", 1790775000000, big.NewInt(9)) }, d.Tapehouse["ReopeningAuction"], &reopeningauction.ReopeningAuctionMetaData, "clear", []any{spy, uint64(1790775000000), big.NewInt(9)}},
		{func() (sdk.Tx, error) { return client.ReopeningAuction().Claim("SPY", 1790775000000, big.NewInt(1)) }, d.Tapehouse["ReopeningAuction"], &reopeningauction.ReopeningAuctionMetaData, "claim", []any{spy, uint64(1790775000000), big.NewInt(1)}},
		{func() (sdk.Tx, error) { return client.ReopeningAuction().Forfeit("SPY", 1790775000000, spy) }, d.Tapehouse["ReopeningAuction"], &reopeningauction.ReopeningAuctionMetaData, "forfeit", []any{spy, uint64(1790775000000), spy}},
		{func() (sdk.Tx, error) { return client.Accounts().AccruePremium() }, d.Tapehouse["MarginAccounts"], &marginaccounts.MarginAccountsMetaData, "accruePremium", []any{}},
		{func() (sdk.Tx, error) { return client.Accounts().Sync("SPY") }, d.Tapehouse["MarginAccounts"], &marginaccounts.MarginAccountsMetaData, "sync", []any{spy}},
		{func() (sdk.Tx, error) { return client.Accounts().Clear(alice, sdk.Cross, "SPY") }, d.Tapehouse["MarginAccounts"], &marginaccounts.MarginAccountsMetaData, "clear", []any{alice, sdk.Cross, spy}},
		{func() (sdk.Tx, error) { return client.Accounts().Settle(alice, sdk.Cross, bob) }, d.Tapehouse["MarginAccounts"], &marginaccounts.MarginAccountsMetaData, "settle", []any{alice, sdk.Cross, bob}},
		{func() (sdk.Tx, error) { return client.LendingVault("SPY").BuyIn(3) }, bob, &stocklendingvault.StockLendingVaultMetaData, "buyIn", []any{big.NewInt(3)}},
		{func() (sdk.Tx, error) { return client.Band().SyncMultiplier("SPY") }, d.Tapehouse["Band"], &band.BandMetaData, "syncMultiplier", []any{spy}},
		{func() (sdk.Tx, error) {
			return client.Band().WriteHalt("SPY", true, 1790775000, 1790775900, []byte{1, 2})
		}, d.Tapehouse["Band"], &band.BandMetaData, "writeHalt", []any{spy, true, uint64(1790775000), uint64(1790775900), []byte{1, 2}}},
		{func() (sdk.Tx, error) { return client.Band().Seal("SPY") }, alice, &bandfeed.BandFeedMetaData, "seal", []any{}},
		{func() (sdk.Tx, error) { return client.GapCover().Void("SPY", 1790460000000) }, d.Tapehouse["GapCover"], &gapcover.GapCoverMetaData, "void", []any{spy, uint64(1790460000000)}},
		{func() (sdk.Tx, error) {
			return client.Basket("PAIR").Rebalance([]*big.Int{big.NewInt(2), big.NewInt(0)}, []*big.Int{big.NewInt(0), big.NewInt(1)}, bob)
		}, alice, &basket.BasketMetaData, "rebalance", []any{[]*big.Int{big.NewInt(2), big.NewInt(0)}, []*big.Int{big.NewInt(0), big.NewInt(1)}, bob}},
	} {
		tx, err := c.pack()
		if c.to == (common.Address{}) {
			if err == nil || !strings.Contains(err.Error(), "the registry has no") {
				t.Errorf("%s: %v", c.method, err)
			}
			continue
		}
		if err != nil || tx.To != c.to {
			t.Errorf("%s: %v, %v", c.method, tx, err)
			continue
		}
		if args := unpack(t, c.metadata, c.method, tx.Data); fmt.Sprint(args) != fmt.Sprint(c.args) {
			t.Errorf("%s args %v, want %v", c.method, args, c.args)
		}
	}
	if _, err := client.LendingVault("NVDA").BuyIn(1); err == nil || err.Error() != "the registry has no .tapehouse.StockLending.NVDA" {
		t.Fatalf("a vault the registry does not name: %v", err)
	}
	if _, err := client.Liquidator().Buy(alice, sdk.Cross, bob, nil, big.NewInt(1), alice); err == nil {
		t.Fatal("a nil amount packed")
	}
}

func TestTheClearingPriceIsTheOneTheReopeningAuctionAccepts(t *testing.T) {
	bid := func(quantity, price int64) reopeningauction.ReopeningAuctionBid {
		return reopeningauction.ReopeningAuctionBid{Quantity: big.NewInt(quantity), Price: big.NewInt(price)}
	}
	for _, c := range []struct {
		name   string
		bids   []reopeningauction.ReopeningAuctionBid
		supply int64
		want   int64
		ok     bool
	}{
		{"no bid clears at the floor", nil, 10, 90, true},
		{"the highest price the bids at or above it take the supply at", []reopeningauction.ReopeningAuctionBid{bid(4, 120), bid(5, 110), bid(3, 100), bid(9, 95)}, 10, 100, true},
		{"the bids above it may fill exactly", []reopeningauction.ReopeningAuctionBid{bid(4, 120), bid(6, 110), bid(3, 100)}, 10, 110, true},
		{"several bids at one price", []reopeningauction.ReopeningAuctionBid{bid(2, 120), bid(4, 105), bid(5, 105), bid(1, 99)}, 10, 105, true},
		{"undersubscribed clears at the lowest bid", []reopeningauction.ReopeningAuctionBid{bid(2, 120), bid(3, 101)}, 10, 101, true},
		{"bids and no supply cannot clear", []reopeningauction.ReopeningAuctionBid{bid(2, 120)}, 0, 0, false},
	} {
		price, ok := sdk.ClearingPrice(c.bids, big.NewInt(c.supply), big.NewInt(90))
		if ok != c.ok || (ok && price.Int64() != c.want) {
			t.Errorf("%s: %v, %v; want %d", c.name, price, ok, c.want)
		}
	}
}

func TestTheBandsHaltAndMultiplierErrorsAreDecoded(t *testing.T) {
	for _, c := range []struct {
		name string
		args []any
		want string
	}{
		{"HaltNotNewer", []any{bytes32(t, "NVDA"), uint64(1790775000), uint64(1790775000)}, "HaltNotNewer(0x4e56444100000000000000000000000000000000000000000000000000000000, 1790775000, 1790775000)"},
		{"HaltOutsideWindow", []any{uint64(1), uint64(3602), uint64(1)}, "HaltOutsideWindow(1, 3602, 1)"},
		{"NoToken", []any{bytes32(t, "TSLA")}, "NoToken(0x54534c4100000000000000000000000000000000000000000000000000000000)"},
		{"UnknownAsset", []any{bytes32(t, "X")}, "UnknownAsset(0x5800000000000000000000000000000000000000000000000000000000000000)"},
		{"PackageNotNewer", []any{bytes32(t, "X"), uint64(2), uint64(2)}, "PackageNotNewer(0x5800000000000000000000000000000000000000000000000000000000000000, 2, 2)"},
	} {
		revert, ok := sdk.DecodeRevertData(encodeError(t, &band.BandMetaData, c.name, c.args...))
		if !ok || revert.Error() != c.want {
			t.Errorf("%s: %v, want %s", c.name, revert, c.want)
		}
	}
}

type auctionChain struct {
	bind.ContractBackend
	t       *testing.T
	started uint64
	closed  bool
	now     uint64
	blocks  map[string]*big.Int
}

func (a *auctionChain) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	if number == nil {
		number = big.NewInt(11)
	}
	return &types.Header{Number: number, Time: a.now}, nil
}

func (a *auctionChain) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	parsed, _ := liquidator.LiquidatorMetaData.ParseABI()
	method, err := parsed.MethodById(call.Data)
	if err != nil {
		a.t.Fatal(err)
	}
	a.blocks[method.Name] = block
	switch method.Name {
	case "auctions":
		return method.Outputs.Pack(a.started, a.closed)
	case "shortfall":
		return method.Outputs.Pack(big.NewInt(-1), big.NewInt(2), true, false)
	case "OPEN_AUCTION_LIFETIME":
		return method.Outputs.Pack(big.NewInt(3600))
	}
	a.t.Fatalf("unexpected %s", method.Name)
	return nil, nil
}

func TestAnAuctionRunsOnlyInTheMarketsStateAndWithinItsLifetime(t *testing.T) {
	chain := &auctionChain{t: t, started: 1790775000, now: 1790778599, blocks: map[string]*big.Int{}}
	liquidation := sdk.NewClient(chain, registry(t, "46630")).Liquidator()
	running, err := liquidation.Running(&bind.CallOpts{}, alice, sdk.Cross)
	if err != nil || !running || chain.blocks["auctions"].Int64() != 11 || chain.blocks["OPEN_AUCTION_LIFETIME"].Int64() != 11 {
		t.Fatalf("an auction an hour old less a second: %v, %v, %v", running, err, chain.blocks)
	}
	chain.now = 1790778600
	if running, err := liquidation.Running(&bind.CallOpts{}, alice, sdk.Cross); err != nil || running {
		t.Fatalf("an auction past its lifetime runs: %v", err)
	}
	chain.now, chain.closed = 1790775001, true
	if running, err := liquidation.Running(&bind.CallOpts{}, alice, sdk.Cross); err != nil || running {
		t.Fatalf("an auction started closed runs while open: %v", err)
	}
	chain.started = 0
	if running, err := liquidation.Running(&bind.CallOpts{}, alice, sdk.Cross); err != nil || running {
		t.Fatalf("no auction runs: %v", err)
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

var uint256Type = mustType("uint256")

func mustType(name string) abi.Type {
	typ, err := abi.NewType(name, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}
