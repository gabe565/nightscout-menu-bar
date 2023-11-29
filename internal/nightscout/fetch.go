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

var (
	u     *url.URL
	token string
)

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

	if token != "" {
		req.Header.Set("Api-Secret", token)
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

func BuildUrl() (*url.URL, error) {
	conf := viper.GetString("url")
	if conf == "" {
		return nil, errors.New("please configure your Nightscout URL")
	}

	newUrl, err := url.Parse(conf)
	if err != nil {
		return nil, err
	}

	return newUrl, err
}

func BuildUrlWithToken() (*url.URL, error) {
	u, err := BuildUrl()
	if err != nil {
		return u, err
	}

	if token != "" {
		query := u.Query()
		query.Set("token", token)
		u.RawQuery = query.Encode()
	}

	return u, nil
}

func UpdateUrl() error {
	newUrl, err := BuildUrl()
	if err != nil {
		return err
	}

	newUrl.Path = path.Join(newUrl.Path, "api", "v2", "properties", "bgnow,buckets,delta,direction")
	u = newUrl

	token = viper.GetString("token")

	return nil
}

func ClearUrl() {
	u = nil
}
