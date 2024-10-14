package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

type History struct {
	*systray.MenuItem
	Subitems []*systray.MenuItem
}

func NewHistory() History {
	item := systray.AddMenuItem("History", "")
	item.SetTemplateIcon(assets.History, assets.History)
	vals := make([]*systray.MenuItem, 0, 4)
	return History{item, vals}
}
