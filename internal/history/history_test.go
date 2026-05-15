package history_test

import (
	"testing"

	"github.com/user/pulsectl/internal/history"
)

func TestAdd_And_Latest(t *testing.T) {
	s := history.New(5)

	s.Add("http://example.com", true, "")
	rec, ok := s.Latest("http://example.com")
	if !ok {
		t.Fatal("expected a record, got none")
	}
	if !rec.Healthy {
		t.Errorf("expected healthy=true, got false")
	}
}

func TestLatest_NoRecord(t *testing.T) {
	s := history.New(5)
	_, ok := s.Latest("http://missing.example.com")
	if ok {
		t.Error("expected ok=false for unknown URL")
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	s := history.New(3)
	url := "http://example.com"

	for i := 0; i < 5; i++ {
		s.Add(url, i%2 == 0, "")
	}

	records := s.All(url)
	if len(records) != 3 {
		t.Errorf("expected 3 records after eviction, got %d", len(records))
	}
}

func TestStatusChanged_DetectsFlap(t *testing.T) {
	s := history.New(10)
	url := "http://flappy.example.com"

	s.Add(url, true, "")
	if s.StatusChanged(url) {
		t.Error("single record should not report a change")
	}

	s.Add(url, false, "connection refused")
	if !s.StatusChanged(url) {
		t.Error("expected status change from healthy to degraded")
	}

	s.Add(url, false, "connection refused")
	if s.StatusChanged(url) {
		t.Error("same status twice should not report a change")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := history.New(5)
	url := "http://example.com"
	s.Add(url, true, "")
	s.Add(url, false, "timeout")

	records := s.All(url)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	// Mutate the copy — original must be unaffected.
	records[0].Healthy = false
	orig := s.All(url)
	if !orig[0].Healthy {
		t.Error("All() should return a copy, not a reference to internal slice")
	}
}

func TestNew_DefaultsMaxLen(t *testing.T) {
	s := history.New(0) // invalid value — should default to 10
	url := "http://example.com"
	for i := 0; i < 15; i++ {
		s.Add(url, true, "")
	}
	if got := len(s.All(url)); got != 10 {
		t.Errorf("expected default maxLen=10, retained %d records", got)
	}
}
