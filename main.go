package main

import (
	"github.com/gabe565/nightscout-systray/internal/config"
	"github.com/gabe565/nightscout-systray/internal/tray"
	"log"
)

func main() {
	if err := config.InitViper(); err != nil {
		log.Fatalln(err)
	}

	tray.Run()
}
