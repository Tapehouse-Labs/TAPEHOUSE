// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/time/rate"
)

// Limits of an opt-in.
const (
	SignatureMaxAge  = 10 * time.Minute
	EmailTokenTTL    = 24 * time.Hour
	TelegramCodeTTL  = time.Hour
	DefaultThreshold = 1250
	MinThreshold     = 1050
	MaxThreshold     = 3000
	codeLength       = 24
	maxBody          = 4096
	codeAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
)

// Links are the URLs an email carries: the unsubscribe link of every alert and the confirmation of a new subscription.
// A subscription's unsubscribe token is derived from its ID and a secret, so each email can carry it while the store
// keeps only its digest.
type Links struct {
	BaseURL string
	Secret  []byte
}

// UnsubscribeToken is the token that unsubscribes subscription id.
func (l Links) UnsubscribeToken(id string) string {
	mac := hmac.New(sha256.New, l.Secret)
	mac.Write([]byte(id))
	return hex.EncodeToString(mac.Sum(nil))
}

// Unsubscribe is the link that unsubscribes subscription id.
func (l Links) Unsubscribe(id string) string {
	return l.BaseURL + "/v1/unsubscribe/" + l.UnsubscribeToken(id)
}

// Digest is the SHA-256 the store keeps of a token or a code.
func Digest(token string) [32]byte {
	return sha256.Sum256([]byte(token))
}

func (l Links) digestOf(id string) [32]byte {
	return Digest(l.UnsubscribeToken(id))
}

// SubscribeText is what an account signs to subscribe to channel.
func SubscribeText(chainID uint64, account common.Address, channel string, issuedAt int64) string {
	return fmt.Sprintf("Tapehouse alerts\naction: subscribe\nchain: %d\naccount: %s\nchannel: %s\nissued: %d",
		chainID, account.Hex(), channel, issuedAt)
}

// UnsubscribeText is what an account signs to delete subscription id.
func UnsubscribeText(chainID uint64, account common.Address, id string, issuedAt int64) string {
	return fmt.Sprintf("Tapehouse alerts\naction: unsubscribe\nchain: %d\naccount: %s\nsubscription: %s\nissued: %d",
		chainID, account.Hex(), id, issuedAt)
}

// Signer is who signed text, an EIP-191 personal_sign signature.
func Signer(text string, signature []byte) (common.Address, error) {
	if len(signature) != 65 {
		return common.Address{}, errors.New("a signature is 65 bytes")
	}
	sig := append([]byte(nil), signature...)
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	key, err := crypto.SigToPub(accounts.TextHash([]byte(text)), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*key), nil
}

// Chatbot is a Telegram bot: a channel, and the name its deep links open.
type Chatbot interface {
	Channel
	Username(ctx context.Context) (string, error)
}

// Server is the opt-in: it takes signed subscriptions, confirms emails by a link and Telegram chats by a code, and
// unsubscribes.
type Server struct {
	Store   *Store
	ChainID uint64
	Links   Links
	Email   Channel
	Bot     Chatbot
	Now     func() time.Time
	Log     *slog.Logger
	limiter *rate.Limiter
}

// NewServer returns the opt-in over store. A nil email or bot leaves its channel closed.
func NewServer(store *Store, chainID uint64, links Links, email Channel, bot Chatbot, log *slog.Logger) *Server {
	return &Server{Store: store, ChainID: chainID, Links: links, Email: email, Bot: bot, Now: time.Now, Log: log,
		limiter: rate.NewLimiter(5, 20)}
}

// Handler serves POST /v1/subscriptions, DELETE /v1/subscriptions/{id}, GET /v1/confirm/{token} and
// GET /v1/unsubscribe/{token}.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/subscriptions", s.subscribe)
	mux.HandleFunc("DELETE /v1/subscriptions/{id}", s.delete)
	mux.HandleFunc("GET /v1/confirm/{token}", s.confirm)
	mux.HandleFunc("GET /v1/unsubscribe/{token}", s.unsubscribe)
	return mux
}

type request struct {
	Account     common.Address  `json:"account"`
	Channel     string          `json:"channel"`
	Destination string          `json:"destination"`
	Threshold   json.RawMessage `json:"threshold"`
	IssuedAt    int64           `json:"issuedAt"`
	Signature   hexutil.Bytes   `json:"signature"`
}

func (s *Server) fresh(issuedAt int64) bool {
	age := s.Now().Sub(time.Unix(issuedAt, 0))
	return age <= SignatureMaxAge && age >= -time.Minute
}

func thousandths(raw json.RawMessage) (int, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return DefaultThreshold, nil
	}
	value, err := strconv.ParseFloat(string(raw), 64)
	if err != nil || math.IsNaN(value) {
		return 0, errors.New("the threshold is not a number")
	}
	scaled := int(math.Round(value * 1000))
	if scaled < MinThreshold || scaled > MaxThreshold {
		return 0, errors.New("the threshold is between 1.05 and 3")
	}
	return scaled, nil
}

func fail(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow() {
		fail(w, http.StatusTooManyRequests, "too many requests")
		return
	}
	var req request
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBody)).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "the body is not a subscription")
		return
	}
	if !s.fresh(req.IssuedAt) {
		fail(w, http.StatusUnauthorized, "issuedAt is not within the last 10 minutes")
		return
	}
	signer, err := Signer(SubscribeText(s.ChainID, req.Account, req.Channel, req.IssuedAt), req.Signature)
	if err != nil || signer != req.Account {
		fail(w, http.StatusUnauthorized, "the signature is not the account's")
		return
	}
	threshold, err := thousandths(req.Threshold)
	if err != nil {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}
	destination := ""
	switch req.Channel {
	case Email:
		address, err := mail.ParseAddress(req.Destination)
		if err != nil || s.Email == nil || address.Name != "" || address.Address != req.Destination || len(req.Destination) > 254 ||
			strings.ContainsAny(req.Destination, "\r\n<>") {
			fail(w, http.StatusBadRequest, "the destination is not an email address, or email is closed")
			return
		}
		destination = req.Destination
	case Telegram:
		if s.Bot == nil {
			fail(w, http.StatusBadRequest, "Telegram is closed")
			return
		}
	default:
		fail(w, http.StatusBadRequest, "the channel is email or telegram")
		return
	}
	id, err := s.Store.Subscribe(r.Context(), req.Account, req.Channel, destination, threshold, s.Links.digestOf)
	if err != nil {
		s.Log.Error("subscribe", "error", err)
		fail(w, http.StatusInternalServerError, "the subscription was not kept")
		return
	}
	out := map[string]any{"id": id, "status": "pending"}
	if req.Channel == Email {
		token, err := s.pending(r.Context(), id, Email, EmailTokenTTL, func() (string, error) { return randomHex(32) })
		if err == nil {
			err = s.Email.Send(r.Context(), destination, Message{Subject: "Confirm your Tapehouse alerts",
				Body: "Follow this link to receive Tapehouse alerts for " + req.Account.Hex() + ":\n\n" + s.Links.BaseURL + "/v1/confirm/" + token +
					"\n\nIt expires in 24 hours. If you did not ask for it, ignore this email.\n\nUnsubscribe: " + s.Links.Unsubscribe(id) + "\n"})
		}
		if err != nil {
			s.Log.Error("confirmation", "subscription", id, "error", err)
			_ = s.Store.Delete(r.Context(), id)
			fail(w, http.StatusBadGateway, "the confirmation email was not sent")
			return
		}
	} else {
		code, err := s.pending(r.Context(), id, Telegram, TelegramCodeTTL, randomCode)
		bot, nameErr := "", error(nil)
		if err == nil {
			bot, nameErr = s.Bot.Username(r.Context())
		}
		if err != nil || nameErr != nil {
			s.Log.Error("telegram link", "subscription", id)
			_ = s.Store.Delete(r.Context(), id)
			fail(w, http.StatusBadGateway, "the Telegram link was not made")
			return
		}
		out["code"] = code
		out["link"] = "https://t.me/" + bot + "?start=" + code
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) pending(ctx context.Context, id, kind string, ttl time.Duration, generate func() (string, error)) (string, error) {
	token, err := generate()
	if err != nil {
		return "", err
	}
	return token, s.Store.AddPending(ctx, id, kind, Digest(token), s.Now().Add(ttl))
}

func randomHex(n int) (string, error) {
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func randomCode() (string, error) {
	raw := make([]byte, codeLength)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	for i, b := range raw {
		raw[i] = codeAlphabet[b&63]
	}
	return string(raw), nil
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBody)).Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "the body is not a signature")
		return
	}
	id := r.PathValue("id")
	sub, found, err := s.Store.Get(r.Context(), id)
	if err != nil {
		fail(w, http.StatusInternalServerError, "the subscription was not read")
		return
	}
	signer, signErr := Signer(UnsubscribeText(s.ChainID, sub.Account, id, req.IssuedAt), req.Signature)
	if !found || !s.fresh(req.IssuedAt) || signErr != nil || signer != sub.Account {
		fail(w, http.StatusUnauthorized, "the signature is not the account's, or is stale")
		return
	}
	if err := s.Store.Delete(r.Context(), id); err != nil {
		fail(w, http.StatusInternalServerError, "the subscription was not deleted")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) confirm(w http.ResponseWriter, r *http.Request) {
	id, ok, err := s.Store.Redeem(r.Context(), Email, Digest(r.PathValue("token")), s.Now())
	if err == nil && ok {
		ok, err = s.Store.Confirm(r.Context(), id, "")
	}
	if err != nil || !ok {
		fail(w, http.StatusNotFound, "this link has expired or was used")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, "Confirmed. You will get Tapehouse alerts by email.\n")
}

func (s *Server) unsubscribe(w http.ResponseWriter, r *http.Request) {
	ok, err := s.Store.DeleteByDigest(r.Context(), Digest(r.PathValue("token")))
	if err != nil || !ok {
		fail(w, http.StatusNotFound, "this link is not a subscription's")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, "Unsubscribed. Nothing is kept.\n")
}

// Chat handles a message a Telegram chat sent the bot: /start with a code links the chat to the subscription the code
// confirms, and /stop deletes every subscription of the chat. Any other message is ignored.
func (s *Server) Chat(ctx context.Context, chat, text string) {
	command, argument, _ := strings.Cut(strings.TrimSpace(text), " ")
	command, _, _ = strings.Cut(command, "@")
	switch command {
	case "/start":
		id, ok, err := s.Store.Redeem(ctx, Telegram, Digest(strings.TrimSpace(argument)), s.Now())
		if err == nil && ok {
			ok, err = s.Store.Confirm(ctx, id, chat)
		}
		reply := "This code has expired or was used. Ask for a new one."
		if err == nil && ok {
			reply = "Linked. You will get Tapehouse alerts here. Send /stop to stop them."
		}
		_ = s.Bot.Send(ctx, chat, Message{Body: reply})
	case "/stop":
		if err := s.Store.DeleteChat(ctx, chat); err == nil {
			_ = s.Bot.Send(ctx, chat, Message{Body: "Stopped. Nothing is kept."})
		}
	}
}
