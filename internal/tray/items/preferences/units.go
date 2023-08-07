package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

type Units struct {
	*systray.MenuItem
}

func (n Units) UpdateTitle() {
	title := "Units: " + viper.GetString(config.UnitsKey)
	n.SetTitle(title)
}

func (n Units) Prompt() error {
	unit, err := ui.PromptUnits()
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	viper.Set(config.UnitsKey, unit)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
