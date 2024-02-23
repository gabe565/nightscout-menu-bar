package ticker

import (
	"errors"
	"log"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const IntervalKey = "interval"

func init() {
	flag.DurationP(IntervalKey, "i", 30*time.Second, "Refresh interval")
	if err := viper.BindPFlag(IntervalKey, flag.Lookup(IntervalKey)); err != nil {
		panic(err)
	}
}

var fetchTimer = time.NewTimer(0)

func BeginFetch() {
	go func() {
		for range fetchTimer.C {
			Fetch()
			fetchTimer.Reset(viper.GetDuration(IntervalKey))
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
		if local_file.Enabled {
			if err := local_file.Write(properties); err != nil {
				log.Println(err)
			}
		}
	}
}
