package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/pprof"
	"gabe565.com/nightscout-menu-bar/internal/tray"
	"gabe565.com/nightscout-menu-bar/internal/util"
)

var version string

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)

	if version == "" {
		version = "beta"
	}
	slog.Info("Nightscout Menu Bar", "version", version, "commit", util.GetCommit())

	pprof.ListenAndServe()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	tray.New(version).Run(ctx)
}
