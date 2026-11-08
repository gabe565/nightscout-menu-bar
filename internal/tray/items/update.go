package items

import (
	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/update"
)

func NewUpdate() *Update {
	item := systray.AddMenuItem("", "")
	item.SetTemplateIcon(assets.Update, assets.Update)
	item.Hide()
	return &Update{MenuItem: item}
}

type Update struct {
	*systray.MenuItem
	URL     string
	visible bool
}

func (u *Update) Show(rel *update.Release) {
	title := "Update available: " + rel.Version
	switch update.InstallMethod() {
	case update.MethodHomebrew:
		title += " (run brew upgrade)"
	case update.MethodPackage:
		title += " (use your package manager)"
	case update.MethodManual:
	}
	u.SetTitle(title)
	u.URL = rel.URL
	u.visible = true
	u.MenuItem.Show()
}

func (u *Update) Hide() {
	u.visible = false
	u.MenuItem.Hide()
}
