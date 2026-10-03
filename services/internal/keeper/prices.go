// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// StatusFeeds are RedStone's New York market status, which the band writes together or not at all and reads its
// session from.
var StatusFeeds = []string{"NY_MARKET_CURRENT_STATUS", "NY_MARKET_NEXT_STATUS", "NY_MARKET_NEXT_CHANGE_TIME"}

// StatusPackage is RedStone's data package that carries the three status feeds at once.
const StatusPackage = "NY_MARKET_STATUS"

const (
	// PackageInterval is how often RedStone signs a package: the next one a write that lost can take.
	PackageInterval = 10 * time.Second
	// PackageRetries is how many later packages a write tries after another writer took the one it sent.
	PackageRetries = 3
	// StatusAlertAge is how old the band's market status may grow before the keeper alerts: well before the hour
	// after which every band degrades and every session is unknown.
	StatusAlertAge = 15 * time.Minute
)

// statusGroup is a package source that reads the three status feeds as their one package, five signatures fewer to
// verify than three packages' fifteen.
type statusGroup struct{ sdk.PackageSource }

func (s statusGroup) Payload(ctx context.Context, feedIDs []string) ([]byte, error) {
	ids := slices.Clone(feedIDs)
	if i := slices.Index(ids, StatusFeeds[0]); i >= 0 && slices.Equal(ids[i:i+len(StatusFeeds)], StatusFeeds) {
		ids = slices.Replace(ids, i, i+len(StatusFeeds), StatusPackage)
	}
	return s.PackageSource.Payload(ctx, ids)
}

// priceFeeds are the feeds the band tracks, in the order a write lists them: the status, then the indexes, which
// anchor only while the stored status says NYSE is in regular hours, then the 24/7 prices.
func (k *Keeper) priceFeeds(ctx context.Context) ([]string, error) {
	assets, err := k.bandAssets(ctx)
	if err != nil {
		return nil, err
	}
	feeds := slices.Clone(StatusFeeds)
	for _, a := range assets {
		if a.indexID != "" && !slices.Contains(feeds, a.indexID) {
			feeds = append(feeds, a.indexID)
		}
	}
	for _, a := range assets {
		if a.feedID != "" && !slices.Contains(feeds, a.feedID) {
			feeds = append(feeds, a.feedID)
		}
	}
	return feeds, nil
}

// writePrices writes every feed the band tracks in one writePrices, from RedStone's latest packages. A write another
// writer beat with the same package, PackageNotNewer, takes the next package. It alerts once the band's status is
// StatusAlertAge old.
func (k *Keeper) writePrices(ctx context.Context) error {
	feeds, err := k.priceFeeds(ctx)
	if err != nil {
		return err
	}
	for attempt := 0; ; attempt++ {
		tx, err := k.client.Band().WritePrices(ctx, k.prices, feeds)
		if err == nil {
			_, err = k.sender.Act(ctx, "writePrices", tx)
		}
		if err == nil || !reverted(err, "PackageNotNewer") || attempt == PackageRetries {
			k.checkStatus(ctx)
			return err
		}
		k.log.Info("another writer sent this package first; taking the next", "error", err.Error())
		if err := k.sleep(ctx, PackageInterval); err != nil {
			return err
		}
	}
}

func (k *Keeper) checkStatus(ctx context.Context) {
	status, err := k.client.Band().Price(k.opts(ctx), StatusFeeds[0])
	if err != nil {
		k.log.Error("the band's market status cannot be read", "error", err)
		return
	}
	age := k.now().Sub(time.UnixMilli(int64(status.PackageTimestampMs)))
	if age >= StatusAlertAge {
		k.log.Error("the band's market status is old: every band degrades once it is an hour old",
			"age", age.Round(time.Second).String())
	}
}

// syncMultipliers calls syncMultiplier for every asset with a Stock Token: once after the keeper starts, after the
// prices, since a band deployed after a step confirms nothing until then; whenever the token's terms change, which is
// how UIMultiplierUpdated shows; and, while a material step waits for Chainlink, after a fresh 24/7 price with each
// new Chainlink round and when the issuer's oracle pause ends.
func (k *Keeper) syncMultipliers(ctx context.Context) error {
	assets, err := k.bandAssets(ctx)
	if err != nil {
		return err
	}
	for _, a := range assets {
		if a.token == (common.Address{}) {
			continue
		}
		if err := k.syncMultiplier(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) syncMultiplier(ctx context.Context, a asset) error {
	band := k.client.Band()
	terms, err := band.Terms(k.opts(ctx), a.name)
	if err != nil {
		return err
	}
	action, err := band.CorporateAction(k.opts(ctx), a.name)
	if err != nil {
		return err
	}
	k.mu.Lock()
	first := !k.started[a.name]
	k.mu.Unlock()
	changed := !k.remember("terms "+a.name, fmt.Sprint(terms.NewMultiplier, terms.EffectiveAt))
	waiting := false
	if action.Status == 2 && !terms.OraclePaused {
		named, err := band.Asset(k.opts(ctx), a.name)
		if err != nil {
			return err
		}
		round, err := k.client.LatestRound(k.opts(ctx), named.ChainlinkFeed)
		if err != nil {
			return err
		}
		waiting = !k.remember("round "+a.name, round.RoundID.String())
	}
	if terms.OraclePaused {
		k.remember("round "+a.name, "paused")
	}
	if !first && !changed && !waiting {
		return nil
	}
	if waiting {
		if err := k.writePrices(ctx); err != nil {
			k.log.Warn("no fresh 24/7 price before the multiplier sync", "asset", a.name, "error", err)
		}
	}
	tx, err := band.SyncMultiplier(a.name)
	if _, err := k.act(ctx, "syncMultiplier "+a.name, tx, err); err != nil {
		k.remember("terms "+a.name, "")
		return err
	}
	k.mu.Lock()
	k.started[a.name] = true
	k.mu.Unlock()
	return nil
}

// HaltLifetime is how long a signed halt holds, and HaltRefresh how often the keeper signs a new one while the asset
// stays halted: a dead keeper's last halt lapses within the lifetime, to a degraded band.
const (
	HaltLifetime = 15 * time.Minute
	HaltRefresh  = 5 * time.Minute
)

// halts signs the band's trading halts with the halt signer's key, from Robinhood's quotes of its Stock Tokens.
type halts struct {
	key    *ecdsa.PrivateKey
	url    string
	client *http.Client
}

// halted reads which Stock Tokens Robinhood shows halted, by symbol.
func (h *halts) halted(ctx context.Context) (map[string]bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return nil, err
	}
	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the halts answered %s", response.Status)
	}
	var quotes struct {
		Quotes []struct {
			Symbol string `json:"tokenSymbol"`
			Halted bool   `json:"isTradingHalt"`
		} `json:"quotes"`
	}
	if err := json.NewDecoder(response.Body).Decode(&quotes); err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(quotes.Quotes))
	for _, quote := range quotes.Quotes {
		out[quote.Symbol] = quote.Halted
	}
	return out, nil
}

// HaltDigest is the EIP-712 digest of HaltState(bytes32 symbol,bool halted,uint64 issuedAt,uint64 expiresAt) under
// the domain "Tapehouse Band", version "1", of the band on chainID, that the halt signer signs.
func HaltDigest(chainID uint64, band common.Address, asset string, halted bool, issuedAt, expiresAt uint64) ([]byte, error) {
	symbol, err := sdk.ToBytes32(asset)
	if err != nil {
		return nil, err
	}
	digest, _, err := apitypes.TypedDataAndHash(apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {{Name: "name", Type: "string"}, {Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"}, {Name: "verifyingContract", Type: "address"}},
			"HaltState": {{Name: "symbol", Type: "bytes32"}, {Name: "halted", Type: "bool"},
				{Name: "issuedAt", Type: "uint64"}, {Name: "expiresAt", Type: "uint64"}},
		},
		PrimaryType: "HaltState",
		Domain: apitypes.TypedDataDomain{Name: "Tapehouse Band", Version: "1",
			ChainId: math.NewHexOrDecimal256(int64(chainID)), VerifyingContract: band.Hex()},
		Message: apitypes.TypedDataMessage{"symbol": hexutil.Encode(symbol[:]), "halted": halted,
			"issuedAt": fmt.Sprint(issuedAt), "expiresAt": fmt.Sprint(expiresAt)},
	})
	return digest, err
}

// SignHalt signs a HaltState with key as the band's halt signer: 65 bytes, r, s and v of 27 or 28.
func SignHalt(key *ecdsa.PrivateKey, chainID uint64, band common.Address, asset string, halted bool, issuedAt, expiresAt uint64) ([]byte, error) {
	digest, err := HaltDigest(chainID, band, asset, halted, issuedAt, expiresAt)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(digest, key)
	if err != nil {
		return nil, err
	}
	signature[64] += 27
	return signature, nil
}

// writeHalts writes a halt of every asset Robinhood shows halted, issued at the latest block's time and holding for
// HaltLifetime, renewed every HaltRefresh, and its lift once the halt ends, which also ends the degraded band a lapsed
// halt leaves. A HaltNotNewer at the message's own issue
// time is the same message relayed first: written.
func (k *Keeper) writeHalts(ctx context.Context) error {
	if k.halts == nil {
		return nil
	}
	halted, err := k.halts.halted(ctx)
	if err != nil {
		return err
	}
	assets, err := k.bandAssets(ctx)
	if err != nil {
		return err
	}
	signer, err := k.client.Band().HaltSigner(k.opts(ctx))
	if err != nil {
		return err
	}
	if signer != crypto.PubkeyToAddress(k.halts.key.PublicKey) {
		return fmt.Errorf("the halt key is not the band's halt signer %s", signer.Hex())
	}
	head, err := k.chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	now := head.Time
	band := k.client.Band()
	for _, a := range assets {
		isHalted, quoted := halted[a.name]
		if !quoted {
			k.log.Warn("Robinhood quotes no such Stock Token; its halt is left as it is", "asset", a.name)
			continue
		}
		stored, err := band.Halt(&bind.CallOpts{Context: ctx, BlockNumber: head.Number}, a.name)
		if err != nil {
			return err
		}
		holding := stored.SignedHalt && stored.Until > now
		renew := isHalted && (!holding || stored.Until-now <= uint64((HaltLifetime-HaltRefresh)/time.Second))
		lift := !isHalted && stored.Until != 0
		if (!renew && !lift) || now <= stored.IssuedAt {
			continue
		}
		expires := now + uint64(HaltLifetime/time.Second)
		if err := k.writeHalt(ctx, a.name, isHalted, now, expires); err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) writeHalt(ctx context.Context, asset string, halted bool, issuedAt, expiresAt uint64) error {
	signature, err := SignHalt(k.halts.key, k.client.Deployments().ChainID, k.client.Deployments().Tapehouse["Band"],
		asset, halted, issuedAt, expiresAt)
	if err != nil {
		return err
	}
	tx, err := k.client.Band().WriteHalt(asset, halted, issuedAt, expiresAt, signature)
	if err != nil {
		return err
	}
	_, err = k.sender.Act(ctx, fmt.Sprintf("writeHalt %s %t", asset, halted), tx)
	if revert, ok := sdk.DecodeRevert(err); ok && revert.Name == "HaltNotNewer" && len(revert.Args) == 3 &&
		revert.Args[1] == issuedAt {
		return nil
	}
	return err
}
