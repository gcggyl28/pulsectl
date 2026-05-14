package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status represents the health status of an endpoint.
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusDegraded Status = "degraded"
)

// Result holds the outcome of a single health check.
type Result struct {
	Endpoint   string
	Status     Status
	StatusCode int
	Latency    time.Duration
	Error      string
	CheckedAt  time.Time
}

// Checker polls HTTP endpoints and returns health results.
type Checker struct {
	client  *http.Client
	timeout time.Duration
}

// New creates a new Checker with the given timeout.
func New(timeout time.Duration) *Checker {
	return &Checker{
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Check performs an HTTP GET against the given URL and returns a Result.
func (c *Checker) Check(ctx context.Context, endpoint string) Result {
	start := time.Now()
	result := Result{
		Endpoint:  endpoint,
		CheckedAt: start,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		result.Status = StatusDegraded
		result.Error = fmt.Sprintf("failed to build request: %v", err)
		return result
	}

	resp, err := c.client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = StatusDegraded
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = StatusHealthy
	} else {
		result.Status = StatusDegraded
		result.Error = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
	}

	return result
}
