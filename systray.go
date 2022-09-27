package main

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"log"
)

var updateChan = make(chan nightscout.Properties)
var errorChan = make(chan error)

func onReady() {
	systray.SetTemplateIcon(assets.IconMenuBar, assets.IconMenuBar)
	systray.SetTitle("Nightscout")
	systray.SetTooltip("Nightscout")

	openNightscout := systray.AddMenuItem("Open Nightscout", "")
	openNightscout.SetTemplateIcon(assets.SquareUpRight, assets.SquareUpRight)

	history := systray.AddMenuItem("History", "")
	history.SetTemplateIcon(assets.RectangleHistory, assets.RectangleHistory)
	historyVals := make([]*systray.MenuItem, 0, 4)

	lastReading := systray.AddMenuItem("Last Reading", "")
	lastReading.SetTemplateIcon(assets.Calendar, assets.Calendar)
	lastReadingVal := lastReading.AddSubMenuItem("", "")
	lastReadingVal.Disable()

	go tick()

	for {
		select {
		case <-openNightscout.ClickedCh:
			if err := open.Run(url); err != nil {
				log.Println(err)
			}
		case properties := <-updateChan:
			systray.SetTitle(properties.String())

			for i, reading := range properties.Buckets {
				var entry *systray.MenuItem
				if i < len(historyVals) {
					entry = historyVals[i]
				} else {
					entry = history.AddSubMenuItem("", "")
					entry.Disable()
					historyVals = append(historyVals, entry)
				}
				entry.SetTitle(reading.String())
			}

			lastReadingVal.SetTitle(properties.Bgnow.Time().String())
		case <-errorChan:
			systray.SetTitle("Error")
		}
	}
}

func onExit() {
	log.Println("exiting")
}
