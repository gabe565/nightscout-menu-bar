package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/getlantern/systray"
)

func NewOpenNightscout() *systray.MenuItem {
	item := systray.AddMenuItem("Open Nightscout", "")
	item.SetTemplateIcon(assets.SquareUpRight, assets.SquareUpRight)
	return item
}
