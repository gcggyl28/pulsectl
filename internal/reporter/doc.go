// Package reporter provides utilities for summarising and formatting
// health-check results collected by the pulsectl poller.
//
// It supports two output formats:
//
//   - Human-readable plain-text via Reporter.Report
//   - Machine-readable JSON via WriteJSON
//
// Example usage:
//
//	statuses := []reporter.Status{
//		{Name: "api", URL: "http://api.example.com", Healthy: true},
//		{Name: "db",  URL: "http://db.example.com",  Healthy: false, Error: "timeout"},
//	}
//
//	// Plain text
//	r := reporter.New(os.Stdout)
//	r.Report(statuses)
//
//	// JSON
//	reporter.WriteJSON(os.Stdout, statuses)
package reporter
