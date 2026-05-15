// Package history provides a bounded, in-memory record of health-check
// results for each monitored endpoint.
//
// It tracks the N most-recent statuses per service, detects status
// transitions (healthy → degraded and vice-versa), and exposes a
// snapshot of all recorded data for reporting purposes.
//
// Typical usage:
//
//	h := history.New(10)
//	h.Add("api", checker.Status{Healthy: true})
//	if h.StatusChanged("api") {
//		// notify or log the transition
//	}
package history
