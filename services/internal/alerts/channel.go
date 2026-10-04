// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Message is an alert as plain text.
type Message struct {
	Subject string
	Body    string
}

// Channel sends a message to a destination: an email address or a Telegram chat ID.
type Channel interface {
	Send(ctx context.Context, destination string, m Message) error
}

const timeout = 30 * time.Second

// SMTP is an email channel over SMTP with STARTTLS: a server that does not offer it is not sent to, and nothing,
// the password included, reaches it.
type SMTP struct {
	// TLS is the configuration of the STARTTLS handshake; its ServerName is the server's host unless set.
	TLS      *tls.Config
	address  string
	host     string
	user     string
	password string
	from     mail.Address
}

// NewEmail returns the email channel of the server at rawURL, smtp://user:password@host:port, sending as from.
func NewEmail(rawURL, from string) (*SMTP, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "smtp" || u.Hostname() == "" {
		return nil, errors.New("the SMTP URL is not smtp://user:password@host:port")
	}
	sender, err := mail.ParseAddress(from)
	if err != nil {
		return nil, errors.New("the sender is not an email address")
	}
	port := u.Port()
	if port == "" {
		port = "587"
	}
	password, _ := u.User.Password()
	return &SMTP{address: net.JoinHostPort(u.Hostname(), port), host: u.Hostname(), user: u.User.Username(), password: password,
		from: *sender}, nil
}

func lineBreak(texts ...string) bool {
	for _, text := range texts {
		if strings.ContainsAny(text, "\r\n") {
			return true
		}
	}
	return false
}

// failure reports what went wrong without what the server said, which may name the recipient.
func failure(stage string, err error) error {
	var reply *textproto.Error
	if errors.As(err, &reply) {
		return fmt.Errorf("email %s: %d", stage, reply.Code)
	}
	return fmt.Errorf("email %s failed", stage)
}

// Send sends m as a plain-text email to destination.
func (s *SMTP) Send(ctx context.Context, destination string, m Message) error {
	if lineBreak(destination, m.Subject) {
		return errors.New("an email address or subject has a line break")
	}
	to, err := mail.ParseAddress(destination)
	if err != nil {
		return errors.New("the destination is not an email address")
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", s.address)
	if err != nil {
		return failure("connect", nil)
	}
	deadline, _ := ctx.Deadline()
	_ = conn.SetDeadline(deadline)
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_ = conn.Close()
		return failure("greeting", err)
	}
	defer func() { _ = client.Close() }()
	if err := client.Hello("localhost"); err != nil {
		return failure("hello", err)
	}
	if ok, _ := client.Extension("STARTTLS"); !ok {
		return errors.New("email: the server does not offer STARTTLS")
	}
	config := &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}
	if s.TLS != nil {
		config = s.TLS.Clone()
		if config.ServerName == "" {
			config.ServerName = s.host
		}
	}
	if err := client.StartTLS(config); err != nil {
		return failure("starttls", err)
	}
	if s.user != "" {
		if err := client.Auth(smtp.PlainAuth("", s.user, s.password, s.host)); err != nil {
			return failure("auth", err)
		}
	}
	if err := client.Mail(s.from.Address); err != nil {
		return failure("sender", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return failure("recipient", err)
	}
	w, err := client.Data()
	if err != nil {
		return failure("data", err)
	}
	if _, err := io.WriteString(w, s.format(to.Address, m)); err != nil {
		return failure("data", err)
	}
	if err := w.Close(); err != nil {
		return failure("data", err)
	}
	_ = client.Quit()
	return nil
}

func (s *SMTP) format(to string, m Message) string {
	body := strings.ReplaceAll(strings.ReplaceAll(m.Body, "\r\n", "\n"), "\n", "\r\n")
	return "From: " + s.from.String() + "\r\nTo: " + (&mail.Address{Address: to}).String() + "\r\nSubject: " +
		mime.QEncoding.Encode("utf-8", m.Subject) + "\r\nDate: " + time.Now().UTC().Format(time.RFC1123Z) +
		"\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n" + body
}

// Bot is a Telegram bot reached over the Bot API.
type Bot struct {
	token  string
	base   string
	client *http.Client
	mu     sync.Mutex
	name   string
}

// NewTelegram returns the Telegram channel of the bot with token, over the Bot API at baseURL, such as
// https://api.telegram.org.
func NewTelegram(token, baseURL string) *Bot {
	return &Bot{token: token, base: strings.TrimRight(baseURL, "/"), client: &http.Client{Timeout: 2 * timeout}}
}

type reply struct {
	OK        bool            `json:"ok"`
	ErrorCode int             `json:"error_code"`
	Result    json.RawMessage `json:"result"`
}

// call posts a Bot API method. Its errors carry the method and the status, never the URL, which holds the token, or
// the description, which may name a chat.
func (b *Bot) call(ctx context.Context, method string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, b.base+"/bot"+b.token+"/"+method, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("telegram %s: bad request", method)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := b.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("telegram %s failed", method)
	}
	defer func() { _ = response.Body.Close() }()
	var out reply
	_ = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&out)
	if response.StatusCode != http.StatusOK || !out.OK {
		return nil, fmt.Errorf("telegram %s: %d", method, response.StatusCode)
	}
	return out.Result, nil
}

// Send sends m as plain text to the chat with ID destination.
func (b *Bot) Send(ctx context.Context, destination string, m Message) error {
	text := m.Body
	if m.Subject != "" {
		text = m.Subject + "\n\n" + m.Body
	}
	_, err := b.call(ctx, "sendMessage", map[string]string{"chat_id": destination, "text": text})
	return err
}

// Username is the bot's name, which its deep links open.
func (b *Bot) Username(ctx context.Context) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.name != "" {
		return b.name, nil
	}
	result, err := b.call(ctx, "getMe", struct{}{})
	if err != nil {
		return "", err
	}
	var me struct {
		Username string `json:"username"`
	}
	if err := json.Unmarshal(result, &me); err != nil || me.Username == "" {
		return "", errors.New("telegram getMe: no username")
	}
	b.name = me.Username
	return b.name, nil
}

// Update is a message sent to the bot: its update ID, the chat's ID and the text. Anything else has no text.
type Update struct {
	ID   int64
	Chat string
	Text string
}

// Updates long-polls the messages sent to the bot from update ID offset, for at most waitS seconds.
func (b *Bot) Updates(ctx context.Context, offset int64, waitS int) ([]Update, error) {
	result, err := b.call(ctx, "getUpdates", map[string]any{"offset": offset, "timeout": waitS, "allowed_updates": []string{"message"}})
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID      int64 `json:"update_id"`
		Message *struct {
			Text string `json:"text"`
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	}
	if err := json.Unmarshal(result, &raw); err != nil {
		return nil, errors.New("telegram getUpdates: unreadable")
	}
	out := make([]Update, len(raw))
	for i, u := range raw {
		out[i] = Update{ID: u.ID}
		if u.Message != nil {
			out[i].Chat, out[i].Text = strconv.FormatInt(u.Message.Chat.ID, 10), u.Message.Text
		}
	}
	return out, nil
}
