package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func NewToken(parent *systray.MenuItem) Token {
	var token Token
	token.MenuItem = parent.AddSubMenuItem(token.GetTitle(), "")
	return token
}

type Token struct {
	*systray.MenuItem
}

func (n Token) GetTitle() string {
	title := "API Token"
	if url := viper.GetString("token"); url != "" {
		title += ": " + url
	}
	return title
}

func (n Token) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n Token) Prompt() error {
	url, err := ui.PromptToken()
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	viper.Set("token", url)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
