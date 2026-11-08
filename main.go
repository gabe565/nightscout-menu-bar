package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/pprof"
	"gabe565.com/nightscout-menu-bar/internal/settings"
	"gabe565.com/nightscout-menu-bar/internal/tray"
	"gabe565.com/nightscout-menu-bar/internal/util"
	"gabe565.com/utils/slogx"
)

var version string

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)

	if version == "" {
		version = "beta"
	}
	slog.Info("Nightscout Menu Bar", "version", version, "commit", util.GetCommit())

	if len(os.Args) > 1 && os.Args[1] == "settings" {
		if err := settings.Run(version, os.Args[2:]); err != nil {
			slog.Error("Settings failed", "error", err)
			os.Exit(1)
		}
		return
	}

	pprof.ListenAndServe()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	tray.New(version).Run(ctx)
}
