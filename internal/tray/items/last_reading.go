package items

import (
	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewLastReading() *fyne.MenuItem {
	item := fyne.NewMenuItem("Last Reading", nil)
	item.Icon = assets.DropletResource
	return item
}
