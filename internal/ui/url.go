package ui

import (
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func PromptURL() (string, error) {
	return zenity.Entry(
		"Enter new Nightscout URL:",
		zenity.Title("Nightscout URL"),
		zenity.EntryText(viper.GetString("url")),
	)
}
