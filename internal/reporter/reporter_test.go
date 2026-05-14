package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pulsectl/pulsectl/internal/reporter"
)

func makeStatus(name, url string, healthy bool, errMsg string) reporter.Status {
	return reporter.Status{
		Name:      name,
		URL:       url,
		Healthy:   healthy,
		CheckedAt: time.Now(),
		Error:     errMsg,
	}
}

func TestReport_ContainsOKForHealthy(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	statuses := []reporter.Status{
		makeStatus("api", "http://api.example.com", true, ""),
	}
	r.Report(statuses)
	output := buf.String()
	if !strings.Contains(output, "[OK]") {
		t.Errorf("expected [OK] in output, got:\n%s", output)
	}
}

func TestReport_ContainsDegradedForUnhealthy(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	statuses := []reporter.Status{
		makeStatus("db", "http://db.example.com", false, "connection refused"),
	}
	r.Report(statuses)
	output := buf.String()
	if !strings.Contains(output, "[DEGRADED]") {
		t.Errorf("expected [DEGRADED] in output, got:\n%s", output)
	}
	if !strings.Contains(output, "connection refused") {
		t.Errorf("expected error message in output, got:\n%s", output)
	}
}

func TestReport_DefaultsToStdout(t *testing.T) {
	// Should not panic when nil writer is provided.
	r := reporter.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestSummary_Counts(t *testing.T) {
	statuses := []reporter.Status{
		makeStatus("a", "http://a", true, ""),
		makeStatus("b", "http://b", false, "err"),
		makeStatus("c", "http://c", true, ""),
	}
	h, d := reporter.Summary(statuses)
	if h != 2 {
		t.Errorf("expected 2 healthy, got %d", h)
	}
	if d != 1 {
		t.Errorf("expected 1 degraded, got %d", d)
	}
}

func TestSummary_AllHealthy(t *testing.T) {
	statuses := []reporter.Status{
		makeStatus("x", "http://x", true, ""),
	}
	h, d := reporter.Summary(statuses)
	if h != 1 || d != 0 {
		t.Errorf("unexpected counts: healthy=%d degraded=%d", h, d)
	}
}
