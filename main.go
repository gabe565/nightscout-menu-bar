package main

import "github.com/getlantern/systray"

func main() {
	systray.Run(onReady, onExit)
}
