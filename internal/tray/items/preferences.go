package items

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items/preferences"
)

func NewPreferences() Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	autostartEnabled, _ := autostart.IsEnabled()
	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)

	return Preferences{
		Item:         item,
		Url:          preferences.NewUrl(item),
		Token:        preferences.NewToken(item),
		Units:        preferences.NewUnits(item),
		StartOnLogin: startOnLogin,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	Url          preferences.Url
	Token        preferences.Token
	Units        preferences.Units
	StartOnLogin *systray.MenuItem
}
