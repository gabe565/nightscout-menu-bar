package tray

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
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
var Error = make(chan error)

func onReady() {
	systray.SetTemplateIcon(assets.IconMenuBar, assets.IconMenuBar)
	systray.SetTitle("Nightscout")
	systray.SetTooltip("Nightscout")

	errorItem := items.NewError()
	openNightscoutItem := items.NewOpenNightscout()
	historyItem, historyVals := items.NewHistory()
	lastReadingItem := items.NewLastReading()
	systray.AddSeparator()
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
			case <-quitItem.ClickedCh:
				systray.Quit()
			case properties := <-Update:
				errorItem.Hide()

				systray.SetTitle(properties.String())

				for i, reading := range properties.Buckets {
					var entry *systray.MenuItem
					if i < len(historyVals) {
						entry = historyVals[i]
					} else {
						entry = historyItem.AddSubMenuItem("", "")
						entry.Disable()
						historyVals = append(historyVals, entry)
					}
					entry.SetTitle(reading.String())
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
