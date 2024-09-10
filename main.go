package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/nightscout-menu-bar/internal/app"
	"github.com/gabe565/nightscout-menu-bar/internal/log"
	"github.com/gabe565/nightscout-menu-bar/internal/pprof"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

var version string

func main() {
	log.Init(os.Stderr, slog.LevelInfo, log.FormatAuto)

	if version == "" {
		version = "beta"
	}
	slog.Info("Nightscout Menu Bar", "version", version, "commit", util.GetCommit())

	pprof.ListenAndServe()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := app.Run(ctx, version); err != nil {
		slog.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
