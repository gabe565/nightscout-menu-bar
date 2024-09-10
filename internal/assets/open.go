package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/square-up-right-solid.svg
	open []byte

	//nolint:gochecknoglobals
	OpenResource = theme.NewThemedResource(fyne.NewStaticResource("open.svg", open))
)
