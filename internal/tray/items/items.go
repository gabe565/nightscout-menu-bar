package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
)

type Items struct {
	LastReading    *systray.MenuItem
	Error          *Error
	History        History
	OpenNightscout *systray.MenuItem
	Settings       *systray.MenuItem
	About          *systray.MenuItem
	Quit           *systray.MenuItem
}

func New(conf *config.Config) Items {
	var items Items

	items.LastReading = NewLastReading()
	items.Error = NewError()
	items.History = NewHistory()
	systray.AddSeparator()

	items.OpenNightscout = NewOpenNightscout(conf.Data().Title)
	items.Settings = NewSettings()
	systray.AddSeparator()

	items.About = NewAbout(conf.Version)
	items.Quit = NewQuit()

	return items
}
