package preferences

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
)

func NewDynamicIcon(conf *config.Config, parent *systray.MenuItem) DynamicIcon {
	item := DynamicIcon{config: conf}
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Enabled",
		"",
		conf.DynamicIcon.Enabled,
	)
	return item
}

type DynamicIcon struct {
	config *config.Config
	*systray.MenuItem
}

func (l DynamicIcon) Toggle() error {
	if l.Checked() {
		l.Uncheck()
	} else {
		l.Check()
	}

	l.config.DynamicIcon.Enabled = l.Checked()
	if err := l.config.Write(); err != nil {
		return err
	}
	return nil
}
