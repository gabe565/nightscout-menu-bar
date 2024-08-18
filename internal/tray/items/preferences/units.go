package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
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
	unit, err := ui.PromptUnits(n.config.Units)
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
