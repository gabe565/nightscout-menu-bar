package items

import (
	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewSettings(app fyne.App) *fyne.MenuItem {
	item := fyne.NewMenuItem("Settings", settings.OpenSettings(app))
	item.Icon = assets.GearResource
	return item
}
