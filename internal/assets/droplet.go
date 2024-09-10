package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed src/droplet-solid.svg
	droplet []byte

	//nolint:gochecknoglobals
	DropletResource = theme.NewThemedResource(fyne.NewStaticResource("droplet.svg", droplet))
)
