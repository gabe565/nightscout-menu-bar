package autostart

import (
	"github.com/emersion/go-autostart"
	"log"
	"os"
)

var app *autostart.App
var initError error

func init() {
	executable, err := os.Executable()
	if err != nil {
		initError = err
		log.Println(err)
	}

	app = &autostart.App{
		Name:        "com.gabe565.nightscout-menu-bar",
		DisplayName: "Nightscout Menu Bar",
		Exec:        []string{executable},
	}
}

func Enable() error {
	if app == nil {
		return initError
	}
	return app.Enable()
}

func Disable() error {
	if app == nil {
		return initError
	}
	return app.Disable()
}

func IsEnabled() bool {
	if app == nil {
		return false
	}
	return app.IsEnabled()
}
