package tray

import (
	"log/slog"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/skratchdot/open-golang/open"
)

func Run() {
	systray.Run(onReady, onExit)
}

func init() {
	config.AddReloader(func() {
		reloadConfig <- struct{}{}
	})
}

var (
	Update       = make(chan *nightscout.Properties)
	reloadConfig = make(chan struct{})
	Error        = make(chan error, 1)
)

func onReady() {
	systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
	systray.SetTitle(config.Default.Title)
	systray.SetTooltip(config.Default.Title)

	lastReadingItem := items.NewLastReading()
	errorItem := items.NewError()
	systray.AddSeparator()
	openNightscoutItem := items.NewOpenNightscout(config.Default.Title)
	historyItem, historyVals := items.NewHistory()
	systray.AddSeparator()
	prefs := items.NewPreferences()
	quitItem := items.NewQuit()

	for {
		select {
		case <-openNightscoutItem.ClickedCh:
			go func() {
				u, err := nightscout.BuildUrlWithToken()
				if err != nil {
					Error <- err
					return
				}
				if err := open.Run(u.String()); err != nil {
					Error <- err
				}
			}()
		case <-prefs.Url.ClickedCh:
			go func() {
				if err := prefs.Url.Prompt(); err != nil {
					Error <- err
				}
			}()
		case <-prefs.Token.ClickedCh:
			go func() {
				if err := prefs.Token.Prompt(); err != nil {
					Error <- err
				}
			}()
		case <-prefs.Units.ClickedCh:
			go func() {
				if err := prefs.Units.Prompt(); err != nil {
					Error <- err
				}
			}()
		case <-reloadConfig:
			prefs.Url.UpdateTitle()
			prefs.Token.UpdateTitle()
			prefs.Units.UpdateTitle()
		case <-prefs.StartOnLogin.ClickedCh:
			if prefs.StartOnLogin.Checked() {
				if err := autostart.Disable(); err != nil {
					Error <- err
				}
				prefs.StartOnLogin.Uncheck()
			} else {
				if err := autostart.Enable(); err != nil {
					Error <- err
				}
				prefs.StartOnLogin.Check()
			}
		case <-prefs.LocalFile.ClickedCh:
			if err := prefs.LocalFile.Toggle(); err != nil {
				Error <- err
			}
		case <-quitItem.ClickedCh:
			systray.Quit()
		case properties := <-Update:
			errorItem.Hide()

			value := properties.String()
			systray.SetTitle(value)
			systray.SetTooltip(value)
			lastReadingItem.SetTitle(value)

			for i, reading := range properties.Buckets {
				if i < len(historyVals) {
					historyVals[i].SetTitle(reading.String())
				} else {
					entry := historyItem.AddSubMenuItem(reading.String(), "")
					entry.Disable()
					historyVals = append(historyVals, entry)
				}
			}
		case err := <-Error:
			errorItem.SetTitle(err.Error())
			errorItem.Show()
		}
	}
}

func onExit() {
	slog.Info("Exiting")
	if err := local_file.Cleanup(); err != nil {
		slog.Error("Failed to cleanup local file", "error", err.Error())
	}
}
