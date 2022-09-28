package main

import (
	"github.com/gabe565/nightscout-systray/internal/config"
	"github.com/gabe565/nightscout-systray/internal/tray"
)

func main() {
	if err := config.InitViper(); err != nil {
		go func() {
			tray.Error <- err
		}()
	}

	tray.Run()
}
