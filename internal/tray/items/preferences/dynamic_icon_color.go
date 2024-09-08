package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
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

//nolint:gosec
func (l DynamicIconColor) Choose() error {
	c, err := zenity.SelectColor(
		zenity.Title("Dynamic Icon Color"),
		zenity.Color(l.config.DynamicIcon.FontColor.RGBA()),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	r, g, b, a := c.RGBA()
	l.config.DynamicIcon.FontColor.R = uint8(r)
	l.config.DynamicIcon.FontColor.G = uint8(g)
	l.config.DynamicIcon.FontColor.B = uint8(b)
	l.config.DynamicIcon.FontColor.A = uint8(a)
	if err := l.config.Write(); err != nil {
		return err
	}
	return nil
}
