package nightscout

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
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

var u *url.URL

func Fetch() (*Properties, error) {
	if u == nil {
		if err := UpdateUrl(); err != nil {
			return nil, err
		}
	}

	// Fetch JSON
	req, err := http.NewRequest("GET", u.String(), nil)
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

func UpdateUrl() error {
	conf := viper.GetString("url")
	if conf == "" {
		return errors.New("please configure your Nightscout URL")
	}

	newUrl, err := url.Parse(conf)
	if err != nil {
		return err
	}

	newUrl.Path = path.Join(newUrl.Path, "api", "v2", "properties", "bgnow,buckets,delta,direction")

	if token := viper.GetString("token"); token != "" {
		query := newUrl.Query()
		query.Set("token", token)
		newUrl.RawQuery = query.Encode()
	}

	u = newUrl
	return nil
}

func ClearUrl() {
	u = nil
}
