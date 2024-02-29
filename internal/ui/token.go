package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func PromptToken() (string, error) {
	return zenity.Entry(
		"Enter new Nightscout API token:",
		zenity.Title("Token"),
		zenity.EntryText(config.Default.Token),
	)
}
