package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/github-brands-solid.svg
	about []byte

	//nolint:gochecknoglobals
	AboutResource = theme.NewThemedResource(fyne.NewStaticResource("about.svg", about))
)
