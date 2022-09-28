package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/getlantern/systray"
)

func NewLastReading() *systray.MenuItem {
	item := systray.AddMenuItem("Last Reading", "")
	item.SetTemplateIcon(assets.Calendar, assets.Calendar)
	val := item.AddSubMenuItem("", "")
	val.Disable()
	return val
}
