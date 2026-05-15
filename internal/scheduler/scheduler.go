// Package scheduler provides a simple ticker-based loop that invokes
// a polling function at a fixed interval until the context is cancelled.
package scheduler

import (
	"context"
	"log"
	"time"
)

// PollFunc is the function called on every tick.
type PollFunc func(ctx context.Context)

// Scheduler drives periodic execution of a PollFunc.
type Scheduler struct {
	interval time.Duration
	poll     PollFunc
	logger   *log.Logger
}

// New creates a Scheduler that calls poll at the given interval.
// If logger is nil, log.Default() is used.
func New(interval time.Duration, poll PollFunc, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		poll:     poll,
		logger:   logger,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
// It fires immediately on the first tick, then waits for the interval.
func (s *Scheduler) Run(ctx context.Context) {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Fire immediately before waiting for the first tick.
	s.poll(ctx)

	for {
		select {
		case <-ticker.C:
			s.poll(ctx)
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return
		}
	}
}
