package ui

import (
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func PromptURL() (string, error) {
	url, err := zenity.Entry(
		"Enter new Nightscout URL:",
		zenity.Title("Nightscout URL"),
		zenity.EntryText(viper.GetString("url")),
	)
	if err != nil && err.Error() != "dialog canceled" {
		return url, err
	}
	return url, nil
}
