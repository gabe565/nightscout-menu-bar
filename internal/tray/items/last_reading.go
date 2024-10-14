package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewLastReading() *systray.MenuItem {
	item := systray.AddMenuItem("Last Reading", "")
	item.SetTemplateIcon(assets.Droplet, assets.Droplet)
	return item
}
