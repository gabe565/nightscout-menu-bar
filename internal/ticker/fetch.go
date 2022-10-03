package ticker

import (
	"errors"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
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

func BeginFetch() {
	go func() {
		ticker = time.NewTicker(viper.GetDuration("interval"))
		Fetch()

		for range ticker.C {
			Fetch()
		}
	}()
}

func Fetch() {
	properties, err := nightscout.Fetch()
	if err != nil && !errors.Is(err, nightscout.ErrNotModified) {
		tray.Error <- err
		return
	}
	if properties != nil {
		RenderCh <- properties
	}
}
