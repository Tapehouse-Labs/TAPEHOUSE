// SPDX-License-Identifier: MIT OR Apache-2.0

package margin

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	testdata       = "../../../stylus/contracts/margin/testdata/"
	parametersFile = "../../../stylus/contracts/margin/parameters.json"
	bandSessions   = "../../../stylus/contracts/band/testdata/session-vectors.json"
)

func load(t *testing.T, path string, into any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, into); err != nil {
		t.Fatal(err)
	}
}

func dec(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("not an integer: %q", s)
	}
	return v
}

func symbols(names []string) [][32]byte {
	out := make([][32]byte, len(names))
	for i, n := range names {
		copy(out[i][:], n)
	}
	return out
}

func TestMatricesAreJudgedAsTheReferenceJudgesThem(t *testing.T) {
	var v struct {
		Cases []struct {
			Name             string
			N                int
			Upper            []uint16
			PositiveDefinite bool
		}
	}
	load(t, testdata+"matrix-vectors.json", &v)
	if len(v.Cases) != 42 {
		t.Fatalf("%d cases", len(v.Cases))
	}
	for _, c := range v.Cases {
		if got := PositiveDefinite(c.N, c.Upper); got != c.PositiveDefinite {
			t.Errorf("%s: got %v", c.Name, got)
		}
	}
}

func TestEverySetIsTheReferenceSet(t *testing.T) {
	type config struct {
		Names        []string
		Volatilities []uint32
		Correlations []uint16
		Gaps         []uint32
		Market       *string
	}
	var v struct {
		Configs map[string]config
		Digests []struct {
			Config  string
			Horizon uint64
			Size    int
			Digest  string
		}
		LaunchRows struct {
			Horizon uint64
			Size    int
			Rows    []struct {
				Index   int
				Returns []int32
			}
		}
	}
	load(t, testdata+"scenario-vectors.json", &v)
	parameters := func(c config) Parameters {
		p := Parameters{Symbols: symbols(c.Names), Volatilities: c.Volatilities, Correlations: c.Correlations, Gaps: c.Gaps, Market: -1}
		for i, n := range c.Names {
			if c.Market != nil && n == *c.Market {
				p.Market = i
			}
		}
		return p
	}
	if len(v.Digests) != 144 {
		t.Fatalf("%d digests", len(v.Digests))
	}
	for _, d := range v.Digests {
		set := NewSet(parameters(v.Configs[d.Config]), d.Size, d.Horizon)
		if got := crypto.Keccak256Hash(set.Encoded()).Hex(); got != d.Digest {
			t.Errorf("%s at %d s over %d points: got %s", d.Config, d.Horizon, d.Size, got)
		}
	}
	set := NewSet(parameters(v.Configs["launch"]), v.LaunchRows.Size, v.LaunchRows.Horizon)
	for _, r := range v.LaunchRows.Rows {
		got := set.Row(r.Index)
		for i := range got {
			if got[i] != r.Returns[i] {
				t.Errorf("row %d: got %v, want %v", r.Index, got, r.Returns)
				break
			}
		}
	}
}

func launchParameters(t *testing.T, names []string) (Parameters, []uint32) {
	t.Helper()
	data, err := os.ReadFile(parametersFile)
	if err != nil {
		t.Fatal(err)
	}
	p, depths, err := Launch(data, names)
	if err != nil {
		t.Fatal(err)
	}
	return p, depths
}

func TestTheRequirementMatchesTheReference(t *testing.T) {
	var v struct {
		Symbols     []string
		Prices      []uint64
		Tokens      map[string]common.Address
		Usdg        common.Address
		Weth        common.Address
		EthUsdRound string
		Pools       map[string]struct {
			Token0    common.Address
			Token1    common.Address
			Fee       uint32
			Observe   string
			Decimals0 uint8
			Decimals1 uint8
		}
		Vectors []struct {
			Label        string
			Pools        string
			Horizon      uint64
			SpansClosure bool
			Quantities   []string
			Requirement  string
			Missing      uint8
		}
		Devnode struct {
			Quantities  []string
			Prices      []uint64
			Requirement string
			Missing     uint8
		}
	}
	load(t, testdata+"requirement-vectors.json", &v)
	parameters, depths := launchParameters(t, v.Symbols)
	int56s, _ := abi.NewType("int56[]", "", nil)
	uint160s, _ := abi.NewType("uint160[]", "", nil)
	observe := abi.Arguments{{Type: int56s}, {Type: uint160s}}
	round := hexutil.MustDecode(v.EthUsdRound)
	ethUsd := new(big.Int).SetBytes(round[32:64])
	pools := make([]Pool, len(v.Symbols))
	for i, n := range v.Symbols {
		p := v.Pools[n]
		stockIsToken0, quoteIsWeth, ok := Classify(p.Token0, p.Token1, v.Tokens[n], v.Usdg, v.Weth)
		if !ok {
			t.Fatalf("%s: the pool does not trade the token", n)
		}
		out, err := observe.Unpack(hexutil.MustDecode(p.Observe))
		if err != nil {
			t.Fatal(err)
		}
		ticks, seconds := out[0].([]*big.Int), out[1].([]*big.Int)
		tick, liquidity, ok := Mean(ticks[0].Int64(), ticks[1].Int64(), seconds[0], seconds[1])
		if !ok {
			t.Fatalf("%s: no mean", n)
		}
		stock, quote := p.Decimals0, p.Decimals1
		if !stockIsToken0 {
			stock, quote = quote, stock
		}
		usd := big.NewInt(100_000_000)
		if quoteIsWeth {
			usd = ethUsd
		}
		terms, ok := Terms(stockIsToken0, stock, quote, tick, liquidity, usd)
		if !ok {
			t.Fatalf("%s: no terms", n)
		}
		terms.Fee = p.Fee
		pools[i] = Pool{Kind: Read, Terms: terms}
	}
	prices := make([]*big.Int, len(v.Prices))
	for i, p := range v.Prices {
		prices[i] = new(big.Int).SetUint64(p)
	}
	if len(v.Vectors) != 120 {
		t.Fatalf("%d vectors", len(v.Vectors))
	}
	for _, c := range v.Vectors {
		quantities := make([]*big.Int, len(c.Quantities))
		for i, q := range c.Quantities {
			quantities[i] = dec(t, q)
		}
		exposures, bad := Exposures(quantities, prices)
		if bad >= 0 {
			t.Fatalf("%s: exposure %d", c.Label, bad)
		}
		with := make([]Pool, len(pools))
		for i, p := range pools {
			switch c.Pools {
			case "read":
				with[i] = p
			case "unread":
				with[i] = Pool{Kind: Unread, Terms: PoolTerms{Fee: p.Terms.Fee}}
			}
		}
		set := NewSet(parameters, 256, c.Horizon)
		got, missing := Requirement(set, exposures, prices, c.SpansClosure, depths, with)
		if got.Cmp(dec(t, c.Requirement)) != 0 || missing != c.Missing {
			t.Errorf("%s (%s, %d s, closure %v): got %v %d, want %s %d", c.Label, c.Pools, c.Horizon, c.SpansClosure, got, missing, c.Requirement, c.Missing)
		}
	}

	three, devDepths := launchParameters(t, []string{"NVDA", "TSLA", "SPY"})
	stub := func(stockIsToken0 bool, quoteDecimals uint8, tick int64, liquidity string, usd int64) Pool {
		perSecond := new(big.Int).Quo(new(big.Int).Lsh(big.NewInt(1), 128), dec(t, liquidity))
		meanTick, meanLiquidity, ok := Mean(0, tick*1_800, new(big.Int), perSecond.Mul(perSecond, big.NewInt(1_800)))
		if !ok {
			t.Fatal("stub mean")
		}
		terms, ok := Terms(stockIsToken0, 18, quoteDecimals, meanTick, meanLiquidity, big.NewInt(usd))
		if !ok {
			t.Fatal("stub terms")
		}
		terms.Fee = 500
		return Pool{Kind: Read, Terms: terms}
	}
	quantities := make([]*big.Int, 3)
	prices = make([]*big.Int, 3)
	for i := range 3 {
		quantities[i] = dec(t, v.Devnode.Quantities[i])
		prices[i] = new(big.Int).SetUint64(v.Devnode.Prices[i])
	}
	exposures, _ := Exposures(quantities, prices)
	pools = []Pool{stub(false, 6, 221_989, "11245526858841909681", 100_000_000), {}, stub(true, 18, -12_513, "16029297629534329325587", 268_330_550_000)}
	got, missing := Requirement(NewSet(three, 256, Horizon), exposures, prices, false, devDepths, pools)
	if got.Cmp(dec(t, v.Devnode.Requirement)) != 0 || missing != v.Devnode.Missing {
		t.Errorf("dev node: got %v %d, want %s %d", got, missing, v.Devnode.Requirement, v.Devnode.Missing)
	}
}

func TestTheShortfallTakesTwoAndFiftySixHundredthsOfTheWorstOf256(t *testing.T) {
	w := newWorst()
	for _, l := range []int64{5, 900, 100, 700, 800, -3} {
		w.push(big.NewInt(l))
	}
	if got := w.shortfall(256).Int64(); got != (100*(900+800)+56*700)/256 {
		t.Errorf("256: %d", got)
	}
	if got := w.shortfall(32).Int64(); got != 900 {
		t.Errorf("32: %d", got)
	}
}

func TestEverySessionTheBandReportsHasItsRegime(t *testing.T) {
	var v struct {
		Sessions []struct {
			Name     string
			NowMs    string `json:"now_ms"`
			Expected struct {
				Known      bool
				Open       bool
				BoundaryMs string `json:"boundary_ms"`
			}
		}
	}
	load(t, bandSessions, &v)
	var seen [4]int
	for _, c := range v.Sessions {
		code := uint8(0)
		if c.Expected.Known {
			code = 1
			if c.Expected.Open {
				code = 2
			}
		}
		r := RegimeAt(code, dec(t, c.Expected.BoundaryMs).Uint64(), dec(t, c.NowMs).Uint64())
		seen[r.Code()]++
	}
	if seen != [4]int{6, 12, 12, 8} {
		t.Errorf("regimes: %v", seen)
	}
}

func TestAMonthOfRealSessionsRaisesARequirementOnlyThroughTheRamp(t *testing.T) {
	var v struct {
		StepMs uint64 `json:"step_ms"`
		Rows   [][6]uint64
	}
	load(t, testdata+"session-timeline.json", &v)
	var seen [4]int
	regimes := make([]Regime, len(v.Rows))
	for i, r := range v.Rows {
		regimes[i] = RegimeAt(uint8(r[1]), r[5], r[0])
		seen[regimes[i].Code()]++
	}
	if seen != [4]int{0, 1_056, 1_880, 135} {
		t.Fatalf("regimes: %v", seen)
	}
	dollar := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	dollars := func(d int64) *big.Int { return new(big.Int).Mul(big.NewInt(d), dollar) }
	open := dollars(1_000)
	for _, closed := range []*big.Int{dollars(1_100), dollars(1_400), dollars(3_000)} {
		for i := 1; i < len(regimes); i++ {
			before := Current(open, closed, new(big.Int), regimes[i-1])
			after := Current(open, closed, new(big.Int), regimes[i])
			if after.Cmp(before) > 0 && regimes[i-1].Code() != 3 && regimes[i].Code() != 3 {
				t.Fatalf("%d: the requirement rose outside the ramp", v.Rows[i][0])
			}
		}
	}
	if got := Current(open, dollars(1_400), new(big.Int), Regime{Kind: Closing, ElapsedMs: RampMs / 4}); got.Cmp(new(big.Int).Add(dollars(1_287), new(big.Int).Quo(dollar, big.NewInt(2)))) != 0 {
		t.Errorf("a quarter into the ramp: %v", got)
	}
	floor := LeverageFloor([]*big.Int{dollars(4_000), dollars(-2_000), dollars(4_000)})
	if floor.Cmp(dollars(2_000)) != 0 {
		t.Errorf("leverage floor: %v", floor)
	}
	if got := Current(big.NewInt(5), new(big.Int), new(big.Int), Regime{Kind: Closed}); got.Int64() != 7 {
		t.Errorf("the buffer rounds up: %v", got)
	}
}

func TestEachAssetsPoolIsNamedAsTheRegistryNamesIt(t *testing.T) {
	data, err := os.ReadFile(parametersFile)
	if err != nil {
		t.Fatal(err)
	}
	pools, err := Pools(data)
	if err != nil || pools["NVDA"] != "NVDA_USDG_500" || pools["SPY"] != "SPY_WETH_500" {
		t.Errorf("%v, %v", pools, err)
	}
}

func TestPairsAreNumberedRowByRow(t *testing.T) {
	k := 0
	for i := range 6 {
		for j := i + 1; j < 6; j++ {
			if Pair(6, i, j) != k {
				t.Fatalf("(%d, %d)", i, j)
			}
			k++
		}
	}
}
