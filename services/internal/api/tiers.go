// SPDX-License-Identifier: MIT OR Apache-2.0

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limit is a tier's allowance: requests a second, the burst above it, and streams open at once.
type Limit struct {
	Rate    rate.Limit
	Burst   int
	Streams int
}

// The tiers' default allowances: the public tier per client address, the keyed tier per key.
var (
	PublicLimit = Limit{Rate: 10, Burst: 20, Streams: 2}
	KeyedLimit  = Limit{Rate: 1000, Burst: 2000, Streams: 50}
)

// Tiers admits each request to the public tier, limited per client address, or, with an X-API-Key the operator
// issued, to the keyed tier, limited per key. The operator configures keys by their SHA-256 digests, so no key is
// stored where the API runs.
type Tiers struct {
	Public, Keyed Limit
	keys          map[[32]byte]bool
	proxies       []netip.Prefix
	mu            sync.Mutex
	clients       map[string]*client
	now           func() time.Time
}

type client struct {
	limiter *rate.Limiter
	streams int
	seen    time.Time
}

// NewTiers returns tiers with the default limits, admitting the keys whose SHA-256 digests are digests, hex, and
// reading the client's address from X-Forwarded-For where the request comes through one of proxies.
func NewTiers(digests []string, proxies []netip.Prefix) (*Tiers, error) {
	keys := map[[32]byte]bool{}
	for _, digest := range digests {
		raw, err := hex.DecodeString(strings.TrimPrefix(digest, "0x"))
		if err != nil || len(raw) != sha256.Size {
			return nil, fmt.Errorf("an API key digest is not 32 bytes of hex")
		}
		keys[[32]byte(raw)] = true
	}
	return &Tiers{Public: PublicLimit, Keyed: KeyedLimit, keys: keys, proxies: proxies, clients: map[string]*client{},
		now: time.Now}, nil
}

// Digest is the SHA-256 digest of an API key, as the operator configures it.
func Digest(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

// Middleware refuses an unknown key with 401 and a request over its tier's rate with 429, and lets every origin read,
// as the API serves public data and takes no credentials but its keys.
func (t *Tiers) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c, tier, ok := t.admit(w, r)
		if !ok {
			return
		}
		w.Header().Set("X-RateLimit-Tier", tier)
		if delay := c.reserve(t.now()); delay > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(delay.Seconds()))))
			writeError(w, http.StatusTooManyRequests, "rate limited: retry later, or use a key of the keyed tier")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// OpenStream counts a stream against the client's tier, refusing one over its limit with 429, and returns the release
// to call once the stream ends.
func (t *Tiers) OpenStream(w http.ResponseWriter, r *http.Request) (func(), bool) {
	c, tier, ok := t.admit(w, r)
	if !ok {
		return nil, false
	}
	limit := t.Public.Streams
	if tier == "keyed" {
		limit = t.Keyed.Streams
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if c.streams >= limit {
		writeError(w, http.StatusTooManyRequests, "too many streams open for this tier")
		return nil, false
	}
	c.streams++
	return func() {
		t.mu.Lock()
		c.streams--
		t.mu.Unlock()
	}, true
}

func (t *Tiers) admit(w http.ResponseWriter, r *http.Request) (*client, string, bool) {
	if key := r.Header.Get("X-API-Key"); key != "" {
		digest := sha256.Sum256([]byte(key))
		if !t.keys[digest] {
			writeError(w, http.StatusUnauthorized, "unknown API key")
			return nil, "", false
		}
		return t.client("key "+hex.EncodeToString(digest[:]), t.Keyed), "keyed", true
	}
	return t.client("address "+t.address(r), t.Public), "public", true
}

func (t *Tiers) client(id string, limit Limit) *client {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for name, c := range t.clients {
		if c.streams == 0 && now.Sub(c.seen) > 10*time.Minute {
			delete(t.clients, name)
		}
	}
	c, ok := t.clients[id]
	if !ok {
		c = &client{limiter: rate.NewLimiter(limit.Rate, limit.Burst)}
		t.clients[id] = c
	}
	c.seen = now
	return c
}

// reserve takes a request from the client's allowance, or returns how long until one is free.
func (c *client) reserve(now time.Time) time.Duration {
	reservation := c.limiter.ReserveN(now, 1)
	if delay := reservation.DelayFrom(now); delay > 0 {
		reservation.CancelAt(now)
		return delay
	}
	return 0
}

// address is the request's client: its remote address, or, through a trusted proxy, the last address in
// X-Forwarded-For that is not a trusted proxy.
func (t *Tiers) address(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	remote, err := netip.ParseAddr(host)
	if err != nil || !t.trusted(remote) {
		return host
	}
	hops := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
	for i := len(hops) - 1; i >= 0; i-- {
		hop, err := netip.ParseAddr(strings.TrimSpace(hops[i]))
		if err != nil {
			break
		}
		if !t.trusted(hop) {
			return hop.String()
		}
	}
	return host
}

func (t *Tiers) trusted(address netip.Addr) bool {
	for _, prefix := range t.proxies {
		if prefix.Contains(address.Unmap()) {
			return true
		}
	}
	return false
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
