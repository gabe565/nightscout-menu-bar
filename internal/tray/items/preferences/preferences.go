package preferences

import "fyne.io/systray"

type Item interface {
	MenuItem() *systray.MenuItem
	GetTitle() string
	UpdateTitle()
	Prompt() error
}
