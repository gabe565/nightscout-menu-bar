package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewError() *Error {
	item := systray.AddMenuItem("", "")
	item.SetTemplateIcon(assets.Error, assets.Error)
	item.Disable()
	item.Hide()
	return &Error{MenuItem: item}
}

type Error struct {
	*systray.MenuItem
	visible bool
}

func (e *Error) Show() {
	e.visible = true
	e.MenuItem.Show()
}

func (e *Error) Hide() {
	e.visible = false
	e.MenuItem.Hide()
}
