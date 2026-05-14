package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Endpoint defines a single service to monitor.
type Endpoint struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
}

// Config holds the full pulsectl configuration.
type Config struct {
	WebhookURL string        `json:"webhook_url"`
	Timeout    time.Duration `json:"timeout"`
	Endpoints  []Endpoint    `json:"endpoints"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	// Apply defaults.
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	for i := range cfg.Endpoints {
		if cfg.Endpoints[i].Interval == 0 {
			cfg.Endpoints[i].Interval = 30 * time.Second
		}
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.WebhookURL == "" {
		return fmt.Errorf("webhook_url is required")
	}
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be defined")
	}
	for i, ep := range c.Endpoints {
		if ep.URL == "" {
			return fmt.Errorf("endpoint[%d]: url is required", i)
		}
		if ep.Name == "" {
			return fmt.Errorf("endpoint[%d]: name is required", i)
		}
	}
	return nil
}
