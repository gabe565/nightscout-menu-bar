package app

import (
	"context"
	"errors"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/fyneutil"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

var ErrUnsupported = errors.New("only desktop is supported")

func Run(ctx context.Context, version string) error {
	fyneapp.SetMetadata(fyne.AppMetadata{
		ID:      "com.gabe565.nightscout-menu-bar",
		Name:    "Nightscout Menu Bar",
		Version: version,
	})

	app, ok := fyneapp.NewWithID("com.gabe565.nightscout-menu-bar").(fyneutil.DesktopApp)
	if !ok {
		return ErrUnsupported
	}

	app.SetSystemTrayIcon(assets.NightscoutResource)
	t := tray.New(app, version)
	app.SetSystemTrayMenu(t.Menu())
	app.Lifecycle().SetOnStarted(func() {
		HideAppIcon()
		t.Run(ctx)
	})

	app.Run()
	return nil
}
