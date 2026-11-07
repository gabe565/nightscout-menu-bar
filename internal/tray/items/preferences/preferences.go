package preferences

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/autostart"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
)

func New(conf *config.Config, store *oauth.Store) Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Preferences, assets.Preferences)

	url := NewURL(conf, item)
	authMenu := item.AddSubMenuItem("Authentication", "")
	token := NewToken(conf, authMenu)
	nocturne := NewNocturne(conf, store, authMenu)
	units := NewUnits(conf, item)
	item.AddSeparator()

	dynamicIconMenu := item.AddSubMenuItem("Dynamic icon", "")
	dynamicIcon := NewDynamicIcon(conf, dynamicIconMenu)
	dynamicIconColor := NewDynamicIconColor(conf, dynamicIconMenu)
	item.AddSeparator()

	autostartEnabled, _ := autostart.IsEnabled()
	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)
	item.AddSeparator()

	socket := NewSocket(conf, item)

	return Preferences{
		MenuItem:         item,
		URL:              url,
		Token:            token,
		Nocturne:         nocturne,
		Units:            units,
		DynamicIcon:      dynamicIcon,
		DynamicIconColor: dynamicIconColor,
		StartOnLogin:     startOnLogin,
		Socket:           socket,
	}
}

type Preferences struct {
	*systray.MenuItem
	URL              URL
	Token            Token
	Nocturne         *Nocturne
	Units            Units
	DynamicIcon      DynamicIcon
	DynamicIconColor DynamicIconColor
	StartOnLogin     *systray.MenuItem
	Socket           Socket
}

type Item interface {
	MenuItem() *systray.MenuItem
	GetTitle() string
	UpdateTitle()
	Prompt() error
}
