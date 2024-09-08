package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func NewURL(conf *config.Config, parent *systray.MenuItem) URL {
	item := URL{config: conf}
	item.MenuItem = parent.AddSubMenuItem("Nightscout URL", "")
	return item
}

type URL struct {
	config *config.Config
	*systray.MenuItem
}

func (n URL) Prompt() error {
	url, err := zenity.Entry(
		"Enter new Nightscout URL:",
		zenity.Title("Nightscout URL"),
		zenity.EntryText(n.config.URL),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	n.config.URL = url
	if err := n.config.Write(); err != nil {
		return err
	}
	return nil
}
