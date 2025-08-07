package preferences

import (
	"errors"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
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
		zenity.Title("Nightscout API Token"),
		zenity.EntryText(n.config.Data().Token),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}

	data := n.config.Data()
	data.Token = token
	if err := n.config.Write(data); err != nil {
		return err
	}
	return nil
}
