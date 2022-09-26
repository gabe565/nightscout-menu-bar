package main

import (
	"github.com/gabe565/nightscout-systray/internal/assets"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"log"
	"time"
)

var updateTitle = make(chan string)
var updateHistory = make(chan []nightscout.Reading)
var updateLastReading = make(chan time.Time)

func onReady() {
	systray.SetTemplateIcon(assets.Icon32, assets.Icon32)
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
		case v := <-updateTitle:
			systray.SetTitle(v)
		case v := <-updateHistory:
			for i, reading := range v {
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
		case v := <-updateLastReading:
			lastReadingVal.SetTitle(v.String())
		}
	}
}

func onExit() {
	log.Println("exiting")
}
