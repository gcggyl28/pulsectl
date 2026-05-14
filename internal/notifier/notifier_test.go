package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/pulsectl/internal/notifier"
)

func TestNotify_Success(t *testing.T) {
	var received notifier.Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notifier.New(ts.URL)
	if err := n.Notify("api", "http://api.example.com/health", "connection refused"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Service != "api" {
		t.Errorf("service = %q, want %q", received.Service, "api")
	}
	if received.Status != "degraded" {
		t.Errorf("status = %q, want %q", received.Status, "degraded")
	}
	if received.Reason != "connection refused" {
		t.Errorf("reason = %q, want %q", received.Reason, "connection refused")
	}
	if received.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
}

func TestNotify_Non2xxResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notifier.New(ts.URL)
	err := n.Notify("db", "http://db.example.com/health", "timeout")
	if err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestNotify_UnreachableWebhook(t *testing.T) {
	n := notifier.New("http://127.0.0.1:0/webhook")
	err := n.Notify("cache", "http://cache.example.com/health", "EOF")
	if err == nil {
		t.Fatal("expected error for unreachable webhook, got nil")
	}
}
