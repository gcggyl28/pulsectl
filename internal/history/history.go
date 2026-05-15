// Package history maintains a rolling in-memory record of health-check
// results so that callers can detect flapping services and avoid sending
// redundant webhook notifications.
package history

import (
	"sync"
	"time"
)

// Record holds a single health-check outcome for one endpoint.
type Record struct {
	Timestamp time.Time
	Healthy   bool
	Error     string
}

// Store keeps the last N records per endpoint URL.
type Store struct {
	mu      sync.RWMutex
	records map[string][]Record
	maxLen  int
}

// New returns a Store that retains at most maxLen records per endpoint.
// If maxLen is less than 1 it defaults to 10.
func New(maxLen int) *Store {
	if maxLen < 1 {
		maxLen = 10
	}
	return &Store{
		records: make(map[string][]Record),
		maxLen:  maxLen,
	}
}

// Add appends a new record for the given endpoint URL, evicting the oldest
// entry when the buffer is full.
func (s *Store) Add(url string, healthy bool, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := Record{Timestamp: time.Now(), Healthy: healthy, Error: errMsg}
	buf := append(s.records[url], r)
	if len(buf) > s.maxLen {
		buf = buf[len(buf)-s.maxLen:]
	}
	s.records[url] = buf
}

// Latest returns the most recent record for url and true, or the zero
// value and false when no record exists yet.
func (s *Store) Latest(url string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf := s.records[url]
	if len(buf) == 0 {
		return Record{}, false
	}
	return buf[len(buf)-1], true
}

// All returns a copy of all records stored for url.
func (s *Store) All(url string) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf := s.records[url]
	out := make([]Record, len(buf))
	copy(out, buf)
	return out
}

// StatusChanged reports whether the health status of url differs between
// the most recent record and the one before it. Returns false when fewer
// than two records exist.
func (s *Store) StatusChanged(url string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	buf := s.records[url]
	if len(buf) < 2 {
		return false
	}
	return buf[len(buf)-1].Healthy != buf[len(buf)-2].Healthy
}
