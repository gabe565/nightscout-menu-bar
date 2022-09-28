package items

import (
	"errors"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func NewPreferences() Preferences {
	item := systray.AddMenuItem("Preferences", "")
	item.SetTemplateIcon(assets.Gear, assets.Gear)

	urlTitle := "Nightscout URL"
	if url := viper.GetString("url"); url != "" {
		urlTitle += ": " + url
	}
	nightscoutUrl := item.AddSubMenuItem(urlTitle, "")

	startOnLogin := item.AddSubMenuItemCheckbox(
		"Start on login",
		"",
		autostart.IsEnabled(),
	)

	return Preferences{
		Item:         item,
		Url:          Url{nightscoutUrl},
		StartOnLogin: startOnLogin,
	}
}

type Preferences struct {
	Item         *systray.MenuItem
	Url          Url
	StartOnLogin *systray.MenuItem
}

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
