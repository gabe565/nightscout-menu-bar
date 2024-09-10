package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/nightscout-transparent.svg
	nightscout []byte

	//nolint:gochecknoglobals
	NightscoutResource = theme.NewThemedResource(fyne.NewStaticResource("nightscout.svg", nightscout))
)
