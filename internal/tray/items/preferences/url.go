package preferences

import (
	"errors"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

type Url struct {
	*systray.MenuItem
}

func (n Url) UpdateTitle() {
	title := "Nightscout URL"
	if url := viper.GetString("url"); url != "" {
		title += ": " + url
	}
	n.SetTitle(title)
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
