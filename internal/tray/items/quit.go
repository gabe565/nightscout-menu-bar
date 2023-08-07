package items

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewQuit() *systray.MenuItem {
	item := systray.AddMenuItem("Quit Nightscout Systray", "")
	item.SetTemplateIcon(assets.Quit, assets.Quit)
	return item
}
