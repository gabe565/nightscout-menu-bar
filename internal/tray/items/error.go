package items

import "github.com/getlantern/systray"

func NewError() *systray.MenuItem {
	item := systray.AddMenuItem("", "")
	item.Hide()
	return item
}
