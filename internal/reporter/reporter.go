package reporter

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Status represents the health status of a service.
type Status struct {
	Name      string
	URL       string
	Healthy   bool
	CheckedAt time.Time
	Error     string
}

// Reporter writes health status summaries to an output writer.
type Reporter struct {
	out io.Writer
}

// New creates a new Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// Report writes a formatted summary of the given statuses.
func (r *Reporter) Report(statuses []Status) {
	fmt.Fprintf(r.out, "=== Health Report [%s] ===\n", time.Now().UTC().Format(time.RFC3339))
	for _, s := range statuses {
		state := "OK"
		if !s.Healthy {
			state = "DEGRADED"
		}
		line := fmt.Sprintf("  [%s] %s (%s)", state, s.Name, s.URL)
		if !s.Healthy && s.Error != "" {
			line += fmt.Sprintf(" — %s", s.Error)
		}
		fmt.Fprintln(r.out, line)
	}
	fmt.Fprintln(r.out, "==============================")
}

// Summary returns counts of healthy and degraded services.
func Summary(statuses []Status) (healthy, degraded int) {
	for _, s := range statuses {
		if s.Healthy {
			healthy++
		} else {
			degraded++
		}
	}
	return
}
