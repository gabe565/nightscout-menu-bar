package nightscout

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.StringP("url", "u", "", "Nightscout base URL")
	if err := viper.BindPFlag("url", flag.Lookup("url")); err != nil {
		panic(err)
	}
}

var lastEtag string

var (
	ErrHttp        = errors.New("unexpected HTTP error")
	ErrNotModified = errors.New("not modified")
)

var client = &http.Client{
	Timeout: time.Minute,
}

func Fetch() (*Properties, error) {
	url := viper.GetString("url")
	if url == "" {
		return nil, errors.New("please configure your Nightscout URL")
	}

	// Fetch JSON
	req, err := http.NewRequest("GET", url+"/api/v2/properties/bgnow,buckets,delta,direction", nil)
	if err != nil {
		return nil, err
	}
	if lastEtag != "" {
		req.Header.Set("If-None-Match", lastEtag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusNotModified:
		return nil, ErrNotModified
	case http.StatusOK:
		// Decode JSON
		var properties Properties
		if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
			return nil, err
		}

		lastEtag = resp.Header.Get("etag")
		return &properties, nil
	default:
		lastEtag = ""
		return nil, fmt.Errorf("%w: %d", ErrHttp, resp.StatusCode)
	}
}
