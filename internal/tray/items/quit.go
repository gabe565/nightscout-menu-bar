package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/getlantern/systray"
)

func NewQuit() *systray.MenuItem {
	item := systray.AddMenuItem("Quit Nightscout Systray", "")
	item.SetTemplateIcon(assets.Quit, assets.Quit)
	return item
}
