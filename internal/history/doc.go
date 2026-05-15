// Package history provides an in-memory, thread-safe rolling buffer of
// health-check results keyed by endpoint URL.
//
// # Overview
//
// A [Store] retains at most N [Record] values per endpoint. Once the buffer
// is full the oldest entry is evicted automatically, keeping memory usage
// bounded regardless of how long pulsectl runs.
//
// # Typical usage
//
//	store := history.New(20)          // keep last 20 checks per endpoint
//
//	// after each check:
//	store.Add(url, healthy, errMsg)
//
//	// before notifying — skip if status hasn't changed:
//	if store.StatusChanged(url) {
//	    notifier.Notify(ctx, payload)
//	}
//
// The zero value is not usable; always construct a Store with [New].
package history
