package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewError() *systray.MenuItem {
	item := systray.AddMenuItem("", "")
	item.SetTemplateIcon(assets.Error, assets.Error)
	item.Disable()
	item.Hide()
	return item
}
