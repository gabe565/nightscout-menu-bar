package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
)

func NewUnits(parent *systray.MenuItem) Units {
	var item Units
	item.MenuItem = parent.AddSubMenuItem(item.GetTitle(), "")
	return item
}

type Units struct {
	*systray.MenuItem
}

func (n Units) GetTitle() string {
	return "Units: " + config.Default.Units
}

func (n Units) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n Units) Prompt() error {
	unit, err := ui.PromptUnits()
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	config.Default.Units = unit
	if err := config.Write(); err != nil {
		return err
	}
	return nil
}
