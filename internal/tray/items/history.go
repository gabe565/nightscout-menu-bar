package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/getlantern/systray"
)

func NewHistory() (*systray.MenuItem, []*systray.MenuItem) {
	item := systray.AddMenuItem("History", "")
	item.SetTemplateIcon(assets.RectangleHistory, assets.RectangleHistory)
	vals := make([]*systray.MenuItem, 0, 4)
	return item, vals
}
