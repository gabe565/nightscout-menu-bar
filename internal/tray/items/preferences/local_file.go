package preferences

import (
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/spf13/viper"
)

func NewLocalFile(parent *systray.MenuItem) LocalFile {
	var item LocalFile
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Write to local file",
		"",
		viper.GetBool(local_file.EnabledKey),
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

	viper.Set(local_file.EnabledKey, l.Checked())
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
