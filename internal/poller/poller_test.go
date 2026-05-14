package poller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/config"
	"github.com/user/pulsectl/internal/poller"
)

func TestPoller_NotifiesOnDegradedService(t *testing.T) {
	var notifyCount int32

	// Webhook receiver
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["status"] == "degraded" {
			atomic.AddInt32(&notifyCount, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer webhook.Close()

	cfg := &config.Config{
		WebhookURL: webhook.URL,
		Interval:   50 * time.Millisecond,
		Timeout:    2 * time.Second,
		Services: []config.Service{
			{Name: "broken", URL: "http://127.0.0.1:0/health"},
		},
	}

	p := poller.New(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	p.Run(ctx)

	if atomic.LoadInt32(&notifyCount) == 0 {
		t.Error("expected at least one degraded notification, got none")
	}
}

func TestPoller_NoNotifyOnHealthyService(t *testing.T) {
	var notifyCount int32

	// Healthy target
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&notifyCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer webhook.Close()

	cfg := &config.Config{
		WebhookURL: webhook.URL,
		Interval:   50 * time.Millisecond,
		Timeout:    2 * time.Second,
		Services: []config.Service{
			{Name: "healthy-svc", URL: target.URL + "/health"},
		},
	}

	p := poller.New(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	p.Run(ctx)

	if atomic.LoadInt32(&notifyCount) != 0 {
		t.Errorf("expected zero notifications for healthy service, got %d", notifyCount)
	}
}
