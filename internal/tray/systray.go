package tray

import (
	"log"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("title", "Nightscout")
}

func Run() {
	systray.Run(onReady, onExit)
}

var (
	Update       = make(chan *nightscout.Properties)
	ReloadConfig = make(chan struct{})
	Error        = make(chan error, 1)
)

func onReady() {
	systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
	systray.SetTitle(viper.GetString("title"))
	systray.SetTooltip(viper.GetString("title"))

	lastReadingItem := items.NewLastReading()
	errorItem := items.NewError()
	systray.AddSeparator()
	openNightscoutItem := items.NewOpenNightscout()
	historyItem, historyVals := items.NewHistory()
	systray.AddSeparator()
	prefs := items.NewPreferences()
	quitItem := items.NewQuit()

	go func() {
		for {
			select {
			case <-openNightscoutItem.ClickedCh:
				go func() {
					u, err := nightscout.BuildUrl()
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
			case <-ReloadConfig:
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
			case <-quitItem.ClickedCh:
				systray.Quit()
			case properties := <-Update:
				errorItem.Hide()

				value := properties.String()
				systray.SetTitle(value)
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
	}()
}

func onExit() {
	log.Println("exiting")
}
