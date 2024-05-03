package ui

import (
	"github.com/ncruces/zenity"
)

func PromptURL(text string) (string, error) {
	return zenity.Entry(
		"Enter new Nightscout URL:",
		zenity.Title("Nightscout URL"),
		zenity.EntryText(text),
	)
}
