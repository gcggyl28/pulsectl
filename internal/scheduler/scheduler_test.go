package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/pulsectl/internal/scheduler"
)

func TestScheduler_CallsPollImmediately(t *testing.T) {
	var count atomic.Int32

	poll := func(_ context.Context) {
		count.Add(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s := scheduler.New(10*time.Second, poll, nil) // long interval — only immediate call expected
	s.Run(ctx)

	if got := count.Load(); got != 1 {
		t.Errorf("expected 1 immediate call, got %d", got)
	}
}

func TestScheduler_CallsPollOnTick(t *testing.T) {
	var count atomic.Int32

	poll := func(_ context.Context) {
		count.Add(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	// 50 ms interval → 1 immediate + ~3 ticks within 180 ms
	s := scheduler.New(50*time.Millisecond, poll, nil)
	s.Run(ctx)

	got := count.Load()
	if got < 2 {
		t.Errorf("expected at least 2 calls, got %d", got)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	var count atomic.Int32

	poll := func(_ context.Context) {
		count.Add(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := scheduler.New(20*time.Millisecond, poll, nil)

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(60 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// scheduler exited cleanly
	case <-time.After(200 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancellation")
	}

	snap := count.Load()
	if snap < 1 {
		t.Errorf("expected at least 1 call before cancel, got %d", snap)
	}
}
