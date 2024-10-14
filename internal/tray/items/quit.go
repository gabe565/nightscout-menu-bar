package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewQuit() *systray.MenuItem {
	item := systray.AddMenuItem("Quit", "")
	item.SetTemplateIcon(assets.Quit, assets.Quit)
	return item
}
