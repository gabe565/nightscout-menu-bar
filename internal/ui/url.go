package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func PromptURL() (string, error) {
	return zenity.Entry(
		"Enter new Nightscout URL:",
		zenity.Title("Nightscout URL"),
		zenity.EntryText(config.Default.URL),
	)
}
