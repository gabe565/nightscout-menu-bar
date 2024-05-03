package preferences

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
)

func NewLocalFile(conf *config.Config, parent *systray.MenuItem) LocalFile {
	item := LocalFile{config: conf}
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Write to local file",
		"",
		conf.LocalFile.Enabled,
	)
	return item
}

type LocalFile struct {
	config *config.Config
	*systray.MenuItem
}

func (l LocalFile) Toggle() error {
	if l.Checked() {
		l.Uncheck()
	} else {
		l.Check()
	}

	l.config.LocalFile.Enabled = l.Checked()
	if err := l.config.Write(); err != nil {
		return err
	}
	return nil
}
