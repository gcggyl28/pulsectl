package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/pulsectl/pulsectl/internal/reporter"
)

func TestWriteJSON_ValidOutput(t *testing.T) {
	statuses := []reporter.Status{
		{Name: "api", URL: "http://api", Healthy: true, CheckedAt: time.Now()},
		{Name: "db", URL: "http://db", Healthy: false, CheckedAt: time.Now(), Error: "timeout"},
	}
	var buf bytes.Buffer
	if err := reporter.WriteJSON(&buf, statuses); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result reporter.JSONReport
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if result.Healthy != 1 {
		t.Errorf("expected 1 healthy, got %d", result.Healthy)
	}
	if result.Degraded != 1 {
		t.Errorf("expected 1 degraded, got %d", result.Degraded)
	}
	if len(result.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(result.Services))
	}
}

func TestWriteJSON_ErrorFieldOmittedWhenEmpty(t *testing.T) {
	statuses := []reporter.Status{
		{Name: "svc", URL: "http://svc", Healthy: true, CheckedAt: time.Now()},
	}
	var buf bytes.Buffer
	if err := reporter.WriteJSON(&buf, statuses); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if bytes.Contains([]byte(output), []byte(`"error"`)) {
		t.Errorf("expected 'error' field to be omitted for healthy service, got:\n%s", output)
	}
}

func TestWriteJSON_EmptyStatuses(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.WriteJSON(&buf, []reporter.Status{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result reporter.JSONReport
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if result.Healthy != 0 || result.Degraded != 0 {
		t.Errorf("expected zero counts for empty statuses")
	}
}
