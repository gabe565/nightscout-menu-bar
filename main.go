package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

func main() {
	if err := config.Load(); err != nil {
		go func() {
			tray.Error <- err
		}()
	}
	if err := config.Watch(); err != nil {
		go func() {
			tray.Error <- err
		}()
	}

	local_file.ReloadConfig()

	ticker.BeginRender()
	ticker.BeginFetch()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	go func() {
		<-ctx.Done()
		systray.Quit()
	}()
	tray.Run()
}
