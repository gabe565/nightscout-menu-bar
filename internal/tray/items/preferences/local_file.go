package preferences

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
)

func NewLocalFile(parent *systray.MenuItem) LocalFile {
	var item LocalFile
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Write to local file",
		"",
		config.Default.LocalFile.Enabled,
	)
	return item
}

type LocalFile struct {
	*systray.MenuItem
}

func (l LocalFile) Toggle() error {
	if l.Checked() {
		l.Uncheck()
	} else {
		l.Check()
	}

	config.Default.LocalFile.Enabled = l.Checked()
	if err := config.Write(); err != nil {
		return err
	}
	return nil
}
