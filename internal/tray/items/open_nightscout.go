package items

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewOpenNightscout() *systray.MenuItem {
	item := systray.AddMenuItem("Open Nightscout", "")
	item.SetTemplateIcon(assets.Open, assets.Open)
	return item
}
