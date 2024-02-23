package ui

import (
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func PromptToken() (string, error) {
	return zenity.Entry(
		"Enter new Nightscout API token:",
		zenity.Title("Token"),
		zenity.EntryText(viper.GetString(nightscout.TokenKey)),
	)
}
