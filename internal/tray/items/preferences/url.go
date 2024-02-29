package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
)

func NewUrl(parent *systray.MenuItem) Url {
	var item Url
	item.MenuItem = parent.AddSubMenuItem(item.GetTitle(), "")
	return item
}

type Url struct {
	*systray.MenuItem
}

func (n Url) GetTitle() string {
	title := "Nightscout URL"
	if config.Default.URL != "" {
		title += ": " + config.Default.URL
	}
	return title
}

func (n Url) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n Url) Prompt() error {
	url, err := ui.PromptURL()
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	config.Default.URL = url
	if err := config.Write(); err != nil {
		return err
	}
	return nil
}
