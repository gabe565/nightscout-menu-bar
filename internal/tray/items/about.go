package items

import (
	"strings"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewAbout(version string) *systray.MenuItem {
	title := "Nightscout Menu Bar"
	if version != "" {
		if strings.HasPrefix(version, "v") {
			title += " " + version
		} else {
			title += " (" + version + ")"
		}
	}
	item := systray.AddMenuItem(title, "")
	item.SetTemplateIcon(assets.About, assets.About)
	return item
}
