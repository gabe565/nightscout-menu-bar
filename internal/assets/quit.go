package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/xmark-solid.svg
	quit []byte

	//nolint:gochecknoglobals
	QuitResource = theme.NewThemedResource(fyne.NewStaticResource("quit.svg", quit))
)
