package checker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/pulsectl/internal/checker"
)

func TestCheck_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := checker.New(5 * time.Second)
	result := c.Check(context.Background(), ts.URL)

	if result.Status != checker.StatusHealthy {
		t.Errorf("expected healthy, got %s (error: %s)", result.Status, result.Error)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", result.StatusCode)
	}
	if result.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestCheck_Degraded_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := checker.New(5 * time.Second)
	result := c.Check(context.Background(), ts.URL)

	if result.Status != checker.StatusDegraded {
		t.Errorf("expected degraded, got %s", result.Status)
	}
	if result.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestCheck_Degraded_Unreachable(t *testing.T) {
	c := checker.New(500 * time.Millisecond)
	result := c.Check(context.Background(), "http://127.0.0.1:19999")

	if result.Status != checker.StatusDegraded {
		t.Errorf("expected degraded, got %s", result.Status)
	}
	if result.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestCheck_Degraded_InvalidURL(t *testing.T) {
	c := checker.New(5 * time.Second)
	result := c.Check(context.Background(), "://bad-url")

	if result.Status != checker.StatusDegraded {
		t.Errorf("expected degraded, got %s", result.Status)
	}
}
