package tray

import (
	"errors"
	"log"

	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/getlantern/systray"
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

	errorItem := items.NewError()
	openNightscoutItem := items.NewOpenNightscout()
	historyItem, historyVals := items.NewHistory()
	lastReadingItem := items.NewLastReading()
	systray.AddSeparator()
	prefs := items.NewPreferences()
	quitItem := items.NewQuit()

	go func() {
		for {
			select {
			case <-openNightscoutItem.ClickedCh:
				url := viper.GetString("url")
				if err := open.Run(url); err != nil {
					Error <- err
				}
			case <-prefs.Url.ClickedCh:
				go func() {
					if err := prefs.Url.Prompt(); err != nil {
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

				systray.SetTitle(properties.String())

				for i, reading := range properties.Buckets {
					if i < len(historyVals) {
						historyVals[i].SetTitle(reading.String())
					} else {
						entry := historyItem.AddSubMenuItem(reading.String(), "")
						entry.Disable()
						historyVals = append(historyVals, entry)
					}
				}

				lastReadingItem.SetTitle(properties.Bgnow.Mills.String())
			case err := <-Error:
				if errors.As(err, &util.SoftError{}) {
					systray.SetTitle(viper.GetString("title"))
				} else {
					systray.SetTitle("Error")
				}
				errorItem.SetTitle(err.Error())
				errorItem.Show()
			}
		}
	}()
}

func onExit() {
	log.Println("exiting")
}
