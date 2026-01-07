package items

import (
	"strings"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
)

func NewAbout(version string) *systray.MenuItem {
	var title strings.Builder
	title.WriteString("Nightscout Menu Bar")
	if version != "" {
		if strings.HasPrefix(version, "v") {
			title.WriteRune(' ')
			title.WriteString(version)
		} else {
			title.WriteString(" (")
			title.WriteString(version)
			title.WriteRune(')')
		}
	}
	item := systray.AddMenuItem(title.String(), "")
	item.SetTemplateIcon(assets.About, assets.About)
	return item
}
