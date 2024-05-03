package items

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items/preferences"
)

func NewPreferences(conf *config.Config) Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	url := preferences.NewUrl(conf, item)
	token := preferences.NewToken(conf, item)
	units := preferences.NewUnits(conf, item)

	autostartEnabled, _ := autostart.IsEnabled()
	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)

	localFile := preferences.NewLocalFile(conf, item)

	return Preferences{
		Item:         item,
		Url:          url,
		Token:        token,
		Units:        units,
		StartOnLogin: startOnLogin,
		LocalFile:    localFile,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	Url          preferences.Url
	Token        preferences.Token
	Units        preferences.Units
	StartOnLogin *systray.MenuItem
	LocalFile    preferences.LocalFile
}
