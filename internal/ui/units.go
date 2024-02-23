package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func PromptUnits() (string, error) {
	return zenity.List(
		"Select units:",
		[]string{nightscout.UnitsMgdl, nightscout.UnitsMmol},
		zenity.Title("Nightscout Units"),
		zenity.DisallowEmpty(),
		zenity.DefaultItems(viper.GetString(nightscout.UnitsKey)),
	)
}
