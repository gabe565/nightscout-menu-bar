package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/rectangle-history-solid.svg
	history []byte

	//nolint:gochecknoglobals
	HistoryResource = theme.NewThemedResource(fyne.NewStaticResource("history.svg", history))
)
