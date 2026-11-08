package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewSettings() *systray.MenuItem {
	item := systray.AddMenuItem("Settings…", "")
	item.SetTemplateIcon(assets.Preferences, assets.Preferences)
	return item
}
