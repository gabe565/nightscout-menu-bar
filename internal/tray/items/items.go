package items

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items/preferences"
)

type Items struct {
	LastReading    *systray.MenuItem
	Error          *systray.MenuItem
	OpenNightscout *systray.MenuItem
	History        History
	Preferences    preferences.Preferences
	Quit           *systray.MenuItem
}

func New(conf *config.Config) Items {
	var items Items

	items.LastReading = NewLastReading()
	items.Error = NewError()
	systray.AddSeparator()

	items.OpenNightscout = NewOpenNightscout(conf.Title)
	items.History = NewHistory()
	systray.AddSeparator()

	items.Preferences = preferences.New(conf)
	items.Quit = NewQuit()

	return items
}
