package tray

import (
	"encoding/json"
	"errors"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func init() {
	flag.StringP("url", "u", "", "Nightscout base URL")
	if err := viper.BindPFlag("url", flag.Lookup("url")); err != nil {
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
