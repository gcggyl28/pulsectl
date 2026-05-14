package poller

import (
	"context"
	"log"
	"time"

	"github.com/user/pulsectl/internal/checker"
	"github.com/user/pulsectl/internal/config"
	"github.com/user/pulsectl/internal/notifier"
)

// Poller periodically checks all configured services and notifies on degradation.
type Poller struct {
	cfg      *config.Config
	checker  *checker.Checker
	notifier *notifier.Notifier
}

// New creates a Poller wired to the provided config.
func New(cfg *config.Config) *Poller {
	return &Poller{
		cfg:      cfg,
		checker:  checker.New(cfg.Timeout),
		notifier: notifier.New(cfg.WebhookURL),
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	log.Printf("poller: starting — interval=%s services=%d", p.cfg.Interval, len(p.cfg.Services))

	for {
		select {
		case <-ctx.Done():
			log.Println("poller: shutting down")
			return
		case <-ticker.C:
			p.runChecks()
		}
	}
}

func (p *Poller) runChecks() {
	for _, svc := range p.cfg.Services {
		ok, reason := p.checker.Check(svc.URL)
		if ok {
			log.Printf("poller: [OK] %s", svc.Name)
			continue
		}

		log.Printf("poller: [DEGRADED] %s — %s", svc.Name, reason)

		if err := p.notifier.Notify(svc.Name, svc.URL, reason); err != nil {
			log.Printf("poller: failed to notify for %s: %v", svc.Name, err)
		}
	}
}
