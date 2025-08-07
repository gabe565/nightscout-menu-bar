package preferences

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
)

func NewDynamicIcon(conf *config.Config, parent *systray.MenuItem) DynamicIcon {
	item := DynamicIcon{config: conf}
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Enabled",
		"",
		conf.Data().DynamicIcon.Enabled,
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

	data := l.config.Data()
	data.DynamicIcon.Enabled = l.Checked()
	if err := l.config.Write(data); err != nil {
		return err
	}
	return nil
}
