package preferences

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
)

func NewSocket(conf *config.Config, parent *systray.MenuItem) Socket {
	item := Socket{config: conf}
	item.MenuItem = parent.AddSubMenuItemCheckbox(
		"Expose readings over local socket",
		"",
		conf.Socket.Enabled,
	)
	return item
}

type Socket struct {
	config *config.Config
	*systray.MenuItem
}

func (s Socket) Toggle() error {
	if s.Checked() {
		s.Uncheck()
	} else {
		s.Check()
	}

	s.config.Socket.Enabled = s.Checked()
	if err := s.config.Write(); err != nil {
		return err
	}
	return nil
}
