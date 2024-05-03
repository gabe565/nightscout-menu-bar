package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
)

func NewURL(conf *config.Config, parent *systray.MenuItem) URL {
	item := URL{config: conf}
	item.MenuItem = parent.AddSubMenuItem(item.GetTitle(), "")
	return item
}

type URL struct {
	config *config.Config
	*systray.MenuItem
}

func (n URL) GetTitle() string {
	title := "Nightscout URL"
	if n.config.URL != "" {
		title += ": " + n.config.URL
	}
	return title
}

func (n URL) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n URL) Prompt() error {
	url, err := ui.PromptURL(n.config.URL)
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
