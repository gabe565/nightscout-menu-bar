package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewOpenNightscout(title string) *systray.MenuItem {
	item := systray.AddMenuItem("Open "+title, "")
	item.SetTemplateIcon(assets.Open, assets.Open)
	return item
}
