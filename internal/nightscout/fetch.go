package nightscout

import (
	"encoding/json"
	"errors"
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

func Fetch() (Properties, error) {
	var properties Properties

	url := viper.GetString("url")
	if url == "" {
		return properties, errors.New("url is required")
	}

	// Fetch JSON
	resp, err := http.Get(url + "/api/v2/properties/bgnow,buckets,delta,direction")
	if err != nil {
		return properties, err
	}

	// Decode JSON
	if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
		return properties, err
	}

	return properties, nil
}
