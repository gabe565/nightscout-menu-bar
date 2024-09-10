package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/gear-solid.svg
	gear []byte

	//nolint:gochecknoglobals
	GearResource = theme.NewThemedResource(fyne.NewStaticResource("gear.svg", gear))
)
