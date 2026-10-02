// SPDX-License-Identifier: MIT OR Apache-2.0

// Command devnode reads SPY's band and its BandFeed on a dev node, sends RedStone packages from a file through a
// PackageSource and decodes the band's refusal of their age, then authorizes a fresh address in the margin accounts,
// which sells 1 SPY short for the account at a QuoterV2 quote and buys it back. Then it mints 1 share of the basket
// PAIR for the fresh address, which deposits it into its cross position, reads it there as NVDA and SPY, unwraps it and
// withdraws them.
//
// Usage: PRIVATE_KEY=0x… go run ./sdk/examples/devnode RPC_URL DEPLOYMENTS_JSON PAYLOAD_FILE
package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tapehouse/tapehouse/services/sdk"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/basket"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/marginaccounts"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/shortpositions"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/stocktoken"
	"github.com/tapehouse/tapehouse/services/sdk/bindings/usdg"
)

type fileSource string

func (f fileSource) Payload(context.Context, []string) ([]byte, error) {
	text, err := os.ReadFile(string(f))
	if err != nil {
		return nil, err
	}
	return hexutil.Decode(strings.TrimSpace(string(text)))
}

func check(ok bool, format string, args ...any) {
	if !ok {
		log.Fatalf("FAIL: "+format, args...)
	}
}

func must[T any](value T, err error) T {
	if err != nil {
		log.Fatalf("FAIL: %v", err)
	}
	return value
}

func refusal(err error) *sdk.Revert {
	var revert *sdk.Revert
	check(errors.As(err, &revert), "the call did not revert with a decodable error: %v", err)
	return revert
}

func main() {
	if len(os.Args) != 4 {
		log.Fatal("usage: PRIVATE_KEY=0x… devnode RPC_URL DEPLOYMENTS_JSON PAYLOAD_FILE")
	}
	secret, ok := os.LookupEnv("PRIVATE_KEY")
	check(ok, "set PRIVATE_KEY to the account owner's key")
	ctx := context.Background()
	backend := must(ethclient.DialContext(ctx, os.Args[1]))
	key := must(crypto.HexToECDSA(strings.TrimPrefix(secret, "0x")))
	deployments := must(sdk.LoadDeployments(os.Args[2]))
	chainID := must(backend.ChainID(ctx))
	check(chainID.Uint64() == deployments.ChainID, "the registry is for chain %d, not %v", deployments.ChainID, chainID)
	client := sdk.NewClient(backend, deployments)
	owner := bind.NewKeyedTransactor(key, chainID)
	account := owner.From
	operatorKey := must(crypto.GenerateKey())
	operator := bind.NewKeyedTransactor(operatorKey, chainID)
	send := func(opts *bind.TransactOpts, tx sdk.Tx, err error) *types.Receipt {
		check(err == nil, "%v", err)
		receipt := must(bind.WaitMined(ctx, backend, must(client.Send(opts, tx)).Hash()))
		check(receipt.Status == types.ReceiptStatusSuccessful, "the transaction to %v reverted", tx.To)
		return receipt
	}

	header := must(backend.HeaderByNumber(ctx, nil))
	at := &bind.CallOpts{Context: ctx, BlockNumber: header.Number}
	band := client.Band()
	quote := must(band.Quote(at, "SPY"))
	session := must(band.Session(at))
	halt := must(band.Halt(at, "SPY"))
	check(quote.State != 0 && !halt.SignedHalt, "SPY's band is halted")
	round := must(band.LatestRound(at, "SPY"))
	check(round.Answer.Uint64() == quote.Low, "the feed answers %v, the band's low edge is %d", round.Answer, quote.Low)
	check(round.RoundID.Uint64() == header.Time && round.UpdatedAt.Uint64() == header.Time, "the round is not the block time")
	whole := must(band.LatestBand(at, "SPY"))
	check(whole.State == quote.State && whole.Mid == quote.Mid, "latestBand is not the band")
	sealed := must(band.Sealed(at, "SPY", session.BoundaryMs))
	sealing := "within its window"
	if err := client.Simulate(&bind.CallOpts{Context: ctx, From: account}, must(band.Seal("SPY"))); err != nil {
		revert := refusal(err)
		check(revert.Name == "NotSealWindow" && revert.Args[0] == session.State && revert.Args[1] == session.BoundaryMs,
			"seal() reverted with %v", revert)
		sealing = revert.Error()
	}

	write := must(band.WritePrices(ctx, fileSource(os.Args[3]), []string{"NVDA---24_7"}))
	stale := refusal(client.Simulate(&bind.CallOpts{Context: ctx, From: account}, write))
	check(strings.HasPrefix(stale.Error(), "TimestampIsTooOld(1790119750, "), "the band took the file's packages: %v", stale)

	fund(ctx, backend, key, chainID, operator.From, big.NewInt(1e16))
	accounts, shorts := client.Accounts(), client.Shorts()
	margin := big.NewInt(1_000_000_000)
	refused := refusal(client.Simulate(&bind.CallOpts{Context: ctx, From: operator.From}, must(shorts.Deposit("SPY", margin, account))))
	check(refused.Error() == fmt.Sprintf("Unauthorized(%v, %v)", operator.From, account), "the shorts did not refuse: %v", refused)
	usdgAddress := deployments.Tokens["USDG"]
	withdrawal := must(accounts.Withdraw(sdk.Cross, usdgAddress, big.NewInt(1), account, account))
	check(refusal(client.Simulate(&bind.CallOpts{Context: ctx, From: operator.From}, withdrawal)).Error() == refused.Error(),
		"the margin accounts did not refuse the operator")

	authorized := send(owner, must(accounts.SetAuthorization(operator.From, true)), nil)
	binding := marginaccounts.NewMarginAccounts()
	set := must(binding.UnpackAuthorizationSetEvent(authorized.Logs[0]))
	check(set.Authorized == operator.From && set.Allowed, "no AuthorizationSet")
	check(must(accounts.IsAuthorized(&bind.CallOpts{Context: ctx}, account, operator.From)), "the operator is not authorized")

	shortsAddress := deployments.Tapehouse["ShortPositions"]
	send(owner, sdk.Tx{To: usdgAddress, Data: must(usdg.NewUsdg().TryPackApprove(shortsAddress, margin))}, nil)
	send(owner, must(shorts.Deposit("SPY", margin, account)), nil)
	one := big.NewInt(1e18)
	latest := &bind.CallOpts{Context: ctx}
	sale := must(shorts.QuoteSale(latest, "SPY", one, 50))
	shortsBinding := shortpositions.NewShortPositions()
	sold := send(operator, must(shorts.Sell("SPY", one, sale.MinProceeds, account)), nil)
	sell := event(sold, shortsBinding.UnpackSellEvent)
	check(sell.Proceeds.Cmp(sale.Proceeds) == 0, "the sale paid %v, QuoterV2 quoted %v", sell.Proceeds, sale.Proceeds)
	open := must(shorts.Position(latest, account, "SPY"))
	check(open.Debt.Cmp(one) >= 0, "the short owes %v, not 1 SPY and its fee", open.Debt)
	health := must(shorts.Health(latest, account, "SPY"))
	check(health.Equity.Cmp(health.Requirement) > 0, "the short falls short")
	buyBack := must(shorts.QuoteCover(latest, "SPY", open.Debt, 50))
	covered := send(operator, must(shorts.Cover("SPY", abi.MaxUint256, buyBack.MaxCost, account)), nil)
	cover := event(covered, shortsBinding.UnpackCoverEvent)
	check(cover.Cost.Cmp(buyBack.Cost) >= 0 && cover.Cost.Cmp(buyBack.MaxCost) <= 0,
		"the buy-back cost %v, QuoterV2 quoted %v, at most %v", cover.Cost, buyBack.Cost, buyBack.MaxCost)
	closed := must(shorts.Position(latest, account, "SPY"))
	check(closed.Shares.Sign() == 0 && closed.Debt.Sign() == 0, "the short is still open")
	send(operator, must(shorts.Withdraw("SPY", closed.UsdgHeld, account, account)), nil)
	send(owner, must(accounts.SetAuthorization(operator.From, false)), nil)
	check(!must(accounts.IsAuthorized(latest, account, operator.From)), "the operator is still authorized")

	pair, holder := client.Basket("PAIR"), operator.From
	basketAddress := deployments.Baskets["PAIR"]
	components := must(pair.Components(latest))
	check(strings.Join(components.Assets, ",") == "NVDA,SPY", "PAIR holds %v", components.Assets)
	paid := must(pair.PreviewMint(latest, one))
	for i, token := range components.Tokens {
		send(owner, sdk.Tx{To: token, Data: must(stocktoken.NewStockToken().TryPackApprove(basketAddress, paid[i]))}, nil)
	}
	send(owner, must(pair.Mint(one, holder, paid)), nil)
	send(operator, sdk.Tx{To: basketAddress, Data: must(basket.NewBasket().TryPackApprove(deployments.Tapehouse["MarginAccounts"], one))}, nil)
	send(operator, must(accounts.Deposit(sdk.Cross, basketAddress, one, holder)), nil)
	through := must(accounts.InBaskets(latest, holder, sdk.Cross))
	margined := must(accounts.Health(latest, holder, sdk.Cross))
	check(through["NVDA"].Cmp(paid[0]) == 0 && through["SPY"].Cmp(paid[1]) == 0, "the share does not hold what it was minted for")
	shares := must(accounts.Collateral(latest, holder, sdk.Cross, basketAddress))
	send(operator, must(accounts.Unwrap(holder, "PAIR", one)), nil)
	unwrapped := make([]*big.Int, len(components.Tokens))
	for i, token := range components.Tokens {
		unwrapped[i] = must(accounts.Collateral(latest, holder, sdk.Cross, token))
	}
	check(unwrapped[0].Cmp(through["NVDA"]) == 0 && unwrapped[1].Cmp(through["SPY"]) == 0,
		"the unwrap did not put the share's tokens in the position")
	for i, token := range components.Tokens {
		send(operator, must(accounts.Withdraw(sdk.Cross, token, unwrapped[i], holder, holder)), nil)
	}

	fmt.Printf("go sdk: SPY band state %d at %d, feed round %v answers %v, seal at %d %s (sealed %d); writePrices %v; "+
		"before authorization %v; sold 1 SPY for %v and bought it back for %v\n",
		quote.State, quote.Mid, round.RoundID, round.Answer, session.BoundaryMs, sealing, sealed.SealedAt, stale,
		refused, sell.Proceeds, cover.Cost)
	fmt.Printf("basket PAIR: %v shares in the cross position hold %v NVDA and %v SPY; equity %v requirement %v\n",
		shares, through["NVDA"], through["SPY"], margined.Equity, margined.Requirement)
	fmt.Printf("basket PAIR: unwrapped into %v NVDA and %v SPY in the cross position, then withdrawn\n", unwrapped[0], unwrapped[1])
	fmt.Println("PASS")
}

func fund(ctx context.Context, backend *ethclient.Client, key *ecdsa.PrivateKey, chainID *big.Int, to common.Address, value *big.Int) {
	from := crypto.PubkeyToAddress(key.PublicKey)
	gas := must(backend.EstimateGas(ctx, ethereum.CallMsg{From: from, To: &to, Value: value}))
	tx := must(types.SignNewTx(key, types.LatestSignerForChainID(chainID), &types.LegacyTx{
		Nonce:    must(backend.PendingNonceAt(ctx, from)),
		GasPrice: must(backend.SuggestGasPrice(ctx)),
		Gas:      gas,
		To:       &to,
		Value:    value,
	}))
	check(backend.SendTransaction(ctx, tx) == nil, "the funding of %v was not sent", to)
	check(must(bind.WaitMined(ctx, backend, tx.Hash())).Status == types.ReceiptStatusSuccessful, "the funding of %v failed", to)
}

func event[T any](receipt *types.Receipt, unpack func(*types.Log) (*T, error)) *T {
	for _, l := range receipt.Logs {
		if value, err := unpack(l); err == nil {
			return value
		}
	}
	log.Fatalf("FAIL: the receipt of %v has no such event", receipt.TxHash)
	return nil
}
