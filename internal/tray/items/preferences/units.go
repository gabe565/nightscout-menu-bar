package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func NewUnits(config *config.Config, parent *systray.MenuItem) Units {
	item := Units{config: config}
	item.MenuItem = parent.AddSubMenuItem(item.GetTitle(), "")
	return item
}

type Units struct {
	config *config.Config
	*systray.MenuItem
}

func (n Units) GetTitle() string {
	return "Units: " + n.config.Units.String()
}

func (n Units) UpdateTitle() {
	n.SetTitle(n.GetTitle())
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
