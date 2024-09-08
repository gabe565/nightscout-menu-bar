package preferences

import (
	"errors"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/ncruces/zenity"
)

func NewToken(config *config.Config, parent *systray.MenuItem) Token {
	token := Token{config: config}
	token.MenuItem = parent.AddSubMenuItem("API Token", "")
	return token
}

type Token struct {
	config *config.Config
	*systray.MenuItem
}

func (n Token) Prompt() error {
	token, err := zenity.Entry(
		"Enter new Nightscout API token:",
		zenity.Title("Token"),
		zenity.EntryText(n.config.Token),
	)
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
