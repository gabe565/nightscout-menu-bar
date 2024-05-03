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

	lastReadingItem := items.NewLastReading()
	errorItem := items.NewError()
	systray.AddSeparator()
	openNightscoutItem := items.NewOpenNightscout(t.config.Title)
	historyItem, historyVals := items.NewHistory()
	systray.AddSeparator()
	prefs := items.NewPreferences(t.config)
	quitItem := items.NewQuit()

	for {
		select {
		case <-openNightscoutItem.ClickedCh:
			u, err := fetch.BuildUrlWithToken(t.config)
			if err != nil {
				t.onError(err)
				return
			}
			if err := open.Run(u.String()); err != nil {
				t.onError(err)
			}
		case <-prefs.Url.ClickedCh:
			go func() {
				if err := prefs.Url.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-prefs.Token.ClickedCh:
			go func() {
				if err := prefs.Token.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-prefs.Units.ClickedCh:
			go func() {
				if err := prefs.Units.Prompt(); err != nil {
					t.onError(err)
				}
			}()
		case <-prefs.StartOnLogin.ClickedCh:
			if prefs.StartOnLogin.Checked() {
				if err := autostart.Disable(); err != nil {
					t.onError(err)
					continue
				}
				prefs.StartOnLogin.Uncheck()
			} else {
				if err := autostart.Enable(); err != nil {
					t.onError(err)
					continue
				}
				prefs.StartOnLogin.Check()
			}
		case <-prefs.LocalFile.ClickedCh:
			if err := prefs.LocalFile.Toggle(); err != nil {
				t.onError(err)
			}
		case <-quitItem.ClickedCh:
			t.Quit()
		case msg := <-t.bus:
			switch msg := msg.(type) {
			case *nightscout.Properties:
				errorItem.Hide()

				value := msg.String(t.config.Units, t.config.Arrows)
				systray.SetTitle(value)
				systray.SetTooltip(value)
				lastReadingItem.SetTitle(value)

				for i, reading := range msg.Buckets {
					if i < len(historyVals) {
						historyVals[i].SetTitle(reading.String(t.config.Units, t.config.Arrows))
					} else {
						entry := historyItem.AddSubMenuItem(reading.String(t.config.Units, t.config.Arrows), "")
						entry.Disable()
						historyVals = append(historyVals, entry)
					}
				}
			case error:
				errorItem.SetTitle(msg.Error())
				errorItem.Show()
			case ReloadConfigMsg:
				prefs.Url.UpdateTitle()
				prefs.Token.UpdateTitle()
				prefs.Units.UpdateTitle()
			}
		}
	}
}

func (t *Tray) onError(err error) {
	select {
	case t.bus <- err:
	default:
	}
}

func (t *Tray) onExit() {
	slog.Info("Exiting")
	t.ticker.Close()
	close(t.bus)
}
