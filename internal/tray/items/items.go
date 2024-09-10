package items

import (
	"fyne.io/fyne/v2"
)

type Items struct {
	LastReading    *fyne.MenuItem
	OpenNightscout *fyne.MenuItem
	History        *fyne.MenuItem
	Settings       *fyne.MenuItem
	About          *fyne.MenuItem
	Quit           *fyne.MenuItem
}

func New(app fyne.App) Items {
	return Items{
		LastReading:    NewLastReading(),
		OpenNightscout: NewOpenNightscout(app),
		History:        NewHistory(),
		Settings:       NewSettings(app),
		About:          NewAbout(app),
		Quit:           NewQuit(),
	}
}
