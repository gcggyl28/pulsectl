package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// JSONReport represents the JSON-serialisable form of a health report.
type JSONReport struct {
	GeneratedAt string         `json:"generated_at"`
	Healthy     int            `json:"healthy"`
	Degraded    int            `json:"degraded"`
	Services    []JSONService  `json:"services"`
}

// JSONService is the JSON-serialisable form of a single service status.
type JSONService struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Healthy   bool   `json:"healthy"`
	CheckedAt string `json:"checked_at"`
	Error     string `json:"error,omitempty"`
}

// WriteJSON encodes the statuses as a JSON report to w.
func WriteJSON(w io.Writer, statuses []Status) error {
	h, d := Summary(statuses)
	services := make([]JSONService, len(statuses))
	for i, s := range statuses {
		services[i] = JSONService{
			Name:      s.Name,
			URL:       s.URL,
			Healthy:   s.Healthy,
			CheckedAt: s.CheckedAt.UTC().Format(time.RFC3339),
			Error:     s.Error,
		}
	}
	report := JSONReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Healthy:     h,
		Degraded:    d,
		Services:    services,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("reporter: json encode: %w", err)
	}
	return nil
}
