package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
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
	if url := viper.GetString("url"); url != "" {
		title += ": " + url
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

	viper.Set("url", url)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
