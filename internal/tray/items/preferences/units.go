package preferences

import (
	"errors"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func NewUnits(config *config.Config, parent *systray.MenuItem) Units {
	item := Units{config: config}
	item.MenuItem = parent.AddSubMenuItem("Units", "")
	return item
}

type Units struct {
	config *config.Config
	*systray.MenuItem
}

func (n Units) Prompt() error {
	unit, err := zenity.List(
		"Select units:",
		config.UnitStrings(),
		zenity.Title("Nightscout Units"),
		zenity.DisallowEmpty(),
		zenity.DefaultItems(n.config.Units.String()),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	if err := n.config.Units.UnmarshalText([]byte(unit)); err != nil {
		return err
	}

	if err := n.config.Write(); err != nil {
		return err
	}
	return nil
}
