package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func PromptUnits(item config.Unit) (string, error) {
	return zenity.List(
		"Select units:",
		config.UnitStrings(),
		zenity.Title("Nightscout Units"),
		zenity.DisallowEmpty(),
		zenity.DefaultItems(item.String()),
	)
}
