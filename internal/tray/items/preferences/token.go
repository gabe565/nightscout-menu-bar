package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/ui"
	"github.com/ncruces/zenity"
)

func NewToken(config *config.Config, parent *systray.MenuItem) Token {
	token := Token{config: config}
	token.MenuItem = parent.AddSubMenuItem(token.GetTitle(), "")
	return token
}

type Token struct {
	config *config.Config
	*systray.MenuItem
}

func (n Token) GetTitle() string {
	title := "API Token"
	if n.config.Token != "" {
		title += ": " + n.config.Token
	}
	return title
}

func (n Token) UpdateTitle() {
	n.SetTitle(n.GetTitle())
}

func (n Token) Prompt() error {
	token, err := ui.PromptToken(n.config.Token)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	n.config.Token = token
	if err := n.config.Write(); err != nil {
		return err
	}
	return nil
}
