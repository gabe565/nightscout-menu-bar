package items

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/getlantern/systray"
)

func NewQuit() *systray.MenuItem {
	item := systray.AddMenuItem("Quit Nightscout Systray", "")
	item.SetTemplateIcon(assets.Xmark, assets.Xmark)
	return item
}
