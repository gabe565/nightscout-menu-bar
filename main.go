package main

import (
	"github.com/getlantern/systray"
	flag "github.com/spf13/pflag"
	"log"
)

func main() {
	flag.Parse()
	if err := initViper(); err != nil {
		log.Fatalln(err)
	}

	systray.Run(onReady, onExit)
}
