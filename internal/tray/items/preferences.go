package items

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/gabe565/nightscout-systray/internal/autostart"
	"github.com/gabe565/nightscout-systray/internal/ui"
	"github.com/getlantern/systray"
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
		return err
	}

	viper.Set("url", url)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
