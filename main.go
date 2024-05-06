package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	tray.New().Run(ctx)
}
