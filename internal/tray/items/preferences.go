package items

import (
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items/preferences"
	"github.com/getlantern/systray"
	"github.com/spf13/viper"
)

func NewPreferences() Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	urlTitle := "Nightscout URL"
	if url := viper.GetString("url"); url != "" {
		urlTitle += ": " + url
	}
	nightscoutUrl := item.AddSubMenuItem(urlTitle, "")

	unitTitle := "Units: " + viper.GetString(config.UnitsKey)
	unit := item.AddSubMenuItem(unitTitle, "")

	autostartEnabled, _ := autostart.IsEnabled()

	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)

	return Preferences{
		Item:         item,
		Url:          preferences.Url{MenuItem: nightscoutUrl},
		Units:        preferences.Units{MenuItem: unit},
		StartOnLogin: startOnLogin,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	Url          preferences.Url
	Units        preferences.Units
	StartOnLogin *systray.MenuItem
}
