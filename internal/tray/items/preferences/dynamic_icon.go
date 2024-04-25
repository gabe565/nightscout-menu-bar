package preferences

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
)

func NewDynamicIcon(parent *systray.MenuItem) DynamicIcon {
	var item DynamicIcon
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Render dynamic icon",
		"",
		config.Default.DynamicIcon.Enabled,
	)
	return item
}

type DynamicIcon struct {
	*systray.MenuItem
}

func (l DynamicIcon) Toggle() error {
	if l.Checked() {
		l.Uncheck()
	} else {
		l.Check()
	}

	config.Default.DynamicIcon.Enabled = l.Checked()
	config.Default.HideTitle = l.Checked()
	if err := config.Write(); err != nil {
		return err
	}
	return nil
}
