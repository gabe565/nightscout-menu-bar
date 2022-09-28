package tray

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"time"
)

func init() {
	flag.DurationP("interval", "i", 30*time.Second, "Refresh interval")
	if err := viper.BindPFlag("interval", flag.Lookup("interval")); err != nil {
		panic(err)
	}
}

var ticker *time.Ticker

func beginTick() {
	ticker = time.NewTicker(viper.GetDuration("interval"))
	go func() {
		if err := fetchFromNightscout(); err != nil {
			log.Println(err)
		}

		for {
			select {
			case <-ticker.C:
				if err := fetchFromNightscout(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
