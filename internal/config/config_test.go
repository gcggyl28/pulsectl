package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/pulsectl/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "pulsectl-config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `{
		"webhook_url": "https://hooks.example.com/notify",
		"timeout": 5000000000,
		"endpoints": [
			{"name": "api", "url": "https://api.example.com/health", "interval": 10000000000}
		]
	}`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WebhookURL != "https://hooks.example.com/notify" {
		t.Errorf("unexpected webhook_url: %s", cfg.WebhookURL)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected timeout: %v", cfg.Timeout)
	}
	if len(cfg.Endpoints) != 1 || cfg.Endpoints[0].Name != "api" {
		t.Errorf("unexpected endpoints: %+v", cfg.Endpoints)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTempConfig(t, `{
		"webhook_url": "https://hooks.example.com/notify",
		"endpoints": [{"name": "svc", "url": "https://svc.example.com/ping"}]
	}`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", cfg.Timeout)
	}
	if cfg.Endpoints[0].Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", cfg.Endpoints[0].Interval)
	}
}

func TestLoad_MissingWebhook(t *testing.T) {
	path := writeTempConfig(t, `{
		"endpoints": [{"name": "svc", "url": "https://svc.example.com/ping"}]
	}`)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
