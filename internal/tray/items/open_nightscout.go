package items

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
)

func NewOpenNightscout(app fyne.App) *fyne.MenuItem {
	item := fyne.NewMenuItem("Open Nightscout", func() {
		u, err := fetch.BuildURLWithToken(app.Preferences())
		if err != nil {
			slog.Error("Failed to build URL", "error", err)
			return
		}
		slog.Debug("Opening Nightscout", "url", u)
		if err := app.OpenURL(u); err != nil {
			slog.Error("Failed to open Nightscout", "error", err)
		}
	})
	item.Icon = assets.OpenResource
	return item
}
