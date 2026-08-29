package handlers_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tuiter.com/api/internal/application/handlers"
)

const (
	testSecret       = "shh-its-a-secret"
	secretHeaderName = "X-Telegram-Bot-Api-Secret-Token"
)

type noopLogger struct{}

func (noopLogger) Printf(_ context.Context, _ string, _ ...any) {}

func newRelay(t *testing.T, upstreamURL string) *handlers.TelegramRelay {
	t.Helper()

	return handlers.NewTelegramRelay(testSecret, upstreamURL, noopLogger{})
}

func TestTelegramRelay_WrongSecret(t *testing.T) {
	relay := newRelay(t, "http://unused.invalid")

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", strings.NewReader(`{}`))
	req.Header.Set(secretHeaderName, "wrong-secret")
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestTelegramRelay_ForwardsBodyAndHeader(t *testing.T) {
	body := `{"update_id":1}`

	var (
		gotBody   []byte
		gotSecret string
	)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		gotSecret = r.Header.Get(secretHeaderName)
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	relay := newRelay(t, upstream.URL)

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", strings.NewReader(body))
	req.Header.Set(secretHeaderName, testSecret)
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if string(gotBody) != body {
		t.Fatalf("expected upstream body %q, got %q", body, gotBody)
	}

	if gotSecret != testSecret {
		t.Fatalf("expected upstream secret header %q, got %q", testSecret, gotSecret)
	}
}

func TestTelegramRelay_UpstreamSuccessReturns204(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()

	relay := newRelay(t, upstream.URL)

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", strings.NewReader(`{}`))
	req.Header.Set(secretHeaderName, testSecret)
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestTelegramRelay_UpstreamFailurePropagatesStatus(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer upstream.Close()

	relay := newRelay(t, upstream.URL)

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", strings.NewReader(`{}`))
	req.Header.Set(secretHeaderName, testSecret)
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestTelegramRelay_ConnectionFailureReturns502(t *testing.T) {
	relay := newRelay(t, "http://127.0.0.1:0")

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", strings.NewReader(`{}`))
	req.Header.Set(secretHeaderName, testSecret)
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestTelegramRelay_OversizedBodyReturns413(t *testing.T) {
	relay := newRelay(t, "http://unused.invalid")

	oversized := bytes.Repeat([]byte("a"), 1<<20+1)

	req := httptest.NewRequest(http.MethodPost, "/v1/telegram/webhook", bytes.NewReader(oversized))
	req.Header.Set(secretHeaderName, testSecret)
	rec := httptest.NewRecorder()

	relay.Relay(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}
