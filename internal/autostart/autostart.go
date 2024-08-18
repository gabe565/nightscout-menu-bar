package autostart

import (
	"log/slog"
	"os"

	"github.com/emersion/go-autostart"
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
	slog.Debug("Enabling autostart")
	return app.Enable()
}

func Disable() error {
	app, err := NewApp()
	if err != nil {
		return err
	}
	slog.Debug("Disabling autostart")
	return app.Disable()
}

func IsEnabled() (bool, error) {
	app, err := NewApp()
	if err != nil {
		return false, err
	}
	v := app.IsEnabled()
	slog.Debug("Detected autostart status", "value", v)
	return v, nil
}
