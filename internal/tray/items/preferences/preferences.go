package preferences

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/autostart"
	"gabe565.com/nightscout-menu-bar/internal/config"
)

func New(conf *config.Config) Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	url := NewURL(conf, item)
	token := NewToken(conf, item)
	units := NewUnits(conf, item)

	autostartEnabled, _ := autostart.IsEnabled()
	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostartEnabled,
	)

	dynamicIconMenu := item.AddSubMenuItem("Dynamic icon", "")
	dynamicIcon := NewDynamicIcon(conf, dynamicIconMenu)
	dynamicIconColor := NewDynamicIconColor(conf, dynamicIconMenu)
	localFile := NewLocalFile(conf, item)

	return Preferences{
		MenuItem:         item,
		URL:              url,
		Token:            token,
		Units:            units,
		StartOnLogin:     startOnLogin,
		DynamicIcon:      dynamicIcon,
		DynamicIconColor: dynamicIconColor,
		LocalFile:        localFile,
	}
}

type Preferences struct {
	*systray.MenuItem
	URL              URL
	Token            Token
	Units            Units
	StartOnLogin     *systray.MenuItem
	DynamicIcon      DynamicIcon
	DynamicIconColor DynamicIconColor
	LocalFile        LocalFile
}

type Item interface {
	MenuItem() *systray.MenuItem
	GetTitle() string
	UpdateTitle()
	Prompt() error
}
