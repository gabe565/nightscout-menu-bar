package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func PromptUnits() (string, error) {
	return zenity.List(
		"Select units:",
		[]string{config.UnitsMgdl, config.UnitsMmol},
		zenity.Title("Nightscout Units"),
		zenity.DisallowEmpty(),
		zenity.DefaultItems(config.Default.Units),
	)
}
