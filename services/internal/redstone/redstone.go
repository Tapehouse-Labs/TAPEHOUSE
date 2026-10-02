// SPDX-License-Identifier: MIT OR Apache-2.0

// Package redstone relays RedStone's latest signed data packages from its gateways, the keyed ones first, as the
// payload the band's writePrices verifies. It serves a package only once its signature recovers to one of the signers
// the band accepts, so it adds nothing a signed package cannot reproduce.
package redstone

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/internal/codec"
)

// Signers are the five nodes of RedStone's redstone-primary-prod data service, the only signers the band accepts.
var Signers = []common.Address{
	common.HexToAddress("0x8BB8F32Df04c8b654987DAaeD53D6B6091e3B774"),
	common.HexToAddress("0xdEB22f54738d54976C4c0fe5ce6d408E40d88499"),
	common.HexToAddress("0x51Ce04Be4b3E32572C4Ec9135221d0691Ba7d202"),
	common.HexToAddress("0xDD682daEC5A90dD295d14DA4b0bec9281017b5bE"),
	common.HexToAddress("0x9c5AE89C4Af6aA32cE58588DBaF90d18a855B6de"),
}

// UniqueSigners is how many of Signers must sign each data package, as the band requires.
const UniqueSigners = 3

// MaxAge is how long the relay keeps a gateway's answer before it reads the gateway again: RedStone signs a package
// every 10 seconds, and its keyed gateways allow one request a second.
const MaxAge = 2 * time.Second

const (
	authenticatedPath = "/v2/data-packages/latest-by-data-feeds/redstone-primary-prod"
	publicPath        = "/data-packages/latest/redstone-primary-prod"
	decimals          = 8
)

var marker = []byte{0x00, 0x00, 0x02, 0xed, 0x57, 0x01, 0x1e, 0x00, 0x00}

// Gateway is a RedStone gateway: a keyed one where Key is set, sent as x-api-key, or a public one.
type Gateway struct {
	URL string
	Key string
}

// Gateways are RedStone's keyed gateways, each with its key from REDSTONE_API_KEY or REDSTONE_BACKUP_API_KEY where
// getenv returns one, then its keyless public gateways.
func Gateways(getenv func(string) string) []Gateway {
	var out []Gateway
	for _, keyed := range []struct{ url, variable string }{
		{"https://oracle-gateway.a.redstone.finance", "REDSTONE_API_KEY"},
		{"https://oracle-gateway.gateway.redstone.vip", "REDSTONE_BACKUP_API_KEY"},
	} {
		if key := getenv(keyed.variable); key != "" {
			out = append(out, Gateway{keyed.url, key})
		}
	}
	return append(out, Gateway{URL: "https://oracle-gateway-1.a.redstone.finance"},
		Gateway{URL: "https://oracle-gateway-2.a.redstone.finance"})
}

// DataPoint is a signed value: its feed ID and the value times 10^8, as the band stores it.
type DataPoint struct {
	DataFeedID string `json:"dataFeedId"`
	Value      string `json:"value"`
}

// Package is a data package whose signature recovered to Signer.
type Package struct {
	DataPackageID string        `json:"dataPackageId"`
	Signer        codec.Address `json:"signer"`
	TimestampMs   uint64        `json:"timestampMs"`
	DataPoints    []DataPoint   `json:"dataPoints"`
	signed        []byte
}

// Signed is the payload writePrices takes for a list of data packages: every package of each, all signed at one
// time, and the packages it carries.
type Signed struct {
	TimestampMs uint64        `json:"timestampMs"`
	Payload     hexutil.Bytes `json:"payload"`
	Packages    []Package     `json:"packages"`
}

// Relay reads the gateways in order and serves the first answer that carries every data package asked for.
type Relay struct {
	gateways []Gateway
	client   *http.Client
	now      func() time.Time
	mu       sync.Mutex
	cache    map[string]cached
}

type cached struct {
	at     time.Time
	latest map[string][]gatewayPackage
}

// New returns a relay of gateways.
func New(gateways []Gateway, client *http.Client) *Relay {
	return &Relay{gateways: gateways, client: client, now: time.Now, cache: map[string]cached{}}
}

// Payload returns the payload of ids, as the SDKs' PackageSource does: the relay is a source an integrator plugs in.
func (r *Relay) Payload(ctx context.Context, ids []string) ([]byte, error) {
	signed, err := r.Signed(ctx, ids)
	if err != nil {
		return nil, err
	}
	return signed.Payload, nil
}

// Signed returns the latest signed packages of ids, each a single feed such as NVDA---24_7 or a group such as
// NY_MARKET_STATUS, and their payload. Each must carry UniqueSigners signatures that recover to Signers, all at one
// timestamp.
func (r *Relay) Signed(ctx context.Context, ids []string) (*Signed, error) {
	if len(ids) == 0 || slices.Contains(ids, "") {
		return nil, errors.New("no data package named")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var failures []string
	for _, gateway := range r.gateways {
		latest, err := r.latest(ctx, gateway, ids)
		if err == nil {
			var signed *Signed
			if signed, err = assemble(latest, ids); err == nil {
				return signed, nil
			}
		}
		failures = append(failures, fmt.Sprintf("%s: %v", gateway.URL, err))
	}
	return nil, fmt.Errorf("no RedStone gateway answered: %s", strings.Join(failures, "; "))
}

// latest returns a gateway's latest packages, read at most once in MaxAge: a public gateway's every package, a keyed
// one's those of ids.
func (r *Relay) latest(ctx context.Context, gateway Gateway, ids []string) (map[string][]gatewayPackage, error) {
	key := gateway.URL
	if gateway.Key != "" {
		key += "?" + strings.Join(ids, ",")
	}
	now := r.now()
	for k, hit := range r.cache {
		if now.Sub(hit.at) >= MaxAge {
			delete(r.cache, k)
		}
	}
	if hit, ok := r.cache[key]; ok {
		return hit.latest, nil
	}
	latest, err := r.fetch(ctx, gateway, ids)
	if err != nil {
		return nil, err
	}
	r.cache[key] = cached{now, latest}
	return latest, nil
}

type gatewayPackage struct {
	TimestampMilliseconds uint64         `json:"timestampMilliseconds"`
	Signature             string         `json:"signature"`
	DataPoints            []gatewayPoint `json:"dataPoints"`
}

type gatewayPoint struct {
	DataFeedID string      `json:"dataFeedId"`
	Value      json.Number `json:"value"`
	Decimals   *int        `json:"decimals"`
}

func (r *Relay) fetch(ctx context.Context, gateway Gateway, ids []string) (map[string][]gatewayPackage, error) {
	target := gateway.URL + publicPath
	if gateway.Key != "" {
		target = gateway.URL + authenticatedPath + "?" + url.Values{"dataFeedIds": {strings.Join(ids, ",")}}.Encode()
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	if gateway.Key != "" {
		request.Header.Set("x-api-key", gateway.Key)
	}
	response, err := r.client.Do(request)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", response.StatusCode)
	}
	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	var latest map[string][]gatewayPackage
	return latest, decoder.Decode(&latest)
}

// assemble verifies the packages of ids in latest and serialises them as RedStone's EVM connector appends them.
func assemble(latest map[string][]gatewayPackage, ids []string) (*Signed, error) {
	out := &Signed{}
	var payload bytes.Buffer
	for _, id := range ids {
		signers := map[codec.Address]bool{}
		for _, raw := range latest[id] {
			p, err := verify(id, raw)
			if err != nil || signers[p.Signer] {
				continue
			}
			signers[p.Signer] = true
			if out.TimestampMs == 0 {
				out.TimestampMs = p.TimestampMs
			}
			if p.TimestampMs != out.TimestampMs {
				return nil, fmt.Errorf("the latest packages carry different timestamps, %d and %d", out.TimestampMs,
					p.TimestampMs)
			}
			out.Packages = append(out.Packages, *p)
			payload.Write(p.signed)
		}
		if len(signers) < UniqueSigners {
			return nil, fmt.Errorf("%s has %d of the %d signatures the band requires", id, len(signers), UniqueSigners)
		}
	}
	payload.Write(binary.BigEndian.AppendUint16(nil, uint16(len(out.Packages))))
	payload.Write([]byte{0, 0, 0})
	payload.Write(marker)
	out.Payload = payload.Bytes()
	return out, nil
}

// verify serialises a package and recovers its signer, refusing one outside Signers.
func verify(id string, raw gatewayPackage) (*Package, error) {
	p := &Package{DataPackageID: id, TimestampMs: raw.TimestampMilliseconds}
	points := slices.Clone(raw.DataPoints)
	slices.SortFunc(points, func(a, b gatewayPoint) int {
		return bytes.Compare(feedID(a.DataFeedID), feedID(b.DataFeedID))
	})
	var signed bytes.Buffer
	for _, point := range points {
		places := decimals
		if point.Decimals != nil {
			places = *point.Decimals
		}
		value, err := scaled(string(point.Value), places)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", point.DataFeedID, err)
		}
		if len(point.DataFeedID) > 31 {
			return nil, fmt.Errorf("%s is longer than 31 bytes", point.DataFeedID)
		}
		signed.Write(feedID(point.DataFeedID))
		signed.Write(common.LeftPadBytes(value.Bytes(), 32))
		p.DataPoints = append(p.DataPoints, DataPoint{point.DataFeedID, value.String()})
	}
	signed.Write(binary.BigEndian.AppendUint64(nil, raw.TimestampMilliseconds)[2:])
	signed.Write(binary.BigEndian.AppendUint32(nil, 32))
	signed.Write(binary.BigEndian.AppendUint32(nil, uint32(len(points)))[1:])
	signature, err := base64.StdEncoding.DecodeString(raw.Signature)
	if err != nil || len(signature) != 65 || (signature[64] != 27 && signature[64] != 28) {
		return nil, errors.New("the signature is not 65 bytes with v of 27 or 28")
	}
	r, s := new(big.Int).SetBytes(signature[:32]), new(big.Int).SetBytes(signature[32:64])
	if !crypto.ValidateSignatureValues(signature[64]-27, r, s, true) {
		return nil, errors.New("the signature is malleable")
	}
	recoverable := append(slices.Clone(signature[:64]), signature[64]-27)
	key, err := crypto.SigToPub(crypto.Keccak256(signed.Bytes()), recoverable)
	if err != nil {
		return nil, err
	}
	signer := crypto.PubkeyToAddress(*key)
	if !slices.Contains(Signers, signer) {
		return nil, fmt.Errorf("%s is not one of redstone-primary-prod's signers", signer)
	}
	p.Signer = codec.Address(signer)
	p.signed = append(signed.Bytes(), signature...)
	return p, nil
}

func feedID(id string) []byte {
	return common.RightPadBytes([]byte(id), 32)
}

// scaled returns value times 10^places, refusing a value with more decimals.
func scaled(value string, places int) (*big.Int, error) {
	rat, ok := new(big.Rat).SetString(value)
	if !ok || rat.Sign() < 0 {
		return nil, fmt.Errorf("%q is not a value", value)
	}
	rat.Mul(rat, new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)))
	if !rat.IsInt() {
		return nil, fmt.Errorf("%q has more than %d decimals", value, places)
	}
	return rat.Num(), nil
}
