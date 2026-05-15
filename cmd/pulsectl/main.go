package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/pulsectl/internal/checker"
	"github.com/example/pulsectl/internal/config"
	"github.com/example/pulsectl/internal/notifier"
	"github.com/example/pulsectl/internal/poller"
	"github.com/example/pulsectl/internal/reporter"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	chk := checker.New(cfg.TimeoutSeconds)
	ntf := notifier.New(cfg.WebhookURL)
	rep := reporter.New(os.Stdout)

	p := poller.New(cfg, chk, ntf, rep)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("pulsectl starting — polling %d service(s) every %ds",
		len(cfg.Services), cfg.IntervalSeconds)

	if err := p.Run(ctx); err != nil {
		log.Fatalf("poller exited with error: %v", err)
	}

	log.Println("pulsectl stopped")
}
