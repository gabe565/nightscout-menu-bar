package tray

import (
	"context"
	"log/slog"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/skratchdot/open-golang/open"
)

func New() *Tray {
	t := &Tray{
		config: config.NewDefault(),
		bus:    make(chan any, 1),
	}

	if err := t.config.Load(); err != nil {
		t.onError(err)
	}

	t.ticker = ticker.New(t.config, t.bus)

	t.config.AddCallback(func() {
		t.bus <- ReloadConfigMsg{}
	})
	return t
}

type Tray struct {
	config *config.Config
	ticker *ticker.Ticker
	bus    chan any
	items  items.Items
}

func (t *Tray) Run(ctx context.Context) {
	t.ticker.Start()
	if err := t.config.Watch(ctx); err != nil {
		t.onError(err)
	}
	go func() {
		<-ctx.Done()
		t.Quit()
	}()
	systray.Run(t.onReady, t.onExit)
}

func (t *Tray) Quit() {
	systray.Quit()
}

func (t *Tray) onReady() {
	systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
	systray.SetTitle(t.config.Title)
	systray.SetTooltip(t.config.Title)

	t.items = items.New(t.config)

	for {
		select {
		case <-t.items.OpenNightscout.ClickedCh:
			u, err := fetch.BuildURLWithToken(t.config)
			if err != nil {
				t.onError(err)
				return
			}
			if err := open.Run(u.String()); err != nil {
				t.onError(err)
			}
		case <-t.items.Preferences.URL.ClickedCh:
			go func() {
				if err := t.items.Preferences.URL.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-t.items.Preferences.Token.ClickedCh:
			go func() {
				if err := t.items.Preferences.Token.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-t.items.Preferences.Units.ClickedCh:
			go func() {
				if err := t.items.Preferences.Units.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-t.items.Preferences.StartOnLogin.ClickedCh:
			if t.items.Preferences.StartOnLogin.Checked() {
				if err := autostart.Disable(); err != nil {
					t.onError(err)
					continue
				}
				t.items.Preferences.StartOnLogin.Uncheck()
			} else {
				if err := autostart.Enable(); err != nil {
					t.onError(err)
					continue
				}
				t.items.Preferences.StartOnLogin.Check()
			}
		case <-t.items.Preferences.LocalFile.ClickedCh:
			if err := t.items.Preferences.LocalFile.Toggle(); err != nil {
				t.onError(err)
			}
		case <-t.items.Quit.ClickedCh:
			t.Quit()
		case msg := <-t.bus:
			switch msg := msg.(type) {
			case *nightscout.Properties:
				t.items.Error.Hide()

				value := msg.String(t.config)
				systray.SetTitle(value)
				systray.SetTooltip(value)
				t.items.LastReading.SetTitle(value)

				for i, reading := range msg.Buckets {
					if i < len(t.items.History.Subitems) {
						t.items.History.Subitems[i].SetTitle(reading.String(t.config))
					} else {
						entry := t.items.History.AddSubMenuItem(reading.String(t.config), "")
						entry.Disable()
						t.items.History.Subitems = append(t.items.History.Subitems, entry)
					}
				}
			case error:
				slog.Error("Displaying error", "error", msg.Error())
				t.items.Error.SetTitle(msg.Error())
				t.items.Error.Show()
			case ReloadConfigMsg:
				t.items.Preferences.URL.UpdateTitle()
				t.items.Preferences.Token.UpdateTitle()
				t.items.Preferences.Units.UpdateTitle()
			}
		}
	}
}

func (t *Tray) onError(err error) {
	select {
	case t.bus <- err:
	default:
		slog.Error("Unable to display error due to full bus", "error", err.Error())
	}
}

func (t *Tray) onExit() {
	slog.Info("Exiting")
	t.ticker.Close()
	close(t.bus)
}
