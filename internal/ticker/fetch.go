package ticker

import (
	"errors"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.DurationP("interval", "i", 30*time.Second, "Refresh interval")
	if err := viper.BindPFlag("interval", flag.Lookup("interval")); err != nil {
		panic(err)
	}
}

var fetchTimer = time.NewTimer(0)

func BeginFetch() {
	go func() {
		for range fetchTimer.C {
			Fetch()
			fetchTimer.Reset(viper.GetDuration("interval"))
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
