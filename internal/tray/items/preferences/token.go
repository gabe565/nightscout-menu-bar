package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
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
	if config.Default.Token != "" {
		title += ": " + config.Default.Token
	}
	return title
}

func (n Token) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n Token) Prompt() error {
	token, err := ui.PromptToken()
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	config.Default.Token = token
	if err := config.Write(); err != nil {
		return err
	}
	return nil
}
