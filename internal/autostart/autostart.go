package autostart

import (
	"os"

	"github.com/emersion/go-autostart"
	"github.com/rs/zerolog/log"
)

func NewApp() (autostart.App, error) {
	executable, err := os.Executable()
	if err != nil {
		return autostart.App{}, err
	}

	return autostart.App{
		Name:        "com.gabe565.nightscout-menu-bar",
		DisplayName: "Nightscout Menu Bar",
		Exec:        []string{executable},
	}, nil
}

func Enable() error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	log.Debug().Msg("Enabling autostart")
	return app.Enable()
}

func Disable() error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	log.Debug().Msg("Disabling autostart")
	return app.Disable()
}

func IsEnabled() (bool, error) {
	app, err := NewApp()
	if err != nil {
		return false, err
	}
	v := app.IsEnabled()
	log.Trace().Bool("value", v).Msg("Detected autostart status")
	return v, nil
}
