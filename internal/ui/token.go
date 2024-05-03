package ui

import (
	"github.com/ncruces/zenity"
)

func PromptToken(text string) (string, error) {
	return zenity.Entry(
		"Enter new Nightscout API token:",
		zenity.Title("Token"),
		zenity.EntryText(text),
	)
}
