package nightscout

import (
	"encoding/json"
	"errors"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
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

var lastEtag string

func Fetch() (*Properties, error) {
	url := viper.GetString("url")
	if url == "" {
		return nil, util.SoftError{Err: errors.New("please configure your Nightscout URL")}
	}

	// Fetch JSON
	req, err := http.NewRequest("GET", url+"/api/v2/properties/bgnow,buckets,delta,direction", nil)
	if err != nil {
		return nil, err
	}
	if lastEtag != "" {
		req.Header.Set("If-None-Match", lastEtag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}

	// Decode JSON
	var properties Properties
	if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
		lastEtag = ""
		return nil, err
	}

	lastEtag = resp.Header.Get("etag")

	return &properties, nil
}
