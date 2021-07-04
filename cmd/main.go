package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	pdmanager "github.com/corsc/pagerduty-manager"
)

const (
	maxExecutionTime = 60 * time.Second
)

func main() {
	cfg := buildConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to init the logger with err: %s", err)
		os.Exit(-1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxExecutionTime)
	defer cancel()

	manager := pdmanager.New(cfg, logger)

	err = manager.Parse(ctx)
	if err != nil {
		logger.Fatal("failed to parse", zap.Error(err))
		return
	}

	err = manager.SyncUsers(ctx)
	if err != nil {
		logger.Fatal("failed to sync users", zap.Error(err))
		return
	}

	err = manager.SyncTeams(ctx)
	if err != nil {
		logger.Fatal("failed to sync teams", zap.Error(err))
		return
	}

	err = manager.SyncServices(ctx)
	if err != nil {
		logger.Fatal("failed to sync services", zap.Error(err))
		return
	}

	err = manager.SyncEscalation(ctx)
	if err != nil {
		logger.Fatal("failed to sync escalation policies", zap.Error(err))
		return
	}

	err = manager.SyncSchedules(ctx)
	if err != nil {
		logger.Fatal("failed to sync on-call schedules", zap.Error(err))
		return
	}
}

func buildConfig() *config {
	cfg := &config{
		accessToken: os.Getenv("PD_TOKEN"),
	}

	flag.BoolVar(&cfg.debug, "debug", false, "enable debug mode")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Please supply a JSON file")
		os.Exit(-1)
	}

	cfg.filename = args[0]

	return cfg
}

type config struct {
	accessToken string
	filename    string
	debug       bool
}

func (c *config) Debug() bool {
	return c.debug
}

func (c *config) Filename() string {
	return c.filename
}

func (c *config) AccessToken() string {
	return c.accessToken
}
