package handlers

import (
	"bytes"
	"crypto/subtle"
	"io"
	"net/http"
	"time"

	"tuiter.com/api/pkg/logging"
)

const (
	maxTelegramBodyBytes  = 1 << 20 // ~1MB
	telegramClientTimeout = 45 * time.Second
	telegramSecretHeader  = "X-Telegram-Bot-Api-Secret-Token"
)

// NewTelegramRelay builds a handler that forwards Telegram webhook updates to Laravel unchanged.
func NewTelegramRelay(secret, upstreamURL string, logger logging.ContextualLogger) *TelegramRelay {
	return &TelegramRelay{
		secret:      secret,
		upstreamURL: upstreamURL,
		client:      &http.Client{Timeout: telegramClientTimeout},
		logger:      logger,
	}
}

type TelegramRelay struct {
	secret      string
	upstreamURL string
	client      *http.Client
	logger      logging.ContextualLogger
}

func (t *TelegramRelay) Relay(w http.ResponseWriter, r *http.Request) {
	// Defense-in-depth: chi only mounts Relay via .Post today, but this
	// keeps the handler safe if it's ever rewired to a non-method-scoped route.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	// ponytail: empty configured secret must never match an empty header.
	secretMatches := subtle.ConstantTimeCompare(
		[]byte(r.Header.Get(telegramSecretHeader)),
		[]byte(t.secret),
	) == 1
	if t.secret == "" || !secretMatches {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxTelegramBodyBytes+1))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(body) > maxTelegramBodyBytes {
		w.WriteHeader(http.StatusRequestEntityTooLarge)

		return
	}

	start := time.Now()

	upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, t.upstreamURL, bytes.NewReader(body))
	if err != nil {
		t.logger.Printf(r.Context(), "telegram relay build error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set(telegramSecretHeader, r.Header.Get(telegramSecretHeader))

	resp, err := t.client.Do(upstreamReq)
	duration := time.Since(start)

	if err != nil {
		t.logger.Printf(r.Context(), "telegram relay upstream error status=%d duration=%s", http.StatusBadGateway, duration)
		w.WriteHeader(http.StatusBadGateway)

		return
	}
	defer func() { _ = resp.Body.Close() }()

	_, _ = io.Copy(io.Discard, resp.Body)

	t.logger.Printf(r.Context(), "telegram relay status=%d duration=%s", resp.StatusCode, duration)

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	w.WriteHeader(resp.StatusCode)
}
