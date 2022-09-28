package ticker

import (
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	"github.com/gabe565/nightscout-systray/internal/tray"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

func init() {
	flag.DurationP("interval", "i", 30*time.Second, "Refresh interval")
	if err := viper.BindPFlag("interval", flag.Lookup("interval")); err != nil {
		panic(err)
	}
}

var ticker *time.Ticker

func BeginTick() {
	ticker = time.NewTicker(viper.GetDuration("interval"))
	go func() {
		Tick()

		for {
			select {
			case <-ticker.C:
				Tick()
			}
		}
	}()
}

func Tick() {
	properties, err := nightscout.Fetch()
	if err != nil {
		tray.Error <- err
		return
	}
	tray.Update <- properties
}
