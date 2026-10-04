// SPDX-License-Identifier: MIT OR Apache-2.0

package alerts

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// smtpDouble is a local SMTP server: it offers STARTTLS, takes AUTH PLAIN only over TLS, and keeps what it was sent.
type smtpDouble struct {
	listener net.Listener
	config   *tls.Config
	pool     *x509.CertPool
	startTLS bool
	rejectTo string
	mu       sync.Mutex
	secure   bool
	user     string
	password string
	from     string
	to       []string
	data     string
	commands []string
}

func newSMTPDouble(t *testing.T, startTLS bool) *smtpDouble {
	t.Helper()
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	template := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "127.0.0.1"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IsCA: true, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	certificate, _ := x509.ParseCertificate(der)
	pool := x509.NewCertPool()
	pool.AddCert(certificate)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &smtpDouble{listener: listener, pool: pool, startTLS: startTLS,
		config: &tls.Config{Certificates: []tls.Certificate{{Certificate: [][]byte{der}, PrivateKey: key}}}}
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go s.serve(conn)
		}
	}()
	return s
}

func (s *smtpDouble) serve(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	reader, secure := bufio.NewReader(conn), false
	say := func(format string, args ...any) { _, _ = fmt.Fprintf(conn, format+"\r\n", args...) }
	say("220 localhost ready")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		command := strings.ToUpper(strings.Fields(line + " x")[0])
		s.mu.Lock()
		s.commands = append(s.commands, command)
		s.mu.Unlock()
		switch command {
		case "EHLO":
			say("250-localhost")
			if s.startTLS && !secure {
				say("250-STARTTLS")
			}
			say("250 AUTH PLAIN")
		case "STARTTLS":
			say("220 go ahead")
			tlsConn := tls.Server(conn, s.config)
			if tlsConn.Handshake() != nil {
				return
			}
			conn, reader, secure = tlsConn, bufio.NewReader(tlsConn), true
			say = func(format string, args ...any) { _, _ = fmt.Fprintf(tlsConn, format+"\r\n", args...) }
		case "AUTH":
			raw, _ := base64.StdEncoding.DecodeString(strings.Fields(line)[2])
			parts := strings.Split(string(raw), "\x00")
			s.mu.Lock()
			s.secure, s.user, s.password = secure, parts[1], parts[2]
			s.mu.Unlock()
			say("235 ok")
		case "MAIL":
			s.mu.Lock()
			s.from = line
			s.mu.Unlock()
			say("250 ok")
		case "RCPT":
			if s.rejectTo != "" {
				say("550 5.1.1 %s: user unknown", s.rejectTo)
				continue
			}
			s.mu.Lock()
			s.to = append(s.to, line)
			s.mu.Unlock()
			say("250 ok")
		case "DATA":
			say("354 go")
			var data strings.Builder
			for {
				text, err := reader.ReadString('\n')
				if err != nil || text == ".\r\n" {
					break
				}
				data.WriteString(text)
			}
			s.mu.Lock()
			s.data = data.String()
			s.mu.Unlock()
			say("250 queued")
		case "QUIT":
			say("221 bye")
			return
		default:
			say("250 ok")
		}
	}
}

func (s *smtpDouble) url(user, password string) string {
	return fmt.Sprintf("smtp://%s:%s@%s", user, password, s.listener.Addr())
}

func TestAnEmailGoesOverSMTPWithSTARTTLS(t *testing.T) {
	server := newSMTPDouble(t, true)
	email, err := NewEmail(server.url("alerts", "s3cret"), "Tapehouse <alerts@tapehouse.example>")
	if err != nil {
		t.Fatal(err)
	}
	email.TLS = &tls.Config{RootCAs: server.pool}
	if err := email.Send(context.Background(), "ada@example.org", Message{Subject: "Weekend: NVDA", Body: "Line one\nLine two\n"}); err != nil {
		t.Fatal(err)
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	if !server.secure || server.user != "alerts" || server.password != "s3cret" {
		t.Fatalf("authenticated over TLS %v as %q", server.secure, server.user)
	}
	if !strings.Contains(server.from, "alerts@tapehouse.example") || len(server.to) != 1 || !strings.Contains(server.to[0], "ada@example.org") {
		t.Fatalf("envelope %q %q", server.from, server.to)
	}
	for _, want := range []string{"From: \"Tapehouse\" <alerts@tapehouse.example>", "To: <ada@example.org>", "Subject: Weekend: NVDA",
		"Content-Type: text/plain; charset=utf-8", "\r\n\r\nLine one\r\nLine two\r\n"} {
		if !strings.Contains(server.data, want) {
			t.Fatalf("the message lacks %q:\n%s", want, server.data)
		}
	}
}

func TestAnEmailIsNotSentWithoutSTARTTLS(t *testing.T) {
	server := newSMTPDouble(t, false)
	email, _ := NewEmail(server.url("alerts", "s3cret"), "alerts@tapehouse.example")
	email.TLS = &tls.Config{RootCAs: server.pool}
	if err := email.Send(context.Background(), "ada@example.org", Message{Subject: "x", Body: "y"}); err == nil {
		t.Fatal("sent in the clear")
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.password != "" || server.data != "" {
		t.Fatalf("the password or a message reached a server without STARTTLS: %q", server.commands)
	}
}

func TestAnEmailsErrorNamesNeitherTheDestinationNorThePassword(t *testing.T) {
	server := newSMTPDouble(t, true)
	server.rejectTo = "ada@example.org"
	email, _ := NewEmail(server.url("alerts", "s3cret"), "alerts@tapehouse.example")
	email.TLS = &tls.Config{RootCAs: server.pool}
	err := email.Send(context.Background(), "ada@example.org", Message{Subject: "x", Body: "y"})
	if err == nil {
		t.Fatal("a rejected recipient is an error")
	}
	if strings.Contains(err.Error(), "ada@example.org") || strings.Contains(err.Error(), "s3cret") || !strings.Contains(err.Error(), "550") {
		t.Fatalf("error %q", err)
	}
	if err := email.Send(context.Background(), "ada@example.org\r\nBcc: eve@example.org", Message{Subject: "x", Body: "y"}); err == nil {
		t.Fatal("a destination with a line break")
	}
	if err := email.Send(context.Background(), "ada@example.org", Message{Subject: "x\r\nBcc: eve@example.org", Body: "y"}); err == nil {
		t.Fatal("a subject with a line break")
	}
	if _, err := NewEmail("http://alerts:s3cret@127.0.0.1:587", "alerts@tapehouse.example"); err == nil || strings.Contains(err.Error(), "s3cret") {
		t.Fatalf("a URL that is not smtp://: %v", err)
	}
}

type botDouble struct {
	t       *testing.T
	server  *httptest.Server
	mu      sync.Mutex
	sent    []map[string]any
	updates []map[string]any
	status  int
	offsets []int64
}

func newBotDouble(t *testing.T) *botDouble {
	b := &botDouble{t: t, status: http.StatusOK}
	b.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b.mu.Lock()
		defer b.mu.Unlock()
		if !strings.HasPrefix(r.URL.Path, "/botTOKEN/") {
			http.NotFound(w, r)
			return
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if b.status != http.StatusOK {
			w.WriteHeader(b.status)
			_, _ = io.WriteString(w, `{"ok":false,"error_code":403,"description":"Forbidden: chat 42 blocked the bot"}`)
			return
		}
		switch strings.TrimPrefix(r.URL.Path, "/botTOKEN/") {
		case "sendMessage":
			b.sent = append(b.sent, body)
			_, _ = io.WriteString(w, `{"ok":true,"result":{}}`)
		case "getMe":
			_, _ = io.WriteString(w, `{"ok":true,"result":{"username":"tapehouse_bot"}}`)
		case "getUpdates":
			b.offsets = append(b.offsets, int64(body["offset"].(float64)))
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": b.updates})
			b.updates = nil
		}
	}))
	t.Cleanup(b.server.Close)
	return b
}

func TestATelegramMessageGoesThroughTheBotAPI(t *testing.T) {
	double := newBotDouble(t)
	bot := NewTelegram("TOKEN", double.server.URL)
	if err := bot.Send(context.Background(), "42", Message{Subject: "Threshold: NVDA", Body: "Ratio 1.20"}); err != nil {
		t.Fatal(err)
	}
	if len(double.sent) != 1 || double.sent[0]["chat_id"] != "42" || double.sent[0]["text"] != "Threshold: NVDA\n\nRatio 1.20" {
		t.Fatalf("sent %v", double.sent)
	}
	if name, err := bot.Username(context.Background()); err != nil || name != "tapehouse_bot" {
		t.Fatalf("username %q, %v", name, err)
	}
	double.updates = []map[string]any{
		{"update_id": 7, "message": map[string]any{"text": "/start abc", "chat": map[string]any{"id": 42}}},
		{"update_id": 8, "edited_message": map[string]any{"text": "ignored"}},
		{"update_id": 9, "message": map[string]any{"text": "/stop", "chat": map[string]any{"id": -1001}}},
	}
	updates, err := bot.Updates(context.Background(), 5, 0)
	if err != nil || len(updates) != 3 {
		t.Fatalf("updates %v, %v", updates, err)
	}
	if updates[0].ID != 7 || updates[0].Chat != "42" || updates[0].Text != "/start abc" || updates[1].Text != "" || updates[2].Chat != "-1001" {
		t.Fatalf("updates %+v", updates)
	}
	if double.offsets[0] != 5 {
		t.Fatalf("offset %v", double.offsets)
	}
}

func TestATelegramErrorNamesNeitherTheTokenNorTheChat(t *testing.T) {
	double := newBotDouble(t)
	double.status = http.StatusForbidden
	bot := NewTelegram("TOKEN", double.server.URL)
	err := bot.Send(context.Background(), "42", Message{Body: "x"})
	if err == nil {
		t.Fatal("a refused message is an error")
	}
	for _, secret := range []string{"TOKEN", "42", double.server.URL} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("error %q names %q", err, secret)
		}
	}
	if !strings.Contains(err.Error(), "403") {
		t.Fatalf("error %q", err)
	}
	dead := NewTelegram("TOKEN", "http://127.0.0.1:1")
	if err := dead.Send(context.Background(), "42", Message{Body: "x"}); err == nil || strings.Contains(err.Error(), "TOKEN") || strings.Contains(err.Error(), "42") {
		t.Fatalf("error %v", err)
	}
}
