package items

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/gabe565/nightscout-systray/internal/autostart"
	"github.com/getlantern/systray"
)

func NewPreferences() Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostart.IsEnabled(),
	)

	return Preferences{
		Item:         item,
		StartOnLogin: startOnLogin,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	StartOnLogin *systray.MenuItem
}
