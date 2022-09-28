package tray

import (
	"encoding/json"
	"errors"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func init() {
	flag.StringP("url", "u", "", "Nightscout base URL")
	if err := viper.BindPFlag("url", flag.Lookup("url")); err != nil {
		panic(err)
	}

	flag.DurationP("interval", "i", 30*time.Second, "Refresh interval")
	if err := viper.BindPFlag("interval", flag.Lookup("interval")); err != nil {
		panic(err)
	}
}

func fetchFromNightscout() error {
	url := viper.GetString("url")
	if url == "" {
		return errors.New("url is required")
	}

	// Fetch JSON
	resp, err := http.Get(url + "/api/v2/properties/bgnow,buckets,delta,direction")
	if err != nil {
		Error <- err
		return err
	}

	// Decode JSON
	var properties nightscout.Properties
	if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
		Error <- err
		return err
	}

	Update <- properties
	return nil
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

func ResetTicker() {
	if ticker != nil {
		ticker.Reset(viper.GetDuration("interval"))
	}
}
