package items

import (
	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
)

func NewHistory() *fyne.MenuItem {
	item := fyne.NewMenuItem("History", nil)
	item.Icon = assets.HistoryResource
	item.ChildMenu = fyne.NewMenu("History")
	return item
}
