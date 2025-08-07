package preferences

import (
	"errors"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func NewDynamicIconColor(conf *config.Config, parent *systray.MenuItem) DynamicIconColor {
	item := DynamicIconColor{config: conf}
	item.MenuItem = parent.AddSubMenuItem("Color", "")
	return item
}

type DynamicIconColor struct {
	config *config.Config
	*systray.MenuItem
}

func (l DynamicIconColor) Choose() error {
	c, err := zenity.SelectColor(
		zenity.Title("Dynamic Icon Color"),
		zenity.Color(l.config.Data().DynamicIcon.FontColor),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	data := l.config.Data()
	data.DynamicIcon.FontColor.Color = c
	if err := l.config.Write(data); err != nil {
		return err
	}
	return nil
}
