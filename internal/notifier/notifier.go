package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the webhook notification body.
type Payload struct {
	Service   string `json:"service"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

// Notifier sends degraded-service alerts to a webhook endpoint.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a Notifier with the given webhook URL.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify posts a degraded-service payload to the configured webhook.
func (n *Notifier) Notify(service, url, reason string) error {
	p := Payload{
		Service:   service,
		URL:       url,
		Status:    "degraded",
		Reason:    reason,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("notifier: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}
