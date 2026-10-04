// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tapehouse/tapehouse/services/internal/margin"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/aggregator"
	bandbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/band"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/liquidator"
	marginbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/margin"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	shortbinding "github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/uniswapv3pool"
)

// chain is a test double of a node: each contract answers the reads of the ABI its SDK binding was generated from,
// at the block a call names, and the head moves on after every call.
type chain struct {
	t         *testing.T
	mu        sync.Mutex
	head      uint64
	time      uint64
	contracts map[common.Address]*contract
	calls     []*big.Int
}

type contract struct {
	abi   abi.ABI
	reads map[string]func(args []any) ([]any, error)
}

func newChain(t *testing.T) *chain {
	return &chain{t: t, head: 100, time: 1_790_775_000, contracts: map[common.Address]*contract{}}
}

func (c *chain) add(address common.Address, metadata *bind.MetaData) *contract {
	parsed, err := metadata.ParseABI()
	if err != nil {
		c.t.Fatal(err)
	}
	k := &contract{abi: *parsed, reads: map[string]func([]any) ([]any, error){}}
	c.contracts[address] = k
	return k
}

func (k *contract) on(method string, answer func(args []any) []any) {
	k.reads[method] = func(args []any) ([]any, error) { return answer(args), nil }
}

func (k *contract) fail(method string, err error) {
	k.reads[method] = func([]any) ([]any, error) { return nil, err }
}

func (c *chain) CallContract(_ context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls = append(c.calls, block)
	c.head++
	k, ok := c.contracts[*call.To]
	if !ok {
		c.t.Fatalf("a call to %s, which the chain does not hold", call.To.Hex())
	}
	method, err := k.abi.MethodById(call.Data)
	if err != nil {
		c.t.Fatal(err)
	}
	args, err := method.Inputs.Unpack(call.Data[4:])
	if err != nil {
		c.t.Fatal(err)
	}
	answer, ok := k.reads[method.Name]
	if !ok {
		c.t.Fatalf("an unexpected read of %s", method.Name)
	}
	out, err := answer(args)
	if err != nil {
		return nil, err
	}
	packed, err := method.Outputs.Pack(out...)
	if err != nil {
		c.t.Fatalf("%s: %v", method.Name, err)
	}
	return packed, nil
}

func (c *chain) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if number == nil {
		number = new(big.Int).SetUint64(c.head)
	}
	return &types.Header{Number: number, Time: c.time}, nil
}

func (c *chain) CodeAt(context.Context, common.Address, *big.Int) ([]byte, error) {
	return []byte{1}, nil
}

func (c *chain) PendingCodeAt(context.Context, common.Address) ([]byte, error) { return []byte{1}, nil }

func (c *chain) PendingNonceAt(context.Context, common.Address) (uint64, error) { return 0, nil }

func (c *chain) SuggestGasPrice(context.Context) (*big.Int, error) { return big.NewInt(1), nil }

func (c *chain) SuggestGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(0), nil }

func (c *chain) EstimateGas(context.Context, ethereum.CallMsg) (uint64, error) {
	return 0, errors.New("no transactions")
}

func (c *chain) SendTransaction(context.Context, *types.Transaction) error {
	return errors.New("no transactions")
}

func (c *chain) TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error) {
	return nil, false, ethereum.NotFound
}

func (c *chain) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	return nil, ethereum.NotFound
}

func (c *chain) FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
	return nil, errors.New("no logs")
}

func (c *chain) SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error) {
	return nil, errors.New("no logs")
}

var (
	addrBand       = common.HexToAddress("0xb1")
	addrEngine     = common.HexToAddress("0xe1")
	addrAccounts   = common.HexToAddress("0xa1")
	addrLiquidator = common.HexToAddress("0xc1")
	addrEthFeed    = common.HexToAddress("0xf1")
	addrUSDG       = common.HexToAddress("0xd1")
	addrWETH       = common.HexToAddress("0xd2")
	addrNVDAToken  = common.HexToAddress("0xd3")
	addrNVDAPool   = common.HexToAddress("0xe2")
	addrSPYToken   = common.HexToAddress("0xd4")
	addrSPYPool    = common.HexToAddress("0xe3")
	addrTSLAToken  = common.HexToAddress("0xd5")
	addrTSLAPool   = common.HexToAddress("0xe4")
	addrShorts     = common.HexToAddress("0xc2")
)

func symbol(name string) [32]byte {
	var out [32]byte
	copy(out[:], name)
	return out
}

// world is a chain of the fixture's three assets with their band, engine and one position, and the SDK client over it.
type world struct {
	*chain
	band, engine, accounts, liquidator, eth *contract
	poolOf                                  map[string]common.Address
	client                                  *sdk.Client
	m                                       *Market
	holdings                                map[string]*big.Int
	requirement                             func(quantities, prices []*big.Int) (*big.Int, uint8)
	ethAnswer, ethUpdated                   *big.Int
	session                                 sessionState
	auction                                 uint64
	short                                   *shortState
}

// sessionState is what the band's session answers: open, regular hours, the next state a trading day, and no boundary.
type sessionState struct {
	state, nyse, nyseNext uint8
	boundaryMs            uint64
}

// shortState is a short the chain holds: its borrow shares and its health.
type shortState struct{ equity, requirement *big.Int }

func newWorld(t *testing.T) *world {
	t.Helper()
	w := &world{chain: newChain(t), m: fixture(t), poolOf: map[string]common.Address{}, holdings: map[string]*big.Int{},
		ethAnswer: units(2000, 8), ethUpdated: new(big.Int).SetUint64(1_790_775_000 - 60)}
	w.band = w.add(addrBand, &bandbinding.BandMetaData)
	w.engine = w.add(addrEngine, &marginbinding.MarginMetaData)
	w.accounts = w.add(addrAccounts, &marginaccounts.MarginAccountsMetaData)
	w.liquidator = w.add(addrLiquidator, &liquidator.LiquidatorMetaData)
	w.eth = w.add(addrEthFeed, &aggregator.AggregatorMetaData)
	tokens := map[string]common.Address{"NVDA": addrNVDAToken, "TSLA": addrTSLAToken, "SPY": addrSPYToken}
	index := func(args []any) int {
		for i, a := range w.m.Assets {
			if symbol(a.Name) == args[0].([32]byte) {
				return i
			}
		}
		w.t.Fatalf("an unknown asset %x", args[0])
		return -1
	}
	w.band.on("sequencerSettled", func([]any) []any { return []any{w.m.Settled} })
	w.session = sessionState{state: 2, nyse: 1, nyseNext: 1}
	w.band.on("session", func([]any) []any {
		return []any{w.session.state, w.session.nyse, w.session.nyseNext, uint64(0), w.session.boundaryMs}
	})
	w.band.on("quote", func(args []any) []any {
		a := w.m.Assets[index(args)]
		return []any{a.State, uint8(2), a.Mid, uint64(30), a.Low, a.High}
	})
	w.band.on("corporateAction", func(args []any) []any {
		status := uint8(0)
		if w.m.Assets[index(args)].Unconfirmed {
			status = 2
		}
		return []any{status, uint64(0), big.NewInt(0), big.NewInt(0)}
	})
	w.band.on("asset", func(args []any) []any {
		a := w.m.Assets[index(args)]
		feed, indexFeed := [32]byte{1, byte(index(args))}, [32]byte{}
		if a.Basis != 0 {
			feed, indexFeed = [32]byte{}, [32]byte{2, byte(index(args))}
		}
		return []any{common.Address{}, feed, indexFeed, tokens[a.Name]}
	})
	w.band.on("variance", func(args []any) []any {
		id := args[0].([32]byte)
		return []any{w.m.Assets[id[1]].Variance}
	})
	w.engine.on("assets", func([]any) []any {
		var out [][32]byte
		for _, a := range w.m.Assets {
			out = append(out, symbol(a.Name))
		}
		return []any{out}
	})
	w.engine.on("market", func([]any) []any { return []any{symbol("SPY")} })
	w.engine.on("volatility", func(args []any) []any {
		return []any{w.m.Set.Parameters.Volatilities[index(args)], uint32(0)}
	})
	w.engine.on("weekendGap", func(args []any) []any { return []any{w.m.Set.Parameters.Gaps[index(args)], uint32(0)} })
	w.engine.on("depth", func(args []any) []any {
		i := index(args)
		return []any{w.m.Depths[2*i], w.m.Depths[2*i+1], uint32(0), uint32(0)}
	})
	w.engine.on("correlation", func(args []any) []any {
		i, j := index(args), 0
		for k, a := range w.m.Assets {
			if symbol(a.Name) == args[1].([32]byte) {
				j = k
			}
		}
		return []any{w.m.Set.Parameters.Correlations[margin.Pair(len(w.m.Assets), i, j)], uint16(0)}
	})
	w.engine.on("pool", func(args []any) []any { return []any{w.poolOf[w.m.Assets[index(args)].Name]} })
	w.engine.on("ethUsdFeed", func([]any) []any { return []any{addrEthFeed} })
	w.liquidator.on("ethUsd", func([]any) []any { return []any{addrEthFeed} })
	w.liquidator.on("recallHaircut", func([]any) []any { return []any{w.m.RecallHaircut} })
	w.eth.on("latestRoundData", func([]any) []any {
		return []any{big.NewInt(1), w.ethAnswer, w.ethUpdated, w.ethUpdated, big.NewInt(1)}
	})
	w.accounts.on("stocks", func([]any) []any {
		var symbols [][32]byte
		var addresses []common.Address
		for _, a := range w.m.Assets {
			symbols, addresses = append(symbols, symbol(a.Name)), append(addresses, tokens[a.Name])
		}
		return []any{symbols, addresses}
	})
	w.accounts.on("inBaskets", func([]any) []any {
		out := make([]*big.Int, len(w.m.Assets))
		for i, a := range w.m.Assets {
			out[i] = n(0)
			if v := w.holdings["basket "+a.Name]; v != nil {
				out[i] = v
			}
		}
		return []any{out}
	})
	held := func(key string) *big.Int {
		if v := w.holdings[key]; v != nil {
			return v
		}
		return n(0)
	}
	w.accounts.on("collateral", func(args []any) []any {
		token := args[2].(common.Address)
		for name, address := range map[string]common.Address{"NVDA": addrNVDAToken, "TSLA": addrTSLAToken, "SPY": addrSPYToken, "USDG": addrUSDG, "WETH": addrWETH} {
			if token == address {
				return []any{held("collateral " + name)}
			}
		}
		return []any{n(0)}
	})
	w.accounts.on("lent", func(args []any) []any {
		for name, address := range tokens {
			if args[2].(common.Address) == address {
				return []any{held("lent " + name)}
			}
		}
		return []any{n(0)}
	})
	w.accounts.on("debt", func([]any) []any { return []any{held("debt")} })
	w.accounts.on("premium", func([]any) []any { return []any{held("premium")} })
	deployments := &sdk.Deployments{ChainID: 412346,
		Tokens:    map[string]common.Address{"USDG": addrUSDG, "WETH": addrWETH},
		Tapehouse: map[string]common.Address{"Band": addrBand, "Margin": addrEngine, "MarginAccounts": addrAccounts, "Liquidator": addrLiquidator}}
	w.client = sdk.NewClient(w.chain, deployments)
	w.engine.on("currentRequirement", func(args []any) []any {
		quantities, prices := args[0].([]*big.Int), args[1].([]*big.Int)
		requirement, regime := w.requirement(quantities, prices)
		return []any{requirement, uint8(0), regime}
	})
	w.requirement = func(quantities, prices []*big.Int) (*big.Int, uint8) {
		regime := margin.RegimeAt(w.session.state, w.session.boundaryMs, w.time*1000)
		requirement, _ := w.m.Requirement(quantities, prices, regime)
		return requirement, regime.Code()
	}
	w.liquidator.on("auctions", func([]any) []any { return []any{w.auction, w.closedNow()} })
	w.liquidator.on("OPEN_AUCTION_LIFETIME", func([]any) []any { return []any{n(3600)} })
	w.liquidator.on("CLOSED_AUCTION_LIFETIME", func([]any) []any { return []any{n(86_400)} })
	w.liquidator.on("shortfall", func(args []any) []any {
		regime := margin.RegimeAt(w.session.state, w.session.boundaryMs, w.time*1000)
		p := w.position(args[0].(common.Address))
		v := w.m.Assess(p, w.m.Current(), w.closedNow(), regime)
		return []any{v.Equity, v.Requirement, v.Short(), w.closedNow()}
	})
	w.accounts.on("liquidationPrice", func([]any) []any { return []any{units(77, 8)} })
	shorts := w.add(addrShorts, &shortbinding.ShortPositionsMetaData)
	shorts.on("position", func([]any) []any {
		if w.short == nil {
			return []any{n(0), n(0), n(0), n(0)}
		}
		return []any{n(0), n(0), n(1), n(0)}
	})
	shorts.on("health", func([]any) []any { return []any{w.short.equity, w.short.requirement, uint8(0), uint8(2)} })
	deployments.Tapehouse["ShortPositions"] = addrShorts
	return w
}

func (w *world) closedNow() bool {
	return closedAt(bandbinding.SessionOutput{State: w.session.state, Nyse: w.session.nyse, NyseNext: w.session.nyseNext,
		BoundaryMs: w.session.boundaryMs}, w.time*1000)
}

// position is the position the double holds, as the reader would read it.
func (w *world) position(account common.Address) *Position {
	p := &Position{Account: account, USDG: w.holding("collateral USDG"), WETH: w.holding("collateral WETH"),
		Debt: w.holding("debt"), Premium: w.holding("premium")}
	for _, a := range w.m.Assets {
		quantity := new(big.Int).Add(w.holding("collateral "+a.Name), w.holding("lent "+a.Name))
		quantity.Add(quantity, w.holding("basket "+a.Name))
		if quantity.Sign() > 0 {
			p.Stocks = append(p.Stocks, Stock{Asset: a.Name, Quantity: quantity, Lent: w.holding("lent " + a.Name)})
		}
	}
	return p
}

func (w *world) holding(key string) *big.Int {
	if v := w.holdings[key]; v != nil {
		return v
	}
	return n(0)
}

func (w *world) pool(asset string, address common.Address, stock, quote common.Address, fee int64, observe func() ([]any, error)) {
	p := w.add(address, &uniswapv3pool.UniswapV3PoolMetaData)
	w.poolOf[asset] = address
	order := func(a, b common.Address) (common.Address, common.Address) {
		if a.Cmp(b) < 0 {
			return a, b
		}
		return b, a
	}
	token0, token1 := order(stock, quote)
	p.on("token0", func([]any) []any { return []any{token0} })
	p.on("token1", func([]any) []any { return []any{token1} })
	p.on("fee", func([]any) []any { return []any{big.NewInt(fee)} })
	p.reads["observe"] = func([]any) ([]any, error) { return observe() }
}
