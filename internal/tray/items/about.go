package items

import (
	"log/slog"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewAbout(app fyne.App) *fyne.MenuItem {
	title := "Nightscout Menu Bar"
	version := app.Metadata().Version
	if version != "" {
		if strings.HasPrefix(version, "v") {
			title += " " + version
		} else {
			title += " (" + version + ")"
		}
	}
	item := fyne.NewMenuItem(title, func() {
		u := &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "/gabe565/nightscout-menu-bar",
		}
		slog.Debug("Opening GitHub repo", "url", u)
		if err := app.OpenURL(u); err != nil {
			slog.Error("Failed to open GitHub repo", "error", err)
		}
	})
	item.Icon = assets.AboutResource
	return item
}
