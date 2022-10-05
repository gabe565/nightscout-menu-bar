package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/getlantern/systray"
)

func NewError() *systray.MenuItem {
	item := systray.AddMenuItem("", "")
	item.SetTemplateIcon(assets.Error, assets.Error)
	item.Disable()
	item.Hide()
	return item
}
