package items

import (
	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewQuit() *fyne.MenuItem {
	item := fyne.NewMenuItem("Quit", nil)
	item.Icon = assets.QuitResource
	item.IsQuit = true
	return item
}
