package main

import (
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

func main() {
	if err := InitViper(); err != nil {
		go func() {
			tray.Error <- err
		}()
	}

	ticker.BeginRender()
	ticker.BeginFetch()
	tray.Run()
}
