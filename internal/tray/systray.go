package tray

import (
	"log"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/skratchdot/open-golang/open"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const TitleKey = "title"

func init() {
	flag.String(TitleKey, "Nightscout", "Title and hover text")
	if err := viper.BindPFlag(TitleKey, flag.Lookup(TitleKey)); err != nil {
		panic(err)
	}
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
	systray.SetTitle(viper.GetString(TitleKey))
	systray.SetTooltip(viper.GetString(TitleKey))

	lastReadingItem := items.NewLastReading()
	errorItem := items.NewError()
	systray.AddSeparator()
	openNightscoutItem := items.NewOpenNightscout()
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
}

func onExit() {
	log.Println("exiting")
	if err := local_file.Cleanup(); err != nil {
		log.Println(err)
	}
}
