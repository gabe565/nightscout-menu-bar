package main

import (
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

func main() {
	if err := InitViper(); err != nil {
		go func() {
			tray.Error <- err
		}()
	}

	local_file.ReloadConfig()

	ticker.BeginRender()
	ticker.BeginFetch()
	tray.Run()
}
