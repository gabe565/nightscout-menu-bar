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

	url := preferences.NewUrl(item)
	token := preferences.NewToken(item)
	units := preferences.NewUnits(item)

	autostartEnabled, _ := autostart.IsEnabled()
	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)

	dynamicIcon := preferences.NewDynamicIcon(item)
	localFile := preferences.NewLocalFile(item)

	return Preferences{
		Item:         item,
		Url:          url,
		Token:        token,
		Units:        units,
		StartOnLogin: startOnLogin,
		DynamicIcon:  dynamicIcon,
		LocalFile:    localFile,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	Url          preferences.Url
	Token        preferences.Token
	Units        preferences.Units
	StartOnLogin *systray.MenuItem
	DynamicIcon  preferences.DynamicIcon
	LocalFile    preferences.LocalFile
}
