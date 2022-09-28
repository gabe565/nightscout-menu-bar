package tray

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/gabe565/nightscout-systray/internal/autostart"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	"github.com/gabe565/nightscout-systray/internal/tray/items"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
	"log"
)

func Run() {
	systray.Run(onReady, onExit)
}

var Update = make(chan nightscout.Properties)
var Error = make(chan error, 1)

func onReady() {
	systray.SetTemplateIcon(assets.IconMenuBar, assets.IconMenuBar)
	systray.SetTitle("Nightscout")
	systray.SetTooltip("Nightscout")

	errorItem := items.NewError()
	openNightscoutItem := items.NewOpenNightscout()
	historyItem, historyVals := items.NewHistory()
	lastReadingItem := items.NewLastReading()
	systray.AddSeparator()
	prefs := items.NewPreferences()
	quitItem := items.NewQuit()

	beginTick()

	go func() {
		for {
			select {
			case <-openNightscoutItem.ClickedCh:
				url := viper.GetString("url")
				if err := open.Run(url); err != nil {
					Error <- err
				}
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

				lastReadingItem.SetTitle(properties.Bgnow.Time().String())
			case err := <-Error:
				systray.SetTitle("Error")
				errorItem.SetTitle(err.Error())
				errorItem.Show()
			}
		}
	}()
}

func onExit() {
	log.Println("exiting")
}
